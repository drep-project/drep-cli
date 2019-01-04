package vm

import (
	"math/big"
	"github.com/drep-project/drepcli/accounts"
	"github.com/drep-project/drepcli/config"
)

var (
	big1      = big.NewInt(1)
	big4      = big.NewInt(4)
	big8      = big.NewInt(8)
	big16     = big.NewInt(16)
	big32     = big.NewInt(32)
	big64     = big.NewInt(64)
	big96     = big.NewInt(96)
	big480    = big.NewInt(480)
	big1024   = big.NewInt(1024)
	big3072   = big.NewInt(3072)
	big199680 = big.NewInt(199680)
)

type Contract struct {
	CallerAddr   accounts.CommonAddress
	ContractAddr accounts.CommonAddress
	ChainId      config.ChainIdType
	ByteCode     accounts.ByteCode
	CodeHash     accounts.Hash
	Input        []byte
	Gas          uint64
	Value        *big.Int
	Jumpdests    destinations
	TxHash       []byte
}

func NewContract(callerAddr accounts.CommonAddress, chainId config.ChainIdType, gas uint64, value *big.Int, jumpdests destinations) *Contract {
	if jumpdests == nil {
		return &Contract{CallerAddr: callerAddr, ChainId: chainId, Gas: gas, Value: value, Jumpdests: NewDest()}
	}
	return &Contract{CallerAddr: callerAddr, Gas: gas, Value: value, Jumpdests: jumpdests}
}

func (c *Contract) SetCode(contractAddr accounts.CommonAddress, byteCode accounts.ByteCode) {
	c.ContractAddr = contractAddr
	c.ByteCode = byteCode
	c.CodeHash = accounts.GetByteCodeHash(byteCode)
}

func (c *Contract) GetOp(n uint64) OpCode {
	return OpCode(c.GetByte(n))
}

func (c *Contract) GetByte(n uint64) byte {
	if n < uint64(len(c.ByteCode)) {
		return c.ByteCode[n]
	}
	return 0
}

func (c *Contract) UseGas(gas uint64) bool {
	if c.Gas < gas {
		return false
	}
	c.Gas -= gas
	return true
}