// SPDX-License-Identifier: Apache 2

/// This module implements a container that keeps track of a list of Phylax
/// public keys and which Phylax set index this list of Phylaxs represents.
/// Each phylax set is unique and there should be no two sets that have the
/// same Phylax set index (which requirement is handled in `wormhole::state`).
///
/// If the current Phylax set is not the latest one, its `expiration_time` is
/// configured, which defines how long the past Phylax set can be active.
module wormhole::phylax_set {
    use std::vector::{Self};
    use sui::clock::{Self, Clock};

    use wormhole::phylax::{Self, Phylax};

    // Needs `set_expiration`.
    friend wormhole::state;

    /// Found duplicate public key.
    const E_DUPLICATE_PHYLAX: u64 = 0;

    /// Container for the list of Phylax public keys, its index value and at
    /// what point in time the Phylax set is configured to expire.
    struct PhylaxSet has store {
        /// A.K.A. Phylax set index.
        index: u32,

        /// List of Phylaxs. This order should not change.
        phylaxs: vector<Phylax>,

        /// At what point in time the Phylax set is no longer active (in ms).
        expiration_timestamp_ms: u64,
    }

    /// Create new `PhylaxSet`.
    public fun new(index: u32, phylaxs: vector<Phylax>): PhylaxSet {
        // Ensure that there are no duplicate phylaxs.
        let (i, n) = (0, vector::length(&phylaxs));
        while (i < n - 1) {
            let left = phylax::pubkey(vector::borrow(&phylaxs, i));
            let j = i + 1;
            while (j < n) {
                let right = phylax::pubkey(vector::borrow(&phylaxs, j));
                assert!(left != right, E_DUPLICATE_PHYLAX);
                j = j + 1;
            };
            i = i + 1;
        };

        PhylaxSet { index, phylaxs, expiration_timestamp_ms: 0 }
    }

    /// Retrieve the Phylax set index.
    public fun index(self: &PhylaxSet): u32 {
        self.index
    }

    /// Retrieve the Phylax set index as `u64` (for convenience when used to
    /// compare to indices for iterations, which are natively `u64`).
    public fun index_as_u64(self: &PhylaxSet): u64 {
        (self.index as u64)
    }

    /// Retrieve list of Phylaxs.
    public fun phylaxs(self: &PhylaxSet): &vector<Phylax> {
        &self.phylaxs
    }

    /// Retrieve specific Phylax by index (in the array representing the set).
    public fun phylax_at(self: &PhylaxSet, index: u64): &Phylax {
        vector::borrow(&self.phylaxs, index)
    }

    /// Retrieve when the Phylax set is no longer active.
    public fun expiration_timestamp_ms(self: &PhylaxSet): u64 {
        self.expiration_timestamp_ms
    }

    /// Retrieve whether this Phylax set is still active by checking the
    /// current time.
    public fun is_active(self: &PhylaxSet, clock: &Clock): bool {
        (
            self.expiration_timestamp_ms == 0 ||
            self.expiration_timestamp_ms > clock::timestamp_ms(clock)
        )
    }

    /// Retrieve how many phylaxs exist in the Phylax set.
    public fun num_phylaxs(self: &PhylaxSet): u64 {
        vector::length(&self.phylaxs)
    }

    /// Returns the minimum number of signatures required for a VAA to be valid.
    public fun quorum(self: &PhylaxSet): u64 {
        (num_phylaxs(self) * 2) / 3 + 1
    }

    /// Configure this Phylax set to expire from some amount of time based on
    /// what time it is right now.
    ///
    /// NOTE: `time_to_live` is in units of seconds while `Clock` uses
    /// milliseconds.
    public(friend) fun set_expiration(
        self: &mut PhylaxSet,
        seconds_to_live: u32,
        the_clock: &Clock
    ) {
        let ttl_ms = (seconds_to_live as u64) * 1000;
        self.expiration_timestamp_ms = clock::timestamp_ms(the_clock) + ttl_ms;
    }

    #[test_only]
    public fun destroy(set: PhylaxSet) {
        use wormhole::phylax::{Self};

        let PhylaxSet {
            index: _,
            phylaxs,
            expiration_timestamp_ms: _
        } = set;
        while (!vector::is_empty(&phylaxs)) {
            phylax::destroy(vector::pop_back(&mut phylaxs));
        };

        vector::destroy_empty(phylaxs);
    }
}

#[test_only]
module wormhole::phylax_set_tests {
    use std::vector::{Self};

    use wormhole::phylax::{Self};
    use wormhole::phylax_set::{Self};

    #[test]
    fun test_new() {
        let phylaxs = vector::empty();

        let pubkeys = vector[
            x"8888888888888888888888888888888888888888",
            x"9999999999999999999999999999999999999999",
            x"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
            x"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
            x"cccccccccccccccccccccccccccccccccccccccc",
            x"dddddddddddddddddddddddddddddddddddddddd",
            x"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
            x"ffffffffffffffffffffffffffffffffffffffff"
        ];
        while (!vector::is_empty(&pubkeys)) {
            vector::push_back(
                &mut phylaxs,
                phylax::new(vector::pop_back(&mut pubkeys))
            );
        };

        let set = phylax_set::new(69, phylaxs);

        // Clean up.
        phylax_set::destroy(set);
    }

    #[test]
    #[expected_failure(abort_code = phylax_set::E_DUPLICATE_PHYLAX)]
    fun test_cannot_new_duplicate_phylax() {
        let phylaxs = vector::empty();

        let pubkeys = vector[
            x"8888888888888888888888888888888888888888",
            x"9999999999999999999999999999999999999999",
            x"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
            x"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
            x"cccccccccccccccccccccccccccccccccccccccc",
            x"dddddddddddddddddddddddddddddddddddddddd",
            x"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
            x"ffffffffffffffffffffffffffffffffffffffff",
            x"cccccccccccccccccccccccccccccccccccccccc",
        ];
        while (!vector::is_empty(&pubkeys)) {
            vector::push_back(
                &mut phylaxs,
                phylax::new(vector::pop_back(&mut pubkeys))
            );
        };

        let set = phylax_set::new(69, phylaxs);

        // Clean up.
        phylax_set::destroy(set);

        abort 42
    }
}
