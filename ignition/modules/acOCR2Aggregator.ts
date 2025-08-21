// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const linkDefault = "0xDA333D5948610a7634202A9420c0c17034B6484A"

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

export const transmitterCertificateHelperModule = buildModule("TransmitterCertificateHelper", (m) => {
  const { certificateVerifier } = m.useModule(certificateVerifierModule);
  const transmitterCertificateHelper = m.contract("TransmitterCertificateHelper", [certificateVerifier]);
  return { transmitterCertificateHelper };
});

const acOffchainAggregator = buildModule("acOffchainAggregator", (m) => {
  const { requesterAC } = m.useModule(requesterACmodule);
  const { billingAC } = m.useModule(billingACmodule);
  const { transmitterCertificateHelper } = m.useModule(transmitterCertificateHelperModule);

  const link = m.getParameter("linkToken", linkDefault);

  const maximumGasPrice = 1000
  const reasonableGasPrice = 1
  const microLinkPerEth = 205305307
  const linkGweiPerObservation = 701978
  const linkGweiPerTransmission = 4212083
  const billingConstructorArgs = [maximumGasPrice, reasonableGasPrice, microLinkPerEth, linkGweiPerObservation, linkGweiPerTransmission, link, billingAC]
  const constructorConfig = [billingConstructorArgs, "10000000000", "1000000000000000", requesterAC, 8, "BTC / USD", transmitterCertificateHelper]
  const acOffchainAggregator = m.contract("AccessControlledOffchainAggregator", constructorConfig);


  const proxy = m.contract("EACAggregatorProxy", [acOffchainAggregator, "0x0000000000000000000000000000000000000000"]);

  m.call(acOffchainAggregator, "addAccess", [proxy])

  return { acOffchainAggregator, proxy, transmitterCertificateHelper };
});

export default acOffchainAggregator;
