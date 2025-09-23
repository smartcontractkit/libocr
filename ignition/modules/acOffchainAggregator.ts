// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

export const linkDefault = "0xDA333D5948610a7634202A9420c0c17034B6484A"

const authorizedSolutionHashDefault = "0x62b71aa4c49e75ed0066dece4c0a28c37d8543c4870a475b9216e0a04bd778b2"

export const defaultConsensusConfig = {
  "deltaProgress": "35000000000",
  "deltaResend": "17000000000",
  "deltaRound": "30000000000",
  "deltaGrace": "12000000000",
  "deltaC": "86400000000000",
  "alphaPPB": "5000000",
  "deltaStage": "60000000000",
  "rMax": "6",
  "encodedConfigVersion": "1"
}

export const LowLevelCallLibModule = buildModule("LowLevelCallLib", (m) => {
  const lowLevelCallLib = m.library("LowLevelCallLib");
  return { lowLevelCallLib };
});

export const requesterACmodule = buildModule("requesterAC", (m) => {
  const requesterAC = m.contract("SimpleWriteAccessController");
  return { requesterAC };
});

export const billingACmodule = buildModule("billingAC", (m) => {
  const billingAC = m.contract("SimpleWriteAccessController");
  return { billingAC };
});

export const certificateVerifierModule = buildModule("CertificateVerifier", (m) => {
  const certificateVerifier = m.contract("CertificateVerifier");
  return { certificateVerifier };
});

export const adminCertificateHelperModule = buildModule("AdminCertificateHelper", (m) => {
  const { certificateVerifier } = m.useModule(certificateVerifierModule);

  const authorizedSolutionHash = m.getParameter("authorizedSolutionHash", authorizedSolutionHashDefault);
  const adminCertificateHelper = m.contract("AdminCertificateHelper", [certificateVerifier, authorizedSolutionHash]);
  return { adminCertificateHelper };
});


const acOffchainAggregator = buildModule("acOffchainAggregator", (m) => {
  const { lowLevelCallLib } = m.useModule(LowLevelCallLibModule);
  const { adminCertificateHelper } = m.useModule(adminCertificateHelperModule);
  const { requesterAC } = m.useModule(requesterACmodule);
  const { billingAC } = m.useModule(billingACmodule);

  const link = m.getParameter("linkToken", linkDefault);

  const maximumGasPrice = 1000
  const reasonableGasPrice = 1
  const microLinkPerEth = 205305307
  const linkGweiPerObservation = 701978
  const linkGweiPerTransmission = 4212083
  const billingConstructorArgs = [maximumGasPrice, reasonableGasPrice, microLinkPerEth, linkGweiPerObservation, linkGweiPerTransmission, link, billingAC]
  const constructorConfig = [billingConstructorArgs, "10000000000", "1000000000000000", requesterAC, 8, "BTC / USD", adminCertificateHelper]
  const acOffchainAggregator = m.contract("AccessControlledOffchainAggregator", constructorConfig, { libraries: { LowLevelCallLib: lowLevelCallLib } });


  const proxy = m.contract("EACAggregatorProxy", [acOffchainAggregator, "0x0000000000000000000000000000000000000000"]);

  m.call(acOffchainAggregator, "addAccess", [proxy])

  const config = [defaultConsensusConfig.deltaProgress, defaultConsensusConfig.deltaResend, defaultConsensusConfig.deltaRound, defaultConsensusConfig.deltaGrace, defaultConsensusConfig.deltaC, defaultConsensusConfig.alphaPPB, defaultConsensusConfig.deltaStage, defaultConsensusConfig.rMax, defaultConsensusConfig.encodedConfigVersion]

  m.call(acOffchainAggregator, "setConsensusConfig", config)

  return { acOffchainAggregator, proxy, adminCertificateHelper };
});

export default acOffchainAggregator;
