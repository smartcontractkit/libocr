import { expect } from "chai";
import { ethers } from "hardhat";
import { Signer } from "ethers";
import { takeSnapshot, SnapshotRestorer } from "@nomicfoundation/hardhat-toolbox/network-helpers";
import { CertificateHelper } from './helpers/certificateHelper';
import { CertVerifierHelper } from './helpers/certVerifierHelper';
import { AdminCertificateHelper } from "../typechain-types";

describe("AdminCertificateHelper", function () {
  let transmitter: Signer, deployer: Signer;
  let adminCertHelper: AdminCertificateHelper;
  let certVerifier: CertVerifierHelper;
  let snapshot: SnapshotRestorer;
  let certsChain: CertificateHelper[]
  let chainId: bigint;

  before(async function () {
    chainId = (await ethers.provider?.getNetwork())?.chainId;

    [deployer, transmitter] = await ethers.getSigners();

    certVerifier = await CertVerifierHelper.create();

    const certificate = await CertificateHelper.gen();
    const certificate2 = await CertificateHelper.gen(certificate);
    certsChain = [certificate2, certificate];

    const factory = await ethers.getContractFactory("AdminCertificateHelper");
    adminCertHelper = await factory.deploy(certVerifier.address, certificate2.cert.userData);

    await certVerifier.setRootCert(certificate.cert);
  });

  beforeEach(async function () {
    snapshot = await takeSnapshot()
  });

  afterEach(async function () {
    await snapshot.restore()
  });


  describe("Constructor", function () {
    it("Set certVerifier", async function () {
      const currentVerifier = await adminCertHelper.getCertVerifier();
      expect(currentVerifier).to.equal(certVerifier.address);
    });

    it("Emit event CertVerifierSet", async function () {
      await expect(adminCertHelper.deploymentTransaction())
        .to.emit(adminCertHelper, "CertVerifierSet")
        .withArgs(ethers.ZeroAddress, certVerifier.address);
    });
  });

  describe("setAuthorizedSolutionHash", function () {
    it("Only owner can setAuthorizedSolutionHash", async function () {

      const newSolutionHash = certsChain[1].cert.userData
      const oldSolutionHash = certsChain[0].cert.userData

      await expect(
        adminCertHelper.connect(transmitter).setAuthorizedSolutionHash(newSolutionHash)
      ).to.be.revertedWith("Only callable by owner");

      await expect(adminCertHelper.setAuthorizedSolutionHash(newSolutionHash))
        .to.emit(adminCertHelper, "SetAuthorizedSolutionHash")
        .withArgs(oldSolutionHash, newSolutionHash);

      expect(await adminCertHelper.getAuthorizedSolutionHash()).to.equal(newSolutionHash);
    });

    it("Should revert if new AuthorizedSolutionHash is zero", async function () {
      await expect(
        adminCertHelper.setAuthorizedSolutionHash(ethers.ZeroHash)
      ).to.be.revertedWith("AuthorizedSolutionHash cannot be zero");
    });
  });

  describe("setCertVerifier", function () {
    it("Only owner can setCertVerifier", async function () {

      const newCertVerifierAddress = await transmitter.getAddress()

      await expect(
        adminCertHelper.connect(transmitter).setCertVerifier(newCertVerifierAddress)
      ).to.be.revertedWith("Only callable by owner");

      await expect(adminCertHelper.setCertVerifier(newCertVerifierAddress))
        .to.emit(adminCertHelper, "CertVerifierSet")
        .withArgs(certVerifier.address, newCertVerifierAddress);

      expect(await adminCertHelper.getCertVerifier()).to.equal(newCertVerifierAddress);
    });

    it("Should revert if new cert verifier address is zero address", async function () {
      await expect(
        adminCertHelper.setCertVerifier(ethers.ZeroAddress)
      ).to.be.revertedWith("certVerifier cannot be zero address");
    });
  });

  describe("addAdmin", function () {
    it("success addAdmin", async function () {
      const transmitterAddress = await transmitter.getAddress();

      const dataHash = ethers.sha256(ethers.AbiCoder.defaultAbiCoder().encode(["address", "uint256", "address", "uint256"], [transmitterAddress, chainId, await adminCertHelper.getAddress(), 0]));
      const signature = certsChain[0].signer.signingKey.sign(dataHash)
      const rawSignature = signature.r + signature.s.slice(2) + signature.v.toString(16);

      const tx = await adminCertHelper.connect(transmitter).addAdmin(
        certsChain.map(x => x.cert),
        0,
        rawSignature
      )
      const txReceipt = await tx.wait()

      expect(await adminCertHelper.getSignatureNonce()).to.equal(1)
      const isInitialized = await adminCertHelper.isAdmin(transmitterAddress);
      expect(isInitialized).to.be.true;

      const events = txReceipt!.logs.map(log => {
        const parsedLog = adminCertHelper.interface.parseLog(log)
        if (parsedLog?.name === "SetAdmin") { return parsedLog }
      });

      expect(events.length).to.be.equal(1);

      const eventArgs = events[0]!.args
      const expectedEventArgs = {
        oldAdmin: ethers.ZeroAddress,
        newAdmin: transmitterAddress,
        certsChain: certsChain.map(x => x.cert),
        rootCertId: 0
      }
      expect(eventArgs.toObject(true)).to.deep.equal(expectedEventArgs);

      let certData = await adminCertHelper.getAdminData();
      const expectedCertData = {
        _admin: transmitterAddress,
        _certsChain: certsChain.map(x => x.cert),
        _rootCertId: 0
      }
      expect(certData.toObject(true)).to.deep.equal(expectedCertData);

      let rawSignature2;
      {
        const dataHash = ethers.sha256(ethers.AbiCoder.defaultAbiCoder().encode(["address", "uint256", "address", "uint256"], [transmitterAddress, chainId, await adminCertHelper.getAddress(), 1]));
        const signature = certsChain[0].signer.signingKey.sign(dataHash)
        rawSignature2 = signature.r + signature.s.slice(2) + signature.v.toString(16);
      }
      // success reinitialize trasmitter
      await adminCertHelper.connect(transmitter).addAdmin(
        certsChain.map(x => x.cert),
        0,
        rawSignature2
      )
      certData = await adminCertHelper.getAdminData();
      expect(certData.toObject(true)).to.deep.equal(expectedCertData);

    });

    it("isAdmin return false for non intialized user", async function () {
      const isInitialized = await adminCertHelper.isAdmin(await deployer.getAddress());
      expect(isInitialized).to.be.false;
    });

    it("Should not addAdmin if certs chain is not valid", async function () {
      const transmitterAddress = await transmitter.getAddress();

      const dataHash = ethers.sha256(ethers.AbiCoder.defaultAbiCoder().encode(["address", "uint256", "address", "uint256"], [transmitterAddress, chainId, await adminCertHelper.getAddress(), 0]));
      const signature = certsChain[0].signer.signingKey.sign(dataHash)
      const rawSignature = signature.r + signature.s.slice(2) + signature.v.toString(16);
      const certs = [certsChain[1].cert, certsChain[0].cert]
      await expect(adminCertHelper.connect(transmitter).addAdmin(
        certs,
        0,
        rawSignature
      )).to.be.revertedWith("Invalid signature");

    });

    it("Should not addAdmin if solution is not authorized", async function () {
      await expect(adminCertHelper.connect(transmitter).addAdmin([certsChain[1].cert], 0, "0x")).to.be.revertedWith("Only authorized solution");
    });

    it("Should not addAdmin if signature is not valid", async function () {
      const transmitterAddress = await transmitter.getAddress();

      const dataHash = ethers.sha256(ethers.AbiCoder.defaultAbiCoder().encode(["address", "uint256", "address", "uint256"], [transmitterAddress, chainId, await adminCertHelper.getAddress(), 0]));
      const wrongSignature = certsChain[1].signer.signingKey.sign(dataHash)
      const rawSignature = wrongSignature.r + wrongSignature.s.slice(2) + wrongSignature.v.toString(16);
      await expect(adminCertHelper.connect(transmitter).addAdmin(
        certsChain.map(x => x.cert),
        0,
        rawSignature
      )).to.be.revertedWith("Invalid data signature");

    });

    it("Should not addAdmin if signatureNonce, chain id or contract address is not valid", async function () {
      const transmitterAddress = await transmitter.getAddress();
      const chainId = (await deployer.provider?.getNetwork())?.chainId
      const contractAddress = await adminCertHelper.getAddress()
      const nonce = 0

      const badAddress = certVerifier.address

      {

        const dataHash = ethers.sha256(ethers.AbiCoder.defaultAbiCoder().encode(["address", "uint256", "address", "uint256"], [badAddress, chainId, contractAddress, nonce]));
        const wrongSignature = certsChain[0].signer.signingKey.sign(dataHash)
        const rawSignature = wrongSignature.r + wrongSignature.s.slice(2) + wrongSignature.v.toString(16);
        await expect(adminCertHelper.connect(transmitter).addAdmin(
          certsChain.map(x => x.cert),
          0,
          rawSignature
        )).to.be.revertedWith("Invalid data signature");
      }
      {
        const dataHash = ethers.sha256(ethers.AbiCoder.defaultAbiCoder().encode(["address", "uint256", "address", "uint256"], [transmitterAddress, 1, badAddress, nonce]));
        const wrongSignature = certsChain[0].signer.signingKey.sign(dataHash)
        const rawSignature = wrongSignature.r + wrongSignature.s.slice(2) + wrongSignature.v.toString(16);
        await expect(adminCertHelper.connect(transmitter).addAdmin(
          certsChain.map(x => x.cert),
          0,
          rawSignature
        )).to.be.revertedWith("Invalid data signature");
      }
      {
        const dataHash = ethers.sha256(ethers.AbiCoder.defaultAbiCoder().encode(["address", "uint256", "address", "uint256"], [transmitterAddress, chainId, contractAddress, 1000]));
        const wrongSignature = certsChain[0].signer.signingKey.sign(dataHash)
        const rawSignature = wrongSignature.r + wrongSignature.s.slice(2) + wrongSignature.v.toString(16);
        await expect(adminCertHelper.connect(transmitter).addAdmin(
          certsChain.map(x => x.cert),
          0,
          rawSignature
        )).to.be.revertedWith("Invalid data signature");
      }
    });
  });
});
