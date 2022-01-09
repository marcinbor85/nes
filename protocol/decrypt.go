package protocol

import (
	"encoding/base64"

	"crypto/rsa"

	"github.com/marcinbor85/pubkey/api"

	"github.com/marcinbor85/nes/crypto/aes"
	r "github.com/marcinbor85/nes/crypto/rsa"
)

func (frame *Frame) Decrypt(privateKey *rsa.PrivateKey, client *api.Client) (*Message, error) {
	
	// get randomKeyEncryptedEncoded from frame
	randomKeyEncryptedEncoded := frame.Cipherkey
	
	// decode messageEncryptedEncoded from base64 
	randomKeyEncrypted, err := base64.URLEncoding.DecodeString(randomKeyEncryptedEncoded)
	if err != nil {
		return nil, err
	}

	// decrypt randomKeyEncrypted using privateKey
	randomKey, err := r.Decrypt(randomKeyEncrypted, privateKey)
	if err != nil {
		return nil, err
	}

	// get messageEncryptedEncoded from frame
	messageEncryptedEncoded := frame.Ciphertext

	// decode messageEncryptedEncoded from base64 
	messageEncrypted, err := base64.URLEncoding.DecodeString(messageEncryptedEncoded)
	if err != nil {
		return nil, err
	}

	// decrypt message
	messageBin := aes.Decrypt(messageEncrypted, randomKey[:32], randomKey[32:])

	messageString := string(messageBin)

	message, err := MessageFromString(messageString)
	if err != nil {
		return nil, err
	}

	publicKey, ee := client.GetPublicKeyByUsername(message.From)
	if ee != nil {
		return nil, ee
	}

	// get signatureEncoded from frame
	messageSignatureEncoded := frame.Signature
	
	// decode signatureEncoded from base64 
	messageSignature, err := base64.URLEncoding.DecodeString(messageSignatureEncoded)

	// verify signature using privateKey
	err = r.Verify(messageBin, messageSignature, publicKey)
	if err != nil {
		return nil, err
	}

	return message, nil
}
