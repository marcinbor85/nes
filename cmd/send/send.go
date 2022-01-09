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
)

type SendContext struct {
	To *string
	Message *string
}

var Cmd = &cmd.Command{
	Name: "send",
	Help: "send message to recipient",
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

	// TODO: check if local user exist

	publicKey, apiErr := common.G.PubkeyClient.GetPublicKeyByUsername(recipient)
	if apiErr != nil {
		fmt.Println("unknown recipient username")
		return
	}

	privateKey, err := rsa.LoadPrivateKey(common.G.Settings.PrivateKeyFile)
	if err != nil {
		fmt.Println("cannot load private key:", err.Error())
		return
	}

	msg := &protocol.Message{
		From: common.G.Settings.Username,
		To: recipient,
		Timestamp: time.Now().UnixMilli(),
		Message: *ctx.Message,
	}

	fmt.Println(msg)

	frame, e := msg.Encrypt(publicKey, privateKey)
	if e != nil {
		fmt.Println("cannot encrypt:", e.Error())
		return
	}

	brokerClient := &broker.Client{
		BrokerAddress: common.G.Settings.MqttBrokerAddress,
	}

	er := brokerClient.Connect()
	if er != nil {
		fmt.Println("cannot connect to broker:", er.Error())
		return
	}
	defer brokerClient.Disconnect()

	brokerClient.Send(frame, recipient)
}
