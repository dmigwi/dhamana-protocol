// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ChatMetaData contains all meta data concerning the Chat contract.
var ChatMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"principal\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"couponRate\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"couponDate\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maturityDate\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"enumBondContract.CurrencyType\",\"name\":\"currency\",\"type\":\"uint8\"}],\"name\":\"BondBodyTerms\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"BondMotivation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumBondContract.StatusChoice\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"BondStatusSigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"holder\",\"type\":\"address\"}],\"name\":\"HolderUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"}],\"name\":\"NewBondCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"NewChatMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumBondContract.StatusChoice\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"StatusChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumBondContract.StatusChoice\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"StatusSigned\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"enumChatContract.MessageTag\",\"name\":\"_tag\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_message\",\"type\":\"string\"}],\"name\":\"addMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createBond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"getBondSecureDetails\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"_security\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_appendix\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"signBondStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_principal\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"_couponRate\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"_couponDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maturityDate\",\"type\":\"uint32\"},{\"internalType\":\"enumBondContract.CurrencyType\",\"name\":\"_currency\",\"type\":\"uint8\"}],\"name\":\"updateBodyInfo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_holder\",\"type\":\"address\"}],\"name\":\"updateBondHolder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"enumBondContract.StatusChoice\",\"name\":\"_status\",\"type\":\"uint8\"}],\"name\":\"updateBondStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ChatABI is the input ABI used to generate the binding from.
// Deprecated: Use ChatMetaData.ABI instead.
var ChatABI = ChatMetaData.ABI

// Chat is an auto generated Go binding around an Ethereum contract.
type Chat struct {
	ChatCaller     // Read-only binding to the contract
	ChatTransactor // Write-only binding to the contract
	ChatFilterer   // Log filterer for contract events
}

// ChatCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChatCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChatTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChatTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChatFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChatFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChatSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChatSession struct {
	Contract     *Chat             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChatCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChatCallerSession struct {
	Contract *ChatCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ChatTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChatTransactorSession struct {
	Contract     *ChatTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChatRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChatRaw struct {
	Contract *Chat // Generic contract binding to access the raw methods on
}

// ChatCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChatCallerRaw struct {
	Contract *ChatCaller // Generic read-only contract binding to access the raw methods on
}

// ChatTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChatTransactorRaw struct {
	Contract *ChatTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChat creates a new instance of Chat, bound to a specific deployed contract.
func NewChat(address common.Address, backend bind.ContractBackend) (*Chat, error) {
	contract, err := bindChat(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chat{ChatCaller: ChatCaller{contract: contract}, ChatTransactor: ChatTransactor{contract: contract}, ChatFilterer: ChatFilterer{contract: contract}}, nil
}

// NewChatCaller creates a new read-only instance of Chat, bound to a specific deployed contract.
func NewChatCaller(address common.Address, caller bind.ContractCaller) (*ChatCaller, error) {
	contract, err := bindChat(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChatCaller{contract: contract}, nil
}

// NewChatTransactor creates a new write-only instance of Chat, bound to a specific deployed contract.
func NewChatTransactor(address common.Address, transactor bind.ContractTransactor) (*ChatTransactor, error) {
	contract, err := bindChat(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChatTransactor{contract: contract}, nil
}

// NewChatFilterer creates a new log filterer instance of Chat, bound to a specific deployed contract.
func NewChatFilterer(address common.Address, filterer bind.ContractFilterer) (*ChatFilterer, error) {
	contract, err := bindChat(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChatFilterer{contract: contract}, nil
}

// bindChat binds a generic wrapper to an already deployed contract.
func bindChat(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChatMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chat *ChatRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chat.Contract.ChatCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chat *ChatRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chat.Contract.ChatTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chat *ChatRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chat.Contract.ChatTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chat *ChatCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chat.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chat *ChatTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chat.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chat *ChatTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chat.Contract.contract.Transact(opts, method, params...)
}

// GetBondSecureDetails is a free data retrieval call binding the contract method 0xc3c95fa3.
//
// Solidity: function getBondSecureDetails(address _contract) view returns(string _security, string _appendix)
func (_Chat *ChatCaller) GetBondSecureDetails(opts *bind.CallOpts, _contract common.Address) (struct {
	Security string
	Appendix string
}, error) {
	var out []interface{}
	err := _Chat.contract.Call(opts, &out, "getBondSecureDetails", _contract)

	outstruct := new(struct {
		Security string
		Appendix string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Security = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Appendix = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// GetBondSecureDetails is a free data retrieval call binding the contract method 0xc3c95fa3.
//
// Solidity: function getBondSecureDetails(address _contract) view returns(string _security, string _appendix)
func (_Chat *ChatSession) GetBondSecureDetails(_contract common.Address) (struct {
	Security string
	Appendix string
}, error) {
	return _Chat.Contract.GetBondSecureDetails(&_Chat.CallOpts, _contract)
}

// GetBondSecureDetails is a free data retrieval call binding the contract method 0xc3c95fa3.
//
// Solidity: function getBondSecureDetails(address _contract) view returns(string _security, string _appendix)
func (_Chat *ChatCallerSession) GetBondSecureDetails(_contract common.Address) (struct {
	Security string
	Appendix string
}, error) {
	return _Chat.Contract.GetBondSecureDetails(&_Chat.CallOpts, _contract)
}

// AddMessage is a paid mutator transaction binding the contract method 0x7efed061.
//
// Solidity: function addMessage(address _contract, uint8 _tag, string _message) returns()
func (_Chat *ChatTransactor) AddMessage(opts *bind.TransactOpts, _contract common.Address, _tag uint8, _message string) (*types.Transaction, error) {
	return _Chat.contract.Transact(opts, "addMessage", _contract, _tag, _message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x7efed061.
//
// Solidity: function addMessage(address _contract, uint8 _tag, string _message) returns()
func (_Chat *ChatSession) AddMessage(_contract common.Address, _tag uint8, _message string) (*types.Transaction, error) {
	return _Chat.Contract.AddMessage(&_Chat.TransactOpts, _contract, _tag, _message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x7efed061.
//
// Solidity: function addMessage(address _contract, uint8 _tag, string _message) returns()
func (_Chat *ChatTransactorSession) AddMessage(_contract common.Address, _tag uint8, _message string) (*types.Transaction, error) {
	return _Chat.Contract.AddMessage(&_Chat.TransactOpts, _contract, _tag, _message)
}

// CreateBond is a paid mutator transaction binding the contract method 0xa091b795.
//
// Solidity: function createBond() returns()
func (_Chat *ChatTransactor) CreateBond(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chat.contract.Transact(opts, "createBond")
}

// CreateBond is a paid mutator transaction binding the contract method 0xa091b795.
//
// Solidity: function createBond() returns()
func (_Chat *ChatSession) CreateBond() (*types.Transaction, error) {
	return _Chat.Contract.CreateBond(&_Chat.TransactOpts)
}

// CreateBond is a paid mutator transaction binding the contract method 0xa091b795.
//
// Solidity: function createBond() returns()
func (_Chat *ChatTransactorSession) CreateBond() (*types.Transaction, error) {
	return _Chat.Contract.CreateBond(&_Chat.TransactOpts)
}

// SignBondStatus is a paid mutator transaction binding the contract method 0x8ec13d50.
//
// Solidity: function signBondStatus(address _contract) returns()
func (_Chat *ChatTransactor) SignBondStatus(opts *bind.TransactOpts, _contract common.Address) (*types.Transaction, error) {
	return _Chat.contract.Transact(opts, "signBondStatus", _contract)
}

// SignBondStatus is a paid mutator transaction binding the contract method 0x8ec13d50.
//
// Solidity: function signBondStatus(address _contract) returns()
func (_Chat *ChatSession) SignBondStatus(_contract common.Address) (*types.Transaction, error) {
	return _Chat.Contract.SignBondStatus(&_Chat.TransactOpts, _contract)
}

// SignBondStatus is a paid mutator transaction binding the contract method 0x8ec13d50.
//
// Solidity: function signBondStatus(address _contract) returns()
func (_Chat *ChatTransactorSession) SignBondStatus(_contract common.Address) (*types.Transaction, error) {
	return _Chat.Contract.SignBondStatus(&_Chat.TransactOpts, _contract)
}

// UpdateBodyInfo is a paid mutator transaction binding the contract method 0xac20a57e.
//
// Solidity: function updateBodyInfo(address _contract, uint32 _principal, uint8 _couponRate, uint32 _couponDate, uint32 _maturityDate, uint8 _currency) returns()
func (_Chat *ChatTransactor) UpdateBodyInfo(opts *bind.TransactOpts, _contract common.Address, _principal uint32, _couponRate uint8, _couponDate uint32, _maturityDate uint32, _currency uint8) (*types.Transaction, error) {
	return _Chat.contract.Transact(opts, "updateBodyInfo", _contract, _principal, _couponRate, _couponDate, _maturityDate, _currency)
}

// UpdateBodyInfo is a paid mutator transaction binding the contract method 0xac20a57e.
//
// Solidity: function updateBodyInfo(address _contract, uint32 _principal, uint8 _couponRate, uint32 _couponDate, uint32 _maturityDate, uint8 _currency) returns()
func (_Chat *ChatSession) UpdateBodyInfo(_contract common.Address, _principal uint32, _couponRate uint8, _couponDate uint32, _maturityDate uint32, _currency uint8) (*types.Transaction, error) {
	return _Chat.Contract.UpdateBodyInfo(&_Chat.TransactOpts, _contract, _principal, _couponRate, _couponDate, _maturityDate, _currency)
}

// UpdateBodyInfo is a paid mutator transaction binding the contract method 0xac20a57e.
//
// Solidity: function updateBodyInfo(address _contract, uint32 _principal, uint8 _couponRate, uint32 _couponDate, uint32 _maturityDate, uint8 _currency) returns()
func (_Chat *ChatTransactorSession) UpdateBodyInfo(_contract common.Address, _principal uint32, _couponRate uint8, _couponDate uint32, _maturityDate uint32, _currency uint8) (*types.Transaction, error) {
	return _Chat.Contract.UpdateBodyInfo(&_Chat.TransactOpts, _contract, _principal, _couponRate, _couponDate, _maturityDate, _currency)
}

// UpdateBondHolder is a paid mutator transaction binding the contract method 0xd5baca22.
//
// Solidity: function updateBondHolder(address _contract, address _holder) returns()
func (_Chat *ChatTransactor) UpdateBondHolder(opts *bind.TransactOpts, _contract common.Address, _holder common.Address) (*types.Transaction, error) {
	return _Chat.contract.Transact(opts, "updateBondHolder", _contract, _holder)
}

// UpdateBondHolder is a paid mutator transaction binding the contract method 0xd5baca22.
//
// Solidity: function updateBondHolder(address _contract, address _holder) returns()
func (_Chat *ChatSession) UpdateBondHolder(_contract common.Address, _holder common.Address) (*types.Transaction, error) {
	return _Chat.Contract.UpdateBondHolder(&_Chat.TransactOpts, _contract, _holder)
}

// UpdateBondHolder is a paid mutator transaction binding the contract method 0xd5baca22.
//
// Solidity: function updateBondHolder(address _contract, address _holder) returns()
func (_Chat *ChatTransactorSession) UpdateBondHolder(_contract common.Address, _holder common.Address) (*types.Transaction, error) {
	return _Chat.Contract.UpdateBondHolder(&_Chat.TransactOpts, _contract, _holder)
}

// UpdateBondStatus is a paid mutator transaction binding the contract method 0x0c3535dc.
//
// Solidity: function updateBondStatus(address _contract, uint8 _status) returns()
func (_Chat *ChatTransactor) UpdateBondStatus(opts *bind.TransactOpts, _contract common.Address, _status uint8) (*types.Transaction, error) {
	return _Chat.contract.Transact(opts, "updateBondStatus", _contract, _status)
}

// UpdateBondStatus is a paid mutator transaction binding the contract method 0x0c3535dc.
//
// Solidity: function updateBondStatus(address _contract, uint8 _status) returns()
func (_Chat *ChatSession) UpdateBondStatus(_contract common.Address, _status uint8) (*types.Transaction, error) {
	return _Chat.Contract.UpdateBondStatus(&_Chat.TransactOpts, _contract, _status)
}

// UpdateBondStatus is a paid mutator transaction binding the contract method 0x0c3535dc.
//
// Solidity: function updateBondStatus(address _contract, uint8 _status) returns()
func (_Chat *ChatTransactorSession) UpdateBondStatus(_contract common.Address, _status uint8) (*types.Transaction, error) {
	return _Chat.Contract.UpdateBondStatus(&_Chat.TransactOpts, _contract, _status)
}

// ChatBondBodyTermsIterator is returned from FilterBondBodyTerms and is used to iterate over the raw logs and unpacked data for BondBodyTerms events raised by the Chat contract.
type ChatBondBodyTermsIterator struct {
	Event *ChatBondBodyTerms // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChatBondBodyTermsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatBondBodyTerms)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChatBondBodyTerms)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChatBondBodyTermsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatBondBodyTermsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatBondBodyTerms represents a BondBodyTerms event raised by the Chat contract.
type ChatBondBodyTerms struct {
	BondAddress  common.Address
	Principal    uint32
	CouponRate   uint8
	CouponDate   uint32
	MaturityDate uint32
	Currency     uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBondBodyTerms is a free log retrieval operation binding the contract event 0xe6f76c15f92de04f9874c9963972f0244d7957b298d1b1bd00ffc99ea8bc1794.
//
// Solidity: event BondBodyTerms(address bondAddress, uint32 principal, uint8 couponRate, uint32 couponDate, uint32 maturityDate, uint8 currency)
func (_Chat *ChatFilterer) FilterBondBodyTerms(opts *bind.FilterOpts) (*ChatBondBodyTermsIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "BondBodyTerms")
	if err != nil {
		return nil, err
	}
	return &ChatBondBodyTermsIterator{contract: _Chat.contract, event: "BondBodyTerms", logs: logs, sub: sub}, nil
}

// WatchBondBodyTerms is a free log subscription operation binding the contract event 0xe6f76c15f92de04f9874c9963972f0244d7957b298d1b1bd00ffc99ea8bc1794.
//
// Solidity: event BondBodyTerms(address bondAddress, uint32 principal, uint8 couponRate, uint32 couponDate, uint32 maturityDate, uint8 currency)
func (_Chat *ChatFilterer) WatchBondBodyTerms(opts *bind.WatchOpts, sink chan<- *ChatBondBodyTerms) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "BondBodyTerms")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatBondBodyTerms)
				if err := _Chat.contract.UnpackLog(event, "BondBodyTerms", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBondBodyTerms is a log parse operation binding the contract event 0xe6f76c15f92de04f9874c9963972f0244d7957b298d1b1bd00ffc99ea8bc1794.
//
// Solidity: event BondBodyTerms(address bondAddress, uint32 principal, uint8 couponRate, uint32 couponDate, uint32 maturityDate, uint8 currency)
func (_Chat *ChatFilterer) ParseBondBodyTerms(log types.Log) (*ChatBondBodyTerms, error) {
	event := new(ChatBondBodyTerms)
	if err := _Chat.contract.UnpackLog(event, "BondBodyTerms", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatBondMotivationIterator is returned from FilterBondMotivation and is used to iterate over the raw logs and unpacked data for BondMotivation events raised by the Chat contract.
type ChatBondMotivationIterator struct {
	Event *ChatBondMotivation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChatBondMotivationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatBondMotivation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChatBondMotivation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChatBondMotivationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatBondMotivationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatBondMotivation represents a BondMotivation event raised by the Chat contract.
type ChatBondMotivation struct {
	BondAddress common.Address
	Message     string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBondMotivation is a free log retrieval operation binding the contract event 0x37e80f8f8deab7dbce22e881bee00b243b0ed1880ff122ed20d655ce0acc4bc8.
//
// Solidity: event BondMotivation(address bondAddress, string message)
func (_Chat *ChatFilterer) FilterBondMotivation(opts *bind.FilterOpts) (*ChatBondMotivationIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "BondMotivation")
	if err != nil {
		return nil, err
	}
	return &ChatBondMotivationIterator{contract: _Chat.contract, event: "BondMotivation", logs: logs, sub: sub}, nil
}

// WatchBondMotivation is a free log subscription operation binding the contract event 0x37e80f8f8deab7dbce22e881bee00b243b0ed1880ff122ed20d655ce0acc4bc8.
//
// Solidity: event BondMotivation(address bondAddress, string message)
func (_Chat *ChatFilterer) WatchBondMotivation(opts *bind.WatchOpts, sink chan<- *ChatBondMotivation) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "BondMotivation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatBondMotivation)
				if err := _Chat.contract.UnpackLog(event, "BondMotivation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBondMotivation is a log parse operation binding the contract event 0x37e80f8f8deab7dbce22e881bee00b243b0ed1880ff122ed20d655ce0acc4bc8.
//
// Solidity: event BondMotivation(address bondAddress, string message)
func (_Chat *ChatFilterer) ParseBondMotivation(log types.Log) (*ChatBondMotivation, error) {
	event := new(ChatBondMotivation)
	if err := _Chat.contract.UnpackLog(event, "BondMotivation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatBondStatusSignedIterator is returned from FilterBondStatusSigned and is used to iterate over the raw logs and unpacked data for BondStatusSigned events raised by the Chat contract.
type ChatBondStatusSignedIterator struct {
	Event *ChatBondStatusSigned // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChatBondStatusSignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatBondStatusSigned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChatBondStatusSigned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChatBondStatusSignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatBondStatusSignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatBondStatusSigned represents a BondStatusSigned event raised by the Chat contract.
type ChatBondStatusSigned struct {
	Sender      common.Address
	BondAddress common.Address
	Status      uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBondStatusSigned is a free log retrieval operation binding the contract event 0xd76eace44117b2647270cc02d05e0bb7f011d97783726d3900c11c952b137367.
//
// Solidity: event BondStatusSigned(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) FilterBondStatusSigned(opts *bind.FilterOpts) (*ChatBondStatusSignedIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "BondStatusSigned")
	if err != nil {
		return nil, err
	}
	return &ChatBondStatusSignedIterator{contract: _Chat.contract, event: "BondStatusSigned", logs: logs, sub: sub}, nil
}

// WatchBondStatusSigned is a free log subscription operation binding the contract event 0xd76eace44117b2647270cc02d05e0bb7f011d97783726d3900c11c952b137367.
//
// Solidity: event BondStatusSigned(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) WatchBondStatusSigned(opts *bind.WatchOpts, sink chan<- *ChatBondStatusSigned) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "BondStatusSigned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatBondStatusSigned)
				if err := _Chat.contract.UnpackLog(event, "BondStatusSigned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBondStatusSigned is a log parse operation binding the contract event 0xd76eace44117b2647270cc02d05e0bb7f011d97783726d3900c11c952b137367.
//
// Solidity: event BondStatusSigned(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) ParseBondStatusSigned(log types.Log) (*ChatBondStatusSigned, error) {
	event := new(ChatBondStatusSigned)
	if err := _Chat.contract.UnpackLog(event, "BondStatusSigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatHolderUpdateIterator is returned from FilterHolderUpdate and is used to iterate over the raw logs and unpacked data for HolderUpdate events raised by the Chat contract.
type ChatHolderUpdateIterator struct {
	Event *ChatHolderUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChatHolderUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatHolderUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChatHolderUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChatHolderUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatHolderUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatHolderUpdate represents a HolderUpdate event raised by the Chat contract.
type ChatHolderUpdate struct {
	BondAddress common.Address
	Holder      common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterHolderUpdate is a free log retrieval operation binding the contract event 0x316e906a1dc3867c66eacd02598322c6e106e4978e6e21aa8bb15a110b86898c.
//
// Solidity: event HolderUpdate(address bondAddress, address holder)
func (_Chat *ChatFilterer) FilterHolderUpdate(opts *bind.FilterOpts) (*ChatHolderUpdateIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "HolderUpdate")
	if err != nil {
		return nil, err
	}
	return &ChatHolderUpdateIterator{contract: _Chat.contract, event: "HolderUpdate", logs: logs, sub: sub}, nil
}

// WatchHolderUpdate is a free log subscription operation binding the contract event 0x316e906a1dc3867c66eacd02598322c6e106e4978e6e21aa8bb15a110b86898c.
//
// Solidity: event HolderUpdate(address bondAddress, address holder)
func (_Chat *ChatFilterer) WatchHolderUpdate(opts *bind.WatchOpts, sink chan<- *ChatHolderUpdate) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "HolderUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatHolderUpdate)
				if err := _Chat.contract.UnpackLog(event, "HolderUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHolderUpdate is a log parse operation binding the contract event 0x316e906a1dc3867c66eacd02598322c6e106e4978e6e21aa8bb15a110b86898c.
//
// Solidity: event HolderUpdate(address bondAddress, address holder)
func (_Chat *ChatFilterer) ParseHolderUpdate(log types.Log) (*ChatHolderUpdate, error) {
	event := new(ChatHolderUpdate)
	if err := _Chat.contract.UnpackLog(event, "HolderUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatNewBondCreatedIterator is returned from FilterNewBondCreated and is used to iterate over the raw logs and unpacked data for NewBondCreated events raised by the Chat contract.
type ChatNewBondCreatedIterator struct {
	Event *ChatNewBondCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChatNewBondCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatNewBondCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChatNewBondCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChatNewBondCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatNewBondCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatNewBondCreated represents a NewBondCreated event raised by the Chat contract.
type ChatNewBondCreated struct {
	Sender      common.Address
	BondAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewBondCreated is a free log retrieval operation binding the contract event 0xff5e78dafd87acba3d62bc893238ff9de367bc49ebc45993c9474fd881b577f2.
//
// Solidity: event NewBondCreated(address sender, address bondAddress)
func (_Chat *ChatFilterer) FilterNewBondCreated(opts *bind.FilterOpts) (*ChatNewBondCreatedIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "NewBondCreated")
	if err != nil {
		return nil, err
	}
	return &ChatNewBondCreatedIterator{contract: _Chat.contract, event: "NewBondCreated", logs: logs, sub: sub}, nil
}

// WatchNewBondCreated is a free log subscription operation binding the contract event 0xff5e78dafd87acba3d62bc893238ff9de367bc49ebc45993c9474fd881b577f2.
//
// Solidity: event NewBondCreated(address sender, address bondAddress)
func (_Chat *ChatFilterer) WatchNewBondCreated(opts *bind.WatchOpts, sink chan<- *ChatNewBondCreated) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "NewBondCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatNewBondCreated)
				if err := _Chat.contract.UnpackLog(event, "NewBondCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewBondCreated is a log parse operation binding the contract event 0xff5e78dafd87acba3d62bc893238ff9de367bc49ebc45993c9474fd881b577f2.
//
// Solidity: event NewBondCreated(address sender, address bondAddress)
func (_Chat *ChatFilterer) ParseNewBondCreated(log types.Log) (*ChatNewBondCreated, error) {
	event := new(ChatNewBondCreated)
	if err := _Chat.contract.UnpackLog(event, "NewBondCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatNewChatMessageIterator is returned from FilterNewChatMessage and is used to iterate over the raw logs and unpacked data for NewChatMessage events raised by the Chat contract.
type ChatNewChatMessageIterator struct {
	Event *ChatNewChatMessage // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChatNewChatMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatNewChatMessage)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChatNewChatMessage)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChatNewChatMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatNewChatMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatNewChatMessage represents a NewChatMessage event raised by the Chat contract.
type ChatNewChatMessage struct {
	BondAddress common.Address
	Sender      common.Address
	Message     string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewChatMessage is a free log retrieval operation binding the contract event 0xbd74e567035eb30415526c1a311d34e5675e23d477622bf0701f23208c6b4bd4.
//
// Solidity: event NewChatMessage(address bondAddress, address sender, string message)
func (_Chat *ChatFilterer) FilterNewChatMessage(opts *bind.FilterOpts) (*ChatNewChatMessageIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "NewChatMessage")
	if err != nil {
		return nil, err
	}
	return &ChatNewChatMessageIterator{contract: _Chat.contract, event: "NewChatMessage", logs: logs, sub: sub}, nil
}

// WatchNewChatMessage is a free log subscription operation binding the contract event 0xbd74e567035eb30415526c1a311d34e5675e23d477622bf0701f23208c6b4bd4.
//
// Solidity: event NewChatMessage(address bondAddress, address sender, string message)
func (_Chat *ChatFilterer) WatchNewChatMessage(opts *bind.WatchOpts, sink chan<- *ChatNewChatMessage) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "NewChatMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatNewChatMessage)
				if err := _Chat.contract.UnpackLog(event, "NewChatMessage", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewChatMessage is a log parse operation binding the contract event 0xbd74e567035eb30415526c1a311d34e5675e23d477622bf0701f23208c6b4bd4.
//
// Solidity: event NewChatMessage(address bondAddress, address sender, string message)
func (_Chat *ChatFilterer) ParseNewChatMessage(log types.Log) (*ChatNewChatMessage, error) {
	event := new(ChatNewChatMessage)
	if err := _Chat.contract.UnpackLog(event, "NewChatMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatStatusChangeIterator is returned from FilterStatusChange and is used to iterate over the raw logs and unpacked data for StatusChange events raised by the Chat contract.
type ChatStatusChangeIterator struct {
	Event *ChatStatusChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChatStatusChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatStatusChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChatStatusChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChatStatusChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatStatusChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatStatusChange represents a StatusChange event raised by the Chat contract.
type ChatStatusChange struct {
	Sender      common.Address
	BondAddress common.Address
	Status      uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStatusChange is a free log retrieval operation binding the contract event 0xc82b702f5d9f3997b9d40d2a61bb841083384d586e4b023ef72e017153beb09a.
//
// Solidity: event StatusChange(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) FilterStatusChange(opts *bind.FilterOpts) (*ChatStatusChangeIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "StatusChange")
	if err != nil {
		return nil, err
	}
	return &ChatStatusChangeIterator{contract: _Chat.contract, event: "StatusChange", logs: logs, sub: sub}, nil
}

// WatchStatusChange is a free log subscription operation binding the contract event 0xc82b702f5d9f3997b9d40d2a61bb841083384d586e4b023ef72e017153beb09a.
//
// Solidity: event StatusChange(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) WatchStatusChange(opts *bind.WatchOpts, sink chan<- *ChatStatusChange) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "StatusChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatStatusChange)
				if err := _Chat.contract.UnpackLog(event, "StatusChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStatusChange is a log parse operation binding the contract event 0xc82b702f5d9f3997b9d40d2a61bb841083384d586e4b023ef72e017153beb09a.
//
// Solidity: event StatusChange(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) ParseStatusChange(log types.Log) (*ChatStatusChange, error) {
	event := new(ChatStatusChange)
	if err := _Chat.contract.UnpackLog(event, "StatusChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatStatusSignedIterator is returned from FilterStatusSigned and is used to iterate over the raw logs and unpacked data for StatusSigned events raised by the Chat contract.
type ChatStatusSignedIterator struct {
	Event *ChatStatusSigned // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChatStatusSignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatStatusSigned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChatStatusSigned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChatStatusSignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatStatusSignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatStatusSigned represents a StatusSigned event raised by the Chat contract.
type ChatStatusSigned struct {
	Sender      common.Address
	BondAddress common.Address
	Status      uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStatusSigned is a free log retrieval operation binding the contract event 0x2988c7b4a31559db7f7be48700d500ea7d4b3c763dc6a23e3c6d1c03bed5f576.
//
// Solidity: event StatusSigned(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) FilterStatusSigned(opts *bind.FilterOpts) (*ChatStatusSignedIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "StatusSigned")
	if err != nil {
		return nil, err
	}
	return &ChatStatusSignedIterator{contract: _Chat.contract, event: "StatusSigned", logs: logs, sub: sub}, nil
}

// WatchStatusSigned is a free log subscription operation binding the contract event 0x2988c7b4a31559db7f7be48700d500ea7d4b3c763dc6a23e3c6d1c03bed5f576.
//
// Solidity: event StatusSigned(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) WatchStatusSigned(opts *bind.WatchOpts, sink chan<- *ChatStatusSigned) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "StatusSigned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatStatusSigned)
				if err := _Chat.contract.UnpackLog(event, "StatusSigned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStatusSigned is a log parse operation binding the contract event 0x2988c7b4a31559db7f7be48700d500ea7d4b3c763dc6a23e3c6d1c03bed5f576.
//
// Solidity: event StatusSigned(address sender, address bondAddress, uint8 status)
func (_Chat *ChatFilterer) ParseStatusSigned(log types.Log) (*ChatStatusSigned, error) {
	event := new(ChatStatusSigned)
	if err := _Chat.contract.UnpackLog(event, "StatusSigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
