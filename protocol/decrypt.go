package protocol

import (
	"crypto/rsa"
	"encoding/base64"
)

func (frm *Frame) Decrypt(privateKey *rsa.PrivateKey) (*Message, error) {
	
	// get messageEncryptedEncoded from frame
	messageEncryptedEncoded := frm.Message
	
	// decode messageEncryptedEncoded from base64 
	messageEncrypted, err := base64.URLEncoding.DecodeString(messageEncryptedEncoded)
	if err != nil {
		return nil, err
	}

	// TODO: decrypt messageEncrypted using private_key (decrypting)
	msg := string(messageEncrypted)

	message, err := MessageFromString(msg)
	if err != nil {
		return nil, err
	}

	// TODO: verify sign

	return message, nil

	// // TODO: decode message using public_key of msg.To (encrypting)
	// messageEncrypted := []byte(message)

	// // encode messageEncrypted to base64
	// messageEncryptedEncoded := base64.URLEncoding.EncodeToString(messageEncrypted)

	// // calculate hash of messageEncrypted
	// hash32 := sha256.Sum256(messageEncrypted)
	// hash := hash32[:]

	// // TODO: encrypt hash using private_key of msg.From (signing)
	// hashEncrypted := []byte(hash)

	// // encode hashEncrypted to base64
	// hashEncryptedEncoded := base64.URLEncoding.EncodeToString(hashEncrypted)
	
	// // create transport frame
	// frame := &Frame{
	// 	Message: messageEncryptedEncoded,
	// 	Hash: hashEncryptedEncoded,
	// }

	// return frame, nil
}
