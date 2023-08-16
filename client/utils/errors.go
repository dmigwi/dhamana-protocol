// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

import "errors"

var (
	// ErrCorruptedConfig error is returned if one of the deployment configs
	// doesn't match the expected values.
	ErrCorruptedConfig = errors.New("deployment config has been corrupted. Regenerate it")

	// ErrInvalidPriKey returns if private key validation fails.
	ErrInvalidPriKey = errors.New("Invalid private key used to sign transactions")

	//-------Server Errors--------

	// serverErrorCodes is mapping of the supported server errors and their
	// respective error codes.
	serverErrorCodes = map[error]uint16{
		ErrInvalidJSON:       1000,
		ErrInvalidReq:        1001,
		ErrInternalFailure:   1002,
		ErrMethodMissing:     1003,
		ErrSenderAddrMissing: 1004,
		ErrExpiredServerKey:  1005,
		ErrSignerKeyMissing:  1006,
		ErrMissingParams:     1007,
		ErrUnknownMethod:     1008,
		ErrUnknownParam:      1009,
		ErrInvalidSigningKey: 1010,
	}

	// ErrInvalidJSON returned if an error occurred while parsing the request JSON
	// due to a malformed request data used.
	ErrInvalidJSON = errors.New("parse error")

	// ErrInvalidReq is returned if the JSON version doesn't match the supported.
	ErrInvalidReq = errors.New("invalid request")

	// ErrInternalFailure is returned if an unexpected error is returned while
	// processing the received request.
	ErrInternalFailure = errors.New("internal error")

	// ErrMethodMissing is returned if a request with a method is received.
	ErrMethodMissing = errors.New("method missing")

	// ErrSenderAddrMissing is returned if the sender's address is missing.
	ErrSenderAddrMissing = errors.New("sender address missing")

	// ErrExpiredServerKey is returned if a sender takes too long to use the server
	// key sent during the tls handshake.
	ErrExpiredServerKey = errors.New("server key expired")

	// ErrSignerKeyMissing is returned if the sender's signer key is missing.
	ErrSignerKeyMissing = errors.New("sender signer key missing")

	// ErrMissingParams is returned if one or more of the expected parameters is
	// missing. Or more than required parameters are returned.
	ErrMissingParams = errors.New("excess or missing param(s)")

	// ErrUnknownMethod is returned if the method provided in the request is
	// not currently supported.
	ErrUnknownMethod = errors.New("unknown method")

	// ErrUnknownParam is returned if one or more of the provided parameters
	// contains unexpected data.
	ErrUnknownParam = errors.New("unknown param(s)")

	// ErrInvalidSigningKey is return if decrypting the actual key from the
	// provided signing key string results into an error.
	ErrInvalidSigningKey = errors.New("invalid client signing key")
)

// GetErrorCode returns the set error code if it exists or max(uint16) if otherwise.
func GetErrorCode(err error) uint16 {
	if code, ok := serverErrorCodes[err]; ok {
		return code
	}
	return 65535 // uint16(2 ^ 16 - 1)
}
