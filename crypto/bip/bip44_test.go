package bip

import (
	"fmt"
	"github.com/btcsuite/btcutil/hdkeychain"
	"math/rand"
	"testing"
)

func TestBip39(t *testing.T) {
	hdkeychain.GenerateSeed(uint8(10))
	extKey, _ := hdkeychain.NewKeyFromString("xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8")
	xxx, _ := extKey.Child(10)
	fmt.Println(xxx)
	yyy, _ := xxx.Child(10)
	fmt.Println(yyy)
}

func TestBip392(t *testing.T) {
	token := make([]byte, 32)
	rand.Read(token)

	mnemonic, _ := NewMnemonic(token)
	masterKey, _ := NewKeyFromMnemonic(mnemonic, "111111", 50, 0, 0, 0)
	//keys,_ := NewKeyFromMasterKey(masterKey,50,0,1,0)
	pub := masterKey.PublicKey()
	pub2, _ := pub.NewChildKey(1000)
	fmt.Println(pub2)

}
