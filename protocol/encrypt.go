package protocol

import (
	"encoding/base64"
	"crypto"
	"crypto/sha256"
	"crypto/rand"
	"crypto/rsa"

	"github.com/marcinbor85/nes/crypto/aes"
)

func (message *Message) Encrypt(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (*Frame, error) {
	// stringinize message struct
	messageString := message.String()

	// binarize messageString
	messageBin := []byte(messageString)

	// generate random key and vector
	randomKey := make([]byte, 48)
	rand.Read(randomKey)

	// encrypt message with randomKey
	messageEncrypted := aes.Encrypt(messageBin, randomKey[:32], randomKey[32:])

	// encode messageEncrypted to base64
	messageEncryptedEncoded := base64.URLEncoding.EncodeToString(messageEncrypted)

	// encrypt randomKey using publicKey
	randomKeyEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, randomKey)
	if err != nil {
		return nil, err
	}

	// encode randomKeyEncrypted to base64
	randomKeyEncryptedEncoded := base64.URLEncoding.EncodeToString(randomKeyEncrypted)

	// calculate hash of messageBin
	messageHash32 := sha256.Sum256(messageBin)
	messageHash := messageHash32[:]

	// sign messageBin hash using privateKey
	messageSignature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, messageHash)
	if err != nil {
		return nil, err
	}

	// encode signature to base64
	messageSignatureEncoded := base64.URLEncoding.EncodeToString(messageSignature)
	
	// create transport frame
	frame := &Frame{
		Cipherkey: randomKeyEncryptedEncoded,
		Ciphertext: messageEncryptedEncoded,
		Signature: messageSignatureEncoded,
	}

	return frame, nil
}
