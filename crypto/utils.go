package crypto

import (
	"errors"

	"encoding/base64"
	"encoding/pem"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

func DecodePublicKey(pemkey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemkey))
	if block == nil {
		return nil, errors.New("decode key error")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key := pub.(*rsa.PublicKey)
	return key, nil
}

func EncodeMessage(msg []byte, pubkey *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()

	output, err := rsa.EncryptOAEP(hash, rand.Reader, pubkey, msg, nil)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func EncodeTextMessage(msg string, pemPublicKey string) (string, error) {
	pubkey, err := DecodePublicKey(pemPublicKey)
	if err != nil {
		return "", err
	}

	input := []byte(msg)
	output, err := EncodeMessage(input, pubkey)
	if err != nil {
		return "", err
	}

	ciphertext := base64.StdEncoding.EncodeToString(output)
	return ciphertext, nil
}
