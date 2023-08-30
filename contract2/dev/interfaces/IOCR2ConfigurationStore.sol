// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;


/**
 * @title IOCR2ConfigurationStore
 * @author Michael Fletcher
 * @notice This contract is used for storing and retrieving OCR related configurations.
 */
interface IOCR2ConfigurationStore{

    struct ExtendedConfiguration {
        uint32 blockNumber;
        address contractAddress;
        bytes32 configDigest;
        Configuration configuration;
    }

    struct Configuration {
        uint64 configCount;
        address[] signers;
        address[] transmitters;
        bytes onchainConfig;
        bytes offchainConfig;
        uint64 offchainConfigVersion;
        uint8 f;
    }

    /**
     * @notice Adds a configuration to the list of OCR configurations stored within contract state
     * @param configurationParams The struct containing the configuration to be added
     * @return The digest of the configuration that was added
     */
    function addConfig(Configuration calldata configurationParams) external returns (bytes32);

    /**
     * @notice Returns the configuration stored at the specified index
     * @param configDigest The id of the configuration to be returned
     * @return The configuration stored at the specified index
     */
    function readConfig(bytes32 configDigest) external view returns (ExtendedConfiguration memory);

    /**
     * @notice Returns the latest configuration stored for a particular contract address
     * @param contractAddress The address of the contract to return the latest configuration for
     * @return The latest configuration stored for a particular contract address
     */
    function latestConfig(address contractAddress) external view returns (ExtendedConfiguration memory);
}
