package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
	"os/user"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/config"
	"github.com/marcinbor85/nes/crypto"
	"github.com/marcinbor85/nes/protocol"

	"github.com/marcinbor85/pubkey/api"
)

type Settings struct {
	MqttBrokerAddress string
	PubKeyAddress     string
	PrivateKeyFile    string
	Username          string
}

var settings = &Settings{}

const (
	MQTT_BROKER_ADDRESS_DEFAULT = "test.mosquitto.org"
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

	listenCmd := parser.NewCommand("listen", "listen to messages")

	sendCmd := parser.NewCommand("send", "send message to recipient")
	toArg := sendCmd.String("t", "to", &argparse.Options{Required: true, Help: "Recipient username"})
	_ = sendCmd.Flag("i", "interactive", &argparse.Options{Help: "Enable interactive mode"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	config.Init(*configArg)

	settings.MqttBrokerAddress = config.Alternate(*brokerArg, "MQTT_BROKER_ADDRESS", MQTT_BROKER_ADDRESS_DEFAULT)
	settings.PubKeyAddress = config.Alternate(*providerArg, "PUBKEY_ADDRESS", PUBKEY_ADDRESS_DEFAULT)
	settings.PrivateKeyFile = config.Alternate(*privateArg, "PRIVATE_KEY_FILE", PRIVATE_KEY_FILE_DEFAULT)

	osUser, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	settings.Username = config.Alternate(*usernameArg, "USERNAME", osUser.Username)

	crypto.Init()

	client := &api.Client{
		Address: settings.PubKeyAddress,
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

		err := client.RegisterNewUsername(settings.Username, *emailArg, publicKeyPem)
		if err != nil {
			fmt.Printf(err.E.Error())
			return
		}
		fmt.Println("username registered. check email for activation.")

	} else if listenCmd.Happened() {

	} else if sendCmd.Happened() {
		recipient := *toArg

		// TODO: check if local user exist

		f, e := os.Open(settings.PrivateKeyFile)
		if e != nil {
			fmt.Println("cannot open private key file")
			return
		}
		defer f.Close()

		bytesBuf := &bytes.Buffer{}
		bytesBuf.ReadFrom(f)
		privateKeyPem := bytesBuf.String()

		privateKey, err := crypto.DecodePrivateKey(privateKeyPem)
		if err != nil {
			fmt.Println("cannot decode private key")
			return
		}

		publicKey, ee := client.GetPublicKeyByUsername(recipient)
		if ee != nil {
			fmt.Println("unknown recipient username")
			return
		}

		msg := &protocol.Message{
			From: settings.Username,
			To: recipient,
			Timestamp: time.Now().UnixMilli(),
			Message: "test",
		}

		fmt.Println(msg)

		frame, e := msg.Encrypt(publicKey, privateKey)
		if e != nil {
			fmt.Println("cannot encrypt:", e.Error())
			return
		}
		fmt.Println(frame)

		msg2, e := frame.Decrypt(privateKey, client)
		if e != nil {
			fmt.Println("cannot decrypt:", e.Error())
			return
		}
		fmt.Println(msg2)

	} else {
		panic("really?")
	}
}
