// SPDX-License-Identifier: MIT
pragma solidity ^0.7.6;

import "./HashSets.sol";

struct BytesBlacklist {
    bytes item;
    bool isBlacklisted;
}

struct Bytes32Blacklist {
    bytes32 item;
    bool isBlacklisted;
}

struct BlackList {
    BytesHashSet bySerialNumber;
    Bytes32HashSet byMrEnclave;
    Bytes32HashSet byMrSigner;
}
