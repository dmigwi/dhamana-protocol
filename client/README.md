# client

- Defines the client implementation that will interact with sapphire paratime

## How to generate solidity contract bindings

### Pre-requisites

1. Install Solidity compiler `solc` compatible with your platform. On Ubuntu run 
```
sudo apt-get solc
```
2. Check installed solidity compiler version using.
```
$ solc --version 
solc, the solidity compiler commandline interface
Version: 0.8.21+commit.d9974bed.Linux.g++
```
3. Download the abigen tool from official [releases](https://geth.ethereum.org/downloads) or build it from [go-ethereum](https://github.com/ethereum/go-ethereum/tree/master#executables)

### Generate Go Bindings

From the project root folder run the following commands:

1. Generate abi files by running the following command:
```
$ solc --abi contracts/chat.sol -o build
Compiler run successful. Artifact(s) can be found in directory "build".
```

2. Generate the go bindings by running the following command:
```
$ abigen --abi ./build/ChatContract.abi --pkg contracts --type Chat --out client/contracts/chat.go
```