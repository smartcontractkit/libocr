// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/TypeAndVersionInterface.sol";
import "./lib/ConfigDigestUtil.sol";
import "./OwnerIsCreator.sol";
import "./OCR2Abstract.sol";

contract OCRConfigurationStoreEVMSimple is TypeAndVersionInterface {

    struct ConfigurationEVMSimple {
        address[] signers;
        address[] transmitters;
        bytes onchainConfig;
        bytes offchainConfig;
        address contractAddress;
        uint64 offchainConfigVersion;
        uint32 configCount;
        uint8 f;
    }

    // @notice a list of configurations keyed by their digest
    mapping(bytes32 => ConfigurationEVMSimple) internal s_configurations;

    // @notice emitted when a new configuration is added
    event NewConfiguration(bytes32 indexed configDigest);

    // @inheritdoc IOCRConfigurationStoreEVMSimple
    function addConfig(ConfigurationEVMSimple calldata configuration) external returns (bytes32) {

        bytes32 configDigest = ConfigDigestUtil.configDigestFromConfigData(
            block.chainid,
            configuration.contractAddress,
            configuration.configCount,
            configuration.signers,
            configuration.transmitters,
            configuration.f,
            configuration.onchainConfig,
            configuration.offchainConfigVersion,
            configuration.offchainConfig
        );

        s_configurations[configDigest] = configuration;

        emit NewConfiguration(configDigest);

        return configDigest;
    }

    // @inheritdoc IOCRConfigurationStoreEVMSimple
    function readConfig(bytes32 configDigest) external view returns (ConfigurationEVMSimple memory) {
        return s_configurations[configDigest];
    }

    function typeAndVersion() external override pure virtual returns (string memory)
    {
        return "OCRConfigurationStoreEVMSimple 1.0.0";
    }
}
