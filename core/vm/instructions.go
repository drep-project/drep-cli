// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"errors"
	"math/big"
	"github.com/drep-project/drepcli/mycrypto"
	"fmt"
    "github.com/drep-project/drepcli/accounts"
	"github.com/drep-project/drepcli/core/common"
)

var (
	bigZero                  = new(big.Int)
	errWriteProtection       = errors.New("evm: write protection")
	errReturnDataOutOfBounds = errors.New("evm: return data out of bounds")
	errExecutionReverted     = errors.New("evm: execution reverted")
	//errMaxCodeSizeExceeded   = errors.New("evm: max code size exceeded")
	ErrCodeStoreOutOfGas        = errors.New("contract creation code storage out of gas")
)

func opAdd(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	common.U256(y.Add(x, y))

	interpreter.IntPool.put(x)
	return nil, nil
}

func opSub(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	common.U256(y.Sub(x, y))

	interpreter.IntPool.put(x)
	return nil, nil
}

func opMul(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	stack.push(common.U256(x.Mul(x, y)))

	interpreter.IntPool.put(y)

	return nil, nil
}

func opDiv(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if y.Sign() != 0 {
		common.U256(y.Div(x, y))
	} else {
		y.SetUint64(0)
	}
	interpreter.IntPool.put(x)
	return nil, nil
}

func opSdiv(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := common.S256(stack.pop()), common.S256(stack.pop())
	res := interpreter.IntPool.getZero()

	if y.Sign() == 0 || x.Sign() == 0 {
		stack.push(res)
	} else {
		if x.Sign() != y.Sign() {
			res.Div(x.Abs(x), y.Abs(y))
			res.Neg(res)
		} else {
			res.Div(x.Abs(x), y.Abs(y))
		}
		stack.push(common.U256(res))
	}
	interpreter.IntPool.put(x, y)
	return nil, nil
}

func opMod(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	if y.Sign() == 0 {
		stack.push(x.SetUint64(0))
	} else {
		stack.push(common.U256(x.Mod(x, y)))
	}
	interpreter.IntPool.put(y)
	return nil, nil
}

func opSmod(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := common.S256(stack.pop()), common.S256(stack.pop())
	res := interpreter.IntPool.getZero()

	if y.Sign() == 0 {
		stack.push(res)
	} else {
		if x.Sign() < 0 {
			res.Mod(x.Abs(x), y.Abs(y))
			res.Neg(res)
		} else {
			res.Mod(x.Abs(x), y.Abs(y))
		}
		stack.push(common.U256(res))
	}
	interpreter.IntPool.put(x, y)
	return nil, nil
}

func opExp(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	base, exponent := stack.pop(), stack.pop()
	stack.push(common.Exp(base, exponent))

	interpreter.IntPool.put(base, exponent)

	return nil, nil
}

func opSignExtend(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	back := stack.pop()
	if back.Cmp(big.NewInt(31)) < 0 {
		bit := uint(back.Uint64()*8 + 7)
		num := stack.pop()
		mask := back.Lsh(common.Big1, bit)
		mask.Sub(mask, common.Big1)
		if num.Bit(int(bit)) > 0 {
			num.Or(num, mask.Not(mask))
		} else {
			num.And(num, mask)
		}

		stack.push(common.U256(num))
	}

	interpreter.IntPool.put(back)
	return nil, nil
}

func opNot(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x := stack.peek()
	common.U256(x.Not(x))
	return nil, nil
}

func opLt(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if x.Cmp(y) < 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	interpreter.IntPool.put(x)
	return nil, nil
}

func opGt(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if x.Cmp(y) > 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	interpreter.IntPool.put(x)
	return nil, nil
}

func opSlt(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()

	xSign := x.Cmp(common.TT255)
	ySign := y.Cmp(common.TT255)

	switch {
	case xSign >= 0 && ySign < 0:
		y.SetUint64(1)

	case xSign < 0 && ySign >= 0:
		y.SetUint64(0)

	default:
		if x.Cmp(y) < 0 {
			y.SetUint64(1)
		} else {
			y.SetUint64(0)
		}
	}
	interpreter.IntPool.put(x)
	return nil, nil
}

func opSgt(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()

	xSign := x.Cmp(common.TT255)
	ySign := y.Cmp(common.TT255)

	switch {
	case xSign >= 0 && ySign < 0:
		y.SetUint64(0)

	case xSign < 0 && ySign >= 0:
		y.SetUint64(1)

	default:
		if x.Cmp(y) > 0 {
			y.SetUint64(1)
		} else {
			y.SetUint64(0)
		}
	}
	interpreter.IntPool.put(x)
	return nil, nil
}

func opEq(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	if x.Cmp(y) == 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	interpreter.IntPool.put(x)
	return nil, nil
}

func opIszero(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x := stack.peek()
	if x.Sign() > 0 {
		x.SetUint64(0)
	} else {
		x.SetUint64(1)
	}
	return nil, nil
}

func opAnd(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	stack.push(x.And(x, y))

	interpreter.IntPool.put(y)
	return nil, nil
}

func opOr(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Or(x, y)

	interpreter.IntPool.put(x)
	return nil, nil
}

func opXor(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.peek()
	y.Xor(x, y)

	interpreter.IntPool.put(x)
	return nil, nil
}

func opByte(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	th, val := stack.pop(), stack.peek()
	if th.Cmp(common.Big32) < 0 {
		b := common.Byte(val, 32, int(th.Int64()))
		val.SetUint64(uint64(b))
	} else {
		val.SetUint64(0)
	}
	interpreter.IntPool.put(th)
	return nil, nil
}

func opAddmod(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	if z.Cmp(bigZero) > 0 {
		x.Add(x, y)
		x.Mod(x, z)
		stack.push(common.U256(x))
	} else {
		stack.push(x.SetUint64(0))
	}
	interpreter.IntPool.put(y, z)
	return nil, nil
}

func opMulmod(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	if z.Cmp(bigZero) > 0 {
		x.Mul(x, y)
		x.Mod(x, z)
		stack.push(common.U256(x))
	} else {
		stack.push(x.SetUint64(0))
	}
	interpreter.IntPool.put(y, z)
	return nil, nil
}

// opSHL implements Shift Left
// The SHL instruction (shift left) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the left by arg1 number of bits.
func opSHL(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := common.U256(stack.pop()), common.U256(stack.peek())
	defer interpreter.IntPool.put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		value.SetUint64(0)
		return nil, nil
	}
	n := uint(shift.Uint64())
	common.U256(value.Lsh(value, n))

	return nil, nil
}

// opSHR implements Logical Shift Right
// The SHR instruction (logical shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with zero fill.
func opSHR(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := common.U256(stack.pop()), common.U256(stack.peek())
	defer interpreter.IntPool.put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		value.SetUint64(0)
		return nil, nil
	}
	n := uint(shift.Uint64())
	common.U256(value.Rsh(value, n))

	return nil, nil
}

// opSAR implements Arithmetic Shift Right
// The SAR instruction (arithmetic shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with sign extension.
func opSAR(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Note, S256 returns (potentially) a new bigint, so we're popping, not peeking this one
	shift, value := common.U256(stack.pop()), common.S256(stack.pop())
	defer interpreter.IntPool.put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		if value.Sign() >= 0 {
			value.SetUint64(0)
		} else {
			value.SetInt64(-1)
		}
		stack.push(common.U256(value))
		return nil, nil
	}
	n := uint(shift.Uint64())
	value.Rsh(value, n)
	stack.push(common.U256(value))

	return nil, nil
}

func opSha3(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	offset, size := stack.pop(), stack.pop()
	data := memory.Get(offset.Int64(), size.Int64())
	//hash := crypto.Keccak256(data)
	hash := mycrypto.Hash256(data)
	//evm := interpreter.evm
	//if evm.vmConfig.EnablePreimageRecording {
	//	evm.StateDB.AddPreimage(BytesToHash(hash), data)
	//}
	fmt.Println("offset: ", offset)
	fmt.Println("size: ", size)
	fmt.Println("sha3 data: ", data)
	stack.push(interpreter.IntPool.get().SetBytes(hash))

	interpreter.IntPool.put(offset, size)
	return nil, nil
}

func opAddress(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x := contract.ContractAddr.Big()
	//stack.push(contract.GetAddress().Big())
	stack.push(x)
	return nil, nil
}

func opBalance(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	slot := stack.peek()
	//slot.Set(interpreter.evm.StateDB.GetBalance(BigToAddress(slot)))
	evm := interpreter.EVM
	balance := evm.State.GetBalance(accounts.Big2Address(slot), contract.ChainId)
	slot.Set(balance)
	return nil, nil
}

func opOrigin(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//stack.push(interpreter.evm.Origin.Big())
	x := interpreter.EVM.Origin.Big()
	stack.push(x)
	return nil, nil
}

func opCaller(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//stack.push(contract.Caller().Big())
	x := contract.CallerAddr.Big()
	stack.push(x)
	return nil, nil
}

func opCallValue(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(interpreter.IntPool.get().Set(contract.Value))
	return nil, nil
}

func opCallDataLoad(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fmt.Println("call data load contract input: ", contract.Input)
	stack.push(interpreter.IntPool.get().SetBytes(common.GetDataBig(contract.Input, stack.pop(), big32)))
	return nil, nil
}

func opCallDataSize(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(interpreter.IntPool.get().SetInt64(int64(len(contract.Input))))
	return nil, nil
}

func opCallDataCopy(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()
	)
	memory.Set(memOffset.Uint64(), length.Uint64(), common.GetDataBig(contract.Input, dataOffset, length))

	interpreter.IntPool.put(memOffset, dataOffset, length)
	return nil, nil
}

func opReturnDataSize(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(interpreter.IntPool.get().SetUint64(uint64(len(interpreter.ReturnData))))
	return nil, nil
}

func opReturnDataCopy(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()

		end = interpreter.IntPool.get().Add(dataOffset, length)
	)
	defer interpreter.IntPool.put(memOffset, dataOffset, length, end)

	if end.BitLen() > 64 || uint64(len(interpreter.ReturnData)) < end.Uint64() {
		return nil, errReturnDataOutOfBounds
	}
	memory.Set(memOffset.Uint64(), length.Uint64(), interpreter.ReturnData[dataOffset.Uint64():end.Uint64()])

	return nil, nil
}

func opExtCodeSize(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	slot := stack.peek()
	//slot.SetUint64(uint64(interpreter.EVM.StateDB.GetCodeSize(BigToAddress(slot))))
	l := interpreter.EVM.State.GetCodeSize(accounts.Big2Address(slot), contract.ChainId)
	slot.SetUint64(uint64(l))
	return nil, nil
}

func opCodeSize(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	l := interpreter.IntPool.get().SetInt64(int64(len(contract.ByteCode)))
	stack.push(l)

	return nil, nil
}

func opCodeCopy(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	codeCopy := common.GetDataBig(contract.ByteCode, codeOffset, length)
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	interpreter.IntPool.put(memOffset, codeOffset, length)
	return nil, nil
}

func opExtCodeCopy(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		addr       = accounts.Big2Address(stack.pop())
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	byteCode := interpreter.EVM.State.GetByteCode(addr, contract.ChainId)
	codeCopy := common.GetDataBig(byteCode, codeOffset, length)
	//codeCopy := getDataBig(interpreter.evm.StateDB.GetCode(addr), codeOffset, length)
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	interpreter.IntPool.put(memOffset, codeOffset, length)
	return nil, nil
}

// opExtCodeHash returns the code hash of a specified account.
// There are several cases when the function is called, while we can relay everything
// to `state.GetCodeHash` function to ensure the correctness.
//   (1) Caller tries to get the code hash of a normal contract account, state
// should return the relative code hash and set it as the result.
//
//   (2) Caller tries to get the code hash of a non-existent account, state should
// return Hash{} and zero will be set as the result.
//
//   (3) Caller tries to get the code hash for an account without contract code,
// state should return emptyCodeHash(0xc5d246...) as the result.
//
//   (4) Caller tries to get the code hash of a precompiled account, the result
// should be zero or emptyCodeHash.
//
// It is worth noting that in order to avoid unnecessary create and clean,
// all precompile accounts on mainnet have been transferred 1 wei, so the return
// here should be emptyCodeHash.
// If the precompile account is not transferred any amount on a private or
// customized chain, the return value will be zero.
//
//   (5) Caller tries to get the code hash for an account which is marked as suicided
// in the current dt, the code hash of this account should be returned.
//
//   (6) Caller tries to get the code hash for an account which is marked as deleted,
// this account should be regarded as a non-existent account and zero should be returned.
func opExtCodeHash(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	slot := stack.peek()
	hash := interpreter.EVM.State.GetCodeHash(accounts.Big2Address(slot), contract.ChainId).Bytes()
	slot.SetBytes(hash)
	//slot.SetBytes(interpreter.evm.StateDB.GetCodeHash(BigToAddress(slot)).Bytes())
	return nil, nil
}

func opGasprice(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(interpreter.IntPool.get().Set(interpreter.EVM.GasPrice))
	return nil, nil
}

func opBlockhash(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fmt.Println("doing blockhash")
	//num := stack.pop()
	//n := interpreter.IntPool.get().Sub(interpreter.evm.BlockNumber, Big257)
	//if num.Cmp(n) > 0 && num.Cmp(interpreter.evm.BlockNumber) < 0 {
	//	stack.push(interpreter.evm.GetHash(num.Uint64()).Big())
	//} else {
	//	stack.push(interpreter.intPool.getZero())
	//}
	//interpreter.IntPool.put(num, n)
	return nil, nil
}

func opCoinbase(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x := interpreter.EVM.CoinBase.Big()
	//stack.push(interpreter.evm.Coinbase.Big())
	stack.push(x)
	return nil, nil
}

func opTimestamp(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(common.U256(interpreter.IntPool.get().Set(interpreter.EVM.Time)))
	return nil, nil
}

func opNumber(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//stack.push(U256(interpreter.IntPool.get().Set(interpreter.evm.BlockNumber)))
	return nil, nil
}

func opDifficulty(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//stack.push(U256(interpreter.IntPool.get().Set(interpreter.evm.Difficulty)))
	return nil, nil
}

func opGasLimit(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(common.U256(interpreter.IntPool.get().SetUint64(interpreter.EVM.GasLimit)))
	return nil, nil
}

func opPop(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	interpreter.IntPool.put(stack.pop())
	return nil, nil
}

func opMload(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	offset := stack.pop()
	val := interpreter.IntPool.get().SetBytes(memory.Get(offset.Int64(), 32))
	stack.push(val)
	interpreter.IntPool.put(offset)
	return nil, nil
}

func opMstore(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// pop value of the stack
	mStart, val := stack.pop(), stack.pop()
	memory.Set32(mStart.Uint64(), val)
	interpreter.IntPool.put(mStart, val)
	return nil, nil
}

func opMstore8(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	off, val := stack.pop().Int64(), stack.pop().Int64()
	memory.store[off] = byte(val & 0xff)

	return nil, nil
}

func opSload(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fmt.Println("is sloading")
	loc := stack.peek()
	fmt.Println("loc: ", loc)
	modifiedLoc := new(big.Int).SetBytes(mycrypto.Hash256(contract.ByteCode, loc.Bytes()))
	//addr, err := interpreter.EVM.State.GetAccountAddress(loc)
	//if err != nil {
	//	loc.Set(new(big.Int))
	//	return nil, err
	//}
	////val := interpreter.evm.StateDB.GetState(contract.GetAddress(), BigToHash(loc))
	////loc.SetBytes(val.Bytes())
	//b, err := AddressToBig(addr)
	//if err != nil {
	//	loc.Set(new(big.Int))
	//	return nil, err
	//}
	//b := interpreter.EVM.State.Load(loc)
	b := interpreter.EVM.State.Load(modifiedLoc)
	loc.SetBytes(b)
	return nil, nil
}

func opSstore(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//loc := BigToHash(stack.pop())
	loc, val := stack.pop(), stack.pop()
	//fmt.Println("loc: ", loc)
	modifiedLoc := new(big.Int).SetBytes(mycrypto.Hash256(contract.ByteCode, loc.Bytes()))
	interpreter.EVM.State.Store(modifiedLoc, val, contract.ChainId)
	//interpreter.EVM.State.Store(loc, val)
	interpreter.IntPool.put(val)
	return nil, nil
}

func opJump(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	pos := stack.pop()
	if !contract.Jumpdests.has(contract.CodeHash, contract.ByteCode, pos) {
		//fmt.Println()
		//fmt.Println("jump pos: ", pos)
		//fmt.Println("jump code: ", contract.GetOp(pos.Uint64()))
		//fmt.Println()
		nop := contract.GetOp(pos.Uint64())
		return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
	}
	*pc = pos.Uint64()

	interpreter.IntPool.put(pos)
	return nil, nil
}

func opJumpi(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//fmt.Println("stack: ", stack)
	pos, cond := stack.pop(), stack.pop()
	//fmt.Println("stack: ", stack)
	//fmt.Println("jumpi pos: ", pos)
	//fmt.Println("jumpi cond: ", cond)
	if cond.Sign() != 0 {
		if !contract.Jumpdests.has(contract.CodeHash, contract.ByteCode, pos) {
			nop := contract.GetOp(pos.Uint64())
			return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
		}
		*pc = pos.Uint64()
	} else {
		*pc++
	}

	interpreter.IntPool.put(pos, cond)
	return nil, nil
}

func opJumpdest(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	return nil, nil
}

func opPc(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(interpreter.IntPool.get().SetUint64(*pc))
	return nil, nil
}

func opMsize(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(interpreter.IntPool.get().SetInt64(int64(memory.Len())))
	return nil, nil
}

func opGas(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(interpreter.IntPool.get().SetUint64(contract.Gas))
	return nil, nil
}

func opCreate(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		value        = stack.pop()
		offset, size = stack.pop(), stack.pop()
		//input        = memory.Get(offset.Int64(), size.Int64())
		gas          = contract.Gas
	)
	//if interpreter.evm.ChainConfig().IsEIP150(interpreter.evm.BlockNumber) {
	//	gas -= gas / 64
	//}

	contract.UseGas(gas)
	//res, addr, returnGas, suberr := interpreter.evm.Create(contract, input, gas, value)
	res, addr, returnGas, suberr := interpreter.EVM.CreateContractCode(contract.CallerAddr, contract.ChainId, contract.ByteCode, gas, value)
	// Push item on the stack based on the returned error. If the ruleset is
	// homestead we must check for CodeStoreOutOfGasError (homestead only
	// rule) and treat as an error, if the ruleset is frontier we must
	// ignore this error and pretend the operation was successful.
	if suberr != nil && suberr != ErrCodeStoreOutOfGas {
		stack.push(interpreter.IntPool.getZero())
	} else {
		//stack.push(addr.Big())
		x:= addr.Big()
		stack.push(x)
	}
	contract.Gas += returnGas
	interpreter.IntPool.put(value, offset, size)

	if suberr == errExecutionReverted {
		interpreter.EVM.State.dt.Discard()
		return res, nil
	}
	return nil, nil
}

func opCreate2(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//var (
	//	endowment    = stack.pop()
	//	offset, size = stack.pop(), stack.pop()
	//	salt         = stack.pop()
	//	input        = memory.Get(offset.Int64(), size.Int64())
	//	gas          = contract.Gas
	//)
	//
	//// Apply EIP150
	//gas -= gas / 64
	//contract.UseGas(gas)
	//res, addr, returnGas, suberr := interpreter.evm.Create2(contract, input, gas, endowment, salt)
	//// Push item on the stack based on the returned error.
	//if suberr != nil {
	//	stack.push(interpreter.IntPool.getZero())
	//} else {
	//	stack.push(addr.Big())
	//}
	//contract.Gas += returnGas
	//interpreter.IntPool.put(endowment, offset, size, salt)
	//
	//if suberr == errExecutionReverted {
	//	return res, nil
	//}
	//return nil, nil
	return opCreate(pc, interpreter, contract, memory, stack)
}

func opCall(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas in in interpreter.evm.callGasTemp.
	interpreter.IntPool.put(stack.pop())
	//gas := interpreter.evm.callGasTemp
	gas := interpreter.EVM.CallGasTemp
	// Pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	//toAddr := BigToAddress(addr)
	value = common.U256(value)
	// Get the arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		gas += CallStipend
	}
	ret, returnGas, err := interpreter.EVM.CallContractCode(contract.CallerAddr, contract.ContractAddr, contract.ChainId, args, gas, value)
	if err != nil {
		stack.push(interpreter.IntPool.getZero())
	} else {
		stack.push(interpreter.IntPool.get().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
		interpreter.EVM.State.dt.Discard()
	}
	contract.Gas += returnGas

	interpreter.IntPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opCallCode(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	interpreter.IntPool.put(stack.pop())
	//gas := interpreter.evm.callGasTemp
	gas := interpreter.EVM.CallGasTemp
	// Pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := accounts.Big2Address(addr)
	value = common.U256(value)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		gas += CallStipend
	}
	ret, returnGas, err := interpreter.EVM.CallContractCode(contract.CallerAddr, toAddr, contract.ChainId, args, gas, value)
	if err != nil {
		stack.push(interpreter.IntPool.getZero())
	} else {
		stack.push(interpreter.IntPool.get().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
		interpreter.EVM.State.dt.Discard()
	}
	contract.Gas += returnGas

	interpreter.IntPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opDelegateCall(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	interpreter.IntPool.put(stack.pop())
	gas := interpreter.EVM.CallGasTemp
	// Pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := accounts.Big2Address(addr)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := interpreter.EVM.DelegateCall(contract, toAddr, args, gas)
	if err != nil {
		stack.push(interpreter.IntPool.getZero())
	} else {
		stack.push(interpreter.IntPool.get().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
		interpreter.EVM.State.dt.Discard()
	}
	contract.Gas += returnGas

	interpreter.IntPool.put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opStaticCall(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.callGasTemp.
	interpreter.IntPool.put(stack.pop())
	//gas := interpreter.evm.callGasTemp
	gas := interpreter.EVM.CallGasTemp
	// Pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := accounts.Big2Address(addr)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := interpreter.EVM.StaticCall(contract.CallerAddr, toAddr, contract.ChainId, args, gas)
	if err != nil {
		stack.push(interpreter.IntPool.getZero())
	} else {
		stack.push(interpreter.IntPool.get().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
		interpreter.EVM.State.dt.Discard()
	}
	contract.Gas += returnGas

	interpreter.IntPool.put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opReturn(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())
	interpreter.IntPool.put(offset, size)
	fmt.Println("ret: ", ret)
	fmt.Println("ret: ", len(ret))
	fmt.Println("ret: ", len(contract.ByteCode))
	return ret, nil
}

func opRevert(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//fmt.Println("is doing revert")
	offset, size := stack.pop(), stack.pop()
	//fmt.Println(offset, size)
	ret := memory.GetPtr(offset.Int64(), size.Int64())
	//fmt.Println("ret: ", ret)

	interpreter.IntPool.put(offset, size)
	return ret, nil
}

func opStop(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	return nil, nil
}

func opSuicide(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	//balance := interpreter.evm.StateDB.GetBalance(contract.GetAddress())
	//interpreter.evm.StateDB.AddBalance(BigToAddress(stack.pop()), balance)
	//interpreter.evm.StateDB.Suicide(contract.GetAddress())

	balance := interpreter.EVM.State.GetBalance(contract.CallerAddr, contract.ChainId)
	interpreter.EVM.State.AddBalance(accounts.Big2Address(stack.pop()), contract.ChainId, balance)
	interpreter.EVM.State.Suicide(contract.CallerAddr, contract.ChainId)
	return nil, nil
}

// following functions are used by the instruction jump  table

// make log instruction function
func makeLog(size int) executionFunc {
	return func(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		//topics := make([]Hash, size)
		topics := make([][]byte, size)
		mStart, mSize := stack.pop(), stack.pop()
		for i := 0; i < size; i++ {
			topics[i] = accounts.Big2Hash(stack.pop()).Bytes()
		}
		d := memory.Get(mStart.Int64(), mSize.Int64())
		//interpreter.evm.StateDB.AddLog(&types.Log{
		//	GetAddress: contract.GetAddress(),
		//	Topics:  topics,
		//	Data:    d,
		//	// This is a non-consensus field, but assigned here because
		//	// core/state doesn't know the current block number.
		//	BlockNumber: interpreter.evm.BlockNumber.Uint64(),
		//})
		//
		//interpreter.IntPool.put(mStart, mSize)
		interpreter.EVM.State.AddLog(contract.CallerAddr, contract.ChainId, contract.TxHash, d, topics)
		return nil, nil
	}
}

// make push instruction function
func makePush(size uint64, pushByteSize int) executionFunc {
	return func(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		codeLen := len(contract.ByteCode)

		//fmt.Println("codeLen: ", codeLen)
		//fmt.Println("*pc + 1: ", *pc+1)

		startMin := codeLen
		if int(*pc+1) < startMin {
			startMin = int(*pc + 1)
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		integer := interpreter.IntPool.get()
		//fmt.Println("integer get: ", integer)
		//fmt.Println("contract code startmin: endmin: ", contract.ByteCode[startMin: endMin])
		stack.push(integer.SetBytes(common.RightPadBytes(contract.ByteCode[startMin:endMin], pushByteSize)))
		//fmt.Println("stack here: ")
		//stack.Print()
		//fmt.Println("intpool here: ")
		//interpreter.IntPool.pool.Print()

		*pc += size
		//fmt.Println("*pc: ", *pc)
		return nil, nil
	}
}

// make dup instruction function
func makeDup(size int64) executionFunc {
	return func(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		stack.dup(interpreter.IntPool, int(size))
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	// switch n + 1 otherwise n would be swapped with n
	size++
	return func(pc *uint64, interpreter *EVMInterpreter, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		stack.swap(int(size))
		return nil, nil
	}
}
