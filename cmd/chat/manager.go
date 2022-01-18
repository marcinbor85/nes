package chat

import (
	"fmt"
	"time"
	"errors"
	"strings"

	"crypto/rsa"

	"github.com/marcinbor85/nes/protocol"
	"github.com/marcinbor85/nes/broker"
	"github.com/marcinbor85/nes/common"
	r "github.com/marcinbor85/nes/crypto/rsa"

	"github.com/marcinbor85/nes/api"
)

type ChatManager struct {
	chatView	    				*ChatView
	inputView						*InputView
	privateKeyMessage				*rsa.PrivateKey
	privateKeySign					*rsa.PrivateKey
	recipientPublicKeyMessage		*rsa.PublicKey
	recipientPublicKeySign			*rsa.PublicKey
	pubkeyClient					*api.Client
	brokerClient					*broker.Client
	recipient						string
}

func NewChatManager(chatView *ChatView, inputView *InputView, recipient string) (*ChatManager, error) {
	pubkeyClient := api.NewClient(common.G.PubKeyAddress, common.G.PubKeyPublicKey)
	
	privateKeyMessage, _, err := r.LoadPrivateKey(common.G.PrivateKeyMessageFile)
	if err != nil {
		return nil, errors.New("Cannot load private key: " + err.Error())
	}

	privateKeySign, _, err := r.LoadPrivateKey(common.G.PrivateKeySignFile)
	if err != nil {
		return nil, errors.New("Cannot load private key: " + err.Error())
	}
	
	recipientPublicKeyMessage, recipientPublicKeySign, err := pubkeyClient.GetPublicKeyByUsername(recipient)
	if err != nil {
		return nil, errors.New("Cannot get recipient public key: " + err.Error())
	}

	self := &ChatManager{
		chatView: chatView,
		inputView: inputView,
		privateKeyMessage: privateKeyMessage,
		privateKeySign: privateKeySign,
		recipientPublicKeyMessage: recipientPublicKeyMessage,
		recipientPublicKeySign: recipientPublicKeySign,
		pubkeyClient: pubkeyClient,
		recipient: recipient,
		brokerClient: &broker.Client{
			BrokerAddress: common.G.MqttBrokerAddress,
			Recipient: common.G.Username,
			OnFrame: func(client *broker.Client, frame *protocol.Frame) {
				msg, e := frame.Decrypt(privateKeyMessage, pubkeyClient)
				if e != nil {
					return
				}
				if msg.From != recipient {
					return
				}
				text := fmt.Sprintf("%s > %s", msg.From, msg.Message)
				chatView.AddMessage(text)
			},
		},
	}

	chatView.SetChatManager(self)
	inputView.SetChatManager(self)

	return self, nil
}

func (chatManager *ChatManager) SendMessage(text string) error {
	text = strings.TrimSuffix(text, "\n")

	sender := common.G.Username
	recipient := chatManager.recipient
	
	msg := &protocol.Message{
		From: sender,
		To: recipient,
		Timestamp: time.Now().UnixMilli(),
		Message: text,
	}

	frame, err := msg.Encrypt(chatManager.recipientPublicKeyMessage, chatManager.privateKeySign)
	if err != nil {
		return errors.New("Cannot encrypt message: " + err.Error())
	}

	chatManager.brokerClient.Send(frame, recipient)

	showText := fmt.Sprintf("%s > %s", sender, text)
	chatManager.chatView.AddMessage(showText)
	return nil
}

func (chatManager *ChatManager) Recipient() string {
	return chatManager.recipient
}

func (chatManager *ChatManager) ScrollUp() {
	chatManager.chatView.ScrollUp();
}

func (chatManager *ChatManager) ScrollDown() {
	chatManager.chatView.ScrollDown();
}

func (chatManager *ChatManager) Connect() error {
	err := chatManager.brokerClient.Connect();
	if err != nil {
		return errors.New("Cannot connect to broker: " + err.Error())
	}
	return nil
}

func (chatManager *ChatManager) Disconnect() error {
	chatManager.brokerClient.Disconnect();
	return nil
}

func (chatManager *ChatManager) Start() error {
	err := chatManager.Connect()
	return err
}