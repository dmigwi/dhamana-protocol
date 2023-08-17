// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"encoding/json"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
)

// ---------Accepted HTTP Request and Response structure-----------

// rpcMessage defines the structure accepted for all requests and responses.
// This struct is compatible with JSON-RPC version 2.0.
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

// ---------HTTP Result Response Types-----------

// serverKeyResp defines the response returned once the server public key is
// requested by a client.
type serverKeyResp struct {
	Pubkey string `json:"pubkey"`
	Expiry uint64 `json:"expiry"` // timestamp in seconds at UTC timezone

	// private fields not exported
	sharedKey []byte // Generate using the remote Pubkey + local private key.
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
