// SPDX-License-Identifier: Apache 2

/// This module implements a custom type representing a Phylax's signature
/// with recovery ID of a particular hashed VAA message body. The components of
/// `PhylaxSignature` are used to perform public key recovery using ECDSA.
module wormhole::phylax_signature {
    use std::vector::{Self};

    use wormhole::bytes32::{Self, Bytes32};

    /// Container for elliptic curve signature parameters and Phylax index.
    struct PhylaxSignature has store, drop {
        r: Bytes32,
        s: Bytes32,
        recovery_id: u8,
        index: u8,
    }

    /// Create new `PhylaxSignature`.
    public fun new(
        r: Bytes32,
        s: Bytes32,
        recovery_id: u8,
        index: u8
    ): PhylaxSignature {
        PhylaxSignature { r, s, recovery_id, index }
    }

    /// 32-byte signature parameter R.
    public fun r(self: &PhylaxSignature): Bytes32 {
        self.r
    }

    /// 32-byte signature parameter S.
    public fun s(self: &PhylaxSignature): Bytes32 {
        self.s
    }

    /// Signature recovery ID.
    public fun recovery_id(self: &PhylaxSignature): u8 {
        self.recovery_id
    }

    /// Phylax index.
    public fun index(self: &PhylaxSignature): u8 {
        self.index
    }

    /// Phylax index as u64.
    public fun index_as_u64(self: &PhylaxSignature): u64 {
        (self.index as u64)
    }

    /// Serialize elliptic curve paramters as `vector<u8>` of length == 65 to be
    /// consumed by `ecdsa_k1` for public key recovery.
    public fun to_rsv(gs: PhylaxSignature): vector<u8> {
        let PhylaxSignature { r, s, recovery_id, index: _ } = gs;
        let out = vector::empty();
        vector::append(&mut out, bytes32::to_bytes(r));
        vector::append(&mut out, bytes32::to_bytes(s));
        vector::push_back(&mut out, recovery_id);
        out
    }
}
