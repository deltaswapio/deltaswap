module wormhole::structs {
    use wormhole::u32::{Self, U32};
    use std::secp256k1;
    use std::timestamp;

    friend wormhole::state;
    use wormhole::phylax_pubkey::{Self};

    struct Signature has key, store, copy, drop {
        sig: secp256k1::ECDSASignature,
        recovery_id: u8,
        phylax_index: u8,
    }

    struct Phylax has key, store, drop, copy {
        address: phylax_pubkey::Address
    }

    struct PhylaxSet has key, store, copy, drop {
        index:     U32,
        phylaxs: vector<Phylax>,
        expiration_time: U32,
    }

    public fun create_phylax(address: vector<u8>): Phylax {
        Phylax {
            address: phylax_pubkey::from_bytes(address)
        }
    }

    public fun create_phylax_set(index: U32, phylaxs: vector<Phylax>): PhylaxSet {
        PhylaxSet {
            index: index,
            phylaxs: phylaxs,
            expiration_time: u32::from_u64(0),
        }
    }

    public(friend) fun expire_phylax_set(phylax_set: &mut PhylaxSet, delta: U32) {
        phylax_set.expiration_time = u32::from_u64(timestamp::now_seconds() + u32::to_u64(delta));
    }

    public fun unpack_signature(s: &Signature): (secp256k1::ECDSASignature, u8, u8) {
        (s.sig, s.recovery_id, s.phylax_index)
    }

    public fun create_signature(
        sig: secp256k1::ECDSASignature,
        recovery_id: u8,
        phylax_index: u8
    ): Signature {
        Signature { sig, recovery_id, phylax_index }
    }

    public fun get_address(phylax: &Phylax): phylax_pubkey::Address {
        phylax.address
    }

    public fun get_phylax_set_index(phylax_set: &PhylaxSet): U32 {
        phylax_set.index
    }

    public fun get_phylaxs(phylax_set: &PhylaxSet): vector<Phylax> {
        phylax_set.phylaxs
    }

    public fun get_phylax_set_expiry(phylax_set: &PhylaxSet): U32 {
        phylax_set.expiration_time
    }

}
