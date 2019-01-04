package core

import (
	"math/big"
	"github.com/drep-project/drepcli/core/vm"
	"testing"
	"fmt"
	"github.com/drep-project/drepcli/database"
	"github.com/drep-project/drepcli/core/abi"
	"strings"
	"encoding/hex"
	"github.com/drep-project/drepcli/accounts"
	"encoding/json"
	"github.com/drep-project/drepcli/bean"
	"github.com/drep-project/drepcli/config"
)

// Execute executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Executes sets up a in memory, temporarily, environment for the execution of
// the given code. It makes sure that it's restored to it's original state afterwards.
func ExecuteCreate(code []byte) {
	t := database.BeginTransaction()
	evm := vm.NewEVM(t)
	s1 := "111111"
	s2 := "222222"
	var chainId config.ChainIdType
	callerAddr1 := accounts.Hex2Address(s1)
	callerAddr2 := accounts.Hex2Address(s2)
	caller1 := &accounts.Account{Address: callerAddr1, Storage: &accounts.Storage{Balance: new(big.Int).SetInt64(100)}}
	caller2 := &accounts.Account{Address: callerAddr2, Storage: &accounts.Storage{Balance: new(big.Int).SetInt64(200)}}
	errPut1 := database.PutStorage(t, callerAddr1, chainId, caller1.Storage)
	errPut2 := database.PutStorage(t, callerAddr2, chainId, caller2.Storage)
	fmt.Println("errPut1: ", errPut1)
	fmt.Println("errPut2: ", errPut2)
	gas := uint64(1000000)
	value := new(big.Int).SetInt64(0)
	ret1, _, _, err1 := evm.CreateContractCode(callerAddr1, chainId, code, gas, value)
	ret2, _, _, err2 := evm.CreateContractCode(callerAddr2, chainId, code, gas, value)
	fmt.Println("err1: ", err1)
	fmt.Println("err2: ", err2)
	fmt.Println("ret1: ", ret1)
	fmt.Println("ret1: ", hex.EncodeToString(ret1))
	fmt.Println("ret2: ", ret2)
	fmt.Println("ret2: ", hex.EncodeToString(ret2))
}

func ExecuteCall(input []byte) {
	t := database.BeginTransaction()
	evm := vm.NewEVM(t)
	s1 := "111111"
	callerAddr := accounts.Hex2Address(s1)
	s2 := "a8eb43d6f487e6fbd2709512a5b8d90417dde6d8"
	gas := uint64(1000000)
	value := new(big.Int).SetInt64(0)
	contractAddr := accounts.Hex2Address(s2)
	var chainId config.ChainIdType
	evm.CallContractCode(callerAddr, contractAddr, chainId, input, gas, value)
}

func TestCreate(t *testing.T) {
	code, err := hex.DecodeString("608060405234801561001057600080fd5b5060b38061001f6000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806326121ff0146044575b600080fd5b604a604c565b005b6000808154809291906001019190505550674563918244f40000341115608057346001600082825401925050819055506085565b600080fd5b5600a165627a7a723058209b857a62cc21900a0effa2e33dffafe6529760e7524927e1b9935750a2371edc0029")
	if err != nil {
		fmt.Println("bad byte code")
	}
	ExecuteCreate(code)
}

func TestCallMyCode(t *testing.T) {
	var mystr = `[
	{
		"constant": false,
		"inputs": [],
		"name": "f",
		"outputs": [],
		"payable": true,
		"stateMutability": "payable",
		"type": "function"
	}
]`

	myabi, err := abi.JSON(strings.NewReader(mystr))
	if err != nil {
		fmt.Println("abi json error: ", err)
	}

	f, err := myabi.Pack("f")
	if err != nil {
		fmt.Println("abi pack error: ", err)
	} else {
		fmt.Println("abi: ", f)
	}
	fmt.Println("test input: ", hex.EncodeToString(f))
}

func TestThis(t *testing.T) {

}

func TestDescribeDatabase(t *testing.T) {
	itr := database.GetItr()
	fmt.Println("itr: ", itr)
	for itr.Next() {
		key := itr.Key()
		value := itr.Value()
		fmt.Println("key: ", key)
		fmt.Println()
		fmt.Println("value: ", value)
		fmt.Println()
		storage := &accounts.Storage{}
		err := json.Unmarshal(value, storage)
		if err == nil {
			fmt.Println()
			fmt.Println("v: ", storage)
			fmt.Println()
			continue
		}
		log := &bean.Log{}
		err = json.Unmarshal(value, log)
		if err == nil {
			fmt.Println("log: ", log)
		}
		fmt.Println()
	}
}