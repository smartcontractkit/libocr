// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package demoblsattestationverifier

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

// DemoBLSAttestationVerifierMetaData contains all meta data concerning the DemoBLSAttestationVerifier contract.
var DemoBLSAttestationVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"DuplicateKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationAttributionBitmask\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAttestationNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNumberOfKeys\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeysOfInvalidSize\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaximumNumberOfKeysExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"configVersion\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"ocr3BlsSignerPublicKeys\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"seqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"attestation\",\"type\":\"bytes\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611906806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636f289d411461003b57806395fd15c814610050575b600080fd5b61004e610049366004611543565b610063565b005b61004e61005e3660046115f3565b610104565b6040805160608101825263ffffffff871680825260ff8781166020840181905290871692909301829052608080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000016909117640100000000909302929092177fffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffff16650100000000009091021790556100fd8483836101d7565b5050505050565b6040805160608101825260805463ffffffff8116825260ff6401000000008204811660208401819052650100000000009092041692820183905290916101519188918891908888886101e9565b805160808054602084015160409094015163ffffffff9093167fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009091161764010000000060ff94851602177fffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffff166501000000000093909216929092021790555050505050565b6101e46000848484610211565b505050565b60006101f6888886610760565b9050610207600087878487876107a9565b5050505050505050565b61021c60a18261172f565b15610253576040517fadd4994500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60ff831661026260a183611772565b14610299576040517fa07f647e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208360ff1611156102d7576040517f1ede571b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102df611440565b61030a6040518060800160405280600081526020016000815260200160008152602001600081525090565b610312611491565b7f79537812dfe48a92fc860b8b010e8d6078b5c19e7037c4cf07f7bed69b54fffc815261033d6114af565b6000805b8860ff168160ff1610156105e15787828861035d6020836117b5565b945061036c92859291906117c8565b610375916117f2565b80865260208086019190915288908390899061039190836117b5565b94506103a092859291906117c8565b6103a9916117f2565b602086810182905260408601919091528890839089906103c990836117b5565b94506103d892859291906117c8565b6103e1916117f2565b6040860181905260608501528782886103fb6020836117b5565b945061040a92859291906117c8565b610413916117f2565b606086018190526080850152600088838961042f6020836117b5565b955061043e92869291906117c8565b610447916117f2565b9050600089898581811061045d5761045d611786565b9050013560f81c60f81b905060018461047691906117b5565b60a087208087527fff000000000000000000000000000000000000000000000000000000000000008316602088015260218720919550906104b88982866108c7565b6104ee576040517f76d4e1e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b88518a60ff87166020811061050557610505611786565b6020020151600001818152505088602001518a8660ff166020811061052c5761052c611786565b6020020151602001818152505088604001518a8660ff166020811061055357610553611786565b6020020151604001818152505088606001518a8660ff166020811061057a5761057a611786565b60200201516060018181525050888e8660ff166020811061059d5761059d611786565b600402016000820151816000015560208201518160010155604082015181600201556060820151816003015590505050505050806105da9061182e565b9050610341565b5060005b8860ff168110156107545760006105fd8260016117b5565b90505b8960ff168110156107435786816020811061061d5761061d611786565b60200201516000015187836020811061063857610638611786565b60200201515114801561067e575086816020811061065857610658611786565b60200201516020015187836020811061067357610673611786565b602002015160200151145b80156106bd575086816020811061069757610697611786565b6020020151604001518783602081106106b2576106b2611786565b602002015160400151145b80156106fc57508681602081106106d6576106d6611786565b6020020151606001518783602081106106f1576106f1611786565b602002015160600151145b15610733576040517f96bfdb4400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61073c8161184d565b9050610600565b5061074d8161184d565b90506105e5565b50505050505050505050565b80516020808301919091206040805180840187905267ffffffffffffffff86168183015260608082019390935281518082039093018352608001905280519101205b9392505050565b602581146107e3576040517f1174ad8500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61080e6040518060800160405280600081526020016000815260200160008152602001600081525090565b600061081f888860ff168686610a47565b909250905061082f866001611885565b60ff16811461086a576040517fddbf0b4400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080602160048701600037600051915086600052602160002090506108918482846108c7565b610754576040517fbd8ba84d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080806108d787610c05565b91965094509250846108f1576000955050505050506107a2565b61091c7fbfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8916610c05565b9196509250905084610936576000955050505050506107a2565b60006040518061018001604052808681526020018581526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec81526020017f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d81526020018481526020018381526020018b6000015181526020018b6020015181526020018b6040015181526020018b606001518152509050610a246114cd565b6020816101808460085afa610a3857600080fd5b519a9950505050505050505050565b610a726040518060800160405280600081526020016000815260200160008152602001600081525090565b600080610a8260048286886117c8565b610a8b9161189e565b60e01c9050801580610aa057506001861b8110155b15610ad7576040517ff4e04eaa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81600116600003610afc5760019190911c90610af58161184d565b9050610ada565b600192506000610b53898360208110610b1757610b17611786565b60040201604051806080016040529081600082015481526020016001820154815260200160028201548152602001600382015481525050610d99565b9050610b5e8261184d565b9150600183901c92505b8215610bee576001831615610bd757610bc9818a8460208110610b8d57610b8d611786565b60040201604051806080016040529081600082015481526020016001820154815260200160028201548152602001600382015481525050610e08565b9050610bd48461184d565b93505b60019290921c91610be78261184d565b9150610b68565b610bf7816110a3565b945050505094509492505050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8216817f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478210610c6057506000915081905080610d92565b60007f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47857f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4787880909089050610d07816002610d007f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760016117b5565b901c61114f565b9150807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4783840914610d4457600080600093509350935050610d92565b60ff85901c600183168114610d8b577f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47610d7e84826118e6565b610d88919061172f565b92505b6001945050505b9193909250565b610dd26040518060c001604052806000815260200160008152602001600081526020016000815260200160008152602001600081525090565b8151815260208083015190820152604080830151908201526060918201519181019190915260006080820152600160a082015290565b610e416040518060c001604052806000815260200160008152602001600081526020016000815260200160008152602001600081525090565b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47608084015160a08501518281830983818208905083848485098503858485090884858286098684860908858684870987038784870908885160208a015188898684098a88840908898a8885098b038b888509088d51935060208e015192508a848c03830891508a838c03820890508a828b0899508a818a089850898b8a8c099a508b8b8c089a508b8c8283098d038d8c8d09089950508a888c038b0899508a878c038a0898508a81830997508a88890897508a8b8384098c038c8384090896508a888c038b0899508a878c038a0898508960808d01528860a08d01528a8860040999508a8760040998508a8b8a84098c8c84090897508a8b8b84098c038c8b84090896508a8b8a86098c8c86090891508a8b8b86098c038c8b860908905060408d0151995060608d015198508a8b868c098c888c090893508a8b878c098c038c878c0908925060408e0151995060608e015198508a8a8c03850895508a898c03840894508a8660020993508a8560020992508a83850995508a86870895508a8b8586098c038c8586090894508a888c03870895508a878c03860894508a8b888c098c8a8c09088b8c8a8d098d038d8a8d09088c826002099b508c816002099a5050508a8260020997508a8160020996508a888c03870897508a878c0386089650878c528660208d01528a888c03830897508a878c0382089650505088898684098a88840908935088898784098a038a878409089250505086868803830895508685880382089450505050508160408501528060608501525050505b92915050565b6110ce6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000806110e384608001518560a001516111b2565b915091506000806110f484846112b6565b9150915061110c82828860000151896020015161137a565b6020870152855261111f8282868661137a565b809250819350505061113b82828860400151896060015161137a565b606087015260408601525092949350505050565b600060405160208152602080820152602060408201528360608201528260808201527f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4760a082015260208160c08360055afa6111aa57600080fd5b519392505050565b600080807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808687097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478687090890506112368161123160027f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd476118e6565b61114f565b90507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4781611284877f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd476118e6565b0992507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478185099150505b9250929050565b6000807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4784840991507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4782830891507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808586097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478586090890509250929050565b6000807f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808488097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478688090891507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47808588097f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4785880908905094509492505050565b6040518061040001604052806020905b61147b6040518060800160405280600081526020016000815260200160008152602001600081525090565b8152602001906001900390816114505790505090565b6040518060a001604052806005906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b803560ff811681146114fc57600080fd5b919050565b60008083601f84011261151357600080fd5b50813567ffffffffffffffff81111561152b57600080fd5b6020830191508360208285010111156112af57600080fd5b60008060008060006080868803121561155b57600080fd5b853563ffffffff8116811461156f57600080fd5b945061157d602087016114eb565b935061158b604087016114eb565b9250606086013567ffffffffffffffff8111156115a757600080fd5b6115b388828901611501565b969995985093965092949392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008060008060006080868803121561160b57600080fd5b85359450602086013567ffffffffffffffff808216821461162b57600080fd5b9094506040870135908082111561164157600080fd5b818801915088601f83011261165557600080fd5b813581811115611667576116676115c4565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156116ad576116ad6115c4565b816040528281528b60208487010111156116c657600080fd5b8260208601602083013760006020848301015280975050505060608801359150808211156116f357600080fd5b506115b388828901611501565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261173e5761173e611700565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008261178157611781611700565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111561109d5761109d611743565b600080858511156117d857600080fd5b838611156117e557600080fd5b5050820193919092039150565b8035602083101561109d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b600060ff821660ff810361184457611844611743565b60010192915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361187e5761187e611743565b5060010190565b60ff818116838216019081111561109d5761109d611743565b7fffffffff0000000000000000000000000000000000000000000000000000000081358181169160048510156118de5780818660040360031b1b83161692505b505092915050565b8181038181111561109d5761109d61174356fea164736f6c6343000813000a",
}

// DemoBLSAttestationVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use DemoBLSAttestationVerifierMetaData.ABI instead.
var DemoBLSAttestationVerifierABI = DemoBLSAttestationVerifierMetaData.ABI

// DemoBLSAttestationVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DemoBLSAttestationVerifierMetaData.Bin instead.
var DemoBLSAttestationVerifierBin = DemoBLSAttestationVerifierMetaData.Bin

// DeployDemoBLSAttestationVerifier deploys a new Ethereum contract, binding an instance of DemoBLSAttestationVerifier to it.
func DeployDemoBLSAttestationVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DemoBLSAttestationVerifier, error) {
	parsed, err := DemoBLSAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DemoBLSAttestationVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DemoBLSAttestationVerifier{DemoBLSAttestationVerifierCaller: DemoBLSAttestationVerifierCaller{contract: contract}, DemoBLSAttestationVerifierTransactor: DemoBLSAttestationVerifierTransactor{contract: contract}, DemoBLSAttestationVerifierFilterer: DemoBLSAttestationVerifierFilterer{contract: contract}}, nil
}

// DemoBLSAttestationVerifier is an auto generated Go binding around an Ethereum contract.
type DemoBLSAttestationVerifier struct {
	DemoBLSAttestationVerifierCaller     // Read-only binding to the contract
	DemoBLSAttestationVerifierTransactor // Write-only binding to the contract
	DemoBLSAttestationVerifierFilterer   // Log filterer for contract events
}

// DemoBLSAttestationVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type DemoBLSAttestationVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoBLSAttestationVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DemoBLSAttestationVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoBLSAttestationVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DemoBLSAttestationVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoBLSAttestationVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DemoBLSAttestationVerifierSession struct {
	Contract     *DemoBLSAttestationVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DemoBLSAttestationVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DemoBLSAttestationVerifierCallerSession struct {
	Contract *DemoBLSAttestationVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// DemoBLSAttestationVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DemoBLSAttestationVerifierTransactorSession struct {
	Contract     *DemoBLSAttestationVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// DemoBLSAttestationVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type DemoBLSAttestationVerifierRaw struct {
	Contract *DemoBLSAttestationVerifier // Generic contract binding to access the raw methods on
}

// DemoBLSAttestationVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DemoBLSAttestationVerifierCallerRaw struct {
	Contract *DemoBLSAttestationVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// DemoBLSAttestationVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DemoBLSAttestationVerifierTransactorRaw struct {
	Contract *DemoBLSAttestationVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDemoBLSAttestationVerifier creates a new instance of DemoBLSAttestationVerifier, bound to a specific deployed contract.
func NewDemoBLSAttestationVerifier(address common.Address, backend bind.ContractBackend) (*DemoBLSAttestationVerifier, error) {
	contract, err := bindDemoBLSAttestationVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DemoBLSAttestationVerifier{DemoBLSAttestationVerifierCaller: DemoBLSAttestationVerifierCaller{contract: contract}, DemoBLSAttestationVerifierTransactor: DemoBLSAttestationVerifierTransactor{contract: contract}, DemoBLSAttestationVerifierFilterer: DemoBLSAttestationVerifierFilterer{contract: contract}}, nil
}

// NewDemoBLSAttestationVerifierCaller creates a new read-only instance of DemoBLSAttestationVerifier, bound to a specific deployed contract.
func NewDemoBLSAttestationVerifierCaller(address common.Address, caller bind.ContractCaller) (*DemoBLSAttestationVerifierCaller, error) {
	contract, err := bindDemoBLSAttestationVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DemoBLSAttestationVerifierCaller{contract: contract}, nil
}

// NewDemoBLSAttestationVerifierTransactor creates a new write-only instance of DemoBLSAttestationVerifier, bound to a specific deployed contract.
func NewDemoBLSAttestationVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*DemoBLSAttestationVerifierTransactor, error) {
	contract, err := bindDemoBLSAttestationVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DemoBLSAttestationVerifierTransactor{contract: contract}, nil
}

// NewDemoBLSAttestationVerifierFilterer creates a new log filterer instance of DemoBLSAttestationVerifier, bound to a specific deployed contract.
func NewDemoBLSAttestationVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*DemoBLSAttestationVerifierFilterer, error) {
	contract, err := bindDemoBLSAttestationVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DemoBLSAttestationVerifierFilterer{contract: contract}, nil
}

// bindDemoBLSAttestationVerifier binds a generic wrapper to an already deployed contract.
func bindDemoBLSAttestationVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DemoBLSAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DemoBLSAttestationVerifier.Contract.DemoBLSAttestationVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.Contract.DemoBLSAttestationVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.Contract.DemoBLSAttestationVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DemoBLSAttestationVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.Contract.contract.Transact(opts, method, params...)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes ocr3BlsSignerPublicKeys) returns()
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierTransactor) SetConfig(opts *bind.TransactOpts, configVersion uint32, n uint8, f uint8, ocr3BlsSignerPublicKeys []byte) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.contract.Transact(opts, "setConfig", configVersion, n, f, ocr3BlsSignerPublicKeys)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes ocr3BlsSignerPublicKeys) returns()
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierSession) SetConfig(configVersion uint32, n uint8, f uint8, ocr3BlsSignerPublicKeys []byte) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.Contract.SetConfig(&_DemoBLSAttestationVerifier.TransactOpts, configVersion, n, f, ocr3BlsSignerPublicKeys)
}

// SetConfig is a paid mutator transaction binding the contract method 0x6f289d41.
//
// Solidity: function setConfig(uint32 configVersion, uint8 n, uint8 f, bytes ocr3BlsSignerPublicKeys) returns()
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierTransactorSession) SetConfig(configVersion uint32, n uint8, f uint8, ocr3BlsSignerPublicKeys []byte) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.Contract.SetConfig(&_DemoBLSAttestationVerifier.TransactOpts, configVersion, n, f, ocr3BlsSignerPublicKeys)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierTransactor) Transmit(opts *bind.TransactOpts, configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.contract.Transact(opts, "transmit", configDigest, seqNr, report, attestation)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierSession) Transmit(configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.Contract.Transmit(&_DemoBLSAttestationVerifier.TransactOpts, configDigest, seqNr, report, attestation)
}

// Transmit is a paid mutator transaction binding the contract method 0x95fd15c8.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation) returns()
func (_DemoBLSAttestationVerifier *DemoBLSAttestationVerifierTransactorSession) Transmit(configDigest [32]byte, seqNr uint64, report []byte, attestation []byte) (*types.Transaction, error) {
	return _DemoBLSAttestationVerifier.Contract.Transmit(&_DemoBLSAttestationVerifier.TransactOpts, configDigest, seqNr, report, attestation)
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

// OCR3BLSAttestationVerifierMetaData contains all meta data concerning the OCR3BLSAttestationVerifier contract.
var OCR3BLSAttestationVerifierMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x6080604052348015600f57600080fd5b50601680601d6000396000f3fe6080604052600080fdfea164736f6c6343000813000a",
}

// OCR3BLSAttestationVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use OCR3BLSAttestationVerifierMetaData.ABI instead.
var OCR3BLSAttestationVerifierABI = OCR3BLSAttestationVerifierMetaData.ABI

// OCR3BLSAttestationVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OCR3BLSAttestationVerifierMetaData.Bin instead.
var OCR3BLSAttestationVerifierBin = OCR3BLSAttestationVerifierMetaData.Bin

// DeployOCR3BLSAttestationVerifier deploys a new Ethereum contract, binding an instance of OCR3BLSAttestationVerifier to it.
func DeployOCR3BLSAttestationVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OCR3BLSAttestationVerifier, error) {
	parsed, err := OCR3BLSAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR3BLSAttestationVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR3BLSAttestationVerifier{OCR3BLSAttestationVerifierCaller: OCR3BLSAttestationVerifierCaller{contract: contract}, OCR3BLSAttestationVerifierTransactor: OCR3BLSAttestationVerifierTransactor{contract: contract}, OCR3BLSAttestationVerifierFilterer: OCR3BLSAttestationVerifierFilterer{contract: contract}}, nil
}

// OCR3BLSAttestationVerifier is an auto generated Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifier struct {
	OCR3BLSAttestationVerifierCaller     // Read-only binding to the contract
	OCR3BLSAttestationVerifierTransactor // Write-only binding to the contract
	OCR3BLSAttestationVerifierFilterer   // Log filterer for contract events
}

// OCR3BLSAttestationVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3BLSAttestationVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3BLSAttestationVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OCR3BLSAttestationVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OCR3BLSAttestationVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OCR3BLSAttestationVerifierSession struct {
	Contract     *OCR3BLSAttestationVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// OCR3BLSAttestationVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OCR3BLSAttestationVerifierCallerSession struct {
	Contract *OCR3BLSAttestationVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// OCR3BLSAttestationVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OCR3BLSAttestationVerifierTransactorSession struct {
	Contract     *OCR3BLSAttestationVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// OCR3BLSAttestationVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierRaw struct {
	Contract *OCR3BLSAttestationVerifier // Generic contract binding to access the raw methods on
}

// OCR3BLSAttestationVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierCallerRaw struct {
	Contract *OCR3BLSAttestationVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// OCR3BLSAttestationVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OCR3BLSAttestationVerifierTransactorRaw struct {
	Contract *OCR3BLSAttestationVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOCR3BLSAttestationVerifier creates a new instance of OCR3BLSAttestationVerifier, bound to a specific deployed contract.
func NewOCR3BLSAttestationVerifier(address common.Address, backend bind.ContractBackend) (*OCR3BLSAttestationVerifier, error) {
	contract, err := bindOCR3BLSAttestationVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR3BLSAttestationVerifier{OCR3BLSAttestationVerifierCaller: OCR3BLSAttestationVerifierCaller{contract: contract}, OCR3BLSAttestationVerifierTransactor: OCR3BLSAttestationVerifierTransactor{contract: contract}, OCR3BLSAttestationVerifierFilterer: OCR3BLSAttestationVerifierFilterer{contract: contract}}, nil
}

// NewOCR3BLSAttestationVerifierCaller creates a new read-only instance of OCR3BLSAttestationVerifier, bound to a specific deployed contract.
func NewOCR3BLSAttestationVerifierCaller(address common.Address, caller bind.ContractCaller) (*OCR3BLSAttestationVerifierCaller, error) {
	contract, err := bindOCR3BLSAttestationVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3BLSAttestationVerifierCaller{contract: contract}, nil
}

// NewOCR3BLSAttestationVerifierTransactor creates a new write-only instance of OCR3BLSAttestationVerifier, bound to a specific deployed contract.
func NewOCR3BLSAttestationVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR3BLSAttestationVerifierTransactor, error) {
	contract, err := bindOCR3BLSAttestationVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR3BLSAttestationVerifierTransactor{contract: contract}, nil
}

// NewOCR3BLSAttestationVerifierFilterer creates a new log filterer instance of OCR3BLSAttestationVerifier, bound to a specific deployed contract.
func NewOCR3BLSAttestationVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR3BLSAttestationVerifierFilterer, error) {
	contract, err := bindOCR3BLSAttestationVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR3BLSAttestationVerifierFilterer{contract: contract}, nil
}

// bindOCR3BLSAttestationVerifier binds a generic wrapper to an already deployed contract.
func bindOCR3BLSAttestationVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OCR3BLSAttestationVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3BLSAttestationVerifier *OCR3BLSAttestationVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3BLSAttestationVerifier.Contract.OCR3BLSAttestationVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3BLSAttestationVerifier *OCR3BLSAttestationVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3BLSAttestationVerifier.Contract.OCR3BLSAttestationVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3BLSAttestationVerifier *OCR3BLSAttestationVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3BLSAttestationVerifier.Contract.OCR3BLSAttestationVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OCR3BLSAttestationVerifier *OCR3BLSAttestationVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR3BLSAttestationVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OCR3BLSAttestationVerifier *OCR3BLSAttestationVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR3BLSAttestationVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OCR3BLSAttestationVerifier *OCR3BLSAttestationVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR3BLSAttestationVerifier.Contract.contract.Transact(opts, method, params...)
}

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
