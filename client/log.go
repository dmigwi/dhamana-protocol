// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btclog"
	"github.com/dmigwi/dhamana-protocol/client/server"
	"github.com/jrick/logrotate/rotator"
)

const logFile = "dhamana.log"

var (
	backendLog = btclog.NewBackend(logWriter{logFile})

	// log is a logger that is initialized with no output filters.  This
	// means the package will not perform any logging by default until the caller
	// requests it.
	log       = backendLog.Logger("MAIN")
	serverLog = backendLog.Logger("SERVER")

	// logRotator is one of the logging outputs.  It should be closed on
	// application shutdown.
	logRotators *rotator.Rotator
)

// Assigns the logger to use.
func init() {
	server.UseLogger(serverLog)
}

// logWriter implements an io.Writer that outputs to both standard output and
// the write-end pipe of an initialized log rotator.
type logWriter struct {
	loggerID string
}

// Write writes the data in p to standard out and the log rotator.
func (l logWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	if logRotators == nil {
		return 1, errors.New("log Rotator not initialised")
	}
	return logRotators.Write(p)
}

// setLogLevel the required log level.
func setLogLevel(level btclog.Level) {
	log.SetLevel(level)
	serverLog.SetLevel(level)
}

// shutdownLog safely triggers the log rotator shutdown.
func shutdownLog() {
	if logRotators != nil {
		logRotators.Close()
	}
}

// initLogRotator initializes the logging rotater to write logs to logFile and
// create roll files in the same directory.  It must be called before the
// package-global log rotater variables are used.
func initLogRotator(logDir string, maxRolls int) error {
	err := os.MkdirAll(logDir, 0o0700)
	if err != nil {
		return fmt.Errorf("failed to create log directory: %v\n", err)
	}

	r, err := rotator.New(filepath.Join(logDir, logFile), 32*1024, false, maxRolls)
	if err != nil {
		return fmt.Errorf("failed to create file rotator: %v\n", err)
	}
	logRotators = r
	return nil
}
