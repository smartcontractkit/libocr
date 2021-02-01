# libocr

libocr consists of a Go library and a set of Solidity smart contracts that implement the *Chainlink Off-chain Reporting Protocol*. The protocol is a Byzantine fault tolerant protocol that allows a set of oracles to generate an aggregate report of the oracles' observations of some underlying data source *off chain*. This report is then transmitted to an on-chain contract in a single transaction.

You may also be interested in [libocr's integration into the actual Chainlink node](https://github.com/smartcontractkit/chainlink/tree/develop/core/services/offchainreporting).


## Protocol Description

Protocol execution mostly happens off chain over a peer to peer network between Chainlink nodes. The nodes regularly elect a new leader node who drives the rest of the protocol.

The leader regularly requests followers to provide freshly signed observations and aggregates them into a report. It then sends this report back to the followers and asks them to verify the report's validity. If a quorum of followers approves the report by sending a signed copy back to the leader, the leader assembles a final report with the quorum's signatures and broadcasts it to all followers.

The nodes then attempt to transmit the final report to the smart contract according to a randomized schedule. Finally, the smart contract verifies that a quorum of nodes signed the report and exposes the median value to consumers.


## Organization
```
.
├── contract: Ethereum smart contracts
├── gethwrappers: go-ethereum bindings for the contracts, generated with abigen
├── networking: libp2p-based p2p networking layer
├── offchainreporting: offchain reporting protocol
├── permutation: helper package for generating permutations
└── subprocesses: helper package for managing go routines
```
