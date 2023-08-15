// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
)

// ZeroAddress defines an empty address value.
var ZeroAddress = common.HexToAddress("")

// A value of this type can a JSON-RPC request, notification, successful response or
// error response. Which one it is depends on the fields.
type rpcMessage struct {
	ID      uint16      `json:"id"`
	Version string      `json:"jsonrpc"`          // required on a request and a response.
	Method  string      `json:"method,omitempty"` // required on a request
	Sender  *senderInfo `json:"sender,omitempty"` // required on a request

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
	Message error       `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// decodeReqBody attempts to extract contents of the request passed, if an error
// occured a response in bytes is returned. isSignerKeyRequired is used to set
// when existence of the signer key should be checked.
func decodeReqBody(req *http.Request, msg *rpcMessage, isSignerKeyRequired bool) (response []byte) {
	var msgError, err error

	// creates the error response to be returned
	defer func() {
		if msgError != nil {
			msg.Error = &rpcError{
				Code:    utils.GetErrorCode(msgError),
				Message: msgError,
				Data:    err,
			}

			// Remove the unnecessary information in the response.
			msg.Sender = nil
			msg.Method = ""
			msg.Params = nil
			msg.Result = nil
			response, _ = json.Marshal(msg)
		}
	}()

	// extract the request body contents
	if err = json.NewDecoder(req.Body).Decode(&msg); err != nil {
		log.Errorf("invalid request body: %v", err)
		msgError = utils.ErrInvalidJSON
		return
	}

	// Checks for JSON-RPC version mismatch
	if msg.Version != utils.JSONRPCVersion {
		msgError = utils.ErrInvalidReq
		err = fmt.Errorf("expected JSON-RPC version %s but found %s",
			utils.JSONRPCVersion, msg.Version)
		log.Error(err)
		return
	}

	// Check for method parameter exists
	if msg.Method == "" {
		msgError = utils.ErrMethodMissing
		err = errors.New("expected a method to be provided")
		log.Error(err)
		return
	}

	// Check the sender's address exists
	if msg.Sender == nil || msg.Sender.Address == ZeroAddress {
		msgError = utils.ErrSenderAddrMissing
		err = errors.New("expected sender address to be provided")
		log.Error(err)
		return
	}

	// Check for the signer key if required.
	if isSignerKeyRequired && (msg.Sender == nil || msg.Sender.SigningKey == "") {
		msgError = utils.ErrSignerKeyMissing
		err = errors.New("expected sender signer key to be provided")
		log.Error(err)
		return
	}
	return
}

// writeResponse writes the response using the provided response writter.
func writeResponse(w http.ResponseWriter, response []byte) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Write(response)
}

// welcomeTextFunc is used to confirm the successful connection to the server.
func (s *ServerConfig) welcomeTextFunc(w http.ResponseWriter, _ *http.Request) {
	writeResponse(w, []byte(utils.WelcomeText))
}

// serverPubkey recieves the client pubkey and sends back session public key
// to be used with diffie-hellman key exchange algorithm.
// This server public key has an expiry date attached to it, after which the
// client must fetch a new server pubkey to keep the communication alive.
// A session server pubkey is mapped to a specific user address.
func (s *ServerConfig) serverPubkey(w http.ResponseWriter, req *http.Request) {
	var msg rpcMessage
	if response := decodeReqBody(req, &msg, false); response != nil {
		writeResponse(w, response)
		return
	}
}

// backendQueryFunc recieves all the requests made to the contracts.
func (s *ServerConfig) backendQueryFunc(w http.ResponseWriter, req *http.Request) {
	// req.Body.Read(p []byte)
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
