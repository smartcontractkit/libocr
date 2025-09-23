# libocr

libocr consists of a Go library and a set of Solidity smart contracts that implements various versions of the *Chainlink Offchain Reporting Protocol*, a [Byzantine fault tolerant](https://en.wikipedia.org/wiki/Byzantine_fault) "consensus" protocol that allows a set of oracles to generate *offchain* an aggregate report of the oracles' observations of some underlying data source. This report is then transmitted to an onchain contract in a single transaction.

You may also be interested in libocr's integration into the actual Chainlink node. ([V1](https://github.com/smartcontractkit/chainlink/tree/develop/core/services/ocr) [V2](https://github.com/smartcontractkit/chainlink/tree/develop/core/services/ocr2) [V3](https://github.com/smartcontractkit/chainlink/tree/develop/core/services/ocr3))


## Protocol Description

Please see the whitepapers available at https://chainlinklabs.com/research for detailed protocol descriptions.

## Protocol Versions

- OCR1 is deprecated and being phased out.
- OCR2 & OCR3 are in production.
- OCR3.1 is in alpha and excluded from any bug bounties at this time. So is the associated Key-Value-Database in `offchainreporting2plus/keyvaluedatabase/`.

## Organization
```
├── bigbigendian: helper package
├── commontypes: shared type definitions
├── contract: OCR1 Ethereum contracts
├── contract2: OCR2 Ethereum contracts
├── contract3: OCR3 Ethereum contracts
├── gethwrappers: go-ethereum bindings for the OCR1 contracts, generated with abigen
├── gethwrappers2: go-ethereum bindings for the OCR2 contracts, generated with abigen
├── gethwrappers3: go-ethereum bindings for the OCR3 contracts, generated with abigen
├── networking: OCR networking layer
├── offchainreporting: OCR1
├── offchainreporting2: OCR2-specific
├── offchainreporting2plus: OCR2 and beyond (These versions share many interface definitions to make integration of new versions easier)
├── permutation: helper package
├── quorumhelper: helper package
├── ragep2p: p2p networking
└── subprocesses: helper package
```

## Setup for contracts usage

### 1. Install dependencies
```sh
npm install
```

### 2. Environment configuration
Create a `.env` file in the root directory:
```sh
cp .env.example .env
```

Add your configuration to `.env`:
```env
DEPLOYER_PRIVATE_KEY=your_private_key_here
RPC_URL=https://opbnb-testnet-rpc.bnbchain.org
SCAN_API_KEY=your_bscscan_api_key_here
```

**Important:** 
- Never commit your `.env` file to version control
- Use a dedicated wallet for testing, not your main wallet
- For mainnet deployment, use `https://opbnb-rpc.bnbchain.org`
## Usage

### Deploy contracts
```sh
# Deploy one AccessControlledOffchainAggregator BTC / USD
npx hardhat ignition deploy  ignition/modules/acOffchainAggregator.ts --network <opbnb/opbnbTestnet/localhost/hardhat>

# multiple deploy and setup
npx hardhat ignition deploy  ignition/modules/multipleDeploy.ts --network <opbnb/opbnbTestnet/localhost/hardhat>
npx hardhat ignition deploy  ignition/modules/multipleSetup.ts --network <opbnb/opbnbTestnet/localhost/hardhat>
```

### Verify contracts
```sh
# add to deploy command --verify flag. Example:
npx hardhat ignition deploy  ignition/modules/acOffchainAggregator.ts --network <opbnb/opbnbTestnet/localhost/hardhat> --verify

# or
npx hardhat ignition verify <chain-5611/chain-204> --include-unrelated-contracts
```

### Build files {description:deployedFeedAddress} and {description:deployedProxyAddress}

Output files: `aggregators-<204/5611>.json` and `proxies-<204/5611>.json`
Custom dapps should use proxies.
```sh
npx hardhat buildDeployedAddressesByDescription --chain-id <204,5611>

```

### Run tests (no setup required)
```sh
npm run test
```

## Networks

- **opBNB Testnet**: Chain ID 5611, RPC: https://opbnb-testnet-rpc.bnbchain.org
- **opBNB Mainnet**: Chain ID 204, RPC: https://opbnb-rpc.bnbchain.org
