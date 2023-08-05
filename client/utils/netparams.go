// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

type NetworkType int

const (
	SapphireMainnet NetworkType = iota
	SapphireTestnet
	SapphireLocalnet
	UnsupportedNet
)

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
