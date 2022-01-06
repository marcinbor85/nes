package main

import (
	"os"
	"os/user"
	"fmt"

	"github.com/akamensky/argparse"
	
	"github.com/marcinbor85/nes/crypto"
	"github.com/marcinbor85/nes/config"
)

type Settings struct {
	MqttBroker string
	KeyProvider string
	PrivateKey string
	Nickname string
}

var settings = &Settings{}

const (
	MQTT_BROKER_ADDRESS_DEFAULT = "test.mosquitto.org"
	PUBKEY_ADDRESS_DEFAULT = "microshell.pl/pubkey"
	PRIVATE_KEY_FILE_DEFAULT = "~/.ssh/id_rsa"
	CONFIG_FILE_DEFAULT = ".env"
)

func main() {
	parser := argparse.NewParser("commands", "NES messenger")

	brokerArg := parser.String("b", "broker", &argparse.Options{
		Required: false,
		Help: `MQTT broker server address. Default: ` + MQTT_BROKER_ADDRESS_DEFAULT,
		Default: nil,
	})

	providerArg := parser.String("p", "provider", &argparse.Options{
		Required: false,
		Help: `Public key provider address. Default: ` + PUBKEY_ADDRESS_DEFAULT,
		Default: nil,
	})

	privateArg := parser.String("k", "key", &argparse.Options{
		Required: false,
		Help: `Private key file. Default: ` + PRIVATE_KEY_FILE_DEFAULT,
		Default: nil,
	})

	nickArg := parser.String("n", "nick", &argparse.Options{
		Required: false,
		Help: "Local nickname. Default: <os_user>",
		Default: nil,
	})

	configArg := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help: `Optional config file. Supported fields: MQTT_BROKER_ADDRESS, PUBKEY_ADDRESS, PRIVATE_KEY_FILE, NICKNAME`,
		Default: CONFIG_FILE_DEFAULT,
	})

	listenCmd := parser.NewCommand("listen", "listen to messages")

	sendCmd := parser.NewCommand("send", "send message to recipient")
	_ = sendCmd.String("t", "to", &argparse.Options{Required: true, Help: "Recipient nickname"})
	sendCmd.Flag("i", "interactive", &argparse.Options{Help: "Enable interactive mode"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	config.Init(*configArg)
	
	settings.MqttBroker = config.Alternate(*brokerArg, "MQTT_BROKER_ADDRESS", MQTT_BROKER_ADDRESS_DEFAULT)
	settings.KeyProvider = config.Alternate(*providerArg, "PUBKEY_ADDRESS", PUBKEY_ADDRESS_DEFAULT)
	settings.PrivateKey = config.Alternate(*privateArg, "PRIVATE_KEY_FILE", PRIVATE_KEY_FILE_DEFAULT)

	osUser, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	settings.Nickname = config.Alternate(*nickArg, "NICKNAME", osUser.Username)

	crypto.Init()

	fmt.Println(settings)

	if listenCmd.Happened() {

	} else if sendCmd.Happened() {

	} else {
		panic("really?")
	}
}
