// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/ChunkedX509Cert.sol";

interface ICertVerifier {
    function verifyCertChain(ChunkedX509Cert[] calldata certsChain, uint256 rootCertId) external view returns (bytes32 leafCertHash);
}
