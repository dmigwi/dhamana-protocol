package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/sapphire"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	serverConf       *ServerConfig
	sampleHexAddress = common.HexToAddress("0x3396FD816Dd81100477c8ea3853039822f36B7ed")
	sampleSigningKey = "8520f0098930a754748b7ddcb43ef75a0dbf3a0d26381af4eba4a98eaa9b4e6a"

	sampleHexAddress1 = common.HexToAddress("0x3396FD816Dd81100477c8ea3853039822f36B7ad")
	sampleHexAddress2 = common.HexToAddress("0x3396FD816Dd81100477c8ea3853039822f36B7bd")
	sampleHexAddress3 = common.HexToAddress("0x3396FD816Dd81100477c8ea3853039822f36B71d")

	pubkey1 = "0x5a59ec90d4cc9f9a60ac42179dcf461b3f3c7caa3b4960b69ce2de7ead3a5417"
	pubkey2 = "0x409f538d156b3ffcec0b457f6f193946ab19e38fd36bd393963102473e583686"
	pubkey3 = "0x673ce33fd92d3ea971412dbffe7c5523dc17174b7a77b4681013adbd0bea4afb"
)

type input struct {
	testName   string
	method     string
	body       interface{}
	needSigner bool
}

type output struct {
	errCode    uint16
	shortErr   error
	longErr    string
	methodType utils.MethodType
}

// Set up the server config before initiating tests.
func TestMain(m *testing.M) {
	var err error
	ctx, cancelFn := context.WithCancel(context.Background())
	serverConf, err = genMockWrapper(ctx)

	if err == nil {
		// Store expired keys
		expiredKey := serverKeyResp{
			Pubkey: pubkey1,
			Expiry: uint64(time.Now().UTC().Unix()),
		}
		serverConf.sessionKeys.Store(sampleHexAddress1, expiredKey)

		// store fresh keys with an expiry of 2 minutes
		freshKey := serverKeyResp{
			Pubkey: pubkey2,
			Expiry: uint64(time.Now().UTC().Add(2 * time.Minute).Unix()),
		}
		serverConf.sessionKeys.Store(sampleHexAddress2, freshKey)
		m.Run()
	} else {
		log.Error("unexpected error: ", err)
	}
	cancelFn()
}

// TestDecodeRequestBody runs tests on function decodeRequestBody.
func TestDecodeRequestBody(t *testing.T) {
	testdata := []struct {
		data input
		val  output
	}{
		{
			data: input{
				testName:   "Test-http-method-support",
				method:     http.MethodGet,
				needSigner: false,
				body: rpcMessage{
					Method: "getServerPubKey",
				},
			},
			val: output{
				errCode:    1001,
				shortErr:   utils.ErrInvalidReq,
				longErr:    "invalid http method GET found expected POST",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-decode-invalid-json-body",
				method:     http.MethodPost,
				needSigner: false,
				body:       "",
			},
			val: output{
				errCode:    1000,
				shortErr:   utils.ErrInvalidJSON,
				longErr:    "EOF",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-json-with-unsupported-rpc-version",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					Version: "1.0",
					Method:  "getServerPubKey",
				},
			},
			val: output{
				errCode:    1001,
				shortErr:   utils.ErrInvalidReq,
				longErr:    "expected JSON-RPC version 2.0 but found 1.0",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-missing-method",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					Version: "2.0",
					Method:  "",
				},
			},
			val: output{
				errCode:    1003,
				shortErr:   utils.ErrMethodMissing,
				longErr:    "expected a method to be provided",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-missing-sender-address",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					Version: "2.0",
					Method:  "getServerPubKey",
					Sender: &senderInfo{
						Address: common.Address{},
					},
				},
			},
			val: output{
				errCode:    1004,
				shortErr:   utils.ErrSenderAddrMissing,
				longErr:    "expected sender address to be provided",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-for-required-signer-key",
				method:     http.MethodPost,
				needSigner: true,
				body: rpcMessage{
					Version: "2.0",
					Method:  "createBond",
					Sender: &senderInfo{
						Address:    sampleHexAddress,
						SigningKey: "",
					},
				},
			},
			val: output{
				errCode:    1006,
				shortErr:   utils.ErrSignerKeyMissing,
				longErr:    "expected sender signer key to be provided",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-for-supportted-methods",
				method:     http.MethodPost,
				needSigner: true,
				body: rpcMessage{
					Version: "2.0",
					Method:  "createBondAndSign",
					Sender: &senderInfo{
						Address:    sampleHexAddress,
						SigningKey: sampleSigningKey,
					},
				},
			},
			val: output{
				errCode:    1008,
				shortErr:   utils.ErrUnknownMethod,
				longErr:    "method createBondAndSign not supportted",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-for-params-count-mismatch",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					Version: "2.0",
					Method:  "getServerPubKey",
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{"client-pub-key", "unsupported-param"},
				},
			},
			val: output{
				errCode:    1007,
				shortErr:   utils.ErrMissingParams,
				longErr:    "method getServerPubKey requires 1 params found 2 params",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-for-params-type-mismatch",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					Version: "2.0",
					Method:  "getServerPubKey",
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{200},
				},
			},
			val: output{
				errCode:    1009,
				shortErr:   utils.ErrUnknownParam,
				longErr:    "expected param 200 to be of type string found it to be int",
				methodType: utils.UnknownType,
			},
		},
		{
			data: input{
				testName:   "Test-for-successful-serverkey-method",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					Version: "2.0",
					Method:  "getServerPubKey",
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
			val: output{
				errCode:    0,
				shortErr:   nil,
				longErr:    "",
				methodType: utils.ServerKeyType,
			},
		},
		{
			data: input{
				testName:   "Test-for-successful-serverkey-method",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  "getServerPubKey",
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
			val: output{
				errCode:    0,
				shortErr:   nil,
				longErr:    "",
				methodType: utils.ServerKeyType,
			},
		},
		{
			data: input{
				testName:   "Test-for-successful-contact-method",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					ID:      21,
					Version: "2.0",
					Method:  "signBondStatus",
					Sender: &senderInfo{
						Address:    sampleHexAddress,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
			val: output{
				errCode:    0,
				shortErr:   nil,
				longErr:    "",
				methodType: utils.ContractType,
			},
		},
		{
			data: input{
				testName:   "Test-for-successful-local-method",
				method:     http.MethodPost,
				needSigner: false,
				body: rpcMessage{
					ID:      21,
					Version: "2.0",
					Method:  "getBondsByStatus",
					Sender: &senderInfo{
						Address:    sampleHexAddress,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{0},
				},
			},
			val: output{
				errCode:    0,
				shortErr:   nil,
				longErr:    "",
				methodType: utils.ContractType,
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.data.testName, func(t *testing.T) {
			msg := rpcMessage{
				ID:      12,
				Version: utils.JSONRPCVersion,
			}

			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(v.data.body) // error ignored since its not being tested.

			req := httptest.NewRequest(v.data.method, "/random-path", &buf)

			// The same struct passed when making is used when retrieving a response.
			retType := decodeRequestBody(req, &msg, v.data.needSigner)

			if retType != v.val.methodType {
				t.Fatalf("expected returned method type to be %q but found %q",
					retType, v.val.methodType)
			}

			// Test for the packed response.

			if msg.Sender != nil {
				t.Fatal("expected sender field to be nil")
			}

			if msg.Method != "" {
				t.Fatal("expected method field to be empty")
			}

			if msg.Params != nil {
				t.Fatal("expected params field to be nil")
			}

			if retType == utils.UnknownType {
				if msg.Error == nil {
					t.Fatal("expected error field not to be nil")
				}

				if msg.Result != nil {
					t.Fatal("expected result field to be nil")
				}
			} else {
				if msg.Error != nil {
					t.Fatal("expected error field to be nil")
				}

				// Prevent further execution since no error is expected
				return
			}

			if msg.Error.Code != v.val.errCode {
				t.Fatalf("expected returned error code to be %q but found %q",
					msg.Error.Code, v.val.errCode)
			}

			err, _ := msg.Error.Data.(error)
			if err != v.val.shortErr {
				t.Fatalf("expected returned short error to be %q but found %q",
					err, v.val.shortErr)
			}

			if msg.Error.Message != v.val.longErr {
				t.Fatalf("expected returned long error to be %q but found %q",
					err, v.val.shortErr)
			}
		})
	}
}

// Create a mock wrapper for use with sapphire backend

type mockWrapper struct{}

func (m *mockWrapper) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}

func (m *mockWrapper) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}

func (m *mockWrapper) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return nil, nil
}

func (m *mockWrapper) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return nil, nil
}

func (m *mockWrapper) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, nil
}

func (m *mockWrapper) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(42), nil
}

func (m *mockWrapper) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(42), nil
}

func (m *mockWrapper) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return 0, nil
}

func (m *mockWrapper) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}

func (m *mockWrapper) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return nil, nil
}

func (m *mockWrapper) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, nil
}

func genMockWrapper(ctx context.Context) (*ServerConfig, error) {
	conn := &mockWrapper{}
	backend, err := sapphire.WrapClient(ctx, conn, utils.SapphireLocalnet,
		func(digest [32]byte, _ []byte) ([]byte, error) { return digest[:], nil })

	// Prevent actual sending of txs
	backend.IsTesting = true

	// Create the chat instance to be used.
	chatInstance, err := contracts.NewChat(sampleHexAddress, backend)
	if err != nil {
		return nil, err
	}

	return &ServerConfig{
		// ctx:          ctx,
		// network:      net,
		// contractAddr: sampleHexAddress,
		// serverURL:    serverURL,
		// datadir:      datadir,
		// tlsCertFile:  certfile,
		// tlsKeyFile:   keyfile,

		backend:  backend,
		bondChat: chatInstance,
	}, nil
}

// TestServerPubkey tests unique functionality implemented in serverPubkey method.
func TestServerPubkey(t *testing.T) {
	testdata := []struct {
		data input
		val  output
	}{
		{
			data: input{
				testName: "Test-for-access-to-non-serverkey-method",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  "getBondsByStatus",
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
			val: output{
				errCode:  1008,
				shortErr: utils.ErrUnknownMethod,
				longErr:  "unsupported method getBondsByStatus found for this route",
			},
		},
		{
			data: input{
				testName: "Test-for-successful-access-to-serverkey-method",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  "getServerPubKey",
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.data.testName, func(t *testing.T) {
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(v.data.body) // error ignored since its not being tested.

			responseWritter := httptest.NewRecorder()

			serverConf.serverPubkey(responseWritter,
				httptest.NewRequest(v.data.method, "/serverpubkey", &buf))

			data, err := io.ReadAll(responseWritter.Body)
			if err != nil {
				t.Fatalf("expected no error but found %q", err)
			}

			msg := rpcMessage{}
			_ = json.Unmarshal(data, &msg)

			if msg.Error == nil && v.val.shortErr != nil {
				t.Fatalf("expected method %q to return an error but found",
					v.val.methodType)
			} else {

				var result serverKeyResp
				_ = json.Unmarshal(msg.Result, &result)

				if result.Pubkey == "" {
					t.Fatal("expected the server pubkey not to be empty")
				}

				now := time.Now().UTC()
				if time.Unix(int64(result.Expiry), 0).UTC().After(now) {
					t.Fatalf("expected the server pubkey expiry to be after %v", now)
				}

				// No error was expected, prevent further error check.
				return
			}

			if msg.Error.Code != v.val.errCode {
				t.Fatalf("expected returned error code to be %q but found %q",
					msg.Error.Code, v.val.errCode)
			}

			err, _ = msg.Error.Data.(error)
			if err != v.val.shortErr {
				t.Fatalf("expected returned short error to be %q but found %q",
					err, v.val.shortErr)
			}

			if msg.Error.Message != v.val.longErr {
				t.Fatalf("expected returned long error to be %q but found %q",
					err, v.val.shortErr)
			}
		})
	}
}

// TestBackendQueryFunc tests unique functionality implemented in backendQueryFunc method.
func TestBackendQueryFunc(t *testing.T) {
	testdata := []struct {
		data input
		val  output
	}{
		{
			data: input{
				testName: "Test-for-access-to-non-contract-method",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  "getServerPubKey",
					Sender: &senderInfo{
						Address:    sampleHexAddress2,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
			val: output{
				errCode:  1008,
				shortErr: utils.ErrUnknownMethod,
				longErr:  "unsupported method getServerPubKey found for this route",
			},
		},
		{
			data: input{
				testName: "Test-for-missing-server-keys",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  "signBondStatus",
					Sender: &senderInfo{
						Address:    sampleHexAddress3,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
			val: output{
				errCode:  1011,
				shortErr: utils.ErrMissingServerKey,
				longErr:  "no server keys found associated with the sender",
			},
		},
		{
			data: input{
				testName: "Test-for-expired-server-keys",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  "signBondStatus",
					Sender: &senderInfo{
						Address:    sampleHexAddress1,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
			val: output{
				errCode:  1005,
				shortErr: utils.ErrExpiredServerKey,
				longErr:  "",
			},
		},
		{
			data: input{
				testName: "Test-for-successful-access-to-contract-method",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  "signBondStatus",
					Sender: &senderInfo{
						Address:    sampleHexAddress2,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{"client-pub-key"},
				},
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.data.testName, func(t *testing.T) {
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(v.data.body) // error ignored since its not being tested.

			responseWritter := httptest.NewRecorder()

			serverConf.serverPubkey(responseWritter,
				httptest.NewRequest(v.data.method, "/serverpubkey", &buf))

			data, err := io.ReadAll(responseWritter.Body)
			if err != nil {
				t.Fatalf("expected no error but found %q", err)
			}

			msg := rpcMessage{}
			_ = json.Unmarshal(data, &msg)

			if msg.Error == nil && v.val.shortErr != nil {
				t.Fatalf("expected method %q to return an error but found",
					v.val.methodType)
			} else {

				var result serverKeyResp
				_ = json.Unmarshal(msg.Result, &result)

				if result.Pubkey == "" {
					t.Fatal("expected the server pubkey not to be empty")
				}

				now := time.Now().UTC()
				if time.Unix(int64(result.Expiry), 0).UTC().After(now) {
					t.Fatalf("expected the server pubkey expiry to be after %v", now)
				}

				// No error was expected, prevent further error check.
				return
			}

			if msg.Error.Code != v.val.errCode {
				t.Fatalf("expected returned error code to be %q but found %q",
					msg.Error.Code, v.val.errCode)
			}

			err, _ = msg.Error.Data.(error)
			if err != v.val.shortErr {
				t.Fatalf("expected returned short error to be %q but found %q",
					err, v.val.shortErr)
			}

			if msg.Error.Message != v.val.longErr {
				t.Fatalf("expected returned long error to be %q but found %q",
					err, v.val.shortErr)
			}
		})
	}
}
