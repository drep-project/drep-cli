package vm

import (
	"github.com/drep-project/drepcli/database"
	"sync"
	"math/big"
	"errors"
	"github.com/drep-project/drepcli/accounts"
	"github.com/drep-project/drepcli/bean"
	"github.com/drep-project/drepcli/config"
)

var (
	state *State
	once sync.Once
	ErrNotAccountAddress = errors.New("a non account address occupied")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrAccountNotExists = errors.New("account not exists")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrCodeAlreadyExists = errors.New("code already exists")
	ErrCodeNotExists = errors.New("code not exists")
	ErrNotLogAddress = errors.New("a non log address occupied")
	ErrLogAlreadyExists = errors.New("log already exists")
)

type State struct {
	dt     database.Transactional
	refund uint64
}

func NewState(dt database.Transactional) *State {
	return &State{dt: dt}
}

func (s *State) CreateContractAccount(callerAddr accounts.CommonAddress, chainId config.ChainIdType, nonce int64) (*accounts.Account, error) {
	account, err := accounts.NewContractAccount(callerAddr, chainId, nonce)
	if err != nil {
		return nil, err
	}
	return account, database.PutStorage(s.dt, account.Address, chainId, account.Storage)
}

func (s *State) SubBalance(addr accounts.CommonAddress, chainId config.ChainIdType, amount *big.Int) error {
	balance := database.GetBalance(addr, chainId)
	return database.PutBalance(s.dt, addr, chainId, new(big.Int).Sub(balance, amount))
}

func (s *State) AddBalance(addr accounts.CommonAddress, chainId config.ChainIdType, amount *big.Int) error {
	balance := database.GetBalance(addr, chainId)
	return database.PutBalance(s.dt, addr, chainId, new(big.Int).Add(balance, amount))
}

func (s *State) GetBalance(addr accounts.CommonAddress, chainId config.ChainIdType) *big.Int {
	return database.GetBalance(addr, chainId)
}

func (s *State) SetNonce(addr accounts.CommonAddress, chainId config.ChainIdType, nonce int64) error {
	return database.PutNonce(s.dt, addr, chainId, nonce)
}

func (s *State) GetNonce(addr accounts.CommonAddress, chainId config.ChainIdType) int64 {
	return database.GetNonce(addr, chainId)
}

func (s *State) Suicide(addr accounts.CommonAddress, chainId config.ChainIdType) error {
	storage := database.GetStorage(addr, chainId)
	storage.Balance = new(big.Int)
	storage.Nonce = 0
	return database.PutStorage(s.dt, addr, chainId, storage)
}

func (s *State) GetByteCode(addr accounts.CommonAddress, chainId config.ChainIdType) accounts.ByteCode {
	return database.GetByteCode(addr, chainId)
}

func (s *State) GetCodeSize(addr accounts.CommonAddress, chainId config.ChainIdType) int {
	byteCode := s.GetByteCode(addr, chainId)
	return len(byteCode)
}

func (s *State) GetCodeHash(addr accounts.CommonAddress, chainId config.ChainIdType) accounts.Hash {
	return database.GetCodeHash(addr, chainId)
}

func (s *State) SetByteCode(addr accounts.CommonAddress, chainId config.ChainIdType, byteCode accounts.ByteCode) error {
	return database.PutByteCode(s.dt, addr, chainId, byteCode)
}

func (s *State) GetLogs(txHash []byte, chainId config.ChainIdType) []*bean.Log {
	return database.GetLogs(txHash, chainId)
}

func (s *State) AddLog(contractAddr accounts.CommonAddress, chainId config.ChainIdType, txHash, data []byte, topics [][]byte) error {
	log := &bean.Log{
		Address: contractAddr,
		ChainId: chainId,
		TxHash: txHash,
		Data: data,
		Topics: topics,
	}
	return database.AddLog(log)
}

func (s *State) AddRefund(gas uint64) {
	s.refund += gas
}

func (s *State) SubRefund(gas uint64) {
	if gas > s.refund {
		panic("refund below zero")
	}
	s.refund -= gas
}

func (s *State) Load(x *big.Int) []byte {
	value := s.dt.Get(x.Bytes())
	if value == nil {
		return new(big.Int).Bytes()
	}
	return value
}

func (s *State) Store(x, y *big.Int, chainId config.ChainIdType) {
	s.dt.Put(chainId, x.Bytes(), y.Bytes())
}