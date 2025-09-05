// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

interface IAdminCertificateHelper {
    function isAdmin(address caller) external view returns(bool);
}
