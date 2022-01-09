package send

import (
	"fmt"
	"time"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/protocol"
	"github.com/marcinbor85/nes/broker"
	"github.com/marcinbor85/nes/common"
	"github.com/marcinbor85/nes/crypto/rsa"

	"github.com/marcinbor85/nes/api"
)

type SendContext struct {
	To *string
	Message *string
}

var Cmd = &cmd.Command{
	Name: "send",
	Help: "Send message to recipient",
	Context: &SendContext{},
	Init: Init,
	Logic: Logic,
}

func Init(c *cmd.Command) {
	ctx := c.Context.(*SendContext)
	ctx.To = c.Cmd.String("t", "to", &argparse.Options{Required: true, Help: "Recipient username"})
	ctx.Message = c.Cmd.String("m", "message", &argparse.Options{Required: true, Help: "Message to send"})
}

func Logic(c *cmd.Command) {
	ctx := c.Context.(*SendContext)

	recipient := *ctx.To

	pubkeyClient := &api.Client{
		Address: common.G.PubKeyAddress,
	}

	publicKey, err := pubkeyClient.GetPublicKeyByUsername(recipient)
	if err != nil {
		fmt.Println("Cannot get recipient public key:", err.Error())
		return
	}

	privateKey, _, err := rsa.LoadPrivateKey(common.G.PrivateKeyFile)
	if err != nil {
		fmt.Println("Cannot load private key:", err.Error())
		return
	}

	msg := &protocol.Message{
		From: common.G.Username,
		To: recipient,
		Timestamp: time.Now().UnixMilli(),
		Message: *ctx.Message,
	}

	frame, err := msg.Encrypt(publicKey, privateKey)
	if err != nil {
		fmt.Println("Cannot encrypt message:", err.Error())
		return
	}

	brokerClient := &broker.Client{
		BrokerAddress: common.G.MqttBrokerAddress,
	}

	err = brokerClient.Connect()
	if err != nil {
		fmt.Println("Cannot connect to broker:", err.Error())
		return
	}
	defer brokerClient.Disconnect()

	brokerClient.Send(frame, recipient)

	fmt.Println("Message sended successfully.")
}
