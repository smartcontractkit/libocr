// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/HashSets.sol";

library Setn {
    function add(Bytes32HashSet storage set, bytes32 item) internal {
        bool isNewUniqElement = set.items.length == 0 || (set.indexes[item] == 0 && set.items[0] != item);
        if (!isNewUniqElement) {
            return;
        }

        set.indexes[item] = set.items.length;
        set.items.push(item);
    }

    function isExists(Bytes32HashSet storage set, bytes32 item) internal view returns (bool) {
        if (set.items.length == 0) {
            return false;
        }
        return set.indexes[item] != 0 || set.items[0] == item;
    }

    function remove(Bytes32HashSet storage set, bytes32 item) internal {
        if (isExists(set, item)) {
            uint256 idx = set.indexes[item];
            require(set.items[idx] == item, "500. Unexpected Bytes32HashSet state");

            bytes32 last = set.items[set.items.length - 1];
            set.items[idx] = last;
            set.indexes[last] = idx;
            set.indexes[item] = 0;
            set.items.pop();
        }
    }

    function clear(Bytes32HashSet storage set) internal {
        uint256 idx = 0;
        for (; idx < set.items.length; idx++) {
            bytes32 item = set.items[idx];
            delete set.indexes[item];
        }
        delete set.items;
    }

    function count(BytesHashSet storage set) internal view returns (uint256) {
        return set.items.length;
    }

    function isEmpty(BytesHashSet storage set) internal view returns (bool) {
        return set.items.length == 0;
    }

    function add(BytesHashSet storage set, bytes memory item) internal {
        bool isNewUniqElement = set.items.length == 0 || (set.indexes[item] == 0 && keccak256(set.items[0]) != keccak256(item));
        if (!isNewUniqElement) {
            return;
        }

        set.indexes[item] = set.items.length;
        set.items.push(item);
    }

    function isExists(BytesHashSet storage set, bytes memory item) internal view returns (bool) {
        if (set.items.length == 0) {
            return false;
        }
        return set.indexes[item] != 0 || keccak256(set.items[0]) == keccak256(item);
    }

    function remove(BytesHashSet storage set, bytes memory item) internal {
        if (isExists(set, item)) {
            uint256 idx = set.indexes[item];
            require(keccak256(set.items[idx]) == keccak256(item), "500. Unexpected UintHashSet state");

            bytes memory last = set.items[set.items.length - 1];
            set.items[idx] = last;
            set.indexes[last] = idx;
            set.indexes[item] = 0;
            set.items.pop();
        }
    }

    function clear(BytesHashSet storage set) internal {
        uint256 itemsLength = set.items.length;
        for (uint256 idx; idx < itemsLength; idx++) {
            bytes memory item = set.items[idx];
            delete set.indexes[item];
        }
        delete set.items;
    }
}
