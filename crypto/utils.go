package crypto

import (
	"os"
	"bytes"
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

func DecodePrivateKey(pemkey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemkey))
	if block == nil {
		return nil, errors.New("decode key error")
	}

	prv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return prv, nil
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

func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytesBuf := &bytes.Buffer{}
	bytesBuf.ReadFrom(file)
	privateKeyPem := bytesBuf.String()

	privateKey, err := DecodePrivateKey(privateKeyPem)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
