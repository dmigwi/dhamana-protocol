// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	// WelcomeText is an easter egg placed in the code for cryptography enthusiasts
	// to attempt solving.
	// Should you be successful in decrypting it, please do what it says!
	WelcomeText = "Pi encrypted  715b48145c501951595355505d194050475a135a5d4251595d5d40021301705647585b5d600c7d527f615760417f63477d5a696045610348437f470d790757"

	// FullDateformat defines the full date format supported
	FullDateformat = "Mon 15:04:05 2006-01-02"

	// FilePerm defines the file permission used to manage the application data.
	FilePerm = os.FileMode(0o0700)

	// JSONRPCVersion defines the JSON version supportted for all the backend requests
	// recieved by the server.
	JSONRPCVersion = "2.0"
)

// PrivateKey is generated using elliptic curve diffie-hellman algorithm. This
// is used to share sensitive information between the server and the client i.e
// signer key.
type PrivateKey struct {
	*ecdh.PrivateKey
}

// GeneratePrivKey() generates a private key using P521 curve.
func GeneratePrivKey() (PrivateKey, error) {
	key, err := ecdh.P521().GenerateKey(rand.Reader)

	return PrivateKey{PrivateKey: key}, err
}

// PubKeyToHexString converts the public key associated with the provided private
// key to a hex text.
func (p *PrivateKey) PubKeyToHexString() string {
	return hexutil.Encode(p.PublicKey().Bytes())
}

// ComputeSharedKey computes the shared key between the remote and the local instances.
func (p *PrivateKey) ComputeSharedKey(remotePubKey string) ([]byte, error) {
	rawPubkey, err := hexutil.Decode(remotePubKey)
	if err != nil {
		return nil, errors.New("unable to decoded public key hex")
	}

	pubkey, err := ecdh.P521().NewPublicKey(rawPubkey)
	if err != nil {
		return nil, errors.New("invalid public key found")
	}

	return p.ECDH(pubkey)
}
