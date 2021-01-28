// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testvalidator

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AggregatorValidatorInterfaceABI is the input ABI used to generate the binding from.
const AggregatorValidatorInterfaceABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"previousRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"previousAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"currentRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"currentAnswer\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AggregatorValidatorInterface is an auto generated Go binding around an Ethereum contract.
type AggregatorValidatorInterface struct {
	AggregatorValidatorInterfaceCaller     // Read-only binding to the contract
	AggregatorValidatorInterfaceTransactor // Write-only binding to the contract
	AggregatorValidatorInterfaceFilterer   // Log filterer for contract events
}

// AggregatorValidatorInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorValidatorInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorValidatorInterfaceSession struct {
	Contract     *AggregatorValidatorInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AggregatorValidatorInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorValidatorInterfaceCallerSession struct {
	Contract *AggregatorValidatorInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// AggregatorValidatorInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorValidatorInterfaceTransactorSession struct {
	Contract     *AggregatorValidatorInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// AggregatorValidatorInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceRaw struct {
	Contract *AggregatorValidatorInterface // Generic contract binding to access the raw methods on
}

// AggregatorValidatorInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceCallerRaw struct {
	Contract *AggregatorValidatorInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorValidatorInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceTransactorRaw struct {
	Contract *AggregatorValidatorInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorValidatorInterface creates a new instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterface(address common.Address, backend bind.ContractBackend) (*AggregatorValidatorInterface, error) {
	contract, err := bindAggregatorValidatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterface{AggregatorValidatorInterfaceCaller: AggregatorValidatorInterfaceCaller{contract: contract}, AggregatorValidatorInterfaceTransactor: AggregatorValidatorInterfaceTransactor{contract: contract}, AggregatorValidatorInterfaceFilterer: AggregatorValidatorInterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorValidatorInterfaceCaller creates a new read-only instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorValidatorInterfaceCaller, error) {
	contract, err := bindAggregatorValidatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceCaller{contract: contract}, nil
}

// NewAggregatorValidatorInterfaceTransactor creates a new write-only instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorValidatorInterfaceTransactor, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceTransactor{contract: contract}, nil
}

// NewAggregatorValidatorInterfaceFilterer creates a new log filterer instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorValidatorInterfaceFilterer, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceFilterer{contract: contract}, nil
}

// bindAggregatorValidatorInterface binds a generic wrapper to an already deployed contract.
func bindAggregatorValidatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorValidatorInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transact(opts, method, params...)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactor) Validate(opts *bind.TransactOpts, previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.contract.Transact(opts, "validate", previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// TestValidatorABI is the input ABI used to generate the binding from.
const TestValidatorABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// TestValidatorBin is the compiled bytecode used for deploying new contracts.
var TestValidatorBin = "0x6080604052348015600f57600080fd5b5060848061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063beed9b5114602d575b600080fd5b605960048036036080811015604157600080fd5b5080359060208101359060408101359060600135606d565b604080519115158252519081900360200190f35b600194935050505056fea164736f6c6343000705000a"

// DeployTestValidator deploys a new Ethereum contract, binding an instance of TestValidator to it.
func DeployTestValidator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestValidator, error) {
	parsed, err := abi.JSON(strings.NewReader(TestValidatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TestValidatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestValidator{TestValidatorCaller: TestValidatorCaller{contract: contract}, TestValidatorTransactor: TestValidatorTransactor{contract: contract}, TestValidatorFilterer: TestValidatorFilterer{contract: contract}}, nil
}

// TestValidator is an auto generated Go binding around an Ethereum contract.
type TestValidator struct {
	TestValidatorCaller     // Read-only binding to the contract
	TestValidatorTransactor // Write-only binding to the contract
	TestValidatorFilterer   // Log filterer for contract events
}

// TestValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestValidatorSession struct {
	Contract     *TestValidator    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestValidatorCallerSession struct {
	Contract *TestValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TestValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestValidatorTransactorSession struct {
	Contract     *TestValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TestValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestValidatorRaw struct {
	Contract *TestValidator // Generic contract binding to access the raw methods on
}

// TestValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestValidatorCallerRaw struct {
	Contract *TestValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// TestValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestValidatorTransactorRaw struct {
	Contract *TestValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestValidator creates a new instance of TestValidator, bound to a specific deployed contract.
func NewTestValidator(address common.Address, backend bind.ContractBackend) (*TestValidator, error) {
	contract, err := bindTestValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestValidator{TestValidatorCaller: TestValidatorCaller{contract: contract}, TestValidatorTransactor: TestValidatorTransactor{contract: contract}, TestValidatorFilterer: TestValidatorFilterer{contract: contract}}, nil
}

// NewTestValidatorCaller creates a new read-only instance of TestValidator, bound to a specific deployed contract.
func NewTestValidatorCaller(address common.Address, caller bind.ContractCaller) (*TestValidatorCaller, error) {
	contract, err := bindTestValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestValidatorCaller{contract: contract}, nil
}

// NewTestValidatorTransactor creates a new write-only instance of TestValidator, bound to a specific deployed contract.
func NewTestValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*TestValidatorTransactor, error) {
	contract, err := bindTestValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestValidatorTransactor{contract: contract}, nil
}

// NewTestValidatorFilterer creates a new log filterer instance of TestValidator, bound to a specific deployed contract.
func NewTestValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*TestValidatorFilterer, error) {
	contract, err := bindTestValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestValidatorFilterer{contract: contract}, nil
}

// bindTestValidator binds a generic wrapper to an already deployed contract.
func bindTestValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestValidatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestValidator *TestValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestValidator.Contract.TestValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestValidator *TestValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestValidator.Contract.TestValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestValidator *TestValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestValidator.Contract.TestValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestValidator *TestValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestValidator *TestValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestValidator *TestValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestValidator.Contract.contract.Transact(opts, method, params...)
}

// Validate is a free data retrieval call binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 , int256 , uint256 , int256 ) pure returns(bool)
func (_TestValidator *TestValidatorCaller) Validate(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int, arg2 *big.Int, arg3 *big.Int) (bool, error) {
	var out []interface{}
	err := _TestValidator.contract.Call(opts, &out, "validate", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Validate is a free data retrieval call binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 , int256 , uint256 , int256 ) pure returns(bool)
func (_TestValidator *TestValidatorSession) Validate(arg0 *big.Int, arg1 *big.Int, arg2 *big.Int, arg3 *big.Int) (bool, error) {
	return _TestValidator.Contract.Validate(&_TestValidator.CallOpts, arg0, arg1, arg2, arg3)
}

// Validate is a free data retrieval call binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 , int256 , uint256 , int256 ) pure returns(bool)
func (_TestValidator *TestValidatorCallerSession) Validate(arg0 *big.Int, arg1 *big.Int, arg2 *big.Int, arg3 *big.Int) (bool, error) {
	return _TestValidator.Contract.Validate(&_TestValidator.CallOpts, arg0, arg1, arg2, arg3)
}
