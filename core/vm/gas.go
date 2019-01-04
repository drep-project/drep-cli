package vm

import (
	"math/big"
)

// Gas costs
const (
	GasQuickStep    uint64 = 2000 //2
	GasFastestStep  uint64 = 3000 //3
	GasFastStep     uint64 = 5000 //5
	GasMidStep      uint64 = 8000 //8
	GasSlowStep     uint64 = 1000 //10
	GasExtStep      uint64 = 2000 //20
	ExtcodeSize     uint64 = 7000 //700
    ExtcodeCopy     uint64 = 7000 //700
    ExtcodeHash     uint64 = 4000 //400
    Balance         uint64 = 4000 //400
    SLoad           uint64 = 2000 //200
    Calls           uint64 = 7000 //700
    Suicide         uint64 = 5000 //5000
    ExpByte         uint64 = 5000 //50
    CreateBySuicide uint64 = 2500 //25000
)

// calcGas returns the actual gas cost of the call.
//
// The cost of gas was changed during the homestead price change HF. To allow for EIP150
// to be implemented. The returned gas is gas - base * 63 / 64.
func callGas(availableGas, base uint64, callCost *big.Int) (uint64, error) {
	availableGas = availableGas - base
	gas := availableGas - availableGas/64
	// If the bit length exceeds 64 bit we know that the newly calculated "gas" for EIP150
	// is smaller than the requested amount. Therefor we return the new gas instead
	// of returning an error.
	if callCost.BitLen() > 64 || gas < callCost.Uint64() {
		return gas, nil
	}
	if callCost.BitLen() > 64 {
		return 0, errGasUintOverflow
	}
	return callCost.Uint64(), nil
}

func bigUint64(v *big.Int) (uint64, bool) {
	return v.Uint64(), v.BitLen() > 64
}