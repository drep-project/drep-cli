package common

import (
	"crypto/hmac"
	"crypto/sha512"
	"github.com/drep-project/drepcli/crypto/sha3"
	"math/rand"
	"net"
	"time"
)

func HmAC(message, key []byte) []byte {
	h := hmac.New(sha512.New, key)
	h.Write(message)
	return h.Sum(nil)
}

func GenUnique() ([]byte, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	uni := ""
	for _, inter := range interfaces {
		mac := inter.HardwareAddr
		uni += mac.String()
	}
	uni += time.Now().String()

	randBytes := make([]byte, 64)
	_, err = rand.Read(randBytes)
	if err != nil {
		panic("key generation: could not read from random source: " + err.Error())
	}

	return sha3.Hash256(append([]byte(uni), randBytes...)), nil
}
