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
const ExchangeBin = `0x608060405234801561001057600080fd5b5060405160408061274483398101604052805160209091015160008054600160a060020a0319908116331790915560018054600160a060020a03948516908316179055600280549390921692169190911790556126d2806100726000396000f3006080604052600436106101485763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630e708203811461014d57806310ac00d81461017e57806313af40351461024f57806313e7c9d81461027257806314df96ee146102935780631778baf4146102b1578063288cdc91146103655780633c9183411461038f5780634b023cf8146103b95780634b57b0be146103da5780635171267f146103ef578063558a729714610605578063673e04811461062b5780638163681e1461067557806386e09c08146106a55780638da5cb5b146106c657806393c1ae09146106db57806398024a8b146107905780639e5ebf5e146107ae578063b4cb2553146107d5578063d581332314610868578063d9a72b5214610880578063e51ad32d146108f7578063f4a8726314610aa3578063ffa1ad7414610aca575b600080fd5b34801561015957600080fd5b50610162610b54565b60408051600160a060020a039092168252519081900360200190f35b34801561018a57600080fd5b506040805161014081810190925261023b91369160049161014491908390600a90839083908082843750506040805160808181019092529497969581810195945092506004915083908390808284375050604080518082018252949786359790969095606082019550935060200191506002908390839080828437505060408051608081810190925294979695818101959450925060049150839083908082843750939650610b6395505050505050565b604080519115158252519081900360200190f35b34801561025b57600080fd5b50610270600160a060020a0360043516610beb565b005b34801561027e57600080fd5b5061023b600160a060020a0360043516610c6a565b34801561029f57600080fd5b5061023b600435602435604435610c7f565b3480156102bd57600080fd5b506040805161014081810190925261023b91369160049161014491908390600a90839083908082843750506040805160808181019092529497969581810195945092506004915083908390808284375050604080518082018252949796958181019594509250600291508390839080828437505060408051608081810190925294979695818101959450925060049150839083908082843750939650610ce895505050505050565b34801561037157600080fd5b5061037d600435610f24565b60408051918252519081900360200190f35b34801561039b57600080fd5b5061023b600160a060020a0360043581169060243516604435610f36565b3480156103c557600080fd5b5061023b600160a060020a0360043516610fe7565b3480156103e657600080fd5b5061016261108e565b3480156103fb57600080fd5b5060408051600480358082013560208181028501810190955280845261023b943694602493909290840191819060009085015b8282101561046b576040805161014081810190925290808402870190600a908390839080828437505050918352505060019091019060200161042e565b50506040805186358801803560208181028401810190945280835296999897830196919550820193509150819060009085015b828210156104da57604080516080818101909252908084028701906004908390839080828437505050918352505060019091019060200161049e565b505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437505060408051873589018035602081810284018101909452808352979a999883019791965082019450925082915060009085015b8282101561058657604080518082018252908084028701906002908390839080828437505050918352505060019091019060200161054d565b50506040805186358801803560208181028401810190945280835296999897830196919550820193509150819060009085015b828210156105f55760408051608081810190925290808402870190600490839083908082843750505091835250506001909101906020016105b9565b5093965061109d95505050505050565b34801561061157600080fd5b5061023b600160a060020a036004351660243515156112f4565b34801561063757600080fd5b50610643600435611395565b60408051948552600160a060020a03938416602086015291909216838201526060830191909152519081900360800190f35b34801561068157600080fd5b5061023b600160a060020a036004351660243560ff604435166064356084356113c9565b3480156106b157600080fd5b5061023b600160a060020a03600435166114f1565b3480156106d257600080fd5b50610162611581565b3480156106e757600080fd5b506040805160808181019092526102709136916004916084919083908190839082908082843750506040805186358801803560208181028481018201909552818452979a99988801979296509082019450925082919085019084908082843750506040805187358901803560208181028481018201909552818452989b9a9989019892975090820195509350839250850190849080828437509497506115909650505050505050565b34801561079c57600080fd5b5061037d6004356024356044356116f1565b3480156107ba57600080fd5b5061037d600160a060020a036004358116906024351661170f565b3480156107e157600080fd5b506040805161014081810190925261084891369160049161014491908390600a90839083908082843750506040805160808181019092529497969581810195945092506004915083908390808284375093965050833594505050602090910135905061173a565b604080519384526020840192909252151582820152519081900360600190f35b34801561087457600080fd5b5061023b600435611eb8565b34801561088c57600080fd5b506040805160c081810190925261023b91369160049160c49190839060069083908390808284375050604080516060818101909252949796958181019594509250600391508390839080828437509396505050823560ff169350505060208101359060400135611ecd565b34801561090357600080fd5b50604080516004803580820135602081810285018101909552808452610270943694602493909290840191819060009085015b82821015610972576040805160c08181019092529080840287019060069083908390808284375050509183525050600190910190602001610936565b50506040805186358801803560208181028401810190945280835296999897830196919550820193509150819060009085015b828210156109e15760408051606081810190925290808402870190600390839083908082843750505091835250506001909101906020016109a5565b50505050509192919290803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a9989019892975090820195509350839250850190849080828437509497506120479650505050505050565b348015610aaf57600080fd5b5061023b600160a060020a03600435811690602435166120e4565b348015610ad657600080fd5b50610adf61211f565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610b19578181015183820152602001610b01565b50505050905090810190601f168015610b465780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600254600160a060020a031681565b6000805481908190600160a060020a0316331480610b9057503360009081526003602052604090205460ff165b1515610b9b57600080fd5b610ba788888787610ce8565b9150811515610bb95760009250610be0565b610bc287612156565b9050610bd08888888461173a565b505050610bde8888886121cf565b505b505095945050505050565b600054600160a060020a03163314610c0257600080fd5b60008054604051600160a060020a03808516939216917fcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c66391a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b60036020526000908152604090205460ff1681565b600080600084801515610c8e57fe5b8685099150811515610ca35760009250610cdf565b610cd5610cb6878663ffffffff6122a116565b610cc984620f424063ffffffff6122a116565b9063ffffffff6122cc16565b90506103e8811192505b50509392505050565b6000610cf26125d7565b610cfa6125d7565b604080516101208101909152600090819080898360209081029190910151600160a060020a031682520189600260209081029190910151600160a060020a031682520189600360209081029190910151600160a060020a03168252018a6000602090810291909101518252018a6001602090810291909101518252018a6002602090810291909101518252018a6003602090810291909101518252018a6008602090810291909101518252018a6009602002015190526040805161012081019091529094508089600160209081029190910151600160a060020a031682520189600260209081029190910151600160a060020a031682520189600360209081029190910151600160a060020a03168252018a6004602090810291909101518252018a6005602090810291909101518252018a6006602090810291909101518252018a6007602090810291909101518252018a6009602090810291909101518252018a6008602002015190529250610e70846122e3565b9150610e7b836122e3565b845188518851929350610e99928591908a60015b60200201516113c9565b1515610edc5760008051602061266783398151915260015b6040805160ff909216825260208201859052818101849052519081900360600190a160009450610f18565b825160208801516040880151610ef7929184918a6003610e8f565b1515610f13576000805160206126678339815191526002610eb1565b600194505b50505050949350505050565b60046020526000908152604090205481565b600080548190600160a060020a03163314610f5057600080fd5b610f5a85856123f1565b60408051608081018252828152600160a060020a039788166020808301918252978916828401908152606083019788526000948552600690985291909220915182555160018201805491881673ffffffffffffffffffffffffffffffffffffffff199283161790559451600282018054919097169516949094179094559051600390920191909155919050565b60008054600160a060020a03163314610fff57600080fd5b600160a060020a038216151561101457600080fd5b60025460408051600160a060020a039283168152918416602083015280517f18d40614e4a77383f4b7337227bdad137b4f3f9b002ef63afd3ddaa142a15f639281900390910190a15060028054600160a060020a03831673ffffffffffffffffffffffffffffffffffffffff199091161790556001919050565b600154600160a060020a031681565b600080546060908190839081908190819081908190600160a060020a03163314806110d757503360009081526003602052604090205460ff165b15156110e257600080fd5b8c5160405190808252806020026020018201604052801561110d578160200160208202803883390190505b5097508c5160405190808252806020026020018201604052801561113b578160200160208202803883390190505b509650600095505b8c51861015611284576111b48e8781518110151561115d57fe5b906020019060200201518e8881518110151561117557fe5b906020019060200201518d8981518110151561118d57fe5b906020019060200201518d8a8151811015156111a557fe5b90602001906020020151610ce8565b94508415156111c657600098506112e3565b6111e68d878151811015156111d757fe5b90602001906020020151612156565b93506112398e878151811015156111f957fe5b906020019060200201518e8881518110151561121157fe5b906020019060200201518e8981518110151561122957fe5b906020019060200201518761173a565b92509250925080156112795782888781518110151561125457fe5b602090810290910101528651829088908890811061126e57fe5b602090810290910101525b600190950194611143565b6112bf8e600081518110151561129657fe5b906020019060200201518e60008151811015156112af57fe5b906020019060200201518e6124bf565b506112e38d60008151811015156112d257fe5b906020019060200201518989611590565b505050505050505095945050505050565b60008054600160a060020a0316331461130c57600080fd5b600160a060020a038316151561132157600080fd5b60408051600160a060020a0385168152831515602082015281517f4af650e9ee9ac50b37ec2cd3ddac7e1c69955ffc871bc3e812563775f3bc0e7d929181900390910190a150600160a060020a0382166000908152600360205260409020805482151560ff19909116179055600192915050565b60066020526000908152604090208054600182015460028301546003909301549192600160a060020a039182169291169084565b600060018560405160200180807f19457468657265756d205369676e6564204d6573736167653a0a333200000000815250601c0182600019166000191681526020019150506040516020818303038152906040526040518082805190602001908083835b6020831061144c5780518252601f19909201916020918201910161142d565b51815160209384036101000a60001901801990921691161790526040805192909401829003822060008084528383018087529190915260ff8c1683860152606083018b9052608083018a9052935160a08084019750919550601f1981019492819003909101925090865af11580156114c8573d6000803e3d6000fd5b50505060206040510351600160a060020a031686600160a060020a031614905095945050505050565b60008054600160a060020a0316331461150957600080fd5b60015460408051600160a060020a039283168152918416602083015280517fb8be72b4c168c2f7d3ea469d9f48ccbc62416784a4f6a69ca93ff13f4f36545b9281900390910190a15060018054600160a060020a03831673ffffffffffffffffffffffffffffffffffffffff19909116178155919050565b600054600160a060020a031681565b60208084015160408086015181516c01000000000000000000000000600160a060020a039485168102828701529390911690920260348301528051602881840301815260489092019081905281519192909182918401908083835b6020831061160a5780518252601f1990920191602091820191016115eb565b51815160209384036101000a60001901801990921691161790526040805192909401829003822084835288519483019490945287519395507fde8acabe30c9bd25d65bb9db28bf46f51dc7500a07b1671f121f1144fbf446fc945087938793508291828101916060840191878101910280838360005b83811015611698578181015183820152602001611680565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156116d75781810151838201526020016116bf565b5050505090500194505050505060405180910390a2505050565b600061170783610cc9868563ffffffff6122a116565b949350505050565b60008061171c84846123f1565b600081815260066020526040902060030154925090505b5092915050565b60008060006117476125d7565b61174f6125d7565b600080548190819081908190600160a060020a031633148061178057503360009081526003602052604090205460ff165b151561178b57600080fd5b604080516101208101909152808e600060209081029190910151600160a060020a03168252018e600260209081029190910151600160a060020a03168252018e600360209081029190910151600160a060020a03168252018f6000602090810291909101518252018f6001602090810291909101518252018f6002602090810291909101518252018f6003602090810291909101518252018f6008602090810291909101518252018f600960200201519052604080516101208101909152909750808e600160209081029190910151600160a060020a03168252018e600260209081029190910151600160a060020a03168252018e600360209081029190910151600160a060020a03168252018f6004602090810291909101518252018f6005602090810291909101518252018f6006602090810291909101518252018f6007602090810291909101518252018f6009602090810291909101518252018f60086020020151905295506118fd876122e3565b9450611908866122e3565b606088015160008781526004602052604090205491955090611930908e63ffffffff6125c816565b11156119795760008051602061266783398151915260075b6040805160ff909216825260208201889052818101879052519081900360600190a184846000995099509950611ea7565b606086015160008581526004602052604090205461199d908e63ffffffff6125c816565b11156119b9576000805160206126678339815191526007611948565b8660a001518660a0015114156119df576000805160206126678339815191526003611948565b60a08701511515611a1057856080015187608001511015611a10576000805160206126678339815191526004611948565b8660a0015160011415611a4357866080015186608001511015611a43576000805160206126678339815191526004611948565b600084815260046020526040902054611a62908d63ffffffff6125c816565b600085815260046020526040808220929092558681522054611a8a908d63ffffffff6125c816565b60008681526004602052604090205560808701518c9350611ab8908c90610cc990869063ffffffff6122a116565b9150611acd8c88606001518960e001516116f1565b90508660a0015160001415611cbe576040808801518851885183516000805160206126878339815191528152600160a060020a03928316600482015290821660248201526044810186905292519116916323b872dd9160648083019260209291908290030181600087803b158015611b4457600080fd5b505af1158015611b58573d6000803e3d6000fd5b505050506040513d6020811015611b6e57600080fd5b50511515611b7b57600080fd5b604080880151885160025483516000805160206126878339815191528152600160a060020a03928316600482015290821660248201526044810185905292519116916323b872dd9160648083019260209291908290030181600087803b158015611be457600080fd5b505af1158015611bf8573d6000803e3d6000fd5b505050506040513d6020811015611c0e57600080fd5b50511515611c1b57600080fd5b60208087015187518951604080516000805160206126878339815191528152600160a060020a0393841660048201529183166024830152604482018890525191909216926323b872dd92606480820193918290030181600087803b158015611c8257600080fd5b505af1158015611c96573d6000803e3d6000fd5b505050506040513d6020811015611cac57600080fd5b50511515611cb957600080fd5b611e9c565b60208088015188518851604080516000805160206126878339815191528152600160a060020a0393841660048201529183166024830152604482018890525191909216926323b872dd92606480820193918290030181600087803b158015611d2557600080fd5b505af1158015611d39573d6000803e3d6000fd5b505050506040513d6020811015611d4f57600080fd5b50511515611d5c57600080fd5b604080870151875160025483516000805160206126878339815191528152600160a060020a03928316600482015290821660248201526044810185905292519116916323b872dd9160648083019260209291908290030181600087803b158015611dc557600080fd5b505af1158015611dd9573d6000803e3d6000fd5b505050506040513d6020811015611def57600080fd5b50511515611dfc57600080fd5b6040808701518751895183516000805160206126878339815191528152600160a060020a0392831660048201529082166024820152848603604482015292519116916323b872dd9160648083019260209291908290030181600087803b158015611e6557600080fd5b505af1158015611e79573d6000803e3d6000fd5b505050506040513d6020811015611e8f57600080fd5b50511515611e9c57600080fd5b848460019950995099505b505050505050509450945094915050565b60056020526000908152604090205460ff1681565b6000611ed76125d7565b5060408051610120810182528651600160a060020a03908116825260208089015182168184015288840151909116828401528851606080840191909152908901516080808401919091529289015160a0808401919091529089015160c083015288015160e0820152908701516101008201526000611f54826122e3565b9050611f6333828888886113c9565b1515611f9b5760408051600081526020810183905290516000805160206126678339815191529181900360600190a160009250610be0565b60608083018051600084815260046020908152604091829020929092558551828701518288015194516080808a015160a0808c015187518c8152600160a060020a03978816998101999099529486168888015294909716978601979097529584019590955282019290925260c0810192909252517fb00984fe824f4973f31e8a414157f54cb4ee29bc2100149ba22a094d0bfd55189181900360e00190a1506001979650505050505050565b60005b84518110156120dc576120d3868281518110151561206457fe5b90602001906020020151868381518110151561207c57fe5b90602001906020020151868481518110151561209457fe5b9060200190602002015186858151811015156120ac57fe5b9060200190602002015186868151811015156120c457fe5b90602001906020020151611ecd565b5060010161204a565b505050505050565b6000806120f184846123f1565b60008181526006602052604090206003015490915015156121155760009150611733565b5060019392505050565b60408051808201909152600581527f312e302e30000000000000000000000000000000000000000000000000000000602082015281565b60008061216161263f565b6040840151606085015161217591906123f1565b6000908152600660209081526040918290208251608081018452815481526001820154600160a060020a03908116938201939093526002820154909216928201929092526003909101546060909101819052949350505050565b608083015161010084015160208401516060850151600093929190846121f68786866116f1565b600254604080516000805160206126878339815191528152600160a060020a0387811660048301529283166024820152604481018490529051929350908416916323b872dd916064808201926020929091908290030181600087803b15801561225e57600080fd5b505af1158015612272573d6000803e3d6000fd5b505050506040513d602081101561228857600080fd5b5051151561229557600080fd5b50505050509392505050565b60008282028315806122bd57508284828115156122ba57fe5b04145b15156122c557fe5b9392505050565b60008082848115156122da57fe5b04949350505050565b80516020808301516040808501516060860151608087015160a088015160c089015160e08a01516101008b015187516c01000000000000000000000000308102828d0152600160a060020a039c8d1681026034830152998c168a0260488201529a909616909702605c8a01526070890193909352609088019190915260b087015260d086015260f0850192909252610110808501929092528051808503909201825261013090930192839052805160009391928291908401908083835b602083106123bf5780518252601f1990920191602091820191016123a0565b5181516020939093036101000a6000190180199091169216919091179052604051920182900390912095945050505050565b600082826040516020018083600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140182600160a060020a0316600160a060020a03166c01000000000000000000000000028152601401925050506040516020818303038152906040526040518082805190602001908083835b6020831061248c5780518252601f19909201916020918201910161246d565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091209695505050505050565b6080830151610100840151602084015160608501516000939291908480805b88518210156125105788828151811015156124f557fe5b906020019060200201518301925081806001019250506124de565b61251b8388886116f1565b600254604080516000805160206126878339815191528152600160a060020a0389811660048301529283166024820152604481018490529051929350908616916323b872dd916064808201926020929091908290030181600087803b15801561258357600080fd5b505af1158015612597573d6000803e3d6000fd5b505050506040513d60208110156125ad57600080fd5b505115156125ba57600080fd5b505050505050509392505050565b6000828201838110156122c557fe5b610120604051908101604052806000600160a060020a031681526020016000600160a060020a031681526020016000600160a060020a031681526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b60408051608081018252600080825260208201819052918101829052606081019190915290560014301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb23b872dd00000000000000000000000000000000000000000000000000000000a165627a7a7230582024ab2ba6ed7ba878e984209e6c7bb00d0c83734cd66d34dbc53a1b53d2a2edfb0029`

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
const OwnedBin = `0x608060405234801561001057600080fd5b5060008054600160a060020a031916331790556101ac806100326000396000f30060806040526004361061004b5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166313af403581146100505780638da5cb5b14610080575b600080fd5b34801561005c57600080fd5b5061007e73ffffffffffffffffffffffffffffffffffffffff600435166100be565b005b34801561008c57600080fd5b50610095610164565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60005473ffffffffffffffffffffffffffffffffffffffff1633146100e257600080fd5b6000805460405173ffffffffffffffffffffffffffffffffffffffff808516939216917fcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c66391a36000805473ffffffffffffffffffffffffffffffffffffffff191673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a7230582076e36060821a6167abb283e930822baede8db6b867d5ad2ae62a2e530ed0aad10029`

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
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a723058201d05756c063faa48d970328dfbd42d1d8e0bab5ab14a3e37f25c5fedefc384b90029`

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
