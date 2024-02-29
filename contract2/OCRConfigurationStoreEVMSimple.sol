// SPDX-License-Identifier: MIT
pragma solidity =0.8.19;

import "./interfaces/TypeAndVersionInterface.sol";
import "./lib/ConfigDigestUtilEVMSimple.sol";
import "./OwnerIsCreator.sol";
import "./OCR2Abstract.sol";

/// @title OCRConfigurationStoreEVMSimple
/// @notice This contract stores configurations for protocol versions OCR2 and
/// above in contract storage. It uses the "EVMSimple" config digester.
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

    /// @notice a list of configurations keyed by their digest
    mapping(bytes32 => ConfigurationEVMSimple) internal s_configurations;

    /// @notice emitted when a new configuration is added
    event NewConfiguration(bytes32 indexed configDigest);

    /// @notice adds a new configuration to the store
    function addConfig(ConfigurationEVMSimple calldata configuration) external returns (bytes32) {

        bytes32 configDigest = ConfigDigestUtilEVMSimple.configDigestFromConfigData(
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

    /// @notice reads a configuration from the store
    function readConfig(bytes32 configDigest) external view returns (ConfigurationEVMSimple memory) {
        return s_configurations[configDigest];
    }

    /// @inheritdoc TypeAndVersionInterface
    function typeAndVersion() external override pure virtual returns (string memory)
    {
        return "OCRConfigurationStoreEVMSimple 1.0.0";
    }
}
