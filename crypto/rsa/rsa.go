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

func GenerateKeysPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.Public()
	return privateKey, publicKey.(*rsa.PublicKey), nil
}

func PrivateKeyPem(key *rsa.PrivateKey) string {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
    pemBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
    )
    return string(pemBytes)
}

func PublicKeyPem(key *rsa.PublicKey) string {
	keyBytes, _ := x509.MarshalPKIXPublicKey(key)
    pemBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: keyBytes,
		},
    )
    return string(pemBytes)
}

func SavePrivateKey(key *rsa.PrivateKey, filename string) error {
	pem := PrivateKeyPem(key)

	file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

	file.WriteString(pem)
	return nil
}

func SavePublicKey(key *rsa.PublicKey, filename string) error {
	pem := PublicKeyPem(key)

	file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

	file.WriteString(pem)
	return nil
}

func DecodePublicKey(pemkey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemkey))
	if block == nil {
		return nil, errors.New("decode public key error")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func DecodePrivateKey(pemkey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemkey))
	if block == nil {
		return nil, errors.New("decode private key error")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func LoadPublicKey(filename string) (*rsa.PublicKey, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	bytesBuf := &bytes.Buffer{}
	bytesBuf.ReadFrom(file)
	pem := bytesBuf.String()

	key, err := DecodePublicKey(pem)
	if err != nil {
		return nil, "", err
	}

	return key, pem, nil
}

func LoadPrivateKey(filename string) (*rsa.PrivateKey, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	bytesBuf := &bytes.Buffer{}
	bytesBuf.ReadFrom(file)
	pem := bytesBuf.String()

	key, err := DecodePrivateKey(pem)
	if err != nil {
		return nil, "", err
	}

	return key, pem, nil
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

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hash, &opts)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func Verify(message []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hash32 := sha256.Sum256(message)
	hash := hash32[:]

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto

	err := rsa.VerifyPSS(publicKey, crypto.SHA256, hash, signature, &opts)
	return err
}