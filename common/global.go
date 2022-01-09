package common

import (
	"github.com/marcinbor85/pubkey/api"
)

type Settings struct {
	MqttBrokerAddress string
	PubKeyAddress     string
	PrivateKeyFile    string
	Username          string
}

type Global struct {
	PubkeyClient	*api.Client
	Settings		*Settings
}

var G = &Global{}