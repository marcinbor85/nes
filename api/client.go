package api

import (
	"crypto/rsa"
)

type PublicKeyMap map[string]*rsa.PublicKey

type Client struct {
	Address string
	PublicKeyCache *PublicKeyMap
}

func NewClient(address string) *Client {
	client := &Client{
		Address: address,
		PublicKeyCache: &PublicKeyMap{},
	}
	return client
}