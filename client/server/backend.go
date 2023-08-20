// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
)

// sessionTime defines the duration when the server public key is valid.
const sessionTime = time.Minute * 10

// ZeroAddress defines an empty address value.
var ZeroAddress = common.HexToAddress("")

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
	if isSignerKeyRequired {
		if msg.Sender == nil || msg.Sender.SigningKey == "" {
			msgError = utils.ErrSignerKeyMissing
			err = errors.New("expected sender signer key to be provided")
			return utils.UnknownType
		}
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
		msg.Params[i], err = castType(p, params[i])
		if err != nil {
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
		err := fmt.Errorf("unsupported method %s found for this route", msg.Method)
		msg.packServerError(utils.ErrUnknownMethod, err)
		writeResponse(w, msg)
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

	// sent the sender before packing the result because its zeroed while preparing
	// the client response.
	sender := msg.Sender.Address

	msg.packServerResult(data)
	writeResponse(w, msg)

	// Appends the sharedkey after sending the user response. This shared key is
	// what this client should use to encrypt information shared with the server.
	data.sharedKey = sharedkey
	s.sessionKeys.Store(sender, data)
}

// backendQueryFunc recieves all the requests made to the contracts.
func (s *ServerConfig) backendQueryFunc(w http.ResponseWriter, req *http.Request) {
	var msg rpcMessage

	methodType := decodeRequestBody(req, &msg, true)
	if msg.Error != nil {
		writeResponse(w, msg)
		return
	}

	// Only allow contract methods to be executed.
	// TODO: Allow execution of locally implementated methods too.
	if methodType != utils.ContractType {
		err := fmt.Errorf("unsupported method %s found for this route", msg.Method)
		msg.packServerError(utils.ErrUnknownMethod, err)
		writeResponse(w, msg)
		return
	}

	sender := msg.Sender.Address
	// Check if the server keys exists.
	data, ok := s.sessionKeys.Load(sender)
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
		s.sessionKeys.Delete(sender)
		return
	}

	sharedKey := data.(serverKeyResp).sharedKey
	if len(sharedKey) == 0 {
		msg.packServerError(utils.ErrInvalidSigningKey, nil)
		writeResponse(w, msg)
		return
	}

	// extracts the private key from the signing key sent. The private key is
	// required to sign all tx by the current sender.
	privKey, err := utils.DecryptAES(sharedKey, msg.Sender.SigningKey)
	if err != nil {
		msg.packServerError(utils.ErrInvalidSigningKey, err)
		writeResponse(w, msg)
		return
	}

	s.backend.SetClientSigningKey(privKey)

	// Create an authorized transactor.
	auth := s.backend.Transactor(sender)
	transactor := contracts.ChatRaw{Contract: s.bondChat}

	tx, err := transactor.Transact(auth, msg.Method, msg.Params...)
	if err != nil {
		msg.packServerError(utils.ErrInternalFailure, err)
		writeResponse(w, msg)
		return
	}

	msg.packServerResult(tx)
	writeResponse(w, msg)
}

// castType returns the parameter cast to the required parameter type.
func castType(param interface{}, pType utils.ParamType) (v interface{}, err error) {
	if pType == utils.UnsupportedType {
		return nil, fmt.Errorf("unexpected type for param %v found", param)
	}

	typeFound := "unsupported"

	switch t := param.(type) {
	case string:
		switch pType {
		case utils.AddressType:
			v = common.HexToAddress(t)
		case utils.StringType:
			v = t
		default:
			typeFound = "string"
		}
	case float64:
		// JSON distinct types do not differentiate between integers and floats.
		// JSON returns all numbers as float64 values.
		// https://www.webdatarocks.com/doc/data-types-in-json/#number
		rawInt := int(t)

		switch pType {
		case utils.Uint8Type:
			if rawInt <= math.MaxUint8 {
				v = uint8(rawInt)
			}
		case utils.Uint16Type:
			if rawInt <= math.MaxUint16 {
				v = uint8(rawInt)
			}
		case utils.Uint32Type:
			if rawInt <= math.MaxUint32 {
				v = uint8(rawInt)
			}
		default:
			typeFound = "number"
		}
	}

	// Casting to the require parameter failed due to use of incorrect parameter value
	if v == nil {
		err = fmt.Errorf("expected param %v to be of type %v but found it to be %s", param, pType, typeFound)
	}
	return
}
