package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/mr-tron/base58/base58"
)

func GenerateEncryptionKey() ([]byte, error) {
	res := make([]byte, 16)

	_, err := rand.Read(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func EncodeEncryptionKey(key []byte) string {
	return base58.FastBase58Encoding(key)
}

func DecodeEncryptionKey(key string) ([]byte, error) {
	return base58.FastBase58Decoding(key)
}

func HashEncryptionKey(key []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func EncryptBytes(key []byte, plainText []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(c, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func DecryptBytes(key []byte, cipherText []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("Ciphertext block size too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(c, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, err
}
