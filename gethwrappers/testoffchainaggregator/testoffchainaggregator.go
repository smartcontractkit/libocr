


package testoffchainaggregator

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


type OffchainAggregatorHotVars struct {
	LatestConfigDigest      [16]byte
	LatestEpochAndRound     *big.Int
	Threshold               uint8
	LatestAggregatorRoundId uint32
}


const AccessControlTestHelperABI = "[{\"anonymous\":false,\"inputs\":[],\"name\":\"Dummy\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_roundID\",\"type\":\"uint256\"}],\"name\":\"readGetAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"},{\"internalType\":\"uint80\",\"name\":\"_roundID\",\"type\":\"uint80\"}],\"name\":\"readGetRoundData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_roundID\",\"type\":\"uint256\"}],\"name\":\"readGetTimestamp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"readLatestAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"readLatestRound\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"readLatestRoundData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"readLatestTimestamp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"testLatestTransmissionDetails\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]"


var AccessControlTestHelperBin = "0x608060405234801561001057600080fd5b506105e3806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063c0c9c7db1161005b578063c0c9c7db14610173578063c9592ab9146101a6578063d2f79c47146101df578063eea2913a1461021257610088565b806304cefda51461008d57806320f2c97c146100c257806395319deb146100f5578063bf5fc18b1461013a575b600080fd5b6100c0600480360360208110156100a357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610245565b005b6100c0600480360360208110156100d857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166102ba565b6100c06004803603604081101561010b57600080fd5b50803573ffffffffffffffffffffffffffffffffffffffff16906020013569ffffffffffffffffffff16610358565b6100c06004803603604081101561015057600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020013561040e565b6100c06004803603602081101561018957600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610489565b6100c0600480360360408110156101bc57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356104f9565b6100c0600480360360208110156101f557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661054a565b6100c06004803603602081101561022857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610590565b8073ffffffffffffffffffffffffffffffffffffffff1663e5fe45776040518163ffffffff1660e01b815260040160a06040518083038186803b15801561028b57600080fd5b505afa15801561029f573d6000803e3d6000fd5b505050506040513d60a08110156102b557600080fd5b505050565b8073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561030057600080fd5b505afa158015610314573d6000803e3d6000fd5b505050506040513d60a081101561032a57600080fd5b50506040517f10e4ab9f2ce395bf5539d7c60c9bfeef0b416602954734c5bb8bfd9433c9ff6890600090a150565b8173ffffffffffffffffffffffffffffffffffffffff16639a6fc8f5826040518263ffffffff1660e01b8152600401808269ffffffffffffffffffff16815260200191505060a06040518083038186803b1580156103b557600080fd5b505afa1580156103c9573d6000803e3d6000fd5b505050506040513d60a08110156103df57600080fd5b50506040517f10e4ab9f2ce395bf5539d7c60c9bfeef0b416602954734c5bb8bfd9433c9ff6890600090a15050565b8173ffffffffffffffffffffffffffffffffffffffff1663b5ab58dc826040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561045f57600080fd5b505afa158015610473573d6000803e3d6000fd5b505050506040513d60208110156103df57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff166350d25bcd6040518163ffffffff1660e01b815260040160206040518083038186803b1580156104cf57600080fd5b505afa1580156104e3573d6000803e3d6000fd5b505050506040513d602081101561032a57600080fd5b8173ffffffffffffffffffffffffffffffffffffffff1663b633620c826040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561045f57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff16638205bf6a6040518163ffffffff1660e01b815260040160206040518083038186803b1580156104cf57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff1663668a0f026040518163ffffffff1660e01b815260040160206040518083038186803b1580156104cf57600080fdfea164736f6c6343000705000a"


func DeployAccessControlTestHelper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AccessControlTestHelper, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlTestHelperABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccessControlTestHelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessControlTestHelper{AccessControlTestHelperCaller: AccessControlTestHelperCaller{contract: contract}, AccessControlTestHelperTransactor: AccessControlTestHelperTransactor{contract: contract}, AccessControlTestHelperFilterer: AccessControlTestHelperFilterer{contract: contract}}, nil
}


type AccessControlTestHelper struct {
	AccessControlTestHelperCaller     
	AccessControlTestHelperTransactor 
	AccessControlTestHelperFilterer   
}


type AccessControlTestHelperCaller struct {
	contract *bind.BoundContract 
}


type AccessControlTestHelperTransactor struct {
	contract *bind.BoundContract 
}


type AccessControlTestHelperFilterer struct {
	contract *bind.BoundContract 
}



type AccessControlTestHelperSession struct {
	Contract     *AccessControlTestHelper 
	CallOpts     bind.CallOpts            
	TransactOpts bind.TransactOpts        
}



type AccessControlTestHelperCallerSession struct {
	Contract *AccessControlTestHelperCaller 
	CallOpts bind.CallOpts                  
}



type AccessControlTestHelperTransactorSession struct {
	Contract     *AccessControlTestHelperTransactor 
	TransactOpts bind.TransactOpts                  
}


type AccessControlTestHelperRaw struct {
	Contract *AccessControlTestHelper 
}


type AccessControlTestHelperCallerRaw struct {
	Contract *AccessControlTestHelperCaller 
}


type AccessControlTestHelperTransactorRaw struct {
	Contract *AccessControlTestHelperTransactor 
}


func NewAccessControlTestHelper(address common.Address, backend bind.ContractBackend) (*AccessControlTestHelper, error) {
	contract, err := bindAccessControlTestHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControlTestHelper{AccessControlTestHelperCaller: AccessControlTestHelperCaller{contract: contract}, AccessControlTestHelperTransactor: AccessControlTestHelperTransactor{contract: contract}, AccessControlTestHelperFilterer: AccessControlTestHelperFilterer{contract: contract}}, nil
}


func NewAccessControlTestHelperCaller(address common.Address, caller bind.ContractCaller) (*AccessControlTestHelperCaller, error) {
	contract, err := bindAccessControlTestHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlTestHelperCaller{contract: contract}, nil
}


func NewAccessControlTestHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControlTestHelperTransactor, error) {
	contract, err := bindAccessControlTestHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlTestHelperTransactor{contract: contract}, nil
}


func NewAccessControlTestHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControlTestHelperFilterer, error) {
	contract, err := bindAccessControlTestHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControlTestHelperFilterer{contract: contract}, nil
}


func bindAccessControlTestHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlTestHelperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_AccessControlTestHelper *AccessControlTestHelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlTestHelper.Contract.AccessControlTestHelperCaller.contract.Call(opts, result, method, params...)
}



func (_AccessControlTestHelper *AccessControlTestHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.AccessControlTestHelperTransactor.contract.Transfer(opts)
}


func (_AccessControlTestHelper *AccessControlTestHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.AccessControlTestHelperTransactor.contract.Transact(opts, method, params...)
}





func (_AccessControlTestHelper *AccessControlTestHelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlTestHelper.Contract.contract.Call(opts, result, method, params...)
}



func (_AccessControlTestHelper *AccessControlTestHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.contract.Transfer(opts)
}


func (_AccessControlTestHelper *AccessControlTestHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.contract.Transact(opts, method, params...)
}




func (_AccessControlTestHelper *AccessControlTestHelperCaller) TestLatestTransmissionDetails(opts *bind.CallOpts, _aggregator common.Address) error {
	var out []interface{}
	err := _AccessControlTestHelper.contract.Call(opts, &out, "testLatestTransmissionDetails", _aggregator)

	if err != nil {
		return err
	}

	return err

}




func (_AccessControlTestHelper *AccessControlTestHelperSession) TestLatestTransmissionDetails(_aggregator common.Address) error {
	return _AccessControlTestHelper.Contract.TestLatestTransmissionDetails(&_AccessControlTestHelper.CallOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperCallerSession) TestLatestTransmissionDetails(_aggregator common.Address) error {
	return _AccessControlTestHelper.Contract.TestLatestTransmissionDetails(&_AccessControlTestHelper.CallOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactor) ReadGetAnswer(opts *bind.TransactOpts, _aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.contract.Transact(opts, "readGetAnswer", _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperSession) ReadGetAnswer(_aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadGetAnswer(&_AccessControlTestHelper.TransactOpts, _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactorSession) ReadGetAnswer(_aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadGetAnswer(&_AccessControlTestHelper.TransactOpts, _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactor) ReadGetRoundData(opts *bind.TransactOpts, _aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.contract.Transact(opts, "readGetRoundData", _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperSession) ReadGetRoundData(_aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadGetRoundData(&_AccessControlTestHelper.TransactOpts, _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactorSession) ReadGetRoundData(_aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadGetRoundData(&_AccessControlTestHelper.TransactOpts, _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactor) ReadGetTimestamp(opts *bind.TransactOpts, _aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.contract.Transact(opts, "readGetTimestamp", _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperSession) ReadGetTimestamp(_aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadGetTimestamp(&_AccessControlTestHelper.TransactOpts, _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactorSession) ReadGetTimestamp(_aggregator common.Address, _roundID *big.Int) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadGetTimestamp(&_AccessControlTestHelper.TransactOpts, _aggregator, _roundID)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactor) ReadLatestAnswer(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.contract.Transact(opts, "readLatestAnswer", _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperSession) ReadLatestAnswer(_aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadLatestAnswer(&_AccessControlTestHelper.TransactOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactorSession) ReadLatestAnswer(_aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadLatestAnswer(&_AccessControlTestHelper.TransactOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactor) ReadLatestRound(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.contract.Transact(opts, "readLatestRound", _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperSession) ReadLatestRound(_aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadLatestRound(&_AccessControlTestHelper.TransactOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactorSession) ReadLatestRound(_aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadLatestRound(&_AccessControlTestHelper.TransactOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactor) ReadLatestRoundData(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.contract.Transact(opts, "readLatestRoundData", _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperSession) ReadLatestRoundData(_aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadLatestRoundData(&_AccessControlTestHelper.TransactOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactorSession) ReadLatestRoundData(_aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadLatestRoundData(&_AccessControlTestHelper.TransactOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactor) ReadLatestTimestamp(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.contract.Transact(opts, "readLatestTimestamp", _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperSession) ReadLatestTimestamp(_aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadLatestTimestamp(&_AccessControlTestHelper.TransactOpts, _aggregator)
}




func (_AccessControlTestHelper *AccessControlTestHelperTransactorSession) ReadLatestTimestamp(_aggregator common.Address) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.ReadLatestTimestamp(&_AccessControlTestHelper.TransactOpts, _aggregator)
}


type AccessControlTestHelperDummyIterator struct {
	Event *AccessControlTestHelperDummy 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlTestHelperDummyIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlTestHelperDummy)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlTestHelperDummy)
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


func (it *AccessControlTestHelperDummyIterator) Error() error {
	return it.fail
}



func (it *AccessControlTestHelperDummyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlTestHelperDummy struct {
	Raw types.Log 
}




func (_AccessControlTestHelper *AccessControlTestHelperFilterer) FilterDummy(opts *bind.FilterOpts) (*AccessControlTestHelperDummyIterator, error) {

	logs, sub, err := _AccessControlTestHelper.contract.FilterLogs(opts, "Dummy")
	if err != nil {
		return nil, err
	}
	return &AccessControlTestHelperDummyIterator{contract: _AccessControlTestHelper.contract, event: "Dummy", logs: logs, sub: sub}, nil
}




func (_AccessControlTestHelper *AccessControlTestHelperFilterer) WatchDummy(opts *bind.WatchOpts, sink chan<- *AccessControlTestHelperDummy) (event.Subscription, error) {

	logs, sub, err := _AccessControlTestHelper.contract.WatchLogs(opts, "Dummy")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlTestHelperDummy)
				if err := _AccessControlTestHelper.contract.UnpackLog(event, "Dummy", log); err != nil {
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




func (_AccessControlTestHelper *AccessControlTestHelperFilterer) ParseDummy(log types.Log) (*AccessControlTestHelperDummy, error) {
	event := new(AccessControlTestHelperDummy)
	if err := _AccessControlTestHelper.contract.UnpackLog(event, "Dummy", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const AccessControlledOffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var AccessControlledOffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b5060405162005c2438038062005c2483398181016040526101808110156200003957600080fd5b815160208301516040808501516060860151608087015160a088015160c089015160e08a01516101008b01516101208c01516101408d01516101608e0180519a519c9e9b9d999c989b979a969995989497939692959194939182019284640100000000821115620000a957600080fd5b908301906020820185811115620000bf57600080fd5b8251640100000000811182820188101715620000da57600080fd5b82525081516020918201929091019080838360005b8381101562000109578181015183820152602001620000ef565b50505050905090810190601f168015620001375780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b03191633179055508b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b88620001728787878787620002ac565b6200017d816200039e565b6001600160601b0319606083901b1660805262000199620004fc565b620001a3620004fc565b60005b601f8160ff161015620001f3576001838260ff16601f8110620001c557fe5b61ffff909216602092909202015260018260ff8316601f8110620001e557fe5b6020020152600101620001a6565b5062000203600483601f6200051b565b5062000213600882601f620005b8565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e052505083516200025c9350602d9250602085019150620005e9565b50620002688662000417565b505050601791820b820b604090811b60a05290820b90910b901b60c0525050602e805460ff1916600117905550620006829f50505050505050505050505050505050565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b0390811690821681146200041357600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b0316331462000477576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b0368010000000000000000909104811690821681146200041357602c8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b604051806103e00160405280601f906020820280368337509192915050565b600283019183908215620005a65791602002820160005b838211156200057457835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030262000532565b8015620005a45782816101000a81549061ffff021916905560020160208160010104928301926001030262000574565b505b50620005b49291506200066b565b5090565b82601f8101928215620005a6579160200282015b82811115620005a6578251825591602001919060010190620005cc565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620006215760008555620005a6565b82601f106200063c57805160ff1916838001178555620005a6565b82800160010185558215620005a65791820182811115620005a6578251825591602001919060010190620005cc565b5b80821115620005b457600081556001016200066c565b60805160601c60a05160401c60c05160401c60e05160f81c61553d620006e760003980610ff7525080611b21528061351b525080610f3e52806134ee525080610f1a528061286c528061295c52806139555280613ff55280614858525061553d6000f3fe608060405234801561001057600080fd5b50600436106102c85760003560e01c80638823da6c1161017b578063c1075329116100d8578063e5fe45771161008c578063f2fde38b11610071578063f2fde38b14610bad578063fbffd2c114610bd3578063feaf968c14610bf9576102c8565b8063e5fe457714610b15578063eb5dcd6c14610b7f576102c8565b8063d09dc339116100bd578063d09dc33914610ac8578063dc7f012414610ad0578063e4902f8214610ad8576102c8565b8063c107532914610988578063c9807539146109b4576102c8565b8063a118f2491161012f578063b5ab58dc11610114578063b5ab58dc14610909578063b633620c14610926578063bd82470614610943576102c8565b8063a118f249146108bd578063b121e147146108e3576102c8565b80638da5cb5b116101605780638da5cb5b146107805780639a6fc8f5146107885780639c849b30146107fb576102c8565b80638823da6c146107345780638ac28d5a1461075a576102c8565b8063585aa7de1161022957806379ba5097116101dd57806381411834116101c2578063814118341461068357806381ff7048146106db5780638205bf6a1461072c576102c8565b806379ba5097146106735780638038e4a11461067b576102c8565b80636b14daf81161020e5780636b14daf81461052457806370da2f67146105ee5780637284e416146105f6576102c8565b8063585aa7de146103ef578063668a0f021461051c576102c8565b806329937268116102805780633a5381b5116102655780633a5381b5146103d757806350d25bcd146103df57806354fd4d50146103e7576102c8565b80632993726814610378578063313ce567146103b9576102c8565b80631327d3d8116102b15780631327d3d81461030f5780631b6b6d231461033557806322adbc7814610359576102c8565b80630a756983146102cd5780630eafb25b146102d7575b600080fd5b6102d5610c01565b005b6102fd600480360360208110156102ed57600080fd5b50356001600160a01b0316610cbf565b60408051918252519081900360200190f35b6102d56004803603602081101561032557600080fd5b50356001600160a01b0316610e1f565b61033d610f18565b604080516001600160a01b039092168252519081900360200190f35b610361610f3c565b6040805160179290920b8252519081900360200190f35b610380610f60565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b6103c1610ff5565b6040805160ff9092168252519081900360200190f35b61033d611019565b6102fd611034565b6102fd6110d5565b6102d5600480360360a081101561040557600080fd5b81019060208101813564010000000081111561042057600080fd5b82018360208201111561043257600080fd5b8035906020019184602083028401116401000000008311171561045457600080fd5b91939092909160208101903564010000000081111561047257600080fd5b82018360208201111561048457600080fd5b803590602001918460208302840111640100000000831117156104a657600080fd5b9193909260ff8335169267ffffffffffffffff6020820135169291906060810190604001356401000000008111156104dd57600080fd5b8201836020820111156104ef57600080fd5b8035906020019184600183028401116401000000008311171561051157600080fd5b5090925090506110da565b6102fd611a5b565b6105da6004803603604081101561053a57600080fd5b6001600160a01b03823516919081019060408101602082013564010000000081111561056557600080fd5b82018360208201111561057757600080fd5b8035906020019184600183028401116401000000008311171561059957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611af7945050505050565b604080519115158252519081900360200190f35b610361611b1f565b6105fe611b43565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610638578181015183820152602001610620565b50505050905090810190601f1680156106655780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102d5611bdf565b6102d5611cad565b61068b611d6c565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156106c75781810151838201526020016106af565b505050509050019250505060405180910390f35b6106e3611dce565b6040805163ffffffff94851681529290931660208301527fffffffffffffffffffffffffffffffff00000000000000000000000000000000168183015290519081900360600190f35b6102fd611def565b6102d56004803603602081101561074a57600080fd5b50356001600160a01b0316611e8b565b6102d56004803603602081101561077057600080fd5b50356001600160a01b0316611f82565b61033d611ff9565b6107b16004803603602081101561079e57600080fd5b503569ffffffffffffffffffff16612008565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102d56004803603604081101561081157600080fd5b81019060208101813564010000000081111561082c57600080fd5b82018360208201111561083e57600080fd5b8035906020019184602083028401116401000000008311171561086057600080fd5b91939092909160208101903564010000000081111561087e57600080fd5b82018360208201111561089057600080fd5b803590602001918460208302840111640100000000831117156108b257600080fd5b5090925090506120bd565b6102d5600480360360208110156108d357600080fd5b50356001600160a01b03166122f6565b6102d5600480360360208110156108f957600080fd5b50356001600160a01b031661235e565b6102fd6004803603602081101561091f57600080fd5b5035612457565b6102fd6004803603602081101561093c57600080fd5b50356124f4565b6102d5600480360360a081101561095957600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516612591565b6102d56004803603604081101561099e57600080fd5b506001600160a01b038135169060200135612730565b6102d5600480360360808110156109ca57600080fd5b8101906020810181356401000000008111156109e557600080fd5b8201836020820111156109f757600080fd5b80359060200191846001830284011164010000000083111715610a1957600080fd5b919390929091602081019035640100000000811115610a3757600080fd5b820183602082011115610a4957600080fd5b80359060200191846020830284011164010000000083111715610a6b57600080fd5b919390929091602081019035640100000000811115610a8957600080fd5b820183602082011115610a9b57600080fd5b80359060200191846020830284011164010000000083111715610abd57600080fd5b919350915035612a5f565b6102fd613950565b6105da613a01565b610afe60048036036020811015610aee57600080fd5b50356001600160a01b0316613a0a565b6040805161ffff9092168252519081900360200190f35b610b1d613ac3565b604080517fffffffffffffffffffffffffffffffff00000000000000000000000000000000909616865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b6102d560048036036040811015610b9557600080fd5b506001600160a01b0381358116916020013516613bb2565b6102d560048036036020811015610bc357600080fd5b50356001600160a01b0316613d0e565b6102d560048036036020811015610be957600080fd5b50356001600160a01b0316613dd6565b6107b1613e3e565b6000546001600160a01b03163314610c60576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602e5460ff1615610cbd57602e80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b6000610cc9615374565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115610d0d57fe5b6002811115610d1857fe5b9052509050600081602001516002811115610d2f57fe5b1415610d3f576000915050610e1a565b610d4761538b565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830181905270010000000000000000000000000000000090910490921660808201528251909160009160019060049060ff16601f8110610dd357fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff16601f8110610e1157fe5b01540301925050505b919050565b6000546001600160a01b03163314610e7e576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b036801000000000000000090910481169082168114610f1457602c80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff16680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000610f7061538b565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b7f000000000000000000000000000000000000000000000000000000000000000081565b602c546801000000000000000090046001600160a01b031690565b6000611077336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611af792505050565b6110c8576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6110d0613ef1565b905090565b600481565b868560ff8616601f831115611136576040805162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b6000811161118b576040805162461bcd60e51b815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b8183146111c95760405162461bcd60e51b815260040180806020018281038252602481526020018061550d6024913960400191505060405180910390fd5b806003028311611220576040805162461bcd60e51b815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b6000546001600160a01b0316331461127f576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6028541561142357602880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810191600091839081106112bc57fe5b6000918252602082200154602980546001600160a01b03909216935090849081106112e357fe5b6000918252602090912001546001600160a01b0316905061130381613f2d565b6001600160a01b0380831660009081526027602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559284168252902080549091169055602880548061135f57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905560298054806113c257fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061127f915050565b60005b8a811015611831576000602760008e8e8581811061144057fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff16600281111561147657fe5b146114c8576040805162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260016020820152602760008e8e858181106114ef57fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561157a57fe5b02179055506000915060069050818c8c8581811061159457fe5b6001600160a01b0360209182029390930135831684528301939093526040909101600020541691909114159050611612576040805162461bcd60e51b815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b6000602760008c8c8581811061162457fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff16600281111561165a57fe5b146116ac576040805162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260026020820152602760008c8c858181106116d357fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561175e57fe5b021790555090505060288c8c8381811061177457fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03959092029390930135939093169290921790555060298a8a838181106117d657fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0396909302949094013594909416179091555001611426565b50602a805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602c80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009091161793849055909104811691166118fd30828f8f8f8f8f8f8f8f61415d565b602a60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602a60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600083820152604051601f909101601f19169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b6000611a9e336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611af792505050565b611aef576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6110d0614261565b6000611b038383614287565b80611b1657506001600160a01b03831632145b90505b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6060611b86336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611af792505050565b611bd7576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6110d06142b7565b6001546001600160a01b03163314611c3e576040805162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b03163314611d0c576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602e5460ff16610cbd57602e80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60606029805480602002602001604051908101604052809291908181526020018280548015611dc457602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611da6575b5050505050905090565b602c54602a5463ffffffff808316926401000000009004169060801b909192565b6000611e32336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611af792505050565b611e83576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6110d0614362565b6000546001600160a01b03163314611eea576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152602f602052604090205460ff1615611f7f576001600160a01b0381166000818152602f602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b6001600160a01b03818116600090815260066020526040902054163314611ff0576040805162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b611f7f81613f2d565b6000546001600160a01b031681565b6000806000806000612051336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611af792505050565b6120a2576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6120ab866143bd565b939a9299509097509550909350915050565b6000546001600160a01b0316331461211c576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b828114612170576040805162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b838110156122ef57600085858381811061218957fe5b905060200201356001600160a01b0316905060008484848181106121a957fe5b6001600160a01b0385811660009081526006602090815260409091205492029390930135831693509091169050801580806121f55750826001600160a01b0316826001600160a01b0316145b612246576040805162461bcd60e51b815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b03848116600090815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016858316908117909155908316146122df57826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506121739050565b5050505050565b6000546001600160a01b03163314612355576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b611f7f81614511565b6001600160a01b038181166000908152600760205260409020541633146123cc576040805162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b6001600160a01b0381811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600061249a336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611af792505050565b6124eb576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611b19826145aa565b6000612537336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611af792505050565b612588576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611b19826145e0565b6003546001600160a01b0316806125ef576040805162461bcd60e51b815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b6000546001600160a01b03163314806126c25750604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b03861694636b14daf8946000939190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b15801561269557600080fd5b505afa1580156126a9573d6000803e3d6000fd5b505050506040513d60208110156126bf57600080fd5b50515b612713576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b61271b614635565b6127288686868686614a11565b505050505050565b6000546001600160a01b031633148061280b5750600354604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b1580156127de57600080fd5b505afa1580156127f2573d6000803e3d6000fd5b505050506040513d602081101561280857600080fd5b50515b61285c576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6000612866614b8b565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b1580156128d757600080fd5b505afa1580156128eb573d6000803e3d6000fd5b505050506040513d602081101561290157600080fd5b505190508181101561295a576040805162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb8561299685850387614d6c565b6040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b1580156129dc57600080fd5b505af11580156129f0573d6000803e3d6000fd5b505050506040513d6020811015612a0657600080fd5b5051612a59576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a9050612a72888888888888614d83565b3614612ac5576040805162461bcd60e51b815260206004820152601960248201527f7472616e736d6974206d65737361676520746f6f206c6f6e6700000000000000604482015290519081900360640190fd5b612acd6153b9565b6040805160808082018352602a549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff1660208301527501000000000000000000000000000000000000000000810460ff169282019290925276010000000000000000000000000000000000000000000090910463ffffffff166060808301919091529082526000908a908a90811015612b8657600080fd5b813591602081013591810190606081016040820135640100000000811115612bad57600080fd5b820183602082011115612bbf57600080fd5b80359060200191846020830284011164010000000083111715612be157600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050505060408801525050506080840182905283515190925060589190911b907fffffffffffffffffffffffffffffffff00000000000000000000000000000000808316911614612ca8576040805162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b608083015183516020015164ffffffffff808316911610612d10576040805162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b83516040015160ff168911612d6c576040805162461bcd60e51b815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b601f891115612dc2576040805162461bcd60e51b815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b868914612e16576040805162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b601f8460400151511115612e71576040805162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b83600001516040015160020260ff1684604001515111612ed8576040805162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b8867ffffffffffffffff81118015612eef57600080fd5b506040519080825280601f01601f191660200182016040528015612f1a576020820181803683370190505b50606085015260005b60ff81168a1115612f8b57868160ff1660208110612f3d57fe5b1a60f81b85606001518260ff1681518110612f5457fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101612f23565b5083604001515167ffffffffffffffff81118015612fa857600080fd5b506040519080825280601f01601f191660200182016040528015612fd3576020820181803683370190505b506020850152612fe16153ed565b60005b8560400151518160ff1610156130e7576000858260ff166020811061300557fe5b1a90508281601f811061301457fe5b60200201511561306b576040805162461bcd60e51b815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b6001838260ff16601f811061307c57fe5b91151560209283029190910152869060ff841690811061309857fe5b1a60f81b87602001518360ff16815181106130af57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101612fe4565b506130f0615374565b336000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561312b57fe5b600281111561313657fe5b905250905060028160200151600281111561314d57fe5b14801561318157506029816000015160ff168154811061316957fe5b6000918252602090912001546001600160a01b031633145b6131d2576040805162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b5050835164ffffffffff90911660209091015250506040516000908a908a90808383808284376040519201829003909120945061321393506153ed92505050565b61321b615374565b60005b898110156134145760006001858760600151848151811061323b57fe5b60209101015160f81c601b018e8e8681811061325357fe5b905060200201358d8d8781811061326657fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156132c1573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526027602090815290849020838501909452835460ff8082168552929650929450840191610100900416600281111561331057fe5b600281111561331b57fe5b905250925060018360200151600281111561333257fe5b14613384576040805162461bcd60e51b815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff16601f811061339557fe5b6020020151156133ec576040805162461bcd60e51b815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16601f811061340157fe5b911515602090920201525060010161321e565b5050505060005b6001826040015151038110156134c55760008260400151826001018151811061344057fe5b602002602001015160170b8360400151838151811061345b57fe5b602002602001015160170b13159050806134bc576040805162461bcd60e51b815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b5060010161341b565b506040810151805160009190600281049081106134de57fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b1315801561354457507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b613595576040805162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602b82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519586523386890181905291860181905260a0988601898152845199870199909952835194909916997ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451998c999298949793969095909492939185019260c086019289820192909102908190849084905b838110156136f85781810151838201526020016136e0565b50505050905001838103825285818151815260200191508051906020019080838360005b8381101561373457818101518382015260200161371c565b50505050905090810190601f1680156137615780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390a281516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a36138168260000151606001518260170b614d9b565b5080518051602a8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff9283160217909155821061393757fe5b613945828260200151614eac565b505050505050505050565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b1580156139c057600080fd5b505afa1580156139d4573d6000803e3d6000fd5b505050506040513d60208110156139ea57600080fd5b5051905060006139f8614b8b565b90910391505090565b602e5460ff1681565b6000613a14615374565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115613a5857fe5b6002811115613a6357fe5b9052509050600081602001516002811115613a7a57fe5b1415613a8a576000915050610e1a565b60016004826000015160ff16601f8110613aa057fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b600080808080333214613b1d576040805162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602a5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602b6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6001600160a01b03828116600090815260066020526040902054163314613c20576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b336001600160a01b0382161415613c7e576040805162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b6001600160a01b03808316600090815260076020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114613d09576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b6000546001600160a01b03163314613d6d576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000546001600160a01b03163314613e35576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b611f7f81615107565b6000806000806000613e87336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611af792505050565b613ed8576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b613ee0615196565b945094509450945094509091929394565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b6020526040902054601790810b900b90565b613f35615374565b6001600160a01b0382166000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115613f7957fe5b6002811115613f8457fe5b90525090506000613f9483610cbf565b90508015613d09576001600160a01b0380841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b15801561403d57600080fd5b505af1158015614051573d6000803e3d6000fd5b505050506040513d602081101561406757600080fd5b50516140ba576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff16601f81106140d057fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16601f811061410b57fe5b0155604080516001600160a01b0380871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b60008a8a8a8a8a8a8a8a8a8a604051602001808b6001600160a01b031681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff1690565b6001600160a01b0382166000908152602f602052604081205460ff1680611b16575050602e5460ff161592915050565b602d8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff610100600188161502019095169490940493840181900481028201810190925282815260609390929091830182828015611dc45780601f1061433657610100808354040283529160200191611dc4565b820191906000526020600020905b81548152906001019060200180831161434457509395945050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e740000000000000000000000000000000000815250906144985760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561445d578181015183820152602001614445565b50505050905090810190601f16801561448a5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b506144a1615374565b5050505063ffffffff83166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b6001600160a01b0381166000908152602f602052604090205460ff16611f7f576001600160a01b0381166000818152602f602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a150565b600063ffffffff8211156145c057506000610e1a565b5063ffffffff166000908152602b6020526040902054601790810b900b90565b600063ffffffff8211156145f657506000610e1a565b5063ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b61463d61538b565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830152700100000000000000000000000000000000900490911660808201526146b46153ed565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116146cd579050505050505090506147146153ed565b604080516103e081019182905290600890601f9082845b81548152602001906001019080831161472b57505050505090506060602980548060200260200160405190810160405280929190818152602001828054801561479d57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161477f575b5050505050905060005b81518110156149f557600060018483601f81106147c057fe5b6020020151039050600060018684601f81106147d857fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca000201905060008111156149ea5760006006600087878151811061481857fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b031690507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb82846040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b1580156148cd57600080fd5b505af11580156148e1573d6000803e3d6000fd5b505050506040513d60208110156148f757600080fd5b505161494a576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018886601f811061495857fe5b61ffff909216602092909202015260018786601f811061497457fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2908790879081106149a957fe5b6020026020010151828460405180846001600160a01b03168152602001836001600160a01b03168152602001828152602001935050505060405180910390a1505b5050506001016147a7565b50614a03600484601f61540c565b506122ef600883601f6154a2565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000614b956153ed565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614bae5790505050505050905060005b601f811015614c1e5760018282601f8110614c0757fe5b60200201510361ffff169290920191600101614bf0565b50614c2761538b565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606080860182905270010000000000000000000000000000000090940490921660808501526029805486518184028101840190975280875297909202633b9aca0002969394929390830182828015614cf757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311614cd9575b50505050509050614d066153ed565b604080516103e081019182905290600890601f9082845b815481526020019060010190808311614d1d575050505050905060005b8251811015614d645760018282601f8110614d5157fe5b6020020151039590950194600101614d3a565b505050505090565b600081831015614d7d575081611b19565b50919050565b602083810286019082020160e4019695505050505050565b602c546801000000000000000090046001600160a01b031680614dbe5750610f14565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602b602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b60248501819052948916604485015260648401889052516001600160a01b0387169363beed9b5193620186a09360848084019491939192918390030190829088803b158015614e7d57600080fd5b5087f193505050508015614ea357506040513d6020811015614e9e57600080fd5b505160015b612728576122ef565b614eb4615374565b336000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115614eef57fe5b6002811115614efa57fe5b9052509050614f0761538b565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000082048116606084015270010000000000000000000000000000000090910416608082015281516103e08101928390529091614fd2918591600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614f905790505050505050615235565b614fe090600490601f61540c565b50600282602001516002811115614ff357fe5b14615045576040805162461bcd60e51b815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b600061506c633b9aca003a04836020015163ffffffff16846000015163ffffffff166152aa565b90506010360260005a9050600061508b8863ffffffff168585856152d0565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff168302816150b957fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff16601f81106150e257fe5b015401016008886000015160ff16601f81106150fa57fe5b0155505050505050505050565b6003546001600160a01b039081169082168114610f1457600380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b602a54760100000000000000000000000000000000000000000000900463ffffffff1660008080806151c6615374565b5050505063ffffffff82166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b61523d6153ed565b60005b83518110156152a257600084828151811061525757fe5b016020015160f81c905061527c8482601f811061527057fe5b6020020151600161535c565b848260ff16601f811061528b57fe5b61ffff909216602092909202015250600101615240565b509092915050565b600083838110156152bd57600285850304015b6152c78184614d6c565b95945050505050565b600081851015615327576040805162461bcd60e51b815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b818503830161179301633b9aca00858202026fffffffffffffffffffffffffffffffff811061535257fe5b9695505050505050565b6000611b168261ffff168461ffff160161ffff614d6c565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518060a001604052806153cc6154d0565b81526060602082018190526040820181905280820152600060809091015290565b604051806103e00160405280601f906020820280368337509192915050565b6002830191839082156154925791602002820160005b8382111561546257835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302615422565b80156154905782816101000a81549061ffff0219169055600201602081600101049283019260010302615462565b505b5061549e9291506154f7565b5090565b82601f8101928215615492579160200282015b828111156154925782518255916020019190600101906154b5565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b8082111561549e57600081556001016154f856fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea164736f6c6343000705000a"


func DeployAccessControlledOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _validator common.Address, _minAnswer *big.Int, _maxAnswer *big.Int, _billingAdminAccessController common.Address, _decimals uint8, description string) (common.Address, *types.Transaction, *AccessControlledOffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlledOffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccessControlledOffchainAggregatorBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _validator, _minAnswer, _maxAnswer, _billingAdminAccessController, _decimals, description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessControlledOffchainAggregator{AccessControlledOffchainAggregatorCaller: AccessControlledOffchainAggregatorCaller{contract: contract}, AccessControlledOffchainAggregatorTransactor: AccessControlledOffchainAggregatorTransactor{contract: contract}, AccessControlledOffchainAggregatorFilterer: AccessControlledOffchainAggregatorFilterer{contract: contract}}, nil
}


type AccessControlledOffchainAggregator struct {
	AccessControlledOffchainAggregatorCaller     
	AccessControlledOffchainAggregatorTransactor 
	AccessControlledOffchainAggregatorFilterer   
}


type AccessControlledOffchainAggregatorCaller struct {
	contract *bind.BoundContract 
}


type AccessControlledOffchainAggregatorTransactor struct {
	contract *bind.BoundContract 
}


type AccessControlledOffchainAggregatorFilterer struct {
	contract *bind.BoundContract 
}



type AccessControlledOffchainAggregatorSession struct {
	Contract     *AccessControlledOffchainAggregator 
	CallOpts     bind.CallOpts                       
	TransactOpts bind.TransactOpts                   
}



type AccessControlledOffchainAggregatorCallerSession struct {
	Contract *AccessControlledOffchainAggregatorCaller 
	CallOpts bind.CallOpts                             
}



type AccessControlledOffchainAggregatorTransactorSession struct {
	Contract     *AccessControlledOffchainAggregatorTransactor 
	TransactOpts bind.TransactOpts                             
}


type AccessControlledOffchainAggregatorRaw struct {
	Contract *AccessControlledOffchainAggregator 
}


type AccessControlledOffchainAggregatorCallerRaw struct {
	Contract *AccessControlledOffchainAggregatorCaller 
}


type AccessControlledOffchainAggregatorTransactorRaw struct {
	Contract *AccessControlledOffchainAggregatorTransactor 
}


func NewAccessControlledOffchainAggregator(address common.Address, backend bind.ContractBackend) (*AccessControlledOffchainAggregator, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregator{AccessControlledOffchainAggregatorCaller: AccessControlledOffchainAggregatorCaller{contract: contract}, AccessControlledOffchainAggregatorTransactor: AccessControlledOffchainAggregatorTransactor{contract: contract}, AccessControlledOffchainAggregatorFilterer: AccessControlledOffchainAggregatorFilterer{contract: contract}}, nil
}


func NewAccessControlledOffchainAggregatorCaller(address common.Address, caller bind.ContractCaller) (*AccessControlledOffchainAggregatorCaller, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCaller{contract: contract}, nil
}


func NewAccessControlledOffchainAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControlledOffchainAggregatorTransactor, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorTransactor{contract: contract}, nil
}


func NewAccessControlledOffchainAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControlledOffchainAggregatorFilterer, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorFilterer{contract: contract}, nil
}


func bindAccessControlledOffchainAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlledOffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}



func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorTransactor.contract.Transfer(opts)
}


func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}





func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlledOffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}



func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.contract.Transfer(opts)
}


func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.contract.Transact(opts, method, params...)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LINK() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.LINK(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.LINK(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) CheckEnabled() (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.CheckEnabled(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.CheckEnabled(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Decimals() (uint8, error) {
	return _AccessControlledOffchainAggregator.Contract.Decimals(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _AccessControlledOffchainAggregator.Contract.Decimals(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Description() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.Description(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Description() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.Description(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetAnswer(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetAnswer(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetBilling(opts *bind.CallOpts) (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(struct {
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPrice = out[0].(uint32)
	outstruct.ReasonableGasPrice = out[1].(uint32)
	outstruct.MicroLinkPerEth = out[2].(uint32)
	outstruct.LinkGweiPerObservation = out[3].(uint32)
	outstruct.LinkGweiPerTransmission = out[4].(uint32)

	return *outstruct, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _AccessControlledOffchainAggregator.Contract.GetBilling(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _AccessControlledOffchainAggregator.Contract.GetBilling(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledOffchainAggregator.Contract.GetRoundData(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledOffchainAggregator.Contract.GetRoundData(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetTimestamp(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetTimestamp(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.HasAccess(&_AccessControlledOffchainAggregator.CallOpts, _user, _calldata)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.HasAccess(&_AccessControlledOffchainAggregator.CallOpts, _user, _calldata)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [16]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = out[0].(uint32)
	outstruct.BlockNumber = out[1].(uint32)
	outstruct.ConfigDigest = out[2].([16]byte)

	return *outstruct, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestConfigDetails(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestConfigDetails(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestRound() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRound(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRound(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRoundData(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRoundData(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTimestamp(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTimestamp(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestTransmissionDetails")

	outstruct := new(struct {
		ConfigDigest    [16]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigDigest = out[0].([16]byte)
	outstruct.Epoch = out[1].(uint32)
	outstruct.Round = out[2].(uint8)
	outstruct.LatestAnswer = out[3].(*big.Int)
	outstruct.LatestTimestamp = out[4].(uint64)

	return *outstruct, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTransmissionDetails(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTransmissionDetails(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LinkAvailableForPayment(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LinkAvailableForPayment(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MaxAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MaxAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MinAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MinAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _AccessControlledOffchainAggregator.Contract.OracleObservationCount(&_AccessControlledOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _AccessControlledOffchainAggregator.Contract.OracleObservationCount(&_AccessControlledOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.OwedPayment(&_AccessControlledOffchainAggregator.CallOpts, _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.OwedPayment(&_AccessControlledOffchainAggregator.CallOpts, _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Owner() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Owner(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Owner(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmitters(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmitters(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "validator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Validator() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Validator(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Validator(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Version() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.Version(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Version() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.Version(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "acceptOwnership")
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptOwnership(&_AccessControlledOffchainAggregator.TransactOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptOwnership(&_AccessControlledOffchainAggregator.TransactOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "acceptPayeeship", _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "addAccess", _user)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AddAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AddAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "disableAccessCheck")
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.DisableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.DisableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "enableAccessCheck")
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.EnableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.EnableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "removeAccess", _user)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RemoveAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RemoveAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBilling(&_AccessControlledOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBilling(&_AccessControlledOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setBillingAccessController", _billingAdminAccessController)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBillingAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _billingAdminAccessController)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBillingAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _billingAdminAccessController)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetConfig(&_AccessControlledOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetConfig(&_AccessControlledOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setPayees", _transmitters, _payees)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetPayees(&_AccessControlledOffchainAggregator.TransactOpts, _transmitters, _payees)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetPayees(&_AccessControlledOffchainAggregator.TransactOpts, _transmitters, _payees)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetValidator(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setValidator", _newValidator)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetValidator(&_AccessControlledOffchainAggregator.TransactOpts, _newValidator)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetValidator(&_AccessControlledOffchainAggregator.TransactOpts, _newValidator)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transferOwnership", _to)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferOwnership(&_AccessControlledOffchainAggregator.TransactOpts, _to)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferOwnership(&_AccessControlledOffchainAggregator.TransactOpts, _to)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter, _proposed)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter, _proposed)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) Transmit(opts *bind.TransactOpts, _report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transmit", _report, _rs, _ss, _rawVs)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmit(&_AccessControlledOffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmit(&_AccessControlledOffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawFunds(&_AccessControlledOffchainAggregator.TransactOpts, _recipient, _amount)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawFunds(&_AccessControlledOffchainAggregator.TransactOpts, _recipient, _amount)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "withdrawPayment", _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawPayment(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawPayment(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}


type AccessControlledOffchainAggregatorAddedAccessIterator struct {
	Event *AccessControlledOffchainAggregatorAddedAccess 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorAddedAccess)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorAddedAccess)
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


func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorAddedAccess struct {
	User common.Address
	Raw  types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorAddedAccessIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorAddedAccessIterator{contract: _AccessControlledOffchainAggregator.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorAddedAccess) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorAddedAccess)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseAddedAccess(log types.Log) (*AccessControlledOffchainAggregatorAddedAccess, error) {
	event := new(AccessControlledOffchainAggregatorAddedAccess)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorAnswerUpdatedIterator struct {
	Event *AccessControlledOffchainAggregatorAnswerUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorAnswerUpdated)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorAnswerUpdated)
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


func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AccessControlledOffchainAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorAnswerUpdatedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorAnswerUpdated)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*AccessControlledOffchainAggregatorAnswerUpdated, error) {
	event := new(AccessControlledOffchainAggregatorAnswerUpdated)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorBillingAccessControllerSetIterator struct {
	Event *AccessControlledOffchainAggregatorBillingAccessControllerSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
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


func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorBillingAccessControllerSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*AccessControlledOffchainAggregatorBillingAccessControllerSet, error) {
	event := new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorBillingSetIterator struct {
	Event *AccessControlledOffchainAggregatorBillingSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorBillingSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorBillingSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorBillingSet)
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


func (it *AccessControlledOffchainAggregatorBillingSetIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorBillingSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorBillingSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorBillingSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseBillingSet(log types.Log) (*AccessControlledOffchainAggregatorBillingSet, error) {
	event := new(AccessControlledOffchainAggregatorBillingSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorCheckAccessDisabledIterator struct {
	Event *AccessControlledOffchainAggregatorCheckAccessDisabled 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorCheckAccessDisabled)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorCheckAccessDisabled)
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


func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorCheckAccessDisabled struct {
	Raw types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorCheckAccessDisabledIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCheckAccessDisabledIterator{contract: _AccessControlledOffchainAggregator.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorCheckAccessDisabled)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseCheckAccessDisabled(log types.Log) (*AccessControlledOffchainAggregatorCheckAccessDisabled, error) {
	event := new(AccessControlledOffchainAggregatorCheckAccessDisabled)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorCheckAccessEnabledIterator struct {
	Event *AccessControlledOffchainAggregatorCheckAccessEnabled 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorCheckAccessEnabled)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorCheckAccessEnabled)
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


func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorCheckAccessEnabled struct {
	Raw types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorCheckAccessEnabledIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCheckAccessEnabledIterator{contract: _AccessControlledOffchainAggregator.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorCheckAccessEnabled)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseCheckAccessEnabled(log types.Log) (*AccessControlledOffchainAggregatorCheckAccessEnabled, error) {
	event := new(AccessControlledOffchainAggregatorCheckAccessEnabled)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorConfigSetIterator struct {
	Event *AccessControlledOffchainAggregatorConfigSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorConfigSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorConfigSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorConfigSet)
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


func (it *AccessControlledOffchainAggregatorConfigSetIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorConfigSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorConfigSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorConfigSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseConfigSet(log types.Log) (*AccessControlledOffchainAggregatorConfigSet, error) {
	event := new(AccessControlledOffchainAggregatorConfigSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorNewRoundIterator struct {
	Event *AccessControlledOffchainAggregatorNewRound 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorNewRoundIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorNewRound)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorNewRound)
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


func (it *AccessControlledOffchainAggregatorNewRoundIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AccessControlledOffchainAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorNewRoundIterator{contract: _AccessControlledOffchainAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorNewRound)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseNewRound(log types.Log) (*AccessControlledOffchainAggregatorNewRound, error) {
	event := new(AccessControlledOffchainAggregatorNewRound)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorNewTransmissionIterator struct {
	Event *AccessControlledOffchainAggregatorNewTransmission 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorNewTransmission)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorNewTransmission)
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


func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	RawReportContext  [32]byte
	Raw               types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*AccessControlledOffchainAggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorNewTransmissionIterator{contract: _AccessControlledOffchainAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorNewTransmission)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseNewTransmission(log types.Log) (*AccessControlledOffchainAggregatorNewTransmission, error) {
	event := new(AccessControlledOffchainAggregatorNewTransmission)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorOraclePaidIterator struct {
	Event *AccessControlledOffchainAggregatorOraclePaid 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorOraclePaid)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorOraclePaid)
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


func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	Raw         types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorOraclePaidIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorOraclePaidIterator{contract: _AccessControlledOffchainAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOraclePaid) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorOraclePaid)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOraclePaid(log types.Log) (*AccessControlledOffchainAggregatorOraclePaid, error) {
	event := new(AccessControlledOffchainAggregatorOraclePaid)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator struct {
	Event *AccessControlledOffchainAggregatorOwnershipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
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


func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*AccessControlledOffchainAggregatorOwnershipTransferRequested, error) {
	event := new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorOwnershipTransferredIterator struct {
	Event *AccessControlledOffchainAggregatorOwnershipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorOwnershipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorOwnershipTransferred)
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


func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AccessControlledOffchainAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorOwnershipTransferredIterator{contract: _AccessControlledOffchainAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorOwnershipTransferred)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*AccessControlledOffchainAggregatorOwnershipTransferred, error) {
	event := new(AccessControlledOffchainAggregatorOwnershipTransferred)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator struct {
	Event *AccessControlledOffchainAggregatorPayeeshipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
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


func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*AccessControlledOffchainAggregatorPayeeshipTransferRequested, error) {
	event := new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorPayeeshipTransferredIterator struct {
	Event *AccessControlledOffchainAggregatorPayeeshipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorPayeeshipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorPayeeshipTransferred)
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


func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*AccessControlledOffchainAggregatorPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorPayeeshipTransferredIterator{contract: _AccessControlledOffchainAggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorPayeeshipTransferred)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*AccessControlledOffchainAggregatorPayeeshipTransferred, error) {
	event := new(AccessControlledOffchainAggregatorPayeeshipTransferred)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorRemovedAccessIterator struct {
	Event *AccessControlledOffchainAggregatorRemovedAccess 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorRemovedAccess)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorRemovedAccess)
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


func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorRemovedAccess struct {
	User common.Address
	Raw  types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorRemovedAccessIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorRemovedAccessIterator{contract: _AccessControlledOffchainAggregator.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorRemovedAccess)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseRemovedAccess(log types.Log) (*AccessControlledOffchainAggregatorRemovedAccess, error) {
	event := new(AccessControlledOffchainAggregatorRemovedAccess)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AccessControlledOffchainAggregatorValidatorUpdatedIterator struct {
	Event *AccessControlledOffchainAggregatorValidatorUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AccessControlledOffchainAggregatorValidatorUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorValidatorUpdated)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorValidatorUpdated)
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


func (it *AccessControlledOffchainAggregatorValidatorUpdatedIterator) Error() error {
	return it.fail
}



func (it *AccessControlledOffchainAggregatorValidatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AccessControlledOffchainAggregatorValidatorUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log 
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterValidatorUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*AccessControlledOffchainAggregatorValidatorUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorValidatorUpdatedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "ValidatorUpdated", logs: logs, sub: sub}, nil
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchValidatorUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorValidatorUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AccessControlledOffchainAggregatorValidatorUpdated)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
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




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseValidatorUpdated(log types.Log) (*AccessControlledOffchainAggregatorValidatorUpdated, error) {
	event := new(AccessControlledOffchainAggregatorValidatorUpdated)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const AccessControllerInterfaceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"


type AccessControllerInterface struct {
	AccessControllerInterfaceCaller     
	AccessControllerInterfaceTransactor 
	AccessControllerInterfaceFilterer   
}


type AccessControllerInterfaceCaller struct {
	contract *bind.BoundContract 
}


type AccessControllerInterfaceTransactor struct {
	contract *bind.BoundContract 
}


type AccessControllerInterfaceFilterer struct {
	contract *bind.BoundContract 
}



type AccessControllerInterfaceSession struct {
	Contract     *AccessControllerInterface 
	CallOpts     bind.CallOpts              
	TransactOpts bind.TransactOpts          
}



type AccessControllerInterfaceCallerSession struct {
	Contract *AccessControllerInterfaceCaller 
	CallOpts bind.CallOpts                    
}



type AccessControllerInterfaceTransactorSession struct {
	Contract     *AccessControllerInterfaceTransactor 
	TransactOpts bind.TransactOpts                    
}


type AccessControllerInterfaceRaw struct {
	Contract *AccessControllerInterface 
}


type AccessControllerInterfaceCallerRaw struct {
	Contract *AccessControllerInterfaceCaller 
}


type AccessControllerInterfaceTransactorRaw struct {
	Contract *AccessControllerInterfaceTransactor 
}


func NewAccessControllerInterface(address common.Address, backend bind.ContractBackend) (*AccessControllerInterface, error) {
	contract, err := bindAccessControllerInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterface{AccessControllerInterfaceCaller: AccessControllerInterfaceCaller{contract: contract}, AccessControllerInterfaceTransactor: AccessControllerInterfaceTransactor{contract: contract}, AccessControllerInterfaceFilterer: AccessControllerInterfaceFilterer{contract: contract}}, nil
}


func NewAccessControllerInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AccessControllerInterfaceCaller, error) {
	contract, err := bindAccessControllerInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceCaller{contract: contract}, nil
}


func NewAccessControllerInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControllerInterfaceTransactor, error) {
	contract, err := bindAccessControllerInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceTransactor{contract: contract}, nil
}


func NewAccessControllerInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControllerInterfaceFilterer, error) {
	contract, err := bindAccessControllerInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceFilterer{contract: contract}, nil
}


func bindAccessControllerInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControllerInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_AccessControllerInterface *AccessControllerInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AccessControllerInterface *AccessControllerInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceTransactor.contract.Transfer(opts)
}


func (_AccessControllerInterface *AccessControllerInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AccessControllerInterface *AccessControllerInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControllerInterface.Contract.contract.Call(opts, result, method, params...)
}



func (_AccessControllerInterface *AccessControllerInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.contract.Transfer(opts)
}


func (_AccessControllerInterface *AccessControllerInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.contract.Transact(opts, method, params...)
}




func (_AccessControllerInterface *AccessControllerInterfaceCaller) HasAccess(opts *bind.CallOpts, user common.Address, data []byte) (bool, error) {
	var out []interface{}
	err := _AccessControllerInterface.contract.Call(opts, &out, "hasAccess", user, data)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_AccessControllerInterface *AccessControllerInterfaceSession) HasAccess(user common.Address, data []byte) (bool, error) {
	return _AccessControllerInterface.Contract.HasAccess(&_AccessControllerInterface.CallOpts, user, data)
}




func (_AccessControllerInterface *AccessControllerInterfaceCallerSession) HasAccess(user common.Address, data []byte) (bool, error) {
	return _AccessControllerInterface.Contract.HasAccess(&_AccessControllerInterface.CallOpts, user, data)
}


const AggregatorInterfaceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"


type AggregatorInterface struct {
	AggregatorInterfaceCaller     
	AggregatorInterfaceTransactor 
	AggregatorInterfaceFilterer   
}


type AggregatorInterfaceCaller struct {
	contract *bind.BoundContract 
}


type AggregatorInterfaceTransactor struct {
	contract *bind.BoundContract 
}


type AggregatorInterfaceFilterer struct {
	contract *bind.BoundContract 
}



type AggregatorInterfaceSession struct {
	Contract     *AggregatorInterface 
	CallOpts     bind.CallOpts        
	TransactOpts bind.TransactOpts    
}



type AggregatorInterfaceCallerSession struct {
	Contract *AggregatorInterfaceCaller 
	CallOpts bind.CallOpts              
}



type AggregatorInterfaceTransactorSession struct {
	Contract     *AggregatorInterfaceTransactor 
	TransactOpts bind.TransactOpts              
}


type AggregatorInterfaceRaw struct {
	Contract *AggregatorInterface 
}


type AggregatorInterfaceCallerRaw struct {
	Contract *AggregatorInterfaceCaller 
}


type AggregatorInterfaceTransactorRaw struct {
	Contract *AggregatorInterfaceTransactor 
}


func NewAggregatorInterface(address common.Address, backend bind.ContractBackend) (*AggregatorInterface, error) {
	contract, err := bindAggregatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterface{AggregatorInterfaceCaller: AggregatorInterfaceCaller{contract: contract}, AggregatorInterfaceTransactor: AggregatorInterfaceTransactor{contract: contract}, AggregatorInterfaceFilterer: AggregatorInterfaceFilterer{contract: contract}}, nil
}


func NewAggregatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorInterfaceCaller, error) {
	contract, err := bindAggregatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceCaller{contract: contract}, nil
}


func NewAggregatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorInterfaceTransactor, error) {
	contract, err := bindAggregatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceTransactor{contract: contract}, nil
}


func NewAggregatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorInterfaceFilterer, error) {
	contract, err := bindAggregatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceFilterer{contract: contract}, nil
}


func bindAggregatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_AggregatorInterface *AggregatorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorInterface.Contract.AggregatorInterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AggregatorInterface *AggregatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.AggregatorInterfaceTransactor.contract.Transfer(opts)
}


func (_AggregatorInterface *AggregatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.AggregatorInterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AggregatorInterface *AggregatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorInterface.Contract.contract.Call(opts, result, method, params...)
}



func (_AggregatorInterface *AggregatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.contract.Transfer(opts)
}


func (_AggregatorInterface *AggregatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.contract.Transact(opts, method, params...)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorInterface *AggregatorInterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorInterface *AggregatorInterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorInterface *AggregatorInterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorInterface *AggregatorInterfaceSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorInterface *AggregatorInterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestTimestamp(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestTimestamp(&_AggregatorInterface.CallOpts)
}


type AggregatorInterfaceAnswerUpdatedIterator struct {
	Event *AggregatorInterfaceAnswerUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AggregatorInterfaceAnswerUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorInterfaceAnswerUpdated)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AggregatorInterfaceAnswerUpdated)
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


func (it *AggregatorInterfaceAnswerUpdatedIterator) Error() error {
	return it.fail
}



func (it *AggregatorInterfaceAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AggregatorInterfaceAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log 
}




func (_AggregatorInterface *AggregatorInterfaceFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AggregatorInterfaceAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorInterface.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceAnswerUpdatedIterator{contract: _AggregatorInterface.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}




func (_AggregatorInterface *AggregatorInterfaceFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AggregatorInterfaceAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorInterface.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AggregatorInterfaceAnswerUpdated)
				if err := _AggregatorInterface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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




func (_AggregatorInterface *AggregatorInterfaceFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorInterfaceAnswerUpdated, error) {
	event := new(AggregatorInterfaceAnswerUpdated)
	if err := _AggregatorInterface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AggregatorInterfaceNewRoundIterator struct {
	Event *AggregatorInterfaceNewRound 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AggregatorInterfaceNewRoundIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorInterfaceNewRound)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AggregatorInterfaceNewRound)
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


func (it *AggregatorInterfaceNewRoundIterator) Error() error {
	return it.fail
}



func (it *AggregatorInterfaceNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AggregatorInterfaceNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log 
}




func (_AggregatorInterface *AggregatorInterfaceFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AggregatorInterfaceNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorInterface.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceNewRoundIterator{contract: _AggregatorInterface.contract, event: "NewRound", logs: logs, sub: sub}, nil
}




func (_AggregatorInterface *AggregatorInterfaceFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AggregatorInterfaceNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorInterface.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AggregatorInterfaceNewRound)
				if err := _AggregatorInterface.contract.UnpackLog(event, "NewRound", log); err != nil {
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




func (_AggregatorInterface *AggregatorInterfaceFilterer) ParseNewRound(log types.Log) (*AggregatorInterfaceNewRound, error) {
	event := new(AggregatorInterfaceNewRound)
	if err := _AggregatorInterface.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const AggregatorV2V3InterfaceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"


type AggregatorV2V3Interface struct {
	AggregatorV2V3InterfaceCaller     
	AggregatorV2V3InterfaceTransactor 
	AggregatorV2V3InterfaceFilterer   
}


type AggregatorV2V3InterfaceCaller struct {
	contract *bind.BoundContract 
}


type AggregatorV2V3InterfaceTransactor struct {
	contract *bind.BoundContract 
}


type AggregatorV2V3InterfaceFilterer struct {
	contract *bind.BoundContract 
}



type AggregatorV2V3InterfaceSession struct {
	Contract     *AggregatorV2V3Interface 
	CallOpts     bind.CallOpts            
	TransactOpts bind.TransactOpts        
}



type AggregatorV2V3InterfaceCallerSession struct {
	Contract *AggregatorV2V3InterfaceCaller 
	CallOpts bind.CallOpts                  
}



type AggregatorV2V3InterfaceTransactorSession struct {
	Contract     *AggregatorV2V3InterfaceTransactor 
	TransactOpts bind.TransactOpts                  
}


type AggregatorV2V3InterfaceRaw struct {
	Contract *AggregatorV2V3Interface 
}


type AggregatorV2V3InterfaceCallerRaw struct {
	Contract *AggregatorV2V3InterfaceCaller 
}


type AggregatorV2V3InterfaceTransactorRaw struct {
	Contract *AggregatorV2V3InterfaceTransactor 
}


func NewAggregatorV2V3Interface(address common.Address, backend bind.ContractBackend) (*AggregatorV2V3Interface, error) {
	contract, err := bindAggregatorV2V3Interface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3Interface{AggregatorV2V3InterfaceCaller: AggregatorV2V3InterfaceCaller{contract: contract}, AggregatorV2V3InterfaceTransactor: AggregatorV2V3InterfaceTransactor{contract: contract}, AggregatorV2V3InterfaceFilterer: AggregatorV2V3InterfaceFilterer{contract: contract}}, nil
}


func NewAggregatorV2V3InterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorV2V3InterfaceCaller, error) {
	contract, err := bindAggregatorV2V3Interface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceCaller{contract: contract}, nil
}


func NewAggregatorV2V3InterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorV2V3InterfaceTransactor, error) {
	contract, err := bindAggregatorV2V3Interface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceTransactor{contract: contract}, nil
}


func NewAggregatorV2V3InterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorV2V3InterfaceFilterer, error) {
	contract, err := bindAggregatorV2V3Interface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceFilterer{contract: contract}, nil
}


func bindAggregatorV2V3Interface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorV2V3InterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceTransactor.contract.Transfer(opts)
}


func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV2V3Interface.Contract.contract.Call(opts, result, method, params...)
}



func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.contract.Transfer(opts)
}


func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.contract.Transact(opts, method, params...)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Decimals() (uint8, error) {
	return _AggregatorV2V3Interface.Contract.Decimals(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorV2V3Interface.Contract.Decimals(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Description() (string, error) {
	return _AggregatorV2V3Interface.Contract.Description(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Description() (string, error) {
	return _AggregatorV2V3Interface.Contract.Description(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetAnswer(&_AggregatorV2V3Interface.CallOpts, roundId)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetAnswer(&_AggregatorV2V3Interface.CallOpts, roundId)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.GetRoundData(&_AggregatorV2V3Interface.CallOpts, _roundId)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.GetRoundData(&_AggregatorV2V3Interface.CallOpts, _roundId)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetTimestamp(&_AggregatorV2V3Interface.CallOpts, roundId)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetTimestamp(&_AggregatorV2V3Interface.CallOpts, roundId)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestAnswer(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestAnswer(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestRound() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestRound(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestRound(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.LatestRoundData(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.LatestRoundData(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestTimestamp(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestTimestamp(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Version() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.Version(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Version() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.Version(&_AggregatorV2V3Interface.CallOpts)
}


type AggregatorV2V3InterfaceAnswerUpdatedIterator struct {
	Event *AggregatorV2V3InterfaceAnswerUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorV2V3InterfaceAnswerUpdated)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AggregatorV2V3InterfaceAnswerUpdated)
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


func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Error() error {
	return it.fail
}



func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AggregatorV2V3InterfaceAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log 
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AggregatorV2V3InterfaceAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorV2V3Interface.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceAnswerUpdatedIterator{contract: _AggregatorV2V3Interface.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AggregatorV2V3InterfaceAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AggregatorV2V3Interface.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AggregatorV2V3InterfaceAnswerUpdated)
				if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorV2V3InterfaceAnswerUpdated, error) {
	event := new(AggregatorV2V3InterfaceAnswerUpdated)
	if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type AggregatorV2V3InterfaceNewRoundIterator struct {
	Event *AggregatorV2V3InterfaceNewRound 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *AggregatorV2V3InterfaceNewRoundIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AggregatorV2V3InterfaceNewRound)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(AggregatorV2V3InterfaceNewRound)
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


func (it *AggregatorV2V3InterfaceNewRoundIterator) Error() error {
	return it.fail
}



func (it *AggregatorV2V3InterfaceNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type AggregatorV2V3InterfaceNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log 
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AggregatorV2V3InterfaceNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorV2V3Interface.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceNewRoundIterator{contract: _AggregatorV2V3Interface.contract, event: "NewRound", logs: logs, sub: sub}, nil
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AggregatorV2V3InterfaceNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AggregatorV2V3Interface.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(AggregatorV2V3InterfaceNewRound)
				if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "NewRound", log); err != nil {
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




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) ParseNewRound(log types.Log) (*AggregatorV2V3InterfaceNewRound, error) {
	event := new(AggregatorV2V3InterfaceNewRound)
	if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const AggregatorV3InterfaceABI = "[{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"


type AggregatorV3Interface struct {
	AggregatorV3InterfaceCaller     
	AggregatorV3InterfaceTransactor 
	AggregatorV3InterfaceFilterer   
}


type AggregatorV3InterfaceCaller struct {
	contract *bind.BoundContract 
}


type AggregatorV3InterfaceTransactor struct {
	contract *bind.BoundContract 
}


type AggregatorV3InterfaceFilterer struct {
	contract *bind.BoundContract 
}



type AggregatorV3InterfaceSession struct {
	Contract     *AggregatorV3Interface 
	CallOpts     bind.CallOpts          
	TransactOpts bind.TransactOpts      
}



type AggregatorV3InterfaceCallerSession struct {
	Contract *AggregatorV3InterfaceCaller 
	CallOpts bind.CallOpts                
}



type AggregatorV3InterfaceTransactorSession struct {
	Contract     *AggregatorV3InterfaceTransactor 
	TransactOpts bind.TransactOpts                
}


type AggregatorV3InterfaceRaw struct {
	Contract *AggregatorV3Interface 
}


type AggregatorV3InterfaceCallerRaw struct {
	Contract *AggregatorV3InterfaceCaller 
}


type AggregatorV3InterfaceTransactorRaw struct {
	Contract *AggregatorV3InterfaceTransactor 
}


func NewAggregatorV3Interface(address common.Address, backend bind.ContractBackend) (*AggregatorV3Interface, error) {
	contract, err := bindAggregatorV3Interface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3Interface{AggregatorV3InterfaceCaller: AggregatorV3InterfaceCaller{contract: contract}, AggregatorV3InterfaceTransactor: AggregatorV3InterfaceTransactor{contract: contract}, AggregatorV3InterfaceFilterer: AggregatorV3InterfaceFilterer{contract: contract}}, nil
}


func NewAggregatorV3InterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorV3InterfaceCaller, error) {
	contract, err := bindAggregatorV3Interface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceCaller{contract: contract}, nil
}


func NewAggregatorV3InterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorV3InterfaceTransactor, error) {
	contract, err := bindAggregatorV3Interface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceTransactor{contract: contract}, nil
}


func NewAggregatorV3InterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorV3InterfaceFilterer, error) {
	contract, err := bindAggregatorV3Interface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceFilterer{contract: contract}, nil
}


func bindAggregatorV3Interface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorV3InterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceTransactor.contract.Transfer(opts)
}


func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AggregatorV3Interface *AggregatorV3InterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV3Interface.Contract.contract.Call(opts, result, method, params...)
}



func (_AggregatorV3Interface *AggregatorV3InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.contract.Transfer(opts)
}


func (_AggregatorV3Interface *AggregatorV3InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.contract.Transact(opts, method, params...)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}




func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Decimals() (uint8, error) {
	return _AggregatorV3Interface.Contract.Decimals(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorV3Interface.Contract.Decimals(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}




func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Description() (string, error) {
	return _AggregatorV3Interface.Contract.Description(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Description() (string, error) {
	return _AggregatorV3Interface.Contract.Description(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_AggregatorV3Interface *AggregatorV3InterfaceSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.GetRoundData(&_AggregatorV3Interface.CallOpts, _roundId)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.GetRoundData(&_AggregatorV3Interface.CallOpts, _roundId)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_AggregatorV3Interface *AggregatorV3InterfaceSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.LatestRoundData(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.LatestRoundData(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Version() (*big.Int, error) {
	return _AggregatorV3Interface.Contract.Version(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Version() (*big.Int, error) {
	return _AggregatorV3Interface.Contract.Version(&_AggregatorV3Interface.CallOpts)
}


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


const LinkTokenInterfaceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"remaining\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimalPlaces\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenName\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalTokensIssued\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


type LinkTokenInterface struct {
	LinkTokenInterfaceCaller     
	LinkTokenInterfaceTransactor 
	LinkTokenInterfaceFilterer   
}


type LinkTokenInterfaceCaller struct {
	contract *bind.BoundContract 
}


type LinkTokenInterfaceTransactor struct {
	contract *bind.BoundContract 
}


type LinkTokenInterfaceFilterer struct {
	contract *bind.BoundContract 
}



type LinkTokenInterfaceSession struct {
	Contract     *LinkTokenInterface 
	CallOpts     bind.CallOpts       
	TransactOpts bind.TransactOpts   
}



type LinkTokenInterfaceCallerSession struct {
	Contract *LinkTokenInterfaceCaller 
	CallOpts bind.CallOpts             
}



type LinkTokenInterfaceTransactorSession struct {
	Contract     *LinkTokenInterfaceTransactor 
	TransactOpts bind.TransactOpts             
}


type LinkTokenInterfaceRaw struct {
	Contract *LinkTokenInterface 
}


type LinkTokenInterfaceCallerRaw struct {
	Contract *LinkTokenInterfaceCaller 
}


type LinkTokenInterfaceTransactorRaw struct {
	Contract *LinkTokenInterfaceTransactor 
}


func NewLinkTokenInterface(address common.Address, backend bind.ContractBackend) (*LinkTokenInterface, error) {
	contract, err := bindLinkTokenInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterface{LinkTokenInterfaceCaller: LinkTokenInterfaceCaller{contract: contract}, LinkTokenInterfaceTransactor: LinkTokenInterfaceTransactor{contract: contract}, LinkTokenInterfaceFilterer: LinkTokenInterfaceFilterer{contract: contract}}, nil
}


func NewLinkTokenInterfaceCaller(address common.Address, caller bind.ContractCaller) (*LinkTokenInterfaceCaller, error) {
	contract, err := bindLinkTokenInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceCaller{contract: contract}, nil
}


func NewLinkTokenInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkTokenInterfaceTransactor, error) {
	contract, err := bindLinkTokenInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceTransactor{contract: contract}, nil
}


func NewLinkTokenInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkTokenInterfaceFilterer, error) {
	contract, err := bindLinkTokenInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceFilterer{contract: contract}, nil
}


func bindLinkTokenInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LinkTokenInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_LinkTokenInterface *LinkTokenInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transfer(opts)
}


func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_LinkTokenInterface *LinkTokenInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.contract.Call(opts, result, method, params...)
}



func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transfer(opts)
}


func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transact(opts, method, params...)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_LinkTokenInterface *LinkTokenInterfaceSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_LinkTokenInterface *LinkTokenInterfaceSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "approve", spender, value)
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Approve(&_LinkTokenInterface.TransactOpts, spender, value)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Approve(&_LinkTokenInterface.TransactOpts, spender, value)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactor) DecreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "decreaseApproval", spender, addedValue)
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) DecreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.DecreaseApproval(&_LinkTokenInterface.TransactOpts, spender, addedValue)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) DecreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.DecreaseApproval(&_LinkTokenInterface.TransactOpts, spender, addedValue)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactor) IncreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "increaseApproval", spender, subtractedValue)
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) IncreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.IncreaseApproval(&_LinkTokenInterface.TransactOpts, spender, subtractedValue)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) IncreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.IncreaseApproval(&_LinkTokenInterface.TransactOpts, spender, subtractedValue)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transfer", to, value)
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Transfer(&_LinkTokenInterface.TransactOpts, to, value)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Transfer(&_LinkTokenInterface.TransactOpts, to, value)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transferAndCall", to, value, data)
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferAndCall(&_LinkTokenInterface.TransactOpts, to, value, data)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferAndCall(&_LinkTokenInterface.TransactOpts, to, value, data)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transferFrom", from, to, value)
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferFrom(&_LinkTokenInterface.TransactOpts, from, to, value)
}




func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferFrom(&_LinkTokenInterface.TransactOpts, from, to, value)
}


const OffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var OffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b50604051620051573803806200515783398181016040526101808110156200003957600080fd5b815160208301516040808501516060860151608087015160a088015160c089015160e08a01516101008b01516101208c01516101408d01516101608e0180519a519c9e9b9d999c989b979a969995989497939692959194939182019284640100000000821115620000a957600080fd5b908301906020820185811115620000bf57600080fd5b8251640100000000811182820188101715620000da57600080fd5b82525081516020918201929091019080838360005b8381101562000109578181015183820152602001620000ef565b50505050905090810190601f168015620001375780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b03191633179055508b8b8b8b8b8b8862000166878787878762000287565b620001718162000379565b6001600160601b0319606083901b166080526200018d620004d7565b62000197620004d7565b60005b601f8160ff161015620001e7576001838260ff16601f8110620001b957fe5b61ffff909216602092909202015260018260ff8316601f8110620001d957fe5b60200201526001016200019a565b50620001f7600483601f620004f6565b5062000207600882601f62000593565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e05250508351620002509350602d9250602085019150620005c4565b506200025c86620003f2565b505050601791820b820b604090811b60a05290820b90910b901b60c052506200065d95505050505050565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b039081169082168114620003ee57600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b0316331462000452576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b036801000000000000000090910481169082168114620003ee57602c8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b604051806103e00160405280601f906020820280368337509192915050565b600283019183908215620005815791602002820160005b838211156200054f57835183826101000a81548161ffff021916908361ffff16021790555092602001926002016020816001010492830192600103026200050d565b80156200057f5782816101000a81549061ffff02191690556002016020816001010492830192600103026200054f565b505b506200058f92915062000646565b5090565b82601f810192821562000581579160200282015b8281111562000581578251825591602001919060010190620005a7565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620005fc576000855562000581565b82601f106200061757805160ff191683800117855562000581565b8280016001018555821562000581579182018281111562000581578251825591602001919060010190620005a7565b5b808211156200058f576000815560010162000647565b60805160601c60a05160401c60c05160401c60e05160f81c614a95620006c260003980610d695250806117905280612f35525080610cb05280612f08525080610c8c52806122865280612376528061336f52806139b65280613e455250614a956000f3fe608060405234801561001057600080fd5b50600436106102265760003560e01c80638ac28d5a1161012a578063c1075329116100bd578063e5fe45771161008c578063f2fde38b11610071578063f2fde38b146109dd578063fbffd2c114610a03578063feaf968c14610a2957610226565b8063e5fe457714610945578063eb5dcd6c146109af57610226565b8063c1075329146107c0578063c9807539146107ec578063d09dc33914610900578063e4902f821461090857610226565b8063b121e147116100f9578063b121e1471461071b578063b5ab58dc14610741578063b633620c1461075e578063bd8247061461077b57610226565b80638ac28d5a146105b85780638da5cb5b146105de5780639a6fc8f5146105e65780639c849b301461065957610226565b806354fd4d50116101bd5780637284e4161161018c5780638141183411610171578063814118341461050757806381ff70481461055f5780638205bf6a146105b057610226565b80637284e4161461048257806379ba5097146104ff57610226565b806354fd4d501461033d578063585aa7de14610345578063668a0f021461047257806370da2f671461047a57610226565b806329937268116101f957806329937268146102ce578063313ce5671461030f5780633a5381b51461032d57806350d25bcd1461033557610226565b80630eafb25b1461022b5780631327d3d8146102635780631b6b6d231461028b57806322adbc78146102af575b600080fd5b6102516004803603602081101561024157600080fd5b50356001600160a01b0316610a31565b60408051918252519081900360200190f35b6102896004803603602081101561027957600080fd5b50356001600160a01b0316610b91565b005b610293610c8a565b604080516001600160a01b039092168252519081900360200190f35b6102b7610cae565b6040805160179290920b8252519081900360200190f35b6102d6610cd2565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b610317610d67565b6040805160ff9092168252519081900360200190f35b610293610d8b565b610251610da6565b610251610de2565b610289600480360360a081101561035b57600080fd5b81019060208101813564010000000081111561037657600080fd5b82018360208201111561038857600080fd5b803590602001918460208302840111640100000000831117156103aa57600080fd5b9193909290916020810190356401000000008111156103c857600080fd5b8201836020820111156103da57600080fd5b803590602001918460208302840111640100000000831117156103fc57600080fd5b9193909260ff8335169267ffffffffffffffff60208201351692919060608101906040013564010000000081111561043357600080fd5b82018360208201111561044557600080fd5b8035906020019184600183028401116401000000008311171561046757600080fd5b509092509050610de7565b610251611768565b6102b761178e565b61048a6117b2565b6040805160208082528351818301528351919283929083019185019080838360005b838110156104c45781810151838201526020016104ac565b50505050905090810190601f1680156104f15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610289611866565b61050f611934565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561054b578181015183820152602001610533565b505050509050019250505060405180910390f35b610567611995565b6040805163ffffffff94851681529290931660208301527fffffffffffffffffffffffffffffffff00000000000000000000000000000000168183015290519081900360600190f35b6102516119b6565b610289600480360360208110156105ce57600080fd5b50356001600160a01b0316611a11565b610293611a8b565b61060f600480360360208110156105fc57600080fd5b503569ffffffffffffffffffff16611a9a565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102896004803603604081101561066f57600080fd5b81019060208101813564010000000081111561068a57600080fd5b82018360208201111561069c57600080fd5b803590602001918460208302840111640100000000831117156106be57600080fd5b9193909290916020810190356401000000008111156106dc57600080fd5b8201836020820111156106ee57600080fd5b8035906020019184602083028401116401000000008311171561071057600080fd5b509092509050611bee565b6102896004803603602081101561073157600080fd5b50356001600160a01b0316611e27565b6102516004803603602081101561075757600080fd5b5035611f20565b6102516004803603602081101561077457600080fd5b5035611f56565b610289600480360360a081101561079157600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516611fab565b610289600480360360408110156107d657600080fd5b506001600160a01b03813516906020013561214a565b6102896004803603608081101561080257600080fd5b81019060208101813564010000000081111561081d57600080fd5b82018360208201111561082f57600080fd5b8035906020019184600183028401116401000000008311171561085157600080fd5b91939092909160208101903564010000000081111561086f57600080fd5b82018360208201111561088157600080fd5b803590602001918460208302840111640100000000831117156108a357600080fd5b9193909290916020810190356401000000008111156108c157600080fd5b8201836020820111156108d357600080fd5b803590602001918460208302840111640100000000831117156108f557600080fd5b919350915035612479565b61025161336a565b61092e6004803603602081101561091e57600080fd5b50356001600160a01b031661341b565b6040805161ffff9092168252519081900360200190f35b61094d6134d4565b604080517fffffffffffffffffffffffffffffffff00000000000000000000000000000000909616865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b610289600480360360408110156109c557600080fd5b506001600160a01b03813581169160200135166135c3565b610289600480360360208110156109f357600080fd5b50356001600160a01b031661371f565b61028960048036036020811015610a1957600080fd5b50356001600160a01b03166137e7565b61060f61384f565b6000610a3b6148cc565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115610a7f57fe5b6002811115610a8a57fe5b9052509050600081602001516002811115610aa157fe5b1415610ab1576000915050610b8c565b610ab96148e3565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830181905270010000000000000000000000000000000090910490921660808201528251909160009160019060049060ff16601f8110610b4557fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff16601f8110610b8357fe5b01540301925050505b919050565b6000546001600160a01b03163314610bf0576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b036801000000000000000090910481169082168114610c8657602c80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff16680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000610ce26148e3565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b7f000000000000000000000000000000000000000000000000000000000000000081565b602c546801000000000000000090046001600160a01b031690565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b6020526040902054601790810b900b90565b600481565b868560ff8616601f831115610e43576040805162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b60008111610e98576040805162461bcd60e51b815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b818314610ed65760405162461bcd60e51b8152600401808060200182810382526024815260200180614a656024913960400191505060405180910390fd5b806003028311610f2d576040805162461bcd60e51b815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b6000546001600160a01b03163314610f8c576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6028541561113057602880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81019160009183908110610fc957fe5b6000918252602082200154602980546001600160a01b0390921693509084908110610ff057fe5b6000918252602090912001546001600160a01b03169050611010816138ee565b6001600160a01b0380831660009081526027602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559284168252902080549091169055602880548061106c57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905560298054806110cf57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610f8c915050565b60005b8a81101561153e576000602760008e8e8581811061114d57fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff16600281111561118357fe5b146111d5576040805162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260016020820152602760008e8e858181106111fc57fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561128757fe5b02179055506000915060069050818c8c858181106112a157fe5b6001600160a01b036020918202939093013583168452830193909352604090910160002054169190911415905061131f576040805162461bcd60e51b815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b6000602760008c8c8581811061133157fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff16600281111561136757fe5b146113b9576040805162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260026020820152602760008c8c858181106113e057fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561146b57fe5b021790555090505060288c8c8381811061148157fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03959092029390930135939093169290921790555060298a8a838181106114e357fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0396909302949094013594909416179091555001611133565b50602a805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602c80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000090911617938490559091048116911661160a30828f8f8f8f8f8f8f8f613b1e565b602a60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602a60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600083820152604051601f909101601f19169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff1690565b7f000000000000000000000000000000000000000000000000000000000000000081565b602d8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61010060018816150201909516949094049384018190048102820181019092528281526060939092909183018282801561185c5780601f106118315761010080835404028352916020019161185c565b820191906000526020600020905b81548152906001019060200180831161183f57829003601f168201915b5050505050905090565b6001546001600160a01b031633146118c5576040805162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060602980548060200260200160405190810160405280929190818152602001828054801561185c57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161196e575050505050905090565b602c54602a5463ffffffff808316926401000000009004169060801b909192565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6001600160a01b03818116600090815260066020526040902054163314611a7f576040805162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b611a88816138ee565b50565b6000546001600160a01b031681565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e74000000000000000000000000000000000081525090611b755760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b83811015611b3a578181015183820152602001611b22565b50505050905090810190601f168015611b675780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50611b7e6148cc565b5050505063ffffffff83166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b6000546001600160a01b03163314611c4d576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b828114611ca1576040805162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015611e20576000858583818110611cba57fe5b905060200201356001600160a01b031690506000848484818110611cda57fe5b6001600160a01b038581166000908152600660209081526040909120549202939093013583169350909116905080158080611d265750826001600160a01b0316826001600160a01b0316145b611d77576040805162461bcd60e51b815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b03848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614611e1057826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505060019092019150611ca49050565b5050505050565b6001600160a01b03818116600090815260076020526040902054163314611e95576040805162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b6001600160a01b0381811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600063ffffffff821115611f3657506000610b8c565b5063ffffffff166000908152602b6020526040902054601790810b900b90565b600063ffffffff821115611f6c57506000610b8c565b5063ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6003546001600160a01b031680612009576040805162461bcd60e51b815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b6000546001600160a01b03163314806120dc5750604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b03861694636b14daf8946000939190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b1580156120af57600080fd5b505afa1580156120c3573d6000803e3d6000fd5b505050506040513d60208110156120d957600080fd5b50515b61212d576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b612135613c22565b6121428686868686613ffe565b505050505050565b6000546001600160a01b03163314806122255750600354604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b1580156121f857600080fd5b505afa15801561220c573d6000803e3d6000fd5b505050506040513d602081101561222257600080fd5b50515b612276576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6000612280614178565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b1580156122f157600080fd5b505afa158015612305573d6000803e3d6000fd5b505050506040513d602081101561231b57600080fd5b5051905081811015612374576040805162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb856123b085850387614359565b6040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b1580156123f657600080fd5b505af115801561240a573d6000803e3d6000fd5b505050506040513d602081101561242057600080fd5b5051612473576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a905061248c888888888888614373565b36146124df576040805162461bcd60e51b815260206004820152601960248201527f7472616e736d6974206d65737361676520746f6f206c6f6e6700000000000000604482015290519081900360640190fd5b6124e7614911565b6040805160808082018352602a549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff1660208301527501000000000000000000000000000000000000000000810460ff169282019290925276010000000000000000000000000000000000000000000090910463ffffffff166060808301919091529082526000908a908a908110156125a057600080fd5b8135916020810135918101906060810160408201356401000000008111156125c757600080fd5b8201836020820111156125d957600080fd5b803590602001918460208302840111640100000000831117156125fb57600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050505060408801525050506080840182905283515190925060589190911b907fffffffffffffffffffffffffffffffff000000000000000000000000000000008083169116146126c2576040805162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b608083015183516020015164ffffffffff80831691161061272a576040805162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b83516040015160ff168911612786576040805162461bcd60e51b815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b601f8911156127dc576040805162461bcd60e51b815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b868914612830576040805162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b601f846040015151111561288b576040805162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b83600001516040015160020260ff16846040015151116128f2576040805162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b8867ffffffffffffffff8111801561290957600080fd5b506040519080825280601f01601f191660200182016040528015612934576020820181803683370190505b50606085015260005b60ff81168a11156129a557868160ff166020811061295757fe5b1a60f81b85606001518260ff168151811061296e57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060010161293d565b5083604001515167ffffffffffffffff811180156129c257600080fd5b506040519080825280601f01601f1916602001820160405280156129ed576020820181803683370190505b5060208501526129fb614945565b60005b8560400151518160ff161015612b01576000858260ff1660208110612a1f57fe5b1a90508281601f8110612a2e57fe5b602002015115612a85576040805162461bcd60e51b815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b6001838260ff16601f8110612a9657fe5b91151560209283029190910152869060ff8416908110612ab257fe5b1a60f81b87602001518360ff1681518110612ac957fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350506001016129fe565b50612b0a6148cc565b336000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115612b4557fe5b6002811115612b5057fe5b9052509050600281602001516002811115612b6757fe5b148015612b9b57506029816000015160ff1681548110612b8357fe5b6000918252602090912001546001600160a01b031633145b612bec576040805162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b5050835164ffffffffff90911660209091015250506040516000908a908a908083838082843760405192018290039091209450612c2d935061494592505050565b612c356148cc565b60005b89811015612e2e57600060018587606001518481518110612c5557fe5b60209101015160f81c601b018e8e86818110612c6d57fe5b905060200201358d8d87818110612c8057fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015612cdb573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526027602090815290849020838501909452835460ff80821685529296509294508401916101009004166002811115612d2a57fe5b6002811115612d3557fe5b9052509250600183602001516002811115612d4c57fe5b14612d9e576040805162461bcd60e51b815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff16601f8110612daf57fe5b602002015115612e06576040805162461bcd60e51b815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16601f8110612e1b57fe5b9115156020909202015250600101612c38565b5050505060005b600182604001515103811015612edf57600082604001518260010181518110612e5a57fe5b602002602001015160170b83604001518381518110612e7557fe5b602002602001015160170b1315905080612ed6576040805162461bcd60e51b815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b50600101612e35565b50604081015180516000919060028104908110612ef857fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b13158015612f5e57507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b612faf576040805162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602b82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519586523386890181905291860181905260a0988601898152845199870199909952835194909916997ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451998c999298949793969095909492939185019260c086019289820192909102908190849084905b838110156131125781810151838201526020016130fa565b50505050905001838103825285818151815260200191508051906020019080838360005b8381101561314e578181015183820152602001613136565b50505050905090810190601f16801561317b5780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390a281516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a36132308260000151606001518260170b61438b565b5080518051602a8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff9283160217909155821061335157fe5b61335f82826020015161449c565b505050505050505050565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b1580156133da57600080fd5b505afa1580156133ee573d6000803e3d6000fd5b505050506040513d602081101561340457600080fd5b505190506000613412614178565b90910391505090565b60006134256148cc565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561346957fe5b600281111561347457fe5b905250905060008160200151600281111561348b57fe5b141561349b576000915050610b8c565b60016004826000015160ff16601f81106134b157fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b60008080808033321461352e576040805162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602a5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602b6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6001600160a01b03828116600090815260066020526040902054163314613631576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b336001600160a01b038216141561368f576040805162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b6001600160a01b03808316600090815260076020526040902080548383167fffffffffffffffffffffffff00000000000000000000000000000000000000008216811790925590911690811461371a576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b6000546001600160a01b0316331461377e576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000546001600160a01b03163314613846576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b611a88816146f7565b602a54760100000000000000000000000000000000000000000000900463ffffffff16600080808061387f6148cc565b5050505063ffffffff82166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b6138f66148cc565b6001600160a01b0382166000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561393a57fe5b600281111561394557fe5b9052509050600061395583610a31565b9050801561371a576001600160a01b0380841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b1580156139fe57600080fd5b505af1158015613a12573d6000803e3d6000fd5b505050506040513d6020811015613a2857600080fd5b5051613a7b576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff16601f8110613a9157fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16601f8110613acc57fe5b0155604080516001600160a01b0380871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b60008a8a8a8a8a8a8a8a8a8a604051602001808b6001600160a01b031681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b613c2a6148e3565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152613ca1614945565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411613cba57905050505050509050613d01614945565b604080516103e081019182905290600890601f9082845b815481526020019060010190808311613d18575050505050905060606029805480602002602001604051908101604052809291908181526020018280548015613d8a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613d6c575b5050505050905060005b8151811015613fe257600060018483601f8110613dad57fe5b6020020151039050600060018684601f8110613dc557fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115613fd757600060066000878781518110613e0557fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b031690507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb82846040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b158015613eba57600080fd5b505af1158015613ece573d6000803e3d6000fd5b505050506040513d6020811015613ee457600080fd5b5051613f37576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018886601f8110613f4557fe5b61ffff909216602092909202015260018786601f8110613f6157fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290879087908110613f9657fe5b6020026020010151828460405180846001600160a01b03168152602001836001600160a01b03168152602001828152602001935050505060405180910390a1505b505050600101613d94565b50613ff0600484601f614964565b50611e20600883601f6149fa565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000614182614945565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff168152602001906002019060208260010104928301926001038202915080841161419b5790505050505050905060005b601f81101561420b5760018282601f81106141f457fe5b60200201510361ffff1692909201916001016141dd565b506142146148e3565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606080860182905270010000000000000000000000000000000090940490921660808501526029805486518184028101840190975280875297909202633b9aca00029693949293908301828280156142e457602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116142c6575b505050505090506142f3614945565b604080516103e081019182905290600890601f9082845b81548152602001906001019080831161430a575050505050905060005b82518110156143515760018282601f811061433e57fe5b6020020151039590950194600101614327565b505050505090565b60008183101561436a57508161436d565b50805b92915050565b602083810286019082020160e4019695505050505050565b602c546801000000000000000090046001600160a01b0316806143ae5750610c86565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602b602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b60248501819052948916604485015260648401889052516001600160a01b0387169363beed9b5193620186a09360848084019491939192918390030190829088803b15801561446d57600080fd5b5087f19350505050801561449357506040513d602081101561448e57600080fd5b505160015b61214257611e20565b6144a46148cc565b336000908152602760209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156144df57fe5b60028111156144ea57fe5b90525090506144f76148e3565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000082048116606084015270010000000000000000000000000000000090910416608082015281516103e081019283905290916145c2918591600490601f90826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116145805790505050505050614786565b6145d090600490601f614964565b506002826020015160028111156145e357fe5b14614635576040805162461bcd60e51b815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b600061465c633b9aca003a04836020015163ffffffff16846000015163ffffffff166147fb565b90506010360260005a9050600061467b8863ffffffff16858585614821565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff168302816146a957fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff16601f81106146d257fe5b015401016008886000015160ff16601f81106146ea57fe5b0155505050505050505050565b6003546001600160a01b039081169082168114610c8657600380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b61478e614945565b60005b83518110156147f35760008482815181106147a857fe5b016020015160f81c90506147cd8482601f81106147c157fe5b602002015160016148ad565b848260ff16601f81106147dc57fe5b61ffff909216602092909202015250600101614791565b509092915050565b6000838381101561480e57600285850304015b6148188184614359565b95945050505050565b600081851015614878576040805162461bcd60e51b815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b818503830161179301633b9aca00858202026fffffffffffffffffffffffffffffffff81106148a357fe5b9695505050505050565b60006148c58261ffff168461ffff160161ffff614359565b9392505050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518060a00160405280614924614a28565b81526060602082018190526040820181905280820152600060809091015290565b604051806103e00160405280601f906020820280368337509192915050565b6002830191839082156149ea5791602002820160005b838211156149ba57835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030261497a565b80156149e85782816101000a81549061ffff02191690556002016020816001010492830192600103026149ba565b505b506149f6929150614a4f565b5090565b82601f81019282156149ea579160200282015b828111156149ea578251825591602001919060010190614a0d565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b808211156149f65760008155600101614a5056fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea164736f6c6343000705000a"


func DeployOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _validator common.Address, _minAnswer *big.Int, _maxAnswer *big.Int, _billingAdminAccessController common.Address, _decimals uint8, _description string) (common.Address, *types.Transaction, *OffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OffchainAggregatorBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _validator, _minAnswer, _maxAnswer, _billingAdminAccessController, _decimals, _description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffchainAggregator{OffchainAggregatorCaller: OffchainAggregatorCaller{contract: contract}, OffchainAggregatorTransactor: OffchainAggregatorTransactor{contract: contract}, OffchainAggregatorFilterer: OffchainAggregatorFilterer{contract: contract}}, nil
}


type OffchainAggregator struct {
	OffchainAggregatorCaller     
	OffchainAggregatorTransactor 
	OffchainAggregatorFilterer   
}


type OffchainAggregatorCaller struct {
	contract *bind.BoundContract 
}


type OffchainAggregatorTransactor struct {
	contract *bind.BoundContract 
}


type OffchainAggregatorFilterer struct {
	contract *bind.BoundContract 
}



type OffchainAggregatorSession struct {
	Contract     *OffchainAggregator 
	CallOpts     bind.CallOpts       
	TransactOpts bind.TransactOpts   
}



type OffchainAggregatorCallerSession struct {
	Contract *OffchainAggregatorCaller 
	CallOpts bind.CallOpts             
}



type OffchainAggregatorTransactorSession struct {
	Contract     *OffchainAggregatorTransactor 
	TransactOpts bind.TransactOpts             
}


type OffchainAggregatorRaw struct {
	Contract *OffchainAggregator 
}


type OffchainAggregatorCallerRaw struct {
	Contract *OffchainAggregatorCaller 
}


type OffchainAggregatorTransactorRaw struct {
	Contract *OffchainAggregatorTransactor 
}


func NewOffchainAggregator(address common.Address, backend bind.ContractBackend) (*OffchainAggregator, error) {
	contract, err := bindOffchainAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregator{OffchainAggregatorCaller: OffchainAggregatorCaller{contract: contract}, OffchainAggregatorTransactor: OffchainAggregatorTransactor{contract: contract}, OffchainAggregatorFilterer: OffchainAggregatorFilterer{contract: contract}}, nil
}


func NewOffchainAggregatorCaller(address common.Address, caller bind.ContractCaller) (*OffchainAggregatorCaller, error) {
	contract, err := bindOffchainAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorCaller{contract: contract}, nil
}


func NewOffchainAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OffchainAggregatorTransactor, error) {
	contract, err := bindOffchainAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorTransactor{contract: contract}, nil
}


func NewOffchainAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OffchainAggregatorFilterer, error) {
	contract, err := bindOffchainAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorFilterer{contract: contract}, nil
}


func bindOffchainAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_OffchainAggregator *OffchainAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregator.Contract.OffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}



func (_OffchainAggregator *OffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.OffchainAggregatorTransactor.contract.Transfer(opts)
}


func (_OffchainAggregator *OffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.OffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}





func (_OffchainAggregator *OffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}



func (_OffchainAggregator *OffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.contract.Transfer(opts)
}


func (_OffchainAggregator *OffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.contract.Transact(opts, method, params...)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) LINK() (common.Address, error) {
	return _OffchainAggregator.Contract.LINK(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _OffchainAggregator.Contract.LINK(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) Decimals() (uint8, error) {
	return _OffchainAggregator.Contract.Decimals(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _OffchainAggregator.Contract.Decimals(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) Description() (string, error) {
	return _OffchainAggregator.Contract.Description(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Description() (string, error) {
	return _OffchainAggregator.Contract.Description(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetAnswer(&_OffchainAggregator.CallOpts, _roundId)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetAnswer(&_OffchainAggregator.CallOpts, _roundId)
}




func (_OffchainAggregator *OffchainAggregatorCaller) GetBilling(opts *bind.CallOpts) (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(struct {
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPrice = out[0].(uint32)
	outstruct.ReasonableGasPrice = out[1].(uint32)
	outstruct.MicroLinkPerEth = out[2].(uint32)
	outstruct.LinkGweiPerObservation = out[3].(uint32)
	outstruct.LinkGweiPerTransmission = out[4].(uint32)

	return *outstruct, err

}




func (_OffchainAggregator *OffchainAggregatorSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _OffchainAggregator.Contract.GetBilling(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _OffchainAggregator.Contract.GetBilling(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_OffchainAggregator *OffchainAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OffchainAggregator.Contract.GetRoundData(&_OffchainAggregator.CallOpts, _roundId)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OffchainAggregator.Contract.GetRoundData(&_OffchainAggregator.CallOpts, _roundId)
}




func (_OffchainAggregator *OffchainAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetTimestamp(&_OffchainAggregator.CallOpts, _roundId)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetTimestamp(&_OffchainAggregator.CallOpts, _roundId)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [16]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = out[0].(uint32)
	outstruct.BlockNumber = out[1].(uint32)
	outstruct.ConfigDigest = out[2].([16]byte)

	return *outstruct, err

}




func (_OffchainAggregator *OffchainAggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _OffchainAggregator.Contract.LatestConfigDetails(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _OffchainAggregator.Contract.LatestConfigDetails(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) LatestRound() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestRound(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestRound(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_OffchainAggregator *OffchainAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OffchainAggregator.Contract.LatestRoundData(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OffchainAggregator.Contract.LatestRoundData(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestTimestamp(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestTimestamp(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestTransmissionDetails")

	outstruct := new(struct {
		ConfigDigest    [16]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigDigest = out[0].([16]byte)
	outstruct.Epoch = out[1].(uint32)
	outstruct.Round = out[2].(uint8)
	outstruct.LatestAnswer = out[3].(*big.Int)
	outstruct.LatestTimestamp = out[4].(uint64)

	return *outstruct, err

}




func (_OffchainAggregator *OffchainAggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _OffchainAggregator.Contract.LatestTransmissionDetails(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _OffchainAggregator.Contract.LatestTransmissionDetails(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregator.Contract.LinkAvailableForPayment(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregator.Contract.LinkAvailableForPayment(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MaxAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MaxAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MinAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MinAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregator.Contract.OracleObservationCount(&_OffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregator.Contract.OracleObservationCount(&_OffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_OffchainAggregator *OffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregator.Contract.OwedPayment(&_OffchainAggregator.CallOpts, _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregator.Contract.OwedPayment(&_OffchainAggregator.CallOpts, _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) Owner() (common.Address, error) {
	return _OffchainAggregator.Contract.Owner(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _OffchainAggregator.Contract.Owner(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _OffchainAggregator.Contract.Transmitters(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _OffchainAggregator.Contract.Transmitters(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "validator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) Validator() (common.Address, error) {
	return _OffchainAggregator.Contract.Validator(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _OffchainAggregator.Contract.Validator(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregator *OffchainAggregatorSession) Version() (*big.Int, error) {
	return _OffchainAggregator.Contract.Version(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Version() (*big.Int, error) {
	return _OffchainAggregator.Contract.Version(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "acceptOwnership")
}




func (_OffchainAggregator *OffchainAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffchainAggregator.Contract.AcceptOwnership(&_OffchainAggregator.TransactOpts)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffchainAggregator.Contract.AcceptOwnership(&_OffchainAggregator.TransactOpts)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "acceptPayeeship", _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.AcceptPayeeship(&_OffchainAggregator.TransactOpts, _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.AcceptPayeeship(&_OffchainAggregator.TransactOpts, _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_OffchainAggregator *OffchainAggregatorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetBilling(&_OffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetBilling(&_OffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setBillingAccessController", _billingAdminAccessController)
}




func (_OffchainAggregator *OffchainAggregatorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetBillingAccessController(&_OffchainAggregator.TransactOpts, _billingAdminAccessController)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetBillingAccessController(&_OffchainAggregator.TransactOpts, _billingAdminAccessController)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_OffchainAggregator *OffchainAggregatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetConfig(&_OffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetConfig(&_OffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setPayees", _transmitters, _payees)
}




func (_OffchainAggregator *OffchainAggregatorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetPayees(&_OffchainAggregator.TransactOpts, _transmitters, _payees)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetPayees(&_OffchainAggregator.TransactOpts, _transmitters, _payees)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) SetValidator(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setValidator", _newValidator)
}




func (_OffchainAggregator *OffchainAggregatorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetValidator(&_OffchainAggregator.TransactOpts, _newValidator)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetValidator(&_OffchainAggregator.TransactOpts, _newValidator)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "transferOwnership", _to)
}




func (_OffchainAggregator *OffchainAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.TransferOwnership(&_OffchainAggregator.TransactOpts, _to)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.TransferOwnership(&_OffchainAggregator.TransactOpts, _to)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}




func (_OffchainAggregator *OffchainAggregatorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.TransferPayeeship(&_OffchainAggregator.TransactOpts, _transmitter, _proposed)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.TransferPayeeship(&_OffchainAggregator.TransactOpts, _transmitter, _proposed)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) Transmit(opts *bind.TransactOpts, _report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "transmit", _report, _rs, _ss, _rawVs)
}




func (_OffchainAggregator *OffchainAggregatorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.Transmit(&_OffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.Transmit(&_OffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}




func (_OffchainAggregator *OffchainAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.WithdrawFunds(&_OffchainAggregator.TransactOpts, _recipient, _amount)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.WithdrawFunds(&_OffchainAggregator.TransactOpts, _recipient, _amount)
}




func (_OffchainAggregator *OffchainAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "withdrawPayment", _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.WithdrawPayment(&_OffchainAggregator.TransactOpts, _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.WithdrawPayment(&_OffchainAggregator.TransactOpts, _transmitter)
}


type OffchainAggregatorAnswerUpdatedIterator struct {
	Event *OffchainAggregatorAnswerUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorAnswerUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorAnswerUpdated)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorAnswerUpdated)
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


func (it *OffchainAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*OffchainAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorAnswerUpdatedIterator{contract: _OffchainAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorAnswerUpdated)
				if err := _OffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*OffchainAggregatorAnswerUpdated, error) {
	event := new(OffchainAggregatorAnswerUpdated)
	if err := _OffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorBillingAccessControllerSetIterator struct {
	Event *OffchainAggregatorBillingAccessControllerSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingAccessControllerSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingAccessControllerSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingAccessControllerSet)
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


func (it *OffchainAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingAccessControllerSetIterator{contract: _OffchainAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingAccessControllerSet)
				if err := _OffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*OffchainAggregatorBillingAccessControllerSet, error) {
	event := new(OffchainAggregatorBillingAccessControllerSet)
	if err := _OffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorBillingSetIterator struct {
	Event *OffchainAggregatorBillingSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingSet)
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


func (it *OffchainAggregatorBillingSetIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*OffchainAggregatorBillingSetIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingSetIterator{contract: _OffchainAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingSet)
				if err := _OffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseBillingSet(log types.Log) (*OffchainAggregatorBillingSet, error) {
	event := new(OffchainAggregatorBillingSet)
	if err := _OffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorConfigSetIterator struct {
	Event *OffchainAggregatorConfigSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorConfigSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorConfigSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorConfigSet)
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


func (it *OffchainAggregatorConfigSetIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OffchainAggregatorConfigSetIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorConfigSetIterator{contract: _OffchainAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorConfigSet)
				if err := _OffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseConfigSet(log types.Log) (*OffchainAggregatorConfigSet, error) {
	event := new(OffchainAggregatorConfigSet)
	if err := _OffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorNewRoundIterator struct {
	Event *OffchainAggregatorNewRound 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorNewRoundIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorNewRound)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorNewRound)
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


func (it *OffchainAggregatorNewRoundIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*OffchainAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorNewRoundIterator{contract: _OffchainAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorNewRound)
				if err := _OffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseNewRound(log types.Log) (*OffchainAggregatorNewRound, error) {
	event := new(OffchainAggregatorNewRound)
	if err := _OffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorNewTransmissionIterator struct {
	Event *OffchainAggregatorNewTransmission 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorNewTransmissionIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorNewTransmission)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorNewTransmission)
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


func (it *OffchainAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	RawReportContext  [32]byte
	Raw               types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*OffchainAggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorNewTransmissionIterator{contract: _OffchainAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorNewTransmission)
				if err := _OffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseNewTransmission(log types.Log) (*OffchainAggregatorNewTransmission, error) {
	event := new(OffchainAggregatorNewTransmission)
	if err := _OffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorOraclePaidIterator struct {
	Event *OffchainAggregatorOraclePaid 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorOraclePaidIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorOraclePaid)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorOraclePaid)
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


func (it *OffchainAggregatorOraclePaidIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	Raw         types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts) (*OffchainAggregatorOraclePaidIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorOraclePaidIterator{contract: _OffchainAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorOraclePaid) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorOraclePaid)
				if err := _OffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseOraclePaid(log types.Log) (*OffchainAggregatorOraclePaid, error) {
	event := new(OffchainAggregatorOraclePaid)
	if err := _OffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorOwnershipTransferRequestedIterator struct {
	Event *OffchainAggregatorOwnershipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorOwnershipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorOwnershipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorOwnershipTransferRequested)
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


func (it *OffchainAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffchainAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorOwnershipTransferRequestedIterator{contract: _OffchainAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorOwnershipTransferRequested)
				if err := _OffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffchainAggregatorOwnershipTransferRequested, error) {
	event := new(OffchainAggregatorOwnershipTransferRequested)
	if err := _OffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorOwnershipTransferredIterator struct {
	Event *OffchainAggregatorOwnershipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorOwnershipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorOwnershipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorOwnershipTransferred)
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


func (it *OffchainAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffchainAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorOwnershipTransferredIterator{contract: _OffchainAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorOwnershipTransferred)
				if err := _OffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*OffchainAggregatorOwnershipTransferred, error) {
	event := new(OffchainAggregatorOwnershipTransferred)
	if err := _OffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorPayeeshipTransferRequestedIterator struct {
	Event *OffchainAggregatorPayeeshipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorPayeeshipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorPayeeshipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorPayeeshipTransferRequested)
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


func (it *OffchainAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*OffchainAggregatorPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorPayeeshipTransferRequestedIterator{contract: _OffchainAggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorPayeeshipTransferRequested)
				if err := _OffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*OffchainAggregatorPayeeshipTransferRequested, error) {
	event := new(OffchainAggregatorPayeeshipTransferRequested)
	if err := _OffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorPayeeshipTransferredIterator struct {
	Event *OffchainAggregatorPayeeshipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorPayeeshipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorPayeeshipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorPayeeshipTransferred)
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


func (it *OffchainAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*OffchainAggregatorPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorPayeeshipTransferredIterator{contract: _OffchainAggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorPayeeshipTransferred)
				if err := _OffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*OffchainAggregatorPayeeshipTransferred, error) {
	event := new(OffchainAggregatorPayeeshipTransferred)
	if err := _OffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorValidatorUpdatedIterator struct {
	Event *OffchainAggregatorValidatorUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorValidatorUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorValidatorUpdated)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorValidatorUpdated)
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


func (it *OffchainAggregatorValidatorUpdatedIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorValidatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorValidatorUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log 
}




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterValidatorUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*OffchainAggregatorValidatorUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorValidatorUpdatedIterator{contract: _OffchainAggregator.contract, event: "ValidatorUpdated", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchValidatorUpdated(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorValidatorUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorValidatorUpdated)
				if err := _OffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
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




func (_OffchainAggregator *OffchainAggregatorFilterer) ParseValidatorUpdated(log types.Log) (*OffchainAggregatorValidatorUpdated, error) {
	event := new(OffchainAggregatorValidatorUpdated)
	if err := _OffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const OffchainAggregatorBillingABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var OffchainAggregatorBillingBin = "0x60a06040523480156200001157600080fd5b506040516200276c3803806200276c833981810160405260e08110156200003757600080fd5b508051602082015160408301516060840151608085015160a086015160c090960151600080546001600160a01b03191633179055949593949293919290919062000085878787878762000136565b620000908162000228565b6001600160601b0319606083901b16608052620000ac620002a1565b620000b6620002a1565b60005b601f8160ff16101562000106576001838260ff16601f8110620000d857fe5b61ffff909216602092909202015260018260ff8316601f8110620000f857fe5b6020020152600101620000b9565b5062000116600483601f620002c0565b5062000126600882601f6200035d565b50505050505050505050620003a5565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b0390811690821681146200029d57600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b604051806103e00160405280601f906020820280368337509192915050565b6002830191839082156200034b5791602002820160005b838211156200031957835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302620002d7565b8015620003495782816101000a81549061ffff021916905560020160208160010104928301926001030262000319565b505b50620003599291506200038e565b5090565b82601f81019282156200034b579160200282015b828111156200034b57825182559160200191906001019062000371565b5b808211156200035957600081556001016200038f565b60805160601c612390620003dc600039806105cb528061101052806111345280611270528061182d5280611c1352506123906000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063b121e14711610097578063e4902f8211610066578063e4902f8214610371578063eb5dcd6c146103bb578063f2fde38b146103f6578063fbffd2c114610429576100f5565b8063b121e147146102b8578063bd824706146102eb578063c107532914610330578063d09dc33914610369576100f5565b806379ba5097116100d357806379ba5097146101b15780638ac28d5a146101bb5780638da5cb5b146101ee5780639c849b30146101f6576100f5565b80630eafb25b146100fa5780631b6b6d231461013f5780632993726814610170575b600080fd5b61012d6004803603602081101561011057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661045c565b60408051918252519081900360200190f35b6101476105c9565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101786105ed565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b6101b9610682565b005b6101b9600480360360208110156101d157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610784565b610147610825565b6101b96004803603604081101561020c57600080fd5b81019060208101813564010000000081111561022757600080fd5b82018360208201111561023957600080fd5b8035906020019184602083028401116401000000008311171561025b57600080fd5b91939092909160208101903564010000000081111561027957600080fd5b82018360208201111561028b57600080fd5b803590602001918460208302840111640100000000831117156102ad57600080fd5b509092509050610841565b6101b9600480360360208110156102ce57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610b3d565b6101b9600480360360a081101561030157600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516610c6a565b6101b96004803603604081101561034657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610e82565b61012d61126b565b6103a46004803603602081101561038757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611336565b6040805161ffff9092168252519081900360200190f35b6101b9600480360360408110156103d157600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160200135166113fc565b6101b96004803603602081101561040c57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166115c0565b6101b96004803603602081101561043f57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166116bc565b6000610466612246565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602760209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156104b757fe5b60028111156104c257fe5b90525090506000816020015160028111156104d957fe5b14156104e95760009150506105c4565b6104f161225d565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830181905270010000000000000000000000000000000090910490921660808201528251909160009160019060049060ff16601f811061057d57fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff16601f81106105bb57fe5b01540301925050505b919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60008060008060006105fd61225d565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b60015473ffffffffffffffffffffffffffffffffffffffff16331461070857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526006602052604090205416331461081957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b6108228161174b565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108c757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b82811461093557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015610b3657600085858381811061094e57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff169050600084848481811061097b57fe5b73ffffffffffffffffffffffffffffffffffffffff85811660009081526006602090815260409091205492029390930135831693509091169050801580806109ee57508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b610a5957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614610b26578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506109389050565b5050505050565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260076020526040902054163314610bd257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60035473ffffffffffffffffffffffffffffffffffffffff1680610cef57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b60005473ffffffffffffffffffffffffffffffffffffffff16331480610dfa5750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610dcd57600080fd5b505afa158015610de1573d6000803e3d6000fd5b505050506040513d6020811015610df757600080fd5b50515b610e6557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b610e6d6119bc565b610e7a8686868686611e1a565b505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331480610f955750600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610f6857600080fd5b505afa158015610f7c573d6000803e3d6000fd5b505050506040513d6020811015610f9257600080fd5b50515b61100057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b600061100a611f94565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561109557600080fd5b505afa1580156110a9573d6000803e3d6000fd5b505050506040513d60208110156110bf57600080fd5b505190508181101561113257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8561117b85850387612182565b6040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156111ce57600080fd5b505af11580156111e2573d6000803e3d6000fd5b505050506040513d60208110156111f857600080fd5b505161126557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156112f557600080fd5b505afa158015611309573d6000803e3d6000fd5b505050506040513d602081101561131f57600080fd5b50519050600061132d611f94565b90910391505090565b6000611340612246565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561139157fe5b600281111561139c57fe5b90525090506000816020015160028111156113b357fe5b14156113c35760009150506105c4565b60016004826000015160ff16601f81106113d957fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526006602052604090205416331461149157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b3373ffffffffffffffffffffffffffffffffffffffff8216141561151657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff0000000000000000000000000000000000000000821681179092559091169081146115bb5760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461164657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461174257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6108228161219c565b611753612246565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602760209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156117a457fe5b60028111156117af57fe5b905250905060006117bf8361045c565b905080156115bb5773ffffffffffffffffffffffffffffffffffffffff80841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b15801561187557600080fd5b505af1158015611889573d6000803e3d6000fd5b505050506040513d602081101561189f57600080fd5b505161190c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff16601f811061192257fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16601f811061195d57fe5b01556040805173ffffffffffffffffffffffffffffffffffffffff80871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b6119c461225d565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152611a3b61228b565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411611a5457905050505050509050611a9b61228b565b604080516103e081019182905290600890601f9082845b815481526020019060010190808311611ab2575050505050905060606029805480602002602001604051908101604052809291908181526020018280548015611b3157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611b06575b5050505050905060005b8151811015611dfe57600060018483601f8110611b5457fe5b6020020151039050600060018684601f8110611b6c57fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115611df357600060066000878781518110611bac57fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015611ca257600080fd5b505af1158015611cb6573d6000803e3d6000fd5b505050506040513d6020811015611ccc57600080fd5b5051611d3957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018886601f8110611d4757fe5b61ffff909216602092909202015260018786601f8110611d6357fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290879087908110611d9857fe5b60200260200101518284604051808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a1505b505050600101611b3b565b50611e0c600484601f6122aa565b50610b36600883601f612340565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000611f9e61228b565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411611fb75790505050505050905060005b601f8110156120275760018282601f811061201057fe5b60200201510361ffff169290920191600101611ff9565b5061203061225d565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606080860182905270010000000000000000000000000000000090940490921660808501526029805486518184028101840190975280875297909202633b9aca000296939492939083018282801561210d57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116120e2575b5050505050905061211c61228b565b604080516103e081019182905290600890601f9082845b815481526020019060010190808311612133575050505050905060005b825181101561217a5760018282601f811061216757fe5b6020020151039590950194600101612150565b505050505090565b600081831015612193575081612196565b50805b92915050565b60035473ffffffffffffffffffffffffffffffffffffffff908116908216811461224257600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b604051806103e00160405280601f906020820280368337509192915050565b6002830191839082156123305791602002820160005b8382111561230057835183826101000a81548161ffff021916908361ffff16021790555092602001926002016020816001010492830192600103026122c0565b801561232e5782816101000a81549061ffff0219169055600201602081600101049283019260010302612300565b505b5061233c92915061236e565b5090565b82601f8101928215612330579160200282015b82811115612330578251825591602001919060010190612353565b5b8082111561233c576000815560010161236f56fea164736f6c6343000705000a"


func DeployOffchainAggregatorBilling(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _billingAdminAccessController common.Address) (common.Address, *types.Transaction, *OffchainAggregatorBilling, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorBillingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OffchainAggregatorBillingBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _billingAdminAccessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffchainAggregatorBilling{OffchainAggregatorBillingCaller: OffchainAggregatorBillingCaller{contract: contract}, OffchainAggregatorBillingTransactor: OffchainAggregatorBillingTransactor{contract: contract}, OffchainAggregatorBillingFilterer: OffchainAggregatorBillingFilterer{contract: contract}}, nil
}


type OffchainAggregatorBilling struct {
	OffchainAggregatorBillingCaller     
	OffchainAggregatorBillingTransactor 
	OffchainAggregatorBillingFilterer   
}


type OffchainAggregatorBillingCaller struct {
	contract *bind.BoundContract 
}


type OffchainAggregatorBillingTransactor struct {
	contract *bind.BoundContract 
}


type OffchainAggregatorBillingFilterer struct {
	contract *bind.BoundContract 
}



type OffchainAggregatorBillingSession struct {
	Contract     *OffchainAggregatorBilling 
	CallOpts     bind.CallOpts              
	TransactOpts bind.TransactOpts          
}



type OffchainAggregatorBillingCallerSession struct {
	Contract *OffchainAggregatorBillingCaller 
	CallOpts bind.CallOpts                    
}



type OffchainAggregatorBillingTransactorSession struct {
	Contract     *OffchainAggregatorBillingTransactor 
	TransactOpts bind.TransactOpts                    
}


type OffchainAggregatorBillingRaw struct {
	Contract *OffchainAggregatorBilling 
}


type OffchainAggregatorBillingCallerRaw struct {
	Contract *OffchainAggregatorBillingCaller 
}


type OffchainAggregatorBillingTransactorRaw struct {
	Contract *OffchainAggregatorBillingTransactor 
}


func NewOffchainAggregatorBilling(address common.Address, backend bind.ContractBackend) (*OffchainAggregatorBilling, error) {
	contract, err := bindOffchainAggregatorBilling(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBilling{OffchainAggregatorBillingCaller: OffchainAggregatorBillingCaller{contract: contract}, OffchainAggregatorBillingTransactor: OffchainAggregatorBillingTransactor{contract: contract}, OffchainAggregatorBillingFilterer: OffchainAggregatorBillingFilterer{contract: contract}}, nil
}


func NewOffchainAggregatorBillingCaller(address common.Address, caller bind.ContractCaller) (*OffchainAggregatorBillingCaller, error) {
	contract, err := bindOffchainAggregatorBilling(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingCaller{contract: contract}, nil
}


func NewOffchainAggregatorBillingTransactor(address common.Address, transactor bind.ContractTransactor) (*OffchainAggregatorBillingTransactor, error) {
	contract, err := bindOffchainAggregatorBilling(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingTransactor{contract: contract}, nil
}


func NewOffchainAggregatorBillingFilterer(address common.Address, filterer bind.ContractFilterer) (*OffchainAggregatorBillingFilterer, error) {
	contract, err := bindOffchainAggregatorBilling(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingFilterer{contract: contract}, nil
}


func bindOffchainAggregatorBilling(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorBillingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingCaller.contract.Call(opts, result, method, params...)
}



func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingTransactor.contract.Transfer(opts)
}


func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingTransactor.contract.Transact(opts, method, params...)
}





func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregatorBilling.Contract.contract.Call(opts, result, method, params...)
}



func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.contract.Transfer(opts)
}


func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.contract.Transact(opts, method, params...)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) LINK() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.LINK(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) LINK() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.LINK(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) GetBilling(opts *bind.CallOpts) (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "getBilling")

	outstruct := new(struct {
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPrice = out[0].(uint32)
	outstruct.ReasonableGasPrice = out[1].(uint32)
	outstruct.MicroLinkPerEth = out[2].(uint32)
	outstruct.LinkGweiPerObservation = out[3].(uint32)
	outstruct.LinkGweiPerTransmission = out[4].(uint32)

	return *outstruct, err

}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _OffchainAggregatorBilling.Contract.GetBilling(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _OffchainAggregatorBilling.Contract.GetBilling(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.LinkAvailableForPayment(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.LinkAvailableForPayment(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregatorBilling.Contract.OracleObservationCount(&_OffchainAggregatorBilling.CallOpts, _signerOrTransmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregatorBilling.Contract.OracleObservationCount(&_OffchainAggregatorBilling.CallOpts, _signerOrTransmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.OwedPayment(&_OffchainAggregatorBilling.CallOpts, _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.OwedPayment(&_OffchainAggregatorBilling.CallOpts, _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) Owner() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.Owner(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) Owner() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.Owner(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "acceptOwnership")
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.AcceptOwnership(&_OffchainAggregatorBilling.TransactOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.AcceptOwnership(&_OffchainAggregatorBilling.TransactOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "acceptPayeeship", _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.AcceptPayeeship(&_OffchainAggregatorBilling.TransactOpts, _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.AcceptPayeeship(&_OffchainAggregatorBilling.TransactOpts, _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetBilling(&_OffchainAggregatorBilling.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetBilling(&_OffchainAggregatorBilling.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "setBillingAccessController", _billingAdminAccessController)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetBillingAccessController(&_OffchainAggregatorBilling.TransactOpts, _billingAdminAccessController)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetBillingAccessController(&_OffchainAggregatorBilling.TransactOpts, _billingAdminAccessController)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "setPayees", _transmitters, _payees)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetPayees(&_OffchainAggregatorBilling.TransactOpts, _transmitters, _payees)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetPayees(&_OffchainAggregatorBilling.TransactOpts, _transmitters, _payees)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "transferOwnership", _to)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.TransferOwnership(&_OffchainAggregatorBilling.TransactOpts, _to)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.TransferOwnership(&_OffchainAggregatorBilling.TransactOpts, _to)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.TransferPayeeship(&_OffchainAggregatorBilling.TransactOpts, _transmitter, _proposed)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.TransferPayeeship(&_OffchainAggregatorBilling.TransactOpts, _transmitter, _proposed)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.WithdrawFunds(&_OffchainAggregatorBilling.TransactOpts, _recipient, _amount)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.WithdrawFunds(&_OffchainAggregatorBilling.TransactOpts, _recipient, _amount)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "withdrawPayment", _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.WithdrawPayment(&_OffchainAggregatorBilling.TransactOpts, _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.WithdrawPayment(&_OffchainAggregatorBilling.TransactOpts, _transmitter)
}


type OffchainAggregatorBillingBillingAccessControllerSetIterator struct {
	Event *OffchainAggregatorBillingBillingAccessControllerSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingBillingAccessControllerSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingBillingAccessControllerSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingBillingAccessControllerSet)
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


func (it *OffchainAggregatorBillingBillingAccessControllerSetIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log 
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorBillingBillingAccessControllerSetIterator, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingBillingAccessControllerSetIterator{contract: _OffchainAggregatorBilling.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingBillingAccessControllerSet)
				if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseBillingAccessControllerSet(log types.Log) (*OffchainAggregatorBillingBillingAccessControllerSet, error) {
	event := new(OffchainAggregatorBillingBillingAccessControllerSet)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorBillingBillingSetIterator struct {
	Event *OffchainAggregatorBillingBillingSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingBillingSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingBillingSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingBillingSet)
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


func (it *OffchainAggregatorBillingBillingSetIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log 
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterBillingSet(opts *bind.FilterOpts) (*OffchainAggregatorBillingBillingSetIterator, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingBillingSetIterator{contract: _OffchainAggregatorBilling.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingBillingSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingBillingSet)
				if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "BillingSet", log); err != nil {
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




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseBillingSet(log types.Log) (*OffchainAggregatorBillingBillingSet, error) {
	event := new(OffchainAggregatorBillingBillingSet)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorBillingOraclePaidIterator struct {
	Event *OffchainAggregatorBillingOraclePaid 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingOraclePaidIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingOraclePaid)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingOraclePaid)
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


func (it *OffchainAggregatorBillingOraclePaidIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	Raw         types.Log 
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterOraclePaid(opts *bind.FilterOpts) (*OffchainAggregatorBillingOraclePaidIterator, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingOraclePaidIterator{contract: _OffchainAggregatorBilling.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingOraclePaid) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.WatchLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingOraclePaid)
				if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseOraclePaid(log types.Log) (*OffchainAggregatorBillingOraclePaid, error) {
	event := new(OffchainAggregatorBillingOraclePaid)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorBillingOwnershipTransferRequestedIterator struct {
	Event *OffchainAggregatorBillingOwnershipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingOwnershipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingOwnershipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingOwnershipTransferRequested)
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


func (it *OffchainAggregatorBillingOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffchainAggregatorBillingOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingOwnershipTransferRequestedIterator{contract: _OffchainAggregatorBilling.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregatorBilling.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingOwnershipTransferRequested)
				if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffchainAggregatorBillingOwnershipTransferRequested, error) {
	event := new(OffchainAggregatorBillingOwnershipTransferRequested)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorBillingOwnershipTransferredIterator struct {
	Event *OffchainAggregatorBillingOwnershipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingOwnershipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingOwnershipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingOwnershipTransferred)
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


func (it *OffchainAggregatorBillingOwnershipTransferredIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffchainAggregatorBillingOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingOwnershipTransferredIterator{contract: _OffchainAggregatorBilling.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregatorBilling.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingOwnershipTransferred)
				if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseOwnershipTransferred(log types.Log) (*OffchainAggregatorBillingOwnershipTransferred, error) {
	event := new(OffchainAggregatorBillingOwnershipTransferred)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorBillingPayeeshipTransferRequestedIterator struct {
	Event *OffchainAggregatorBillingPayeeshipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingPayeeshipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingPayeeshipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingPayeeshipTransferRequested)
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


func (it *OffchainAggregatorBillingPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log 
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*OffchainAggregatorBillingPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingPayeeshipTransferRequestedIterator{contract: _OffchainAggregatorBilling.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _OffchainAggregatorBilling.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingPayeeshipTransferRequested)
				if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParsePayeeshipTransferRequested(log types.Log) (*OffchainAggregatorBillingPayeeshipTransferRequested, error) {
	event := new(OffchainAggregatorBillingPayeeshipTransferRequested)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OffchainAggregatorBillingPayeeshipTransferredIterator struct {
	Event *OffchainAggregatorBillingPayeeshipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OffchainAggregatorBillingPayeeshipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorBillingPayeeshipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorBillingPayeeshipTransferred)
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


func (it *OffchainAggregatorBillingPayeeshipTransferredIterator) Error() error {
	return it.fail
}



func (it *OffchainAggregatorBillingPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OffchainAggregatorBillingPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log 
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*OffchainAggregatorBillingPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingPayeeshipTransferredIterator{contract: _OffchainAggregatorBilling.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorBillingPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OffchainAggregatorBilling.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OffchainAggregatorBillingPayeeshipTransferred)
				if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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




func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParsePayeeshipTransferred(log types.Log) (*OffchainAggregatorBillingPayeeshipTransferred, error) {
	event := new(OffchainAggregatorBillingPayeeshipTransferred)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const OwnedABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var OwnedBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556102db806100326000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b14610081575b600080fd5b61004e6100b4565b005b6100586101b6565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e6004803603602081101561009757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166101d2565b60015473ffffffffffffffffffffffffffffffffffffffff16331461013a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461025857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a35056fea164736f6c6343000705000a"


func DeployOwned(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Owned, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}


type Owned struct {
	OwnedCaller     
	OwnedTransactor 
	OwnedFilterer   
}


type OwnedCaller struct {
	contract *bind.BoundContract 
}


type OwnedTransactor struct {
	contract *bind.BoundContract 
}


type OwnedFilterer struct {
	contract *bind.BoundContract 
}



type OwnedSession struct {
	Contract     *Owned            
	CallOpts     bind.CallOpts     
	TransactOpts bind.TransactOpts 
}



type OwnedCallerSession struct {
	Contract *OwnedCaller  
	CallOpts bind.CallOpts 
}



type OwnedTransactorSession struct {
	Contract     *OwnedTransactor  
	TransactOpts bind.TransactOpts 
}


type OwnedRaw struct {
	Contract *Owned 
}


type OwnedCallerRaw struct {
	Contract *OwnedCaller 
}


type OwnedTransactorRaw struct {
	Contract *OwnedTransactor 
}


func NewOwned(address common.Address, backend bind.ContractBackend) (*Owned, error) {
	contract, err := bindOwned(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}


func NewOwnedCaller(address common.Address, caller bind.ContractCaller) (*OwnedCaller, error) {
	contract, err := bindOwned(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedCaller{contract: contract}, nil
}


func NewOwnedTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTransactor, error) {
	contract, err := bindOwned(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTransactor{contract: contract}, nil
}


func NewOwnedFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedFilterer, error) {
	contract, err := bindOwned(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedFilterer{contract: contract}, nil
}


func bindOwned(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}



func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}


func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}





func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}



func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}


func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}




func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Owned.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_Owned *OwnedSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}




func (_Owned *OwnedCallerSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}




func (_Owned *OwnedTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "acceptOwnership")
}




func (_Owned *OwnedSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}




func (_Owned *OwnedTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}




func (_Owned *OwnedTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "transferOwnership", _to)
}




func (_Owned *OwnedSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Owned.Contract.TransferOwnership(&_Owned.TransactOpts, _to)
}




func (_Owned *OwnedTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Owned.Contract.TransferOwnership(&_Owned.TransactOpts, _to)
}


type OwnedOwnershipTransferRequestedIterator struct {
	Event *OwnedOwnershipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OwnedOwnershipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedOwnershipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OwnedOwnershipTransferRequested)
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


func (it *OwnedOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *OwnedOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OwnedOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_Owned *OwnedFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnedOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Owned.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnedOwnershipTransferRequestedIterator{contract: _Owned.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}




func (_Owned *OwnedFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OwnedOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Owned.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OwnedOwnershipTransferRequested)
				if err := _Owned.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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




func (_Owned *OwnedFilterer) ParseOwnershipTransferRequested(log types.Log) (*OwnedOwnershipTransferRequested, error) {
	event := new(OwnedOwnershipTransferRequested)
	if err := _Owned.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type OwnedOwnershipTransferredIterator struct {
	Event *OwnedOwnershipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *OwnedOwnershipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedOwnershipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(OwnedOwnershipTransferred)
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


func (it *OwnedOwnershipTransferredIterator) Error() error {
	return it.fail
}



func (it *OwnedOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type OwnedOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_Owned *OwnedFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OwnedOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Owned.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OwnedOwnershipTransferredIterator{contract: _Owned.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}




func (_Owned *OwnedFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnedOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Owned.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(OwnedOwnershipTransferred)
				if err := _Owned.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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




func (_Owned *OwnedFilterer) ParseOwnershipTransferred(log types.Log) (*OwnedOwnershipTransferred, error) {
	event := new(OwnedOwnershipTransferred)
	if err := _Owned.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const SimpleReadAccessControllerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var SimpleReadAccessControllerBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556001805460ff60a01b1916600160a01b1790556109c4806100456000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638823da6c11610076578063a118f2491161005b578063a118f249146101fd578063dc7f012414610230578063f2fde38b14610238576100a3565b80638823da6c146101995780638da5cb5b146101cc576100a3565b80630a756983146100a85780636b14daf8146100b257806379ba5097146101895780638038e4a114610191575b600080fd5b6100b061026b565b005b610175600480360360408110156100c857600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561010057600080fd5b82018360208201111561011257600080fd5b8035906020019184600183028401116401000000008311171561013457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610368945050505050565b604080519115158252519081900360200190f35b6100b061039b565b6100b061049d565b6100b0600480360360208110156101af57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166105af565b6101d46106e7565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100b06004803603602081101561021357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610703565b610175610792565b6100b06004803603602081101561024e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166107b3565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff161561036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b600061037483836108af565b80610394575073ffffffffffffffffffffffffffffffffffffffff831632145b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461042157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461052357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff1661036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60005473ffffffffffffffffffffffffffffffffffffffff16331461063557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff16156106e45773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461078957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6106e481610904565b60015474010000000000000000000000000000000000000000900460ff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461083957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff168061039457505060015474010000000000000000000000000000000000000000900460ff161592915050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff166106e45773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a15056fea164736f6c6343000705000a"


func DeploySimpleReadAccessController(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleReadAccessController, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleReadAccessControllerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SimpleReadAccessControllerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleReadAccessController{SimpleReadAccessControllerCaller: SimpleReadAccessControllerCaller{contract: contract}, SimpleReadAccessControllerTransactor: SimpleReadAccessControllerTransactor{contract: contract}, SimpleReadAccessControllerFilterer: SimpleReadAccessControllerFilterer{contract: contract}}, nil
}


type SimpleReadAccessController struct {
	SimpleReadAccessControllerCaller     
	SimpleReadAccessControllerTransactor 
	SimpleReadAccessControllerFilterer   
}


type SimpleReadAccessControllerCaller struct {
	contract *bind.BoundContract 
}


type SimpleReadAccessControllerTransactor struct {
	contract *bind.BoundContract 
}


type SimpleReadAccessControllerFilterer struct {
	contract *bind.BoundContract 
}



type SimpleReadAccessControllerSession struct {
	Contract     *SimpleReadAccessController 
	CallOpts     bind.CallOpts               
	TransactOpts bind.TransactOpts           
}



type SimpleReadAccessControllerCallerSession struct {
	Contract *SimpleReadAccessControllerCaller 
	CallOpts bind.CallOpts                     
}



type SimpleReadAccessControllerTransactorSession struct {
	Contract     *SimpleReadAccessControllerTransactor 
	TransactOpts bind.TransactOpts                     
}


type SimpleReadAccessControllerRaw struct {
	Contract *SimpleReadAccessController 
}


type SimpleReadAccessControllerCallerRaw struct {
	Contract *SimpleReadAccessControllerCaller 
}


type SimpleReadAccessControllerTransactorRaw struct {
	Contract *SimpleReadAccessControllerTransactor 
}


func NewSimpleReadAccessController(address common.Address, backend bind.ContractBackend) (*SimpleReadAccessController, error) {
	contract, err := bindSimpleReadAccessController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessController{SimpleReadAccessControllerCaller: SimpleReadAccessControllerCaller{contract: contract}, SimpleReadAccessControllerTransactor: SimpleReadAccessControllerTransactor{contract: contract}, SimpleReadAccessControllerFilterer: SimpleReadAccessControllerFilterer{contract: contract}}, nil
}


func NewSimpleReadAccessControllerCaller(address common.Address, caller bind.ContractCaller) (*SimpleReadAccessControllerCaller, error) {
	contract, err := bindSimpleReadAccessController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCaller{contract: contract}, nil
}


func NewSimpleReadAccessControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleReadAccessControllerTransactor, error) {
	contract, err := bindSimpleReadAccessController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerTransactor{contract: contract}, nil
}


func NewSimpleReadAccessControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleReadAccessControllerFilterer, error) {
	contract, err := bindSimpleReadAccessController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerFilterer{contract: contract}, nil
}


func bindSimpleReadAccessController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleReadAccessControllerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerCaller.contract.Call(opts, result, method, params...)
}



func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerTransactor.contract.Transfer(opts)
}


func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerTransactor.contract.Transact(opts, method, params...)
}





func (_SimpleReadAccessController *SimpleReadAccessControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleReadAccessController.Contract.contract.Call(opts, result, method, params...)
}



func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.contract.Transfer(opts)
}


func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.contract.Transact(opts, method, params...)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) CheckEnabled() (bool, error) {
	return _SimpleReadAccessController.Contract.CheckEnabled(&_SimpleReadAccessController.CallOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) CheckEnabled() (bool, error) {
	return _SimpleReadAccessController.Contract.CheckEnabled(&_SimpleReadAccessController.CallOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _SimpleReadAccessController.Contract.HasAccess(&_SimpleReadAccessController.CallOpts, _user, _calldata)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _SimpleReadAccessController.Contract.HasAccess(&_SimpleReadAccessController.CallOpts, _user, _calldata)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) Owner() (common.Address, error) {
	return _SimpleReadAccessController.Contract.Owner(&_SimpleReadAccessController.CallOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) Owner() (common.Address, error) {
	return _SimpleReadAccessController.Contract.Owner(&_SimpleReadAccessController.CallOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "acceptOwnership")
}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AcceptOwnership(&_SimpleReadAccessController.TransactOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AcceptOwnership(&_SimpleReadAccessController.TransactOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "addAccess", _user)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AddAccess(&_SimpleReadAccessController.TransactOpts, _user)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AddAccess(&_SimpleReadAccessController.TransactOpts, _user)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "disableAccessCheck")
}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.DisableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.DisableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "enableAccessCheck")
}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.EnableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.EnableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "removeAccess", _user)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.RemoveAccess(&_SimpleReadAccessController.TransactOpts, _user)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.RemoveAccess(&_SimpleReadAccessController.TransactOpts, _user)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "transferOwnership", _to)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.TransferOwnership(&_SimpleReadAccessController.TransactOpts, _to)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.TransferOwnership(&_SimpleReadAccessController.TransactOpts, _to)
}


type SimpleReadAccessControllerAddedAccessIterator struct {
	Event *SimpleReadAccessControllerAddedAccess 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleReadAccessControllerAddedAccessIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerAddedAccess)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleReadAccessControllerAddedAccess)
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


func (it *SimpleReadAccessControllerAddedAccessIterator) Error() error {
	return it.fail
}



func (it *SimpleReadAccessControllerAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleReadAccessControllerAddedAccess struct {
	User common.Address
	Raw  types.Log 
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*SimpleReadAccessControllerAddedAccessIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerAddedAccessIterator{contract: _SimpleReadAccessController.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerAddedAccess) (event.Subscription, error) {

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleReadAccessControllerAddedAccess)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseAddedAccess(log types.Log) (*SimpleReadAccessControllerAddedAccess, error) {
	event := new(SimpleReadAccessControllerAddedAccess)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleReadAccessControllerCheckAccessDisabledIterator struct {
	Event *SimpleReadAccessControllerCheckAccessDisabled 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerCheckAccessDisabled)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleReadAccessControllerCheckAccessDisabled)
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


func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Error() error {
	return it.fail
}



func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleReadAccessControllerCheckAccessDisabled struct {
	Raw types.Log 
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*SimpleReadAccessControllerCheckAccessDisabledIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCheckAccessDisabledIterator{contract: _SimpleReadAccessController.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleReadAccessControllerCheckAccessDisabled)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseCheckAccessDisabled(log types.Log) (*SimpleReadAccessControllerCheckAccessDisabled, error) {
	event := new(SimpleReadAccessControllerCheckAccessDisabled)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleReadAccessControllerCheckAccessEnabledIterator struct {
	Event *SimpleReadAccessControllerCheckAccessEnabled 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerCheckAccessEnabled)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleReadAccessControllerCheckAccessEnabled)
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


func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Error() error {
	return it.fail
}



func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleReadAccessControllerCheckAccessEnabled struct {
	Raw types.Log 
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*SimpleReadAccessControllerCheckAccessEnabledIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCheckAccessEnabledIterator{contract: _SimpleReadAccessController.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleReadAccessControllerCheckAccessEnabled)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseCheckAccessEnabled(log types.Log) (*SimpleReadAccessControllerCheckAccessEnabled, error) {
	event := new(SimpleReadAccessControllerCheckAccessEnabled)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleReadAccessControllerOwnershipTransferRequestedIterator struct {
	Event *SimpleReadAccessControllerOwnershipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerOwnershipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleReadAccessControllerOwnershipTransferRequested)
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


func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleReadAccessControllerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleReadAccessControllerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerOwnershipTransferRequestedIterator{contract: _SimpleReadAccessController.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleReadAccessControllerOwnershipTransferRequested)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseOwnershipTransferRequested(log types.Log) (*SimpleReadAccessControllerOwnershipTransferRequested, error) {
	event := new(SimpleReadAccessControllerOwnershipTransferRequested)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleReadAccessControllerOwnershipTransferredIterator struct {
	Event *SimpleReadAccessControllerOwnershipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerOwnershipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleReadAccessControllerOwnershipTransferred)
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


func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Error() error {
	return it.fail
}



func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleReadAccessControllerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleReadAccessControllerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerOwnershipTransferredIterator{contract: _SimpleReadAccessController.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleReadAccessControllerOwnershipTransferred)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseOwnershipTransferred(log types.Log) (*SimpleReadAccessControllerOwnershipTransferred, error) {
	event := new(SimpleReadAccessControllerOwnershipTransferred)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleReadAccessControllerRemovedAccessIterator struct {
	Event *SimpleReadAccessControllerRemovedAccess 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleReadAccessControllerRemovedAccessIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleReadAccessControllerRemovedAccess)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleReadAccessControllerRemovedAccess)
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


func (it *SimpleReadAccessControllerRemovedAccessIterator) Error() error {
	return it.fail
}



func (it *SimpleReadAccessControllerRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleReadAccessControllerRemovedAccess struct {
	User common.Address
	Raw  types.Log 
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*SimpleReadAccessControllerRemovedAccessIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerRemovedAccessIterator{contract: _SimpleReadAccessController.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *SimpleReadAccessControllerRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _SimpleReadAccessController.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleReadAccessControllerRemovedAccess)
				if err := _SimpleReadAccessController.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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




func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseRemovedAccess(log types.Log) (*SimpleReadAccessControllerRemovedAccess, error) {
	event := new(SimpleReadAccessControllerRemovedAccess)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const SimpleWriteAccessControllerABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var SimpleWriteAccessControllerBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556001805460ff60a01b1916600160a01b179055610992806100456000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638823da6c11610076578063a118f2491161005b578063a118f249146101fd578063dc7f012414610230578063f2fde38b14610238576100a3565b80638823da6c146101995780638da5cb5b146101cc576100a3565b80630a756983146100a85780636b14daf8146100b257806379ba5097146101895780638038e4a114610191575b600080fd5b6100b061026b565b005b610175600480360360408110156100c857600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561010057600080fd5b82018360208201111561011257600080fd5b8035906020019184600183028401116401000000008311171561013457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610368945050505050565b604080519115158252519081900360200190f35b6100b06103be565b6100b06104c0565b6100b0600480360360208110156101af57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166105d2565b6101d461070a565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100b06004803603602081101561021357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610726565b6101756107b5565b6100b06004803603602081101561024e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166107d6565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff161561036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff16806103b7575060015474010000000000000000000000000000000000000000900460ff16155b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461044457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461054657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff1661036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60005473ffffffffffffffffffffffffffffffffffffffff16331461065857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff16156107075773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146107ac57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b610707816108d2565b60015474010000000000000000000000000000000000000000900460ff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461085c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff166107075773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a15056fea164736f6c6343000705000a"


func DeploySimpleWriteAccessController(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleWriteAccessController, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleWriteAccessControllerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SimpleWriteAccessControllerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleWriteAccessController{SimpleWriteAccessControllerCaller: SimpleWriteAccessControllerCaller{contract: contract}, SimpleWriteAccessControllerTransactor: SimpleWriteAccessControllerTransactor{contract: contract}, SimpleWriteAccessControllerFilterer: SimpleWriteAccessControllerFilterer{contract: contract}}, nil
}


type SimpleWriteAccessController struct {
	SimpleWriteAccessControllerCaller     
	SimpleWriteAccessControllerTransactor 
	SimpleWriteAccessControllerFilterer   
}


type SimpleWriteAccessControllerCaller struct {
	contract *bind.BoundContract 
}


type SimpleWriteAccessControllerTransactor struct {
	contract *bind.BoundContract 
}


type SimpleWriteAccessControllerFilterer struct {
	contract *bind.BoundContract 
}



type SimpleWriteAccessControllerSession struct {
	Contract     *SimpleWriteAccessController 
	CallOpts     bind.CallOpts                
	TransactOpts bind.TransactOpts            
}



type SimpleWriteAccessControllerCallerSession struct {
	Contract *SimpleWriteAccessControllerCaller 
	CallOpts bind.CallOpts                      
}



type SimpleWriteAccessControllerTransactorSession struct {
	Contract     *SimpleWriteAccessControllerTransactor 
	TransactOpts bind.TransactOpts                      
}


type SimpleWriteAccessControllerRaw struct {
	Contract *SimpleWriteAccessController 
}


type SimpleWriteAccessControllerCallerRaw struct {
	Contract *SimpleWriteAccessControllerCaller 
}


type SimpleWriteAccessControllerTransactorRaw struct {
	Contract *SimpleWriteAccessControllerTransactor 
}


func NewSimpleWriteAccessController(address common.Address, backend bind.ContractBackend) (*SimpleWriteAccessController, error) {
	contract, err := bindSimpleWriteAccessController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessController{SimpleWriteAccessControllerCaller: SimpleWriteAccessControllerCaller{contract: contract}, SimpleWriteAccessControllerTransactor: SimpleWriteAccessControllerTransactor{contract: contract}, SimpleWriteAccessControllerFilterer: SimpleWriteAccessControllerFilterer{contract: contract}}, nil
}


func NewSimpleWriteAccessControllerCaller(address common.Address, caller bind.ContractCaller) (*SimpleWriteAccessControllerCaller, error) {
	contract, err := bindSimpleWriteAccessController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerCaller{contract: contract}, nil
}


func NewSimpleWriteAccessControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleWriteAccessControllerTransactor, error) {
	contract, err := bindSimpleWriteAccessController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerTransactor{contract: contract}, nil
}


func NewSimpleWriteAccessControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleWriteAccessControllerFilterer, error) {
	contract, err := bindSimpleWriteAccessController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerFilterer{contract: contract}, nil
}


func bindSimpleWriteAccessController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleWriteAccessControllerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerCaller.contract.Call(opts, result, method, params...)
}



func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerTransactor.contract.Transfer(opts)
}


func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerTransactor.contract.Transact(opts, method, params...)
}





func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleWriteAccessController.Contract.contract.Call(opts, result, method, params...)
}



func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.contract.Transfer(opts)
}


func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.contract.Transact(opts, method, params...)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SimpleWriteAccessController.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) CheckEnabled() (bool, error) {
	return _SimpleWriteAccessController.Contract.CheckEnabled(&_SimpleWriteAccessController.CallOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerSession) CheckEnabled() (bool, error) {
	return _SimpleWriteAccessController.Contract.CheckEnabled(&_SimpleWriteAccessController.CallOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) HasAccess(opts *bind.CallOpts, _user common.Address, arg1 []byte) (bool, error) {
	var out []interface{}
	err := _SimpleWriteAccessController.contract.Call(opts, &out, "hasAccess", _user, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) HasAccess(_user common.Address, arg1 []byte) (bool, error) {
	return _SimpleWriteAccessController.Contract.HasAccess(&_SimpleWriteAccessController.CallOpts, _user, arg1)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerSession) HasAccess(_user common.Address, arg1 []byte) (bool, error) {
	return _SimpleWriteAccessController.Contract.HasAccess(&_SimpleWriteAccessController.CallOpts, _user, arg1)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleWriteAccessController.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) Owner() (common.Address, error) {
	return _SimpleWriteAccessController.Contract.Owner(&_SimpleWriteAccessController.CallOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerSession) Owner() (common.Address, error) {
	return _SimpleWriteAccessController.Contract.Owner(&_SimpleWriteAccessController.CallOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "acceptOwnership")
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.AcceptOwnership(&_SimpleWriteAccessController.TransactOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.AcceptOwnership(&_SimpleWriteAccessController.TransactOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "addAccess", _user)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.AddAccess(&_SimpleWriteAccessController.TransactOpts, _user)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.AddAccess(&_SimpleWriteAccessController.TransactOpts, _user)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "disableAccessCheck")
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.DisableAccessCheck(&_SimpleWriteAccessController.TransactOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.DisableAccessCheck(&_SimpleWriteAccessController.TransactOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "enableAccessCheck")
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.EnableAccessCheck(&_SimpleWriteAccessController.TransactOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.EnableAccessCheck(&_SimpleWriteAccessController.TransactOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "removeAccess", _user)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.RemoveAccess(&_SimpleWriteAccessController.TransactOpts, _user)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.RemoveAccess(&_SimpleWriteAccessController.TransactOpts, _user)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "transferOwnership", _to)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.TransferOwnership(&_SimpleWriteAccessController.TransactOpts, _to)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.TransferOwnership(&_SimpleWriteAccessController.TransactOpts, _to)
}


type SimpleWriteAccessControllerAddedAccessIterator struct {
	Event *SimpleWriteAccessControllerAddedAccess 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleWriteAccessControllerAddedAccessIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleWriteAccessControllerAddedAccess)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleWriteAccessControllerAddedAccess)
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


func (it *SimpleWriteAccessControllerAddedAccessIterator) Error() error {
	return it.fail
}



func (it *SimpleWriteAccessControllerAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleWriteAccessControllerAddedAccess struct {
	User common.Address
	Raw  types.Log 
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*SimpleWriteAccessControllerAddedAccessIterator, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerAddedAccessIterator{contract: _SimpleWriteAccessController.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *SimpleWriteAccessControllerAddedAccess) (event.Subscription, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleWriteAccessControllerAddedAccess)
				if err := _SimpleWriteAccessController.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseAddedAccess(log types.Log) (*SimpleWriteAccessControllerAddedAccess, error) {
	event := new(SimpleWriteAccessControllerAddedAccess)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleWriteAccessControllerCheckAccessDisabledIterator struct {
	Event *SimpleWriteAccessControllerCheckAccessDisabled 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleWriteAccessControllerCheckAccessDisabledIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleWriteAccessControllerCheckAccessDisabled)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleWriteAccessControllerCheckAccessDisabled)
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


func (it *SimpleWriteAccessControllerCheckAccessDisabledIterator) Error() error {
	return it.fail
}



func (it *SimpleWriteAccessControllerCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleWriteAccessControllerCheckAccessDisabled struct {
	Raw types.Log 
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*SimpleWriteAccessControllerCheckAccessDisabledIterator, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerCheckAccessDisabledIterator{contract: _SimpleWriteAccessController.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *SimpleWriteAccessControllerCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleWriteAccessControllerCheckAccessDisabled)
				if err := _SimpleWriteAccessController.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseCheckAccessDisabled(log types.Log) (*SimpleWriteAccessControllerCheckAccessDisabled, error) {
	event := new(SimpleWriteAccessControllerCheckAccessDisabled)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleWriteAccessControllerCheckAccessEnabledIterator struct {
	Event *SimpleWriteAccessControllerCheckAccessEnabled 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleWriteAccessControllerCheckAccessEnabledIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleWriteAccessControllerCheckAccessEnabled)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleWriteAccessControllerCheckAccessEnabled)
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


func (it *SimpleWriteAccessControllerCheckAccessEnabledIterator) Error() error {
	return it.fail
}



func (it *SimpleWriteAccessControllerCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleWriteAccessControllerCheckAccessEnabled struct {
	Raw types.Log 
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*SimpleWriteAccessControllerCheckAccessEnabledIterator, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerCheckAccessEnabledIterator{contract: _SimpleWriteAccessController.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *SimpleWriteAccessControllerCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleWriteAccessControllerCheckAccessEnabled)
				if err := _SimpleWriteAccessController.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseCheckAccessEnabled(log types.Log) (*SimpleWriteAccessControllerCheckAccessEnabled, error) {
	event := new(SimpleWriteAccessControllerCheckAccessEnabled)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleWriteAccessControllerOwnershipTransferRequestedIterator struct {
	Event *SimpleWriteAccessControllerOwnershipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleWriteAccessControllerOwnershipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleWriteAccessControllerOwnershipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleWriteAccessControllerOwnershipTransferRequested)
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


func (it *SimpleWriteAccessControllerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *SimpleWriteAccessControllerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleWriteAccessControllerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleWriteAccessControllerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerOwnershipTransferRequestedIterator{contract: _SimpleWriteAccessController.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *SimpleWriteAccessControllerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleWriteAccessController.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleWriteAccessControllerOwnershipTransferRequested)
				if err := _SimpleWriteAccessController.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseOwnershipTransferRequested(log types.Log) (*SimpleWriteAccessControllerOwnershipTransferRequested, error) {
	event := new(SimpleWriteAccessControllerOwnershipTransferRequested)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleWriteAccessControllerOwnershipTransferredIterator struct {
	Event *SimpleWriteAccessControllerOwnershipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleWriteAccessControllerOwnershipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleWriteAccessControllerOwnershipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleWriteAccessControllerOwnershipTransferred)
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


func (it *SimpleWriteAccessControllerOwnershipTransferredIterator) Error() error {
	return it.fail
}



func (it *SimpleWriteAccessControllerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleWriteAccessControllerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleWriteAccessControllerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerOwnershipTransferredIterator{contract: _SimpleWriteAccessController.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimpleWriteAccessControllerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleWriteAccessController.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleWriteAccessControllerOwnershipTransferred)
				if err := _SimpleWriteAccessController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseOwnershipTransferred(log types.Log) (*SimpleWriteAccessControllerOwnershipTransferred, error) {
	event := new(SimpleWriteAccessControllerOwnershipTransferred)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type SimpleWriteAccessControllerRemovedAccessIterator struct {
	Event *SimpleWriteAccessControllerRemovedAccess 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *SimpleWriteAccessControllerRemovedAccessIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleWriteAccessControllerRemovedAccess)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SimpleWriteAccessControllerRemovedAccess)
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


func (it *SimpleWriteAccessControllerRemovedAccessIterator) Error() error {
	return it.fail
}



func (it *SimpleWriteAccessControllerRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type SimpleWriteAccessControllerRemovedAccess struct {
	User common.Address
	Raw  types.Log 
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*SimpleWriteAccessControllerRemovedAccessIterator, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerRemovedAccessIterator{contract: _SimpleWriteAccessController.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *SimpleWriteAccessControllerRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SimpleWriteAccessControllerRemovedAccess)
				if err := _SimpleWriteAccessController.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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




func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseRemovedAccess(log types.Log) (*SimpleWriteAccessControllerRemovedAccess, error) {
	event := new(SimpleWriteAccessControllerRemovedAccess)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


const TestOffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"billingData\",\"outputs\":[{\"internalType\":\"uint16[31]\",\"name\":\"observationsCounts\",\"type\":\"uint16[31]\"},{\"internalType\":\"uint256[31]\",\"name\":\"gasReimbursements\",\"type\":\"uint256[31]\"},{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDigest\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testAccountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"testBurnLINK\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"testDecodeReport\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"int192[]\",\"name\":\"\",\"type\":\"int192[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes16\",\"name\":\"latestConfigDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint40\",\"name\":\"latestEpochAndRound\",\"type\":\"uint40\"},{\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"latestAggregatorRoundId\",\"type\":\"uint32\"}],\"internalType\":\"structOffchainAggregator.HotVars\",\"name\":\"h\",\"type\":\"tuple\"}],\"name\":\"testExposeHotvarsDummy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"txGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reasonableGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maximumGasPrice\",\"type\":\"uint256\"}],\"name\":\"testImpliedGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"testPayee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_x\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"_y\",\"type\":\"uint16\"}],\"name\":\"testSaturatingAddUint16\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitterOrSigner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amountLinkWei\",\"type\":\"uint256\"}],\"name\":\"testSetGasReimbursements\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_amount\",\"type\":\"uint16\"}],\"name\":\"testSetOracleObservationCount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testTotalLinkDue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"linkDue\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callDataCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"}],\"name\":\"testTransmitterGasCostEthWei\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var TestOffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b50604051620064e9380380620064e983398101604081905262000035916200060f565b604080518082019091526004815263151154d560e21b6020820152600080546001600160a01b031916331781558b918b918b918b918b918b918b918b918b918b91908b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b88620000998787878787620001f4565b620000a481620002e6565b6001600160601b0319606083901b16608052620000c062000444565b620000ca62000444565b60005b601f8160ff1610156200011a576001838260ff16601f8110620000ec57fe5b61ffff909216602092909202015260018260ff8316601f81106200010c57fe5b6020020152600101620000cd565b506200012a600483601f62000463565b506200013a600882601f62000500565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e05250508351620001839350602d925060208501915062000531565b506200018f866200035f565b8460170b60a08160170b60401b815250508360170b60c08160170b60401b815250505050505050505050505050506001602e60006101000a81548160ff02191690831515021790555050505050505050505050505050505050505050505050620006dd565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b0390811690821681146200035b57600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b03163314620003bf576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b0368010000000000000000909104811690821681146200035b57602c8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b604051806103e00160405280601f906020820280368337509192915050565b600283019183908215620004ee5791602002820160005b83821115620004bc57835183826101000a81548161ffff021916908361ffff16021790555092602001926002016020816001010492830192600103026200047a565b8015620004ec5782816101000a81549061ffff0219169055600201602081600101049283019260010302620004bc565b505b50620004fc929150620005b3565b5090565b82601f8101928215620004ee579160200282015b82811115620004ee57825182559160200191906001019062000514565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620005695760008555620004ee565b82601f106200058457805160ff1916838001178555620004ee565b82800160010185558215620004ee5791820182811115620004ee57825182559160200191906001019062000514565b5b80821115620004fc5760008155600101620005b4565b80516001600160a01b0381168114620005e257600080fd5b919050565b8051601781900b8114620005e257600080fd5b805163ffffffff81168114620005e257600080fd5b6000806000806000806000806000806101408b8d0312156200062f578586fd5b6200063a8b620005fa565b99506200064a60208c01620005fa565b98506200065a60408c01620005fa565b97506200066a60608c01620005fa565b96506200067a60808c01620005fa565b95506200068a60a08c01620005ca565b94506200069a60c08c01620005ca565b9350620006aa60e08c01620005e7565b9250620006bb6101008c01620005e7565b9150620006cc6101208c01620005ca565b90509295989b9194979a5092959850565b60805160601c60a05160401c60c05160401c60e05160f81c615da26200074760003980610b345250806117c152806132b5525080610a7b5280613288525080610a57528061206e528061260652806126f652806136ef52806141265280614a2d5250615da26000f3fe608060405234801561001057600080fd5b50600436106103575760003560e01c80638205bf6a116101c8578063bd82470611610104578063e4902f82116100a2578063f2fde38b1161007c578063f2fde38b146106e8578063fa98a1c7146106fb578063fbffd2c11461070e578063feaf968c1461072157610357565b8063e4902f82146106a9578063e5fe4577146106bc578063eb5dcd6c146106d557610357565b8063d09dc339116100de578063d09dc33914610673578063d18bf87e1461067b578063dc7f01241461068e578063e28519111461069657610357565b8063bd8247061461063a578063c10753291461064d578063c98075391461066057610357565b80639c849b3011610171578063acfe7f9c1161014b578063acfe7f9c146105ee578063b121e14714610601578063b5ab58dc14610614578063b633620c1461062757610357565b80639c849b30146105b55780639eb6e060146105c8578063a118f249146105db57610357565b80638da5cb5b116101a25780638da5cb5b146105695780639a6fc8f5146105715780639b764d971461059557610357565b80638205bf6a1461053b5780638823da6c146105435780638ac28d5a1461055657610357565b806350d25bcd1161029757806370da2f671161024057806379ba50971161021a57806379ba5097146104ff5780638038e4a114610507578063814118341461050f57806381ff70481461052457610357565b806370da2f67146104c25780637284e416146104ca57806377096177146104df57610357565b8063668a0f0211610271578063668a0f021461048557806366cfeaf11461048d5780636b14daf8146104a257610357565b806350d25bcd1461046257806354fd4d501461046a578063585aa7de1461047257610357565b806322adbc781161030457806335a2e492116102de57806335a2e4921461040a5780633a5381b51461041d5780633b5cdfa2146104255780633c04967b1461044757610357565b806322adbc78146103c757806329937268146103dc578063313ce567146103f557610357565b8063102a474b11610335578063102a474b146103975780631327d3d81461039f5780631b6b6d23146103b257610357565b80630a7569831461035c5780630b69df86146103665780630eafb25b14610384575b600080fd5b610364610729565b005b61036e6107e7565b60405161037b9190615c00565b60405180910390f35b61036e610392366004615593565b6107ed565b61036e61094d565b6103646103ad366004615593565b61095c565b6103ba610a55565b60405161037b91906159df565b6103cf610a79565b60405161037b9190615bf2565b6103e4610a9d565b60405161037b959493929190615d01565b6103fd610b32565b60405161037b9190615d63565b6103646104183660046158a7565b610b56565b6103ba610b59565b610438610433366004615874565b610b74565b60405161037b93929190615b9b565b61044f610b8f565b60405161037b9796959493929190615a59565b61036e610cd0565b61036e610d6c565b6103646104803660046156e6565b610d71565b61036e6116f2565b61049561178e565b60405161037b9190615b11565b6104b56104b03660046155df565b611797565b60405161037b9190615b06565b6103cf6117bf565b6104d26117e3565b60405161037b9190615c09565b6104f26104ed366004615916565b61187f565b60405161037b9190615c93565b610364611896565b610364611964565b610517611a23565b60405161037b9190615a0c565b61052c611a85565b60405161037b93929190615cbf565b61036e611aa6565b610364610551366004615593565b611b42565b610364610564366004615593565b611c38565b6103ba611caf565b61058461057f3660046159ab565b611cbe565b60405161037b959493929190615d30565b6105a86105a33660046158b8565b611d73565b60405161037b9190615cb0565b6103646105c336600461567d565b611d7f565b6103ba6105d6366004615593565b611fb8565b6103646105e9366004615593565b611fd6565b6103646105fc3660046158d3565b61203e565b61036461060f366004615593565b6120f8565b61036e6106223660046158d3565b6121f1565b61036e6106353660046158d3565b61228e565b610364610648366004615947565b61232b565b61036461065b366004615654565b6124ca565b61036461066e3660046157d5565b6127f9565b61036e6136ea565b61036461068936600461562b565b61379b565b6104b56137f2565b6103646106a4366004615654565b6137fb565b6105a86106b7366004615593565b613884565b6106c461393d565b60405161037b959493929190615b3e565b6103646106e33660046155ad565b613a2c565b6103646106f6366004615593565b613b88565b61036e6107093660046158eb565b613c50565b61036461071c366004615593565b613c65565b610584613ccd565b6000546001600160a01b03163314610788576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602e5460ff16156107e557602e80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b61179390565b60006107f76152bb565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561083b57fe5b600281111561084657fe5b905250905060008160200151600281111561085d57fe5b141561086d576000915050610948565b6108756152d2565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830181905270010000000000000000000000000000000090910490921660808201528251909160009160019060049060ff16601f811061090157fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff16601f811061093f57fe5b01540301925050505b919050565b6000610957613d80565b905090565b6000546001600160a01b031633146109bb576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b036801000000000000000090910481169082168114610a5157602c80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff16680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000610aad6152d2565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b7f000000000000000000000000000000000000000000000000000000000000000081565b50565b602c546801000000000000000090046001600160a01b031690565b6000806060610b8284613f61565b9250925092509193909250565b610b97615300565b610b9f615300565b6000806000806000610baf6152d2565b506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483168587018190526c0100000000000000000000000085048416606087018190527001000000000000000000000000000000009095049093166080860181905286516103e081019788905295966004966008969495939492918890601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411610c48575050604080516103e0810191829052959c508b9450601f93509150839050845b815481526020019060010190808311610c9e575050505050955097509750975097509750975097505090919293949596565b6000610d13336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179792505050565b610d64576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b610957614022565b600481565b868560ff8616601f831115610dcd576040805162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b60008111610e22576040805162461bcd60e51b815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b818314610e605760405162461bcd60e51b8152600401808060200182810382526024815260200180615d726024913960400191505060405180910390fd5b806003028311610eb7576040805162461bcd60e51b815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b6000546001600160a01b03163314610f16576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602854156110ba57602880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81019160009183908110610f5357fe5b6000918252602082200154602980546001600160a01b0390921693509084908110610f7a57fe5b6000918252602090912001546001600160a01b03169050610f9a8161405e565b6001600160a01b0380831660009081526027602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592841682529020805490911690556028805480610ff657fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055602980548061105957fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610f16915050565b60005b8a8110156114c8576000602760008e8e858181106110d757fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff16600281111561110d57fe5b1461115f576040805162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260016020820152602760008e8e8581811061118657fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561121157fe5b02179055506000915060069050818c8c8581811061122b57fe5b6001600160a01b03602091820293909301358316845283019390935260409091016000205416919091141590506112a9576040805162461bcd60e51b815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b6000602760008c8c858181106112bb57fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff1660028111156112f157fe5b14611343576040805162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260026020820152602760008c8c8581811061136a57fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101008360028111156113f557fe5b021790555090505060288c8c8381811061140b57fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03959092029390930135939093169290921790555060298a8a8381811061146d57fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03969093029490940135949094161790915550016110bd565b50602a805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602c80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000090911617938490559091048116911661159430828f8f8f8f8f8f8f8f61428e565b602a60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602a60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600083820152604051601f909101601f19169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b6000611735336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179792505050565b611786576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b610957614392565b602a5460801b90565b60006117a383836143b8565b806117b657506001600160a01b03831632145b90505b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6060611826336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179792505050565b611877576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6109576143e8565b600061188d85858585614493565b95945050505050565b6001546001600160a01b031633146118f5576040805162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b031633146119c3576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602e5460ff166107e557602e80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60606029805480602002602001604051908101604052809291908181526020018280548015611a7b57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611a5d575b5050505050905090565b602c54602a5463ffffffff808316926401000000009004169060801b909192565b6000611ae9336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179792505050565b611b3a576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61095761451f565b6000546001600160a01b03163314611ba1576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152602f602052604090205460ff1615610b56576001600160a01b0381166000818152602f602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a150565b6001600160a01b03818116600090815260066020526040902054163314611ca6576040805162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b610b568161405e565b6000546001600160a01b031681565b6000806000806000611d07336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179792505050565b611d58576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611d618661457a565b939a9299509097509550909350915050565b60006117b683836146ce565b6000546001600160a01b03163314611dde576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b828114611e32576040805162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015611fb1576000858583818110611e4b57fe5b905060200201356001600160a01b031690506000848484818110611e6b57fe5b6001600160a01b038581166000908152600660209081526040909120549202939093013583169350909116905080158080611eb75750826001600160a01b0316826001600160a01b0316145b611f08576040805162461bcd60e51b815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b03848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614611fa157826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505060019092019150611e359050565b5050505050565b6001600160a01b039081166000908152600660205260409020541690565b6000546001600160a01b03163314612035576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b610b56816146e6565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000081526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906120a69060019085906004016159f3565b602060405180830381600087803b1580156120c057600080fd5b505af11580156120d4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a5191906157ae565b6001600160a01b03818116600090815260076020526040902054163314612166576040805162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b6001600160a01b0381811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b6000612234336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179792505050565b612285576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6117b98261477f565b60006122d1336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179792505050565b612322576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6117b9826147b5565b6003546001600160a01b031680612389576040805162461bcd60e51b815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b6000546001600160a01b031633148061245c5750604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b03861694636b14daf8946000939190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b15801561242f57600080fd5b505afa158015612443573d6000803e3d6000fd5b505050506040513d602081101561245957600080fd5b50515b6124ad576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6124b561480a565b6124c28686868686614be6565b505050505050565b6000546001600160a01b03163314806125a55750600354604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b15801561257857600080fd5b505afa15801561258c573d6000803e3d6000fd5b505050506040513d60208110156125a257600080fd5b50515b6125f6576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6000612600613d80565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b15801561267157600080fd5b505afa158015612685573d6000803e3d6000fd5b505050506040513d602081101561269b57600080fd5b50519050818110156126f4576040805162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb8561273085850387614d60565b6040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b15801561277657600080fd5b505af115801561278a573d6000803e3d6000fd5b505050506040513d60208110156127a057600080fd5b50516127f3576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a905061280c888888888888614d77565b361461285f576040805162461bcd60e51b815260206004820152601960248201527f7472616e736d6974206d65737361676520746f6f206c6f6e6700000000000000604482015290519081900360640190fd5b61286761531f565b6040805160808082018352602a549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff1660208301527501000000000000000000000000000000000000000000810460ff169282019290925276010000000000000000000000000000000000000000000090910463ffffffff166060808301919091529082526000908a908a9081101561292057600080fd5b81359160208101359181019060608101604082013564010000000081111561294757600080fd5b82018360208201111561295957600080fd5b8035906020019184602083028401116401000000008311171561297b57600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050505060408801525050506080840182905283515190925060589190911b907fffffffffffffffffffffffffffffffff00000000000000000000000000000000808316911614612a42576040805162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b608083015183516020015164ffffffffff808316911610612aaa576040805162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b83516040015160ff168911612b06576040805162461bcd60e51b815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b601f891115612b5c576040805162461bcd60e51b815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b868914612bb0576040805162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b601f8460400151511115612c0b576040805162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b83600001516040015160020260ff1684604001515111612c72576040805162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b8867ffffffffffffffff81118015612c8957600080fd5b506040519080825280601f01601f191660200182016040528015612cb4576020820181803683370190505b50606085015260005b60ff81168a1115612d2557868160ff1660208110612cd757fe5b1a60f81b85606001518260ff1681518110612cee57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101612cbd565b5083604001515167ffffffffffffffff81118015612d4257600080fd5b506040519080825280601f01601f191660200182016040528015612d6d576020820181803683370190505b506020850152612d7b615300565b60005b8560400151518160ff161015612e81576000858260ff1660208110612d9f57fe5b1a90508281601f8110612dae57fe5b602002015115612e05576040805162461bcd60e51b815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b6001838260ff16601f8110612e1657fe5b91151560209283029190910152869060ff8416908110612e3257fe5b1a60f81b87602001518360ff1681518110612e4957fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101612d7e565b50612e8a6152bb565b336000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115612ec557fe5b6002811115612ed057fe5b9052509050600281602001516002811115612ee757fe5b148015612f1b57506029816000015160ff1681548110612f0357fe5b6000918252602090912001546001600160a01b031633145b612f6c576040805162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b5050835164ffffffffff90911660209091015250506040516000908a908a908083838082843760405192018290039091209450612fad935061530092505050565b612fb56152bb565b60005b898110156131ae57600060018587606001518481518110612fd557fe5b60209101015160f81c601b018e8e86818110612fed57fe5b905060200201358d8d8781811061300057fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801561305b573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526027602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156130aa57fe5b60028111156130b557fe5b90525092506001836020015160028111156130cc57fe5b1461311e576040805162461bcd60e51b815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff16601f811061312f57fe5b602002015115613186576040805162461bcd60e51b815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16601f811061319b57fe5b9115156020909202015250600101612fb8565b5050505060005b60018260400151510381101561325f576000826040015182600101815181106131da57fe5b602002602001015160170b836040015183815181106131f557fe5b602002602001015160170b1315905080613256576040805162461bcd60e51b815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b506001016131b5565b5060408101518051600091906002810490811061327857fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b131580156132de57507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b61332f576040805162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602b82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519586523386890181905291860181905260a0988601898152845199870199909952835194909916997ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451998c999298949793969095909492939185019260c086019289820192909102908190849084905b8381101561349257818101518382015260200161347a565b50505050905001838103825285818151815260200191508051906020019080838360005b838110156134ce5781810151838201526020016134b6565b50505050905090810190601f1680156134fb5780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390a281516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a36135b08260000151606001518260170b614d8f565b5080518051602a8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff928316021790915582106136d157fe5b6136df828260200151614ea0565b505050505050505050565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b15801561375a57600080fd5b505afa15801561376e573d6000803e3d6000fd5b505050506040513d602081101561378457600080fd5b505190506000613792613d80565b90910391505090565b6001600160a01b038216600090815260276020526040902054600182019060049060ff16601f81106137c957fe5b601091828204019190066002026101000a81548161ffff021916908361ffff1602179055505050565b602e5460ff1681565b60006001600160a01b038316600090815260276020526040902054610100900460ff16600281111561382957fe5b14156138505760405162461bcd60e51b815260040161384790615c5c565b60405180910390fd5b6001600160a01b038216600090815260276020526040902054600182019060089060ff16601f811061387e57fe5b01555050565b600061388e6152bb565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156138d257fe5b60028111156138dd57fe5b90525090506000816020015160028111156138f457fe5b1415613904576000915050610948565b60016004826000015160ff16601f811061391a57fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b600080808080333214613997576040805162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602a5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602b6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6001600160a01b03828116600090815260066020526040902054163314613a9a576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b336001600160a01b0382161415613af8576040805162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b6001600160a01b03808316600090815260076020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114613b83576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b6000546001600160a01b03163314613be7576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613c5d8484846150fb565b949350505050565b6000546001600160a01b03163314613cc4576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b610b5681615118565b6000806000806000613d16336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179792505050565b613d67576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b613d6f6151a7565b945094509450945094509091929394565b6000613d8a615300565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411613da35790505050505050905060005b601f811015613e135760018282601f8110613dfc57fe5b60200201510361ffff169290920191600101613de5565b50613e1c6152d2565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606080860182905270010000000000000000000000000000000090940490921660808501526029805486518184028101840190975280875297909202633b9aca0002969394929390830182828015613eec57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613ece575b50505050509050613efb615300565b604080516103e081019182905290600890601f9082845b815481526020019060010190808311613f12575050505050905060005b8251811015613f595760018282601f8110613f4657fe5b6020020151039590950194600101613f2f565b505050505090565b6000806060838060200190516060811015613f7b57600080fd5b81516020830151604080850180519151939592948301929184640100000000821115613fa657600080fd5b908301906020820185811115613fbb57600080fd5b8251866020820283011164010000000082111715613fd857600080fd5b82525081516020918201928201910280838360005b83811015614005578181015183820152602001613fed565b505050509190910160405250949993985091965091945050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b6020526040902054601790810b900b90565b6140666152bb565b6001600160a01b0382166000908152602760209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156140aa57fe5b60028111156140b557fe5b905250905060006140c5836107ed565b90508015613b83576001600160a01b0380841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b15801561416e57600080fd5b505af1158015614182573d6000803e3d6000fd5b505050506040513d602081101561419857600080fd5b50516141eb576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff16601f811061420157fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16601f811061423c57fe5b0155604080516001600160a01b0380871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b60008a8a8a8a8a8a8a8a8a8a604051602001808b6001600160a01b031681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff1690565b6001600160a01b0382166000908152602f602052604081205460ff16806117b6575050602e5460ff161592915050565b602d8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff610100600188161502019095169490940493840181900481028201810190925282815260609390929091830182828015611a7b5780601f1061446757610100808354040283529160200191611a7b565b820191906000526020600020905b81548152906001019060200180831161447557509395945050505050565b6000818510156144ea576040805162461bcd60e51b815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b818503830161179301633b9aca00858202026fffffffffffffffffffffffffffffffff811061451557fe5b9695505050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e740000000000000000000000000000000000815250906146555760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561461a578181015183820152602001614602565b50505050905090810190601f1680156146475780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5061465e6152bb565b5050505063ffffffff83166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b60006117b68261ffff168461ffff160161ffff614d60565b6001600160a01b0381166000908152602f602052604090205460ff16610b56576001600160a01b0381166000818152602f602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a150565b600063ffffffff82111561479557506000610948565b5063ffffffff166000908152602b6020526040902054601790810b900b90565b600063ffffffff8211156147cb57506000610948565b5063ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6148126152d2565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152614889615300565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116148a2579050505050505090506148e9615300565b604080516103e081019182905290600890601f9082845b81548152602001906001019080831161490057505050505090506060602980548060200260200160405190810160405280929190818152602001828054801561497257602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311614954575b5050505050905060005b8151811015614bca57600060018483601f811061499557fe5b6020020151039050600060018684601f81106149ad57fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115614bbf576000600660008787815181106149ed57fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b031690507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb82846040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b158015614aa257600080fd5b505af1158015614ab6573d6000803e3d6000fd5b505050506040513d6020811015614acc57600080fd5b5051614b1f576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018886601f8110614b2d57fe5b61ffff909216602092909202015260018786601f8110614b4957fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290879087908110614b7e57fe5b6020026020010151828460405180846001600160a01b03168152602001836001600160a01b03168152602001828152602001935050505060405180910390a1505b50505060010161497c565b50614bd8600484601f615353565b50611fb1600883601f6153e9565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b600081831015614d715750816117b9565b50919050565b602083810286019082020160e4019695505050505050565b602c546801000000000000000090046001600160a01b031680614db25750610a51565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602b602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b60248501819052948916604485015260648401889052516001600160a01b0387169363beed9b5193620186a09360848084019491939192918390030190829088803b158015614e7157600080fd5b5087f193505050508015614e9757506040513d6020811015614e9257600080fd5b505160015b6124c257611fb1565b614ea86152bb565b336000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115614ee357fe5b6002811115614eee57fe5b9052509050614efb6152d2565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000082048116606084015270010000000000000000000000000000000090910416608082015281516103e08101928390529091614fc6918591600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614f845790505050505050615246565b614fd490600490601f615353565b50600282602001516002811115614fe757fe5b14615039576040805162461bcd60e51b815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b6000615060633b9aca003a04836020015163ffffffff16846000015163ffffffff166150fb565b90506010360260005a9050600061507f8863ffffffff16858585614493565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff168302816150ad57fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff16601f81106150d657fe5b015401016008886000015160ff16601f81106150ee57fe5b0155505050505050505050565b6000838381101561510e57600285850304015b61188d8184614d60565b6003546001600160a01b039081169082168114610a5157600380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b602a54760100000000000000000000000000000000000000000000900463ffffffff1660008080806151d76152bb565b5050505063ffffffff82166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b61524e615300565b60005b83518110156152b357600084828151811061526857fe5b016020015160f81c905061528d8482601f811061528157fe5b602002015160016146ce565b848260ff16601f811061529c57fe5b61ffff909216602092909202015250600101615251565b509092915050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b604051806103e00160405280601f906020820280368337509192915050565b6040518060a00160405280615332615417565b81526060602082018190526040820181905280820152600060809091015290565b6002830191839082156153d95791602002820160005b838211156153a957835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302615369565b80156153d75782816101000a81549061ffff02191690556002016020816001010492830192600103026153a9565b505b506153e592915061543e565b5090565b82601f81019282156153d9579160200282015b828111156153d95782518255916020019190600101906153fc565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b808211156153e5576000815560010161543f565b80356001600160a01b038116811461094857600080fd5b60008083601f84011261547b578182fd5b50813567ffffffffffffffff811115615492578182fd5b60208301915083602080830285010111156154ac57600080fd5b9250929050565b60008083601f8401126154c4578182fd5b50813567ffffffffffffffff8111156154db578182fd5b6020830191508360208285010111156154ac57600080fd5b600082601f830112615503578081fd5b813567ffffffffffffffff8082111561551857fe5b6040516020601f19601f850116820101818110838211171561553657fe5b60405282815292508284830160200186101561555157600080fd5b8260208601602083013760006020848301015250505092915050565b803561ffff8116811461094857600080fd5b803563ffffffff8116811461094857600080fd5b6000602082840312156155a4578081fd5b6117b682615453565b600080604083850312156155bf578081fd5b6155c883615453565b91506155d660208401615453565b90509250929050565b600080604083850312156155f1578182fd5b6155fa83615453565b9150602083013567ffffffffffffffff811115615615578182fd5b615621858286016154f3565b9150509250929050565b6000806040838503121561563d578182fd5b61564683615453565b91506155d66020840161556d565b60008060408385031215615666578182fd5b61566f83615453565b946020939093013593505050565b60008060008060408587031215615692578182fd5b843567ffffffffffffffff808211156156a9578384fd5b6156b58883890161546a565b909650945060208701359150808211156156cd578384fd5b506156da8782880161546a565b95989497509550505050565b60008060008060008060008060a0898b031215615701578384fd5b883567ffffffffffffffff80821115615718578586fd5b6157248c838d0161546a565b909a50985060208b013591508082111561573c578586fd5b6157488c838d0161546a565b909850965060408b0135915060ff82168214615762578586fd5b90945060608a0135908082168214615778578485fd5b90935060808a0135908082111561578d578384fd5b5061579a8b828c016154b3565b999c989b5096995094979396929594505050565b6000602082840312156157bf578081fd5b815180151581146157ce578182fd5b9392505050565b60008060008060008060006080888a0312156157ef578283fd5b873567ffffffffffffffff80821115615806578485fd5b6158128b838c016154b3565b909950975060208a013591508082111561582a578485fd5b6158368b838c0161546a565b909750955060408a013591508082111561584e578485fd5b5061585b8a828b0161546a565b989b979a50959894979596606090950135949350505050565b600060208284031215615885578081fd5b813567ffffffffffffffff81111561589b578182fd5b613c5d848285016154f3565b600060808284031215614d71578081fd5b600080604083850312156158ca578182fd5b6156468361556d565b6000602082840312156158e4578081fd5b5035919050565b6000806000606084860312156158ff578081fd5b505081359360208301359350604090920135919050565b6000806000806080858703121561592b578182fd5b5050823594602084013594506040840135936060013592509050565b600080600080600060a0868803121561595e578283fd5b6159678661557f565b94506159756020870161557f565b93506159836040870161557f565b92506159916060870161557f565b915061599f6080870161557f565b90509295509295909350565b6000602082840312156159bc578081fd5b813569ffffffffffffffffffff811681146157ce578182fd5b63ffffffff169052565b6001600160a01b0391909116815260200190565b6001600160a01b03929092168252602082015260400190565b6020808252825182820181905260009190848201906040850190845b81811015615a4d5783516001600160a01b031683529284019291840191600101615a28565b50909695505050505050565b6108608101818960005b601f811015615a8657815161ffff16835260209283019290910190600101615a63565b5050506103e082018860005b601f811015615ab1578151835260209283019290910190600101615a92565b505050615ac26107c08301886159d5565b615ad06107e08301876159d5565b615ade6108008301866159d5565b615aec6108208301856159d5565b615afa6108408301846159d5565b98975050505050505050565b901515815260200190565b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000091909116815260200190565b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000095909516855263ffffffff93909316602085015260ff91909116604084015260170b606083015267ffffffffffffffff16608082015260a00190565b60006060820185835260208581850152606060408501528185518084526080860191508287019350845b81811015615be457845160170b83529383019391830191600101615bc5565b509098975050505050505050565b60179190910b815260200190565b90815260200190565b6000602080835283518082850152825b81811015615c3557858101830151858201604001528201615c19565b81811115615c465783604083870101525b50601f01601f1916929092016040019392505050565b6020808252600f908201527f6164647265737320756e6b6e6f776e0000000000000000000000000000000000604082015260600190565b6fffffffffffffffffffffffffffffffff91909116815260200190565b61ffff91909116815260200190565b63ffffffff93841681529190921660208201527fffffffffffffffffffffffffffffffff00000000000000000000000000000000909116604082015260600190565b63ffffffff95861681529385166020850152918416604084015283166060830152909116608082015260a00190565b69ffffffffffffffffffff9586168152602081019490945260408401929092526060830152909116608082015260a00190565b60ff9190911681526020019056fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea164736f6c6343000705000a"


func DeployTestOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _validator common.Address, _minAnswer *big.Int, _maxAnswer *big.Int, _billingAdminAccessController common.Address) (common.Address, *types.Transaction, *TestOffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(TestOffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TestOffchainAggregatorBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _validator, _minAnswer, _maxAnswer, _billingAdminAccessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestOffchainAggregator{TestOffchainAggregatorCaller: TestOffchainAggregatorCaller{contract: contract}, TestOffchainAggregatorTransactor: TestOffchainAggregatorTransactor{contract: contract}, TestOffchainAggregatorFilterer: TestOffchainAggregatorFilterer{contract: contract}}, nil
}


type TestOffchainAggregator struct {
	TestOffchainAggregatorCaller     
	TestOffchainAggregatorTransactor 
	TestOffchainAggregatorFilterer   
}


type TestOffchainAggregatorCaller struct {
	contract *bind.BoundContract 
}


type TestOffchainAggregatorTransactor struct {
	contract *bind.BoundContract 
}


type TestOffchainAggregatorFilterer struct {
	contract *bind.BoundContract 
}



type TestOffchainAggregatorSession struct {
	Contract     *TestOffchainAggregator 
	CallOpts     bind.CallOpts           
	TransactOpts bind.TransactOpts       
}



type TestOffchainAggregatorCallerSession struct {
	Contract *TestOffchainAggregatorCaller 
	CallOpts bind.CallOpts                 
}



type TestOffchainAggregatorTransactorSession struct {
	Contract     *TestOffchainAggregatorTransactor 
	TransactOpts bind.TransactOpts                 
}


type TestOffchainAggregatorRaw struct {
	Contract *TestOffchainAggregator 
}


type TestOffchainAggregatorCallerRaw struct {
	Contract *TestOffchainAggregatorCaller 
}


type TestOffchainAggregatorTransactorRaw struct {
	Contract *TestOffchainAggregatorTransactor 
}


func NewTestOffchainAggregator(address common.Address, backend bind.ContractBackend) (*TestOffchainAggregator, error) {
	contract, err := bindTestOffchainAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregator{TestOffchainAggregatorCaller: TestOffchainAggregatorCaller{contract: contract}, TestOffchainAggregatorTransactor: TestOffchainAggregatorTransactor{contract: contract}, TestOffchainAggregatorFilterer: TestOffchainAggregatorFilterer{contract: contract}}, nil
}


func NewTestOffchainAggregatorCaller(address common.Address, caller bind.ContractCaller) (*TestOffchainAggregatorCaller, error) {
	contract, err := bindTestOffchainAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorCaller{contract: contract}, nil
}


func NewTestOffchainAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*TestOffchainAggregatorTransactor, error) {
	contract, err := bindTestOffchainAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorTransactor{contract: contract}, nil
}


func NewTestOffchainAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*TestOffchainAggregatorFilterer, error) {
	contract, err := bindTestOffchainAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorFilterer{contract: contract}, nil
}


func bindTestOffchainAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestOffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}



func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorTransactor.contract.Transfer(opts)
}


func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}





func (_TestOffchainAggregator *TestOffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestOffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}



func (_TestOffchainAggregator *TestOffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.contract.Transfer(opts)
}


func (_TestOffchainAggregator *TestOffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.contract.Transact(opts, method, params...)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LINK() (common.Address, error) {
	return _TestOffchainAggregator.Contract.LINK(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _TestOffchainAggregator.Contract.LINK(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) BillingData(opts *bind.CallOpts) (struct {
	ObservationsCounts      [31]uint16
	GasReimbursements       [31]*big.Int
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "billingData")

	outstruct := new(struct {
		ObservationsCounts      [31]uint16
		GasReimbursements       [31]*big.Int
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ObservationsCounts = out[0].([31]uint16)
	outstruct.GasReimbursements = out[1].([31]*big.Int)
	outstruct.MaximumGasPrice = out[2].(uint32)
	outstruct.ReasonableGasPrice = out[3].(uint32)
	outstruct.MicroLinkPerEth = out[4].(uint32)
	outstruct.LinkGweiPerObservation = out[5].(uint32)
	outstruct.LinkGweiPerTransmission = out[6].(uint32)

	return *outstruct, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) BillingData() (struct {
	ObservationsCounts      [31]uint16
	GasReimbursements       [31]*big.Int
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _TestOffchainAggregator.Contract.BillingData(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) BillingData() (struct {
	ObservationsCounts      [31]uint16
	GasReimbursements       [31]*big.Int
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _TestOffchainAggregator.Contract.BillingData(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) CheckEnabled() (bool, error) {
	return _TestOffchainAggregator.Contract.CheckEnabled(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _TestOffchainAggregator.Contract.CheckEnabled(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Decimals() (uint8, error) {
	return _TestOffchainAggregator.Contract.Decimals(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _TestOffchainAggregator.Contract.Decimals(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Description() (string, error) {
	return _TestOffchainAggregator.Contract.Description(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Description() (string, error) {
	return _TestOffchainAggregator.Contract.Description(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetAnswer(&_TestOffchainAggregator.CallOpts, _roundId)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetAnswer(&_TestOffchainAggregator.CallOpts, _roundId)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetBilling(opts *bind.CallOpts) (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(struct {
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPrice = out[0].(uint32)
	outstruct.ReasonableGasPrice = out[1].(uint32)
	outstruct.MicroLinkPerEth = out[2].(uint32)
	outstruct.LinkGweiPerObservation = out[3].(uint32)
	outstruct.LinkGweiPerTransmission = out[4].(uint32)

	return *outstruct, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _TestOffchainAggregator.Contract.GetBilling(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _TestOffchainAggregator.Contract.GetBilling(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetConfigDigest(opts *bind.CallOpts) ([16]byte, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getConfigDigest")

	if err != nil {
		return *new([16]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([16]byte)).(*[16]byte)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetConfigDigest() ([16]byte, error) {
	return _TestOffchainAggregator.Contract.GetConfigDigest(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetConfigDigest() ([16]byte, error) {
	return _TestOffchainAggregator.Contract.GetConfigDigest(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _TestOffchainAggregator.Contract.GetRoundData(&_TestOffchainAggregator.CallOpts, _roundId)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _TestOffchainAggregator.Contract.GetRoundData(&_TestOffchainAggregator.CallOpts, _roundId)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetTimestamp(&_TestOffchainAggregator.CallOpts, _roundId)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetTimestamp(&_TestOffchainAggregator.CallOpts, _roundId)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _TestOffchainAggregator.Contract.HasAccess(&_TestOffchainAggregator.CallOpts, _user, _calldata)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _TestOffchainAggregator.Contract.HasAccess(&_TestOffchainAggregator.CallOpts, _user, _calldata)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [16]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = out[0].(uint32)
	outstruct.BlockNumber = out[1].(uint32)
	outstruct.ConfigDigest = out[2].([16]byte)

	return *outstruct, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _TestOffchainAggregator.Contract.LatestConfigDetails(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _TestOffchainAggregator.Contract.LatestConfigDetails(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestRound() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestRound(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestRound(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = out[0].(*big.Int)
	outstruct.Answer = out[1].(*big.Int)
	outstruct.StartedAt = out[2].(*big.Int)
	outstruct.UpdatedAt = out[3].(*big.Int)
	outstruct.AnsweredInRound = out[4].(*big.Int)

	return *outstruct, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _TestOffchainAggregator.Contract.LatestRoundData(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _TestOffchainAggregator.Contract.LatestRoundData(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestTimestamp(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestTimestamp(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestTransmissionDetails")

	outstruct := new(struct {
		ConfigDigest    [16]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigDigest = out[0].([16]byte)
	outstruct.Epoch = out[1].(uint32)
	outstruct.Round = out[2].(uint8)
	outstruct.LatestAnswer = out[3].(*big.Int)
	outstruct.LatestTimestamp = out[4].(uint64)

	return *outstruct, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _TestOffchainAggregator.Contract.LatestTransmissionDetails(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _TestOffchainAggregator.Contract.LatestTransmissionDetails(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LinkAvailableForPayment(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LinkAvailableForPayment(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MaxAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MaxAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MinAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MinAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _TestOffchainAggregator.Contract.OracleObservationCount(&_TestOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _TestOffchainAggregator.Contract.OracleObservationCount(&_TestOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.OwedPayment(&_TestOffchainAggregator.CallOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.OwedPayment(&_TestOffchainAggregator.CallOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Owner() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Owner(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Owner(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestAccountingGasCost(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testAccountingGasCost")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestAccountingGasCost() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestAccountingGasCost(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestAccountingGasCost() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestAccountingGasCost(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestDecodeReport(opts *bind.CallOpts, report []byte) ([32]byte, [32]byte, []*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testDecodeReport", report)

	if err != nil {
		return *new([32]byte), *new([32]byte), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	out2 := *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, out2, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestDecodeReport(report []byte) ([32]byte, [32]byte, []*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestDecodeReport(&_TestOffchainAggregator.CallOpts, report)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestDecodeReport(report []byte) ([32]byte, [32]byte, []*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestDecodeReport(&_TestOffchainAggregator.CallOpts, report)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestImpliedGasPrice(opts *bind.CallOpts, txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testImpliedGasPrice", txGasPrice, reasonableGasPrice, maximumGasPrice)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestImpliedGasPrice(txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestImpliedGasPrice(&_TestOffchainAggregator.CallOpts, txGasPrice, reasonableGasPrice, maximumGasPrice)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestImpliedGasPrice(txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestImpliedGasPrice(&_TestOffchainAggregator.CallOpts, txGasPrice, reasonableGasPrice, maximumGasPrice)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestPayee(opts *bind.CallOpts, _transmitter common.Address) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testPayee", _transmitter)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestPayee(_transmitter common.Address) (common.Address, error) {
	return _TestOffchainAggregator.Contract.TestPayee(&_TestOffchainAggregator.CallOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestPayee(_transmitter common.Address) (common.Address, error) {
	return _TestOffchainAggregator.Contract.TestPayee(&_TestOffchainAggregator.CallOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestSaturatingAddUint16(opts *bind.CallOpts, _x uint16, _y uint16) (uint16, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testSaturatingAddUint16", _x, _y)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestSaturatingAddUint16(_x uint16, _y uint16) (uint16, error) {
	return _TestOffchainAggregator.Contract.TestSaturatingAddUint16(&_TestOffchainAggregator.CallOpts, _x, _y)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestSaturatingAddUint16(_x uint16, _y uint16) (uint16, error) {
	return _TestOffchainAggregator.Contract.TestSaturatingAddUint16(&_TestOffchainAggregator.CallOpts, _x, _y)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestTotalLinkDue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testTotalLinkDue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestTotalLinkDue() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTotalLinkDue(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestTotalLinkDue() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTotalLinkDue(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestTransmitterGasCostEthWei(opts *bind.CallOpts, initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testTransmitterGasCostEthWei", initialGas, gasPrice, callDataCost, gasLeft)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestTransmitterGasCostEthWei(initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTransmitterGasCostEthWei(&_TestOffchainAggregator.CallOpts, initialGas, gasPrice, callDataCost, gasLeft)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestTransmitterGasCostEthWei(initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTransmitterGasCostEthWei(&_TestOffchainAggregator.CallOpts, initialGas, gasPrice, callDataCost, gasLeft)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _TestOffchainAggregator.Contract.Transmitters(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _TestOffchainAggregator.Contract.Transmitters(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "validator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Validator() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Validator(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Validator(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Version() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.Version(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Version() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.Version(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "acceptOwnership")
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AcceptOwnership(&_TestOffchainAggregator.TransactOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AcceptOwnership(&_TestOffchainAggregator.TransactOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "acceptPayeeship", _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AcceptPayeeship(&_TestOffchainAggregator.TransactOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AcceptPayeeship(&_TestOffchainAggregator.TransactOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "addAccess", _user)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AddAccess(&_TestOffchainAggregator.TransactOpts, _user)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AddAccess(&_TestOffchainAggregator.TransactOpts, _user)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "disableAccessCheck")
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.DisableAccessCheck(&_TestOffchainAggregator.TransactOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.DisableAccessCheck(&_TestOffchainAggregator.TransactOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "enableAccessCheck")
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.EnableAccessCheck(&_TestOffchainAggregator.TransactOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.EnableAccessCheck(&_TestOffchainAggregator.TransactOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "removeAccess", _user)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.RemoveAccess(&_TestOffchainAggregator.TransactOpts, _user)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.RemoveAccess(&_TestOffchainAggregator.TransactOpts, _user)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetBilling(&_TestOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetBilling(&_TestOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setBillingAccessController", _billingAdminAccessController)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetBillingAccessController(&_TestOffchainAggregator.TransactOpts, _billingAdminAccessController)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetBillingAccessController(&_TestOffchainAggregator.TransactOpts, _billingAdminAccessController)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetConfig(&_TestOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetConfig(&_TestOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setPayees", _transmitters, _payees)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetPayees(&_TestOffchainAggregator.TransactOpts, _transmitters, _payees)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetPayees(&_TestOffchainAggregator.TransactOpts, _transmitters, _payees)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetValidator(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setValidator", _newValidator)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetValidator(&_TestOffchainAggregator.TransactOpts, _newValidator)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetValidator(&_TestOffchainAggregator.TransactOpts, _newValidator)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TestBurnLINK(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "testBurnLINK", amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestBurnLINK(amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestBurnLINK(&_TestOffchainAggregator.TransactOpts, amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TestBurnLINK(amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestBurnLINK(&_TestOffchainAggregator.TransactOpts, amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TestExposeHotvarsDummy(opts *bind.TransactOpts, h OffchainAggregatorHotVars) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "testExposeHotvarsDummy", h)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestExposeHotvarsDummy(h OffchainAggregatorHotVars) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestExposeHotvarsDummy(&_TestOffchainAggregator.TransactOpts, h)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TestExposeHotvarsDummy(h OffchainAggregatorHotVars) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestExposeHotvarsDummy(&_TestOffchainAggregator.TransactOpts, h)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TestSetGasReimbursements(opts *bind.TransactOpts, _transmitterOrSigner common.Address, _amountLinkWei *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "testSetGasReimbursements", _transmitterOrSigner, _amountLinkWei)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestSetGasReimbursements(_transmitterOrSigner common.Address, _amountLinkWei *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestSetGasReimbursements(&_TestOffchainAggregator.TransactOpts, _transmitterOrSigner, _amountLinkWei)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TestSetGasReimbursements(_transmitterOrSigner common.Address, _amountLinkWei *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestSetGasReimbursements(&_TestOffchainAggregator.TransactOpts, _transmitterOrSigner, _amountLinkWei)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TestSetOracleObservationCount(opts *bind.TransactOpts, _oracle common.Address, _amount uint16) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "testSetOracleObservationCount", _oracle, _amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestSetOracleObservationCount(_oracle common.Address, _amount uint16) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestSetOracleObservationCount(&_TestOffchainAggregator.TransactOpts, _oracle, _amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TestSetOracleObservationCount(_oracle common.Address, _amount uint16) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestSetOracleObservationCount(&_TestOffchainAggregator.TransactOpts, _oracle, _amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "transferOwnership", _to)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TransferOwnership(&_TestOffchainAggregator.TransactOpts, _to)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TransferOwnership(&_TestOffchainAggregator.TransactOpts, _to)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TransferPayeeship(&_TestOffchainAggregator.TransactOpts, _transmitter, _proposed)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TransferPayeeship(&_TestOffchainAggregator.TransactOpts, _transmitter, _proposed)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) Transmit(opts *bind.TransactOpts, _report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "transmit", _report, _rs, _ss, _rawVs)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.Transmit(&_TestOffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.Transmit(&_TestOffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.WithdrawFunds(&_TestOffchainAggregator.TransactOpts, _recipient, _amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.WithdrawFunds(&_TestOffchainAggregator.TransactOpts, _recipient, _amount)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "withdrawPayment", _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.WithdrawPayment(&_TestOffchainAggregator.TransactOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.WithdrawPayment(&_TestOffchainAggregator.TransactOpts, _transmitter)
}


type TestOffchainAggregatorAddedAccessIterator struct {
	Event *TestOffchainAggregatorAddedAccess 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorAddedAccessIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorAddedAccess)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorAddedAccess)
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


func (it *TestOffchainAggregatorAddedAccessIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorAddedAccess struct {
	User common.Address
	Raw  types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*TestOffchainAggregatorAddedAccessIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorAddedAccessIterator{contract: _TestOffchainAggregator.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorAddedAccess) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorAddedAccess)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseAddedAccess(log types.Log) (*TestOffchainAggregatorAddedAccess, error) {
	event := new(TestOffchainAggregatorAddedAccess)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorAnswerUpdatedIterator struct {
	Event *TestOffchainAggregatorAnswerUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorAnswerUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorAnswerUpdated)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorAnswerUpdated)
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


func (it *TestOffchainAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*TestOffchainAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorAnswerUpdatedIterator{contract: _TestOffchainAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorAnswerUpdated)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*TestOffchainAggregatorAnswerUpdated, error) {
	event := new(TestOffchainAggregatorAnswerUpdated)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorBillingAccessControllerSetIterator struct {
	Event *TestOffchainAggregatorBillingAccessControllerSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorBillingAccessControllerSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorBillingAccessControllerSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorBillingAccessControllerSet)
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


func (it *TestOffchainAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*TestOffchainAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorBillingAccessControllerSetIterator{contract: _TestOffchainAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorBillingAccessControllerSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*TestOffchainAggregatorBillingAccessControllerSet, error) {
	event := new(TestOffchainAggregatorBillingAccessControllerSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorBillingSetIterator struct {
	Event *TestOffchainAggregatorBillingSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorBillingSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorBillingSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorBillingSet)
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


func (it *TestOffchainAggregatorBillingSetIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*TestOffchainAggregatorBillingSetIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorBillingSetIterator{contract: _TestOffchainAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorBillingSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseBillingSet(log types.Log) (*TestOffchainAggregatorBillingSet, error) {
	event := new(TestOffchainAggregatorBillingSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorCheckAccessDisabledIterator struct {
	Event *TestOffchainAggregatorCheckAccessDisabled 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorCheckAccessDisabledIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorCheckAccessDisabled)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorCheckAccessDisabled)
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


func (it *TestOffchainAggregatorCheckAccessDisabledIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorCheckAccessDisabled struct {
	Raw types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*TestOffchainAggregatorCheckAccessDisabledIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorCheckAccessDisabledIterator{contract: _TestOffchainAggregator.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorCheckAccessDisabled)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseCheckAccessDisabled(log types.Log) (*TestOffchainAggregatorCheckAccessDisabled, error) {
	event := new(TestOffchainAggregatorCheckAccessDisabled)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorCheckAccessEnabledIterator struct {
	Event *TestOffchainAggregatorCheckAccessEnabled 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorCheckAccessEnabledIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorCheckAccessEnabled)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorCheckAccessEnabled)
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


func (it *TestOffchainAggregatorCheckAccessEnabledIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorCheckAccessEnabled struct {
	Raw types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*TestOffchainAggregatorCheckAccessEnabledIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorCheckAccessEnabledIterator{contract: _TestOffchainAggregator.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorCheckAccessEnabled)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseCheckAccessEnabled(log types.Log) (*TestOffchainAggregatorCheckAccessEnabled, error) {
	event := new(TestOffchainAggregatorCheckAccessEnabled)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorConfigSetIterator struct {
	Event *TestOffchainAggregatorConfigSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorConfigSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorConfigSet)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorConfigSet)
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


func (it *TestOffchainAggregatorConfigSetIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*TestOffchainAggregatorConfigSetIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorConfigSetIterator{contract: _TestOffchainAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorConfigSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseConfigSet(log types.Log) (*TestOffchainAggregatorConfigSet, error) {
	event := new(TestOffchainAggregatorConfigSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorNewRoundIterator struct {
	Event *TestOffchainAggregatorNewRound 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorNewRoundIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorNewRound)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorNewRound)
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


func (it *TestOffchainAggregatorNewRoundIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*TestOffchainAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorNewRoundIterator{contract: _TestOffchainAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorNewRound)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseNewRound(log types.Log) (*TestOffchainAggregatorNewRound, error) {
	event := new(TestOffchainAggregatorNewRound)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorNewTransmissionIterator struct {
	Event *TestOffchainAggregatorNewTransmission 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorNewTransmissionIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorNewTransmission)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorNewTransmission)
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


func (it *TestOffchainAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	RawReportContext  [32]byte
	Raw               types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*TestOffchainAggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorNewTransmissionIterator{contract: _TestOffchainAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorNewTransmission)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseNewTransmission(log types.Log) (*TestOffchainAggregatorNewTransmission, error) {
	event := new(TestOffchainAggregatorNewTransmission)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorOraclePaidIterator struct {
	Event *TestOffchainAggregatorOraclePaid 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorOraclePaidIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorOraclePaid)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorOraclePaid)
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


func (it *TestOffchainAggregatorOraclePaidIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	Raw         types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts) (*TestOffchainAggregatorOraclePaidIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorOraclePaidIterator{contract: _TestOffchainAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOraclePaid) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorOraclePaid)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseOraclePaid(log types.Log) (*TestOffchainAggregatorOraclePaid, error) {
	event := new(TestOffchainAggregatorOraclePaid)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorOwnershipTransferRequestedIterator struct {
	Event *TestOffchainAggregatorOwnershipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorOwnershipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorOwnershipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorOwnershipTransferRequested)
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


func (it *TestOffchainAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestOffchainAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorOwnershipTransferRequestedIterator{contract: _TestOffchainAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorOwnershipTransferRequested)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*TestOffchainAggregatorOwnershipTransferRequested, error) {
	event := new(TestOffchainAggregatorOwnershipTransferRequested)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorOwnershipTransferredIterator struct {
	Event *TestOffchainAggregatorOwnershipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorOwnershipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorOwnershipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorOwnershipTransferred)
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


func (it *TestOffchainAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestOffchainAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorOwnershipTransferredIterator{contract: _TestOffchainAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorOwnershipTransferred)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*TestOffchainAggregatorOwnershipTransferred, error) {
	event := new(TestOffchainAggregatorOwnershipTransferred)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorPayeeshipTransferRequestedIterator struct {
	Event *TestOffchainAggregatorPayeeshipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorPayeeshipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorPayeeshipTransferRequested)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorPayeeshipTransferRequested)
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


func (it *TestOffchainAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*TestOffchainAggregatorPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorPayeeshipTransferRequestedIterator{contract: _TestOffchainAggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorPayeeshipTransferRequested)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*TestOffchainAggregatorPayeeshipTransferRequested, error) {
	event := new(TestOffchainAggregatorPayeeshipTransferRequested)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorPayeeshipTransferredIterator struct {
	Event *TestOffchainAggregatorPayeeshipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorPayeeshipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorPayeeshipTransferred)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorPayeeshipTransferred)
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


func (it *TestOffchainAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*TestOffchainAggregatorPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorPayeeshipTransferredIterator{contract: _TestOffchainAggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorPayeeshipTransferred)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*TestOffchainAggregatorPayeeshipTransferred, error) {
	event := new(TestOffchainAggregatorPayeeshipTransferred)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorRemovedAccessIterator struct {
	Event *TestOffchainAggregatorRemovedAccess 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorRemovedAccessIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorRemovedAccess)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorRemovedAccess)
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


func (it *TestOffchainAggregatorRemovedAccessIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorRemovedAccess struct {
	User common.Address
	Raw  types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*TestOffchainAggregatorRemovedAccessIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorRemovedAccessIterator{contract: _TestOffchainAggregator.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorRemovedAccess)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseRemovedAccess(log types.Log) (*TestOffchainAggregatorRemovedAccess, error) {
	event := new(TestOffchainAggregatorRemovedAccess)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}


type TestOffchainAggregatorValidatorUpdatedIterator struct {
	Event *TestOffchainAggregatorValidatorUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *TestOffchainAggregatorValidatorUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorValidatorUpdated)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(TestOffchainAggregatorValidatorUpdated)
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


func (it *TestOffchainAggregatorValidatorUpdatedIterator) Error() error {
	return it.fail
}



func (it *TestOffchainAggregatorValidatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type TestOffchainAggregatorValidatorUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log 
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterValidatorUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*TestOffchainAggregatorValidatorUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorValidatorUpdatedIterator{contract: _TestOffchainAggregator.contract, event: "ValidatorUpdated", logs: logs, sub: sub}, nil
}




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchValidatorUpdated(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorValidatorUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(TestOffchainAggregatorValidatorUpdated)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
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




func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseValidatorUpdated(log types.Log) (*TestOffchainAggregatorValidatorUpdated, error) {
	event := new(TestOffchainAggregatorValidatorUpdated)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
