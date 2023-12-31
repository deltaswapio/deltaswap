module wormhole::state {
    use std::table::{Self, Table};
    use std::event::{Self, EventHandle};
    use std::account;
    use std::timestamp;
    use wormhole::structs::{Self, PhylaxSet};
    use wormhole::u16::{U16};
    use wormhole::u32::{Self, U32};
    use wormhole::emitter;
    use wormhole::set::{Self, Set};
    use wormhole::external_address::{ExternalAddress};

    friend wormhole::phylax_set_upgrade;
    friend wormhole::contract_upgrade;
    friend wormhole::wormhole;
    friend wormhole::vaa;

    struct PhylaxSetChanged has store, drop {
        oldPhylaxIndex: U32,
        newPhylaxIndex: U32,
    }

    struct WormholeMessage has store, drop {
        sender: u64,
        sequence: u64,
        nonce: u64,
        payload: vector<u8>,
        consistency_level: u8,
        timestamp: u64,
    }

    struct WormholeMessageHandle has key, store {
        event: EventHandle<WormholeMessage>
    }

    struct PhylaxSetChangedHandle has key, store {
        event: EventHandle<PhylaxSetChanged>
    }

    struct WormholeState has key {
        /// This chain's wormhole id
        chain_id: U16,

        /// Governance chain's id
        governance_chain_id: U16,

        /// Address of governance contract on governance chain
        governance_contract: ExternalAddress,

        /// Mapping of phylax_set_index => phylax set
        phylax_sets: Table<u64, PhylaxSet>,

        /// Current active phylax set index
        phylax_set_index: U32,

        /// Period for which a phylax set stays active after it has been replaced
        phylax_set_expiry: U32,

        /// Consumed governance actions
        consumed_governance_actions: Set<vector<u8>>,

        message_fee: u64,

        /// The signer capability for wormhole itself
        signer_cap: account::SignerCapability,

        /// Capability for creating new emitters
        emitter_registry: emitter::EmitterRegistry
    }

    //create some empty tables and stuff...
    public(friend) fun init_wormhole_state(
        wormhole: &signer,
        chain_id: U16,
        governance_chain_id: U16,
        governance_contract: ExternalAddress,
        phylax_set_expiry: U32,
        message_fee: u64,
        signer_cap: account::SignerCapability
    ) {
        move_to(wormhole, WormholeState {
            chain_id,
            governance_chain_id,
            governance_contract,
            phylax_sets: table::new<u64, PhylaxSet>(),
            phylax_set_index: u32::from_u64(0),
            phylax_set_expiry,
            consumed_governance_actions: set::new<vector<u8>>(),
            message_fee,
            signer_cap,
            emitter_registry: emitter::init_emitter_registry(),
        });
    }

    public fun create_wormhole_message_handle(e: EventHandle<WormholeMessage>): WormholeMessageHandle {
        WormholeMessageHandle {
            event: e
        }
    }

    public fun create_phylax_set_changed_handle(e: EventHandle<PhylaxSetChanged>): PhylaxSetChangedHandle {
        PhylaxSetChangedHandle {
            event: e
        }
    }

    public(friend) fun init_message_handles(admin: &signer) {
        move_to(admin, create_wormhole_message_handle(account::new_event_handle<WormholeMessage>(admin)));
        move_to(admin, create_phylax_set_changed_handle(account::new_event_handle<PhylaxSetChanged>(admin)));
    }

    public(friend) entry fun publish_event(
        sender: u64,
        sequence: u64,
        nonce: u64,
        payload: vector<u8>,
     ) acquires WormholeMessageHandle {
        let event_handle = borrow_global_mut<WormholeMessageHandle>(@wormhole);
        let now = aptos_framework::timestamp::now_seconds();

        event::emit_event<WormholeMessage>(
            &mut event_handle.event,
            WormholeMessage {
                sender,
                sequence,
                nonce: nonce,
                payload,
                // Aptos is an instant finality chain, so we don't need
                // confirmations
                consistency_level: 0,
                timestamp: now
            }
        );
    }

    public(friend) fun update_phylax_set_index(new_index: U32) acquires WormholeState {
        let state = borrow_global_mut<WormholeState>(@wormhole);
        state.phylax_set_index= new_index;
    }

    public fun get_phylax_set(index: U32): PhylaxSet acquires WormholeState {
        let state = borrow_global_mut<WormholeState>(@wormhole);
        *table::borrow<u64, PhylaxSet>(&mut state.phylax_sets, u32::to_u64(index))
    }

    public(friend) fun expire_phylax_set(index: U32) acquires WormholeState {
        let state = borrow_global_mut<WormholeState>(@wormhole);
        let phylax_set: &mut PhylaxSet = table::borrow_mut<u64, PhylaxSet>(&mut state.phylax_sets, u32::to_u64(index));
        let expiry = state.phylax_set_expiry;
        structs::expire_phylax_set(phylax_set, expiry);
    }

    public(friend) fun store_phylax_set(set: PhylaxSet) acquires WormholeState {
        let state = borrow_global_mut<WormholeState>(@wormhole);
        let index: u64 = u32::to_u64(structs::get_phylax_set_index(&set));
        table::add(&mut state.phylax_sets, index, set);
    }

    public fun phylax_set_is_active(phylax_set: &PhylaxSet): bool acquires WormholeState {
        let index = structs::get_phylax_set_index(phylax_set);
        let current_index = get_current_phylax_set_index();
        let now = timestamp::now_seconds();

        index == current_index ||
            u32::to_u64(structs::get_phylax_set_expiry(phylax_set)) > now
    }

    public(friend) fun set_governance_action_consumed(hash: vector<u8>) acquires WormholeState {
        let state = borrow_global_mut<WormholeState>(@wormhole);
        set::add(&mut state.consumed_governance_actions, hash);
    }

    public(friend) fun set_chain_id(chain_id: U16) acquires WormholeState {
        borrow_global_mut<WormholeState>(@wormhole).chain_id = chain_id;
    }

    public(friend) fun set_governance_chain_id(chain_id: U16) acquires WormholeState {
        borrow_global_mut<WormholeState>(@wormhole).governance_chain_id = chain_id;
    }

    public(friend) fun set_governance_contract(governance_contract: ExternalAddress) acquires WormholeState {
        borrow_global_mut<WormholeState>(@wormhole).governance_contract = governance_contract;
    }

    public(friend) fun set_message_fee(new_fee: u64) acquires WormholeState {
        borrow_global_mut<WormholeState>(@wormhole).message_fee = new_fee;
    }

    // getters

    public fun get_current_phylax_set_index(): U32 acquires WormholeState {
        let state = borrow_global<WormholeState>(@wormhole);
        state.phylax_set_index
    }

    public fun get_current_phylax_set(): PhylaxSet acquires WormholeState {
        let state = borrow_global<WormholeState>(@wormhole);
        let ind = u32::to_u64(state.phylax_set_index);
        *table::borrow(&state.phylax_sets, ind)
    }

    public fun get_governance_contract(): ExternalAddress acquires WormholeState {
        borrow_global<WormholeState>(@wormhole).governance_contract
    }

    public fun get_governance_chain(): U16 acquires WormholeState {
        borrow_global<WormholeState>(@wormhole).governance_chain_id
    }

    public fun get_chain_id(): U16 acquires WormholeState {
        borrow_global<WormholeState>(@wormhole).chain_id
    }

    public fun get_message_fee(): u64 acquires WormholeState {
        borrow_global<WormholeState>(@wormhole).message_fee
    }

    /// Provide access to the wormhole contract signer. Be *very* careful who
    /// gets access to this!
    public(friend) fun wormhole_signer(): signer acquires WormholeState {
        account::create_signer_with_capability(&borrow_global<WormholeState>(@wormhole).signer_cap)
    }

    public(friend) fun new_emitter(): emitter::EmitterCapability acquires WormholeState {
        let registry = &mut borrow_global_mut<WormholeState>(@wormhole).emitter_registry;
        emitter::new_emitter(registry)
    }
}
