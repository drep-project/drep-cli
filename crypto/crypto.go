// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package crypto

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/drep-project/drepcli/common"
	"github.com/drep-project/drepcli/crypto/secp256k1"
	"github.com/drep-project/drepcli/crypto/sha3"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))

	errInvalidPubkey = errors.New("invalid secp256k1 public key")
)

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h Hash) {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress(b CommonAddress, nonce uint64) CommonAddress {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, nonce)
	data := append(b.Bytes(), buf...)
	return Bytes2Address(Keccak256(data)[12:])
}

// CreateAddress2 creates an ethereum address given the address bytes, initial
// contract code and a salt.
func CreateAddress2(b CommonAddress, salt [32]byte, code []byte) CommonAddress {
	return Bytes2Address(Keccak256([]byte{0xff}, b.Bytes(), salt[:], Keccak256(code))[12:])
}

// LoadECDSA loads a secp256k1 private key from the given file.
func LoadECDSA(file string) (*secp256k1.PrivateKey, error) {
	buf := make([]byte, 64)
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	if _, err := io.ReadFull(fd, buf); err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(string(buf))
	if err != nil {
		return nil, err
	}
	return secp256k1.NewPrivateKey(new(big.Int).SetBytes(key)), nil
}

// SaveECDSA saves a secp256k1 private key to the given file with
// restrictive permissions. The key data is saved hex-encoded.
func SaveECDSA(file string, key *secp256k1.PrivateKey) error {
	k := hex.EncodeToString(key.Serialize())
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// Generate an elliptic curve public / private keypair. If params is nil,
// the recommended default parameters for the key will be chosen.
func GenerateKey(rand io.Reader) (prv *secp256k1.PrivateKey, err error) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		return nil, err
	}
	return (*secp256k1.PrivateKey)(key), nil
}

func PubkeyToAddress(p *secp256k1.PublicKey) CommonAddress {
	return Bytes2Address(Keccak256(p.Serialize()[1:])[12:])
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
		return false
	}
	// reject upper range of s values (ECDSA malleability)
	// see discussion in secp256k1/libsecp256k1/include/secp256k1.h
	if homestead && s.Cmp(secp256k1halfN) > 0 {
		return false
	}
	// Frontier: allow s to be in full N range
	return r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 && (v == 0 || v == 1)
}
