// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

const (
//	func (*ChatTransactor).AddMessage(opts *bind.TransactOpts, _contract common.Address, _tag uint8, _message string) (*types.Transaction, error)
//
// func (*ChatTransactor).CreateBond(opts *bind.TransactOpts) (*types.Transaction, error)
// func (*ChatTransactor).SignBondStatus(opts *bind.TransactOpts, _contract common.Address) (*types.Transaction, error)
// func (*ChatTransactor).UpdateBodyInfo(opts *bind.TransactOpts, _contract common.Address, _principal uint32, _couponRate uint8, _couponDate uint32, _maturityDate uint32, _currency uint8) (*types.Transaction, error)
// func (*ChatTransactor).UpdateBondHolder(opts *bind.TransactOpts, _contract common.Address, _holder common.Address) (*types.Transaction, error)
// func (*ChatTransactor).UpdateBondStatus(opts *bind.TransactOpts, _contract common.Address, _status uint8) (*types.Transaction, error)
)

// ---------HTTP Request Types------------

// serverKeyReq defines the request recieved when the client wants the server to
// share its session public key.
type serverKeyReq struct {
	Pubkey string `json:"pubkey"`
}

// isValid checks the validity of the serverKeyReq recieved request.
func (s *serverKeyReq) isValid() bool {
	return s.Pubkey != ""
}

// ---------HTTP Response Types-----------

// serverKeyResp defines the response returned once the server public key is
// requested by a client.
type serverKeyResp struct {
	Pubkey string `json:"pubkey"`
	Expiry uint64 `json:"expiry"` // timestamp in seconds at UTC timezone
}
