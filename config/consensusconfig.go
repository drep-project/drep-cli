package config

import (
	"github.com/drep-project/drepcli/mycrypto"
)

type ConsensusConfig struct {
	ConsensusMode string			`json:"consensusMode"`
	Producers []*Produce		`json:"producers"`
	MyPk 	 *mycrypto.Point		`json:"mypk"`
}

//TODO how to identify a mine pk or pr&addr
type Produce struct {
	Public  *mycrypto.Point
	Ip string
	Port int
}