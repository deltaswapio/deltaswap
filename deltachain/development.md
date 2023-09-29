# Develop

## prerequsites

- Go >= 1.16
- nodejs >= 16

## Building the blockchain

Run

```shell
make
```

This command creates a `build` directory and in particular, the
`build/deltachaind` binary, which can be used to run and interact with the
blockchain.

You can start a local development instance by running

```shell
make run
```

Or equivalently

```shell
./build/deltachaind --home build
```

If you want to reset the blockchain, just run

```shell
make clean
```

Then you can `make run` again.

## Running tests

Golang tests

    make test

Client tests, run against the chain. Deltachain must be running via `make run` or `tilt up -- --deltachain`

    cd ./ts-sdk
    npm ci
    npm run build
    cd ../testing/js
    npm ci
    npm run test

## Interacting with the blockchain

You can interact with the blockchain by using the go binary:

```shell
./build/deltachaind tx --from tiltPhylax --home build
```

Note the flags `--from tiltPhylax --home build`. These have to be passed
in each time you make a transaction (the `tiltPhylax` account is created in
`config.yml`). Queries don't need the `--from` flag.

## Scaffolding stuff with Ignite

Deltachain was initially scaffolded using the [Ignite CLI](https://github.com/ignite) (formerly Starport). Now, we only use Ignite for generating code from protobuf definitions.

To avoid system compatibility issues, we run Ignite using docker. The below commands should be run using the ignite docker container (see the Makefile recipes for examples).

```shell
ignite scaffold type phylax-key key:string --module deltaswap --no-message
```

modify `proto/deltaswap/phylax_key.proto` (string -> bytes)

```shell
ignite scaffold message register-account-as-phylax phylax-pubkey:PhylaxKey address-bech32:string signature:string --desc "Register a phylax public key with a deltaswap chain address." --module deltaswap --signer signer
```

Scaffold a query:

```shell
ignite scaffold query latest_phylax_set_index --response LatestPhylaxSetIndex --module wormhole
```

(then modify "deltachain/x/deltaswap/types/query.pb.go" to change the response type)
