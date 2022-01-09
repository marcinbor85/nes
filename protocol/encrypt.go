package protocol

import (
	"encoding/base64"
	"crypto/rand"
	"crypto/rsa"

	"github.com/marcinbor85/nes/crypto/aes"
	r "github.com/marcinbor85/nes/crypto/rsa"
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
	randomKeyEncrypted, err := r.Encrypt(randomKey, publicKey)
	if err != nil {
		return nil, err
	}

	// encode randomKeyEncrypted to base64
	randomKeyEncryptedEncoded := base64.URLEncoding.EncodeToString(randomKeyEncrypted)

	// sign messageBin hash using privateKey
	messageSignature, err := r.Sign(messageBin, privateKey)
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
