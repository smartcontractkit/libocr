


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


const AccessControlTestHelperABI = "[{\"anonymous\":false,\"inputs\":[],\"name\":\"Dummy\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_roundID\",\"type\":\"uint256\"}],\"name\":\"readGetAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"},{\"internalType\":\"uint80\",\"name\":\"_roundID\",\"type\":\"uint80\"}],\"name\":\"readGetRoundData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_roundID\",\"type\":\"uint256\"}],\"name\":\"readGetTimestamp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"readLatestAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"readLatestRound\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"readLatestRoundData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"readLatestTimestamp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var AccessControlTestHelperBin = "0x608060405234801561001057600080fd5b50610559806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063c0c9c7db1161005b578063c0c9c7db14610135578063c9592ab914610168578063d2f79c47146101a1578063eea2913a146101d45761007d565b806320f2c97c1461008257806395319deb146100b7578063bf5fc18b146100fc575b600080fd5b6100b56004803603602081101561009857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610207565b005b6100b5600480360360408110156100cd57600080fd5b50803573ffffffffffffffffffffffffffffffffffffffff16906020013569ffffffffffffffffffff166102a5565b6100b56004803603604081101561011257600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020013561035b565b6100b56004803603602081101561014b57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166103d6565b6100b56004803603604081101561017e57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610446565b6100b5600480360360208110156101b757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610497565b6100b5600480360360208110156101ea57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166104dd565b8073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561024d57600080fd5b505afa158015610261573d6000803e3d6000fd5b505050506040513d60a081101561027757600080fd5b50506040517f10e4ab9f2ce395bf5539d7c60c9bfeef0b416602954734c5bb8bfd9433c9ff6890600090a150565b8173ffffffffffffffffffffffffffffffffffffffff16639a6fc8f5826040518263ffffffff1660e01b8152600401808269ffffffffffffffffffff16815260200191505060a06040518083038186803b15801561030257600080fd5b505afa158015610316573d6000803e3d6000fd5b505050506040513d60a081101561032c57600080fd5b50506040517f10e4ab9f2ce395bf5539d7c60c9bfeef0b416602954734c5bb8bfd9433c9ff6890600090a15050565b8173ffffffffffffffffffffffffffffffffffffffff1663b5ab58dc826040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b1580156103ac57600080fd5b505afa1580156103c0573d6000803e3d6000fd5b505050506040513d602081101561032c57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff166350d25bcd6040518163ffffffff1660e01b815260040160206040518083038186803b15801561041c57600080fd5b505afa158015610430573d6000803e3d6000fd5b505050506040513d602081101561027757600080fd5b8173ffffffffffffffffffffffffffffffffffffffff1663b633620c826040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b1580156103ac57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff16638205bf6a6040518163ffffffff1660e01b815260040160206040518083038186803b15801561041c57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff1663668a0f026040518163ffffffff1660e01b815260040160206040518083038186803b15801561041c57600080fdfea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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





func (_AccessControlTestHelper *AccessControlTestHelperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccessControlTestHelper.Contract.AccessControlTestHelperCaller.contract.Call(opts, result, method, params...)
}



func (_AccessControlTestHelper *AccessControlTestHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.AccessControlTestHelperTransactor.contract.Transfer(opts)
}


func (_AccessControlTestHelper *AccessControlTestHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.AccessControlTestHelperTransactor.contract.Transact(opts, method, params...)
}





func (_AccessControlTestHelper *AccessControlTestHelperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccessControlTestHelper.Contract.contract.Call(opts, result, method, params...)
}



func (_AccessControlTestHelper *AccessControlTestHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.contract.Transfer(opts)
}


func (_AccessControlTestHelper *AccessControlTestHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlTestHelper.Contract.contract.Transact(opts, method, params...)
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
	return event, nil
}


const AccessControlledOffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"availableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var AccessControlledOffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b50604051620062463803806200624683398181016040526101808110156200003957600080fd5b815160208301516040808501516060860151608087015160a088015160c089015160e08a01516101008b01516101208c01516101408d01516101608e0180519a519c9e9b9d999c989b979a969995989497939692959194939182019284640100000000821115620000a957600080fd5b908301906020820185811115620000bf57600080fd5b8251640100000000811182820188101715620000da57600080fd5b82525081516020918201929091019080838360005b8381101562000109578181015183820152602001620000ef565b50505050905090810190601f168015620001375780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b03191633179055508b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b88620001728787878787620002b1565b6200017d81620003a3565b6001600160601b0319606083901b166080526200019962000501565b620001a362000501565b60005b60208160ff161015620001f8576001838260ff1660208110620001c557fe5b602002019061ffff16908161ffff16815250506001828260ff1660208110620001ea57fe5b6020020152600101620001a6565b5062000208600483602062000520565b50620002186008826020620005bd565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e05250508351620002619350602e9250602085019150620005fc565b506200026d866200041c565b505050601791820b820b604090811b60a05290820b90910b901b60c0525050602f805460ff1916600117905550620006a19f50505050505050505050505050505050565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b0390811690821681146200041857600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b031633146200047c576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d546001600160a01b0368010000000000000000909104811690821681146200041857602d8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215620005ab5791602002820160005b838211156200057957835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030262000537565b8015620005a95782816101000a81549061ffff021916905560020160208160010104928301926001030262000579565b505b50620005b99291506200066e565b5090565b8260208101928215620005ee579160200282015b82811115620005ee578251825591602001919060010190620005d1565b50620005b99291506200068a565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200063f57805160ff1916838001178555620005ee565b82800160010185558215620005ee5791820182811115620005ee578251825591602001919060010190620005d1565b5b80821115620005b957805461ffff191681556001016200066f565b5b80821115620005b957600081556001016200068b565b60805160601c60a05160401c60c05160401c60e05160f81c615b4062000706600039806110f5525080611d2352806138d552508061103c52806138a85250806110185280612c025280612cc45280613d0f528061445a5280614d6e5250615b406000f3fe608060405234801561001057600080fd5b50600436106102d35760003560e01c80638205bf6a11610186578063bd824706116100e3578063e4902f8211610097578063f2fde38b11610071578063f2fde38b14610c4f578063fbffd2c114610c82578063feaf968c14610cb5576102d3565b8063e4902f8214610b60578063e5fe457714610baa578063eb5dcd6c14610c14576102d3565b8063c9807539116100c8578063c980753914610a3c578063d09dc33914610b50578063dc7f012414610b58576102d3565b8063bd824706146109be578063c107532914610a03576102d3565b80639c849b301161013a578063b121e1471161011f578063b121e14714610951578063b5ab58dc14610984578063b633620c146109a1576102d3565b80639c849b301461085c578063a118f2491461091e576102d3565b80638ac28d5a1161016b5780638ac28d5a146107ae5780638da5cb5b146107e15780639a6fc8f5146107e9576102d3565b80638205bf6a146107735780638823da6c1461077b576102d3565b8063585aa7de1161023457806373f666f8116101e85780638038e4a1116101cd5780638038e4a1146106c257806381411834146106ca57806381ff704814610722576102d3565b806373f666f8146106b257806379ba5097146106ba576102d3565b80636b14daf8116102195780636b14daf81461055657806370da2f671461062d5780637284e41614610635576102d3565b8063585aa7de14610421578063668a0f021461054e576102d3565b8063299372681161028b5780633a5381b5116102705780633a5381b51461040957806350d25bcd1461041157806354fd4d5014610419576102d3565b806329937268146103aa578063313ce567146103eb576102d3565b80631327d3d8116102bc5780631327d3d8146103275780631b6b6d231461035a57806322adbc781461038b576102d3565b80630a756983146102d85780630eafb25b146102e2575b600080fd5b6102e0610cbd565b005b610315600480360360208110156102f857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610d88565b60408051918252519081900360200190f35b6102e06004803603602081101561033d57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610ef6565b610362611016565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61039361103a565b6040805160179290920b8252519081900360200190f35b6103b261105e565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b6103f36110f3565b6040805160ff9092168252519081900360200190f35b610362611117565b610315611140565b6103156111e1565b6102e0600480360360a081101561043757600080fd5b81019060208101813564010000000081111561045257600080fd5b82018360208201111561046457600080fd5b8035906020019184602083028401116401000000008311171561048657600080fd5b9193909290916020810190356401000000008111156104a457600080fd5b8201836020820111156104b657600080fd5b803590602001918460208302840111640100000000831117156104d857600080fd5b9193909260ff8335169267ffffffffffffffff60208201351692919060608101906040013564010000000081111561050f57600080fd5b82018360208201111561052157600080fd5b8035906020019184600183028401116401000000008311171561054357600080fd5b5090925090506111e6565b610315611c50565b6106196004803603604081101561056c57600080fd5b73ffffffffffffffffffffffffffffffffffffffff82351691908101906040810160208201356401000000008111156105a457600080fd5b8201836020820111156105b657600080fd5b803590602001918460018302840111640100000000831117156105d857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611cec945050505050565b604080519115158252519081900360200190f35b610393611d21565b61063d611d45565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561067757818101518382015260200161065f565b50505050905090810190601f1680156106a45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610315611de1565b6102e0611de7565b6102e0611ecf565b6106d2611f9b565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561070e5781810151838201526020016106f6565b505050509050019250505060405180910390f35b61072a61200a565b6040805163ffffffff94851681529290931660208301527fffffffffffffffffffffffffffffffff00000000000000000000000000000000168183015290519081900360600190f35b61031561202b565b6102e06004803603602081101561079157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166120c7565b6102e0600480360360208110156107c457600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166121e5565b610362612269565b610812600480360360208110156107ff57600080fd5b503569ffffffffffffffffffff16612285565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102e06004803603604081101561087257600080fd5b81019060208101813564010000000081111561088d57600080fd5b82018360208201111561089f57600080fd5b803590602001918460208302840111640100000000831117156108c157600080fd5b9193909290916020810190356401000000008111156108df57600080fd5b8201836020820111156108f157600080fd5b8035906020019184602083028401116401000000008311171561091357600080fd5b50909250905061233a565b6102e06004803603602081101561093457600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166125e8565b6102e06004803603602081101561096757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661265d565b6103156004803603602081101561099a57600080fd5b5035612770565b610315600480360360208110156109b757600080fd5b503561280d565b6102e0600480360360a08110156109d457600080fd5b5063ffffffff8135811691602081013582169160408201358116916060810135821691608090910135166128aa565b6102e060048036036040811015610a1957600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135612a8e565b6102e060048036036080811015610a5257600080fd5b810190602081018135640100000000811115610a6d57600080fd5b820183602082011115610a7f57600080fd5b80359060200191846001830284011164010000000083111715610aa157600080fd5b919390929091602081019035640100000000811115610abf57600080fd5b820183602082011115610ad157600080fd5b80359060200191846020830284011164010000000083111715610af357600080fd5b919390929091602081019035640100000000811115610b1157600080fd5b820183602082011115610b2357600080fd5b80359060200191846020830284011164010000000083111715610b4557600080fd5b919350915035612de1565b610315613d0a565b610619613de9565b610b9360048036036020811015610b7657600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613df2565b6040805161ffff9092168252519081900360200190f35b610bb2613eb3565b604080517fffffffffffffffffffffffffffffffff00000000000000000000000000000000909616865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b6102e060048036036040811015610c2a57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516613fa2565b6102e060048036036020811015610c6557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16614132565b6102e060048036036020811015610c9857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16614214565b610812614289565b60005473ffffffffffffffffffffffffffffffffffffffff163314610d29576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602f5460ff1615610d8657602f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b6000610d9261590b565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602860209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115610de357fe5b6002811115610dee57fe5b9052509050600081602001516002811115610e0557fe5b1415610e15576000915050610ef1565b610e1d615922565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216948401949094526c0100000000000000000000000082048116606084018190527001000000000000000000000000000000009092041660808301528351919260009260019160049160ff16908110610eaa57fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff1660208110610ee857fe5b01540301925050505b919050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610f62576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d5473ffffffffffffffffffffffffffffffffffffffff680100000000000000009091048116908216811461101257602d80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff166801000000000000000073ffffffffffffffffffffffffffffffffffffffff85811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b600080600080600061106e615922565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b7f000000000000000000000000000000000000000000000000000000000000000081565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff165b90565b6000611183336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cec92505050565b6111d4576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6111dc61433c565b905090565b600481565b868560ff86166020831115611242576040805162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b60008111611297576040805162461bcd60e51b815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b8183146112d55760405162461bcd60e51b8152600401808060200182810382526024815260200180615ae76024913960400191505060405180910390fd5b80600302831161132c576040805162461bcd60e51b815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b60005473ffffffffffffffffffffffffffffffffffffffff163314611398576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6029541561156357602980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810191600091839081106113d557fe5b6000918252602082200154602a805473ffffffffffffffffffffffffffffffffffffffff9092169350908490811061140957fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16905061143681614378565b73ffffffffffffffffffffffffffffffffffffffff80831660009081526028602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559284168252902080549091169055602980548061149f57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055602a80548061150257fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611398915050565b60005b8a8110156119cc576000602860008e8e8581811061158057fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081019190915260400160002054610100900460ff1660028111156115c357fe5b14611615576040805162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260016020820152602860008e8e8581811061163c57fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101008360028111156116d457fe5b02179055506000915060069050818c8c858181106116ee57fe5b73ffffffffffffffffffffffffffffffffffffffff60209182029390930135831684528301939093526040909101600020541691909114159050611779576040805162461bcd60e51b815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b6000602860008c8c8581811061178b57fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081019190915260400160002054610100900460ff1660028111156117ce57fe5b14611820576040805162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260026020820152602860008c8c8581811061184757fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101008360028111156118df57fe5b021790555090505060298c8c838181106118f557fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9590920293909301359390931692909217905550602a8a8a8381811061196457fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff96909302949094013594909416179091555001611566565b50602b805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602d80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000909116179384905590910481169116611a9830828f8f8f8f8f8f8f8f6145cf565b602b60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602b60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528681526020019050868680828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b6000611c93336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cec92505050565b611ce4576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6111dc61471c565b6000611cf88383614742565b80611d18575073ffffffffffffffffffffffffffffffffffffffff831632145b90505b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6060611d88336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cec92505050565b611dd9576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6111dc61477f565b6159d881565b60015473ffffffffffffffffffffffffffffffffffffffff163314611e53576040805162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314611f3b576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602f5460ff16610d8657602f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b6060602a80548060200260200160405190810160405280929190818152602001828054801561200057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611fd5575b5050505050905090565b602d54602b5463ffffffff808316926401000000009004169060801b909192565b600061206e336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cec92505050565b6120bf576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6111dc61482a565b60005473ffffffffffffffffffffffffffffffffffffffff163314612133576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526030602052604090205460ff16156121e25773ffffffffffffffffffffffffffffffffffffffff811660008181526030602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260066020526040902054163314612260576040805162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b6121e281614378565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008060006122ce336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cec92505050565b61231f576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61232886614885565b939a9299509097509550909350915050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146123a6576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b8281146123fa576040805162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b838110156125e157600085858381811061241357fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff169050600084848481811061244057fe5b73ffffffffffffffffffffffffffffffffffffffff85811660009081526006602090815260409091205492029390930135831693509091169050801580806124b357508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b612504576040805162461bcd60e51b815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016858316908117909155908316146125d1578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506123fd9050565b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612654576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6121e2816149d9565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152600760205260409020541633146126d8576040805162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60006127b3336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cec92505050565b612804576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611d1b82614a8c565b6000612850336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cec92505050565b6128a1576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611d1b82614ac2565b60035473ffffffffffffffffffffffffffffffffffffffff1680612915576040805162461bcd60e51b815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b1580156129d257600080fd5b505afa1580156129e6573d6000803e3d6000fd5b505050506040513d60208110156129fc57600080fd5b505180612a20575060005473ffffffffffffffffffffffffffffffffffffffff1633145b612a71576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b612a79614b17565b612a868686868686614f60565b505050505050565b600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015612b5357600080fd5b505afa158015612b67573d6000803e3d6000fd5b505050506040513d6020811015612b7d57600080fd5b505180612ba1575060005473ffffffffffffffffffffffffffffffffffffffff1633145b612bf2576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b6000612bfc6150da565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015612c8757600080fd5b505afa158015612c9b573d6000803e3d6000fd5b505050506040513d6020811015612cb157600080fd5b5051905081811015612cc257600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85612d0b858503876152c8565b6040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015612d5e57600080fd5b505af1158015612d72573d6000803e3d6000fd5b505050506040513d6020811015612d8857600080fd5b5051612ddb576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a9050612df48888888888886152df565b3614612e47576040805162461bcd60e51b815260206004820152601960248201527f7472616e736d6974206d65737361676520746f6f206c6f6e6700000000000000604482015290519081900360640190fd5b612e4f615950565b6040805160808082018352602b549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff1660208301527501000000000000000000000000000000000000000000810460ff169282019290925276010000000000000000000000000000000000000000000090910463ffffffff166060808301919091529082526000908a908a90811015612f0857600080fd5b813591602081013591810190606081016040820135640100000000811115612f2f57600080fd5b820183602082011115612f4157600080fd5b80359060200191846020830284011164010000000083111715612f6357600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050505060408801525050506080840182905283515190925060589190911b907fffffffffffffffffffffffffffffffff0000000000000000000000000000000080831691161461302a576040805162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b608083015183516020015164ffffffffff808316911610613092576040805162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b83516040015160ff1689116130ee576040805162461bcd60e51b815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b6020891115613144576040805162461bcd60e51b815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b868914613198576040805162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b602084604001515111156131f3576040805162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b83600001516040015160020260ff168460400151511161325a576040805162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b8867ffffffffffffffff8111801561327157600080fd5b506040519080825280601f01601f19166020018201604052801561329c576020820181803683370190505b50606085015260005b60ff81168a111561330d57868160ff16602081106132bf57fe5b1a60f81b85606001518260ff16815181106132d657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506001016132a5565b5083604001515167ffffffffffffffff8111801561332a57600080fd5b506040519080825280601f01601f191660200182016040528015613355576020820181803683370190505b506020850152613363615984565b60005b8560400151518160ff161015613469576000858260ff166020811061338757fe5b1a905082816020811061339657fe5b6020020151156133ed576040805162461bcd60e51b815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b6001838260ff16602081106133fe57fe5b91151560209283029190910152869060ff841690811061341a57fe5b1a60f81b87602001518360ff168151811061343157fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101613366565b5061347261590b565b336000908152602860209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156134ad57fe5b60028111156134b857fe5b90525090506002816020015160028111156134cf57fe5b1480156135105750602a816000015160ff16815481106134eb57fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b613561576040805162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b5050835164ffffffffff90911660209091015250506040516000908a908a9080838380828437604051920182900390912094506135a2935061598492505050565b6135aa61590b565b60005b898110156137ce576000600185876060015184815181106135ca57fe5b60209101015160f81c601b018e8e868181106135e257fe5b905060200201358d8d878181106135f557fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015613650573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526028602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156136ca57fe5b60028111156136d557fe5b90525092506001836020015160028111156136ec57fe5b1461373e576040805162461bcd60e51b815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff166020811061374f57fe5b6020020151156137a6576040805162461bcd60e51b815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16602081106137bb57fe5b91151560209092020152506001016135ad565b5050505060005b60018260400151510381101561387f576000826040015182600101815181106137fa57fe5b602002602001015160170b8360400151838151811061381557fe5b602002602001015160170b1315905080613876576040805162461bcd60e51b815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b506001016137d5565b5060408101518051600091906002810490811061389857fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b131580156138fe57507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b61394f576040805162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602c82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519586523386890181905291860181905260a0988601898152845199870199909952835194909916997ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451998c999298949793969095909492939185019260c086019289820192909102908190849084905b83811015613ab2578181015183820152602001613a9a565b50505050905001838103825285818151815260200191508051906020019080838360005b83811015613aee578181015183820152602001613ad6565b50505050905090810190601f168015613b1b5780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390a281516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a3613bd08260000151606001518260170b6152f7565b5080518051602b8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff92831602179091558210613cf157fe5b613cff828260200151615422565b505050505050505050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015613d9457600080fd5b505afa158015613da8573d6000803e3d6000fd5b505050506040513d6020811015613dbe57600080fd5b505190506000613dcc6150da565b9050818111613dde579003905061113d565b60009250505061113d565b602f5460ff1681565b6000613dfc61590b565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602860209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115613e4d57fe5b6002811115613e5857fe5b9052509050600081602001516002811115613e6f57fe5b1415613e7f576000915050610ef1565b805160049060ff1660208110613e9157fe5b601091828204019190066002029054906101000a900461ffff16915050919050565b600080808080333214613f0d576040805162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602b5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602c6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526006602052604090205416331461401d576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b3373ffffffffffffffffffffffffffffffffffffffff82161415614088576040805162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff00000000000000000000000000000000000000008216811790925590911690811461412d5760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461419e576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314614280576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6121e281615681565b60008060008060006142d2336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cec92505050565b614323576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61432b61572a565b945094509450945094509091929394565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c6020526040902054601790810b900b90565b61438061590b565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602860209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156143d157fe5b60028111156143dc57fe5b905250905060006143ec83610d88565b9050801561412d5773ffffffffffffffffffffffffffffffffffffffff80841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b1580156144a257600080fd5b505af11580156144b6573d6000803e3d6000fd5b505050506040513d60208110156144cc57600080fd5b505161451f576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff166020811061453557fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff166020811061457057fe5b01556040805173ffffffffffffffffffffffffffffffffffffffff80871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b60008a8a8a8a8a8a8a8a8a8a604051602001808b73ffffffffffffffffffffffffffffffffffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b602b54760100000000000000000000000000000000000000000000900463ffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff821660009081526030602052604081205460ff1680611d18575050602f5460ff161592915050565b602e8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156120005780601f106147fe57610100808354040283529160200191612000565b820191906000526020600020905b81548152906001019060200180831161480c57509395945050505050565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e740000000000000000000000000000000000815250906149605760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561492557818101518382015260200161490d565b50505050905090810190601f1680156149525780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5061496961590b565b5050505063ffffffff83166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b73ffffffffffffffffffffffffffffffffffffffff811660009081526030602052604090205460ff166121e25773ffffffffffffffffffffffffffffffffffffffff811660008181526030602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a150565b600063ffffffff821115614aa257506000610ef1565b5063ffffffff166000908152602c6020526040902054601790810b900b90565b600063ffffffff821115614ad857506000610ef1565b5063ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b614b1f615922565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152614b96615984565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614baf57905050505050509050614bf6615984565b604080516104008101918290529060089060209082845b815481526020019060010190808311614c0d57505050505090506060602a805480602002602001604051908101604052809291908181526020018280548015614c8c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311614c61575b5050505050905060005b8151811015614f445760006001848360208110614caf57fe5b602002015103905060006001868460208110614cc757fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115614f3957600060066000878781518110614d0757fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015614dfd57600080fd5b505af1158015614e11573d6000803e3d6000fd5b505050506040513d6020811015614e2757600080fd5b5051614e7a576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b6001888660208110614e8857fe5b602002019061ffff16908161ffff16815250506001878660208110614ea957fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290879087908110614ede57fe5b60200260200101518284604051808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a1505b505050600101614c96565b50614f5260048460206159a3565b506125e16008836020615a39565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b60006150e4615984565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116150fd5790505050505050905060005b602081101561516d57600182826020811061515657fe5b60200201510361ffff16929092019160010161513f565b50615176615922565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608086018290527001000000000000000000000000000000009094049092166080850152602a805486518184028101840190975280875297909202633b9aca000296939492939083018282801561525357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311615228575b50505050509050615262615984565b604080516104008101918290529060089060209082845b815481526020019060010190808311615279575050505050905060005b82518110156152c05760018282602081106152ad57fe5b6020020151039590950194600101615296565b505050505090565b6000818310156152d9575081611d1b565b50919050565b602083810286019082020160e4019695505050505050565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff16806153275750611012565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602c602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b602485018190529489166044850152606484018890525173ffffffffffffffffffffffffffffffffffffffff87169363beed9b5193620186a09360848084019491939192918390030190829088803b1580156153f357600080fd5b5087f19350505050801561541957506040513d602081101561541457600080fd5b505160015b612a86576125e1565b61542a61590b565b336000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561546557fe5b600281111561547057fe5b905250905061547d615922565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608501527001000000000000000000000000000000009092041660808301528251610400810193849052919261554c928692909160049190826000855b82829054906101000a900461ffff1661ffff168152602001906002019060208260010104928301926001038202915080841161550a57905050505050506157c9565b61555a9060049060206159a3565b5060028260200151600281111561556d57fe5b146155bf576040805162461bcd60e51b815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b60006155e6633b9aca003a04836020015163ffffffff16846000015163ffffffff16615841565b90506010360260005a905060006156058863ffffffff16858585615867565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff1683028161563357fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff166020811061565c57fe5b015401016008886000015160ff166020811061567457fe5b0155505050505050505050565b60035473ffffffffffffffffffffffffffffffffffffffff908116908216811461101257600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b602b54760100000000000000000000000000000000000000000000900463ffffffff16600080808061575a61590b565b5050505063ffffffff82166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b6157d1615984565b60005b83518110156158395760008482815181106157eb57fe5b0160209081015160f81c91506158139085908390811061580757fe5b602002015160016158f3565b848260ff166020811061582257fe5b61ffff9092166020929092020152506001016157d4565b509092915050565b6000838381101561585457600285850304015b61585e81846152c8565b95945050505050565b6000818510156158be576040805162461bcd60e51b815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b81850383016159d801633b9aca00858202026fffffffffffffffffffffffffffffffff81106158e957fe5b9695505050505050565b6000611d188261ffff168461ffff160161ffff6152c8565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518060a00160405280615963615a73565b81526060602082018190526040820181905280820152600060809091015290565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215615a295791602002820160005b838211156159f957835183826101000a81548161ffff021916908361ffff16021790555092602001926002016020816001010492830192600103026159b9565b8015615a275782816101000a81549061ffff02191690556002016020816001010492830192600103026159f9565b505b50615a35929150615a9a565b5090565b8260208101928215615a67579160200282015b82811115615a67578251825591602001919060010190615a4c565b50615a35929150615ad1565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b80821115615a355780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000168155600101615a9b565b5b80821115615a355760008155600101615ad256fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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





func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}



func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorTransactor.contract.Transfer(opts)
}


func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}





func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccessControlledOffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}



func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.contract.Transfer(opts)
}


func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.contract.Transact(opts, method, params...)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "LINK")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LINK() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.LINK(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.LINK(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) AccountingGasCost(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "accountingGasCost")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AccountingGasCost() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.AccountingGasCost(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) AccountingGasCost() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.AccountingGasCost(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "checkEnabled")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) CheckEnabled() (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.CheckEnabled(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.CheckEnabled(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "decimals")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Decimals() (uint8, error) {
	return _AccessControlledOffchainAggregator.Contract.Decimals(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _AccessControlledOffchainAggregator.Contract.Decimals(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "description")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Description() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.Description(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Description() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.Description(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "getAnswer", _roundId)
	return *ret0, err
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
	ret := new(struct {
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	out := ret
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "getBilling")
	return *ret, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "getRoundData", _roundId)
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "getTimestamp", _roundId)
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetTimestamp(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetTimestamp(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "hasAccess", _user, _calldata)
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.HasAccess(&_AccessControlledOffchainAggregator.CallOpts, _user, _calldata)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.HasAccess(&_AccessControlledOffchainAggregator.CallOpts, _user, _calldata)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
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
	ret := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [16]byte
	})
	out := ret
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "latestConfigDetails")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "latestRound")
	return *ret0, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "latestRoundData")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
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
	ret := new(struct {
		ConfigDigest    [16]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	out := ret
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "latestTransmissionDetails")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "linkAvailableForPayment")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LinkAvailableForPayment(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LinkAvailableForPayment(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "maxAnswer")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MaxAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MaxAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "minAnswer")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MinAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MinAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "oracleObservationCount", _signerOrTransmitter)
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _AccessControlledOffchainAggregator.Contract.OracleObservationCount(&_AccessControlledOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _AccessControlledOffchainAggregator.Contract.OracleObservationCount(&_AccessControlledOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "owedPayment", _transmitter)
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.OwedPayment(&_AccessControlledOffchainAggregator.CallOpts, _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.OwedPayment(&_AccessControlledOffchainAggregator.CallOpts, _transmitter)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "owner")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Owner() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Owner(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Owner(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "transmitters")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmitters(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmitters(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "validator")
	return *ret0, err
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Validator() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Validator(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Validator(&_AccessControlledOffchainAggregator.CallOpts)
}




func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AccessControlledOffchainAggregator.contract.Call(opts, out, "version")
	return *ret0, err
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





func (_AccessControllerInterface *AccessControllerInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AccessControllerInterface *AccessControllerInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceTransactor.contract.Transfer(opts)
}


func (_AccessControllerInterface *AccessControllerInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AccessControllerInterface *AccessControllerInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccessControllerInterface.Contract.contract.Call(opts, result, method, params...)
}



func (_AccessControllerInterface *AccessControllerInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.contract.Transfer(opts)
}


func (_AccessControllerInterface *AccessControllerInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.contract.Transact(opts, method, params...)
}




func (_AccessControllerInterface *AccessControllerInterfaceCaller) HasAccess(opts *bind.CallOpts, user common.Address, data []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _AccessControllerInterface.contract.Call(opts, out, "hasAccess", user, data)
	return *ret0, err
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





func (_AggregatorInterface *AggregatorInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AggregatorInterface.Contract.AggregatorInterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AggregatorInterface *AggregatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.AggregatorInterfaceTransactor.contract.Transfer(opts)
}


func (_AggregatorInterface *AggregatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.AggregatorInterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AggregatorInterface *AggregatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AggregatorInterface.Contract.contract.Call(opts, result, method, params...)
}



func (_AggregatorInterface *AggregatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.contract.Transfer(opts)
}


func (_AggregatorInterface *AggregatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.contract.Transact(opts, method, params...)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "getAnswer", roundId)
	return *ret0, err
}




func (_AggregatorInterface *AggregatorInterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "getTimestamp", roundId)
	return *ret0, err
}




func (_AggregatorInterface *AggregatorInterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
}




func (_AggregatorInterface *AggregatorInterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "latestRound")
	return *ret0, err
}




func (_AggregatorInterface *AggregatorInterfaceSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}




func (_AggregatorInterface *AggregatorInterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorInterface.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
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





func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceTransactor.contract.Transfer(opts)
}


func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AggregatorV2V3Interface.Contract.contract.Call(opts, result, method, params...)
}



func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.contract.Transfer(opts)
}


func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.contract.Transact(opts, method, params...)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "decimals")
	return *ret0, err
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Decimals() (uint8, error) {
	return _AggregatorV2V3Interface.Contract.Decimals(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorV2V3Interface.Contract.Decimals(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "description")
	return *ret0, err
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Description() (string, error) {
	return _AggregatorV2V3Interface.Contract.Description(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Description() (string, error) {
	return _AggregatorV2V3Interface.Contract.Description(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "getAnswer", roundId)
	return *ret0, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "getRoundData", _roundId)
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "getTimestamp", roundId)
	return *ret0, err
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetTimestamp(&_AggregatorV2V3Interface.CallOpts, roundId)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetTimestamp(&_AggregatorV2V3Interface.CallOpts, roundId)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestAnswer(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestAnswer(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "latestRound")
	return *ret0, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "latestRoundData")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestTimestamp(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestTimestamp(&_AggregatorV2V3Interface.CallOpts)
}




func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorV2V3Interface.contract.Call(opts, out, "version")
	return *ret0, err
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





func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceTransactor.contract.Transfer(opts)
}


func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AggregatorV3Interface *AggregatorV3InterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AggregatorV3Interface.Contract.contract.Call(opts, result, method, params...)
}



func (_AggregatorV3Interface *AggregatorV3InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.contract.Transfer(opts)
}


func (_AggregatorV3Interface *AggregatorV3InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.contract.Transact(opts, method, params...)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _AggregatorV3Interface.contract.Call(opts, out, "decimals")
	return *ret0, err
}




func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Decimals() (uint8, error) {
	return _AggregatorV3Interface.Contract.Decimals(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorV3Interface.Contract.Decimals(&_AggregatorV3Interface.CallOpts)
}




func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AggregatorV3Interface.contract.Call(opts, out, "description")
	return *ret0, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _AggregatorV3Interface.contract.Call(opts, out, "getRoundData", _roundId)
	return *ret, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _AggregatorV3Interface.contract.Call(opts, out, "latestRoundData")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AggregatorV3Interface.contract.Call(opts, out, "version")
	return *ret0, err
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





func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transfer(opts)
}


func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
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





func (_LinkTokenInterface *LinkTokenInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceCaller.contract.Call(opts, result, method, params...)
}



func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transfer(opts)
}


func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transact(opts, method, params...)
}





func (_LinkTokenInterface *LinkTokenInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.contract.Call(opts, result, method, params...)
}



func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transfer(opts)
}


func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transact(opts, method, params...)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "decimals")
	return *ret0, err
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "name")
	return *ret0, err
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "symbol")
	return *ret0, err
}




func (_LinkTokenInterface *LinkTokenInterfaceSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}




func (_LinkTokenInterface *LinkTokenInterfaceCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LinkTokenInterface.contract.Call(opts, out, "totalSupply")
	return *ret0, err
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


const OffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"availableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var OffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b5060405162005b3e38038062005b3e83398181016040526101808110156200003957600080fd5b815160208301516040808501516060860151608087015160a088015160c089015160e08a01516101008b01516101208c01516101408d01516101608e0180519a519c9e9b9d999c989b979a969995989497939692959194939182019284640100000000821115620000a957600080fd5b908301906020820185811115620000bf57600080fd5b8251640100000000811182820188101715620000da57600080fd5b82525081516020918201929091019080838360005b8381101562000109578181015183820152602001620000ef565b50505050905090810190601f168015620001375780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b03191633179055508b8b8b8b8b8b886200016687878787876200028c565b62000171816200037e565b6001600160601b0319606083901b166080526200018d620004dc565b62000197620004dc565b60005b60208160ff161015620001ec576001838260ff1660208110620001b957fe5b602002019061ffff16908161ffff16815250506001828260ff1660208110620001de57fe5b60200201526001016200019a565b50620001fc6004836020620004fb565b506200020c600882602062000598565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e05250508351620002559350602e9250602085019150620005d7565b506200026186620003f7565b505050601791820b820b604090811b60a05290820b90910b901b60c052506200067c95505050505050565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b039081169082168114620003f357600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b0316331462000457576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d546001600160a01b036801000000000000000090910481169082168114620003f357602d8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215620005865791602002820160005b838211156200055457835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030262000512565b8015620005845782816101000a81549061ffff021916905560020160208160010104928301926001030262000554565b505b506200059492915062000649565b5090565b8260208101928215620005c9579160200282015b82811115620005c9578251825591602001919060010190620005ac565b506200059492915062000665565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200061a57805160ff1916838001178555620005c9565b82800160010185558215620005c95791820182811115620005c9578251825591602001919060010190620005ac565b5b808211156200059457805461ffff191681556001016200064a565b5b8082111562000594576000815560010162000666565b60805160601c60a05160401c60c05160401c60e05160f81c61545d620006e160003980610e5d525080611a4b52806135d7525080610da452806135aa525080610d805280612798528061285a5280613a2b528061419f52806146d2525061545d6000f3fe608060405234801561001057600080fd5b50600436106102415760003560e01c80638205bf6a11610145578063c1075329116100bd578063e5fe45771161008c578063f2fde38b11610071578063f2fde38b14610a68578063fbffd2c114610a9b578063feaf968c14610ace57610241565b8063e5fe4577146109c3578063eb5dcd6c14610a2d57610241565b8063c107532914610824578063c98075391461085d578063d09dc33914610971578063e4902f821461097957610241565b80639c849b3011610114578063b5ab58dc116100f9578063b5ab58dc146107a5578063b633620c146107c2578063bd824706146107df57610241565b80639c849b30146106b0578063b121e1471461077257610241565b80638205bf6a146105fa5780638ac28d5a146106025780638da5cb5b146106355780639a6fc8f51461063d57610241565b806354fd4d50116101d85780637284e416116101a757806379ba50971161018c57806379ba509714610549578063814118341461055157806381ff7048146105a957610241565b80637284e416146104c457806373f666f81461054157610241565b806354fd4d501461037f578063585aa7de14610387578063668a0f02146104b457806370da2f67146104bc57610241565b806329937268116102145780632993726814610310578063313ce567146103515780633a5381b51461036f57806350d25bcd1461037757610241565b80630eafb25b146102465780631327d3d81461028b5780631b6b6d23146102c057806322adbc78146102f1575b600080fd5b6102796004803603602081101561025c57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610ad6565b60408051918252519081900360200190f35b6102be600480360360208110156102a157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610c44565b005b6102c8610d7e565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6102f9610da2565b6040805160179290920b8252519081900360200190f35b610318610dc6565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b610359610e5b565b6040805160ff9092168252519081900360200190f35b6102c8610e7f565b610279610ea8565b610279610ee4565b6102be600480360360a081101561039d57600080fd5b8101906020810181356401000000008111156103b857600080fd5b8201836020820111156103ca57600080fd5b803590602001918460208302840111640100000000831117156103ec57600080fd5b91939092909160208101903564010000000081111561040a57600080fd5b82018360208201111561041c57600080fd5b8035906020019184602083028401116401000000008311171561043e57600080fd5b9193909260ff8335169267ffffffffffffffff60208201351692919060608101906040013564010000000081111561047557600080fd5b82018360208201111561048757600080fd5b803590602001918460018302840111640100000000831117156104a957600080fd5b509092509050610ee9565b610279611a23565b6102f9611a49565b6104cc611a6d565b6040805160208082528351818301528351919283929083019185019080838360005b838110156105065781810151838201526020016104ee565b50505050905090810190601f1680156105335780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610279611b21565b6102be611b27565b610559611c29565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561059557818101518382015260200161057d565b505050509050019250505060405180910390f35b6105b1611c97565b6040805163ffffffff94851681529290931660208301527fffffffffffffffffffffffffffffffff00000000000000000000000000000000168183015290519081900360600190f35b610279611cb8565b6102be6004803603602081101561061857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611d13565b6102c8611db4565b6106666004803603602081101561065357600080fd5b503569ffffffffffffffffffff16611dd0565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102be600480360360408110156106c657600080fd5b8101906020810181356401000000008111156106e157600080fd5b8201836020820111156106f357600080fd5b8035906020019184602083028401116401000000008311171561071557600080fd5b91939092909160208101903564010000000081111561073357600080fd5b82018360208201111561074557600080fd5b8035906020019184602083028401116401000000008311171561076757600080fd5b509092509050611f3e565b6102be6004803603602081101561078857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661223a565b610279600480360360208110156107bb57600080fd5b5035612367565b610279600480360360208110156107d857600080fd5b503561239d565b6102be600480360360a08110156107f557600080fd5b5063ffffffff8135811691602081013582169160408201358116916060810135821691608090910135166123f2565b6102be6004803603604081101561083a57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020013561260a565b6102be6004803603608081101561087357600080fd5b81019060208101813564010000000081111561088e57600080fd5b8201836020820111156108a057600080fd5b803590602001918460018302840111640100000000831117156108c257600080fd5b9193909290916020810190356401000000008111156108e057600080fd5b8201836020820111156108f257600080fd5b8035906020019184602083028401116401000000008311171561091457600080fd5b91939092909160208101903564010000000081111561093257600080fd5b82018360208201111561094457600080fd5b8035906020019184602083028401116401000000008311171561096657600080fd5b919350915035612991565b610279613a26565b6109ac6004803603602081101561098f57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613b05565b6040805161ffff9092168252519081900360200190f35b6109cb613bc6565b604080517fffffffffffffffffffffffffffffffff00000000000000000000000000000000909616865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b6102be60048036036040811015610a4357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516613ccf565b6102be60048036036020811015610a7e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613e93565b6102be60048036036020811015610ab157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613f8f565b61066661401e565b6000610ae0615228565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602860209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115610b3157fe5b6002811115610b3c57fe5b9052509050600081602001516002811115610b5357fe5b1415610b63576000915050610c3f565b610b6b61523f565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216948401949094526c0100000000000000000000000082048116606084018190527001000000000000000000000000000000009092041660808301528351919260009260019160049160ff16908110610bf857fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff1660208110610c3657fe5b01540301925050505b919050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610cca57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d5473ffffffffffffffffffffffffffffffffffffffff6801000000000000000090910481169082168114610d7a57602d80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff166801000000000000000073ffffffffffffffffffffffffffffffffffffffff85811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000610dd661523f565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b7f000000000000000000000000000000000000000000000000000000000000000081565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff165b90565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c6020526040902054601790810b900b90565b600481565b868560ff86166020831115610f5f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b60008111610fce57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b818314611026576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260248152602001806154046024913960400191505060405180910390fd5b80600302831161109757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b60005473ffffffffffffffffffffffffffffffffffffffff16331461111d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602954156112e857602980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101916000918390811061115a57fe5b6000918252602082200154602a805473ffffffffffffffffffffffffffffffffffffffff9092169350908490811061118e57fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1690506111bb816140bd565b73ffffffffffffffffffffffffffffffffffffffff80831660009081526028602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559284168252902080549091169055602980548061122457fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055602a80548061128757fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061111d915050565b60005b8a81101561179f576000602860008e8e8581811061130557fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081019190915260400160002054610100900460ff16600281111561134857fe5b146113b457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260016020820152602860008e8e858181106113db57fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561147357fe5b02179055506000915060069050818c8c8581811061148d57fe5b73ffffffffffffffffffffffffffffffffffffffff6020918202939093013583168452830193909352604090910160002054169190911415905061153257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b6000602860008c8c8581811061154457fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081019190915260400160002054610100900460ff16600281111561158757fe5b146115f357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260026020820152602860008c8c8581811061161a57fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff1683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101008360028111156116b257fe5b021790555090505060298c8c838181106116c857fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9590920293909301359390931692909217905550602a8a8a8381811061173757fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff969093029490940135949094161790915550016112eb565b50602b805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602d80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000090911617938490559091048116911661186b30828f8f8f8f8f8f8f8f61432e565b602b60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602b60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528681526020019050868680828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b602b54760100000000000000000000000000000000000000000000900463ffffffff1690565b7f000000000000000000000000000000000000000000000000000000000000000081565b602e8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff610100600188161502019095169490940493840181900481028201810190925282815260609390929091830182828015611b175780601f10611aec57610100808354040283529160200191611b17565b820191906000526020600020905b815481529060010190602001808311611afa57829003601f168201915b5050505050905090565b6159d881565b60015473ffffffffffffffffffffffffffffffffffffffff163314611bad57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060602a805480602002602001604051908101604052809291908181526020018280548015611b1757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611c63575050505050905090565b602d54602b5463ffffffff808316926401000000009004169060801b909192565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260066020526040902054163314611da857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b611db1816140bd565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e74000000000000000000000000000000000081525090611ec5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015611e8a578181015183820152602001611e72565b50505050905090810190601f168015611eb75780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50611ece615228565b5050505063ffffffff83166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b60005473ffffffffffffffffffffffffffffffffffffffff163314611fc457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b82811461203257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b8381101561223357600085858381811061204b57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff169050600084848481811061207857fe5b73ffffffffffffffffffffffffffffffffffffffff85811660009081526006602090815260409091205492029390930135831693509091169050801580806120eb57508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b61215657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614612223578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506120359050565b5050505050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152600760205260409020541633146122cf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600063ffffffff82111561237d57506000610c3f565b5063ffffffff166000908152602c6020526040902054601790810b900b90565b600063ffffffff8211156123b357506000610c3f565b5063ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b60035473ffffffffffffffffffffffffffffffffffffffff168061247757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561253457600080fd5b505afa158015612548573d6000803e3d6000fd5b505050506040513d602081101561255e57600080fd5b505180612582575060005473ffffffffffffffffffffffffffffffffffffffff1633145b6125ed57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b6125f561447b565b61260286868686866148de565b505050505050565b600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b1580156126cf57600080fd5b505afa1580156126e3573d6000803e3d6000fd5b505050506040513d60208110156126f957600080fd5b50518061271d575060005473ffffffffffffffffffffffffffffffffffffffff1633145b61278857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b6000612792614a58565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561281d57600080fd5b505afa158015612831573d6000803e3d6000fd5b505050506040513d602081101561284757600080fd5b505190508181101561285857600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb856128a185850387614c46565b6040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156128f457600080fd5b505af1158015612908573d6000803e3d6000fd5b505050506040513d602081101561291e57600080fd5b505161298b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a90506129a4888888888888614c60565b3614612a1157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f7472616e736d6974206d65737361676520746f6f206c6f6e6700000000000000604482015290519081900360640190fd5b612a1961526d565b6040805160808082018352602b549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff1660208301527501000000000000000000000000000000000000000000810460ff169282019290925276010000000000000000000000000000000000000000000090910463ffffffff166060808301919091529082526000908a908a90811015612ad257600080fd5b813591602081013591810190606081016040820135640100000000811115612af957600080fd5b820183602082011115612b0b57600080fd5b80359060200191846020830284011164010000000083111715612b2d57600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050505060408801525050506080840182905283515190925060589190911b907fffffffffffffffffffffffffffffffff00000000000000000000000000000000808316911614612c0e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b608083015183516020015164ffffffffff808316911610612c9057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b83516040015160ff168911612d0657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b6020891115612d7657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b868914612de457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b60208460400151511115612e5957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b83600001516040015160020260ff1684604001515111612eda57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b8867ffffffffffffffff81118015612ef157600080fd5b506040519080825280601f01601f191660200182016040528015612f1c576020820181803683370190505b50606085015260005b60ff81168a1115612f8d57868160ff1660208110612f3f57fe5b1a60f81b85606001518260ff1681518110612f5657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101612f25565b5083604001515167ffffffffffffffff81118015612faa57600080fd5b506040519080825280601f01601f191660200182016040528015612fd5576020820181803683370190505b506020850152612fe36152a1565b60005b8560400151518160ff161015613103576000858260ff166020811061300757fe5b1a905082816020811061301657fe5b60200201511561308757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b6001838260ff166020811061309857fe5b91151560209283029190910152869060ff84169081106130b457fe5b1a60f81b87602001518360ff16815181106130cb57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101612fe6565b5061310c615228565b336000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561314757fe5b600281111561315257fe5b905250905060028160200151600281111561316957fe5b1480156131aa5750602a816000015160ff168154811061318557fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b61321557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b5050835164ffffffffff90911660209091015250506040516000908a908a90808383808284376040519201829003909120945061325693506152a192505050565b61325e615228565b60005b898110156134b65760006001858760600151848151811061327e57fe5b60209101015160f81c601b018e8e8681811061329657fe5b905060200201358d8d878181106132a957fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015613304573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526028602090815290849020838501909452835460ff8082168552929650929450840191610100900416600281111561337e57fe5b600281111561338957fe5b90525092506001836020015160028111156133a057fe5b1461340c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff166020811061341d57fe5b60200201511561348e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16602081106134a357fe5b9115156020909202015250600101613261565b5050505060005b600182604001515103811015613581576000826040015182600101815181106134e257fe5b602002602001015160170b836040015183815181106134fd57fe5b602002602001015160170b131590508061357857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b506001016134bd565b5060408101518051600091906002810490811061359a57fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b1315801561360057507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b61366b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602c82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519586523386890181905291860181905260a0988601898152845199870199909952835194909916997ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451998c999298949793969095909492939185019260c086019289820192909102908190849084905b838110156137ce5781810151838201526020016137b6565b50505050905001838103825285818151815260200191508051906020019080838360005b8381101561380a5781810151838201526020016137f2565b50505050905090810190601f1680156138375780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390a281516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a36138ec8260000151606001518260170b614c78565b5080518051602b8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff92831602179091558210613a0d57fe5b613a1b828260200151614da3565b505050505050505050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015613ab057600080fd5b505afa158015613ac4573d6000803e3d6000fd5b505050506040513d6020811015613ada57600080fd5b505190506000613ae8614a58565b9050818111613afa5790039050610ea5565b600092505050610ea5565b6000613b0f615228565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602860209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115613b6057fe5b6002811115613b6b57fe5b9052509050600081602001516002811115613b8257fe5b1415613b92576000915050610c3f565b805160049060ff1660208110613ba457fe5b601091828204019190066002029054906101000a900461ffff16915050919050565b600080808080333214613c3a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602b5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602c6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff828116600090815260066020526040902054163314613d6457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b3373ffffffffffffffffffffffffffffffffffffffff82161415613de957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114613e8e5760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314613f1957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461401557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b611db18161501c565b602b54760100000000000000000000000000000000000000000000900463ffffffff16600080808061404e615228565b5050505063ffffffff82166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b6140c5615228565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561411657fe5b600281111561412157fe5b9052509050600061413183610ad6565b90508015613e8e5773ffffffffffffffffffffffffffffffffffffffff80841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b1580156141e757600080fd5b505af11580156141fb573d6000803e3d6000fd5b505050506040513d602081101561421157600080fd5b505161427e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff166020811061429457fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16602081106142cf57fe5b01556040805173ffffffffffffffffffffffffffffffffffffffff80871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b60008a8a8a8a8a8a8a8a8a8a604051602001808b73ffffffffffffffffffffffffffffffffffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b61448361523f565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830152700100000000000000000000000000000000900490911660808201526144fa6152a1565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116145135790505050505050905061455a6152a1565b604080516104008101918290529060089060209082845b81548152602001906001019080831161457157505050505090506060602a8054806020026020016040519081016040528092919081815260200182805480156145f057602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116145c5575b5050505050905060005b81518110156148c2576000600184836020811061461357fe5b60200201510390506000600186846020811061462b57fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca000201905060008111156148b75760006006600087878151811061466b57fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561476157600080fd5b505af1158015614775573d6000803e3d6000fd5b505050506040513d602081101561478b57600080fd5b50516147f857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b600188866020811061480657fe5b602002019061ffff16908161ffff1681525050600187866020811061482757fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29087908790811061485c57fe5b60200260200101518284604051808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a1505b5050506001016145fa565b506148d060048460206152c0565b506122336008836020615356565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000614a626152a1565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614a7b5790505050505050905060005b6020811015614aeb576001828260208110614ad457fe5b60200201510361ffff169290920191600101614abd565b50614af461523f565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608086018290527001000000000000000000000000000000009094049092166080850152602a805486518184028101840190975280875297909202633b9aca0002969394929390830182828015614bd157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311614ba6575b50505050509050614be06152a1565b604080516104008101918290529060089060209082845b815481526020019060010190808311614bf7575050505050905060005b8251811015614c3e576001828260208110614c2b57fe5b6020020151039590950194600101614c14565b505050505090565b600081831015614c57575081614c5a565b50805b92915050565b602083810286019082020160e4019695505050505050565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1680614ca85750610d7a565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602c602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b602485018190529489166044850152606484018890525173ffffffffffffffffffffffffffffffffffffffff87169363beed9b5193620186a09360848084019491939192918390030190829088803b158015614d7457600080fd5b5087f193505050508015614d9a57506040513d6020811015614d9557600080fd5b505160015b61260257612233565b614dab615228565b336000908152602860209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115614de657fe5b6002811115614df157fe5b9052509050614dfe61523f565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606085015270010000000000000000000000000000000090920416608083015282516104008101938490529192614ecd928692909160049190826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614e8b57905050505050506150c5565b614edb9060049060206152c0565b50600282602001516002811115614eee57fe5b14614f5a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b6000614f81633b9aca003a04836020015163ffffffff16846000015163ffffffff1661513d565b90506010360260005a90506000614fa08863ffffffff16858585615163565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff16830281614fce57fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff1660208110614ff757fe5b015401016008886000015160ff166020811061500f57fe5b0155505050505050505050565b60035473ffffffffffffffffffffffffffffffffffffffff9081169082168114610d7a57600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b6150cd6152a1565b60005b83518110156151355760008482815181106150e757fe5b0160209081015160f81c915061510f9085908390811061510357fe5b60200201516001615209565b848260ff166020811061511e57fe5b61ffff9092166020929092020152506001016150d0565b509092915050565b6000838381101561515057600285850304015b61515a8184614c46565b95945050505050565b6000818510156151d457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b81850383016159d801633b9aca00858202026fffffffffffffffffffffffffffffffff81106151ff57fe5b9695505050505050565b60006152218261ffff168461ffff160161ffff614c46565b9392505050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518060a00160405280615280615390565b81526060602082018190526040820181905280820152600060809091015290565b6040518061040001604052806020906020820280368337509192915050565b6002830191839082156153465791602002820160005b8382111561531657835183826101000a81548161ffff021916908361ffff16021790555092602001926002016020816001010492830192600103026152d6565b80156153445782816101000a81549061ffff0219169055600201602081600101049283019260010302615316565b505b506153529291506153b7565b5090565b8260208101928215615384579160200282015b82811115615384578251825591602001919060010190615369565b506153529291506153ee565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b808211156153525780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001681556001016153b8565b5b8082111561535257600081556001016153ef56fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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





func (_OffchainAggregator *OffchainAggregatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OffchainAggregator.Contract.OffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}



func (_OffchainAggregator *OffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.OffchainAggregatorTransactor.contract.Transfer(opts)
}


func (_OffchainAggregator *OffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.OffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}





func (_OffchainAggregator *OffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}



func (_OffchainAggregator *OffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.contract.Transfer(opts)
}


func (_OffchainAggregator *OffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.contract.Transact(opts, method, params...)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "LINK")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) LINK() (common.Address, error) {
	return _OffchainAggregator.Contract.LINK(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _OffchainAggregator.Contract.LINK(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) AccountingGasCost(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "accountingGasCost")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) AccountingGasCost() (*big.Int, error) {
	return _OffchainAggregator.Contract.AccountingGasCost(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) AccountingGasCost() (*big.Int, error) {
	return _OffchainAggregator.Contract.AccountingGasCost(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "decimals")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) Decimals() (uint8, error) {
	return _OffchainAggregator.Contract.Decimals(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _OffchainAggregator.Contract.Decimals(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "description")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) Description() (string, error) {
	return _OffchainAggregator.Contract.Description(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Description() (string, error) {
	return _OffchainAggregator.Contract.Description(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "getAnswer", _roundId)
	return *ret0, err
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
	ret := new(struct {
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	out := ret
	err := _OffchainAggregator.contract.Call(opts, out, "getBilling")
	return *ret, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _OffchainAggregator.contract.Call(opts, out, "getRoundData", _roundId)
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "getTimestamp", _roundId)
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetTimestamp(&_OffchainAggregator.CallOpts, _roundId)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetTimestamp(&_OffchainAggregator.CallOpts, _roundId)
}




func (_OffchainAggregator *OffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
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
	ret := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [16]byte
	})
	out := ret
	err := _OffchainAggregator.contract.Call(opts, out, "latestConfigDetails")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "latestRound")
	return *ret0, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _OffchainAggregator.contract.Call(opts, out, "latestRoundData")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
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
	ret := new(struct {
		ConfigDigest    [16]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	out := ret
	err := _OffchainAggregator.contract.Call(opts, out, "latestTransmissionDetails")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "linkAvailableForPayment")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregator.Contract.LinkAvailableForPayment(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregator.Contract.LinkAvailableForPayment(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "maxAnswer")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MaxAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MaxAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "minAnswer")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MinAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MinAnswer(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "oracleObservationCount", _signerOrTransmitter)
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregator.Contract.OracleObservationCount(&_OffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregator.Contract.OracleObservationCount(&_OffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_OffchainAggregator *OffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "owedPayment", _transmitter)
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregator.Contract.OwedPayment(&_OffchainAggregator.CallOpts, _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregator.Contract.OwedPayment(&_OffchainAggregator.CallOpts, _transmitter)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "owner")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) Owner() (common.Address, error) {
	return _OffchainAggregator.Contract.Owner(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _OffchainAggregator.Contract.Owner(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "transmitters")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _OffchainAggregator.Contract.Transmitters(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _OffchainAggregator.Contract.Transmitters(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "validator")
	return *ret0, err
}




func (_OffchainAggregator *OffchainAggregatorSession) Validator() (common.Address, error) {
	return _OffchainAggregator.Contract.Validator(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _OffchainAggregator.Contract.Validator(&_OffchainAggregator.CallOpts)
}




func (_OffchainAggregator *OffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregator.contract.Call(opts, out, "version")
	return *ret0, err
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
	return event, nil
}


const OffchainAggregatorBillingABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"availableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var OffchainAggregatorBillingBin = "0x60a06040523480156200001157600080fd5b50604051620027d2380380620027d2833981810160405260e08110156200003757600080fd5b508051602082015160408301516060840151608085015160a086015160c090960151600080546001600160a01b0319163317905594959394929391929091906200008587878787876200013b565b62000090816200022d565b6001600160601b0319606083901b16608052620000ac620002a6565b620000b6620002a6565b60005b60208160ff1610156200010b576001838260ff1660208110620000d857fe5b602002019061ffff16908161ffff16815250506001828260ff1660208110620000fd57fe5b6020020152600101620000b9565b506200011b6004836020620002c5565b506200012b600882602062000362565b50505050505050505050620003d4565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b039081169082168114620002a257600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215620003505791602002820160005b838211156200031e57835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302620002dc565b80156200034e5782816101000a81549061ffff02191690556002016020816001010492830192600103026200031e565b505b506200035e929150620003a1565b5090565b826020810192821562000393579160200282015b828111156200039357825182559160200191906001019062000376565b506200035e929150620003bd565b5b808211156200035e57805461ffff19168155600101620003a2565b5b808211156200035e5760008155600101620003be565b60805160601c6123c76200040b600039806105df528061102a52806110ec528061122852806117f35280611bd952506123c76000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063b121e14711610097578063e4902f8211610066578063e4902f8214610384578063eb5dcd6c146103ce578063f2fde38b14610409578063fbffd2c11461043c57610100565b8063b121e147146102cb578063bd824706146102fe578063c107532914610343578063d09dc3391461037c57610100565b806379ba5097116100d357806379ba5097146101c45780638ac28d5a146101ce5780638da5cb5b146102015780639c849b301461020957610100565b80630eafb25b146101055780631b6b6d231461014a578063299372681461017b57806373f666f8146101bc575b600080fd5b6101386004803603602081101561011b57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661046f565b60408051918252519081900360200190f35b6101526105dd565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b610183610601565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b610138610696565b6101cc61069c565b005b6101cc600480360360208110156101e457600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661079e565b61015261083f565b6101cc6004803603604081101561021f57600080fd5b81019060208101813564010000000081111561023a57600080fd5b82018360208201111561024c57600080fd5b8035906020019184602083028401116401000000008311171561026e57600080fd5b91939092909160208101903564010000000081111561028c57600080fd5b82018360208201111561029e57600080fd5b803590602001918460208302840111640100000000831117156102c057600080fd5b50909250905061085b565b6101cc600480360360208110156102e157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610b57565b6101cc600480360360a081101561031457600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516610c84565b6101cc6004803603604081101561035957600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610e9c565b610138611223565b6103b76004803603602081101561039a57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611301565b6040805161ffff9092168252519081900360200190f35b6101cc600480360360408110156103e457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160200135166113c2565b6101cc6004803603602081101561041f57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611586565b6101cc6004803603602081101561045257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611682565b6000610479612211565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602860209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156104ca57fe5b60028111156104d557fe5b90525090506000816020015160028111156104ec57fe5b14156104fc5760009150506105d8565b610504612228565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216948401949094526c0100000000000000000000000082048116606084018190527001000000000000000000000000000000009092041660808301528351919260009260019160049160ff1690811061059157fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff16602081106105cf57fe5b01540301925050505b919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000610611612228565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b6159d881565b60015473ffffffffffffffffffffffffffffffffffffffff16331461072257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526006602052604090205416331461083357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b61083c81611711565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b82811461094f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015610b5057600085858381811061096857fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff169050600084848481811061099557fe5b73ffffffffffffffffffffffffffffffffffffffff8581166000908152600660209081526040909120549202939093013583169350909116905080158080610a0857508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b610a7357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614610b40578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506109529050565b5050505050565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260076020526040902054163314610bec57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60035473ffffffffffffffffffffffffffffffffffffffff1680610d0957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610dc657600080fd5b505afa158015610dda573d6000803e3d6000fd5b505050506040513d6020811015610df057600080fd5b505180610e14575060005473ffffffffffffffffffffffffffffffffffffffff1633145b610e7f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b610e87611982565b610e948686868686611de5565b505050505050565b600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610f6157600080fd5b505afa158015610f75573d6000803e3d6000fd5b505050506040513d6020811015610f8b57600080fd5b505180610faf575060005473ffffffffffffffffffffffffffffffffffffffff1633145b61101a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b6000611024611f5f565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156110af57600080fd5b505afa1580156110c3573d6000803e3d6000fd5b505050506040513d60208110156110d957600080fd5b50519050818110156110ea57600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb856111338585038761214d565b6040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561118657600080fd5b505af115801561119a573d6000803e3d6000fd5b505050506040513d60208110156111b057600080fd5b505161121d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156112ad57600080fd5b505afa1580156112c1573d6000803e3d6000fd5b505050506040513d60208110156112d757600080fd5b5051905060006112e5611f5f565b90508181116112f757900390506112fe565b6000925050505b90565b600061130b612211565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561135c57fe5b600281111561136757fe5b905250905060008160200151600281111561137e57fe5b141561138e5760009150506105d8565b805160049060ff16602081106113a057fe5b601091828204019190066002029054906101000a900461ffff16915050919050565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526006602052604090205416331461145757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b3373ffffffffffffffffffffffffffffffffffffffff821614156114dc57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff0000000000000000000000000000000000000000821681179092559091169081146115815760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461160c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461170857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b61083c81612167565b611719612211565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561176a57fe5b600281111561177557fe5b905250905060006117858361046f565b905080156115815773ffffffffffffffffffffffffffffffffffffffff80841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b15801561183b57600080fd5b505af115801561184f573d6000803e3d6000fd5b505050506040513d602081101561186557600080fd5b50516118d257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff16602081106118e857fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff166020811061192357fe5b01556040805173ffffffffffffffffffffffffffffffffffffffff80871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b61198a612228565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152611a01612256565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411611a1a57905050505050509050611a61612256565b604080516104008101918290529060089060209082845b815481526020019060010190808311611a7857505050505090506060602a805480602002602001604051908101604052809291908181526020018280548015611af757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611acc575b5050505050905060005b8151811015611dc95760006001848360208110611b1a57fe5b602002015103905060006001868460208110611b3257fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115611dbe57600060066000878781518110611b7257fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015611c6857600080fd5b505af1158015611c7c573d6000803e3d6000fd5b505050506040513d6020811015611c9257600080fd5b5051611cff57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b6001888660208110611d0d57fe5b602002019061ffff16908161ffff16815250506001878660208110611d2e57fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290879087908110611d6357fe5b60200260200101518284604051808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a1505b505050600101611b01565b50611dd76004846020612275565b50610b50600883602061230b565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000611f69612256565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411611f825790505050505050905060005b6020811015611ff2576001828260208110611fdb57fe5b60200201510361ffff169290920191600101611fc4565b50611ffb612228565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608086018290527001000000000000000000000000000000009094049092166080850152602a805486518184028101840190975280875297909202633b9aca00029693949293908301828280156120d857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116120ad575b505050505090506120e7612256565b604080516104008101918290529060089060209082845b8154815260200190600101908083116120fe575050505050905060005b825181101561214557600182826020811061213257fe5b602002015103959095019460010161211b565b505050505090565b60008183101561215e575081612161565b50805b92915050565b60035473ffffffffffffffffffffffffffffffffffffffff908116908216811461220d57600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518061040001604052806020906020820280368337509192915050565b6002830191839082156122fb5791602002820160005b838211156122cb57835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030261228b565b80156122f95782816101000a81549061ffff02191690556002016020816001010492830192600103026122cb565b505b50612307929150612345565b5090565b8260208101928215612339579160200282015b8281111561233957825182559160200191906001019061231e565b5061230792915061237c565b5b808211156123075780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000168155600101612346565b5b80821115612307576000815560010161237d56fea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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





func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingCaller.contract.Call(opts, result, method, params...)
}



func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingTransactor.contract.Transfer(opts)
}


func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingTransactor.contract.Transact(opts, method, params...)
}





func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OffchainAggregatorBilling.Contract.contract.Call(opts, result, method, params...)
}



func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.contract.Transfer(opts)
}


func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.contract.Transact(opts, method, params...)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OffchainAggregatorBilling.contract.Call(opts, out, "LINK")
	return *ret0, err
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) LINK() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.LINK(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) LINK() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.LINK(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) AccountingGasCost(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregatorBilling.contract.Call(opts, out, "accountingGasCost")
	return *ret0, err
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) AccountingGasCost() (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.AccountingGasCost(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) AccountingGasCost() (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.AccountingGasCost(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) GetBilling(opts *bind.CallOpts) (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	ret := new(struct {
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	out := ret
	err := _OffchainAggregatorBilling.contract.Call(opts, out, "getBilling")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregatorBilling.contract.Call(opts, out, "linkAvailableForPayment")
	return *ret0, err
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.LinkAvailableForPayment(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.LinkAvailableForPayment(&_OffchainAggregatorBilling.CallOpts)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _OffchainAggregatorBilling.contract.Call(opts, out, "oracleObservationCount", _signerOrTransmitter)
	return *ret0, err
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregatorBilling.Contract.OracleObservationCount(&_OffchainAggregatorBilling.CallOpts, _signerOrTransmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregatorBilling.Contract.OracleObservationCount(&_OffchainAggregatorBilling.CallOpts, _signerOrTransmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OffchainAggregatorBilling.contract.Call(opts, out, "owedPayment", _transmitter)
	return *ret0, err
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.OwedPayment(&_OffchainAggregatorBilling.CallOpts, _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.OwedPayment(&_OffchainAggregatorBilling.CallOpts, _transmitter)
}




func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OffchainAggregatorBilling.contract.Call(opts, out, "owner")
	return *ret0, err
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
	return event, nil
}


const OwnedABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var OwnedBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b03191633179055610304806100326000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b14610081575b600080fd5b61004e6100b4565b005b6100586101b6565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e6004803603602081101561009757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166101d2565b60015473ffffffffffffffffffffffffffffffffffffffff16331461013a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461025857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a35056fea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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





func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}



func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}


func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}





func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}



func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}


func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}




func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Owned.contract.Call(opts, out, "owner")
	return *ret0, err
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
	return event, nil
}


const SimpleReadAccessControllerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var SimpleReadAccessControllerBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556001805460ff60a01b1916600160a01b1790556109ed806100456000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638823da6c11610076578063a118f2491161005b578063a118f249146101fd578063dc7f012414610230578063f2fde38b14610238576100a3565b80638823da6c146101995780638da5cb5b146101cc576100a3565b80630a756983146100a85780636b14daf8146100b257806379ba5097146101895780638038e4a114610191575b600080fd5b6100b061026b565b005b610175600480360360408110156100c857600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561010057600080fd5b82018360208201111561011257600080fd5b8035906020019184600183028401116401000000008311171561013457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610368945050505050565b604080519115158252519081900360200190f35b6100b061039b565b6100b061049d565b6100b0600480360360208110156101af57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166105af565b6101d46106e7565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100b06004803603602081101561021357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610703565b610175610792565b6100b06004803603602081101561024e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166107b3565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff161561036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b600061037483836108af565b80610394575073ffffffffffffffffffffffffffffffffffffffff831632145b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461042157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461052357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff1661036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60005473ffffffffffffffffffffffffffffffffffffffff16331461063557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff16156106e45773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461078957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6106e481610904565b60015474010000000000000000000000000000000000000000900460ff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461083957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff168061039457505060015474010000000000000000000000000000000000000000900460ff161592915050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff166106e45773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a15056fea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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





func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerCaller.contract.Call(opts, result, method, params...)
}



func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerTransactor.contract.Transfer(opts)
}


func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerTransactor.contract.Transact(opts, method, params...)
}





func (_SimpleReadAccessController *SimpleReadAccessControllerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleReadAccessController.Contract.contract.Call(opts, result, method, params...)
}



func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.contract.Transfer(opts)
}


func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.contract.Transact(opts, method, params...)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SimpleReadAccessController.contract.Call(opts, out, "checkEnabled")
	return *ret0, err
}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) CheckEnabled() (bool, error) {
	return _SimpleReadAccessController.Contract.CheckEnabled(&_SimpleReadAccessController.CallOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) CheckEnabled() (bool, error) {
	return _SimpleReadAccessController.Contract.CheckEnabled(&_SimpleReadAccessController.CallOpts)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SimpleReadAccessController.contract.Call(opts, out, "hasAccess", _user, _calldata)
	return *ret0, err
}




func (_SimpleReadAccessController *SimpleReadAccessControllerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _SimpleReadAccessController.Contract.HasAccess(&_SimpleReadAccessController.CallOpts, _user, _calldata)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _SimpleReadAccessController.Contract.HasAccess(&_SimpleReadAccessController.CallOpts, _user, _calldata)
}




func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleReadAccessController.contract.Call(opts, out, "owner")
	return *ret0, err
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
	return event, nil
}


const SimpleWriteAccessControllerABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var SimpleWriteAccessControllerBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556001805460ff60a01b1916600160a01b1790556109bb806100456000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638823da6c11610076578063a118f2491161005b578063a118f249146101fd578063dc7f012414610230578063f2fde38b14610238576100a3565b80638823da6c146101995780638da5cb5b146101cc576100a3565b80630a756983146100a85780636b14daf8146100b257806379ba5097146101895780638038e4a114610191575b600080fd5b6100b061026b565b005b610175600480360360408110156100c857600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561010057600080fd5b82018360208201111561011257600080fd5b8035906020019184600183028401116401000000008311171561013457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610368945050505050565b604080519115158252519081900360200190f35b6100b06103be565b6100b06104c0565b6100b0600480360360208110156101af57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166105d2565b6101d461070a565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100b06004803603602081101561021357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610726565b6101756107b5565b6100b06004803603602081101561024e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166107d6565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff161561036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff16806103b7575060015474010000000000000000000000000000000000000000900460ff16155b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461044457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461054657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff1661036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60005473ffffffffffffffffffffffffffffffffffffffff16331461065857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff16156107075773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146107ac57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b610707816108d2565b60015474010000000000000000000000000000000000000000900460ff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461085c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff166107075773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a15056fea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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





func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerCaller.contract.Call(opts, result, method, params...)
}



func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerTransactor.contract.Transfer(opts)
}


func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerTransactor.contract.Transact(opts, method, params...)
}





func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleWriteAccessController.Contract.contract.Call(opts, result, method, params...)
}



func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.contract.Transfer(opts)
}


func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.contract.Transact(opts, method, params...)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SimpleWriteAccessController.contract.Call(opts, out, "checkEnabled")
	return *ret0, err
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) CheckEnabled() (bool, error) {
	return _SimpleWriteAccessController.Contract.CheckEnabled(&_SimpleWriteAccessController.CallOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerSession) CheckEnabled() (bool, error) {
	return _SimpleWriteAccessController.Contract.CheckEnabled(&_SimpleWriteAccessController.CallOpts)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) HasAccess(opts *bind.CallOpts, _user common.Address, arg1 []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SimpleWriteAccessController.contract.Call(opts, out, "hasAccess", _user, arg1)
	return *ret0, err
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) HasAccess(_user common.Address, arg1 []byte) (bool, error) {
	return _SimpleWriteAccessController.Contract.HasAccess(&_SimpleWriteAccessController.CallOpts, _user, arg1)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerSession) HasAccess(_user common.Address, arg1 []byte) (bool, error) {
	return _SimpleWriteAccessController.Contract.HasAccess(&_SimpleWriteAccessController.CallOpts, _user, arg1)
}




func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleWriteAccessController.contract.Call(opts, out, "owner")
	return *ret0, err
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
	return event, nil
}


const TestOffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"billingData\",\"outputs\":[{\"internalType\":\"uint16[32]\",\"name\":\"osbservationsCounts\",\"type\":\"uint16[32]\"},{\"internalType\":\"uint256[32]\",\"name\":\"gasReimbursements\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDigest\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"availableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"testDecodeReport\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"int192[]\",\"name\":\"\",\"type\":\"int192[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes16\",\"name\":\"latestConfigDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint40\",\"name\":\"latestEpochAndRound\",\"type\":\"uint40\"},{\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"latestAggregatorRoundId\",\"type\":\"uint32\"}],\"internalType\":\"structOffchainAggregator.HotVars\",\"name\":\"h\",\"type\":\"tuple\"}],\"name\":\"testExposeHotvarsDummy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"txGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reasonableGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maximumGasPrice\",\"type\":\"uint256\"}],\"name\":\"testImpliedGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"testPayee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_x\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"_y\",\"type\":\"uint16\"}],\"name\":\"testSaturatingAddUint16\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitterOrSigner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amountLinkWei\",\"type\":\"uint256\"}],\"name\":\"testSetGasReimbursements\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_amount\",\"type\":\"uint16\"}],\"name\":\"testSetOracleObservationCount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testTotalLinkDue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"linkDue\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callDataCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"}],\"name\":\"testTransmitterGasCostEthWei\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var TestOffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b506040516200652038038062006520833981016040819052620000359162000606565b604080518082019091526004815263151154d560e21b6020820152600080546001600160a01b031916331781558b918b918b918b918b918b918b918b918b918b91908b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b88620000998787878787620001f9565b620000a481620002f1565b6001600160601b0319606083901b16608052620000c062000420565b620000ca62000420565b60005b60208160ff1610156200011f576001838260ff1660208110620000ec57fe5b602002019061ffff16908161ffff16815250506001828260ff16602081106200011157fe5b6020020152600101620000cd565b506200012f60048360206200043f565b506200013f6008826020620004dc565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e05250508351620001889350602e92506020850191506200051b565b50620001948662000365565b8460170b60a08160170b60401b815250508360170b60c08160170b60401b815250505050505050505050505050506001602f60006101000a81548160ff021916908315150217905550505050505050505050505050505050505050505050506200075e565b6040805160a08101825263ffffffff878116808352878216602084018190528783168486018190528784166060860181905293871660809095018590526002805463ffffffff191690931763ffffffff60201b19166401000000009092029190911763ffffffff60401b1916680100000000000000009091021763ffffffff60601b19166c010000000000000000000000009092029190911763ffffffff60801b1916600160801b909202919091179055517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b690620002e290879087908790879087906200072f565b60405180910390a15050505050565b6003546001600160a01b0390811690821681146200036157600380546001600160a01b0319166001600160a01b0384161790556040517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891290620003589083908590620006de565b60405180910390a15b5050565b6000546001600160a01b031633146200039b5760405162461bcd60e51b81526004016200039290620006f8565b60405180910390fd5b602d546001600160a01b0368010000000000000000909104811690821681146200036157602d8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215620004ca5791602002820160005b838211156200049857835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030262000456565b8015620004c85782816101000a81549061ffff021916905560020160208160010104928301926001030262000498565b505b50620004d89291506200058d565b5090565b82602081019282156200050d579160200282015b828111156200050d578251825591602001919060010190620004f0565b50620004d8929150620005a9565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200055e57805160ff19168380011785556200050d565b828001600101855582156200050d57918201828111156200050d578251825591602001919060010190620004f0565b5b80821115620004d857805461ffff191681556001016200058e565b5b80821115620004d85760008155600101620005aa565b80516001600160a01b0381168114620005d857600080fd5b92915050565b8051601781900b8114620005d857600080fd5b805163ffffffff81168114620005d857600080fd5b6000806000806000806000806000806101408b8d03121562000626578586fd5b620006328c8c620005f1565b9950620006438c60208d01620005f1565b9850620006548c60408d01620005f1565b9750620006658c60608d01620005f1565b9650620006768c60808d01620005f1565b9550620006878c60a08d01620005c0565b9450620006988c60c08d01620005c0565b9350620006a98c60e08d01620005de565b9250620006bb8c6101008d01620005de565b9150620006cd8c6101208d01620005c0565b90509295989b9194979a5092959850565b6001600160a01b0392831681529116602082015260400190565b60208082526016908201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604082015260600190565b63ffffffff95861681529385166020850152918416604084015283166060830152909116608082015260a00190565b60805160601c60a05160401c60c05160401c60e05160f81c615d5d620007c360003980610af2525080611562528061299e525080610a395280612971525080610a155280612082528061213a5280612cdd52806135f65280613dc95250615d5d6000f3fe608060405234801561001057600080fd5b506004361061034c5760003560e01c80638205bf6a116101bd578063c1075329116100f9578063e4902f82116100a2578063f2fde38b1161007c578063f2fde38b146106ca578063fa98a1c7146106dd578063fbffd2c1146106f0578063feaf968c146107035761034c565b8063e4902f821461068b578063e5fe45771461069e578063eb5dcd6c146106b75761034c565b8063d18bf87e116100d3578063d18bf87e1461065d578063dc7f012414610670578063e2851911146106785761034c565b8063c10753291461062f578063c980753914610642578063d09dc339146106555761034c565b80639c849b3011610166578063b121e14711610140578063b121e147146105e3578063b5ab58dc146105f6578063b633620c14610609578063bd8247061461061c5761034c565b80639c849b30146105aa5780639eb6e060146105bd578063a118f249146105d05761034c565b80638da5cb5b116101975780638da5cb5b1461055e5780639a6fc8f5146105665780639b764d971461058a5761034c565b80638205bf6a146105305780638823da6c146105385780638ac28d5a1461054b5761034c565b806354fd4d501161028c5780637284e4161161023557806379ba50971161020f57806379ba5097146104f45780638038e4a1146104fc578063814118341461050457806381ff7048146105195761034c565b80637284e416146104b757806373f666f8146104cc57806377096177146104d45761034c565b806366cfeaf11161026657806366cfeaf11461047a5780636b14daf81461048f57806370da2f67146104af5761034c565b806354fd4d5014610457578063585aa7de1461045f578063668a0f02146104725761034c565b806329937268116102f95780633a5381b5116102d35780633a5381b51461040a5780633b5cdfa2146104125780633c04967b1461043457806350d25bcd1461044f5761034c565b806329937268146103c9578063313ce567146103e257806335a2e492146103f75761034c565b80631327d3d81161032a5780631327d3d81461038c5780631b6b6d231461039f57806322adbc78146103b45761034c565b80630a756983146103515780630eafb25b1461035b578063102a474b14610384575b600080fd5b61035961070b565b005b61036e610369366004614962565b6107aa565b60405161037b9190615413565b60405180910390f35b61036e610918565b61035961039a366004614962565b610928565b6103a7610a13565b60405161037b919061506c565b6103bc610a37565b60405161037b91906153a8565b6103d1610a5b565b60405161037b959493929190615c30565b6103ea610af0565b60405161037b9190615cc3565b610359610405366004614dce565b610b14565b6103a7610b17565b610425610420366004614d9b565b610b3f565b60405161037b9392919061536b565b61043c610b5a565b60405161037b9796959493929190615221565b61036e610c9b565b61036e610d02565b61035961046d366004614abc565b610d07565b61036e6114bb565b610482611522565b60405161037b91906152e1565b6104a261049d3660046149b1565b61152b565b60405161037b91906152d6565b6103bc611560565b6104bf611584565b60405161037b919061541c565b61036e6115eb565b6104e76104e2366004614e56565b6115f1565b60405161037b9190615b9b565b610359611608565b6103596116bb565b61050c611752565b60405161037b91906151c7565b6105216117c1565b60405161037b93929190615bee565b61036e6117e2565b610359610546366004614962565b611849565b610359610559366004614962565b611934565b6103a7611983565b610579610574366004614ef0565b61199f565b60405161037b959493929190615c90565b61059d610598366004614ddf565b611a1f565b60405161037b9190615bb8565b6103596105b8366004614a53565b611a2b565b6103a76105cb366004614962565b611c38565b6103596105de366004614962565b611c63565b6103596105f1366004614962565b611ca3565b61036e610604366004614dfb565b611d81565b61036e610617366004614dfb565b611de9565b61035961062a366004614e87565b611e51565b61035961063d366004614a29565b611f89565b610359610650366004614cfc565b612212565b61036e612cd8565b61035961066b3660046149fe565b612dad565b6104a2612e10565b610359610686366004614a29565b612e19565b61059d610699366004614962565b612eb5565b6106a6612f76565b60405161037b95949392919061530e565b6103596106c536600461497d565b613030565b6103596106d8366004614962565b613156565b61036e6106eb366004614e2b565b613203565b6103596106fe366004614962565b613218565b610579613258565b60005473ffffffffffffffffffffffffffffffffffffffff16331461074b5760405162461bcd60e51b815260040161074290615679565b60405180910390fd5b602f5460ff16156107a857602f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b60006107b461462d565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561080557fe5b600281111561081057fe5b905250905060008160200151600281111561082757fe5b1415610837576000915050610913565b61083f614644565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216948401949094526c0100000000000000000000000082048116606084018190527001000000000000000000000000000000009092041660808301528351919260009260019160049160ff169081106108cc57fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff166020811061090a57fe5b01540301925050505b919050565b60006109226132d6565b90505b90565b60005473ffffffffffffffffffffffffffffffffffffffff16331461095f5760405162461bcd60e51b815260040161074290615679565b602d5473ffffffffffffffffffffffffffffffffffffffff6801000000000000000090910481169082168114610a0f57602d80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff166801000000000000000073ffffffffffffffffffffffffffffffffffffffff85811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000610a6b614644565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b7f000000000000000000000000000000000000000000000000000000000000000081565b50565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b6000806060610b4d846134c4565b9250925092509193909250565b610b62614672565b610b6a614672565b6000806000806000610b7a614644565b506040805160a08101825260025463ffffffff80821680845264010000000083048216602080860182905268010000000000000000850484168688018190526c010000000000000000000000008604851660608801819052700100000000000000000000000000000000909604909416608087018190528751610400810198899052969760049760089795969495939091899190826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411610c1357505060408051610400810191829052959c508b9450602093509150839050845b815481526020019060010190808311610c69575050505050955097509750975097509750975097505090919293949596565b6000610cde336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061152b92505050565b610cfa5760405162461bcd60e51b8152600401610742906159ae565b6109226134ea565b600481565b868560ff8616818314610d2c5760405162461bcd60e51b815260040161074290615577565b60005473ffffffffffffffffffffffffffffffffffffffff163314610d635760405162461bcd60e51b815260040161074290615679565b60295415610f2e57602980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81019160009183908110610da057fe5b6000918252602082200154602a805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110610dd457fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff169050610e0181613526565b73ffffffffffffffffffffffffffffffffffffffff80831660009081526028602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090811690915592841682529020805490911690556029805480610e6a57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055602a805480610ecd57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610d63915050565b60005b8a811015611344576000602860008e8e85818110610f4b57fe5b9050602002016020810190610f609190614962565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002054610100900460ff166002811115610f9a57fe5b14610fb75760405162461bcd60e51b815260040161074290615866565b6040805180820190915260ff8216815260016020820152602860008e8e85818110610fde57fe5b9050602002016020810190610ff39190614962565b73ffffffffffffffffffffffffffffffffffffffff168152602080820192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561108257fe5b02179055506000915060069050818c8c8581811061109c57fe5b90506020020160208101906110b19190614962565b73ffffffffffffffffffffffffffffffffffffffff90811682526020820192909252604001600020541614156110f95760405162461bcd60e51b815260040161074290615642565b6000602860008c8c8581811061110b57fe5b90506020020160208101906111209190614962565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002054610100900460ff16600281111561115a57fe5b146111775760405162461bcd60e51b81526004016107429061550b565b6040805180820190915260ff8216815260026020820152602860008c8c8581811061119e57fe5b90506020020160208101906111b39190614962565b73ffffffffffffffffffffffffffffffffffffffff168152602080820192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561124257fe5b021790555090505060298c8c8381811061125857fe5b905060200201602081019061126d9190614962565b81546001810183556000928352602090922090910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909216919091179055602a8a8a838181106112d257fe5b90506020020160208101906112e79190614962565b815460018082018455600093845260209093200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905501610f31565b50602b805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602d80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000090911617938490559091048116911661141030828f8f8f8f8f8f8f8f613730565b602b60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602b60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f6040516114a49a99989796959493929190615c5f565b60405180910390a150505050505050505050505050565b60006114fe336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061152b92505050565b61151a5760405162461bcd60e51b8152600401610742906159ae565b61092261377b565b602b5460801b90565b600061153783836137a1565b80611557575073ffffffffffffffffffffffffffffffffffffffff831632145b90505b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60606115c7336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061152b92505050565b6115e35760405162461bcd60e51b8152600401610742906159ae565b6109226137de565b6159d881565b60006115ff85858585613889565b95945050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461163f5760405162461bcd60e51b815260040161074290615466565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1633146116f25760405162461bcd60e51b815260040161074290615679565b602f5460ff166107a857602f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b6060602a8054806020026020016040519081016040528092919081815260200182805480156117b757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161178c575b5050505050905090565b602d54602b5463ffffffff808316926401000000009004169060801b909192565b6000611825336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061152b92505050565b6118415760405162461bcd60e51b8152600401610742906159ae565b6109226138e0565b60005473ffffffffffffffffffffffffffffffffffffffff1633146118805760405162461bcd60e51b815260040161074290615679565b73ffffffffffffffffffffffffffffffffffffffff811660009081526030602052604090205460ff1615610b145773ffffffffffffffffffffffffffffffffffffffff81166000908152603060205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19061192990839061506c565b60405180910390a150565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526006602052604090205416331461197a5760405162461bcd60e51b815260040161074290615a53565b610b1481613526565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008060006119e8336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061152b92505050565b611a045760405162461bcd60e51b8152600401610742906159ae565b611a0d8661393b565b939a9299509097509550909350915050565b60006115578383613a24565b60005473ffffffffffffffffffffffffffffffffffffffff163314611a625760405162461bcd60e51b815260040161074290615679565b828114611a815760405162461bcd60e51b815260040161074290615af8565b60005b83811015611c31576000858583818110611a9a57fe5b9050602002016020810190611aaf9190614962565b90506000848484818110611abf57fe5b9050602002016020810190611ad49190614962565b73ffffffffffffffffffffffffffffffffffffffff8084166000908152600660205260409020549192501680158080611b3857508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b611b545760405162461bcd60e51b8152600401610742906159e5565b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614611c21578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505060019092019150611a849050565b5050505050565b73ffffffffffffffffffffffffffffffffffffffff9081166000908152600660205260409020541690565b60005473ffffffffffffffffffffffffffffffffffffffff163314611c9a5760405162461bcd60e51b815260040161074290615679565b610b1481613a3c565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260076020526040902054163314611ce95760405162461bcd60e51b81526004016107429061542f565b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b6000611dc4336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061152b92505050565b611de05760405162461bcd60e51b8152600401610742906159ae565b61155a82613ae7565b6000611e2c336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061152b92505050565b611e485760405162461bcd60e51b8152600401610742906159ae565b61155a82613b1d565b60035473ffffffffffffffffffffffffffffffffffffffff1680611e875760405162461bcd60e51b8152600401610742906158d2565b6040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636b14daf890611ede903390600090369060040161508d565b60206040518083038186803b158015611ef657600080fd5b505afa158015611f0a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f2e9190614b84565b80611f50575060005473ffffffffffffffffffffffffffffffffffffffff1633145b611f6c5760405162461bcd60e51b81526004016107429061571e565b611f74613b72565b611f818686868686613f3e565b505050505050565b6003546040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690636b14daf890611fe4903390600090369060040161508d565b60206040518083038186803b158015611ffc57600080fd5b505afa158015612010573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120349190614b84565b80612056575060005473ffffffffffffffffffffffffffffffffffffffff1633145b6120725760405162461bcd60e51b81526004016107429061571e565b600061207c6132d6565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016120d9919061506c565b60206040518083038186803b1580156120f157600080fd5b505afa158015612105573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121299190614e13565b90508181101561213857600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85612181858503876140bc565b6040518363ffffffff1660e01b815260040161219e929190615115565b602060405180830381600087803b1580156121b857600080fd5b505af11580156121cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121f09190614b84565b61220c5760405162461bcd60e51b815260040161074290615b2d565b50505050565b60005a90506122258888888888886140d3565b36146122435760405162461bcd60e51b81526004016107429061582f565b61224b614691565b6040805160808082018352602b549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff1660208301527501000000000000000000000000000000000000000000810460ff169282019290925276010000000000000000000000000000000000000000000090910463ffffffff166060820152815260006122fa8a8a018b614bab565b60408501526080840182905283515190925060589190911b907fffffffffffffffffffffffffffffffff000000000000000000000000000000008083169116146123565760405162461bcd60e51b815260040161074290615a8a565b608083015183516020015164ffffffffff8083169116106123895760405162461bcd60e51b8152600401610742906155d4565b83516040015160ff1689116123b05760405162461bcd60e51b8152600401610742906156b0565b60208911156123d15760405162461bcd60e51b815260040161074290615753565b8689146123f05760405162461bcd60e51b815260040161074290615b64565b602084604001515111156124165760405162461bcd60e51b815260040161074290615909565b83600001516040015160020260ff16846040015151116124485760405162461bcd60e51b8152600401610742906157f8565b8867ffffffffffffffff8111801561245f57600080fd5b506040519080825280601f01601f19166020018201604052801561248a576020820181803683370190505b50606085015260005b60ff81168a11156124fb57868160ff16602081106124ad57fe5b1a60f81b85606001518260ff16815181106124c457fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101612493565b5083604001515167ffffffffffffffff8111801561251857600080fd5b506040519080825280601f01601f191660200182016040528015612543576020820181803683370190505b506020850152612551614672565b60005b8560400151518160ff161015612622576000858260ff166020811061257557fe5b1a905082816020811061258457fe5b6020020151156125a65760405162461bcd60e51b81526004016107429061578a565b6001838260ff16602081106125b757fe5b91151560209283029190910152869060ff84169081106125d357fe5b1a60f81b87602001518360ff16815181106125ea57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101612554565b5061262b61462d565b336000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561266657fe5b600281111561267157fe5b905250905060028160200151600281111561268857fe5b1480156126c95750602a816000015160ff16815481106126a457fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b6126e55760405162461bcd60e51b815260040161074290615a1c565b5050835164ffffffffff909116602090910152505060405160009061270d908b908b9061505c565b6040518091039020905061271f614672565b61272761462d565b60005b898110156128cc5760006001858760600151848151811061274757fe5b60209101015160f81c601b018e8e8681811061275f57fe5b905060200201358d8d8781811061277257fe5b9050602002013560405160008152602001604052604051612796949392919061538a565b6020604051602081039080840390855afa1580156127b8573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526028602090815290849020838501909452835460ff8082168552929650929450840191610100900416600281111561283257fe5b600281111561283d57fe5b905250925060018360200151600281111561285457fe5b146128715760405162461bcd60e51b8152600401610742906154d4565b8251849060ff166020811061288257fe5b6020020151156128a45760405162461bcd60e51b8152600401610742906156e7565b600184846000015160ff16602081106128b957fe5b911515602090920201525060010161272a565b5050505060005b600182604001515103811015612948576000826040015182600101815181106128f857fe5b602002602001015160170b8360400151838151811061291357fe5b602002602001015160170b131590508061293f5760405162461bcd60e51b815260040161074290615977565b506001016128d3565b5060408101518051600091906002810490811061296157fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b131580156129c757507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b6129e35760405162461bcd60e51b815260040161074290615940565b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b825267ffffffffffffffff4281166020808501918252895188015187166000908152602c82528690209451855492519093167801000000000000000000000000000000000000000000000000029290930b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff00000000000000000000000000000000000000000000000090921691909117161790915585519093015181860151938601516080870151925191909316937ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d45193612aef93879333939291906153b6565b60405180910390a281516060015160405160009163ffffffff16907f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac6027190612b37904290615413565b60405180910390a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f42604051612b819190615413565b60405180910390a3612b9e8260000151606001518260170b6140eb565b5080518051602b8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff92831602179091558210612cbf57fe5b612ccd82826020015161422d565b505050505050505050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401612d34919061506c565b60206040518083038186803b158015612d4c57600080fd5b505afa158015612d60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d849190614e13565b90506000612d906132d6565b9050818111612da25790039050610925565b600092505050610925565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260286020908152604090912054829160049160ff16908110612de757fe5b601091828204019190066002026101000a81548161ffff021916908361ffff1602179055505050565b602f5460ff1681565b600073ffffffffffffffffffffffffffffffffffffffff8316600090815260286020526040902054610100900460ff166002811115612e5457fe5b1415612e725760405162461bcd60e51b81526004016107429061560b565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260286020908152604090912054600183019160089160ff16908110612eaf57fe5b01555050565b6000612ebf61462d565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602860209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115612f1057fe5b6002811115612f1b57fe5b9052509050600081602001516002811115612f3257fe5b1415612f42576000915050610913565b805160049060ff1660208110612f5457fe5b601091828204019190066002029054906101000a900461ffff16915050919050565b600080808080333214612f9b5760405162461bcd60e51b8152600401610742906157c1565b5050602b5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602c6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600660205260409020541633146130765760405162461bcd60e51b81526004016107429061549d565b3373ffffffffffffffffffffffffffffffffffffffff821614156130ac5760405162461bcd60e51b815260040161074290615ac1565b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff0000000000000000000000000000000000000000821681179092559091169081146131515760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461318d5760405162461bcd60e51b815260040161074290615679565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613210848484614457565b949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461324f5760405162461bcd60e51b815260040161074290615679565b610b1481614474565b60008060008060006132a1336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061152b92505050565b6132bd5760405162461bcd60e51b8152600401610742906159ae565b6132c5614516565b945094509450945094509091929394565b60006132e0614672565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116132f95790505050505050905060005b602081101561336957600182826020811061335257fe5b60200201510361ffff16929092019160010161333b565b50613372614644565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608086018290527001000000000000000000000000000000009094049092166080850152602a805486518184028101840190975280875297909202633b9aca000296939492939083018282801561344f57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613424575b5050505050905061345e614672565b604080516104008101918290529060089060209082845b815481526020019060010190808311613475575050505050905060005b82518110156134bc5760018282602081106134a957fe5b6020020151039590950194600101613492565b505050505090565b6000806060838060200190518101906134dd9190614c5e565b9196909550909350915050565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c6020526040902054601790810b900b90565b61352e61462d565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561357f57fe5b600281111561358a57fe5b9052509050600061359a836107aa565b905080156131515773ffffffffffffffffffffffffffffffffffffffff808416600090815260066020526040908190205490517fa9059cbb000000000000000000000000000000000000000000000000000000008152908216917f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb9061362d9084908690600401615115565b602060405180830381600087803b15801561364757600080fd5b505af115801561365b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061367f9190614b84565b61369b5760405162461bcd60e51b815260040161074290615b2d565b60016004846000015160ff16602081106136b157fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16602081106136ec57fe5b01556040517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290613722908690849086906150e4565b60405180910390a150505050565b60008a8a8a8a8a8a8a8a8a8a6040516020016137559a9998979695949392919061513b565b6040516020818303038152906040528051906020012090509a9950505050505050505050565b602b54760100000000000000000000000000000000000000000000900463ffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff821660009081526030602052604081205460ff1680611557575050602f5460ff161592915050565b602e8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156117b75780601f1061385d576101008083540402835291602001916117b7565b820191906000526020600020905b81548152906001019060200180831161386b57509395945050505050565b6000818510156138ab5760405162461bcd60e51b815260040161074290615542565b81850383016159d801633b9aca00858202026fffffffffffffffffffffffffffffffff81106138d657fe5b9695505050505050565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e740000000000000000000000000000000000815250906139ab5760405162461bcd60e51b8152600401610742919061541c565b506139b461462d565b5050505063ffffffff83166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b60006115578261ffff168461ffff160161ffff6140bc565b73ffffffffffffffffffffffffffffffffffffffff811660009081526030602052604090205460ff16610b145773ffffffffffffffffffffffffffffffffffffffff81166000908152603060205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49061192990839061506c565b600063ffffffff821115613afd57506000610913565b5063ffffffff166000908152602c6020526040902054601790810b900b90565b600063ffffffff821115613b3357506000610913565b5063ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b613b7a614644565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152613bf1614672565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411613c0a57905050505050509050613c51614672565b604080516104008101918290529060089060209082845b815481526020019060010190808311613c6857505050505090506060602a805480602002602001604051908101604052809291908181526020018280548015613ce757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613cbc575b5050505050905060005b8151811015613f225760006001848360208110613d0a57fe5b602002015103905060006001868460208110613d2257fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115613f1757600060066000878781518110613d6257fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82846040518363ffffffff1660e01b8152600401613e22929190615115565b602060405180830381600087803b158015613e3c57600080fd5b505af1158015613e50573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e749190614b84565b613e905760405162461bcd60e51b815260040161074290615b2d565b6001888660208110613e9e57fe5b602002019061ffff16908161ffff16815250506001878660208110613ebf57fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290879087908110613ef457fe5b60200260200101518284604051613f0d939291906150e4565b60405180910390a1505b505050600101613cf1565b50613f3060048460206146c5565b50611c31600883602061475b565b6040805160a08101825263ffffffff87811680835287821660208401819052878316848601819052878416606086018190529387166080909501859052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000169093177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff16640100000000909202919091177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff1668010000000000000000909102177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c01000000000000000000000000909202919091177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000909202919091179055517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6906140ad9087908790879087908790615c30565b60405180910390a15050505050565b6000818310156140cd57508161155a565b50919050565b602083810286019082020160e4019695505050505050565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff168061411b5750610a0f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff81166000908152602c6020526040908190205490517fbeed9b51000000000000000000000000000000000000000000000000000000008152601791820b90910b9073ffffffffffffffffffffffffffffffffffffffff84169063beed9b5190620186a0906141bb90869086908b908b90600401615bc7565b602060405180830381600088803b1580156141d557600080fd5b5087f193505050508015614224575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820190925261422191810190614b84565b60015b611f8157611c31565b61423561462d565b336000908152602860209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561427057fe5b600281111561427b57fe5b9052509050614288614644565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606085015270010000000000000000000000000000000090920416608083015282516104008101938490529192614357928692909160049190826000855b82829054906101000a900461ffff1661ffff168152602001906002019060208260010104928301926001038202915080841161431557905050505050506145b5565b6143659060049060206146c5565b5060028260200151600281111561437857fe5b146143955760405162461bcd60e51b81526004016107429061589d565b60006143bc633b9aca003a04836020015163ffffffff16846000015163ffffffff16614457565b90506010360260005a905060006143db8863ffffffff16858585613889565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff1683028161440957fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff166020811061443257fe5b015401016008886000015160ff166020811061444a57fe5b0155505050505050505050565b6000838381101561446a57600285850304015b6115ff81846140bc565b60035473ffffffffffffffffffffffffffffffffffffffff9081169082168114610a0f57600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84161790556040517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129061450a90839085906150bd565b60405180910390a15050565b602b54760100000000000000000000000000000000000000000000900463ffffffff16600080808061454661462d565b5050505063ffffffff82166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b6145bd614672565b60005b83518110156146255760008482815181106145d757fe5b0160209081015160f81c91506145ff908590839081106145f357fe5b60200201516001613a24565b848260ff166020811061460e57fe5b61ffff9092166020929092020152506001016145c0565b509092915050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518061040001604052806020906020820280368337509192915050565b6040518060a001604052806146a4614795565b81526060602082018190526040820181905280820152600060809091015290565b60028301918390821561474b5791602002820160005b8382111561471b57835183826101000a81548161ffff021916908361ffff16021790555092602001926002016020816001010492830192600103026146db565b80156147495782816101000a81549061ffff021916905560020160208160010104928301926001030261471b565b505b506147579291506147bc565b5090565b8260208101928215614789579160200282015b8281111561478957825182559160200191906001019061476e565b506147579291506147f3565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b808211156147575780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001681556001016147bd565b5b8082111561475757600081556001016147f4565b803573ffffffffffffffffffffffffffffffffffffffff8116811461155a57600080fd5b60008083601f84011261483d578182fd5b50813567ffffffffffffffff811115614854578182fd5b602083019150836020808302850101111561486e57600080fd5b9250929050565b60008083601f840112614886578182fd5b50813567ffffffffffffffff81111561489d578182fd5b60208301915083602082850101111561486e57600080fd5b600082601f8301126148c5578081fd5b813567ffffffffffffffff8111156148db578182fd5b61490c60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601615cd1565b915080825283602082850101111561492357600080fd5b8060208401602084013760009082016020015292915050565b803561ffff8116811461155a57600080fd5b803563ffffffff8116811461155a57600080fd5b600060208284031215614973578081fd5b6115578383614808565b6000806040838503121561498f578081fd5b6149998484614808565b91506149a88460208501614808565b90509250929050565b600080604083850312156149c3578182fd5b6149cd8484614808565b9150602083013567ffffffffffffffff8111156149e8578182fd5b6149f4858286016148b5565b9150509250929050565b60008060408385031215614a10578182fd5b614a1a8484614808565b91506149a8846020850161493c565b60008060408385031215614a3b578182fd5b614a458484614808565b946020939093013593505050565b60008060008060408587031215614a68578182fd5b843567ffffffffffffffff80821115614a7f578384fd5b614a8b8883890161482c565b90965094506020870135915080821115614aa3578384fd5b50614ab08782880161482c565b95989497509550505050565b60008060008060008060008060a0898b031215614ad7578586fd5b883567ffffffffffffffff80821115614aee578788fd5b614afa8c838d0161482c565b909a50985060208b0135915080821115614b12578788fd5b614b1e8c838d0161482c565b909850965060408b0135915060ff82168214614b38578586fd5b90945060608a0135908082168214614b4e578485fd5b90935060808a01359080821115614b63578384fd5b50614b708b828c01614875565b999c989b5096995094979396929594505050565b600060208284031215614b95578081fd5b81518015158114614ba4578182fd5b9392505050565b600080600060608486031215614bbf578081fd5b833592506020808501359250604085013567ffffffffffffffff811115614be4578283fd5b8501601f81018713614bf4578283fd5b8035614c07614c0282615cf8565b615cd1565b81815283810190838501858402850186018b1015614c23578687fd5b8694505b83851015614c4e578035614c3a81615d18565b835260019490940193918501918501614c27565b5080955050505050509250925092565b600080600060608486031215614c72578081fd5b835192506020808501519250604085015167ffffffffffffffff811115614c97578283fd5b8501601f81018713614ca7578283fd5b8051614cb5614c0282615cf8565b81815283810190838501858402850186018b1015614cd1578687fd5b8694505b83851015614c4e578051614ce881615d18565b835260019490940193918501918501614cd5565b60008060008060008060006080888a031215614d16578081fd5b873567ffffffffffffffff80821115614d2d578283fd5b614d398b838c01614875565b909950975060208a0135915080821115614d51578283fd5b614d5d8b838c0161482c565b909750955060408a0135915080821115614d75578283fd5b50614d828a828b0161482c565b989b979a50959894979596606090950135949350505050565b600060208284031215614dac578081fd5b813567ffffffffffffffff811115614dc2578182fd5b613210848285016148b5565b6000608082840312156140cd578081fd5b60008060408385031215614df1578182fd5b614a1a848461493c565b600060208284031215614e0c578081fd5b5035919050565b600060208284031215614e24578081fd5b5051919050565b600080600060608486031215614e3f578081fd5b505081359360208301359350604090920135919050565b60008060008060808587031215614e6b578182fd5b5050823594602084013594506040840135936060013592509050565b600080600080600060a08688031215614e9e578283fd5b614ea8878761494e565b9450614eb7876020880161494e565b9350614ec6876040880161494e565b9250614ed5876060880161494e565b9150614ee4876080880161494e565b90509295509295909350565b600060208284031215614f01578081fd5b813569ffffffffffffffffffff81168114614ba4578182fd5b60008284526020808501945082825b85811015614f645782820173ffffffffffffffffffffffffffffffffffffffff614f538285614808565b168852968301969150600101614f29565b509495945050505050565b6000815180845260208085019450808401835b83811015614f6457815160170b87529582019590820190600101614f82565b600082845282826020860137806020848601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f85011685010190509392505050565b60008151808452815b8181101561500e57602081850181015186830182015201614ff2565b8181111561501f5782602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b63ffffffff169052565b6000828483379101908152919050565b73ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b600073ffffffffffffffffffffffffffffffffffffffff85168252604060208301526115ff604083018486614fa1565b73ffffffffffffffffffffffffffffffffffffffff92831681529116602082015260400190565b73ffffffffffffffffffffffffffffffffffffffff9384168152919092166020820152604081019190915260600190565b73ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b600073ffffffffffffffffffffffffffffffffffffffff8c16825267ffffffffffffffff808c16602084015260e0604084015261517c60e084018b8d614f1a565b838103606085015261518f818a8c614f1a565b905060ff8816608085015281871660a085015283810360c08501526151b5818688614fa1565b9e9d5050505050505050505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561521557835173ffffffffffffffffffffffffffffffffffffffff16835292840192918401916001016151e3565b50909695505050505050565b6108a08101818960005b602080821061523a5750615252565b825161ffff168452928301929091019060010161522b565b50505061040082018860005b602080821061526d5750615281565b82518452928301929091019060010161525e565b505050615292610800830188615052565b6152a0610820830187615052565b6152ae610840830186615052565b6152bc610860830185615052565b6152ca610880830184615052565b98975050505050505050565b901515815260200190565b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000091909116815260200190565b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000095909516855263ffffffff93909316602085015260ff91909116604084015260170b606083015267ffffffffffffffff16608082015260a00190565b6000848252836020830152606060408301526115ff6060830184614f6f565b93845260ff9290921660208401526040830152606082015260800190565b60179190910b815260200190565b60008660170b825273ffffffffffffffffffffffffffffffffffffffff8616602083015260a060408301526153ee60a0830186614f6f565b82810360608401526154008186614fe9565b9150508260808301529695505050505050565b90815260200190565b6000602082526115576020830184614fe9565b6020808252601f908201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604082015260600190565b60208082526016908201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604082015260600190565b6020808252601d908201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604082015260600190565b6020808252601e908201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604082015260600190565b6020808252601c908201527f7265706561746564207472616e736d6974746572206164647265737300000000604082015260600190565b6020808252818101527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604082015260600190565b60208082526024908201527f6f7261636c6520616464726573736573206f7574206f6620726567697374726160408201527f74696f6e00000000000000000000000000000000000000000000000000000000606082015260800190565b6020808252600c908201527f7374616c65207265706f72740000000000000000000000000000000000000000604082015260600190565b6020808252600f908201527f6164647265737320756e6b6e6f776e0000000000000000000000000000000000604082015260600190565b60208082526011908201527f7061796565206d75737420626520736574000000000000000000000000000000604082015260600190565b60208082526016908201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604082015260600190565b60208082526015908201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604082015260600190565b60208082526014908201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604082015260600190565b6020808252818101527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604082015260600190565b60208082526013908201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604082015260600190565b60208082526017908201527f6f6273657276657220696e646578207265706561746564000000000000000000604082015260600190565b60208082526014908201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604082015260600190565b6020808252601e908201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604082015260600190565b60208082526019908201527f7472616e736d6974206d65737361676520746f6f206c6f6e6700000000000000604082015260600190565b60208082526017908201527f7265706561746564207369676e65722061646472657373000000000000000000604082015260600190565b6020808252818101527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604082015260600190565b6020808252601d908201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604082015260600190565b6020808252601e908201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604082015260600190565b6020808252601e908201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604082015260600190565b60208082526017908201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604082015260600190565b60208082526009908201527f4e6f206163636573730000000000000000000000000000000000000000000000604082015260600190565b60208082526011908201527f706179656520616c726561647920736574000000000000000000000000000000604082015260600190565b60208082526018908201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604082015260600190565b60208082526017908201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604082015260600190565b60208082526015908201527f636f6e666967446967657374206d69736d617463680000000000000000000000604082015260600190565b60208082526017908201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604082015260600190565b6020808252818101527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604082015260600190565b60208082526012908201527f696e73756666696369656e742066756e64730000000000000000000000000000604082015260600190565b6020808252601e908201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604082015260600190565b6fffffffffffffffffffffffffffffffff91909116815260200190565b61ffff91909116815260200190565b63ffffffff9485168152602081019390935292166040820152606081019190915260800190565b63ffffffff93841681529190921660208201527fffffffffffffffffffffffffffffffff00000000000000000000000000000000909116604082015260600190565b63ffffffff95861681529385166020850152918416604084015283166060830152909116608082015260a00190565b600063ffffffff8c16825267ffffffffffffffff808c16602084015260e0604084015261517c60e084018b8d614f1a565b69ffffffffffffffffffff9586168152602081019490945260408401929092526060830152909116608082015260a00190565b60ff91909116815260200190565b60405181810167ffffffffffffffff81118282101715615cf057600080fd5b604052919050565b600067ffffffffffffffff821115615d0e578081fd5b5060209081020190565b8060170b8114610b1457600080fdfea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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





func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}



func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorTransactor.contract.Transfer(opts)
}


func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}





func (_TestOffchainAggregator *TestOffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TestOffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}



func (_TestOffchainAggregator *TestOffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.contract.Transfer(opts)
}


func (_TestOffchainAggregator *TestOffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.contract.Transact(opts, method, params...)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "LINK")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LINK() (common.Address, error) {
	return _TestOffchainAggregator.Contract.LINK(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _TestOffchainAggregator.Contract.LINK(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) AccountingGasCost(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "accountingGasCost")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) AccountingGasCost() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.AccountingGasCost(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) AccountingGasCost() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.AccountingGasCost(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) BillingData(opts *bind.CallOpts) (struct {
	OsbservationsCounts     [32]uint16
	GasReimbursements       [32]*big.Int
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	ret := new(struct {
		OsbservationsCounts     [32]uint16
		GasReimbursements       [32]*big.Int
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	out := ret
	err := _TestOffchainAggregator.contract.Call(opts, out, "billingData")
	return *ret, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) BillingData() (struct {
	OsbservationsCounts     [32]uint16
	GasReimbursements       [32]*big.Int
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _TestOffchainAggregator.Contract.BillingData(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) BillingData() (struct {
	OsbservationsCounts     [32]uint16
	GasReimbursements       [32]*big.Int
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _TestOffchainAggregator.Contract.BillingData(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "checkEnabled")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) CheckEnabled() (bool, error) {
	return _TestOffchainAggregator.Contract.CheckEnabled(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _TestOffchainAggregator.Contract.CheckEnabled(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "decimals")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Decimals() (uint8, error) {
	return _TestOffchainAggregator.Contract.Decimals(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _TestOffchainAggregator.Contract.Decimals(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "description")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Description() (string, error) {
	return _TestOffchainAggregator.Contract.Description(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Description() (string, error) {
	return _TestOffchainAggregator.Contract.Description(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "getAnswer", _roundId)
	return *ret0, err
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
	ret := new(struct {
		MaximumGasPrice         uint32
		ReasonableGasPrice      uint32
		MicroLinkPerEth         uint32
		LinkGweiPerObservation  uint32
		LinkGweiPerTransmission uint32
	})
	out := ret
	err := _TestOffchainAggregator.contract.Call(opts, out, "getBilling")
	return *ret, err
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
	var (
		ret0 = new([16]byte)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "getConfigDigest")
	return *ret0, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _TestOffchainAggregator.contract.Call(opts, out, "getRoundData", _roundId)
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "getTimestamp", _roundId)
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetTimestamp(&_TestOffchainAggregator.CallOpts, _roundId)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetTimestamp(&_TestOffchainAggregator.CallOpts, _roundId)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "hasAccess", _user, _calldata)
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _TestOffchainAggregator.Contract.HasAccess(&_TestOffchainAggregator.CallOpts, _user, _calldata)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _TestOffchainAggregator.Contract.HasAccess(&_TestOffchainAggregator.CallOpts, _user, _calldata)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
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
	ret := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [16]byte
	})
	out := ret
	err := _TestOffchainAggregator.contract.Call(opts, out, "latestConfigDetails")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "latestRound")
	return *ret0, err
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
	ret := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	out := ret
	err := _TestOffchainAggregator.contract.Call(opts, out, "latestRoundData")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
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
	ret := new(struct {
		ConfigDigest    [16]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	out := ret
	err := _TestOffchainAggregator.contract.Call(opts, out, "latestTransmissionDetails")
	return *ret, err
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
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "linkAvailableForPayment")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LinkAvailableForPayment(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LinkAvailableForPayment(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "maxAnswer")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MaxAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MaxAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "minAnswer")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MinAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MinAnswer(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "oracleObservationCount", _signerOrTransmitter)
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _TestOffchainAggregator.Contract.OracleObservationCount(&_TestOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _TestOffchainAggregator.Contract.OracleObservationCount(&_TestOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "owedPayment", _transmitter)
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.OwedPayment(&_TestOffchainAggregator.CallOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.OwedPayment(&_TestOffchainAggregator.CallOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "owner")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Owner() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Owner(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Owner(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestDecodeReport(opts *bind.CallOpts, report []byte) ([32]byte, [32]byte, []*big.Int, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new([32]byte)
		ret2 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _TestOffchainAggregator.contract.Call(opts, out, "testDecodeReport", report)
	return *ret0, *ret1, *ret2, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestDecodeReport(report []byte) ([32]byte, [32]byte, []*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestDecodeReport(&_TestOffchainAggregator.CallOpts, report)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestDecodeReport(report []byte) ([32]byte, [32]byte, []*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestDecodeReport(&_TestOffchainAggregator.CallOpts, report)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestImpliedGasPrice(opts *bind.CallOpts, txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "testImpliedGasPrice", txGasPrice, reasonableGasPrice, maximumGasPrice)
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestImpliedGasPrice(txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestImpliedGasPrice(&_TestOffchainAggregator.CallOpts, txGasPrice, reasonableGasPrice, maximumGasPrice)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestImpliedGasPrice(txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestImpliedGasPrice(&_TestOffchainAggregator.CallOpts, txGasPrice, reasonableGasPrice, maximumGasPrice)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestPayee(opts *bind.CallOpts, _transmitter common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "testPayee", _transmitter)
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestPayee(_transmitter common.Address) (common.Address, error) {
	return _TestOffchainAggregator.Contract.TestPayee(&_TestOffchainAggregator.CallOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestPayee(_transmitter common.Address) (common.Address, error) {
	return _TestOffchainAggregator.Contract.TestPayee(&_TestOffchainAggregator.CallOpts, _transmitter)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestSaturatingAddUint16(opts *bind.CallOpts, _x uint16, _y uint16) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "testSaturatingAddUint16", _x, _y)
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestSaturatingAddUint16(_x uint16, _y uint16) (uint16, error) {
	return _TestOffchainAggregator.Contract.TestSaturatingAddUint16(&_TestOffchainAggregator.CallOpts, _x, _y)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestSaturatingAddUint16(_x uint16, _y uint16) (uint16, error) {
	return _TestOffchainAggregator.Contract.TestSaturatingAddUint16(&_TestOffchainAggregator.CallOpts, _x, _y)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestTotalLinkDue(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "testTotalLinkDue")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestTotalLinkDue() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTotalLinkDue(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestTotalLinkDue() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTotalLinkDue(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestTransmitterGasCostEthWei(opts *bind.CallOpts, initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "testTransmitterGasCostEthWei", initialGas, gasPrice, callDataCost, gasLeft)
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestTransmitterGasCostEthWei(initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTransmitterGasCostEthWei(&_TestOffchainAggregator.CallOpts, initialGas, gasPrice, callDataCost, gasLeft)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestTransmitterGasCostEthWei(initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTransmitterGasCostEthWei(&_TestOffchainAggregator.CallOpts, initialGas, gasPrice, callDataCost, gasLeft)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "transmitters")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _TestOffchainAggregator.Contract.Transmitters(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _TestOffchainAggregator.Contract.Transmitters(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "validator")
	return *ret0, err
}




func (_TestOffchainAggregator *TestOffchainAggregatorSession) Validator() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Validator(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Validator(&_TestOffchainAggregator.CallOpts)
}




func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TestOffchainAggregator.contract.Call(opts, out, "version")
	return *ret0, err
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
	return event, nil
}
