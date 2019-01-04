package mycrypto

import (
	"math/big"
)

var Zero = new(big.Int)

type Curve interface {
	Params() *CurveParams
	IsOnCurve(*Point) bool
	Add(*Point, *Point) *Point
	Double(*Point) *Point
	ScalarMultiply(*Point, []byte) *Point
	ScalarBaseMultiply([]byte) *Point
}

// Y^2 == X^3 + AX + B (mod p), with a == 0
type CurveParams struct {
	P *big.Int
	N *big.Int
	B *big.Int
	G *Point
	BitSize int
	Name string
}

type JacobiCoordinate struct {
	X *big.Int
	Y *big.Int
	Z *big.Int
}

func (curveParams *CurveParams) Params() *CurveParams {
	return curveParams
}

// Y^2 == X^3 + 7 (mod p)
func (curveParams *CurveParams) IsOnCurve(point *Point) bool {
	x, y := point.Int()
	P := curveParams.P
	B := curveParams.B
	ySquare := new(big.Int).Mul(y, y)
	ySquare.Mod(ySquare, P)
	xCube := new(big.Int).Mul(x, x)
	xCube.Mod(xCube, P)
	xCube.Mul(xCube, x)
	xPolynomial := new(big.Int).Add(xCube, B)
	xPolynomial.Mod(xPolynomial, P)
	return ySquare.Cmp(xPolynomial) == 0
}

func JacobiAffine(point *Point) *JacobiCoordinate {
	x, y := point.Int()
	z := new(big.Int)
	if x.Sign() != 0 || y.Sign() != 0 {
		z.SetInt64(1)
	}
	return &JacobiCoordinate{x, y, z}
}

func (curveParams *CurveParams) InverseJacobiAffine(jc *JacobiCoordinate) *Point {
	x, y, z := jc.X, jc.Y, jc.Z
	if z.Sign() == 0 {
		return &Point{X: new(big.Int).Bytes(), Y: new(big.Int).Bytes()}
	}
	P := curveParams.P
	zInv := new(big.Int).ModInverse(z, P)
	zInvSquare := new(big.Int).Mul(zInv, zInv)
	zInvSquare.Mod(zInvSquare, P)
	zInvCube := new(big.Int).Mul(zInvSquare, zInv)
	zInvCube.Mod(zInvCube, P)
	xOut := new(big.Int).Mul(x, zInvSquare)
	xOut.Mod(xOut, P)
	yOut := new(big.Int).Mul(y, zInvCube)
	yOut.Mod(yOut, P)
	return &Point{X: xOut.Bytes(), Y: yOut.Bytes()}
}

// add-2007-bl addition
// Cost: 11M + 5S + 9add + 4*2
// Cost: 10M + 4S + 9add + 4*2 dependent upon the first point
// Source: 2007 Bernstein–Lange; note that the improvement from 12M+4S to 11M+5S was already mentioned in 2001 Bernstein http://cr.yp.to/talks.html#2001.10.29
// Explicit formulas:
// Explicit formulas:
//      Z1Z1 = Z1^2
//      Z2Z2 = Z2^2
//      U1 = X1*Z2Z2
//      U2 = X2*Z1Z1
//      S1 = Y1*Z2*Z2Z2
//      S2 = Y2*Z1*Z1Z1
//      H = U2-U1
//      I = (2*H)^2
//      J = H*I
//      R = 2*(S2-S1)
//      V = U1*I
//      X3 = R^2-J-2*V
//      Y3 = R*(V-X3)-2*S1*J
//      Z3 = ((Z1+Z2)^2-Z1Z1-Z2Z2)*H
func (curveParams *CurveParams) JacobiAddition(jc1, jc2 *JacobiCoordinate) *JacobiCoordinate {
	x1, y1, z1 := jc1.X, jc1.Y, jc1.Z
	x2, y2, z2 := jc2.X, jc2.Y, jc2.Z
	x3, y3, z3 := new(big.Int), new(big.Int), new(big.Int)
	if z1.Sign() == 0 {
		return jc2
	}
	if z2.Sign() == 0 {
		return jc1
	}
	P := curveParams.P
	z1Square := new(big.Int).Mul(z1, z1)
	z1Square.Mod(z1Square, P)
	z1Cube := new(big.Int).Mul(z1Square, z1)
	z1Cube.Mod(z1Cube, P)
	z2Square := new(big.Int).Mul(z2, z2)
	z2Square.Mod(z2Square, P)
	z2Cube := new(big.Int).Mul(z2Square, z2)
	z2Cube.Mod(z2Cube, P)
	u1 := new(big.Int).Mul(x1, z2Square)
	u1.Mod(u1, P)
	u2 := new(big.Int).Mul(x2, z1Square)
	u2.Mod(u2, P)
	s1 := new(big.Int).Mul(y1, z2Cube)
	s1.Mod(s1, P)
	s2 := new(big.Int).Mul(y2, z1Cube)
	s2.Mod(s2, P)
	h := new(big.Int).Sub(u2, u1)
	h.Mod(h, P)
	i := new(big.Int).Lsh(h, 1)
	i.Mul(i, i)
	i.Mod(i, P)
	j := new(big.Int).Mul(h, i)
	j.Mod(j, P)
	r := new(big.Int).Sub(s2, s1)
	r.Lsh(r, 1)
	r.Mod(r, P)
	v := new(big.Int).Mul(u1, i)
	v.Mod(v, P)
	rSquare := new(big.Int).Mul(r, r)
	rSquare.Mod(rSquare, P)
	vDouble := new(big.Int).Lsh(v, 1)
	x3.Add(x3, rSquare)
	x3.Sub(x3, j)
	x3.Sub(x3, vDouble)
	x3.Mod(x3, P)
	y3Item1 := new(big.Int).Sub(v, x3)
	y3Item1.Mul(y3Item1, r)
	y3Item1.Mod(y3Item1, P)
	y3Item2 := new(big.Int).Lsh(s1, 1)
	y3Item2.Mul(y3Item2, j)
	y3Item2.Mod(y3Item2, P)
	y3.Sub(y3Item1, y3Item2)
	y3.Mod(y3, P)
	z3.Mul(z1, z2)
	z3.Mul(z3, h)
	z3.Lsh(z3, 1)
	z3.Mod(z3, P)
	return &JacobiCoordinate{x3, y3, z3}
}

// dbl-2007-bl
// Cost: 1M + 8S + 1*A + 10add + 2*2 + 1*3 + 1*8
// Source: 2007 Bernstein–Lange
// Explicit formulas:
//      XX = X1^2
//      YY = Y1^2
//      YYYY = YY^2
//      ZZ = Z1^2
//      S = 2*((X1+YY)^2-XX-YYYY)
//      M = 3*XX
//      T = M^2-2*S
//      X3 = T
//      Y3 = M*(S-T)-8*YYYY
//      Z3 = (Y1+Z1)^2-YY-ZZ
func (curveParams *CurveParams) JacobiDoubling(jc *JacobiCoordinate) *JacobiCoordinate {
	x1, y1, z1 := jc.X, jc.Y, jc.Z
	x2, y2, z2 := new(big.Int), new(big.Int), new(big.Int)
	if z1.Sign() == 0 {
		return jc
	}
	P := curveParams.P
	x1Square := new(big.Int).Mul(x1, x1)
	x1Square.Mod(x1Square, P)
	y1Square := new(big.Int).Mul(y1, y1)
	y1Square.Mod(y1Square, P)
	y1Biquadratic := new(big.Int).Mul(y1Square, y1Square)
	y1Biquadratic.Mod(y1Biquadratic, P)
	z1Square := new(big.Int).Mul(z1, z1)
	z1Square.Mod(z1Square, P)
	s := new(big.Int).Mul(x1, y1Square)
	s.Lsh(s, 2)
	s.Mod(s, P)
	m := new(big.Int).Lsh(x1Square, 1)
	m.Add(m, x1Square)
	m.Mod(m, P)
	t := new(big.Int).Mul(m, m)
	sDouble := new(big.Int).Lsh(s, 1)
	t.Sub(t, sDouble)
	t.Mod(t, P)
	x2.Set(t)
	y2Item1 := new(big.Int).Sub(s, t)
	y2Item1.Mul(y2Item1, m)
	y2Item1.Mod(y2Item1, P)
	y2Item2 := new(big.Int).Lsh(y1Biquadratic, 3)
	y2.Sub(y2Item1, y2Item2)
	y2.Mod(y2, P)
	z2.Mul(y1, z1)
	z2.Lsh(z2, 1)
	z2.Mod(z2, P)
	return &JacobiCoordinate{x2, y2, z2}
}

func (curveParams *CurveParams) Add(pt1, pt2 *Point) *Point {
	jc1 := JacobiAffine(pt1)
	jc2 := JacobiAffine(pt2)
	return curveParams.InverseJacobiAffine(curveParams.JacobiAddition(jc1, jc2))
}

func (curveParams *CurveParams) Double(point *Point) *Point {
	jc := JacobiAffine(point)
	return curveParams.InverseJacobiAffine(curveParams.JacobiDoubling(jc))
}

func (curveParams *CurveParams) ScalarMultiply(point *Point, k []byte) *Point {
	jc0 := JacobiAffine(point)
	jc := &JacobiCoordinate{new(big.Int), new(big.Int), new(big.Int)}
	for _, byt := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			jc = curveParams.JacobiDoubling(jc)
			if byt & 0x80 == 0x80 {
				jc = curveParams.JacobiAddition(jc, jc0)
			}
			byt <<= 1
		}
	}
	return curveParams.InverseJacobiAffine(jc)
}

func (curveParams *CurveParams) ScalarBaseMultiply(k []byte) *Point {
	return curveParams.ScalarMultiply(curveParams.G, k)
}

func (curveParams *CurveParams) ScalarBaseMultiplyByFormula(k int) *Point {
	Gx, Gy := curveParams.G.Int()
	Fx, Fy := new(big.Int).Set(Gx), new(big.Int).Set(Gy)
	P := curveParams.P
	for i := 1; i < k; i ++ {
		x1 := new(big.Int).Mul(Fx, Fx)
		x1.Mod(x1, P)
		x2 := new(big.Int).Mul(Gx, Gx)
		x2.Mod(x2, P)
		x3 := new(big.Int).Mul(Fx, Gx)
		x3.Mod(x3, P)
		x4 := new(big.Int).Add(Fy, Gy)
		x4.ModInverse(x4, P)
		t := new(big.Int).Set(x1)
		t.Add(t, x2)
		t.Add(t, x3)
		t.Mul(t, x4)
		t.Mod(t, P)
		tSquare := new(big.Int).Mul(t, t)
		tSquare.Mod(tSquare, P)
		xSum := new(big.Int).Add(Fx, Gx)
		u := new(big.Int).Sub(tSquare, xSum)
		u.Mod(u, P)
		v := new(big.Int).Sub(Fx, u)
		v.Mul(v, t)
		v.Sub(v, Fy)
		v.Mod(v, P)
		Fx.Set(u)
		Fy.Set(v)
	}
	return &Point{X: Fx.Bytes(), Y:Fy.Bytes()}
}