// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./OCR3ECDSAAttestationVerifierLib.sol";

/// @title Shim for the core ECDSA attestation verifier library, allowing it to be pre-deployed separately.
/// @dev The function modifiers of the main interface functions are updated from internal to external.
library OCR3DynamicallyDispatchedECDSAAttestationVerifierLib {
    function setVerificationKeys(uint256[32] storage s_keys, uint8 n, bytes calldata keys) external {
        OCR3ECDSAAttestationVerifierLib.setVerificationKeys(s_keys, n, keys);
    }

    function verifyAttestation(
        uint256[32] storage s_keys,
        uint8 n,
        uint8 f,
        bytes32 reportHash,
        bytes calldata attestation
    ) external view {
        OCR3ECDSAAttestationVerifierLib.verifyAttestation(s_keys, n, f, reportHash, attestation);
    }

    // Function to initialize the selectors for delegate-calling into this library.
    // Derived using keccak256 from the function signatures (without parameter names):
    //  - keccak256("setVerificationKeys(uint256[32] storage,uint8,bytes)")[:4]
    //  - keccak256("verifyAttestation(uint256[32] storage,uint8,uint8,bytes32,bytes)")[:4]
    function getSelectors() external pure returns (bytes4, bytes4) {
        return (
            OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.setVerificationKeys.selector,
            OCR3DynamicallyDispatchedECDSAAttestationVerifierLib.verifyAttestation.selector
        );
    }
}
