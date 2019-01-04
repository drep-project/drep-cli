package vm

import (
	"sync/atomic"
	"fmt"
	"github.com/drep-project/drepcli/accounts"
	"github.com/drep-project/drepcli/core/common"
)

type EVMInterpreter struct {
	IntPool *intPool
	EVM *EVM
	JumpTable [256]operation
	ReturnData []byte
	ReadOnly bool
}

func NewEVMInterpreter(evm *EVM) *EVMInterpreter {
	return &EVMInterpreter{EVM: evm, JumpTable: constantinopleInstructionSet}
}

func (in *EVMInterpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {
	if in.IntPool == nil {
		in.IntPool = poolOfIntPools.get()
		defer func() {
			poolOfIntPools.put(in.IntPool)
			in.IntPool = nil
		}()
	}

	if readOnly && !in.ReadOnly {
		in.ReadOnly = true
		defer func() { in.ReadOnly = false }()
	}

	// Reset the previous call's return data. It's unimportant to preserve the old buffer
	// as every returning call will return new data anyway.
	in.ReturnData = nil

	// Don't bother with the execution if there's no code.
	if len(contract.ByteCode) == 0 {
		return nil, nil
	}

	var (
		op    OpCode        // current opcode
		mem   = NewMemory() // bound memory
		stack = newstack()  // local stack
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		pc   = uint64(0) // program counter
		cost uint64
	)
	contract.Input = input

	// Reclaim the stack as an int pool when the execution stops
	defer func() { in.IntPool.put(stack.data...) }()


	// The Interpreter main run loop (contextual). This loop runs until either an
	// explicit STOP, RETURN or SELFDESTRUCT is executed, an error occurred during
	// the execution of one of the operations or until the done flag is set by the
	// parent context.

	c := contract.ByteCode
	co := make([]string, len(c))
	for ii := 0; ii < len(c); ii++ {
		co[ii] = OpCode(c[ii]).String()
	}
	fmt.Println()
	fmt.Println("c len: ", len(c))
	fmt.Println("c: ", c)
	fmt.Println("co: ", co)
	fmt.Println()

	opCount := 0
	for atomic.LoadInt32(&in.EVM.Abort) == 0 {
		// Get the operation from the jump table and validate the stack to ensure there are
		// enough stack items available to perform the operation.
		opCount += 1
		fmt.Println()
		fmt.Println("opCount: ", opCount)

		op = contract.GetOp(pc)
		operation := in.JumpTable[op]
		if !operation.valid {
			return nil, fmt.Errorf("invalid opcode 0x%x", int(op))
		}
		if err := operation.validateStack(stack); err != nil {
			return nil, err
		}
		if err := in.enforceRestrictions(op, operation, stack); err != nil {
			return nil, err
		}

		var memorySize uint64
		// calculate the new memory size and expand the memory to fit
		// the operation
		if operation.memorySize != nil {
			memSize, overflow := bigUint64(operation.memorySize(stack))
			if overflow {
				return nil, errGasUintOverflow
			}
			// memory is expanded in words of 32 bytes. Gas
			// is also calculated in words.
			if memorySize, overflow = common.SafeMul(common.ToWordSize(memSize), 32); overflow {
				return nil, errGasUintOverflow
			}
		}
		// consume the gas and return an error if not enough gas is available.
		// cost is explicitly set so that the capture state defer method can get the proper cost
		cost, err = operation.gasCost(in.EVM, contract, stack, mem, memorySize)
		fmt.Println("operation: ", operation)
		fmt.Println("pc: ", pc)
		fmt.Println("op: ", op.String())
		fmt.Print("before")
		stack.Print()
		fmt.Println("cost: ", cost)
		fmt.Println("err: ", err)
		if err != nil || !contract.UseGas(cost) {
			return nil, ErrOutOfGas
		}
		if memorySize > 0 {
			mem.Resize(memorySize)
		}

		// execute the operation
		res, err := operation.execute(&pc, in, contract, mem, stack)
		fmt.Print("after")
		stack.Print()
		fmt.Print()

		// if the operation clears the return data (e.g. it has returning data)
		// set the last return to the result of the operation.
		if operation.returns {
			in.ReturnData = res
		}

		switch {
		case err != nil:
			return nil, err
		case operation.reverts:
			return res, errExecutionReverted
		case operation.halts:
			return res, nil
		case !operation.jumps:
			pc++
		}
	}
	return nil, nil
}

func (in *EVMInterpreter) canRun(byteCode accounts.ByteCode) bool {
	return true
}

func (in *EVMInterpreter) enforceRestrictions(op OpCode, operation operation, stack *Stack) error {
	if in.ReadOnly {
		// If the interpreter is operating in readonly mode, make sure no
		// state-modifying operation is performed. The 3rd stack item
		// for a call operation is the value. Transferring value from one
		// account to the others means the state is modified and should also
		// return with an error.
		if operation.writes || (op == CALL && stack.Back(2).BitLen() > 0) {
			return errWriteProtection
		}
	}
	return nil
}