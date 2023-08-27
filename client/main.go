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
		log.Errorf("loadConfig error: %v", err)
		return
	}

	log.Infof("Using data directory=%s", config.DataDirPath)

	// Initialize the logger while creating the data dir if it doesn't exists.
	if err := initLogRotator(config.DataDirPath, 50); err != nil {
		log.Errorf("initLogRotator error: %v", err)
		return
	}

	level, _ := btclog.LevelFromString(config.LogLevel)
	setLogLevel(level)

	s, err := server.NewServer(ctx, config.DbPort, config.TLSCertFile,
		config.TLSKeyFile, config.DataDirPath, config.Network, config.ServerURL,
		config.DbHost, config.DbName, config.DbUser, config.DbPassword)
	if err != nil {
		log.Errorf("Server Config error: %v", err)
		return
	}

	// Initiate the data syncer
	if err = s.SyncData(); err != nil {
		log.Errorf("SyncData failed error: %v", err)
		return
	}

	// Run the server
	if err = s.Run(); err != nil {
		log.Errorf("Server failed error: %v", err)
		return
	}
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
