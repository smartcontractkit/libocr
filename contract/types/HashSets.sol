// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

struct Bytes32HashSet {
    mapping(bytes32 => uint256) indexes;
    bytes32[] items;
}

struct BytesHashSet {
    mapping(bytes => uint256) indexes;
    bytes[] items;
}
