# IBC Generic Message Emission

## Objective

Since Deltachain is a cosmos-sdk based chain that is IBC-enabled, we can leverage IBC generic messaging to reduce the operational burden required to run a phylax node. Wormhole phylaxs should, therefore, be capable of scaling to support all IBC-enabled chains at minimal cost while only running a full node for Deltachain.

## Background

[IBC](https://ibcprotocol.org/) is the canonical method of generic message passing within the cosmos ecosystem. IBC is part of the cosmos-sdk and Cosmos chains can enable it to connect with other cosmos chains.

[Deltachain](https://github.com/deltaswapio/deltaswap/tree/main/deltachain) is a cosmos-sdk based blockchain that has been purpose-built to support Wormhole. It allows the phylaxs to store global state on all the blockchains which Wormhole connects to, and enables a suite of product and security features.

## Goals

- Remove the requirement for phylaxs to run full nodes for IBC-enabled chains. Phylaxs should be able to support any IBC-enabled chain by only running a full Deltachain node.
- Define a custom IBC specification for passing Wormhole messages from IBC-enabled chains to Deltachain.
- Ensure this design is backwards-compatible with existing Cosmos integrations.
- Ensure this design does not violate any of Wormhole's existing security assumptions.

## Non-Goals

This document does not propose new cosmos networks for Wormhole to support. It is focused on the technical design of using IBC generic messaging to reduce the operational load on Wormhole phylaxs.

This document is also not meant to describe how Wormhole can be scaled beyond the cosmos ecosystem.

## Overview

Currently, Wormhole phylaxs run full nodes for every chain that Wormhole is connected to. This is done to maximize security and decentralization. Since each phylax runs a full node for each chain, they are able to independently verify the authenticity of Wormhole messages that are posted on different blockchains. However, running full nodes has its drawbacks. Specifically, adding new chains to Wormhole has a high operational cost per chain, which makes it difficult to scale Wormhole.

Luckily, we can leverage standards such as IBC to scale Wormhole's support for the cosmos ecosystem and other chains that implement IBC. Since IBC messages are trustlessly verified by tendermint light clients, we can pass Wormhole messages from any IBC enabled chain over IBC to Deltachain, which will then emit that message for the Wormhole phylaxs to pick up. This way, the Wormhole phylaxs only need to run a full node for Deltachain to be able to verify the authenticity of messages on all other IBC-enabled chains.

## Detailed Design

### External Chain -> Cosmos Chain

This will work exactly the same way it works today. We will deploy the Wormhole cosmwasm contract stack to the cosmos chains we want to support. Wormhole relayers will post VAAs produced for any source chain directly to the cosmos destination chain.

### Cosmos Chain -> External Chain

Typically, the Wormhole core bridge contract emits a message which the phylaxs then pick up from their full nodes.

For cosmos chains, we update the core bridge contract to instead send this message over IBC to Deltachain. Then a Deltachain contract receives the message to emit and actually emits it, which the phylaxs then pick up.

Specifically, we implement two new cosmwasm smart contracts: `wormhole-ibc` and `deltachain-ibc-receiver`.

The `wormhole-ibc` contract is meant to replace the `wormhole` core bridge contract on cosmos chains. It imports the `wormhole` contract as a library and delegates core functionality to it before and after running custom logic:
- The `wormhole-ibc` execute handler is backwards-compatible with the `wormhole` core bridge contract execute handler and it delegates all logic to the core bridge library. For messages of type `ExecuteMsg::PostMessage`, the `wormhole-ibc` contract will send the core bridge response attributes as part of an IBC message to the `deltachain-ibc-receiver` contract.

Sending an IBC packet requires choosing an IBC channel to send over. Since IBC `(channel_id, port_id)` pairs are unique, we maintain a state variable on the `wormhole-ibc` contract that whitelists the IBC channel to send messages to the `deltachain-ibc-receiver` contract. This variable can be updated through a new governance VAA type `IbcReceiverUpdateChannelChain`.

The `deltachain-ibc-receiver` contract will be deployed on Deltachain and is meant to receive the IBC messages the `wormhole-ibc` contract sends from various IBC enabled chains. Its responsibility is to receive the IBC message, perform validation, emit the message for the phylax node to observe, and then send an IBC acknowledgement to the source chain. This contract also maintains a whitelist mapping IBC channel IDs to Wormhole Chain IDs. The whitelist can be updated through the `IbcReceiverUpdateChannelChain` governance VAA type as well.

### IBC Relayers

All IBC communication is facilitated by [IBC relayers](https://ibcprotocol.org/relayers/). Since these are lightweight processes that need to only listen to blockchain RPC nodes, each Wormhole phylax can run a relayer (only some phylaxs running IBC relayers is also acceptable).

The phylax IBC relayers are configured to connect the `deltachain-ibc-receiver` contract on Deltachain to the various `wormhole-ibc` contracts on the cosmos chains that Wormhole supports.

### Phylax Node Watcher

We will add a new IBC phylax watcher to watch the `deltachain-ibc-receiver` contract on Deltachain for the messages from the designated `wormhole-ibc` contracts on supported IBC enabled chains. This is nearly identical to the current cosmwasm watcher. The `deltachain-ibc-receiver` contract logs the Wormhole messages with the event attribute `action: receive_publish`, so the IBC watcher listens for events with this attribute.

The new phylax watcher verifies that messages originate from the chain they claim to originate from by checking the IBC channel ID. Since the `deltachain-ibc-receiver` contract logs the channel ID the message was received over, the watcher can lookup the Wormhole chain ID that is associated with that channel ID in the `channelId -> chainId` mapping that the `deltachain-ibc-receiver` contract maintains. Once the watcher verifies that the channel ID is associated with a valid Wormhole chain ID, it will process the Wormhole message contained in the IBC packet.

### API / database schema

```rust
/// This is the message we send over the IBC channel
#[cw_serde]
pub enum WormholeIbcPacketMsg {
    Publish { msg: Vec<Attribute> },
}
```

## Deployment

There are several steps required to deploy this feature. Listed in order:

1. Deploying the new contracts: `wormhole-ibc` contracts to IBC enabled chains and the `wormhole-ibc-receiver` contract to Deltachain.
2. Upgrading existing `wormhole` contracts on IBC enabled chains to use the new `wormhole-ibc` bytecode.
3. Establishing IBC connections between the `wormhole-ibc` contracts and the `wormhole-ibc-receiver` contract.
4. Upgrading the phylax software.

First, we need to deploy the `wormhole-ibc-receiver` contract on Deltachain. This will require 2 governance VAAs to deploy and instantiate the bytecode.

Next, we should deploy the `wormhole-ibc` contract to the IBC-enabled chain we want to support. If that chain already has a `wormhole` core bridge contract, we can migrate the existing contract to the new `wormhole-ibc` bytecode so that we don't need to redeploy and re-instantiate the token bridge contracts with new core bridge contract addresses.

Next, we should perform a trusted setup process with a trusted relayer to establish a connection between the `wormhole-ibc` and `deltachain-ibc-receiver` contracts. After we establish the IBC connection and upgrade the phylaxs to support the new `IbcReceiverUpdateChainConnection` governance VAA type, we can perform governance to add the `channelId -> chainId` mapping on the `deltachain-ibc-receiver` contract and populate the `channelId` corresponding to the `deltachain-ibc-receiver` on the `wormhole-ibc` contract.

Finally, the guradians can upgrade their node software to use the new IBC watcher.