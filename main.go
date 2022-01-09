package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/akamensky/argparse"

	"github.com/marcinbor85/nes/config"
	"github.com/marcinbor85/nes/crypto"
	"github.com/marcinbor85/nes/common"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/cmd/send"
	"github.com/marcinbor85/nes/cmd/listen"
	"github.com/marcinbor85/nes/cmd/register"

	"github.com/marcinbor85/pubkey/api"
)

const (
	MQTT_BROKER_ADDRESS_DEFAULT = "tcp://test.mosquitto.org:1883"
	PUBKEY_ADDRESS_DEFAULT      = "https://microshell.pl/pubkey"
	CONFIG_FILE_DEFAULT         = ".env"

	APP_SETTINGS_HOME_DIR       = ".nes"
)

func main() {
	parser := argparse.NewParser("nes", "NES messenger")

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
		Help:     `Private key file. Default: ~/` + APP_SETTINGS_HOME_DIR + `/<user>_rsa`,
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

	commands := []*cmd.Command{
		register.Cmd,
		listen.Cmd,
		send.Cmd,
	}

	for _, c := range commands {
		c.Register(parser)
	}

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	config.Init(*configArg)

	common.G.Settings = &common.Settings{}
	common.G.Settings.MqttBrokerAddress = config.Alternate(*brokerArg, "MQTT_BROKER_ADDRESS", MQTT_BROKER_ADDRESS_DEFAULT)
	common.G.Settings.PubKeyAddress = config.Alternate(*providerArg, "PUBKEY_ADDRESS", PUBKEY_ADDRESS_DEFAULT)
	
	osUser, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	
	common.G.Settings.Username = config.Alternate(*usernameArg, "USERNAME", osUser.Username)

	keyFilename := strings.Join([]string{common.G.Settings.Username, "rsa"}, "-")
	defPrivateKeyFile := path.Join(osUser.HomeDir, APP_SETTINGS_HOME_DIR, keyFilename)

	common.G.Settings.PrivateKeyFile = config.Alternate(*privateArg, "PRIVATE_KEY_FILE", defPrivateKeyFile)

	crypto.Init()

	common.G.PubkeyClient = &api.Client{
		Address: common.G.Settings.PubKeyAddress,
	}

	for _, c := range commands {
		if c.IsInvoked() {
			c.Service()
			break;
		}
	}
}
