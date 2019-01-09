package common

import (
	"fmt"
	"math"
	"strings"
)

type BlockNumber int64

const (
	PendingBlockNumber  = BlockNumber(-2)
	LatestBlockNumber   = BlockNumber(-1)
	EarliestBlockNumber = BlockNumber(0)
)

// UnmarshalJSON parses the given JSON fragment into a BlockNumber. It supports:
// - "latest", "earliest" or "pending" as string arguments
// - the block number
// Returned errors:
// - an invalid block number error when the given argument isn't a known strings
// - an out of range error when the given block number is either too little or too large
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case "earliest":
		*bn = EarliestBlockNumber
		return nil
	case "latest":
		*bn = LatestBlockNumber
		return nil
	case "pending":
		*bn = PendingBlockNumber
		return nil
	}

	blckNum, err := DecodeUint64(input)
	if err != nil {
		return err
	}
	if blckNum > math.MaxInt64 {
		return fmt.Errorf("Blocknumber too high")
	}

	*bn = BlockNumber(blckNum)
	return nil
}

func (bn BlockNumber) Int64() int64 {
	return (int64)(bn)
}
