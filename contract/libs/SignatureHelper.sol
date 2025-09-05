// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library SignatureHelper {

    function checkCertSignature(bytes memory signature, bytes memory pubKey, bytes memory data) internal pure returns (bytes32 dataHash) {
        address signatoryAddress = address(uint160(uint256(keccak256(pubKey))));
        if (signature.length == 64) {
            bytes32 r;
            bytes32 s;
            /// @solidity memory-safe-assembly
            assembly {
                r := mload(add(signature, 0x20))
                s := mload(add(signature, 0x40))
            }
            // Do not check 's' by EIP-2 because
            // certificates can generate 's' greater than secp256k1n ÷ 2 + 1
            return checkSignature(r, s, signatoryAddress, data);
        } else {
            revert("Invalid signature length");
        }
    }

    function checkSignature(bytes32 r, bytes32 s, address signatoryAddress, bytes memory data) internal pure returns (bytes32 dataHash) {
        dataHash = sha256(data);

        if (
            signatoryAddress != ecrecover(dataHash, 27, r, s) &&
            signatoryAddress != ecrecover(dataHash, 28, r, s) &&
            signatoryAddress != ecrecover(dataHash, 29, r, s) &&
            signatoryAddress != ecrecover(dataHash, 30, r, s)
        ) {
            revert("Invalid signature");
        }
    }

}
