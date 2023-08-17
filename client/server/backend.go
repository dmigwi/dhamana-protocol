// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
)

// sessionTime defines the duration when the server public key is valid.
const sessionTime = time.Minute * 10

var (
	// ZeroAddress defines an empty address value.
	ZeroAddress = common.HexToAddress("")

	sessionKeys sync.Map
)

// decodeReqBody attempts to extract contents of the request passed, if an error
// occured a response in bytes is returned. isSignerKeyRequired is used to set
// when existence of the signer key should be checked.
func decodeReqBody(req *http.Request, msg *rpcMessage, isSignerKeyRequired bool) {
	var msgError, err error

	// creates the error response to be returned
	defer func() {
		if msgError != nil {
			msg.packServerError(msgError, err)
		}
	}()

	// extract the request body contents
	if err = json.NewDecoder(req.Body).Decode(&msg); err != nil {
		msgError = utils.ErrInvalidJSON
		return
	}

	// Checks for JSON-RPC version mismatch
	if msg.Version != utils.JSONRPCVersion {
		msgError = utils.ErrInvalidReq
		err = fmt.Errorf("expected JSON-RPC version %s but found %s",
			utils.JSONRPCVersion, msg.Version)
		return
	}

	// Check for method parameter exists
	if msg.Method == "" {
		msgError = utils.ErrMethodMissing
		err = errors.New("expected a method to be provided")
		return
	}

	// Check the sender's address exists
	if msg.Sender == nil || msg.Sender.Address == ZeroAddress {
		msgError = utils.ErrSenderAddrMissing
		err = errors.New("expected sender address to be provided")
		return
	}

	// Check for the signer key if required.
	if isSignerKeyRequired && (msg.Sender == nil || msg.Sender.SigningKey == "") {
		msgError = utils.ErrSignerKeyMissing
		err = errors.New("expected sender signer key to be provided")
		return
	}
	return
}

// writeResponse writes the response using the provided response writter.
func writeResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("response writter failed: %v", err)
	}
}

// welcomeTextFunc is used to confirm the successful connection to the server.
func (s *ServerConfig) welcomeTextFunc(w http.ResponseWriter, req *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	writeResponse(w, utils.WelcomeText)
}

// serverPubkey recieves the client pubkey and sends back session public key
// to be used with diffie-hellman key exchange algorithm.
// This server public key has an expiry date attached to it, after which the
// client must fetch a new server pubkey to keep the communication alive.
// A session server pubkey is mapped to a specific user address.
func (s *ServerConfig) serverPubkey(w http.ResponseWriter, req *http.Request) {
	var msg rpcMessage
	decodeReqBody(req, &msg, false)
	if msg.Error != nil {
		writeResponse(w, msg)
		return
	}

	// decode param
	var param []string
	if err := json.Unmarshal(msg.Params, &param); err != nil {
		msg.packServerError(utils.ErrUnknownParam, err)
		writeResponse(w, msg)
		return
	}

	// Only one parameter is of type string is expected.
	if len(param) != 1 {
		msg.packServerError(utils.ErrMissingParams, nil)
		writeResponse(w, msg)
		return
	}

	privKey, err := utils.GeneratePrivKey()
	if err != nil {
		msg.packServerError(utils.ErrInternalFailure, nil)
		writeResponse(w, msg)
		return
	}

	sharedkey, err := privKey.ComputeSharedKey(param[0])
	if err != nil {
		err := errors.New("invalid public key used")
		msg.packServerError(utils.ErrInternalFailure, err)
		writeResponse(w, msg)
		return
	}

	// server public key is valid for 10 minutes after which a new public key
	// must be requested.
	data := serverKeyResp{
		Pubkey: privKey.PubKeyToHexString(),
		Expiry: uint64(time.Now().UTC().Add(sessionTime).Unix()),
	}

	msg.packServerResult(data)
	writeResponse(w, msg)

	// Appends the sharedkey after sending the user response. This shared key is
	// what this client should use to encrypt information shared with the server.
	data.sharedKey = sharedkey
	sessionKeys.Store(msg.Sender.Address, data)
}

// backendQueryFunc recieves all the requests made to the contracts.
func (s *ServerConfig) backendQueryFunc(w http.ResponseWriter, req *http.Request) {
	var msg rpcMessage

	decodeReqBody(req, &msg, false)
	if msg.Error != nil {
		writeResponse(w, msg)
		return
	}

	data, ok := sessionKeys.Load(msg.Sender.Address)
	if !ok {
		msg.packServerError(utils.ErrExpiredServerKey, nil)
		writeResponse(w, msg)
		return
	}

	// extracts the private key from the signing key sent. The private key is
	// required to sign all tx by the current sender.
	privKey, err := utils.Decrypt(data.(serverKeyResp).sharedKey, msg.Sender.SigningKey)
	if err != nil {
		msg.packServerError(utils.ErrInvalidSigningKey, err)
		writeResponse(w, msg)
		return
	}

	s.backend.SetClientSigningKey(privKey)

	// Create an authorized transactor and call the store function
	auth := s.backend.Transactor(msg.Sender.Address)

	transactor := contracts.ChatTransactorRaw{
		Contract: &s.bondChat.ChatTransactor,
	}

	tx, err := transactor.Transact(auth, msg.Method, msg.Params)
	if err != nil {
		msg.packServerError(utils.ErrInternalFailure, err)
		writeResponse(w, msg)
		return
	}

	msg.packServerResult(tx)
	writeResponse(w, msg)
}
