// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./OCR3AttestationVerifierErrors.sol";

/// @title Internal library for attestation verification using ECDSA signatures
/// @dev Developers should inherit from OCR3ECDSAAttestationVerifier for this functionality to be compiled with the
///      application contract or see OCR3DynamicallyDispatchedAttestationVerifier for using it as a pre-deployed
///      library.
library OCR3ECDSAAttestationVerifierLib {
    /// @notice Verifies and stores the provided ECDSA public keys for later use within `verifyAttestation(...)`.
    ///         Reverts if the decoded key count does not match `n`, any key is found invalid, or a duplicate key is
    ///         provided.
    /// @dev Bytes-typed overload of `setAttestationVerificationKeys`. The keys are ABI-decoded into an `address[]`
    ///      and forwarded to the address-typed overload, which performs all validation and storage. This exists for
    ///      callers that hold the keys as opaque bytes (for example, dispatching through the scheme-agnostic
    ///      `_setAttestationVerificationKeys(uint8, bytes)` interface).
    /// @param s_ocr3EcdsaSignerPublicKeys Storage space for holding all ECDSA verification keys (in the form of
    ///        addresses). The first `n` values are populated. The size of 32 is the maximum number of keys supported
    ///        (based on the width of the attribution bitmask). Storage costs are only paid for the number of keys used.
    /// @param n The number of keys expected to be set by this call. Must match the number of keys decoded from the
    ///        `ocr3EcdsaSignerPublicKeys` parameter.
    /// @param ocr3EcdsaSignerPublicKeys The ABI-encoding of the signers' ECDSA public keys, i.e.
    ///        `abi.encode(address[])`.
    function setAttestationVerificationKeys(
        uint256[32] storage s_ocr3EcdsaSignerPublicKeys,
        uint8 n,
        bytes calldata ocr3EcdsaSignerPublicKeys
    ) internal {
        address[] memory signers = abi.decode(ocr3EcdsaSignerPublicKeys, (address[]));
        // Verify that `n` is consistent with the number of keys decoded from `ocr3EcdsaSignerPublicKeys`.
        if (signers.length != n) {
            revert InvalidNumberOfAttestationVerificationKeys();
        }
        setAttestationVerificationKeys(s_ocr3EcdsaSignerPublicKeys, signers);
    }

    /// @notice Verifies and stores the provided ECDSA public keys (as addresses) for later use within
    ///         `verifyAttestation(...)`. Reverts if too many keys are provided, any key is found invalid, or a
    ///         duplicate key is provided.
    /// @dev Address-typed convenience overload of `setAttestationVerificationKeys`. Callers holding an `address[]`
    ///      (the common case on EVM chains) can use this directly, without ABI-encoding the keys into the `bytes`
    ///      layout and without supplying a redundant key count. The `bytes` overload remains available for callers
    ///      that hold the keys as an ABI-encoded `address[]`.
    /// @param s_ocr3EcdsaSignerPublicKeys Storage space for holding all ECDSA verification keys (in the form of
    ///        addresses). The first `ocr3EcdsaSignerPublicKeys.length` values are populated. The size of 32 is the
    ///        maximum number of keys supported (based on the width of the attribution bitmask). Storage costs are only
    ///        paid for the number of keys used.
    /// @param ocr3EcdsaSignerPublicKeys The ECDSA public keys (i.e., addresses) of the OCR3 signers.
    function setAttestationVerificationKeys(
        uint256[32] storage s_ocr3EcdsaSignerPublicKeys,
        address[] memory ocr3EcdsaSignerPublicKeys
    ) internal {
        uint256 n = ocr3EcdsaSignerPublicKeys.length;

        // The maximum of 32 keys is based on the width of the attribution bitmask (currently set to 32 bits).
        if (n > 32) {
            revert MaximumNumberOfAttestationVerificationKeysExceeded();
        }

        // Copy the provided keys to storage (s_ocr3EcdsaSignerPublicKeys). After copying, the i-th 32 byte storage
        // slot holds the i-th 20 byte key (i.e., the signer's address) in its lower bytes.
        for (uint256 i = 0; i < n; ++i) {
            address key = ocr3EcdsaSignerPublicKeys[i];
            // Ensure the key/address is non-zero. Clearly the value 0x0000000000000000000000000000000000000000 is an
            // invalid key/address to sign attestations from.
            if (key == address(0)) {
                revert InvalidAttestationVerificationKey();
            }
            // Reject duplicate keys. Each key must be unique to ensure that no two oracles share the same signing
            // identity, which could otherwise lead to incorrect attribution during attestation verification. `n` is at
            // most 32, so the quadratic scan over the already-checked keys is cheap.
            for (uint256 j = 0; j < i; ++j) {
                if (ocr3EcdsaSignerPublicKeys[j] == key) {
                    revert DuplicateAttestationVerificationKey();
                }
            }
            s_ocr3EcdsaSignerPublicKeys[i] = uint256(uint160(key));
        }
    }

    /// @notice Verifies the attestation for the given report hash. Reverts on verification failure.
    /// @param s_ocr3EcdsaSignerPublicKeys A list (stored in the application contract) holding the ECDSA public keys
    ///        (i.e., addresses) of all oracles.
    /// @param n The total number of oracles. Used to verify the attribution bitmask as part of the attestation data.
    /// @param f Maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly.
    ///        Signatures from exactly `f + 1` oracles are expected.
    /// @param reportHash The hash of the report data with context information to be attested.
    /// @param attestation A concatenation of the attribution bitmask, and `f + 1` ECDSA signatures.
    ///
    /// @dev Attestation format:
    ///       a) attribution bitmask (32 bits), and
    ///       b) `f+1` signatures
    ///         - 64 bytes per signature, a tuple of (r, s) values
    ///         - `s` is normalized such that recovery id `v` for the public key is 0 (=27 for the ecrecover call),
    ///         - the signatures are sorted by oracle index (ascending order), the least significant set bit of the
    ///           attribution bitmask corresponds to the first signature provided
    function verifyAttestation(
        uint256[32] storage s_ocr3EcdsaSignerPublicKeys,
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
        // of the main verification loop below assumes this fact to be checked.
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
            for { let i := s_ocr3EcdsaSignerPublicKeys.slot } gt(attributionBitmask, 0) {} {
                // The variable `i` holds the storage index of the oracle's public key. If it is set, the next
                // signature is verified against the i-th public key.
                if and(attributionBitmask, 1) {
                    // numAttributionBitsSet += 1
                    numAttributionBitsSet := add(numAttributionBitsSet, 1)

                    // Copy the next signature (64 bytes, r and s values) to the input data for ecRecover.
                    // Potentially, the number of bits set in the attribution bitmask exceeds the number of signatures
                    // provided. In this corner case, calldatacopy() would return
                    //  a) zero-bytes (no other calldata after attestation)
                    //  b) or arbitrary data (extra calldata after attestation).
                    // This may lead to over-counting the number of valid signatures, which however, is not a safety
                    // concern in this particular case, because, at the end of this function, we only accept an
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
                    // storedKey = s_ocr3EcdsaSignerPublicKeys[i]
                    let recoveredKey := mload(0) // zero if ecRecover failed
                    let storedKey := sload(i) // always non-zero, ensured by `setAttestationVerificationKeys`

                    // As a second step for signature verification, compare the recovered public key with the stored
                    // one for the current index `i` and increment the `numValidSignatures` counter if the keys match.
                    //
                    // Attention: Special care needs to be taken in case the above call to ecRecover failed, in which
                    // case `recoveredKey` holds the value zero. The `gt(recoveredKey, 0)` guard ensures such a failed
                    // recovery is never counted as valid. Under normal operation this is redundant with the equality
                    // check, because `storedKey` is guaranteed non-zero by `setAttestationVerificationKeys` and a zero
                    // `recoveredKey` therefore cannot match it. The guard is defense in depth: should that invariant
                    // ever be violated and a stored slot hold zero, a failed ecRecover (also zero) would otherwise
                    // match the zero slot and be counted as valid. The guard adds ~15 gas per signature.
                    numValidSignatures :=
                        add(numValidSignatures, and(eq(recoveredKey, storedKey), gt(recoveredKey, 0)))
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
