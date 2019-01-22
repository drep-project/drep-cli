package aes

import (
	"crypto/aes"
	"crypto/cipher"
)

func AesCBCEncrypt(plainText, key, iv []byte) []byte {
	key = key[:16]
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	paddingText := PKCS5Padding(plainText, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return cipherText
}

func AesCBCDecrypt(cipherText, key, iv []byte) []byte {
	key = key[:16]
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	paddingText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(paddingText, cipherText)
	plainText := PKCS5UnPadding(paddingText)
	return plainText
}
