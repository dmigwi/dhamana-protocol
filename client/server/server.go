// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"context"
	"errors"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/sapphire"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const fullDateformat = "Mon 15:04:05 2006-01-02"

type ServerConfig struct {
	contractAddr common.Address
	network      utils.NetworkType
	ctx          context.Context
	cancelFunc   context.CancelFunc

	backend bind.ContractBackend
}

var (
	// contractDeployedBlock = 2091593
	// userAddress           = common.HexToAddress("0xe1e2A376FEab01145F8Fb5679D964360cDd1B331")
	// privatekey            = hexutil.MustDecode("0x61e91868454365a28f4f9724ef3aaa7df0c09c16883338900a1b3dac197c89f0")

	// ErrCorruptedConfig error is returned if one of the deployment configs
	// doesn't match the expected values.
	ErrCorruptedConfig = errors.New("Deployment config has been corrupted. Regenerate it!")

	// ErrInvalidPriKey returns if private key validation fails.
	ErrInvalidPriKey = errors.New("Invalid private key used to sign transactions")

	// privateKeyChan is used to pass the received
	privateKeyChan = make([]byte, 1)
)

// NewServer validates the deployment configuration information before
// creating a sapphire client wrapped around an eth client.
func NewServer(ctx context.Context, network string) (*ServerConfig, error) {
	// Validate deployment information first.
	net := utils.ToNetType(network)
	if !isDeployedNetMatching(net) {
		log.Error("Network mismatch")
		return nil, ErrCorruptedConfig // network mismatch
	}

	log.Infof("Running on the network: %s", net)

	address := getContractAddress(net)
	if address == common.HexToAddress("") {
		log.Error("Empty Address found")
		return nil, ErrCorruptedConfig // Address mismatch
	}

	log.Infof("Deployed contract address found: %s", address.String())

	// query the deployed time.
	deployedTime := getDeploymentTime(net)
	if deployedTime.Equal(time.Time{}) {
		log.Error("Zero deployment timestamp found")
		return nil, ErrCorruptedConfig // timestamp mismatch
	}

	log.Infof("Contract in use was deployed on: %s", deployedTime.Format(fullDateformat))

	// query the network params
	networkParams, err := utils.GetNetworkConfig(net)
	if err != nil {
		return nil, err
	}

	if networkParams.ChainID.Cmp(getChainID(net)) != 0 {
		log.Error("chainId mismatch")
		return nil, ErrCorruptedConfig // ChainId mismatch
	}

	// Create RPC connection to a remote node and instantiate a contract binding
	conn, err := ethclient.Dial(networkParams.DefaultGateway)
	if err != nil {
		log.Errorf("failed to connect to the Sapphire Paratime client: %v", err)
		return nil, err
	}

	// generate a new context using the parent context passed.
	ctx, cancelfn := context.WithCancel(ctx)

	log.Info("Creating a sapphire client wrapped over an eth client")

	backend, err := sapphire.WrapClient(ctx, *conn, net, func(digest [32]byte, privateKey []byte) ([]byte, error) {
		key, err := crypto.ToECDSA(privateKey)
		if err != nil {
			log.Errorf("invalid private key: %v", err)
			return nil, ErrInvalidPriKey
		}
		return crypto.Sign(digest[:], key)
	})

	return &ServerConfig{
		contractAddr: address,
		network:      net,
		ctx:          ctx,
		cancelFunc:   cancelfn,

		backend: backend,
	}, nil
}

// Run runs the actual server instance.
func (s *ServerConfig) Run() {
}

// func (s *ServerConfig) Connection() error {
// 	network, err := utils.GetNetworkConfig(s.network)
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}

// 	// Create RPC connection to a remote node and instantiate a contract binding
// 	conn, err := ethclient.Dial(network.DefaultGateway)
// 	if err != nil {
// 		log.Errorf("failed to connect to the Sapphire Paratime client: %v", err)
// 		return err
// 	}

// 	backend, err := sapphire.WrapClient(s.ctx, *conn, s.network, func(digest [32]byte) ([]byte, error) {
// 		// Pass in a custom signing function to interact with the signer
// 		key, err := crypto.ToECDSA(privatekey)
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid private key: %v", err)
// 		}
// 		return crypto.Sign(digest[:], key)
// 	})

// 	chatInstance, err := contracts.NewChat(s.contractAddr, backend)
// 	if err != nil {
// 		log.Errorf("failed to instantiate a Chat contract: %v", err)
// 		return err
// 	}

// 	// Create an authorized transactor and call the store function
// 	auth := backend.Transactor(userAddress)

// 	tx, err := chatInstance.CreateBond(auth)
// 	if err != nil {
// 		log.Errorf("failed to update value: %v", err)
// 		return err
// 	}

// 	log.Infof("Update pending: 0x%x", tx.Hash())
// 	end := uint64(2149450)

// 	ops := &bind.FilterOpts{
// 		Start:   2149400,
// 		End:     &end,
// 		Context: s.ctx,
// 	}

// 	events, err := chatInstance.FilterNewBondCreated(ops)
// 	if err != nil {
// 		log.Error("Filter new bonds created events failed: ", err)
// 		return err
// 	}

// 	for events.Next() {
// 		fmt.Printf(" >>>> Bond Address: %v Sender Address: %v Timestamp: %v \n", events.Event.BondAddress, events.Event.Sender, events.Event.Timestamp)
// 	}
// }
