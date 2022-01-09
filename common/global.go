package common

import (
	"fmt"
)

type Settings struct {
	MqttBrokerAddress string
	PubKeyAddress     string
	PrivateKeyFile    string
	PublicKeyFile     string
	Username          string
	ConfigFile        string
}

var G = &Settings{}

func (s *Settings) String() string {
	vals := []struct{Key, Val string}{
		{"MQTT_BROKER_ADDRESS", s.MqttBrokerAddress},
		{"PUBKEY_ADDRESS", s.PubKeyAddress},
		{"PRIVATE_KEY_FILE", s.PrivateKeyFile},
		{"PUBLIC_KEY_FILE", s.PublicKeyFile},
		{"USERNAME", s.Username},
	}
	ret := ""
	for _, v := range vals {
        ret += fmt.Sprintf("%s = %s\n", v.Key, v.Val)
    }
	return ret
}
