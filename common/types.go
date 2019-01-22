package common

import (
	"encoding/hex"
)

const (
	ChainIdSize = 64
)

type ChainIdType [ChainIdSize]byte

func (c ChainIdType) Hex() string {
	return hex.EncodeToString(c[:])
}

func (c *ChainIdType) SetBytes(b []byte) {
	if len(b) > len(c) {
		copy(c[:], b[len(b)-ChainIdSize:])
	} else {
		copy(c[ChainIdSize-len(b):], b)
	}
}

func Bytes2ChainId(b []byte) ChainIdType {
	if b == nil {
		return ChainIdType{}
	}
	var chainId ChainIdType
	chainId.SetBytes(b)
	return chainId
}

func Hex2ChainId(s string) ChainIdType {
	if s == "" {
		return ChainIdType{}
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return ChainIdType{}
	}
	return Bytes2ChainId(b)
}

func (c ChainIdType) MarshalText() ([]byte, error) {
	return []byte(c.Hex()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ChainIdType) UnmarshalJSON(input []byte) error {
	return c.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler
func (c *ChainIdType) UnmarshalText(input []byte) error {
	c.SetBytes(input)
	return nil
}
