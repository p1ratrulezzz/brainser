package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"
)

var randomBytes = pbkdf2.Key([]byte("3E4pEj2NDzzrCbxfvTqta9uZ"), []byte("salt"), 1000, 16, sha1.New)
var secretBytes = pbkdf2.Key([]byte("sdsgWgQcyqaGrstXAxfFsMYmKyhA9eLwsaHzRbGxZq"), []byte("salt"), 1000, 24, sha1.New)
var secret = string(secretBytes)

func Encrypt(data []byte) []byte {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		panic("Error" + err.Error())
	}
	cfb := cipher.NewCFBEncrypter(block, randomBytes)
	cipherText := make([]byte, len(data))
	cfb.XORKeyStream(cipherText, data)
	return cipherText
}

func Decrypt(data []byte) []byte {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		panic("Error" + err.Error())
	}

	if len(data) == 0 {
		panic("encrypted data is empty")
	}

	cfb := cipher.NewCFBDecrypter(block, randomBytes)
	plainData := make([]byte, len(data))
	cfb.XORKeyStream(plainData, data)
	return plainData
}
