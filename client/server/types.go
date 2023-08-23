// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"encoding/json"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
)

// ---------Accepted HTTP Request and Response structure-----------

// rpcMessage defines the structure accepted for all requests and responses.
// This struct is compatible with JSON-RPC version 2.0.
type rpcMessage struct {
	ID      uint16       `json:"id"`
	Version string       `json:"jsonrpc"`          // required on a request and a response.
	Method  utils.Method `json:"method,omitempty"` // required on a request
	Sender  *senderInfo  `json:"sender,omitempty"` // required on a request

	Params []interface{}   `json:"params,omitempty"`
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

// BondResp defines the response in an array form that is returned when get bonds
// local method is queried by the client.
type bondResp struct {
	BondAddress common.Address `json:"bond_address"`
	CreatedTime time.Time      `json:"created_time"`
	CouponRate  uint8          `json:"coupon_rate"`
	Currency    uint8          `json:"currency"`
	LastStatus  uint8          `json:"last_status"`
}

// bondByAddressResp defines the complete bond details excluding the secure
// details. Secure bond details require a separate request to access them.
type bondByAddressResp struct {
	*bondResp
	Issuer         common.Address `json:"issuer_address"`
	Holder         common.Address `json:"holder_address"`
	TxHash         string         `json:"tx_hash"`
	CreatedAtBlock uint32         `json:"created_at_block"`
	Principal      uint64         `json:"principal"`
	CouponDate     time.Time      `json:"coupon_date"`
	MaturityDate   time.Time      `json:"maturity_date"`
	IntroMessage   string         `json:"intro_msg"`
	LastUpdate     time.Time      `json:"last_update"`
}

// packServerError packs the errors identified into a response ready to be sent
// to the client.
func (msg *rpcMessage) packServerError(shortErr, desc error) {
	msg.Error = &rpcError{
		Code:    utils.GetErrorCode(shortErr),
		Message: shortErr.Error(),
	}

	// Data field is optional
	if desc != nil {
		msg.Error.Data = desc.Error()
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

	// encode interface to bytes
	b, _ := json.Marshal(data)

	// push the data bytes into the msg.Result.
	_ = json.Unmarshal(b, &msg.Result)
}

// Reader interface implementation for type bondResp.
func (r *bondResp) Read(fn func(fields ...any) error) (interface{}, error) {
	var resp bondResp
	err := fn(&resp.BondAddress, &resp.CreatedTime, &resp.CouponRate,
		&resp.Currency, &resp.LastStatus,
	)
	return &resp, err
}

// Reader interface implementation for type bondByAddressResp.
func (r *bondByAddressResp) Read(fn func(fields ...any) error) (interface{}, error) {
	var resp bondByAddressResp
	err := fn(&resp.BondAddress, &resp.Issuer, &resp.Holder, &resp.CreatedTime,
		&resp.TxHash, &resp.CreatedAtBlock, &resp.Principal, &resp.CouponRate,
		&resp.CouponDate, &resp.MaturityDate, &resp.Currency, &resp.IntroMessage,
		&resp.LastStatus, &resp.LastUpdate,
	)
	return &resp, err
}
