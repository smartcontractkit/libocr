// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
import fs from 'fs';
import path from "path";
import { billingACmodule, linkDefault, requesterACmodule, transmitterCertificateHelperModule } from "./acOffchainAggregator";
export const uniquePairs = JSON.parse(fs.readFileSync(path.join(__dirname, '../' + "data-feeds.json")).toString()).uniquePair;

const multipleDeploy = buildModule("multipleDeploy", (m) => {

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

  const aggregatorsDeployments: any = {};
  for (const description in uniquePairs) {
    const decimals = uniquePairs[description].decimals
    const uuid = description.replace(new RegExp('[+ /().-]', 'g'), '_');
    const accessControlledOffchainAggregatorUuid = "accessControlledOffchainAggregator_" + uuid;
    const eacAggregatorProxyUuid = "Proxy_" + uuid;
    const constructorConfig = [
      billingConstructorArgs,
      "1",
      "0xffffffffffffffffffffffffffffffffffffffff",
      requesterAC,
      decimals,
      description,
      transmitterCertificateHelper
    ];

    const acOffchainAggregator = m.contract("AccessControlledOffchainAggregator", constructorConfig, { id: accessControlledOffchainAggregatorUuid });
    const proxy = m.contract("EACAggregatorProxy", [acOffchainAggregator, "0x0000000000000000000000000000000000000000"], { id: eacAggregatorProxyUuid });
    m.call(acOffchainAggregator, "addAccess", [proxy])
    aggregatorsDeployments['aggregator_' + description] = acOffchainAggregator;
    aggregatorsDeployments['proxy_' + description] = proxy;
  }

  return { ...aggregatorsDeployments };
});

export default multipleDeploy;
