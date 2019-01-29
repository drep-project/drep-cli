# DREP CLI

DREP Command Line Interface, a JSON-RPC based client for connecting to the DREP chain to query information on the chain
and send transactions.

# Features

1. Provides a decentralized ID system and data privacy protection
2. Integrates the underlying chain with both smart contracts and novel smart reputation pipes.
3. Provides the DREP reputation protocols and reputation templates.
4. Has high performance and high compatibilities.

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

# Usage

To start the DREP Cli

```
 drepcli http://127.0.0.1:15645
```

After the DREP Cli is started, it will connect to the DREP chain.
Meanwhile, you can input commands into the interface and perform operations to the chain.

# APIs

## Blocks and balances

```
[{
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
 }...]
```
* getBalance

|Method|getBalance|
|---|---|
|Parameters| 1: Address<br>2: Chain id|
|Description|Get the account's balance with respect to the address|
|Returns|Number|
|Example|db.getBalance("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b",0)|
|Example Return| 0|

* getBlock

|Method|getBlock|
|---|---|
|Parameters|Block height|
|Description|Get the corresponding block|
|Returns|Object|
|Example|db.getBlock(10)|

```json
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

* getBlocksFrom

|Method|getBlocksFrom|
|---|---|
|Parameters|1:The height of the first block to get<br>2:The number of the blocks to get|
|Description|Get the blocks between two heights|
|Returns|Array|
|Example| db.getBlocksFrom(10,2)|

```json
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
* getByteCode

|Method|getByteCode|
|---|---|
|Parameters|1:Address<br>2:Chain id|
|Description| get bytecode of a stored smart contract by its address and chainID |
|Returns|Array|
|Example|  chain.check("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00") |

* getCodeHash

|Method|getCodeHash|
|---|---|
|Parameters|1:Address<br>2:Chain id|
|Description|get the bytecode hash value of a smart contract by its address and chainID| 
|Returns|Array|
|Example|  chain.check("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00") |


* getHighestBlock

|Method|getHighestBlock|
|---|---|
|Parameters|None|
|Description|Get latest block|
|Returns|Object|
|Example | db.getHighestBlock()|

```json
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

* getMaxHeight

|Method|getMaxHeight|
|---|---|
|Parameters|None|
|Description|Get current block's height|
|Returns|Number|
|Example Return| db.getMaxHeight()|

* getMostRecentBlocks

|Method|getMostRecentBlocks|
|---|---|
|Parameters|The number of the newest blocks |
|Description|Get newest blocks data via the number parameter|
|Returns|Array|
|Example | db.getMostRecentBlocks(2)|

```json
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
* getNonce

|Method|getNonce|
|---|---|
|Parameters|1:Address<br>2:Chain id|
|Description|Get value fo nonce in the chain via address and chain id|
|Returns|Number|
|Example| db.getNonce("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00")|

### Chain
* me

|Method|me|
|---|---|
|Parameters|None|
|Description|Get information about yourself|
|Returns|Object|
|Example Return| chain.me()|

```json
{
  addr: "0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b",
  balance: 0,
  chainId: 0,
  nonce: 0
}
```

* check

|Method|check|
|---|---|
|Parameters|1: Address<br>2: Chain id|
|Description|Get information of the chain via address|
|Returns|Obeject|
|Example|  chain.check("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00") |

```json
{
  Balance: 0,
  ByteCode: null,
  CodeHash: "0x0000000000000000000000000000000000000000000000000000000000000000",
  Nonce: 0,
  Reputation: 0
}
```
* checkNonce

|Method|checkNonce|
|---|---|
|Parameters|Address|
|Description|Get nonce via address|
|Returns|Number|
|Example| chain.checkNonce("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b")|

* checkBalance

|Method|checkBalance|
|---|---|
|Parameters|Address|
|Description|Get balance via address|
|Returns|Number|
|Example| chain.checkBalance("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b")|

* create

|Method|create|
|---|---|
|Parameters|The code of contract|
|Description|Create contract|
|Returns|String|
|Example| chain.create("60806040523480...")|

* call

|Method|call|
|---|---|
|Parameters|1:Address<br>2:Chain id<br>3:Contract address<br>4:Amount<br>5:Whether the read-only|
|Description|Invoke contract|
|Returns|String|
|Example Return| chain.call("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00","511bede065df9c06e98ca75df9f0584f48a0a4f31a66104a5c8095377d69a60f","10",false)|

* send

|Method|getAllBlocks|
|---|---|
|Parameters|1:Adress<br>2:Chain id<br>3:Amount|
|Description|Get data of all the blocks|
|Returns|Array|
|Example| chain.send("0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b","0x00","10")|

### Account
* addressList

|Method|addressList|
|---|---|
|Parameters|None|
|Description|Print the list  of local accounts|
|Returns|Array|
|Example | account.addressList()|

```json
["0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b"]
```
* createAccount

|Method|createAccount|
|---|---|
|Parameters|None|
|Description|Create account|
|Returns|String|
|Example| account.createAccount()|

```json
0x772dec19e0b0b2d63a57a3a7fb03fc066d915e6b
```
* dumpPrikey

|Method|dumpPrikey|
|---|---|
|Parameters|Address|
|Description|Export privekey|
|Returns|String|
|Example| account.dumpPrikey()|

* open

|Method|open|
|---|---|
|Parameters|Password|
|Description|Open wallet|
|Returns|None|
|Example| account.open()|

* close

|Method|close|
|---|---|
|Parameters|None|
|Description|Close wallet|
|Returns|None|
|Example| account.close()|

* lock

|Method|lock|
|---|---|
|Parameters|None|
|Description|Lock password|
|Returns|None|
|Example| account.lock()|

* unLock

|Method|open|
|---|---|
|Parameters|Password|
|Description|Unlock wallet|
|Returns|None|
|Example| account.unLock()|

# License

The drep-cli library is licensed under the GNU Lesser General Public License v3.0.