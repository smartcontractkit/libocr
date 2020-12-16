


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


var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)


const AggregatorValidatorInterfaceABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"previousRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"previousAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"currentRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"currentAnswer\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


type AggregatorValidatorInterface struct {
	AggregatorValidatorInterfaceCaller     
	AggregatorValidatorInterfaceTransactor 
	AggregatorValidatorInterfaceFilterer   
}


type AggregatorValidatorInterfaceCaller struct {
	contract *bind.BoundContract 
}


type AggregatorValidatorInterfaceTransactor struct {
	contract *bind.BoundContract 
}


type AggregatorValidatorInterfaceFilterer struct {
	contract *bind.BoundContract 
}



type AggregatorValidatorInterfaceSession struct {
	Contract     *AggregatorValidatorInterface 
	CallOpts     bind.CallOpts                 
	TransactOpts bind.TransactOpts             
}



type AggregatorValidatorInterfaceCallerSession struct {
	Contract *AggregatorValidatorInterfaceCaller 
	CallOpts bind.CallOpts                       
}



type AggregatorValidatorInterfaceTransactorSession struct {
	Contract     *AggregatorValidatorInterfaceTransactor 
	TransactOpts bind.TransactOpts                       
}


type AggregatorValidatorInterfaceRaw struct {
	Contract *AggregatorValidatorInterface 
}


type AggregatorValidatorInterfaceCallerRaw struct {
	Contract *AggregatorValidatorInterfaceCaller 
}


type AggregatorValidatorInterfaceTransactorRaw struct {
	Contract *AggregatorValidatorInterfaceTransactor 
}


func NewAggregatorValidatorInterface(address common.Address, backend bind.ContractBackend) (*AggregatorValidatorInterface, error) {
	contract, err := bindAggregatorValidatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterface{AggregatorValidatorInterfaceCaller: AggregatorValidatorInterfaceCaller{contract: contract}, AggregatorValidatorInterfaceTransactor: AggregatorValidatorInterfaceTransactor{contract: contract}, AggregatorValidatorInterfaceFilterer: AggregatorValidatorInterfaceFilterer{contract: contract}}, nil
}


func NewAggregatorValidatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorValidatorInterfaceCaller, error) {
	contract, err := bindAggregatorValidatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceCaller{contract: contract}, nil
}


func NewAggregatorValidatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorValidatorInterfaceTransactor, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceTransactor{contract: contract}, nil
}


func NewAggregatorValidatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorValidatorInterfaceFilterer, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceFilterer{contract: contract}, nil
}


func bindAggregatorValidatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorValidatorInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transfer(opts)
}


func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.contract.Call(opts, result, method, params...)
}



func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transfer(opts)
}


func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transact(opts, method, params...)
}




func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactor) Validate(opts *bind.TransactOpts, previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.contract.Transact(opts, "validate", previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}




func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}




func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}


const TestValidatorABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"


var TestValidatorBin = "0x6080604052348015600f57600080fd5b5060848061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063beed9b5114602d575b600080fd5b605960048036036080811015604157600080fd5b5080359060208101359060408101359060600135606d565b604080519115158252519081900360200190f35b600194935050505056fea164736f6c6343000705000a"


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


type TestValidator struct {
	TestValidatorCaller     
	TestValidatorTransactor 
	TestValidatorFilterer   
}


type TestValidatorCaller struct {
	contract *bind.BoundContract 
}


type TestValidatorTransactor struct {
	contract *bind.BoundContract 
}


type TestValidatorFilterer struct {
	contract *bind.BoundContract 
}



type TestValidatorSession struct {
	Contract     *TestValidator    
	CallOpts     bind.CallOpts     
	TransactOpts bind.TransactOpts 
}



type TestValidatorCallerSession struct {
	Contract *TestValidatorCaller 
	CallOpts bind.CallOpts        
}



type TestValidatorTransactorSession struct {
	Contract     *TestValidatorTransactor 
	TransactOpts bind.TransactOpts        
}


type TestValidatorRaw struct {
	Contract *TestValidator 
}


type TestValidatorCallerRaw struct {
	Contract *TestValidatorCaller 
}


type TestValidatorTransactorRaw struct {
	Contract *TestValidatorTransactor 
}


func NewTestValidator(address common.Address, backend bind.ContractBackend) (*TestValidator, error) {
	contract, err := bindTestValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestValidator{TestValidatorCaller: TestValidatorCaller{contract: contract}, TestValidatorTransactor: TestValidatorTransactor{contract: contract}, TestValidatorFilterer: TestValidatorFilterer{contract: contract}}, nil
}


func NewTestValidatorCaller(address common.Address, caller bind.ContractCaller) (*TestValidatorCaller, error) {
	contract, err := bindTestValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestValidatorCaller{contract: contract}, nil
}


func NewTestValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*TestValidatorTransactor, error) {
	contract, err := bindTestValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestValidatorTransactor{contract: contract}, nil
}


func NewTestValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*TestValidatorFilterer, error) {
	contract, err := bindTestValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestValidatorFilterer{contract: contract}, nil
}


func bindTestValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestValidatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_TestValidator *TestValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestValidator.Contract.TestValidatorCaller.contract.Call(opts, result, method, params...)
}



func (_TestValidator *TestValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestValidator.Contract.TestValidatorTransactor.contract.Transfer(opts)
}


func (_TestValidator *TestValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestValidator.Contract.TestValidatorTransactor.contract.Transact(opts, method, params...)
}





func (_TestValidator *TestValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestValidator.Contract.contract.Call(opts, result, method, params...)
}



func (_TestValidator *TestValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestValidator.Contract.contract.Transfer(opts)
}


func (_TestValidator *TestValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestValidator.Contract.contract.Transact(opts, method, params...)
}




func (_TestValidator *TestValidatorCaller) Validate(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int, arg2 *big.Int, arg3 *big.Int) (bool, error) {
	var out []interface{}
	err := _TestValidator.contract.Call(opts, &out, "validate", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_TestValidator *TestValidatorSession) Validate(arg0 *big.Int, arg1 *big.Int, arg2 *big.Int, arg3 *big.Int) (bool, error) {
	return _TestValidator.Contract.Validate(&_TestValidator.CallOpts, arg0, arg1, arg2, arg3)
}




func (_TestValidator *TestValidatorCallerSession) Validate(arg0 *big.Int, arg1 *big.Int, arg2 *big.Int, arg3 *big.Int) (bool, error) {
	return _TestValidator.Contract.Validate(&_TestValidator.CallOpts, arg0, arg1, arg2, arg3)
}
