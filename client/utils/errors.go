// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

import "errors"

var (
	// ErrCorruptedConfig error is returned if one of the deployment configs
	// doesn't match the expected values.
	ErrCorruptedConfig = errors.New("Deployment config has been corrupted. Regenerate it!")

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
		ErrSenderInfoMissing: 1004,
		ErrExpiredServerKey:  1005,
		ErrMissingParams:     1006,
		ErrUnknownMethod:     1007,
		ErrUnknownParam:      1008,
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

	// ErrSenderInfoMissing is returned if partly or all of the sender's information
	// is missing.
	ErrSenderInfoMissing = errors.New("sender info incomplete")

	// ErrExpiredServerKey is returned if a sender takes too long to use the server
	// key sent during the tls handshake.
	ErrExpiredServerKey = errors.New("server key expired")

	// ErrMissingParams is returned if one or more of the expected parameters is
	// missing.
	ErrMissingParams = errors.New("missing params")

	// ErrUnknownMethod is returned if the method provided in the request is
	// not currently supported.
	ErrUnknownMethod = errors.New("unknown method")

	// ErrUnknownParam is returned if one or more of the provided parameters
	// contains unexpected data.
	ErrUnknownParam = errors.New("unknown params")
)

// GetErrorCode returns the set error code if it exists or max(uint16) if otherwise.
func GetErrorCode(err error) uint16 {
	if code, ok := serverErrorCodes[err]; ok {
		return code
	}
	return 65535 // uint16(2 ^ 16 - 1)
}
