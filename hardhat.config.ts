import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import 'hardhat-contract-sizer';

import * as dotenv from 'dotenv';
import path from 'path';
import { AccessControlledOffchainAggregator } from "./typechain-types";
dotenv.config({ path: path.join(__dirname, '/.env') });

task("setConfig").addParam('acOffchainAggregator', 'address of AccessControlledOffchainAggregator')
  .addFlag('replenishTransmitters', 'replenish transmitters with native asset')
  .setAction(async ({ acOffchainAggregator, replenishTransmitters }, hre) => {
    const ethers = hre.ethers
    const [deployer] = await ethers.getSigners();
    const contract: AccessControlledOffchainAggregator = await ethers.getContractAt("AccessControlledOffchainAggregator", acOffchainAggregator)

    const signers = ["0x22F2a420CD141A8b6d47c6e740cb9FBef7edfEd3", "0x034e1d4051Eff8389D0B0D37838D4D745Cb526FB", "0xe149095A6a0aF66652f8274306e8Ad83576F9658", "0xCDbd9bea9b32D0158f18E22daC818D40196908ef"]
    const transmitters = ["0x2d7686F4A6a9e3260236c0C22742f4cF5F4C5999", "0xaB8225862e3DE15dD0C0fb0B492a5Bc9e1D87694", "0x687469481cd54d00FDe81b0f668C116D20432B90", "0x9CB3aF1872E12F26D9279C48f8aDE4C06817D685"]
    const payees = new Array(transmitters.length).fill(deployer.address)

    if (replenishTransmitters) {
      for (const trans of transmitters) {
        const tx = await deployer.sendTransaction({ to: trans, value: ethers.parseEther("0.0005") })
        await tx.wait()
      }
    }

    let tx;
    tx = await contract.setPayees(transmitters, payees)
    await tx.wait()

    // node config
    const threshold = 1
    const encodedConfigVersion = 1
    const emptyConfig = {
      deltaProgress: 23000000000,
      deltaResend: 10000000000,
      deltaRound: 5000000000 * 4,
      deltaGrace: 3000000000,
      deltaC: 10000000000 * 6 * 10,
      alphaPPB: 500000,
      deltaStage: 10000000000,
      rMax: 100,
      s: [1, 1, 2],
      offchainPublicKeys: ["0x04c5774f6e04d2c03fed8e3d6894010fd263a3260944f2b888db6262b0c39885", "0xd74ec89b45cff0f480b9339037230753114308309e5deada41093da5ec5e8486", "0xab9813b58e6a847771bdae2199e16db3c130d159d58b1928671ea2a54889c1c8", "0x292bb623f861773c29fb8d8c34197b964e1790acaf0ac2573382ac3f168ac7fc"],
      peerIDs: "12D3KooWCLniDRotFVwBPC79319c6Ty5dZHnwmk2CdJPm3Zuq5zJ,12D3KooWBwb35i74FBdC5sZ6NsFTyJ8jzfwEvFyJhgrWGwBYSvVQ,12D3KooWLosbXRQqSVEKrXR1PnKqNviHkjhGbmTKF9FEiYZbaYR9,12D3KooWR8uj6pyqgPnyc1kdqWw7kvqifomA6cPLCh3kFbeVQVGd",
      sharedSecretEncryptions: {
        diffieHellmanPoint: "0x4ad1570774b185c731939a180354699474543754735952d0527e413eea450077",
        sharedSecretHash: "0xddb2241a7dc5e604a10301f24dd773cebd261a631eda8b1f22ffb5c00f8f605b",
        encryptions: [
          "0x38a9d5ae28f9f31c7198f8eef0c00987",
          "0x43943d8d5239500bc01291945f25c391",
          "0x9952b3eb6083681425e6f8831f0147a9",
          "0xfdb2fa7be97eede9a4a5075c55623feb"
        ]
      }
    };

    const codingStructure = ["tuple(int64 deltaProgress,int64 deltaResend,int64 deltaRound,int64 deltaGrace,int64 deltaC,uint64 alphaPPB,int64 deltaStage,uint8 rMax,uint8[] s,bytes32[] offchainPublicKeys,string peerIDs,tuple(bytes32 diffieHellmanPoint,bytes32 sharedSecretHash,bytes16[] encryptions) sharedSecretEncryptions)"]
    const encoded = ethers.AbiCoder.defaultAbiCoder().encode(codingStructure, [emptyConfig]);

    const mainConfig = [signers, transmitters, threshold, encodedConfigVersion, encoded]

    tx = await contract.setConfig(...mainConfig)
    await tx.wait()
  });

task("helper").setAction(async (args, hre) => { });


const DEFAULT_HARDHAT_PRIVATE_KEY = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

const URL_ACCOUNTS_SETTINGS = {
  url: process.env.RPC_URL ?? "",
  accounts: [process.env.DEPLOYER_PRIVATE_KEY ?? DEFAULT_HARDHAT_PRIVATE_KEY],
}

const config: HardhatUserConfig = {
  solidity: {
    compilers: [{
      version: "0.6.6",
      settings: {
        optimizer: {
          enabled: true,
          runs: 1000000
        }
      }
    },
    {
      version: "0.8.19",
      settings: {
        optimizer: {
          enabled: true,
          runs: 10000
        }, evmVersion: "paris"
      },
    },
    {
      version: "0.7.6",
      settings: {
        optimizer: {
          enabled: true,
          runs: 20000
        }
      },
    },
    ]
  },
  paths: { sources: "./contract" },
  networks: {
    opbnbTestnet: {
      chainId: 5611,
      ...URL_ACCOUNTS_SETTINGS
    },
    opbnb: {
      chainId: 204,
      ...URL_ACCOUNTS_SETTINGS
    },
    arbitrum: {
      chainId: 42161,
      ...URL_ACCOUNTS_SETTINGS
    },
  },
  contractSizer: {
    alphaSort: false,
    disambiguatePaths: true,
    runOnCompile: true,
    strict: false,
  },
};

export default config;
