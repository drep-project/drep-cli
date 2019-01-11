package types

import "reflect"

type RestDescription struct {
	Api reflect.Value
}

type Request struct {
	Method string `json:"method"`
	Params string `json:"params"`
}

type Response struct {
	Success  bool        `json:"success"`
	ErrorMsg string      `json:"errMsg"`
	Data     interface{} `json:"body"`
}
