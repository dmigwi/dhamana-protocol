[![Contracts Tests](https://github.com/dmigwi/dhamana-protocol/actions/workflows/contract-tests.yaml/badge.svg)](https://github.com/dmigwi/dhamana-protocol/actions/workflows/contract-tests.yaml)

<h1 align="center"> Dhamana Protocol </h1>
<p align="center">
    <img alt="dhamana-protocol" src="https://github.com/dmigwi/dhamana-protocol/assets/22055953/0c0709c9-cb94-41b7-8463-dbd8f9dfb258" width="150">
  </a>
</p>
<p align="center">Bond protocol that runs on the computing layer of the sapphire OSL (Oasis Private Layer).</p>

<!--
<p align="center">
  <a href="https://itunes.apple.com/us/app/gitpoint/id1251245162?mt=8">
    <img alt="Download on the App Store" title="App Store" src="http://i.imgur.com/0n2zqHD.png" width="140">
  </a>
  <a href="https://play.google.com/store/apps/details?id=com.gitpoint">
    <img alt="Get it on Google Play" title="Google Play" src="http://i.imgur.com/mtGRPuM.png" width="140">
  </a>-->
</p>

### Introduction

[Dhamana](https://en.wiktionary.org/wiki/dhamana) is a swahili word with arabic origins that means __"*to offer guarantee*"__. **Dhamana-protocol** is a system where people unknown to each other but approved by a reputable 3rd party, engage in a discussion on how one of them will offer a (bond) security in exchange for financial support **(Crypto/Fiat)** on a project or an idea they are passionate about.
The bond issuer, initiates the conversation by describing the principal, interest rate, security e.t.c.(bond details) with a view that one of the people interested will agree on their terms and subscribe to their bond.
The system allows people to make legal binding agreements that can be used to seek legal redress should one of the parties fall short of their expectations.
Once the terms are agreed upon, financial transactions are done outside dhamana-protocol just as parties to the bond have agreed.

### Table of Contents
- [Getting Started](#getting-started)
- [Features](#features)
- [Feedback](#feedback)
- [Acknowledgments](#acknowledgments)

### Getting Started

#### => Generate Chat contract ABI golang bindings
<details>
<summary>Step only necessary if you've modified the contracts functionality.</summary>
The latest binds update are always uploaded online so no need to run this command if the original code is still intact. This go bindings enable the use of <a href="https://github.com/ethereum/go-ethereum">go-ethereum</a> library when interacting with the deployed contracts.

1. Download Abigen from <a href="https://geth.ethereum.org/downloads">Geth & Tools release</a> for your platform. Extract the zipped file and `abigen` will be one of the zipped files.
2. Install `solc`, the solidity compiler in one of the ways described <a href="https://docs.soliditylang.org/en/v0.8.9/installing-solidity.html">here</a>
3. Confirm the two installation done above were successful by:
    - **abigen**
    ```
    $abigen --version
    abigen version 1.12.0-stable-e501b3b0
    ```
    - **solc**
    ```
    $ solc --version
    solc, the solidity compiler commandline interface
    Version: 0.8.21+commit.d9974bed.Linux.g++
    ```
4. Generate ABI JSON file using:
    ```
    $ solc --abi contracts/chat.sol -o build --overwrite
    Compiler run successful. Artifact(s) can be found in directory "build".
    ```
5. Generate Contract bindings in golang using
    ```
    $ abigen --abi ./build/ChatContract.abi --pkg contracts --type Chat --out client/contracts/chat.go
    ```
</details>

#### => Contract Deployment
To deploy a contract; several networks are supported. They include `development`, `sapphire_localnet`, `sapphire_testnet` and `sapphire_mainnet`.
<details>
<summary>1. <strong>development</strong> network interacts with the raw ganache instance.</summary>

- To deploy; <a href="https://trufflesuite.com/docs/ganache/quickstart/#1-install-ganache">install ganache</a> then check for its installation. <strong>(Terminal Window 1)</strong>
    ```
    $ which ganache 
    ~/.nvm/versions/node/v18.0.0/bin/ganache 
    ```
- Start up ganache:
    ```
    $ ganache
    ganache v7.9.0 (@ganache/cli: 0.10.0, @ganache/core: 0.10.0)
    Starting RPC server

    Available Accounts
    ==================
    (0) 0xcf7c8a3f10504Df1407088Ad9eFc6D0466645E73 (1000 ETH)
    (1) 0x696aE3BA64a757deE6039C848bA831697379d1C4 (1000 ETH)
    (2) 0xA21e897831A9D00F74bD9d72CbF17FeddA8E4845 (1000 ETH)
    ...
    ```

- Run the contract deployment.  It will deploy the contract, export the deployment log and generate go config bindings for the development network. <strong>(Terminal Window 2)</strong>
    ```
    $ pnpm run deploy_development
    > dhamana-protocol@0.0.1 deploy_development ~/golang/src/github.com/dmigwi/dhamana-protocol
    > ./node_modules/.bin/truffle migrate --network development --describe-json >> build/deploy_development.log && pnpm run gen_deployment_config build/deploy_development.log

    > dhamana-protocol@0.0.1 gen_deployment_config ~/golang/src/github.com/dmigwi/dhamana-protocol
    > go run ./deployInfo  "build/deploy_development.log"

    2023/08/20 21:38:12 Log deployment information has me written to : client/deployment/development/deployment.go
    ```
</details>
<details>
<summary>2. <strong>sapphire_localnet</strong> network interacts with a local deployment of the OASIS network.</summary>

- To deploy install the backend by running a docker container. <strong>(Terminal Window 1)</strong>
    ```
    $ ./setup_localnet.sh 
    sapphire-dev 2023-07-10-gitbacd168 (oasis-core: 22.2.8, sapphire-paratime: 0.5.2, oasis-web3-gateway: 3.3.0-gitbacd168)

    Starting oasis-net-runner with sapphire...
    Starting postgresql...
    Starting oasis-web3-gateway...
    Bootstrapping network and populating account(s) (this might take a minute)...

    Available Accounts
    ==================
    (0) 0x9Bde9f59ef7b76A5283eB15F93f73aeD9aF044aF (1000 TEST)
    (1) 0xdf46bb474947756741dc1257ed0f54848606670D (1000 TEST)

    ...
    Private Keys
    ==================
    (0) 0xe32c076b6bafa297fa50d0a8dcee339edea7b52153bbe5860faa0f2424623ba4
    (1) 0xa87903873df7e7157010f341af7692326a54fa012e175bb0c91b5ed89395eb8b
    ...
    ```
- Export the private key on the command line. This Private key will be from the account funding the contract deployment. <strong>(Terminal Window 2)</strong>

    ```
    export PRIVATE_KEY="0xe32c076b6bafa297fa50d0a8dcee339edea7b52153bbe5860faa0f2424623ba4"
    ```
- Run the contract deployment.  It will deploy the contract, export the deployment log and generate go config bindings for the sapphire_localnet network.
    ```
    $ pnpm run deploy_sapphire_localnet
    > dhamana-protocol@0.0.1 deploy_sapphire_localnet ~/golang/src/github.com/dmigwi/dhamana-protocol
    > ./node_modules/.bin/truffle migrate --network sapphire_localnet --describe-json >> build/deploy_sapphire_localnet.log && pnpm run gen_deployment_config build/deploy_sapphire_localnet.log

      ⠏ Blocks: 0            Seconds: 0
    > dhamana-protocol@0.0.1 gen_deployment_config ~/golang/src/github.com/dmigwi/dhamana-protocol
    > go run ./deployInfo  "build/deploy_sapphire_localnet.log"

    2023/08/20 21:58:09 Log deployment information has me written to : client/deployment/sapphirelocalnet/deployment.go
    ```
</details>
<details>
<summary>3. <strong>sapphire_mainnet</strong> and <strong>sapphire_testnet</strong> networks deployments are way simpler, similar and straight forward. They interacts with the respective remote nodes of the OASIS network.</summary>

- Export the respectve network private key on the command line. This Private key will be from the account funding the contract deployment. <strong>(Terminal Window 1)</strong>
    ```
    export PRIVATE_KEY="0xe.....4"
    ```
- Run the contract deployment.  It will deploy the contract, export the deployment log and generate go config bindings for the sapphire_testnet network.
    ```
    $ pnpm run deploy_sapphire_testnet
    > dhamana-protocol@0.0.1 deploy_sapphire_testnet ~/golang/src/github.com/dmigwi/dhamana-protocol
    > ./node_modules/.bin/truffle migrate --network sapphire_testnet --describe-json >> build/deploy_sapphire_testnet.log && pnpm run gen_deployment_config build/deploy_sapphire_testnet.log

      ⠏ Blocks: 2            Seconds: 21
    > dhamana-protocol@0.0.1 gen_deployment_config ~/golang/src/github.com/dmigwi/dhamana-protocol
    > go run ./deployInfo  "build/deploy_sapphire_testnet.log"

    2023/08/20 22:08:44 Log deployment information has me written to : client/deployment/sapphiretestnet/deployment.go
    ```
</details>

#### => Build Binary
Depending on the network you intend to run dhamana-protocol client on, deployment configuration with go bindings have already been generated at `client/deployment/*` path. The binary generated here will be tied to those binds.
This implies that a single binary generated can interact with the same contract but running on different machine instances without the need of sharing the deployment logs.

- Generate the binary called `dhamana-protocol` using.
    ```
    $ go build -o dhamana-protocol ./client
    ```
- Initiate the server using:
    ```
    $  ./dhamana-protocol --network=sapphire_testnet
    2023-08-20 22:47:00.410 [INF] MAIN: Loading command configurations
    2023-08-20 22:47:00.411 [INF] MAIN: Using data directory: ~/.config/dhamana-protocol
    2023-08-20 22:47:00.411 [INF] SERVER: Running on the network: "SapphireTestnet"
    2023-08-20 22:47:00.411 [INF] SERVER: Deployed contract address found: "0x65043857F998FD9d7a910fCE23E692eBd1dc6878"
    2023-08-20 22:47:00.411 [INF] SERVER: Contract in use was deployed on Date: "Wed 05:36:49 2023-08-09"
    2023-08-20 22:47:00.411 [INF] SERVER: Contract in use was deployed on Tx: "0xbb502808fa8e30e49567f466f30cd300ce6b4759ec0d06e8fa16fa1f94816852"
    2023-08-20 22:47:00.411 [INF] SERVER: Creating a sapphire client wrapped over an eth client
    2023-08-20 22:47:02.364 [INF] SERVER: Initiating the server on: "https://0.0.0.0:30443"
    ```

And **Dhamana Protocol** is now live via **https://0.0.0.0:30443**

### Features

### Feedback

### Acknowledgments

