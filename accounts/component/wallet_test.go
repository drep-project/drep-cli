package component

import (
	accountTypes "github.com/drep-project/drepcli/accounts/types"
	"BlockChainTest/core/crypto/secp256k1"
	"BlockChainTest/core/crypto/sha3"
	"fmt"
	"log"
	"testing"
)

func TestWallet(t *testing.T) {
	password := string(sha3.Hash256([]byte("AAAAAAAAAAAAAAAA")))

	newNode := accountTypes.NewNode(nil, accountTypes.RootChain)
	fmt.Println(newNode)

	accountConfig := &accountTypes.Config{
		KeyStoreDir:"TestWallet",
	}
	wallet, err := NewWallet(accountConfig, accountTypes.RootChain)
	wallet.chainId = accountTypes.RootChain
	if err != nil {
		log.Fatal("NewWallet error")
	}

	err = wallet.Open(password)
	if err != nil {
		log.Fatal("open wallet error")
	}

	nodes := []*accountTypes.Node{}
	for i:=0;i<10;i++ {
		node,err := wallet.NewAccount()
		pk := node.PrivateKey.PubKey()
		isOnCurve := secp256k1.S256().IsOnCurve(pk.X,pk.Y)
		if !isOnCurve {
			log.Fatal("error public key")
		}
		if err != nil {
			log.Fatal("open wallet error")
		}
		nodes = append(nodes, node)
	}

	wallet.Lock()
	_, err = wallet.NewAccount()
	if err == nil {
		log.Fatal("Lock not effect")
	}

	wallet.UnLock(password)
	_, err = wallet.NewAccount()
	if err != nil {
		log.Fatal("UnLock not effect")
	}

	wallet.Close()

	wallet.Open(password)

	for _, node := range  nodes {
		reloadNode, err := wallet.GetAccountByAddress(node.Address)
		if err != nil {
			log.Fatal("reload wallet error")
		}
		pk := reloadNode.PrivateKey.PubKey()
		isOnCurve := secp256k1.S256().IsOnCurve(pk.X,pk.Y)
		if !isOnCurve {
			log.Fatal("error public key")
		}
		if reloadNode.PrivateKey == nil {
			log.Fatal("privateKey wallet error")
		}
	}
}