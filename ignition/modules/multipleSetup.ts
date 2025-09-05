import { IgnitionModuleBuilder } from "@nomicfoundation/ignition-core/dist/src/types/module-builder";
import { NamedArtifactContractDeploymentFuture } from "@nomicfoundation/ignition-core/dist/src/types/module";
import { buildModule, } from "@nomicfoundation/hardhat-ignition/modules";
import multipleDeploy, { uniquePairs } from "./multipleDeploy";

import fs from 'fs'
import path from 'path';

const args = process.argv.slice(2);
const paramsFlagIndex = args.indexOf('--parameters');
const paramsPath = paramsFlagIndex !== -1 ? args[paramsFlagIndex + 1] : "ignition/parameters.json";
const multipleSetupParams = JSON.parse(fs.readFileSync(path.join(__dirname, '../../' + paramsPath)).toString()).multipleSetup;

async function setConfig(m: IgnitionModuleBuilder, aggregator: NamedArtifactContractDeploymentFuture<"AccessControlledOffchainAggregator">, options?: { deltaC: string, alphaPPB: string }) {
  // node config
  const { deltaProgress, deltaResend, deltaRound, deltaGrace, deltaC: deltaCdefault, alphaPPB: alphaPPBdefault, deltaStage, rMax, encodedConfigVersion } = multipleSetupParams.nodesConfig

  const deltaC = options?.deltaC !== undefined ? options.deltaC : deltaCdefault
  const alphaPPB = options?.alphaPPB !== undefined ? options.alphaPPB : alphaPPBdefault

  const config = [deltaProgress, deltaResend, deltaRound, deltaGrace, deltaC, alphaPPB, deltaStage, rMax, encodedConfigVersion]

  m.call(aggregator, "setConsensusConfig", config, { id: aggregator.id.replace("#", "_") + "_setConsensusConfig" })
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
