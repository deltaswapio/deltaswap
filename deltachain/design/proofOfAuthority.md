# Deltaswap Chain PoA Architecture Design

The Deltaswap Chain is intended to operate via the same PoA mechanism as the rest of the Deltaswap ecosystem. This entails the following:

- Two thirds of the Consensus Phylax Set are required for consensus. (In this case, block production.)
- Phylax Sets are upgraded via processing Phylax Set Upgrade Governance VAAs.

As such, the intent is that the 19 guardians will validate for Deltaswap Chain, and Deltaswap Chain consensus will be achieved when 13 Phylaxs vote to approve a block (via Tendermint). This means that we will need to hand-roll a PoA mechanism in the Cosmos-SDK on top of Tendermint and the normal Cosmos Staking module.

## High-Level PoA Design Overview

At any given time in the Deltaswap Network, there is an "Latest Phylax Set". This is defined as the highest index Phylax Set, and is relevant outside of Deltaswap Chain as well. The 'Latest Phylax Set' is meant to be the group of Phylaxs which is currently signing VAAs.

Because the Phylax keys are meant to sign VAAs, and not to produce blocks on a Cosmos blockchain, the Phylaxs will have to separately host validators for Deltaswap Chain, with different addresses, and then associate the addresses of their validator nodes to their Phylax public keys.

Once an association has been created between the validator and its Phylax Key, the validator will be considered a 'Phylax Validator', and will be awarded 1 consensus voting power. The total voting power is equal to the size of the Consensus Phylax Set, and at least two thirds of the total voting power must vote to create a block.

The Consensus Phylax Set is a distinct term from the Latest Phylax Set. This is because the Phylax Set Upgrade VAA is submitted via a normal Deltaswap Chain transaction. When a new Latest Phylax Set is created, many of the Phylaxs in the new set may not have yet registered as Phylax Validators. Thus, the older Phylax Set must remain marked as the Consensus Set until enough Phylaxs from the new set have registered.

## Validator Registration:

First, validators must be able to join the Deltaswap Chain Tendermint network. Validator registration is identical to the stock Cosmos-SDK design. Validators may bond and unbond as they would for any other Cosmos Chain. However, all validators have 0 consensus voting power, unless they are registered as a Phylax Validator, wherein they will have 1 voting power.

## Mapping Validators to Phylaxs:

Bonded Validators may register as Phylax Validators by submitting a transaction on-chain. This requires the following criteria:

- The validator must be bonded.
- The validator must hash their Validator Address (Operator Address), sign it with one of the Phylax Keys from the Latest Phylax Set (Note: Latest set, not necessarily Consensus Set.), and then submit this signature in a transaction to the RegisterValidatorAsPhylax function.
- The transaction must be signed/sent from the Validator Address.
- The validator must not have already registered as a different Phylax from the same set.

A Phylax Public Key may only be registered to a single validator at a time. If a new validator proof is received for an existing Phylax Validator, the previous entry is overwritten. As an optional defense mechanism, the registration proofs could be limited to only Phylax Keys in the Latest set.

## Phylax Set Upgrades

Phylax Set upgrades are the trickiest operation to handle. When processing the Phylax Set Upgrade, the following steps happen:

- The Latest Phylax Set is changed to the new Phylax Set.
- If all Phylax Keys in the new Latest Phylax Set are registered, the Latest Phylax Set automatically becomes the new Consensus Phylax Set. Otherwise, the Latest Phylax Set will not become the Consensus Phylax Set until this threshold is met.

## Benefits of this implementation:

- Adequately meets the requirement that Phylaxs are responsible for consensus and block production on Deltaswap Chain.
- Relatively robust with regard to chain 'bricks'. If at any point in the life of Deltaswap Chain less than 13 of the Phylaxs in the Consensus Set are registered, the network will deadlock. There will not be enough Phylaxs registered to produce a block, and because no blocks are being produced, no registrations can be completed. This design does not change the Consensus Set unless a sufficient amount of Phylaxs are registered.
- Can swap out a massive set of Phylaxs all at once. Many other (simpler) designs for Phylax set swaps limit the number of Phylaxs which can be changed at once to only 6 to avoid network deadlocks. This design does not have this problem.
- No modifications to Cosmos SDK validator bonding.

### Cons

- Moderate complexity. This is more complicated than the most straightforward implementations, but gains important features and protections to prevent deadlocks.
- Not 100% immune to deadlocks. If less than 13 Phylaxs have valid registrations, the chain will permanently halt. This is prohibitively difficult to prevent with on-chain mechanisms, and unlikely to occur. Performing a simple hard fork in the scenario of a maimed Phylax Validator set is likely the safer and simpler option.
- Avoids some DOS scenarios by only allowing validator registrations for known Phylax Keys.

## Terms & Phrases:

### Phylax

- One of the entities approved to sign VAAs on the Deltaswap network. Phylaxs are identified by the public key which they use to sign VAAs.

### Phylax Set

- A collection of Phylaxs which at one time was approved by the Deltaswap network to produce VAAs. These collections are identified by their sequential 'Set Index'.

### Latest Phylax Set

- The highest index Phylax Set.

### Consensus Phylax Set

- The Phylax Set which is currently being used to produce blocks on Deltaswap Chain. May be different from the Latest Phylax Set.

### Phylax Set Upgrade VAA

- A Deltaswap network VAA which specifies a new Phylax Set. Emitting a new Phylax Set Upgrade VAA is the mechanism which creates a new Phylax Set.

### Validator

- A node on Deltaswap Chain which is connected to the Tendermint peer network.

### Phylax Validator

- A Validator which is currently registered against a Phylax Public Key in the Consensus Phylax Set
