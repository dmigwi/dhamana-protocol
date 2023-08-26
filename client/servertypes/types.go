// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package servertypes

import (
	"encoding/json"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
)

// ---------Accepted HTTP Request and Response structure-----------

// RPCMessage defines the structure accepted for all requests and responses.
// This struct is compatible with JSON-RPC version 2.0.
type RPCMessage struct {
	ID      uint16       `json:"id"`
	Version string       `json:"jsonrpc"`          // required on a request and a response.
	Method  utils.Method `json:"method,omitempty"` // required on a request
	Sender  *SenderInfo  `json:"sender,omitempty"` // required on a request

	Params []interface{}   `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *RPCError       `json:"error,omitempty"`
}

// SenderInfo defines the required sender information attached in every request.
type SenderInfo struct {
	Address common.Address `json:"address"`
	// SigningKey must be encrypted with the session's public key before being sent.
	// Failure to do so could expose the actual user key to hackers.
	// It must be signed via the diffie-hellman passed pubkey.
	SigningKey string `json:"signingkey,omitempty"`
}

// RPCError defines the error message information sent to the user on happening.
type RPCError struct {
	Code    uint16      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ---------HTTP Result Response Types-----------

// ServerKeyResp defines the response returned once the server public key is
// requested by a client.
type ServerKeyResp struct {
	Pubkey string `json:"pubkey"`
	Expiry uint64 `json:"expiry"` // timestamp in seconds at UTC timezone

	// private field ignored by the JSON encoder.
	SharedKey []byte `json:"-"` // Generate using the remote Pubkey + local private key.
}

// BondResp defines the response returned in an array form that is returned
// when get bonds local type method is queried by the client.
type BondResp struct {
	BondAddress common.Address `json:"bond_address"`
	Issuer      common.Address `json:"issuer_address"`
	CreatedTime time.Time      `json:"created_time"`
	CouponRate  uint8          `json:"coupon_rate"`
	Currency    uint8          `json:"currency"`
	LastStatus  uint8          `json:"last_status"`
}

// BondByAddressResp defines the complete bond details excluding the secure
// details. Secure bond details require a separate request to access them.
type BondByAddressResp struct {
	*BondResp
	Holder          common.Address `json:"holder_address"`
	CreatedAtBlock  uint64         `json:"created_at_block"`
	Principal       uint64         `json:"principal"`
	CouponDate      uint8          `json:"coupon_date"`
	MaturityDate    time.Time      `json:"maturity_date"`
	IntroMessage    string         `json:"intro_msg"`
	LastUpdate      time.Time      `json:"last_update"`
	LastSyncedBlock uint64         `json:"last_synced_block"`
}

// LastSyncedBlockResp defines the block last synced.
type LastSyncedBlockResp uint64

// ChatMsgsResp defines the response returned in an array form when get
// chats local type method is queried by the client.
type ChatMsgsResp struct {
	Sender          common.Address `json:"sender"`
	BondAddress     common.Address `json:"bond_address"`
	Message         string         `json:"chat_msg"`
	CreatedTime     time.Time      `json:"created_at"`
	LastSyncedBlock uint64         `json:"last_synced_block"`
}

// packServerError packs the errors identified into a response ready to be sent
// to the client.
func (msg *RPCMessage) PackServerError(shortErr, desc error) {
	msg.Error = &RPCError{
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
func (msg *RPCMessage) PackServerResult(data interface{}) {
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

// Reader interface implementation for type BondResp.
func (r *BondResp) Read(fn func(fields ...any) error) (interface{}, error) {
	var resp BondResp
	err := fn(&resp.BondAddress, &resp.Issuer, &resp.CreatedTime, &resp.CouponRate,
		&resp.Currency, &resp.LastStatus,
	)
	return &resp, err
}

// Reader interface implementation for type BondByAddressResp.
func (r *BondByAddressResp) Read(fn func(fields ...any) error) (interface{}, error) {
	var resp BondByAddressResp
	err := fn(&resp.BondAddress, &resp.Issuer, &resp.Holder, &resp.CreatedTime,
		&resp.CreatedAtBlock, &resp.Principal, &resp.CouponRate,
		&resp.CouponDate, &resp.MaturityDate, &resp.Currency, &resp.IntroMessage,
		&resp.LastStatus, &resp.LastUpdate, &resp.LastSyncedBlock,
	)
	return &resp, err
}

// Reader interface implementation for type LastSyncedBlockResp.
func (r *LastSyncedBlockResp) Read(fn func(fields ...any) error) (interface{}, error) {
	var resp uint32

	err := fn(&resp)
	return &resp, err
}

// Reader interface implementation for type ChatMsgsResp.
func (r *ChatMsgsResp) Read(fn func(fields ...any) error) (interface{}, error) {
	var resp ChatMsgsResp
	err := fn(&resp.Sender, &resp.BondAddress, &resp.Message,
		&resp.CreatedTime, &resp.LastSyncedBlock,
	)
	return &resp, err
}
