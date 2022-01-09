package main

import (
	"bytes"
	"fmt"
	"os"
	"os/user"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/config"
	"github.com/marcinbor85/nes/crypto"
	"github.com/marcinbor85/nes/common"

	"github.com/marcinbor85/nes/cmd/send"
	"github.com/marcinbor85/nes/cmd/listen"

	"github.com/marcinbor85/pubkey/api"
)

const (
	MQTT_BROKER_ADDRESS_DEFAULT = "tcp://test.mosquitto.org:1883"
	PUBKEY_ADDRESS_DEFAULT      = "https://microshell.pl/pubkey"
	PRIVATE_KEY_FILE_DEFAULT    = "~/.ssh/id_rsa"
	CONFIG_FILE_DEFAULT         = ".env"
)

func main() {
	parser := argparse.NewParser("commands", "NES messenger")

	brokerArg := parser.String("b", "broker", &argparse.Options{
		Required: false,
		Help:     `MQTT broker server address. Default: ` + MQTT_BROKER_ADDRESS_DEFAULT,
		Default:  nil,
	})

	providerArg := parser.String("p", "provider", &argparse.Options{
		Required: false,
		Help:     `Public key provider address. Default: ` + PUBKEY_ADDRESS_DEFAULT,
		Default:  nil,
	})

	privateArg := parser.String("k", "key", &argparse.Options{
		Required: false,
		Help:     `Private key file. Default: ` + PRIVATE_KEY_FILE_DEFAULT,
		Default:  nil,
	})

	usernameArg := parser.String("u", "user", &argparse.Options{
		Required: false,
		Help:     "Local username. Default: <os_user>",
		Default:  nil,
	})

	configArg := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     `Optional config file. Supported fields: MQTT_BROKER_ADDRESS, PUBKEY_ADDRESS, PRIVATE_KEY_FILE, USERNAME`,
		Default:  CONFIG_FILE_DEFAULT,
	})

	registerCmd := parser.NewCommand("register", "register username")
	publicArg := registerCmd.String("P", "public", &argparse.Options{
		Required: true,
		Help:     `Public key file.`,
	})
	emailArg := registerCmd.String("e", "email", &argparse.Options{
		Required: true,
		Help:     `User email (need to activate username).`,
	})

	listen.ListenCmd.Register(parser)
	send.SendCmd.Register(parser)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	config.Init(*configArg)


	common.G.Settings = &common.Settings{}
	common.G.Settings.MqttBrokerAddress = config.Alternate(*brokerArg, "MQTT_BROKER_ADDRESS", MQTT_BROKER_ADDRESS_DEFAULT)
	common.G.Settings.PubKeyAddress = config.Alternate(*providerArg, "PUBKEY_ADDRESS", PUBKEY_ADDRESS_DEFAULT)
	common.G.Settings.PrivateKeyFile = config.Alternate(*privateArg, "PRIVATE_KEY_FILE", PRIVATE_KEY_FILE_DEFAULT)

	osUser, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	common.G.Settings.Username = config.Alternate(*usernameArg, "USERNAME", osUser.Username)

	crypto.Init()

	common.G.PubkeyClient = &api.Client{
		Address: common.G.Settings.PubKeyAddress,
	}

	if registerCmd.Happened() {

		f, e := os.Open(*publicArg)
		if e != nil {
			fmt.Println("cannot open public key file")
			return
		}
		defer f.Close()

		bytesBuf := &bytes.Buffer{}
		bytesBuf.ReadFrom(f)
		publicKeyPem := bytesBuf.String()

		err := common.G.PubkeyClient.RegisterNewUsername(common.G.Settings.Username, *emailArg, publicKeyPem)
		if err != nil {
			fmt.Printf(err.E.Error())
			return
		}
		fmt.Println("username registered. check email for activation.")

	} else if listen.ListenCmd.IsInvoked() {
		listen.ListenCmd.Service()
	} else if send.SendCmd.IsInvoked() {
		send.SendCmd.Service()
	} else {
		panic("really?")
	}
}
