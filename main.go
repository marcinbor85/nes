package main

import (
	"fmt"
	"os"
	"os/user"
	path "path/filepath"
	"strings"

	"github.com/akamensky/argparse"

	cfg "github.com/marcinbor85/nes/config"
	"github.com/marcinbor85/nes/crypto"
	r "github.com/marcinbor85/nes/crypto/rsa"
	"github.com/marcinbor85/nes/common"
	"github.com/marcinbor85/nes/keys"

	"github.com/marcinbor85/nes/cmd"
	"github.com/marcinbor85/nes/cmd/send"
	"github.com/marcinbor85/nes/cmd/listen"
	"github.com/marcinbor85/nes/cmd/register"
	"github.com/marcinbor85/nes/cmd/config"
	"github.com/marcinbor85/nes/cmd/generate"
	"github.com/marcinbor85/nes/cmd/chat"
	"github.com/marcinbor85/nes/cmd/version"
)

const (
	MQTT_BROKER_ADDRESS_DEFAULT = "tcp://test.mosquitto.org:1883"
	PUBKEY_ADDRESS_DEFAULT      = "https://microshell.pl/pubkey"

	APP_SETTINGS_HOME_DIR       = ".nes"
	APP_SETTINGS_CONFIG_FILE    = "config"

	KEY_MESSAGE_SUFFIX			= "-rsa-message"
	KEY_SIGN_SUFFIX				= "-rsa-sign"
	KEY_PUBLIC_SUFFIX			= ".pub"
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
		Help:     `Public keys provider address. Default: ` + PUBKEY_ADDRESS_DEFAULT,
		Default:  nil,
	})

	providerPublicKeyArg := parser.String("P", "provider_key", &argparse.Options{
		Required: false,
		Help:     `Provider public key. Default: PUBLIC_KEY_OF(` + PUBKEY_ADDRESS_DEFAULT + `)`,
		Default:  nil,
	})

	privateMessageArg := parser.String("k", "private_message", &argparse.Options{
		Required: false,
		Help:     `Private key file for message. Default: ~/` + APP_SETTINGS_HOME_DIR + `/<user>` + KEY_MESSAGE_SUFFIX,
		Default:  nil,
	})

	publicMessageArg := parser.String("K", "public_message", &argparse.Options{
		Required: false,
		Help:     `Public key file for message. Default: ~/` + APP_SETTINGS_HOME_DIR + `/<user>` + KEY_MESSAGE_SUFFIX + KEY_PUBLIC_SUFFIX,
		Default:  nil,
	})

	privateSignArg := parser.String("j", "private_sign", &argparse.Options{
		Required: false,
		Help:     `Private key file for signing. Default: ~/` + APP_SETTINGS_HOME_DIR + `/<user>` + KEY_SIGN_SUFFIX,
		Default:  nil,
	})

	publicSignArg := parser.String("J", "public_sign", &argparse.Options{
		Required: false,
		Help:     `Public key file for signature verification. Default: ~/` + APP_SETTINGS_HOME_DIR + `/<user>` + KEY_SIGN_SUFFIX + KEY_PUBLIC_SUFFIX,
		Default:  nil,
	})

	usernameArg := parser.String("u", "user", &argparse.Options{
		Required: false,
		Help:     "Local username. Default: <os_user>",
		Default:  nil,
	})

	configArg := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     `Optional config file. Supported fields: MQTT_BROKER_ADDRESS, PUBKEY_ADDRESS, PUBKEY_PUBLIC_KEY_FILE, PRIVATE_KEY_MESSAGE_FILE, PUBLIC_KEY_MESSAGE_FILE, PRIVATE_KEY_SIGN_FILE, PUBLIC_KEY_SIGN_FILE, NES_USERNAME. Default: ~/` + APP_SETTINGS_HOME_DIR + `/` + APP_SETTINGS_CONFIG_FILE,
		Default:  nil,
	})

	commands := []*cmd.Command{
		register.Cmd,
		listen.Cmd,
		send.Cmd,
		config.Cmd,
		generate.Cmd,
		chat.Cmd,
		version.Cmd,
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
	username := cfg.Alternate(*usernameArg, "NES_USERNAME", osUser.Username)
	common.G.Username = strings.ToLower(username)

	keyFilename := common.G.Username + KEY_MESSAGE_SUFFIX
	defKeyFile := path.Join(osUser.HomeDir, APP_SETTINGS_HOME_DIR, keyFilename)
	common.G.PrivateKeyMessageFile = cfg.Alternate(*privateMessageArg, "PRIVATE_KEY_MESSAGE_FILE", defKeyFile)

	keyFilename = common.G.Username + KEY_MESSAGE_SUFFIX + KEY_PUBLIC_SUFFIX
	defKeyFile = path.Join(osUser.HomeDir, APP_SETTINGS_HOME_DIR, keyFilename)
	common.G.PublicKeyMessageFile = cfg.Alternate(*publicMessageArg, "PUBLIC_KEY_MESSAGE_FILE", defKeyFile)

	keyFilename = common.G.Username + KEY_SIGN_SUFFIX
	defKeyFile = path.Join(osUser.HomeDir, APP_SETTINGS_HOME_DIR, keyFilename)
	common.G.PrivateKeySignFile = cfg.Alternate(*privateSignArg, "PRIVATE_KEY_SIGN_FILE", defKeyFile)

	keyFilename = common.G.Username + KEY_SIGN_SUFFIX + KEY_PUBLIC_SUFFIX
	defKeyFile = path.Join(osUser.HomeDir, APP_SETTINGS_HOME_DIR, keyFilename)
	common.G.PublicKeySignFile = cfg.Alternate(*publicSignArg, "PUBLIC_KEY_SIGN_FILE", defKeyFile)

	pubKeyPublicKeyFile := cfg.Alternate(*providerPublicKeyArg, "PUBKEY_PUBLIC_KEY_FILE", "")
	if pubKeyPublicKeyFile != "" {
		common.G.PubKeyPublicKeyFile = pubKeyPublicKeyFile
		key, _, err := r.LoadPublicKey(pubKeyPublicKeyFile)
		if err != nil {
			panic(err.Error())
		}
		common.G.PubKeyPublicKey = key
	} else {
		common.G.PubKeyPublicKey, _ = r.DecodePublicKey(keys.MICROSHELL_PL_PUBKEY)
	}

	crypto.Init()

	for _, c := range commands {
		if c.IsInvoked() {
			c.Service()
			break;
		}
	}
}
