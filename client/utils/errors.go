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
)
