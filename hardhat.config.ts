import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import 'hardhat-contract-sizer';

import * as dotenv from 'dotenv';
import path from 'path';
dotenv.config({ path: path.join(__dirname, '/.env') });

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
      version: "0.7.6",
      settings: {
        optimizer: {
          enabled: true,
          runs: 20000
        }
      },
    },
    {
      version: "0.8.30",
      settings: {
        viaIR: true,
        optimizer: {
          enabled: true,
          runs: 5000
        }
      },
    },
    ]
  },
  paths: { sources: "./contract" },
  networks: {
    opbnbTestnet: {
      chainId: 5611,
      ...URL_ACCOUNTS_SETTINGS,
      ignition: {
        maxFeePerGas: 1_000_000n,
        maxPriorityFeePerGas: 1n,
        disableFeeBumping: true,
      },
    },
    opbnb: {
      chainId: 204,
      ...URL_ACCOUNTS_SETTINGS,
      ignition: {
        maxFeePerGas: 1_000_000n,
        maxPriorityFeePerGas: 1n,
        disableFeeBumping: true,
      },
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
  etherscan: {
    apiKey: process.env.SCAN_API_KEY,
    customChains: [
      {
        network: 'opbnbTestnet',
        chainId: 5611,
        urls: {
          apiURL: 'https://api.etherscan.io/v2/api',
          browserURL: 'https://opbnb-testnet.bscscan.com/',
        },
      },
      {
        network: 'opbnb',
        chainId: 204,
        urls: {
          apiURL: 'https://api.etherscan.io/v2/api',
          browserURL: 'https://opbnb.bscscan.com/',
        },
      },
    ],
  },
};

export default config;
