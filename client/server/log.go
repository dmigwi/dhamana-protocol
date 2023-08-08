// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import "github.com/btcsuite/btclog"

var log btclog.Logger

// UseLogger sets the subsystem logs to use the provided loggers.
func UseLogger(sLogger btclog.Logger) {
	log = sLogger
}
