// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btclog"
	"github.com/decred/dcrd/certgen"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultDataDir = ".dhamana-protocol"
)

type config struct {
	Network     string `long:"network" description:"Network to use; Supported networks: SapphireMainnet, SapphireTestnet and SapphireLocalnet" default:"SapphireTestnet" required:"required"`
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

	h := &bytes.Buffer{}
	parser.WriteHelp(h)

	if net := utils.ToNetType(conf.Network); net == utils.UnsupportedNet {
		return nil, fmt.Errorf("unsupported network used: (%v) \n %v", conf.Network, h.String())
	}

	if _, ok := btclog.LevelFromString(conf.LogLevel); !ok {
		return nil, fmt.Errorf("invalid LogLevel found: (%v) \n %v", conf.LogLevel, h.String())
	}

	if conf.DataDirPath == "" {
		return nil, fmt.Errorf("empty datadir path found \n %v", h.String())
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

// generateTLSConfig creates the TLS certificate and key files if don't exist
// or are expired.
func generateTLSConfig(datadir string) error {
	keyPath := filepath.Join(datadir, utils.TLSKeyFile)
	certPath := filepath.Join(datadir, utils.TLSCertFile)
	keyFile, _ := os.ReadFile(keyPath)
	certFile, _ := os.ReadFile(certPath)

	// Only exit the config regeneration if cert and key exit and are not expired.
	if keyFile != nil && certFile != nil {
		// exit regeneration if cert and keys are found to be valid.
		if isTLSConfigValid(certFile, keyFile) {
			return nil
		}
	}

	// certificate is valid for 6 months. Then a certificate regenerate on next
	// next restart will be necessary.
	validUntil := time.Unix(time.Now().Add(30*24*time.Hour).Unix(), 0)

	// This details don't change often thus no need to pass them as function variables.
	org := "dhamana-protocol cert"
	extraHosts := []string{"localhost", "127.0.0.1"}

	cert, key, err := certgen.NewEd25519TLSCertPair(org, validUntil, extraHosts)
	if err != nil {
		return fmt.Errorf("unable the generate tls key and cert: %v", err)
	}

	if err := os.WriteFile(keyPath, key, utils.FilePerm); err != nil {
		return fmt.Errorf("unable to store tls key file at: %s Error: %v", keyPath, err)
	}

	if err := os.WriteFile(certPath, cert, utils.FilePerm); err != nil {
		return fmt.Errorf("unable to store tls cert file at: %s Error: %v", certPath, err)
	}

	return nil
}
