// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/btcsuite/btclog"
	// "github.com/dmigwi/dhamana-protocol/client/utils"
)

func run(ctx context.Context, cancelFunc context.CancelFunc) {
	// initiate app shutdown
	defer cancelFunc()

	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Initialize the logger while while creating the data dir if it doesn't exists.
	if err := initLogRotator(config.DataDirPath, 50); err != nil {
		fmt.Println(err)
		return
	}

	level, _ := btclog.LevelFromString(config.LogLevel)
	setLogLevel(level)

	// net := utils.ToNetType(config.Network)
}

// shutdown initiates the shutdown sequence.
func shutdown() {
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
