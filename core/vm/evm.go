package vm

import (
	"math/big"
	"errors"
	"github.com/drep-project/drepcli/accounts"
	"fmt"
	"github.com/drep-project/drepcli/database"
	"github.com/drep-project/drepcli/config"
	"github.com/ethereum/go-ethereum/params"
)

type EVM struct {
	State *State
	Interpreter *EVMInterpreter
	CallGasTemp uint64
	GasLimit uint64
	Origin accounts.CommonAddress
	GasPrice *big.Int
	CoinBase accounts.CommonAddress
	Time *big.Int
	Abort int32
}

var (
	ErrNoCompatibleInterpreter  = errors.New("no compatible interpreter")
)

func NewEVM(dt database.Transactional) *EVM {
	evm := &EVM{}
	subDt := dt.BeginTransaction()
	evm.State = NewState(subDt)
	evm.Interpreter = NewEVMInterpreter(evm)
	return evm
}

func (evm *EVM) CreateContractCode(callerAddr accounts.CommonAddress, chainId config.ChainIdType, byteCode accounts.ByteCode, gas uint64, value *big.Int) ([]byte, accounts.CommonAddress, uint64, error) {
	if !evm.CanTransfer(callerAddr, chainId, value) {
		return nil, accounts.CommonAddress{}, gas, ErrInsufficientBalance
	}

	nonce := evm.State.GetNonce(callerAddr, chainId)
	account, err := evm.State.CreateContractAccount(callerAddr, chainId, nonce)
	if err != nil {
		return nil, accounts.CommonAddress{}, gas, err
	}

	contractAddr := account.Address
	evm.Transfer(callerAddr, contractAddr, chainId, value)
	contract := NewContract(callerAddr, chainId, gas, value, nil)
	contract.SetCode(contractAddr, byteCode)
	ret, err := run(evm, contract, nil, false)

	createDataGas := uint64(len(ret)) * params.CreateDataGas
	if contract.UseGas(createDataGas) {
		err = evm.State.SetByteCode(contractAddr, chainId, ret)
	} else {
		err = ErrCodeStoreOutOfGas
	}

	if err != nil && err != ErrCodeStoreOutOfGas {
		evm.State.dt.Discard()
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	} else {
		evm.State.dt.Commit()
	}

	fmt.Println("contract address: ", contractAddr.Hex())
	fmt.Println("contract gas: ", contract.Gas)
	return ret, contractAddr, contract.Gas, err
}

func (evm *EVM) CallContractCode(callerAddr, contractAddr accounts.CommonAddress, chainId config.ChainIdType, input []byte, gas uint64, value *big.Int) (ret []byte, returnGas uint64, err error) {
	if !evm.CanTransfer(callerAddr, chainId, value) {
		return nil, gas, ErrInsufficientBalance
	}

	byteCode := evm.State.GetByteCode(contractAddr, chainId)
	if byteCode == nil {
		return nil, gas, ErrCodeNotExists
	}

	evm.Transfer(callerAddr, contractAddr, chainId, value)
	contract := NewContract(callerAddr, chainId, gas, value, nil)
	contract.SetCode(contractAddr, byteCode)

	ret, err = run(evm, contract, input, false)
	if err != nil {
		evm.State.dt.Discard()
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	} else {
		evm.State.dt.Commit()
	}

	return ret, contract.Gas, err
}

func (evm *EVM) StaticCall(callerAddr, contractAddr accounts.CommonAddress, chainId config.ChainIdType, input []byte, gas uint64) (ret []byte, returnGas uint64, err error) {
	byteCode := evm.State.GetByteCode(contractAddr, chainId)
	if byteCode == nil {
		return nil, gas, ErrCodeNotExists
	}

	contract := NewContract(callerAddr, chainId, gas, new(big.Int), nil)
	contract.SetCode(contractAddr, byteCode)

	ret, err = run(evm, contract, input, true)
	if err != nil {
		evm.State.dt.Discard()
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	} else {
		evm.State.dt.Commit()
	}

	return ret, contract.Gas, err
}

func (evm *EVM) DelegateCall(con *Contract, contractAddr accounts.CommonAddress, input []byte, gas uint64) (ret []byte, leftGas uint64, err error) {
	callerAddr := con.CallerAddr
	chainId := con.ChainId
	jumpdests := con.Jumpdests

	byteCode := evm.State.GetByteCode(contractAddr, chainId)
	if byteCode == nil {
		return nil, gas, ErrCodeNotExists
	}

	contract := NewContract(callerAddr, chainId, gas, new(big.Int), jumpdests)
	contract.SetCode(contractAddr, byteCode)

	ret, err = run(evm, contract, input, false)
	if err != nil {
		evm.State.dt.Discard()
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	} else {
		evm.State.dt.Commit()
	}

	return ret, con.Gas, err
}


func run(evm *EVM, contract *Contract, input []byte, readOnly bool) ([]byte, error) {
	if !contract.ContractAddr.IsEmpty() {
		precompiles := PrecompiledContracts
		if p := precompiles[contract.ContractAddr]; p != nil {
			return RunPrecompiledContract(p, input, contract)
		}
	}
	interpreter := evm.Interpreter
	if interpreter.canRun(contract.ByteCode) {
		return interpreter.Run(contract, input, readOnly)
	}
	return nil, ErrNoCompatibleInterpreter
}

func (evm *EVM) CanTransfer(addr accounts.CommonAddress, chainId config.ChainIdType, amount *big.Int) bool {
	balance := evm.State.GetBalance(addr, chainId)
	return balance.Cmp(amount) >= 0
}

func (evm *EVM) Transfer(from, to accounts.CommonAddress, chainId config.ChainIdType, amount *big.Int) error {
	err := evm.State.SubBalance(from, chainId, amount)
	if err != nil {
		return err
	}
	err = evm.State.AddBalance(to, chainId, amount)
	if err != nil {
		evm.State.AddBalance(from, chainId, amount)
		return err
	}
	return nil
}

