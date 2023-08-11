// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"encoding/json"
	"net/http"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
)

// A value of this type can a JSON-RPC request, notification, successful response or
// error response. Which one it is depends on the fields.
type rpcMessage struct {
	Version string     `json:"jsonrpc"` // required
	Method  string     `json:"method"`  // required
	Sender  senderInfo `json:"sender"`  // required

	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *rpcError       `json:"error,omitempty"`
}

// senderInfo defines the required sender information attached in every request.
type senderInfo struct {
	Address common.Address `json:"address"`
	// SigningKey must be encrypted with the session's public key before being sent.
	// Failure to do so could expose the actual user key to hackers.
	// It must be signed via the diffie-hellman passed pubkey.
	SigningKey string `json:"signingkey,omitempty"`
}

// rpcError defines the error message information sent to the user on happening.
type rpcError struct {
	Code    uint16      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// welcomeTextFunc is used to confirm the successful connection to the server.
func (s *ServerConfig) welcomeTextFunc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Write([]byte(utils.WelcomeText))
}

// serverPubkey recieves the client pubkey and sends back session public key
// to be used with diffie-hellman key exchange algorithm.
// This server public key has an expiry date attached to it, after which the
// client must fetch a new server pubkey to keep the communication alive.
// A session server pubkey is mapped to a specific user address.
func (s *ServerConfig) serverPubkey(w http.ResponseWriter, req *http.Request) {
}

// backendQueryFunc recieves all the requests made to the contracts.
func (s *ServerConfig) backendQueryFunc(w http.ResponseWriter, req *http.Request) {
	// chatInstance, err := contracts.NewChat(s.contractAddr, backend)
	// if err != nil {
	// 	log.Errorf("failed to instantiate a Chat contract: %v", err)
	// 	return err
	// }

	// // Create an authorized transactor and call the store function
	// auth := backend.Transactor(userAddress)

	// tx, err := chatInstance.CreateBond(auth)
	// if err != nil {
	// 	log.Errorf("failed to update value: %v", err)
	// 	return err
	// }

	// log.Infof("Update pending: 0x%x", tx.Hash())
	// end := uint64(2149450)

	// ops := &bind.FilterOpts{
	// 	Start:   2149400,
	// 	End:     &end,
	// 	Context: s.ctx,
	// }

	// events, err := chatInstance.FilterNewBondCreated(ops)
	// if err != nil {
	// 	log.Error("Filter new bonds created events failed: ", err)
	// 	return err
	// }

	// for events.Next() {
	// 	fmt.Printf(" >>>> Bond Address: %v Sender Address: %v Timestamp: %v \n", events.Event.BondAddress, events.Event.Sender, events.Event.Timestamp)
	// }
}
