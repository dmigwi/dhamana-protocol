// Copyright (c) 2023 Migwi Ndung'u
// Copyright (c) 2023 Sapphire-Paratime Authors
// See LICENSE for details.

package sapphire

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	DefaultGasPrice = 100_000_000_000
	// DefaultGasLimit is set on all transactions without explicit gas limit to avoid being set on signed queries by the web3 gateway.
	DefaultGasLimit   = 30_000_000
	DefaultBlockRange = 15
)

type WrappedBackend struct {
	bind.ContractBackend
	chainID    big.Int
	cipher     Cipher
	signerFunc SignerFn
	ctx        context.Context

	privateKey []byte
}

// Confirm that WrappedBacked implements the bind.ContractBackend interface.
var _ bind.ContractBackend = (*WrappedBackend)(nil)

// SignerFn is a function that produces secp256k1 signatures in RSV format.
type SignerFn = func(digest [32]byte, privateKey []byte) ([]byte, error)

// NewCipher creates a default cipher.
func NewCipher(ctx context.Context, net utils.NetworkType) (Cipher, error) {
	runtimePublicKey, err := getRuntimePublicKey(ctx, net)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch runtime callata public key: %w", err)
	}

	keypair, err := NewCurve25519KeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral keypair: %w", err)
	}

	// convert the slice to an array
	var pubkey [32]byte
	copy(pubkey[:], runtimePublicKey)

	cipher, err := NewX25519DeoxysIICipher(*keypair, pubkey)
	if err != nil {
		return nil, fmt.Errorf("failed to create default cipher: %w", err)
	}
	return cipher, nil
}

// WrapClient wraps an ethclient.Client so that it can talk to Sapphire.
func WrapClient(ctx context.Context, c ethclient.Client, net utils.NetworkType, sign SignerFn) (*WrappedBackend, error) {
	network, err := utils.GetNetworkConfig(net)
	if err != nil {
		return nil, err
	}

	cipher, err := NewCipher(ctx, net)
	if err != nil {
		return nil, err
	}

	return &WrappedBackend{
		ContractBackend: &c,
		chainID:         network.ChainID,
		cipher:          cipher,
		signerFunc:      sign,
	}, nil
}

func (b *WrappedBackend) SetClientSigningKey(privateKey []byte) {
	// Set the private key
	b.privateKey = privateKey
}

// Transactor returns a TransactOpts that can be used with Sapphire.
func (b *WrappedBackend) Transactor(from common.Address) *bind.TransactOpts {
	signer := types.LatestSignerForChainID(&b.chainID)
	signFn := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if addr != from {
			return nil, bind.ErrNotAuthorized
		}

		packedTx := types.NewTx(&types.LegacyTx{
			Nonce:    tx.Nonce(),
			GasPrice: tx.GasPrice(),
			Gas:      tx.Gas(),
			To:       tx.To(),
			Value:    tx.Value(),
			Data:     b.cipher.EncryptEncode(tx.Data()),
		})

		var signedTxBytes [32]byte
		copy(signedTxBytes[:], signer.Hash(packedTx).Bytes())

		sig, err := b.signerFunc(signedTxBytes, b.privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to sign tx: %w", err)
		}

		return packedTx.WithSignature(signer, sig)
	}

	return &bind.TransactOpts{
		From:     from,
		Signer:   signFn,
		GasPrice: big.NewInt(DefaultGasPrice),
		GasLimit: DefaultGasLimit,
	}
}

// CallContract executes a Sapphire paratime contract call with the specified
// data as the input. CallContract implements ContractCaller.
func (b *WrappedBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var err error
	var packedCall *ethereum.CallMsg

	if call.From == [common.AddressLength]byte{} {
		// prepares call.Data for being sent to Sapphire. The call will be
		// end-to-end encrypted, but the `from` address will be zero.
		packedCall.Data = b.cipher.EncryptEncode(call.Data)
	} else {
		leashBlockNumber := big.NewInt(0)
		if blockNumber != nil {
			leashBlockNumber.Sub(blockNumber, big.NewInt(1))
		} else {
			latestHeader, err := b.ContractBackend.HeaderByNumber(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch latest block number: %w", err)
			}
			leashBlockNumber.Sub(latestHeader.Number, big.NewInt(1))
		}

		header, err := b.ContractBackend.HeaderByNumber(ctx, leashBlockNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch leash block header: %w", err)
		}

		blockHash := header.Hash()
		leash := NewLeash(header.Nonce.Uint64(), header.Number.Uint64(), blockHash[:], DefaultBlockRange)

		// prepares call.Data in-place for being sent to Sapphire. The call will be
		// end-to-end encrypted and a signature will be used to authenticate the `from` address.
		dataPack, err := NewDataPack(b.signerFunc, b.privateKey, b.chainID.Uint64(), call.From[:],
			call.To[:], DefaultGasLimit, call.GasPrice, call.Value, call.Data, leash)
		if err != nil {
			return nil, fmt.Errorf("failed to create signed call data back: %w", err)
		}
		packedCall.Data = dataPack.EncryptEncode(b.cipher)
	}

	res, err := b.ContractBackend.CallContract(ctx, *packedCall, blockNumber)
	if err != nil {
		return nil, err
	}
	return b.cipher.DecryptEncoded(res)
}

// EstimateGas implements ContractTransactor.
func (b *WrappedBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return DefaultGasLimit, nil
}
