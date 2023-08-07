// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ServerConfig struct {
	ContractAddr common.Address
	Network      utils.NetworkType
}

const privateKey = "61e91868454365a28f4f9724ef3aaa7df0c09c16883338900a1b3dac197c89f0"

func NewServer(contractAddr string, network string) *ServerConfig {
	return &ServerConfig{
		ContractAddr: common.HexToAddress(contractAddr),
		Network:      utils.ToNetType(network),
	}
}

func (s *ServerConfig) Connection() {
	// Create an IPC based RPC connection to a remote node and instantiate a contract binding
	conn, err := ethclient.Dial(utils.Networks[s.Network].DefaultGateway)
	if err != nil {
		log.Fatalf("Failed to connect to the Sapphire Paratime client: %v", err)
	}

	chatInstance, err := contracts.NewChat(s.ContractAddr, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Chat contract: %v", err)
	}

	// Create an authorized transactor and call the store function
	auth, err := bind.NewTransactorWithChainID(strings.NewReader(privateKey), "strong_password", big.NewInt(420))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	// Call the store() function
	tx, err := chatInstance.CreateBond(auth)
	if err != nil {
		log.Fatalf("Failed to update value: %v", err)
	}

	fmt.Printf("Update pending: 0x%x\n", tx.Hash())
}
