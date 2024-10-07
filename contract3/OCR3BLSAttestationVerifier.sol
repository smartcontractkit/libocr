// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./OCR3AttestationVerifierBase.sol";
import "./OCR3BLSAttestationVerifierLib.sol";

/// @title Base contract for OCR3 report verification using BLS signatures
/// @dev An application contract should inherit from this contract to compile the provided functionality into the
///      application contract.
/// @dev To instead use a dynamically-dispatched (pre-deployed) variant, see OCR3DynamicallyDispatchedAttestationVerifier.
contract OCR3BLSAttestationVerifier is OCR3AttestationVerifierBase {
    // Reserve storage for up to 32 BLS public keys. The current implementation supports up to 32 keys, limited by the
    // width of the attribution bitmask used. Keeping the array size fixed at 32 entries is fine for smaller
    // configurations, storage costs are only payed for the used number of keys.
    OCR3BLSAttestationVerifierLib.G2PointAffine[32] s_keys;

    function _setVerificationKeys(uint8 n, bytes calldata keys) internal override {
        OCR3BLSAttestationVerifierLib.setVerificationKeys(s_keys, n, keys);
    }

    function _verifyAttestation(
        bytes32 configDigest,
        uint64 seqNr,
        uint8 n,
        uint8 f,
        bytes memory report,
        bytes calldata attestation
    ) internal view override {
        bytes32 reportHash = _hashReport(configDigest, seqNr, report);
        OCR3BLSAttestationVerifierLib.verifyAttestation(s_keys, n, f, reportHash, attestation);
    }
}
