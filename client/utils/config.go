// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	// WelcomeText is an easter egg placed in the code for cryptography enthusiasts
	// to attempt solving.
	// Should you be successful in decrypting it, please do what it says!
	WelcomeText = "Pi encrypted  715b48145c501951595355505d194050475a135a5d4251595d5d" +
		"40021301705647585b5d600c7d527f615760417f63477d5a696045610348437f470d790757"

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

// GeneratePrivKey generates a private key using a P256 curve. P256 is used
// because it provides keys of size 32 bit which are the maximum allowed by AES.
func GeneratePrivKey(randGen io.Reader) (PrivateKey, error) {
	if randGen == nil {
		randGen = cryptorand.Reader
	}

	key, err := ecdh.P256().GenerateKey(randGen)

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

	pubkey, err := ecdh.P256().NewPublicKey(rawPubkey)
	if err != nil {
		return nil, fmt.Errorf("invalid public key found: %v", err)
	}

	return p.ECDH(pubkey)
}

// EncryptAES implements the AES algorithm in encrypting data.
// To encrypt data, we need the following:
// 1. sharedkey (32-bytes for AES-256 encryption)
// 2. nonce (a random no., only used only once)
// 3. plaintext data to encrypt

// Its outputs a hexutils encode string and an error
func EncryptAES(sharedKey []byte, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(sharedKey)
	if err != nil {
		return "", errors.New("unable to generate new aes cipher")
	}

	// Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("gcm or Galois/Counter Mode creation failed")
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	// Encrypt the data using aesGCM.Seal. Since we don't want to save the nonce
	// somewhere else in this case, we add it as a prefix to the encrypted data.
	// The first nonce argument in Seal is the prefix.
	return hex.EncodeToString(aesGCM.Seal(nonce, nonce, plaintext, nil)), nil
}

// DecryptAES implements the AES algorithm in decoding the provided ciphertext.
// To decrypt the ciphertext, we need the following:
// 1. sharedkey (32-bytes for AES-256 encryption)
// 2. ciphertext data to decrypt.

// Decode the ciphertext first into a bytes array.
// The decoded output returned is in bytes form
func DecryptAES(sharedKey []byte, ciphertext string) ([]byte, error) {
	txt, err := hex.DecodeString(ciphertext)
	if err != nil {
		return nil, errors.New("unable to decode the hex string")
	}

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(sharedKey)
	if err != nil {
		return nil, errors.New("unable to generate new aes cipher")
	}

	// Create a new GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("gcm or Galois/Counter Mode creation failed")
	}

	// Get the nonce size
	nonceSize := gcm.NonceSize()
	if len(txt) < nonceSize {
		return nil, errors.New("invalid decoded cipher text size")
	}

	// Extract the nonce from the encrypted data
	nonce, txt := txt[:nonceSize], txt[nonceSize:]

	// Decrypt the data
	return gcm.Open(nil, nonce, txt, nil)
}
