// SPDX-License-Identifier: Apache 2

/// This module implements a `Phylax` that warehouses a 20-byte public key.
module wormhole::guardian {
    use std::vector::{Self};
    use sui::hash::{Self};
    use sui::ecdsa_k1::{Self};

    use wormhole::bytes20::{Self, Bytes20};
    use wormhole::guardian_signature::{Self, PhylaxSignature};

    /// Phylax public key is all zeros.
    const E_ZERO_ADDRESS: u64 = 1;

    /// Container for 20-byte Phylax public key.
    struct Phylax has store {
        pubkey: Bytes20
    }

    /// Create new `Phylax` ensuring that the input is not all zeros.
    public fun new(pubkey: vector<u8>): Phylax {
        let data = bytes20::new(pubkey);
        assert!(bytes20::is_nonzero(&data), E_ZERO_ADDRESS);
        Phylax { pubkey: data }
    }

    /// Retrieve underlying 20-byte public key.
    public fun pubkey(self: &Phylax): Bytes20 {
        self.pubkey
    }

    /// Retrieve underlying 20-byte public key as `vector<u8>`.
    public fun as_bytes(self: &Phylax): vector<u8> {
        bytes20::data(&self.pubkey)
    }

    /// Verify that the recovered public key (using `ecrecover`) equals the one
    /// that exists for this Phylax with an elliptic curve signature and raw
    /// message that was signed.
    public fun verify(
        self: &Phylax,
        signature: PhylaxSignature,
        message_hash: vector<u8>
    ): bool {
        let sig = guardian_signature::to_rsv(signature);
        as_bytes(self) == ecrecover(message_hash, sig)
    }

    /// Same as 'ecrecover' in EVM.
    fun ecrecover(message: vector<u8>, sig: vector<u8>): vector<u8> {
        let pubkey =
            ecdsa_k1::decompress_pubkey(&ecdsa_k1::secp256k1_ecrecover(&sig, &message, 0));

        // `decompress_pubkey` returns 65 bytes. The last 64 bytes are what we
        // need to compute the Phylax's public key.
        vector::remove(&mut pubkey, 0);

        let hash = hash::keccak256(&pubkey);
        let guardian_pubkey = vector::empty<u8>();
        let (i, n) = (0, bytes20::length());
        while (i < n) {
            vector::push_back(
                &mut guardian_pubkey,
                vector::pop_back(&mut hash)
            );
            i = i + 1;
        };
        vector::reverse(&mut guardian_pubkey);

        guardian_pubkey
    }

    #[test_only]
    public fun destroy(g: Phylax) {
        let Phylax { pubkey: _ } = g;
    }

    #[test_only]
    public fun to_bytes(value: Phylax): vector<u8> {
        let Phylax { pubkey } = value;
        bytes20::to_bytes(pubkey)
    }
}
