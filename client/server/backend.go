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

	// sessionKeys holds the sessional access keys associated with a given user.
	sessionKeys sync.Map
)

// decodeRequestBody attempts to extract contents of the request passed, if an error
// occured a response in bytes is returned. isSignerKeyRequired is used to set
// when existence of the signer key should be checked.
// Its returns the method type depending on how it is implemented.
func decodeRequestBody(req *http.Request, msg *rpcMessage, isSignerKeyRequired bool) utils.MethodType {
	var msgError, err error

	// creates the error response to be returned
	defer func() {
		if msgError != nil {
			msg.packServerError(msgError, err)
		}
	}()

	if req.Method != http.MethodPost {
		msgError = utils.ErrInvalidReq
		err = fmt.Errorf("invalid http method %s found expected %s",
			req.Method, http.MethodPost)
		return utils.UnknownType
	}

	// extract the request body contents
	if err = json.NewDecoder(req.Body).Decode(&msg); err != nil {
		msgError = utils.ErrInvalidJSON
		return utils.UnknownType
	}

	// Checks for JSON-RPC version mismatch
	if msg.Version != utils.JSONRPCVersion {
		msgError = utils.ErrInvalidReq
		err = fmt.Errorf("expected JSON-RPC version %s but found %s",
			utils.JSONRPCVersion, msg.Version)
		return utils.UnknownType
	}

	// Check for method parameter exists
	if msg.Method == "" {
		msgError = utils.ErrMethodMissing
		err = errors.New("expected a method to be provided")
		return utils.UnknownType
	}

	// Check the sender's address exists
	if msg.Sender == nil || msg.Sender.Address == ZeroAddress {
		msgError = utils.ErrSenderAddrMissing
		err = errors.New("expected sender address to be provided")
		return utils.UnknownType
	}

	// Check for the signer key if required.
	if isSignerKeyRequired && (msg.Sender == nil || msg.Sender.SigningKey == "") {
		msgError = utils.ErrSignerKeyMissing
		err = errors.New("expected sender signer key to be provided")
		return utils.UnknownType
	}

	// validate the method passed.
	methodType, params := utils.GetMethodParams(msg.Method)
	if methodType == utils.UnknownType {
		err = fmt.Errorf("method %s not supportted", msg.Method)
		msgError = utils.ErrUnknownMethod
		return utils.UnknownType
	}

	if len(msg.Params) != len(params) {
		err = fmt.Errorf("method %s requires %d params found %d params",
			msg.Method, len(params), len(msg.Params))
		msgError = utils.ErrMissingParams
		return utils.UnknownType
	}

	// Confirm the required param types are used.
	for i, p := range msg.Params {
		paramType := utils.GetParamType(p)
		if paramType != params[i] {
			err = fmt.Errorf("expected param %s to be of type %s found it to be %s",
				p, paramType, params[i])
			msgError = utils.ErrUnknownParam
			return utils.UnknownType
		}
	}
	return methodType
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
	methodType := decodeRequestBody(req, &msg, false)
	if msg.Error != nil {
		writeResponse(w, msg)
		return
	}

	// confirm server key methods match.
	if methodType != utils.ServerKeyType {
		err := fmt.Errorf("unsupported method %s found in route %s",
			msg.Method, req.URL.Path)
		msg.packServerError(utils.ErrUnknownMethod, err)
		return
	}

	// Pass nil so that the default rand reader can be used.
	privKey, err := utils.GeneratePrivKey(nil)
	if err != nil {
		msg.packServerError(utils.ErrInternalFailure, nil)
		writeResponse(w, msg)
		return
	}

	sharedkey, err := privKey.ComputeSharedKey(msg.Params[0].(string))
	if err != nil {
		err := errors.New("invalid client public key used")
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

	methodType := decodeRequestBody(req, &msg, false)
	if msg.Error != nil {
		writeResponse(w, msg)
		return
	}

	// Only allow contract methods to be executed.
	if methodType != utils.ContractType {
		err := fmt.Errorf("unsupported method %s found in route %s",
			msg.Method, req.URL.Path)
		msg.packServerError(utils.ErrUnknownMethod, err)
		return
	}

	// Check if the server keys exists.
	data, ok := sessionKeys.Load(msg.Sender.Address)
	if !ok {
		err := errors.New("no server keys found associated with the sender")
		msg.packServerError(utils.ErrMissingServerKey, err)
		writeResponse(w, msg)
		return
	}

	// check for the server keys expiry.
	expiryTime := time.Unix(int64(data.(serverKeyResp).Expiry), 0).UTC()
	if time.Now().UTC().After(expiryTime) {
		msg.packServerError(utils.ErrExpiredServerKey, nil)
		writeResponse(w, msg)

		// Delete expired keys
		sessionKeys.Delete(msg.Sender.Address)
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
