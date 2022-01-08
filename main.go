package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
	"bufio"
	"os/user"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/config"
	"github.com/marcinbor85/nes/crypto"
	"github.com/marcinbor85/nes/protocol"
	"github.com/marcinbor85/nes/broker"

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

	pubkeyClient := &api.Client{
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

		err := pubkeyClient.RegisterNewUsername(settings.Username, *emailArg, publicKeyPem)
		if err != nil {
			fmt.Printf(err.E.Error())
			return
		}
		fmt.Println("username registered. check email for activation.")

	} else if listenCmd.Happened() {
		// TODO: check if local user exist

		privateKey, err := crypto.LoadPrivateKey(settings.PrivateKeyFile)
		if err != nil {
			fmt.Println("cannot load private key:", err.Error())
			return
		}

		brokerClient := &broker.Client{
			BrokerAddress: settings.MqttBrokerAddress,
			OnFrame: func(client *broker.Client, frame *protocol.Frame) {
				msg, e := frame.Decrypt(privateKey, pubkeyClient)
				if e != nil {
					fmt.Println("cannot decrypt:", e.Error())
					return
				}
				t := time.UnixMilli(msg.Timestamp)
				tm := t.Format("2006-01-02 15:04:05")
				fmt.Printf("[%s] %s > %s\n", tm, msg.From, msg.Message)
			},
		}

		er := brokerClient.Connect()
		if er != nil {
			fmt.Println("cannot connect to broker:", er.Error())
			return
		}
		defer brokerClient.Disconnect()

		er = brokerClient.Recv(settings.Username)
		if er != nil {
			fmt.Println("cannot subscribe:", er.Error())
			return
		}

		fmt.Println("Press the Enter Key to exit.")
    	fmt.Scanln()

	} else if sendCmd.Happened() {
		recipient := *toArg

		// TODO: check if local user exist

		publicKey, apiErr := pubkeyClient.GetPublicKeyByUsername(recipient)
		if apiErr != nil {
			fmt.Println("unknown recipient username")
			return
		}

		privateKey, err := crypto.LoadPrivateKey(settings.PrivateKeyFile)
		if err != nil {
			fmt.Println("cannot load private key:", err.Error())
			return
		}

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		message := scanner.Text()

		msg := &protocol.Message{
			From: settings.Username,
			To: recipient,
			Timestamp: time.Now().UnixMilli(),
			Message: message,
		}

		fmt.Println(msg)

		frame, e := msg.Encrypt(publicKey, privateKey)
		if e != nil {
			fmt.Println("cannot encrypt:", e.Error())
			return
		}

		brokerClient := &broker.Client{
			BrokerAddress: settings.MqttBrokerAddress,
		}

		er := brokerClient.Connect()
		if er != nil {
			fmt.Println("cannot connect to broker:", er.Error())
			return
		}
		defer brokerClient.Disconnect()

		brokerClient.Send(frame, recipient)

		

	} else {
		panic("really?")
	}
}
