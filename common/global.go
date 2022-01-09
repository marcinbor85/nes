package common

import (
	"fmt"
)

type Settings struct {
	MqttBrokerAddress string
	PubKeyAddress     string
	PrivateKeyFile    string
	Username          string
}

var G = &Settings{}

func (s *Settings) String() string {
	vals := []struct{Key, Val string}{
		{"MQTT_BROKER_ADDRESS", s.MqttBrokerAddress},
		{"PUBKEY_ADDRESS", s.PubKeyAddress},
		{"PRIVATE_KEY_FILE", s.PrivateKeyFile},
		{"USERNAME", s.Username},
	}
	ret := ""
	for _, v := range vals {
        ret += fmt.Sprintf("%s = %s\n", v.Key, v.Val)
    }
	return ret
}
