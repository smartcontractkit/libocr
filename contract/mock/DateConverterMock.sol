// SPDX-License-Identifier: MIT
pragma solidity ^0.7.6;
import "../libs/DateConverter.sol";

contract DateConverterMock {
    function convertDateToUnixTimestamp(string calldata date) external pure returns (uint256) {
        return DateConverter.convertDateToUnixTimestamp(date);
    }

    function isLeapYear(uint256 year) external pure returns (bool) {
        return DateConverter.isLeapYear(year);
    }

    function getDaysInMonth(uint256 month, uint256 year) external pure returns (uint256) {
        return DateConverter.getDaysInMonth(month, year);
    }
}
