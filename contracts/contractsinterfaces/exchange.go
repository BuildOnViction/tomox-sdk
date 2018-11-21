// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractsinterfaces

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalTokenSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20Bin is the compiled bytecode used for deploying new contracts.
const ERC20Bin = `0x`

// DeployERC20 deploys a new Ethereum contract, binding an instance of ERC20 to it.
func DeployERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// ERC20 is an auto generated Go binding around an Ethereum contract.
type ERC20 struct {
	ERC20Caller     // Read-only binding to the contract
	ERC20Transactor // Write-only binding to the contract
	ERC20Filterer   // Log filterer for contract events
}

// ERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20Session struct {
	Contract     *ERC20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20CallerSession struct {
	Contract *ERC20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20TransactorSession struct {
	Contract     *ERC20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20Raw struct {
	Contract *ERC20 // Generic contract binding to access the raw methods on
}

// ERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20CallerRaw struct {
	Contract *ERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20TransactorRaw struct {
	Contract *ERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20 creates a new instance of ERC20, bound to a specific deployed contract.
func NewERC20(address common.Address, backend bind.ContractBackend) (*ERC20, error) {
	contract, err := bindERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// NewERC20Caller creates a new read-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Caller(address common.Address, caller bind.ContractCaller) (*ERC20Caller, error) {
	contract, err := bindERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Caller{contract: contract}, nil
}

// NewERC20Transactor creates a new write-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC20Transactor, error) {
	contract, err := bindERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Transactor{contract: contract}, nil
}

// NewERC20Filterer creates a new log filterer instance of ERC20, bound to a specific deployed contract.
func NewERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC20Filterer, error) {
	contract, err := bindERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20Filterer{contract: contract}, nil
}

// bindERC20 binds a generic wrapper to an already deployed contract.
func bindERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.ERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_ERC20 *ERC20Caller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_ERC20 *ERC20Session) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_ERC20 *ERC20CallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ERC20 *ERC20Session) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ERC20 *ERC20CallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20 *ERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20 *ERC20Session) Decimals() (uint8, error) {
	return _ERC20.Contract.Decimals(&_ERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20 *ERC20CallerSession) Decimals() (uint8, error) {
	return _ERC20.Contract.Decimals(&_ERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20 *ERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20 *ERC20Session) Name() (string, error) {
	return _ERC20.Contract.Name(&_ERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20 *ERC20CallerSession) Name() (string, error) {
	return _ERC20.Contract.Name(&_ERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20 *ERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20 *ERC20Session) Symbol() (string, error) {
	return _ERC20.Contract.Symbol(&_ERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20 *ERC20CallerSession) Symbol() (string, error) {
	return _ERC20.Contract.Symbol(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Session) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalTokenSupply is a free data retrieval call binding the contract method 0x1ca8b6cb.
//
// Solidity: function totalTokenSupply() constant returns(uint256)
func (_ERC20 *ERC20Caller) TotalTokenSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "totalTokenSupply")
	return *ret0, err
}

// TotalTokenSupply is a free data retrieval call binding the contract method 0x1ca8b6cb.
//
// Solidity: function totalTokenSupply() constant returns(uint256)
func (_ERC20 *ERC20Session) TotalTokenSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalTokenSupply(&_ERC20.CallOpts)
}

// TotalTokenSupply is a free data retrieval call binding the contract method 0x1ca8b6cb.
//
// Solidity: function totalTokenSupply() constant returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalTokenSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalTokenSupply(&_ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Session) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_ERC20 *ERC20TransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Session) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20TransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Session) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20TransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, _from, _to, _value)
}

// ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20 contract.
type ERC20ApprovalIterator struct {
	Event *ERC20Approval // Event containing the contract specifics and raw log

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
func (it *ERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Approval)
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
		it.Event = new(ERC20Approval)
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
func (it *ERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Approval represents a Approval event raised by the ERC20 contract.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_ERC20 *ERC20Filterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*ERC20ApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ApprovalIterator{contract: _ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_ERC20 *ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20Approval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Approval)
				if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20 contract.
type ERC20TransferIterator struct {
	Event *ERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Transfer)
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
		it.Event = new(ERC20Transfer)
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
func (it *ERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_ERC20 *ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*ERC20TransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20TransferIterator{contract: _ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_ERC20 *ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20Transfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Transfer)
				if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ExchangeABI is the input ABI used to generate the binding from.
const ExchangeABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"rewardAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[10]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8[2]\"},{\"name\":\"rs\",\"type\":\"bytes32[4]\"}],\"name\":\"executeSingleTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"numerator\",\"type\":\"uint256\"},{\"name\":\"denominator\",\"type\":\"uint256\"},{\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"isRoundingError\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[10]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"v\",\"type\":\"uint8[2]\"},{\"name\":\"rs\",\"type\":\"bytes32[4]\"}],\"name\":\"validateSignatures\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_baseToken\",\"type\":\"address\"},{\"name\":\"_quoteToken\",\"type\":\"address\"},{\"name\":\"_pricepointMultiplier\",\"type\":\"uint256\"}],\"name\":\"registerPair\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_rewardAccount\",\"type\":\"address\"}],\"name\":\"setFeeAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"wethToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[10][]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4][]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"name\":\"v\",\"type\":\"uint8[2][]\"},{\"name\":\"rs\",\"type\":\"bytes32[4][]\"}],\"name\":\"executeBatchTrades\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"},{\"name\":\"_isOperator\",\"type\":\"bool\"}],\"name\":\"setOperator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"pairs\",\"outputs\":[{\"name\":\"pairID\",\"type\":\"bytes32\"},{\"name\":\"baseToken\",\"type\":\"address\"},{\"name\":\"quoteToken\",\"type\":\"address\"},{\"name\":\"pricepointMultiplier\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"},{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_wethToken\",\"type\":\"address\"}],\"name\":\"setWethToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"makerOrderHashes\",\"type\":\"bytes32[]\"},{\"name\":\"takerOrderHashes\",\"type\":\"bytes32[]\"}],\"name\":\"emitLog\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"numerator\",\"type\":\"uint256\"},{\"name\":\"denominator\",\"type\":\"uint256\"},{\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"getPartialAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_baseToken\",\"type\":\"address\"},{\"name\":\"_quoteToken\",\"type\":\"address\"}],\"name\":\"getPairPricepointMultiplier\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[10]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"pricepointMultiplier\",\"type\":\"uint256\"}],\"name\":\"executeTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"traded\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[6]\"},{\"name\":\"orderAddresses\",\"type\":\"address[3]\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[6][]\"},{\"name\":\"orderAddresses\",\"type\":\"address[3][]\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_baseToken\",\"type\":\"address\"},{\"name\":\"_quoteToken\",\"type\":\"address\"}],\"name\":\"pairIsRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_wethToken\",\"type\":\"address\"},{\"name\":\"_rewardAccount\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"oldWethToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newWethToken\",\"type\":\"address\"}],\"name\":\"LogWethTokenUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"oldRewardAccount\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newRewardAccount\",\"type\":\"address\"}],\"name\":\"LogRewardAccountUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isOperator\",\"type\":\"bool\"}],\"name\":\"LogOperatorUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"makerOrderHashes\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"name\":\"takerOrderHashes\",\"type\":\"bytes32[]\"},{\"indexed\":true,\"name\":\"tokenPairHash\",\"type\":\"bytes32\"}],\"name\":\"LogBatchTrades\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"filledAmountSell\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"filledAmountBuy\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"paidFeeMake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"paidFeeTake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"tradeHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"tokenPairHash\",\"type\":\"bytes32\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"errorId\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"makerOrderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"takerOrderHash\",\"type\":\"bytes32\"}],\"name\":\"LogError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"userAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"baseToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"quoteToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"pricepoint\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"side\",\"type\":\"uint256\"}],\"name\":\"LogCancelOrder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"}]"

// ExchangeBin is the compiled bytecode used for deploying new contracts.
const ExchangeBin = "0x608060405234801561001057600080fd5b5060405160208061408183398101806040528101908080519060200190929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050613fbd806100c46000396000f300608060405260043610610133576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630e7082031461013857806310ac00d81461018f57806313af40351461029157806313e7c9d8146102d457806314df96ee1461032f5780631778baf414610388578063288cdc91146104805780633c918341146104c55780634b023cf81461054a5780635171267f146105a5578063558a729714610820578063673e0481146108875780638163681e146109415780638da5cb5b146109d357806393c1ae0914610a2a57806398024a8b14610b025780639e5ebf5e14610b57578063b4cb255314610bce578063d581332314610c9a578063d9a72b5214610ce3578063e51ad32d14610da5578063f4a8726314610f8f578063ffa1ad741461100a575b600080fd5b34801561014457600080fd5b5061014d61109a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561019b57600080fd5b50610277600480360381019080806101400190600a806020026040519081016040528092919082600a6020028082843782019150505050509192919290806080019060048060200260405190810160405280929190826004602002808284378201915050505050919291929080359060200190929190806040019060028060200260405190810160405280929190826002602002808284378201915050505050919291929080608001906004806020026040519081016040528092919082600460200280828437820191505050505091929192905050506110c0565b604051808215151515815260200191505060405180910390f35b34801561029d57600080fd5b506102d2600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506111c4565b005b3480156102e057600080fd5b50610315600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506112dd565b604051808215151515815260200191505060405180910390f35b34801561033b57600080fd5b5061036e6004803603810190808035906020019092919080359060200190929190803590602001909291905050506112fd565b604051808215151515815260200191505060405180910390f35b34801561039457600080fd5b50610466600480360381019080806101400190600a806020026040519081016040528092919082600a6020028082843782019150505050509192919290806080019060048060200260405190810160405280929190826004602002808284378201915050505050919291929080604001906002806020026040519081016040528092919082600260200280828437820191505050505091929192908060800190600480602002604051908101604052809291908260046020028082843782019150505050509192919290505050611370565b604051808215151515815260200191505060405180910390f35b34801561048c57600080fd5b506104af6004803603810190808035600019169060200190929190505050611790565b6040518082815260200191505060405180910390f35b3480156104d157600080fd5b50610530600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506117a8565b604051808215151515815260200191505060405180910390f35b34801561055657600080fd5b5061058b600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611932565b604051808215151515815260200191505060405180910390f35b3480156105b157600080fd5b5061080660048036038101908080359060200190820180359060200190808060200260200160405190810160405280939291908181526020016000905b828210156106345784848390506101400201600a806020026040519081016040528092919082600a602002808284378201915050505050815260200190600101906105ee565b5050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020016000905b828210156106b357848483905060800201600480602002604051908101604052809291908260046020028082843782019150505050508152602001906001019061066e565b505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020016000905b82821015610775578484839050604002016002806020026040519081016040528092919082600260200280828437820191505050505081526020019060010190610730565b5050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020016000905b828210156107f45784848390506080020160048060200260405190810160405280929190826004602002808284378201915050505050815260200190600101906107af565b50505050509192919290505050611ace565b604051808215151515815260200191505060405180910390f35b34801561082c57600080fd5b5061086d600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803515159060200190929190505050611dc5565b604051808215151515815260200191505060405180910390f35b34801561089357600080fd5b506108b66004803603810190808035600019169060200190929190505050611f2e565b6040518085600019166000191681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200194505050505060405180910390f35b34801561094d57600080fd5b506109b9600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035600019169060200190929190803560ff16906020019092919080356000191690602001909291908035600019169060200190929190505050611f9e565b604051808215151515815260200191505060405180910390f35b3480156109df57600080fd5b506109e861210b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b348015610a3657600080fd5b50610b0060048036038101908080608001906004806020026040519081016040528092919082600460200280828437820191505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509192919290505050612130565b005b348015610b0e57600080fd5b50610b41600480360381019080803590602001909291908035906020019092919080359060200190929190505050612320565b6040518082815260200191505060405180910390f35b348015610b6357600080fd5b50610bb8600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612350565b6040518082815260200191505060405180910390f35b348015610bda57600080fd5b50610c62600480360381019080806101400190600a806020026040519081016040528092919082600a602002808284378201915050505050919291929080608001906004806020026040519081016040528092919082600460200280828437820191505050505091929192908035906020019092919080359060200190929190505050612387565b604051808460001916600019168152602001836000191660001916815260200182151515158152602001935050505060405180910390f35b348015610ca657600080fd5b50610cc9600480360381019080803560001916906020019092919050505061321f565b604051808215151515815260200191505060405180910390f35b348015610cef57600080fd5b50610d8b6004803603810190808060c001906006806020026040519081016040528092919082600660200280828437820191505050505091929192908060600190600380602002604051908101604052809291908260036020028082843782019150505050509192919290803560ff1690602001909291908035600019169060200190929190803560001916906020019092919050505061323f565b604051808215151515815260200191505060405180910390f35b348015610db157600080fd5b50610f8d60048036038101908080359060200190820180359060200190808060200260200160405190810160405280939291908181526020016000905b82821015610e3357848483905060c002016006806020026040519081016040528092919082600660200280828437820191505050505081526020019060010190610dee565b5050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020016000905b82821015610eb2578484839050606002016003806020026040519081016040528092919082600360200280828437820191505050505081526020019060010190610e6d565b5050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509192919290803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929050505061353d565b005b348015610f9b57600080fd5b50610ff0600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506135e2565b604051808215151515815260200191505060405180910390f35b34801561101657600080fd5b5061101f61362d565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561105f578082015181840152602081019050611044565b50505050905090810190601f16801561108c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806111695750600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff165b151561117457600080fd5b61118088888787611370565b915081151561119257600092506111b9565b61119b87613666565b90506111a988888884612387565b5050506111b78888886137a2565b505b505095945050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561121f57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c66360405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60026020528060005260406000206000915054906101000a900460ff1681565b60008060008480151561130c57fe5b868509915060008214156113235760009250611367565b61135d611339858861395b90919063ffffffff16565b61134f620f42408561395b90919063ffffffff16565b61398e90919063ffffffff16565b90506103e8811192505b50509392505050565b600061137a613eaa565b611382613eaa565b600080610120604051908101604052808960006004811015156113a157fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018960026004811015156113d057fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018960036004811015156113ff57fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018a6000600a8110151561142e57fe5b602002015181526020018a6001600a8110151561144757fe5b602002015181526020018a6002600a8110151561146057fe5b602002015181526020018a6003600a8110151561147957fe5b602002015181526020018a6008600a8110151561149257fe5b602002015181526020018a6009600a811015156114ab57fe5b60200201518152509350610120604051908101604052808960016004811015156114d157fe5b602002015173ffffffffffffffffffffffffffffffffffffffff16815260200189600260048110151561150057fe5b602002015173ffffffffffffffffffffffffffffffffffffffff16815260200189600360048110151561152f57fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018a6004600a8110151561155e57fe5b602002015181526020018a6005600a8110151561157757fe5b602002015181526020018a6006600a8110151561159057fe5b602002015181526020018a6007600a811015156115a957fe5b602002015181526020018a6009600a811015156115c257fe5b602002015181526020018a6008600a811015156115db57fe5b602002015181525092506115ee846139a9565b91506115f9836139a9565b905061164584600001518389600060028110151561161357fe5b602002015189600060048110151561162757fe5b60200201518a600160048110151561163b57fe5b6020020151611f9e565b15156116bd577f14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb6001600881111561167957fe5b8383604051808460ff1660ff16815260200183600019166000191681526020018260001916600019168152602001935050505060405180910390a160009450611784565b6117078360000151828960016002811015156116d557fe5b60200201518960026004811015156116e957fe5b60200201518a60036004811015156116fd57fe5b6020020151611f9e565b151561177f577f14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb6002600881111561173b57fe5b8383604051808460ff1660ff16815260200183600019166000191681526020018260001916600019168152602001935050505060405180910390a160009450611784565b600194505b50505050949350505050565b60036020528060005260406000206000915090505481565b6000806000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561180657600080fd5b6118108585613b8d565b9050608060405190810160405280826000191681526020018673ffffffffffffffffffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff168152602001848152506005600083600019166000191681526020019081526020016000206000820151816000019060001916905560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060608201518160030155905050509392505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561198f57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141515156119cb57600080fd5b7f18d40614e4a77383f4b7337227bdad137b4f3f9b002ef63afd3ddaa142a15f63600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a181600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060019050919050565b60006060806000806000806000806000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480611b805750600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff165b1515611b8b57600080fd5b8c51604051908082528060200260200182016040528015611bbb5781602001602082028038833980820191505090505b5097508c51604051908082528060200260200182016040528015611bee5781602001602082028038833980820191505090505b509650600095505b8c51861015611d5457611c678e87815181101515611c1057fe5b906020019060200201518e88815181101515611c2857fe5b906020019060200201518d89815181101515611c4057fe5b906020019060200201518d8a815181101515611c5857fe5b90602001906020020151611370565b9450841515611c795760009850611db4565b611c998d87815181101515611c8a57fe5b90602001906020020151613666565b9350611cec8e87815181101515611cac57fe5b906020019060200201518e88815181101515611cc457fe5b906020019060200201518e89815181101515611cdc57fe5b9060200190602002015187612387565b9250925092508015611d4757828887815181101515611d0757fe5b906020019060200201906000191690816000191681525050818787815181101515611d2e57fe5b9060200190602002019060001916908160001916815250505b8580600101965050611bf6565b611d8f8e6000815181101515611d6657fe5b906020019060200201518e6000815181101515611d7f57fe5b906020019060200201518e613c97565b50611db38d6000815181101515611da257fe5b906020019060200201518989612130565b5b505050505050505095945050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515611e2257600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515611e5e57600080fd5b7f4af650e9ee9ac50b37ec2cd3ddac7e1c69955ffc871bc3e812563775f3bc0e7d8383604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001821515151581526020019250505060405180910390a181600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001905092915050565b60056020528060005260406000206000915090508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060030154905084565b600060018560405160200180807f19457468657265756d205369676e6564204d6573736167653a0a333200000000815250601c0182600019166000191681526020019150506040516020818303038152906040526040518082805190602001908083835b6020831015156120275780518252602082019150602081019050602083039250612002565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040518091039020858585604051600081526020016040526040518085600019166000191681526020018460ff1660ff1681526020018360001916600019168152602001826000191660001916815260200194505050505060206040516020810390808403906000865af11580156120c8573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1614905095945050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b82600160048110151561213f57fe5b602002015183600260048110151561215357fe5b6020020151604051602001808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c010000000000000000000000000281526014018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c01000000000000000000000000028152601401925050506040516020818303038152906040526040518082805190602001908083835b6020831015156122295780518252602082019150602081019050602083039250612204565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040518091039020600019167fde8acabe30c9bd25d65bb9db28bf46f51dc7500a07b1671f121f1144fbf446fc8383604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156122c45780820151818401526020810190506122a9565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156123065780820151818401526020810190506122eb565b5050505090500194505050505060405180910390a2505050565b600061234783612339848761395b90919063ffffffff16565b61398e90919063ffffffff16565b90509392505050565b60008061235d8484613b8d565b90506005600082600019166000191681526020019081526020016000206003015491505092915050565b6000806000612394613eaa565b61239c613eaa565b60008060008060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806124485750600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff165b151561245357600080fd5b610120604051908101604052808e600060048110151561246f57fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018e600260048110151561249e57fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018e60036004811015156124cd57fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018f6000600a811015156124fc57fe5b602002015181526020018f6001600a8110151561251557fe5b602002015181526020018f6002600a8110151561252e57fe5b602002015181526020018f6003600a8110151561254757fe5b602002015181526020018f6008600a8110151561256057fe5b602002015181526020018f6009600a8110151561257957fe5b60200201518152509650610120604051908101604052808e600160048110151561259f57fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018e60026004811015156125ce57fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018e60036004811015156125fd57fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018f6004600a8110151561262c57fe5b602002015181526020018f6005600a8110151561264557fe5b602002015181526020018f6006600a8110151561265e57fe5b602002015181526020018f6007600a8110151561267757fe5b602002015181526020018f6009600a8110151561269057fe5b602002015181526020018f6008600a811015156126a957fe5b602002015181525095506126bc876139a9565b94506126c7866139a9565b935086606001516126fc8d60036000896000191660001916815260200190815260200160002054613e8c90919063ffffffff16565b111561277a577f14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb6007600881111561273057fe5b8686604051808460ff1660ff16815260200183600019166000191681526020018260001916600019168152602001935050505060405180910390a18484600099509950995061320e565b85606001516127ad8d60036000886000191660001916815260200190815260200160002054613e8c90919063ffffffff16565b111561282b577f14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb600760088111156127e157fe5b8686604051808460ff1660ff16815260200183600019166000191681526020018260001916600019168152602001935050505060405180910390a18484600099509950995061320e565b8660a001518660a0015114156128b3577f14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb6003600881111561286957fe5b8686604051808460ff1660ff16815260200183600019166000191681526020018260001916600019168152602001935050505060405180910390a18484600099509950995061320e565b60008760a00151141561294957856080015187608001511015612948577f14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb600460088111156128fe57fe5b8686604051808460ff1660ff16815260200183600019166000191681526020018260001916600019168152602001935050505060405180910390a18484600099509950995061320e565b5b60018760a0015114156129df578660800151866080015110156129de577f14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb6004600881111561299457fe5b8686604051808460ff1660ff16815260200183600019166000191681526020018260001916600019168152602001935050505060405180910390a18484600099509950995061320e565b5b612a0d8c60036000876000191660001916815260200190815260200160002054613e8c90919063ffffffff16565b60036000866000191660001916815260200190815260200160002081905550612a5a8c60036000886000191660001916815260200190815260200160002054613e8c90919063ffffffff16565b600360008760001916600019168152602001908152602001600020819055508b9250612aa58b612a9789608001518f61395b90919063ffffffff16565b61398e90919063ffffffff16565b9150612aba8c88606001518960e00151612320565b905060008760a001511415612e6757866040015173ffffffffffffffffffffffffffffffffffffffff166323b872dd88600001518860000151856040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b158015612bac57600080fd5b505af1158015612bc0573d6000803e3d6000fd5b505050506040513d6020811015612bd657600080fd5b81019080805190602001909291905050501515612bf257600080fd5b866040015173ffffffffffffffffffffffffffffffffffffffff166323b872dd8860000151600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b158015612cf357600080fd5b505af1158015612d07573d6000803e3d6000fd5b505050506040513d6020811015612d1d57600080fd5b81019080805190602001909291905050501515612d3957600080fd5b856020015173ffffffffffffffffffffffffffffffffffffffff166323b872dd87600001518960000151866040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b158015612e1c57600080fd5b505af1158015612e30573d6000803e3d6000fd5b505050506040513d6020811015612e4657600080fd5b81019080805190602001909291905050501515612e6257600080fd5b613203565b866020015173ffffffffffffffffffffffffffffffffffffffff166323b872dd88600001518860000151866040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b158015612f4a57600080fd5b505af1158015612f5e573d6000803e3d6000fd5b505050506040513d6020811015612f7457600080fd5b81019080805190602001909291905050501515612f9057600080fd5b856040015173ffffffffffffffffffffffffffffffffffffffff166323b872dd8760000151600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b15801561309157600080fd5b505af11580156130a5573d6000803e3d6000fd5b505050506040513d60208110156130bb57600080fd5b810190808051906020019092919050505015156130d757600080fd5b856040015173ffffffffffffffffffffffffffffffffffffffff166323b872dd876000015189600001518486036040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b1580156131bc57600080fd5b505af11580156131d0573d6000803e3d6000fd5b505050506040513d60208110156131e657600080fd5b8101908080519060200190929190505050151561320257600080fd5b5b848460019950995099505b505050505050509450945094915050565b60046020528060005260406000206000915054906101000a900460ff1681565b6000613249613eaa565b60006101206040519081016040528088600060038110151561326757fe5b602002015173ffffffffffffffffffffffffffffffffffffffff16815260200188600160038110151561329657fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018860026003811015156132c557fe5b602002015173ffffffffffffffffffffffffffffffffffffffff1681526020018960006006811015156132f457fe5b6020020151815260200189600160068110151561330d57fe5b6020020151815260200189600260068110151561332657fe5b6020020151815260200189600360068110151561333f57fe5b6020020151815260200189600560068110151561335857fe5b6020020151815260200189600460068110151561337157fe5b60200201518152509150613384826139a9565b90506133933382888888611f9e565b15156133fe577f14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb600060088111156133c757fe5b82604051808360ff1660ff16815260200182600019166000191681526020016020019250505060405180910390a160009250613532565b8160600151600360008360001916600019168152602001908152602001600020819055507fb00984fe824f4973f31e8a414157f54cb4ee29bc2100149ba22a094d0bfd551881836000015184602001518560400151866060015187608001518860a001516040518088600019166000191681526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200184815260200183815260200182815260200197505050505050505060405180910390a1600192505b505095945050505050565b60008090505b84518110156135da576135cc868281518110151561355d57fe5b90602001906020020151868381518110151561357557fe5b90602001906020020151868481518110151561358d57fe5b9060200190602002015186858151811015156135a557fe5b9060200190602002015186868151811015156135bd57fe5b9060200190602002015161323f565b508080600101915050613543565b505050505050565b6000806135ef8484613b8d565b905060006005600083600019166000191681526020019081526020016000206003015414156136215760009150613626565b600191505b5092915050565b6040805190810160405280600581526020017f312e302e3000000000000000000000000000000000000000000000000000000081525081565b600080613671613f39565b6136a184600260048110151561368357fe5b602002015185600360048110151561369757fe5b6020020151613b8d565b9150600560008360001916600019168152602001908152602001600020608060405190810160405290816000820154600019166000191681526020016001820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016003820154815250509050806060015192505050919050565b600080600080600080886004600a811015156137ba57fe5b60200201519450886008600a811015156137d057fe5b602002015193508760016004811015156137e657fe5b602002015192508760036004811015156137fc57fe5b6020020151915061380e878686612320565b90508173ffffffffffffffffffffffffffffffffffffffff166323b872dd84600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b15801561390957600080fd5b505af115801561391d573d6000803e3d6000fd5b505050506040513d602081101561393357600080fd5b8101908080519060200190929190505050151561394f57600080fd5b50505050509392505050565b6000808284029050600084148061397c575082848281151561397957fe5b04145b151561398457fe5b8091505092915050565b600080828481151561399c57fe5b0490508091505092915050565b600030826000015183602001518460400151856060015186608001518760a001518860c001518960e001518a6101000151604051602001808b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c010000000000000000000000000281526014018a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c010000000000000000000000000281526014018973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c010000000000000000000000000281526014018873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c010000000000000000000000000281526014018781526020018681526020018581526020018481526020018381526020018281526020019a50505050505050505050506040516020818303038152906040526040518082805190602001908083835b602083101515613b595780518252602082019150602081019050602083039250613b34565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390209050919050565b60008282604051602001808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c010000000000000000000000000281526014018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c01000000000000000000000000028152601401925050506040516020818303038152906040526040518082805190602001908083835b602083101515613c625780518252602082019150602081019050602083039250613c3d565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040518091039020905092915050565b6000806000806000806000808a6004600a81101515613cb257fe5b602002015196508a6008600a81101515613cc857fe5b60200201519550896001600481101515613cde57fe5b60200201519450896003600481101515613cf457fe5b60200201519350600091505b8851821015613d32578882815181101515613d1757fe5b90602001906020020151830192508180600101925050613d00565b613d3d838888612320565b90508373ffffffffffffffffffffffffffffffffffffffff166323b872dd86600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b158015613e3857600080fd5b505af1158015613e4c573d6000803e3d6000fd5b505050506040513d6020811015613e6257600080fd5b81019080805190602001909291905050501515613e7e57600080fd5b505050505050509392505050565b6000808284019050838110151515613ea057fe5b8091505092915050565b61012060405190810160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff1681526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b60806040519081016040528060008019168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff1681526020016000815250905600a165627a7a72305820e6cdb27cd2482f21469dfea0737440db5b391c7ec7feecaff4747ad06e3768f50029"


// DeployExchange deploys a new Ethereum contract, binding an instance of Exchange to it.
func DeployExchange(auth *bind.TransactOpts, backend bind.ContractBackend, _wethToken common.Address, _rewardAccount common.Address) (common.Address, *types.Transaction, *Exchange, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ExchangeBin), backend, _wethToken, _rewardAccount)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Exchange{ExchangeCaller: ExchangeCaller{contract: contract}, ExchangeTransactor: ExchangeTransactor{contract: contract}, ExchangeFilterer: ExchangeFilterer{contract: contract}}, nil
}

// Exchange is an auto generated Go binding around an Ethereum contract.
type Exchange struct {
	ExchangeCaller     // Read-only binding to the contract
	ExchangeTransactor // Write-only binding to the contract
	ExchangeFilterer   // Log filterer for contract events
}

// ExchangeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExchangeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangeSession struct {
	Contract     *Exchange         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExchangeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangeCallerSession struct {
	Contract *ExchangeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ExchangeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangeTransactorSession struct {
	Contract     *ExchangeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ExchangeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangeRaw struct {
	Contract *Exchange // Generic contract binding to access the raw methods on
}

// ExchangeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangeCallerRaw struct {
	Contract *ExchangeCaller // Generic read-only contract binding to access the raw methods on
}

// ExchangeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangeTransactorRaw struct {
	Contract *ExchangeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExchange creates a new instance of Exchange, bound to a specific deployed contract.
func NewExchange(address common.Address, backend bind.ContractBackend) (*Exchange, error) {
	contract, err := bindExchange(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Exchange{ExchangeCaller: ExchangeCaller{contract: contract}, ExchangeTransactor: ExchangeTransactor{contract: contract}, ExchangeFilterer: ExchangeFilterer{contract: contract}}, nil
}

// NewExchangeCaller creates a new read-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeCaller(address common.Address, caller bind.ContractCaller) (*ExchangeCaller, error) {
	contract, err := bindExchange(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeCaller{contract: contract}, nil
}

// NewExchangeTransactor creates a new write-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangeTransactor, error) {
	contract, err := bindExchange(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeTransactor{contract: contract}, nil
}

// NewExchangeFilterer creates a new log filterer instance of Exchange, bound to a specific deployed contract.
func NewExchangeFilterer(address common.Address, filterer bind.ContractFilterer) (*ExchangeFilterer, error) {
	contract, err := bindExchange(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExchangeFilterer{contract: contract}, nil
}

// bindExchange binds a generic wrapper to an already deployed contract.
func bindExchange(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchange *ExchangeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchange.Contract.ExchangeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchange *ExchangeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.Contract.ExchangeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchange *ExchangeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchange.Contract.ExchangeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchange *ExchangeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchange.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchange *ExchangeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchange *ExchangeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchange.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeSession) VERSION() (string, error) {
	return _Exchange.Contract.VERSION(&_Exchange.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeCallerSession) VERSION() (string, error) {
	return _Exchange.Contract.VERSION(&_Exchange.CallOpts)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCaller) Filled(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "filled", arg0)
	return *ret0, err
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// GetPairPricepointMultiplier is a free data retrieval call binding the contract method 0x9e5ebf5e.
//
// Solidity: function getPairPricepointMultiplier(_baseToken address, _quoteToken address) constant returns(uint256)
func (_Exchange *ExchangeCaller) GetPairPricepointMultiplier(opts *bind.CallOpts, _baseToken common.Address, _quoteToken common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getPairPricepointMultiplier", _baseToken, _quoteToken)
	return *ret0, err
}

// GetPairPricepointMultiplier is a free data retrieval call binding the contract method 0x9e5ebf5e.
//
// Solidity: function getPairPricepointMultiplier(_baseToken address, _quoteToken address) constant returns(uint256)
func (_Exchange *ExchangeSession) GetPairPricepointMultiplier(_baseToken common.Address, _quoteToken common.Address) (*big.Int, error) {
	return _Exchange.Contract.GetPairPricepointMultiplier(&_Exchange.CallOpts, _baseToken, _quoteToken)
}

// GetPairPricepointMultiplier is a free data retrieval call binding the contract method 0x9e5ebf5e.
//
// Solidity: function getPairPricepointMultiplier(_baseToken address, _quoteToken address) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) GetPairPricepointMultiplier(_baseToken common.Address, _quoteToken common.Address) (*big.Int, error) {
	return _Exchange.Contract.GetPairPricepointMultiplier(&_Exchange.CallOpts, _baseToken, _quoteToken)
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeCaller) GetPartialAmount(opts *bind.CallOpts, numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getPartialAmount", numerator, denominator, target)
	return *ret0, err
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeSession) GetPartialAmount(numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	return _Exchange.Contract.GetPartialAmount(&_Exchange.CallOpts, numerator, denominator, target)
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) GetPartialAmount(numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	return _Exchange.Contract.GetPartialAmount(&_Exchange.CallOpts, numerator, denominator, target)
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeCaller) IsRoundingError(opts *bind.CallOpts, numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "isRoundingError", numerator, denominator, target)
	return *ret0, err
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeSession) IsRoundingError(numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	return _Exchange.Contract.IsRoundingError(&_Exchange.CallOpts, numerator, denominator, target)
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeCallerSession) IsRoundingError(numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	return _Exchange.Contract.IsRoundingError(&_Exchange.CallOpts, numerator, denominator, target)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeCaller) IsValidSignature(opts *bind.CallOpts, signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "isValidSignature", signer, hash, v, r, s)
	return *ret0, err
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeSession) IsValidSignature(signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, signer, hash, v, r, s)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeCallerSession) IsValidSignature(signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, signer, hash, v, r, s)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Exchange *ExchangeCaller) Operators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "operators", arg0)
	return *ret0, err
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Exchange *ExchangeSession) Operators(arg0 common.Address) (bool, error) {
	return _Exchange.Contract.Operators(&_Exchange.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Exchange *ExchangeCallerSession) Operators(arg0 common.Address) (bool, error) {
	return _Exchange.Contract.Operators(&_Exchange.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeSession) Owner() (common.Address, error) {
	return _Exchange.Contract.Owner(&_Exchange.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeCallerSession) Owner() (common.Address, error) {
	return _Exchange.Contract.Owner(&_Exchange.CallOpts)
}

// PairIsRegistered is a free data retrieval call binding the contract method 0xf4a87263.
//
// Solidity: function pairIsRegistered(_baseToken address, _quoteToken address) constant returns(bool)
func (_Exchange *ExchangeCaller) PairIsRegistered(opts *bind.CallOpts, _baseToken common.Address, _quoteToken common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "pairIsRegistered", _baseToken, _quoteToken)
	return *ret0, err
}

// PairIsRegistered is a free data retrieval call binding the contract method 0xf4a87263.
//
// Solidity: function pairIsRegistered(_baseToken address, _quoteToken address) constant returns(bool)
func (_Exchange *ExchangeSession) PairIsRegistered(_baseToken common.Address, _quoteToken common.Address) (bool, error) {
	return _Exchange.Contract.PairIsRegistered(&_Exchange.CallOpts, _baseToken, _quoteToken)
}

// PairIsRegistered is a free data retrieval call binding the contract method 0xf4a87263.
//
// Solidity: function pairIsRegistered(_baseToken address, _quoteToken address) constant returns(bool)
func (_Exchange *ExchangeCallerSession) PairIsRegistered(_baseToken common.Address, _quoteToken common.Address) (bool, error) {
	return _Exchange.Contract.PairIsRegistered(&_Exchange.CallOpts, _baseToken, _quoteToken)
}

// Pairs is a free data retrieval call binding the contract method 0x673e0481.
//
// Solidity: function pairs( bytes32) constant returns(pairID bytes32, baseToken address, quoteToken address, pricepointMultiplier uint256)
func (_Exchange *ExchangeCaller) Pairs(opts *bind.CallOpts, arg0 [32]byte) (struct {
	PairID               [32]byte
	BaseToken            common.Address
	QuoteToken           common.Address
	PricepointMultiplier *big.Int
}, error) {
	ret := new(struct {
		PairID               [32]byte
		BaseToken            common.Address
		QuoteToken           common.Address
		PricepointMultiplier *big.Int
	})
	out := ret
	err := _Exchange.contract.Call(opts, out, "pairs", arg0)
	return *ret, err
}

// Pairs is a free data retrieval call binding the contract method 0x673e0481.
//
// Solidity: function pairs( bytes32) constant returns(pairID bytes32, baseToken address, quoteToken address, pricepointMultiplier uint256)
func (_Exchange *ExchangeSession) Pairs(arg0 [32]byte) (struct {
	PairID               [32]byte
	BaseToken            common.Address
	QuoteToken           common.Address
	PricepointMultiplier *big.Int
}, error) {
	return _Exchange.Contract.Pairs(&_Exchange.CallOpts, arg0)
}

// Pairs is a free data retrieval call binding the contract method 0x673e0481.
//
// Solidity: function pairs( bytes32) constant returns(pairID bytes32, baseToken address, quoteToken address, pricepointMultiplier uint256)
func (_Exchange *ExchangeCallerSession) Pairs(arg0 [32]byte) (struct {
	PairID               [32]byte
	BaseToken            common.Address
	QuoteToken           common.Address
	PricepointMultiplier *big.Int
}, error) {
	return _Exchange.Contract.Pairs(&_Exchange.CallOpts, arg0)
}

// RewardAccount is a free data retrieval call binding the contract method 0x0e708203.
//
// Solidity: function rewardAccount() constant returns(address)
func (_Exchange *ExchangeCaller) RewardAccount(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "rewardAccount")
	return *ret0, err
}

// RewardAccount is a free data retrieval call binding the contract method 0x0e708203.
//
// Solidity: function rewardAccount() constant returns(address)
func (_Exchange *ExchangeSession) RewardAccount() (common.Address, error) {
	return _Exchange.Contract.RewardAccount(&_Exchange.CallOpts)
}

// RewardAccount is a free data retrieval call binding the contract method 0x0e708203.
//
// Solidity: function rewardAccount() constant returns(address)
func (_Exchange *ExchangeCallerSession) RewardAccount() (common.Address, error) {
	return _Exchange.Contract.RewardAccount(&_Exchange.CallOpts)
}

// Traded is a free data retrieval call binding the contract method 0xd5813323.
//
// Solidity: function traded( bytes32) constant returns(bool)
func (_Exchange *ExchangeCaller) Traded(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "traded", arg0)
	return *ret0, err
}

// Traded is a free data retrieval call binding the contract method 0xd5813323.
//
// Solidity: function traded( bytes32) constant returns(bool)
func (_Exchange *ExchangeSession) Traded(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Traded(&_Exchange.CallOpts, arg0)
}

// Traded is a free data retrieval call binding the contract method 0xd5813323.
//
// Solidity: function traded( bytes32) constant returns(bool)
func (_Exchange *ExchangeCallerSession) Traded(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Traded(&_Exchange.CallOpts, arg0)
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Exchange *ExchangeCaller) WethToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "wethToken")
	return *ret0, err
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Exchange *ExchangeSession) WethToken() (common.Address, error) {
	return _Exchange.Contract.WethToken(&_Exchange.CallOpts)
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Exchange *ExchangeCallerSession) WethToken() (common.Address, error) {
	return _Exchange.Contract.WethToken(&_Exchange.CallOpts)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xe51ad32d.
//
// Solidity: function batchCancelOrders(orderValues uint256[6][], orderAddresses address[3][], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactor) BatchCancelOrders(opts *bind.TransactOpts, orderValues [][6]*big.Int, orderAddresses [][3]common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchCancelOrders", orderValues, orderAddresses, v, r, s)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xe51ad32d.
//
// Solidity: function batchCancelOrders(orderValues uint256[6][], orderAddresses address[3][], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeSession) BatchCancelOrders(orderValues [][6]*big.Int, orderAddresses [][3]common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xe51ad32d.
//
// Solidity: function batchCancelOrders(orderValues uint256[6][], orderAddresses address[3][], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactorSession) BatchCancelOrders(orderValues [][6]*big.Int, orderAddresses [][3]common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd9a72b52.
//
// Solidity: function cancelOrder(orderValues uint256[6], orderAddresses address[3], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeTransactor) CancelOrder(opts *bind.TransactOpts, orderValues [6]*big.Int, orderAddresses [3]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelOrder", orderValues, orderAddresses, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd9a72b52.
//
// Solidity: function cancelOrder(orderValues uint256[6], orderAddresses address[3], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeSession) CancelOrder(orderValues [6]*big.Int, orderAddresses [3]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd9a72b52.
//
// Solidity: function cancelOrder(orderValues uint256[6], orderAddresses address[3], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeTransactorSession) CancelOrder(orderValues [6]*big.Int, orderAddresses [3]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// EmitLog is a paid mutator transaction binding the contract method 0x93c1ae09.
//
// Solidity: function emitLog(orderAddresses address[4], makerOrderHashes bytes32[], takerOrderHashes bytes32[]) returns()
func (_Exchange *ExchangeTransactor) EmitLog(opts *bind.TransactOpts, orderAddresses [4]common.Address, makerOrderHashes [][32]byte, takerOrderHashes [][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "emitLog", orderAddresses, makerOrderHashes, takerOrderHashes)
}

// EmitLog is a paid mutator transaction binding the contract method 0x93c1ae09.
//
// Solidity: function emitLog(orderAddresses address[4], makerOrderHashes bytes32[], takerOrderHashes bytes32[]) returns()
func (_Exchange *ExchangeSession) EmitLog(orderAddresses [4]common.Address, makerOrderHashes [][32]byte, takerOrderHashes [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.EmitLog(&_Exchange.TransactOpts, orderAddresses, makerOrderHashes, takerOrderHashes)
}

// EmitLog is a paid mutator transaction binding the contract method 0x93c1ae09.
//
// Solidity: function emitLog(orderAddresses address[4], makerOrderHashes bytes32[], takerOrderHashes bytes32[]) returns()
func (_Exchange *ExchangeTransactorSession) EmitLog(orderAddresses [4]common.Address, makerOrderHashes [][32]byte, takerOrderHashes [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.EmitLog(&_Exchange.TransactOpts, orderAddresses, makerOrderHashes, takerOrderHashes)
}

// ExecuteBatchTrades is a paid mutator transaction binding the contract method 0x5171267f.
//
// Solidity: function executeBatchTrades(orderValues uint256[10][], orderAddresses address[4][], amounts uint256[], v uint8[2][], rs bytes32[4][]) returns(bool)
func (_Exchange *ExchangeTransactor) ExecuteBatchTrades(opts *bind.TransactOpts, orderValues [][10]*big.Int, orderAddresses [][4]common.Address, amounts []*big.Int, v [][2]uint8, rs [][4][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "executeBatchTrades", orderValues, orderAddresses, amounts, v, rs)
}

// ExecuteBatchTrades is a paid mutator transaction binding the contract method 0x5171267f.
//
// Solidity: function executeBatchTrades(orderValues uint256[10][], orderAddresses address[4][], amounts uint256[], v uint8[2][], rs bytes32[4][]) returns(bool)
func (_Exchange *ExchangeSession) ExecuteBatchTrades(orderValues [][10]*big.Int, orderAddresses [][4]common.Address, amounts []*big.Int, v [][2]uint8, rs [][4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteBatchTrades(&_Exchange.TransactOpts, orderValues, orderAddresses, amounts, v, rs)
}

// ExecuteBatchTrades is a paid mutator transaction binding the contract method 0x5171267f.
//
// Solidity: function executeBatchTrades(orderValues uint256[10][], orderAddresses address[4][], amounts uint256[], v uint8[2][], rs bytes32[4][]) returns(bool)
func (_Exchange *ExchangeTransactorSession) ExecuteBatchTrades(orderValues [][10]*big.Int, orderAddresses [][4]common.Address, amounts []*big.Int, v [][2]uint8, rs [][4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteBatchTrades(&_Exchange.TransactOpts, orderValues, orderAddresses, amounts, v, rs)
}

// ExecuteSingleTrade is a paid mutator transaction binding the contract method 0x10ac00d8.
//
// Solidity: function executeSingleTrade(orderValues uint256[10], orderAddresses address[4], amount uint256, v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeTransactor) ExecuteSingleTrade(opts *bind.TransactOpts, orderValues [10]*big.Int, orderAddresses [4]common.Address, amount *big.Int, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "executeSingleTrade", orderValues, orderAddresses, amount, v, rs)
}

// ExecuteSingleTrade is a paid mutator transaction binding the contract method 0x10ac00d8.
//
// Solidity: function executeSingleTrade(orderValues uint256[10], orderAddresses address[4], amount uint256, v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeSession) ExecuteSingleTrade(orderValues [10]*big.Int, orderAddresses [4]common.Address, amount *big.Int, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteSingleTrade(&_Exchange.TransactOpts, orderValues, orderAddresses, amount, v, rs)
}

// ExecuteSingleTrade is a paid mutator transaction binding the contract method 0x10ac00d8.
//
// Solidity: function executeSingleTrade(orderValues uint256[10], orderAddresses address[4], amount uint256, v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeTransactorSession) ExecuteSingleTrade(orderValues [10]*big.Int, orderAddresses [4]common.Address, amount *big.Int, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteSingleTrade(&_Exchange.TransactOpts, orderValues, orderAddresses, amount, v, rs)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0xb4cb2553.
//
// Solidity: function executeTrade(orderValues uint256[10], orderAddresses address[4], amount uint256, pricepointMultiplier uint256) returns(bytes32, bytes32, bool)
func (_Exchange *ExchangeTransactor) ExecuteTrade(opts *bind.TransactOpts, orderValues [10]*big.Int, orderAddresses [4]common.Address, amount *big.Int, pricepointMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "executeTrade", orderValues, orderAddresses, amount, pricepointMultiplier)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0xb4cb2553.
//
// Solidity: function executeTrade(orderValues uint256[10], orderAddresses address[4], amount uint256, pricepointMultiplier uint256) returns(bytes32, bytes32, bool)
func (_Exchange *ExchangeSession) ExecuteTrade(orderValues [10]*big.Int, orderAddresses [4]common.Address, amount *big.Int, pricepointMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTrade(&_Exchange.TransactOpts, orderValues, orderAddresses, amount, pricepointMultiplier)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0xb4cb2553.
//
// Solidity: function executeTrade(orderValues uint256[10], orderAddresses address[4], amount uint256, pricepointMultiplier uint256) returns(bytes32, bytes32, bool)
func (_Exchange *ExchangeTransactorSession) ExecuteTrade(orderValues [10]*big.Int, orderAddresses [4]common.Address, amount *big.Int, pricepointMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTrade(&_Exchange.TransactOpts, orderValues, orderAddresses, amount, pricepointMultiplier)
}

// RegisterPair is a paid mutator transaction binding the contract method 0x3c918341.
//
// Solidity: function registerPair(_baseToken address, _quoteToken address, _pricepointMultiplier uint256) returns(bool)
func (_Exchange *ExchangeTransactor) RegisterPair(opts *bind.TransactOpts, _baseToken common.Address, _quoteToken common.Address, _pricepointMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "registerPair", _baseToken, _quoteToken, _pricepointMultiplier)
}

// RegisterPair is a paid mutator transaction binding the contract method 0x3c918341.
//
// Solidity: function registerPair(_baseToken address, _quoteToken address, _pricepointMultiplier uint256) returns(bool)
func (_Exchange *ExchangeSession) RegisterPair(_baseToken common.Address, _quoteToken common.Address, _pricepointMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.RegisterPair(&_Exchange.TransactOpts, _baseToken, _quoteToken, _pricepointMultiplier)
}

// RegisterPair is a paid mutator transaction binding the contract method 0x3c918341.
//
// Solidity: function registerPair(_baseToken address, _quoteToken address, _pricepointMultiplier uint256) returns(bool)
func (_Exchange *ExchangeTransactorSession) RegisterPair(_baseToken common.Address, _quoteToken common.Address, _pricepointMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.RegisterPair(&_Exchange.TransactOpts, _baseToken, _quoteToken, _pricepointMultiplier)
}

// SetFeeAccount is a paid mutator transaction binding the contract method 0x4b023cf8.
//
// Solidity: function setFeeAccount(_rewardAccount address) returns(bool)
func (_Exchange *ExchangeTransactor) SetFeeAccount(opts *bind.TransactOpts, _rewardAccount common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setFeeAccount", _rewardAccount)
}

// SetFeeAccount is a paid mutator transaction binding the contract method 0x4b023cf8.
//
// Solidity: function setFeeAccount(_rewardAccount address) returns(bool)
func (_Exchange *ExchangeSession) SetFeeAccount(_rewardAccount common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetFeeAccount(&_Exchange.TransactOpts, _rewardAccount)
}

// SetFeeAccount is a paid mutator transaction binding the contract method 0x4b023cf8.
//
// Solidity: function setFeeAccount(_rewardAccount address) returns(bool)
func (_Exchange *ExchangeTransactorSession) SetFeeAccount(_rewardAccount common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetFeeAccount(&_Exchange.TransactOpts, _rewardAccount)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(_operator address, _isOperator bool) returns(bool)
func (_Exchange *ExchangeTransactor) SetOperator(opts *bind.TransactOpts, _operator common.Address, _isOperator bool) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setOperator", _operator, _isOperator)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(_operator address, _isOperator bool) returns(bool)
func (_Exchange *ExchangeSession) SetOperator(_operator common.Address, _isOperator bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetOperator(&_Exchange.TransactOpts, _operator, _isOperator)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(_operator address, _isOperator bool) returns(bool)
func (_Exchange *ExchangeTransactorSession) SetOperator(_operator common.Address, _isOperator bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetOperator(&_Exchange.TransactOpts, _operator, _isOperator)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Exchange *ExchangeTransactor) SetOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setOwner", newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Exchange *ExchangeSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetOwner(&_Exchange.TransactOpts, newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Exchange *ExchangeTransactorSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetOwner(&_Exchange.TransactOpts, newOwner)
}

// SetWethToken is a paid mutator transaction binding the contract method 0x86e09c08.
//
// Solidity: function setWethToken(_wethToken address) returns(bool)
func (_Exchange *ExchangeTransactor) SetWethToken(opts *bind.TransactOpts, _wethToken common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setWethToken", _wethToken)
}

// SetWethToken is a paid mutator transaction binding the contract method 0x86e09c08.
//
// Solidity: function setWethToken(_wethToken address) returns(bool)
func (_Exchange *ExchangeSession) SetWethToken(_wethToken common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetWethToken(&_Exchange.TransactOpts, _wethToken)
}

// SetWethToken is a paid mutator transaction binding the contract method 0x86e09c08.
//
// Solidity: function setWethToken(_wethToken address) returns(bool)
func (_Exchange *ExchangeTransactorSession) SetWethToken(_wethToken common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetWethToken(&_Exchange.TransactOpts, _wethToken)
}

// ValidateSignatures is a paid mutator transaction binding the contract method 0x1778baf4.
//
// Solidity: function validateSignatures(orderValues uint256[10], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeTransactor) ValidateSignatures(opts *bind.TransactOpts, orderValues [10]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "validateSignatures", orderValues, orderAddresses, v, rs)
}

// ValidateSignatures is a paid mutator transaction binding the contract method 0x1778baf4.
//
// Solidity: function validateSignatures(orderValues uint256[10], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeSession) ValidateSignatures(orderValues [10]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ValidateSignatures(&_Exchange.TransactOpts, orderValues, orderAddresses, v, rs)
}

// ValidateSignatures is a paid mutator transaction binding the contract method 0x1778baf4.
//
// Solidity: function validateSignatures(orderValues uint256[10], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeTransactorSession) ValidateSignatures(orderValues [10]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ValidateSignatures(&_Exchange.TransactOpts, orderValues, orderAddresses, v, rs)
}

// ExchangeLogBatchTradesIterator is returned from FilterLogBatchTrades and is used to iterate over the raw logs and unpacked data for LogBatchTrades events raised by the Exchange contract.
type ExchangeLogBatchTradesIterator struct {
	Event *ExchangeLogBatchTrades // Event containing the contract specifics and raw log

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
func (it *ExchangeLogBatchTradesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogBatchTrades)
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
		it.Event = new(ExchangeLogBatchTrades)
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
func (it *ExchangeLogBatchTradesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogBatchTradesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogBatchTrades represents a LogBatchTrades event raised by the Exchange contract.
type ExchangeLogBatchTrades struct {
	MakerOrderHashes [][32]byte
	TakerOrderHashes [][32]byte
	TokenPairHash    [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogBatchTrades is a free log retrieval operation binding the contract event 0xde8acabe30c9bd25d65bb9db28bf46f51dc7500a07b1671f121f1144fbf446fc.
//
// Solidity: e LogBatchTrades(makerOrderHashes bytes32[], takerOrderHashes bytes32[], tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) FilterLogBatchTrades(opts *bind.FilterOpts, tokenPairHash [][32]byte) (*ExchangeLogBatchTradesIterator, error) {

	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogBatchTrades", tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeLogBatchTradesIterator{contract: _Exchange.contract, event: "LogBatchTrades", logs: logs, sub: sub}, nil
}

// WatchLogBatchTrades is a free log subscription operation binding the contract event 0xde8acabe30c9bd25d65bb9db28bf46f51dc7500a07b1671f121f1144fbf446fc.
//
// Solidity: e LogBatchTrades(makerOrderHashes bytes32[], takerOrderHashes bytes32[], tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) WatchLogBatchTrades(opts *bind.WatchOpts, sink chan<- *ExchangeLogBatchTrades, tokenPairHash [][32]byte) (event.Subscription, error) {

	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogBatchTrades", tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogBatchTrades)
				if err := _Exchange.contract.UnpackLog(event, "LogBatchTrades", log); err != nil {
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

// ExchangeLogCancelOrderIterator is returned from FilterLogCancelOrder and is used to iterate over the raw logs and unpacked data for LogCancelOrder events raised by the Exchange contract.
type ExchangeLogCancelOrderIterator struct {
	Event *ExchangeLogCancelOrder // Event containing the contract specifics and raw log

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
func (it *ExchangeLogCancelOrderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogCancelOrder)
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
		it.Event = new(ExchangeLogCancelOrder)
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
func (it *ExchangeLogCancelOrderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogCancelOrderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogCancelOrder represents a LogCancelOrder event raised by the Exchange contract.
type ExchangeLogCancelOrder struct {
	OrderHash   [32]byte
	UserAddress common.Address
	BaseToken   common.Address
	QuoteToken  common.Address
	Amount      *big.Int
	Pricepoint  *big.Int
	Side        *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLogCancelOrder is a free log retrieval operation binding the contract event 0xb00984fe824f4973f31e8a414157f54cb4ee29bc2100149ba22a094d0bfd5518.
//
// Solidity: e LogCancelOrder(orderHash bytes32, userAddress address, baseToken address, quoteToken address, amount uint256, pricepoint uint256, side uint256)
func (_Exchange *ExchangeFilterer) FilterLogCancelOrder(opts *bind.FilterOpts) (*ExchangeLogCancelOrderIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogCancelOrder")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogCancelOrderIterator{contract: _Exchange.contract, event: "LogCancelOrder", logs: logs, sub: sub}, nil
}

// WatchLogCancelOrder is a free log subscription operation binding the contract event 0xb00984fe824f4973f31e8a414157f54cb4ee29bc2100149ba22a094d0bfd5518.
//
// Solidity: e LogCancelOrder(orderHash bytes32, userAddress address, baseToken address, quoteToken address, amount uint256, pricepoint uint256, side uint256)
func (_Exchange *ExchangeFilterer) WatchLogCancelOrder(opts *bind.WatchOpts, sink chan<- *ExchangeLogCancelOrder) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogCancelOrder")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogCancelOrder)
				if err := _Exchange.contract.UnpackLog(event, "LogCancelOrder", log); err != nil {
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

// ExchangeLogErrorIterator is returned from FilterLogError and is used to iterate over the raw logs and unpacked data for LogError events raised by the Exchange contract.
type ExchangeLogErrorIterator struct {
	Event *ExchangeLogError // Event containing the contract specifics and raw log

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
func (it *ExchangeLogErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogError)
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
		it.Event = new(ExchangeLogError)
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
func (it *ExchangeLogErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogError represents a LogError event raised by the Exchange contract.
type ExchangeLogError struct {
	ErrorId        uint8
	MakerOrderHash [32]byte
	TakerOrderHash [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterLogError is a free log retrieval operation binding the contract event 0x14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb.
//
// Solidity: e LogError(errorId uint8, makerOrderHash bytes32, takerOrderHash bytes32)
func (_Exchange *ExchangeFilterer) FilterLogError(opts *bind.FilterOpts) (*ExchangeLogErrorIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogError")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogErrorIterator{contract: _Exchange.contract, event: "LogError", logs: logs, sub: sub}, nil
}

// WatchLogError is a free log subscription operation binding the contract event 0x14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb.
//
// Solidity: e LogError(errorId uint8, makerOrderHash bytes32, takerOrderHash bytes32)
func (_Exchange *ExchangeFilterer) WatchLogError(opts *bind.WatchOpts, sink chan<- *ExchangeLogError) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogError)
				if err := _Exchange.contract.UnpackLog(event, "LogError", log); err != nil {
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

// ExchangeLogOperatorUpdateIterator is returned from FilterLogOperatorUpdate and is used to iterate over the raw logs and unpacked data for LogOperatorUpdate events raised by the Exchange contract.
type ExchangeLogOperatorUpdateIterator struct {
	Event *ExchangeLogOperatorUpdate // Event containing the contract specifics and raw log

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
func (it *ExchangeLogOperatorUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogOperatorUpdate)
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
		it.Event = new(ExchangeLogOperatorUpdate)
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
func (it *ExchangeLogOperatorUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogOperatorUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogOperatorUpdate represents a LogOperatorUpdate event raised by the Exchange contract.
type ExchangeLogOperatorUpdate struct {
	Operator   common.Address
	IsOperator bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogOperatorUpdate is a free log retrieval operation binding the contract event 0x4af650e9ee9ac50b37ec2cd3ddac7e1c69955ffc871bc3e812563775f3bc0e7d.
//
// Solidity: e LogOperatorUpdate(operator address, isOperator bool)
func (_Exchange *ExchangeFilterer) FilterLogOperatorUpdate(opts *bind.FilterOpts) (*ExchangeLogOperatorUpdateIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogOperatorUpdate")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogOperatorUpdateIterator{contract: _Exchange.contract, event: "LogOperatorUpdate", logs: logs, sub: sub}, nil
}

// WatchLogOperatorUpdate is a free log subscription operation binding the contract event 0x4af650e9ee9ac50b37ec2cd3ddac7e1c69955ffc871bc3e812563775f3bc0e7d.
//
// Solidity: e LogOperatorUpdate(operator address, isOperator bool)
func (_Exchange *ExchangeFilterer) WatchLogOperatorUpdate(opts *bind.WatchOpts, sink chan<- *ExchangeLogOperatorUpdate) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogOperatorUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogOperatorUpdate)
				if err := _Exchange.contract.UnpackLog(event, "LogOperatorUpdate", log); err != nil {
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

// ExchangeLogRewardAccountUpdateIterator is returned from FilterLogRewardAccountUpdate and is used to iterate over the raw logs and unpacked data for LogRewardAccountUpdate events raised by the Exchange contract.
type ExchangeLogRewardAccountUpdateIterator struct {
	Event *ExchangeLogRewardAccountUpdate // Event containing the contract specifics and raw log

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
func (it *ExchangeLogRewardAccountUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogRewardAccountUpdate)
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
		it.Event = new(ExchangeLogRewardAccountUpdate)
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
func (it *ExchangeLogRewardAccountUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogRewardAccountUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogRewardAccountUpdate represents a LogRewardAccountUpdate event raised by the Exchange contract.
type ExchangeLogRewardAccountUpdate struct {
	OldRewardAccount common.Address
	NewRewardAccount common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogRewardAccountUpdate is a free log retrieval operation binding the contract event 0x18d40614e4a77383f4b7337227bdad137b4f3f9b002ef63afd3ddaa142a15f63.
//
// Solidity: e LogRewardAccountUpdate(oldRewardAccount address, newRewardAccount address)
func (_Exchange *ExchangeFilterer) FilterLogRewardAccountUpdate(opts *bind.FilterOpts) (*ExchangeLogRewardAccountUpdateIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogRewardAccountUpdate")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogRewardAccountUpdateIterator{contract: _Exchange.contract, event: "LogRewardAccountUpdate", logs: logs, sub: sub}, nil
}

// WatchLogRewardAccountUpdate is a free log subscription operation binding the contract event 0x18d40614e4a77383f4b7337227bdad137b4f3f9b002ef63afd3ddaa142a15f63.
//
// Solidity: e LogRewardAccountUpdate(oldRewardAccount address, newRewardAccount address)
func (_Exchange *ExchangeFilterer) WatchLogRewardAccountUpdate(opts *bind.WatchOpts, sink chan<- *ExchangeLogRewardAccountUpdate) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogRewardAccountUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogRewardAccountUpdate)
				if err := _Exchange.contract.UnpackLog(event, "LogRewardAccountUpdate", log); err != nil {
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

// ExchangeLogTradeIterator is returned from FilterLogTrade and is used to iterate over the raw logs and unpacked data for LogTrade events raised by the Exchange contract.
type ExchangeLogTradeIterator struct {
	Event *ExchangeLogTrade // Event containing the contract specifics and raw log

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
func (it *ExchangeLogTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogTrade)
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
		it.Event = new(ExchangeLogTrade)
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
func (it *ExchangeLogTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogTrade represents a LogTrade event raised by the Exchange contract.
type ExchangeLogTrade struct {
	Maker            common.Address
	Taker            common.Address
	TokenSell        common.Address
	TokenBuy         common.Address
	FilledAmountSell *big.Int
	FilledAmountBuy  *big.Int
	PaidFeeMake      *big.Int
	PaidFeeTake      *big.Int
	OrderHash        [32]byte
	TradeHash        [32]byte
	TokenPairHash    [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogTrade is a free log retrieval operation binding the contract event 0x174a42d8fdc3a48bf80a4e95ac4b280ef69189e4603105caac770bf9771357fc.
//
// Solidity: e LogTrade(maker indexed address, taker indexed address, tokenSell address, tokenBuy address, filledAmountSell uint256, filledAmountBuy uint256, paidFeeMake uint256, paidFeeTake uint256, orderHash bytes32, tradeHash bytes32, tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) FilterLogTrade(opts *bind.FilterOpts, maker []common.Address, taker []common.Address, tokenPairHash [][32]byte) (*ExchangeLogTradeIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogTrade", makerRule, takerRule, tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeLogTradeIterator{contract: _Exchange.contract, event: "LogTrade", logs: logs, sub: sub}, nil
}

// WatchLogTrade is a free log subscription operation binding the contract event 0x174a42d8fdc3a48bf80a4e95ac4b280ef69189e4603105caac770bf9771357fc.
//
// Solidity: e LogTrade(maker indexed address, taker indexed address, tokenSell address, tokenBuy address, filledAmountSell uint256, filledAmountBuy uint256, paidFeeMake uint256, paidFeeTake uint256, orderHash bytes32, tradeHash bytes32, tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) WatchLogTrade(opts *bind.WatchOpts, sink chan<- *ExchangeLogTrade, maker []common.Address, taker []common.Address, tokenPairHash [][32]byte) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogTrade", makerRule, takerRule, tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogTrade)
				if err := _Exchange.contract.UnpackLog(event, "LogTrade", log); err != nil {
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

// ExchangeLogWethTokenUpdateIterator is returned from FilterLogWethTokenUpdate and is used to iterate over the raw logs and unpacked data for LogWethTokenUpdate events raised by the Exchange contract.
type ExchangeLogWethTokenUpdateIterator struct {
	Event *ExchangeLogWethTokenUpdate // Event containing the contract specifics and raw log

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
func (it *ExchangeLogWethTokenUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogWethTokenUpdate)
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
		it.Event = new(ExchangeLogWethTokenUpdate)
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
func (it *ExchangeLogWethTokenUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogWethTokenUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogWethTokenUpdate represents a LogWethTokenUpdate event raised by the Exchange contract.
type ExchangeLogWethTokenUpdate struct {
	OldWethToken common.Address
	NewWethToken common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLogWethTokenUpdate is a free log retrieval operation binding the contract event 0xb8be72b4c168c2f7d3ea469d9f48ccbc62416784a4f6a69ca93ff13f4f36545b.
//
// Solidity: e LogWethTokenUpdate(oldWethToken address, newWethToken address)
func (_Exchange *ExchangeFilterer) FilterLogWethTokenUpdate(opts *bind.FilterOpts) (*ExchangeLogWethTokenUpdateIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogWethTokenUpdate")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogWethTokenUpdateIterator{contract: _Exchange.contract, event: "LogWethTokenUpdate", logs: logs, sub: sub}, nil
}

// WatchLogWethTokenUpdate is a free log subscription operation binding the contract event 0xb8be72b4c168c2f7d3ea469d9f48ccbc62416784a4f6a69ca93ff13f4f36545b.
//
// Solidity: e LogWethTokenUpdate(oldWethToken address, newWethToken address)
func (_Exchange *ExchangeFilterer) WatchLogWethTokenUpdate(opts *bind.WatchOpts, sink chan<- *ExchangeLogWethTokenUpdate) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogWethTokenUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogWethTokenUpdate)
				if err := _Exchange.contract.UnpackLog(event, "LogWethTokenUpdate", log); err != nil {
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

// ExchangeSetOwnerIterator is returned from FilterSetOwner and is used to iterate over the raw logs and unpacked data for SetOwner events raised by the Exchange contract.
type ExchangeSetOwnerIterator struct {
	Event *ExchangeSetOwner // Event containing the contract specifics and raw log

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
func (it *ExchangeSetOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeSetOwner)
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
		it.Event = new(ExchangeSetOwner)
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
func (it *ExchangeSetOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeSetOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeSetOwner represents a SetOwner event raised by the Exchange contract.
type ExchangeSetOwner struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSetOwner is a free log retrieval operation binding the contract event 0xcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c663.
//
// Solidity: e SetOwner(previousOwner indexed address, newOwner indexed address)
func (_Exchange *ExchangeFilterer) FilterSetOwner(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExchangeSetOwnerIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "SetOwner", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeSetOwnerIterator{contract: _Exchange.contract, event: "SetOwner", logs: logs, sub: sub}, nil
}

// WatchSetOwner is a free log subscription operation binding the contract event 0xcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c663.
//
// Solidity: e SetOwner(previousOwner indexed address, newOwner indexed address)
func (_Exchange *ExchangeFilterer) WatchSetOwner(opts *bind.WatchOpts, sink chan<- *ExchangeSetOwner, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "SetOwner", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeSetOwner)
				if err := _Exchange.contract.UnpackLog(event, "SetOwner", log); err != nil {
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

// OwnedABI is the input ABI used to generate the binding from.
const OwnedABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"}]"

// OwnedBin is the compiled bytecode used for deploying new contracts.
const OwnedBin = "0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610255806100606000396000f30060806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806313af4035146100515780638da5cb5b14610094575b600080fd5b34801561005d57600080fd5b50610092600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506100eb565b005b3480156100a057600080fd5b506100a9610204565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561014657600080fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c66360405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a72305820f75d0afdbb87452350bbe83e9576bb2229ae569d65e90a12d64c5c64042d781e0029"

// DeployOwned deploys a new Ethereum contract, binding an instance of Owned to it.
func DeployOwned(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Owned, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// Owned is an auto generated Go binding around an Ethereum contract.
type Owned struct {
	OwnedCaller     // Read-only binding to the contract
	OwnedTransactor // Write-only binding to the contract
	OwnedFilterer   // Log filterer for contract events
}

// OwnedCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedSession struct {
	Contract     *Owned            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedCallerSession struct {
	Contract *OwnedCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OwnedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedTransactorSession struct {
	Contract     *OwnedTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedRaw struct {
	Contract *Owned // Generic contract binding to access the raw methods on
}

// OwnedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedCallerRaw struct {
	Contract *OwnedCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedTransactorRaw struct {
	Contract *OwnedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwned creates a new instance of Owned, bound to a specific deployed contract.
func NewOwned(address common.Address, backend bind.ContractBackend) (*Owned, error) {
	contract, err := bindOwned(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// NewOwnedCaller creates a new read-only instance of Owned, bound to a specific deployed contract.
func NewOwnedCaller(address common.Address, caller bind.ContractCaller) (*OwnedCaller, error) {
	contract, err := bindOwned(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedCaller{contract: contract}, nil
}

// NewOwnedTransactor creates a new write-only instance of Owned, bound to a specific deployed contract.
func NewOwnedTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTransactor, error) {
	contract, err := bindOwned(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTransactor{contract: contract}, nil
}

// NewOwnedFilterer creates a new log filterer instance of Owned, bound to a specific deployed contract.
func NewOwnedFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedFilterer, error) {
	contract, err := bindOwned(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedFilterer{contract: contract}, nil
}

// bindOwned binds a generic wrapper to an already deployed contract.
func bindOwned(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Owned.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCallerSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Owned *OwnedTransactor) SetOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "setOwner", newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Owned *OwnedSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.SetOwner(&_Owned.TransactOpts, newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Owned *OwnedTransactorSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.SetOwner(&_Owned.TransactOpts, newOwner)
}

// OwnedSetOwnerIterator is returned from FilterSetOwner and is used to iterate over the raw logs and unpacked data for SetOwner events raised by the Owned contract.
type OwnedSetOwnerIterator struct {
	Event *OwnedSetOwner // Event containing the contract specifics and raw log

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
func (it *OwnedSetOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedSetOwner)
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
		it.Event = new(OwnedSetOwner)
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
func (it *OwnedSetOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedSetOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedSetOwner represents a SetOwner event raised by the Owned contract.
type OwnedSetOwner struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSetOwner is a free log retrieval operation binding the contract event 0xcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c663.
//
// Solidity: e SetOwner(previousOwner indexed address, newOwner indexed address)
func (_Owned *OwnedFilterer) FilterSetOwner(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnedSetOwnerIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Owned.contract.FilterLogs(opts, "SetOwner", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnedSetOwnerIterator{contract: _Owned.contract, event: "SetOwner", logs: logs, sub: sub}, nil
}

// WatchSetOwner is a free log subscription operation binding the contract event 0xcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c663.
//
// Solidity: e SetOwner(previousOwner indexed address, newOwner indexed address)
func (_Owned *OwnedFilterer) WatchSetOwner(opts *bind.WatchOpts, sink chan<- *OwnedSetOwner, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Owned.contract.WatchLogs(opts, "SetOwner", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedSetOwner)
				if err := _Owned.contract.UnpackLog(event, "SetOwner", log); err != nil {
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

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820771f68f066c4bcb1a6844fa47607bbe2dfbe08557f5e86e9aafdead2287d13d60029`

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}
