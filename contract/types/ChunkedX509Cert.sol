// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

struct ChunkedX509Cert {
    bytes[] nonSerializedParts;
    bytes expirationDate;
    bytes ca; // certificate authority
    bytes32 userData;
    bytes publicKey;
    bytes serialNumber;
    bytes32 mrEnclave;
    bytes32 mrSigner;
    bytes signature;
}
