package component

import (
	"github.com/drep-project/drepcli/common"
	"github.com/drep-project/drepcli/crypto"
	"github.com/drep-project/drepcli/crypto/secp256k1"
	"github.com/drep-project/drepcli/crypto/sha3"
	accountTypes "github.com/drep-project/drepcli/accounts/types"
	"github.com/pkg/errors"
	"sync/atomic"
)

const (
	RPERMISSION = iota   //read
	WPERMISSION			//write
)

const (
	LOCKED = iota    	//locked
	UNLOCKED			//unlocked
)

type Wallet struct {
	cacheStore *accountCache

	chainId common.ChainIdType
	config *accountTypes.Config

	isLock int32
	password string
}

func NewWallet(config *accountTypes.Config, chainId  common.ChainIdType) (*Wallet, error) {
	wallet := &Wallet{
		config : config,
		chainId : chainId,
	}
	return wallet, nil
}

func (wallet *Wallet) Open(password string) error  {
	if wallet.cacheStore != nil {
		return errors.New("wallet is already open")
	}
	cryptedPassword := wallet.cryptoPassword(password)
	accountCacheStore, err  := NewAccountCache(wallet.config.KeyStoreDir, cryptedPassword)
	if err != nil {
		return err
	}
	wallet.cacheStore = accountCacheStore
	wallet.unLock(password)
	return nil
}

func (wallet *Wallet) Close()   {
	wallet.Lock()
	wallet.cacheStore = nil
	wallet.password = ""
}

func (wallet *Wallet) NewAccount() ( *accountTypes.Node, error){
	if err := wallet.checkWallet(WPERMISSION); err!= nil {
		return nil, err
	}

	newNode := accountTypes.NewNode(nil, wallet.chainId)
    wallet.cacheStore.StoreKey(newNode, wallet.password)
	return newNode, nil
}

func (wallet *Wallet) GetAccountByAddress(addr *crypto.CommonAddress) (*accountTypes.Node, error){
	if err := wallet.checkWallet(RPERMISSION); err!= nil {
		return nil, errors.New("wallet is not open")
	}
	return wallet.cacheStore.GetKey(addr, wallet.password)
}

func (wallet *Wallet) GetAccountByPubkey(pubkey *secp256k1.PublicKey) (*accountTypes.Node, error){
	if err := wallet.checkWallet(RPERMISSION); err!= nil {
		return nil, errors.New("wallet is not open")
	}
	addr := crypto.PubKey2Address(pubkey)
	return wallet.GetAccountByAddress(&addr)
}

func (wallet *Wallet) ListAddress() ([]*crypto.CommonAddress, error){
	if err := wallet.checkWallet(RPERMISSION); err!= nil {
		return nil, errors.New("wallet is not open")
	}
	nodes, err := wallet.cacheStore.ExportKey(wallet.password)
	if err != nil {
		return nil, err
	}
	addreses := []*crypto.CommonAddress{}
	for _, node := range nodes {
		addreses = append(addreses, node.Address)
	}
	return addreses, nil
}

func (wallet *Wallet) DumpPrivateKey(addr *crypto.CommonAddress) (*secp256k1.PrivateKey, error){
	if err := wallet.checkWallet(WPERMISSION); err!= nil {
		return nil, err
	}

	node, err :=  wallet.cacheStore.GetKey(addr, wallet.password)
	if err != nil {
		return  nil, err
	}
	return node.PrivateKey, nil
}

// 0 is locked  1 is unlock
func (wallet *Wallet) IsLock() bool {
	return atomic.LoadInt32(&wallet.isLock) == LOCKED
}

func (wallet *Wallet) IsOpen() bool {
	return wallet.cacheStore != nil
}
func (wallet *Wallet) Lock() error {
	atomic.StoreInt32(&wallet.isLock, LOCKED)
	wallet.cacheStore.ClearKeys()
	return nil
}

func (wallet *Wallet) UnLock(password string) error {
	if wallet.cacheStore == nil {
		wallet.Open(password)
	}else{
		wallet.unLock(password)
	}
	return nil
}

func (wallet *Wallet) unLock(password string) error {
	atomic.StoreInt32(&wallet.isLock, UNLOCKED)
	wallet.password = wallet.cryptoPassword(password)
	wallet.cacheStore.ReloadKeys(wallet.password)
	return nil
}

func (wallet *Wallet) checkWallet(op int) error {
	if wallet.cacheStore == nil {
		return errors.New("wallet is not open")
	}
	if 	op == WPERMISSION {
		if wallet.IsLock() {
			return errors.New("wallet is locked")
		}
	}
	return nil
}

func (wallet *Wallet) cryptoPassword(password string) string {
	return string(sha3.Hash256([]byte(password)))
}
