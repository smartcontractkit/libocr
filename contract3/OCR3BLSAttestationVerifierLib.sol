// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./OCR3AttestationVerifierErrors.sol";

/// @title Internal library for attestation verification using BLS aggregate signatures over BN254
/// @dev Developers should inherit from OCR3BLSAttestationVerifier for this functionality to be compiled with the
///      application contract or see OCR3DynamicallyDispatchedAttestationVerifier for using it as a pre-deployed
///      library.
library OCR3BLSAttestationVerifierLib {
    // #################################################################################################################
    // # Primary interface functions                                                                                   #
    // #################################################################################################################

    /// @notice Verifies and stores the provided BLS public keys for later use within `verifyAttestation(...)`.
    /// @param s_keys A list (stored in the application contract) holding the BLS public keys of all oracles in
    ///        uncompressed, affine form. The size of 32 is the maximum number of keys supported (based on the width of
    //         the attribution bitmask). Storage costs are only payed for the number of keys used `n`.
    /// @param n The number of keys expected to be set by this call. Must match the actual number of keys present in
    ///        the `keys` parameter.
    /// @param keys A concatenation of `n` BLS public keys (uncompressed affine format with appended
    ///             proof-of-possession).
    ///
    /// @dev Detailed format for parameter `keys`:
    ///       - concatenation of `n` values
    ///       - each value is composed of
    ///          - BLS public key in uncompressed affine format (128 bytes)
    ///          - proof-of-possession, i.e., a BLS signature (32 bytes compressed point, 1 byte counter)
    ///
    /// @dev The proof-of-possession is computed over:
    ///       - a domain separation tag (dst)
    ///       - the public key data (pk),
    ///      using the following format: keccak256(keccak256(abi.encode(dst, pk)), counter_byte).
    ///                                            ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^... innerHash
    ///                                  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^... outerHash
    ///
    function setVerificationKeys(G2PointAffine[32] storage s_keys, uint8 n, bytes calldata keys) internal {
        // Verify that `n` is consistent with the amount of data being passed in the `keys` parameter.
        // The maximum of 32 keys is based on the width of the attribution bitmask (currently set to 32 bits).
        if (keys.length % KEYSIZE_WITH_POP != 0) {
            revert KeysOfInvalidSize();
        }
        if (keys.length / KEYSIZE_WITH_POP != n) {
            revert InvalidNumberOfKeys();
        }
        if (n > 32) {
            revert MaximumNumberOfKeysExceeded();
        }

        // Temporary variable for the key to be verified in memory before it is written to storage. This value is
        // updated in each loop iteration and copied to storage after successful verification.
        G2PointAffine memory key;

        // Reserve space for 160 bytes (=32x5 bytes) holding the domain separation tag and public key.
        // Initialize with the domain separation tag. The public key is set/updated in each loop iteration.
        uint256[5] memory innerHashInput;
        innerHashInput[0] = DOMAIN_SEPARATION_TAG_BLS_PROOF_OF_POSSESSION;

        // Reserve space for 33 bytes holding the input for the outer hash.
        bytes32[2] memory outerHashInput;

        // Main loop to split and verify all keys.
        uint256 inputPosition = 0;
        for (uint8 i = 0; i < n; ++i) {
            // Read the next key from the input data, and update the input position accordingly. The validity of the
            // loaded values (i.e., that they describe a valid group element of the elliptic curve) is not verified
            // here. The verification is performed implicitly when the proof-of-possession is verified, because the
            // involved call to the ecPairing precompiles can only succeed if all inputs are indeed valid points in the
            // respective groups. In this regard, we (and the auditors) checked the specification (EIP-197, and the
            // Ethereum Yellowpaper) as well as the geth source code implementing the precompile.
            innerHashInput[1] = key.x_imag = uint256(bytes32(keys[inputPosition:inputPosition += 32]));
            innerHashInput[2] = key.x_real = uint256(bytes32(keys[inputPosition:inputPosition += 32]));
            innerHashInput[3] = key.y_imag = uint256(bytes32(keys[inputPosition:inputPosition += 32]));
            innerHashInput[4] = key.y_real = uint256(bytes32(keys[inputPosition:inputPosition += 32]));

            // Read the proof-of-possession, i.e., the signature value and signature counter byte.
            // The signature counter byte (type: bytes1) is, considered as 32 byte word, is left-aligned in memory.
            // (The most significant byte holds the counter value.) This alignment is critical for the hash computation
            // of `outerHash := keccak256(inputHash || counter byte)`.
            bytes32 popSignature = bytes32(keys[inputPosition:inputPosition += 32]);
            bytes1 popSignatureCounter = keys[inputPosition];
            inputPosition += 1;

            // Prepare the signature verification by computing the actual value that was signed, i.e., `outerHash`,
            // itself being based on `innerHash`.
            bytes32 innerHash;
            assembly {
                innerHash := keccak256(innerHashInput, 160)
            }

            bytes32 outerHash;
            outerHashInput[0] = innerHash;
            outerHashInput[1] = popSignatureCounter;
            assembly {
                outerHash := keccak256(outerHashInput, 33)
            }

            // Actually verify the signature against the computed outer hash and revert on failure.
            if (!_verifySignature(key, outerHash, popSignature)) {
                revert InvalidKey();
            }

            // Write the verified key to the application contract's storage.
            s_keys[i] = key;
        }
    }

    /// @notice Verifies the attestation for the given report hash. Reverts on verification failure.
    /// @param s_keys A list (stored in the application contract) holding the BLS public keys of all oracles in
    ///        uncompressed, affine form.
    /// @param n The total number of oracles. Used to verify the attribution bitmask as part of the attestation data.
    /// @param f Maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly.
    ///        An aggregate signature from exactly `f + 1` oracles is expected.
    /// @param reportHash hash of the report data with context information to be attested
    /// @param attestation A concatenation of the attribution bitmask (4 bytes), a BLS aggregate signature (32 bytes),
    //         and the counter byte for the BLS signature.
    function verifyAttestation(
        G2PointAffine[32] storage s_keys,
        uint8 n,
        uint8 f,
        bytes32 reportHash,
        bytes calldata attestation
    ) internal view {
        // Ensure the attestation has the proper format.
        //  - 4 bytes: bitmask for attribution data
        //  - 33 bytes: aggregate signature with counter
        if (attestation.length != 37) {
            revert InvalidAttestationLength();
        }

        // Compute the aggregate public key for verification from the attribution data in the attestation. The call
        // reverts if the attribution data is invalid or contains less than the required number of signers.
        OCR3BLSAttestationVerifierLib.G2PointAffine memory verificationKey;
        uint256 numSigners;
        (verificationKey, numSigners) = _computeVerificationKey(s_keys, n, attestation);

        // Ensure the correct number of nodes participated in generating the attestation, f+1 signers are required.
        if (numSigners != f + 1) {
            revert InvalidAttestationNumberOfSignatures();
        }

        bytes32 signature;
        bytes32 reportHashWithCounterByte;
        assembly {
            // scratchSpace[0:33] = signature || counter byte
            calldatacopy(0, add(attestation.offset, 4), 33)

            // signature = scratchSpace[0:32] = attestation[0:32] = attestation without counter byte
            signature := mload(0)

            // scratchSpace[0:32] = reportHash
            mstore(0, reportHash)

            // Here, scratchSpace[0:33] = reportHash || attestation counter byte
            // reportHashWithCounterByte = keccak256(scratchSpace[0:33])
            reportHashWithCounterByte := keccak256(0, 33)
        }

        if (!_verifySignature(verificationKey, reportHashWithCounterByte, signature)) {
            revert InvalidAttestation();
        }
    }

    // #################################################################################################################
    // # Library wide constants                                                                                        #
    // #################################################################################################################

    // Various constants defined below. These values are carefully chosen, and the implementation assumes that these
    // values are set as they are - be careful when changing them. Defined here to improving readability of the code.

    // Size of a BLS public key in compressed affine format in bytes.
    uint256 private constant KEYSIZE = 128;

    // Size of a BLS public key in compressed affine format including a proof-of-possession in bytes.
    uint256 private constant KEYSIZE_WITH_POP = 161;

    // // Size of the domain separation tag used for the proof-of-possession in bytes.
    uint256 private constant DOMAIN_SEPARATION_TAG_SIZE = 32;

    // Size of a BLS signature in compressed format with appended counter byte in bytes.
    uint256 private constant SIGNATURE_SIZE = 33;

    // Base field for G1 is 𝔽ₚ
    // 0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47
    // https://github.com/ethereum/EIPs/blob/master/EIPS/eip-196.md#specification
    uint256 private constant P = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

    // since p % 4 = 3, x^{(p+1)/4} = sqrt(x), if sqrt(x) exists
    uint256 private constant SQRT_POWER = (P + 1) >> 2;

    /// In 𝔽ₚ, modExp(x, P - 2) can be used to compute the modular inverse element.
    /// I.e., x * modExp(x, P - 2) == 1 (mod P).
    uint256 private constant INVERSE_POWER = P - 2;

    /// Domain separation tag used for verifying proof-of-possessions.
    /// Picked value: SHA3-256(b"DOMAIN_SEPARATION_TAG_BLS_PROOF_OF_POSSESSION")
    uint256 private constant DOMAIN_SEPARATION_TAG_BLS_PROOF_OF_POSSESSION =
        0x79537812dfe48a92fc860b8b010e8d6078b5c19e7037c4cf07f7bed69b54fffc;

    // G2 points is on curve y² = x³ + twistB in GFp2, where TWIST_B = TWIST_B_REAL + i*TWIST+B_IMAG
    uint256 private constant TWIST_B_IMAG = 0x009713b03af0fed4cd2cafadeed8fdf4a74fa084e52d1852e4a2bd0685c315d2;
    uint256 private constant TWIST_B_REAL = 0x2b149d40ceb8aaae81be18991be06ac3b5b4c5e559dbefa33267e6dc24a138e5;

    // #################################################################################################################
    // # Data structures                                                                                               #
    // #################################################################################################################

    // Regular/affine representation of G2 point over y² = x³ + twistB.
    // This format is required for the Ethereum precompile for pairing verification.
    struct G2PointAffine {
        uint256 x_imag;
        uint256 x_real;
        uint256 y_imag;
        uint256 y_real;
    }

    // Jacobian representation of a point on G2 point over y² = x³ + twistB * z³.
    // The jacobian format is used for more gas-efficient computation on intermediate values.
    // A point (x, y) is represented by the jacobian coordinates (X, Y, Z) such that the following equations hold:
    //   x = X / Z²
    //   y = Y / Z³
    struct G2PointJacobian {
        uint256 x_imag;
        uint256 x_real;
        uint256 y_imag;
        uint256 y_real;
        uint256 z_imag;
        uint256 z_real;
    }

    // #################################################################################################################
    // # Core BLS signature verification and BLS public key aggregation functions                                      #
    // #################################################################################################################

    /// @notice Uses `verificationKey` to check if the provided signature is valid signature for the given data.
    ///         Reverts if the signature is invalid.
    /// @param verificationKey the public key used for signature verification, an aggregated public key in our case
    /// @param reportHash the keccak256 hash of the underlying (claimed) message the signature should attest
    /// @param signature the signature data, composed of a G1 point in compressed form and a 1 byte counter required
    ///                  for hashing the data to elliptic curve
    /// @return ok true if and only if the signature verification was successful
    function _verifySignature(G2PointAffine memory verificationKey, bytes32 reportHash, bytes32 signature)
        private
        view
        returns (bool)
    {
        bool ok;
        uint256 sX;
        uint256 sY;
        uint256 hX;
        uint256 hY;

        // Unpack the signature data into a point on G1.
        // Return false if the unpacked failed.
        (ok, sX, sY) = _unpackToG1(uint256(signature));
        if (!ok) {
            return false;
        }

        // Unpack the hash (with its 2nd most significant bit set to zero) into a point on G1.
        //  - The field modulus is a 254 bit prime.
        //  - The most significant bit of the hash encodes the 'sign' of the y-coordinate of the corresponding point on
        //    G1.
        //  - The 2nd highest bit serves no purpose and is cleared to reduce the hashToCurve rejection rate.
        //    0xbfff... == 0b1011_1111_1111_1111...
        (ok, hX, hY) =
            _unpackToG1(uint256(reportHash) & 0xbfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
        if (!ok) {
            return false;
        }

        uint256[12] memory pairingCheckInput = [
            sX,
            sY,
            0x198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2, // (-G2).x.imag
            0x1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed, // (-G2).x.real
            0x275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec, // (-G2).y.imag
            0x1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d, // (-G2).y.real
            hX,
            hY,
            verificationKey.x_imag,
            verificationKey.x_real,
            verificationKey.y_imag,
            verificationKey.y_real
        ];

        bool[1] memory pairingCheckOk;
        assembly {
            // Call the ecPairing precompiled using staticcall to verify the pairing equation and check that the
            // operation was successful.
            if iszero(
                staticcall(
                    gas(), //             gas provided to the ecPairing call, 113000 needed, but use gas() for robustness
                    0x08, //              address of the ecPairing pre-compile
                    pairingCheckInput, // input
                    384, //               input size in bytes
                    pairingCheckOk, //    address to store the result of the ecPairing call (0 == failed, 1 == ok)
                    32 //                 size of result in bytes
                )
            ) {
                // The staticcall failed unexpectedly (e.g., due to a gas misconfiguration).
                // This should not happen during normal operation, even if a signature is found invalid.
                revert(0, 0)
            }
        }
        return pairingCheckOk[0];
    }

    /// @notice Given a report's attestation, extract the bitmask indicating which oracles signed the report to produce
    //          the attestation, then load and aggregate all corresponding public keys to derive the verification key.
    /// @param n the number of oracles
    /// @param attestation the report's attestation data, composed of a bitmask (32 bit) of attribution data and an
    ///        aggregate signature.
    /// @return verificationKey the computed verification key
    /// @return numSigners the number of keys aggregated to form the verification key
    function _computeVerificationKey(G2PointAffine[32] storage s_keys, uint256 n, bytes calldata attestation)
        private
        view
        returns (G2PointAffine memory verificationKey, uint256 numSigners)
    {
        // Ensure the attribution bitmask is non-zero and that at the highest possibly set bit cannot exceed the bit for
        // the oracle with highest index. The non-zero check is needed to ensure we have a valid starting point for the
        // aggregation of the keys.
        //
        // Example for n = 7
        //  - allowed bitmask: 00000000_00000000_00000000_0xxxxxxx, x from {0, 1}
        //  - 1 << 7:          00000000_00000000_00000000_10000000
        uint256 bitmask = uint256(uint32(bytes4(attestation[:4])));
        if (bitmask == 0 || bitmask >= (1 << n)) {
            revert InvalidAttestationAttributionBitmask();
        }

        // Move index `i` to point to the least significant bit set the attribution data.
        uint256 i = 0;
        while ((bitmask & 1) == 0) {
            bitmask >>= 1;
            ++i;
        }

        // Set the basis for the aggregation to the public key corresponding to the lowest bit set.
        numSigners = 1;
        G2PointJacobian memory vkJacobian = _jacobian(s_keys[i]);

        // Move to the next possibly set bit.
        ++i;
        bitmask >>= 1;

        while (bitmask > 0) {
            if ((bitmask & 1) > 0) {
                vkJacobian = _addPoints(vkJacobian, s_keys[i]); // add the public key of the next oracle
                ++numSigners; // keep track of the number of keys we aggregated
            }
            bitmask >>= 1;
            ++i;
        }

        // Finalize the verification key by converting into its affine representation.
        verificationKey = _affine(vkJacobian);

        return (verificationKey, numSigners);
    }

    /// @notice Given a value v, representing (1) the x coordinate of a point on G1 and (2) the sign of the
    ///         corresponding y coordinate in its most significant bit, return the point on G1. Revert if the resulting
    ///         point is not a valid point on G1.
    /// @param v the compressed representation of a point on G1
    /// @return ok true if and only if v could be unpacked to a Point on G1 successfully
    /// @return x the x coordinate of the corresponding point on G1
    /// @return y the y coordinate of the corresponding point on G1
    function _unpackToG1(uint256 v) private view returns (bool ok, uint256 x, uint256 y) {
        // Clear top bit of v, save the result to x.
        x = v & 0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff;
        if (x >= P) {
            return (false, 0, 0);
        }

        // Compute rhs = x³ + 3; curve equation: y^2 = x³ + 3 (mod P).
        uint256 rhs = addmod(mulmod(mulmod(x, x, P), x, P), 3, P);

        // Solve for y by computing y = sqrt(x³ + 3) = sqrt(rhs).
        // Ensure that the sqrt actually exists.
        y = _modExp(rhs, SQRT_POWER);
        if (mulmod(y, y, P) != rhs) {
            return (false, 0, 0);
        }

        // Invert the value of y, depending on the original lsb of y encoded into the msb of v.
        uint256 y_lsb = v >> 255;
        if ((y & 1) != y_lsb) {
            y = (P - y) % P;
        }

        return (true, x, y);
    }

    function _gfp2_add(uint256 a_imag, uint256 a_real, uint256 b_imag, uint256 b_real)
        private
        pure
        returns (uint256 r_imag, uint256 r_real)
    {
        assembly {
            // r_imag = a_imag + b_imag  (mod P)
            r_imag := addmod(a_imag, b_imag, P)

            // r_real = a_real + b_real  (mod P)
            r_real := addmod(a_real, b_real, P)
        }
        return (r_imag, r_real);
    }

    function _gfp2_sub(uint256 a_imag, uint256 a_real, uint256 b_imag, uint256 b_real)
        private
        pure
        returns (uint256 r_imag, uint256 r_real)
    {
        assembly {
            // r_imag = a_imag - b_imag  (mod P)
            r_imag := addmod(a_imag, sub(P, b_imag), P)

            // r_real = a_real - b_real  (mod P)
            r_real := addmod(a_real, sub(P, b_real), P)
        }
        return (r_imag, r_real);
    }

    function _gfp2_mul(uint256 a_imag, uint256 a_real, uint256 b_imag, uint256 b_real)
        private
        pure
        returns (uint256 r_imag, uint256 r_real)
    {
        assembly {
            // r_imag = (a_real * b_imag) + (a_imag * b_real)  (mod P)
            r_imag := addmod(mulmod(a_real, b_imag, P), mulmod(a_imag, b_real, P), P)

            // r_real = (a_real * b_real) - (a_imag * b_imag)  (mod P)
            r_real := addmod(mulmod(a_real, b_real, P), sub(P, mulmod(a_imag, b_imag, P)), P)
        }
        return (r_imag, r_real);
    }

    function _gfp2_square(uint256 a_imag, uint256 a_real) private pure returns (uint256 r_imag, uint256 r_real) {
        assembly {
            // r_imag = 2 * (s_real * s_imag)  (mod P)
            r_imag := mulmod(a_real, a_imag, P)
            r_imag := addmod(r_imag, r_imag, P)

            // r_real = a_real² - a_imag²  (mod P)
            r_real := addmod(mulmod(a_real, a_real, P), sub(P, mulmod(a_imag, a_imag, P)), P)
        }
        return (r_imag, r_real);
    }

    function _gfp2_scalar_mul(uint256 s, uint256 a_imag, uint256 a_real)
        private
        pure
        returns (uint256 r_imag, uint256 r_real)
    {
        assembly {
            // r_imag = s * a_imag  (mod P)
            r_imag := mulmod(s, a_imag, P)

            // r_real = s * a_real  (mod P)
            r_real := mulmod(s, a_real, P)
        }
        return (r_imag, r_real);
    }

    /// @notice Return the multiplicative inverse of x in GFp2
    /// @param x_imag value to return inverse of
    /// @param x_real value to return inverse of
    /// @return inv_imag such that x * inv = 1 in GFp2
    /// @return inv_real such that x * inv = 1 in GFp2
    /// @dev bigModExp call is expensive, about 1349 gas. (See gas_cost_modexp.py)
    function _gfp2_inverse(uint256 x_imag, uint256 x_real) private view returns (uint256 inv_imag, uint256 inv_real) {
        // Note that 1/(a+ib) = (a-ib)/(a-ib)(a+ib) = a/(a²+b²) - ib/(a²+b²)
        uint256 denom = addmod(mulmod(x_real, x_real, P), mulmod(x_imag, x_imag, P), P);
        denom = _modExp(denom, INVERSE_POWER);

        inv_imag = mulmod(P - x_imag, denom, P); // -b/(a²+b²)
        inv_real = mulmod(x_real, denom, P); // a/(a²+b²)

        return (inv_imag, inv_real);
    }

    /// @notice Return a+b, in Jacobian coordinates
    /// @param a first summand, in Jacobian coordinates; must be on curve
    /// @param b second summand, in affine coordinates; must be on curve
    /// @return result a+b in Jacobian coordinates
    ///
    /// @dev No checking is done to ensure the points are on the curve, etc. This is expected to be used in a loop.
    ///      Validate the inputs prior to using this!
    ///
    /// @dev Will fail if a and b represent the same point on the curve, or if they sum to zero. The value of a must be
    ///      non-zero (b, an affine representation, cannot be zero). It is expected that this will be used in
    ///      cryptographic contexts, where the above failure conditions are vanishingly unlikely.
    ///
    /// @dev The sum is kept in Jacobian coordinates so that this can be used to sum multiple points, avoiding expensive
    ///      inversions in conversion to affine coordinates until the sum is completed.
    function _addPoints(G2PointJacobian memory a, G2PointAffine memory b)
        private
        pure
        returns (G2PointJacobian memory result)
    {
        // Implementation based on:
        // https://hyperelliptic.org/EFD/g1p/auto-code/shortw/jacobian/addition/madd-2007-bl.op3
        //
        // The input points a and b are given in the following format:
        //  - a is given in jacobian coordinates (x.imag, x.real, y.imag, y.real, z.imag, z.real)
        //  - b is given in affine coordinates (x.imag, x.real, y.imag, y.real, z.imag, z.real)
        //
        // The input points a and b translate to the variables from hyperelliptic.org as follows:
        //  - X1 <=> (a.x.imag, a.x.real)
        //  - Y1 <=> (a.y.imag, a.y.real)
        //  - Z1 <=> (a.z.imag, a.z.real)
        //  - X2 <=> (b.x.imag, b.x.real)
        //  - Y2 <=> (b.y.imag, b.y.real)
        //  - X3 <=> (result.x.imag, result.x.real)
        //  - Y3 <=> (result.y.imag, result.y.real)
        //  - Z3 <=> (result.z.imag, result.z.real)
        //
        // Instructions as given by hyperelliptic.org; optimized order to reduce stack depth:
        // Z1 = a.z
        // Z1Z1 = Z1 ^ 2
        // t0 = Z1 * Z1Z1
        // X2 = b.x
        // U2 = X2 * Z1Z1
        // X1 = a.x
        // H = U2 - X1
        // t9 = Z1 + H
        // t10 = t9 ^ 2
        // t11 = t10 - Z1Z1
        // HH = H ^ 2
        // Z3 = t11 - HH
        // result.z = Z3
        // I = 4 * HH
        // J = H * I
        // V = X1 * I
        // Y2 = b.y
        // S2 = Y2 * t0
        // Y1 = a.y
        // t1 = S2 - Y1
        // r = 2 * t1
        // t2 = r ^ 2
        // t4 = t2 - J
        // t6 = Y1 * J
        // t7 = 2 * t6
        // t3 = 2 * V
        // X3 = t4 - t3
        // result.x = X3
        // t5 = V - X3
        // t8 = r * t5
        // Y3 = t8 - t7
        // result.y = Y3

        assembly {
            // Load P into local variable p to reduce bytecode size.
            let p := P

            // Z1 = a.z
            let r0i := mload(add(a, 0x80))
            let r0r := mload(add(a, 0xa0))

            // Z1Z1 = Z1 ^ 2
            let r1i := mulmod(r0i, r0r, p)
            r1i := addmod(r1i, r1i, p)
            let r1r := addmod(mulmod(r0r, r0r, p), sub(p, mulmod(r0i, r0i, p)), p)

            // t0 = Z1 * Z1Z1
            let r2i := addmod(mulmod(r0r, r1i, p), mulmod(r0i, r1r, p), p)
            let r2r := addmod(mulmod(r0r, r1r, p), sub(p, mulmod(r0i, r1i, p)), p)

            // X2 = b.x
            let r3i := mload(b)
            let r3r := mload(add(b, 0x20))

            // U2 = X2 * Z1Z1
            let r4i := addmod(mulmod(r3r, r1i, p), mulmod(r3i, r1r, p), p)
            let r4r := addmod(mulmod(r3r, r1r, p), sub(p, mulmod(r3i, r1i, p)), p)

            // X1 = a.x
            r3i := mload(a)
            r3r := mload(add(a, 0x20))

            // H = U2 - X1
            r4i := addmod(r4i, sub(p, r3i), p)
            r4r := addmod(r4r, sub(p, r3r), p)

            // t9 = Z1 + H
            r0i := addmod(r0i, r4i, p)
            r0r := addmod(r0r, r4r, p)

            // t10 = t9 ^ 2
            let tmp := r0i
            r0i := mulmod(r0i, r0r, p)
            r0i := addmod(r0i, r0i, p)
            r0r := addmod(mulmod(r0r, r0r, p), sub(p, mulmod(tmp, tmp, p)), p)

            // t11 = t10 - Z1Z1
            r0i := addmod(r0i, sub(p, r1i), p)
            r0r := addmod(r0r, sub(p, r1r), p)

            // HH = H ^ 2
            r1i := mulmod(r4i, r4r, p)
            r1i := addmod(r1i, r1i, p)
            r1r := addmod(mulmod(r4r, r4r, p), sub(p, mulmod(r4i, r4i, p)), p)

            // Z3 = t11 - HH
            r0i := addmod(r0i, sub(p, r1i), p)
            r0r := addmod(r0r, sub(p, r1r), p)

            // result.z = Z3
            mstore(add(result, 0x80), r0i)
            mstore(add(result, 0xa0), r0r)

            // I = 4 * HH
            r0i := mulmod(4, r1i, p)
            r0r := mulmod(4, r1r, p)

            // J = H * I
            r1i := addmod(mulmod(r4r, r0i, p), mulmod(r4i, r0r, p), p)
            r1r := addmod(mulmod(r4r, r0r, p), sub(p, mulmod(r4i, r0i, p)), p)

            // V = X1 * I
            r4i := addmod(mulmod(r3r, r0i, p), mulmod(r3i, r0r, p), p)
            r4r := addmod(mulmod(r3r, r0r, p), sub(p, mulmod(r3i, r0i, p)), p)

            // Y2 = b.y
            r0i := mload(add(b, 0x40))
            r0r := mload(add(b, 0x60))

            // S2 = Y2 * t0
            r3i := addmod(mulmod(r0r, r2i, p), mulmod(r0i, r2r, p), p)
            r3r := addmod(mulmod(r0r, r2r, p), sub(p, mulmod(r0i, r2i, p)), p)

            // Y1 = a.y
            r0i := mload(add(a, 0x40))
            r0r := mload(add(a, 0x60))

            // t1 = S2 - Y1
            r2i := addmod(r3i, sub(p, r0i), p)
            r2r := addmod(r3r, sub(p, r0r), p)

            // r = 2 * t1
            r3i := mulmod(2, r2i, p)
            r3r := mulmod(2, r2r, p)

            // t2 = r ^ 2
            r2i := mulmod(r3i, r3r, p)
            r2i := addmod(r2i, r2i, p)
            r2r := addmod(mulmod(r3r, r3r, p), sub(p, mulmod(r3i, r3i, p)), p)

            // t4 = t2 - J
            r2i := addmod(r2i, sub(p, r1i), p)
            r2r := addmod(r2r, sub(p, r1r), p)

            // t6 = Y1 * J
            let r5i := addmod(mulmod(r0r, r1i, p), mulmod(r0i, r1r, p), p)
            let r5r := addmod(mulmod(r0r, r1r, p), sub(p, mulmod(r0i, r1i, p)), p)

            // t7 = 2 * t6
            r0i := mulmod(2, r5i, p)
            r0r := mulmod(2, r5r, p)

            // t3 = 2 * V
            r1i := mulmod(2, r4i, p)
            r1r := mulmod(2, r4r, p)

            // X3 = t4 - t3
            r1i := addmod(r2i, sub(p, r1i), p)
            r1r := addmod(r2r, sub(p, r1r), p)

            // result.x = X3
            mstore(result, r1i)
            mstore(add(result, 0x20), r1r)

            // t5 = V - X3
            r1i := addmod(r4i, sub(p, r1i), p)
            r1r := addmod(r4r, sub(p, r1r), p)

            // t8 = r * t5
            r2i := addmod(mulmod(r3r, r1i, p), mulmod(r3i, r1r, p), p)
            r2r := addmod(mulmod(r3r, r1r, p), sub(p, mulmod(r3i, r1i, p)), p)

            // Y3 = t8 - t7
            r0i := addmod(r2i, sub(p, r0i), p)
            r0r := addmod(r2r, sub(p, r0r), p)

            // result.y = Y3
            mstore(add(result, 0x40), r0i)
            mstore(add(result, 0x60), r0r)
        }

        return result;
    }

    /// @notice Return the affine representation of p. The value of p must not be zero, which cannot happen within the
    ///         context this function _is used (i.e. for converting aggregated public keys into affine format).
    /// @param p Jacobian-represented point of which to return affine representation.
    /// @return q affine representation of p
    ///
    /// @dev This is a bit expensive, as it involves a field inversion.
    function _affine(G2PointJacobian memory p) private view returns (G2PointAffine memory q) {
        // Compute 1 / p.z
        (uint256 zi_imag, uint256 zi_real) = _gfp2_inverse(p.z_imag, p.z_real);

        // Compute 1 / (p.z)²
        (uint256 t_imag, uint256 t_real) = _gfp2_square(zi_imag, zi_real);

        // Compute q.x = p.x / (p.z)²
        (q.x_imag, q.x_real) = _gfp2_mul(t_imag, t_real, p.x_imag, p.x_real);

        // Compute 1 / (p.z)³
        (t_imag, t_real) = _gfp2_mul(t_imag, t_real, zi_imag, zi_real);

        // Compute q.y = p.y / (p.z)³
        (q.y_imag, q.y_real) = _gfp2_mul(t_imag, t_real, p.y_imag, p.y_real);

        return q;
    }

    /// @notice Return the jacobian representation of p.
    /// @param p an affine-represented point of which to return Jacobian representation
    /// @return q the Jacobian representation of p
    function _jacobian(G2PointAffine memory p) private pure returns (G2PointJacobian memory q) {
        q.x_imag = p.x_imag;
        q.x_real = p.x_real;
        q.y_imag = p.y_imag;
        q.y_real = p.y_real;
        q.z_imag = 0;
        q.z_real = 1;
        return q;
    }

    /// @notice Expose the the modExp precompile, which computes base^exponent (mod modulus). The fixed value of
    ///         modulus = P is used.
    /// @param base value to be raised to power exponent
    /// @param exponent power to raise base to
    /// @return result (base^exponent) % P
    ///
    /// @dev Loosely based on https://medium.com/@rbkhmrcr/precompiles-solidity-e5d29bd428c4
    function _modExp(uint256 base, uint256 exponent) private view returns (uint256 result) {
        assembly {
            // Get pointer to free memory.
            let ptr := mload(0x40)

            // Prepare arguments for invocation of the bigModExp precompile.
            mstore(ptr, 32) //                  length of base
            mstore(add(ptr, 0x20), 32) //       length of exponent
            mstore(add(ptr, 0x40), 32) //       length of modulus
            mstore(add(ptr, 0x60), base) //     base
            mstore(add(ptr, 0x80), exponent) // exponent
            mstore(add(ptr, 0xa0), P) //        modulus

            // Invoke the bipModExp precompile, revert on failure.
            if iszero(
                staticcall(
                    gas(), // gas cost: no limit
                    0x05, //  precompile address for bigModExp
                    ptr, //  pointer to the input arguments
                    192, //  size of the input arguments (6x uint256)
                    ptr, //  pointer to the where the result is to be stored
                    32 //    size of the result
                )
            ) { revert(0, 0) }

            // Set return value of the function to the result of the bigModExp call.
            result := mload(ptr)
        }
        return result;
    }
}
