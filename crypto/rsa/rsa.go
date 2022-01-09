package rsa

import (
	"os"
	"bytes"
	"errors"

	"encoding/pem"

	"crypto"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
)

func DecodePublicKey(pemkey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemkey))
	if block == nil {
		return nil, errors.New("decode public key error")
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
		return nil, errors.New("decode private key error")
	}

	prv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return prv, nil
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

func Encrypt(input []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	output, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func Decrypt(input []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	output, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func Sign(message []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash32 := sha256.Sum256(message)
	hash := hash32[:]

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func Verify(message []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hash32 := sha256.Sum256(message)
	hash := hash32[:]

	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, signature)
	return err
}
