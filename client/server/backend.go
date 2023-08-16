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
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// packServerError packs the errors identified into a response ready to be sent
// to the client.
func (msg *rpcMessage) packServerError(shortErr, desc error) {
	msg.Error = &rpcError{
		Code:    utils.GetErrorCode(shortErr),
		Message: shortErr.Error(),
		Data:    desc.Error(),
	}

	// Remove the unnecessary information in the response.
	msg.Sender = nil
	msg.Method = ""
	msg.Params = nil
	msg.Result = nil
}

// packServerResult packs the successful result queried into a response ready
// to be sent to the client.
func (msg *rpcMessage) packServerResult(data interface{}) {
	// Remove the unnecessary information in the response.
	msg.Error = nil
	msg.Sender = nil
	msg.Method = ""
	msg.Params = nil

	res, _ := json.Marshal(data)
	msg.Result.UnmarshalJSON(res)
}

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

		sharedKey: sharedkey,
	}

	sessionKeys.Store(msg.Sender.Address, data)

	msg.packServerResult(data)
	writeResponse(w, msg)
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

	privKey, err := utils.Decrypt(data.(serverKeyResp).sharedKey, msg.Sender.SigningKey)
	if err != nil {
		msg.packServerError(utils.ErrInvalidSigningKey, err)
		writeResponse(w, msg)
		return
	}

	s.backend.SetClientSigningKey(privKey)

	chatInstance, err := contracts.NewChat(s.contractAddr, s.backend)
	if err != nil {
		log.Errorf("failed to instantiate a Chat contract: %v", err)
		msg.packServerError(utils.ErrInternalFailure, err)
		writeResponse(w, msg)
		return
	}

	// Create an authorized transactor and call the store function
	auth := s.backend.Transactor(msg.Sender.Address)

	transactor := contracts.ChatTransactorRaw{
		Contract: &chatInstance.ChatTransactor,
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
