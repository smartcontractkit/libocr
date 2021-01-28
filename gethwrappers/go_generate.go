// Package gethwrappers keeps track of the golang wrappers of the solidity contracts
package gethwrappers

//go:generate ./compile.sh 15000 ../../contract/src/AccessControlledOffchainAggregator.sol
//go:generate ./compile.sh 15000 ../../contract/src/OffchainAggregator.sol
//go:generate ./compile.sh 15000 ../../contract/src/ExposedOffchainAggregator.sol

//go:generate ./compile.sh 1000 ../../contract/src/TestOffchainAggregator.sol
//go:generate ./compile.sh 1000 ../../contract/src/TestValidator.sol
//go:generate ./compile.sh 1000 ../../contract/src/AccessControlTestHelper.sol
