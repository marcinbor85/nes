package listen

import (
	"fmt"
	"time"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/protocol"
	"github.com/marcinbor85/nes/broker"
	"github.com/marcinbor85/nes/common"
	"github.com/marcinbor85/nes/crypto/rsa"

	"github.com/marcinbor85/nes/api"
)

type ListenContext struct {
}

var Cmd = &cmd.Command{
	Name: "listen",
	Help: "Listen to messages",
	Context: &ListenContext{},
	Init: Init,
	Logic: Logic,
}

func Init(c *cmd.Command) {
}

func Logic(c *cmd.Command) {
	pubkeyClient := api.NewClient(common.G.PubKeyAddress, common.G.PubKeyPublicKey)

	privateKeyMessage, _, err := rsa.LoadPrivateKey(common.G.PrivateKeyMessageFile)
	if err != nil {
		fmt.Println("Cannot load private key:", err.Error())
		return
	}

	brokerClient := &broker.Client{
		BrokerAddress: common.G.MqttBrokerAddress,
		Recipient: common.G.Username,
		OnFrame: func(client *broker.Client, frame *protocol.Frame) {
			msg, e := frame.Decrypt(privateKeyMessage, pubkeyClient)
			if e != nil {
				fmt.Printf("%v\n",e)
				return
			}
			t := time.UnixMilli(msg.Timestamp)
			tm := t.Format("15:04:05")
			fmt.Printf("[%s] %s > %s\n", tm, msg.From, msg.Message)
		},
	}

	er := brokerClient.Connect()
	if er != nil {
		fmt.Println("Cannot connect to broker:", er.Error())
		return
	}
	defer brokerClient.Disconnect()

	fmt.Println("Press the Enter Key to exit.")
	var s string
	fmt.Scanln(&s)
}
