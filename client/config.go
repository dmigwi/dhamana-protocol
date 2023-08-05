// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btclog"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultDataDir = ".dhamana-protocol"
)

type config struct {
	Network     string `long:"network" description:"Network to use; Supported networks: SapphireMainnet, SapphireTestnet and SapphireLocalnet" default:"SapphireTestnet"`
	DataDirPath string `long:"datadir" description:"Directory path to where the app data is stored"`
	LogLevel    string `long:"loglevel" description:"Logging level {trace, debug, info, warn, error, critical, off}" default:"info"`
}

// defaultDataDir returns the default
func defaultDataDirPath() (string, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("unable to fetch default config path: %v", err)
	}

	return filepath.Join(configPath, defaultDataDir), nil
}

func loadConfig() (*config, error) {
	defaultDataDirPath, err := defaultDataDirPath()
	if err != nil {
		return nil, err
	}

	conf := config{
		DataDirPath: defaultDataDirPath,
	}

	parser := flags.NewParser(&conf, flags.Default)
	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	if net := utils.ToNetType(conf.Network); net == utils.UnsupportedNet {
		return nil, fmt.Errorf("unsupported network used: %v", parser.Usage)
	}

	if _, ok := btclog.LevelFromString(conf.LogLevel); !ok {
		return nil, fmt.Errorf("invalid LogLevel found: %v", parser.Usage)
	}

	if conf.DataDirPath == "" {
		return nil, fmt.Errorf("empty datadir path found: %v", parser.Usage)
	}

	return &conf, nil
}
