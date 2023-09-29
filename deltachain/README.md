# Wormchain

**Wormchain** is a blockchain built using Cosmos SDK and Tendermint and initially created with [Ignite](https://github.com/ignite).

## Building

Build and install deltachain.  You will need golang version 1.16+ installed.

```
make build/deltachaind
```

## Develop

See [development.md](./development.md)

## How to run the tests

    run "tilt up -- --deltachain"
    cd ./ts-sdk
    npm ci
    npm run build
    cd ../testing/js
    npm ci
    npm run test

## Learn more about Cosmos & Ignite

- [Ignite](https://github.com/ignite)
- [Ignite Docs](https://docs.ignite.com/)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Discord](https://discord.gg/cosmosnetwork)

## Allowlists

Accounts on deltachain are allowlisted.  To be able to submit a tx on deltachain, you must have an account that is either:
* A validator on deltachain that is part of a current or future phylax set, or
* An account that is allowlisted by a current validator on deltachain.

To create or delete an allowlist entry, you use a validator account.  Allowlist entries can become stale,
meaning the owning validators are no longer part of the validator set.  Any validator can delete or replace stale entries.
To manage allowlists, use the `deltachaind` client.
