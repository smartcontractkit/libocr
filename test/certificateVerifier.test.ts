import { expect } from 'chai';
import { BlacklistHelper, CertificateHelper } from './helpers/certificateHelper';
import { CertVerifierHelper } from './helpers/certVerifierHelper';
import { BytesBlacklistStruct } from '../typechain-types/CertificateVerifier.sol/CertificateVerifier';
import { ethers } from 'hardhat';
import { Result } from "ethers"
import { time, takeSnapshot, SnapshotRestorer } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe('Certificate verifier', function () {
    let certVerifier: CertVerifierHelper;
    let snapshot: SnapshotRestorer;

    before(async function () {
        certVerifier = await CertVerifierHelper.create();
    });

    beforeEach(async function () {
        snapshot = await takeSnapshot()
    });

    afterEach(async function () {
        await snapshot.restore()
    });

    describe('setRootCert', function () {
        it('Only diamond owner can setRootCert', async function () {
            const certificate = await CertificateHelper.gen();
            const randomUser = (await ethers.getSigners())[10]
            await expect(certVerifier.setRootCert(certificate.cert, randomUser)).to.be.rejectedWith('Only callable by owner');
        });

        it('Root certificate must be intermediate', async function () {
            const certificate = await CertificateHelper.gen();
            certificate.cert.ca = '0x0000000000';
            await expect(certVerifier.setRootCert(certificate.cert)).to.be.rejectedWith('Root certificate cannot be non-intermediate');
        });

        it('Root certificate must be valid', async function () {
            const certificate = await CertificateHelper.gen();
            const expirationDate = certificate.cert.expirationDate;
            certificate.cert.expirationDate = CertificateHelper.encodeDate(0);
            await expect(certVerifier.setRootCert(certificate.cert)).to.be.rejectedWith('Certificate is expired');

            const signer = await CertificateHelper.genRandomWallet();
            const { signature: invalidSignature } = await CertificateHelper.signCert(certificate.cert, signer);
            certificate.cert.expirationDate = expirationDate;
            certificate.cert.signature = invalidSignature;
            await expect(certVerifier.setRootCert(certificate.cert)).to.be.rejectedWith('Invalid signature');
        });

        it('setRootCert update storage', async function () {
            const certificate = await CertificateHelper.gen();
            const id = await certVerifier.totalRootCertificates();
            await expect(certVerifier.setRootCert(certificate.cert))
                .to.emit(certVerifier.contract, 'SetRootCert')
                .withArgs(id);

            expect(((await certVerifier.getRootCert(id)) as any as Result).toObject()).to.deep.equal(certificate.cert);
            expect(await certVerifier.totalRootCertificates()).to.equal(id + 1n);
        });
    });

    describe('setCertsBlacklist', function () {
        it('Only cert owner can setCertsBlacklist', async function () {
            const randomUser = (await ethers.getSigners())[10]
            await expect(certVerifier.setCertsBlacklist([], [], [], randomUser)).to.be.rejectedWith('Only callable by owner');
        });

        it('setCertsBlacklist update storage', async function () {
            const bySerialNumberBL = BlacklistHelper.gen();
            const byMrEnclaveBL = BlacklistHelper.gen();
            const byMrSignerBL = BlacklistHelper.gen();

            const tx = await certVerifier.setCertsBlacklist(bySerialNumberBL, byMrEnclaveBL, byMrSignerBL);
            await tx.wait();

            const { bySerialNumber, byMrEnclave, byMrSigner } = await certVerifier.getCertsBlacklist();

            async function checkBlacklist(toSet: BytesBlacklistStruct[], fromGetter: string[], eventName: string) {
                let count = 0;
                for (let i = 0; i < toSet.length; i++) {
                    await expect(tx).to.emit(certVerifier.contract, eventName).withArgs(toSet[i].item, toSet[i].isBlacklisted);

                    if (toSet[i].isBlacklisted) {
                        expect(fromGetter[count]).to.equal(toSet[i].item);
                        count++;
                    }
                }
                expect(count).to.equal(fromGetter.length);
            }
            await checkBlacklist(bySerialNumberBL, bySerialNumber, 'SerialNumberBlacklisted');
            await checkBlacklist(byMrEnclaveBL, byMrEnclave, 'MrEnclaveBlacklisted');
            await checkBlacklist(byMrSignerBL, byMrSigner, 'MrSignerBlacklisted');
        });
    });

    describe('verifyCert', function () {
        it('Revert if current time is more than cert expiration date', async function () {
            const certificate = await CertificateHelper.gen();
            await time.increaseTo(CertificateHelper.decodeDate(certificate.cert.expirationDate));
            await expect(certVerifier.verifyCert(certificate.cert, certificate.cert.publicKey)).to.be.rejectedWith('Certificate is expired');
        });

        it('Revert if signature is invalid', async function () {
            const certificate = await CertificateHelper.gen();
            const signer = CertificateHelper.genRandomWallet();
            await expect(certVerifier.verifyCert(certificate.cert, signer.publicKey)).to.be.rejectedWith('Invalid signature');
            certificate.cert.signature = '0x';
            await expect(certVerifier.verifyCert(certificate.cert, certificate.cert.publicKey)).to.be.rejectedWith('Invalid signature');
        });

        it('Successful verification certificate', async function () {
            const certificate = await CertificateHelper.gen();
            const certHash = await certVerifier.verifyCert(certificate.cert, certificate.cert.publicKey);
            expect(certHash).to.equal(certificate.certHash);
        });

        it('Successful verification examples sertificates', async function () {
            const rootCert = (await import('./certsExamples/rootCert.json')).default;
            const subroot1Cert = (await import('./certsExamples/subroot1Cert.json')).default;
            const subroot2Cert = (await import('./certsExamples/subroot2Cert.json')).default;
            const wrongSignatureCert = (await import('./certsExamples/wrongSignatureCert.json')).default;

            await certVerifier.verifyCert(rootCert, rootCert.publicKey);
            await certVerifier.verifyCert(subroot1Cert, rootCert.publicKey);
            await certVerifier.verifyCert(subroot2Cert, subroot1Cert.publicKey);
            await certVerifier.setRootCert(rootCert);
            await expect(certVerifier.verifyCertChain([subroot2Cert, subroot1Cert], 0)).to.rejectedWith(
                'Intermediate certificate does not have ca flag'
            );
            await expect(certVerifier.verifyCert(wrongSignatureCert, wrongSignatureCert.publicKey)).to.rejectedWith('Invalid signature');
        });
    });

    describe('verifyCertChain', function () {
        let certsCash: CertificateHelper[] = [];
        let certs: CertificateHelper[] = [];
        let snapshot2: SnapshotRestorer;
        before(async function () {
            const certificate = await CertificateHelper.gen();
            const certificate2 = await CertificateHelper.gen(certificate);
            const certificate3 = await CertificateHelper.gen(certificate2);
            const certificate4 = await CertificateHelper.gen(certificate3);
            certsCash = [certificate4, certificate3, certificate2, certificate];
            await certVerifier.setRootCert(certificate.cert);
        });
        beforeEach(async function () {
            certs = certsCash.map(x => {
                return { ...x };
            });
            snapshot2 = await takeSnapshot()
        });

        afterEach(async function () {
            await snapshot2.restore()
        });

        async function checkVerify(certs: CertificateHelper[]) {
            const certsChain = certs.map(x => x.cert);
            const certHash = await certVerifier.verifyCertChain(certsChain, 0);
            expect(certHash).to.equal(certs[0].certHash);
        }

        it('Revert if certsChain is empty array', async function () {
            await expect(certVerifier.verifyCertChain([], 0)).to.be.rejectedWith('certsChain cannot be empty array');
        });

        it('Revert if root cert is not exist on id', async function () {
            await expect(certVerifier.verifyCertChain([certs[0].cert], 2)).to.be.rejectedWith('Root certificate is not exist');
        });

        it('Revert if root cert is expired', async function () {
            await time.increaseTo(CertificateHelper.decodeDate(certs[0].cert.expirationDate));
            await expect(certVerifier.verifyCertChain([certs[0].cert], 0)).to.be.rejectedWith('Root certificate is expired');
        });

        it('Revert if intermediate cert does not have ca flag', async function () {
            certs[1].cert = { ...certs[1].cert, ca: '0x0000000000' };
            await expect(checkVerify(certs)).to.be.rejectedWith('Intermediate certificate does not have ca flag');
        });

        it('Revert if one of the cert is blacklisted', async function () {
            const bySerialNumber: BytesBlacklistStruct[] = [{ item: certs[2].cert.serialNumber, isBlacklisted: false }];
            const byMrEnclave: BytesBlacklistStruct[] = [{ item: certs[2].cert.mrEnclave, isBlacklisted: false }];
            const byMrSigner: BytesBlacklistStruct[] = [{ item: certs[2].cert.mrSigner, isBlacklisted: false }];
            const blacklists = [bySerialNumber, byMrEnclave, byMrSigner];
            for (let i = 0; i < 3; i++) {
                blacklists.map(x => (x[0].isBlacklisted = false));
                blacklists[i][0].isBlacklisted = true;
                const tx = await certVerifier.setCertsBlacklist(bySerialNumber, byMrEnclave, byMrSigner);
                await tx.wait();
                await expect(checkVerify(certs)).to.be.rejectedWith('Certificate is blacklisted');
            }

            // latest root certificate is not necessary in function argument
            certs.pop();
            await expect(checkVerify(certs)).to.be.rejectedWith('Certificate is blacklisted');
        });

        it('Successful verification of certificates', async function () {
            await checkVerify(certs);

            // latest root certificate is not necessary in function argument
            certs.pop();
            await checkVerify(certs);

            await checkVerify(certs.slice(-2));
            await checkVerify(certs.slice(-1));
        });
    });
});
