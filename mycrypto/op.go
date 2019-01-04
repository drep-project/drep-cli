package mycrypto

import (
    "math/big"
    "crypto/rand"
    "errors"
)

const (
    MaximumGenerateKeyRetry = 100
)

var curveInstance *CurveParams

func init() {
    curveParams := &CurveParams{}
    curveParams.P = new(big.Int).SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
        0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFC, 0x2F})
    curveParams.N = new(big.Int).SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
        0xFF, 0xFF, 0xFF, 0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x41})
    curveParams.B = new(big.Int).SetBytes([]byte{0x07})
    Gx := []byte{0x79, 0xBE, 0x66, 0x7E, 0xF9, 0xDC, 0xBB, 0xAC, 0x55, 0xA0, 0x62, 0x95, 0xCE, 0x87, 0x0B, 0x07, 0x02, 0x9B, 0xFC,
        0xDB, 0x2D, 0xCE, 0x28, 0xD9, 0x59, 0xF2, 0x81, 0x5B, 0x16, 0xF8, 0x17, 0x98}
    Gy := []byte{0x48, 0x3A, 0xDA, 0x77, 0x26, 0xA3, 0xC4, 0x65, 0x5D, 0xA4, 0xFB, 0xFC, 0x0E, 0x11, 0x08, 0xA8, 0xFD, 0x17, 0xB4,
        0x48, 0xA6, 0x85, 0x54, 0x19, 0x9C, 0x47, 0xD0, 0x8F, 0xFB, 0x10, 0xD4, 0xB8}
    curveParams.G = &Point{X: Gx, Y: Gy}
    curveParams.BitSize = 256
    curveParams.Name = "Secp256-k1"
    curveInstance = curveParams
}

func GetCurve() *CurveParams {
    return curveInstance
}
func GetRandomKQ() ([]byte, *Point, error) {
    curve := curveInstance
    mask := []byte{0xff, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f}
    N := curve.Params().N
    BitSize := curve.Params().BitSize
    byteLen := (BitSize + 7) >> 3
    ok := false
    try := 0
    var k []byte
    for !ok {
        if try > MaximumGenerateKeyRetry {
            break
        }
        k = make([]byte, byteLen)
        if _, err := rand.Read(k); err != nil {
            try += 1
            continue
        }
        k[0] &= mask[BitSize % 8]
        kInt := new(big.Int).SetBytes(k)
        if kInt.Cmp(Zero) > 0 || kInt.Cmp(N) < 0 {
            ok = true
        } else {
            try += 1
        }
    }
    if ok {
        return k, curve.ScalarBaseMultiply(k), nil
    } else {
        return nil, nil, errors.New("random fail")
    }
}

func GeneratePrivateKey() (*PrivateKey, error) {
    prv, pubKey, err0 := GetRandomKQ()
    if err0 != nil {
        return nil, errors.New("error")
    } else {
        return &PrivateKey{Prv: prv, PubKey: pubKey}, nil
    }
}

func Sign(prvKey *PrivateKey, b []byte) (*Signature, error) {
    curve := curveInstance
    r, s := new(big.Int), new(big.Int)
    prvInt := new(big.Int).SetBytes(prvKey.Prv)
    for r.Cmp(Zero) == 0 || s.Cmp(Zero) == 0 {
        k, Q, err := GetRandomKQ()
        if err != nil {
            return nil, err
        }
        N := curve.Params().N
        r = new(big.Int).SetBytes(ConcatHash256(Q.Bytes(), prvKey.PubKey.Bytes(), b))
        r.Mod(r, N)
        s = new(big.Int).Mul(r, prvInt)
        s.Mod(s, N)
        s.Sub(new(big.Int).SetBytes(k), s)
        s.Mod(s, N)
    }
    sig := &Signature{}
    sig.R = r.Bytes()
    sig.S = s.Bytes()
    return sig, nil
}

func Verify(sig *Signature, pubKey *Point, b []byte) bool {
    //fmt.Println("Validate 3")
    curve := curveInstance
    r, s := new(big.Int).SetBytes(sig.R), new(big.Int).SetBytes(sig.S)
    if r.Cmp(Zero) <= 0 || r.Cmp(curve.Params().N) >= 0 || s.Cmp(Zero) <=0 || s.Cmp(curve.Params().N) >=0 {
        //fmt.Println("Validate 4", r.Cmp(Zero) <= 0, r.Cmp(curve.Params().N) >= 0 , s.Cmp(Zero) <=0 , s.Cmp(curve.Params().N) >=0 )
        return false
    }
    N := curve.Params().N
    sG := curve.ScalarBaseMultiply(sig.S)
    rP := curve.ScalarMultiply(pubKey, sig.R)
    Q:= curve.Add(sG, rP)
    Qx, Qy := Q.Int()
    //fmt.Println("Validate 5")
    if Qx.Cmp(Zero) == 0 && Qy.Cmp(Zero) == 0 {
        //fmt.Println("Validate 6")
        return false
    }
    v := new(big.Int).SetBytes(ConcatHash256(Q.Bytes(), pubKey.Bytes(), b))
    v.Mod(v, N)
    //fmt.Println("Validate 7")
    if v.Cmp(r) == 0{
        //fmt.Println("Validate 8")
        return true
    } else {
        //fmt.Println("Validate 9")
        return false
    }
}

func Encrypt(pubKey *Point, b []byte) ([]byte, error) {
    //curve := GetCurve()
    //k, p1, err := GetRandomKQ()
    //if err != nil {
    //    return nil, err
    //}
    //c1 := p1.Bytes()
    //p2 := curve.ScalarMultiply(pubKey, k)
    //j2 := p2.Bytes()
    //t := new(big.Int).SetBytes(hash.KDF(j2))
    //m := new(big.Int).SetBytes(b)
    //c2 := new(big.Int).Xor(m, t).Bytes()
    //buf := make([]byte, len(j2) + len(b))
    //copy(buf[:len(j2)], j2)
    //copy(buf[len(j2):], b)
    //c3 := hash.Hash256(buf)
    //cipher := make([]byte, 3 * hash.ByteLen + len(c2))
    //copy(cipher[: 2 * hash.ByteLen], c1)
    //copy(cipher[2 * hash.ByteLen: 3 * hash.ByteLen], c3)
    //copy(cipher[3 * hash.ByteLen:], c2)
    //return cipher, nil
    return b, nil
}

func Decrypt(cipher []byte) ([]byte, error) {
    //curve := GetCurve()
    //prvKey, err := GetPrivateKey()
    //if err != nil {
    //    return nil, err
    //}
    //p1 := &bean.Point{X: cipher[:hash.ByteLen], Y: cipher[hash.ByteLen: 2 * hash.ByteLen]}
    //if !curve.IsOnCurve(p1) {
    //    return nil, errors.New("point not on CurveInstance")
    //}
    //p2 := curve.ScalarMultiply(p1, prvKey.Prv)
    //j2 := p2.Bytes()
    //t := new(big.Int).SetBytes(hash.KDF(j2))
    //c2 := cipher[3 * hash.ByteLen:]
    //c := new(big.Int).SetBytes(c2)
    //m := new(big.Int).Xor(c, t)
    //msg := m.Bytes()
    //b := make([]byte, len(j2) + len(msg))
    //copy(b[:len(j2)], j2)
    //copy(b[len(j2):], msg)
    //u := hash.Hash256(b)
    //c3 := cipher[2 * hash.ByteLen: 3 * hash.ByteLen]
    //if bytes.Equal(u, c3) {
    //    return msg, nil
    //} else {
    //    return nil, errors.New("cipher wrong")
    //}
    return cipher, nil
}