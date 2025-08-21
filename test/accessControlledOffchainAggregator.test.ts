import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";
import { expect } from "chai";
import { ethers, ignition } from "hardhat";
import { Result } from "ethers";

import acOffchainAggregator from "../ignition/modules/acOCR2Aggregator"
import { AccessControlledOffchainAggregator } from "../typechain-types/AccessControlledOffchainAggregator"

import { EACAggregatorProxy, TransmitterCertificateHelper } from "../typechain-types";
import { CertificateHelper } from "./helpers/certificateHelper";

describe("AccessControlledOffchainAggregator", function () {

  async function deployFixture() {
    const deployedContracts = await ignition.deploy(acOffchainAggregator);
    const aggregator = deployedContracts.acOffchainAggregator as any as AccessControlledOffchainAggregator
    const proxy = deployedContracts.proxy as any as EACAggregatorProxy
    const TCHelper = deployedContracts.transmitterCertificateHelper as any as TransmitterCertificateHelper

    return { aggregator, proxy, TCHelper };
  }

  async function transmit(aggregator: AccessControlledOffchainAggregator, deployerAddress: string) {
    const signers = [CertificateHelper.genRandomWallet()]
    const transmitters = [deployerAddress]
    for (let i = 1; i < 4; i++) {
      signers.push(CertificateHelper.genRandomWallet())
      transmitters.push(CertificateHelper.genRandomWallet().address)
    }
    const payees = transmitters
    const setPayees = await aggregator.setPayees(transmitters, payees)
    await setPayees.wait()

    const threshold = 1
    const encodedConfigVersion = 1
    const encoded = "0x"
    const mainConfig = [signers.map(x => x.address), transmitters, threshold, encodedConfigVersion, encoded]

    const setConfig = await aggregator.setConfig(...mainConfig)
    await setConfig.wait()

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
    return observations
  }

  it("Success update data in AccessControlledOffchainAggregator", async function () {
    const hardhatSigners = await ethers.getSigners()
    const deployer = hardhatSigners[0]
    const deployerAddress = await deployer.getAddress()

    const { aggregator, proxy, TCHelper } = await loadFixture(deployFixture);
    const certVerifier = await ethers.getContractAt("CertificateVerifier", await TCHelper.getCertVerifier())

    const certificate = await CertificateHelper.gen();
    const certificate2 = await CertificateHelper.gen(certificate);

    await certVerifier.setRootCert(certificate.cert)
    const hash = ethers.keccak256(ethers.AbiCoder.defaultAbiCoder().encode(["address"], [deployerAddress]));

    const signature = await certificate2.signer.signingKey.sign(hash);
    const rawSignature = signature.r + signature.s.slice(2) + signature.v.toString(16);
    await TCHelper.initializeTransmitter([certificate2.cert], 0, rawSignature)

    const observations = await transmit(aggregator, deployer.address)

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

  it("Should revert if Transmitter is not initialized in TransmitterCertificateHelper contract", async function () {
    const hardhatSigners = await ethers.getSigners()
    const deployer = hardhatSigners[0]

    const { aggregator } = await loadFixture(deployFixture);
    await expect(transmit(aggregator, deployer.address)).to.be.rejectedWith("Transmitter must be intialized")

  });

});
