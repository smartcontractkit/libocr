// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// Raised when the number of provided verifications keys does not match the expected number of keys (parameter: n).
error InvalidNumberOfKeys();

// Raised when the provided verifications keys are of invalid size.
error KeysOfInvalidSize();

// Raised when an attempt to set more than 32 verification keys is made.
// An upper limit of 32 keys is enforced by the width of the bitmask used for the attribution data.
error MaximumNumberOfKeysExceeded();

// Raised when a provided verification key is found invalid.
// Potential causes for invalid keys are, for example:
//  - ECDSA: the value 0x0000000000000000000000000000000000000000
//  - BLS: a key with an invalid proof-of-possession
error InvalidKey();

// Raised when the signature verification failed for the provided attestation.
error InvalidAttestation();

// Raised when the provided attestation contains (or is composed of) an invalid number of signatures.
error InvalidAttestationNumberOfSignatures();

// Raised when a provided attestation is of invalid size.
// The expected size depends on the signature scheme used, for example:
//  - ECDSA: 4 + 64*(f+1) bytes
//  - BLS:             37 bytes
error InvalidAttestationLength();

// Raised when the attribution bitmask of a given attestation is invalid.
// Exactly f+1 of the least significant n bits must be set.
error InvalidAttestationAttributionBitmask();
