// SPDX-License-Identifier: MIT
pragma solidity ^0.7.6;
pragma abicoder v2;

import "../types/ChunkedX509Cert.sol";

interface ICertVerifier {
    function verifyCertChain(ChunkedX509Cert[] calldata certsChain, uint256 rootCertId) external view returns (bytes32 leafCertHash);
}
