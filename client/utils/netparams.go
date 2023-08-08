// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

import (
	"fmt"
	"math/big"
)

type NetworkType int

const (
	SapphireMainnet NetworkType = iota
	SapphireTestnet
	SapphireLocalnet
	UnsupportedNet
)

type NetworkParams struct {
	Name           string
	ChainID        big.Int
	DefaultGateway string
	RuntimeID      string
}

// Networks defines the configurations mappings to the various networks supported.
var networks = map[NetworkType]NetworkParams{
	SapphireMainnet: {
		Name:           "mainnet",
		ChainID:        *big.NewInt(0x5afe),
		DefaultGateway: "https://sapphire.oasis.io",
		RuntimeID:      "0x000000000000000000000000000000000000000000000000f80306c9858e7279",
	},
	SapphireTestnet: {
		Name:           "testnet",
		ChainID:        *big.NewInt(0x5aff),
		DefaultGateway: "https://testnet.sapphire.oasis.dev",
		RuntimeID:      "0x000000000000000000000000000000000000000000000000a6d1e3ebf60dff6c",
	},
	SapphireLocalnet: {
		Name:           "localnet",
		ChainID:        *big.NewInt(0x5afd),
		DefaultGateway: "http://localhost:8545",
		RuntimeID:      "0x8000000000000000000000000000000000000000000000000000000000000000",
	},
}

// GetNetworkConfig returns the configured sapphire network params.
func GetNetworkConfig(net NetworkType) (*NetworkParams, error) {
	params, ok := networks[net]
	if !ok {
		return nil, fmt.Errorf("could not fetch %v network", net.String())
	}
	return &params, nil
}

// String defines the default stringer for NetworkType.
func (n NetworkType) String() string {
	switch n {
	case SapphireMainnet:
		return "SapphireMainnet"
	case SapphireTestnet:
		return "SapphireTestnet"
	case SapphireLocalnet:
		return "SapphireLocalnet"
	default:
		return "UnsupportedNet"
	}
}

// ToNetType maps a network type to the user network type input either in
// in camel case or snake case.
func ToNetType(net string) NetworkType {
	switch net {
	case "SapphireMainnet", "sapphire_mainnet":
		return SapphireMainnet
	case "SapphireTestnet", "sapphire_testnet":
		return SapphireTestnet
	case "SapphireLocalnet", "sapphire_localnet":
		return SapphireLocalnet
	default:
		return UnsupportedNet
	}
}
