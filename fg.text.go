package main

import (
	"BlockChainTest/core/ethcrypto/sha3"
	"math/big"

	"crypto/ecdsa"
	"fmt"
	"github.com/drep-project/drepcli/mycrypto"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

func main(){
	randSign := "1111111111111111"
	pri1 ,_ := mycrypto.GeneratePrivateKey()
	//bytes := pri1.Prv
	//aaaaaa := crypto.S256().IsOnCurve(new (big.Int).SetBytes(pri1.PubKey.X), new (big.Int).SetBytes(pri1.PubKey.Y))
	hexStr := pri1.Hex()

	crypto.HexToECDSA(hexStr[2:])

	pri2, _ := crypto.ToECDSA(pri1.Prv)
	data  := []byte{1,2,3,4,5}
	hash :=  sha3.Hash256(data)

	r,s, _ := ecdsa.Sign(strings.NewReader(randSign), pri2, hash)
	fmt.Println(r)
	fmt.Println(s)

	xxx := ecdsa.Verify(&pri2.PublicKey,hash,r,s)
	fmt.Println(xxx)

	sign2, _ := mycrypto.Sign(pri1, data)
	fmt.Println(sign2.R)
	fmt.Println(sign2.S)

	xxx = ecdsa.Verify(&pri2.PublicKey, hash, new(big.Int).SetBytes(sign2.R), new(big.Int).SetBytes(sign2.S))
	fmt.Println(xxx)

	xxx = mycrypto.Verify(sign2, pri1.PubKey, data)
	fmt.Println(xxx)

	sigTemp := &mycrypto.Signature{R:r.Bytes(),S:s.Bytes()}
	xxx = mycrypto.Verify(sigTemp, pri1.PubKey, data)
	fmt.Println(xxx)
//	priv := ecdsa.PrivateKey{
//		PublicKey:pub1,
//		D: new (big.Int).SetBytes(pri1.Bytes()),
//	}
}