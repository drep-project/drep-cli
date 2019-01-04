package common

import (
	"math/big"
)

func CalcMemSize(off, l *big.Int) *big.Int {
	if l.Sign() == 0 {
		return Big0
	}
	return new(big.Int).Add(off, l)
}

func ToWordSize(size uint64) uint64 {
	if size > MaxUint64-31 {
		return MaxUint64/32 + 1
	}
	return (size + 31) / 32
}

func GetData(data []byte, start uint64, size uint64) []byte {
	length := uint64(len(data))
	if start > length {
		start = length
	}
	end := start + size
	if end > length {
		end = length
	}
	return RightPadBytes(data[start:end], int(size))
}

func AllZero(b []byte) bool {
	for _, byt := range b {
		if byt != 0 {
			return false
		}
	}
	return true
}

func GetDataBig(data []byte, start *big.Int, size *big.Int) []byte {
	dlen := big.NewInt(int64(len(data)))

	s := BigMin(start, dlen)
	e := BigMin(new(big.Int).Add(s, size), dlen)
	return RightPadBytes(data[s.Uint64():e.Uint64()], int(size.Uint64()))
}

func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}

func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded, slice)

	return padded
}

func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}
//
//func MustParseBig256(s string) *big.Int {
//	v, ok := ParseBig256(s)
//	if !ok {
//		panic("invalid 256 bit integer: " + s)
//	}
//	return v
//}
//
//func ParseBig256(s string) (*big.Int, bool) {
//	if s == "" {
//		return new(big.Int), true
//	}
//	var bigint *big.Int
//	var ok bool
//	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
//		bigint, ok = new(big.Int).SetString(s[2:], 16)
//	} else {
//		bigint, ok = new(big.Int).SetString(s, 10)
//	}
//	if ok && bigint.BitLen() > 256 {
//		bigint, ok = nil, false
//	}
//	return bigint, ok
//}