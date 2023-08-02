const HDWalletProvider = require("@truffle/hdwallet-provider");
const sapphire = require("@oasisprotocol/sapphire-paratime");
const mnemonicPhrase = "refuse sport student pond weekend treat churn foot uniform bonus gaze science";

module.exports = {
    contracts_directory: "./contracts",
    migrations_directory: "./contracts/migrations",
    networks: {
      development: {
        provider: () =>
          sapphire.wrap(
            new HDWalletProvider({
              // 5 accounts are created and all are necessary for the successful running of tests.
              // To startup the local sapphire network, the following command is used:
              // $ ./setup_localnet.sh
              mnemonic: { phrase: mnemonicPhrase },
              providerOrUrl: "http://localhost:8545"
            }),
          ),
        network_id: 0x5afd,
      },
      sapphire_mainnet: {
        provider: () =>
          sapphire.wrap(
            new HDWalletProvider({
              privateKeys: [ process.env.PRIVATE_KEY ],
              providerOrUrl: "https://sapphire.oasis.io"
             })
          ),
        network_id: 0x5afe,
      },
      sapphire_testnet: {
        provider: () =>
          sapphire.wrap(
            new HDWalletProvider({
              privateKeys: [ process.env.PRIVATE_KEY ], 
              providerOrUrl: "https://testnet.sapphire.oasis.dev"
            })
          ),
        network_id: 0x5aff,
      }
    },
    compilers: {
      solc: {
        version: "0.8.13"
      }
    }
  }

