package aes

import "bytes"

func PKCS5Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	newText := append(plainText, padText...)
	return newText
}

func PKCS5UnPadding(plainText []byte) []byte {
	length := len(plainText)
	number := int(plainText[length-1])
	return plainText[:length-number]
}
