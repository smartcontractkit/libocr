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
	ABI: "[{\"inputs\":[],\"name\":\"AttestationVerificationKeysOfInvalidSize\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateAttestationVerificationKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationAttributionBitmask\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationVerificationKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNumberOfAttestationVerificationKeys\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaximumNumberOfAttestationVerificationKeysExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"getSelectors\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"},{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6116c961003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061004b5760003560e01c806322c95de3146100505780634b503f0b146100655780635af994ad146100be575b600080fd5b61006361005e3660046113f1565b6100de565b005b604080517f5af994ad0000000000000000000000000000000000000000000000000000000081527f22c95de300000000000000000000000000000000000000000000000000000000602082015281519081900390910190f35b8180156100ca57600080fd5b506100636100d9366004611469565b6100f4565b6100ec868686868686610106565b505050505050565b61010084848484610230565b50505050565b60258114610140576040517f1174ad8500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61016b6040518060800160405280600081526020016000815260200160008152602001600081525090565b600061017c888860ff168686610773565b909250905061018c8660016114f2565b60ff1681146101c7576040517fddbf0b4400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080602160048701600037600051915086600052602160002090506101ee848284610931565b610224576040517fbd8ba84d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050505050505050565b61023b60a18261153a565b15610272576040517ff065284400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60ff831661028160a18361154e565b146102b8576040517f680d418700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208360ff1611156102f6576040517ffd6c8dce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102fe6112ee565b6103296040518060800160405280600081526020016000815260200160008152602001600081525090565b61033161133f565b7f79537812dfe48a92fc860b8b010e8d6078b5c19e7037c4cf07f7bed69b54fffc815261035c61135d565b6000805b8860ff168160ff1610156106005787828861037c602083611591565b945061038b92859291906115a4565b610394916115ce565b8086526020808601919091528890839089906103b09083611591565b94506103bf92859291906115a4565b6103c8916115ce565b602086810182905260408601919091528890839089906103e89083611591565b94506103f792859291906115a4565b610400916115ce565b60408601819052606085015287828861041a602083611591565b945061042992859291906115a4565b610432916115ce565b606086018190526080850152600088838961044e602083611591565b955061045d92869291906115a4565b610466916115ce565b9050600089898581811061047c5761047c611562565b9050013560f81c60f81b90506001846104959190611591565b60a087208087527fff000000000000000000000000000000000000000000000000000000000000008316602088015260218720919550906104d7898286610931565b61050d576040517f51e8c2fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b88518a60ff87166020811061052457610524611562565b6020020151600001818152505088602001518a8660ff166020811061054b5761054b611562565b6020020151602001818152505088604001518a8660ff166020811061057257610572611562565b6020020151604001818152505088606001518a8660ff166020811061059957610599611562565b60200201516060018181525050888e8660ff16602081106105bc576105bc611562565b600402016000820151816000015560208201518160010155604082015181600201556060820151816003015590505050505050806105f99061160a565b9050610360565b5060005b8860ff1681101561022457600061061c826001611591565b90505b8960ff168110156107625786816020811061063c5761063c611562565b60200201516000015187836020811061065757610657611562565b60200201515114801561069d575086816020811061067757610677611562565b60200201516020015187836020811061069257610692611562565b602002015160200151145b80156106dc57508681602081106106b6576106b6611562565b6020020151604001518783602081106106d1576106d1611562565b602002015160400151145b801561071b57508681602081106106f5576106f5611562565b60200201516060015187836020811061071057610710611562565b602002015160600151145b15610752576040517fb717a60a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61075b81611629565b905061061f565b5061076c81611629565b9050610604565b61079e6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000806107ae60048286886115a4565b6107b791611661565b60e01c90508015806107cc57506001861b8110155b15610803576040517ff4e04eaa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b816001166000036108285760019190911c9061082181611629565b9050610806565b60019250600061087f89836020811061084357610843611562565b60040201604051806080016040529081600082015481526020016001820154815260200160028201548152602001600382015481525050610ab3565b905061088a82611629565b9150600183901c92505b821561091a576001831615610903576108f5818a84602081106108b9576108b9611562565b60040201604051806080016040529081600082015481526020016001820154815260200160028201548152602001600382015481525050610b22565b905061090084611629565b93505b60019290921c9161091382611629565b9150610894565b61092381610dbd565b945050505094509492505050565b6000808080808061094187610e69565b919650945092508461095b57600095505050505050610aac565b6109867fbfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8916610e69565b91965092509050846109a057600095505050505050610aac565b60006040518061018001604052808681526020018581526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec81526020017f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d81526020018481526020018381526020018b6000015181526020018b6020015181526020018b6040015181526020018b606001518152509050610a8e61137b565b6020816101808460085afa610aa257600080fd5b5196505050505050505b9392505050565b610aec6040518060c001604052806000815260200160008152602001600081526020016000815260200160008152602001600081525090565b8151815260208083015190820152604080830151908201526060918201519181019190915260006080820152600160a082015290565b610b5b6040518060c001604052806000815260200160008152602001600081526020016000815260200160008152602001600081525090565b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47608084015160a08501518281830983818208905083848485098503858485090884858286098684860908858684870987038784870908885160208a015188898684098a88840908898a8885098b038b888509088d51935060208e015192508a848c03830891508a838c03820890508a828b0899508a818a089850898b8a8c099a508b8b8c089a508b8c8283098d038d8c8d09089950508a888c038b0899508a878c038a0898508a81830997508a88890897508a8b8384098c038c8384090896508a888c038b0899508a878c038a0898508960808d01528860a08d01528a8860040999508a8760040998508a8b8a84098c8c84090897508a8b8b84098c038c8b84090896508a8b8a86098c8c86090891508a8b8b86098c038c8b860908905060408d0151995060608d015198508a8b868c098c888c090893508a8b878c098c038c878c0908925060408e0151995060608e015198508a8a8c03850895508a898c03840894508a8660020993508a8560020992508a83850995508a86870895508a8b8586098c038c8586090894508a888c03870895508a878c03860894508a8b888c098c8a8c09088b8c8a8d098d038d8a8d09088c826002099b508c816002099a5050508a8260020997508a8160020996508a888c03870897508a878c0386089650878c528660208d01528a888c03830897508a878c0382089650505088898684098a88840908935088898784098a038a878409089250505086868803830895508685880382089450505050508160408501528060608501525050505b92915050565b610de86040518060800160405280600081526020016000815260200160008152602001600081525090565b600080610dfd84608001518560a00151610ffd565b91509150600080610e0e8484611101565b91509150610e268282886000015189602001516111c5565b60208701528552610e39828286866111c5565b8092508193505050610e558282886040015189606001516111c5565b606087015260408601525092949350505050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8216817f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478210610ec457506000915081905080610ff6565b60007f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47857f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4787880909089050610f6b816002610f647f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd476001611591565b901c61128b565b9150807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4783840914610fa857600080600093509350935050610ff6565b60ff85901c600183168114610fef577f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47610fe284826116a9565b610fec919061153a565b92505b6001945050505b9193909250565b600080807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808687097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478687090890506110818161107c60027f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd476116a9565b61128b565b90507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47816110cf877f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd476116a9565b0992507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478185099150505b9250929050565b6000807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4784840991507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4782830891507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808586097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478586090890509250929050565b6000807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808488097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478688090891507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808588097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4785880908905094509492505050565b600060405160208152602080820152602060408201528360608201528260808201527f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760a082015260208160c08360055afa6112e657600080fd5b519392505050565b6040518061040001604052806020905b6113296040518060800160405280600081526020016000815260200160008152602001600081525090565b8152602001906001900390816112fe5790505090565b6040518060a001604052806005906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b803560ff811681146113aa57600080fd5b919050565b60008083601f8401126113c157600080fd5b50813567ffffffffffffffff8111156113d957600080fd5b6020830191508360208285010111156110fa57600080fd5b60008060008060008060a0878903121561140a57600080fd5b8635955061141a60208801611399565b945061142860408801611399565b935060608701359250608087013567ffffffffffffffff81111561144b57600080fd5b61145789828a016113af565b979a9699509497509295939492505050565b6000806000806060858703121561147f57600080fd5b8435935061148f60208601611399565b9250604085013567ffffffffffffffff8111156114ab57600080fd5b6114b7878288016113af565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff8181168382160190811115610db757610db76114c3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826115495761154961150b565b500690565b60008261155d5761155d61150b565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80820180821115610db757610db76114c3565b600080858511156115b457600080fd5b838611156115c157600080fd5b5050820193919092039150565b80356020831015610db7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b600060ff821660ff8103611620576116206114c3565b60010192915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361165a5761165a6114c3565b5060010190565b7fffffffff0000000000000000000000000000000000000000000000000000000081358181169160048510156116a15780818660040360031b1b83161692505b505092915050565b81810381811115610db757610db76114c356fea164736f6c6343000813000a",
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
