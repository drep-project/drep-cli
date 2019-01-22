package component

import (
	"encoding/json"
	"fmt"
	accountTypes "github.com/drep-project/drepcli/accounts/types"
	"github.com/drep-project/drepcli/common"
	"github.com/drep-project/drepcli/crypto"
	"github.com/drep-project/drepcli/crypto/aes"
	"github.com/drep-project/drepcli/crypto/secp256k1"
	"github.com/drep-project/drepcli/log"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type keyStore interface {
	// Loads and decrypts the key from disk.
	GetKey(addr *crypto.CommonAddress, auth string) (*accountTypes.Node, error)
	// Writes and encrypts the key.
	StoreKey(k *accountTypes.Node, auth string) error
	// Writes and encrypts the key.
	ExportKey(auth string) ([]*accountTypes.Node, error)
	// Joins filename with the key directory unless it is already absolute.
	JoinPath(filename string) string
}

type CryptedNode struct {
	CryptoPrivateKey []byte                `json:"cryptoPrivateKey"`
	PrivateKey       *secp256k1.PrivateKey `json:"-"`
	ChainId          common.ChainIdType    `json:"chainId"`
	ChainCode        []byte                `json:"chainCode"`

	Key []byte `json:"-"`
	Iv  []byte `json:"iv"`
}

func (cryptedNode *CryptedNode) EnCrypt() {
	cryptedNode.CryptoPrivateKey = aes.AesCBCEncrypt(cryptedNode.PrivateKey.Serialize(), cryptedNode.Key, cryptedNode.Iv)
}

func (cryptedNode *CryptedNode) DeCrypt() *accountTypes.Node {
	privKeyBytes := aes.AesCBCDecrypt(cryptedNode.CryptoPrivateKey, cryptedNode.Key, cryptedNode.Iv)
	privkey, pubkey := secp256k1.PrivKeyFromBytes(privKeyBytes)
	address := crypto.PubKey2Address(pubkey)
	return &accountTypes.Node{
		Address:    &address,
		PrivateKey: privkey,
		ChainId:    cryptedNode.ChainId,
		ChainCode:  cryptedNode.ChainCode,
	}
}

type FileStore struct {
	keysDirPath string
}

func NewFileStore(keyStoreDir string) FileStore {
	if !common.IsDirExists(keyStoreDir) {
		err := os.Mkdir(keyStoreDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return FileStore{
		keysDirPath: keyStoreDir,
	}
}

// GetKey read key in file
func (fs FileStore) GetKey(addr *crypto.CommonAddress, auth string) (*accountTypes.Node, error) {
	contents, err := ioutil.ReadFile(fs.JoinPath(addr.Hex()))
	if err != nil {
		return nil, err
	}

	node, err := bytesToCryptoNode(contents, auth)
	if err != nil {
		return nil, err
	}

	//ensure ressult after read and decrypto correct
	if node.Address.Hex() != addr.Hex() {
		return nil, fmt.Errorf("key content mismatch: have address %x, want %x", node.Address, addr)
	}
	return node, nil
}

// store the key in file encrypto
func (fs FileStore) StoreKey(key *accountTypes.Node, auth string) error {
	iv, err := common.GenUnique()
	if err != nil {
		return err
	}
	cryptoNode := &CryptedNode{
		PrivateKey: key.PrivateKey,
		ChainId:    key.ChainId,
		ChainCode:  key.ChainCode,
		Key:        []byte(auth),
		Iv:         iv[:16],
	}
	cryptoNode.EnCrypt()
	content, err := json.Marshal(cryptoNode)
	if err != nil {
		return err
	}
	return writeKeyFile(fs.JoinPath(key.Address.Hex()), content)
}

// ExportKey export all key in file by password
func (fs FileStore) ExportKey(auth string) ([]*accountTypes.Node, error) {
	persistedNodes := []*accountTypes.Node{}
	err := common.EachChildFile(fs.keysDirPath, func(path string) (bool, error) {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			log.Error("read key store error ", "Msg", err.Error())
			return false, err
		}

		node, err := bytesToCryptoNode(contents, auth)
		if err != nil {
			return false, err
		}

		if err != nil {
			log.Error("read key store error ", "Msg", err.Error())
			return false, err
		}
		persistedNodes = append(persistedNodes, node)
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return persistedNodes, nil
}

// JoinPath return keystore directory
func (fs FileStore) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(fs.keysDirPath, filename)
}

func writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

func writeKeyFile(file string, content []byte) error {
	name, err := writeTemporaryKeyFile(file, content)
	if err != nil {
		return err
	}
	return os.Rename(name, file)
}

// DbStore use leveldb as the storegae
type DbStore struct {
	dbDirPath string
	db        *leveldb.DB
}

func NewDbStore(dbStoreDir string) DbStore {
	if !common.IsDirExists(dbStoreDir) {
		err := os.Mkdir(dbStoreDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	db, err := leveldb.OpenFile("account_db", nil)
	if err != nil {
		panic(err)
	}
	return DbStore{
		dbDirPath: dbStoreDir,
		db:        db,
	}
}

// GetKey read key in db
func (db *DbStore) GetKey(addr crypto.CommonAddress, auth string) (*accountTypes.Node, error) {
	bytes := []byte{0}
	node, err := bytesToCryptoNode(bytes, auth)
	if err != nil {
		return nil, err
	}

	//ensure ressult after read and decrypto correct
	if node.Address.Hex() != addr.Hex() {
		return nil, fmt.Errorf("key content mismatch: have address %x, want %x", node.Address, addr)
	}
	return node, nil
}

// store the key in db after encrypto
func (dbStore *DbStore) StoreKey(key *accountTypes.Node, auth string) error {
	iv, err := common.GenUnique()
	if err != nil {
		return err
	}
	cryptoNode := &CryptedNode{
		PrivateKey: key.PrivateKey,
		ChainId:    key.ChainId,
		ChainCode:  key.ChainCode,
		Key:        []byte(auth),
		Iv:         iv[:16],
	}
	cryptoNode.EnCrypt()
	content, err := json.Marshal(cryptoNode)
	if err != nil {
		return err
	}
	addr := crypto.PubKey2Address(key.PrivateKey.PubKey()).Hex()
	return dbStore.db.Put([]byte(addr), content, nil)
}

// ExportKey export all key in db by password
func (dbStore *DbStore) ExportKey(auth string) ([]*accountTypes.Node, error) {
	dbStore.db.NewIterator(nil, nil)
	iter := dbStore.db.NewIterator(nil, nil)
	persistedNodes := []*accountTypes.Node{}
	for iter.Next() {
		value := iter.Value()

		node, err := bytesToCryptoNode(value, auth)
		if err != nil {
			log.Error("read key store error ", "Msg", err)
			continue
		}
		persistedNodes = append(persistedNodes, node)
	}
	return persistedNodes, nil
}

// JoinPath return the db file path
func (dbStore *DbStore) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(dbStore.dbDirPath, "db") //dbfile fixed datadir
}

// bytesToCryptoNode cocnvert given bytes and password to a node
func bytesToCryptoNode(data []byte, auth string) (node *accountTypes.Node, errRef error) {
	defer func() {
		if err := recover(); err != nil {
			errRef = errors.New("decryption failed")
		}
	}()
	cryptoNode := new(CryptedNode)
	if err := json.Unmarshal(data, cryptoNode); err != nil {
		return nil, err
	}
	cryptoNode.Key = []byte(auth)
	node = cryptoNode.DeCrypt()
	return
}

// accountCache This is used for buffering real storage and upper applications to speed up reading.
// TODO If the write speed becomes a bottleneck, write caching can be added
type accountCache struct {
	store       keyStore //  This points to a de facto storage.
	keyStoreDir string
	nodes       []*accountTypes.Node
	rlock       sync.RWMutex
}

// NewAccountCache receive an path and password as argument
// path refer to  the file that contain all key
// password used to decrypto content in key file
func NewAccountCache(keyStoreDir string, password string) (*accountCache, error) {
	ac := &accountCache{
		keyStoreDir: keyStoreDir,
		store:       NewFileStore(keyStoreDir),
	}
	persistedNodes, err := ac.store.ExportKey(password)
	if err != nil {
		return nil, err
	}
	ac.nodes = persistedNodes
	return ac, nil
}

// GetKey Get the private key by address and password
// Notice if you wallet is locked ,private key cant be found
func (ac *accountCache) GetKey(addr *crypto.CommonAddress, auth string) (*accountTypes.Node, error) {
	ac.rlock.RLock()
	defer ac.rlock.RUnlock()

	for _, node := range ac.nodes {
		if node.Address.Hex() == addr.Hex() {
			return node, nil
		}
	}
	return nil, errors.New("key not found")
}

// ExportKey export all key in cache by password
func (ac *accountCache) ExportKey(auth string) ([]*accountTypes.Node, error) {
	return ac.nodes, nil
}

// StoreKey store key local storage medium
func (ac *accountCache) StoreKey(k *accountTypes.Node, auth string) error {
	ac.rlock.Lock()
	defer ac.rlock.Unlock()

	err := ac.store.StoreKey(k, auth)
	if err != nil {
		return errors.New("save key failed" + err.Error())
	}
	ac.nodes = append(ac.nodes, k)
	return nil
}

func (ac *accountCache) ReloadKeys(auth string) error {
	ac.rlock.Lock()
	defer ac.rlock.Unlock()

	for _, node := range ac.nodes {
		if node.PrivateKey == nil {
			key, err := ac.store.GetKey(node.Address, auth)
			if err != nil {
				return err
			} else {
				node.PrivateKey = key.PrivateKey
			}
		}
	}
	return nil
}

func (ac *accountCache) ClearKeys() {
	ac.rlock.Lock()
	defer ac.rlock.Unlock()

	for _, node := range ac.nodes {
		node.PrivateKey = nil
	}
}

// JoinPath refer to local file
func (ac *accountCache) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(ac.keyStoreDir, filename)
}
