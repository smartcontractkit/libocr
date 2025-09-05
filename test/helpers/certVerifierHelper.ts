import { ethers } from 'hardhat';
import { BigNumberish, BytesLike } from "ethers"
import { CertificateVerifier } from '../../typechain-types';
import { ChunkedX509CertStruct, BytesBlacklistStruct } from '../../typechain-types/contract/CertificateVerifier.sol/CertificateVerifier';
import { HardhatEthersSigner } from '@nomicfoundation/hardhat-ethers/signers';

export class CertVerifierHelper {
    contract: CertificateVerifier;
    address: string;
    deployer: HardhatEthersSigner;

    constructor(contract: CertificateVerifier, deployer: HardhatEthersSigner, address: string) {
        this.contract = contract;
        this.deployer = deployer;
        this.address = address
    }

    static async create(certVerifierAddress?: string): Promise<CertVerifierHelper> {
        let contract;
        if (!certVerifierAddress) {
            const factory = await ethers.getContractFactory('CertificateVerifier');
            contract = await factory.deploy()
            certVerifierAddress = await contract.getAddress()
        } else {
            contract = await ethers.getContractAt("CertificateVerifier", certVerifierAddress)
        }

        const deployer = (await ethers.getSigners())[0];
        return new CertVerifierHelper(contract, deployer, certVerifierAddress);
    }

    async setRootCert(cert: ChunkedX509CertStruct, signer: HardhatEthersSigner = this.deployer) {
        return this.contract.connect(signer).setRootCert(cert);
    }

    async setCertsBlacklist(
        serialNumberChanges: BytesBlacklistStruct[],
        mrEnclaveChanges: BytesBlacklistStruct[],
        mrSignerChanges: BytesBlacklistStruct[],
        signer: HardhatEthersSigner = this.deployer
    ) {
        return this.contract.connect(signer).setCertsBlacklist(serialNumberChanges, mrEnclaveChanges, mrSignerChanges);
    }

    async getCertsBlacklist(): Promise<{ bySerialNumber: string[]; byMrEnclave: string[]; byMrSigner: string[] }> {
        return await this.contract.getCertsBlacklist();
    }

    async totalRootCertificates(): Promise<bigint> {
        return await this.contract.totalRootCertificates();
    }

    async verifyCert(cert: ChunkedX509CertStruct, rootPubKey: BytesLike): Promise<string> {
        return await this.contract.verifyCert(cert, rootPubKey);
    }

    async verifyCertChain(certsChain: ChunkedX509CertStruct[], rootCertId: number): Promise<string> {
        return await this.contract.verifyCertChain(certsChain, rootCertId);
    }

    async getRootCert(rootCertId: BigNumberish): Promise<ChunkedX509CertStruct> {
        return await this.contract.getRootCert(rootCertId);
    }
}
