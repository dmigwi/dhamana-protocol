// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"fmt"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/sapphire"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ServerConfig struct {
	ContractAddr common.Address
	Network      utils.NetworkType
}

var userAddress = common.HexToAddress("0xe1e2A376FEab01145F8Fb5679D964360cDd1B331")

func NewServer(contractAddr string, network string) *ServerConfig {
	s := &ServerConfig{
		ContractAddr: common.HexToAddress(contractAddr),
		Network:      utils.ToNetType(network),
	}

	log.Infof("Running on the network: (%v)", s.Network)
	return s
}

func (s *ServerConfig) Connection() {
	// Create RPC connection to a remote node and instantiate a contract binding
	conn, err := ethclient.Dial(utils.Networks[s.Network].DefaultGateway)
	if err != nil {
		log.Errorf("failed to connect to the Sapphire Paratime client: %v", err)
		return
	}

	backend, err := sapphire.WrapClient(*conn, s.Network, func(digest [32]byte) ([]byte, error) {
		// Pass in a custom signing function to interact with the signer
		key, err := crypto.ToECDSA(userAddress.Bytes())
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %v", err)
		}
		return crypto.Sign(digest[:], key)
	})

	chatInstance, err := contracts.NewChat(s.ContractAddr, conn)
	if err != nil {
		log.Errorf("failed to instantiate a Chat contract: %v", err)
		return
	}

	// Create an authorized transactor and call the store function
	auth := backend.Transactor(userAddress)

	tx, err := chatInstance.CreateBond(auth)
	if err != nil {
		log.Errorf("failed to update value: %v", err)
		return
	}

	log.Infof("Update pending: 0x%x\n", tx.Hash())
}
