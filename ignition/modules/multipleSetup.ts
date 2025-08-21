import { IgnitionModuleBuilder } from "@nomicfoundation/ignition-core/dist/src/types/module-builder";
import { NamedArtifactContractDeploymentFuture } from "@nomicfoundation/ignition-core/dist/src/types/module";
import { buildModule, } from "@nomicfoundation/hardhat-ignition/modules";
import multipleDeploy, { uniquePairs } from "./multipleDeploy";
import { AbiCoder } from 'ethers';

import fs from 'fs'
import path from 'path';

const args = process.argv.slice(2);
const paramsFlagIndex = args.indexOf('--parameters');
const paramsPath = paramsFlagIndex !== -1 ? args[paramsFlagIndex + 1] : "ignition/parameters.json";
const multipleSetupParams = JSON.parse(fs.readFileSync(path.join(__dirname, '../../' + paramsPath)).toString()).multipleSetup;

const codingStructure = ["tuple(int64 deltaProgress,int64 deltaResend,int64 deltaRound,int64 deltaGrace,int64 deltaC,uint64 alphaPPB,int64 deltaStage,uint8 rMax,uint8[] s,bytes32[] offchainPublicKeys,string peerIDs,tuple(bytes32 diffieHellmanPoint,bytes32 sharedSecretHash,bytes16[] encryptions) sharedSecretEncryptions)"]

async function setConfig(m: IgnitionModuleBuilder, aggregator: NamedArtifactContractDeploymentFuture<"AccessControlledOffchainAggregator">, options?: { deltaC: string, alphaPPB: string }) {
  const deployer = m.getAccount(0)
  const signers = multipleSetupParams.signers
  const transmitters = multipleSetupParams.transmitters
  const payees = multipleSetupParams.payees ?? new Array(transmitters.length).fill(deployer)

  const setPayees = m.call(aggregator, "setPayees", [transmitters, payees], { id: aggregator.id.replace("#", "_") + "_setPayees" })

  // node config
  const threshold = 1
  const encodedConfigVersion = 1
  const nodesConfig = { ...multipleSetupParams.nodesConfig }

  if (options?.deltaC !== undefined) {
    nodesConfig.deltaC = options.deltaC
  }
  if (options?.alphaPPB !== undefined) {
    nodesConfig.alphaPPB = options.alphaPPB
  }
  const encoded = AbiCoder.defaultAbiCoder().encode(codingStructure, [nodesConfig]);

  const mainConfig = [signers, transmitters, threshold, encodedConfigVersion, encoded]

  m.call(aggregator, "setConfig", mainConfig, { after: [setPayees], id: aggregator.id.replace("#", "_") + "_setConfig" })
}

export const multipleSetup = buildModule("multipleSetup", (m) => {
  const aggregators = m.useModule(multipleDeploy);

  const prefix = "aggregator_"
  for (const key in aggregators) {
    if (key.startsWith(prefix)) {
      let options: { deltaC: string, alphaPPB: string } | undefined = undefined;
      const pair = uniquePairs[key.slice(prefix.length,)]
      if (pair !== undefined) {
        options = {
          alphaPPB: pair.alphaPPB,
          deltaC: pair.deltaC
        }
      }
      setConfig(m, aggregators[key], options)
    }
  }

  return {};
});

export default multipleSetup;
