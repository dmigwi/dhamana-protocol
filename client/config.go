// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btclog"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	flags "github.com/jessevdk/go-flags"
)

const (
	// defaultDataDir sets the default data directory name appended on the
	// user config file path based on the os in use.
	defaultDataDir = "dhamana-protocol"
)

type config struct {
	Network     string `long:"network" description:"Network to use; Supported networks: SapphireMainnet, SapphireTestnet and SapphireLocalnet" default:"SapphireTestnet" required:"required"`
	DataDirPath string `long:"datadir" description:"Directory path to where the app data is stored"`
	LogLevel    string `long:"loglevel" description:"Logging level {trace, debug, info, warn, error, critical, off}" default:"info"`
	TLSCertFile string `long:"certfile" description:"tls certificate file name" default:"server.crt"`
	TLSKeyFile  string `long:"keyfile" description:"tls key file name" default:"server.key"`
	ServerURL   string `long:"url" description:"Server url to server content using" default:"https://0.0.0.0:30443"`

	// DB configuration
	DbPort     uint16 `long:"db_port" description:"Port to use when connecting to the db" default:"5432"`
	DbDriver   string `long:"db_driver" description:"Driver name of the database to use {postgres,sqlite}" default:"postgres"`
	DbHost     string `long:"db_host" description:"Host to use in connecting to the db" default:"localhost"`
	DbUser     string `long:"db_user" description:"Username to use in connecting to the db" default:"ana"`
	DbPassword string `long:"db_password" description:"Password of the database username to use" default:"ana"`
	DbName     string `long:"db_name" description:"Name of the database to connect to" default:"dhamana"`
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

	h := &bytes.Buffer{}
	parser.WriteHelp(h)

	if net := utils.ToNetType(conf.Network); net == utils.UnsupportedNet {
		return nil, fmt.Errorf("unsupported network used: %q \n %s", conf.Network, h.String())
	}

	if _, ok := btclog.LevelFromString(conf.LogLevel); !ok {
		return nil, fmt.Errorf("invalid LogLevel found: %q \n %s", conf.LogLevel, h.String())
	}

	if conf.DataDirPath == "" {
		return nil, fmt.Errorf("empty datadir path found \n %s", h.String())
	}

	// validateTLSCerts confirms TLS certificates exists and are valid.
	if err := validateTLSCerts(&conf); err != nil {
		return nil, fmt.Errorf("validateTLSCerts error: %v \n %s", err, h.String())
	}

	if _, err := url.Parse(conf.ServerURL); err != nil {
		return nil, fmt.Errorf("invalid server url found: %q \n %s", conf.ServerURL, h.String())
	}

	// confirm all the db configurations have non-zero values
	if !isDbConfig(&conf) {
		return nil, fmt.Errorf("missing some db configurations \n %s", h.String())
	}

	return &conf, nil
}

// isTLSConfigValid confirms if the cert and key contents are valid.
func isTLSConfigValid(cert, key []byte) bool {
	pemKey, _ := pem.Decode(key)
	pemCert, _ := pem.Decode(cert)
	// Ensure the PEM-encoded cert and key that is returned can be decoded.
	// otherwise return false.
	if pemCert == nil || pemKey == nil {
		return false
	}

	// Ensure the DER-encoded key bytes can be successfully parsed.
	_, err := x509.ParseECPrivateKey(pemKey.Bytes)
	if err != nil {
		return false
	}

	// Ensure the DER-encoded cert bytes can be successfully into an X.509
	// certificate.
	x509Cert, err := x509.ParseCertificate(pemCert.Bytes)
	if err != nil {
		return false
	}

	// Confirm the certificate is still with the validity period.
	return x509Cert.NotAfter.After(time.Now())
}

// validateTLSCerts confirm the TLS certificates files exist and can be decoded.
func validateTLSCerts(conf *config) error {
	keyPath := filepath.Join(conf.DataDirPath, conf.TLSKeyFile)
	certPath := filepath.Join(conf.DataDirPath, conf.TLSCertFile)
	keyFile, _ := os.ReadFile(keyPath)
	certFile, _ := os.ReadFile(certPath)

	// Exit if cert and key don't exist.
	if keyFile == nil || certFile == nil {
		return fmt.Errorf("missing %s or %s files at datadir %s",
			conf.TLSCertFile, conf.TLSKeyFile, conf.DataDirPath)
	}
	// Confirm if cert and keys are found to be valid.
	if isTLSConfigValid(certFile, keyFile) {
		return errors.New("unable to decode cert and key files. Regenerate them")
	}

	return nil
}

// isDbConfig confirms that the provided db config is valid.
func isDbConfig(conf *config) bool {
	if conf.DbPort == 0 {
		return false
	}

	// confirm the db driver is among the supported ones.
	if _, ok := utils.SupportedDbDrivers[conf.DbDriver]; !ok {
		return false
	}

	if conf.DbHost == "" || conf.DbName == "" ||
		conf.DbPassword == "" || conf.DbUser == "" {
		return false
	}
	return true
}
