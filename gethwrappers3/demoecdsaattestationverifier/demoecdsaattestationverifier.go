// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package demoecdsaattestationverifier

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

// DemoECDSAAttestationVerifierMetaData contains all meta data concerning the DemoECDSAAttestationVerifier contract.
var DemoECDSAAttestationVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"DuplicateAttestationVerificationKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationAttributionBitmask\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationVerificationKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNumberOfAttestationVerificationKeys\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaximumNumberOfAttestationVerificationKeysExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"ocr3EcdsaSignerPublicKeys\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"seqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610a56806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636f289d411461003b57806395fd15c814610050575b600080fd5b61004e61004936600461063f565b610063565b005b61004e61005e36600461073e565b610103565b6040805160608101825263ffffffff871680825260ff87811660208085018290529188169390940183905280547fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000016909117640100000000909302929092177fffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffff16650100000000009091021790556100fc8483836101d5565b5050505050565b604080516060810182526020805463ffffffff8116835260ff64010000000082048116928401839052650100000000009091041692820183905290916101509188918891908888886101e7565b8051602080548184015160409094015163ffffffff9093167fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009091161764010000000060ff94851602177fffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffff166501000000000093909216929092021790555050505050565b6101e26000848484610241565b505050565b8251602080850191909120604080518084018b905267ffffffffffffffff8a1681830152606080820193909352815180820390930183526080019052805191012061023760008787848787610298565b5050505050505050565b600061024f8284018461082c565b90508360ff1681511461028e576040517f680d418700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6100fc8582610438565b60006102a5856001610927565b60ff1690506102b5816040610946565b6102c090600461095d565b82146102f8576040517f1174ad8500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006103076004828587610970565b6103109161099a565b60e01c9050600160ff88161b8110610354576040517ff4e04eaa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080604051878152601b602082015260408101600488018c5b86156103b65760018716156103aa576001850194506040828437604082019150600080526020600060808660015afa5060005181541495909501945b600196871c960161036e565b505050508382146103f3576040517fddbf0b4400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83811461042c576040517ff4e04eaa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050505050505050565b80516020811115610475576040517ffd6c8dce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b818110156105da576000838281518110610494576104946109e2565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610504576040517f51e8c2fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8281101561059d578173ffffffffffffffffffffffffffffffffffffffff16858281518110610538576105386109e2565b602002602001015173ffffffffffffffffffffffffffffffffffffffff160361058d576040517fb717a60a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61059681610a11565b9050610507565b508073ffffffffffffffffffffffffffffffffffffffff168583602081106105c7576105c76109e2565b0155506105d381610a11565b9050610478565b50505050565b803560ff811681146105f157600080fd5b919050565b60008083601f84011261060857600080fd5b50813567ffffffffffffffff81111561062057600080fd5b60208301915083602082850101111561063857600080fd5b9250929050565b60008060008060006080868803121561065757600080fd5b853563ffffffff8116811461066b57600080fd5b9450610679602087016105e0565b9350610687604087016105e0565b9250606086013567ffffffffffffffff8111156106a357600080fd5b6106af888289016105f6565b969995985093965092949392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610736576107366106c0565b604052919050565b60008060008060006080868803121561075657600080fd5b8535945060208087013567ffffffffffffffff808216821461077757600080fd5b9095506040880135908082111561078d57600080fd5b818901915089601f8301126107a157600080fd5b8135818111156107b3576107b36106c0565b6107e3847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016106ef565b8181528b858386010111156107f757600080fd5b81858501868301376000918101909401529194506060880135918083111561081e57600080fd5b50506106af888289016105f6565b6000602080838503121561083f57600080fd5b823567ffffffffffffffff8082111561085757600080fd5b818501915085601f83011261086b57600080fd5b81358181111561087d5761087d6106c0565b8060051b915061088e8483016106ef565b81815291830184019184810190888411156108a857600080fd5b938501935b838510156108ec578435925073ffffffffffffffffffffffffffffffffffffffff831683146108dc5760008081fd5b82825293850193908501906108ad565b98975050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff8181168382160190811115610940576109406108f8565b92915050565b8082028115828204841417610940576109406108f8565b80820180821115610940576109406108f8565b6000808585111561098057600080fd5b8386111561098d57600080fd5b5050820193919092039150565b7fffffffff0000000000000000000000000000000000000000000000000000000081358181169160048510156109da5780818660040360031b1b83161692505b505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610a4257610a426108f8565b506001019056fea164736f6c6343000813000a",
}

// DemoECDSAAttestationVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use DemoECDSAAttestationVerifierMetaData.ABI instead.
var DemoECDSAAttestationVerifierABI = DemoECDSAAttestationVerifierMetaData.ABI

// DemoECDSAAttestationVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DemoECDSAAttestationVerifierMetaData.Bin instead.
var DemoECDSAAttestationVerifierBin = DemoECDSAAttestationVerifierMetaData.Bin

// DeployDemoECDSAAttestationVerifier deploys a new Ethereum contract, binding an instance of DemoECDSAAttestationVerifier to it.
func DeployDemoECDSAAttestationVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DemoECDSAAttestationVerifier, error) {
	parsed, err := DemoECDSAAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DemoECDSAAttestationVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DemoECDSAAttestationVerifier{DemoECDSAAttestationVerifierCaller: DemoECDSAAttestationVerifierCaller{contract: contract}, DemoECDSAAttestationVerifierTransactor: DemoECDSAAttestationVerifierTransactor{contract: contract}, DemoECDSAAttestationVerifierFilterer: DemoECDSAAttestationVerifierFilterer{contract: contract}}, nil
}

// DemoECDSAAttestationVerifier is an auto generated Go binding around an Ethereum contract.
type DemoECDSAAttestationVerifier struct {
	DemoECDSAAttestationVerifierCaller     // Read-only binding to the contract
	DemoECDSAAttestationVerifierTransactor // Write-only binding to the contract
	DemoECDSAAttestationVerifierFilterer   // Log filterer for contract events
}

// DemoECDSAAttestationVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type DemoECDSAAttestationVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoECDSAAttestationVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DemoECDSAAttestationVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoECDSAAttestationVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DemoECDSAAttestationVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoECDSAAttestationVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DemoECDSAAttestationVerifierSession struct {
	Contract     *DemoECDSAAttestationVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DemoECDSAAttestationVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DemoECDSAAttestationVerifierCallerSession struct {
	Contract *DemoECDSAAttestationVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// DemoECDSAAttestationVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DemoECDSAAttestationVerifierTransactorSession struct {
	Contract     *DemoECDSAAttestationVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// DemoECDSAAttestationVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type DemoECDSAAttestationVerifierRaw struct {
	Contract *DemoECDSAAttestationVerifier // Generic contract binding to access the raw methods on
}

// DemoECDSAAttestationVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DemoECDSAAttestationVerifierCallerRaw struct {
	Contract *DemoECDSAAttestationVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// DemoECDSAAttestationVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DemoECDSAAttestationVerifierTransactorRaw struct {
	Contract *DemoECDSAAttestationVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDemoECDSAAttestationVerifier creates a new instance of DemoECDSAAttestationVerifier, bound to a specific deployed contract.
func NewDemoECDSAAttestationVerifier(address common.Address, backend bind.ContractBackend) (*DemoECDSAAttestationVerifier, error) {
	contract, err := bindDemoECDSAAttestationVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DemoECDSAAttestationVerifier{DemoECDSAAttestationVerifierCaller: DemoECDSAAttestationVerifierCaller{contract: contract}, DemoECDSAAttestationVerifierTransactor: DemoECDSAAttestationVerifierTransactor{contract: contract}, DemoECDSAAttestationVerifierFilterer: DemoECDSAAttestationVerifierFilterer{contract: contract}}, nil
}

// NewDemoECDSAAttestationVerifierCaller creates a new read-only instance of DemoECDSAAttestationVerifier, bound to a specific deployed contract.
func NewDemoECDSAAttestationVerifierCaller(address common.Address, caller bind.ContractCaller) (*DemoECDSAAttestationVerifierCaller, error) {
	contract, err := bindDemoECDSAAttestationVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DemoECDSAAttestationVerifierCaller{contract: contract}, nil
}

// NewDemoECDSAAttestationVerifierTransactor creates a new write-only instance of DemoECDSAAttestationVerifier, bound to a specific deployed contract.
func NewDemoECDSAAttestationVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*DemoECDSAAttestationVerifierTransactor, error) {
	contract, err := bindDemoECDSAAttestationVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DemoECDSAAttestationVerifierTransactor{contract: contract}, nil
}

// NewDemoECDSAAttestationVerifierFilterer creates a new log filterer instance of DemoECDSAAttestationVerifier, bound to a specific deployed contract.
func NewDemoECDSAAttestationVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*DemoECDSAAttestationVerifierFilterer, error) {
	contract, err := bindDemoECDSAAttestationVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DemoECDSAAttestationVerifierFilterer{contract: contract}, nil
}

// bindDemoECDSAAttestationVerifier binds a generic wrapper to an already deployed contract.
func bindDemoECDSAAttestationVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DemoECDSAAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DemoECDSAAttestationVerifier.Contract.DemoECDSAAttestationVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.Contract.DemoECDSAAttestationVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.Contract.DemoECDSAAttestationVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DemoECDSAAttestationVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.Contract.contract.Transact(opts, method, params...)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes ocr3EcdsaSignerPublicKeys) returns()
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierTransactor) SetConfig(opts *bind.TransactOpts, configVersion uint32, n uint8, f uint8, ocr3EcdsaSignerPublicKeys []byte) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.contract.Transact(opts, "setConfig", configVersion, n, f, ocr3EcdsaSignerPublicKeys)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes ocr3EcdsaSignerPublicKeys) returns()
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierSession) SetConfig(configVersion uint32, n uint8, f uint8, ocr3EcdsaSignerPublicKeys []byte) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.Contract.SetConfig(&_DemoECDSAAttestationVerifier.TransactOpts, configVersion, n, f, ocr3EcdsaSignerPublicKeys)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes ocr3EcdsaSignerPublicKeys) returns()
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierTransactorSession) SetConfig(configVersion uint32, n uint8, f uint8, ocr3EcdsaSignerPublicKeys []byte) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.Contract.SetConfig(&_DemoECDSAAttestationVerifier.TransactOpts, configVersion, n, f, ocr3EcdsaSignerPublicKeys)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierTransactor) Transmit(opts *bind.TransactOpts, configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.contract.Transact(opts, "transmit", configDigest, seqNr, report, attestation)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierSession) Transmit(configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.Contract.Transmit(&_DemoECDSAAttestationVerifier.TransactOpts, configDigest, seqNr, report, attestation)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoECDSAAttestationVerifier *DemoECDSAAttestationVerifierTransactorSession) Transmit(configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoECDSAAttestationVerifier.Contract.Transmit(&_DemoECDSAAttestationVerifier.TransactOpts, configDigest, seqNr, report, attestation)
}

// OCR3AttestationVerifierBaseMetaData contains all meta data concerning the OCR3AttestationVerifierBase contract.
var OCR3AttestationVerifierBaseMetaData = &bind.MetaData{
	ABI: "[]",
}

// OCR3AttestationVerifierBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3AttestationVerifierBaseMetaData.ABI instead.
var OCR3AttestationVerifierBaseABI = OCR3AttestationVerifierBaseMetaData.ABI

// OCR3AttestationVerifierBase is an auto generated Go binding around an Ethereum contract.
type OCR3AttestationVerifierBase struct {
	OCR3AttestationVerifierBaseCaller     // Read-only binding to the contract
	OCR3AttestationVerifierBaseTransactor // Write-only binding to the contract
	OCR3AttestationVerifierBaseFilterer   // Log filterer for contract events
}

// OCR3AttestationVerifierBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3AttestationVerifierBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3AttestationVerifierBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3AttestationVerifierBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3AttestationVerifierBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3AttestationVerifierBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3AttestationVerifierBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3AttestationVerifierBaseSession struct {
	Contract     *OCR3AttestationVerifierBase // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// OCR3AttestationVerifierBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3AttestationVerifierBaseCallerSession struct {
	Contract *OCR3AttestationVerifierBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// OCR3AttestationVerifierBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3AttestationVerifierBaseTransactorSession struct {
	Contract     *OCR3AttestationVerifierBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// OCR3AttestationVerifierBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3AttestationVerifierBaseRaw struct {
	Contract *OCR3AttestationVerifierBase // Generic contract binding to access the raw methods on
}

// OCR3AttestationVerifierBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3AttestationVerifierBaseCallerRaw struct {
	Contract *OCR3AttestationVerifierBaseCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3AttestationVerifierBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3AttestationVerifierBaseTransactorRaw struct {
	Contract *OCR3AttestationVerifierBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3AttestationVerifierBase creates a new instance of OCR3AttestationVerifierBase, bound to a specific deployed contract.
func NewOCR3AttestationVerifierBase(address common.Address, backend bind.ContractBackend) (*OCR3AttestationVerifierBase, error) {
	contract, err := bindOCR3AttestationVerifierBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3AttestationVerifierBase{OCR3AttestationVerifierBaseCaller: OCR3AttestationVerifierBaseCaller{contract: contract}, OCR3AttestationVerifierBaseTransactor: OCR3AttestationVerifierBaseTransactor{contract: contract}, OCR3AttestationVerifierBaseFilterer: OCR3AttestationVerifierBaseFilterer{contract: contract}}, nil
}

// NewOCR3AttestationVerifierBaseCaller creates a new read-only instance of OCR3AttestationVerifierBase, bound to a specific deployed contract.
func NewOCR3AttestationVerifierBaseCaller(address common.Address, caller bind.ContractCaller) (*OCR3AttestationVerifierBaseCaller, error) {
	contract, err := bindOCR3AttestationVerifierBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3AttestationVerifierBaseCaller{contract: contract}, nil
}

// NewOCR3AttestationVerifierBaseTransactor creates a new write-only instance of OCR3AttestationVerifierBase, bound to a specific deployed contract.
func NewOCR3AttestationVerifierBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3AttestationVerifierBaseTransactor, error) {
	contract, err := bindOCR3AttestationVerifierBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3AttestationVerifierBaseTransactor{contract: contract}, nil
}

// NewOCR3AttestationVerifierBaseFilterer creates a new log filterer instance of OCR3AttestationVerifierBase, bound to a specific deployed contract.
func NewOCR3AttestationVerifierBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3AttestationVerifierBaseFilterer, error) {
	contract, err := bindOCR3AttestationVerifierBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3AttestationVerifierBaseFilterer{contract: contract}, nil
}

// bindOCR3AttestationVerifierBase binds a generic wrapper to an already deployed contract.
func bindOCR3AttestationVerifierBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3AttestationVerifierBaseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3AttestationVerifierBase *OCR3AttestationVerifierBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3AttestationVerifierBase.Contract.OCR3AttestationVerifierBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3AttestationVerifierBase *OCR3AttestationVerifierBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3AttestationVerifierBase.Contract.OCR3AttestationVerifierBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3AttestationVerifierBase *OCR3AttestationVerifierBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3AttestationVerifierBase.Contract.OCR3AttestationVerifierBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3AttestationVerifierBase *OCR3AttestationVerifierBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3AttestationVerifierBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3AttestationVerifierBase *OCR3AttestationVerifierBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3AttestationVerifierBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3AttestationVerifierBase *OCR3AttestationVerifierBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3AttestationVerifierBase.Contract.contract.Transact(opts, method, params...)
}

// OCR3ECDSAAttestationVerifierMetaData contains all meta data concerning the OCR3ECDSAAttestationVerifier contract.
var OCR3ECDSAAttestationVerifierMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x6080604052348015600f57600080fd5b50601680601d6000396000f3fe6080604052600080fdfea164736f6c6343000813000a",
}

// OCR3ECDSAAttestationVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3ECDSAAttestationVerifierMetaData.ABI instead.
var OCR3ECDSAAttestationVerifierABI = OCR3ECDSAAttestationVerifierMetaData.ABI

// OCR3ECDSAAttestationVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR3ECDSAAttestationVerifierMetaData.Bin instead.
var OCR3ECDSAAttestationVerifierBin = OCR3ECDSAAttestationVerifierMetaData.Bin

// DeployOCR3ECDSAAttestationVerifier deploys a new Ethereum contract, binding an instance of OCR3ECDSAAttestationVerifier to it.
func DeployOCR3ECDSAAttestationVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR3ECDSAAttestationVerifier, error) {
	parsed, err := OCR3ECDSAAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR3ECDSAAttestationVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR3ECDSAAttestationVerifier{OCR3ECDSAAttestationVerifierCaller: OCR3ECDSAAttestationVerifierCaller{contract: contract}, OCR3ECDSAAttestationVerifierTransactor: OCR3ECDSAAttestationVerifierTransactor{contract: contract}, OCR3ECDSAAttestationVerifierFilterer: OCR3ECDSAAttestationVerifierFilterer{contract: contract}}, nil
}

// OCR3ECDSAAttestationVerifier is an auto generated Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifier struct {
	OCR3ECDSAAttestationVerifierCaller     // Read-only binding to the contract
	OCR3ECDSAAttestationVerifierTransactor // Write-only binding to the contract
	OCR3ECDSAAttestationVerifierFilterer   // Log filterer for contract events
}

// OCR3ECDSAAttestationVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3ECDSAAttestationVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3ECDSAAttestationVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3ECDSAAttestationVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3ECDSAAttestationVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3ECDSAAttestationVerifierSession struct {
	Contract     *OCR3ECDSAAttestationVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OCR3ECDSAAttestationVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3ECDSAAttestationVerifierCallerSession struct {
	Contract *OCR3ECDSAAttestationVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// OCR3ECDSAAttestationVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3ECDSAAttestationVerifierTransactorSession struct {
	Contract     *OCR3ECDSAAttestationVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// OCR3ECDSAAttestationVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierRaw struct {
	Contract *OCR3ECDSAAttestationVerifier // Generic contract binding to access the raw methods on
}

// OCR3ECDSAAttestationVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierCallerRaw struct {
	Contract *OCR3ECDSAAttestationVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3ECDSAAttestationVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierTransactorRaw struct {
	Contract *OCR3ECDSAAttestationVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3ECDSAAttestationVerifier creates a new instance of OCR3ECDSAAttestationVerifier, bound to a specific deployed contract.
func NewOCR3ECDSAAttestationVerifier(address common.Address, backend bind.ContractBackend) (*OCR3ECDSAAttestationVerifier, error) {
	contract, err := bindOCR3ECDSAAttestationVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3ECDSAAttestationVerifier{OCR3ECDSAAttestationVerifierCaller: OCR3ECDSAAttestationVerifierCaller{contract: contract}, OCR3ECDSAAttestationVerifierTransactor: OCR3ECDSAAttestationVerifierTransactor{contract: contract}, OCR3ECDSAAttestationVerifierFilterer: OCR3ECDSAAttestationVerifierFilterer{contract: contract}}, nil
}

// NewOCR3ECDSAAttestationVerifierCaller creates a new read-only instance of OCR3ECDSAAttestationVerifier, bound to a specific deployed contract.
func NewOCR3ECDSAAttestationVerifierCaller(address common.Address, caller bind.ContractCaller) (*OCR3ECDSAAttestationVerifierCaller, error) {
	contract, err := bindOCR3ECDSAAttestationVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3ECDSAAttestationVerifierCaller{contract: contract}, nil
}

// NewOCR3ECDSAAttestationVerifierTransactor creates a new write-only instance of OCR3ECDSAAttestationVerifier, bound to a specific deployed contract.
func NewOCR3ECDSAAttestationVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3ECDSAAttestationVerifierTransactor, error) {
	contract, err := bindOCR3ECDSAAttestationVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3ECDSAAttestationVerifierTransactor{contract: contract}, nil
}

// NewOCR3ECDSAAttestationVerifierFilterer creates a new log filterer instance of OCR3ECDSAAttestationVerifier, bound to a specific deployed contract.
func NewOCR3ECDSAAttestationVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3ECDSAAttestationVerifierFilterer, error) {
	contract, err := bindOCR3ECDSAAttestationVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3ECDSAAttestationVerifierFilterer{contract: contract}, nil
}

// bindOCR3ECDSAAttestationVerifier binds a generic wrapper to an already deployed contract.
func bindOCR3ECDSAAttestationVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3ECDSAAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3ECDSAAttestationVerifier *OCR3ECDSAAttestationVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3ECDSAAttestationVerifier.Contract.OCR3ECDSAAttestationVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3ECDSAAttestationVerifier *OCR3ECDSAAttestationVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3ECDSAAttestationVerifier.Contract.OCR3ECDSAAttestationVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3ECDSAAttestationVerifier *OCR3ECDSAAttestationVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3ECDSAAttestationVerifier.Contract.OCR3ECDSAAttestationVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3ECDSAAttestationVerifier *OCR3ECDSAAttestationVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3ECDSAAttestationVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3ECDSAAttestationVerifier *OCR3ECDSAAttestationVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3ECDSAAttestationVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3ECDSAAttestationVerifier *OCR3ECDSAAttestationVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3ECDSAAttestationVerifier.Contract.contract.Transact(opts, method, params...)
}

// OCR3ECDSAAttestationVerifierLibMetaData contains all meta data concerning the OCR3ECDSAAttestationVerifierLib contract.
var OCR3ECDSAAttestationVerifierLibMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x602d6037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea164736f6c6343000813000a",
}

// OCR3ECDSAAttestationVerifierLibABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3ECDSAAttestationVerifierLibMetaData.ABI instead.
var OCR3ECDSAAttestationVerifierLibABI = OCR3ECDSAAttestationVerifierLibMetaData.ABI

// OCR3ECDSAAttestationVerifierLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR3ECDSAAttestationVerifierLibMetaData.Bin instead.
var OCR3ECDSAAttestationVerifierLibBin = OCR3ECDSAAttestationVerifierLibMetaData.Bin

// DeployOCR3ECDSAAttestationVerifierLib deploys a new Ethereum contract, binding an instance of OCR3ECDSAAttestationVerifierLib to it.
func DeployOCR3ECDSAAttestationVerifierLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR3ECDSAAttestationVerifierLib, error) {
	parsed, err := OCR3ECDSAAttestationVerifierLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR3ECDSAAttestationVerifierLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR3ECDSAAttestationVerifierLib{OCR3ECDSAAttestationVerifierLibCaller: OCR3ECDSAAttestationVerifierLibCaller{contract: contract}, OCR3ECDSAAttestationVerifierLibTransactor: OCR3ECDSAAttestationVerifierLibTransactor{contract: contract}, OCR3ECDSAAttestationVerifierLibFilterer: OCR3ECDSAAttestationVerifierLibFilterer{contract: contract}}, nil
}

// OCR3ECDSAAttestationVerifierLib is an auto generated Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierLib struct {
	OCR3ECDSAAttestationVerifierLibCaller     // Read-only binding to the contract
	OCR3ECDSAAttestationVerifierLibTransactor // Write-only binding to the contract
	OCR3ECDSAAttestationVerifierLibFilterer   // Log filterer for contract events
}

// OCR3ECDSAAttestationVerifierLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3ECDSAAttestationVerifierLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3ECDSAAttestationVerifierLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3ECDSAAttestationVerifierLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3ECDSAAttestationVerifierLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3ECDSAAttestationVerifierLibSession struct {
	Contract     *OCR3ECDSAAttestationVerifierLib // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                    // Call options to use throughout this session
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// OCR3ECDSAAttestationVerifierLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3ECDSAAttestationVerifierLibCallerSession struct {
	Contract *OCR3ECDSAAttestationVerifierLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                          // Call options to use throughout this session
}

// OCR3ECDSAAttestationVerifierLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3ECDSAAttestationVerifierLibTransactorSession struct {
	Contract     *OCR3ECDSAAttestationVerifierLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                          // Transaction auth options to use throughout this session
}

// OCR3ECDSAAttestationVerifierLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierLibRaw struct {
	Contract *OCR3ECDSAAttestationVerifierLib // Generic contract binding to access the raw methods on
}

// OCR3ECDSAAttestationVerifierLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierLibCallerRaw struct {
	Contract *OCR3ECDSAAttestationVerifierLibCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3ECDSAAttestationVerifierLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3ECDSAAttestationVerifierLibTransactorRaw struct {
	Contract *OCR3ECDSAAttestationVerifierLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3ECDSAAttestationVerifierLib creates a new instance of OCR3ECDSAAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3ECDSAAttestationVerifierLib(address common.Address, backend bind.ContractBackend) (*OCR3ECDSAAttestationVerifierLib, error) {
	contract, err := bindOCR3ECDSAAttestationVerifierLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3ECDSAAttestationVerifierLib{OCR3ECDSAAttestationVerifierLibCaller: OCR3ECDSAAttestationVerifierLibCaller{contract: contract}, OCR3ECDSAAttestationVerifierLibTransactor: OCR3ECDSAAttestationVerifierLibTransactor{contract: contract}, OCR3ECDSAAttestationVerifierLibFilterer: OCR3ECDSAAttestationVerifierLibFilterer{contract: contract}}, nil
}

// NewOCR3ECDSAAttestationVerifierLibCaller creates a new read-only instance of OCR3ECDSAAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3ECDSAAttestationVerifierLibCaller(address common.Address, caller bind.ContractCaller) (*OCR3ECDSAAttestationVerifierLibCaller, error) {
	contract, err := bindOCR3ECDSAAttestationVerifierLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3ECDSAAttestationVerifierLibCaller{contract: contract}, nil
}

// NewOCR3ECDSAAttestationVerifierLibTransactor creates a new write-only instance of OCR3ECDSAAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3ECDSAAttestationVerifierLibTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3ECDSAAttestationVerifierLibTransactor, error) {
	contract, err := bindOCR3ECDSAAttestationVerifierLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3ECDSAAttestationVerifierLibTransactor{contract: contract}, nil
}

// NewOCR3ECDSAAttestationVerifierLibFilterer creates a new log filterer instance of OCR3ECDSAAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3ECDSAAttestationVerifierLibFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3ECDSAAttestationVerifierLibFilterer, error) {
	contract, err := bindOCR3ECDSAAttestationVerifierLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3ECDSAAttestationVerifierLibFilterer{contract: contract}, nil
}

// bindOCR3ECDSAAttestationVerifierLib binds a generic wrapper to an already deployed contract.
func bindOCR3ECDSAAttestationVerifierLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3ECDSAAttestationVerifierLibMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3ECDSAAttestationVerifierLib *OCR3ECDSAAttestationVerifierLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3ECDSAAttestationVerifierLib.Contract.OCR3ECDSAAttestationVerifierLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3ECDSAAttestationVerifierLib *OCR3ECDSAAttestationVerifierLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3ECDSAAttestationVerifierLib.Contract.OCR3ECDSAAttestationVerifierLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3ECDSAAttestationVerifierLib *OCR3ECDSAAttestationVerifierLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3ECDSAAttestationVerifierLib.Contract.OCR3ECDSAAttestationVerifierLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3ECDSAAttestationVerifierLib *OCR3ECDSAAttestationVerifierLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3ECDSAAttestationVerifierLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3ECDSAAttestationVerifierLib *OCR3ECDSAAttestationVerifierLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3ECDSAAttestationVerifierLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3ECDSAAttestationVerifierLib *OCR3ECDSAAttestationVerifierLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3ECDSAAttestationVerifierLib.Contract.contract.Transact(opts, method, params...)
}
