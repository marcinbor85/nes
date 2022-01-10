package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/akamensky/argparse"

	cfg "github.com/marcinbor85/nes/config"
	"github.com/marcinbor85/nes/crypto"
	"github.com/marcinbor85/nes/common"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/cmd/send"
	"github.com/marcinbor85/nes/cmd/listen"
	"github.com/marcinbor85/nes/cmd/register"
	"github.com/marcinbor85/nes/cmd/config"
	"github.com/marcinbor85/nes/cmd/generate"
	"github.com/marcinbor85/nes/cmd/chat"
)

const (
	MQTT_BROKER_ADDRESS_DEFAULT = "tcp://test.mosquitto.org:1883"
	PUBKEY_ADDRESS_DEFAULT      = "https://microshell.pl/pubkey"

	APP_SETTINGS_HOME_DIR       = ".nes"
	APP_SETTINGS_CONFIG_FILE    = "config"
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

	privateArg := parser.String("k", "private", &argparse.Options{
		Required: false,
		Help:     `Private key file. Default: ~/` + APP_SETTINGS_HOME_DIR + `/<user>-rsa`,
		Default:  nil,
	})

	publicArg := parser.String("K", "public", &argparse.Options{
		Required: false,
		Help:     `Public key file. Default: ~/` + APP_SETTINGS_HOME_DIR + `/<user>-rsa.pub`,
		Default:  nil,
	})

	usernameArg := parser.String("u", "user", &argparse.Options{
		Required: false,
		Help:     "Local username. Default: <os_user>",
		Default:  nil,
	})

	configArg := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     `Optional config file. Supported fields: MQTT_BROKER_ADDRESS, PUBKEY_ADDRESS, PRIVATE_KEY_FILE, PUBLIC_KEY_FILE, USERNAME. Default: ~/` + APP_SETTINGS_HOME_DIR + `/` + APP_SETTINGS_CONFIG_FILE,
		Default:  nil,
	})

	commands := []*cmd.Command{
		register.Cmd,
		listen.Cmd,
		send.Cmd,
		config.Cmd,
		generate.Cmd,
		chat.Cmd,
	}

	for _, c := range commands {
		c.Register(parser)
	}

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	osUser, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	appDir := path.Join(osUser.HomeDir, APP_SETTINGS_HOME_DIR)
	err = os.MkdirAll(appDir, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}

	configFile := path.Join(appDir, APP_SETTINGS_CONFIG_FILE)
	if *configArg != "" {
		configFile = *configArg
	}
	common.G.ConfigFile = configFile

	cfg.Init(common.G.ConfigFile)

	common.G.MqttBrokerAddress = cfg.Alternate(*brokerArg, "MQTT_BROKER_ADDRESS", MQTT_BROKER_ADDRESS_DEFAULT)
	common.G.PubKeyAddress = cfg.Alternate(*providerArg, "PUBKEY_ADDRESS", PUBKEY_ADDRESS_DEFAULT)
	common.G.Username = cfg.Alternate(*usernameArg, "USERNAME", osUser.Username)

	keyFilename := strings.Join([]string{common.G.Username, "rsa"}, "-")
	defKeyFile := path.Join(osUser.HomeDir, APP_SETTINGS_HOME_DIR, keyFilename)
	common.G.PrivateKeyFile = cfg.Alternate(*privateArg, "PRIVATE_KEY_FILE", defKeyFile)

	keyFilename = strings.Join([]string{common.G.Username, "rsa.pub"}, "-")
	defKeyFile = path.Join(osUser.HomeDir, APP_SETTINGS_HOME_DIR, keyFilename)
	common.G.PublicKeyFile = cfg.Alternate(*publicArg, "PUBLIC_KEY_FILE", defKeyFile)

	crypto.Init()

	for _, c := range commands {
		if c.IsInvoked() {
			c.Service()
			break;
		}
	}
}
