// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"net/http"

	"github.com/dmigwi/dhamana-protocol/client/utils"
)

var (
// contractDeployedBlock = 2091593
// userAddress           = common.HexToAddress("0xe1e2A376FEab01145F8Fb5679D964360cDd1B331")
// privatekey            = hexutil.MustDecode("0x61e91868454365a28f4f9724ef3aaa7df0c09c16883338900a1b3dac197c89f0")

)

// welcomeTextFunc is used to confirm the successful connection to the server.
func (s *ServerConfig) welcomeTextFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Write([]byte(utils.WelcomeText))
}

// contractQueryFunc recieves all the requests made to the contracts.
func (s *ServerConfig) contractQueryFunc(w http.ResponseWriter, req *http.Request) {
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
