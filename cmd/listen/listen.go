package listen

import (
	"fmt"
	"time"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/protocol"
	"github.com/marcinbor85/nes/broker"
	"github.com/marcinbor85/nes/common"
	"github.com/marcinbor85/nes/crypto/rsa"

	"github.com/marcinbor85/pubkey/api"
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
	// TODO: check if local user exist

	pubkeyClient := &api.Client{
		Address: common.G.PubKeyAddress,
	}

	privateKey, err := rsa.LoadPrivateKey(common.G.PrivateKeyFile)
	if err != nil {
		fmt.Println("cannot load private key:", err.Error())
		return
	}

	brokerClient := &broker.Client{
		BrokerAddress: common.G.MqttBrokerAddress,
		Recipient: common.G.Username,
		OnFrame: func(client *broker.Client, frame *protocol.Frame) {
			msg, e := frame.Decrypt(privateKey, pubkeyClient)
			if e != nil {
				fmt.Println("cannot decrypt:", e.Error())
				return
			}
			t := time.UnixMilli(msg.Timestamp)
			tm := t.Format("2006-01-02 15:04:05")
			fmt.Printf("\x1B[2K\r")
			fmt.Printf("[%s] %s > %s\r\n", tm, msg.From, msg.Message)
		},
	}

	er := brokerClient.Connect()
	if er != nil {
		fmt.Println("cannot connect to broker:", er.Error())
		return
	}
	defer brokerClient.Disconnect()

	fmt.Println("Press the Enter Key to exit.")
	var s string
	fmt.Scanln(&s)

}
