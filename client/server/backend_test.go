package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/sapphire"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	serverConf = &ServerConfig{
		sessionKeys: new(sync.Map),
	}

	sampleHexAddress = common.HexToAddress("0x3396FD816Dd81100477c8ea3853039822f36B7ed")
	sampleSigningKey = "e33204e6138743d6908cc5caabda374c8dc7716912693764fcbf9ba7cf56dc8349bd51d416d750745a26441b341b834ced6fe2a874b50b91f26b405bbfe93d052b9e81f97033b88f25f52711675e4b88dfb165b5843604356cd8a88e4ed6"

	sampleHexAddress1 = common.HexToAddress("0x3396FD816Dd81100477c8ea3853039822f36B7ad")
	sampleHexAddress2 = common.HexToAddress("0x3396FD816Dd81100477c8ea3853039822f36B7bd")
	sampleHexAddress3 = common.HexToAddress("0x3396FD816Dd81100477c8ea3853039822f36B71d")

	pubkey1 = "0x041ebfc6b4cc5797953c4c95791fc67089912f69e081e28df62b7885e296df950" +
		"756633381d1f3869859b7203c243f852821e8121b3483edee45e75bf8727f02ff"
	pubkey2 = "0x04e00754d4e029a54e2d9be882521a25c428e8a559b98e6ea46b60fec95c0c8957f5b95522d9cb74881c49cf6ff3c28c1ea0743ff4e649b182cbce168173e6380a"
	pubkey3 = "0x0444a59dd97719656b4960a1a74340fc1b41bbc5435cf20198cc726f8ab4fe68b3" +
		"aab494052a57f1915cc6d10f4323862e608cba658241de921b6fcaeb46869726"

	sharedKey1 = "0x1e933d207f79abbfe644471e3a5e2d1907091f01a282dbf0861c89c8c79bb925"
	sharedKey2 = "0xef0f3ee401b103a5b535673a1d5f7638903e59774d204d6784faeb674a6fb068"
	sharedKey3 = "0x43915b03de29fe619623309517993af22c71b0bdf94b9994700337f07a4682f6"
)

type input struct {
	testName   string
	method     utils.Method
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

	err = genMockWrapper(ctx)
	if err == nil {
		key1, _ := hexutil.Decode(sharedKey1)
		key2, _ := hexutil.Decode(sharedKey2)

		// Store expired keys
		expiredKey := serverKeyResp{
			Pubkey:    pubkey1,
			Expiry:    uint64(time.Now().UTC().Unix()),
			sharedKey: key1,
		}
		serverConf.sessionKeys.Store(sampleHexAddress1, expiredKey)

		// store fresh keys with an expiry of 2 minutes
		freshKey := serverKeyResp{
			Pubkey:    pubkey2,
			Expiry:    uint64(time.Now().UTC().Add(10 * time.Minute).Unix()),
			sharedKey: key2,
		}
		serverConf.sessionKeys.Store(sampleHexAddress2, freshKey)

		m.Run()
	} else {
		fmt.Printf("unexpected error: %v \n", err)
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
					Method: utils.GetServerPubKey,
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
				longErr:    "json: cannot unmarshal string into Go value of type server.rpcMessage",
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
					Method:  utils.GetServerPubKey,
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
					Method:  utils.GetServerPubKey,
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
					Method:  utils.CreateBond,
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
					Method:  utils.Method("createBondAndSign"),
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
					Method:  utils.GetServerPubKey,
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
					Method:  utils.GetServerPubKey,
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{int(200)},
				},
			},
			val: output{
				errCode:    1009,
				shortErr:   utils.ErrUnknownParam,
				longErr:    "expected param 200 to be of type string but found it to be number",
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
					Method:  utils.GetServerPubKey,
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
					Method:  utils.SignBondStatus,
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
					Method:  utils.GetBondByAddress,
					Sender: &senderInfo{
						Address:    sampleHexAddress,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{sampleHexAddress},
				},
			},
			val: output{
				errCode:    1008,
				shortErr:   utils.ErrUnknownMethod,
				longErr:    "unsupported method getBondsByStatus found for this route",
				methodType: utils.LocalType,
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.data.testName, func(t *testing.T) {
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(v.data.body) // error ignored since its not being tested.

			req := httptest.NewRequest(string(v.data.method), "/random-path", &buf)

			msg := rpcMessage{}
			retType := decodeRequestBody(req, &msg, v.data.needSigner)

			if retType != v.val.methodType {
				t.Fatalf("expected returned method type to be %v but found %v",
					int(retType), int(v.val.methodType))
			}

			// Test for the packed response.
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

			if msg.Sender != nil {
				t.Fatal("expected sender field to be nil")
			}

			if msg.Method != "" {
				t.Fatal("expected method field to be empty")
			}

			if msg.Params != nil {
				t.Fatal("expected params field to be nil")
			}

			if msg.Error.Code != v.val.errCode {
				t.Fatalf("expected returned error code to be %q but found %q",
					msg.Error.Code, v.val.errCode)
			}

			if msg.Error.Message != v.val.shortErr.Error() {
				t.Fatalf("expected returned short error to be %q but found %q",
					msg.Error.Message, v.val.shortErr)
			}

			errStr, _ := msg.Error.Data.(string)
			if errStr != v.val.longErr {
				t.Fatalf("expected returned long error to be %q but found %q",
					errStr, v.val.longErr)
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

func genMockWrapper(ctx context.Context) error {
	conn := &mockWrapper{}
	backend, err := sapphire.WrapClient(ctx, conn, utils.LocalTesting,
		func(digest [32]byte, _ []byte) ([]byte, error) { return digest[:], nil })
	if err != nil {
		return err
	}

	// Create the chat instance to be used.
	chatInstance, err := contracts.NewChat(sampleHexAddress, backend)
	if err != nil {
		return err
	}

	serverConf.backend = backend
	serverConf.bondChat = chatInstance

	return nil
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
					Method:  utils.GetBondByAddress,
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{sampleHexAddress1},
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
					Method:  utils.GetServerPubKey,
					Sender: &senderInfo{
						Address: sampleHexAddress,
					},
					Params: []interface{}{pubkey2},
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
				httptest.NewRequest(string(v.data.method), "/serverpubkey", &buf))

			data, err := io.ReadAll(responseWritter.Body)
			if err != nil {
				t.Fatalf("expected no error but found %q", err)
			}

			msg := rpcMessage{}
			_ = json.Unmarshal(data, &msg)

			if msg.Error == nil && v.val.shortErr != nil {
				t.Fatalf("expected method %d to return an error but found none",
					int(v.val.methodType))
			}

			if msg.Error == nil {
				var result serverKeyResp
				_ = json.Unmarshal(msg.Result, &result)

				if result.Pubkey == "" {
					t.Fatal("expected the server pubkey not to be empty")
				}

				now := time.Now().UTC()
				expired := time.Unix(int64(result.Expiry), 0).UTC()
				if now.After(expired) {
					t.Fatalf("expected the server pubkey expiry %q to be before %q", expired, now)
				}

				// No error was expected, prevent further error check.
				return
			}

			if msg.Error.Code != v.val.errCode {
				t.Fatalf("expected returned error code to be %d but found %d",
					msg.Error.Code, v.val.errCode)
			}

			if msg.Error.Message != v.val.shortErr.Error() {
				t.Fatalf("expected returned short error to be %q but found %q",
					msg.Error.Message, v.val.shortErr)
			}

			errStr, _ := msg.Error.Data.(string)
			if errStr != v.val.longErr {
				t.Fatalf("expected returned long error to be %q but found %q",
					errStr, v.val.longErr)
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
					Method:  utils.GetServerPubKey,
					Sender: &senderInfo{
						Address:    sampleHexAddress2,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{pubkey2},
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
					Method:  utils.SignBondStatus,
					Sender: &senderInfo{
						Address:    sampleHexAddress3,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{sampleHexAddress},
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
					Method:  utils.SignBondStatus,
					Sender: &senderInfo{
						Address:    sampleHexAddress1,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{sampleHexAddress},
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
				testName: "Test-for-successful-access-to-contract-method-with-no-param",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  utils.CreateBond,
					Sender: &senderInfo{
						Address:    sampleHexAddress2,
						SigningKey: sampleSigningKey,
					},
				},
			},
		},
		{
			data: input{
				testName: "Test-for-successful-access-to-contract-method-with-one-param",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  utils.SignBondStatus,
					Sender: &senderInfo{
						Address:    sampleHexAddress2,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{sampleHexAddress},
				},
			},
		},
		{
			data: input{
				testName: "Test-for-successful-access-to-contract-method-with-multiple-params",
				method:   http.MethodPost,
				body: rpcMessage{
					ID:      20,
					Version: "2.0",
					Method:  utils.UpdateBondStatus,
					Sender: &senderInfo{
						Address:    sampleHexAddress2,
						SigningKey: sampleSigningKey,
					},
					Params: []interface{}{sampleHexAddress, 2},
				},
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.data.testName, func(t *testing.T) {
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(v.data.body) // error ignored since its not being tested.

			responseWritter := httptest.NewRecorder()

			serverConf.backendQueryFunc(responseWritter,
				httptest.NewRequest(string(v.data.method), "/backend", &buf))

			data, err := io.ReadAll(responseWritter.Body)
			if err != nil {
				t.Fatalf("expected no error but found %q", err)
			}

			msg := rpcMessage{}
			_ = json.Unmarshal(data, &msg)

			if msg.Error == nil && v.val.shortErr != nil {
				t.Fatalf("expected method %d to return an error but found none",
					int(v.val.methodType))
			}

			if msg.Error == nil {
				if len(msg.Result) == 0 {
					t.Fatal("expected the Result data not to be empty")
				}

				// No error was expected, prevent further error check.
				return
			}

			if msg.Error.Code != v.val.errCode {
				t.Fatalf("expected returned error code to be %d but found %d",
					msg.Error.Code, v.val.errCode)
			}

			if msg.Error.Message != v.val.shortErr.Error() {
				t.Fatalf("expected returned short error to be %q but found %q",
					msg.Error.Message, v.val.shortErr)
			}

			errStr, _ := msg.Error.Data.(string)
			if errStr != v.val.longErr {
				t.Fatalf("expected returned long error to be %q but found %q",
					errStr, v.val.longErr)
			}
		})
	}
}
