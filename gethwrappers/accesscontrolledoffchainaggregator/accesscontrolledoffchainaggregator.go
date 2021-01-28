// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accesscontrolledoffchainaggregator

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AccessControlledOffchainAggregatorABI is the input ABI used to generate the binding from.
const AccessControlledOffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_requesterAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"billingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requesterAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_requesterAccessController\",\"type\":\"address\"}],\"name\":\"setRequesterAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AccessControlledOffchainAggregatorBin is the compiled bytecode used for deploying new contracts.
var AccessControlledOffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b50604051620060a7380380620060a783398181016040526101a08110156200003957600080fd5b815160208301516040808501516060860151608087015160a088015160c089015160e08a01516101008b01516101208c01516101408d01516101608e01516101808f0180519b519d9f9c9e9a9d999c989b979a96999598949793969295919491939282019284640100000000821115620000b257600080fd5b908301906020820185811115620000c857600080fd5b8251640100000000811182820188101715620000e357600080fd5b82525081516020918201929091019080838360005b8381101562000112578181015183820152602001620000f8565b50505050905090810190601f168015620001405780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b03191633179055508c8c8c8c8c8c8c8c8c8c8c8c8c8c8c8c8c8c8c896200017c8787878787620002da565b6200018781620003cc565b6001600160601b0319606083901b16608052620001a362000602565b620001ad62000602565b60005b601f8160ff161015620001fd576001838260ff16601f8110620001cf57fe5b61ffff909216602092909202015260018260ff8316601f8110620001ef57fe5b6020020152600101620001b0565b506200020d600483601f62000621565b506200021d600882601f620006be565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e05250508351620002669350602e9250602085019150620006ef565b50620002728362000445565b6200027d876200051d565b8560170b60a08160170b60401b815250508460170b60c08160170b60401b81525050505050505050505050505050506001602f60006101000a81548160ff0219169083151502179055505050505050505050505050505062000788565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b0390811690821681146200044157600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b03163314620004a5576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d546001600160a01b0390811690821681146200044157602d80546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349281900390910190a15050565b6000546001600160a01b031633146200057d576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b0368010000000000000000909104811690821681146200044157602c8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b604051806103e00160405280601f906020820280368337509192915050565b600283019183908215620006ac5791602002820160005b838211156200067a57835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030262000638565b8015620006aa5782816101000a81549061ffff02191690556002016020816001010492830192600103026200067a565b505b50620006ba92915062000771565b5090565b82601f8101928215620006ac579160200282015b82811115620006ac578251825591602001919060010190620006d2565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620007275760008555620006ac565b82601f106200074257805160ff1916838001178555620006ac565b82800160010185558215620006ac5791820182811115620006ac578251825591602001919060010190620006d2565b5b80821115620006ba576000815560010162000772565b60805160601c60a05160401c60c05160401c60e05160f81c6158ba620007ed60003980611080525080611baa5280613898525080610fc7528061386b525080610fa35280612be95280612cd95280613cd252806143725280614bd552506158ba6000f3fe608060405234801561001057600080fd5b50600436106102f45760003560e01c80638ac28d5a11610191578063bd824706116100e3578063e4902f8211610097578063f2fde38b11610071578063f2fde38b14610c36578063fbffd2c114610c5c578063feaf968c14610c82576102f4565b8063e4902f8214610b61578063e5fe457714610b9e578063eb5dcd6c14610c08576102f4565b8063c9807539116100c8578063c980753914610a3d578063d09dc33914610b51578063dc7f012414610b59576102f4565b8063bd824706146109cc578063c107532914610a11576102f4565b80639c849b3011610145578063b121e1471161011f578063b121e1471461096c578063b5ab58dc14610992578063b633620c146109af576102f4565b80639c849b301461085e5780639e3ceeab14610920578063a118f24914610946576102f4565b806398e5b12a1161017657806398e5b12a146107bc578063996e8298146107e35780639a6fc8f5146107eb576102f4565b80638ac28d5a1461078e5780638da5cb5b146107b4576102f4565b8063668a0f021161024a57806379ba5097116101fe57806381ff7048116101d857806381ff70481461070f5780638205bf6a146107605780638823da6c14610768576102f4565b806379ba5097146106a75780638038e4a1146106af57806381411834146106b7576102f4565b806370da2f671161022f57806370da2f671461061a57806370efdf2d146106225780637284e4161461062a576102f4565b8063668a0f02146105485780636b14daf814610550576102f4565b806329937268116102ac57806350d25bcd1161028657806350d25bcd1461040b57806354fd4d5014610413578063585aa7de1461041b576102f4565b806329937268146103a4578063313ce567146103e55780633a5381b514610403576102f4565b80631327d3d8116102dd5780631327d3d81461033b5780631b6b6d231461036157806322adbc7814610385576102f4565b80630a756983146102f95780630eafb25b14610303575b600080fd5b610301610c8a565b005b6103296004803603602081101561031957600080fd5b50356001600160a01b0316610d48565b60408051918252519081900360200190f35b6103016004803603602081101561035157600080fd5b50356001600160a01b0316610ea8565b610369610fa1565b604080516001600160a01b039092168252519081900360200190f35b61038d610fc5565b6040805160179290920b8252519081900360200190f35b6103ac610fe9565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b6103ed61107e565b6040805160ff9092168252519081900360200190f35b6103696110a2565b6103296110bd565b61032961115e565b610301600480360360a081101561043157600080fd5b81019060208101813564010000000081111561044c57600080fd5b82018360208201111561045e57600080fd5b8035906020019184602083028401116401000000008311171561048057600080fd5b91939092909160208101903564010000000081111561049e57600080fd5b8201836020820111156104b057600080fd5b803590602001918460208302840111640100000000831117156104d257600080fd5b9193909260ff8335169267ffffffffffffffff60208201351692919060608101906040013564010000000081111561050957600080fd5b82018360208201111561051b57600080fd5b8035906020019184600183028401116401000000008311171561053d57600080fd5b509092509050611163565b610329611ae4565b6106066004803603604081101561056657600080fd5b6001600160a01b03823516919081019060408101602082013564010000000081111561059157600080fd5b8201836020820111156105a357600080fd5b803590602001918460018302840111640100000000831117156105c557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611b80945050505050565b604080519115158252519081900360200190f35b61038d611ba8565b610369611bcc565b610632611bdb565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561066c578181015183820152602001610654565b50505050905090810190601f1680156106995780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610301611c77565b610301611d45565b6106bf611e04565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156106fb5781810151838201526020016106e3565b505050509050019250505060405180910390f35b610717611e66565b6040805163ffffffff94851681529290931660208301527fffffffffffffffffffffffffffffffff00000000000000000000000000000000168183015290519081900360600190f35b610329611e87565b6103016004803603602081101561077e57600080fd5b50356001600160a01b0316611f23565b610301600480360360208110156107a457600080fd5b50356001600160a01b031661201a565b610369612091565b6107c46120a0565b6040805169ffffffffffffffffffff9092168252519081900360200190f35b6103696122df565b6108146004803603602081101561080157600080fd5b503569ffffffffffffffffffff166122ee565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6103016004803603604081101561087457600080fd5b81019060208101813564010000000081111561088f57600080fd5b8201836020820111156108a157600080fd5b803590602001918460208302840111640100000000831117156108c357600080fd5b9193909290916020810190356401000000008111156108e157600080fd5b8201836020820111156108f357600080fd5b8035906020019184602083028401116401000000008311171561091557600080fd5b5090925090506123a3565b6103016004803603602081101561093657600080fd5b50356001600160a01b03166125dc565b6103016004803603602081101561095c57600080fd5b50356001600160a01b03166126ca565b6103016004803603602081101561098257600080fd5b50356001600160a01b0316612732565b610329600480360360208110156109a857600080fd5b503561282b565b610329600480360360208110156109c557600080fd5b50356128c8565b610301600480360360a08110156109e257600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516612965565b61030160048036036040811015610a2757600080fd5b506001600160a01b038135169060200135612aad565b61030160048036036080811015610a5357600080fd5b810190602081018135640100000000811115610a6e57600080fd5b820183602082011115610a8057600080fd5b80359060200191846001830284011164010000000083111715610aa257600080fd5b919390929091602081019035640100000000811115610ac057600080fd5b820183602082011115610ad257600080fd5b80359060200191846020830284011164010000000083111715610af457600080fd5b919390929091602081019035640100000000811115610b1257600080fd5b820183602082011115610b2457600080fd5b80359060200191846020830284011164010000000083111715610b4657600080fd5b919350915035612ddc565b610329613ccd565b610606613d7e565b610b8760048036036020811015610b7757600080fd5b50356001600160a01b0316613d87565b6040805161ffff9092168252519081900360200190f35b610ba6613e40565b604080517fffffffffffffffffffffffffffffffff00000000000000000000000000000000909616865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b61030160048036036040811015610c1e57600080fd5b506001600160a01b0381358116916020013516613f2f565b61030160048036036020811015610c4c57600080fd5b50356001600160a01b031661408b565b61030160048036036020811015610c7257600080fd5b50356001600160a01b0316614153565b6108146141bb565b6000546001600160a01b03163314610ce9576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602f5460ff1615610d4657602f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b6000610d526156f1565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115610d9657fe5b6002811115610da157fe5b9052509050600081602001516002811115610db857fe5b1415610dc8576000915050610ea3565b610dd0615708565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830181905270010000000000000000000000000000000090910490921660808201528251909160009160019060049060ff16601f8110610e5c57fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff16601f8110610e9a57fe5b01540301925050505b919050565b6000546001600160a01b03163314610f07576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b036801000000000000000090910481169082168114610f9d57602c80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff16680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000610ff9615708565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b7f000000000000000000000000000000000000000000000000000000000000000081565b602c546801000000000000000090046001600160a01b031690565b6000611100336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b8092505050565b611151576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61115961426e565b905090565b600481565b868560ff8616601f8311156111bf576040805162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b60008111611214576040805162461bcd60e51b815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b8183146112525760405162461bcd60e51b815260040180806020018281038252602481526020018061588a6024913960400191505060405180910390fd5b8060030283116112a9576040805162461bcd60e51b815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b6000546001600160a01b03163314611308576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602854156114ac57602880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101916000918390811061134557fe5b6000918252602082200154602980546001600160a01b039092169350908490811061136c57fe5b6000918252602090912001546001600160a01b0316905061138c816142aa565b6001600160a01b0380831660009081526027602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000908116909155928416825290208054909116905560288054806113e857fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055602980548061144b57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611308915050565b60005b8a8110156118ba576000602760008e8e858181106114c957fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff1660028111156114ff57fe5b14611551576040805162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260016020820152602760008e8e8581811061157857fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561160357fe5b02179055506000915060069050818c8c8581811061161d57fe5b6001600160a01b036020918202939093013583168452830193909352604090910160002054169190911415905061169b576040805162461bcd60e51b815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b6000602760008c8c858181106116ad57fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff1660028111156116e357fe5b14611735576040805162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260026020820152602760008c8c8581811061175c57fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101008360028111156117e757fe5b021790555090505060288c8c838181106117fd57fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03959092029390930135939093169290921790555060298a8a8381811061185f57fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03969093029490940135949094161790915550016114af565b50602a805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602c80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000090911617938490559091048116911661198630828f8f8f8f8f8f8f8f6144da565b602a60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602a60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600083820152604051601f909101601f19169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b6000611b27336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b8092505050565b611b78576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6111596145de565b6000611b8c8383614604565b80611b9f57506001600160a01b03831632145b90505b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b602d546001600160a01b031690565b6060611c1e336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b8092505050565b611c6f576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611159614634565b6001546001600160a01b03163314611cd6576040805162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b03163314611da4576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602f5460ff16610d4657602f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60606029805480602002602001604051908101604052809291908181526020018280548015611e5c57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611e3e575b5050505050905090565b602c54602a5463ffffffff808316926401000000009004169060801b909192565b6000611eca336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b8092505050565b611f1b576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6111596146df565b6000546001600160a01b03163314611f82576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6001600160a01b03811660009081526030602052604090205460ff1615612017576001600160a01b03811660008181526030602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b6001600160a01b03818116600090815260066020526040902054163314612088576040805162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b612017816142aa565b6000546001600160a01b031681565b600080546001600160a01b031633148061217c5750602d54604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b15801561214f57600080fd5b505afa158015612163573d6000803e3d6000fd5b505050506040513d602081101561217957600080fd5b50515b6121cd576040805162461bcd60e51b815260206004820152601d60248201527f4f6e6c79206f776e6572267265717565737465722063616e2063616c6c000000604482015290519081900360640190fd5b6121d5615736565b506040805160808082018352602a549081901b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000016808352700100000000000000000000000000000000820464ffffffffff81166020808601919091527501000000000000000000000000000000000000000000840460ff9081168688015276010000000000000000000000000000000000000000000090940463ffffffff9081166060808801919091528751948552600884901c909116918401919091529216818501529251919233927f3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037929181900390910190a2806060015160010163ffffffff1691505090565b6003546001600160a01b031690565b6000806000806000612337336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b8092505050565b612388576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6123918661473a565b939a9299509097509550909350915050565b6000546001600160a01b03163314612402576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b828114612456576040805162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b838110156125d557600085858381811061246f57fe5b905060200201356001600160a01b03169050600084848481811061248f57fe5b6001600160a01b0385811660009081526006602090815260409091205492029390930135831693509091169050801580806124db5750826001600160a01b0316826001600160a01b0316145b61252c576040805162461bcd60e51b815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b03848116600090815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016858316908117909155908316146125c557826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506124599050565b5050505050565b6000546001600160a01b0316331461263b576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d546001600160a01b039081169082168114610f9d57602d80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349281900390910190a15050565b6000546001600160a01b03163314612729576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6120178161488e565b6001600160a01b038181166000908152600760205260409020541633146127a0576040805162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b6001600160a01b0381811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600061286e336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b8092505050565b6128bf576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611ba282614927565b600061290b336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b8092505050565b61295c576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611ba28261495d565b6003546000546001600160a01b039182169116331480612a3f5750604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b03861694636b14daf8946000939190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b158015612a1257600080fd5b505afa158015612a26573d6000803e3d6000fd5b505050506040513d6020811015612a3c57600080fd5b50515b612a90576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b612a986149b2565b612aa58686868686614d8e565b505050505050565b6000546001600160a01b0316331480612b885750600354604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b158015612b5b57600080fd5b505afa158015612b6f573d6000803e3d6000fd5b505050506040513d6020811015612b8557600080fd5b50515b612bd9576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6000612be3614f08565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b158015612c5457600080fd5b505afa158015612c68573d6000803e3d6000fd5b505050506040513d6020811015612c7e57600080fd5b5051905081811015612cd7576040805162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb85612d13858503876150e9565b6040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b158015612d5957600080fd5b505af1158015612d6d573d6000803e3d6000fd5b505050506040513d6020811015612d8357600080fd5b5051612dd6576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a9050612def888888888888615100565b3614612e42576040805162461bcd60e51b815260206004820152601960248201527f7472616e736d6974206d65737361676520746f6f206c6f6e6700000000000000604482015290519081900360640190fd5b612e4a61575d565b6040805160808082018352602a549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff1660208301527501000000000000000000000000000000000000000000810460ff169282019290925276010000000000000000000000000000000000000000000090910463ffffffff166060808301919091529082526000908a908a90811015612f0357600080fd5b813591602081013591810190606081016040820135640100000000811115612f2a57600080fd5b820183602082011115612f3c57600080fd5b80359060200191846020830284011164010000000083111715612f5e57600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050505060408801525050506080840182905283515190925060589190911b907fffffffffffffffffffffffffffffffff00000000000000000000000000000000808316911614613025576040805162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b608083015183516020015164ffffffffff80831691161061308d576040805162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b83516040015160ff1689116130e9576040805162461bcd60e51b815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b601f89111561313f576040805162461bcd60e51b815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b868914613193576040805162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b601f84604001515111156131ee576040805162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b83600001516040015160020260ff1684604001515111613255576040805162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b8867ffffffffffffffff8111801561326c57600080fd5b506040519080825280601f01601f191660200182016040528015613297576020820181803683370190505b50606085015260005b60ff81168a111561330857868160ff16602081106132ba57fe5b1a60f81b85606001518260ff16815181106132d157fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506001016132a0565b5083604001515167ffffffffffffffff8111801561332557600080fd5b506040519080825280601f01601f191660200182016040528015613350576020820181803683370190505b50602085015261335e615791565b60005b8560400151518160ff161015613464576000858260ff166020811061338257fe5b1a90508281601f811061339157fe5b6020020151156133e8576040805162461bcd60e51b815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b6001838260ff16601f81106133f957fe5b91151560209283029190910152869060ff841690811061341557fe5b1a60f81b87602001518360ff168151811061342c57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101613361565b5061346d6156f1565b336000908152602760209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156134a857fe5b60028111156134b357fe5b90525090506002816020015160028111156134ca57fe5b1480156134fe57506029816000015160ff16815481106134e657fe5b6000918252602090912001546001600160a01b031633145b61354f576040805162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b5050835164ffffffffff90911660209091015250506040516000908a908a908083838082843760405192018290039091209450613590935061579192505050565b6135986156f1565b60005b89811015613791576000600185876060015184815181106135b857fe5b60209101015160f81c601b018e8e868181106135d057fe5b905060200201358d8d878181106135e357fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801561363e573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526027602090815290849020838501909452835460ff8082168552929650929450840191610100900416600281111561368d57fe5b600281111561369857fe5b90525092506001836020015160028111156136af57fe5b14613701576040805162461bcd60e51b815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff16601f811061371257fe5b602002015115613769576040805162461bcd60e51b815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16601f811061377e57fe5b911515602090920201525060010161359b565b5050505060005b600182604001515103811015613842576000826040015182600101815181106137bd57fe5b602002602001015160170b836040015183815181106137d857fe5b602002602001015160170b1315905080613839576040805162461bcd60e51b815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b50600101613798565b5060408101518051600091906002810490811061385b57fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b131580156138c157507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b613912576040805162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602b82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519586523386890181905291860181905260a0988601898152845199870199909952835194909916997ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451998c999298949793969095909492939185019260c086019289820192909102908190849084905b83811015613a75578181015183820152602001613a5d565b50505050905001838103825285818151815260200191508051906020019080838360005b83811015613ab1578181015183820152602001613a99565b50505050905090810190601f168015613ade5780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390a281516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a3613b938260000151606001518260170b615118565b5080518051602a8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff92831602179091558210613cb457fe5b613cc2828260200151615229565b505050505050505050565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b158015613d3d57600080fd5b505afa158015613d51573d6000803e3d6000fd5b505050506040513d6020811015613d6757600080fd5b505190506000613d75614f08565b90910391505090565b602f5460ff1681565b6000613d916156f1565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115613dd557fe5b6002811115613de057fe5b9052509050600081602001516002811115613df757fe5b1415613e07576000915050610ea3565b60016004826000015160ff16601f8110613e1d57fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b600080808080333214613e9a576040805162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602a5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602b6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6001600160a01b03828116600090815260066020526040902054163314613f9d576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b336001600160a01b0382161415613ffb576040805162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b6001600160a01b03808316600090815260076020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114614086576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b6000546001600160a01b031633146140ea576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000546001600160a01b031633146141b2576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b61201781615484565b6000806000806000614204336000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611b8092505050565b614255576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61425d615513565b945094509450945094509091929394565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b6020526040902054601790810b900b90565b6142b26156f1565b6001600160a01b0382166000908152602760209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156142f657fe5b600281111561430157fe5b9052509050600061431183610d48565b90508015614086576001600160a01b0380841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b1580156143ba57600080fd5b505af11580156143ce573d6000803e3d6000fd5b505050506040513d60208110156143e457600080fd5b5051614437576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff16601f811061444d57fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16601f811061448857fe5b0155604080516001600160a01b0380871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b60008a8a8a8a8a8a8a8a8a8a604051602001808b6001600160a01b031681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff1690565b6001600160a01b03821660009081526030602052604081205460ff1680611b9f575050602f5460ff161592915050565b602e8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff610100600188161502019095169490940493840181900481028201810190925282815260609390929091830182828015611e5c5780601f106146b357610100808354040283529160200191611e5c565b820191906000526020600020905b8154815290600101906020018083116146c157509395945050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e740000000000000000000000000000000000815250906148155760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b838110156147da5781810151838201526020016147c2565b50505050905090810190601f1680156148075780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5061481e6156f1565b5050505063ffffffff83166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b6001600160a01b03811660009081526030602052604090205460ff16612017576001600160a01b03811660008181526030602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a150565b600063ffffffff82111561493d57506000610ea3565b5063ffffffff166000908152602b6020526040902054601790810b900b90565b600063ffffffff82111561497357506000610ea3565b5063ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6149ba615708565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c0100000000000000000000000081048316606083015270010000000000000000000000000000000090049091166080820152614a31615791565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614a4a57905050505050509050614a91615791565b604080516103e081019182905290600890601f9082845b815481526020019060010190808311614aa8575050505050905060606029805480602002602001604051908101604052809291908181526020018280548015614b1a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311614afc575b5050505050905060005b8151811015614d7257600060018483601f8110614b3d57fe5b6020020151039050600060018684601f8110614b5557fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115614d6757600060066000878781518110614b9557fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b031690507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb82846040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b158015614c4a57600080fd5b505af1158015614c5e573d6000803e3d6000fd5b505050506040513d6020811015614c7457600080fd5b5051614cc7576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018886601f8110614cd557fe5b61ffff909216602092909202015260018786601f8110614cf157fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290879087908110614d2657fe5b6020026020010151828460405180846001600160a01b03168152602001836001600160a01b03168152602001828152602001935050505060405180910390a1505b505050600101614b24565b50614d80600484601f6157b0565b506125d5600883601f615846565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000614f12615791565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614f2b5790505050505050905060005b601f811015614f9b5760018282601f8110614f8457fe5b60200201510361ffff169290920191600101614f6d565b50614fa4615708565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606080860182905270010000000000000000000000000000000090940490921660808501526029805486518184028101840190975280875297909202633b9aca000296939492939083018282801561507457602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311615056575b50505050509050615083615791565b604080516103e081019182905290600890601f9082845b81548152602001906001019080831161509a575050505050905060005b82518110156150e15760018282601f81106150ce57fe5b60200201510395909501946001016150b7565b505050505090565b6000818310156150fa575081611ba2565b50919050565b602083810286019082020160e4019695505050505050565b602c546801000000000000000090046001600160a01b03168061513b5750610f9d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602b602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b60248501819052948916604485015260648401889052516001600160a01b0387169363beed9b5193620186a09360848084019491939192918390030190829088803b1580156151fa57600080fd5b5087f19350505050801561522057506040513d602081101561521b57600080fd5b505160015b612aa5576125d5565b6152316156f1565b336000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561526c57fe5b600281111561527757fe5b9052509050615284615708565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000082048116606084015270010000000000000000000000000000000090910416608082015281516103e0810192839052909161534f918591600490601f90826000855b82829054906101000a900461ffff1661ffff168152602001906002019060208260010104928301926001038202915080841161530d57905050505050506155b2565b61535d90600490601f6157b0565b5060028260200151600281111561537057fe5b146153c2576040805162461bcd60e51b815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b60006153e9633b9aca003a04836020015163ffffffff16846000015163ffffffff16615627565b90506010360260005a905060006154088863ffffffff1685858561564d565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff1683028161543657fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff16601f811061545f57fe5b015401016008886000015160ff16601f811061547757fe5b0155505050505050505050565b6003546001600160a01b039081169082168114610f9d57600380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b602a54760100000000000000000000000000000000000000000000900463ffffffff1660008080806155436156f1565b5050505063ffffffff82166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b6155ba615791565b60005b835181101561561f5760008482815181106155d457fe5b016020015160f81c90506155f98482601f81106155ed57fe5b602002015160016156d9565b848260ff16601f811061560857fe5b61ffff9092166020929092020152506001016155bd565b509092915050565b6000838381101561563a57600285850304015b61564481846150e9565b95945050505050565b6000818510156156a4576040805162461bcd60e51b815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b818503830161179301633b9aca00858202026fffffffffffffffffffffffffffffffff81106156cf57fe5b9695505050505050565b6000611b9f8261ffff168461ffff160161ffff6150e9565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b60408051608081018252600080825260208201819052918101829052606081019190915290565b6040518060a00160405280615770615736565b81526060602082018190526040820181905280820152600060809091015290565b604051806103e00160405280601f906020820280368337509192915050565b6002830191839082156158365791602002820160005b8382111561580657835183826101000a81548161ffff021916908361ffff16021790555092602001926002016020816001010492830192600103026157c6565b80156158345782816101000a81549061ffff0219169055600201602081600101049283019260010302615806565b505b50615842929150615874565b5090565b82601f8101928215615836579160200282015b82811115615836578251825591602001919060010190615859565b5b80821115615842576000815560010161587556fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea164736f6c6343000705000a"

// DeployAccessControlledOffchainAggregator deploys a new Ethereum contract, binding an instance of AccessControlledOffchainAggregator to it.
func DeployAccessControlledOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _validator common.Address, _minAnswer *big.Int, _maxAnswer *big.Int, _billingAccessController common.Address, _requesterAccessController common.Address, _decimals uint8, description string) (common.Address, *types.Transaction, *AccessControlledOffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlledOffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccessControlledOffchainAggregatorBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _validator, _minAnswer, _maxAnswer, _billingAccessController, _requesterAccessController, _decimals, description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessControlledOffchainAggregator{AccessControlledOffchainAggregatorCaller: AccessControlledOffchainAggregatorCaller{contract: contract}, AccessControlledOffchainAggregatorTransactor: AccessControlledOffchainAggregatorTransactor{contract: contract}, AccessControlledOffchainAggregatorFilterer: AccessControlledOffchainAggregatorFilterer{contract: contract}}, nil
}

// AccessControlledOffchainAggregator is an auto generated Go binding around an Ethereum contract.
type AccessControlledOffchainAggregator struct {
	AccessControlledOffchainAggregatorCaller     // Read-only binding to the contract
	AccessControlledOffchainAggregatorTransactor // Write-only binding to the contract
	AccessControlledOffchainAggregatorFilterer   // Log filterer for contract events
}

// AccessControlledOffchainAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccessControlledOffchainAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlledOffchainAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccessControlledOffchainAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlledOffchainAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccessControlledOffchainAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlledOffchainAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccessControlledOffchainAggregatorSession struct {
	Contract     *AccessControlledOffchainAggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                       // Call options to use throughout this session
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// AccessControlledOffchainAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccessControlledOffchainAggregatorCallerSession struct {
	Contract *AccessControlledOffchainAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                             // Call options to use throughout this session
}

// AccessControlledOffchainAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccessControlledOffchainAggregatorTransactorSession struct {
	Contract     *AccessControlledOffchainAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// AccessControlledOffchainAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccessControlledOffchainAggregatorRaw struct {
	Contract *AccessControlledOffchainAggregator // Generic contract binding to access the raw methods on
}

// AccessControlledOffchainAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccessControlledOffchainAggregatorCallerRaw struct {
	Contract *AccessControlledOffchainAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// AccessControlledOffchainAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccessControlledOffchainAggregatorTransactorRaw struct {
	Contract *AccessControlledOffchainAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccessControlledOffchainAggregator creates a new instance of AccessControlledOffchainAggregator, bound to a specific deployed contract.
func NewAccessControlledOffchainAggregator(address common.Address, backend bind.ContractBackend) (*AccessControlledOffchainAggregator, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregator{AccessControlledOffchainAggregatorCaller: AccessControlledOffchainAggregatorCaller{contract: contract}, AccessControlledOffchainAggregatorTransactor: AccessControlledOffchainAggregatorTransactor{contract: contract}, AccessControlledOffchainAggregatorFilterer: AccessControlledOffchainAggregatorFilterer{contract: contract}}, nil
}

// NewAccessControlledOffchainAggregatorCaller creates a new read-only instance of AccessControlledOffchainAggregator, bound to a specific deployed contract.
func NewAccessControlledOffchainAggregatorCaller(address common.Address, caller bind.ContractCaller) (*AccessControlledOffchainAggregatorCaller, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCaller{contract: contract}, nil
}

// NewAccessControlledOffchainAggregatorTransactor creates a new write-only instance of AccessControlledOffchainAggregator, bound to a specific deployed contract.
func NewAccessControlledOffchainAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControlledOffchainAggregatorTransactor, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorTransactor{contract: contract}, nil
}

// NewAccessControlledOffchainAggregatorFilterer creates a new log filterer instance of AccessControlledOffchainAggregator, bound to a specific deployed contract.
func NewAccessControlledOffchainAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControlledOffchainAggregatorFilterer, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorFilterer{contract: contract}, nil
}

// bindAccessControlledOffchainAggregator binds a generic wrapper to an already deployed contract.
func bindAccessControlledOffchainAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlledOffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlledOffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.contract.Transact(opts, method, params...)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LINK() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.LINK(&_AccessControlledOffchainAggregator.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.LINK(&_AccessControlledOffchainAggregator.CallOpts)
}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) BillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "billingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) BillingAccessController() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.BillingAccessController(&_AccessControlledOffchainAggregator.CallOpts)
}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) BillingAccessController() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.BillingAccessController(&_AccessControlledOffchainAggregator.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) CheckEnabled() (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.CheckEnabled(&_AccessControlledOffchainAggregator.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.CheckEnabled(&_AccessControlledOffchainAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Decimals() (uint8, error) {
	return _AccessControlledOffchainAggregator.Contract.Decimals(&_AccessControlledOffchainAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _AccessControlledOffchainAggregator.Contract.Decimals(&_AccessControlledOffchainAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Description() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.Description(&_AccessControlledOffchainAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Description() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.Description(&_AccessControlledOffchainAggregator.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetAnswer(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetAnswer(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
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

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _AccessControlledOffchainAggregator.Contract.GetBilling(&_AccessControlledOffchainAggregator.CallOpts)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _AccessControlledOffchainAggregator.Contract.GetBilling(&_AccessControlledOffchainAggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
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

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledOffchainAggregator.Contract.GetRoundData(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledOffchainAggregator.Contract.GetRoundData(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetTimestamp(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetTimestamp(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.HasAccess(&_AccessControlledOffchainAggregator.CallOpts, _user, _calldata)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.HasAccess(&_AccessControlledOffchainAggregator.CallOpts, _user, _calldata)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes16 configDigest)
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

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes16 configDigest)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestConfigDetails(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes16 configDigest)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestConfigDetails(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestRound() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRound(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRound(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
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

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRoundData(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRoundData(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTimestamp(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTimestamp(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes16 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp)
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

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes16 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTransmissionDetails(&_AccessControlledOffchainAggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes16 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTransmissionDetails(&_AccessControlledOffchainAggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LinkAvailableForPayment(&_AccessControlledOffchainAggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LinkAvailableForPayment(&_AccessControlledOffchainAggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MaxAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MaxAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MinAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MinAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _AccessControlledOffchainAggregator.Contract.OracleObservationCount(&_AccessControlledOffchainAggregator.CallOpts, _signerOrTransmitter)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _AccessControlledOffchainAggregator.Contract.OracleObservationCount(&_AccessControlledOffchainAggregator.CallOpts, _signerOrTransmitter)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.OwedPayment(&_AccessControlledOffchainAggregator.CallOpts, _transmitter)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.OwedPayment(&_AccessControlledOffchainAggregator.CallOpts, _transmitter)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Owner() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Owner(&_AccessControlledOffchainAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Owner(&_AccessControlledOffchainAggregator.CallOpts)
}

// RequesterAccessController is a free data retrieval call binding the contract method 0x70efdf2d.
//
// Solidity: function requesterAccessController() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) RequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "requesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RequesterAccessController is a free data retrieval call binding the contract method 0x70efdf2d.
//
// Solidity: function requesterAccessController() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) RequesterAccessController() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.RequesterAccessController(&_AccessControlledOffchainAggregator.CallOpts)
}

// RequesterAccessController is a free data retrieval call binding the contract method 0x70efdf2d.
//
// Solidity: function requesterAccessController() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) RequesterAccessController() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.RequesterAccessController(&_AccessControlledOffchainAggregator.CallOpts)
}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmitters(&_AccessControlledOffchainAggregator.CallOpts)
}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmitters(&_AccessControlledOffchainAggregator.CallOpts)
}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "validator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Validator() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Validator(&_AccessControlledOffchainAggregator.CallOpts)
}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Validator(&_AccessControlledOffchainAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Version() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.Version(&_AccessControlledOffchainAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Version() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.Version(&_AccessControlledOffchainAggregator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptOwnership(&_AccessControlledOffchainAggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptOwnership(&_AccessControlledOffchainAggregator.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "acceptPayeeship", _transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "addAccess", _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AddAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AddAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "disableAccessCheck")
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.DisableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.DisableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "enableAccessCheck")
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.EnableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.EnableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "removeAccess", _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RemoveAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RemoveAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "requestNewRound")
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RequestNewRound(&_AccessControlledOffchainAggregator.TransactOpts)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RequestNewRound(&_AccessControlledOffchainAggregator.TransactOpts)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBilling(&_AccessControlledOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBilling(&_AccessControlledOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBillingAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBillingAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _billingAccessController)
}

// SetConfig is a paid mutator transaction binding the contract method 0x585aa7de.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _threshold, uint64 _encodedConfigVersion, bytes _encoded) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}

// SetConfig is a paid mutator transaction binding the contract method 0x585aa7de.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _threshold, uint64 _encodedConfigVersion, bytes _encoded) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetConfig(&_AccessControlledOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}

// SetConfig is a paid mutator transaction binding the contract method 0x585aa7de.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _threshold, uint64 _encodedConfigVersion, bytes _encoded) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetConfig(&_AccessControlledOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setPayees", _transmitters, _payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetPayees(&_AccessControlledOffchainAggregator.TransactOpts, _transmitters, _payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetPayees(&_AccessControlledOffchainAggregator.TransactOpts, _transmitters, _payees)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address _requesterAccessController) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, _requesterAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setRequesterAccessController", _requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address _requesterAccessController) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetRequesterAccessController(_requesterAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetRequesterAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address _requesterAccessController) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetRequesterAccessController(_requesterAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetRequesterAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _requesterAccessController)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetValidator(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setValidator", _newValidator)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetValidator(&_AccessControlledOffchainAggregator.TransactOpts, _newValidator)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetValidator(&_AccessControlledOffchainAggregator.TransactOpts, _newValidator)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferOwnership(&_AccessControlledOffchainAggregator.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferOwnership(&_AccessControlledOffchainAggregator.TransactOpts, _to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter, _proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter, _proposed)
}

// Transmit is a paid mutator transaction binding the contract method 0xc9807539.
//
// Solidity: function transmit(bytes _report, bytes32[] _rs, bytes32[] _ss, bytes32 _rawVs) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) Transmit(opts *bind.TransactOpts, _report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transmit", _report, _rs, _ss, _rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xc9807539.
//
// Solidity: function transmit(bytes _report, bytes32[] _rs, bytes32[] _ss, bytes32 _rawVs) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmit(&_AccessControlledOffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xc9807539.
//
// Solidity: function transmit(bytes _report, bytes32[] _rs, bytes32[] _ss, bytes32 _rawVs) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmit(&_AccessControlledOffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawFunds(&_AccessControlledOffchainAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawFunds(&_AccessControlledOffchainAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "withdrawPayment", _transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawPayment(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawPayment(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}

// AccessControlledOffchainAggregatorAddedAccessIterator is returned from FilterAddedAccess and is used to iterate over the raw logs and unpacked data for AddedAccess events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorAddedAccessIterator struct {
	Event *AccessControlledOffchainAggregatorAddedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorAddedAccess represents a AddedAccess event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorAddedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedAccess is a free log retrieval operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorAddedAccessIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorAddedAccessIterator{contract: _AccessControlledOffchainAggregator.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

// WatchAddedAccess is a free log subscription operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
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
				// New log arrived, parse the event and forward to the user
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

// ParseAddedAccess is a log parse operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseAddedAccess(log types.Log) (*AccessControlledOffchainAggregatorAddedAccess, error) {
	event := new(AccessControlledOffchainAggregatorAddedAccess)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorAnswerUpdatedIterator struct {
	Event *AccessControlledOffchainAggregatorAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorAnswerUpdated represents a AnswerUpdated event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
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

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
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
				// New log arrived, parse the event and forward to the user
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*AccessControlledOffchainAggregatorAnswerUpdated, error) {
	event := new(AccessControlledOffchainAggregatorAnswerUpdated)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorBillingAccessControllerSetIterator is returned from FilterBillingAccessControllerSet and is used to iterate over the raw logs and unpacked data for BillingAccessControllerSet events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorBillingAccessControllerSetIterator struct {
	Event *AccessControlledOffchainAggregatorBillingAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorBillingAccessControllerSet represents a BillingAccessControllerSet event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBillingAccessControllerSet is a free log retrieval operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorBillingAccessControllerSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchBillingAccessControllerSet is a free log subscription operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
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
				// New log arrived, parse the event and forward to the user
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

// ParseBillingAccessControllerSet is a log parse operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*AccessControlledOffchainAggregatorBillingAccessControllerSet, error) {
	event := new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorBillingSetIterator is returned from FilterBillingSet and is used to iterate over the raw logs and unpacked data for BillingSet events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorBillingSetIterator struct {
	Event *AccessControlledOffchainAggregatorBillingSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorBillingSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorBillingSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorBillingSet represents a BillingSet event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterBillingSet is a free log retrieval operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorBillingSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorBillingSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

// WatchBillingSet is a free log subscription operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
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
				// New log arrived, parse the event and forward to the user
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

// ParseBillingSet is a log parse operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseBillingSet(log types.Log) (*AccessControlledOffchainAggregatorBillingSet, error) {
	event := new(AccessControlledOffchainAggregatorBillingSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorCheckAccessDisabledIterator is returned from FilterCheckAccessDisabled and is used to iterate over the raw logs and unpacked data for CheckAccessDisabled events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorCheckAccessDisabledIterator struct {
	Event *AccessControlledOffchainAggregatorCheckAccessDisabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorCheckAccessDisabled represents a CheckAccessDisabled event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorCheckAccessDisabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessDisabled is a free log retrieval operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorCheckAccessDisabledIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCheckAccessDisabledIterator{contract: _AccessControlledOffchainAggregator.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessDisabled is a free log subscription operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
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
				// New log arrived, parse the event and forward to the user
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

// ParseCheckAccessDisabled is a log parse operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseCheckAccessDisabled(log types.Log) (*AccessControlledOffchainAggregatorCheckAccessDisabled, error) {
	event := new(AccessControlledOffchainAggregatorCheckAccessDisabled)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorCheckAccessEnabledIterator is returned from FilterCheckAccessEnabled and is used to iterate over the raw logs and unpacked data for CheckAccessEnabled events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorCheckAccessEnabledIterator struct {
	Event *AccessControlledOffchainAggregatorCheckAccessEnabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorCheckAccessEnabled represents a CheckAccessEnabled event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorCheckAccessEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessEnabled is a free log retrieval operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorCheckAccessEnabledIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCheckAccessEnabledIterator{contract: _AccessControlledOffchainAggregator.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessEnabled is a free log subscription operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
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
				// New log arrived, parse the event and forward to the user
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

// ParseCheckAccessEnabled is a log parse operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseCheckAccessEnabled(log types.Log) (*AccessControlledOffchainAggregatorCheckAccessEnabled, error) {
	event := new(AccessControlledOffchainAggregatorCheckAccessEnabled)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorConfigSetIterator struct {
	Event *AccessControlledOffchainAggregatorConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorConfigSet represents a ConfigSet event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb9.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, uint64 configCount, address[] signers, address[] transmitters, uint8 threshold, uint64 encodedConfigVersion, bytes encoded)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorConfigSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorConfigSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb9.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, uint64 configCount, address[] signers, address[] transmitters, uint8 threshold, uint64 encodedConfigVersion, bytes encoded)
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
				// New log arrived, parse the event and forward to the user
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

// ParseConfigSet is a log parse operation binding the contract event 0x25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb9.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, uint64 configCount, address[] signers, address[] transmitters, uint8 threshold, uint64 encodedConfigVersion, bytes encoded)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseConfigSet(log types.Log) (*AccessControlledOffchainAggregatorConfigSet, error) {
	event := new(AccessControlledOffchainAggregatorConfigSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorNewRoundIterator struct {
	Event *AccessControlledOffchainAggregatorNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorNewRound represents a NewRound event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
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

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
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
				// New log arrived, parse the event and forward to the user
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseNewRound(log types.Log) (*AccessControlledOffchainAggregatorNewRound, error) {
	event := new(AccessControlledOffchainAggregatorNewRound)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorNewTransmissionIterator is returned from FilterNewTransmission and is used to iterate over the raw logs and unpacked data for NewTransmission events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorNewTransmissionIterator struct {
	Event *AccessControlledOffchainAggregatorNewTransmission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorNewTransmission represents a NewTransmission event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	RawReportContext  [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterNewTransmission is a free log retrieval operation binding the contract event 0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, int192[] observations, bytes observers, bytes32 rawReportContext)
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

// WatchNewTransmission is a free log subscription operation binding the contract event 0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, int192[] observations, bytes observers, bytes32 rawReportContext)
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
				// New log arrived, parse the event and forward to the user
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

// ParseNewTransmission is a log parse operation binding the contract event 0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, int192[] observations, bytes observers, bytes32 rawReportContext)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseNewTransmission(log types.Log) (*AccessControlledOffchainAggregatorNewTransmission, error) {
	event := new(AccessControlledOffchainAggregatorNewTransmission)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorOraclePaidIterator is returned from FilterOraclePaid and is used to iterate over the raw logs and unpacked data for OraclePaid events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorOraclePaidIterator struct {
	Event *AccessControlledOffchainAggregatorOraclePaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorOraclePaid represents a OraclePaid event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOraclePaid is a free log retrieval operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorOraclePaidIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorOraclePaidIterator{contract: _AccessControlledOffchainAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

// WatchOraclePaid is a free log subscription operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOraclePaid is a log parse operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOraclePaid(log types.Log) (*AccessControlledOffchainAggregatorOraclePaid, error) {
	event := new(AccessControlledOffchainAggregatorOraclePaid)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator struct {
	Event *AccessControlledOffchainAggregatorOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*AccessControlledOffchainAggregatorOwnershipTransferRequested, error) {
	event := new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorOwnershipTransferredIterator struct {
	Event *AccessControlledOffchainAggregatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*AccessControlledOffchainAggregatorOwnershipTransferred, error) {
	event := new(AccessControlledOffchainAggregatorOwnershipTransferred)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator struct {
	Event *AccessControlledOffchainAggregatorPayeeshipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorPayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
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

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
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
				// New log arrived, parse the event and forward to the user
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

// ParsePayeeshipTransferRequested is a log parse operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*AccessControlledOffchainAggregatorPayeeshipTransferRequested, error) {
	event := new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorPayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorPayeeshipTransferredIterator struct {
	Event *AccessControlledOffchainAggregatorPayeeshipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorPayeeshipTransferred represents a PayeeshipTransferred event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
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

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
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
				// New log arrived, parse the event and forward to the user
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

// ParsePayeeshipTransferred is a log parse operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*AccessControlledOffchainAggregatorPayeeshipTransferred, error) {
	event := new(AccessControlledOffchainAggregatorPayeeshipTransferred)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorRemovedAccessIterator is returned from FilterRemovedAccess and is used to iterate over the raw logs and unpacked data for RemovedAccess events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorRemovedAccessIterator struct {
	Event *AccessControlledOffchainAggregatorRemovedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorRemovedAccess represents a RemovedAccess event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorRemovedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedAccess is a free log retrieval operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorRemovedAccessIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorRemovedAccessIterator{contract: _AccessControlledOffchainAggregator.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

// WatchRemovedAccess is a free log subscription operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
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
				// New log arrived, parse the event and forward to the user
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

// ParseRemovedAccess is a log parse operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseRemovedAccess(log types.Log) (*AccessControlledOffchainAggregatorRemovedAccess, error) {
	event := new(AccessControlledOffchainAggregatorRemovedAccess)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator is returned from FilterRequesterAccessControllerSet and is used to iterate over the raw logs and unpacked data for RequesterAccessControllerSet events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator struct {
	Event *AccessControlledOffchainAggregatorRequesterAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorRequesterAccessControllerSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorRequesterAccessControllerSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorRequesterAccessControllerSet represents a RequesterAccessControllerSet event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRequesterAccessControllerSet is a free log retrieval operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchRequesterAccessControllerSet is a free log subscription operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledOffchainAggregatorRequesterAccessControllerSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
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

// ParseRequesterAccessControllerSet is a log parse operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*AccessControlledOffchainAggregatorRequesterAccessControllerSet, error) {
	event := new(AccessControlledOffchainAggregatorRequesterAccessControllerSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorRoundRequestedIterator is returned from FilterRoundRequested and is used to iterate over the raw logs and unpacked data for RoundRequested events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorRoundRequestedIterator struct {
	Event *AccessControlledOffchainAggregatorRoundRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorRoundRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorRoundRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlledOffchainAggregatorRoundRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorRoundRequested represents a RoundRequested event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [16]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRoundRequested is a free log retrieval operation binding the contract event 0x3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037.
//
// Solidity: event RoundRequested(address indexed requester, bytes16 configDigest, uint32 epoch, uint8 round)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*AccessControlledOffchainAggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorRoundRequestedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

// WatchRoundRequested is a free log subscription operation binding the contract event 0x3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037.
//
// Solidity: event RoundRequested(address indexed requester, bytes16 configDigest, uint32 epoch, uint8 round)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlledOffchainAggregatorRoundRequested)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
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

// ParseRoundRequested is a log parse operation binding the contract event 0x3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037.
//
// Solidity: event RoundRequested(address indexed requester, bytes16 configDigest, uint32 epoch, uint8 round)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseRoundRequested(log types.Log) (*AccessControlledOffchainAggregatorRoundRequested, error) {
	event := new(AccessControlledOffchainAggregatorRoundRequested)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControlledOffchainAggregatorValidatorUpdatedIterator is returned from FilterValidatorUpdated and is used to iterate over the raw logs and unpacked data for ValidatorUpdated events raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorValidatorUpdatedIterator struct {
	Event *AccessControlledOffchainAggregatorValidatorUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlledOffchainAggregatorValidatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlledOffchainAggregatorValidatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlledOffchainAggregatorValidatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlledOffchainAggregatorValidatorUpdated represents a ValidatorUpdated event raised by the AccessControlledOffchainAggregator contract.
type AccessControlledOffchainAggregatorValidatorUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterValidatorUpdated is a free log retrieval operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
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

// WatchValidatorUpdated is a free log subscription operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
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
				// New log arrived, parse the event and forward to the user
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

// ParseValidatorUpdated is a log parse operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseValidatorUpdated(log types.Log) (*AccessControlledOffchainAggregatorValidatorUpdated, error) {
	event := new(AccessControlledOffchainAggregatorValidatorUpdated)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccessControllerInterfaceABI is the input ABI used to generate the binding from.
const AccessControllerInterfaceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AccessControllerInterface is an auto generated Go binding around an Ethereum contract.
type AccessControllerInterface struct {
	AccessControllerInterfaceCaller     // Read-only binding to the contract
	AccessControllerInterfaceTransactor // Write-only binding to the contract
	AccessControllerInterfaceFilterer   // Log filterer for contract events
}

// AccessControllerInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccessControllerInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControllerInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccessControllerInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControllerInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccessControllerInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControllerInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccessControllerInterfaceSession struct {
	Contract     *AccessControllerInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AccessControllerInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccessControllerInterfaceCallerSession struct {
	Contract *AccessControllerInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// AccessControllerInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccessControllerInterfaceTransactorSession struct {
	Contract     *AccessControllerInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// AccessControllerInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccessControllerInterfaceRaw struct {
	Contract *AccessControllerInterface // Generic contract binding to access the raw methods on
}

// AccessControllerInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccessControllerInterfaceCallerRaw struct {
	Contract *AccessControllerInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AccessControllerInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccessControllerInterfaceTransactorRaw struct {
	Contract *AccessControllerInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccessControllerInterface creates a new instance of AccessControllerInterface, bound to a specific deployed contract.
func NewAccessControllerInterface(address common.Address, backend bind.ContractBackend) (*AccessControllerInterface, error) {
	contract, err := bindAccessControllerInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterface{AccessControllerInterfaceCaller: AccessControllerInterfaceCaller{contract: contract}, AccessControllerInterfaceTransactor: AccessControllerInterfaceTransactor{contract: contract}, AccessControllerInterfaceFilterer: AccessControllerInterfaceFilterer{contract: contract}}, nil
}

// NewAccessControllerInterfaceCaller creates a new read-only instance of AccessControllerInterface, bound to a specific deployed contract.
func NewAccessControllerInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AccessControllerInterfaceCaller, error) {
	contract, err := bindAccessControllerInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceCaller{contract: contract}, nil
}

// NewAccessControllerInterfaceTransactor creates a new write-only instance of AccessControllerInterface, bound to a specific deployed contract.
func NewAccessControllerInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControllerInterfaceTransactor, error) {
	contract, err := bindAccessControllerInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceTransactor{contract: contract}, nil
}

// NewAccessControllerInterfaceFilterer creates a new log filterer instance of AccessControllerInterface, bound to a specific deployed contract.
func NewAccessControllerInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControllerInterfaceFilterer, error) {
	contract, err := bindAccessControllerInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControllerInterfaceFilterer{contract: contract}, nil
}

// bindAccessControllerInterface binds a generic wrapper to an already deployed contract.
func bindAccessControllerInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControllerInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControllerInterface *AccessControllerInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControllerInterface *AccessControllerInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControllerInterface *AccessControllerInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.AccessControllerInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControllerInterface *AccessControllerInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControllerInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControllerInterface *AccessControllerInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControllerInterface *AccessControllerInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControllerInterface.Contract.contract.Transact(opts, method, params...)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address user, bytes data) view returns(bool)
func (_AccessControllerInterface *AccessControllerInterfaceCaller) HasAccess(opts *bind.CallOpts, user common.Address, data []byte) (bool, error) {
	var out []interface{}
	err := _AccessControllerInterface.contract.Call(opts, &out, "hasAccess", user, data)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address user, bytes data) view returns(bool)
func (_AccessControllerInterface *AccessControllerInterfaceSession) HasAccess(user common.Address, data []byte) (bool, error) {
	return _AccessControllerInterface.Contract.HasAccess(&_AccessControllerInterface.CallOpts, user, data)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address user, bytes data) view returns(bool)
func (_AccessControllerInterface *AccessControllerInterfaceCallerSession) HasAccess(user common.Address, data []byte) (bool, error) {
	return _AccessControllerInterface.Contract.HasAccess(&_AccessControllerInterface.CallOpts, user, data)
}

// AggregatorInterfaceABI is the input ABI used to generate the binding from.
const AggregatorInterfaceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AggregatorInterface is an auto generated Go binding around an Ethereum contract.
type AggregatorInterface struct {
	AggregatorInterfaceCaller     // Read-only binding to the contract
	AggregatorInterfaceTransactor // Write-only binding to the contract
	AggregatorInterfaceFilterer   // Log filterer for contract events
}

// AggregatorInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorInterfaceSession struct {
	Contract     *AggregatorInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// AggregatorInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorInterfaceCallerSession struct {
	Contract *AggregatorInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// AggregatorInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorInterfaceTransactorSession struct {
	Contract     *AggregatorInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// AggregatorInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorInterfaceRaw struct {
	Contract *AggregatorInterface // Generic contract binding to access the raw methods on
}

// AggregatorInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorInterfaceCallerRaw struct {
	Contract *AggregatorInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorInterfaceTransactorRaw struct {
	Contract *AggregatorInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorInterface creates a new instance of AggregatorInterface, bound to a specific deployed contract.
func NewAggregatorInterface(address common.Address, backend bind.ContractBackend) (*AggregatorInterface, error) {
	contract, err := bindAggregatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterface{AggregatorInterfaceCaller: AggregatorInterfaceCaller{contract: contract}, AggregatorInterfaceTransactor: AggregatorInterfaceTransactor{contract: contract}, AggregatorInterfaceFilterer: AggregatorInterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorInterfaceCaller creates a new read-only instance of AggregatorInterface, bound to a specific deployed contract.
func NewAggregatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorInterfaceCaller, error) {
	contract, err := bindAggregatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceCaller{contract: contract}, nil
}

// NewAggregatorInterfaceTransactor creates a new write-only instance of AggregatorInterface, bound to a specific deployed contract.
func NewAggregatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorInterfaceTransactor, error) {
	contract, err := bindAggregatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceTransactor{contract: contract}, nil
}

// NewAggregatorInterfaceFilterer creates a new log filterer instance of AggregatorInterface, bound to a specific deployed contract.
func NewAggregatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorInterfaceFilterer, error) {
	contract, err := bindAggregatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorInterfaceFilterer{contract: contract}, nil
}

// bindAggregatorInterface binds a generic wrapper to an already deployed contract.
func bindAggregatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorInterface *AggregatorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorInterface.Contract.AggregatorInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorInterface *AggregatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.AggregatorInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorInterface *AggregatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.AggregatorInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorInterface *AggregatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorInterface *AggregatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorInterface *AggregatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorInterface.Contract.contract.Transact(opts, method, params...)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetAnswer(&_AggregatorInterface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorInterface.Contract.GetTimestamp(&_AggregatorInterface.CallOpts, roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestAnswer(&_AggregatorInterface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestRound(&_AggregatorInterface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorInterface.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestTimestamp(&_AggregatorInterface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorInterface *AggregatorInterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorInterface.Contract.LatestTimestamp(&_AggregatorInterface.CallOpts)
}

// AggregatorInterfaceAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AggregatorInterface contract.
type AggregatorInterfaceAnswerUpdatedIterator struct {
	Event *AggregatorInterfaceAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AggregatorInterfaceAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AggregatorInterfaceAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorInterfaceAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorInterfaceAnswerUpdated represents a AnswerUpdated event raised by the AggregatorInterface contract.
type AggregatorInterfaceAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
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

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
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
				// New log arrived, parse the event and forward to the user
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorInterface *AggregatorInterfaceFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorInterfaceAnswerUpdated, error) {
	event := new(AggregatorInterfaceAnswerUpdated)
	if err := _AggregatorInterface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorInterfaceNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AggregatorInterface contract.
type AggregatorInterfaceNewRoundIterator struct {
	Event *AggregatorInterfaceNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AggregatorInterfaceNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AggregatorInterfaceNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorInterfaceNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorInterfaceNewRound represents a NewRound event raised by the AggregatorInterface contract.
type AggregatorInterfaceNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
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

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
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
				// New log arrived, parse the event and forward to the user
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorInterface *AggregatorInterfaceFilterer) ParseNewRound(log types.Log) (*AggregatorInterfaceNewRound, error) {
	event := new(AggregatorInterfaceNewRound)
	if err := _AggregatorInterface.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorV2V3InterfaceABI is the input ABI used to generate the binding from.
const AggregatorV2V3InterfaceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AggregatorV2V3Interface is an auto generated Go binding around an Ethereum contract.
type AggregatorV2V3Interface struct {
	AggregatorV2V3InterfaceCaller     // Read-only binding to the contract
	AggregatorV2V3InterfaceTransactor // Write-only binding to the contract
	AggregatorV2V3InterfaceFilterer   // Log filterer for contract events
}

// AggregatorV2V3InterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV2V3InterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV2V3InterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorV2V3InterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV2V3InterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorV2V3InterfaceSession struct {
	Contract     *AggregatorV2V3Interface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AggregatorV2V3InterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorV2V3InterfaceCallerSession struct {
	Contract *AggregatorV2V3InterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// AggregatorV2V3InterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorV2V3InterfaceTransactorSession struct {
	Contract     *AggregatorV2V3InterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// AggregatorV2V3InterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceRaw struct {
	Contract *AggregatorV2V3Interface // Generic contract binding to access the raw methods on
}

// AggregatorV2V3InterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceCallerRaw struct {
	Contract *AggregatorV2V3InterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorV2V3InterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorV2V3InterfaceTransactorRaw struct {
	Contract *AggregatorV2V3InterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorV2V3Interface creates a new instance of AggregatorV2V3Interface, bound to a specific deployed contract.
func NewAggregatorV2V3Interface(address common.Address, backend bind.ContractBackend) (*AggregatorV2V3Interface, error) {
	contract, err := bindAggregatorV2V3Interface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3Interface{AggregatorV2V3InterfaceCaller: AggregatorV2V3InterfaceCaller{contract: contract}, AggregatorV2V3InterfaceTransactor: AggregatorV2V3InterfaceTransactor{contract: contract}, AggregatorV2V3InterfaceFilterer: AggregatorV2V3InterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorV2V3InterfaceCaller creates a new read-only instance of AggregatorV2V3Interface, bound to a specific deployed contract.
func NewAggregatorV2V3InterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorV2V3InterfaceCaller, error) {
	contract, err := bindAggregatorV2V3Interface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceCaller{contract: contract}, nil
}

// NewAggregatorV2V3InterfaceTransactor creates a new write-only instance of AggregatorV2V3Interface, bound to a specific deployed contract.
func NewAggregatorV2V3InterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorV2V3InterfaceTransactor, error) {
	contract, err := bindAggregatorV2V3Interface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceTransactor{contract: contract}, nil
}

// NewAggregatorV2V3InterfaceFilterer creates a new log filterer instance of AggregatorV2V3Interface, bound to a specific deployed contract.
func NewAggregatorV2V3InterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorV2V3InterfaceFilterer, error) {
	contract, err := bindAggregatorV2V3Interface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorV2V3InterfaceFilterer{contract: contract}, nil
}

// bindAggregatorV2V3Interface binds a generic wrapper to an already deployed contract.
func bindAggregatorV2V3Interface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorV2V3InterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.AggregatorV2V3InterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV2V3Interface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV2V3Interface.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Decimals() (uint8, error) {
	return _AggregatorV2V3Interface.Contract.Decimals(&_AggregatorV2V3Interface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorV2V3Interface.Contract.Decimals(&_AggregatorV2V3Interface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Description() (string, error) {
	return _AggregatorV2V3Interface.Contract.Description(&_AggregatorV2V3Interface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Description() (string, error) {
	return _AggregatorV2V3Interface.Contract.Description(&_AggregatorV2V3Interface.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetAnswer(&_AggregatorV2V3Interface.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetAnswer(&_AggregatorV2V3Interface.CallOpts, roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
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

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.GetRoundData(&_AggregatorV2V3Interface.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.GetRoundData(&_AggregatorV2V3Interface.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetTimestamp(&_AggregatorV2V3Interface.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.GetTimestamp(&_AggregatorV2V3Interface.CallOpts, roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestAnswer(&_AggregatorV2V3Interface.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestAnswer() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestAnswer(&_AggregatorV2V3Interface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestRound() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestRound(&_AggregatorV2V3Interface.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestRound() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestRound(&_AggregatorV2V3Interface.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
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

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.LatestRoundData(&_AggregatorV2V3Interface.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV2V3Interface.Contract.LatestRoundData(&_AggregatorV2V3Interface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestTimestamp(&_AggregatorV2V3Interface.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.LatestTimestamp(&_AggregatorV2V3Interface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV2V3Interface.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceSession) Version() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.Version(&_AggregatorV2V3Interface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceCallerSession) Version() (*big.Int, error) {
	return _AggregatorV2V3Interface.Contract.Version(&_AggregatorV2V3Interface.CallOpts)
}

// AggregatorV2V3InterfaceAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the AggregatorV2V3Interface contract.
type AggregatorV2V3InterfaceAnswerUpdatedIterator struct {
	Event *AggregatorV2V3InterfaceAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorV2V3InterfaceAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorV2V3InterfaceAnswerUpdated represents a AnswerUpdated event raised by the AggregatorV2V3Interface contract.
type AggregatorV2V3InterfaceAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
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

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
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
				// New log arrived, parse the event and forward to the user
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) ParseAnswerUpdated(log types.Log) (*AggregatorV2V3InterfaceAnswerUpdated, error) {
	event := new(AggregatorV2V3InterfaceAnswerUpdated)
	if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorV2V3InterfaceNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the AggregatorV2V3Interface contract.
type AggregatorV2V3InterfaceNewRoundIterator struct {
	Event *AggregatorV2V3InterfaceNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AggregatorV2V3InterfaceNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AggregatorV2V3InterfaceNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AggregatorV2V3InterfaceNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AggregatorV2V3InterfaceNewRound represents a NewRound event raised by the AggregatorV2V3Interface contract.
type AggregatorV2V3InterfaceNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
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

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
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
				// New log arrived, parse the event and forward to the user
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_AggregatorV2V3Interface *AggregatorV2V3InterfaceFilterer) ParseNewRound(log types.Log) (*AggregatorV2V3InterfaceNewRound, error) {
	event := new(AggregatorV2V3InterfaceNewRound)
	if err := _AggregatorV2V3Interface.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AggregatorV3InterfaceABI is the input ABI used to generate the binding from.
const AggregatorV3InterfaceABI = "[{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AggregatorV3Interface is an auto generated Go binding around an Ethereum contract.
type AggregatorV3Interface struct {
	AggregatorV3InterfaceCaller     // Read-only binding to the contract
	AggregatorV3InterfaceTransactor // Write-only binding to the contract
	AggregatorV3InterfaceFilterer   // Log filterer for contract events
}

// AggregatorV3InterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorV3InterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV3InterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorV3InterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV3InterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorV3InterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorV3InterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorV3InterfaceSession struct {
	Contract     *AggregatorV3Interface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AggregatorV3InterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorV3InterfaceCallerSession struct {
	Contract *AggregatorV3InterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// AggregatorV3InterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorV3InterfaceTransactorSession struct {
	Contract     *AggregatorV3InterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// AggregatorV3InterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorV3InterfaceRaw struct {
	Contract *AggregatorV3Interface // Generic contract binding to access the raw methods on
}

// AggregatorV3InterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorV3InterfaceCallerRaw struct {
	Contract *AggregatorV3InterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorV3InterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorV3InterfaceTransactorRaw struct {
	Contract *AggregatorV3InterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorV3Interface creates a new instance of AggregatorV3Interface, bound to a specific deployed contract.
func NewAggregatorV3Interface(address common.Address, backend bind.ContractBackend) (*AggregatorV3Interface, error) {
	contract, err := bindAggregatorV3Interface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3Interface{AggregatorV3InterfaceCaller: AggregatorV3InterfaceCaller{contract: contract}, AggregatorV3InterfaceTransactor: AggregatorV3InterfaceTransactor{contract: contract}, AggregatorV3InterfaceFilterer: AggregatorV3InterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorV3InterfaceCaller creates a new read-only instance of AggregatorV3Interface, bound to a specific deployed contract.
func NewAggregatorV3InterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorV3InterfaceCaller, error) {
	contract, err := bindAggregatorV3Interface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceCaller{contract: contract}, nil
}

// NewAggregatorV3InterfaceTransactor creates a new write-only instance of AggregatorV3Interface, bound to a specific deployed contract.
func NewAggregatorV3InterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorV3InterfaceTransactor, error) {
	contract, err := bindAggregatorV3Interface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceTransactor{contract: contract}, nil
}

// NewAggregatorV3InterfaceFilterer creates a new log filterer instance of AggregatorV3Interface, bound to a specific deployed contract.
func NewAggregatorV3InterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorV3InterfaceFilterer, error) {
	contract, err := bindAggregatorV3Interface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorV3InterfaceFilterer{contract: contract}, nil
}

// bindAggregatorV3Interface binds a generic wrapper to an already deployed contract.
func bindAggregatorV3Interface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorV3InterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorV3Interface *AggregatorV3InterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.AggregatorV3InterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorV3Interface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorV3Interface *AggregatorV3InterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorV3Interface *AggregatorV3InterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorV3Interface.Contract.contract.Transact(opts, method, params...)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Decimals() (uint8, error) {
	return _AggregatorV3Interface.Contract.Decimals(&_AggregatorV3Interface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Decimals() (uint8, error) {
	return _AggregatorV3Interface.Contract.Decimals(&_AggregatorV3Interface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Description() (string, error) {
	return _AggregatorV3Interface.Contract.Description(&_AggregatorV3Interface.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Description() (string, error) {
	return _AggregatorV3Interface.Contract.Description(&_AggregatorV3Interface.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
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

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.GetRoundData(&_AggregatorV3Interface.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.GetRoundData(&_AggregatorV3Interface.CallOpts, _roundId)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
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

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.LatestRoundData(&_AggregatorV3Interface.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _AggregatorV3Interface.Contract.LatestRoundData(&_AggregatorV3Interface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV3Interface *AggregatorV3InterfaceCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AggregatorV3Interface.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV3Interface *AggregatorV3InterfaceSession) Version() (*big.Int, error) {
	return _AggregatorV3Interface.Contract.Version(&_AggregatorV3Interface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_AggregatorV3Interface *AggregatorV3InterfaceCallerSession) Version() (*big.Int, error) {
	return _AggregatorV3Interface.Contract.Version(&_AggregatorV3Interface.CallOpts)
}

// AggregatorValidatorInterfaceABI is the input ABI used to generate the binding from.
const AggregatorValidatorInterfaceABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"previousRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"previousAnswer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"currentRoundId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"currentAnswer\",\"type\":\"int256\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AggregatorValidatorInterface is an auto generated Go binding around an Ethereum contract.
type AggregatorValidatorInterface struct {
	AggregatorValidatorInterfaceCaller     // Read-only binding to the contract
	AggregatorValidatorInterfaceTransactor // Write-only binding to the contract
	AggregatorValidatorInterfaceFilterer   // Log filterer for contract events
}

// AggregatorValidatorInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AggregatorValidatorInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AggregatorValidatorInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AggregatorValidatorInterfaceSession struct {
	Contract     *AggregatorValidatorInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AggregatorValidatorInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AggregatorValidatorInterfaceCallerSession struct {
	Contract *AggregatorValidatorInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// AggregatorValidatorInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AggregatorValidatorInterfaceTransactorSession struct {
	Contract     *AggregatorValidatorInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// AggregatorValidatorInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceRaw struct {
	Contract *AggregatorValidatorInterface // Generic contract binding to access the raw methods on
}

// AggregatorValidatorInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceCallerRaw struct {
	Contract *AggregatorValidatorInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// AggregatorValidatorInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AggregatorValidatorInterfaceTransactorRaw struct {
	Contract *AggregatorValidatorInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAggregatorValidatorInterface creates a new instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterface(address common.Address, backend bind.ContractBackend) (*AggregatorValidatorInterface, error) {
	contract, err := bindAggregatorValidatorInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterface{AggregatorValidatorInterfaceCaller: AggregatorValidatorInterfaceCaller{contract: contract}, AggregatorValidatorInterfaceTransactor: AggregatorValidatorInterfaceTransactor{contract: contract}, AggregatorValidatorInterfaceFilterer: AggregatorValidatorInterfaceFilterer{contract: contract}}, nil
}

// NewAggregatorValidatorInterfaceCaller creates a new read-only instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceCaller(address common.Address, caller bind.ContractCaller) (*AggregatorValidatorInterfaceCaller, error) {
	contract, err := bindAggregatorValidatorInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceCaller{contract: contract}, nil
}

// NewAggregatorValidatorInterfaceTransactor creates a new write-only instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*AggregatorValidatorInterfaceTransactor, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceTransactor{contract: contract}, nil
}

// NewAggregatorValidatorInterfaceFilterer creates a new log filterer instance of AggregatorValidatorInterface, bound to a specific deployed contract.
func NewAggregatorValidatorInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*AggregatorValidatorInterfaceFilterer, error) {
	contract, err := bindAggregatorValidatorInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AggregatorValidatorInterfaceFilterer{contract: contract}, nil
}

// bindAggregatorValidatorInterface binds a generic wrapper to an already deployed contract.
func bindAggregatorValidatorInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AggregatorValidatorInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.AggregatorValidatorInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AggregatorValidatorInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.contract.Transact(opts, method, params...)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactor) Validate(opts *bind.TransactOpts, previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.contract.Transact(opts, "validate", previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// Validate is a paid mutator transaction binding the contract method 0xbeed9b51.
//
// Solidity: function validate(uint256 previousRoundId, int256 previousAnswer, uint256 currentRoundId, int256 currentAnswer) returns(bool)
func (_AggregatorValidatorInterface *AggregatorValidatorInterfaceTransactorSession) Validate(previousRoundId *big.Int, previousAnswer *big.Int, currentRoundId *big.Int, currentAnswer *big.Int) (*types.Transaction, error) {
	return _AggregatorValidatorInterface.Contract.Validate(&_AggregatorValidatorInterface.TransactOpts, previousRoundId, previousAnswer, currentRoundId, currentAnswer)
}

// LinkTokenInterfaceABI is the input ABI used to generate the binding from.
const LinkTokenInterfaceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"remaining\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimalPlaces\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenName\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalTokensIssued\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// LinkTokenInterface is an auto generated Go binding around an Ethereum contract.
type LinkTokenInterface struct {
	LinkTokenInterfaceCaller     // Read-only binding to the contract
	LinkTokenInterfaceTransactor // Write-only binding to the contract
	LinkTokenInterfaceFilterer   // Log filterer for contract events
}

// LinkTokenInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type LinkTokenInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LinkTokenInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LinkTokenInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LinkTokenInterfaceSession struct {
	Contract     *LinkTokenInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// LinkTokenInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LinkTokenInterfaceCallerSession struct {
	Contract *LinkTokenInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// LinkTokenInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LinkTokenInterfaceTransactorSession struct {
	Contract     *LinkTokenInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// LinkTokenInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type LinkTokenInterfaceRaw struct {
	Contract *LinkTokenInterface // Generic contract binding to access the raw methods on
}

// LinkTokenInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LinkTokenInterfaceCallerRaw struct {
	Contract *LinkTokenInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// LinkTokenInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LinkTokenInterfaceTransactorRaw struct {
	Contract *LinkTokenInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLinkTokenInterface creates a new instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterface(address common.Address, backend bind.ContractBackend) (*LinkTokenInterface, error) {
	contract, err := bindLinkTokenInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterface{LinkTokenInterfaceCaller: LinkTokenInterfaceCaller{contract: contract}, LinkTokenInterfaceTransactor: LinkTokenInterfaceTransactor{contract: contract}, LinkTokenInterfaceFilterer: LinkTokenInterfaceFilterer{contract: contract}}, nil
}

// NewLinkTokenInterfaceCaller creates a new read-only instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceCaller(address common.Address, caller bind.ContractCaller) (*LinkTokenInterfaceCaller, error) {
	contract, err := bindLinkTokenInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceCaller{contract: contract}, nil
}

// NewLinkTokenInterfaceTransactor creates a new write-only instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkTokenInterfaceTransactor, error) {
	contract, err := bindLinkTokenInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceTransactor{contract: contract}, nil
}

// NewLinkTokenInterfaceFilterer creates a new log filterer instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkTokenInterfaceFilterer, error) {
	contract, err := bindLinkTokenInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceFilterer{contract: contract}, nil
}

// bindLinkTokenInterface binds a generic wrapper to an already deployed contract.
func bindLinkTokenInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LinkTokenInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenInterface *LinkTokenInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Approve(&_LinkTokenInterface.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Approve(&_LinkTokenInterface.TransactOpts, spender, value)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) DecreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "decreaseApproval", spender, addedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) DecreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.DecreaseApproval(&_LinkTokenInterface.TransactOpts, spender, addedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) DecreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.DecreaseApproval(&_LinkTokenInterface.TransactOpts, spender, addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) IncreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "increaseApproval", spender, subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceSession) IncreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.IncreaseApproval(&_LinkTokenInterface.TransactOpts, spender, subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) IncreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.IncreaseApproval(&_LinkTokenInterface.TransactOpts, spender, subtractedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Transfer(&_LinkTokenInterface.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Transfer(&_LinkTokenInterface.TransactOpts, to, value)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transferAndCall", to, value, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferAndCall(&_LinkTokenInterface.TransactOpts, to, value, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferAndCall(&_LinkTokenInterface.TransactOpts, to, value, data)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferFrom(&_LinkTokenInterface.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferFrom(&_LinkTokenInterface.TransactOpts, from, to, value)
}

// OffchainAggregatorABI is the input ABI used to generate the binding from.
const OffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_requesterAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"billingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requesterAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"_encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_encoded\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_requesterAccessController\",\"type\":\"address\"}],\"name\":\"setRequesterAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newValidator\",\"type\":\"address\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"_rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"_ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"_rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OffchainAggregatorBin is the compiled bytecode used for deploying new contracts.
var OffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b50604051620056013803806200560183398181016040526101a08110156200003957600080fd5b815160208301516040808501516060860151608087015160a088015160c089015160e08a01516101008b01516101208c01516101408d01516101608e01516101808f0180519b519d9f9c9e9a9d999c989b979a96999598949793969295919491939282019284640100000000821115620000b257600080fd5b908301906020820185811115620000c857600080fd5b8251640100000000811182820188101715620000e357600080fd5b82525081516020918201929091019080838360005b8381101562000112578181015183820152602001620000f8565b50505050905090810190601f168015620001405780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b03191633179055508c8c8c8c8c8c896200016f87878787876200029c565b6200017a816200038e565b6001600160601b0319606083901b1660805262000196620005c4565b620001a0620005c4565b60005b601f8160ff161015620001f0576001838260ff16601f8110620001c257fe5b61ffff909216602092909202015260018260ff8316601f8110620001e257fe5b6020020152600101620001a3565b5062000200600483601f620005e3565b5062000210600882601f62000680565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e05250508351620002599350602e9250602085019150620006b1565b50620002658362000407565b6200027087620004df565b50505050601791820b820b604090811b60a05290820b90910b901b60c052506200074a95505050505050565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b0390811690821681146200040357600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b0316331462000467576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d546001600160a01b0390811690821681146200040357602d80546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349281900390910190a15050565b6000546001600160a01b031633146200053f576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b0368010000000000000000909104811690821681146200040357602c8054600160401b600160e01b031916680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35050565b604051806103e00160405280601f906020820280368337509192915050565b6002830191839082156200066e5791602002820160005b838211156200063c57835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302620005fa565b80156200066c5782816101000a81549061ffff02191690556002016020816001010492830192600103026200063c565b505b506200067c92915062000733565b5090565b82601f81019282156200066e579160200282015b828111156200066e57825182559160200191906001019062000694565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620006e957600085556200066e565b82601f106200070457805160ff19168380011785556200066e565b828001600101855582156200066e57918201828111156200066e57825182559160200191906001019062000694565b5b808211156200067c576000815560010162000734565b60805160601c60a05160401c60c05160401c60e05160f81c614e52620007af60003980610e3252508061185952806132f2525080610d7952806132c5525080610d5552806126435280612733528061372c5280613d7352806142025250614e526000f3fe608060405234801561001057600080fd5b50600436106102925760003560e01c80638da5cb5b11610160578063bd824706116100d8578063e5fe45771161008c578063f2fde38b11610071578063f2fde38b14610aa6578063fbffd2c114610acc578063feaf968c14610af257610292565b8063e5fe457714610a0e578063eb5dcd6c14610a7857610292565b8063c9807539116100bd578063c9807539146108b5578063d09dc339146109c9578063e4902f82146109d157610292565b8063bd82470614610844578063c10753291461088957610292565b80639c849b301161012f578063b121e14711610114578063b121e147146107e4578063b5ab58dc1461080a578063b633620c1461082757610292565b80639c849b30146106fc5780639e3ceeab146107be57610292565b80638da5cb5b1461065257806398e5b12a1461065a578063996e8298146106815780639a6fc8f51461068957610292565b8063585aa7de1161020e57806379ba5097116101c257806381ff7048116101a757806381ff7048146105d35780638205bf6a146106245780638ac28d5a1461062c57610292565b806379ba509714610573578063814118341461057b57610292565b806370da2f67116101f357806370da2f67146104e657806370efdf2d146104ee5780637284e416146104f657610292565b8063585aa7de146103b1578063668a0f02146104de57610292565b806329937268116102655780633a5381b51161024a5780633a5381b51461039957806350d25bcd146103a157806354fd4d50146103a957610292565b8063299372681461033a578063313ce5671461037b57610292565b80630eafb25b146102975780631327d3d8146102cf5780631b6b6d23146102f757806322adbc781461031b575b600080fd5b6102bd600480360360208110156102ad57600080fd5b50356001600160a01b0316610afa565b60408051918252519081900360200190f35b6102f5600480360360208110156102e557600080fd5b50356001600160a01b0316610c5a565b005b6102ff610d53565b604080516001600160a01b039092168252519081900360200190f35b610323610d77565b6040805160179290920b8252519081900360200190f35b610342610d9b565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b610383610e30565b6040805160ff9092168252519081900360200190f35b6102ff610e54565b6102bd610e6f565b6102bd610eab565b6102f5600480360360a08110156103c757600080fd5b8101906020810181356401000000008111156103e257600080fd5b8201836020820111156103f457600080fd5b8035906020019184602083028401116401000000008311171561041657600080fd5b91939092909160208101903564010000000081111561043457600080fd5b82018360208201111561044657600080fd5b8035906020019184602083028401116401000000008311171561046857600080fd5b9193909260ff8335169267ffffffffffffffff60208201351692919060608101906040013564010000000081111561049f57600080fd5b8201836020820111156104b157600080fd5b803590602001918460018302840111640100000000831117156104d357600080fd5b509092509050610eb0565b6102bd611831565b610323611857565b6102ff61187b565b6104fe61188a565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610538578181015183820152602001610520565b50505050905090810190601f1680156105655780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102f561193e565b610583611a0c565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156105bf5781810151838201526020016105a7565b505050509050019250505060405180910390f35b6105db611a6d565b6040805163ffffffff94851681529290931660208301527fffffffffffffffffffffffffffffffff00000000000000000000000000000000168183015290519081900360600190f35b6102bd611a8e565b6102f56004803603602081101561064257600080fd5b50356001600160a01b0316611ae9565b6102ff611b63565b610662611b72565b6040805169ffffffffffffffffffff9092168252519081900360200190f35b6102ff611db1565b6106b26004803603602081101561069f57600080fd5b503569ffffffffffffffffffff16611dc0565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6102f56004803603604081101561071257600080fd5b81019060208101813564010000000081111561072d57600080fd5b82018360208201111561073f57600080fd5b8035906020019184602083028401116401000000008311171561076157600080fd5b91939092909160208101903564010000000081111561077f57600080fd5b82018360208201111561079157600080fd5b803590602001918460208302840111640100000000831117156107b357600080fd5b509092509050611f14565b6102f5600480360360208110156107d457600080fd5b50356001600160a01b031661214d565b6102f5600480360360208110156107fa57600080fd5b50356001600160a01b031661223b565b6102bd6004803603602081101561082057600080fd5b5035612334565b6102bd6004803603602081101561083d57600080fd5b503561236a565b6102f5600480360360a081101561085a57600080fd5b5063ffffffff8135811691602081013582169160408201358116916060810135821691608090910135166123bf565b6102f56004803603604081101561089f57600080fd5b506001600160a01b038135169060200135612507565b6102f5600480360360808110156108cb57600080fd5b8101906020810181356401000000008111156108e657600080fd5b8201836020820111156108f857600080fd5b8035906020019184600183028401116401000000008311171561091a57600080fd5b91939092909160208101903564010000000081111561093857600080fd5b82018360208201111561094a57600080fd5b8035906020019184602083028401116401000000008311171561096c57600080fd5b91939092909160208101903564010000000081111561098a57600080fd5b82018360208201111561099c57600080fd5b803590602001918460208302840111640100000000831117156109be57600080fd5b919350915035612836565b6102bd613727565b6109f7600480360360208110156109e757600080fd5b50356001600160a01b03166137d8565b6040805161ffff9092168252519081900360200190f35b610a16613891565b604080517fffffffffffffffffffffffffffffffff00000000000000000000000000000000909616865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b6102f560048036036040811015610a8e57600080fd5b506001600160a01b0381358116916020013516613980565b6102f560048036036020811015610abc57600080fd5b50356001600160a01b0316613adc565b6102f560048036036020811015610ae257600080fd5b50356001600160a01b0316613ba4565b6106b2613c0c565b6000610b04614c89565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115610b4857fe5b6002811115610b5357fe5b9052509050600081602001516002811115610b6a57fe5b1415610b7a576000915050610c55565b610b82614ca0565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830181905270010000000000000000000000000000000090910490921660808201528251909160009160019060049060ff16601f8110610c0e57fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff16601f8110610c4c57fe5b01540301925050505b919050565b6000546001600160a01b03163314610cb9576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602c546001600160a01b036801000000000000000090910481169082168114610d4f57602c80547fffffffff0000000000000000000000000000000000000000ffffffffffffffff16680100000000000000006001600160a01b0385811691820292909217909255604051908316907fcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d90600090a35b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000610dab614ca0565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b7f000000000000000000000000000000000000000000000000000000000000000081565b602c546801000000000000000090046001600160a01b031690565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b6020526040902054601790810b900b90565b600481565b868560ff8616601f831115610f0c576040805162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b60008111610f61576040805162461bcd60e51b815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b818314610f9f5760405162461bcd60e51b8152600401808060200182810382526024815260200180614e226024913960400191505060405180910390fd5b806003028311610ff6576040805162461bcd60e51b815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b6000546001600160a01b03163314611055576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602854156111f957602880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101916000918390811061109257fe5b6000918252602082200154602980546001600160a01b03909216935090849081106110b957fe5b6000918252602090912001546001600160a01b031690506110d981613cab565b6001600160a01b0380831660009081526027602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559284168252902080549091169055602880548061113557fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055602980548061119857fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611055915050565b60005b8a811015611607576000602760008e8e8581811061121657fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff16600281111561124c57fe5b1461129e576040805162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260016020820152602760008e8e858181106112c557fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561135057fe5b02179055506000915060069050818c8c8581811061136a57fe5b6001600160a01b03602091820293909301358316845283019390935260409091016000205416919091141590506113e8576040805162461bcd60e51b815260206004820152601160248201527f7061796565206d75737420626520736574000000000000000000000000000000604482015290519081900360640190fd5b6000602760008c8c858181106113fa57fe5b602090810292909201356001600160a01b031683525081019190915260400160002054610100900460ff16600281111561143057fe5b14611482576040805162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260026020820152602760008c8c858181106114a957fe5b602090810292909201356001600160a01b031683525081810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561153457fe5b021790555090505060288c8c8381811061154a57fe5b835460018101855560009485526020948590200180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03959092029390930135939093169290921790555060298a8a838181106115ac57fe5b835460018181018655600095865260209586902090910180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03969093029490940135949094161790915550016111fc565b50602a805460ff89167501000000000000000000000000000000000000000000027fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff909116179055602c80544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009091161793849055909104811691166116d330828f8f8f8f8f8f8f8f613edb565b602a60000160006101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506000602a60000160106101000a81548164ffffffffff021916908364ffffffffff1602179055507f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb982828f8f8f8f8f8f8f8f604051808b63ffffffff1681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600083820152604051601f909101601f19169092018290039f50909d5050505050505050505050505050a150505050505050505050505050565b602a54760100000000000000000000000000000000000000000000900463ffffffff1690565b7f000000000000000000000000000000000000000000000000000000000000000081565b602d546001600160a01b031690565b602e8054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156119345780601f1061190957610100808354040283529160200191611934565b820191906000526020600020905b81548152906001019060200180831161191757829003601f168201915b5050505050905090565b6001546001600160a01b0316331461199d576040805162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060602980548060200260200160405190810160405280929190818152602001828054801561193457602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611a46575050505050905090565b602c54602a5463ffffffff808316926401000000009004169060801b909192565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6001600160a01b03818116600090815260066020526040902054163314611b57576040805162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b611b6081613cab565b50565b6000546001600160a01b031681565b600080546001600160a01b0316331480611c4e5750602d54604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b158015611c2157600080fd5b505afa158015611c35573d6000803e3d6000fd5b505050506040513d6020811015611c4b57600080fd5b50515b611c9f576040805162461bcd60e51b815260206004820152601d60248201527f4f6e6c79206f776e6572267265717565737465722063616e2063616c6c000000604482015290519081900360640190fd5b611ca7614cce565b506040805160808082018352602a549081901b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000016808352700100000000000000000000000000000000820464ffffffffff81166020808601919091527501000000000000000000000000000000000000000000840460ff9081168688015276010000000000000000000000000000000000000000000090940463ffffffff9081166060808801919091528751948552600884901c909116918401919091529216818501529251919233927f3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037929181900390910190a2806060015160010163ffffffff1691505090565b6003546001600160a01b031690565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e74000000000000000000000000000000000081525090611e9b5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b83811015611e60578181015183820152602001611e48565b50505050905090810190601f168015611e8d5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50611ea4614c89565b5050505063ffffffff83166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b6000546001600160a01b03163314611f73576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b828114611fc7576040805162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015612146576000858583818110611fe057fe5b905060200201356001600160a01b03169050600084848481811061200057fe5b6001600160a01b03858116600090815260066020908152604090912054920293909301358316935090911690508015808061204c5750826001600160a01b0316826001600160a01b0316145b61209d576040805162461bcd60e51b815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b03848116600090815260066020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168583169081179091559083161461213657826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505060019092019150611fca9050565b5050505050565b6000546001600160a01b031633146121ac576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b602d546001600160a01b039081169082168114610d4f57602d80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349281900390910190a15050565b6001600160a01b038181166000908152600760205260409020541633146122a9576040805162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b6001600160a01b0381811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600063ffffffff82111561234a57506000610c55565b5063ffffffff166000908152602b6020526040902054601790810b900b90565b600063ffffffff82111561238057506000610c55565b5063ffffffff166000908152602b60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6003546000546001600160a01b0391821691163314806124995750604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b03861694636b14daf8946000939190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b15801561246c57600080fd5b505afa158015612480573d6000803e3d6000fd5b505050506040513d602081101561249657600080fd5b50515b6124ea576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6124f2613fdf565b6124ff86868686866143bb565b505050505050565b6000546001600160a01b03163314806125e25750600354604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b1580156125b557600080fd5b505afa1580156125c9573d6000803e3d6000fd5b505050506040513d60208110156125df57600080fd5b50515b612633576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b600061263d614535565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b1580156126ae57600080fd5b505afa1580156126c2573d6000803e3d6000fd5b505050506040513d60208110156126d857600080fd5b5051905081811015612731576040805162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb8561276d85850387614716565b6040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b1580156127b357600080fd5b505af11580156127c7573d6000803e3d6000fd5b505050506040513d60208110156127dd57600080fd5b5051612830576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b60005a9050612849888888888888614730565b361461289c576040805162461bcd60e51b815260206004820152601960248201527f7472616e736d6974206d65737361676520746f6f206c6f6e6700000000000000604482015290519081900360640190fd5b6128a4614cf5565b6040805160808082018352602a549081901b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000168252700100000000000000000000000000000000810464ffffffffff1660208301527501000000000000000000000000000000000000000000810460ff169282019290925276010000000000000000000000000000000000000000000090910463ffffffff166060808301919091529082526000908a908a9081101561295d57600080fd5b81359160208101359181019060608101604082013564010000000081111561298457600080fd5b82018360208201111561299657600080fd5b803590602001918460208302840111640100000000831117156129b857600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050505060408801525050506080840182905283515190925060589190911b907fffffffffffffffffffffffffffffffff00000000000000000000000000000000808316911614612a7f576040805162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b608083015183516020015164ffffffffff808316911610612ae7576040805162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b83516040015160ff168911612b43576040805162461bcd60e51b815260206004820152601560248201527f6e6f7420656e6f756768207369676e6174757265730000000000000000000000604482015290519081900360640190fd5b601f891115612b99576040805162461bcd60e51b815260206004820152601360248201527f746f6f206d616e79207369676e61747572657300000000000000000000000000604482015290519081900360640190fd5b868914612bed576040805162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b601f8460400151511115612c48576040805162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b83600001516040015160020260ff1684604001515111612caf576040805162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b8867ffffffffffffffff81118015612cc657600080fd5b506040519080825280601f01601f191660200182016040528015612cf1576020820181803683370190505b50606085015260005b60ff81168a1115612d6257868160ff1660208110612d1457fe5b1a60f81b85606001518260ff1681518110612d2b57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600101612cfa565b5083604001515167ffffffffffffffff81118015612d7f57600080fd5b506040519080825280601f01601f191660200182016040528015612daa576020820181803683370190505b506020850152612db8614d29565b60005b8560400151518160ff161015612ebe576000858260ff1660208110612ddc57fe5b1a90508281601f8110612deb57fe5b602002015115612e42576040805162461bcd60e51b815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b6001838260ff16601f8110612e5357fe5b91151560209283029190910152869060ff8416908110612e6f57fe5b1a60f81b87602001518360ff1681518110612e8657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101612dbb565b50612ec7614c89565b336000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115612f0257fe5b6002811115612f0d57fe5b9052509050600281602001516002811115612f2457fe5b148015612f5857506029816000015160ff1681548110612f4057fe5b6000918252602090912001546001600160a01b031633145b612fa9576040805162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b5050835164ffffffffff90911660209091015250506040516000908a908a908083838082843760405192018290039091209450612fea9350614d2992505050565b612ff2614c89565b60005b898110156131eb5760006001858760600151848151811061301257fe5b60209101015160f81c601b018e8e8681811061302a57fe5b905060200201358d8d8781811061303d57fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015613098573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526027602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156130e757fe5b60028111156130f257fe5b905250925060018360200151600281111561310957fe5b1461315b576040805162461bcd60e51b815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff16601f811061316c57fe5b6020020151156131c3576040805162461bcd60e51b815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16601f81106131d857fe5b9115156020909202015250600101612ff5565b5050505060005b60018260400151510381101561329c5760008260400151826001018151811061321757fe5b602002602001015160170b8360400151838151811061323257fe5b602002602001015160170b1315905080613293576040805162461bcd60e51b815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b506001016131f2565b506040810151805160009190600281049081106132b557fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b1315801561331b57507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b61336c576040805162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b81516060908101805163ffffffff60019091018116909152604080518082018252601785810b80835267ffffffffffffffff42811660208086019182528a5189015188166000908152602b82528781209651875493519094167801000000000000000000000000000000000000000000000000029390950b77ffffffffffffffffffffffffffffffffffffffffffffffff9081167fffffffffffffffff0000000000000000000000000000000000000000000000009093169290921790911691909117909355875186015184890151848a01516080808c015188519586523386890181905291860181905260a0988601898152845199870199909952835194909916997ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451998c999298949793969095909492939185019260c086019289820192909102908190849084905b838110156134cf5781810151838201526020016134b7565b50505050905001838103825285818151815260200191508051906020019080838360005b8381101561350b5781810151838201526020016134f3565b50505050905090810190601f1680156135385780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390a281516060015160408051428152905160009263ffffffff16917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271919081900360200190a381600001516060015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a36135ed8260000151606001518260170b614748565b5080518051602a8054602084015160408501516060909501517fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660809490941c939093177fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000064ffffffffff90941693909302929092177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000060ff90941693909302929092177fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1676010000000000000000000000000000000000000000000063ffffffff9283160217909155821061370e57fe5b61371c828260200151614859565b505050505050505050565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b15801561379757600080fd5b505afa1580156137ab573d6000803e3d6000fd5b505050506040513d60208110156137c157600080fd5b5051905060006137cf614535565b90910391505090565b60006137e2614c89565b6001600160a01b0383166000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561382657fe5b600281111561383157fe5b905250905060008160200151600281111561384857fe5b1415613858576000915050610c55565b60016004826000015160ff16601f811061386e57fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b6000808080803332146138eb576040805162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050602a5463ffffffff760100000000000000000000000000000000000000000000820481166000908152602b6020526040902054608083901b96700100000000000000000000000000000000909304600881901c909216955064ffffffffff9091169350601781900b92507801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6001600160a01b038281166000908152600660205260409020541633146139ee576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b336001600160a01b0382161415613a4c576040805162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b6001600160a01b03808316600090815260076020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114613ad7576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b6000546001600160a01b03163314613b3b576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000546001600160a01b03163314613c03576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b611b6081614ab4565b602a54760100000000000000000000000000000000000000000000900463ffffffff166000808080613c3c614c89565b5050505063ffffffff82166000908152602b6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052939493900b9290915081908490565b613cb3614c89565b6001600160a01b0382166000908152602760209081526040918290208251808401909352805460ff80821685529192840191610100909104166002811115613cf757fe5b6002811115613d0257fe5b90525090506000613d1283610afa565b90508015613ad7576001600160a01b0380841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b158015613dbb57600080fd5b505af1158015613dcf573d6000803e3d6000fd5b505050506040513d6020811015613de557600080fd5b5051613e38576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff16601f8110613e4e57fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16601f8110613e8957fe5b0155604080516001600160a01b0380871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b60008a8a8a8a8a8a8a8a8a8a604051602001808b6001600160a01b031681526020018a67ffffffffffffffff16815260200180602001806020018760ff1681526020018667ffffffffffffffff1681526020018060200184810384528c8c82818152602001925060200280828437600083820152601f01601f191690910185810384528a8152602090810191508b908b0280828437600083820152601f01601f191690910185810383528681526020019050868680828437600081840152601f19601f8201169050808301925050509d50505050505050505050505050506040516020818303038152906040528051906020012090509a9950505050505050505050565b613fe7614ca0565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c010000000000000000000000008104831660608301527001000000000000000000000000000000009004909116608082015261405e614d29565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411614077579050505050505090506140be614d29565b604080516103e081019182905290600890601f9082845b8154815260200190600101908083116140d557505050505090506060602980548060200260200160405190810160405280929190818152602001828054801561414757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311614129575b5050505050905060005b815181101561439f57600060018483601f811061416a57fe5b6020020151039050600060018684601f811061418257fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115614394576000600660008787815181106141c257fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b031690507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb82846040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b15801561427757600080fd5b505af115801561428b573d6000803e3d6000fd5b505050506040513d60208110156142a157600080fd5b50516142f4576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018886601f811061430257fe5b61ffff909216602092909202015260018786601f811061431e57fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29087908790811061435357fe5b6020026020010151828460405180846001600160a01b03168152602001836001600160a01b03168152602001828152602001935050505060405180910390a1505b505050600101614151565b506143ad600484601f614d48565b50612146600883601f614dde565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b600061453f614d29565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116145585790505050505050905060005b601f8110156145c85760018282601f81106145b157fe5b60200201510361ffff16929092019160010161459a565b506145d1614ca0565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606080860182905270010000000000000000000000000000000090940490921660808501526029805486518184028101840190975280875297909202633b9aca00029693949293908301828280156146a157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311614683575b505050505090506146b0614d29565b604080516103e081019182905290600890601f9082845b8154815260200190600101908083116146c7575050505050905060005b825181101561470e5760018282601f81106146fb57fe5b60200201510395909501946001016146e4565b505050505090565b60008183101561472757508161472a565b50805b92915050565b602083810286019082020160e4019695505050505050565b602c546801000000000000000090046001600160a01b03168061476b5750610d4f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602b602090815260408083205481517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810195909552601790810b900b60248501819052948916604485015260648401889052516001600160a01b0387169363beed9b5193620186a09360848084019491939192918390030190829088803b15801561482a57600080fd5b5087f19350505050801561485057506040513d602081101561484b57600080fd5b505160015b6124ff57612146565b614861614c89565b336000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561489c57fe5b60028111156148a757fe5b90525090506148b4614ca0565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000082048116606084015270010000000000000000000000000000000090910416608082015281516103e0810192839052909161497f918591600490601f90826000855b82829054906101000a900461ffff1661ffff168152602001906002019060208260010104928301926001038202915080841161493d5790505050505050614b43565b61498d90600490601f614d48565b506002826020015160028111156149a057fe5b146149f2576040805162461bcd60e51b815260206004820181905260248201527f73656e7420627920756e64657369676e61746564207472616e736d6974746572604482015290519081900360640190fd5b6000614a19633b9aca003a04836020015163ffffffff16846000015163ffffffff16614bb8565b90506010360260005a90506000614a388863ffffffff16858585614bde565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff16830281614a6657fe5b049050856080015163ffffffff16633b9aca0002816008896000015160ff16601f8110614a8f57fe5b015401016008886000015160ff16601f8110614aa757fe5b0155505050505050505050565b6003546001600160a01b039081169082168114610d4f57600380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b614b4b614d29565b60005b8351811015614bb0576000848281518110614b6557fe5b016020015160f81c9050614b8a8482601f8110614b7e57fe5b60200201516001614c6a565b848260ff16601f8110614b9957fe5b61ffff909216602092909202015250600101614b4e565b509092915050565b60008383811015614bcb57600285850304015b614bd58184614716565b95945050505050565b600081851015614c35576040805162461bcd60e51b815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b818503830161179301633b9aca00858202026fffffffffffffffffffffffffffffffff8110614c6057fe5b9695505050505050565b6000614c828261ffff168461ffff160161ffff614716565b9392505050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b60408051608081018252600080825260208201819052918101829052606081019190915290565b6040518060a00160405280614d08614cce565b81526060602082018190526040820181905280820152600060809091015290565b604051806103e00160405280601f906020820280368337509192915050565b600283019183908215614dce5791602002820160005b83821115614d9e57835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302614d5e565b8015614dcc5782816101000a81549061ffff0219169055600201602081600101049283019260010302614d9e565b505b50614dda929150614e0c565b5090565b82601f8101928215614dce579160200282015b82811115614dce578251825591602001919060010190614df1565b5b80821115614dda5760008155600101614e0d56fe6f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea164736f6c6343000705000a"

// DeployOffchainAggregator deploys a new Ethereum contract, binding an instance of OffchainAggregator to it.
func DeployOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _validator common.Address, _minAnswer *big.Int, _maxAnswer *big.Int, _billingAccessController common.Address, _requesterAccessController common.Address, _decimals uint8, _description string) (common.Address, *types.Transaction, *OffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OffchainAggregatorBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _validator, _minAnswer, _maxAnswer, _billingAccessController, _requesterAccessController, _decimals, _description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffchainAggregator{OffchainAggregatorCaller: OffchainAggregatorCaller{contract: contract}, OffchainAggregatorTransactor: OffchainAggregatorTransactor{contract: contract}, OffchainAggregatorFilterer: OffchainAggregatorFilterer{contract: contract}}, nil
}

// OffchainAggregator is an auto generated Go binding around an Ethereum contract.
type OffchainAggregator struct {
	OffchainAggregatorCaller     // Read-only binding to the contract
	OffchainAggregatorTransactor // Write-only binding to the contract
	OffchainAggregatorFilterer   // Log filterer for contract events
}

// OffchainAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OffchainAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OffchainAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OffchainAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OffchainAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OffchainAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OffchainAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OffchainAggregatorSession struct {
	Contract     *OffchainAggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OffchainAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OffchainAggregatorCallerSession struct {
	Contract *OffchainAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// OffchainAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OffchainAggregatorTransactorSession struct {
	Contract     *OffchainAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OffchainAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OffchainAggregatorRaw struct {
	Contract *OffchainAggregator // Generic contract binding to access the raw methods on
}

// OffchainAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OffchainAggregatorCallerRaw struct {
	Contract *OffchainAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// OffchainAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OffchainAggregatorTransactorRaw struct {
	Contract *OffchainAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOffchainAggregator creates a new instance of OffchainAggregator, bound to a specific deployed contract.
func NewOffchainAggregator(address common.Address, backend bind.ContractBackend) (*OffchainAggregator, error) {
	contract, err := bindOffchainAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregator{OffchainAggregatorCaller: OffchainAggregatorCaller{contract: contract}, OffchainAggregatorTransactor: OffchainAggregatorTransactor{contract: contract}, OffchainAggregatorFilterer: OffchainAggregatorFilterer{contract: contract}}, nil
}

// NewOffchainAggregatorCaller creates a new read-only instance of OffchainAggregator, bound to a specific deployed contract.
func NewOffchainAggregatorCaller(address common.Address, caller bind.ContractCaller) (*OffchainAggregatorCaller, error) {
	contract, err := bindOffchainAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorCaller{contract: contract}, nil
}

// NewOffchainAggregatorTransactor creates a new write-only instance of OffchainAggregator, bound to a specific deployed contract.
func NewOffchainAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OffchainAggregatorTransactor, error) {
	contract, err := bindOffchainAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorTransactor{contract: contract}, nil
}

// NewOffchainAggregatorFilterer creates a new log filterer instance of OffchainAggregator, bound to a specific deployed contract.
func NewOffchainAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OffchainAggregatorFilterer, error) {
	contract, err := bindOffchainAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorFilterer{contract: contract}, nil
}

// bindOffchainAggregator binds a generic wrapper to an already deployed contract.
func bindOffchainAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OffchainAggregator *OffchainAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregator.Contract.OffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OffchainAggregator *OffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.OffchainAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OffchainAggregator *OffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.OffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OffchainAggregator *OffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OffchainAggregator *OffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OffchainAggregator *OffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.contract.Transact(opts, method, params...)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_OffchainAggregator *OffchainAggregatorSession) LINK() (common.Address, error) {
	return _OffchainAggregator.Contract.LINK(&_OffchainAggregator.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCallerSession) LINK() (common.Address, error) {
	return _OffchainAggregator.Contract.LINK(&_OffchainAggregator.CallOpts)
}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCaller) BillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "billingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_OffchainAggregator *OffchainAggregatorSession) BillingAccessController() (common.Address, error) {
	return _OffchainAggregator.Contract.BillingAccessController(&_OffchainAggregator.CallOpts)
}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCallerSession) BillingAccessController() (common.Address, error) {
	return _OffchainAggregator.Contract.BillingAccessController(&_OffchainAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_OffchainAggregator *OffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_OffchainAggregator *OffchainAggregatorSession) Decimals() (uint8, error) {
	return _OffchainAggregator.Contract.Decimals(&_OffchainAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_OffchainAggregator *OffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _OffchainAggregator.Contract.Decimals(&_OffchainAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_OffchainAggregator *OffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_OffchainAggregator *OffchainAggregatorSession) Description() (string, error) {
	return _OffchainAggregator.Contract.Description(&_OffchainAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_OffchainAggregator *OffchainAggregatorCallerSession) Description() (string, error) {
	return _OffchainAggregator.Contract.Description(&_OffchainAggregator.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_OffchainAggregator *OffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_OffchainAggregator *OffchainAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetAnswer(&_OffchainAggregator.CallOpts, _roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 _roundId) view returns(int256)
func (_OffchainAggregator *OffchainAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetAnswer(&_OffchainAggregator.CallOpts, _roundId)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
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

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_OffchainAggregator *OffchainAggregatorSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _OffchainAggregator.Contract.GetBilling(&_OffchainAggregator.CallOpts)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_OffchainAggregator *OffchainAggregatorCallerSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _OffchainAggregator.Contract.GetBilling(&_OffchainAggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
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

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OffchainAggregator *OffchainAggregatorSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OffchainAggregator.Contract.GetRoundData(&_OffchainAggregator.CallOpts, _roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OffchainAggregator *OffchainAggregatorCallerSession) GetRoundData(_roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OffchainAggregator.Contract.GetRoundData(&_OffchainAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetTimestamp(&_OffchainAggregator.CallOpts, _roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 _roundId) view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _OffchainAggregator.Contract.GetTimestamp(&_OffchainAggregator.CallOpts, _roundId)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_OffchainAggregator *OffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_OffchainAggregator *OffchainAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestAnswer(&_OffchainAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestAnswer(&_OffchainAggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes16 configDigest)
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

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes16 configDigest)
func (_OffchainAggregator *OffchainAggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _OffchainAggregator.Contract.LatestConfigDetails(&_OffchainAggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes16 configDigest)
func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [16]byte
}, error) {
	return _OffchainAggregator.Contract.LatestConfigDetails(&_OffchainAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorSession) LatestRound() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestRound(&_OffchainAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestRound(&_OffchainAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
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

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OffchainAggregator *OffchainAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OffchainAggregator.Contract.LatestRoundData(&_OffchainAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _OffchainAggregator.Contract.LatestRoundData(&_OffchainAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestTimestamp(&_OffchainAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _OffchainAggregator.Contract.LatestTimestamp(&_OffchainAggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes16 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp)
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

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes16 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp)
func (_OffchainAggregator *OffchainAggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _OffchainAggregator.Contract.LatestTransmissionDetails(&_OffchainAggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes16 configDigest, uint32 epoch, uint8 round, int192 latestAnswer, uint64 latestTimestamp)
func (_OffchainAggregator *OffchainAggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [16]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _OffchainAggregator.Contract.LatestTransmissionDetails(&_OffchainAggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OffchainAggregator *OffchainAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OffchainAggregator *OffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregator.Contract.LinkAvailableForPayment(&_OffchainAggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OffchainAggregator *OffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregator.Contract.LinkAvailableForPayment(&_OffchainAggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_OffchainAggregator *OffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_OffchainAggregator *OffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MaxAnswer(&_OffchainAggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int192)
func (_OffchainAggregator *OffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MaxAnswer(&_OffchainAggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_OffchainAggregator *OffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_OffchainAggregator *OffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MinAnswer(&_OffchainAggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int192)
func (_OffchainAggregator *OffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _OffchainAggregator.Contract.MinAnswer(&_OffchainAggregator.CallOpts)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_OffchainAggregator *OffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_OffchainAggregator *OffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregator.Contract.OracleObservationCount(&_OffchainAggregator.CallOpts, _signerOrTransmitter)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_OffchainAggregator *OffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregator.Contract.OracleObservationCount(&_OffchainAggregator.CallOpts, _signerOrTransmitter)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregator.Contract.OwedPayment(&_OffchainAggregator.CallOpts, _transmitter)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregator.Contract.OwedPayment(&_OffchainAggregator.CallOpts, _transmitter)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OffchainAggregator *OffchainAggregatorSession) Owner() (common.Address, error) {
	return _OffchainAggregator.Contract.Owner(&_OffchainAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _OffchainAggregator.Contract.Owner(&_OffchainAggregator.CallOpts)
}

// RequesterAccessController is a free data retrieval call binding the contract method 0x70efdf2d.
//
// Solidity: function requesterAccessController() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCaller) RequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "requesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RequesterAccessController is a free data retrieval call binding the contract method 0x70efdf2d.
//
// Solidity: function requesterAccessController() view returns(address)
func (_OffchainAggregator *OffchainAggregatorSession) RequesterAccessController() (common.Address, error) {
	return _OffchainAggregator.Contract.RequesterAccessController(&_OffchainAggregator.CallOpts)
}

// RequesterAccessController is a free data retrieval call binding the contract method 0x70efdf2d.
//
// Solidity: function requesterAccessController() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCallerSession) RequesterAccessController() (common.Address, error) {
	return _OffchainAggregator.Contract.RequesterAccessController(&_OffchainAggregator.CallOpts)
}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OffchainAggregator *OffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OffchainAggregator *OffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _OffchainAggregator.Contract.Transmitters(&_OffchainAggregator.CallOpts)
}

// Transmitters is a free data retrieval call binding the contract method 0x81411834.
//
// Solidity: function transmitters() view returns(address[])
func (_OffchainAggregator *OffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _OffchainAggregator.Contract.Transmitters(&_OffchainAggregator.CallOpts)
}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "validator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_OffchainAggregator *OffchainAggregatorSession) Validator() (common.Address, error) {
	return _OffchainAggregator.Contract.Validator(&_OffchainAggregator.CallOpts)
}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_OffchainAggregator *OffchainAggregatorCallerSession) Validator() (common.Address, error) {
	return _OffchainAggregator.Contract.Validator(&_OffchainAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorSession) Version() (*big.Int, error) {
	return _OffchainAggregator.Contract.Version(&_OffchainAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_OffchainAggregator *OffchainAggregatorCallerSession) Version() (*big.Int, error) {
	return _OffchainAggregator.Contract.Version(&_OffchainAggregator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OffchainAggregator *OffchainAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffchainAggregator.Contract.AcceptOwnership(&_OffchainAggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffchainAggregator.Contract.AcceptOwnership(&_OffchainAggregator.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "acceptPayeeship", _transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_OffchainAggregator *OffchainAggregatorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.AcceptPayeeship(&_OffchainAggregator.TransactOpts, _transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.AcceptPayeeship(&_OffchainAggregator.TransactOpts, _transmitter)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_OffchainAggregator *OffchainAggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "requestNewRound")
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_OffchainAggregator *OffchainAggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _OffchainAggregator.Contract.RequestNewRound(&_OffchainAggregator.TransactOpts)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_OffchainAggregator *OffchainAggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _OffchainAggregator.Contract.RequestNewRound(&_OffchainAggregator.TransactOpts)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_OffchainAggregator *OffchainAggregatorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetBilling(&_OffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetBilling(&_OffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OffchainAggregator *OffchainAggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetBillingAccessController(&_OffchainAggregator.TransactOpts, _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetBillingAccessController(&_OffchainAggregator.TransactOpts, _billingAccessController)
}

// SetConfig is a paid mutator transaction binding the contract method 0x585aa7de.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _threshold, uint64 _encodedConfigVersion, bytes _encoded) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}

// SetConfig is a paid mutator transaction binding the contract method 0x585aa7de.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _threshold, uint64 _encodedConfigVersion, bytes _encoded) returns()
func (_OffchainAggregator *OffchainAggregatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetConfig(&_OffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}

// SetConfig is a paid mutator transaction binding the contract method 0x585aa7de.
//
// Solidity: function setConfig(address[] _signers, address[] _transmitters, uint8 _threshold, uint64 _encodedConfigVersion, bytes _encoded) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _encodedConfigVersion uint64, _encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetConfig(&_OffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _encodedConfigVersion, _encoded)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setPayees", _transmitters, _payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_OffchainAggregator *OffchainAggregatorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetPayees(&_OffchainAggregator.TransactOpts, _transmitters, _payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetPayees(&_OffchainAggregator.TransactOpts, _transmitters, _payees)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address _requesterAccessController) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, _requesterAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setRequesterAccessController", _requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address _requesterAccessController) returns()
func (_OffchainAggregator *OffchainAggregatorSession) SetRequesterAccessController(_requesterAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetRequesterAccessController(&_OffchainAggregator.TransactOpts, _requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address _requesterAccessController) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetRequesterAccessController(_requesterAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetRequesterAccessController(&_OffchainAggregator.TransactOpts, _requesterAccessController)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) SetValidator(opts *bind.TransactOpts, _newValidator common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "setValidator", _newValidator)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_OffchainAggregator *OffchainAggregatorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetValidator(&_OffchainAggregator.TransactOpts, _newValidator)
}

// SetValidator is a paid mutator transaction binding the contract method 0x1327d3d8.
//
// Solidity: function setValidator(address _newValidator) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) SetValidator(_newValidator common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.SetValidator(&_OffchainAggregator.TransactOpts, _newValidator)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_OffchainAggregator *OffchainAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.TransferOwnership(&_OffchainAggregator.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.TransferOwnership(&_OffchainAggregator.TransactOpts, _to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_OffchainAggregator *OffchainAggregatorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.TransferPayeeship(&_OffchainAggregator.TransactOpts, _transmitter, _proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.TransferPayeeship(&_OffchainAggregator.TransactOpts, _transmitter, _proposed)
}

// Transmit is a paid mutator transaction binding the contract method 0xc9807539.
//
// Solidity: function transmit(bytes _report, bytes32[] _rs, bytes32[] _ss, bytes32 _rawVs) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) Transmit(opts *bind.TransactOpts, _report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "transmit", _report, _rs, _ss, _rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xc9807539.
//
// Solidity: function transmit(bytes _report, bytes32[] _rs, bytes32[] _ss, bytes32 _rawVs) returns()
func (_OffchainAggregator *OffchainAggregatorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.Transmit(&_OffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xc9807539.
//
// Solidity: function transmit(bytes _report, bytes32[] _rs, bytes32[] _ss, bytes32 _rawVs) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) Transmit(_report []byte, _rs [][32]byte, _ss [][32]byte, _rawVs [32]byte) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.Transmit(&_OffchainAggregator.TransactOpts, _report, _rs, _ss, _rawVs)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_OffchainAggregator *OffchainAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.WithdrawFunds(&_OffchainAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.WithdrawFunds(&_OffchainAggregator.TransactOpts, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_OffchainAggregator *OffchainAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.contract.Transact(opts, "withdrawPayment", _transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_OffchainAggregator *OffchainAggregatorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.WithdrawPayment(&_OffchainAggregator.TransactOpts, _transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_OffchainAggregator *OffchainAggregatorTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregator.Contract.WithdrawPayment(&_OffchainAggregator.TransactOpts, _transmitter)
}

// OffchainAggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the OffchainAggregator contract.
type OffchainAggregatorAnswerUpdatedIterator struct {
	Event *OffchainAggregatorAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorAnswerUpdated represents a AnswerUpdated event raised by the OffchainAggregator contract.
type OffchainAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
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

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
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
				// New log arrived, parse the event and forward to the user
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*OffchainAggregatorAnswerUpdated, error) {
	event := new(OffchainAggregatorAnswerUpdated)
	if err := _OffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingAccessControllerSetIterator is returned from FilterBillingAccessControllerSet and is used to iterate over the raw logs and unpacked data for BillingAccessControllerSet events raised by the OffchainAggregator contract.
type OffchainAggregatorBillingAccessControllerSetIterator struct {
	Event *OffchainAggregatorBillingAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingAccessControllerSet represents a BillingAccessControllerSet event raised by the OffchainAggregator contract.
type OffchainAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBillingAccessControllerSet is a free log retrieval operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_OffchainAggregator *OffchainAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingAccessControllerSetIterator{contract: _OffchainAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchBillingAccessControllerSet is a free log subscription operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
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
				// New log arrived, parse the event and forward to the user
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

// ParseBillingAccessControllerSet is a log parse operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*OffchainAggregatorBillingAccessControllerSet, error) {
	event := new(OffchainAggregatorBillingAccessControllerSet)
	if err := _OffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingSetIterator is returned from FilterBillingSet and is used to iterate over the raw logs and unpacked data for BillingSet events raised by the OffchainAggregator contract.
type OffchainAggregatorBillingSetIterator struct {
	Event *OffchainAggregatorBillingSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingSet represents a BillingSet event raised by the OffchainAggregator contract.
type OffchainAggregatorBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterBillingSet is a free log retrieval operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_OffchainAggregator *OffchainAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*OffchainAggregatorBillingSetIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingSetIterator{contract: _OffchainAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

// WatchBillingSet is a free log subscription operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
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
				// New log arrived, parse the event and forward to the user
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

// ParseBillingSet is a log parse operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseBillingSet(log types.Log) (*OffchainAggregatorBillingSet, error) {
	event := new(OffchainAggregatorBillingSet)
	if err := _OffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the OffchainAggregator contract.
type OffchainAggregatorConfigSetIterator struct {
	Event *OffchainAggregatorConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorConfigSet represents a ConfigSet event raised by the OffchainAggregator contract.
type OffchainAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb9.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, uint64 configCount, address[] signers, address[] transmitters, uint8 threshold, uint64 encodedConfigVersion, bytes encoded)
func (_OffchainAggregator *OffchainAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OffchainAggregatorConfigSetIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorConfigSetIterator{contract: _OffchainAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb9.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, uint64 configCount, address[] signers, address[] transmitters, uint8 threshold, uint64 encodedConfigVersion, bytes encoded)
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
				// New log arrived, parse the event and forward to the user
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

// ParseConfigSet is a log parse operation binding the contract event 0x25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb9.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, uint64 configCount, address[] signers, address[] transmitters, uint8 threshold, uint64 encodedConfigVersion, bytes encoded)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseConfigSet(log types.Log) (*OffchainAggregatorConfigSet, error) {
	event := new(OffchainAggregatorConfigSet)
	if err := _OffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the OffchainAggregator contract.
type OffchainAggregatorNewRoundIterator struct {
	Event *OffchainAggregatorNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorNewRound represents a NewRound event raised by the OffchainAggregator contract.
type OffchainAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
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

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
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
				// New log arrived, parse the event and forward to the user
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseNewRound(log types.Log) (*OffchainAggregatorNewRound, error) {
	event := new(OffchainAggregatorNewRound)
	if err := _OffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorNewTransmissionIterator is returned from FilterNewTransmission and is used to iterate over the raw logs and unpacked data for NewTransmission events raised by the OffchainAggregator contract.
type OffchainAggregatorNewTransmissionIterator struct {
	Event *OffchainAggregatorNewTransmission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorNewTransmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorNewTransmission represents a NewTransmission event raised by the OffchainAggregator contract.
type OffchainAggregatorNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	RawReportContext  [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterNewTransmission is a free log retrieval operation binding the contract event 0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, int192[] observations, bytes observers, bytes32 rawReportContext)
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

// WatchNewTransmission is a free log subscription operation binding the contract event 0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, int192[] observations, bytes observers, bytes32 rawReportContext)
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
				// New log arrived, parse the event and forward to the user
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

// ParseNewTransmission is a log parse operation binding the contract event 0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, int192[] observations, bytes observers, bytes32 rawReportContext)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseNewTransmission(log types.Log) (*OffchainAggregatorNewTransmission, error) {
	event := new(OffchainAggregatorNewTransmission)
	if err := _OffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorOraclePaidIterator is returned from FilterOraclePaid and is used to iterate over the raw logs and unpacked data for OraclePaid events raised by the OffchainAggregator contract.
type OffchainAggregatorOraclePaidIterator struct {
	Event *OffchainAggregatorOraclePaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorOraclePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorOraclePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorOraclePaid represents a OraclePaid event raised by the OffchainAggregator contract.
type OffchainAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOraclePaid is a free log retrieval operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
func (_OffchainAggregator *OffchainAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts) (*OffchainAggregatorOraclePaidIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorOraclePaidIterator{contract: _OffchainAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

// WatchOraclePaid is a free log subscription operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOraclePaid is a log parse operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseOraclePaid(log types.Log) (*OffchainAggregatorOraclePaid, error) {
	event := new(OffchainAggregatorOraclePaid)
	if err := _OffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OffchainAggregator contract.
type OffchainAggregatorOwnershipTransferRequestedIterator struct {
	Event *OffchainAggregatorOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OffchainAggregator contract.
type OffchainAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffchainAggregatorOwnershipTransferRequested, error) {
	event := new(OffchainAggregatorOwnershipTransferRequested)
	if err := _OffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OffchainAggregator contract.
type OffchainAggregatorOwnershipTransferredIterator struct {
	Event *OffchainAggregatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the OffchainAggregator contract.
type OffchainAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*OffchainAggregatorOwnershipTransferred, error) {
	event := new(OffchainAggregatorOwnershipTransferred)
	if err := _OffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorPayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the OffchainAggregator contract.
type OffchainAggregatorPayeeshipTransferRequestedIterator struct {
	Event *OffchainAggregatorPayeeshipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorPayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorPayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the OffchainAggregator contract.
type OffchainAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
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

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
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
				// New log arrived, parse the event and forward to the user
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

// ParsePayeeshipTransferRequested is a log parse operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*OffchainAggregatorPayeeshipTransferRequested, error) {
	event := new(OffchainAggregatorPayeeshipTransferRequested)
	if err := _OffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorPayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the OffchainAggregator contract.
type OffchainAggregatorPayeeshipTransferredIterator struct {
	Event *OffchainAggregatorPayeeshipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorPayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorPayeeshipTransferred represents a PayeeshipTransferred event raised by the OffchainAggregator contract.
type OffchainAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
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

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
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
				// New log arrived, parse the event and forward to the user
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

// ParsePayeeshipTransferred is a log parse operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*OffchainAggregatorPayeeshipTransferred, error) {
	event := new(OffchainAggregatorPayeeshipTransferred)
	if err := _OffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorRequesterAccessControllerSetIterator is returned from FilterRequesterAccessControllerSet and is used to iterate over the raw logs and unpacked data for RequesterAccessControllerSet events raised by the OffchainAggregator contract.
type OffchainAggregatorRequesterAccessControllerSetIterator struct {
	Event *OffchainAggregatorRequesterAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorRequesterAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorRequesterAccessControllerSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorRequesterAccessControllerSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorRequesterAccessControllerSet represents a RequesterAccessControllerSet event raised by the OffchainAggregator contract.
type OffchainAggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRequesterAccessControllerSet is a free log retrieval operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_OffchainAggregator *OffchainAggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorRequesterAccessControllerSetIterator{contract: _OffchainAggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchRequesterAccessControllerSet is a free log subscription operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_OffchainAggregator *OffchainAggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OffchainAggregatorRequesterAccessControllerSet)
				if err := _OffchainAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
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

// ParseRequesterAccessControllerSet is a log parse operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*OffchainAggregatorRequesterAccessControllerSet, error) {
	event := new(OffchainAggregatorRequesterAccessControllerSet)
	if err := _OffchainAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorRoundRequestedIterator is returned from FilterRoundRequested and is used to iterate over the raw logs and unpacked data for RoundRequested events raised by the OffchainAggregator contract.
type OffchainAggregatorRoundRequestedIterator struct {
	Event *OffchainAggregatorRoundRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorRoundRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorRoundRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OffchainAggregatorRoundRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorRoundRequested represents a RoundRequested event raised by the OffchainAggregator contract.
type OffchainAggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [16]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRoundRequested is a free log retrieval operation binding the contract event 0x3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037.
//
// Solidity: event RoundRequested(address indexed requester, bytes16 configDigest, uint32 epoch, uint8 round)
func (_OffchainAggregator *OffchainAggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*OffchainAggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OffchainAggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorRoundRequestedIterator{contract: _OffchainAggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

// WatchRoundRequested is a free log subscription operation binding the contract event 0x3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037.
//
// Solidity: event RoundRequested(address indexed requester, bytes16 configDigest, uint32 epoch, uint8 round)
func (_OffchainAggregator *OffchainAggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OffchainAggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OffchainAggregatorRoundRequested)
				if err := _OffchainAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
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

// ParseRoundRequested is a log parse operation binding the contract event 0x3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037.
//
// Solidity: event RoundRequested(address indexed requester, bytes16 configDigest, uint32 epoch, uint8 round)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseRoundRequested(log types.Log) (*OffchainAggregatorRoundRequested, error) {
	event := new(OffchainAggregatorRoundRequested)
	if err := _OffchainAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorValidatorUpdatedIterator is returned from FilterValidatorUpdated and is used to iterate over the raw logs and unpacked data for ValidatorUpdated events raised by the OffchainAggregator contract.
type OffchainAggregatorValidatorUpdatedIterator struct {
	Event *OffchainAggregatorValidatorUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorValidatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorValidatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorValidatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorValidatorUpdated represents a ValidatorUpdated event raised by the OffchainAggregator contract.
type OffchainAggregatorValidatorUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterValidatorUpdated is a free log retrieval operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
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

// WatchValidatorUpdated is a free log subscription operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
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
				// New log arrived, parse the event and forward to the user
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

// ParseValidatorUpdated is a log parse operation binding the contract event 0xcfac5dc75b8d9a7e074162f59d9adcd33da59f0fe8dfb21580db298fc0fdad0d.
//
// Solidity: event ValidatorUpdated(address indexed previous, address indexed current)
func (_OffchainAggregator *OffchainAggregatorFilterer) ParseValidatorUpdated(log types.Log) (*OffchainAggregatorValidatorUpdated, error) {
	event := new(OffchainAggregatorValidatorUpdated)
	if err := _OffchainAggregator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingABI is the input ABI used to generate the binding from.
const OffchainAggregatorBillingABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"billingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OffchainAggregatorBillingBin is the compiled bytecode used for deploying new contracts.
var OffchainAggregatorBillingBin = "0x60a06040523480156200001157600080fd5b506040516200271d3803806200271d833981810160405260e08110156200003757600080fd5b508051602082015160408301516060840151608085015160a086015160c090960151600080546001600160a01b03191633179055949593949293919290919062000085878787878762000136565b620000908162000228565b6001600160601b0319606083901b16608052620000ac620002a1565b620000b6620002a1565b60005b601f8160ff16101562000106576001838260ff16601f8110620000d857fe5b61ffff909216602092909202015260018260ff8316601f8110620000f857fe5b6020020152600101620000b9565b5062000116600483601f620002c0565b5062000126600882601f6200035d565b50505050505050505050620003a5565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526002805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6003546001600160a01b0390811690821681146200029d57600380546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b604051806103e00160405280601f906020820280368337509192915050565b6002830191839082156200034b5791602002820160005b838211156200031957835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302620002d7565b8015620003495782816101000a81549061ffff021916905560020160208160010104928301926001030262000319565b505b50620003599291506200038e565b5090565b82601f81019282156200034b579160200282015b828111156200034b57825182559160200191906001019062000371565b5b808211156200035957600081556001016200038f565b60805160601c612341620003dc600039806105de5280610fc152806110e5528061122152806117de5280611bc452506123416000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063b121e14711610097578063e4902f8211610066578063e4902f8214610384578063eb5dcd6c146103ce578063f2fde38b14610409578063fbffd2c11461043c57610100565b8063b121e147146102cb578063bd824706146102fe578063c107532914610343578063d09dc3391461037c57610100565b80638ac28d5a116100d35780638ac28d5a146101c65780638da5cb5b146101f9578063996e8298146102015780639c849b301461020957610100565b80630eafb25b146101055780631b6b6d231461014a578063299372681461017b57806379ba5097146101bc575b600080fd5b6101386004803603602081101561011b57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661046f565b60408051918252519081900360200190f35b6101526105dc565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b610183610600565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b6101c4610695565b005b6101c4600480360360208110156101dc57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610797565b610152610838565b610152610854565b6101c46004803603604081101561021f57600080fd5b81019060208101813564010000000081111561023a57600080fd5b82018360208201111561024c57600080fd5b8035906020019184602083028401116401000000008311171561026e57600080fd5b91939092909160208101903564010000000081111561028c57600080fd5b82018360208201111561029e57600080fd5b803590602001918460208302840111640100000000831117156102c057600080fd5b509092509050610870565b6101c4600480360360208110156102e157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610b6c565b6101c4600480360360a081101561031457600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516610c99565b6101c46004803603604081101561035957600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610e33565b61013861121c565b6103b76004803603602081101561039a57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166112e7565b6040805161ffff9092168252519081900360200190f35b6101c4600480360360408110156103e457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160200135166113ad565b6101c46004803603602081101561041f57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611571565b6101c46004803603602081101561045257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661166d565b60006104796121f7565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602760209081526040918290208251808401909352805460ff808216855291928401916101009091041660028111156104ca57fe5b60028111156104d557fe5b90525090506000816020015160028111156104ec57fe5b14156104fc5760009150506105d7565b61050461220e565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830181905270010000000000000000000000000000000090910490921660808201528251909160009160019060049060ff16601f811061059057fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca0002905060016008846000015160ff16601f81106105ce57fe5b01540301925050505b919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600080600080600061061061220e565b50506040805160a08101825260025463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529890975092955093509150565b60015473ffffffffffffffffffffffffffffffffffffffff16331461071b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526006602052604090205416331461082c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b610835816116fc565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60035473ffffffffffffffffffffffffffffffffffffffff1690565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108f657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b82811461096457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015610b6557600085858381811061097d57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff16905060008484848181106109aa57fe5b73ffffffffffffffffffffffffffffffffffffffff8581166000908152600660209081526040909120549202939093013583169350909116905080158080610a1d57508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b610a8857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260066020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614610b55578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506109679050565b5050505050565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260076020526040902054163314610c0157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526006602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556007909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60035460005473ffffffffffffffffffffffffffffffffffffffff9182169116331480610dab5750604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610d7e57600080fd5b505afa158015610d92573d6000803e3d6000fd5b505050506040513d6020811015610da857600080fd5b50515b610e1657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b610e1e61196d565b610e2b8686868686611dcb565b505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331480610f465750600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015610f1957600080fd5b505afa158015610f2d573d6000803e3d6000fd5b505050506040513d6020811015610f4357600080fd5b50515b610fb157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6000610fbb611f45565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561104657600080fd5b505afa15801561105a573d6000803e3d6000fd5b505050506040513d602081101561107057600080fd5b50519050818110156110e357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8561112c85850387612133565b6040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b15801561117f57600080fd5b505af1158015611193573d6000803e3d6000fd5b505050506040513d60208110156111a957600080fd5b505161121657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156112a657600080fd5b505afa1580156112ba573d6000803e3d6000fd5b505050506040513d60208110156112d057600080fd5b5051905060006112de611f45565b90910391505090565b60006112f16121f7565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561134257fe5b600281111561134d57fe5b905250905060008160200151600281111561136457fe5b14156113745760009150506105d7565b60016004826000015160ff16601f811061138a57fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526006602052604090205416331461144257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b3373ffffffffffffffffffffffffffffffffffffffff821614156114c757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260076020526040902080548383167fffffffffffffffffffffffff00000000000000000000000000000000000000008216811790925590911690811461156c5760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146115f757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005473ffffffffffffffffffffffffffffffffffffffff1633146116f357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6108358161214d565b6117046121f7565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602760209081526040918290208251808401909352805460ff8082168552919284019161010090910416600281111561175557fe5b600281111561176057fe5b905250905060006117708361046f565b9050801561156c5773ffffffffffffffffffffffffffffffffffffffff80841660009081526006602090815260408083205481517fa9059cbb0000000000000000000000000000000000000000000000000000000081529085166004820181905260248201879052915191947f0000000000000000000000000000000000000000000000000000000000000000169363a9059cbb9360448084019491939192918390030190829087803b15801561182657600080fd5b505af115801561183a573d6000803e3d6000fd5b505050506040513d602081101561185057600080fd5b50516118bd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60016004846000015160ff16601f81106118d357fe5b601091828204019190066002026101000a81548161ffff021916908361ffff16021790555060016008846000015160ff16601f811061190e57fe5b01556040805173ffffffffffffffffffffffffffffffffffffffff80871682528316602082015280820184905290517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b29181900360600190a150505050565b61197561220e565b506040805160a08101825260025463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c01000000000000000000000000810483166060830152700100000000000000000000000000000000900490911660808201526119ec61223c565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411611a0557905050505050509050611a4c61223c565b604080516103e081019182905290600890601f9082845b815481526020019060010190808311611a63575050505050905060606029805480602002602001604051908101604052809291908181526020018280548015611ae257602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611ab7575b5050505050905060005b8151811015611daf57600060018483601f8110611b0557fe5b6020020151039050600060018684601f8110611b1d57fe5b60200201510361ffff169050600082886060015163ffffffff168302633b9aca00020190506000811115611da457600060066000878781518110611b5d57fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015611c5357600080fd5b505af1158015611c67573d6000803e3d6000fd5b505050506040513d6020811015611c7d57600080fd5b5051611cea57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018886601f8110611cf857fe5b61ffff909216602092909202015260018786601f8110611d1457fe5b602002015285517fe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b290879087908110611d4957fe5b60200260200101518284604051808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a1505b505050600101611aec565b50611dbd600484601f61225b565b50610b65600883601f6122f1565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b6000611f4f61223c565b604080516103e081019182905290600490601f90826000855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411611f685790505050505050905060005b601f811015611fd85760018282601f8110611fc157fe5b60200201510361ffff169290920191600101611faa565b50611fe161220e565b506040805160a08101825260025463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c0100000000000000000000000083048216606080860182905270010000000000000000000000000000000090940490921660808501526029805486518184028101840190975280875297909202633b9aca00029693949293908301828280156120be57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612093575b505050505090506120cd61223c565b604080516103e081019182905290600890601f9082845b8154815260200190600101908083116120e4575050505050905060005b825181101561212b5760018282601f811061211857fe5b6020020151039590950194600101612101565b505050505090565b600081831015612144575081612147565b50805b92915050565b60035473ffffffffffffffffffffffffffffffffffffffff90811690821681146121f357600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b604080518082019091526000808252602082015290565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b604051806103e00160405280601f906020820280368337509192915050565b6002830191839082156122e15791602002820160005b838211156122b157835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302612271565b80156122df5782816101000a81549061ffff02191690556002016020816001010492830192600103026122b1565b505b506122ed92915061231f565b5090565b82601f81019282156122e1579160200282015b828111156122e1578251825591602001919060010190612304565b5b808211156122ed576000815560010161232056fea164736f6c6343000705000a"

// DeployOffchainAggregatorBilling deploys a new Ethereum contract, binding an instance of OffchainAggregatorBilling to it.
func DeployOffchainAggregatorBilling(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _billingAccessController common.Address) (common.Address, *types.Transaction, *OffchainAggregatorBilling, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorBillingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OffchainAggregatorBillingBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _billingAccessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffchainAggregatorBilling{OffchainAggregatorBillingCaller: OffchainAggregatorBillingCaller{contract: contract}, OffchainAggregatorBillingTransactor: OffchainAggregatorBillingTransactor{contract: contract}, OffchainAggregatorBillingFilterer: OffchainAggregatorBillingFilterer{contract: contract}}, nil
}

// OffchainAggregatorBilling is an auto generated Go binding around an Ethereum contract.
type OffchainAggregatorBilling struct {
	OffchainAggregatorBillingCaller     // Read-only binding to the contract
	OffchainAggregatorBillingTransactor // Write-only binding to the contract
	OffchainAggregatorBillingFilterer   // Log filterer for contract events
}

// OffchainAggregatorBillingCaller is an auto generated read-only Go binding around an Ethereum contract.
type OffchainAggregatorBillingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OffchainAggregatorBillingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OffchainAggregatorBillingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OffchainAggregatorBillingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OffchainAggregatorBillingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OffchainAggregatorBillingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OffchainAggregatorBillingSession struct {
	Contract     *OffchainAggregatorBilling // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// OffchainAggregatorBillingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OffchainAggregatorBillingCallerSession struct {
	Contract *OffchainAggregatorBillingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// OffchainAggregatorBillingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OffchainAggregatorBillingTransactorSession struct {
	Contract     *OffchainAggregatorBillingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// OffchainAggregatorBillingRaw is an auto generated low-level Go binding around an Ethereum contract.
type OffchainAggregatorBillingRaw struct {
	Contract *OffchainAggregatorBilling // Generic contract binding to access the raw methods on
}

// OffchainAggregatorBillingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OffchainAggregatorBillingCallerRaw struct {
	Contract *OffchainAggregatorBillingCaller // Generic read-only contract binding to access the raw methods on
}

// OffchainAggregatorBillingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OffchainAggregatorBillingTransactorRaw struct {
	Contract *OffchainAggregatorBillingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOffchainAggregatorBilling creates a new instance of OffchainAggregatorBilling, bound to a specific deployed contract.
func NewOffchainAggregatorBilling(address common.Address, backend bind.ContractBackend) (*OffchainAggregatorBilling, error) {
	contract, err := bindOffchainAggregatorBilling(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBilling{OffchainAggregatorBillingCaller: OffchainAggregatorBillingCaller{contract: contract}, OffchainAggregatorBillingTransactor: OffchainAggregatorBillingTransactor{contract: contract}, OffchainAggregatorBillingFilterer: OffchainAggregatorBillingFilterer{contract: contract}}, nil
}

// NewOffchainAggregatorBillingCaller creates a new read-only instance of OffchainAggregatorBilling, bound to a specific deployed contract.
func NewOffchainAggregatorBillingCaller(address common.Address, caller bind.ContractCaller) (*OffchainAggregatorBillingCaller, error) {
	contract, err := bindOffchainAggregatorBilling(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingCaller{contract: contract}, nil
}

// NewOffchainAggregatorBillingTransactor creates a new write-only instance of OffchainAggregatorBilling, bound to a specific deployed contract.
func NewOffchainAggregatorBillingTransactor(address common.Address, transactor bind.ContractTransactor) (*OffchainAggregatorBillingTransactor, error) {
	contract, err := bindOffchainAggregatorBilling(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingTransactor{contract: contract}, nil
}

// NewOffchainAggregatorBillingFilterer creates a new log filterer instance of OffchainAggregatorBilling, bound to a specific deployed contract.
func NewOffchainAggregatorBillingFilterer(address common.Address, filterer bind.ContractFilterer) (*OffchainAggregatorBillingFilterer, error) {
	contract, err := bindOffchainAggregatorBilling(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingFilterer{contract: contract}, nil
}

// bindOffchainAggregatorBilling binds a generic wrapper to an already deployed contract.
func bindOffchainAggregatorBilling(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OffchainAggregatorBillingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OffchainAggregatorBilling *OffchainAggregatorBillingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.OffchainAggregatorBillingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregatorBilling.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.contract.Transact(opts, method, params...)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) LINK() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.LINK(&_OffchainAggregatorBilling.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) LINK() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.LINK(&_OffchainAggregatorBilling.CallOpts)
}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) BillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "billingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) BillingAccessController() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.BillingAccessController(&_OffchainAggregatorBilling.CallOpts)
}

// BillingAccessController is a free data retrieval call binding the contract method 0x996e8298.
//
// Solidity: function billingAccessController() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) BillingAccessController() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.BillingAccessController(&_OffchainAggregatorBilling.CallOpts)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
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

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _OffchainAggregatorBilling.Contract.GetBilling(&_OffchainAggregatorBilling.CallOpts)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) GetBilling() (struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}, error) {
	return _OffchainAggregatorBilling.Contract.GetBilling(&_OffchainAggregatorBilling.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.LinkAvailableForPayment(&_OffchainAggregatorBilling.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.LinkAvailableForPayment(&_OffchainAggregatorBilling.CallOpts)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregatorBilling.Contract.OracleObservationCount(&_OffchainAggregatorBilling.CallOpts, _signerOrTransmitter)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address _signerOrTransmitter) view returns(uint16)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _OffchainAggregatorBilling.Contract.OracleObservationCount(&_OffchainAggregatorBilling.CallOpts, _signerOrTransmitter)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.OwedPayment(&_OffchainAggregatorBilling.CallOpts, _transmitter)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address _transmitter) view returns(uint256)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _OffchainAggregatorBilling.Contract.OwedPayment(&_OffchainAggregatorBilling.CallOpts, _transmitter)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffchainAggregatorBilling.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) Owner() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.Owner(&_OffchainAggregatorBilling.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingCallerSession) Owner() (common.Address, error) {
	return _OffchainAggregatorBilling.Contract.Owner(&_OffchainAggregatorBilling.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.AcceptOwnership(&_OffchainAggregatorBilling.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.AcceptOwnership(&_OffchainAggregatorBilling.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "acceptPayeeship", _transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.AcceptPayeeship(&_OffchainAggregatorBilling.TransactOpts, _transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address _transmitter) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.AcceptPayeeship(&_OffchainAggregatorBilling.TransactOpts, _transmitter)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetBilling(&_OffchainAggregatorBilling.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBilling is a paid mutator transaction binding the contract method 0xbd824706.
//
// Solidity: function setBilling(uint32 _maximumGasPrice, uint32 _reasonableGasPrice, uint32 _microLinkPerEth, uint32 _linkGweiPerObservation, uint32 _linkGweiPerTransmission) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetBilling(&_OffchainAggregatorBilling.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetBillingAccessController(&_OffchainAggregatorBilling.TransactOpts, _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetBillingAccessController(&_OffchainAggregatorBilling.TransactOpts, _billingAccessController)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "setPayees", _transmitters, _payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetPayees(&_OffchainAggregatorBilling.TransactOpts, _transmitters, _payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] _transmitters, address[] _payees) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.SetPayees(&_OffchainAggregatorBilling.TransactOpts, _transmitters, _payees)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.TransferOwnership(&_OffchainAggregatorBilling.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.TransferOwnership(&_OffchainAggregatorBilling.TransactOpts, _to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.TransferPayeeship(&_OffchainAggregatorBilling.TransactOpts, _transmitter, _proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address _transmitter, address _proposed) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.TransferPayeeship(&_OffchainAggregatorBilling.TransactOpts, _transmitter, _proposed)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.WithdrawFunds(&_OffchainAggregatorBilling.TransactOpts, _recipient, _amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address _recipient, uint256 _amount) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.WithdrawFunds(&_OffchainAggregatorBilling.TransactOpts, _recipient, _amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.contract.Transact(opts, "withdrawPayment", _transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.WithdrawPayment(&_OffchainAggregatorBilling.TransactOpts, _transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address _transmitter) returns()
func (_OffchainAggregatorBilling *OffchainAggregatorBillingTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorBilling.Contract.WithdrawPayment(&_OffchainAggregatorBilling.TransactOpts, _transmitter)
}

// OffchainAggregatorBillingBillingAccessControllerSetIterator is returned from FilterBillingAccessControllerSet and is used to iterate over the raw logs and unpacked data for BillingAccessControllerSet events raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingBillingAccessControllerSetIterator struct {
	Event *OffchainAggregatorBillingBillingAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingBillingAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingBillingAccessControllerSet represents a BillingAccessControllerSet event raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBillingAccessControllerSet is a free log retrieval operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorBillingBillingAccessControllerSetIterator, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingBillingAccessControllerSetIterator{contract: _OffchainAggregatorBilling.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchBillingAccessControllerSet is a free log subscription operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
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
				// New log arrived, parse the event and forward to the user
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

// ParseBillingAccessControllerSet is a log parse operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseBillingAccessControllerSet(log types.Log) (*OffchainAggregatorBillingBillingAccessControllerSet, error) {
	event := new(OffchainAggregatorBillingBillingAccessControllerSet)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingBillingSetIterator is returned from FilterBillingSet and is used to iterate over the raw logs and unpacked data for BillingSet events raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingBillingSetIterator struct {
	Event *OffchainAggregatorBillingBillingSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingBillingSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingBillingSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingBillingSet represents a BillingSet event raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterBillingSet is a free log retrieval operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterBillingSet(opts *bind.FilterOpts) (*OffchainAggregatorBillingBillingSetIterator, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingBillingSetIterator{contract: _OffchainAggregatorBilling.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

// WatchBillingSet is a free log subscription operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
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
				// New log arrived, parse the event and forward to the user
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

// ParseBillingSet is a log parse operation binding the contract event 0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6.
//
// Solidity: event BillingSet(uint32 maximumGasPrice, uint32 reasonableGasPrice, uint32 microLinkPerEth, uint32 linkGweiPerObservation, uint32 linkGweiPerTransmission)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseBillingSet(log types.Log) (*OffchainAggregatorBillingBillingSet, error) {
	event := new(OffchainAggregatorBillingBillingSet)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingOraclePaidIterator is returned from FilterOraclePaid and is used to iterate over the raw logs and unpacked data for OraclePaid events raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingOraclePaidIterator struct {
	Event *OffchainAggregatorBillingOraclePaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingOraclePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingOraclePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingOraclePaid represents a OraclePaid event raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOraclePaid is a free log retrieval operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) FilterOraclePaid(opts *bind.FilterOpts) (*OffchainAggregatorBillingOraclePaidIterator, error) {

	logs, sub, err := _OffchainAggregatorBilling.contract.FilterLogs(opts, "OraclePaid")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorBillingOraclePaidIterator{contract: _OffchainAggregatorBilling.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

// WatchOraclePaid is a free log subscription operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOraclePaid is a log parse operation binding the contract event 0xe8ec50e5150ae28ae37e493ff389ffab7ffaec2dc4dccfca03f12a3de29d12b2.
//
// Solidity: event OraclePaid(address transmitter, address payee, uint256 amount)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseOraclePaid(log types.Log) (*OffchainAggregatorBillingOraclePaid, error) {
	event := new(OffchainAggregatorBillingOraclePaid)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingOwnershipTransferRequestedIterator struct {
	Event *OffchainAggregatorBillingOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffchainAggregatorBillingOwnershipTransferRequested, error) {
	event := new(OffchainAggregatorBillingOwnershipTransferRequested)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingOwnershipTransferredIterator struct {
	Event *OffchainAggregatorBillingOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingOwnershipTransferred represents a OwnershipTransferred event raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParseOwnershipTransferred(log types.Log) (*OffchainAggregatorBillingOwnershipTransferred, error) {
	event := new(OffchainAggregatorBillingOwnershipTransferred)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingPayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingPayeeshipTransferRequestedIterator struct {
	Event *OffchainAggregatorBillingPayeeshipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingPayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingPayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
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

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
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
				// New log arrived, parse the event and forward to the user
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

// ParsePayeeshipTransferRequested is a log parse operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParsePayeeshipTransferRequested(log types.Log) (*OffchainAggregatorBillingPayeeshipTransferRequested, error) {
	event := new(OffchainAggregatorBillingPayeeshipTransferRequested)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OffchainAggregatorBillingPayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingPayeeshipTransferredIterator struct {
	Event *OffchainAggregatorBillingPayeeshipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OffchainAggregatorBillingPayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OffchainAggregatorBillingPayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OffchainAggregatorBillingPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OffchainAggregatorBillingPayeeshipTransferred represents a PayeeshipTransferred event raised by the OffchainAggregatorBilling contract.
type OffchainAggregatorBillingPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
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

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
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
				// New log arrived, parse the event and forward to the user
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

// ParsePayeeshipTransferred is a log parse operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_OffchainAggregatorBilling *OffchainAggregatorBillingFilterer) ParsePayeeshipTransferred(log types.Log) (*OffchainAggregatorBillingPayeeshipTransferred, error) {
	event := new(OffchainAggregatorBillingPayeeshipTransferred)
	if err := _OffchainAggregatorBilling.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnedABI is the input ABI used to generate the binding from.
const OwnedABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OwnedBin is the compiled bytecode used for deploying new contracts.
var OwnedBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556102db806100326000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806379ba5097146100465780638da5cb5b14610050578063f2fde38b14610081575b600080fd5b61004e6100b4565b005b6100586101b6565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61004e6004803603602081101561009757600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166101d2565b60015473ffffffffffffffffffffffffffffffffffffffff16331461013a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461025857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a35056fea164736f6c6343000705000a"

// DeployOwned deploys a new Ethereum contract, binding an instance of Owned to it.
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

// Owned is an auto generated Go binding around an Ethereum contract.
type Owned struct {
	OwnedCaller     // Read-only binding to the contract
	OwnedTransactor // Write-only binding to the contract
	OwnedFilterer   // Log filterer for contract events
}

// OwnedCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedSession struct {
	Contract     *Owned            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedCallerSession struct {
	Contract *OwnedCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OwnedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedTransactorSession struct {
	Contract     *OwnedTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedRaw struct {
	Contract *Owned // Generic contract binding to access the raw methods on
}

// OwnedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedCallerRaw struct {
	Contract *OwnedCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedTransactorRaw struct {
	Contract *OwnedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwned creates a new instance of Owned, bound to a specific deployed contract.
func NewOwned(address common.Address, backend bind.ContractBackend) (*Owned, error) {
	contract, err := bindOwned(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// NewOwnedCaller creates a new read-only instance of Owned, bound to a specific deployed contract.
func NewOwnedCaller(address common.Address, caller bind.ContractCaller) (*OwnedCaller, error) {
	contract, err := bindOwned(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedCaller{contract: contract}, nil
}

// NewOwnedTransactor creates a new write-only instance of Owned, bound to a specific deployed contract.
func NewOwnedTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTransactor, error) {
	contract, err := bindOwned(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTransactor{contract: contract}, nil
}

// NewOwnedFilterer creates a new log filterer instance of Owned, bound to a specific deployed contract.
func NewOwnedFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedFilterer, error) {
	contract, err := bindOwned(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedFilterer{contract: contract}, nil
}

// bindOwned binds a generic wrapper to an already deployed contract.
func bindOwned(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Owned.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Owned *OwnedSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Owned *OwnedCallerSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Owned *OwnedTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Owned.Contract.AcceptOwnership(&_Owned.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Owned *OwnedTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Owned *OwnedSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Owned.Contract.TransferOwnership(&_Owned.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Owned *OwnedTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Owned.Contract.TransferOwnership(&_Owned.TransactOpts, _to)
}

// OwnedOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Owned contract.
type OwnedOwnershipTransferRequestedIterator struct {
	Event *OwnedOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnedOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnedOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Owned contract.
type OwnedOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Owned *OwnedFilterer) ParseOwnershipTransferRequested(log types.Log) (*OwnedOwnershipTransferRequested, error) {
	event := new(OwnedOwnershipTransferRequested)
	if err := _Owned.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnedOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Owned contract.
type OwnedOwnershipTransferredIterator struct {
	Event *OwnedOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnedOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnedOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedOwnershipTransferred represents a OwnershipTransferred event raised by the Owned contract.
type OwnedOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Owned *OwnedFilterer) ParseOwnershipTransferred(log types.Log) (*OwnedOwnershipTransferred, error) {
	event := new(OwnedOwnershipTransferred)
	if err := _Owned.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleReadAccessControllerABI is the input ABI used to generate the binding from.
const SimpleReadAccessControllerABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SimpleReadAccessControllerBin is the compiled bytecode used for deploying new contracts.
var SimpleReadAccessControllerBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556001805460ff60a01b1916600160a01b1790556109c4806100456000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638823da6c11610076578063a118f2491161005b578063a118f249146101fd578063dc7f012414610230578063f2fde38b14610238576100a3565b80638823da6c146101995780638da5cb5b146101cc576100a3565b80630a756983146100a85780636b14daf8146100b257806379ba5097146101895780638038e4a114610191575b600080fd5b6100b061026b565b005b610175600480360360408110156100c857600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561010057600080fd5b82018360208201111561011257600080fd5b8035906020019184600183028401116401000000008311171561013457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610368945050505050565b604080519115158252519081900360200190f35b6100b061039b565b6100b061049d565b6100b0600480360360208110156101af57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166105af565b6101d46106e7565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100b06004803603602081101561021357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610703565b610175610792565b6100b06004803603602081101561024e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166107b3565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff161561036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b600061037483836108af565b80610394575073ffffffffffffffffffffffffffffffffffffffff831632145b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461042157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461052357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff1661036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60005473ffffffffffffffffffffffffffffffffffffffff16331461063557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff16156106e45773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461078957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6106e481610904565b60015474010000000000000000000000000000000000000000900460ff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461083957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff168061039457505060015474010000000000000000000000000000000000000000900460ff161592915050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff166106e45773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a15056fea164736f6c6343000705000a"

// DeploySimpleReadAccessController deploys a new Ethereum contract, binding an instance of SimpleReadAccessController to it.
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

// SimpleReadAccessController is an auto generated Go binding around an Ethereum contract.
type SimpleReadAccessController struct {
	SimpleReadAccessControllerCaller     // Read-only binding to the contract
	SimpleReadAccessControllerTransactor // Write-only binding to the contract
	SimpleReadAccessControllerFilterer   // Log filterer for contract events
}

// SimpleReadAccessControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleReadAccessControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleReadAccessControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleReadAccessControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleReadAccessControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleReadAccessControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleReadAccessControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleReadAccessControllerSession struct {
	Contract     *SimpleReadAccessController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SimpleReadAccessControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleReadAccessControllerCallerSession struct {
	Contract *SimpleReadAccessControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// SimpleReadAccessControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleReadAccessControllerTransactorSession struct {
	Contract     *SimpleReadAccessControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// SimpleReadAccessControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleReadAccessControllerRaw struct {
	Contract *SimpleReadAccessController // Generic contract binding to access the raw methods on
}

// SimpleReadAccessControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleReadAccessControllerCallerRaw struct {
	Contract *SimpleReadAccessControllerCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleReadAccessControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleReadAccessControllerTransactorRaw struct {
	Contract *SimpleReadAccessControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleReadAccessController creates a new instance of SimpleReadAccessController, bound to a specific deployed contract.
func NewSimpleReadAccessController(address common.Address, backend bind.ContractBackend) (*SimpleReadAccessController, error) {
	contract, err := bindSimpleReadAccessController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessController{SimpleReadAccessControllerCaller: SimpleReadAccessControllerCaller{contract: contract}, SimpleReadAccessControllerTransactor: SimpleReadAccessControllerTransactor{contract: contract}, SimpleReadAccessControllerFilterer: SimpleReadAccessControllerFilterer{contract: contract}}, nil
}

// NewSimpleReadAccessControllerCaller creates a new read-only instance of SimpleReadAccessController, bound to a specific deployed contract.
func NewSimpleReadAccessControllerCaller(address common.Address, caller bind.ContractCaller) (*SimpleReadAccessControllerCaller, error) {
	contract, err := bindSimpleReadAccessController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCaller{contract: contract}, nil
}

// NewSimpleReadAccessControllerTransactor creates a new write-only instance of SimpleReadAccessController, bound to a specific deployed contract.
func NewSimpleReadAccessControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleReadAccessControllerTransactor, error) {
	contract, err := bindSimpleReadAccessController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerTransactor{contract: contract}, nil
}

// NewSimpleReadAccessControllerFilterer creates a new log filterer instance of SimpleReadAccessController, bound to a specific deployed contract.
func NewSimpleReadAccessControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleReadAccessControllerFilterer, error) {
	contract, err := bindSimpleReadAccessController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerFilterer{contract: contract}, nil
}

// bindSimpleReadAccessController binds a generic wrapper to an already deployed contract.
func bindSimpleReadAccessController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleReadAccessControllerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleReadAccessController *SimpleReadAccessControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.SimpleReadAccessControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleReadAccessController *SimpleReadAccessControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleReadAccessController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.contract.Transact(opts, method, params...)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) CheckEnabled() (bool, error) {
	return _SimpleReadAccessController.Contract.CheckEnabled(&_SimpleReadAccessController.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) CheckEnabled() (bool, error) {
	return _SimpleReadAccessController.Contract.CheckEnabled(&_SimpleReadAccessController.CallOpts)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _SimpleReadAccessController.Contract.HasAccess(&_SimpleReadAccessController.CallOpts, _user, _calldata)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _SimpleReadAccessController.Contract.HasAccess(&_SimpleReadAccessController.CallOpts, _user, _calldata)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleReadAccessController *SimpleReadAccessControllerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleReadAccessController.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) Owner() (common.Address, error) {
	return _SimpleReadAccessController.Contract.Owner(&_SimpleReadAccessController.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleReadAccessController *SimpleReadAccessControllerCallerSession) Owner() (common.Address, error) {
	return _SimpleReadAccessController.Contract.Owner(&_SimpleReadAccessController.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AcceptOwnership(&_SimpleReadAccessController.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AcceptOwnership(&_SimpleReadAccessController.TransactOpts)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "addAccess", _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AddAccess(&_SimpleReadAccessController.TransactOpts, _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.AddAccess(&_SimpleReadAccessController.TransactOpts, _user)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "disableAccessCheck")
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.DisableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.DisableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "enableAccessCheck")
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.EnableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.EnableAccessCheck(&_SimpleReadAccessController.TransactOpts)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "removeAccess", _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.RemoveAccess(&_SimpleReadAccessController.TransactOpts, _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.RemoveAccess(&_SimpleReadAccessController.TransactOpts, _user)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.TransferOwnership(&_SimpleReadAccessController.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_SimpleReadAccessController *SimpleReadAccessControllerTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleReadAccessController.Contract.TransferOwnership(&_SimpleReadAccessController.TransactOpts, _to)
}

// SimpleReadAccessControllerAddedAccessIterator is returned from FilterAddedAccess and is used to iterate over the raw logs and unpacked data for AddedAccess events raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerAddedAccessIterator struct {
	Event *SimpleReadAccessControllerAddedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleReadAccessControllerAddedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleReadAccessControllerAddedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleReadAccessControllerAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleReadAccessControllerAddedAccess represents a AddedAccess event raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerAddedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedAccess is a free log retrieval operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*SimpleReadAccessControllerAddedAccessIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerAddedAccessIterator{contract: _SimpleReadAccessController.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

// WatchAddedAccess is a free log subscription operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
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
				// New log arrived, parse the event and forward to the user
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

// ParseAddedAccess is a log parse operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseAddedAccess(log types.Log) (*SimpleReadAccessControllerAddedAccess, error) {
	event := new(SimpleReadAccessControllerAddedAccess)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleReadAccessControllerCheckAccessDisabledIterator is returned from FilterCheckAccessDisabled and is used to iterate over the raw logs and unpacked data for CheckAccessDisabled events raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerCheckAccessDisabledIterator struct {
	Event *SimpleReadAccessControllerCheckAccessDisabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleReadAccessControllerCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleReadAccessControllerCheckAccessDisabled represents a CheckAccessDisabled event raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerCheckAccessDisabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessDisabled is a free log retrieval operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*SimpleReadAccessControllerCheckAccessDisabledIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCheckAccessDisabledIterator{contract: _SimpleReadAccessController.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessDisabled is a free log subscription operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
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
				// New log arrived, parse the event and forward to the user
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

// ParseCheckAccessDisabled is a log parse operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseCheckAccessDisabled(log types.Log) (*SimpleReadAccessControllerCheckAccessDisabled, error) {
	event := new(SimpleReadAccessControllerCheckAccessDisabled)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleReadAccessControllerCheckAccessEnabledIterator is returned from FilterCheckAccessEnabled and is used to iterate over the raw logs and unpacked data for CheckAccessEnabled events raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerCheckAccessEnabledIterator struct {
	Event *SimpleReadAccessControllerCheckAccessEnabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleReadAccessControllerCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleReadAccessControllerCheckAccessEnabled represents a CheckAccessEnabled event raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerCheckAccessEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessEnabled is a free log retrieval operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*SimpleReadAccessControllerCheckAccessEnabledIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerCheckAccessEnabledIterator{contract: _SimpleReadAccessController.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessEnabled is a free log subscription operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
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
				// New log arrived, parse the event and forward to the user
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

// ParseCheckAccessEnabled is a log parse operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseCheckAccessEnabled(log types.Log) (*SimpleReadAccessControllerCheckAccessEnabled, error) {
	event := new(SimpleReadAccessControllerCheckAccessEnabled)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleReadAccessControllerOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerOwnershipTransferRequestedIterator struct {
	Event *SimpleReadAccessControllerOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleReadAccessControllerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleReadAccessControllerOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseOwnershipTransferRequested(log types.Log) (*SimpleReadAccessControllerOwnershipTransferRequested, error) {
	event := new(SimpleReadAccessControllerOwnershipTransferRequested)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleReadAccessControllerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerOwnershipTransferredIterator struct {
	Event *SimpleReadAccessControllerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleReadAccessControllerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleReadAccessControllerOwnershipTransferred represents a OwnershipTransferred event raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseOwnershipTransferred(log types.Log) (*SimpleReadAccessControllerOwnershipTransferred, error) {
	event := new(SimpleReadAccessControllerOwnershipTransferred)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleReadAccessControllerRemovedAccessIterator is returned from FilterRemovedAccess and is used to iterate over the raw logs and unpacked data for RemovedAccess events raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerRemovedAccessIterator struct {
	Event *SimpleReadAccessControllerRemovedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleReadAccessControllerRemovedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleReadAccessControllerRemovedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleReadAccessControllerRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleReadAccessControllerRemovedAccess represents a RemovedAccess event raised by the SimpleReadAccessController contract.
type SimpleReadAccessControllerRemovedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedAccess is a free log retrieval operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*SimpleReadAccessControllerRemovedAccessIterator, error) {

	logs, sub, err := _SimpleReadAccessController.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleReadAccessControllerRemovedAccessIterator{contract: _SimpleReadAccessController.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

// WatchRemovedAccess is a free log subscription operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
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
				// New log arrived, parse the event and forward to the user
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

// ParseRemovedAccess is a log parse operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_SimpleReadAccessController *SimpleReadAccessControllerFilterer) ParseRemovedAccess(log types.Log) (*SimpleReadAccessControllerRemovedAccess, error) {
	event := new(SimpleReadAccessControllerRemovedAccess)
	if err := _SimpleReadAccessController.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleWriteAccessControllerABI is the input ABI used to generate the binding from.
const SimpleWriteAccessControllerABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SimpleWriteAccessControllerBin is the compiled bytecode used for deploying new contracts.
var SimpleWriteAccessControllerBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556001805460ff60a01b1916600160a01b179055610992806100456000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638823da6c11610076578063a118f2491161005b578063a118f249146101fd578063dc7f012414610230578063f2fde38b14610238576100a3565b80638823da6c146101995780638da5cb5b146101cc576100a3565b80630a756983146100a85780636b14daf8146100b257806379ba5097146101895780638038e4a114610191575b600080fd5b6100b061026b565b005b610175600480360360408110156100c857600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561010057600080fd5b82018360208201111561011257600080fd5b8035906020019184600183028401116401000000008311171561013457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610368945050505050565b604080519115158252519081900360200190f35b6100b06103be565b6100b06104c0565b6100b0600480360360208110156101af57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166105d2565b6101d461070a565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100b06004803603602081101561021357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610726565b6101756107b5565b6100b06004803603602081101561024e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166107d6565b60005473ffffffffffffffffffffffffffffffffffffffff1633146102f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff161561036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff16806103b7575060015474010000000000000000000000000000000000000000900460ff16155b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461044457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff16331461054657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff1661036657600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60005473ffffffffffffffffffffffffffffffffffffffff16331461065857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff16156107075773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146107ac57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b610707816108d2565b60015474010000000000000000000000000000000000000000900460ff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461085c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff166107075773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a15056fea164736f6c6343000705000a"

// DeploySimpleWriteAccessController deploys a new Ethereum contract, binding an instance of SimpleWriteAccessController to it.
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

// SimpleWriteAccessController is an auto generated Go binding around an Ethereum contract.
type SimpleWriteAccessController struct {
	SimpleWriteAccessControllerCaller     // Read-only binding to the contract
	SimpleWriteAccessControllerTransactor // Write-only binding to the contract
	SimpleWriteAccessControllerFilterer   // Log filterer for contract events
}

// SimpleWriteAccessControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleWriteAccessControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleWriteAccessControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleWriteAccessControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleWriteAccessControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleWriteAccessControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleWriteAccessControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleWriteAccessControllerSession struct {
	Contract     *SimpleWriteAccessController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// SimpleWriteAccessControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleWriteAccessControllerCallerSession struct {
	Contract *SimpleWriteAccessControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// SimpleWriteAccessControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleWriteAccessControllerTransactorSession struct {
	Contract     *SimpleWriteAccessControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// SimpleWriteAccessControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleWriteAccessControllerRaw struct {
	Contract *SimpleWriteAccessController // Generic contract binding to access the raw methods on
}

// SimpleWriteAccessControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleWriteAccessControllerCallerRaw struct {
	Contract *SimpleWriteAccessControllerCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleWriteAccessControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleWriteAccessControllerTransactorRaw struct {
	Contract *SimpleWriteAccessControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleWriteAccessController creates a new instance of SimpleWriteAccessController, bound to a specific deployed contract.
func NewSimpleWriteAccessController(address common.Address, backend bind.ContractBackend) (*SimpleWriteAccessController, error) {
	contract, err := bindSimpleWriteAccessController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessController{SimpleWriteAccessControllerCaller: SimpleWriteAccessControllerCaller{contract: contract}, SimpleWriteAccessControllerTransactor: SimpleWriteAccessControllerTransactor{contract: contract}, SimpleWriteAccessControllerFilterer: SimpleWriteAccessControllerFilterer{contract: contract}}, nil
}

// NewSimpleWriteAccessControllerCaller creates a new read-only instance of SimpleWriteAccessController, bound to a specific deployed contract.
func NewSimpleWriteAccessControllerCaller(address common.Address, caller bind.ContractCaller) (*SimpleWriteAccessControllerCaller, error) {
	contract, err := bindSimpleWriteAccessController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerCaller{contract: contract}, nil
}

// NewSimpleWriteAccessControllerTransactor creates a new write-only instance of SimpleWriteAccessController, bound to a specific deployed contract.
func NewSimpleWriteAccessControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleWriteAccessControllerTransactor, error) {
	contract, err := bindSimpleWriteAccessController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerTransactor{contract: contract}, nil
}

// NewSimpleWriteAccessControllerFilterer creates a new log filterer instance of SimpleWriteAccessController, bound to a specific deployed contract.
func NewSimpleWriteAccessControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleWriteAccessControllerFilterer, error) {
	contract, err := bindSimpleWriteAccessController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerFilterer{contract: contract}, nil
}

// bindSimpleWriteAccessController binds a generic wrapper to an already deployed contract.
func bindSimpleWriteAccessController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleWriteAccessControllerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleWriteAccessController *SimpleWriteAccessControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.SimpleWriteAccessControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleWriteAccessController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.contract.Transact(opts, method, params...)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SimpleWriteAccessController.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) CheckEnabled() (bool, error) {
	return _SimpleWriteAccessController.Contract.CheckEnabled(&_SimpleWriteAccessController.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerSession) CheckEnabled() (bool, error) {
	return _SimpleWriteAccessController.Contract.CheckEnabled(&_SimpleWriteAccessController.CallOpts)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes ) view returns(bool)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) HasAccess(opts *bind.CallOpts, _user common.Address, arg1 []byte) (bool, error) {
	var out []interface{}
	err := _SimpleWriteAccessController.contract.Call(opts, &out, "hasAccess", _user, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes ) view returns(bool)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) HasAccess(_user common.Address, arg1 []byte) (bool, error) {
	return _SimpleWriteAccessController.Contract.HasAccess(&_SimpleWriteAccessController.CallOpts, _user, arg1)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes ) view returns(bool)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerSession) HasAccess(_user common.Address, arg1 []byte) (bool, error) {
	return _SimpleWriteAccessController.Contract.HasAccess(&_SimpleWriteAccessController.CallOpts, _user, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleWriteAccessController.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) Owner() (common.Address, error) {
	return _SimpleWriteAccessController.Contract.Owner(&_SimpleWriteAccessController.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerCallerSession) Owner() (common.Address, error) {
	return _SimpleWriteAccessController.Contract.Owner(&_SimpleWriteAccessController.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.AcceptOwnership(&_SimpleWriteAccessController.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.AcceptOwnership(&_SimpleWriteAccessController.TransactOpts)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "addAccess", _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.AddAccess(&_SimpleWriteAccessController.TransactOpts, _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.AddAccess(&_SimpleWriteAccessController.TransactOpts, _user)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "disableAccessCheck")
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.DisableAccessCheck(&_SimpleWriteAccessController.TransactOpts)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.DisableAccessCheck(&_SimpleWriteAccessController.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "enableAccessCheck")
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.EnableAccessCheck(&_SimpleWriteAccessController.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.EnableAccessCheck(&_SimpleWriteAccessController.TransactOpts)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "removeAccess", _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.RemoveAccess(&_SimpleWriteAccessController.TransactOpts, _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.RemoveAccess(&_SimpleWriteAccessController.TransactOpts, _user)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.TransferOwnership(&_SimpleWriteAccessController.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _SimpleWriteAccessController.Contract.TransferOwnership(&_SimpleWriteAccessController.TransactOpts, _to)
}

// SimpleWriteAccessControllerAddedAccessIterator is returned from FilterAddedAccess and is used to iterate over the raw logs and unpacked data for AddedAccess events raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerAddedAccessIterator struct {
	Event *SimpleWriteAccessControllerAddedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleWriteAccessControllerAddedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleWriteAccessControllerAddedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleWriteAccessControllerAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleWriteAccessControllerAddedAccess represents a AddedAccess event raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerAddedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedAccess is a free log retrieval operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*SimpleWriteAccessControllerAddedAccessIterator, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerAddedAccessIterator{contract: _SimpleWriteAccessController.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

// WatchAddedAccess is a free log subscription operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
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
				// New log arrived, parse the event and forward to the user
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

// ParseAddedAccess is a log parse operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseAddedAccess(log types.Log) (*SimpleWriteAccessControllerAddedAccess, error) {
	event := new(SimpleWriteAccessControllerAddedAccess)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleWriteAccessControllerCheckAccessDisabledIterator is returned from FilterCheckAccessDisabled and is used to iterate over the raw logs and unpacked data for CheckAccessDisabled events raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerCheckAccessDisabledIterator struct {
	Event *SimpleWriteAccessControllerCheckAccessDisabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleWriteAccessControllerCheckAccessDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleWriteAccessControllerCheckAccessDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleWriteAccessControllerCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleWriteAccessControllerCheckAccessDisabled represents a CheckAccessDisabled event raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerCheckAccessDisabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessDisabled is a free log retrieval operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*SimpleWriteAccessControllerCheckAccessDisabledIterator, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerCheckAccessDisabledIterator{contract: _SimpleWriteAccessController.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessDisabled is a free log subscription operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
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
				// New log arrived, parse the event and forward to the user
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

// ParseCheckAccessDisabled is a log parse operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseCheckAccessDisabled(log types.Log) (*SimpleWriteAccessControllerCheckAccessDisabled, error) {
	event := new(SimpleWriteAccessControllerCheckAccessDisabled)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleWriteAccessControllerCheckAccessEnabledIterator is returned from FilterCheckAccessEnabled and is used to iterate over the raw logs and unpacked data for CheckAccessEnabled events raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerCheckAccessEnabledIterator struct {
	Event *SimpleWriteAccessControllerCheckAccessEnabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleWriteAccessControllerCheckAccessEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleWriteAccessControllerCheckAccessEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleWriteAccessControllerCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleWriteAccessControllerCheckAccessEnabled represents a CheckAccessEnabled event raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerCheckAccessEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessEnabled is a free log retrieval operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*SimpleWriteAccessControllerCheckAccessEnabledIterator, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerCheckAccessEnabledIterator{contract: _SimpleWriteAccessController.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessEnabled is a free log subscription operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
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
				// New log arrived, parse the event and forward to the user
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

// ParseCheckAccessEnabled is a log parse operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseCheckAccessEnabled(log types.Log) (*SimpleWriteAccessControllerCheckAccessEnabled, error) {
	event := new(SimpleWriteAccessControllerCheckAccessEnabled)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleWriteAccessControllerOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerOwnershipTransferRequestedIterator struct {
	Event *SimpleWriteAccessControllerOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleWriteAccessControllerOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleWriteAccessControllerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleWriteAccessControllerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleWriteAccessControllerOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseOwnershipTransferRequested(log types.Log) (*SimpleWriteAccessControllerOwnershipTransferRequested, error) {
	event := new(SimpleWriteAccessControllerOwnershipTransferRequested)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleWriteAccessControllerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerOwnershipTransferredIterator struct {
	Event *SimpleWriteAccessControllerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleWriteAccessControllerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleWriteAccessControllerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleWriteAccessControllerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleWriteAccessControllerOwnershipTransferred represents a OwnershipTransferred event raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
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
				// New log arrived, parse the event and forward to the user
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseOwnershipTransferred(log types.Log) (*SimpleWriteAccessControllerOwnershipTransferred, error) {
	event := new(SimpleWriteAccessControllerOwnershipTransferred)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleWriteAccessControllerRemovedAccessIterator is returned from FilterRemovedAccess and is used to iterate over the raw logs and unpacked data for RemovedAccess events raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerRemovedAccessIterator struct {
	Event *SimpleWriteAccessControllerRemovedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleWriteAccessControllerRemovedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleWriteAccessControllerRemovedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleWriteAccessControllerRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleWriteAccessControllerRemovedAccess represents a RemovedAccess event raised by the SimpleWriteAccessController contract.
type SimpleWriteAccessControllerRemovedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedAccess is a free log retrieval operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*SimpleWriteAccessControllerRemovedAccessIterator, error) {

	logs, sub, err := _SimpleWriteAccessController.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &SimpleWriteAccessControllerRemovedAccessIterator{contract: _SimpleWriteAccessController.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

// WatchRemovedAccess is a free log subscription operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
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
				// New log arrived, parse the event and forward to the user
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

// ParseRemovedAccess is a log parse operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_SimpleWriteAccessController *SimpleWriteAccessControllerFilterer) ParseRemovedAccess(log types.Log) (*SimpleWriteAccessControllerRemovedAccess, error) {
	event := new(SimpleWriteAccessControllerRemovedAccess)
	if err := _SimpleWriteAccessController.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
