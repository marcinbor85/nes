package protocol

import (
	"encoding/base64"
	"crypto"
	"crypto/sha256"
	"crypto/rand"
	"crypto/rsa"

	"github.com/marcinbor85/pubkey/api"
)

func (frame *Frame) Decrypt(privateKey *rsa.PrivateKey, client *api.Client) (*Message, error) {
	
	// get messageEncryptedEncoded from frame
	messageEncryptedEncoded := frame.Ciphertext
	
	// decode messageEncryptedEncoded from base64 
	messageEncrypted, err := base64.URLEncoding.DecodeString(messageEncryptedEncoded)
	if err != nil {
		return nil, err
	}

	// decrypt messageEncrypted using privateKey
	messageBin, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, messageEncrypted)
	if err != nil {
		return nil, err
	}

	messageString := string(messageBin)

	message, err := MessageFromString(messageString)
	if err != nil {
		return nil, err
	}

	publicKey, ee := client.GetPublicKeyByUsername(message.From)
	if ee != nil {
		return nil, ee
	}

	// calculate hash of messageBin
	hash32 := sha256.Sum256(messageBin)
	hash := hash32[:]

	// get signatureEncoded from frame
	signatureEncoded := frame.Signature
	
	// decode signatureEncoded from base64 
	signature, err := base64.URLEncoding.DecodeString(signatureEncoded)

	// verify signature using privateKey
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, signature)
	if err != nil {
		return nil, err
	}

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
