package mycrypto

type Point struct {
	X                    []byte
	Y                    []byte
}

type PrivateKey struct {
	Prv                  []byte
	PubKey               *Point
}

type Signature struct {
	R                    []byte
	S                    []byte
}