// Package gethwrappers keeps track of the golang wrappers of the solidity contracts
package gethwrappers

//go:generate ./compile.sh 15000 ../../contract2/src/OCR2TitleRequest.sol

//go:generate ./compile.sh 15000 ../../contract2/src/OCR2Aggregator.sol
//go:generate ./compile.sh 15000 ../../contract2/src/AccessControlledOCR2Aggregator.sol
//go:generate ./compile.sh 15000 ../../contract2/src/ExposedOCR2Aggregator.sol

//go:generate ./compile.sh 1000 ../../contract2/src/TestOCR2Aggregator.sol
//go:generate ./compile.sh 1000 ../../contract2/src/TestValidator.sol
//go:generate ./compile.sh 1000 ../../contract2/src/AccessControlTestHelper.sol
