// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;

interface ITransmitterCertificateHelper {
    function isTransmitterInitialized(address transmitter) external view returns(bool);
}
