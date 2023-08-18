// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/sapphire"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ServerConfig defines the configuration needed to run a TLS enabled server
// that interacts with the contract backend.
type ServerConfig struct {
	serverURL    string
	datadir      string
	tlsCertFile  string
	tlsKeyFile   string
	contractAddr common.Address
	network      utils.NetworkType
	ctx          context.Context

	backend  *sapphire.WrappedBackend
	bondChat *contracts.Chat

	// sessionKeys holds the sessional access keys associated with a given user.
	sessionKeys sync.Map
}

// NewServer validates the deployment configuration information before
// creating a sapphire client wrapped around an eth client.
func NewServer(ctx context.Context, certfile, keyfile, datadir, network, serverURL string) (*ServerConfig, error) {
	// Validate deployment information first.
	net := utils.ToNetType(network)
	if !isDeployedNetMatching(net) {
		log.Error("Network mismatch")
		return nil, utils.ErrCorruptedConfig // network mismatch
	}

	log.Infof("Running on the network: %q", net)

	address := getContractAddress(net)
	if address == common.HexToAddress("") {
		log.Error("Empty Address found")
		return nil, utils.ErrCorruptedConfig // Address mismatch
	}

	log.Infof("Deployed contract address found: %q", address.String())

	// query the deployed time.
	deployedTime := getDeploymentTime(net)
	if deployedTime.Equal(time.Time{}) {
		log.Error("Zero deployment timestamp found")
		return nil, utils.ErrCorruptedConfig // timestamp mismatch
	}

	log.Infof("Contract in use was deployed on Date: %q",
		deployedTime.Format(utils.FullDateformat))

	// query the deployed transaction hash
	txHash := getDeployedTxHash(net)
	if txHash == "" {
		log.Error("Missing transaction hash")
		return nil, utils.ErrCorruptedConfig // tx hash mismatch
	}

	log.Infof("Contract in use was deployed on Tx: %q", txHash)

	// query the network params
	networkParams, err := utils.GetNetworkConfig(net)
	if err != nil {
		return nil, err
	}

	if networkParams.ChainID.Cmp(getChainID(net)) != 0 {
		log.Error("chainId mismatch")
		return nil, utils.ErrCorruptedConfig // ChainId mismatch
	}

	// Create RPC connection to a remote node and instantiate a contract binding
	conn, err := ethclient.DialContext(ctx, networkParams.DefaultGateway)
	if err != nil {
		log.Errorf("failed to connect to the Sapphire Paratime client: %v", err)
		return nil, err
	}

	log.Info("Creating a sapphire client wrapped over an eth client")

	backend, err := sapphire.WrapClient(ctx, conn, net,
		func(digest [32]byte, privateKey []byte) ([]byte, error) {
			key, err := crypto.ToECDSA(privateKey)
			if err != nil {
				log.Errorf("invalid private key: %v", err)
				return nil, utils.ErrInvalidPriKey
			}
			return crypto.Sign(digest[:], key)
		})

	// Create the chat instance to be used.
	chatInstance, err := contracts.NewChat(address, backend)
	if err != nil {
		log.Errorf("failed to instantiate a Chat contract: %v", err)
		return nil, err
	}

	return &ServerConfig{
		ctx:          ctx,
		network:      net,
		contractAddr: address,
		serverURL:    serverURL,
		datadir:      datadir,
		tlsCertFile:  certfile,
		tlsKeyFile:   keyfile,

		backend:  backend,
		bondChat: chatInstance,
	}, nil
}

// Run the actual TLS server instance using mTLS where both server and client
// must share their certificates.
func (s *ServerConfig) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.welcomeTextFunc)
	mux.HandleFunc("/backend", s.backendQueryFunc)
	mux.HandleFunc("/serverpubkey", s.serverPubkey)

	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAnyClientCert,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// Ignore the error because the url has already been validated.
	serverURL, _ := url.Parse(s.serverURL)
	srv := &http.Server{
		Addr:         serverURL.Host,
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// Generate the complete path to the cert and key files.
	certPath := filepath.Join(s.datadir, s.tlsCertFile)
	keyPath := filepath.Join(s.datadir, s.tlsKeyFile)

	log.Infof("Initiating the server on: %q", s.serverURL)

	return srv.ListenAndServeTLS(certPath, keyPath)
}
