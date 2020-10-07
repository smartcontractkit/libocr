


package exposedoffchainaggregator

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


const ExposedOffchainAggregatorABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encodedConfig\",\"type\":\"bytes\"}],\"name\":\"exposedConfigDigestFromConfigData\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"availableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var ExposedOffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b506040805160208101909152600080825280546001600160a01b03191633178155908190819081908190819081908190819081908190818080808080806200005e818080808062000184565b620000698162000276565b6001600160601b0319606083901b1660805262000085620003d4565b6200008f620003d4565b60005b60208160ff161015620000e4576001838260ff1660208110620000b157fe5b602002019061ffff16908161ffff16815250506001828260ff1660208110620000d657fe5b602002015260010162000092565b50620000f46004836020620003f3565b5062000104600882602062000490565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e052505083516200014d9350602e9250602085019150620004cf565b506200015986620002ef565b505050601791820b820b604090811b60a05290820b90910b901b60c052506200057495505050505050565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b039081169082168114620002eb57600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b031633146200034f576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d546001600160a01b036801000000000000000090910481169082168114620002eb57602d8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b6040518061040001604052806020906020820280368337509192915050565b6002830191839082156200047e5791602002820160005b838211156200044c57835183826101000a81548161ffff021916908361ffff16021790555092602001926002016020816001010492830192600103026200040a565b80156200047c5782816101000a81549061ffff02191690556002016020816001010492830192600103026200044c565b505b506200048c92915062000541565b5090565b8260208101928215620004c1579160200282015b82811115620004c1578251825591602001919060010190620004a4565b506200048c9291506200055d565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200051257805160ff1916838001178555620004c1565b82800160010185558215620004c15791820182811115620004c1578251825591602001919060010190620004a4565b5b808211156200048c57805461ffff1916815560010162000542565b5b808211156200048c57600081556001016200055e565b60805160601c60a05160401c60c05160401c60e05160f81c6150ba620005d960003980610ee95250806118815280613133525080610ec55280613106525080610ea152806125ce528061269052806135d65280613df0528061415152506150ba6000f3fe608060405234801561001057600080fd5b50600436106102415760003560e01c80638205bf6a11610145578063c1075329116100bd578063e5fe45771161008c578063f2fde38b11610071578063f2fde38b14610bb5578063fbffd2c114610be8578063feaf968c14610c1b57610241565b8063e5fe457714610b10578063eb5dcd6c14610b7a57610241565b8063c107532914610971578063c9807539146109aa578063d09dc33914610abe578063e4902f8214610ac657610241565b80639c849b3011610114578063b5ab58dc116100f9578063b5ab58dc146108f2578063b633620c1461090f578063bd8247061461092c57610241565b80639c849b30146107fd578063b121e147146108bf57610241565b80638205bf6a146107475780638ac28d5a1461074f5780638da5cb5b146107825780639a6fc8f51461078a57610241565b806354fd4d50116101d85780637284e416116101a757806379ba50971161018c57806379ba509714610696578063814118341461069e57806381ff7048146106f657610241565b80637284e4161461061157806373f666f81461068e57610241565b806354fd4d50146104cc578063585aa7de146104d4578063668a0f021461060157806370da2f671461060957610241565b8063313ce56711610214578063313ce567146103105780633a5381b51461032e5780634428ef6c1461033657806350d25bcd146104c457610241565b80630eafb25b146102465780631327d3d81461028b5780631b6b6d23146102c057806322adbc78146102f1575b600080fd5b6102796004803603602081101561025c57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610c23565b60408051918252519081900360200190f35b6102be600480360360208110156102a157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610d65565b005b6102c8610e9f565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6102f9610ec3565b6040805160179290920b8252519081900360200190f35b610318610ee7565b6040805160ff9092168252519081900360200190f35b6102c8610f0b565b61048f600480360360e081101561034c57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169167ffffffffffffffff6020820135169181019060608101604082013564010000000081111561039357600080fd5b8201836020820111156103a557600080fd5b803590602001918460208302840111640100000000831117156103c757600080fd5b9193909290916020810190356401000000008111156103e557600080fd5b8201836020820111156103f757600080fd5b8035906020019184602083028401116401000000008311171561041957600080fd5b9193909260ff8335169267ffffffffffffffff60208201351692919060608101906040013564010000000081111561045057600080fd5b82018360208201111561046257600080fd5b8035906020019184600183028401116401000000008311171561048457600080fd5b509092509050610f33565b604080517fffffffffffffffffffffffffffffffff000000000000000000000000000000009092168252519081900360200190f35b610279610f56565b610279610f92565b6102be600480360360a08110156104ea57600080fd5b81019060208101813564010000000081111561050557600080fd5b82018360208201111561051757600080fd5b8035906020019184602083028401116401000000008311171561053957600080fd5b91939092909160208101903564010000000081111561055757600080fd5b82018360208201111561056957600080fd5b8035906020019184602083028401116401000000008311171561058b57600080fd5b9193909260ff8335169267ffffffffffffffff6020820135169291906060810190604001356401000000008111156105c257600080fd5b8201836020820111156105d457600080fd5b803590602001918460018302840111640100000000831117156105f657600080fd5b509092509050610f97565b610279611859565b6102f961187f565b6106196118a3565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561065357818101518382015260200161063b565b50505050905090810190601f1680156106805780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610279611957565b6102be61195d565b6106a6611a5f565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156106e25781810151838201526020016106ca565b505050509050019250505060405180910390f35b6106fe611acd565b6040805163ffffffff94851681529290931660208301527fffffffffffffffffffffffffffffffff00000000000000000000000000000000168183015290519081900360600190f35b610279611aee565b6102be6004803603602081101561076557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611b49565b6102c8611bea565b6107b3600480360360208110156107a057600080fd5b503569ffffffffffffffffffff16611c06565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102be6004803603604081101561081357600080fd5b81019060208101813564010000000081111561082e57600080fd5b82018360208201111561084057600080fd5b8035906020019184602083028401116401000000008311171561086257600080fd5b91939092909160208101903564010000000081111561088057600080fd5b82018360208201111561089257600080fd5b803590602001918460208302840111640100000000831117156108b457600080fd5b509092509050611d74565b6102be600480360360208110156108d557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16612070565b6102796004803603602081101561090857600080fd5b503561219d565b6102796004803603602081101561092557600080fd5b50356121d3565b6102be600480360360a081101561094257600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516612228565b6102be6004803603604081101561098757600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135612440565b6102be600480360360808110156109c057600080fd5b8101906020810181356401000000008111156109db57600080fd5b8201836020820111156109ed57600080fd5b80359060200191846001830284011164010000000083111715610a0f57600080fd5b919390929091602081019035640100000000811115610a2d57600080fd5b820183602082011115610a3f57600080fd5b80359060200191846020830284011164010000000083111715610a6157600080fd5b919390929091602081019035640100000000811115610a7f57600080fd5b820183602082011115610a9157600080fd5b80359060200191846020830284011164010000000083111715610ab357600080fd5b9193509150356127c7565b610279613586565b610af960048036036020811015610adc57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661364e565b6040805161ffff9092168252519081900360200190f35b610b186136e3565b604080517fffffffffffffffffffffffffffffffff00000000000000000000000000000000909616865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b6102be60048036036040811015610b9057600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160200135166137ec565b6102be60048036036020811015610bcb57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166139b0565b6102be60048036036020811015610bfe57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613aac565b6107b3613b3b565b6000610c2d614e85565b5073ffffffffffffffffffffffffffffffffffffffff821660009081526028602090815260409182902082518084019093525460ff8082168452610100909104161515908201819052610c84576000915050610d60565b610c8c614e9c565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216948401949094526c0100000000000000000000000082048116606084018190527001000000000000000000000000000000009092041660808301528351919260009260019160049160ff16908110610d1957fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff1660208110610d5757fe5b01540301925050505b919050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610deb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d5473ffffffffffffffffffffffffffffffffffffffff6801000000000000000090910481169082168114610e9b57602d80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff166801000000000000000073ffffffffffffffffffffffffffffffffffffffff85811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b6000610f478b8b8b8b8b8b8b8b8b8b613bda565b9b9a5050505050505050505050565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c6020526040902054601790810b900b90565b600481565b868585818314610ff2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260248152602001806150616024913960400191505060405180910390fd5b8060030260ff16831161106657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b60005473ffffffffffffffffffffffffffffffffffffffff1633146110ec57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602954156112b757602980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101916000918390811061112957fe5b6000918252602082200154602a805473ffffffffffffffffffffffffffffffffffffffff9092169350908490811061115d57fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16905061118a81613d27565b73ffffffffffffffffffffffffffffffffffffffff80831660009081526028602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000908116909155928416825290208054909116905560298054806111f357fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055602a80548061125657fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055506110ec915050565b60005b8a8110156115d4576112ca614e85565b506040805180820190915260ff821681526001602082015280602860008f8f868181106112f357fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff168352508181019290925260400160009081208351815494909301517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0090941660ff909316929092177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100931515939093029290921790556006818d8d8681811061139c57fe5b73ffffffffffffffffffffffffffffffffffffffff6020918202939093013583168452830193909352604090910160002054169190911415905061144157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b80602860008d8d8681811061145257fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff16835250818101929092526040016000208251815493909201517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0090931660ff909216919091177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101009215159290920291909117905560298d8d848181106114f857fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9590920293909301359390931692909217905550602a8b8b8481811061156757fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9690930294909401359490941617909155509190910190506112ba565b50602b805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602d80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff8416178083166001810184167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009092169190911790945590910416906116a130828f8f8f8f8f8f8f8f613bda565b602b60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602b60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528681526020019050868680828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b602b54760100000000000000000000000000000000000000000000900463ffffffff1690565b7f000000000000000000000000000000000000000000000000000000000000000081565b602e8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61010060018816150201909516949094049384018190048102820181019092528281526060939092909183018282801561194d5780601f106119225761010080835404028352916020019161194d565b820191906000526020600020905b81548152906001019060200180831161193057829003601f168201915b5050505050905090565b6159d881565b60015473ffffffffffffffffffffffffffffffffffffffff1633146119e357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060602a80548060200260200160405190810160405280929190818152602001828054801561194d57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611a99575050505050905090565b602d54602b5463ffffffff808316926401000000009004169060801b909192565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260066020526040902054163314611bde57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b611be781613d27565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e74000000000000000000000000000000000081525090611cfb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015611cc0578181015183820152602001611ca8565b50505050905090810190601f168015611ced5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50611d04614e85565b5050505063ffffffff83166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b60005473ffffffffffffffffffffffffffffffffffffffff163314611dfa57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b828114611e6857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015612069576000858583818110611e8157fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1690506000848484818110611eae57fe5b73ffffffffffffffffffffffffffffffffffffffff8581166000908152600660209081526040909120549202939093013583169350909116905080158080611f2157508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b611f8c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614612059578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505060019092019150611e6b9050565b5050505050565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526007602052604090205416331461210557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600063ffffffff8211156121b357506000610d60565b5063ffffffff166000908152602c6020526040902054601790810b900b90565b600063ffffffff8211156121e957506000610d60565b5063ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b60035473ffffffffffffffffffffffffffffffffffffffff16806122ad57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561236a57600080fd5b505afa15801561237e573d6000803e3d6000fd5b505050506040513d602081101561239457600080fd5b5051806123b8575060005473ffffffffffffffffffffffffffffffffffffffff1633145b61242357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b61242b613f7a565b61243886868686866142e5565b505050505050565b600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561250557600080fd5b505afa158015612519573d6000803e3d6000fd5b505050506040513d602081101561252f57600080fd5b505180612553575060005473ffffffffffffffffffffffffffffffffffffffff1633145b6125be57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b60006125c861445f565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561265357600080fd5b505afa158015612667573d6000803e3d6000fd5b505050506040513d602081101561267d57600080fd5b505190508181101561268e57600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb856126d78585038761464d565b6040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561272a57600080fd5b505af115801561273e573d6000803e3d6000fd5b505050506040513d602081101561275457600080fd5b50516127c157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a90506127d4614eca565b6000806128168b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061466792505050565b6040880152608087019290925291855290925090508767ffffffffffffffff8111801561284257600080fd5b506040519080825280601f01601f19166020018201604052801561286d576020820181803683370190505b50606084015260005b60ff81168911156128de57858160ff166020811061289057fe5b1a60f81b84606001518260ff16815181106128a757fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101612876565b5082604001515167ffffffffffffffff811180156128fb57600080fd5b506040519080825280601f01601f191660200182016040528015612926576020820181803683370190505b506020840152612934614efe565b60005b8460400151518160ff161015612a54576000838260ff166020811061295857fe5b1a905082816020811061296757fe5b6020020151156129d857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6f6273657276657220696e646578207265706561746564210000000000000000604482015290519081900360640190fd5b6001838260ff16602081106129e957fe5b91151560209283029190910152849060ff8416908110612a0557fe5b1a60f81b86602001518360ff1681518110612a1c57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101612937565b50612a5d614e85565b503360009081526028602090815260409182902082518084019093525460ff808216845261010090910416158015918301829052612acf5750602a816000015160ff1681548110612aaa57fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b612b3a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b84516040015160ff168a11612bb057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b60208a1115612c2057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b878a14612c8e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b846040015151600211158015612caa5750602085604001515111155b612d1557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b84600001516040015160020260ff1685604001515111612d9657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b5050825164ffffffffff909216602090920191909152506040516000908a908a908083838082843760405192018290039091209450612dd99350614efe92505050565b612de1614e85565b60005b8981101561301257600060018587606001518481518110612e0157fe5b60209101015160f81c601b018e8e86818110612e1957fe5b905060200201358d8d87818110612e2c57fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015612e87573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526028602090815290849020838501909452925460ff808216845261010090910416151592820183905290955092509050612f6857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f7261636c6520776173206e6f742070726573656e7421000000000000000000604482015290519081900360640190fd5b8251849060ff1660208110612f7957fe5b602002015115612fea57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff1660208110612fff57fe5b9115156020909202015250600101612de4565b5050505060005b6001826040015151038110156130dd5760008260400151826001018151811061303e57fe5b602002602001015160170b8360400151838151811061305957fe5b602002602001015160170b13159050806130d457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b50600101613019565b506040810151805160009190600281049081106130f657fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b1315801561315c57507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b6131c757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602c82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519985168a52898801959095523397890188905260a0890185905260c09989018a815283519a8a019a909a5282517ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d4519a94998c99989497939695909490939285019260e08601928983019202908190849084905b8381101561332d578181015183820152602001613315565b50505050905001838103825285818151815260200191508051906020019080838360005b83811015613369578181015183820152602001613351565b50505050905090810190601f1680156133965780820380516001836020036101000a031916815260200191505b509850505050505050505060405180910390a181516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a361344c8260000151606001518260170b6148fd565b5080518051602b8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff9283160217909155821061356d57fe5b61357b828260200151614a28565b505050505050505050565b600061359061445f565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016916370a08231916024808301926020929190829003018186803b15801561361c57600080fd5b505afa158015613630573d6000803e3d6000fd5b505050506040513d602081101561364657600080fd5b505103905090565b6000613658614e85565b5073ffffffffffffffffffffffffffffffffffffffff821660009081526028602090815260409182902082518084019093525460ff80821684526101009091041615159082018190526136af576000915050610d60565b805160049060ff16602081106136c157fe5b601091828204019190066002029054906101000a900461ffff16915050919050565b60008080808033321461375757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602b5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602c6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526006602052604090205416331461388157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b3373ffffffffffffffffffffffffffffffffffffffff8216141561390657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff0000000000000000000000000000000000000000821681179092559091169081146139ab5760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314613a3657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314613b3257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b611be781614c79565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000808080613b6b614e85565b5050505063ffffffff82166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b60008a8a8a8a8a8a8a8a8a8a604051602001808b73ffffffffffffffffffffffffffffffffffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b613d2f614e85565b5073ffffffffffffffffffffffffffffffffffffffff8116600090815260286020908152604080832081518083019092525460ff80821683526101009091041615159181019190915290613d8283610c23565b905080156139ab5773ffffffffffffffffffffffffffffffffffffffff80841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b158015613e3857600080fd5b505af1158015613e4c573d6000803e3d6000fd5b505050506040513d6020811015613e6257600080fd5b5051613ecf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff1660208110613ee557fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff1660208110613f2057fe5b01556040805173ffffffffffffffffffffffffffffffffffffffff80871682528316602082015281517fe403f9d01c1c84f7381a468229c5bdb3f454c492688c83a6cc62fcfd14811610929181900390910190a150505050565b613f82614e9c565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152613ff9614efe565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff168152602001906002019060208260010104928301926001038202915080841161401257905050505050509050614059614efe565b604080516104008101918290529060089060209082845b81548152602001906001019080831161407057505050505090506060602a8054806020026020016040519081016040528092919081815260200182805480156140ef57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116140c4575b5050505050905060005b81518110156142c9576000600184836020811061411257fe5b60200201510390506000600186846020811061412a57fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb86868151811061419757fe5b6020026020010151836040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156141f357600080fd5b505af1158015614207573d6000803e3d6000fd5b505050506040513d602081101561421d57600080fd5b505161428a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b600187856020811061429857fe5b602002019061ffff16908161ffff168152505060018685602081106142b957fe5b60200201525050506001016140f9565b506142d76004846020614f1d565b506120696008836020614fb3565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000614469614efe565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116144825790505050505050905060005b60208110156144f25760018282602081106144db57fe5b60200201510361ffff1692909201916001016144c4565b506144fb614e9c565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608086018290527001000000000000000000000000000000009094049092166080850152602a805486518184028101840190975280875297909202633b9aca00029693949293908301828280156145d857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116145ad575b505050505090506145e7614efe565b604080516104008101918290529060089060209082845b8154815260200190600101908083116145fe575050505050905060005b825181101561464557600182826020811061463257fe5b602002015103959095019460010161461b565b505050505090565b60008183101561465e575081614661565b50805b92915050565b61466f614fed565b506040805160808082018352602b549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff166020808401919091527501000000000000000000000000000000000000000000820460ff1693830193909352760100000000000000000000000000000000000000000000900463ffffffff1660608083019190915283519192600092839283929091908701908281101561473457600080fd5b8151602083015160408085018051915193959294830192918464010000000082111561475f57600080fd5b90830190602082018581111561477457600080fd5b825186602082028301116401000000008211171561479157600080fd5b82525081516020918201928201910280838360005b838110156147be5781810151838201526020016147a6565b5050505091909101604052505089519498508897509295509093505050605884901b907fffffffffffffffffffffffffffffffff0000000000000000000000000000000080831691161461487357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b8364ffffffffff16866020015164ffffffffff16106148f357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b5091939590929450565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff168061492d5750610e9b565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602c602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b602485018190529489166044850152606484018890525173ffffffffffffffffffffffffffffffffffffffff87169363beed9b5193620186a09360848084019491939192918390030190829088803b1580156149f957600080fd5b5087f193505050508015614a1f57506040513d6020811015614a1a57600080fd5b505160015b61243857612069565b614a30614e85565b503360009081526028602090815260409182902082518084019093525460ff808216845261010090910416151590820152614a69614e9c565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606085015270010000000000000000000000000000000090920416608083015282516104008101938490529192614b38928692909160049190826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614af65790505050505050614d22565b614b46906004906020614f1d565b508160200151614bb757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b6000614bde633b9aca003a04836020015163ffffffff16846000015163ffffffff16614d9a565b90506010360260005a90506000614bfd8863ffffffff16858585614dc0565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff16830281614c2b57fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff1660208110614c5457fe5b015401016008886000015160ff1660208110614c6c57fe5b0155505050505050505050565b60035473ffffffffffffffffffffffffffffffffffffffff9081169082168114610e9b57600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b614d2a614efe565b60005b8351811015614d92576000848281518110614d4457fe5b0160209081015160f81c9150614d6c90859083908110614d6057fe5b60200201516001614e66565b848260ff1660208110614d7b57fe5b61ffff909216602092909202015250600101614d2d565b509092915050565b60008383811015614dad57600285850304015b614db7818461464d565b95945050505050565b600081851015614e3157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b81850383016159d801633b9aca00858202026fffffffffffffffffffffffffffffffff8110614e5c57fe5b9695505050505050565b6000614e7e8261ffff168461ffff160161ffff61464d565b9392505050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518060a00160405280614edd614fed565b81526060602082018190526040820181905280820152600060809091015290565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215614fa35791602002820160005b83821115614f7357835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302614f33565b8015614fa15782816101000a81549061ffff0219169055600201602081600101049283019260010302614f73565b505b50614faf929150615014565b5090565b8260208101928215614fe1579160200282015b82811115614fe1578251825591602001919060010190614fc6565b50614faf92915061504b565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b80821115614faf5780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000168155600101615015565b5b80821115614faf576000815560010161504c56fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


func DeployExposedOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ExposedOffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(ExposedOffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ExposedOffchainAggregatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExposedOffchainAggregator{ExposedOffchainAggregatorCaller: ExposedOffchainAggregatorCaller{contract: contract}, ExposedOffchainAggregatorTransactor: ExposedOffchainAggregatorTransactor{contract: contract}, ExposedOffchainAggregatorFilterer: ExposedOffchainAggregatorFilterer{contract: contract}}, nil
}


type ExposedOffchainAggregator struct {
	ExposedOffchainAggregatorCaller     
	ExposedOffchainAggregatorTransactor 
	ExposedOffchainAggregatorFilterer   
}


type ExposedOffchainAggregatorCaller struct {
	contract *bind.BoundContract 
}


type ExposedOffchainAggregatorTransactor struct {
	contract *bind.BoundContract 
}


type ExposedOffchainAggregatorFilterer struct {
	contract *bind.BoundContract 
}



type ExposedOffchainAggregatorSession struct {
	Contract     *ExposedOffchainAggregator 
	CallOpts     bind.CallOpts              
	TransactOpts bind.TransactOpts          
}



type ExposedOffchainAggregatorCallerSession struct {
	Contract *ExposedOffchainAggregatorCaller 
	CallOpts bind.CallOpts                    
}



type ExposedOffchainAggregatorTransactorSession struct {
	Contract     *ExposedOffchainAggregatorTransactor 
	TransactOpts bind.TransactOpts                    
}


type ExposedOffchainAggregatorRaw struct {
	Contract *ExposedOffchainAggregator 
}


type ExposedOffchainAggregatorCallerRaw struct {
	Contract *ExposedOffchainAggregatorCaller 
}


type ExposedOffchainAggregatorTransactorRaw struct {
	Contract *ExposedOffchainAggregatorTransactor 
}


func NewExposedOffchainAggregator(address common.Address, backend bind.ContractBackend) (*ExposedOffchainAggregator, error) {
	contract, err := bindExposedOffchainAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregator{ExposedOffchainAggregatorCaller: ExposedOffchainAggregatorCaller{contract: contract}, ExposedOffchainAggregatorTransactor: ExposedOffchainAggregatorTransactor{contract: contract}, ExposedOffchainAggregatorFilterer: ExposedOffchainAggregatorFilterer{contract: contract}}, nil
}


func NewExposedOffchainAggregatorCaller(address common.Address, caller bind.ContractCaller) (*ExposedOffchainAggregatorCaller, error) {
	contract, err := bindExposedOffchainAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorCaller{contract: contract}, nil
}


func NewExposedOffchainAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*ExposedOffchainAggregatorTransactor, error) {
	contract, err := bindExposedOffchainAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorTransactor{contract: contract}, nil
}


func NewExposedOffchainAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*ExposedOffchainAggregatorFilterer, error) {
	contract, err := bindExposedOffchainAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorFilterer{contract: contract}, nil
}


func bindExposedOffchainAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExposedOffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}





func (_ExposedOffchainAggregator *ExposedOffchainAggregatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExposedOffchainAggregator.Contract.ExposedOffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}



func (_ExposedOffchainAggregator *ExposedOffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.ExposedOffchainAggregatorTransactor.contract.Transfer(opts)
}


func (_ExposedOffchainAggregator *ExposedOffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.ExposedOffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}





func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExposedOffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}



func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.contract.Transfer(opts)
}


func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.contract.Transact(opts, method, params...)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "LINK")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) LINK() (common.Address, error) {
	return _ExposedOffchainAggregator.Contract.LINK(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _ExposedOffchainAggregator.Contract.LINK(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) AccountingGasCost(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "accountingGasCost")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) AccountingGasCost() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.AccountingGasCost(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) AccountingGasCost() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.AccountingGasCost(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "decimals")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) Decimals() (uint8, error) {
	return _ExposedOffchainAggregator.Contract.Decimals(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _ExposedOffchainAggregator.Contract.Decimals(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "description")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) Description() (string, error) {
	return _ExposedOffchainAggregator.Contract.Description(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) Description() (string, error) {
	return _ExposedOffchainAggregator.Contract.Description(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) ExposedConfigDigestFromConfigData(opts *bind.CallOpts, _contractAddress common.Address, _configCount uint64, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encodedConfig []byte) ([16]byte, error) {
	var (
		ret0 = new([16]byte)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "exposedConfigDigestFromConfigData", _contractAddress, _configCount, _signers, _transmitters, _threshold, _encodedConfigVersion, _encodedConfig)
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) ExposedConfigDigestFromConfigData(_contractAddress common.Address, _configCount uint64, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encodedConfig []byte) ([16]byte, error) {
	return _ExposedOffchainAggregator.Contract.ExposedConfigDigestFromConfigData(&_ExposedOffchainAggregator.CallOpts, _contractAddress, _configCount, _signers, _transmitters, _threshold, _encodedConfigVersion, _encodedConfig)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) ExposedConfigDigestFromConfigData(_contractAddress common.Address, _configCount uint64, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encodedConfig []byte) ([16]byte, error) {
	return _ExposedOffchainAggregator.Contract.ExposedConfigDigestFromConfigData(&_ExposedOffchainAggregator.CallOpts, _contractAddress, _configCount, _signers, _transmitters, _threshold, _encodedConfigVersion, _encodedConfig)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "getAnswer", _roundId)
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.GetAnswer(&_ExposedOffchainAggregator.CallOpts, _roundId)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.GetAnswer(&_ExposedOffchainAggregator.CallOpts, _roundId)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (struct {
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
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "getRoundData", _roundId)
	return *ret, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _ExposedOffchainAggregator.Contract.GetRoundData(&_ExposedOffchainAggregator.CallOpts, _roundId)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _ExposedOffchainAggregator.Contract.GetRoundData(&_ExposedOffchainAggregator.CallOpts, _roundId)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "getTimestamp", _roundId)
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.GetTimestamp(&_ExposedOffchainAggregator.CallOpts, _roundId)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.GetTimestamp(&_ExposedOffchainAggregator.CallOpts, _roundId)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "latestAnswer")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.LatestAnswer(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.LatestAnswer(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
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
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "latestConfigDetails")
	return *ret, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _ExposedOffchainAggregator.Contract.LatestConfigDetails(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _ExposedOffchainAggregator.Contract.LatestConfigDetails(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "latestRound")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) LatestRound() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.LatestRound(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.LatestRound(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
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
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "latestRoundData")
	return *ret, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _ExposedOffchainAggregator.Contract.LatestRoundData(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _ExposedOffchainAggregator.Contract.LatestRoundData(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "latestTimestamp")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.LatestTimestamp(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.LatestTimestamp(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (struct {
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
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "latestTransmissionDetails")
	return *ret, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _ExposedOffchainAggregator.Contract.LatestTransmissionDetails(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _ExposedOffchainAggregator.Contract.LatestTransmissionDetails(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "linkAvailableForPayment")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.LinkAvailableForPayment(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.LinkAvailableForPayment(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "maxAnswer")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.MaxAnswer(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.MaxAnswer(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "minAnswer")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.MinAnswer(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.MinAnswer(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "oracleObservationCount", _signerOrTransmitter)
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _ExposedOffchainAggregator.Contract.OracleObservationCount(&_ExposedOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _ExposedOffchainAggregator.Contract.OracleObservationCount(&_ExposedOffchainAggregator.CallOpts, _signerOrTransmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "owedPayment", _transmitter)
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.OwedPayment(&_ExposedOffchainAggregator.CallOpts, _transmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.OwedPayment(&_ExposedOffchainAggregator.CallOpts, _transmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "owner")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) Owner() (common.Address, error) {
	return _ExposedOffchainAggregator.Contract.Owner(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _ExposedOffchainAggregator.Contract.Owner(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "transmitters")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _ExposedOffchainAggregator.Contract.Transmitters(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _ExposedOffchainAggregator.Contract.Transmitters(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "validator")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) Validator() (common.Address, error) {
	return _ExposedOffchainAggregator.Contract.Validator(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _ExposedOffchainAggregator.Contract.Validator(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExposedOffchainAggregator.contract.Call(opts, out, "version")
	return *ret0, err
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) Version() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.Version(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorCallerSession) Version() (*big.Int, error) {
	return _ExposedOffchainAggregator.Contract.Version(&_ExposedOffchainAggregator.CallOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "acceptOwnership")
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.AcceptOwnership(&_ExposedOffchainAggregator.TransactOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.AcceptOwnership(&_ExposedOffchainAggregator.TransactOpts)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "acceptPayeeship", _transmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.AcceptPayeeship(&_ExposedOffchainAggregator.TransactOpts, _transmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.AcceptPayeeship(&_ExposedOffchainAggregator.TransactOpts, _transmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetBilling(&_ExposedOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetBilling(&_ExposedOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "setBillingAccessController", _billingAdminAccessController)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetBillingAccessController(&_ExposedOffchainAggregator.TransactOpts, _billingAdminAccessController)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) SetBillingAccessController(_billingAdminAccessController common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetBillingAccessController(&_ExposedOffchainAggregator.TransactOpts, _billingAdminAccessController)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetConfig(&_ExposedOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetConfig(&_ExposedOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "setPayees", _transmitters, _payees)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetPayees(&_ExposedOffchainAggregator.TransactOpts, _transmitters, _payees)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetPayees(&_ExposedOffchainAggregator.TransactOpts, _transmitters, _payees)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) SetValidator(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "setValidator", _newValidator)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetValidator(&_ExposedOffchainAggregator.TransactOpts, _newValidator)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.SetValidator(&_ExposedOffchainAggregator.TransactOpts, _newValidator)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "transferOwnership", _to)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.TransferOwnership(&_ExposedOffchainAggregator.TransactOpts, _to)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.TransferOwnership(&_ExposedOffchainAggregator.TransactOpts, _to)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.TransferPayeeship(&_ExposedOffchainAggregator.TransactOpts, _transmitter, _proposed)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.TransferPayeeship(&_ExposedOffchainAggregator.TransactOpts, _transmitter, _proposed)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) Transmit(opts *bind.TransactOpts, _report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "transmit", _report, _rs, _ss, _rawVs)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.Transmit(&_ExposedOffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.Transmit(&_ExposedOffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.WithdrawFunds(&_ExposedOffchainAggregator.TransactOpts, _recipient, _amount)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.WithdrawFunds(&_ExposedOffchainAggregator.TransactOpts, _recipient, _amount)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.contract.Transact(opts, "withdrawPayment", _transmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.WithdrawPayment(&_ExposedOffchainAggregator.TransactOpts, _transmitter)
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _ExposedOffchainAggregator.Contract.WithdrawPayment(&_ExposedOffchainAggregator.TransactOpts, _transmitter)
}


type ExposedOffchainAggregatorAnswerUpdatedIterator struct {
	Event *ExposedOffchainAggregatorAnswerUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorAnswerUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorAnswerUpdated)
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
		it.Event = new(ExposedOffchainAggregatorAnswerUpdated)
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


func (it *ExposedOffchainAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*ExposedOffchainAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorAnswerUpdatedIterator{contract: _ExposedOffchainAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorAnswerUpdated)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*ExposedOffchainAggregatorAnswerUpdated, error) {
	event := new(ExposedOffchainAggregatorAnswerUpdated)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorBillingAccessControllerSetIterator struct {
	Event *ExposedOffchainAggregatorBillingAccessControllerSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorBillingAccessControllerSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorBillingAccessControllerSet)
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
		it.Event = new(ExposedOffchainAggregatorBillingAccessControllerSet)
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


func (it *ExposedOffchainAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*ExposedOffchainAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorBillingAccessControllerSetIterator{contract: _ExposedOffchainAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorBillingAccessControllerSet)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*ExposedOffchainAggregatorBillingAccessControllerSet, error) {
	event := new(ExposedOffchainAggregatorBillingAccessControllerSet)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorBillingSetIterator struct {
	Event *ExposedOffchainAggregatorBillingSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorBillingSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorBillingSet)
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
		it.Event = new(ExposedOffchainAggregatorBillingSet)
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


func (it *ExposedOffchainAggregatorBillingSetIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*ExposedOffchainAggregatorBillingSetIterator, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorBillingSetIterator{contract: _ExposedOffchainAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorBillingSet)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseBillingSet(log types.Log) (*ExposedOffchainAggregatorBillingSet, error) {
	event := new(ExposedOffchainAggregatorBillingSet)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorConfigSetIterator struct {
	Event *ExposedOffchainAggregatorConfigSet 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorConfigSetIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorConfigSet)
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
		it.Event = new(ExposedOffchainAggregatorConfigSet)
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


func (it *ExposedOffchainAggregatorConfigSetIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*ExposedOffchainAggregatorConfigSetIterator, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorConfigSetIterator{contract: _ExposedOffchainAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorConfigSet)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseConfigSet(log types.Log) (*ExposedOffchainAggregatorConfigSet, error) {
	event := new(ExposedOffchainAggregatorConfigSet)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorNewRoundIterator struct {
	Event *ExposedOffchainAggregatorNewRound 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorNewRoundIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorNewRound)
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
		it.Event = new(ExposedOffchainAggregatorNewRound)
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


func (it *ExposedOffchainAggregatorNewRoundIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*ExposedOffchainAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorNewRoundIterator{contract: _ExposedOffchainAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorNewRound)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseNewRound(log types.Log) (*ExposedOffchainAggregatorNewRound, error) {
	event := new(ExposedOffchainAggregatorNewRound)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorNewTransmissionIterator struct {
	Event *ExposedOffchainAggregatorNewTransmission 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorNewTransmissionIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorNewTransmission)
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
		it.Event = new(ExposedOffchainAggregatorNewTransmission)
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


func (it *ExposedOffchainAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	RawReportContext  [32]byte
	Raw               types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts) (*ExposedOffchainAggregatorNewTransmissionIterator, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "NewTransmission")
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorNewTransmissionIterator{contract: _ExposedOffchainAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorNewTransmission) (event.Subscription, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "NewTransmission")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorNewTransmission)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseNewTransmission(log types.Log) (*ExposedOffchainAggregatorNewTransmission, error) {
	event := new(ExposedOffchainAggregatorNewTransmission)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorOraclePaidIterator struct {
	Event *ExposedOffchainAggregatorOraclePaid 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorOraclePaidIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorOraclePaid)
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
		it.Event = new(ExposedOffchainAggregatorOraclePaid)
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


func (it *ExposedOffchainAggregatorOraclePaidIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Raw         types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts) (*ExposedOffchainAggregatorOraclePaidIterator, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorOraclePaidIterator{contract: _ExposedOffchainAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorOraclePaid) (event.Subscription, error) {

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorOraclePaid)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseOraclePaid(log types.Log) (*ExposedOffchainAggregatorOraclePaid, error) {
	event := new(ExposedOffchainAggregatorOraclePaid)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorOwnershipTransferRequestedIterator struct {
	Event *ExposedOffchainAggregatorOwnershipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorOwnershipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorOwnershipTransferRequested)
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
		it.Event = new(ExposedOffchainAggregatorOwnershipTransferRequested)
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


func (it *ExposedOffchainAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExposedOffchainAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorOwnershipTransferRequestedIterator{contract: _ExposedOffchainAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorOwnershipTransferRequested)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*ExposedOffchainAggregatorOwnershipTransferRequested, error) {
	event := new(ExposedOffchainAggregatorOwnershipTransferRequested)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorOwnershipTransferredIterator struct {
	Event *ExposedOffchainAggregatorOwnershipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorOwnershipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorOwnershipTransferred)
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
		it.Event = new(ExposedOffchainAggregatorOwnershipTransferred)
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


func (it *ExposedOffchainAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ExposedOffchainAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorOwnershipTransferredIterator{contract: _ExposedOffchainAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorOwnershipTransferred)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*ExposedOffchainAggregatorOwnershipTransferred, error) {
	event := new(ExposedOffchainAggregatorOwnershipTransferred)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorPayeeshipTransferRequestedIterator struct {
	Event *ExposedOffchainAggregatorPayeeshipTransferRequested 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorPayeeshipTransferRequestedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorPayeeshipTransferRequested)
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
		it.Event = new(ExposedOffchainAggregatorPayeeshipTransferRequested)
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


func (it *ExposedOffchainAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*ExposedOffchainAggregatorPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorPayeeshipTransferRequestedIterator{contract: _ExposedOffchainAggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorPayeeshipTransferRequested)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*ExposedOffchainAggregatorPayeeshipTransferRequested, error) {
	event := new(ExposedOffchainAggregatorPayeeshipTransferRequested)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorPayeeshipTransferredIterator struct {
	Event *ExposedOffchainAggregatorPayeeshipTransferred 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorPayeeshipTransferredIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorPayeeshipTransferred)
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
		it.Event = new(ExposedOffchainAggregatorPayeeshipTransferred)
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


func (it *ExposedOffchainAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*ExposedOffchainAggregatorPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorPayeeshipTransferredIterator{contract: _ExposedOffchainAggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorPayeeshipTransferred)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*ExposedOffchainAggregatorPayeeshipTransferred, error) {
	event := new(ExposedOffchainAggregatorPayeeshipTransferred)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}


type ExposedOffchainAggregatorValidatorUpdatedIterator struct {
	Event *ExposedOffchainAggregatorValidatorUpdated 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}




func (it *ExposedOffchainAggregatorValidatorUpdatedIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExposedOffchainAggregatorValidatorUpdated)
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
		it.Event = new(ExposedOffchainAggregatorValidatorUpdated)
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


func (it *ExposedOffchainAggregatorValidatorUpdatedIterator) Error() error {
	return it.fail
}



func (it *ExposedOffchainAggregatorValidatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}


type ExposedOffchainAggregatorValidatorUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log 
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) FilterValidatorUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*ExposedOffchainAggregatorValidatorUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.FilterLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &ExposedOffchainAggregatorValidatorUpdatedIterator{contract: _ExposedOffchainAggregator.contract, event: "ValidatorUpdated", logs: logs, sub: sub}, nil
}




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) WatchValidatorUpdated(opts *bind.WatchOpts, sink chan<- *ExposedOffchainAggregatorValidatorUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _ExposedOffchainAggregator.contract.WatchLogs(opts, "ValidatorUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(ExposedOffchainAggregatorValidatorUpdated)
				if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
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




func (_ExposedOffchainAggregator *ExposedOffchainAggregatorFilterer) ParseValidatorUpdated(log types.Log) (*ExposedOffchainAggregatorValidatorUpdated, error) {
	event := new(ExposedOffchainAggregatorValidatorUpdated)
	if err := _ExposedOffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
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


const OffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"availableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var OffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b50604051620055cf380380620055cf83398181016040526101808110156200003957600080fd5b815160208301516040808501516060860151608087015160a088015160c089015160e08a01516101008b01516101208c01516101408d01516101608e0180519a519c9e9b9d999c989b979a969995989497939692959194939182019284640100000000821115620000a957600080fd5b908301906020820185811115620000bf57600080fd5b8251640100000000811182820188101715620000da57600080fd5b82525081516020918201929091019080838360005b8381101562000109578181015183820152602001620000ef565b50505050905090810190601f168015620001375780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b03191633179055508b8b8b8b8b8b886200016687878787876200028c565b62000171816200037e565b6001600160601b0319606083901b166080526200018d620004dc565b62000197620004dc565b60005b60208160ff161015620001ec576001838260ff1660208110620001b957fe5b602002019061ffff16908161ffff16815250506001828260ff1660208110620001de57fe5b60200201526001016200019a565b50620001fc6004836020620004fb565b506200020c600882602062000598565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e05250508351620002559350602e9250602085019150620005d7565b506200026186620003f7565b505050601791820b820b604090811b60a05290820b90910b901b60c052506200067c95505050505050565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b039081169082168114620003f357600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b0316331462000457576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d546001600160a01b036801000000000000000090910481169082168114620003f357602d8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215620005865791602002820160005b838211156200055457835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030262000512565b8015620005845782816101000a81549061ffff021916905560020160208160010104928301926001030262000554565b505b506200059492915062000649565b5090565b8260208101928215620005c9579160200282015b82811115620005c9578251825591602001919060010190620005ac565b506200059492915062000665565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200061a57805160ff1916838001178555620005c9565b82800160010185558215620005c95791820182811115620005c9578251825591602001919060010190620005ac565b5b808211156200059457805461ffff191681556001016200064a565b5b8082111562000594576000815560010162000666565b60805160601c60a05160401c60c05160401c60e05160f81c614eee620006e160003980610d405250806116b55280612f67525080610d1c5280612f3a525080610cf8528061240252806124c4528061340a5280613ad75280613f855250614eee6000f3fe608060405234801561001057600080fd5b50600436106102265760003560e01c80638ac28d5a1161012a578063c1075329116100bd578063e5fe45771161008c578063f2fde38b11610071578063f2fde38b14610a0c578063fbffd2c114610a3f578063feaf968c14610a7257610226565b8063e5fe457714610967578063eb5dcd6c146109d157610226565b8063c1075329146107c8578063c980753914610801578063d09dc33914610915578063e4902f821461091d57610226565b8063b121e147116100f9578063b121e14714610716578063b5ab58dc14610749578063b633620c14610766578063bd8247061461078357610226565b80638ac28d5a146105a65780638da5cb5b146105d95780639a6fc8f5146105e15780639c849b301461065457610226565b8063585aa7de116101bd57806373f666f81161018c578063814118341161017157806381411834146104f557806381ff70481461054d5780638205bf6a1461059e57610226565b806373f666f8146104e557806379ba5097146104ed57610226565b8063585aa7de1461032b578063668a0f021461045857806370da2f67146104605780637284e4161461046857610226565b8063313ce567116101f9578063313ce567146102f55780633a5381b51461031357806350d25bcd1461031b57806354fd4d501461032357610226565b80630eafb25b1461022b5780631327d3d8146102705780631b6b6d23146102a557806322adbc78146102d6575b600080fd5b61025e6004803603602081101561024157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610a7a565b60408051918252519081900360200190f35b6102a36004803603602081101561028657600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610bbc565b005b6102ad610cf6565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6102de610d1a565b6040805160179290920b8252519081900360200190f35b6102fd610d3e565b6040805160ff9092168252519081900360200190f35b6102ad610d62565b61025e610d8a565b61025e610dc6565b6102a3600480360360a081101561034157600080fd5b81019060208101813564010000000081111561035c57600080fd5b82018360208201111561036e57600080fd5b8035906020019184602083028401116401000000008311171561039057600080fd5b9193909290916020810190356401000000008111156103ae57600080fd5b8201836020820111156103c057600080fd5b803590602001918460208302840111640100000000831117156103e257600080fd5b9193909260ff8335169267ffffffffffffffff60208201351692919060608101906040013564010000000081111561041957600080fd5b82018360208201111561042b57600080fd5b8035906020019184600183028401116401000000008311171561044d57600080fd5b509092509050610dcb565b61025e61168d565b6102de6116b3565b6104706116d7565b6040805160208082528351818301528351919283929083019185019080838360005b838110156104aa578181015183820152602001610492565b50505050905090810190601f1680156104d75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61025e61178b565b6102a3611791565b6104fd611893565b60408051602080825283518183015283519192839290830191858101910280838360005b83811015610539578181015183820152602001610521565b505050509050019250505060405180910390f35b610555611901565b6040805163ffffffff94851681529290931660208301527fffffffffffffffffffffffffffffffff00000000000000000000000000000000168183015290519081900360600190f35b61025e611922565b6102a3600480360360208110156105bc57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661197d565b6102ad611a1e565b61060a600480360360208110156105f757600080fd5b503569ffffffffffffffffffff16611a3a565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102a36004803603604081101561066a57600080fd5b81019060208101813564010000000081111561068557600080fd5b82018360208201111561069757600080fd5b803590602001918460208302840111640100000000831117156106b957600080fd5b9193909290916020810190356401000000008111156106d757600080fd5b8201836020820111156106e957600080fd5b8035906020019184602083028401116401000000008311171561070b57600080fd5b509092509050611ba8565b6102a36004803603602081101561072c57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611ea4565b61025e6004803603602081101561075f57600080fd5b5035611fd1565b61025e6004803603602081101561077c57600080fd5b5035612007565b6102a3600480360360a081101561079957600080fd5b5063ffffffff81358116916020810135821691604082013581169160608101358216916080909101351661205c565b6102a3600480360360408110156107de57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135612274565b6102a36004803603608081101561081757600080fd5b81019060208101813564010000000081111561083257600080fd5b82018360208201111561084457600080fd5b8035906020019184600183028401116401000000008311171561086657600080fd5b91939092909160208101903564010000000081111561088457600080fd5b82018360208201111561089657600080fd5b803590602001918460208302840111640100000000831117156108b857600080fd5b9193909290916020810190356401000000008111156108d657600080fd5b8201836020820111156108e857600080fd5b8035906020019184602083028401116401000000008311171561090a57600080fd5b9193509150356125fb565b61025e6133ba565b6109506004803603602081101561093357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16613482565b6040805161ffff9092168252519081900360200190f35b61096f613517565b604080517fffffffffffffffffffffffffffffffff00000000000000000000000000000000909616865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b6102a3600480360360408110156109e757600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516613620565b6102a360048036036020811015610a2257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166137e4565b6102a360048036036020811015610a5557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166138e0565b61060a61396f565b6000610a84614cb9565b5073ffffffffffffffffffffffffffffffffffffffff821660009081526028602090815260409182902082518084019093525460ff8082168452610100909104161515908201819052610adb576000915050610bb7565b610ae3614cd0565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216948401949094526c0100000000000000000000000082048116606084018190527001000000000000000000000000000000009092041660808301528351919260009260019160049160ff16908110610b7057fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff1660208110610bae57fe5b01540301925050505b919050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610c4257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d5473ffffffffffffffffffffffffffffffffffffffff6801000000000000000090910481169082168114610cf257602d80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff166801000000000000000073ffffffffffffffffffffffffffffffffffffffff85811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c6020526040902054601790810b900b90565b600481565b868585818314610e26576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526024815260200180614e956024913960400191505060405180910390fd5b8060030260ff168311610e9a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b60005473ffffffffffffffffffffffffffffffffffffffff163314610f2057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602954156110eb57602980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81019160009183908110610f5d57fe5b6000918252602082200154602a805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110610f9157fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff169050610fbe81613a0e565b73ffffffffffffffffffffffffffffffffffffffff80831660009081526028602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559284168252902080549091169055602980548061102757fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055602a80548061108a57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550610f20915050565b60005b8a811015611408576110fe614cb9565b506040805180820190915260ff821681526001602082015280602860008f8f8681811061112757fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff168352508181019290925260400160009081208351815494909301517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0090941660ff909316929092177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100931515939093029290921790556006818d8d868181106111d057fe5b73ffffffffffffffffffffffffffffffffffffffff6020918202939093013583168452830193909352604090910160002054169190911415905061127557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b80602860008d8d8681811061128657fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff16835250818101929092526040016000208251815493909201517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0090931660ff909216919091177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101009215159290920291909117905560298d8d8481811061132c57fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9590920293909301359390931692909217905550602a8b8b8481811061139b57fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9690930294909401359490941617909155509190910190506110ee565b50602b805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602d80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff8416178083166001810184167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009092169190911790945590910416906114d530828f8f8f8f8f8f8f8f613c61565b602b60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602b60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528681526020019050868680828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b602b54760100000000000000000000000000000000000000000000900463ffffffff1690565b7f000000000000000000000000000000000000000000000000000000000000000081565b602e8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156117815780601f1061175657610100808354040283529160200191611781565b820191906000526020600020905b81548152906001019060200180831161176457829003601f168201915b5050505050905090565b6159d881565b60015473ffffffffffffffffffffffffffffffffffffffff16331461181757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060602a80548060200260200160405190810160405280929190818152602001828054801561178157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116118cd575050505050905090565b602d54602b5463ffffffff808316926401000000009004169060801b909192565b602b54760100000000000000000000000000000000000000000000900463ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260066020526040902054163314611a1257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b611a1b81613a0e565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e74000000000000000000000000000000000081525090611b2f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015611af4578181015183820152602001611adc565b50505050905090810190601f168015611b215780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50611b38614cb9565b5050505063ffffffff83166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b60005473ffffffffffffffffffffffffffffffffffffffff163314611c2e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b828114611c9c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015611e9d576000858583818110611cb557fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1690506000848484818110611ce257fe5b73ffffffffffffffffffffffffffffffffffffffff8581166000908152600660209081526040909120549202939093013583169350909116905080158080611d5557508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b611dc057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614611e8d578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505060019092019150611c9f9050565b5050505050565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260076020526040902054163314611f3957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600063ffffffff821115611fe757506000610bb7565b5063ffffffff166000908152602c6020526040902054601790810b900b90565b600063ffffffff82111561201d57506000610bb7565b5063ffffffff166000908152602c60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b60035473ffffffffffffffffffffffffffffffffffffffff16806120e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561219e57600080fd5b505afa1580156121b2573d6000803e3d6000fd5b505050506040513d60208110156121c857600080fd5b5051806121ec575060005473ffffffffffffffffffffffffffffffffffffffff1633145b61225757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b61225f613dae565b61226c8686868686614119565b505050505050565b600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561233957600080fd5b505afa15801561234d573d6000803e3d6000fd5b505050506040513d602081101561236357600080fd5b505180612387575060005473ffffffffffffffffffffffffffffffffffffffff1633145b6123f257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b60006123fc614293565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561248757600080fd5b505afa15801561249b573d6000803e3d6000fd5b505050506040513d60208110156124b157600080fd5b50519050818110156124c257600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8561250b85850387614481565b6040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561255e57600080fd5b505af1158015612572573d6000803e3d6000fd5b505050506040513d602081101561258857600080fd5b50516125f557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a9050612608614cfe565b60008061264a8b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061449b92505050565b6040880152608087019290925291855290925090508767ffffffffffffffff8111801561267657600080fd5b506040519080825280601f01601f1916602001820160405280156126a1576020820181803683370190505b50606084015260005b60ff811689111561271257858160ff16602081106126c457fe5b1a60f81b84606001518260ff16815181106126db57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506001016126aa565b5082604001515167ffffffffffffffff8111801561272f57600080fd5b506040519080825280601f01601f19166020018201604052801561275a576020820181803683370190505b506020840152612768614d32565b60005b8460400151518160ff161015612888576000838260ff166020811061278c57fe5b1a905082816020811061279b57fe5b60200201511561280c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6f6273657276657220696e646578207265706561746564210000000000000000604482015290519081900360640190fd5b6001838260ff166020811061281d57fe5b91151560209283029190910152849060ff841690811061283957fe5b1a60f81b86602001518360ff168151811061285057fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505060010161276b565b50612891614cb9565b503360009081526028602090815260409182902082518084019093525460ff8082168452610100909104161580159183018290526129035750602a816000015160ff16815481106128de57fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b61296e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b84516040015160ff168a116129e457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b60208a1115612a5457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b878a14612ac257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b846040015151600211158015612ade5750602085604001515111155b612b4957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b84600001516040015160020260ff1685604001515111612bca57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b5050825164ffffffffff909216602090920191909152506040516000908a908a908083838082843760405192018290039091209450612c0d9350614d3292505050565b612c15614cb9565b60005b89811015612e4657600060018587606001518481518110612c3557fe5b60209101015160f81c601b018e8e86818110612c4d57fe5b905060200201358d8d87818110612c6057fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015612cbb573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526028602090815290849020838501909452925460ff808216845261010090910416151592820183905290955092509050612d9c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f7261636c6520776173206e6f742070726573656e7421000000000000000000604482015290519081900360640190fd5b8251849060ff1660208110612dad57fe5b602002015115612e1e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff1660208110612e3357fe5b9115156020909202015250600101612c18565b5050505060005b600182604001515103811015612f1157600082604001518260010181518110612e7257fe5b602002602001015160170b83604001518381518110612e8d57fe5b602002602001015160170b1315905080612f0857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b50600101612e4d565b50604081015180516000919060028104908110612f2a57fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b13158015612f9057507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b612ffb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602c82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519985168a52898801959095523397890188905260a0890185905260c09989018a815283519a8a019a909a5282517ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d4519a94998c99989497939695909490939285019260e08601928983019202908190849084905b83811015613161578181015183820152602001613149565b50505050905001838103825285818151815260200191508051906020019080838360005b8381101561319d578181015183820152602001613185565b50505050905090810190601f1680156131ca5780820380516001836020036101000a031916815260200191505b509850505050505050505060405180910390a181516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a36132808260000151606001518260170b614731565b5080518051602b8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff928316021790915582106133a157fe5b6133af82826020015161485c565b505050505050505050565b60006133c4614293565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016916370a08231916024808301926020929190829003018186803b15801561345057600080fd5b505afa158015613464573d6000803e3d6000fd5b505050506040513d602081101561347a57600080fd5b505103905090565b600061348c614cb9565b5073ffffffffffffffffffffffffffffffffffffffff821660009081526028602090815260409182902082518084019093525460ff80821684526101009091041615159082018190526134e3576000915050610bb7565b805160049060ff16602081106134f557fe5b601091828204019190066002029054906101000a900461ffff16915050919050565b60008080808033321461358b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602b5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602c6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600660205260409020541633146136b557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b3373ffffffffffffffffffffffffffffffffffffffff8216141561373a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff0000000000000000000000000000000000000000821681179092559091169081146137df5760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461386a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461396657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b611a1b81614aad565b602b54760100000000000000000000000000000000000000000000900463ffffffff16600080808061399f614cb9565b5050505063ffffffff82166000908152602c6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b613a16614cb9565b5073ffffffffffffffffffffffffffffffffffffffff8116600090815260286020908152604080832081518083019092525460ff80821683526101009091041615159181019190915290613a6983610a7a565b905080156137df5773ffffffffffffffffffffffffffffffffffffffff80841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b158015613b1f57600080fd5b505af1158015613b33573d6000803e3d6000fd5b505050506040513d6020811015613b4957600080fd5b5051613bb657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff1660208110613bcc57fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff1660208110613c0757fe5b01556040805173ffffffffffffffffffffffffffffffffffffffff80871682528316602082015281517fe403f9d01c1c84f7381a468229c5bdb3f454c492688c83a6cc62fcfd14811610929181900390910190a150505050565b60008a8a8a8a8a8a8a8a8a8a604051602001808b73ffffffffffffffffffffffffffffffffffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810384528a8152602090810191508b908b0280828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b613db6614cd0565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152613e2d614d32565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411613e4657905050505050509050613e8d614d32565b604080516104008101918290529060089060209082845b815481526020019060010190808311613ea457505050505090506060602a805480602002602001604051908101604052809291908181526020018280548015613f2357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613ef8575b5050505050905060005b81518110156140fd5760006001848360208110613f4657fe5b602002015103905060006001868460208110613f5e57fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb868681518110613fcb57fe5b6020026020010151836040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561402757600080fd5b505af115801561403b573d6000803e3d6000fd5b505050506040513d602081101561405157600080fd5b50516140be57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018785602081106140cc57fe5b602002019061ffff16908161ffff168152505060018685602081106140ed57fe5b6020020152505050600101613f2d565b5061410b6004846020614d51565b50611e9d6008836020614de7565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b600061429d614d32565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116142b65790505050505050905060005b602081101561432657600182826020811061430f57fe5b60200201510361ffff1692909201916001016142f8565b5061432f614cd0565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608086018290527001000000000000000000000000000000009094049092166080850152602a805486518184028101840190975280875297909202633b9aca000296939492939083018282801561440c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116143e1575b5050505050905061441b614d32565b604080516104008101918290529060089060209082845b815481526020019060010190808311614432575050505050905060005b825181101561447957600182826020811061446657fe5b602002015103959095019460010161444f565b505050505090565b600081831015614492575081614495565b50805b92915050565b6144a3614e21565b506040805160808082018352602b549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff166020808401919091527501000000000000000000000000000000000000000000820460ff1693830193909352760100000000000000000000000000000000000000000000900463ffffffff1660608083019190915283519192600092839283929091908701908281101561456857600080fd5b8151602083015160408085018051915193959294830192918464010000000082111561459357600080fd5b9083019060208201858111156145a857600080fd5b82518660208202830111640100000000821117156145c557600080fd5b82525081516020918201928201910280838360005b838110156145f25781810151838201526020016145da565b5050505091909101604052505089519498508897509295509093505050605884901b907fffffffffffffffffffffffffffffffff000000000000000000000000000000008083169116146146a757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b8364ffffffffff16866020015164ffffffffff161061472757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b5091939590929450565b602d5468010000000000000000900473ffffffffffffffffffffffffffffffffffffffff16806147615750610cf2565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602c602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b602485018190529489166044850152606484018890525173ffffffffffffffffffffffffffffffffffffffff87169363beed9b5193620186a09360848084019491939192918390030190829088803b15801561482d57600080fd5b5087f19350505050801561485357506040513d602081101561484e57600080fd5b505160015b61226c57611e9d565b614864614cb9565b503360009081526028602090815260409182902082518084019093525460ff80821684526101009091041615159082015261489d614cd0565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608501527001000000000000000000000000000000009092041660808301528251610400810193849052919261496c928692909160049190826000855b82829054906101000a900461ffff1661ffff168152602001906002019060208260010104928301926001038202915080841161492a5790505050505050614b56565b61497a906004906020614d51565b5081602001516149eb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b6000614a12633b9aca003a04836020015163ffffffff16846000015163ffffffff16614bce565b90506010360260005a90506000614a318863ffffffff16858585614bf4565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff16830281614a5f57fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff1660208110614a8857fe5b015401016008886000015160ff1660208110614aa057fe5b0155505050505050505050565b60035473ffffffffffffffffffffffffffffffffffffffff9081169082168114610cf257600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b614b5e614d32565b60005b8351811015614bc6576000848281518110614b7857fe5b0160209081015160f81c9150614ba090859083908110614b9457fe5b60200201516001614c9a565b848260ff1660208110614baf57fe5b61ffff909216602092909202015250600101614b61565b509092915050565b60008383811015614be157600285850304015b614beb8184614481565b95945050505050565b600081851015614c6557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b81850383016159d801633b9aca00858202026fffffffffffffffffffffffffffffffff8110614c9057fe5b9695505050505050565b6000614cb28261ffff168461ffff160161ffff614481565b9392505050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518060a00160405280614d11614e21565b81526060602082018190526040820181905280820152600060809091015290565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215614dd75791602002820160005b83821115614da757835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302614d67565b8015614dd55782816101000a81549061ffff0219169055600201602081600101049283019260010302614da7565b505b50614de3929150614e48565b5090565b8260208101928215614e15579160200282015b82811115614e15578251825591602001919060010190614dfa565b50614de3929150614e7f565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b80821115614de35780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000168155600101614e49565b5b80821115614de35760008155600101614e8056fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


func DeployOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _validator common.Address, _minAnswer *big.Int, _maxAnswer *big.Int, _billingAdminAccessController common.Address, _decimals uint8, description string) (common.Address, *types.Transaction, *OffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OffchainAggregatorBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _validator, _minAnswer, _maxAnswer, _billingAdminAccessController, _decimals, description)
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




func (_OffchainAggregator *OffchainAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts) (*OffchainAggregatorNewTransmissionIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "NewTransmission")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorNewTransmissionIterator{contract: _OffchainAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}




func (_OffchainAggregator *OffchainAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorNewTransmission) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "NewTransmission")
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


const OffchainAggregatorBillingABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"availableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_billingAdminAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"


var OffchainAggregatorBillingBin = "0x60a06040523480156200001157600080fd5b506040516200256d3803806200256d833981810160405260e08110156200003757600080fd5b508051602082015160408301516060840151608085015160a086015160c090960151600080546001600160a01b0319163317905594959394929391929091906200008587878787876200013b565b62000090816200022d565b6001600160601b0319606083901b16608052620000ac620002a6565b620000b6620002a6565b60005b60208160ff1610156200010b576001838260ff1660208110620000d857fe5b602002019061ffff16908161ffff16815250506001828260ff1660208110620000fd57fe5b6020020152600101620000b9565b506200011b6004836020620002c5565b506200012b600882602062000362565b50505050505050505050620003d4565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b039081169082168114620002a257600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6040518061040001604052806020906020820280368337509192915050565b600283019183908215620003505791602002820160005b838211156200031e57835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302620002dc565b80156200034e5782816101000a81549061ffff02191690556002016020816001010492830192600103026200031e565b505b506200035e929150620003a1565b5090565b826020810192821562000393579160200282015b828111156200039357825182559160200191906001019062000376565b506200035e929150620003bd565b5b808211156200035e57805461ffff19168155600101620003a2565b5b808211156200035e5760008155600101620003be565b60805160601c6121626200040b600039806105675280610f1d5280610fdf5280611166528061168b52806119ec52506121626000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063b121e14711610097578063e4902f8211610066578063e4902f8214610338578063eb5dcd6c14610382578063f2fde38b146103bd578063fbffd2c1146103f0576100f5565b8063b121e1471461027f578063bd824706146102b2578063c1075329146102f7578063d09dc33914610330576100f5565b806379ba5097116100d357806379ba5097146101785780638ac28d5a146101825780638da5cb5b146101b55780639c849b30146101bd576100f5565b80630eafb25b146100fa5780631b6b6d231461013f57806373f666f814610170575b600080fd5b61012d6004803603602081101561011057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610423565b60408051918252519081900360200190f35b610147610565565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61012d610589565b61018061058f565b005b6101806004803603602081101561019857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610691565b610147610732565b610180600480360360408110156101d357600080fd5b8101906020810181356401000000008111156101ee57600080fd5b82018360208201111561020057600080fd5b8035906020019184602083028401116401000000008311171561022257600080fd5b91939092909160208101903564010000000081111561024057600080fd5b82018360208201111561025257600080fd5b8035906020019184602083028401116401000000008311171561027457600080fd5b50909250905061074e565b6101806004803603602081101561029557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610a4a565b610180600480360360a08110156102c857600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516610b77565b6101806004803603604081101561030d57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610d8f565b61012d611116565b61036b6004803603602081101561034e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166111de565b6040805161ffff9092168252519081900360200190f35b6101806004803603604081101561039857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516611273565b610180600480360360208110156103d357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611437565b6101806004803603602081101561040657600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611533565b600061042d611fac565b5073ffffffffffffffffffffffffffffffffffffffff821660009081526028602090815260409182902082518084019093525460ff8082168452610100909104161515908201819052610484576000915050610560565b61048c611fc3565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216948401949094526c0100000000000000000000000082048116606084018190527001000000000000000000000000000000009092041660808301528351919260009260019160049160ff1690811061051957fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff166020811061055757fe5b01540301925050505b919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6159d881565b60015473ffffffffffffffffffffffffffffffffffffffff16331461061557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526006602052604090205416331461072657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b61072f816115c2565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146107d457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b82811461084257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015610a4357600085858381811061085b57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff169050600084848481811061088857fe5b73ffffffffffffffffffffffffffffffffffffffff85811660009081526006602090815260409091205492029390930135831693509091169050801580806108fb57508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b61096657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614610a33578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506108459050565b5050505050565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260076020526040902054163314610adf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60035473ffffffffffffffffffffffffffffffffffffffff1680610bfc57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f61636365737320636f6e74726f6c6c6572206d75737420626520736574000000604482015290519081900360640190fd5b604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610cb957600080fd5b505afa158015610ccd573d6000803e3d6000fd5b505050506040513d6020811015610ce357600080fd5b505180610d07575060005473ffffffffffffffffffffffffffffffffffffffff1633145b610d7257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b610d7a611815565b610d878686868686611b80565b505050505050565b600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610e5457600080fd5b505afa158015610e68573d6000803e3d6000fd5b505050506040513d6020811015610e7e57600080fd5b505180610ea2575060005473ffffffffffffffffffffffffffffffffffffffff1633145b610f0d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c792062696c6c696e6741646d696e266f776e65722063616e2063616c6c604482015290519081900360640190fd5b6000610f17611cfa565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015610fa257600080fd5b505afa158015610fb6573d6000803e3d6000fd5b505050506040513d6020811015610fcc57600080fd5b5051905081811015610fdd57600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8561102685850387611ee8565b6040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561107957600080fd5b505af115801561108d573d6000803e3d6000fd5b505050506040513d60208110156110a357600080fd5b505161111057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b6000611120611cfa565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016916370a08231916024808301926020929190829003018186803b1580156111ac57600080fd5b505afa1580156111c0573d6000803e3d6000fd5b505050506040513d60208110156111d657600080fd5b505103905090565b60006111e8611fac565b5073ffffffffffffffffffffffffffffffffffffffff821660009081526028602090815260409182902082518084019093525460ff808216845261010090910416151590820181905261123f576000915050610560565b805160049060ff166020811061125157fe5b601091828204019190066002029054906101000a900461ffff16915050919050565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526006602052604090205416331461130857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b3373ffffffffffffffffffffffffffffffffffffffff8216141561138d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff0000000000000000000000000000000000000000821681179092559091169081146114325760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146114bd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff1633146115b957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b61072f81611f02565b6115ca611fac565b5073ffffffffffffffffffffffffffffffffffffffff8116600090815260286020908152604080832081518083019092525460ff8082168352610100909104161515918101919091529061161d83610423565b905080156114325773ffffffffffffffffffffffffffffffffffffffff80841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b1580156116d357600080fd5b505af11580156116e7573d6000803e3d6000fd5b505050506040513d60208110156116fd57600080fd5b505161176a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff166020811061178057fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16602081106117bb57fe5b01556040805173ffffffffffffffffffffffffffffffffffffffff80871682528316602082015281517fe403f9d01c1c84f7381a468229c5bdb3f454c492688c83a6cc62fcfd14811610929181900390910190a150505050565b61181d611fc3565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152611894611ff1565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116118ad579050505050505090506118f4611ff1565b604080516104008101918290529060089060209082845b81548152602001906001019080831161190b57505050505090506060602a80548060200260200160405190810160405280929190818152602001828054801561198a57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161195f575b5050505050905060005b8151811015611b6457600060018483602081106119ad57fe5b6020020151039050600060018684602081106119c557fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb868681518110611a3257fe5b6020026020010151836040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015611a8e57600080fd5b505af1158015611aa2573d6000803e3d6000fd5b505050506040513d6020811015611ab857600080fd5b5051611b2557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b6001878560208110611b3357fe5b602002019061ffff16908161ffff16815250506001868560208110611b5457fe5b6020020152505050600101611994565b50611b726004846020612010565b50610a4360088360206120a6565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000611d04611ff1565b6040805161040081019182905290600490602090826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411611d1d5790505050505050905060005b6020811015611d8d576001828260208110611d7657fe5b60200201510361ffff169290920191600101611d5f565b50611d96611fc3565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c010000000000000000000000008304821660608086018290527001000000000000000000000000000000009094049092166080850152602a805486518184028101840190975280875297909202633b9aca0002969394929390830182828015611e7357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611e48575b50505050509050611e82611ff1565b604080516104008101918290529060089060209082845b815481526020019060010190808311611e99575050505050905060005b8251811015611ee0576001828260208110611ecd57fe5b6020020151039590950194600101611eb6565b505050505090565b600081831015611ef9575081611efc565b50805b92915050565b60035473ffffffffffffffffffffffffffffffffffffffff9081169082168114611fa857600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b6040518061040001604052806020906020820280368337509192915050565b6002830191839082156120965791602002820160005b8382111561206657835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302612026565b80156120945782816101000a81549061ffff0219169055600201602081600101049283019260010302612066565b505b506120a29291506120e0565b5090565b82602081019282156120d4579160200282015b828111156120d45782518255916020019190600101906120b9565b506120a2929150612117565b5b808211156120a25780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001681556001016120e1565b5b808211156120a2576000815560010161211856fea2646970667358220000000000000000000000000000000000000000000000000000000000000000000064736f6c63430000000033"


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
