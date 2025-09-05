// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library DateConverter {
    uint256 constant UNIX_YEAR = 1970;
    uint256 constant FIRST_LEAP_YEAR_BEFORE_UNIX = 1968;
    uint256 constant FIRST_NON_LEAP_CENTURY_YEAR_BEFORE_UNIX = 1900;
    uint256 constant FIRST_EXCLUSIVE_LEAP_YEAR_BEFORE_UNIX = 1600;

    uint256 constant PASSED_SECONDS_BEFORE_JANUARY = 0 days;
    uint256 constant PASSED_SECONDS_BEFORE_FEBRUARY = 31 days;
    uint256 constant PASSED_SECONDS_BEFORE_MARCH = 59 days;
    uint256 constant PASSED_SECONDS_BEFORE_APRIL = 90 days;
    uint256 constant PASSED_SECONDS_BEFORE_MAY = 120 days;
    uint256 constant PASSED_SECONDS_BEFORE_JUNE = 151 days;
    uint256 constant PASSED_SECONDS_BEFORE_JULY = 181 days;
    uint256 constant PASSED_SECONDS_BEFORE_AUGUST = 212 days;
    uint256 constant PASSED_SECONDS_BEFORE_SEPTEMBER = 243 days;
    uint256 constant PASSED_SECONDS_BEFORE_OCTOBER = 273 days;
    uint256 constant PASSED_SECONDS_BEFORE_NOVEMBER = 304 days;
    uint256 constant PASSED_SECONDS_BEFORE_DECEMBER = 334 days;

    uint256 constant MONTH_NOT_FOUND = 1;
    uint256 constant DIGIT_NOT_FOUND = 10;

    function isLeapYear(uint256 year) internal pure returns (bool) {
        return (year % 4 == 0 && year % 100 != 0) || (year % 400 == 0);
    }

    function getDaysInMonth(uint256 month, uint256 year) internal pure returns (uint256) {
        if (month == 2) {
            return isLeapYear(year) ? 29 : 28;
        } else if (month == 4 || month == 6 || month == 9 || month == 11) {
            return 30;
        } else {
            return 31;
        }
    }

    function convertDateToUnixTimestamp(string memory date) internal pure returns (uint256 timestamp) {
        return convertDateToUnixTimestamp(bytes(date));
    }

    function convertDateToUnixTimestamp(bytes memory date) internal pure returns (uint256 timestamp) {
        uint256 dateBytesLength = date.length;
        uint256 startIndex;
        uint256 year;
        if (dateBytesLength == 13) {
            startIndex = 2;
            year = _stringToUint(date, 0, 2);
            year += year < 50 ? 2000 : 1900;
        } else if (dateBytesLength == 19) {
            startIndex = 4;
            year = _stringToUint(date, 0, 4);
        } else {
            revert("Only 13 and 19 bytes are supported date formats");
        }
        uint256 month = _stringToUint(date, startIndex, startIndex + 2);
        uint256 day = _stringToUint(date, startIndex + 2, startIndex + 4);
        uint256 hour = _stringToUint(date, startIndex + 4, startIndex + 6);
        uint256 minute = _stringToUint(date, startIndex + 6, startIndex + 8);
        uint256 second = _stringToUint(date, startIndex + 8, startIndex + 10);

        require(day >= 1 && day <= getDaysInMonth(month, year), "Invalid day");
        require(hour < 24, "Invalid hours");
        require(minute < 60, "Invalid minutes");
        require(second < 60, "Invalid seconds");

        timestamp = _getYearPassedSeconds(year);
        timestamp += _getMonthPassedSeconds(month, year);
        timestamp += (day - 1) * 1 days + hour * 1 hours + minute * 1 minutes + second;
    }

    function _stringToUint(bytes memory b, uint256 startIndex, uint256 endIndex) internal pure returns (uint256 result) {
        for (uint256 i = startIndex; i < endIndex; i++) {
            result = result * 10 + _decodeSymbolToDigit(bytes(b)[i]);
        }
    }

    function _getYearPassedSeconds(uint256 year) internal pure returns (uint256 passedSeconds) {
        require(year >= UNIX_YEAR, "Year must be >= 1970");

        // (year - 1) is used because the number of leap years that have already passed is calculated.
        uint256 nonLeapCenturyCount = (year - 1 - FIRST_NON_LEAP_CENTURY_YEAR_BEFORE_UNIX) / 100;
        uint256 exclusiveLeapYears = (year - 1 - FIRST_EXCLUSIVE_LEAP_YEAR_BEFORE_UNIX) / 400;
        uint256 countOfLeapYears = (year - 1 - FIRST_LEAP_YEAR_BEFORE_UNIX) / 4 - nonLeapCenturyCount + exclusiveLeapYears;
        passedSeconds = (year - UNIX_YEAR) * 365 days + countOfLeapYears * 1 days;
    }

    function _getMonthPassedSeconds(uint256 month, uint256 year) internal pure returns (uint256 passedSeconds) {
        assembly {
            switch month
            case 1 {
                passedSeconds := PASSED_SECONDS_BEFORE_JANUARY
            }
            case 2 {
                passedSeconds := PASSED_SECONDS_BEFORE_FEBRUARY
            }
            case 3 {
                passedSeconds := PASSED_SECONDS_BEFORE_MARCH
            }
            case 4 {
                passedSeconds := PASSED_SECONDS_BEFORE_APRIL
            }
            case 5 {
                passedSeconds := PASSED_SECONDS_BEFORE_MAY
            }
            case 6 {
                passedSeconds := PASSED_SECONDS_BEFORE_JUNE
            }
            case 7 {
                passedSeconds := PASSED_SECONDS_BEFORE_JULY
            }
            case 8 {
                passedSeconds := PASSED_SECONDS_BEFORE_AUGUST
            }
            case 9 {
                passedSeconds := PASSED_SECONDS_BEFORE_SEPTEMBER
            }
            case 10 {
                passedSeconds := PASSED_SECONDS_BEFORE_OCTOBER
            }
            case 11 {
                passedSeconds := PASSED_SECONDS_BEFORE_NOVEMBER
            }
            case 12 {
                passedSeconds := PASSED_SECONDS_BEFORE_DECEMBER
            }
            default {
                passedSeconds := MONTH_NOT_FOUND
            }
        }
        if (passedSeconds == MONTH_NOT_FOUND) {
            revert("Invalid month");
        }
        if (month > 2 && isLeapYear(year)) {
            passedSeconds += 1 days;
        }
    }

    function _decodeSymbolToDigit(bytes1 symbol) internal pure returns (uint256 digit) {
        assembly {
            switch symbol
            case "0" {
                digit := 0
            }
            case "1" {
                digit := 1
            }
            case "2" {
                digit := 2
            }
            case "3" {
                digit := 3
            }
            case "4" {
                digit := 4
            }
            case "5" {
                digit := 5
            }
            case "6" {
                digit := 6
            }
            case "7" {
                digit := 7
            }
            case "8" {
                digit := 8
            }
            case "9" {
                digit := 9
            }
            default {
                digit := DIGIT_NOT_FOUND
            }
        }
        if (digit == DIGIT_NOT_FOUND) {
            revert("ASCII symbol is not a number");
        }
    }
}
