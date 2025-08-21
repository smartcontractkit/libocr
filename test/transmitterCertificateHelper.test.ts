import { expect } from "chai";
import { ethers } from "hardhat";
import { Signer } from "ethers";
import { takeSnapshot, SnapshotRestorer } from "@nomicfoundation/hardhat-toolbox/network-helpers";
import { CertificateHelper } from './helpers/certificateHelper';
import { CertVerifierHelper } from './helpers/certVerifierHelper';
import { TransmitterCertificateHelper } from "../typechain-types";

describe("TransmitterCertificateHelper", function () {
  let transmitter: Signer, deployer: Signer;
  let transmitterCertHelper: TransmitterCertificateHelper;
  let certVerifier: CertVerifierHelper;
  let snapshot: SnapshotRestorer;
  let certsChain: CertificateHelper[]

  before(async function () {
    [deployer, transmitter] = await ethers.getSigners();

    certVerifier = await CertVerifierHelper.create();

    const factory = await ethers.getContractFactory("TransmitterCertificateHelper");
    transmitterCertHelper = await factory.deploy(certVerifier.address);


    const certificate = await CertificateHelper.gen();
    const certificate2 = await CertificateHelper.gen(certificate);
    certsChain = [certificate2, certificate];
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
      const currentVerifier = await transmitterCertHelper.getCertVerifier();
      expect(currentVerifier).to.equal(certVerifier.address);
    });

    it("Emit event CertVerifierSet", async function () {
      await expect(transmitterCertHelper.deploymentTransaction())
        .to.emit(transmitterCertHelper, "CertVerifierSet")
        .withArgs(ethers.ZeroAddress, certVerifier.address);
    });
  });

  describe("setCertVerifier", function () {
    it("Only owner can setCertVerifier", async function () {

      const newCertVerifierAddress = await transmitter.getAddress()

      await expect(
        transmitterCertHelper.connect(transmitter).setCertVerifier(newCertVerifierAddress)
      ).to.be.revertedWith("Only callable by owner");

      await expect(transmitterCertHelper.setCertVerifier(newCertVerifierAddress))
        .to.emit(transmitterCertHelper, "CertVerifierSet")
        .withArgs(certVerifier.address, newCertVerifierAddress);

      expect(await transmitterCertHelper.getCertVerifier()).to.equal(newCertVerifierAddress);
    });

    it("Should revert if new cert verifier address is zero address", async function () {
      await expect(
        transmitterCertHelper.setCertVerifier(ethers.ZeroAddress)
      ).to.be.revertedWith("certVerifier cannot be zero address");
    });
  });

  describe("initializeTransmitter", function () {
    it("success initializeTransmitter", async function () {
      const transmitterAddress = await transmitter.getAddress();

      const dataHash = ethers.keccak256(ethers.AbiCoder.defaultAbiCoder().encode(["address"], [transmitterAddress]));
      const signature = certsChain[0].signer.signingKey.sign(dataHash)
      const rawSignature = signature.r + signature.s.slice(2) + signature.v.toString(16);
      const tx = await transmitterCertHelper.connect(transmitter).initializeTransmitter(
        certsChain.map(x => x.cert),
        0,
        rawSignature
      )

      const txReceipt = await tx.wait()

      const events = txReceipt!.logs.map(log => {
        const parsedLog = transmitterCertHelper.interface.parseLog(log)
        if (parsedLog?.name === "InitializeTransmitter") { return parsedLog }
      });

      expect(events.length).to.be.equal(1);

      const eventArgs = events[0]!.args
      const expectedEventArgs = {
        transmitter: transmitterAddress,
        certsChain: certsChain.map(x => x.cert),
        rootCertId: 0
      }
      expect(eventArgs.toObject(true)).to.deep.equal(expectedEventArgs);

      let certData = await transmitterCertHelper.getCertData(transmitterAddress);
      const expectedCertData = {
        initialized: true,
        certsChain: certsChain.map(x => x.cert),
        rootCertId: 0
      }
      expect(certData.toObject(true)).to.deep.equal(expectedCertData);


      // success reinitialize trasmitter
      await transmitterCertHelper.connect(transmitter).initializeTransmitter(
        certsChain.map(x => x.cert),
        0,
        rawSignature
      )
      certData = await transmitterCertHelper.getCertData(transmitterAddress);
      expect(certData.toObject(true)).to.deep.equal(expectedCertData);

    });

    it("isTransmitterInitialized return false for non intialized user", async function () {
      const isInitialized = await transmitterCertHelper.isTransmitterInitialized(await deployer.getAddress());
      expect(isInitialized).to.be.false;
    });

    it("Should not initializeTransmitter if certs chain is not valid", async function () {
      const transmitterAddress = await transmitter.getAddress();

      const dataHash = ethers.keccak256(ethers.AbiCoder.defaultAbiCoder().encode(["address"], [transmitterAddress]));
      const signature = certsChain[0].signer.signingKey.sign(dataHash)
      const rawSignature = signature.r + signature.s.slice(2) + signature.v.toString(16);
      const certs = [certsChain[1].cert, certsChain[0].cert]
      await expect(transmitterCertHelper.connect(transmitter).initializeTransmitter(
        certs,
        0,
        rawSignature
      )).to.be.revertedWith("Invalid signature");

    });

    it("Should not initializeTransmitter if trasnmitter address signature is not valid", async function () {
      const transmitterAddress = await transmitter.getAddress();

      const dataHash = ethers.keccak256(ethers.AbiCoder.defaultAbiCoder().encode(["address"], [transmitterAddress]));
      const wrongSignature = certsChain[1].signer.signingKey.sign(dataHash)
      const rawSignature = wrongSignature.r + wrongSignature.s.slice(2) + wrongSignature.v.toString(16);
      await expect(transmitterCertHelper.connect(transmitter).initializeTransmitter(
        certsChain.map(x => x.cert),
        0,
        rawSignature
      )).to.be.revertedWith("Invalid transmitter signature");

      const wrongDataHash = ethers.keccak256(ethers.AbiCoder.defaultAbiCoder().encode(["address"], [certVerifier.address]));
      const signature = certsChain[0].signer.signingKey.sign(wrongDataHash)
      const rawSignature2 = signature.r + signature.s.slice(2) + signature.v.toString(16);
      await expect(transmitterCertHelper.connect(transmitter).initializeTransmitter(
        certsChain.map(x => x.cert),
        0,
        rawSignature2
      )).to.be.revertedWith("Invalid transmitter signature")
    });

  });

});
