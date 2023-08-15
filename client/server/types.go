// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

// ---------HTTP Response Types-----------

// serverKeyResp defines the response returned once the server public key is
// requested by a client.
type serverKeyResp struct {
	Pubkey string `json:"pubkey"`
	Expiry uint64 `json:"expiry"` // timestamp in seconds at UTC timezone

	// private fields not exported
	sharedKey []byte // Generate using the remote Pubkey + local private key.
}
