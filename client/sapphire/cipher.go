// Copyright (c) 2023 Migwi Ndung'u
// Copyright (c) 2023 Sapphire-Paratime Authors
// See LICENSE for details.

package sapphire

import (
	"bytes"
	"context"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/oasisprotocol/deoxysii"
	"github.com/oasisprotocol/oasis-core/go/common/cbor"
	mraeApi "github.com/oasisprotocol/oasis-core/go/common/crypto/mrae/api"
	mrae "github.com/oasisprotocol/oasis-core/go/common/crypto/mrae/deoxysii"
	"golang.org/x/crypto/curve25519"
)

type Kind uint64

const (
	Plain          = iota
	X25519DeoxysII = 1
)

var (
	ErrCallFailed       = errors.New("call failed in module")
	ErrCallResultDecode = errors.New("could not decode call result")
)

type CallResult struct {
	Fail    *Failure      `json:"failure,omitempty"`
	OK      []byte        `json:"ok,omitempty"`
	Unknown *AeadEnvelope `json:"unknown,omitempty"`
}

type Inner struct {
	Fail *Failure `json:"fail"`
	OK   []byte   `json:"ok"`
}

type Failure struct {
	Module  string `json:"module"`
	Code    uint64 `json:"code"`
	Message string `json:"message,omitempty"`
}

type AeadEnvelope struct {
	Nonce []byte `json:"nonce"`
	Data  []byte `json:"data"`
}

type Request struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Response struct {
	Error  *Error          `json:"error"`
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
}

// CallDataPublicKey is the public key alongside the key manager's signature.
type CallDataPublicKey struct {
	// PublicKey is the requested public key.
	PublicKey hexutil.Bytes `json:"key"`
	// Checksum is the checksum of the key manager state.
	Checksum hexutil.Bytes `json:"checksum"`
	// Signature is the Sign(sk, (key || checksum)) from the key manager.
	Signature hexutil.Bytes `json:"signature"`
}

type Cipher interface {
	Kind() uint64
	Encrypt(plaintext []byte) (ciphertext []byte, nonce []byte)
	Decrypt(nonce []byte, ciphertext []byte) (plaintext []byte, err error)
	EncryptEncode(plaintext []byte) []byte
	DecryptEncoded(result []byte) ([]byte, error)
	DecryptCallResult(result []byte) ([]byte, error)
}

// X25519DeoxysIICipher is the default cipher that does what it says on the tin.
type X25519DeoxysIICipher struct {
	cipher  cipher.AEAD
	keypair Curve25519KeyPair
}

type Curve25519KeyPair struct {
	PublicKey [curve25519.PointSize]byte
	SecretKey [curve25519.ScalarSize]byte
}

// NewCurve25519KeyPair generates a random keypair suitable for use with the X25519DeoxysII cipher.
func NewCurve25519KeyPair() (*Curve25519KeyPair, error) {
	public, private, err := mraeApi.GenerateKeyPair(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Curve25519KeyPair{
		PublicKey: *public,
		SecretKey: *private,
	}, nil
}

func NewX25519DeoxysIICipher(keypair Curve25519KeyPair, peerPublicKey [curve25519.PointSize]byte) (*X25519DeoxysIICipher, error) {
	var sharedKey [deoxysii.KeySize]byte
	mrae.Box.DeriveSymmetricKey(sharedKey[:], &peerPublicKey, &keypair.SecretKey)

	cipher, err := deoxysii.New(sharedKey[:])
	// Set the slice to zero before exiting.
	mraeApi.Bzero(sharedKey[:])
	if err != nil {
		return nil, err
	}

	return &X25519DeoxysIICipher{
		cipher:  cipher,
		keypair: keypair,
	}, nil
}

func (c X25519DeoxysIICipher) Kind() uint64 {
	return X25519DeoxysII
}

func (c X25519DeoxysIICipher) Encrypt(plaintext []byte) (ciphertext []byte, nonce []byte) {
	nonce = make([]byte, deoxysii.NonceSize)
	if _, err := rand.Reader.Read(nonce); err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: %v", err))
	}

	res := c.cipher.Seal(ciphertext, nonce, plaintext, []byte{})
	return res, nonce
}

func (c X25519DeoxysIICipher) Decrypt(nonce []byte, ciphertext []byte) ([]byte, error) {
	meta := make([]byte, 0)
	return c.cipher.Open(ciphertext[:0], nonce, ciphertext, meta)
}

func (c X25519DeoxysIICipher) encryptCallData(plaintext []byte) (ciphertext []byte, nonce []byte) {
	return c.Encrypt(cbor.Marshal(Data{
		Body: plaintext,
	}))
}

func (c X25519DeoxysIICipher) EncryptEnvelope(plaintext []byte) *EncryptedBodyEnvelope {
	// Txs without data are just balance transfers, and all data in those is public.
	if len(plaintext) == 0 {
		return nil
	}

	data, nonce := c.encryptCallData(plaintext)
	return &EncryptedBodyEnvelope{
		Body: Body{
			Nonce: nonce,
			Data:  data,
			PK:    c.keypair.PublicKey[:],
		},
		Format: c.Kind(),
	}
}

func (c X25519DeoxysIICipher) EncryptEncode(plaintext []byte) []byte {
	envelope := c.EncryptEnvelope(plaintext)
	return cbor.Marshal(envelope)
}

func (c X25519DeoxysIICipher) DecryptCallResult(response []byte) ([]byte, error) {
	var callResult CallResult
	cbor.MustUnmarshal(response, &callResult)

	// TODO: actually decode and return failure
	if callResult.Fail != nil {
		return nil, ErrCallFailed
	}

	var aeadEnvelope AeadEnvelope
	if callResult.OK != nil {
		if err := cbor.Unmarshal(callResult.OK, &aeadEnvelope); err != nil {
			return callResult.OK, nil
		}
	} else if callResult.Unknown != nil {
		aeadEnvelope = *callResult.Unknown
	} else {
		return nil, ErrCallResultDecode
	}

	decrypted, err := c.Decrypt(aeadEnvelope.Nonce, aeadEnvelope.Data)
	if err != nil {
		return nil, err
	}

	var innerResult Inner
	cbor.MustUnmarshal(decrypted, &innerResult)

	if innerResult.OK != nil {
		return innerResult.OK, nil
	}

	if innerResult.Fail != nil {
		msg := innerResult.Fail.Message
		if len(msg) == 0 {
			msg = fmt.Sprintf("Call failed in module %s with code %d", innerResult.Fail.Module, innerResult.Fail.Code)
		}
		return nil, errors.New(msg)
	}

	return nil, errors.New("Unexpected inner call result:" + string(callResult.Unknown.Data))
}

func (c X25519DeoxysIICipher) DecryptEncoded(response []byte) ([]byte, error) {
	return c.DecryptCallResult(response)
}

// getRuntimePublicKey fetches the runtime calldata public key from the default Sapphire gateway.
func getRuntimePublicKey(ctx context.Context, net utils.NetworkType) ([]byte, error) {
	network, err := utils.GetNetworkConfig(net)
	if err != nil {
		return nil, err
	}

	request := Request{
		Version: "2.0",
		Method:  "oasis_callDataPublicKey",
		ID:      1,
	}
	rawReq, _ := json.Marshal(request)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, network.DefaultGateway, bytes.NewBuffer(rawReq))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for runtime calldata public key: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request runtime calldata public key: %w", err)
	}

	decoder := json.NewDecoder(res.Body)
	rpcRes := new(Response)
	if err := decoder.Decode(&rpcRes); err != nil {
		return nil, fmt.Errorf("unexpected response to request for runtime calldata public key: %w", err)
	}
	res.Body.Close()

	var pubKey CallDataPublicKey
	if err := json.Unmarshal(rpcRes.Result, &pubKey); err != nil {
		return nil, fmt.Errorf("invalid response when fetching runtime calldata public key: %w", err)
	}
	if len(pubKey.PublicKey) != 32 {
		return nil, fmt.Errorf("invalid public key length")
	}

	return pubKey.PublicKey, nil
}
