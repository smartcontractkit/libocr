// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/*
 * @title ConfigDigestUtil
 * @author Michael Fletcher
 * @notice ConfigDigest related utility functions
 */
library ConfigDigestUtil {

    function configDigestFromConfigData(
        uint256 chainId,
        address contractAddress,
        uint64 configCount,
        address[] memory signers,
        address[] memory transmitters,
        uint8 f,
        bytes memory onchainConfig,
        uint64 offchainConfigVersion,
        bytes memory offchainConfig
    ) internal pure returns (bytes32)
    {
        uint256 hash = uint256(
            keccak256(
                abi.encode(
                    chainId,
                    contractAddress,
                    configCount,
                    signers,
                    transmitters,
                    f,
                    onchainConfig,
                    offchainConfigVersion,
                    offchainConfig
        )));

        uint256 prefixMask = type(uint256).max << (256-16); // 0xFFFF00..00
        uint256 prefix = 0x0001 << (256-16); // 0x000100..00

        return bytes32((prefix & prefixMask) | (hash & ~prefixMask));
    }

}