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

// ChatContractMessageInfo is an auto generated low-level Go binding around an user-defined struct.
type ChatContractMessageInfo struct {
	Sender    common.Address
	Message   string
	Timestamp *big.Int
}

// ChatMetaData contains all meta data concerning the Chat contract.
var ChatMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"}],\"name\":\"BondDisputeResolved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"BondMotivation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"}],\"name\":\"BondUnderDispute\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"principal\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"couponRate\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"couponDate\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maturityDate\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"enumBondContract.CurrencyType\",\"name\":\"currency\",\"type\":\"uint8\"}],\"name\":\"FinalBondTerms\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"NewBondCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bondAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structChatContract.MessageInfo\",\"name\":\"chat\",\"type\":\"tuple\"}],\"name\":\"NewChatMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"enumChatContract.MessageTag\",\"name\":\"_tag\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_message\",\"type\":\"string\"}],\"name\":\"addMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createBond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"getBondSecureDetails\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"_security\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_appendix\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"signBondStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_principal\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"_couponRate\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"_couponDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_maturityDate\",\"type\":\"uint32\"},{\"internalType\":\"enumBondContract.CurrencyType\",\"name\":\"_currency\",\"type\":\"uint8\"}],\"name\":\"updateBodyInfo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_holder\",\"type\":\"address\"}],\"name\":\"updateBondHolder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"enumBondContract.StatusChoice\",\"name\":\"_status\",\"type\":\"uint8\"}],\"name\":\"updateBondStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// ChatBondDisputeResolvedIterator is returned from FilterBondDisputeResolved and is used to iterate over the raw logs and unpacked data for BondDisputeResolved events raised by the Chat contract.
type ChatBondDisputeResolvedIterator struct {
	Event *ChatBondDisputeResolved // Event containing the contract specifics and raw log

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
func (it *ChatBondDisputeResolvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatBondDisputeResolved)
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
		it.Event = new(ChatBondDisputeResolved)
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
func (it *ChatBondDisputeResolvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatBondDisputeResolvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatBondDisputeResolved represents a BondDisputeResolved event raised by the Chat contract.
type ChatBondDisputeResolved struct {
	Sender      common.Address
	BondAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBondDisputeResolved is a free log retrieval operation binding the contract event 0x67a381a6a07926a6859febbf3faccb07a5d19deef5b03561a631757579f1c89c.
//
// Solidity: event BondDisputeResolved(address sender, address bondAddress)
func (_Chat *ChatFilterer) FilterBondDisputeResolved(opts *bind.FilterOpts) (*ChatBondDisputeResolvedIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "BondDisputeResolved")
	if err != nil {
		return nil, err
	}
	return &ChatBondDisputeResolvedIterator{contract: _Chat.contract, event: "BondDisputeResolved", logs: logs, sub: sub}, nil
}

// WatchBondDisputeResolved is a free log subscription operation binding the contract event 0x67a381a6a07926a6859febbf3faccb07a5d19deef5b03561a631757579f1c89c.
//
// Solidity: event BondDisputeResolved(address sender, address bondAddress)
func (_Chat *ChatFilterer) WatchBondDisputeResolved(opts *bind.WatchOpts, sink chan<- *ChatBondDisputeResolved) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "BondDisputeResolved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatBondDisputeResolved)
				if err := _Chat.contract.UnpackLog(event, "BondDisputeResolved", log); err != nil {
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

// ParseBondDisputeResolved is a log parse operation binding the contract event 0x67a381a6a07926a6859febbf3faccb07a5d19deef5b03561a631757579f1c89c.
//
// Solidity: event BondDisputeResolved(address sender, address bondAddress)
func (_Chat *ChatFilterer) ParseBondDisputeResolved(log types.Log) (*ChatBondDisputeResolved, error) {
	event := new(ChatBondDisputeResolved)
	if err := _Chat.contract.UnpackLog(event, "BondDisputeResolved", log); err != nil {
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
	Sender      common.Address
	BondAddress common.Address
	Message     string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBondMotivation is a free log retrieval operation binding the contract event 0x9c2db9547a0e615dbb8f26a29f24669e9cef62f8dd222509384887635c254d23.
//
// Solidity: event BondMotivation(address sender, address bondAddress, string message)
func (_Chat *ChatFilterer) FilterBondMotivation(opts *bind.FilterOpts) (*ChatBondMotivationIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "BondMotivation")
	if err != nil {
		return nil, err
	}
	return &ChatBondMotivationIterator{contract: _Chat.contract, event: "BondMotivation", logs: logs, sub: sub}, nil
}

// WatchBondMotivation is a free log subscription operation binding the contract event 0x9c2db9547a0e615dbb8f26a29f24669e9cef62f8dd222509384887635c254d23.
//
// Solidity: event BondMotivation(address sender, address bondAddress, string message)
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

// ParseBondMotivation is a log parse operation binding the contract event 0x9c2db9547a0e615dbb8f26a29f24669e9cef62f8dd222509384887635c254d23.
//
// Solidity: event BondMotivation(address sender, address bondAddress, string message)
func (_Chat *ChatFilterer) ParseBondMotivation(log types.Log) (*ChatBondMotivation, error) {
	event := new(ChatBondMotivation)
	if err := _Chat.contract.UnpackLog(event, "BondMotivation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatBondUnderDisputeIterator is returned from FilterBondUnderDispute and is used to iterate over the raw logs and unpacked data for BondUnderDispute events raised by the Chat contract.
type ChatBondUnderDisputeIterator struct {
	Event *ChatBondUnderDispute // Event containing the contract specifics and raw log

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
func (it *ChatBondUnderDisputeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatBondUnderDispute)
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
		it.Event = new(ChatBondUnderDispute)
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
func (it *ChatBondUnderDisputeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatBondUnderDisputeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatBondUnderDispute represents a BondUnderDispute event raised by the Chat contract.
type ChatBondUnderDispute struct {
	Sender      common.Address
	BondAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBondUnderDispute is a free log retrieval operation binding the contract event 0x5edcf956463f8d45ca26b3d5fb4fc9cd24aed980e33f4ccb131805819ddff187.
//
// Solidity: event BondUnderDispute(address sender, address bondAddress)
func (_Chat *ChatFilterer) FilterBondUnderDispute(opts *bind.FilterOpts) (*ChatBondUnderDisputeIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "BondUnderDispute")
	if err != nil {
		return nil, err
	}
	return &ChatBondUnderDisputeIterator{contract: _Chat.contract, event: "BondUnderDispute", logs: logs, sub: sub}, nil
}

// WatchBondUnderDispute is a free log subscription operation binding the contract event 0x5edcf956463f8d45ca26b3d5fb4fc9cd24aed980e33f4ccb131805819ddff187.
//
// Solidity: event BondUnderDispute(address sender, address bondAddress)
func (_Chat *ChatFilterer) WatchBondUnderDispute(opts *bind.WatchOpts, sink chan<- *ChatBondUnderDispute) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "BondUnderDispute")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatBondUnderDispute)
				if err := _Chat.contract.UnpackLog(event, "BondUnderDispute", log); err != nil {
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

// ParseBondUnderDispute is a log parse operation binding the contract event 0x5edcf956463f8d45ca26b3d5fb4fc9cd24aed980e33f4ccb131805819ddff187.
//
// Solidity: event BondUnderDispute(address sender, address bondAddress)
func (_Chat *ChatFilterer) ParseBondUnderDispute(log types.Log) (*ChatBondUnderDispute, error) {
	event := new(ChatBondUnderDispute)
	if err := _Chat.contract.UnpackLog(event, "BondUnderDispute", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChatFinalBondTermsIterator is returned from FilterFinalBondTerms and is used to iterate over the raw logs and unpacked data for FinalBondTerms events raised by the Chat contract.
type ChatFinalBondTermsIterator struct {
	Event *ChatFinalBondTerms // Event containing the contract specifics and raw log

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
func (it *ChatFinalBondTermsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChatFinalBondTerms)
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
		it.Event = new(ChatFinalBondTerms)
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
func (it *ChatFinalBondTermsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChatFinalBondTermsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChatFinalBondTerms represents a FinalBondTerms event raised by the Chat contract.
type ChatFinalBondTerms struct {
	Principal    uint32
	CouponRate   uint8
	CouponDate   uint32
	MaturityDate uint32
	Currency     uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFinalBondTerms is a free log retrieval operation binding the contract event 0x7e5631fad24e4cbdaaf2d85082725b7cc657fb9d32a0025131310dca36d2eb41.
//
// Solidity: event FinalBondTerms(uint32 principal, uint8 couponRate, uint32 couponDate, uint32 maturityDate, uint8 currency)
func (_Chat *ChatFilterer) FilterFinalBondTerms(opts *bind.FilterOpts) (*ChatFinalBondTermsIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "FinalBondTerms")
	if err != nil {
		return nil, err
	}
	return &ChatFinalBondTermsIterator{contract: _Chat.contract, event: "FinalBondTerms", logs: logs, sub: sub}, nil
}

// WatchFinalBondTerms is a free log subscription operation binding the contract event 0x7e5631fad24e4cbdaaf2d85082725b7cc657fb9d32a0025131310dca36d2eb41.
//
// Solidity: event FinalBondTerms(uint32 principal, uint8 couponRate, uint32 couponDate, uint32 maturityDate, uint8 currency)
func (_Chat *ChatFilterer) WatchFinalBondTerms(opts *bind.WatchOpts, sink chan<- *ChatFinalBondTerms) (event.Subscription, error) {

	logs, sub, err := _Chat.contract.WatchLogs(opts, "FinalBondTerms")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChatFinalBondTerms)
				if err := _Chat.contract.UnpackLog(event, "FinalBondTerms", log); err != nil {
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

// ParseFinalBondTerms is a log parse operation binding the contract event 0x7e5631fad24e4cbdaaf2d85082725b7cc657fb9d32a0025131310dca36d2eb41.
//
// Solidity: event FinalBondTerms(uint32 principal, uint8 couponRate, uint32 couponDate, uint32 maturityDate, uint8 currency)
func (_Chat *ChatFilterer) ParseFinalBondTerms(log types.Log) (*ChatFinalBondTerms, error) {
	event := new(ChatFinalBondTerms)
	if err := _Chat.contract.UnpackLog(event, "FinalBondTerms", log); err != nil {
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
	Timestamp   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewBondCreated is a free log retrieval operation binding the contract event 0x74bc6db54501e99f138a7c95f5207903e8ff7c77a541e7599688bac39d534f96.
//
// Solidity: event NewBondCreated(address sender, address bondAddress, uint256 timestamp)
func (_Chat *ChatFilterer) FilterNewBondCreated(opts *bind.FilterOpts) (*ChatNewBondCreatedIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "NewBondCreated")
	if err != nil {
		return nil, err
	}
	return &ChatNewBondCreatedIterator{contract: _Chat.contract, event: "NewBondCreated", logs: logs, sub: sub}, nil
}

// WatchNewBondCreated is a free log subscription operation binding the contract event 0x74bc6db54501e99f138a7c95f5207903e8ff7c77a541e7599688bac39d534f96.
//
// Solidity: event NewBondCreated(address sender, address bondAddress, uint256 timestamp)
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

// ParseNewBondCreated is a log parse operation binding the contract event 0x74bc6db54501e99f138a7c95f5207903e8ff7c77a541e7599688bac39d534f96.
//
// Solidity: event NewBondCreated(address sender, address bondAddress, uint256 timestamp)
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
	Chat        ChatContractMessageInfo
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewChatMessage is a free log retrieval operation binding the contract event 0xe375e64f8bee7abfbf28664e9a49830e9bc77a3d866bb3e3bcc69a65dc054e24.
//
// Solidity: event NewChatMessage(address bondAddress, (address,string,uint256) chat)
func (_Chat *ChatFilterer) FilterNewChatMessage(opts *bind.FilterOpts) (*ChatNewChatMessageIterator, error) {

	logs, sub, err := _Chat.contract.FilterLogs(opts, "NewChatMessage")
	if err != nil {
		return nil, err
	}
	return &ChatNewChatMessageIterator{contract: _Chat.contract, event: "NewChatMessage", logs: logs, sub: sub}, nil
}

// WatchNewChatMessage is a free log subscription operation binding the contract event 0xe375e64f8bee7abfbf28664e9a49830e9bc77a3d866bb3e3bcc69a65dc054e24.
//
// Solidity: event NewChatMessage(address bondAddress, (address,string,uint256) chat)
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

// ParseNewChatMessage is a log parse operation binding the contract event 0xe375e64f8bee7abfbf28664e9a49830e9bc77a3d866bb3e3bcc69a65dc054e24.
//
// Solidity: event NewChatMessage(address bondAddress, (address,string,uint256) chat)
func (_Chat *ChatFilterer) ParseNewChatMessage(log types.Log) (*ChatNewChatMessage, error) {
	event := new(ChatNewChatMessage)
	if err := _Chat.contract.UnpackLog(event, "NewChatMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
