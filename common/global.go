package common

import (
	"fmt"

	"crypto/rsa"
)

type Settings struct {
	MqttBrokerAddress 		string
	PubKeyAddress     		string
	PubKeyPublicKeyFile		string
	PrivateKeyMessageFile   string
	PublicKeyMessageFile    string
	PrivateKeySignFile    	string
	PublicKeySignFile     	string
	Username          		string
	ConfigFile        		string
	PubKeyPublicKey			*rsa.PublicKey
}

var G = &Settings{}

func (s *Settings) String() string {
	vals := []struct{Key, Val string}{
		{"MQTT_BROKER_ADDRESS", s.MqttBrokerAddress},
		{"PUBKEY_ADDRESS", s.PubKeyAddress},
		{"PUBKEY_PUBLIC_KEY_FILE", s.PubKeyPublicKeyFile},
		{"PRIVATE_KEY_MESSAGE_FILE", s.PrivateKeyMessageFile},
		{"PUBLIC_KEY_MESSAGE_FILE", s.PublicKeyMessageFile},
		{"PRIVATE_KEY_SIGN_FILE", s.PrivateKeySignFile},
		{"PUBLIC_KEY_SIGN_FILE", s.PublicKeySignFile},
		{"NES_USERNAME", s.Username},
	}
	ret := ""
	for _, v := range vals {
        ret += fmt.Sprintf("%s = %s\n", v.Key, v.Val)
    }
	return ret
}
