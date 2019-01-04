package mycrypto

import (
    "github.com/drep-project/drepcli/core/ethhexutil"
    "bytes"
    "github.com/pkg/errors"
    "math/big"
)

const (
    ByteLen = 32
)

func (p *Point) Bytes() []byte {
    j := make([]byte, 2 * ByteLen)
    copy(j[ByteLen - len(p.X): ByteLen], p.X)
    copy(j[2 * ByteLen - len(p.Y):], p.Y)
    return j
}

func (p *Point)SetString(hexBytesStr string)error {
    bytes, err := ethhexutil.Decode(hexBytesStr)
    if err != nil {
        return err
    }
    if len(bytes) !=  2 * ByteLen {
        return errors.New("mistake ")
    }
    p.X = make([]byte,ByteLen)
    p.Y = make([]byte,ByteLen)
    copy(p.X[:], bytes[0: ByteLen])
    copy(p.Y[:], bytes[ByteLen:])
    return nil
}

func (p *Point) Hex() string {
   return ethhexutil.Encode(p.Bytes())
}

func (p *Point) Equal(q *Point) bool {
    if !bytes.Equal(p.X, q.X) {
        return false
    }
    if !bytes.Equal(p.Y, q.Y) {
        return false
    }
    return true
}

func (p *Point) Int() (*big.Int, *big.Int) {
    return new(big.Int).SetBytes(p.X), new(big.Int).SetBytes(p.Y)
}

// MarshalText returns the hex representation of a.
func (p *Point)  MarshalText() ([]byte, error) {
    return ethhexutil.Bytes(p.Bytes()).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (p *Point)  UnmarshalText(input []byte) error {
    return p.SetString(string(input))
}

// UnmarshalJSON parses a hash in hex syntax.
func (p *Point)  UnmarshalJSON(input []byte) error {
    return p.UnmarshalText(input[1:len(input)-1])
}

func (priv *PrivateKey) Bytes() []byte{
    return priv.Prv
}

func (priv *PrivateKey) SetString(hexBytesStr string) error{
    bytes, err := ethhexutil.Decode(hexBytesStr)
    if err != nil {
        return err
    }
    cur := GetCurve()
    pubKey := cur.ScalarBaseMultiply(bytes)
    priv.Prv = bytes
    priv.PubKey = pubKey
    return nil
}

func (priv *PrivateKey) Hex() string {
    return ethhexutil.Encode(priv.Bytes())
}

// MarshalText returns the hex representation of a.
func (priv *PrivateKey)  MarshalText() ([]byte, error) {
    return ethhexutil.Bytes(priv.Prv).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (priv *PrivateKey)  UnmarshalText(input []byte) error {
    return priv.SetString(string(input))
}

// UnmarshalJSON parses a hash in hex syntax.
func (priv *PrivateKey)  UnmarshalJSON(input []byte) error {
    return priv.UnmarshalText(input[1:len(input)-1])
}
