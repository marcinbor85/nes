package api

import (
	"crypto/rsa"
)

type KeyCache struct {
	PublicKeyMessage *rsa.PublicKey
	PublicKeySign	 *rsa.PublicKey
}

type PublicKeyMap map[string]*KeyCache

type Client struct {
	Address 		string
	ServerPublicKey	*rsa.PublicKey
	PublicKeyCache 	*PublicKeyMap
}

func NewClient(address string, serverPublicKey *rsa.PublicKey) *Client {
	client := &Client{
		Address: address,
		ServerPublicKey: serverPublicKey,
		PublicKeyCache: &PublicKeyMap{},
	}
	return client
}