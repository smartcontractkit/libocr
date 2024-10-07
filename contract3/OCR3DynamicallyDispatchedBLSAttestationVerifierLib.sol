// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./OCR3BLSAttestationVerifierLib.sol";

/// @title Shim for the core BLS attestation verifier library, allowing it to be pre-deployed separately.
/// @dev The function modifiers of the main interface functions are updated from internal to external.
library OCR3DynamicallyDispatchedBLSAttestationVerifierLib {
    function setVerificationKeys(
        OCR3BLSAttestationVerifierLib.G2PointAffine[32] storage s_verificationKeys,
        uint8 n,
        bytes calldata keys
    ) external {
        OCR3BLSAttestationVerifierLib.setVerificationKeys(s_verificationKeys, n, keys);
    }

    function verifyAttestation(
        OCR3BLSAttestationVerifierLib.G2PointAffine[32] storage s_verificationKeys,
        uint8 n,
        uint8 f,
        bytes32 reportHash,
        bytes calldata attestation
    ) external view {
        OCR3BLSAttestationVerifierLib.verifyAttestation(s_verificationKeys, n, f, reportHash, attestation);
    }

    // Function to initialize the selectors for delegate-calling into this library.
    // Derived using keccak256 from the function signatures (without parameter names):
    //  - keccak256("setVerificationKeys(OCR3BLSAttestationVerifierLib.G2PointAffine[32] storage,uint8,bytes)")[:4]
    //  - keccak256(
    //        "verifyAttestation(OCR3BLSAttestationVerifierLib.G2PointAffine[32] storage,uint8,uint8,bytes32,bytes)"
    //    )[:4]
    function getSelectors() external pure returns (bytes4, bytes4) {
        return (
            OCR3DynamicallyDispatchedBLSAttestationVerifierLib.setVerificationKeys.selector,
            OCR3DynamicallyDispatchedBLSAttestationVerifierLib.verifyAttestation.selector
        );
    }
}
