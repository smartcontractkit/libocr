// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/IOCR2ConfigurationStore.sol";
import "./interfaces/TypeAndVersionInterface.sol";
import "./lib/ConfigDigestUtil.sol";
import "./OwnerIsCreator.sol";
import "./OCR2Abstract.sol";

contract OCR2ConfigurationStore is IOCR2ConfigurationStore, OwnerIsCreator, TypeAndVersionInterface {

    // @notice a list of configurations keyed by their digest
    mapping(bytes32 => ExtendedConfiguration) public s_configurations;

    // @notice keep track of the latest config for each contract
    mapping(address => bytes32) public s_latestConfigurationDigest;

    // @inheritdoc IOCR2ConfigurationStore
    function addConfig(Configuration calldata configuration) external override returns (bytes32) {

        //calculate the digest
        bytes32 configDigest = ConfigDigestUtil.configDigestFromConfigData(
            block.chainid,
            msg.sender,
            configuration.configCount,
            configuration.signers,
            configuration.transmitters,
            configuration.f,
            configuration.onchainConfig,
            configuration.offchainConfigVersion,
            configuration.offchainConfig
        );

        //create the extended configuration containing all the information to calculate a digest
        s_configurations[configDigest] = ExtendedConfiguration({
            blockNumber: uint32(block.number),
            contractAddress: msg.sender,
            configDigest: configDigest,
            configuration: configuration
        });

        s_latestConfigurationDigest[msg.sender] = configDigest;

        return configDigest;
    }

    // @inheritdoc IOCR2ConfigurationStore
    function readConfig(bytes32 configDigest) external view override returns (ExtendedConfiguration memory) {
        return s_configurations[configDigest];
    }

    // @inheritdoc IOCR2ConfigurationStore
    function latestConfig(address contractAddress) external view override returns (ExtendedConfiguration memory) {
        return s_configurations[s_latestConfigurationDigest[contractAddress]];
    }

    function typeAndVersion() external override pure virtual returns (string memory)
    {
        return "OCR2ConfigurationStore 1.0.0";
    }
}
