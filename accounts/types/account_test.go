package types

import (
    "encoding/hex"
    "fmt"
    "github.com/drep-project/drepcli/common"
    "testing"
)

func TestNewRootAccount(t *testing.T) {
    var parent *Node
    var chainId common.ChainIdType
    account, err := NewNormalAccount(parent, chainId)
    fmt.Println("err: ", err)
    fmt.Println("account: ", account)
    fmt.Println("prv:   ", hex.EncodeToString(account.Node.PrivateKey.Serialize()))
    fmt.Println("pub x: ", hex.EncodeToString(account.Node.PrivateKey.PubKey().X.Bytes()))
    fmt.Println("pub y: ", hex.EncodeToString(account.Node.PrivateKey.PubKey().Y.Bytes()))
    fmt.Println("address: ", account.Node.Address.Hex())
    fmt.Println("save err: ", err)
}

func TestNewChildAccount(t *testing.T) {
    var parent *Node
    var chainId common.ChainIdType
    root, err := NewNormalAccount(parent, chainId)
    fmt.Println("root err: ", err)
    fmt.Println("root: ")
    fmt.Println("address: ", root.Address.Hex())
    fmt.Println("chainId: ", root.Node.ChainId)
    fmt.Println("balance: ", root.Storage.Balance)
    fmt.Println("nonce: ", root.Storage.Nonce)
    fmt.Println("byteCode: ", root.Storage.ByteCode)
    fmt.Println("codeHash: ", root.Storage.CodeHash)
    fmt.Println()
    var cid common.ChainIdType = [common.ChainIdSize]byte{1, 2, 3}
    child, err := NewNormalAccount(root.Node, cid)
    fmt.Println("child err: ", err)
    fmt.Println("child: ", child)
    fmt.Println("address: ", child.Address.Hex())
    fmt.Println("chainId: ", child.Node.ChainId)
    fmt.Println("balance: ", child.Storage.Balance)
    fmt.Println("nonce: ", child.Storage.Nonce)
    fmt.Println("byteCode: ", child.Storage.ByteCode)
    fmt.Println("codeHash: ", child.Storage.CodeHash)
}