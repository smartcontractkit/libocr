// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr3dynamicallydispatchedecdsaattestationverifierlib

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

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibMetaData contains all meta data concerning the OCR3DynamicallyDispatchedECDSAAttestationVerifierLib contract.
var OCR3DynamicallyDispatchedECDSAAttestationVerifierLibMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"DuplicateKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationAttributionBitmask\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNumberOfKeys\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeysOfInvalidSize\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaximumNumberOfKeysExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"getSelectors\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"},{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6108ad61003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061004b5760003560e01c806349e333a6146100505780634b503f0b146100725780635cc513a6146100cb575b600080fd5b81801561005c57600080fd5b5061007061006b3660046105e0565b6100de565b005b604080517f49e333a60000000000000000000000000000000000000000000000000000000081527f5cc513a600000000000000000000000000000000000000000000000000000000602082015281519081900390910190f35b6100706100d936600461063a565b6100f0565b6100ea84848484610106565b50505050565b6100fe8686868686866103c2565b505050505050565b6101116014826106e1565b15610148576040517fadd4994500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60ff8316610157601483610724565b1461018e576040517fa07f647e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208360ff1611156101cc576040517f1ede571b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101d4610562565b6000805b8560ff168110156102dd5760008583866101f3826014610738565b9261020093929190610751565b6102099161077b565b60601c90506000819003610249576040517f76d4e1e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8084836020811061025c5761025c6107c3565b602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff168883602081106102bc576102bc6107c3565b01556102c9601484610738565b925050806102d6906107f2565b90506101d8565b5060005b8560ff168110156103b95760006102f9826001610738565b90505b8660ff168110156103a857838160208110610319576103196107c3565b602002015173ffffffffffffffffffffffffffffffffffffffff16848360208110610346576103466107c3565b602002015173ffffffffffffffffffffffffffffffffffffffff1603610398576040517f96bfdb4400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103a1816107f2565b90506102fc565b506103b2816107f2565b90506102e1565b50505050505050565b60006103cf85600161082a565b60ff1690506103df816040610843565b6103ea906004610738565b8214610422576040517f1174ad8500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006104316004828587610751565b61043a9161085a565b60e01c9050600160ff88161b811061047e576040517ff4e04eaa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080604051878152601b602082015260408101600488018c5b86156104e05760018716156104d4576001850194506040828437604082019150600080526020600060808660015afa5060005181541495909501945b600196871c9601610498565b5050505083821461051d576040517fddbf0b4400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b838114610556576040517ff4e04eaa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050505050505050565b6040518061040001604052806020906020820280368337509192915050565b803560ff8116811461059257600080fd5b919050565b60008083601f8401126105a957600080fd5b50813567ffffffffffffffff8111156105c157600080fd5b6020830191508360208285010111156105d957600080fd5b9250929050565b600080600080606085870312156105f657600080fd5b8435935061060660208601610581565b9250604085013567ffffffffffffffff81111561062257600080fd5b61062e87828801610597565b95989497509550505050565b60008060008060008060a0878903121561065357600080fd5b8635955061066360208801610581565b945061067160408801610581565b935060608701359250608087013567ffffffffffffffff81111561069457600080fd5b6106a089828a01610597565b979a9699509497509295939492505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826106f0576106f06106b2565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600082610733576107336106b2565b500490565b8082018082111561074b5761074b6106f5565b92915050565b6000808585111561076157600080fd5b8386111561076e57600080fd5b5050820193919092039150565b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000081358181169160148510156107bb5780818660140360031b1b83161692505b505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610823576108236106f5565b5060010190565b60ff818116838216019081111561074b5761074b6106f5565b808202811582820484141761074b5761074b6106f5565b7fffffffff0000000000000000000000000000000000000000000000000000000081358181169160048510156107bb5760049490940360031b84901b169092169291505056fea164736f6c6343000813000a",
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3DynamicallyDispatchedECDSAAttestationVerifierLibMetaData.ABI instead.
var OCR3DynamicallyDispatchedECDSAAttestationVerifierLibABI = OCR3DynamicallyDispatchedECDSAAttestationVerifierLibMetaData.ABI

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR3DynamicallyDispatchedECDSAAttestationVerifierLibMetaData.Bin instead.
var OCR3DynamicallyDispatchedECDSAAttestationVerifierLibBin = OCR3DynamicallyDispatchedECDSAAttestationVerifierLibMetaData.Bin

// DeployOCR3DynamicallyDispatchedECDSAAttestationVerifierLib deploys a new Ethereum contract, binding an instance of OCR3DynamicallyDispatchedECDSAAttestationVerifierLib to it.
func DeployOCR3DynamicallyDispatchedECDSAAttestationVerifierLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR3DynamicallyDispatchedECDSAAttestationVerifierLib, error) {
	parsed, err := OCR3DynamicallyDispatchedECDSAAttestationVerifierLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR3DynamicallyDispatchedECDSAAttestationVerifierLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR3DynamicallyDispatchedECDSAAttestationVerifierLib{OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller: OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller{contract: contract}, OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor: OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor{contract: contract}, OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer: OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer{contract: contract}}, nil
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLib is an auto generated Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLib struct {
	OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller     // Read-only binding to the contract
	OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor // Write-only binding to the contract
	OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer   // Log filterer for contract events
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibSession struct {
	Contract     *OCR3DynamicallyDispatchedECDSAAttestationVerifierLib // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                                         // Call options to use throughout this session
	TransactOpts bind.TransactOpts                                     // Transaction auth options to use throughout this session
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCallerSession struct {
	Contract *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                                               // Call options to use throughout this session
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactorSession struct {
	Contract     *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                                               // Transaction auth options to use throughout this session
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibRaw struct {
	Contract *OCR3DynamicallyDispatchedECDSAAttestationVerifierLib // Generic contract binding to access the raw methods on
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCallerRaw struct {
	Contract *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactorRaw struct {
	Contract *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3DynamicallyDispatchedECDSAAttestationVerifierLib creates a new instance of OCR3DynamicallyDispatchedECDSAAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedECDSAAttestationVerifierLib(address common.Address, backend bind.ContractBackend) (*OCR3DynamicallyDispatchedECDSAAttestationVerifierLib, error) {
	contract, err := bindOCR3DynamicallyDispatchedECDSAAttestationVerifierLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedECDSAAttestationVerifierLib{OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller: OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller{contract: contract}, OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor: OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor{contract: contract}, OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer: OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer{contract: contract}}, nil
}

// NewOCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller creates a new read-only instance of OCR3DynamicallyDispatchedECDSAAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller(address common.Address, caller bind.ContractCaller) (*OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller, error) {
	contract, err := bindOCR3DynamicallyDispatchedECDSAAttestationVerifierLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller{contract: contract}, nil
}

// NewOCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor creates a new write-only instance of OCR3DynamicallyDispatchedECDSAAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor, error) {
	contract, err := bindOCR3DynamicallyDispatchedECDSAAttestationVerifierLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor{contract: contract}, nil
}

// NewOCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer creates a new log filterer instance of OCR3DynamicallyDispatchedECDSAAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer, error) {
	contract, err := bindOCR3DynamicallyDispatchedECDSAAttestationVerifierLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedECDSAAttestationVerifierLibFilterer{contract: contract}, nil
}

// bindOCR3DynamicallyDispatchedECDSAAttestationVerifierLib binds a generic wrapper to an already deployed contract.
func bindOCR3DynamicallyDispatchedECDSAAttestationVerifierLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3DynamicallyDispatchedECDSAAttestationVerifierLibMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.Contract.OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.Contract.OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.Contract.OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.Contract.contract.Transact(opts, method, params...)
}

// GetSelectors is a free data retrieval call binding the contract method 0x4b503f0b.
//
// Solidity: function getSelectors() pure returns(bytes4, bytes4)
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCaller) GetSelectors(opts *bind.CallOpts) ([4]byte, [4]byte, error) {
	var out []interface{}
	err := _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.contract.Call(opts, &out, "getSelectors")

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
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibSession) GetSelectors() ([4]byte, [4]byte, error) {
	return _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.Contract.GetSelectors(&_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.CallOpts)
}

// GetSelectors is a free data retrieval call binding the contract method 0x4b503f0b.
//
// Solidity: function getSelectors() pure returns(bytes4, bytes4)
func (_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib *OCR3DynamicallyDispatchedECDSAAttestationVerifierLibCallerSession) GetSelectors() ([4]byte, [4]byte, error) {
	return _OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.Contract.GetSelectors(&_OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.CallOpts)
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
