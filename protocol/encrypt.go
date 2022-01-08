package protocol

import (
	"encoding/base64"
	"crypto"
	"crypto/sha256"
	"crypto/rand"
	"crypto/rsa"
)

func (message *Message) Encrypt(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (*Frame, error) {
	// stringinize message struct
	messageString := message.String()

	// convert to binary data
	messageBin := []byte(messageString)

	// encrypt message using publicKey
	messageEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, messageBin)
	if err != nil {
		return nil, err
	}

	// encode messageEncrypted to base64
	messageEncryptedEncoded := base64.URLEncoding.EncodeToString(messageEncrypted)

	// calculate hash of messageBin
	hash32 := sha256.Sum256(messageBin)
	hash := hash32[:]

	// sign messageBin hash using privateKey
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		return nil, err
	}

	// encode signature to base64
	signatureEncoded := base64.URLEncoding.EncodeToString(signature)
	
	// create transport frame
	frame := &Frame{
		Ciphertext: messageEncryptedEncoded,
		Signature: signatureEncoded,
	}

	return frame, nil
}
