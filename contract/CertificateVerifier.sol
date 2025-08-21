// SPDX-License-Identifier: MIT
pragma solidity ^0.7.6;
pragma abicoder v2;

import "./Owned.sol";
import "./libs/DateConverter.sol";
import "./libs/SignatureHelper.sol";
import "./libs/Setn.sol";
import "./interfaces/ICertVerifier.sol";
import "./types/ChunkedX509Cert.sol";
import "./types/Blacklist.sol";

contract CertificateVerifierStorageAccessor is Owned {
    mapping(uint256 => ChunkedX509Cert) rootCerts;
    uint256 rootCertCounter;
    BlackList blacklist;

    function _blacklist() internal view returns (BlackList storage) {
        return blacklist;
    }

    function _rootCert(uint256 rootCertId) internal view returns (ChunkedX509Cert storage) {
        return rootCerts[rootCertId];
    }

    function _setRootCert(uint256 rootCertId, ChunkedX509Cert calldata cert) internal {
        rootCerts[rootCertId] = cert;
    }

}

contract CertificateVerifier is ICertVerifier, CertificateVerifierStorageAccessor {
    using Setn for Bytes32HashSet;
    using Setn for BytesHashSet;

    bytes32 constant CA_HASH = keccak256(hex"30030101ff");

    event SetRootCert(uint256 indexed serialNumber);
    event SerialNumberBlacklisted(bytes indexed serialNumber, bool indexed isBlacklisted);
    event MrEnclaveBlacklisted(bytes32 indexed mrEnclave, bool indexed isBlacklisted);
    event MrSignerBlacklisted(bytes32 indexed mrSigner, bool indexed isBlacklisted);

    modifier onlyExistedRootCert(uint256 rootCertId) {
        require(rootCertId < totalRootCertificates(), "Root certificate is not exist");
        _;
    }

    function setRootCert(ChunkedX509Cert calldata cert) external onlyOwner() {
        uint256 id = rootCertCounter;
        require(keccak256(cert.ca) == CA_HASH, "Root certificate cannot be non-intermediate");
        verifyCert(cert, cert.publicKey);

        _setRootCert(id, cert);
        rootCertCounter++;

        emit SetRootCert(id);
    }

    function setCertsBlacklist(
        BytesBlacklist[] calldata serialNumberChanges,
        Bytes32Blacklist[] calldata mrEnclaveChanges,
        Bytes32Blacklist[] calldata mrSignerChanges
    ) external onlyOwner() {
        BlackList storage blacklist = _blacklist();

        for (uint256 i; i < serialNumberChanges.length; i++) {
            BytesBlacklist calldata blacklistItem = serialNumberChanges[i];
            if (blacklistItem.isBlacklisted) {
                blacklist.bySerialNumber.add(blacklistItem.item);
            } else {
                blacklist.bySerialNumber.remove(blacklistItem.item);
            }
            emit SerialNumberBlacklisted(blacklistItem.item, blacklistItem.isBlacklisted);
        }

        for (uint256 i; i < mrEnclaveChanges.length; i++) {
            Bytes32Blacklist calldata blacklistItem = mrEnclaveChanges[i];
            if (blacklistItem.isBlacklisted) {
                blacklist.byMrEnclave.add(blacklistItem.item);
            } else {
                blacklist.byMrEnclave.remove(blacklistItem.item);
            }
            emit MrEnclaveBlacklisted(blacklistItem.item, blacklistItem.isBlacklisted);
        }

        for (uint256 i; i < mrSignerChanges.length; i++) {
            Bytes32Blacklist calldata blacklistItem = mrSignerChanges[i];
            if (blacklistItem.isBlacklisted) {
                blacklist.byMrSigner.add(blacklistItem.item);
            } else {
                blacklist.byMrSigner.remove(blacklistItem.item);
            }
            emit MrSignerBlacklisted(blacklistItem.item, blacklistItem.isBlacklisted);
        }
    }

    function totalRootCertificates() public view returns (uint256) {
        return rootCertCounter;
    }

    function getRootCert(uint256 rootCertId) public view onlyExistedRootCert(rootCertId) returns (ChunkedX509Cert memory) {
        return _rootCert(rootCertId);
    }

    function isCertBlacklisted(bytes memory serialNumber, bytes32 mrEnclave, bytes32 mrSigner) public view returns (bool) {
        BlackList storage blacklist = _blacklist();

        if (blacklist.bySerialNumber.isExists(serialNumber)) {
            return true;
        }
        if (mrEnclave != 0 && blacklist.byMrEnclave.isExists(mrEnclave)) {
            return true;
        }
        if (mrSigner != 0 && blacklist.byMrSigner.isExists(mrSigner)) {
            return true;
        }
        return false;
    }

    function getCertsBlacklist() public view returns (bytes[] memory bySerialNumber, bytes32[] memory byMrEnclave, bytes32[] memory byMrSigner) {
        BlackList storage blacklist = _blacklist();
        return (blacklist.bySerialNumber.items, blacklist.byMrEnclave.items, blacklist.byMrSigner.items);
    }

    function verifyCert(ChunkedX509Cert calldata cert, bytes memory caPubKey) public view returns (bytes32 certHash) {
        require(!isCertBlacklisted(cert.serialNumber, cert.mrEnclave, cert.mrSigner), "Certificate is blacklisted");
        require(DateConverter.convertDateToUnixTimestamp(cert.expirationDate) > block.timestamp, "Certificate is expired");

        bytes memory packedData = abi.encodePacked(
            cert.nonSerializedParts[1],
            cert.serialNumber,
            cert.nonSerializedParts[2],
            cert.expirationDate,
            cert.nonSerializedParts[3]
        );
        packedData = abi.encodePacked(packedData, cert.publicKey, cert.nonSerializedParts[4], cert.ca, cert.nonSerializedParts[5]);

        if (cert.userData != 0) {
            packedData = abi.encodePacked(packedData, cert.userData, cert.nonSerializedParts[6]);
        }
        if (cert.mrEnclave != 0) {
            packedData = abi.encodePacked(packedData, cert.mrEnclave);
        }
        if (cert.mrSigner != 0) {
            packedData = abi.encodePacked(packedData, cert.nonSerializedParts[7], cert.mrSigner);
        }

        return SignatureHelper.checkCertSignature(cert.signature, caPubKey, packedData);
    }

    function verifyCertChain(
        ChunkedX509Cert[] calldata certsChain,
        uint256 rootCertId
    ) public view override onlyExistedRootCert(rootCertId) returns (bytes32 leafCertHash) {
        require(certsChain.length != 0, "certsChain cannot be empty array");
        ChunkedX509Cert storage rootCert = _rootCert(rootCertId);
        require(!isCertBlacklisted(rootCert.serialNumber, rootCert.mrEnclave, rootCert.mrSigner), "Root certificate is blacklisted");
        require(DateConverter.convertDateToUnixTimestamp(rootCert.expirationDate) > block.timestamp, "Root certificate is expired");

        bytes memory rootCertPubKey = rootCert.publicKey;
        bytes memory caPubKey;
        for (uint256 i; i < certsChain.length; i++) {
            if (i != certsChain.length - 1) {
                ChunkedX509Cert calldata cert = certsChain[i + 1];
                caPubKey = cert.publicKey;
                require(keccak256(cert.ca) == CA_HASH, "Intermediate certificate does not have ca flag");
            } else {
                caPubKey = rootCertPubKey;
            }

            if (i == 0) {
                leafCertHash = verifyCert(certsChain[i], caPubKey);
            } else {
                verifyCert(certsChain[i], caPubKey);
            }
        }
    }
}

