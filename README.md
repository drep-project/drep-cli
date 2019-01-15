# DREP Cli

DREP Command Line Interface, a JSON-RPC based client for connecting to the DREP chain to query information on the chain
and send transactions.

# Features

// TODO:

# Installation

## Binaries (Windows/Linux/macOS)

Executable binaries can be downloaded from [our GitHub page](https://github.com/drep-project/drepcli/releases)

## Build from source code (all platforms)

Clone the repository into the appropriate $GOPATH/src/github.com/drep-project directory.


**NOTE:**  requires `Go >= 1.1.11`.

```
  git clone https://github.com/drep-project/drepcli.git
  cd drepcli
  go install
```

After building the source code sucessfully, you should see :

// TODO


# Usage

To start the DREP Cli

```
 drepcli http://127.0.0.1:15645
```

After the DREP Cli is started, it will connect to the DREP chain.
Meanwhile, you can input commands into the interface and perform operations to the chain.

# APIs

## Blocks and balances

* getBalance: function()

|   |   |
|---|---|
|Method|getBalance|
|Parameters| 1: 地址<br>2: 链id|
|Description|获取地址在链上的token|
|Returns|Number|
|Example|db.getBalance("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b",0)|
|Example Return| 0|

* getBlock: function()

|   |   |
|---|---|
|Method|getBlock|
|Parameters|区块号|
|Description|获取区块数据|
|Returns|Object|
|Example|db.getBlock(10)|

```$xslt
{
  Data: {
    TxCount: 1,
    TxList: [{
        Data: {...},
        Sig: null
    }]
  },
  Header: {
    ChainId: "33333330333333303333333033333330333333303333333033333330333333303333333033333330333333303333333033333330333333303333333033333330",
    GasLimit: "RWORgkT0AAA=",
    GasUsed: "",
    Height: 10,
    LeaderPubKey: "0x0258fc6797a561c701ca2c1336b4ea641745bd427319c50e82fa044c8a85b18616",
    MerkleRoot: "miMn/oiEs+F7TYCzjRmm7Yj9SvqRK1UxdW8pNJnyXYU=",
    MinorPubKeys: ["0x0242901062c10d36702cdf75330f74a597ef8058245027a5d2d993418c089fd704", "0x027d6958df8109037eab3fb5f2ca2b03ddf5f1952eed1348f8019c29437770b150"],
    PreviousHash: "e94imdsJSiPBxu3GqRJJ0DsJwOjzLePIx99aEcnMUqQ=",
    StateRoot: "kovoW3vQnm4jgiwOeMxNzwlAwzpEdvraCpdrLqiKUDI=",
    Timestamp: 1547280876,
    TxHashes: ["miMn/oiEs+F7TYCzjRmm7Yj9SvqRK1UxdW8pNJnyXYU="],
    Version: 1
  },
  MultiSig: {
    Bitmap: "AQE=",
    Sig: {
      R: 3.7142117789744075e+76,
      S: 9.782885030459428e+75
    }
  }
}
```

* getBlocksFrom: function()

|---|---|
|Method|getBlocksFrom|
|Parameters|1:the height of the first block to get<br>2:the number of the blocks to get|
|Description|get the blocks between two heights|
|Returns|Array|
|Example| db.getBlocksFrom(10,2)|

```$xslt
[{
    Data: {
      TxCount: 1,
      TxList: [{...}]
    },
    Header: {
      ChainId: "33333330333333303333333033333330333333303333333033333330333333303333333033333330333333303333333033333330333333303333333033333330",
      GasLimit: "RWORgkT0AAA=",
      GasUsed: "",
      Height: 10,
      LeaderPubKey: "0x0258fc6797a561c701ca2c1336b4ea641745bd427319c50e82fa044c8a85b18616",
      MerkleRoot: "miMn/oiEs+F7TYCzjRmm7Yj9SvqRK1UxdW8pNJnyXYU=",
      MinorPubKeys: ["0x0242901062c10d36702cdf75330f74a597ef8058245027a5d2d993418c089fd704", "0x027d6958df8109037eab3fb5f2ca2b03ddf5f1952eed1348f8019c29437770b150"],
      PreviousHash: "e94imdsJSiPBxu3GqRJJ0DsJwOjzLePIx99aEcnMUqQ=",
      StateRoot: "kovoW3vQnm4jgiwOeMxNzwlAwzpEdvraCpdrLqiKUDI=",
      Timestamp: 1547280876,
      TxHashes: ["miMn/oiEs+F7TYCzjRmm7Yj9SvqRK1UxdW8pNJnyXYU="],
      Version: 1
    },
    MultiSig: {
      Bitmap: "AQE=",
      Sig: {
        R: 3.7142117789744075e+76,
        S: 9.782885030459428e+75
      }
    }
}, ...]
```
* getByteCode: function()

|   |   |
|---|---|
|Method|getAllBlocks|
|Parameters|None|
|Description|获取所有的区块数据|
|Returns|Array|
|Example Return| |

* getCodeHash: function()

|   |   |
|---|---|
|Method|getAllBlocks|
|Parameters|None|
|Description|获取所有的区块数据|
|Returns|Array|
|Example Return| |


* getHighestBlock: function()

|   |   |
|---|---|
|Method|getHighestBlock|
|Parameters|None|
|Description|获取最新区块数据|
|Returns|Object|
|Example | db.getHighestBlock()|

```$xslt
{
  Data: {
    TxCount: 1,
    TxList: [{
        Data: {...},
        Sig: null
    }]
  },
  Header: {
    ChainId: "33333330333333303333333033333330333333303333333033333330333333303333333033333330333333303333333033333330333333303333333033333330",
    GasLimit: "RWORgkT0AAA=",
    GasUsed: "",
    Height: 1018,
    LeaderPubKey: "0x0258fc6797a561c701ca2c1336b4ea641745bd427319c50e82fa044c8a85b18616",
    MerkleRoot: "v1M3UKxS9xoaIMraw/UlYlBb7o+4vMkgLgGuUfUXlf8=",
    MinorPubKeys: ["0x0242901062c10d36702cdf75330f74a597ef8058245027a5d2d993418c089fd704", "0x027d6958df8109037eab3fb5f2ca2b03ddf5f1952eed1348f8019c29437770b150"],
    PreviousHash: "mCK7fbNwQpGUcVwqgYnH2Ea0a8gwbf2irAHdZqH3xYQ=",
    StateRoot: "1oGeQV/LDejnVqb9Q4uovyAiQvoYzan3FRFqJGHkq/U=",
    Timestamp: 1547288016,
    TxHashes: ["v1M3UKxS9xoaIMraw/UlYlBb7o+4vMkgLgGuUfUXlf8="],
    Version: 1
  },
  MultiSig: {
    Bitmap: "AQE=",
    Sig: {
      R: 8.467045644602432e+76,
      S: 9.513484500254109e+75
    }
  }
}
```

* getMaxHeight: function()

|   |   |
|---|---|
|Method|getMaxHeight|
|Parameters|None|
|Description|获取当前最大区块高度|
|Returns|Number|
|Example Return| db.getMaxHeight()|

* getMostRecentBlocks: function()

|   |   |
|---|---|
|Method|getMostRecentBlocks|
|Parameters|最新区块数量|
|Description|获得最新的n个区块数据|
|Returns|Array|
|Example | db.getMostRecentBlocks(2)|

```$xslt
[{
    Data: {
      TxCount: 1,
      TxList: [{...}]
    },
    Header: {
      ChainId: "33333330333333303333333033333330333333303333333033333330333333303333333033333330333333303333333033333330333333303333333033333330",
      GasLimit: "RWORgkT0AAA=",
      GasUsed: "",
      Height: 1,
      LeaderPubKey: "0x0258fc6797a561c701ca2c1336b4ea641745bd427319c50e82fa044c8a85b18616",
      MerkleRoot: "oRDZuhaSx/h2IlTXgEvSdGIUzln/V50MZIxK/PlVgrg=",
      MinorPubKeys: ["0x0242901062c10d36702cdf75330f74a597ef8058245027a5d2d993418c089fd704", "0x027d6958df8109037eab3fb5f2ca2b03ddf5f1952eed1348f8019c29437770b150"],
      PreviousHash: "IsZzW/wWnPD9zt1b+gQAd0mvyNzx/fmEtHsjEDiX+QE=",
      StateRoot: "jdME0UFb2ooUd84ymQMAQeTMhl/xuQZmb5uDlFv6RaA=",
      Timestamp: 1547291143,
      TxHashes: ["oRDZuhaSx/h2IlTXgEvSdGIUzln/V50MZIxK/PlVgrg="],
      Version: 1
    },
    MultiSig: {
      Bitmap: "AQE=",
      Sig: {
        R: 7.159555199205425e+76,
        S: 6.394261149473947e+76
      }
    }
}]
```
* getNonce: function()

|   |   |
|---|---|
|Method|getNonce|
|Parameters|1:地址<br>2:链id|
|Description|获取地址在链上的nonce值|
|Returns|Number|
|Example| db.getNonce("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00")|

### 链
* me

|   |   |
|---|---|
|Method|me|
|Parameters|None|
|Description|获取当前用户信息|
|Returns|Object|
|Example Return| chain.me()|

```$xslt
{
  addr: "0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b",
  balance: 0,
  chainId: 0,
  nonce: 0
}
```

* check

|   |   |
|---|---|
|Method|check|
|Parameters|1: 地址<br>2: 链id|
|Description|获取地址在链上的信息|
|Returns|Obeject|
|Example|  chain.check("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00") |

```$xslt
{
  Balance: 0,
  ByteCode: null,
  CodeHash: "0x0000000000000000000000000000000000000000000000000000000000000000",
  Nonce: 0,
  Reputation: 0
}
```
* checkNonce

|   |   |
|---|---|
|Method|checkNonce|
|Parameters|地址|
|Description|获取地址的nonce|
|Returns|Number|
|Example| chain.checkNonce("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b")|

* checkBalance

|   |   |
|---|---|
|Method|checkBalance|
|Parameters|地址|
|Description|获取地址的balance|
|Returns|Number|
|Example| chain.checkBalance("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b")|

* create

|   |   |
|---|---|
|Method|create|
|Parameters|合约代码|
|Description|部署合约|
|Returns|String|
|Example| chain.create("60806040523480...")|

* call

|   |   |
|---|---|
|Method|call|
|Parameters|1:地址<br>2:链id<br>3:合约地址<br>4:金额<br>5:是否只读|
|Description|调用合约|
|Returns|String|
|Example Return| chain.call("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00","511bede065df9c06e98ca75df9f0584f48a0a4f31a66104a5c8095377d69a60f","10",false)|

* send

|   |   |
|---|---|
|Method|getAllBlocks|
|Parameters|1:地址<br>2:链id<br>3:转账数量|
|Description|获取所有的区块数据|
|Returns|Array|
|Example| chain.send("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00","10")|

### 账号
* accountList

|   |   |
|---|---|
|Method|accountList|
|Parameters|None|
|Description|列出本地账号|
|Returns|Array|
|Example | account.accountList()|

* createAccount

|   |   |
|---|---|
|Method|getAllBlocks|
|Parameters|None|
|Description|创建账号|
|Returns|String|
|Example| account.createAccount()|

# Contact

# License