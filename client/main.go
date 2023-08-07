// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/btcsuite/btclog"
	"github.com/dmigwi/dhamana-protocol/client/server"
)

func run(ctx context.Context, cancelFunc context.CancelFunc) {
	// initiate app shutdown
	defer cancelFunc()

	log.Infof("Loading command configurations")
	config, err := loadConfig()
	if err != nil {
		log.Errorf("loadConfig Error: %v", err)
		return
	}

	log.Infof("Initiating the log rotator")
	// Initialize the logger while while creating the data dir if it doesn't exists.
	if err := initLogRotator(config.DataDirPath, 50); err != nil {
		log.Errorf("initLogRotator Error: %v", err)
		return
	}

	level, _ := btclog.LevelFromString(config.LogLevel)
	setLogLevel(level)

	s := server.NewServer(config.Contract, config.Network)
	log.Infof("Attempting to make a connection to the sapphire network via contract %v", config.Contract)
	// Attempt to make connection.
	s.Connection()
}

// shutdown initiates the shutdown sequence.
func shutdown() {
	shutdownLog()

	log.Info("Shutdown sequence successfully completed!")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	// initiates the app logic execution in a goroutine so as to keep the main
	// goroutine waiting for events and shutdown requests
	go run(ctx, cancel)

	select {
	case <-ctx.Done():
	case <-exit:
		cancel()
	}

	// trigger the shutdown of the background processes
	shutdown()

	close(exit)
}
