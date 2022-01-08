package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func Encrypt(src []byte, key []byte, vec []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	ecb := cipher.NewCBCEncrypter(block, vec)
	content := src
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return crypted
}

func Decrypt(crypt []byte, key []byte, vec []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, vec)
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)
	return PKCS5Trimming(decrypted)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
