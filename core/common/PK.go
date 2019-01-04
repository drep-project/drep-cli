package common

/*
import (
    "github.com/drep-project/drepcli/mycrypto"
    "encoding/hex"
)

type PK struct {
    X string    `json:"x"`
    Y string    `json:"y"`
}

func FormatPubKey(pubKey *mycrypto.Point) *PK {
    return &PK{
        X: hex.EncodeToString(pubKey.X),
        Y: hex.EncodeToString(pubKey.Y),
    }
}

func ParsePK(pk *PK) *mycrypto.Point {
    x, _ := hex.DecodeString(pk.X)
    y, _ := hex.DecodeString(pk.Y)
    return &mycrypto.Point{
        X: x,
        Y: y,
    }
}

type BootNode struct {
    PubKey  *PK         `json:"pubKey"`
    Address string      `json:"address"`
    IP      string      `json:"ip"`
    Port    int         `json:"port"`
}

*/