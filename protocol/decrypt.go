package protocol

import (
	"encoding/base64"
	"crypto"
	"crypto/sha256"
	"crypto/rand"
	"crypto/rsa"

	"github.com/marcinbor85/pubkey/api"

	"github.com/marcinbor85/nes/crypto/aes"
)

func (frame *Frame) Decrypt(privateKey *rsa.PrivateKey, client *api.Client) (*Message, error) {
	
	// get randomKeyEncryptedEncoded from frame
	randomKeyEncryptedEncoded := frame.Cipherkey
	
	// decode messageEncryptedEncoded from base64 
	randomKeyEncrypted, err := base64.URLEncoding.DecodeString(randomKeyEncryptedEncoded)
	if err != nil {
		return nil, err
	}

	// decrypt messageEncrypted using privateKey
	randomKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, randomKeyEncrypted)
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

	// calculate hash of messageBin
	messageHash32 := sha256.Sum256(messageBin)
	messageHash := messageHash32[:]

	// get signatureEncoded from frame
	messageSignatureEncoded := frame.Signature
	
	// decode signatureEncoded from base64 
	messageSignature, err := base64.URLEncoding.DecodeString(messageSignatureEncoded)

	// verify signature using privateKey
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, messageHash, messageSignature)
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
