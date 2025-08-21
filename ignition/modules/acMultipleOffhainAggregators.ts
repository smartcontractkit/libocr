// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
import fs from 'fs';
import path from "path";
import { strict as assert } from 'assert';
import { billingACmodule, linkDefault, requesterACmodule, transmitterCertificateHelperModule } from "./acOCR2Aggregator";

const acMultipleOffhainAggregators = buildModule("acMultipleOffhainAggregators", (m) => {
  const args = process.argv.slice(2);
  const parametersFlagIndex = args.indexOf('--parameters');
  assert(args.length > parametersFlagIndex + 1, "pairs file should be specified");
  const parametersFilePath = args[parametersFlagIndex + 1];
  const parametersJson = JSON.parse(fs.readFileSync(path.join(__dirname, '../../' + parametersFilePath)).toString());
  const pairs = Object.entries(parametersJson['uniquePair'])
    .map((e, idx) => {
      const params: { decimals: string } = e[1] as any;
      const description = (e[0] as string);
      const uuid = idx.toString() + '_' + description.replace(new RegExp('[ /]+', 'g'), '_');
      return { uuid, description, decimals: params.decimals };
    });


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
  for (const pair of pairs) {
    const accessControlledOffchainAggregatorUuid = "accessControlledOffchainAggregator" + pair.uuid;
    const eacAggregatorProxyUuid = "Proxy" + pair.uuid;
    const constructorConfig = [
      billingConstructorArgs,
      "1",
      "0xffffffffffffffffffffffffffffffffffffffff",
      requesterAC,
      pair.decimals,
      pair.description,
      transmitterCertificateHelper
    ];
    const acOffchainAggregator = m.contract("AccessControlledOffchainAggregator", constructorConfig, { id: accessControlledOffchainAggregatorUuid });
    const proxy = m.contract("EACAggregatorProxy", [acOffchainAggregator, "0x0000000000000000000000000000000000000000"], { id: eacAggregatorProxyUuid });
    m.call(acOffchainAggregator, "addAccess", [proxy])
    aggregatorsDeployments['aggregator_' + pair.uuid] = acOffchainAggregator;
    aggregatorsDeployments['proxy_' + pair.uuid] = proxy;
  }

  return { ...aggregatorsDeployments, transmitterCertificateHelper };
});

export default acMultipleOffhainAggregators;
