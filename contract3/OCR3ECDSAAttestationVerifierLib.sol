// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./OCR3AttestationVerifierErrors.sol";

/// @title Internal library for attestation verification using ECDSA signatures
/// @dev Developers should inherit from OCR3eCDSAAttestationVerifier for this functionality to be compiled with the
///      application contract or see OCR3DynamicallyDispatchedAttestationVerifier for using it as a pre-deployed
///      library.
library OCR3ECDSAAttestationVerifierLib {
    /// @notice Verifies and stores the provided ECDSA public keys for later use within `verifyAttestation(...)`.
    /// @param s_keys Storage space for holding all ECDSA verification keys (in the form of addresses). The first `n`
    ///        values are populated with values from the `keys` parameter. The size of 32 is the maximum number of keys
    ///        supported (based on the width of the attribution bitmask). Storage costs are only payed for the number of
    ///        keys used `n`.
    /// @param n The number of keys expected to be set by this call. Must match the actual number of keys present in
    ///        the `keys` parameter.
    /// @param keys A concatenation of `n` ECDSA public keys (i.e., addresses, 20 bytes each).
    function setVerificationKeys(uint256[32] storage s_keys, uint8 n, bytes calldata keys) internal {
        // Verify that `n` is consistent with the amount of data being passed in the `keys` parameter.
        // The maximum of 32 keys is based on the width of the attribution bitmask (currently set to 32 bits).
        if (keys.length % 20 != 0) {
            revert KeysOfInvalidSize();
        }
        if (keys.length / 20 != n) {
            revert InvalidNumberOfKeys();
        }
        if (n > 32) {
            revert MaximumNumberOfKeysExceeded();
        }

        // Copy the provided keys from calldata (keys) to storage (s_keys). After copying, the i-th 32 byte storage slot
        // holds the i-th 20 byte key (i.e., the signer's address) in its lower bytes.
        uint256 pos = 0;
        for (uint256 i = 0; i < n; ++i) {
            // Read the next 20 byte key/address from the keys parameter and ensure the key/address is non-zero.
            // Clearly the value 0x0000000000000000000000000000000000000000 is an invalid key/address, however, during
            // signature verification the call to the ecRecover precompile returns zero on failure, and its return
            // value is directly compared to s_keys[i], which therefore must never be zero.
            uint160 key = uint160(bytes20(keys[pos:pos + 20]));
            if (key == 0) {
                revert InvalidKey();
            }
            s_keys[i] = key;
            pos += 20;
        }
    }

    /// @notice Verifies the attestation for the given report hash. Reverts on verification failure.
    /// @param s_keys A list (stored in the application contract) holding the ECDSA public keys (i.e., addresses) of all
    ///               oracles.
    /// @param n The total number of oracles. Used to verify the attribution bitmask as part of the attestation data.
    /// @param f Maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly.
    ///        Signatures from exactly `f + 1` oracles are expected.
    /// @param reportHash The hash of the report data with context information to be attested.
    /// @param attestation A concatenation of the attribution bitmask, and `f + 1` ECDSA signatures.
    ///
    /// @dev Attestation format:
    ///       a) attribution bitmask (32 bits), and
    ///       b) `n` signatures
    ///         - 64 bytes per signature, a tuple of (r, s) values
    ///         - `s` is normalized such that recovery id `v` for the public key is 0 (=27 for the ecrecover call),
    ///         - the signatures are sorted by oracle index (ascending order), the least significant set bit of the
    ///           attribution bitmask corresponds to the first signature provided
    function verifyAttestation(
        uint256[32] storage s_keys,
        uint8 n,
        uint8 f,
        bytes32 reportHash,
        bytes calldata attestation
    ) internal view {
        // The number of signatures required.
        uint256 t = f + 1;

        // An attestation is rejected if it does not contain exactly a 4 byte attribution bitmask and `t = f + 1`
        // signatures (64 bytes each).
        if (attestation.length != 4 + 64 * t) {
            revert InvalidAttestationLength();
        }

        // Extract the attribution bitmask from the first 4 bytes of the attestation.
        uint256 attributionBitmask = uint32(bytes4(attestation[:4]));

        // The attribution bitmask contains at most 32 bits. For any given configuration (number of nodes: n), only
        // the bottom n bits may be set. The following check ensure that no higher bit is set.
        if (attributionBitmask >= (1 << n)) {
            revert InvalidAttestationAttributionBitmask();
        }

        // Keep track of the number of successfully verified signatures.
        uint256 numValidSignatures = 0;

        // At this point we already know we that data for f + 1 signatures was passed. However, more or less than f + 1
        // bits may be set in the attribution bitmask. The final check `numValidSignatures == t` captures the
        // "less than" case. To handle the "more than" case according, we track the exact number of bits set in the
        // attribution bitmask and enforce it to be equal to f + 1.
        //
        // Caution: At a first glance the "more than" case may not seem safety critical. However, the implementation
        // of the main verification loop below, assume this fact to be checked!.
        uint256 numAttributionBitsSet = 0;

        assembly {
            // We use 4 words (128 bytes) from the start of the free memory as temporary space to store the input
            // data for the ecRecover call. For this temporary use, there is no need to update the free memory pointer.

            // The input for ecRecover is composed of four 32 bytes words:
            //  - hash: set once
            //  - v: signature value, fixed 27, set once
            //  - r: signature value, different for each signature, updated in the loop
            //  - s: signature value, different for each signature, updated in the loop

            // Get the start position of the free memory (mload(0x40)) and set the values `hash` and `v`.
            let ptr_ecrecoverInput_start := mload(0x40)
            mstore(ptr_ecrecoverInput_start, reportHash)
            mstore(add(ptr_ecrecoverInput_start, 32), 27)

            // Get a pointer to the location where r and s will be placed using calldatacopy.
            let ptr_ecrecoverInput_RS := add(ptr_ecrecoverInput_start, 64)

            // Pointer to the next signature, a (r, s) tuple, from the attestation.
            // Initialized to skip the 32 bit of attribution data.
            // Incremented by 64 whenever a new (r, s) tuple was read.
            let ptr_nextSignatureFromAttestation := add(attestation.offset, 4)

            // Loop over the bits set in the attribution bitmask. Conceptually, the variable `i` holds the index of the
            // i-th bit in the attribution bitmask. As optimization, it actually holds the storage slot of the i-th
            // oracles' public key.
            // At the end of each loop iteration, the attribution bitmask is right-shifted by a single bit, so in each
            // iteration the least significant bit can be used to check if the i-th bit of the original bitmask was set.
            for { let i := s_keys.slot } gt(attributionBitmask, 0) {} {
                // The variable `i` holds the storage index of the oracle's public public key. If the it is set, the
                // next signature is verified against the i-th public key.
                if and(attributionBitmask, 1) {
                    // numAttributionBitsSet += 1
                    numAttributionBitsSet := add(numAttributionBitsSet, 1)

                    // Copy the next signature (64 bytes, r and s values) to the input data for ecRecover.
                    // Potentially, the number of bits set in the attribution bitmask exceeds the number of signatures
                    // provided. In this corner case, calldatacopy() would return
                    //  a) zero-bytes (no other calldata after attestation)
                    //  b) or arbitrary data (extra calldata after attestation).
                    // This may lead to over-counting the number of valid signatures, which however, is not a safety
                    // consider in this particular case, because, at the end of this functions, we only accept an
                    // attestation if exactly f + 1 bits are set in the attribution bitmask.
                    calldatacopy(ptr_ecrecoverInput_RS, ptr_nextSignatureFromAttestation, 64)

                    // Update the attestation input pointer to point to the next signature.
                    ptr_nextSignatureFromAttestation := add(ptr_nextSignatureFromAttestation, 64)

                    // Clear the memory where the result of the following ecRecover call will be stored. This ensures
                    // that in case the ecRecover call fails, and therefore does not update the output memory, it will
                    // still be well initialized to the value zero.
                    mstore(0x00, 0) // mem[0:32] = 0

                    // As a first step for signature verification, recover the signer's address from the inputs.
                    // On success, the corresponding recovered address is written to the 1st word of the scratch space.
                    // mem[0:32] = ecrecover(hash, 27, r, s)
                    //  - mem[0:12]... zero bytes
                    //  - mem[12:32]... recovered address
                    pop(
                        staticcall(
                            gas(), //                    gas
                            0x01, //                     address of ecrecover precompile
                            ptr_ecrecoverInput_start, // input memory pointer
                            0x80, //                     input size
                            0x00, //                     output memory pointer
                            0x20 //                      output size
                        )
                    )

                    // Load the recovered and stored public keys.
                    // recoveredKey = mem[0:32]
                    // storedKey = s_keys[i]
                    let recoveredKey := mload(0) // zero if ecRecover failed
                    let storedKey := sload(i) // always non-zero, ensured by `setVerificationKeys`

                    // As a second step for signature verification, compare the recovered public key with the stored
                    // one for the current index `i` and increment the `numValidSignatures` counter if the keys match.
                    //
                    // Attention: Special care needs to be taken in case the above call to the ecRecover failed and thus
                    // `recoveredKey` holds the value zero.
                    //
                    // The optimized check below relies on the fact that `storedKey` can never be zero - an invariant
                    // ensured by the implementation of `setVerificationKeys`. Therefore, the equality instruction
                    // correctly yields zero (and the counter `numValidSignatures` is not incremented) in the particular
                    // case where the call to ecRecover precompile failed.
                    numValidSignatures := add(numValidSignatures, eq(recoveredKey, storedKey))
                }

                // i += 1
                i := add(i, 1)

                // attributionBitmask >>= 1
                attributionBitmask := shr(1, attributionBitmask)
            }
        }

        if (numValidSignatures != t) {
            revert InvalidAttestationNumberOfSignatures();
        }
        if (numAttributionBitsSet != t) {
            revert InvalidAttestationAttributionBitmask();
        }
    }
}
