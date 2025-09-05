// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

import "./Owned.sol";
import "./interfaces/ICertVerifier.sol";
import "./interfaces/IAdminCertificateHelper.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract AdminCertificateHelper is IAdminCertificateHelper, Owned {

  ICertVerifier internal certVerifier;

  address internal admin;
  ChunkedX509Cert[] internal certsChain;
  uint256 internal rootCertId;
  uint256 internal signatureNonce;
  bytes32 internal authorizedSolutionHash;

  event CertVerifierSet(ICertVerifier indexed oldCertVerifier, ICertVerifier indexed newCertVerifier);
  event SetAdmin(address indexed oldAdmin, address indexed newAdmin, ChunkedX509Cert[] certsChain, uint256 rootCertId);
  event SetAuthorizedSolutionHash(bytes32 indexed oldAuthorizedSolutionHash, bytes32 indexed newAuthorizedSolutionHash);

  constructor(ICertVerifier _certVerifier,bytes32 _authorizedSolutionHash){
    _setCertVerifier(_certVerifier);
    _setAuthorizedSolutionHash(_authorizedSolutionHash);
  }

  function setCertVerifier(ICertVerifier _certVerifier) external onlyOwner() {
    _setCertVerifier(_certVerifier);
  }
  
  function setAuthorizedSolutionHash(bytes32 _authorizedSolutionHash) external onlyOwner() {
    _setAuthorizedSolutionHash(_authorizedSolutionHash);
  }

  function addAdmin(ChunkedX509Cert[] calldata _certsChain, uint256 _rootCertId,bytes calldata signature) external {
    certVerifier.verifyCertChain(_certsChain, _rootCertId);
    require(_certsChain[0].userData==authorizedSolutionHash,"Only authorized solution");

    address newAdmin=msg.sender;
    bytes32 dataHash=sha256(abi.encode(newAdmin,block.chainid,address(this),signatureNonce));
    address signer=ECDSA.recover(dataHash, signature);
    require(signer==address(uint160(uint256(keccak256(_certsChain[0].publicKey)))),"Invalid data signature");
    signatureNonce++;

    address oldAdmin=admin;
    admin=newAdmin;
    certsChain=_certsChain;
    rootCertId=_rootCertId;
    emit SetAdmin(oldAdmin,newAdmin, _certsChain, _rootCertId);
  }

  function isAdmin(address caller) external view override returns(bool){
    return admin==caller;
  }

  function getAdminData() external view returns(address _admin,ChunkedX509Cert[] memory _certsChain, uint256 _rootCertId){
    return (admin,certsChain,rootCertId);
  }

  function getCertVerifier() external view returns(ICertVerifier){
    return certVerifier;
  }

  function getAuthorizedSolutionHash() external view returns(bytes32){
    return authorizedSolutionHash;
  }

  function getSignatureNonce() external view returns(uint256){
    return signatureNonce;
  }

  function _setCertVerifier(ICertVerifier _certVerifier) internal {
    require(address(_certVerifier) != address(0),"certVerifier cannot be zero address");
    emit CertVerifierSet(certVerifier, _certVerifier);
    certVerifier = _certVerifier;
  }

  function _setAuthorizedSolutionHash(bytes32 _authorizedSolutionHash) internal {
    require(_authorizedSolutionHash != 0,"AuthorizedSolutionHash cannot be zero");
    emit SetAuthorizedSolutionHash(authorizedSolutionHash, _authorizedSolutionHash);
    authorizedSolutionHash = _authorizedSolutionHash;
  }

}
