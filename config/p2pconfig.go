package config

import (
	"github.com/drep-project/drepcli/mycrypto"
)
type P2pConfig struct {
	PrvKey *mycrypto.PrivateKey `json:"omitempty"`
	ListerAddr string	`json:"omitempty"`
	Port int
	BootNodes []BootNode  //pub@Addr
}
