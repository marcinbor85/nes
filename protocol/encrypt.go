package protocol

import (
	"encoding/base64"
	"crypto/sha256"
)

func (msg *Message) Encrypt() (*Frame, error) {
	// stringinize message struct
	message := msg.String()

	// TODO: encrypt message using public_key of msg.To (encrypting)
	messageEncrypted := []byte(message)

	// encode messageEncrypted to base64
	messageEncryptedEncoded := base64.URLEncoding.EncodeToString(messageEncrypted)

	// calculate hash of messageEncrypted
	hash32 := sha256.Sum256(messageEncrypted)
	hash := hash32[:]

	// TODO: encrypt hash using private_key of msg.From (signing)
	hashEncrypted := []byte(hash)

	// encode hashEncrypted to base64
	hashEncryptedEncoded := base64.URLEncoding.EncodeToString(hashEncrypted)
	
	// create transport frame
	frame := &Frame{
		Message: messageEncryptedEncoded,
		Hash: hashEncryptedEncoded,
	}

	return frame, nil
}
