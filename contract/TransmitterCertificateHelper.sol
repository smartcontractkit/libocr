// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma abicoder v2;

import "./Owned.sol";
import "./interfaces/ICertVerifier.sol";
import "./interfaces/ITransmitterCertificateHelper.sol";
import "@openzeppelin/contracts/cryptography/ECDSA.sol";

contract TransmitterCertificateHelper is ITransmitterCertificateHelper, Owned {

  struct CertData {
    bool initialized;
    ChunkedX509Cert[] certsChain;
    uint256 rootCertId;
  }

  ICertVerifier internal s_certVerifier;
  mapping (address /* transmitter address */ => CertData) internal s_certData;

  event CertVerifierSet(ICertVerifier indexed oldCertVerifier, ICertVerifier indexed newCertVerifier);
  event InitializeTransmitter(address indexed transmitter, ChunkedX509Cert[] certsChain, uint256 rootCertId);

  constructor(ICertVerifier certVerifier){
    _setCertVerifier(certVerifier);
  }

  function setCertVerifier(ICertVerifier certVerifier) external onlyOwner() {
    _setCertVerifier(certVerifier);
  }

  function initializeTransmitter(ChunkedX509Cert[] calldata certsChain, uint256 rootCertId,bytes calldata signature) external {
    s_certVerifier.verifyCertChain(certsChain, rootCertId);

    address transmitter=msg.sender;
    bytes32 dataHash=keccak256(abi.encode(transmitter));
    address signer=ECDSA.recover(dataHash, signature);
    require(signer==address(uint160(uint256(keccak256(certsChain[0].publicKey)))),"Invalid transmitter signature");
  
    CertData storage certData=s_certData[transmitter];
    certData.initialized = true;
    certData.rootCertId = rootCertId;
    delete certData.certsChain;
    for(uint256 i=0;i<certsChain.length;i++){
      certData.certsChain.push(certsChain[i]);
    }
    emit InitializeTransmitter(transmitter, certsChain, rootCertId);
  }

  function isTransmitterInitialized(address transmitter) external view override returns(bool){
    return s_certData[transmitter].initialized;
  }

  function getCertData(address transmitter) external view returns(CertData memory){
    return s_certData[transmitter];
  }

  function getCertVerifier() external view returns(ICertVerifier){
    return s_certVerifier;
  }

  function _setCertVerifier(ICertVerifier certVerifier) internal {
    require(address(certVerifier) != address(0),"certVerifier cannot be zero address");
    emit CertVerifierSet(s_certVerifier, certVerifier);
    s_certVerifier = certVerifier;
  }

}
