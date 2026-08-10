// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./OCR3AttestationVerifierBase.sol";

/// @title Base contract for on-deployment linking of a ECDSA or BLS attestation verifier library
/// @dev The constructor takes a library address which must point to a deployed library of the following types:
///        - OCR3DynamicallyDispatchedECDSAAttestationVerifierLib, or
///        - OCR3DynamicallyDispatchedBLSAttestationVerifierLib.
contract OCR3DynamicallyDispatchedAttestationVerifier is OCR3AttestationVerifierBase {
    // Address of the pre-deployed library contract.
    address immutable i_verifierLibraryAddress;

    // Function selectors for the setAttestationVerificationKeys(...) and verifyAttestation(...) functions of the
    // library.
    bytes4 immutable i_selectorSetAttestationVerificationKeys;
    bytes4 immutable i_selectorVerifyAttestation;

    // Placeholder for reserving storage for up to 32 verification keys. The used library stores an implementation
    // specific data structure within these reserved storages slots. 128 words are required for 32 keys, as each key is
    // composed of 4 words in the BLS case.
    uint256[128] s_ocr3SignerPublicKeys;

    // Constructor used to specify the address of the predeployed verifier library to which the calls should be
    // forwarded to.
    constructor(address verifierLibraryAddress) {
        i_verifierLibraryAddress = verifierLibraryAddress;

        // The following call to getSelector() is not only needed to get the correct selector values for performing a
        // delegatecall into the library, but is also security critical, it protects against potential misconfiguration,
        // where i_verifierLibraryAddress does not point to a contract. Additional details are provided in the comment
        // in the _delegatecall(...) helper below.
        (i_selectorSetAttestationVerificationKeys, i_selectorVerifyAttestation) =
        (OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface(verifierLibraryAddress).getSelectors());
    }

    // Helper function to perform a delegate call to the underlying verifier library.
    // Checks if the call was successful, and, on error, reverts using the internally error (if provided).
    function _delegatecall(bytes memory encodedCalldata) private {
        // Delegate call into the library contract.
        //
        // Quote from: https://docs.soliditylang.org/en/latest/control-structures.html
        // The low-level functions call, delegatecall and staticcall return true as their first return value if the
        // account called is non-existent, as part of the design of the EVM. Account existence must be checked prior to
        // calling if needed.
        //
        // We use the low-level delegatecall below, but are still safe against a potential misconfiguration where
        // i_verifierLibraryAddress does not point to a contract.
        // Reason: We actually perform a high-level call in the constructor above. In this case, the compiler
        // automatically adds a safeguard using extcodesize (which works unless the library would SELFDESTRUCT, but no
        // reasonable verifier lib would even contain that opcode).
        (bool success, bytes memory returnData) = i_verifierLibraryAddress.delegatecall(encodedCalldata);
        if (!success) {
            assembly {
                revert(add(32, returnData), mload(returnData))
            }
        }
    }

    function _setAttestationVerificationKeys(uint8 n, bytes calldata ocr3SignerPublicKeys) internal override {
        uint256 storagePtr;
        assembly {
            storagePtr := s_ocr3SignerPublicKeys.slot
        }
        _delegatecall(
            abi.encodeWithSelector(i_selectorSetAttestationVerificationKeys, storagePtr, n, ocr3SignerPublicKeys)
        );
    }

    function _verifyAttestation(
        bytes32 configDigest,
        uint64 seqNr,
        uint8 n,
        uint8 f,
        bytes memory report,
        bytes calldata attestation
    ) internal override {
        bytes32 reportHash = _hashReport(configDigest, seqNr, report);
        uint256 storagePtr;
        assembly {
            storagePtr := s_ocr3SignerPublicKeys.slot
        }
        _delegatecall(abi.encodeWithSelector(i_selectorVerifyAttestation, storagePtr, n, f, reportHash, attestation));
    }
}

/// @title Internal selector interface for dynamically dispatched OCR3 attestation verifier libraries
/// @dev Exposes the selectors for the `setAttestationVerificationKeys(...)`, and `verifyAttestation(...)` functions.
/// @dev Required for delegate-calling into a pre-deployed library, implemented in the `DynamicallyDispatched` shims.
interface OCR3DynamicallyDispatchedAttestationVerifierSelectorInterface {
    function getSelectors() external pure returns (bytes4, bytes4);
}
