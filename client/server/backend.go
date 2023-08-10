// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"encoding/json"
	"net/http"

	"github.com/dmigwi/dhamana-protocol/client/utils"
)

// A value of this type can a JSON-RPC request, notification, successful response or
// error response. Which one it is depends on the fields.
type rpcMessage struct {
	ID      uint8      `json:"id"`      // required
	Version string     `json:"jsonrpc"` // required
	Method  string     `json:"method"`  // required
	Sender  senderInfo `json:"sender"`  // required

	Params json.RawMessage `json:"params,omitempty"`
	Error  *rpcError       `json:"error,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

type senderInfo struct {
	Address string `json:"address"`
	// SigningKey must be encrypted with the session's public key before being sent.
	// Failure to do so could expose the actual user key to any network hyjacker.
	SigningKey string `json:"signingkey"`
}

type rpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// welcomeTextFunc is used to confirm the successful connection to the server.
func (s *ServerConfig) welcomeTextFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Write([]byte(utils.WelcomeText))
}

// contractQueryFunc recieves all the requests made to the contracts.
func (s *ServerConfig) contractQueryFunc(w http.ResponseWriter, req *http.Request) {
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
