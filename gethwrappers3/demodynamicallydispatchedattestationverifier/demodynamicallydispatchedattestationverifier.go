// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package demodynamicallydispatchedattestationverifier

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

// DemoDynamicallyDispatchedAttestationVerifierMetaData contains all meta data concerning the DemoDynamicallyDispatchedAttestationVerifier contract.
var DemoDynamicallyDispatchedAttestationVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierLibraryAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"keys\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"seqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801561001057600080fd5b5060405161083938038061083983398101604081905261002f916100b6565b6001600160a01b038116608081905260408051634b503f0b60e01b81528151849392634b503f0b92600480820193918290030181865afa158015610077573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061009b9190610103565b6001600160e01b031990811660c0521660a052506101369050565b6000602082840312156100c857600080fd5b81516001600160a01b03811681146100df57600080fd5b9392505050565b80516001600160e01b0319811681146100fe57600080fd5b919050565b6000806040838503121561011657600080fd5b61011f836100e6565b915061012d602084016100e6565b90509250929050565b60805160a05160c0516106d461016560003960006102e4015260006101de0152600061032e01526106d46000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636f289d411461003b57806395fd15c814610050575b600080fd5b61004e610049366004610427565b610063565b005b61004e61005e3660046104d7565b610104565b6040805160608101825263ffffffff871680825260ff8781166020840181905290871692909301829052608080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000016909117640100000000909302929092177fffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffff16650100000000009091021790556100fd8483836101d7565b5050505050565b6040805160608101825260805463ffffffff8116825260ff64010000000082048116602084018190526501000000000090920416928201839052909161015191889188919088888861029c565b805160808054602084015160409094015163ffffffff9093167fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009091161764010000000060ff94851602177fffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffff166501000000000093909216929092021790555050505050565b60006102967f000000000000000000000000000000000000000000000000000000000000000082868686604051602401610214949392919061062d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152610329565b50505050565b8251602080850191909120604080518084018b905267ffffffffffffffff8a16818301526060808201939093528151808203909301835260800190528051910120600061031e7f00000000000000000000000000000000000000000000000000000000000000008289898689896040516024016102149695949392919061065a565b505050505050505050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16836040516103719190610698565b600060405180830381855af49150503d80600081146103ac576040519150601f19603f3d011682016040523d82523d6000602084013e6103b1565b606091505b5091509150816103c357805181602001fd5b505050565b803560ff811681146103d957600080fd5b919050565b60008083601f8401126103f057600080fd5b50813567ffffffffffffffff81111561040857600080fd5b60208301915083602082850101111561042057600080fd5b9250929050565b60008060008060006080868803121561043f57600080fd5b853563ffffffff8116811461045357600080fd5b9450610461602087016103c8565b935061046f604087016103c8565b9250606086013567ffffffffffffffff81111561048b57600080fd5b610497888289016103de565b969995985093965092949392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806000806000608086880312156104ef57600080fd5b85359450602086013567ffffffffffffffff808216821461050f57600080fd5b9094506040870135908082111561052557600080fd5b818801915088601f83011261053957600080fd5b81358181111561054b5761054b6104a8565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610591576105916104a8565b816040528281528b60208487010111156105aa57600080fd5b8260208601602083013760006020848301015280975050505060608801359150808211156105d757600080fd5b50610497888289016103de565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b84815260ff841660208201526060604082015260006106506060830184866105e4565b9695505050505050565b86815260ff8616602082015260ff8516604082015283606082015260a06080820152600061068c60a0830184866105e4565b98975050505050505050565b6000825160005b818110156106b9576020818601810151858301520161069f565b50600092019182525091905056fea164736f6c6343000813000a",
}

// DemoDynamicallyDispatchedAttestationVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use DemoDynamicallyDispatchedAttestationVerifierMetaData.ABI instead.
var DemoDynamicallyDispatchedAttestationVerifierABI = DemoDynamicallyDispatchedAttestationVerifierMetaData.ABI

// DemoDynamicallyDispatchedAttestationVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DemoDynamicallyDispatchedAttestationVerifierMetaData.Bin instead.
var DemoDynamicallyDispatchedAttestationVerifierBin = DemoDynamicallyDispatchedAttestationVerifierMetaData.Bin

// DeployDemoDynamicallyDispatchedAttestationVerifier deploys a new Ethereum contract, binding an instance of DemoDynamicallyDispatchedAttestationVerifier to it.
func DeployDemoDynamicallyDispatchedAttestationVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, verifierLibraryAddress common.Address) (common.Address, *types.Transaction, *DemoDynamicallyDispatchedAttestationVerifier, error) {
	parsed, err := DemoDynamicallyDispatchedAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DemoDynamicallyDispatchedAttestationVerifierBin), backend, verifierLibraryAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DemoDynamicallyDispatchedAttestationVerifier{DemoDynamicallyDispatchedAttestationVerifierCaller: DemoDynamicallyDispatchedAttestationVerifierCaller{contract: contract}, DemoDynamicallyDispatchedAttestationVerifierTransactor: DemoDynamicallyDispatchedAttestationVerifierTransactor{contract: contract}, DemoDynamicallyDispatchedAttestationVerifierFilterer: DemoDynamicallyDispatchedAttestationVerifierFilterer{contract: contract}}, nil
}

// DemoDynamicallyDispatchedAttestationVerifier is an auto generated Go binding around an Ethereum contract.
type DemoDynamicallyDispatchedAttestationVerifier struct {
	DemoDynamicallyDispatchedAttestationVerifierCaller     // Read-only binding to the contract
	DemoDynamicallyDispatchedAttestationVerifierTransactor // Write-only binding to the contract
	DemoDynamicallyDispatchedAttestationVerifierFilterer   // Log filterer for contract events
}

// DemoDynamicallyDispatchedAttestationVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type DemoDynamicallyDispatchedAttestationVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoDynamicallyDispatchedAttestationVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DemoDynamicallyDispatchedAttestationVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoDynamicallyDispatchedAttestationVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DemoDynamicallyDispatchedAttestationVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoDynamicallyDispatchedAttestationVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DemoDynamicallyDispatchedAttestationVerifierSession struct {
	Contract     *DemoDynamicallyDispatchedAttestationVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// DemoDynamicallyDispatchedAttestationVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DemoDynamicallyDispatchedAttestationVerifierCallerSession struct {
	Contract *DemoDynamicallyDispatchedAttestationVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                                       // Call options to use throughout this session
}

// DemoDynamicallyDispatchedAttestationVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DemoDynamicallyDispatchedAttestationVerifierTransactorSession struct {
	Contract     *DemoDynamicallyDispatchedAttestationVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                                       // Transaction auth options to use throughout this session
}

// DemoDynamicallyDispatchedAttestationVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type DemoDynamicallyDispatchedAttestationVerifierRaw struct {
	Contract *DemoDynamicallyDispatchedAttestationVerifier // Generic contract binding to access the raw methods on
}

// DemoDynamicallyDispatchedAttestationVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DemoDynamicallyDispatchedAttestationVerifierCallerRaw struct {
	Contract *DemoDynamicallyDispatchedAttestationVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// DemoDynamicallyDispatchedAttestationVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DemoDynamicallyDispatchedAttestationVerifierTransactorRaw struct {
	Contract *DemoDynamicallyDispatchedAttestationVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDemoDynamicallyDispatchedAttestationVerifier creates a new instance of DemoDynamicallyDispatchedAttestationVerifier, bound to a specific deployed contract.
func NewDemoDynamicallyDispatchedAttestationVerifier(address common.Address, backend bind.ContractBackend) (*DemoDynamicallyDispatchedAttestationVerifier, error) {
	contract, err := bindDemoDynamicallyDispatchedAttestationVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DemoDynamicallyDispatchedAttestationVerifier{DemoDynamicallyDispatchedAttestationVerifierCaller: DemoDynamicallyDispatchedAttestationVerifierCaller{contract: contract}, DemoDynamicallyDispatchedAttestationVerifierTransactor: DemoDynamicallyDispatchedAttestationVerifierTransactor{contract: contract}, DemoDynamicallyDispatchedAttestationVerifierFilterer: DemoDynamicallyDispatchedAttestationVerifierFilterer{contract: contract}}, nil
}

// NewDemoDynamicallyDispatchedAttestationVerifierCaller creates a new read-only instance of DemoDynamicallyDispatchedAttestationVerifier, bound to a specific deployed contract.
func NewDemoDynamicallyDispatchedAttestationVerifierCaller(address common.Address, caller bind.ContractCaller) (*DemoDynamicallyDispatchedAttestationVerifierCaller, error) {
	contract, err := bindDemoDynamicallyDispatchedAttestationVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DemoDynamicallyDispatchedAttestationVerifierCaller{contract: contract}, nil
}

// NewDemoDynamicallyDispatchedAttestationVerifierTransactor creates a new write-only instance of DemoDynamicallyDispatchedAttestationVerifier, bound to a specific deployed contract.
func NewDemoDynamicallyDispatchedAttestationVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*DemoDynamicallyDispatchedAttestationVerifierTransactor, error) {
	contract, err := bindDemoDynamicallyDispatchedAttestationVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DemoDynamicallyDispatchedAttestationVerifierTransactor{contract: contract}, nil
}

// NewDemoDynamicallyDispatchedAttestationVerifierFilterer creates a new log filterer instance of DemoDynamicallyDispatchedAttestationVerifier, bound to a specific deployed contract.
func NewDemoDynamicallyDispatchedAttestationVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*DemoDynamicallyDispatchedAttestationVerifierFilterer, error) {
	contract, err := bindDemoDynamicallyDispatchedAttestationVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DemoDynamicallyDispatchedAttestationVerifierFilterer{contract: contract}, nil
}

// bindDemoDynamicallyDispatchedAttestationVerifier binds a generic wrapper to an already deployed contract.
func bindDemoDynamicallyDispatchedAttestationVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DemoDynamicallyDispatchedAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.DemoDynamicallyDispatchedAttestationVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.DemoDynamicallyDispatchedAttestationVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.DemoDynamicallyDispatchedAttestationVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.contract.Transact(opts, method, params...)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes keys) returns()
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierTransactor) SetConfig(opts *bind.TransactOpts, configVersion uint32, n uint8, f uint8, keys []byte) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.contract.Transact(opts, "setConfig", configVersion, n, f, keys)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes keys) returns()
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierSession) SetConfig(configVersion uint32, n uint8, f uint8, keys []byte) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.SetConfig(&_DemoDynamicallyDispatchedAttestationVerifier.TransactOpts, configVersion, n, f, keys)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes keys) returns()
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierTransactorSession) SetConfig(configVersion uint32, n uint8, f uint8, keys []byte) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.SetConfig(&_DemoDynamicallyDispatchedAttestationVerifier.TransactOpts, configVersion, n, f, keys)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierTransactor) Transmit(opts *bind.TransactOpts, configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.contract.Transact(opts, "transmit", configDigest, seqNr, report, attestation)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierSession) Transmit(configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.Transmit(&_DemoDynamicallyDispatchedAttestationVerifier.TransactOpts, configDigest, seqNr, report, attestation)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoDynamicallyDispatchedAttestationVerifier *DemoDynamicallyDispatchedAttestationVerifierTransactorSession) Transmit(configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoDynamicallyDispatchedAttestationVerifier.Contract.Transmit(&_DemoDynamicallyDispatchedAttestationVerifier.TransactOpts, configDigest, seqNr, report, attestation)
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

// OCR3DynamicallyDispatchedAttestationVerifierMetaData contains all meta data concerning the OCR3DynamicallyDispatchedAttestationVerifier contract.
var OCR3DynamicallyDispatchedAttestationVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifierLibraryAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]",
	Bin: "0x60e060405234801561001057600080fd5b5060405161016d38038061016d83398101604081905261002f916100b3565b6001600160a01b038116608081905260408051634b503f0b60e01b81528151634b503f0b926004808401939192918290030181865afa158015610076573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061009a9190610100565b6001600160e01b031990811660c0521660a05250610133565b6000602082840312156100c557600080fd5b81516001600160a01b03811681146100dc57600080fd5b9392505050565b80516001600160e01b0319811681146100fb57600080fd5b919050565b6000806040838503121561011357600080fd5b61011c836100e3565b915061012a602084016100e3565b90509250929050565b60805160a05160c051601661015760003960005050600050506000505060166000f3fe6080604052600080fdfea164736f6c6343000813000a",
}

// OCR3DynamicallyDispatchedAttestationVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3DynamicallyDispatchedAttestationVerifierMetaData.ABI instead.
var OCR3DynamicallyDispatchedAttestationVerifierABI = OCR3DynamicallyDispatchedAttestationVerifierMetaData.ABI

// OCR3DynamicallyDispatchedAttestationVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR3DynamicallyDispatchedAttestationVerifierMetaData.Bin instead.
var OCR3DynamicallyDispatchedAttestationVerifierBin = OCR3DynamicallyDispatchedAttestationVerifierMetaData.Bin

// DeployOCR3DynamicallyDispatchedAttestationVerifier deploys a new Ethereum contract, binding an instance of OCR3DynamicallyDispatchedAttestationVerifier to it.
func DeployOCR3DynamicallyDispatchedAttestationVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, verifierLibraryAddress common.Address) (common.Address, *types.Transaction, *OCR3DynamicallyDispatchedAttestationVerifier, error) {
	parsed, err := OCR3DynamicallyDispatchedAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR3DynamicallyDispatchedAttestationVerifierBin), backend, verifierLibraryAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR3DynamicallyDispatchedAttestationVerifier{OCR3DynamicallyDispatchedAttestationVerifierCaller: OCR3DynamicallyDispatchedAttestationVerifierCaller{contract: contract}, OCR3DynamicallyDispatchedAttestationVerifierTransactor: OCR3DynamicallyDispatchedAttestationVerifierTransactor{contract: contract}, OCR3DynamicallyDispatchedAttestationVerifierFilterer: OCR3DynamicallyDispatchedAttestationVerifierFilterer{contract: contract}}, nil
}

// OCR3DynamicallyDispatchedAttestationVerifier is an auto generated Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifier struct {
	OCR3DynamicallyDispatchedAttestationVerifierCaller     // Read-only binding to the contract
	OCR3DynamicallyDispatchedAttestationVerifierTransactor // Write-only binding to the contract
	OCR3DynamicallyDispatchedAttestationVerifierFilterer   // Log filterer for contract events
}

// OCR3DynamicallyDispatchedAttestationVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedAttestationVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedAttestationVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3DynamicallyDispatchedAttestationVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedAttestationVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3DynamicallyDispatchedAttestationVerifierSession struct {
	Contract     *OCR3DynamicallyDispatchedAttestationVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// OCR3DynamicallyDispatchedAttestationVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3DynamicallyDispatchedAttestationVerifierCallerSession struct {
	Contract *OCR3DynamicallyDispatchedAttestationVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                                       // Call options to use throughout this session
}

// OCR3DynamicallyDispatchedAttestationVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3DynamicallyDispatchedAttestationVerifierTransactorSession struct {
	Contract     *OCR3DynamicallyDispatchedAttestationVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                                       // Transaction auth options to use throughout this session
}

// OCR3DynamicallyDispatchedAttestationVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierRaw struct {
	Contract *OCR3DynamicallyDispatchedAttestationVerifier // Generic contract binding to access the raw methods on
}

// OCR3DynamicallyDispatchedAttestationVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierCallerRaw struct {
	Contract *OCR3DynamicallyDispatchedAttestationVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3DynamicallyDispatchedAttestationVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierTransactorRaw struct {
	Contract *OCR3DynamicallyDispatchedAttestationVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3DynamicallyDispatchedAttestationVerifier creates a new instance of OCR3DynamicallyDispatchedAttestationVerifier, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedAttestationVerifier(address common.Address, backend bind.ContractBackend) (*OCR3DynamicallyDispatchedAttestationVerifier, error) {
	contract, err := bindOCR3DynamicallyDispatchedAttestationVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedAttestationVerifier{OCR3DynamicallyDispatchedAttestationVerifierCaller: OCR3DynamicallyDispatchedAttestationVerifierCaller{contract: contract}, OCR3DynamicallyDispatchedAttestationVerifierTransactor: OCR3DynamicallyDispatchedAttestationVerifierTransactor{contract: contract}, OCR3DynamicallyDispatchedAttestationVerifierFilterer: OCR3DynamicallyDispatchedAttestationVerifierFilterer{contract: contract}}, nil
}

// NewOCR3DynamicallyDispatchedAttestationVerifierCaller creates a new read-only instance of OCR3DynamicallyDispatchedAttestationVerifier, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedAttestationVerifierCaller(address common.Address, caller bind.ContractCaller) (*OCR3DynamicallyDispatchedAttestationVerifierCaller, error) {
	contract, err := bindOCR3DynamicallyDispatchedAttestationVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedAttestationVerifierCaller{contract: contract}, nil
}

// NewOCR3DynamicallyDispatchedAttestationVerifierTransactor creates a new write-only instance of OCR3DynamicallyDispatchedAttestationVerifier, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedAttestationVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3DynamicallyDispatchedAttestationVerifierTransactor, error) {
	contract, err := bindOCR3DynamicallyDispatchedAttestationVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedAttestationVerifierTransactor{contract: contract}, nil
}

// NewOCR3DynamicallyDispatchedAttestationVerifierFilterer creates a new log filterer instance of OCR3DynamicallyDispatchedAttestationVerifier, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedAttestationVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3DynamicallyDispatchedAttestationVerifierFilterer, error) {
	contract, err := bindOCR3DynamicallyDispatchedAttestationVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedAttestationVerifierFilterer{contract: contract}, nil
}

// bindOCR3DynamicallyDispatchedAttestationVerifier binds a generic wrapper to an already deployed contract.
func bindOCR3DynamicallyDispatchedAttestationVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3DynamicallyDispatchedAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3DynamicallyDispatchedAttestationVerifier *OCR3DynamicallyDispatchedAttestationVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3DynamicallyDispatchedAttestationVerifier.Contract.OCR3DynamicallyDispatchedAttestationVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3DynamicallyDispatchedAttestationVerifier *OCR3DynamicallyDispatchedAttestationVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifier.Contract.OCR3DynamicallyDispatchedAttestationVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3DynamicallyDispatchedAttestationVerifier *OCR3DynamicallyDispatchedAttestationVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifier.Contract.OCR3DynamicallyDispatchedAttestationVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3DynamicallyDispatchedAttestationVerifier *OCR3DynamicallyDispatchedAttestationVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3DynamicallyDispatchedAttestationVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3DynamicallyDispatchedAttestationVerifier *OCR3DynamicallyDispatchedAttestationVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3DynamicallyDispatchedAttestationVerifier *OCR3DynamicallyDispatchedAttestationVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifier.Contract.contract.Transact(opts, method, params...)
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceMetaData contains all meta data concerning the OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface contract.
var OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getSelectors\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"},{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceMetaData.ABI instead.
var OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceABI = OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceMetaData.ABI

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface is an auto generated Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface struct {
	OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller     // Read-only binding to the contract
	OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor // Write-only binding to the contract
	OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer   // Log filterer for contract events
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceSession struct {
	Contract     *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                                                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts                                              // Transaction auth options to use throughout this session
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCallerSession struct {
	Contract *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                                                        // Call options to use throughout this session
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactorSession struct {
	Contract     *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                                                        // Transaction auth options to use throughout this session
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceRaw struct {
	Contract *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface // Generic contract binding to access the raw methods on
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCallerRaw struct {
	Contract *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactorRaw struct {
	Contract *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3DynamicallyDispatchedAttestationVerifierSelectorInterface creates a new instance of OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedAttestationVerifierSelectorInterface(address common.Address, backend bind.ContractBackend) (*OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface, error) {
	contract, err := bindOCR3DynamicallyDispatchedAttestationVerifierSelectorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface{OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller: OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller{contract: contract}, OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor: OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor{contract: contract}, OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer: OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer{contract: contract}}, nil
}

// NewOCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller creates a new read-only instance of OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller, error) {
	contract, err := bindOCR3DynamicallyDispatchedAttestationVerifierSelectorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller{contract: contract}, nil
}

// NewOCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor creates a new write-only instance of OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor, error) {
	contract, err := bindOCR3DynamicallyDispatchedAttestationVerifierSelectorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor{contract: contract}, nil
}

// NewOCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer creates a new log filterer instance of OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer, error) {
	contract, err := bindOCR3DynamicallyDispatchedAttestationVerifierSelectorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceFilterer{contract: contract}, nil
}

// bindOCR3DynamicallyDispatchedAttestationVerifierSelectorInterface binds a generic wrapper to an already deployed contract.
func bindOCR3DynamicallyDispatchedAttestationVerifierSelectorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.Contract.OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.Contract.OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.Contract.OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.Contract.contract.Transact(opts, method, params...)
}

// GetSelectors is a free data retrieval call binding the contract method 0x4b503f0b.
//
// Solidity: function getSelectors() pure returns(bytes4, bytes4)
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCaller) GetSelectors(opts *bind.CallOpts) ([4]byte, [4]byte, error) {
	var out []interface{}
	err := _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.contract.Call(opts, &out, "getSelectors")

	if err != nil {
		return *new([4]byte), *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)
	out1 := *abi.ConvertType(out[1], new([4]byte)).(*[4]byte)

	return out0, out1, err

}

// GetSelectors is a free data retrieval call binding the contract method 0x4b503f0b.
//
// Solidity: function getSelectors() pure returns(bytes4, bytes4)
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceSession) GetSelectors() ([4]byte, [4]byte, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.Contract.GetSelectors(&_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.CallOpts)
}

// GetSelectors is a free data retrieval call binding the contract method 0x4b503f0b.
//
// Solidity: function getSelectors() pure returns(bytes4, bytes4)
func (_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface *OCR3DynamicallyDispatchedAttestationVerifierSelectorInterfaceCallerSession) GetSelectors() ([4]byte, [4]byte, error) {
	return _OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.Contract.GetSelectors(&_OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface.CallOpts)
}
