// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"math/big"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/deployment/sapphirelocalnet"
	"github.com/dmigwi/dhamana-protocol/client/deployment/sapphiremainnet"
	"github.com/dmigwi/dhamana-protocol/client/deployment/sapphiretestnet"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
)

// isDeployedNetMatching confirms that the provided network matches the deployed
// network.
func isDeployedNetMatching(net utils.NetworkType) bool {
	switch net {
	case utils.SapphireTestnet:
		return net == utils.ToNetType(sapphiretestnet.GetNetwork())
	case utils.SapphireLocalnet:
		return net == utils.ToNetType(sapphirelocalnet.GetNetwork())
	case utils.SapphireMainnet:
		return net == utils.ToNetType(sapphiremainnet.GetNetwork())
	default:
		return false
	}
}

// getContractAddress returns the deployed contract address for the provided
// network.
func getContractAddress(net utils.NetworkType) (address common.Address) {
	switch net {
	case utils.SapphireTestnet:
		address = common.HexToAddress(sapphiretestnet.GetContractAddress())
	case utils.SapphireLocalnet:
		address = common.HexToAddress(sapphirelocalnet.GetContractAddress())
	case utils.SapphireMainnet:
		address = common.HexToAddress(sapphiremainnet.GetContractAddress())
	}
	return
}

// getChainID returns the chain ID of the network associated with the deployed contract.
func getChainID(net utils.NetworkType) (chainID *big.Int) {
	switch net {
	case utils.SapphireTestnet:
		chainID = big.NewInt(int64(sapphiretestnet.GetChainID()))
	case utils.SapphireLocalnet:
		chainID = big.NewInt(int64(sapphirelocalnet.GetChainID()))
	case utils.SapphireMainnet:
		chainID = big.NewInt(int64(sapphiremainnet.GetChainID()))
	}
	return
}

// getDeploymentTime returns when the deployed contract was deployed.
func getDeploymentTime(net utils.NetworkType) (timestamp time.Time) {
	switch net {
	case utils.SapphireTestnet:
		timestamp = time.Unix(int64(sapphiretestnet.GetDeploymentTime()), 0)
	case utils.SapphireLocalnet:
		timestamp = time.Unix(int64(sapphirelocalnet.GetDeploymentTime()), 0)
	case utils.SapphireMainnet:
		timestamp = time.Unix(int64(sapphiremainnet.GetDeploymentTime()), 0)
	}
	return
}

// getDeployedTxHash returns the transaction hash when the current contract was deployed.
func getDeployedTxHash(net utils.NetworkType) (tx common.Hash) {
	switch net {
	case utils.SapphireTestnet:
		tx = common.HexToHash(sapphiretestnet.GetTransactionHash())
	case utils.SapphireLocalnet:
		tx = common.HexToHash(sapphirelocalnet.GetTransactionHash())
	case utils.SapphireMainnet:
		tx = common.HexToHash(sapphiremainnet.GetTransactionHash())
	}
	return
}

// getDeployedBlock returns the block number when the current contract was deployed.
func getDeployedBlock(net utils.NetworkType) (block uint64) {
	switch net {
	case utils.SapphireTestnet:
		block = sapphiretestnet.GetDeploymentBlock()
	case utils.SapphireLocalnet:
		block = sapphirelocalnet.GetDeploymentBlock()
	case utils.SapphireMainnet:
		block = sapphiremainnet.GetDeploymentBlock()
	}
	return
}
