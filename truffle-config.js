const HDWalletProvider = require("@truffle/hdwallet-provider");
const sapphire = require("@oasisprotocol/sapphire-paratime");

module.exports = {
    contracts_directory: "./contracts",
    migrations_directory: "./contracts/migrations",
    networks: {
      development: {
        host: "127.0.0.1",
        port: 8545,
        network_id: "*"
      },
      sapphire_mainnet: {
        provider: () =>
          sapphire.wrap(
            new HDWalletProvider(
              [process.env.PRIVATE_KEY],
              "https://sapphire.oasis.io"
            )
          ),
        network_id: 0x5afe,
      },
      sapphire_testnet: {
        provider: () =>
          sapphire.wrap(
            new HDWalletProvider(
              [process.env.PRIVATE_KEY],
              "https://testnet.sapphire.oasis.dev"
            )
          ),
        network_id: 0x5aff,
      },
      sapphire_localnet: {
        provider: () =>
          sapphire.wrap(
            new HDWalletProvider(
              [process.env.PRIVATE_KEY],
              "http://localhost:8545"
            )
          ),
        network_id: 0x5afd,
      }
    },
    compilers: {
      solc: {
        version: "0.8.10"
      }
    }
  }