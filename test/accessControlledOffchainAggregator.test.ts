import { expect } from "chai";
import { ethers, ignition } from "hardhat";
import { Result, Signer } from "ethers";
import { takeSnapshot, SnapshotRestorer } from "@nomicfoundation/hardhat-toolbox/network-helpers";

import acOffchainAggregator, { defaultConsensusConfig } from "../ignition/modules/acOffchainAggregator"
import { AccessControlledOffchainAggregator, OffchainAggregator } from "../typechain-types/contract/AccessControlledOffchainAggregator"
import { AdminCertificateHelper, EACAggregatorProxy } from "../typechain-types";

import { CertificateHelper } from "./helpers/certificateHelper";
import { CertVerifierHelper } from './helpers/certVerifierHelper';


describe("AccessControlledOffchainAggregator", function () {
  let randomUser: Signer, deployer: Signer;
  let snapshot: SnapshotRestorer;
  let proxy: EACAggregatorProxy;
  let aggregator: AccessControlledOffchainAggregator;
  let adminCertHelper: AdminCertificateHelper;

  before(async function () {
    const deployedContracts = await ignition.deploy(acOffchainAggregator);
    aggregator = deployedContracts.acOffchainAggregator as any as AccessControlledOffchainAggregator
    proxy = deployedContracts.proxy as any as EACAggregatorProxy
    adminCertHelper = deployedContracts.adminCertificateHelper as any as AdminCertificateHelper

    [deployer, randomUser] = await ethers.getSigners();

    const certVerifier = await CertVerifierHelper.create(await adminCertHelper.getCertVerifier());
    const certificate = await CertificateHelper.gen();
    const certificate2 = await CertificateHelper.gen(certificate);
    await certVerifier.setRootCert(certificate.cert);

    const chainId = (await deployer.provider?.getNetwork())?.chainId
    const dataHash = ethers.sha256(ethers.AbiCoder.defaultAbiCoder().encode(["address", "uint256", "address", "uint256"], [await deployer.getAddress(), chainId, await adminCertHelper.getAddress(), 0]));
    const signature = certificate2.signer.signingKey.sign(dataHash)
    const rawSignature = signature.r + signature.s.slice(2) + signature.v.toString(16);
    await adminCertHelper.setAuthorizedSolutionHash(certificate2.cert.userData)
    await adminCertHelper.addAdmin(
      [certificate2.cert],
      0,
      rawSignature
    )

  });

  beforeEach(async function () {
    snapshot = await takeSnapshot()
  });

  afterEach(async function () {
    await snapshot.restore()
  });

  async function getConfigParam() {
    const hardhatSigners = await ethers.getSigners()
    const deployer = hardhatSigners[0]
    const deployerAddress = await deployer.getAddress()
    const signers = [CertificateHelper.genRandomWallet()]
    const transmitters = [deployerAddress]
    for (let i = 1; i < 4; i++) {
      signers.push(CertificateHelper.genRandomWallet())
      transmitters.push(CertificateHelper.genRandomWallet().address)
    }
    const bytes16Value = ethers.zeroPadBytes("0x01", 16)
    const donConfig: OffchainAggregator.ConsensusDonConfigStruct = {
      signers: signers.map(x => x.address),
      transmitters: transmitters,
      threshold: 1,
      s: [1, 2, 2, 2],
      offchainPublicKeys: [ethers.ZeroHash, ethers.ZeroHash, ethers.ZeroHash, ethers.ZeroHash],
      peerIDs: "aaa,bbb,ccc,ddd",
      sharedSecretEncryptions: {
        diffieHellmanPoint: ethers.ZeroHash,
        sharedSecretHash: ethers.ZeroHash,
        encryptions: [bytes16Value, bytes16Value, bytes16Value, bytes16Value]
      }
    }
    return { signers, donConfig }
  }

  it("Success update data in AccessControlledOffchainAggregator", async function () {
    const { signers, donConfig } = await getConfigParam()
    await aggregator.setConfig(donConfig)

    const configDetails = await aggregator.latestConfigDetails()

    const rawReportContext = "0x0000000000000000000000" + configDetails.configDigest.toString().slice(2) + "0000000001"//<4-byte epoch><1-byte round>
    const observations = ["10000000000", "20101010101", "30000000000"]
    const rawObservers = "0x0302010000000000000000000000000000000000000000000000000000000000"
    const report = ethers.AbiCoder.defaultAbiCoder().encode(["bytes32", "bytes32", "int192[]"], [rawReportContext, rawObservers, observations])

    const hashForSign = ethers.keccak256(report)

    const sign0 = signers[0].signingKey.sign(ethers.toBeArray(hashForSign))
    const sign1 = signers[1].signingKey.sign(ethers.toBeArray(hashForSign))

    const sigObj0 = ethers.Signature.from(sign0);
    const sigObj1 = ethers.Signature.from(sign1);


    const rs = [sigObj0.r, sigObj1.r]
    const ss = [sigObj0.s, sigObj1.s]
    const rawVs = "0x" + "0" + sigObj0.yParity.toString() + "0" + sigObj1.yParity.toString() + "000000000000000000000000000000000000000000000000000000000000"

    const transmitArgs = [report, rs, ss, rawVs]

    await aggregator.transmit(...transmitArgs)

    const latestRoundData = await aggregator.latestRoundData()
    const expectedAnswer = observations[Number(BigInt(observations.length) / 2n)]
    const timestamp = (await ethers.provider.getBlock("latest"))?.timestamp
    const expectedLatestRoundData = {
      roundId: 1,
      answer: expectedAnswer,
      startedAt: timestamp,
      updatedAt: timestamp,
      answeredInRound: 1
    }

    expect((latestRoundData as any as Result).toObject()).to.deep.equal(expectedLatestRoundData)

    expect(await proxy.latestAnswer()).to.equal(expectedAnswer)
    expect(await aggregator.latestAnswer()).to.equal(expectedAnswer)
  });

  describe("setConfig", function () {
    it("Only confidential admin can setConfig", async function () {
      const { donConfig } = await getConfigParam()
      await expect(aggregator.connect(randomUser).setConfig(donConfig)).to.be.rejectedWith("Only confidential admin")
    });

    it("setConfig emit event ConfigSet", async function () {
      const { donConfig } = await getConfigParam()
      const tx = await aggregator.setConfig(donConfig)
      const txReceipt = await tx.wait()


      const events = txReceipt!.logs.map(log => {
        const parsedLog = aggregator.interface.parseLog(log)
        if (parsedLog?.name === "ConfigSet") { return parsedLog }
      });

      expect(events.length).to.be.equal(1);

      const eventArgs = events[0]!.args

      const codingStructure = ["tuple(int64 deltaProgress,int64 deltaResend,int64 deltaRound,int64 deltaGrace,int64 deltaC,uint64 alphaPPB,int64 deltaStage,uint8 rMax,uint8[] s,bytes32[] offchainPublicKeys,string peerIDs,tuple(bytes32 diffieHellmanPoint,bytes32 sharedSecretHash,bytes16[] encryptions) sharedSecretEncryptions)"]
      const nodesConfig = {
        deltaProgress: defaultConsensusConfig.deltaProgress,
        deltaResend: defaultConsensusConfig.deltaResend,
        deltaRound: defaultConsensusConfig.deltaRound,
        deltaGrace: defaultConsensusConfig.deltaGrace,
        deltaC: defaultConsensusConfig.deltaC,
        alphaPPB: defaultConsensusConfig.alphaPPB,
        deltaStage: defaultConsensusConfig.deltaStage,
        rMax: defaultConsensusConfig.rMax,
        s: donConfig.s,
        offchainPublicKeys: donConfig.offchainPublicKeys,
        peerIDs: donConfig.peerIDs,
        sharedSecretEncryptions: donConfig.sharedSecretEncryptions
      }
      const encoded = ethers.AbiCoder.defaultAbiCoder().encode(codingStructure, [nodesConfig]);
      const expectedEventArgs = {
        previousConfigBlockNumber: 0,
        configCount: 1,
        signers: donConfig.signers,
        transmitters: donConfig.transmitters,
        threshold: donConfig.threshold,
        encodedConfigVersion: defaultConsensusConfig.encodedConfigVersion,
        encoded: encoded
      }

      expect(eventArgs.toObject(true)).to.deep.equal(expectedEventArgs);
    });
  })

  describe("setConsensusConfig", function () {
    it("Only contract owner can setConsensusConfig", async function () {
      await expect(aggregator.connect(randomUser).setConsensusConfig(0, 0, 0, 0, 0, 0, 0, 0, 0)).to.be.rejectedWith("Only callable by owner")
    });

    it("setConsensusConfig emit event ConfigSet only if there was setConfig", async function () {
      {
        const tx = await aggregator.setConsensusConfig(0, 0, 0, 0, 0, 0, 0, 0, 0)
        const txReceipt = await tx.wait()
        const events = txReceipt!.logs.map(log => {
          const parsedLog = aggregator.interface.parseLog(log)
          if (parsedLog?.name === "ConfigSet") { return parsedLog }
        });
        expect(events.length).to.be.equal(0);
      }

      const { donConfig } = await getConfigParam()
      await aggregator.setConfig(donConfig)
      const blockNumber = (await deployer.provider?.getBlock("latest"))?.number

      const nodesConfig = {
        deltaProgress: 5,
        deltaResend: 6,
        deltaRound: 7,
        deltaGrace: 8,
        deltaC: 9,
        alphaPPB: 10,
        deltaStage: 11,
        rMax: 12,
        s: donConfig.s,
        offchainPublicKeys: donConfig.offchainPublicKeys,
        peerIDs: donConfig.peerIDs,
        sharedSecretEncryptions: donConfig.sharedSecretEncryptions
      }

      const encodedConfigVersion = 9
      const tx = await aggregator.setConsensusConfig(nodesConfig.deltaProgress, nodesConfig.deltaResend, nodesConfig.deltaRound, nodesConfig.deltaGrace, nodesConfig.deltaC, nodesConfig.alphaPPB, nodesConfig.deltaStage, nodesConfig.rMax, encodedConfigVersion)
      const txReceipt = await tx.wait()

      const events = txReceipt!.logs.map(log => {
        const parsedLog = aggregator.interface.parseLog(log)
        if (parsedLog?.name === "ConfigSet") { return parsedLog }
      });

      expect(events.length).to.be.equal(1);

      const eventArgs = events[0]!.args

      const codingStructure = ["tuple(int64 deltaProgress,int64 deltaResend,int64 deltaRound,int64 deltaGrace,int64 deltaC,uint64 alphaPPB,int64 deltaStage,uint8 rMax,uint8[] s,bytes32[] offchainPublicKeys,string peerIDs,tuple(bytes32 diffieHellmanPoint,bytes32 sharedSecretHash,bytes16[] encryptions) sharedSecretEncryptions)"]
      const encoded = ethers.AbiCoder.defaultAbiCoder().encode(codingStructure, [nodesConfig]);
      const expectedEventArgs = {
        previousConfigBlockNumber: blockNumber,
        configCount: 2,
        signers: donConfig.signers,
        transmitters: donConfig.transmitters,
        threshold: donConfig.threshold,
        encodedConfigVersion: encodedConfigVersion,
        encoded: encoded
      }

      expect(eventArgs.toObject(true)).to.deep.equal(expectedEventArgs);
    });
  })

  describe("setAdminCertificateHelper", function () {
    it("Only owner can setAdminCertificateHelper", async function () {

      const newAdminCertificateHelper = await randomUser.getAddress()

      await expect(
        aggregator.connect(randomUser).setAdminCertificateHelper(newAdminCertificateHelper)
      ).to.be.revertedWith("Only callable by owner");

      await expect(aggregator.setAdminCertificateHelper(newAdminCertificateHelper))
        .to.emit(aggregator, "AdminCertificateHelperSet")
        .withArgs(await adminCertHelper.getAddress(), newAdminCertificateHelper);

      expect(await aggregator.adminCertificateHelper()).to.equal(newAdminCertificateHelper);
    });

    it("Should revert if new AdminCertificateHelper is zero address", async function () {
      await expect(
        aggregator.setAdminCertificateHelper(ethers.ZeroAddress)
      ).to.be.revertedWith("AdminCertificateHelper cannot be zero address");
    });
  });
});
