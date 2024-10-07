// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr3dynamicallydispatchedblsattestationverifierlib

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

// OCR3BLSAttestationVerifierLibMetaData contains all meta data concerning the OCR3BLSAttestationVerifierLib contract.
var OCR3BLSAttestationVerifierLibMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x602d6037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea164736f6c6343000813000a",
}

// OCR3BLSAttestationVerifierLibABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3BLSAttestationVerifierLibMetaData.ABI instead.
var OCR3BLSAttestationVerifierLibABI = OCR3BLSAttestationVerifierLibMetaData.ABI

// OCR3BLSAttestationVerifierLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR3BLSAttestationVerifierLibMetaData.Bin instead.
var OCR3BLSAttestationVerifierLibBin = OCR3BLSAttestationVerifierLibMetaData.Bin

// DeployOCR3BLSAttestationVerifierLib deploys a new Ethereum contract, binding an instance of OCR3BLSAttestationVerifierLib to it.
func DeployOCR3BLSAttestationVerifierLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR3BLSAttestationVerifierLib, error) {
	parsed, err := OCR3BLSAttestationVerifierLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR3BLSAttestationVerifierLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR3BLSAttestationVerifierLib{OCR3BLSAttestationVerifierLibCaller: OCR3BLSAttestationVerifierLibCaller{contract: contract}, OCR3BLSAttestationVerifierLibTransactor: OCR3BLSAttestationVerifierLibTransactor{contract: contract}, OCR3BLSAttestationVerifierLibFilterer: OCR3BLSAttestationVerifierLibFilterer{contract: contract}}, nil
}

// OCR3BLSAttestationVerifierLib is an auto generated Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierLib struct {
	OCR3BLSAttestationVerifierLibCaller     // Read-only binding to the contract
	OCR3BLSAttestationVerifierLibTransactor // Write-only binding to the contract
	OCR3BLSAttestationVerifierLibFilterer   // Log filterer for contract events
}

// OCR3BLSAttestationVerifierLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3BLSAttestationVerifierLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3BLSAttestationVerifierLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3BLSAttestationVerifierLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3BLSAttestationVerifierLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3BLSAttestationVerifierLibSession struct {
	Contract     *OCR3BLSAttestationVerifierLib // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// OCR3BLSAttestationVerifierLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3BLSAttestationVerifierLibCallerSession struct {
	Contract *OCR3BLSAttestationVerifierLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// OCR3BLSAttestationVerifierLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3BLSAttestationVerifierLibTransactorSession struct {
	Contract     *OCR3BLSAttestationVerifierLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// OCR3BLSAttestationVerifierLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierLibRaw struct {
	Contract *OCR3BLSAttestationVerifierLib // Generic contract binding to access the raw methods on
}

// OCR3BLSAttestationVerifierLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierLibCallerRaw struct {
	Contract *OCR3BLSAttestationVerifierLibCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3BLSAttestationVerifierLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierLibTransactorRaw struct {
	Contract *OCR3BLSAttestationVerifierLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3BLSAttestationVerifierLib creates a new instance of OCR3BLSAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3BLSAttestationVerifierLib(address common.Address, backend bind.ContractBackend) (*OCR3BLSAttestationVerifierLib, error) {
	contract, err := bindOCR3BLSAttestationVerifierLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3BLSAttestationVerifierLib{OCR3BLSAttestationVerifierLibCaller: OCR3BLSAttestationVerifierLibCaller{contract: contract}, OCR3BLSAttestationVerifierLibTransactor: OCR3BLSAttestationVerifierLibTransactor{contract: contract}, OCR3BLSAttestationVerifierLibFilterer: OCR3BLSAttestationVerifierLibFilterer{contract: contract}}, nil
}

// NewOCR3BLSAttestationVerifierLibCaller creates a new read-only instance of OCR3BLSAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3BLSAttestationVerifierLibCaller(address common.Address, caller bind.ContractCaller) (*OCR3BLSAttestationVerifierLibCaller, error) {
	contract, err := bindOCR3BLSAttestationVerifierLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3BLSAttestationVerifierLibCaller{contract: contract}, nil
}

// NewOCR3BLSAttestationVerifierLibTransactor creates a new write-only instance of OCR3BLSAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3BLSAttestationVerifierLibTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3BLSAttestationVerifierLibTransactor, error) {
	contract, err := bindOCR3BLSAttestationVerifierLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3BLSAttestationVerifierLibTransactor{contract: contract}, nil
}

// NewOCR3BLSAttestationVerifierLibFilterer creates a new log filterer instance of OCR3BLSAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3BLSAttestationVerifierLibFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3BLSAttestationVerifierLibFilterer, error) {
	contract, err := bindOCR3BLSAttestationVerifierLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3BLSAttestationVerifierLibFilterer{contract: contract}, nil
}

// bindOCR3BLSAttestationVerifierLib binds a generic wrapper to an already deployed contract.
func bindOCR3BLSAttestationVerifierLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3BLSAttestationVerifierLibMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3BLSAttestationVerifierLib *OCR3BLSAttestationVerifierLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3BLSAttestationVerifierLib.Contract.OCR3BLSAttestationVerifierLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3BLSAttestationVerifierLib *OCR3BLSAttestationVerifierLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3BLSAttestationVerifierLib.Contract.OCR3BLSAttestationVerifierLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3BLSAttestationVerifierLib *OCR3BLSAttestationVerifierLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3BLSAttestationVerifierLib.Contract.OCR3BLSAttestationVerifierLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3BLSAttestationVerifierLib *OCR3BLSAttestationVerifierLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3BLSAttestationVerifierLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3BLSAttestationVerifierLib *OCR3BLSAttestationVerifierLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3BLSAttestationVerifierLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3BLSAttestationVerifierLib *OCR3BLSAttestationVerifierLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3BLSAttestationVerifierLib.Contract.contract.Transact(opts, method, params...)
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibMetaData contains all meta data concerning the OCR3DynamicallyDispatchedBLSAttestationVerifierLib contract.
var OCR3DynamicallyDispatchedBLSAttestationVerifierLibMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidAttestation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationAttributionBitmask\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNumberOfKeys\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeysOfInvalidSize\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaximumNumberOfKeysExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"getSelectors\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"},{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x61146f61003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061004b5760003560e01c806322c95de3146100505780634b503f0b14610065578063e5eae2f6146100be575b600080fd5b61006361005e366004611197565b6100de565b005b604080517fe5eae2f60000000000000000000000000000000000000000000000000000000081527f22c95de300000000000000000000000000000000000000000000000000000000602082015281519081900390910190f35b8180156100ca57600080fd5b506100636100d936600461120f565b6100f4565b6100ec868686868686610106565b505050505050565b61010084848484610230565b50505050565b60258114610140576040517f1174ad8500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61016b6040518060800160405280600081526020016000815260200160008152602001600081525090565b600061017c888860ff16868661056a565b909250905061018c866001611298565b60ff1681146101c7576040517fddbf0b4400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080602160048701600037600051915086600052602160002090506101ee848284610728565b610224576040517fbd8ba84d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050505050505050565b61023b60a1826112e0565b15610272576040517fadd4994500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60ff831661028160a1836112f4565b146102b8576040517fa07f647e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208360ff1611156102f6576040517f1ede571b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103216040518060800160405280600081526020016000815260200160008152602001600081525090565b6103296110e5565b7f79537812dfe48a92fc860b8b010e8d6078b5c19e7037c4cf07f7bed69b54fffc8152610354611103565b6000805b8760ff168160ff16101561055f57868287610374602083611337565b9450610383928592919061134a565b61038c91611374565b8086526020808601919091528790839088906103a89083611337565b94506103b7928592919061134a565b6103c091611374565b602086810182905260408601919091528790839088906103e09083611337565b94506103ef928592919061134a565b6103f891611374565b604086018190526060850152868287610412602083611337565b9450610421928592919061134a565b61042a91611374565b6060860181905260808501526000878388610446602083611337565b9550610455928692919061134a565b61045e91611374565b9050600088888581811061047457610474611308565b9050013560f81c60f81b905060018461048d9190611337565b60a087208087527fff000000000000000000000000000000000000000000000000000000000000008316602088015260218720919550906104cf898286610728565b610505576040517f76d4e1e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b888d8660ff166020811061051b5761051b611308565b60040201600082015181600001556020820151816001015560408201518160020155606082015181600301559050505050505080610558906113b0565b9050610358565b505050505050505050565b6105956040518060800160405280600081526020016000815260200160008152602001600081525090565b6000806105a5600482868861134a565b6105ae916113cf565b60e01c90508015806105c357506001861b8110155b156105fa576040517ff4e04eaa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8160011660000361061f5760019190911c9061061881611417565b90506105fd565b60019250600061067689836020811061063a5761063a611308565b600402016040518060800160405290816000820154815260200160018201548152602001600282015481526020016003820154815250506108aa565b905061068182611417565b9150600183901c92505b82156107115760018316156106fa576106ec818a84602081106106b0576106b0611308565b60040201604051806080016040529081600082015481526020016001820154815260200160028201548152602001600382015481525050610919565b90506106f784611417565b93505b60019290921c9161070a82611417565b915061068b565b61071a81610bb4565b945050505094509492505050565b6000808080808061073887610c60565b9196509450925084610752576000955050505050506108a3565b61077d7fbfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8916610c60565b9196509250905084610797576000955050505050506108a3565b60006040518061018001604052808681526020018581526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec81526020017f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d81526020018481526020018381526020018b6000015181526020018b6020015181526020018b6040015181526020018b606001518152509050610885611121565b6020816101808460085afa61089957600080fd5b5196505050505050505b9392505050565b6108e36040518060c001604052806000815260200160008152602001600081526020016000815260200160008152602001600081525090565b8151815260208083015190820152604080830151908201526060918201519181019190915260006080820152600160a082015290565b6109526040518060c001604052806000815260200160008152602001600081526020016000815260200160008152602001600081525090565b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47608084015160a08501518281830983818208905083848485098503858485090884858286098684860908858684870987038784870908885160208a015188898684098a88840908898a8885098b038b888509088d51935060208e015192508a848c03830891508a838c03820890508a828b0899508a818a089850898b8a8c099a508b8b8c089a508b8c8283098d038d8c8d09089950508a888c038b0899508a878c038a0898508a81830997508a88890897508a8b8384098c038c8384090896508a888c038b0899508a878c038a0898508960808d01528860a08d01528a8860040999508a8760040998508a8b8a84098c8c84090897508a8b8b84098c038c8b84090896508a8b8a86098c8c86090891508a8b8b86098c038c8b860908905060408d0151995060608d015198508a8b868c098c888c090893508a8b878c098c038c878c0908925060408e0151995060608e015198508a8a8c03850895508a898c03840894508a8660020993508a8560020992508a83850995508a86870895508a8b8586098c038c8586090894508a888c03870895508a878c03860894508a8b888c098c8a8c09088b8c8a8d098d038d8a8d09088c826002099b508c816002099a5050508a8260020997508a8160020996508a888c03870897508a878c0386089650878c528660208d01528a888c03830897508a878c0382089650505088898684098a88840908935088898784098a038a878409089250505086868803830895508685880382089450505050508160408501528060608501525050505b92915050565b610bdf6040518060800160405280600081526020016000815260200160008152602001600081525090565b600080610bf484608001518560a00151610df4565b91509150600080610c058484610ef8565b91509150610c1d828288600001518960200151610fbc565b60208701528552610c3082828686610fbc565b8092508193505050610c4c828288604001518960600151610fbc565b606087015260408601525092949350505050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8216817f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478210610cbb57506000915081905080610ded565b60007f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47857f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4787880909089050610d62816002610d5b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd476001611337565b901c611082565b9150807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4783840914610d9f57600080600093509350935050610ded565b60ff85901c600183168114610de6577f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47610dd9848261144f565b610de391906112e0565b92505b6001945050505b9193909250565b600080807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808687097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47868709089050610e7881610e7360027f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4761144f565b611082565b90507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4781610ec6877f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4761144f565b0992507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478185099150505b9250929050565b6000807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4784840991507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4782830891507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808586097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478586090890509250929050565b6000807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808488097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478688090891507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808588097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4785880908905094509492505050565b600060405160208152602080820152602060408201528360608201528260808201527f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760a082015260208160c08360055afa6110dd57600080fd5b519392505050565b6040518060a001604052806005906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b803560ff8116811461115057600080fd5b919050565b60008083601f84011261116757600080fd5b50813567ffffffffffffffff81111561117f57600080fd5b602083019150836020828501011115610ef157600080fd5b60008060008060008060a087890312156111b057600080fd5b863595506111c06020880161113f565b94506111ce6040880161113f565b935060608701359250608087013567ffffffffffffffff8111156111f157600080fd5b6111fd89828a01611155565b979a9699509497509295939492505050565b6000806000806060858703121561122557600080fd5b843593506112356020860161113f565b9250604085013567ffffffffffffffff81111561125157600080fd5b61125d87828801611155565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff8181168382160190811115610bae57610bae611269565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826112ef576112ef6112b1565b500690565b600082611303576113036112b1565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80820180821115610bae57610bae611269565b6000808585111561135a57600080fd5b8386111561136757600080fd5b5050820193919092039150565b80356020831015610bae577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b600060ff821660ff81036113c6576113c6611269565b60010192915050565b7fffffffff00000000000000000000000000000000000000000000000000000000813581811691600485101561140f5780818660040360031b1b83161692505b505092915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361144857611448611269565b5060010190565b81810381811115610bae57610bae61126956fea164736f6c6343000813000a",
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3DynamicallyDispatchedBLSAttestationVerifierLibMetaData.ABI instead.
var OCR3DynamicallyDispatchedBLSAttestationVerifierLibABI = OCR3DynamicallyDispatchedBLSAttestationVerifierLibMetaData.ABI

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR3DynamicallyDispatchedBLSAttestationVerifierLibMetaData.Bin instead.
var OCR3DynamicallyDispatchedBLSAttestationVerifierLibBin = OCR3DynamicallyDispatchedBLSAttestationVerifierLibMetaData.Bin

// DeployOCR3DynamicallyDispatchedBLSAttestationVerifierLib deploys a new Ethereum contract, binding an instance of OCR3DynamicallyDispatchedBLSAttestationVerifierLib to it.
func DeployOCR3DynamicallyDispatchedBLSAttestationVerifierLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR3DynamicallyDispatchedBLSAttestationVerifierLib, error) {
	parsed, err := OCR3DynamicallyDispatchedBLSAttestationVerifierLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR3DynamicallyDispatchedBLSAttestationVerifierLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR3DynamicallyDispatchedBLSAttestationVerifierLib{OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller: OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller{contract: contract}, OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor: OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor{contract: contract}, OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer: OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer{contract: contract}}, nil
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLib is an auto generated Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLib struct {
	OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller     // Read-only binding to the contract
	OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor // Write-only binding to the contract
	OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer   // Log filterer for contract events
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibSession struct {
	Contract     *OCR3DynamicallyDispatchedBLSAttestationVerifierLib // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                                       // Call options to use throughout this session
	TransactOpts bind.TransactOpts                                   // Transaction auth options to use throughout this session
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibCallerSession struct {
	Contract *OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                                             // Call options to use throughout this session
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactorSession struct {
	Contract     *OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                                             // Transaction auth options to use throughout this session
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibRaw struct {
	Contract *OCR3DynamicallyDispatchedBLSAttestationVerifierLib // Generic contract binding to access the raw methods on
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibCallerRaw struct {
	Contract *OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactorRaw struct {
	Contract *OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3DynamicallyDispatchedBLSAttestationVerifierLib creates a new instance of OCR3DynamicallyDispatchedBLSAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedBLSAttestationVerifierLib(address common.Address, backend bind.ContractBackend) (*OCR3DynamicallyDispatchedBLSAttestationVerifierLib, error) {
	contract, err := bindOCR3DynamicallyDispatchedBLSAttestationVerifierLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedBLSAttestationVerifierLib{OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller: OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller{contract: contract}, OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor: OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor{contract: contract}, OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer: OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer{contract: contract}}, nil
}

// NewOCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller creates a new read-only instance of OCR3DynamicallyDispatchedBLSAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller(address common.Address, caller bind.ContractCaller) (*OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller, error) {
	contract, err := bindOCR3DynamicallyDispatchedBLSAttestationVerifierLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller{contract: contract}, nil
}

// NewOCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor creates a new write-only instance of OCR3DynamicallyDispatchedBLSAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor, error) {
	contract, err := bindOCR3DynamicallyDispatchedBLSAttestationVerifierLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor{contract: contract}, nil
}

// NewOCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer creates a new log filterer instance of OCR3DynamicallyDispatchedBLSAttestationVerifierLib, bound to a specific deployed contract.
func NewOCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer, error) {
	contract, err := bindOCR3DynamicallyDispatchedBLSAttestationVerifierLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3DynamicallyDispatchedBLSAttestationVerifierLibFilterer{contract: contract}, nil
}

// bindOCR3DynamicallyDispatchedBLSAttestationVerifierLib binds a generic wrapper to an already deployed contract.
func bindOCR3DynamicallyDispatchedBLSAttestationVerifierLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3DynamicallyDispatchedBLSAttestationVerifierLibMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.Contract.OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.Contract.OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.Contract.OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.Contract.contract.Transact(opts, method, params...)
}

// GetSelectors is a free data retrieval call binding the contract method 0x4b503f0b.
//
// Solidity: function getSelectors() pure returns(bytes4, bytes4)
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibCaller) GetSelectors(opts *bind.CallOpts) ([4]byte, [4]byte, error) {
	var out []interface{}
	err := _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.contract.Call(opts, &out, "getSelectors")

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
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibSession) GetSelectors() ([4]byte, [4]byte, error) {
	return _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.Contract.GetSelectors(&_OCR3DynamicallyDispatchedBLSAttestationVerifierLib.CallOpts)
}

// GetSelectors is a free data retrieval call binding the contract method 0x4b503f0b.
//
// Solidity: function getSelectors() pure returns(bytes4, bytes4)
func (_OCR3DynamicallyDispatchedBLSAttestationVerifierLib *OCR3DynamicallyDispatchedBLSAttestationVerifierLibCallerSession) GetSelectors() ([4]byte, [4]byte, error) {
	return _OCR3DynamicallyDispatchedBLSAttestationVerifierLib.Contract.GetSelectors(&_OCR3DynamicallyDispatchedBLSAttestationVerifierLib.CallOpts)
}
