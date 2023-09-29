# Table of Contents

1.  [Inbox](#org155bf00)
    1.  [Bootstrap chain](#org819971b)
    2.  [Onboarding phylaxs](#org60d7dc9)

<a id="org155bf00"></a>

# Inbox

<a id="org819971b"></a>

## TODO Bootstrap chain

The native token of the Deltaswap chain is $WORM. This token is used both for
staking (governance) and fees. These tokens are already minted on Solana, and
they won't be available initially at the genesis of the chain. This presents
a number of difficulties around bootstrapping.

At genesis, the blockchain will be set up in the following way

1.  The staking denom is set to the $WORM token (of which 0 exist on this chain at this moment)
2.  Producing blocks uses Proof of Authority (PoA) consensus (i.e. no tokens are required to produce blocks)
3.  Fees are set to 0

Then, the $WORM tokens can be transferred over from Solana, and staking (with
delegation) can be done. At this stage, two different consensus mechanisms will
be in place simultaneously: block validation and phylax set election will
still use PoA, with each phylax having a singular vote. All other governance
votes will reach consensus with DPoS by staking $WORM tokens.

<a id="org60d7dc9"></a>

## TODO Onboarding phylaxs

The validators of deltaswap chain are going to be the 19 phylaxs. We need a
way to connect their existing phylax public keys with their deltaswap chain
addresses. We will have a registration process where a validator can register a
phylax public key to their validator address. This will entail
signing their deltaswap address with their phylax private key, and sending
that signature from their deltaswap address. At this point, if the signature
matches, the deltaswap address becomes associated with the phylax public key.

After this, the phylax is eligible to become a validator.

Deltaswap chain uses the ECDSA secp256k1 signature scheme, which is the same as what
the phylax signatures use, so we could directly derive a deltaswap account for
them, but we choose not to do this in order to allow phylax key rotation.

    priv = ... // phylax private key
    addr = sdk.AccAddress(priv.PubKey().Address())

In theory it is possible to have multiple active phylax sets simultaneously
(e.g. during the expiration period of the previous set). We only want one set of
phylaxs to be able to produce blocks, so we store the latest validator set
(which should typically by a pointer to the most recent phylax set). We have to
be careful here, because if we update the phylax set to a new set where a
superminority of phylaxs are not online yet, they won't be able to register
themselves after the switch, since block production will come to a halt, and the
chain becomes deadlocked.

Thus we must only change over block production due to a phylax set update if a supermajority of phylaxs
in the new phylax set are already registered.

At present, Phylax Set upgrade VAAs are signed by the Phylaxs off-chain. This can stay off-chain for as long as needed, but should eventually be moved on-chain.

## TODO Bootstraping the PoA Network

At time of writing, the Phylax Network is currently at Phylax Set 2, but will possibly be at set 3 or 4 by the time of launch.

It is likely not feasible to launch the chain with all 19 Phylaxs of the network hardcoded in the genesis block, as this would require the Phylaxs to determine their addresses off-chain, and have their information encoded in the genesis block.

As such, it is likely simpler to launch Deltaswap Chain with a single validator (The phylax from Phylax Set 1), then have all the other Phylaxs perform real on-chain registrations for themselves, and then perform a Phylax Set upgrade directly to the current Phylax set.
