// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package did

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

// EthereumDIDRegistryMetaData contains all meta data concerning the EthereumDIDRegistry contract.
var EthereumDIDRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"validTo\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousChange\",\"type\":\"uint256\"}],\"name\":\"DIDAttributeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"delegateType\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"validTo\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousChange\",\"type\":\"uint256\"}],\"name\":\"DIDDelegateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousChange\",\"type\":\"uint256\"}],\"name\":\"DIDOwnerChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"delegateType\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"validity\",\"type\":\"uint256\"}],\"name\":\"addDelegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"sigV\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"sigR\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sigS\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"delegateType\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"validity\",\"type\":\"uint256\"}],\"name\":\"addDelegateSigned\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"sigV\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"sigR\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sigS\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"changeOwnerSigned\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"changed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"delegates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"}],\"name\":\"identityOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"owners\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"revokeAttribute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"sigV\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"sigR\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sigS\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"revokeAttributeSigned\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"delegateType\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"name\":\"revokeDelegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"sigV\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"sigR\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sigS\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"delegateType\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"name\":\"revokeDelegateSigned\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validity\",\"type\":\"uint256\"}],\"name\":\"setAttribute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"sigV\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"sigR\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sigS\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validity\",\"type\":\"uint256\"}],\"name\":\"setAttributeSigned\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"identity\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"delegateType\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"name\":\"validDelegate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// [{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"sigV","type":"uint8"},{"name":"sigR","type":"bytes32"},{"name":"sigS","type":"bytes32"},{"name":"delegateType","type":"bytes32"},{"name":"delegate","type":"address"},{"name":"validity","type":"uint256"}],"name":"addDelegateSigned","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"sigV","type":"uint8"},{"name":"sigR","type":"bytes32"},{"name":"sigS","type":"bytes32"},{"name":"name","type":"bytes32"},{"name":"value","type":"bytes"},{"name":"validTo","type":"uint256"}],"name":"setAttributeSigned","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"sigV","type":"uint8"},{"name":"sigR","type":"bytes32"},{"name":"sigS","type":"bytes32"},{"name":"delegateType","type":"bytes32"},{"name":"delegate","type":"address"}],"name":"revokeDelegateSigned","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"delegateType","type":"bytes32"},{"name":"delegate","type":"address"},{"name":"validity","type":"uint256"}],"name":"addDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"changed","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"name","type":"bytes32"},{"name":"value","type":"bytes"},{"name":"validTo","type":"uint256"}],"name":"setAttribute","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"sigV","type":"uint8"},{"name":"sigR","type":"bytes32"},{"name":"sigS","type":"bytes32"},{"name":"newOwner","type":"address"}],"name":"changeOwnerSigned","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"delegateType","type":"bytes32"},{"name":"delegate","type":"address"}],"name":"revokeDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"identity","type":"address"},{"name":"delegateType","type":"bytes32"},{"name":"delegate","type":"address"}],"name":"validDelegate","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"bytes32"},{"name":"","type":"address"}],"name":"delegates","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"newOwner","type":"address"}],"name":"changeOwner","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"nonce","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"owners","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"name","type":"bytes32"},{"name":"value","type":"bytes"}],"name":"revokeAttribute","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"identity","type":"address"},{"name":"sigV","type":"uint8"},{"name":"sigR","type":"bytes32"},{"name":"sigS","type":"bytes32"},{"name":"name","type":"bytes32"},{"name":"value","type":"bytes"}],"name":"revokeAttributeSigned","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"identity","type":"address"},{"indexed":false,"name":"owner","type":"address"},{"indexed":false,"name":"previousChange","type":"uint256"}],"name":"DIDOwnerChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"identity","type":"address"},{"indexed":false,"name":"delegateType","type":"bytes32"},{"indexed":false,"name":"delegate","type":"address"},{"indexed":false,"name":"validTo","type":"uint256"},{"indexed":false,"name":"previousChange","type":"uint256"}],"name":"DIDDelegateChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"identity","type":"address"},{"indexed":false,"name":"name","type":"bytes32"},{"indexed":false,"name":"value","type":"bytes"},{"indexed":false,"name":"validTo","type":"uint256"},{"indexed":false,"name":"previousChange","type":"uint256"}],"name":"DIDAttributeChanged","type":"event"}]
// EthereumDIDRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use EthereumDIDRegistryMetaData.ABI instead.
var EthereumDIDRegistryABI = EthereumDIDRegistryMetaData.ABI

// EthereumDIDRegistry is an auto generated Go binding around an Ethereum contract.
type EthereumDIDRegistry struct {
	EthereumDIDRegistryCaller     // Read-only binding to the contract
	EthereumDIDRegistryTransactor // Write-only binding to the contract
	EthereumDIDRegistryFilterer   // Log filterer for contract events
}

// EthereumDIDRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthereumDIDRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumDIDRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthereumDIDRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumDIDRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthereumDIDRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumDIDRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthereumDIDRegistrySession struct {
	Contract     *EthereumDIDRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// EthereumDIDRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthereumDIDRegistryCallerSession struct {
	Contract *EthereumDIDRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// EthereumDIDRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthereumDIDRegistryTransactorSession struct {
	Contract     *EthereumDIDRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// EthereumDIDRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthereumDIDRegistryRaw struct {
	Contract *EthereumDIDRegistry // Generic contract binding to access the raw methods on
}

// EthereumDIDRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthereumDIDRegistryCallerRaw struct {
	Contract *EthereumDIDRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// EthereumDIDRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthereumDIDRegistryTransactorRaw struct {
	Contract *EthereumDIDRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthereumDIDRegistry creates a new instance of EthereumDIDRegistry, bound to a specific deployed contract.
func NewEthereumDIDRegistry(address common.Address, backend bind.ContractBackend) (*EthereumDIDRegistry, error) {
	contract, err := bindEthereumDIDRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthereumDIDRegistry{EthereumDIDRegistryCaller: EthereumDIDRegistryCaller{contract: contract}, EthereumDIDRegistryTransactor: EthereumDIDRegistryTransactor{contract: contract}, EthereumDIDRegistryFilterer: EthereumDIDRegistryFilterer{contract: contract}}, nil
}

// NewEthereumDIDRegistryCaller creates a new read-only instance of EthereumDIDRegistry, bound to a specific deployed contract.
func NewEthereumDIDRegistryCaller(address common.Address, caller bind.ContractCaller) (*EthereumDIDRegistryCaller, error) {
	contract, err := bindEthereumDIDRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumDIDRegistryCaller{contract: contract}, nil
}

// NewEthereumDIDRegistryTransactor creates a new write-only instance of EthereumDIDRegistry, bound to a specific deployed contract.
func NewEthereumDIDRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*EthereumDIDRegistryTransactor, error) {
	contract, err := bindEthereumDIDRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumDIDRegistryTransactor{contract: contract}, nil
}

// NewEthereumDIDRegistryFilterer creates a new log filterer instance of EthereumDIDRegistry, bound to a specific deployed contract.
func NewEthereumDIDRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*EthereumDIDRegistryFilterer, error) {
	contract, err := bindEthereumDIDRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthereumDIDRegistryFilterer{contract: contract}, nil
}

// bindEthereumDIDRegistry binds a generic wrapper to an already deployed contract.
func bindEthereumDIDRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EthereumDIDRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthereumDIDRegistry *EthereumDIDRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthereumDIDRegistry.Contract.EthereumDIDRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthereumDIDRegistry *EthereumDIDRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.EthereumDIDRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthereumDIDRegistry *EthereumDIDRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.EthereumDIDRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthereumDIDRegistry *EthereumDIDRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthereumDIDRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.contract.Transact(opts, method, params...)
}

// Changed is a free data retrieval call binding the contract method 0xf96d0f9f.
//
// Solidity: function changed(address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistryCaller) Changed(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EthereumDIDRegistry.contract.Call(opts, &out, "changed", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Changed is a free data retrieval call binding the contract method 0xf96d0f9f.
//
// Solidity: function changed(address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) Changed(arg0 common.Address) (*big.Int, error) {
	return _EthereumDIDRegistry.Contract.Changed(&_EthereumDIDRegistry.CallOpts, arg0)
}

// Changed is a free data retrieval call binding the contract method 0xf96d0f9f.
//
// Solidity: function changed(address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistryCallerSession) Changed(arg0 common.Address) (*big.Int, error) {
	return _EthereumDIDRegistry.Contract.Changed(&_EthereumDIDRegistry.CallOpts, arg0)
}

// Delegates is a free data retrieval call binding the contract method 0x0d44625b.
//
// Solidity: function delegates(address , bytes32 , address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistryCaller) Delegates(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EthereumDIDRegistry.contract.Call(opts, &out, "delegates", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Delegates is a free data retrieval call binding the contract method 0x0d44625b.
//
// Solidity: function delegates(address , bytes32 , address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) Delegates(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _EthereumDIDRegistry.Contract.Delegates(&_EthereumDIDRegistry.CallOpts, arg0, arg1, arg2)
}

// Delegates is a free data retrieval call binding the contract method 0x0d44625b.
//
// Solidity: function delegates(address , bytes32 , address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistryCallerSession) Delegates(arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return _EthereumDIDRegistry.Contract.Delegates(&_EthereumDIDRegistry.CallOpts, arg0, arg1, arg2)
}

// IdentityOwner is a free data retrieval call binding the contract method 0x8733d4e8.
//
// Solidity: function identityOwner(address identity) view returns(address)
func (_EthereumDIDRegistry *EthereumDIDRegistryCaller) IdentityOwner(opts *bind.CallOpts, identity common.Address) (common.Address, error) {
	var out []interface{}
	err := _EthereumDIDRegistry.contract.Call(opts, &out, "identityOwner", identity)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IdentityOwner is a free data retrieval call binding the contract method 0x8733d4e8.
//
// Solidity: function identityOwner(address identity) view returns(address)
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) IdentityOwner(identity common.Address) (common.Address, error) {
	return _EthereumDIDRegistry.Contract.IdentityOwner(&_EthereumDIDRegistry.CallOpts, identity)
}

// IdentityOwner is a free data retrieval call binding the contract method 0x8733d4e8.
//
// Solidity: function identityOwner(address identity) view returns(address)
func (_EthereumDIDRegistry *EthereumDIDRegistryCallerSession) IdentityOwner(identity common.Address) (common.Address, error) {
	return _EthereumDIDRegistry.Contract.IdentityOwner(&_EthereumDIDRegistry.CallOpts, identity)
}

// Nonce is a free data retrieval call binding the contract method 0x70ae92d2.
//
// Solidity: function nonce(address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistryCaller) Nonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EthereumDIDRegistry.contract.Call(opts, &out, "nonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0x70ae92d2.
//
// Solidity: function nonce(address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) Nonce(arg0 common.Address) (*big.Int, error) {
	return _EthereumDIDRegistry.Contract.Nonce(&_EthereumDIDRegistry.CallOpts, arg0)
}

// Nonce is a free data retrieval call binding the contract method 0x70ae92d2.
//
// Solidity: function nonce(address ) view returns(uint256)
func (_EthereumDIDRegistry *EthereumDIDRegistryCallerSession) Nonce(arg0 common.Address) (*big.Int, error) {
	return _EthereumDIDRegistry.Contract.Nonce(&_EthereumDIDRegistry.CallOpts, arg0)
}

// Owners is a free data retrieval call binding the contract method 0x022914a7.
//
// Solidity: function owners(address ) view returns(address)
func (_EthereumDIDRegistry *EthereumDIDRegistryCaller) Owners(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _EthereumDIDRegistry.contract.Call(opts, &out, "owners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owners is a free data retrieval call binding the contract method 0x022914a7.
//
// Solidity: function owners(address ) view returns(address)
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) Owners(arg0 common.Address) (common.Address, error) {
	return _EthereumDIDRegistry.Contract.Owners(&_EthereumDIDRegistry.CallOpts, arg0)
}

// Owners is a free data retrieval call binding the contract method 0x022914a7.
//
// Solidity: function owners(address ) view returns(address)
func (_EthereumDIDRegistry *EthereumDIDRegistryCallerSession) Owners(arg0 common.Address) (common.Address, error) {
	return _EthereumDIDRegistry.Contract.Owners(&_EthereumDIDRegistry.CallOpts, arg0)
}

// ValidDelegate is a free data retrieval call binding the contract method 0x622b2a3c.
//
// Solidity: function validDelegate(address identity, bytes32 delegateType, address delegate) view returns(bool)
func (_EthereumDIDRegistry *EthereumDIDRegistryCaller) ValidDelegate(opts *bind.CallOpts, identity common.Address, delegateType [32]byte, delegate common.Address) (bool, error) {
	var out []interface{}
	err := _EthereumDIDRegistry.contract.Call(opts, &out, "validDelegate", identity, delegateType, delegate)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidDelegate is a free data retrieval call binding the contract method 0x622b2a3c.
//
// Solidity: function validDelegate(address identity, bytes32 delegateType, address delegate) view returns(bool)
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) ValidDelegate(identity common.Address, delegateType [32]byte, delegate common.Address) (bool, error) {
	return _EthereumDIDRegistry.Contract.ValidDelegate(&_EthereumDIDRegistry.CallOpts, identity, delegateType, delegate)
}

// ValidDelegate is a free data retrieval call binding the contract method 0x622b2a3c.
//
// Solidity: function validDelegate(address identity, bytes32 delegateType, address delegate) view returns(bool)
func (_EthereumDIDRegistry *EthereumDIDRegistryCallerSession) ValidDelegate(identity common.Address, delegateType [32]byte, delegate common.Address) (bool, error) {
	return _EthereumDIDRegistry.Contract.ValidDelegate(&_EthereumDIDRegistry.CallOpts, identity, delegateType, delegate)
}

// AddDelegate is a paid mutator transaction binding the contract method 0xa7068d66.
//
// Solidity: function addDelegate(address identity, bytes32 delegateType, address delegate, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) AddDelegate(opts *bind.TransactOpts, identity common.Address, delegateType [32]byte, delegate common.Address, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "addDelegate", identity, delegateType, delegate, validity)
}

// AddDelegate is a paid mutator transaction binding the contract method 0xa7068d66.
//
// Solidity: function addDelegate(address identity, bytes32 delegateType, address delegate, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) AddDelegate(identity common.Address, delegateType [32]byte, delegate common.Address, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.AddDelegate(&_EthereumDIDRegistry.TransactOpts, identity, delegateType, delegate, validity)
}

// AddDelegate is a paid mutator transaction binding the contract method 0xa7068d66.
//
// Solidity: function addDelegate(address identity, bytes32 delegateType, address delegate, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) AddDelegate(identity common.Address, delegateType [32]byte, delegate common.Address, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.AddDelegate(&_EthereumDIDRegistry.TransactOpts, identity, delegateType, delegate, validity)
}

// AddDelegateSigned is a paid mutator transaction binding the contract method 0x9c2c1b2b.
//
// Solidity: function addDelegateSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 delegateType, address delegate, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) AddDelegateSigned(opts *bind.TransactOpts, identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, delegateType [32]byte, delegate common.Address, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "addDelegateSigned", identity, sigV, sigR, sigS, delegateType, delegate, validity)
}

// AddDelegateSigned is a paid mutator transaction binding the contract method 0x9c2c1b2b.
//
// Solidity: function addDelegateSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 delegateType, address delegate, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) AddDelegateSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, delegateType [32]byte, delegate common.Address, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.AddDelegateSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, delegateType, delegate, validity)
}

// AddDelegateSigned is a paid mutator transaction binding the contract method 0x9c2c1b2b.
//
// Solidity: function addDelegateSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 delegateType, address delegate, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) AddDelegateSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, delegateType [32]byte, delegate common.Address, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.AddDelegateSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, delegateType, delegate, validity)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xf00d4b5d.
//
// Solidity: function changeOwner(address identity, address newOwner) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) ChangeOwner(opts *bind.TransactOpts, identity common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "changeOwner", identity, newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xf00d4b5d.
//
// Solidity: function changeOwner(address identity, address newOwner) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) ChangeOwner(identity common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.ChangeOwner(&_EthereumDIDRegistry.TransactOpts, identity, newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xf00d4b5d.
//
// Solidity: function changeOwner(address identity, address newOwner) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) ChangeOwner(identity common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.ChangeOwner(&_EthereumDIDRegistry.TransactOpts, identity, newOwner)
}

// ChangeOwnerSigned is a paid mutator transaction binding the contract method 0x240cf1fa.
//
// Solidity: function changeOwnerSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, address newOwner) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) ChangeOwnerSigned(opts *bind.TransactOpts, identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, newOwner common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "changeOwnerSigned", identity, sigV, sigR, sigS, newOwner)
}

// ChangeOwnerSigned is a paid mutator transaction binding the contract method 0x240cf1fa.
//
// Solidity: function changeOwnerSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, address newOwner) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) ChangeOwnerSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, newOwner common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.ChangeOwnerSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, newOwner)
}

// ChangeOwnerSigned is a paid mutator transaction binding the contract method 0x240cf1fa.
//
// Solidity: function changeOwnerSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, address newOwner) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) ChangeOwnerSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, newOwner common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.ChangeOwnerSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, newOwner)
}

// RevokeAttribute is a paid mutator transaction binding the contract method 0x00c023da.
//
// Solidity: function revokeAttribute(address identity, bytes32 name, bytes value) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) RevokeAttribute(opts *bind.TransactOpts, identity common.Address, name [32]byte, value []byte) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "revokeAttribute", identity, name, value)
}

// RevokeAttribute is a paid mutator transaction binding the contract method 0x00c023da.
//
// Solidity: function revokeAttribute(address identity, bytes32 name, bytes value) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) RevokeAttribute(identity common.Address, name [32]byte, value []byte) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.RevokeAttribute(&_EthereumDIDRegistry.TransactOpts, identity, name, value)
}

// RevokeAttribute is a paid mutator transaction binding the contract method 0x00c023da.
//
// Solidity: function revokeAttribute(address identity, bytes32 name, bytes value) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) RevokeAttribute(identity common.Address, name [32]byte, value []byte) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.RevokeAttribute(&_EthereumDIDRegistry.TransactOpts, identity, name, value)
}

// RevokeAttributeSigned is a paid mutator transaction binding the contract method 0xe476af5c.
//
// Solidity: function revokeAttributeSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 name, bytes value) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) RevokeAttributeSigned(opts *bind.TransactOpts, identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, name [32]byte, value []byte) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "revokeAttributeSigned", identity, sigV, sigR, sigS, name, value)
}

// RevokeAttributeSigned is a paid mutator transaction binding the contract method 0xe476af5c.
//
// Solidity: function revokeAttributeSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 name, bytes value) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) RevokeAttributeSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, name [32]byte, value []byte) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.RevokeAttributeSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, name, value)
}

// RevokeAttributeSigned is a paid mutator transaction binding the contract method 0xe476af5c.
//
// Solidity: function revokeAttributeSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 name, bytes value) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) RevokeAttributeSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, name [32]byte, value []byte) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.RevokeAttributeSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, name, value)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0x80b29f7c.
//
// Solidity: function revokeDelegate(address identity, bytes32 delegateType, address delegate) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) RevokeDelegate(opts *bind.TransactOpts, identity common.Address, delegateType [32]byte, delegate common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "revokeDelegate", identity, delegateType, delegate)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0x80b29f7c.
//
// Solidity: function revokeDelegate(address identity, bytes32 delegateType, address delegate) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) RevokeDelegate(identity common.Address, delegateType [32]byte, delegate common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.RevokeDelegate(&_EthereumDIDRegistry.TransactOpts, identity, delegateType, delegate)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0x80b29f7c.
//
// Solidity: function revokeDelegate(address identity, bytes32 delegateType, address delegate) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) RevokeDelegate(identity common.Address, delegateType [32]byte, delegate common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.RevokeDelegate(&_EthereumDIDRegistry.TransactOpts, identity, delegateType, delegate)
}

// RevokeDelegateSigned is a paid mutator transaction binding the contract method 0x93072684.
//
// Solidity: function revokeDelegateSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 delegateType, address delegate) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) RevokeDelegateSigned(opts *bind.TransactOpts, identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, delegateType [32]byte, delegate common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "revokeDelegateSigned", identity, sigV, sigR, sigS, delegateType, delegate)
}

// RevokeDelegateSigned is a paid mutator transaction binding the contract method 0x93072684.
//
// Solidity: function revokeDelegateSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 delegateType, address delegate) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) RevokeDelegateSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, delegateType [32]byte, delegate common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.RevokeDelegateSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, delegateType, delegate)
}

// RevokeDelegateSigned is a paid mutator transaction binding the contract method 0x93072684.
//
// Solidity: function revokeDelegateSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 delegateType, address delegate) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) RevokeDelegateSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, delegateType [32]byte, delegate common.Address) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.RevokeDelegateSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, delegateType, delegate)
}

// SetAttribute is a paid mutator transaction binding the contract method 0x7ad4b0a4.
//
// Solidity: function setAttribute(address identity, bytes32 name, bytes value, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) SetAttribute(opts *bind.TransactOpts, identity common.Address, name [32]byte, value []byte, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "setAttribute", identity, name, value, validity)
}

// SetAttribute is a paid mutator transaction binding the contract method 0x7ad4b0a4.
//
// Solidity: function setAttribute(address identity, bytes32 name, bytes value, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) SetAttribute(identity common.Address, name [32]byte, value []byte, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.SetAttribute(&_EthereumDIDRegistry.TransactOpts, identity, name, value, validity)
}

// SetAttribute is a paid mutator transaction binding the contract method 0x7ad4b0a4.
//
// Solidity: function setAttribute(address identity, bytes32 name, bytes value, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) SetAttribute(identity common.Address, name [32]byte, value []byte, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.SetAttribute(&_EthereumDIDRegistry.TransactOpts, identity, name, value, validity)
}

// SetAttributeSigned is a paid mutator transaction binding the contract method 0x123b5e98.
//
// Solidity: function setAttributeSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 name, bytes value, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactor) SetAttributeSigned(opts *bind.TransactOpts, identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, name [32]byte, value []byte, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.contract.Transact(opts, "setAttributeSigned", identity, sigV, sigR, sigS, name, value, validity)
}

// SetAttributeSigned is a paid mutator transaction binding the contract method 0x123b5e98.
//
// Solidity: function setAttributeSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 name, bytes value, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistrySession) SetAttributeSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, name [32]byte, value []byte, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.SetAttributeSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, name, value, validity)
}

// SetAttributeSigned is a paid mutator transaction binding the contract method 0x123b5e98.
//
// Solidity: function setAttributeSigned(address identity, uint8 sigV, bytes32 sigR, bytes32 sigS, bytes32 name, bytes value, uint256 validity) returns()
func (_EthereumDIDRegistry *EthereumDIDRegistryTransactorSession) SetAttributeSigned(identity common.Address, sigV uint8, sigR [32]byte, sigS [32]byte, name [32]byte, value []byte, validity *big.Int) (*types.Transaction, error) {
	return _EthereumDIDRegistry.Contract.SetAttributeSigned(&_EthereumDIDRegistry.TransactOpts, identity, sigV, sigR, sigS, name, value, validity)
}

// EthereumDIDRegistryDIDAttributeChangedIterator is returned from FilterDIDAttributeChanged and is used to iterate over the raw logs and unpacked data for DIDAttributeChanged events raised by the EthereumDIDRegistry contract.
type EthereumDIDRegistryDIDAttributeChangedIterator struct {
	Event *EthereumDIDRegistryDIDAttributeChanged // Event containing the contract specifics and raw log

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
func (it *EthereumDIDRegistryDIDAttributeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumDIDRegistryDIDAttributeChanged)
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
		it.Event = new(EthereumDIDRegistryDIDAttributeChanged)
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
func (it *EthereumDIDRegistryDIDAttributeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumDIDRegistryDIDAttributeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumDIDRegistryDIDAttributeChanged represents a DIDAttributeChanged event raised by the EthereumDIDRegistry contract.
type EthereumDIDRegistryDIDAttributeChanged struct {
	Identity       common.Address
	Name           [32]byte
	Value          []byte
	ValidTo        *big.Int
	PreviousChange *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDIDAttributeChanged is a free log retrieval operation binding the contract event 0x18ab6b2ae3d64306c00ce663125f2bd680e441a098de1635bd7ad8b0d44965e4.
//
// Solidity: event DIDAttributeChanged(address indexed identity, bytes32 name, bytes value, uint256 validTo, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) FilterDIDAttributeChanged(opts *bind.FilterOpts, identity []common.Address) (*EthereumDIDRegistryDIDAttributeChangedIterator, error) {

	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _EthereumDIDRegistry.contract.FilterLogs(opts, "DIDAttributeChanged", identityRule)
	if err != nil {
		return nil, err
	}
	return &EthereumDIDRegistryDIDAttributeChangedIterator{contract: _EthereumDIDRegistry.contract, event: "DIDAttributeChanged", logs: logs, sub: sub}, nil
}

// WatchDIDAttributeChanged is a free log subscription operation binding the contract event 0x18ab6b2ae3d64306c00ce663125f2bd680e441a098de1635bd7ad8b0d44965e4.
//
// Solidity: event DIDAttributeChanged(address indexed identity, bytes32 name, bytes value, uint256 validTo, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) WatchDIDAttributeChanged(opts *bind.WatchOpts, sink chan<- *EthereumDIDRegistryDIDAttributeChanged, identity []common.Address) (event.Subscription, error) {

	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _EthereumDIDRegistry.contract.WatchLogs(opts, "DIDAttributeChanged", identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumDIDRegistryDIDAttributeChanged)
				if err := _EthereumDIDRegistry.contract.UnpackLog(event, "DIDAttributeChanged", log); err != nil {
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

// ParseDIDAttributeChanged is a log parse operation binding the contract event 0x18ab6b2ae3d64306c00ce663125f2bd680e441a098de1635bd7ad8b0d44965e4.
//
// Solidity: event DIDAttributeChanged(address indexed identity, bytes32 name, bytes value, uint256 validTo, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) ParseDIDAttributeChanged(log types.Log) (*EthereumDIDRegistryDIDAttributeChanged, error) {
	event := new(EthereumDIDRegistryDIDAttributeChanged)
	if err := _EthereumDIDRegistry.contract.UnpackLog(event, "DIDAttributeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthereumDIDRegistryDIDDelegateChangedIterator is returned from FilterDIDDelegateChanged and is used to iterate over the raw logs and unpacked data for DIDDelegateChanged events raised by the EthereumDIDRegistry contract.
type EthereumDIDRegistryDIDDelegateChangedIterator struct {
	Event *EthereumDIDRegistryDIDDelegateChanged // Event containing the contract specifics and raw log

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
func (it *EthereumDIDRegistryDIDDelegateChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumDIDRegistryDIDDelegateChanged)
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
		it.Event = new(EthereumDIDRegistryDIDDelegateChanged)
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
func (it *EthereumDIDRegistryDIDDelegateChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumDIDRegistryDIDDelegateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumDIDRegistryDIDDelegateChanged represents a DIDDelegateChanged event raised by the EthereumDIDRegistry contract.
type EthereumDIDRegistryDIDDelegateChanged struct {
	Identity       common.Address
	DelegateType   [32]byte
	Delegate       common.Address
	ValidTo        *big.Int
	PreviousChange *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDIDDelegateChanged is a free log retrieval operation binding the contract event 0x5a5084339536bcab65f20799fcc58724588145ca054bd2be626174b27ba156f7.
//
// Solidity: event DIDDelegateChanged(address indexed identity, bytes32 delegateType, address delegate, uint256 validTo, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) FilterDIDDelegateChanged(opts *bind.FilterOpts, identity []common.Address) (*EthereumDIDRegistryDIDDelegateChangedIterator, error) {

	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _EthereumDIDRegistry.contract.FilterLogs(opts, "DIDDelegateChanged", identityRule)
	if err != nil {
		return nil, err
	}
	return &EthereumDIDRegistryDIDDelegateChangedIterator{contract: _EthereumDIDRegistry.contract, event: "DIDDelegateChanged", logs: logs, sub: sub}, nil
}

// WatchDIDDelegateChanged is a free log subscription operation binding the contract event 0x5a5084339536bcab65f20799fcc58724588145ca054bd2be626174b27ba156f7.
//
// Solidity: event DIDDelegateChanged(address indexed identity, bytes32 delegateType, address delegate, uint256 validTo, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) WatchDIDDelegateChanged(opts *bind.WatchOpts, sink chan<- *EthereumDIDRegistryDIDDelegateChanged, identity []common.Address) (event.Subscription, error) {

	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _EthereumDIDRegistry.contract.WatchLogs(opts, "DIDDelegateChanged", identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumDIDRegistryDIDDelegateChanged)
				if err := _EthereumDIDRegistry.contract.UnpackLog(event, "DIDDelegateChanged", log); err != nil {
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

// ParseDIDDelegateChanged is a log parse operation binding the contract event 0x5a5084339536bcab65f20799fcc58724588145ca054bd2be626174b27ba156f7.
//
// Solidity: event DIDDelegateChanged(address indexed identity, bytes32 delegateType, address delegate, uint256 validTo, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) ParseDIDDelegateChanged(log types.Log) (*EthereumDIDRegistryDIDDelegateChanged, error) {
	event := new(EthereumDIDRegistryDIDDelegateChanged)
	if err := _EthereumDIDRegistry.contract.UnpackLog(event, "DIDDelegateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthereumDIDRegistryDIDOwnerChangedIterator is returned from FilterDIDOwnerChanged and is used to iterate over the raw logs and unpacked data for DIDOwnerChanged events raised by the EthereumDIDRegistry contract.
type EthereumDIDRegistryDIDOwnerChangedIterator struct {
	Event *EthereumDIDRegistryDIDOwnerChanged // Event containing the contract specifics and raw log

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
func (it *EthereumDIDRegistryDIDOwnerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumDIDRegistryDIDOwnerChanged)
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
		it.Event = new(EthereumDIDRegistryDIDOwnerChanged)
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
func (it *EthereumDIDRegistryDIDOwnerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumDIDRegistryDIDOwnerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumDIDRegistryDIDOwnerChanged represents a DIDOwnerChanged event raised by the EthereumDIDRegistry contract.
type EthereumDIDRegistryDIDOwnerChanged struct {
	Identity       common.Address
	Owner          common.Address
	PreviousChange *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDIDOwnerChanged is a free log retrieval operation binding the contract event 0x38a5a6e68f30ed1ab45860a4afb34bcb2fc00f22ca462d249b8a8d40cda6f7a3.
//
// Solidity: event DIDOwnerChanged(address indexed identity, address owner, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) FilterDIDOwnerChanged(opts *bind.FilterOpts, identity []common.Address) (*EthereumDIDRegistryDIDOwnerChangedIterator, error) {

	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _EthereumDIDRegistry.contract.FilterLogs(opts, "DIDOwnerChanged", identityRule)
	if err != nil {
		return nil, err
	}
	return &EthereumDIDRegistryDIDOwnerChangedIterator{contract: _EthereumDIDRegistry.contract, event: "DIDOwnerChanged", logs: logs, sub: sub}, nil
}

// WatchDIDOwnerChanged is a free log subscription operation binding the contract event 0x38a5a6e68f30ed1ab45860a4afb34bcb2fc00f22ca462d249b8a8d40cda6f7a3.
//
// Solidity: event DIDOwnerChanged(address indexed identity, address owner, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) WatchDIDOwnerChanged(opts *bind.WatchOpts, sink chan<- *EthereumDIDRegistryDIDOwnerChanged, identity []common.Address) (event.Subscription, error) {

	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _EthereumDIDRegistry.contract.WatchLogs(opts, "DIDOwnerChanged", identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumDIDRegistryDIDOwnerChanged)
				if err := _EthereumDIDRegistry.contract.UnpackLog(event, "DIDOwnerChanged", log); err != nil {
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

// ParseDIDOwnerChanged is a log parse operation binding the contract event 0x38a5a6e68f30ed1ab45860a4afb34bcb2fc00f22ca462d249b8a8d40cda6f7a3.
//
// Solidity: event DIDOwnerChanged(address indexed identity, address owner, uint256 previousChange)
func (_EthereumDIDRegistry *EthereumDIDRegistryFilterer) ParseDIDOwnerChanged(log types.Log) (*EthereumDIDRegistryDIDOwnerChanged, error) {
	event := new(EthereumDIDRegistryDIDOwnerChanged)
	if err := _EthereumDIDRegistry.contract.UnpackLog(event, "DIDOwnerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
