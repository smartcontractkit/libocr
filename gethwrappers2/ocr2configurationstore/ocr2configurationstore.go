// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2configurationstore

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

// IOCR2ConfigurationStoreConfiguration is an auto generated low-level Go binding around an user-defined struct.
type IOCR2ConfigurationStoreConfiguration struct {
	ConfigCount           uint64
	Signers               []common.Address
	Transmitters          []common.Address
	OnchainConfig         []byte
	OffchainConfig        []byte
	OffchainConfigVersion uint64
	F                     uint8
}

// IOCR2ConfigurationStoreExtendedConfiguration is an auto generated low-level Go binding around an user-defined struct.
type IOCR2ConfigurationStoreExtendedConfiguration struct {
	BlockNumber     uint32
	ContractAddress common.Address
	ConfigDigest    [32]byte
	Configuration   IOCR2ConfigurationStoreConfiguration
}

// ConfigDigestUtilMetaData contains all meta data concerning the ConfigDigestUtil contract.
var ConfigDigestUtilMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x602d6037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea164736f6c6343000806000a",
}

// ConfigDigestUtilABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfigDigestUtilMetaData.ABI instead.
var ConfigDigestUtilABI = ConfigDigestUtilMetaData.ABI

// ConfigDigestUtilBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfigDigestUtilMetaData.Bin instead.
var ConfigDigestUtilBin = ConfigDigestUtilMetaData.Bin

// DeployConfigDigestUtil deploys a new Ethereum contract, binding an instance of ConfigDigestUtil to it.
func DeployConfigDigestUtil(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConfigDigestUtil, error) {
	parsed, err := ConfigDigestUtilMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfigDigestUtilBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfigDigestUtil{ConfigDigestUtilCaller: ConfigDigestUtilCaller{contract: contract}, ConfigDigestUtilTransactor: ConfigDigestUtilTransactor{contract: contract}, ConfigDigestUtilFilterer: ConfigDigestUtilFilterer{contract: contract}}, nil
}

// ConfigDigestUtil is an auto generated Go binding around an Ethereum contract.
type ConfigDigestUtil struct {
	ConfigDigestUtilCaller     // Read-only binding to the contract
	ConfigDigestUtilTransactor // Write-only binding to the contract
	ConfigDigestUtilFilterer   // Log filterer for contract events
}

// ConfigDigestUtilCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfigDigestUtilCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigDigestUtilTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfigDigestUtilTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigDigestUtilFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfigDigestUtilFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigDigestUtilSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfigDigestUtilSession struct {
	Contract     *ConfigDigestUtil // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConfigDigestUtilCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfigDigestUtilCallerSession struct {
	Contract *ConfigDigestUtilCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ConfigDigestUtilTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfigDigestUtilTransactorSession struct {
	Contract     *ConfigDigestUtilTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ConfigDigestUtilRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfigDigestUtilRaw struct {
	Contract *ConfigDigestUtil // Generic contract binding to access the raw methods on
}

// ConfigDigestUtilCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfigDigestUtilCallerRaw struct {
	Contract *ConfigDigestUtilCaller // Generic read-only contract binding to access the raw methods on
}

// ConfigDigestUtilTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfigDigestUtilTransactorRaw struct {
	Contract *ConfigDigestUtilTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfigDigestUtil creates a new instance of ConfigDigestUtil, bound to a specific deployed contract.
func NewConfigDigestUtil(address common.Address, backend bind.ContractBackend) (*ConfigDigestUtil, error) {
	contract, err := bindConfigDigestUtil(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfigDigestUtil{ConfigDigestUtilCaller: ConfigDigestUtilCaller{contract: contract}, ConfigDigestUtilTransactor: ConfigDigestUtilTransactor{contract: contract}, ConfigDigestUtilFilterer: ConfigDigestUtilFilterer{contract: contract}}, nil
}

// NewConfigDigestUtilCaller creates a new read-only instance of ConfigDigestUtil, bound to a specific deployed contract.
func NewConfigDigestUtilCaller(address common.Address, caller bind.ContractCaller) (*ConfigDigestUtilCaller, error) {
	contract, err := bindConfigDigestUtil(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfigDigestUtilCaller{contract: contract}, nil
}

// NewConfigDigestUtilTransactor creates a new write-only instance of ConfigDigestUtil, bound to a specific deployed contract.
func NewConfigDigestUtilTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfigDigestUtilTransactor, error) {
	contract, err := bindConfigDigestUtil(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfigDigestUtilTransactor{contract: contract}, nil
}

// NewConfigDigestUtilFilterer creates a new log filterer instance of ConfigDigestUtil, bound to a specific deployed contract.
func NewConfigDigestUtilFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfigDigestUtilFilterer, error) {
	contract, err := bindConfigDigestUtil(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfigDigestUtilFilterer{contract: contract}, nil
}

// bindConfigDigestUtil binds a generic wrapper to an already deployed contract.
func bindConfigDigestUtil(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfigDigestUtilMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfigDigestUtil *ConfigDigestUtilRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfigDigestUtil.Contract.ConfigDigestUtilCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfigDigestUtil *ConfigDigestUtilRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfigDigestUtil.Contract.ConfigDigestUtilTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfigDigestUtil *ConfigDigestUtilRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfigDigestUtil.Contract.ConfigDigestUtilTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfigDigestUtil *ConfigDigestUtilCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfigDigestUtil.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfigDigestUtil *ConfigDigestUtilTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfigDigestUtil.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfigDigestUtil *ConfigDigestUtilTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfigDigestUtil.Contract.contract.Transact(opts, method, params...)
}

// ConfirmedOwnerMetaData contains all meta data concerning the ConfirmedOwner contract.
var ConfirmedOwnerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161051638038061051683398101604081905261002f9161016f565b8060006001600160a01b03821661008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100bd576100bd816100c5565b50505061019f565b6001600160a01b03811633141561011e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561018157600080fd5b81516001600160a01b038116811461019857600080fd5b9392505050565b610368806101ae6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031e565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156102a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561033057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035457600080fd5b939250505056fea164736f6c6343000806000a",
}

// ConfirmedOwnerABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfirmedOwnerMetaData.ABI instead.
var ConfirmedOwnerABI = ConfirmedOwnerMetaData.ABI

// ConfirmedOwnerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfirmedOwnerMetaData.Bin instead.
var ConfirmedOwnerBin = ConfirmedOwnerMetaData.Bin

// DeployConfirmedOwner deploys a new Ethereum contract, binding an instance of ConfirmedOwner to it.
func DeployConfirmedOwner(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address) (common.Address, *types.Transaction, *ConfirmedOwner, error) {
	parsed, err := ConfirmedOwnerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfirmedOwnerBin), backend, newOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfirmedOwner{ConfirmedOwnerCaller: ConfirmedOwnerCaller{contract: contract}, ConfirmedOwnerTransactor: ConfirmedOwnerTransactor{contract: contract}, ConfirmedOwnerFilterer: ConfirmedOwnerFilterer{contract: contract}}, nil
}

// ConfirmedOwner is an auto generated Go binding around an Ethereum contract.
type ConfirmedOwner struct {
	ConfirmedOwnerCaller     // Read-only binding to the contract
	ConfirmedOwnerTransactor // Write-only binding to the contract
	ConfirmedOwnerFilterer   // Log filterer for contract events
}

// ConfirmedOwnerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfirmedOwnerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfirmedOwnerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfirmedOwnerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfirmedOwnerSession struct {
	Contract     *ConfirmedOwner   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConfirmedOwnerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfirmedOwnerCallerSession struct {
	Contract *ConfirmedOwnerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ConfirmedOwnerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfirmedOwnerTransactorSession struct {
	Contract     *ConfirmedOwnerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ConfirmedOwnerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfirmedOwnerRaw struct {
	Contract *ConfirmedOwner // Generic contract binding to access the raw methods on
}

// ConfirmedOwnerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfirmedOwnerCallerRaw struct {
	Contract *ConfirmedOwnerCaller // Generic read-only contract binding to access the raw methods on
}

// ConfirmedOwnerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfirmedOwnerTransactorRaw struct {
	Contract *ConfirmedOwnerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfirmedOwner creates a new instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwner(address common.Address, backend bind.ContractBackend) (*ConfirmedOwner, error) {
	contract, err := bindConfirmedOwner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwner{ConfirmedOwnerCaller: ConfirmedOwnerCaller{contract: contract}, ConfirmedOwnerTransactor: ConfirmedOwnerTransactor{contract: contract}, ConfirmedOwnerFilterer: ConfirmedOwnerFilterer{contract: contract}}, nil
}

// NewConfirmedOwnerCaller creates a new read-only instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerCaller(address common.Address, caller bind.ContractCaller) (*ConfirmedOwnerCaller, error) {
	contract, err := bindConfirmedOwner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerCaller{contract: contract}, nil
}

// NewConfirmedOwnerTransactor creates a new write-only instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfirmedOwnerTransactor, error) {
	contract, err := bindConfirmedOwner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerTransactor{contract: contract}, nil
}

// NewConfirmedOwnerFilterer creates a new log filterer instance of ConfirmedOwner, bound to a specific deployed contract.
func NewConfirmedOwnerFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfirmedOwnerFilterer, error) {
	contract, err := bindConfirmedOwner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerFilterer{contract: contract}, nil
}

// bindConfirmedOwner binds a generic wrapper to an already deployed contract.
func bindConfirmedOwner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfirmedOwnerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwner.Contract.ConfirmedOwnerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.ConfirmedOwnerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwner *ConfirmedOwnerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.ConfirmedOwnerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwner *ConfirmedOwnerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwner.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwner *ConfirmedOwnerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwner *ConfirmedOwnerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfirmedOwner.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerSession) Owner() (common.Address, error) {
	return _ConfirmedOwner.Contract.Owner(&_ConfirmedOwner.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwner *ConfirmedOwnerCallerSession) Owner() (common.Address, error) {
	return _ConfirmedOwner.Contract.Owner(&_ConfirmedOwner.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwner.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.AcceptOwnership(&_ConfirmedOwner.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.AcceptOwnership(&_ConfirmedOwner.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.TransferOwnership(&_ConfirmedOwner.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwner *ConfirmedOwnerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwner.Contract.TransferOwnership(&_ConfirmedOwner.TransactOpts, to)
}

// ConfirmedOwnerOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferRequestedIterator struct {
	Event *ConfirmedOwnerOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerOwnershipTransferRequested)
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
		it.Event = new(ConfirmedOwnerOwnershipTransferRequested)
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
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerOwnershipTransferRequestedIterator{contract: _ConfirmedOwner.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerOwnershipTransferRequested)
				if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) ParseOwnershipTransferRequested(log types.Log) (*ConfirmedOwnerOwnershipTransferRequested, error) {
	event := new(ConfirmedOwnerOwnershipTransferRequested)
	if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferredIterator struct {
	Event *ConfirmedOwnerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ConfirmedOwnerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerOwnershipTransferred)
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
		it.Event = new(ConfirmedOwnerOwnershipTransferred)
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
func (it *ConfirmedOwnerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerOwnershipTransferred represents a OwnershipTransferred event raised by the ConfirmedOwner contract.
type ConfirmedOwnerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerOwnershipTransferredIterator{contract: _ConfirmedOwner.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwner.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerOwnershipTransferred)
				if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwner *ConfirmedOwnerFilterer) ParseOwnershipTransferred(log types.Log) (*ConfirmedOwnerOwnershipTransferred, error) {
	event := new(ConfirmedOwnerOwnershipTransferred)
	if err := _ConfirmedOwner.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerWithProposalMetaData contains all meta data concerning the ConfirmedOwnerWithProposal contract.
var ConfirmedOwnerWithProposalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161053138038061053183398101604081905261002f91610187565b6001600160a01b03821661008a5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100ba576100ba816100c1565b50506101ba565b6001600160a01b03811633141561011a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610081565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b038116811461018257600080fd5b919050565b6000806040838503121561019a57600080fd5b6101a38361016b565b91506101b16020840161016b565b90509250929050565b610368806101c96000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031e565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156102a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561033057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035457600080fd5b939250505056fea164736f6c6343000806000a",
}

// ConfirmedOwnerWithProposalABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfirmedOwnerWithProposalMetaData.ABI instead.
var ConfirmedOwnerWithProposalABI = ConfirmedOwnerWithProposalMetaData.ABI

// ConfirmedOwnerWithProposalBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfirmedOwnerWithProposalMetaData.Bin instead.
var ConfirmedOwnerWithProposalBin = ConfirmedOwnerWithProposalMetaData.Bin

// DeployConfirmedOwnerWithProposal deploys a new Ethereum contract, binding an instance of ConfirmedOwnerWithProposal to it.
func DeployConfirmedOwnerWithProposal(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address, pendingOwner common.Address) (common.Address, *types.Transaction, *ConfirmedOwnerWithProposal, error) {
	parsed, err := ConfirmedOwnerWithProposalMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfirmedOwnerWithProposalBin), backend, newOwner, pendingOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfirmedOwnerWithProposal{ConfirmedOwnerWithProposalCaller: ConfirmedOwnerWithProposalCaller{contract: contract}, ConfirmedOwnerWithProposalTransactor: ConfirmedOwnerWithProposalTransactor{contract: contract}, ConfirmedOwnerWithProposalFilterer: ConfirmedOwnerWithProposalFilterer{contract: contract}}, nil
}

// ConfirmedOwnerWithProposal is an auto generated Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposal struct {
	ConfirmedOwnerWithProposalCaller     // Read-only binding to the contract
	ConfirmedOwnerWithProposalTransactor // Write-only binding to the contract
	ConfirmedOwnerWithProposalFilterer   // Log filterer for contract events
}

// ConfirmedOwnerWithProposalCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfirmedOwnerWithProposalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfirmedOwnerWithProposalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfirmedOwnerWithProposalSession struct {
	Contract     *ConfirmedOwnerWithProposal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ConfirmedOwnerWithProposalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfirmedOwnerWithProposalCallerSession struct {
	Contract *ConfirmedOwnerWithProposalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// ConfirmedOwnerWithProposalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfirmedOwnerWithProposalTransactorSession struct {
	Contract     *ConfirmedOwnerWithProposalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// ConfirmedOwnerWithProposalRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalRaw struct {
	Contract *ConfirmedOwnerWithProposal // Generic contract binding to access the raw methods on
}

// ConfirmedOwnerWithProposalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalCallerRaw struct {
	Contract *ConfirmedOwnerWithProposalCaller // Generic read-only contract binding to access the raw methods on
}

// ConfirmedOwnerWithProposalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfirmedOwnerWithProposalTransactorRaw struct {
	Contract *ConfirmedOwnerWithProposalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfirmedOwnerWithProposal creates a new instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposal(address common.Address, backend bind.ContractBackend) (*ConfirmedOwnerWithProposal, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposal{ConfirmedOwnerWithProposalCaller: ConfirmedOwnerWithProposalCaller{contract: contract}, ConfirmedOwnerWithProposalTransactor: ConfirmedOwnerWithProposalTransactor{contract: contract}, ConfirmedOwnerWithProposalFilterer: ConfirmedOwnerWithProposalFilterer{contract: contract}}, nil
}

// NewConfirmedOwnerWithProposalCaller creates a new read-only instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalCaller(address common.Address, caller bind.ContractCaller) (*ConfirmedOwnerWithProposalCaller, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalCaller{contract: contract}, nil
}

// NewConfirmedOwnerWithProposalTransactor creates a new write-only instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfirmedOwnerWithProposalTransactor, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalTransactor{contract: contract}, nil
}

// NewConfirmedOwnerWithProposalFilterer creates a new log filterer instance of ConfirmedOwnerWithProposal, bound to a specific deployed contract.
func NewConfirmedOwnerWithProposalFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfirmedOwnerWithProposalFilterer, error) {
	contract, err := bindConfirmedOwnerWithProposal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalFilterer{contract: contract}, nil
}

// bindConfirmedOwnerWithProposal binds a generic wrapper to an already deployed contract.
func bindConfirmedOwnerWithProposal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfirmedOwnerWithProposalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.ConfirmedOwnerWithProposalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfirmedOwnerWithProposal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfirmedOwnerWithProposal.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) Owner() (common.Address, error) {
	return _ConfirmedOwnerWithProposal.Contract.Owner(&_ConfirmedOwnerWithProposal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalCallerSession) Owner() (common.Address, error) {
	return _ConfirmedOwnerWithProposal.Contract.Owner(&_ConfirmedOwnerWithProposal.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.AcceptOwnership(&_ConfirmedOwnerWithProposal.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.AcceptOwnership(&_ConfirmedOwnerWithProposal.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.TransferOwnership(&_ConfirmedOwnerWithProposal.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _ConfirmedOwnerWithProposal.Contract.TransferOwnership(&_ConfirmedOwnerWithProposal.TransactOpts, to)
}

// ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator struct {
	Event *ConfirmedOwnerWithProposalOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
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
		it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
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
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerWithProposalOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalOwnershipTransferRequestedIterator{contract: _ConfirmedOwnerWithProposal.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerWithProposalOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
				if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) ParseOwnershipTransferRequested(log types.Log) (*ConfirmedOwnerWithProposalOwnershipTransferRequested, error) {
	event := new(ConfirmedOwnerWithProposalOwnershipTransferRequested)
	if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfirmedOwnerWithProposalOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferredIterator struct {
	Event *ConfirmedOwnerWithProposalOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferred)
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
		it.Event = new(ConfirmedOwnerWithProposalOwnershipTransferred)
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
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfirmedOwnerWithProposalOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfirmedOwnerWithProposalOwnershipTransferred represents a OwnershipTransferred event raised by the ConfirmedOwnerWithProposal contract.
type ConfirmedOwnerWithProposalOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfirmedOwnerWithProposalOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfirmedOwnerWithProposalOwnershipTransferredIterator{contract: _ConfirmedOwnerWithProposal.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfirmedOwnerWithProposalOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfirmedOwnerWithProposal.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfirmedOwnerWithProposalOwnershipTransferred)
				if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_ConfirmedOwnerWithProposal *ConfirmedOwnerWithProposalFilterer) ParseOwnershipTransferred(log types.Log) (*ConfirmedOwnerWithProposalOwnershipTransferred, error) {
	event := new(ConfirmedOwnerWithProposalOwnershipTransferred)
	if err := _ConfirmedOwnerWithProposal.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IOCR2ConfigurationStoreMetaData contains all meta data concerning the IOCR2ConfigurationStore contract.
var IOCR2ConfigurationStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"internalType\":\"structIOCR2ConfigurationStore.Configuration\",\"name\":\"configurationParams\",\"type\":\"tuple\"}],\"name\":\"addConfig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"latestConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"internalType\":\"structIOCR2ConfigurationStore.Configuration\",\"name\":\"configuration\",\"type\":\"tuple\"}],\"internalType\":\"structIOCR2ConfigurationStore.ExtendedConfiguration\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"readConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"internalType\":\"structIOCR2ConfigurationStore.Configuration\",\"name\":\"configuration\",\"type\":\"tuple\"}],\"internalType\":\"structIOCR2ConfigurationStore.ExtendedConfiguration\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IOCR2ConfigurationStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use IOCR2ConfigurationStoreMetaData.ABI instead.
var IOCR2ConfigurationStoreABI = IOCR2ConfigurationStoreMetaData.ABI

// IOCR2ConfigurationStore is an auto generated Go binding around an Ethereum contract.
type IOCR2ConfigurationStore struct {
	IOCR2ConfigurationStoreCaller     // Read-only binding to the contract
	IOCR2ConfigurationStoreTransactor // Write-only binding to the contract
	IOCR2ConfigurationStoreFilterer   // Log filterer for contract events
}

// IOCR2ConfigurationStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type IOCR2ConfigurationStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOCR2ConfigurationStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IOCR2ConfigurationStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOCR2ConfigurationStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IOCR2ConfigurationStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOCR2ConfigurationStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IOCR2ConfigurationStoreSession struct {
	Contract     *IOCR2ConfigurationStore // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IOCR2ConfigurationStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IOCR2ConfigurationStoreCallerSession struct {
	Contract *IOCR2ConfigurationStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// IOCR2ConfigurationStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IOCR2ConfigurationStoreTransactorSession struct {
	Contract     *IOCR2ConfigurationStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// IOCR2ConfigurationStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type IOCR2ConfigurationStoreRaw struct {
	Contract *IOCR2ConfigurationStore // Generic contract binding to access the raw methods on
}

// IOCR2ConfigurationStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IOCR2ConfigurationStoreCallerRaw struct {
	Contract *IOCR2ConfigurationStoreCaller // Generic read-only contract binding to access the raw methods on
}

// IOCR2ConfigurationStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IOCR2ConfigurationStoreTransactorRaw struct {
	Contract *IOCR2ConfigurationStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIOCR2ConfigurationStore creates a new instance of IOCR2ConfigurationStore, bound to a specific deployed contract.
func NewIOCR2ConfigurationStore(address common.Address, backend bind.ContractBackend) (*IOCR2ConfigurationStore, error) {
	contract, err := bindIOCR2ConfigurationStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IOCR2ConfigurationStore{IOCR2ConfigurationStoreCaller: IOCR2ConfigurationStoreCaller{contract: contract}, IOCR2ConfigurationStoreTransactor: IOCR2ConfigurationStoreTransactor{contract: contract}, IOCR2ConfigurationStoreFilterer: IOCR2ConfigurationStoreFilterer{contract: contract}}, nil
}

// NewIOCR2ConfigurationStoreCaller creates a new read-only instance of IOCR2ConfigurationStore, bound to a specific deployed contract.
func NewIOCR2ConfigurationStoreCaller(address common.Address, caller bind.ContractCaller) (*IOCR2ConfigurationStoreCaller, error) {
	contract, err := bindIOCR2ConfigurationStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IOCR2ConfigurationStoreCaller{contract: contract}, nil
}

// NewIOCR2ConfigurationStoreTransactor creates a new write-only instance of IOCR2ConfigurationStore, bound to a specific deployed contract.
func NewIOCR2ConfigurationStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*IOCR2ConfigurationStoreTransactor, error) {
	contract, err := bindIOCR2ConfigurationStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IOCR2ConfigurationStoreTransactor{contract: contract}, nil
}

// NewIOCR2ConfigurationStoreFilterer creates a new log filterer instance of IOCR2ConfigurationStore, bound to a specific deployed contract.
func NewIOCR2ConfigurationStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*IOCR2ConfigurationStoreFilterer, error) {
	contract, err := bindIOCR2ConfigurationStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IOCR2ConfigurationStoreFilterer{contract: contract}, nil
}

// bindIOCR2ConfigurationStore binds a generic wrapper to an already deployed contract.
func bindIOCR2ConfigurationStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IOCR2ConfigurationStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOCR2ConfigurationStore.Contract.IOCR2ConfigurationStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOCR2ConfigurationStore.Contract.IOCR2ConfigurationStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOCR2ConfigurationStore.Contract.IOCR2ConfigurationStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOCR2ConfigurationStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOCR2ConfigurationStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOCR2ConfigurationStore.Contract.contract.Transact(opts, method, params...)
}

// LatestConfig is a free data retrieval call binding the contract method 0x9d386827.
//
// Solidity: function latestConfig(address contractAddress) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreCaller) LatestConfig(opts *bind.CallOpts, contractAddress common.Address) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	var out []interface{}
	err := _IOCR2ConfigurationStore.contract.Call(opts, &out, "latestConfig", contractAddress)

	if err != nil {
		return *new(IOCR2ConfigurationStoreExtendedConfiguration), err
	}

	out0 := *abi.ConvertType(out[0], new(IOCR2ConfigurationStoreExtendedConfiguration)).(*IOCR2ConfigurationStoreExtendedConfiguration)

	return out0, err

}

// LatestConfig is a free data retrieval call binding the contract method 0x9d386827.
//
// Solidity: function latestConfig(address contractAddress) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreSession) LatestConfig(contractAddress common.Address) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	return _IOCR2ConfigurationStore.Contract.LatestConfig(&_IOCR2ConfigurationStore.CallOpts, contractAddress)
}

// LatestConfig is a free data retrieval call binding the contract method 0x9d386827.
//
// Solidity: function latestConfig(address contractAddress) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreCallerSession) LatestConfig(contractAddress common.Address) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	return _IOCR2ConfigurationStore.Contract.LatestConfig(&_IOCR2ConfigurationStore.CallOpts, contractAddress)
}

// ReadConfig is a free data retrieval call binding the contract method 0xbc4215dc.
//
// Solidity: function readConfig(bytes32 configDigest) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreCaller) ReadConfig(opts *bind.CallOpts, configDigest [32]byte) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	var out []interface{}
	err := _IOCR2ConfigurationStore.contract.Call(opts, &out, "readConfig", configDigest)

	if err != nil {
		return *new(IOCR2ConfigurationStoreExtendedConfiguration), err
	}

	out0 := *abi.ConvertType(out[0], new(IOCR2ConfigurationStoreExtendedConfiguration)).(*IOCR2ConfigurationStoreExtendedConfiguration)

	return out0, err

}

// ReadConfig is a free data retrieval call binding the contract method 0xbc4215dc.
//
// Solidity: function readConfig(bytes32 configDigest) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreSession) ReadConfig(configDigest [32]byte) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	return _IOCR2ConfigurationStore.Contract.ReadConfig(&_IOCR2ConfigurationStore.CallOpts, configDigest)
}

// ReadConfig is a free data retrieval call binding the contract method 0xbc4215dc.
//
// Solidity: function readConfig(bytes32 configDigest) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreCallerSession) ReadConfig(configDigest [32]byte) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	return _IOCR2ConfigurationStore.Contract.ReadConfig(&_IOCR2ConfigurationStore.CallOpts, configDigest)
}

// AddConfig is a paid mutator transaction binding the contract method 0x23e48b7d.
//
// Solidity: function addConfig((uint64,address[],address[],bytes,bytes,uint64,uint8) configurationParams) returns(bytes32)
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreTransactor) AddConfig(opts *bind.TransactOpts, configurationParams IOCR2ConfigurationStoreConfiguration) (*types.Transaction, error) {
	return _IOCR2ConfigurationStore.contract.Transact(opts, "addConfig", configurationParams)
}

// AddConfig is a paid mutator transaction binding the contract method 0x23e48b7d.
//
// Solidity: function addConfig((uint64,address[],address[],bytes,bytes,uint64,uint8) configurationParams) returns(bytes32)
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreSession) AddConfig(configurationParams IOCR2ConfigurationStoreConfiguration) (*types.Transaction, error) {
	return _IOCR2ConfigurationStore.Contract.AddConfig(&_IOCR2ConfigurationStore.TransactOpts, configurationParams)
}

// AddConfig is a paid mutator transaction binding the contract method 0x23e48b7d.
//
// Solidity: function addConfig((uint64,address[],address[],bytes,bytes,uint64,uint8) configurationParams) returns(bytes32)
func (_IOCR2ConfigurationStore *IOCR2ConfigurationStoreTransactorSession) AddConfig(configurationParams IOCR2ConfigurationStoreConfiguration) (*types.Transaction, error) {
	return _IOCR2ConfigurationStore.Contract.AddConfig(&_IOCR2ConfigurationStore.TransactOpts, configurationParams)
}

// OCR2AbstractMetaData contains all meta data concerning the OCR2Abstract contract.
var OCR2AbstractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// OCR2AbstractABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR2AbstractMetaData.ABI instead.
var OCR2AbstractABI = OCR2AbstractMetaData.ABI

// OCR2Abstract is an auto generated Go binding around an Ethereum contract.
type OCR2Abstract struct {
	OCR2AbstractCaller     // Read-only binding to the contract
	OCR2AbstractTransactor // Write-only binding to the contract
	OCR2AbstractFilterer   // Log filterer for contract events
}

// OCR2AbstractCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR2AbstractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR2AbstractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR2AbstractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2AbstractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR2AbstractSession struct {
	Contract     *OCR2Abstract     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OCR2AbstractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR2AbstractCallerSession struct {
	Contract *OCR2AbstractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// OCR2AbstractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR2AbstractTransactorSession struct {
	Contract     *OCR2AbstractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OCR2AbstractRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR2AbstractRaw struct {
	Contract *OCR2Abstract // Generic contract binding to access the raw methods on
}

// OCR2AbstractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR2AbstractCallerRaw struct {
	Contract *OCR2AbstractCaller // Generic read-only contract binding to access the raw methods on
}

// OCR2AbstractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR2AbstractTransactorRaw struct {
	Contract *OCR2AbstractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR2Abstract creates a new instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2Abstract(address common.Address, backend bind.ContractBackend) (*OCR2Abstract, error) {
	contract, err := bindOCR2Abstract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2Abstract{OCR2AbstractCaller: OCR2AbstractCaller{contract: contract}, OCR2AbstractTransactor: OCR2AbstractTransactor{contract: contract}, OCR2AbstractFilterer: OCR2AbstractFilterer{contract: contract}}, nil
}

// NewOCR2AbstractCaller creates a new read-only instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractCaller(address common.Address, caller bind.ContractCaller) (*OCR2AbstractCaller, error) {
	contract, err := bindOCR2Abstract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractCaller{contract: contract}, nil
}

// NewOCR2AbstractTransactor creates a new write-only instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2AbstractTransactor, error) {
	contract, err := bindOCR2Abstract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractTransactor{contract: contract}, nil
}

// NewOCR2AbstractFilterer creates a new log filterer instance of OCR2Abstract, bound to a specific deployed contract.
func NewOCR2AbstractFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2AbstractFilterer, error) {
	contract, err := bindOCR2Abstract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractFilterer{contract: contract}, nil
}

// bindOCR2Abstract binds a generic wrapper to an already deployed contract.
func bindOCR2Abstract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR2AbstractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Abstract *OCR2AbstractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Abstract.Contract.OCR2AbstractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Abstract *OCR2AbstractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.OCR2AbstractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Abstract *OCR2AbstractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.OCR2AbstractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2Abstract *OCR2AbstractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Abstract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2Abstract *OCR2AbstractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2Abstract *OCR2AbstractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.contract.Transact(opts, method, params...)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDetails(&_OCR2Abstract.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_OCR2Abstract *OCR2AbstractCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDetails(&_OCR2Abstract.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDigestAndEpoch(&_OCR2Abstract.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _OCR2Abstract.Contract.LatestConfigDigestAndEpoch(&_OCR2Abstract.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Abstract.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractSession) TypeAndVersion() (string, error) {
	return _OCR2Abstract.Contract.TypeAndVersion(&_OCR2Abstract.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2Abstract *OCR2AbstractCallerSession) TypeAndVersion() (string, error) {
	return _OCR2Abstract.Contract.TypeAndVersion(&_OCR2Abstract.CallOpts)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.SetConfig(&_OCR2Abstract.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_OCR2Abstract *OCR2AbstractTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.SetConfig(&_OCR2Abstract.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.Transmit(&_OCR2Abstract.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OCR2Abstract *OCR2AbstractTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Abstract.Contract.Transmit(&_OCR2Abstract.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// OCR2AbstractConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the OCR2Abstract contract.
type OCR2AbstractConfigSetIterator struct {
	Event *OCR2AbstractConfigSet // Event containing the contract specifics and raw log

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
func (it *OCR2AbstractConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AbstractConfigSet)
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
		it.Event = new(OCR2AbstractConfigSet)
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
func (it *OCR2AbstractConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AbstractConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AbstractConfigSet represents a ConfigSet event raised by the OCR2Abstract contract.
type OCR2AbstractConfigSet struct {
	BlockNumber           uint32
	ConfigDigest          [32]byte
	ConfigCount           uint64
	Signers               []common.Address
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 blockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2AbstractConfigSetIterator, error) {

	logs, sub, err := _OCR2Abstract.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractConfigSetIterator{contract: _OCR2Abstract.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 blockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AbstractConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Abstract.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AbstractConfigSet)
				if err := _OCR2Abstract.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 blockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_OCR2Abstract *OCR2AbstractFilterer) ParseConfigSet(log types.Log) (*OCR2AbstractConfigSet, error) {
	event := new(OCR2AbstractConfigSet)
	if err := _OCR2Abstract.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2AbstractTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the OCR2Abstract contract.
type OCR2AbstractTransmittedIterator struct {
	Event *OCR2AbstractTransmitted // Event containing the contract specifics and raw log

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
func (it *OCR2AbstractTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AbstractTransmitted)
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
		it.Event = new(OCR2AbstractTransmitted)
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
func (it *OCR2AbstractTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2AbstractTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2AbstractTransmitted represents a Transmitted event raised by the OCR2Abstract contract.
type OCR2AbstractTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2AbstractTransmittedIterator, error) {

	logs, sub, err := _OCR2Abstract.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2AbstractTransmittedIterator{contract: _OCR2Abstract.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2AbstractTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2Abstract.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2AbstractTransmitted)
				if err := _OCR2Abstract.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_OCR2Abstract *OCR2AbstractFilterer) ParseTransmitted(log types.Log) (*OCR2AbstractTransmitted, error) {
	event := new(OCR2AbstractTransmitted)
	if err := _OCR2Abstract.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2ConfigurationStoreMetaData contains all meta data concerning the OCR2ConfigurationStore contract.
var OCR2ConfigurationStoreMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"internalType\":\"structIOCR2ConfigurationStore.Configuration\",\"name\":\"configuration\",\"type\":\"tuple\"}],\"name\":\"addConfig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"latestConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"internalType\":\"structIOCR2ConfigurationStore.Configuration\",\"name\":\"configuration\",\"type\":\"tuple\"}],\"internalType\":\"structIOCR2ConfigurationStore.ExtendedConfiguration\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"readConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"internalType\":\"structIOCR2ConfigurationStore.Configuration\",\"name\":\"configuration\",\"type\":\"tuple\"}],\"internalType\":\"structIOCR2ConfigurationStore.ExtendedConfiguration\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_configurations\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"}],\"internalType\":\"structIOCR2ConfigurationStore.Configuration\",\"name\":\"configuration\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"s_latestConfigurationDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610149565b6001600160a01b0381163314156100f85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611816806101586000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638da5cb5b11610076578063bc4215dc1161005b578063bc4215dc14610183578063dc2e47cd14610196578063f2fde38b146101b957600080fd5b80638da5cb5b1461013b5780639d3868271461016357600080fd5b8063181f5a77146100a857806323e48b7d146100f0578063586505cb1461011157806379ba509714610131575b600080fd5b604080518082018252601c81527f4f435232436f6e66696775726174696f6e53746f726520312e302e3000000000602082015290516100e791906113f0565b60405180910390f35b6101036100fe366004611225565b6101cc565b6040519081526020016100e7565b61010361011f3660046111ea565b60036020526000908152604090205481565b6101396104e8565b005b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e7565b6101766101713660046111ea565b6105ea565b6040516100e79190611403565b61017661019136600461120c565b6108f4565b6101a96101a436600461120c565b610aa4565b6040516100e79493929190611505565b6101396101c73660046111ea565b610d32565b60008061030846336101e16020870187611260565b6101ee6020880188611550565b8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525061022d925050506040890189611550565b8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525061026f9250505060e08a0160c08b0161127b565b61027c60608b018b6115bf565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506102c19250505060c08c0160a08d01611260565b6102ce60808d018d6115bf565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610d4692505050565b905060405180608001604052804363ffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020018461034d9061169c565b9052600082815260026020818152604092839020845181548684015173ffffffffffffffffffffffffffffffffffffffff16640100000000027fffffffffffffffff00000000000000000000000000000000000000000000000090911663ffffffff9092169190911717815592840151600184015560608401518051928401805467ffffffffffffffff9094167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090941693909317835580820151805191939261041f92600387019290910190610f6a565b506040820151805161043b916002840191602090910190610f6a565b5060608201518051610457916003840191602090910190610ff4565b5060808201518051610473916004840191602090910190610ff4565b5060a08201516005909101805460c09093015160ff1668010000000000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000000000000090931667ffffffffffffffff90921691909117919091179055505033600090815260036020526040902081905592915050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461056e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b604080516080808201835260008083526020808401829052838501829052845160e081018652828152606091810182905294850181905280850181905291840182905260a0840181905260c084015281019190915273ffffffffffffffffffffffffffffffffffffffff80831660009081526003602081815260408084205484526002808352938190208151608081018352815463ffffffff81168252640100000000900490961686840152600181015486830152815160e081018352948101805467ffffffffffffffff168652938101805483518186028101860190945280845291956060880195909490938582019390929183018282801561072457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116106f9575b505050505081526020016002820180548060200260200160405190810160405280929190818152602001828054801561079357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610768575b505050505081526020016003820180546107ac90611786565b80601f01602080910402602001604051908101604052809291908181526020018280546107d890611786565b80156108255780601f106107fa57610100808354040283529160200191610825565b820191906000526020600020905b81548152906001019060200180831161080857829003601f168201915b5050505050815260200160048201805461083e90611786565b80601f016020809104026020016040519081016040528092919081815260200182805461086a90611786565b80156108b75780601f1061088c576101008083540402835291602001916108b7565b820191906000526020600020905b81548152906001019060200180831161089a57829003601f168201915b50505091835250506005919091015467ffffffffffffffff8116602083015268010000000000000000900460ff1660409091015290525092915050565b604080516080808201835260008083526020808401829052838501829052845160e081018652828152606091810182905294850181905280850181905291840182905260a0840181905260c08401528101919091526000828152600260208181526040928390208351608081018552815463ffffffff81168252640100000000900473ffffffffffffffffffffffffffffffffffffffff1681840152600182015481860152845160e081018652938201805467ffffffffffffffff168552600383018054875181870281018701909852808852929693956060880195909492938582019392918301828280156107245760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116106f95750505050508152602001600282018054806020026020016040519081016040528092919081815260200182805480156107935760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161076857505050505081526020016003820180546107ac90611786565b60026020818152600092835260409283902080546001820154855160e081018752948301805467ffffffffffffffff16865260038401805488518188028101880190995280895263ffffffff85169864010000000090950473ffffffffffffffffffffffffffffffffffffffff169793969394929385810193929190830182828015610b6657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610b3b575b5050505050815260200160028201805480602002602001604051908101604052809291908181526020018280548015610bd557602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610baa575b50505050508152602001600382018054610bee90611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610c1a90611786565b8015610c675780601f10610c3c57610100808354040283529160200191610c67565b820191906000526020600020905b815481529060010190602001808311610c4a57829003601f168201915b50505050508152602001600482018054610c8090611786565b80601f0160208091040260200160405190810160405280929190818152602001828054610cac90611786565b8015610cf95780601f10610cce57610100808354040283529160200191610cf9565b820191906000526020600020905b815481529060010190602001808311610cdc57829003601f168201915b50505091835250506005919091015467ffffffffffffffff8116602083015268010000000000000000900460ff16604090910152905084565b610d3a610df1565b610d4381610e74565b50565b6000808a8a8a8a8a8a8a8a8a604051602001610d6a99989796959493929190611460565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e72576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610565565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610ef4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610565565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610fe4579160200282015b82811115610fe457825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610f8a565b50610ff0929150611068565b5090565b82805461100090611786565b90600052602060002090601f0160209004810192826110225760008555610fe4565b82601f1061103b57805160ff1916838001178555610fe4565b82800160010185558215610fe4579182015b82811115610fe457825182559160200191906001019061104d565b5b80821115610ff05760008155600101611069565b803573ffffffffffffffffffffffffffffffffffffffff811681146110a157600080fd5b919050565b600082601f8301126110b757600080fd5b8135602067ffffffffffffffff8211156110d3576110d36117da565b8160051b6110e282820161164d565b8381528281019086840183880185018910156110fd57600080fd5b600093505b85841015611127576111138161107d565b835260019390930192918401918401611102565b50979650505050505050565b600082601f83011261114457600080fd5b813567ffffffffffffffff81111561115e5761115e6117da565b61118f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161164d565b8181528460208386010111156111a457600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff811681146110a157600080fd5b803560ff811681146110a157600080fd5b6000602082840312156111fc57600080fd5b6112058261107d565b9392505050565b60006020828403121561121e57600080fd5b5035919050565b60006020828403121561123757600080fd5b813567ffffffffffffffff81111561124e57600080fd5b820160e0818503121561120557600080fd5b60006020828403121561127257600080fd5b611205826111c1565b60006020828403121561128d57600080fd5b611205826111d9565b600081518084526020808501945080840160005b838110156112dc57815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016112aa565b509495945050505050565b6000815180845260005b8181101561130d576020818501810151868301820152016112f1565b8181111561131f576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600067ffffffffffffffff808351168452602083015160e0602086015261137c60e0860182611296565b9050604084015185820360408701526113958282611296565b915050606084015185820360608701526113af82826112e7565b915050608084015185820360808701526113c982826112e7565b9150508160a08501511660a086015260ff60c08501511660c0860152809250505092915050565b60208152600061120560208301846112e7565b6020815263ffffffff825116602082015273ffffffffffffffffffffffffffffffffffffffff6020830151166040820152604082015160608201526000606083015160808084015261145860a0840182611352565b949350505050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526114a78285018b611296565b915083820360808501526114bb828a611296565b915060ff881660a085015283820360c08501526114d882886112e7565b90861660e085015283810361010085015290506114f581856112e7565b9c9b505050505050505050505050565b63ffffffff8516815273ffffffffffffffffffffffffffffffffffffffff841660208201528260408201526080606082015260006115466080830184611352565b9695505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261158557600080fd5b83018035915067ffffffffffffffff8211156115a057600080fd5b6020019150600581901b36038213156115b857600080fd5b9250929050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126115f457600080fd5b83018035915067ffffffffffffffff82111561160f57600080fd5b6020019150368190038213156115b857600080fd5b60405160e0810167ffffffffffffffff81118282101715611647576116476117da565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611694576116946117da565b604052919050565b600060e082360312156116ae57600080fd5b6116b6611624565b6116bf836111c1565b8152602083013567ffffffffffffffff808211156116dc57600080fd5b6116e8368387016110a6565b6020840152604085013591508082111561170157600080fd5b61170d368387016110a6565b6040840152606085013591508082111561172657600080fd5b61173236838701611133565b6060840152608085013591508082111561174b57600080fd5b5061175836828601611133565b60808301525061176a60a084016111c1565b60a082015261177b60c084016111d9565b60c082015292915050565b600181811c9082168061179a57607f821691505b602082108114156117d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

// OCR2ConfigurationStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR2ConfigurationStoreMetaData.ABI instead.
var OCR2ConfigurationStoreABI = OCR2ConfigurationStoreMetaData.ABI

// OCR2ConfigurationStoreBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR2ConfigurationStoreMetaData.Bin instead.
var OCR2ConfigurationStoreBin = OCR2ConfigurationStoreMetaData.Bin

// DeployOCR2ConfigurationStore deploys a new Ethereum contract, binding an instance of OCR2ConfigurationStore to it.
func DeployOCR2ConfigurationStore(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR2ConfigurationStore, error) {
	parsed, err := OCR2ConfigurationStoreMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2ConfigurationStoreBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2ConfigurationStore{OCR2ConfigurationStoreCaller: OCR2ConfigurationStoreCaller{contract: contract}, OCR2ConfigurationStoreTransactor: OCR2ConfigurationStoreTransactor{contract: contract}, OCR2ConfigurationStoreFilterer: OCR2ConfigurationStoreFilterer{contract: contract}}, nil
}

// OCR2ConfigurationStore is an auto generated Go binding around an Ethereum contract.
type OCR2ConfigurationStore struct {
	OCR2ConfigurationStoreCaller     // Read-only binding to the contract
	OCR2ConfigurationStoreTransactor // Write-only binding to the contract
	OCR2ConfigurationStoreFilterer   // Log filterer for contract events
}

// OCR2ConfigurationStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR2ConfigurationStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2ConfigurationStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR2ConfigurationStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2ConfigurationStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR2ConfigurationStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR2ConfigurationStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR2ConfigurationStoreSession struct {
	Contract     *OCR2ConfigurationStore // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OCR2ConfigurationStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR2ConfigurationStoreCallerSession struct {
	Contract *OCR2ConfigurationStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// OCR2ConfigurationStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR2ConfigurationStoreTransactorSession struct {
	Contract     *OCR2ConfigurationStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// OCR2ConfigurationStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR2ConfigurationStoreRaw struct {
	Contract *OCR2ConfigurationStore // Generic contract binding to access the raw methods on
}

// OCR2ConfigurationStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR2ConfigurationStoreCallerRaw struct {
	Contract *OCR2ConfigurationStoreCaller // Generic read-only contract binding to access the raw methods on
}

// OCR2ConfigurationStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR2ConfigurationStoreTransactorRaw struct {
	Contract *OCR2ConfigurationStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR2ConfigurationStore creates a new instance of OCR2ConfigurationStore, bound to a specific deployed contract.
func NewOCR2ConfigurationStore(address common.Address, backend bind.ContractBackend) (*OCR2ConfigurationStore, error) {
	contract, err := bindOCR2ConfigurationStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2ConfigurationStore{OCR2ConfigurationStoreCaller: OCR2ConfigurationStoreCaller{contract: contract}, OCR2ConfigurationStoreTransactor: OCR2ConfigurationStoreTransactor{contract: contract}, OCR2ConfigurationStoreFilterer: OCR2ConfigurationStoreFilterer{contract: contract}}, nil
}

// NewOCR2ConfigurationStoreCaller creates a new read-only instance of OCR2ConfigurationStore, bound to a specific deployed contract.
func NewOCR2ConfigurationStoreCaller(address common.Address, caller bind.ContractCaller) (*OCR2ConfigurationStoreCaller, error) {
	contract, err := bindOCR2ConfigurationStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2ConfigurationStoreCaller{contract: contract}, nil
}

// NewOCR2ConfigurationStoreTransactor creates a new write-only instance of OCR2ConfigurationStore, bound to a specific deployed contract.
func NewOCR2ConfigurationStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2ConfigurationStoreTransactor, error) {
	contract, err := bindOCR2ConfigurationStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2ConfigurationStoreTransactor{contract: contract}, nil
}

// NewOCR2ConfigurationStoreFilterer creates a new log filterer instance of OCR2ConfigurationStore, bound to a specific deployed contract.
func NewOCR2ConfigurationStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2ConfigurationStoreFilterer, error) {
	contract, err := bindOCR2ConfigurationStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2ConfigurationStoreFilterer{contract: contract}, nil
}

// bindOCR2ConfigurationStore binds a generic wrapper to an already deployed contract.
func bindOCR2ConfigurationStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR2ConfigurationStoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2ConfigurationStore.Contract.OCR2ConfigurationStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.OCR2ConfigurationStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.OCR2ConfigurationStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2ConfigurationStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.contract.Transact(opts, method, params...)
}

// LatestConfig is a free data retrieval call binding the contract method 0x9d386827.
//
// Solidity: function latestConfig(address contractAddress) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCaller) LatestConfig(opts *bind.CallOpts, contractAddress common.Address) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	var out []interface{}
	err := _OCR2ConfigurationStore.contract.Call(opts, &out, "latestConfig", contractAddress)

	if err != nil {
		return *new(IOCR2ConfigurationStoreExtendedConfiguration), err
	}

	out0 := *abi.ConvertType(out[0], new(IOCR2ConfigurationStoreExtendedConfiguration)).(*IOCR2ConfigurationStoreExtendedConfiguration)

	return out0, err

}

// LatestConfig is a free data retrieval call binding the contract method 0x9d386827.
//
// Solidity: function latestConfig(address contractAddress) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) LatestConfig(contractAddress common.Address) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	return _OCR2ConfigurationStore.Contract.LatestConfig(&_OCR2ConfigurationStore.CallOpts, contractAddress)
}

// LatestConfig is a free data retrieval call binding the contract method 0x9d386827.
//
// Solidity: function latestConfig(address contractAddress) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCallerSession) LatestConfig(contractAddress common.Address) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	return _OCR2ConfigurationStore.Contract.LatestConfig(&_OCR2ConfigurationStore.CallOpts, contractAddress)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2ConfigurationStore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) Owner() (common.Address, error) {
	return _OCR2ConfigurationStore.Contract.Owner(&_OCR2ConfigurationStore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCallerSession) Owner() (common.Address, error) {
	return _OCR2ConfigurationStore.Contract.Owner(&_OCR2ConfigurationStore.CallOpts)
}

// ReadConfig is a free data retrieval call binding the contract method 0xbc4215dc.
//
// Solidity: function readConfig(bytes32 configDigest) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCaller) ReadConfig(opts *bind.CallOpts, configDigest [32]byte) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	var out []interface{}
	err := _OCR2ConfigurationStore.contract.Call(opts, &out, "readConfig", configDigest)

	if err != nil {
		return *new(IOCR2ConfigurationStoreExtendedConfiguration), err
	}

	out0 := *abi.ConvertType(out[0], new(IOCR2ConfigurationStoreExtendedConfiguration)).(*IOCR2ConfigurationStoreExtendedConfiguration)

	return out0, err

}

// ReadConfig is a free data retrieval call binding the contract method 0xbc4215dc.
//
// Solidity: function readConfig(bytes32 configDigest) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) ReadConfig(configDigest [32]byte) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	return _OCR2ConfigurationStore.Contract.ReadConfig(&_OCR2ConfigurationStore.CallOpts, configDigest)
}

// ReadConfig is a free data retrieval call binding the contract method 0xbc4215dc.
//
// Solidity: function readConfig(bytes32 configDigest) view returns((uint32,address,bytes32,(uint64,address[],address[],bytes,bytes,uint64,uint8)))
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCallerSession) ReadConfig(configDigest [32]byte) (IOCR2ConfigurationStoreExtendedConfiguration, error) {
	return _OCR2ConfigurationStore.Contract.ReadConfig(&_OCR2ConfigurationStore.CallOpts, configDigest)
}

// SConfigurations is a free data retrieval call binding the contract method 0xdc2e47cd.
//
// Solidity: function s_configurations(bytes32 ) view returns(uint32 blockNumber, address contractAddress, bytes32 configDigest, (uint64,address[],address[],bytes,bytes,uint64,uint8) configuration)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCaller) SConfigurations(opts *bind.CallOpts, arg0 [32]byte) (struct {
	BlockNumber     uint32
	ContractAddress common.Address
	ConfigDigest    [32]byte
	Configuration   IOCR2ConfigurationStoreConfiguration
}, error) {
	var out []interface{}
	err := _OCR2ConfigurationStore.contract.Call(opts, &out, "s_configurations", arg0)

	outstruct := new(struct {
		BlockNumber     uint32
		ContractAddress common.Address
		ConfigDigest    [32]byte
		Configuration   IOCR2ConfigurationStoreConfiguration
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockNumber = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ContractAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Configuration = *abi.ConvertType(out[3], new(IOCR2ConfigurationStoreConfiguration)).(*IOCR2ConfigurationStoreConfiguration)

	return *outstruct, err

}

// SConfigurations is a free data retrieval call binding the contract method 0xdc2e47cd.
//
// Solidity: function s_configurations(bytes32 ) view returns(uint32 blockNumber, address contractAddress, bytes32 configDigest, (uint64,address[],address[],bytes,bytes,uint64,uint8) configuration)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) SConfigurations(arg0 [32]byte) (struct {
	BlockNumber     uint32
	ContractAddress common.Address
	ConfigDigest    [32]byte
	Configuration   IOCR2ConfigurationStoreConfiguration
}, error) {
	return _OCR2ConfigurationStore.Contract.SConfigurations(&_OCR2ConfigurationStore.CallOpts, arg0)
}

// SConfigurations is a free data retrieval call binding the contract method 0xdc2e47cd.
//
// Solidity: function s_configurations(bytes32 ) view returns(uint32 blockNumber, address contractAddress, bytes32 configDigest, (uint64,address[],address[],bytes,bytes,uint64,uint8) configuration)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCallerSession) SConfigurations(arg0 [32]byte) (struct {
	BlockNumber     uint32
	ContractAddress common.Address
	ConfigDigest    [32]byte
	Configuration   IOCR2ConfigurationStoreConfiguration
}, error) {
	return _OCR2ConfigurationStore.Contract.SConfigurations(&_OCR2ConfigurationStore.CallOpts, arg0)
}

// SLatestConfigurationDigest is a free data retrieval call binding the contract method 0x586505cb.
//
// Solidity: function s_latestConfigurationDigest(address ) view returns(bytes32)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCaller) SLatestConfigurationDigest(opts *bind.CallOpts, arg0 common.Address) ([32]byte, error) {
	var out []interface{}
	err := _OCR2ConfigurationStore.contract.Call(opts, &out, "s_latestConfigurationDigest", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SLatestConfigurationDigest is a free data retrieval call binding the contract method 0x586505cb.
//
// Solidity: function s_latestConfigurationDigest(address ) view returns(bytes32)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) SLatestConfigurationDigest(arg0 common.Address) ([32]byte, error) {
	return _OCR2ConfigurationStore.Contract.SLatestConfigurationDigest(&_OCR2ConfigurationStore.CallOpts, arg0)
}

// SLatestConfigurationDigest is a free data retrieval call binding the contract method 0x586505cb.
//
// Solidity: function s_latestConfigurationDigest(address ) view returns(bytes32)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCallerSession) SLatestConfigurationDigest(arg0 common.Address) ([32]byte, error) {
	return _OCR2ConfigurationStore.Contract.SLatestConfigurationDigest(&_OCR2ConfigurationStore.CallOpts, arg0)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2ConfigurationStore.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) TypeAndVersion() (string, error) {
	return _OCR2ConfigurationStore.Contract.TypeAndVersion(&_OCR2ConfigurationStore.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreCallerSession) TypeAndVersion() (string, error) {
	return _OCR2ConfigurationStore.Contract.TypeAndVersion(&_OCR2ConfigurationStore.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.AcceptOwnership(&_OCR2ConfigurationStore.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.AcceptOwnership(&_OCR2ConfigurationStore.TransactOpts)
}

// AddConfig is a paid mutator transaction binding the contract method 0x23e48b7d.
//
// Solidity: function addConfig((uint64,address[],address[],bytes,bytes,uint64,uint8) configuration) returns(bytes32)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreTransactor) AddConfig(opts *bind.TransactOpts, configuration IOCR2ConfigurationStoreConfiguration) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.contract.Transact(opts, "addConfig", configuration)
}

// AddConfig is a paid mutator transaction binding the contract method 0x23e48b7d.
//
// Solidity: function addConfig((uint64,address[],address[],bytes,bytes,uint64,uint8) configuration) returns(bytes32)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) AddConfig(configuration IOCR2ConfigurationStoreConfiguration) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.AddConfig(&_OCR2ConfigurationStore.TransactOpts, configuration)
}

// AddConfig is a paid mutator transaction binding the contract method 0x23e48b7d.
//
// Solidity: function addConfig((uint64,address[],address[],bytes,bytes,uint64,uint8) configuration) returns(bytes32)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreTransactorSession) AddConfig(configuration IOCR2ConfigurationStoreConfiguration) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.AddConfig(&_OCR2ConfigurationStore.TransactOpts, configuration)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.TransferOwnership(&_OCR2ConfigurationStore.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2ConfigurationStore.Contract.TransferOwnership(&_OCR2ConfigurationStore.TransactOpts, to)
}

// OCR2ConfigurationStoreOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OCR2ConfigurationStore contract.
type OCR2ConfigurationStoreOwnershipTransferRequestedIterator struct {
	Event *OCR2ConfigurationStoreOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *OCR2ConfigurationStoreOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2ConfigurationStoreOwnershipTransferRequested)
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
		it.Event = new(OCR2ConfigurationStoreOwnershipTransferRequested)
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
func (it *OCR2ConfigurationStoreOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2ConfigurationStoreOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2ConfigurationStoreOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OCR2ConfigurationStore contract.
type OCR2ConfigurationStoreOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2ConfigurationStoreOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2ConfigurationStore.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2ConfigurationStoreOwnershipTransferRequestedIterator{contract: _OCR2ConfigurationStore.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2ConfigurationStoreOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2ConfigurationStore.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2ConfigurationStoreOwnershipTransferRequested)
				if err := _OCR2ConfigurationStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2ConfigurationStoreOwnershipTransferRequested, error) {
	event := new(OCR2ConfigurationStoreOwnershipTransferRequested)
	if err := _OCR2ConfigurationStore.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OCR2ConfigurationStoreOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OCR2ConfigurationStore contract.
type OCR2ConfigurationStoreOwnershipTransferredIterator struct {
	Event *OCR2ConfigurationStoreOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OCR2ConfigurationStoreOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2ConfigurationStoreOwnershipTransferred)
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
		it.Event = new(OCR2ConfigurationStoreOwnershipTransferred)
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
func (it *OCR2ConfigurationStoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OCR2ConfigurationStoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OCR2ConfigurationStoreOwnershipTransferred represents a OwnershipTransferred event raised by the OCR2ConfigurationStore contract.
type OCR2ConfigurationStoreOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2ConfigurationStoreOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2ConfigurationStore.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2ConfigurationStoreOwnershipTransferredIterator{contract: _OCR2ConfigurationStore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2ConfigurationStoreOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2ConfigurationStore.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OCR2ConfigurationStoreOwnershipTransferred)
				if err := _OCR2ConfigurationStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OCR2ConfigurationStore *OCR2ConfigurationStoreFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2ConfigurationStoreOwnershipTransferred, error) {
	event := new(OCR2ConfigurationStoreOwnershipTransferred)
	if err := _OCR2ConfigurationStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnableInterfaceMetaData contains all meta data concerning the OwnableInterface contract.
var OwnableInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OwnableInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableInterfaceMetaData.ABI instead.
var OwnableInterfaceABI = OwnableInterfaceMetaData.ABI

// OwnableInterface is an auto generated Go binding around an Ethereum contract.
type OwnableInterface struct {
	OwnableInterfaceCaller     // Read-only binding to the contract
	OwnableInterfaceTransactor // Write-only binding to the contract
	OwnableInterfaceFilterer   // Log filterer for contract events
}

// OwnableInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableInterfaceSession struct {
	Contract     *OwnableInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableInterfaceCallerSession struct {
	Contract *OwnableInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// OwnableInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableInterfaceTransactorSession struct {
	Contract     *OwnableInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// OwnableInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableInterfaceRaw struct {
	Contract *OwnableInterface // Generic contract binding to access the raw methods on
}

// OwnableInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableInterfaceCallerRaw struct {
	Contract *OwnableInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableInterfaceTransactorRaw struct {
	Contract *OwnableInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnableInterface creates a new instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterface(address common.Address, backend bind.ContractBackend) (*OwnableInterface, error) {
	contract, err := bindOwnableInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnableInterface{OwnableInterfaceCaller: OwnableInterfaceCaller{contract: contract}, OwnableInterfaceTransactor: OwnableInterfaceTransactor{contract: contract}, OwnableInterfaceFilterer: OwnableInterfaceFilterer{contract: contract}}, nil
}

// NewOwnableInterfaceCaller creates a new read-only instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceCaller(address common.Address, caller bind.ContractCaller) (*OwnableInterfaceCaller, error) {
	contract, err := bindOwnableInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceCaller{contract: contract}, nil
}

// NewOwnableInterfaceTransactor creates a new write-only instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableInterfaceTransactor, error) {
	contract, err := bindOwnableInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceTransactor{contract: contract}, nil
}

// NewOwnableInterfaceFilterer creates a new log filterer instance of OwnableInterface, bound to a specific deployed contract.
func NewOwnableInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableInterfaceFilterer, error) {
	contract, err := bindOwnableInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableInterfaceFilterer{contract: contract}, nil
}

// bindOwnableInterface binds a generic wrapper to an already deployed contract.
func bindOwnableInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OwnableInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableInterface *OwnableInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableInterface.Contract.OwnableInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableInterface *OwnableInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.Contract.OwnableInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableInterface *OwnableInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableInterface.Contract.OwnableInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableInterface *OwnableInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableInterface *OwnableInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableInterface *OwnableInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableInterface.Contract.contract.Transact(opts, method, params...)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnableInterface.Contract.AcceptOwnership(&_OwnableInterface.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnableInterface *OwnableInterfaceTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnableInterface.Contract.AcceptOwnership(&_OwnableInterface.TransactOpts)
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceTransactor) Owner(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "owner")
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceSession) Owner() (*types.Transaction, error) {
	return _OwnableInterface.Contract.Owner(&_OwnableInterface.TransactOpts)
}

// Owner is a paid mutator transaction binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() returns(address)
func (_OwnableInterface *OwnableInterfaceTransactorSession) Owner() (*types.Transaction, error) {
	return _OwnableInterface.Contract.Owner(&_OwnableInterface.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceTransactor) TransferOwnership(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.contract.Transact(opts, "transferOwnership", recipient)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceSession) TransferOwnership(recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.Contract.TransferOwnership(&_OwnableInterface.TransactOpts, recipient)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address recipient) returns()
func (_OwnableInterface *OwnableInterfaceTransactorSession) TransferOwnership(recipient common.Address) (*types.Transaction, error) {
	return _OwnableInterface.Contract.TransferOwnership(&_OwnableInterface.TransactOpts, recipient)
}

// OwnerIsCreatorMetaData contains all meta data concerning the OwnerIsCreator contract.
var OwnerIsCreatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610149565b6001600160a01b0381163314156100f85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610368806101586000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b1461007c575b600080fd5b61004e61008f565b005b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e61008a36600461031e565b610191565b60015473ffffffffffffffffffffffffffffffffffffffff163314610115576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6101996101a5565b6101a281610228565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610226576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161010c565b565b73ffffffffffffffffffffffffffffffffffffffff81163314156102a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161010c565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006020828403121561033057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461035457600080fd5b939250505056fea164736f6c6343000806000a",
}

// OwnerIsCreatorABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnerIsCreatorMetaData.ABI instead.
var OwnerIsCreatorABI = OwnerIsCreatorMetaData.ABI

// OwnerIsCreatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OwnerIsCreatorMetaData.Bin instead.
var OwnerIsCreatorBin = OwnerIsCreatorMetaData.Bin

// DeployOwnerIsCreator deploys a new Ethereum contract, binding an instance of OwnerIsCreator to it.
func DeployOwnerIsCreator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OwnerIsCreator, error) {
	parsed, err := OwnerIsCreatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OwnerIsCreatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OwnerIsCreator{OwnerIsCreatorCaller: OwnerIsCreatorCaller{contract: contract}, OwnerIsCreatorTransactor: OwnerIsCreatorTransactor{contract: contract}, OwnerIsCreatorFilterer: OwnerIsCreatorFilterer{contract: contract}}, nil
}

// OwnerIsCreator is an auto generated Go binding around an Ethereum contract.
type OwnerIsCreator struct {
	OwnerIsCreatorCaller     // Read-only binding to the contract
	OwnerIsCreatorTransactor // Write-only binding to the contract
	OwnerIsCreatorFilterer   // Log filterer for contract events
}

// OwnerIsCreatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnerIsCreatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnerIsCreatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnerIsCreatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerIsCreatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnerIsCreatorSession struct {
	Contract     *OwnerIsCreator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnerIsCreatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnerIsCreatorCallerSession struct {
	Contract *OwnerIsCreatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OwnerIsCreatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnerIsCreatorTransactorSession struct {
	Contract     *OwnerIsCreatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OwnerIsCreatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnerIsCreatorRaw struct {
	Contract *OwnerIsCreator // Generic contract binding to access the raw methods on
}

// OwnerIsCreatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnerIsCreatorCallerRaw struct {
	Contract *OwnerIsCreatorCaller // Generic read-only contract binding to access the raw methods on
}

// OwnerIsCreatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnerIsCreatorTransactorRaw struct {
	Contract *OwnerIsCreatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnerIsCreator creates a new instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreator(address common.Address, backend bind.ContractBackend) (*OwnerIsCreator, error) {
	contract, err := bindOwnerIsCreator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreator{OwnerIsCreatorCaller: OwnerIsCreatorCaller{contract: contract}, OwnerIsCreatorTransactor: OwnerIsCreatorTransactor{contract: contract}, OwnerIsCreatorFilterer: OwnerIsCreatorFilterer{contract: contract}}, nil
}

// NewOwnerIsCreatorCaller creates a new read-only instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorCaller(address common.Address, caller bind.ContractCaller) (*OwnerIsCreatorCaller, error) {
	contract, err := bindOwnerIsCreator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorCaller{contract: contract}, nil
}

// NewOwnerIsCreatorTransactor creates a new write-only instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnerIsCreatorTransactor, error) {
	contract, err := bindOwnerIsCreator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorTransactor{contract: contract}, nil
}

// NewOwnerIsCreatorFilterer creates a new log filterer instance of OwnerIsCreator, bound to a specific deployed contract.
func NewOwnerIsCreatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnerIsCreatorFilterer, error) {
	contract, err := bindOwnerIsCreator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorFilterer{contract: contract}, nil
}

// bindOwnerIsCreator binds a generic wrapper to an already deployed contract.
func bindOwnerIsCreator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OwnerIsCreatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerIsCreator.Contract.OwnerIsCreatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.OwnerIsCreatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerIsCreator *OwnerIsCreatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.OwnerIsCreatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerIsCreator *OwnerIsCreatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerIsCreator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerIsCreator *OwnerIsCreatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerIsCreator *OwnerIsCreatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OwnerIsCreator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorSession) Owner() (common.Address, error) {
	return _OwnerIsCreator.Contract.Owner(&_OwnerIsCreator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnerIsCreator *OwnerIsCreatorCallerSession) Owner() (common.Address, error) {
	return _OwnerIsCreator.Contract.Owner(&_OwnerIsCreator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerIsCreator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.AcceptOwnership(&_OwnerIsCreator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.AcceptOwnership(&_OwnerIsCreator.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.TransferOwnership(&_OwnerIsCreator.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_OwnerIsCreator *OwnerIsCreatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OwnerIsCreator.Contract.TransferOwnership(&_OwnerIsCreator.TransactOpts, to)
}

// OwnerIsCreatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferRequestedIterator struct {
	Event *OwnerIsCreatorOwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerIsCreatorOwnershipTransferRequested)
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
		it.Event = new(OwnerIsCreatorOwnershipTransferRequested)
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
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerIsCreatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerIsCreatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnerIsCreatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorOwnershipTransferRequestedIterator{contract: _OwnerIsCreator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OwnerIsCreatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerIsCreatorOwnershipTransferRequested)
				if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OwnerIsCreatorOwnershipTransferRequested, error) {
	event := new(OwnerIsCreatorOwnershipTransferRequested)
	if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnerIsCreatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferredIterator struct {
	Event *OwnerIsCreatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OwnerIsCreatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerIsCreatorOwnershipTransferred)
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
		it.Event = new(OwnerIsCreatorOwnershipTransferred)
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
func (it *OwnerIsCreatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerIsCreatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerIsCreatorOwnershipTransferred represents a OwnershipTransferred event raised by the OwnerIsCreator contract.
type OwnerIsCreatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnerIsCreatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnerIsCreatorOwnershipTransferredIterator{contract: _OwnerIsCreator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnerIsCreatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OwnerIsCreator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerIsCreatorOwnershipTransferred)
				if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OwnerIsCreator *OwnerIsCreatorFilterer) ParseOwnershipTransferred(log types.Log) (*OwnerIsCreatorOwnershipTransferred, error) {
	event := new(OwnerIsCreatorOwnershipTransferred)
	if err := _OwnerIsCreator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TypeAndVersionInterfaceMetaData contains all meta data concerning the TypeAndVersionInterface contract.
var TypeAndVersionInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// TypeAndVersionInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use TypeAndVersionInterfaceMetaData.ABI instead.
var TypeAndVersionInterfaceABI = TypeAndVersionInterfaceMetaData.ABI

// TypeAndVersionInterface is an auto generated Go binding around an Ethereum contract.
type TypeAndVersionInterface struct {
	TypeAndVersionInterfaceCaller     // Read-only binding to the contract
	TypeAndVersionInterfaceTransactor // Write-only binding to the contract
	TypeAndVersionInterfaceFilterer   // Log filterer for contract events
}

// TypeAndVersionInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TypeAndVersionInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypeAndVersionInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TypeAndVersionInterfaceSession struct {
	Contract     *TypeAndVersionInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TypeAndVersionInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TypeAndVersionInterfaceCallerSession struct {
	Contract *TypeAndVersionInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// TypeAndVersionInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TypeAndVersionInterfaceTransactorSession struct {
	Contract     *TypeAndVersionInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// TypeAndVersionInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type TypeAndVersionInterfaceRaw struct {
	Contract *TypeAndVersionInterface // Generic contract binding to access the raw methods on
}

// TypeAndVersionInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceCallerRaw struct {
	Contract *TypeAndVersionInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// TypeAndVersionInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TypeAndVersionInterfaceTransactorRaw struct {
	Contract *TypeAndVersionInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTypeAndVersionInterface creates a new instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterface(address common.Address, backend bind.ContractBackend) (*TypeAndVersionInterface, error) {
	contract, err := bindTypeAndVersionInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterface{TypeAndVersionInterfaceCaller: TypeAndVersionInterfaceCaller{contract: contract}, TypeAndVersionInterfaceTransactor: TypeAndVersionInterfaceTransactor{contract: contract}, TypeAndVersionInterfaceFilterer: TypeAndVersionInterfaceFilterer{contract: contract}}, nil
}

// NewTypeAndVersionInterfaceCaller creates a new read-only instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceCaller(address common.Address, caller bind.ContractCaller) (*TypeAndVersionInterfaceCaller, error) {
	contract, err := bindTypeAndVersionInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceCaller{contract: contract}, nil
}

// NewTypeAndVersionInterfaceTransactor creates a new write-only instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*TypeAndVersionInterfaceTransactor, error) {
	contract, err := bindTypeAndVersionInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceTransactor{contract: contract}, nil
}

// NewTypeAndVersionInterfaceFilterer creates a new log filterer instance of TypeAndVersionInterface, bound to a specific deployed contract.
func NewTypeAndVersionInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*TypeAndVersionInterfaceFilterer, error) {
	contract, err := bindTypeAndVersionInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TypeAndVersionInterfaceFilterer{contract: contract}, nil
}

// bindTypeAndVersionInterface binds a generic wrapper to an already deployed contract.
func bindTypeAndVersionInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TypeAndVersionInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersionInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TypeAndVersionInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TypeAndVersionInterface *TypeAndVersionInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TypeAndVersionInterface.Contract.contract.Transact(opts, method, params...)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TypeAndVersionInterface.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceSession) TypeAndVersion() (string, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersion(&_TypeAndVersionInterface.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() pure returns(string)
func (_TypeAndVersionInterface *TypeAndVersionInterfaceCallerSession) TypeAndVersion() (string, error) {
	return _TypeAndVersionInterface.Contract.TypeAndVersion(&_TypeAndVersionInterface.CallOpts)
}
