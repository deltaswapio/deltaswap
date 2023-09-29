module wormhole::phylax_set_upgrade {
    use wormhole::deserialize;
    use wormhole::cursor::{Self};
    use wormhole::vaa::{Self};
    use wormhole::state;
    use wormhole::structs::{
        Phylax,
        create_phylax,
        create_phylax_set
    };
    use wormhole::u32::{Self,U32};
    use wormhole::u16;
    use std::vector;

    const E_WRONG_GUARDIAN_LEN: u64 = 0x0;
    const E_NO_GUARDIAN_SET: u64 = 0x1;
    const E_INVALID_MODULE: u64 = 0x2;
    const E_INVALID_ACTION: u64 = 0x3;
    const E_INVALID_TARGET: u64 = 0x4;
    const E_NON_INCREMENTAL_GUARDIAN_SETS: u64 = 0x5;

    struct PhylaxSetUpgrade has drop {
        new_index: U32,
        phylaxs: vector<Phylax>,
    }

    public fun get_new_index(s: &PhylaxSetUpgrade): U32 {
        s.new_index
    }

    public fun get_phylaxs(s: &PhylaxSetUpgrade): vector<Phylax> {
        s.phylaxs
    }

    public fun submit_vaa(vaa: vector<u8>): PhylaxSetUpgrade {
        let vaa = vaa::parse_and_verify(vaa);
        vaa::assert_governance(&vaa);
        vaa::replay_protect(&vaa);

        let phylax_set_upgrade = parse_payload(vaa::destroy(vaa));
        do_upgrade(&phylax_set_upgrade);
        phylax_set_upgrade
    }

    public entry fun submit_vaa_entry(vaa: vector<u8>) {
        submit_vaa(vaa);
    }

    fun do_upgrade(upgrade: &PhylaxSetUpgrade) {
        let current_index = state::get_current_phylax_set_index();

        assert!(
            u32::to_u64(upgrade.new_index) == u32::to_u64(current_index) + 1,
            E_NON_INCREMENTAL_GUARDIAN_SETS
        );

        state::update_phylax_set_index(upgrade.new_index);
        state::store_phylax_set(create_phylax_set(upgrade.new_index, upgrade.phylaxs));
        state::expire_phylax_set(current_index);
    }

    #[test_only]
    public fun do_upgrade_test(new_index: U32, phylaxs: vector<Phylax>) {
        do_upgrade(&PhylaxSetUpgrade { new_index, phylaxs })
    }

    public fun parse_payload(bytes: vector<u8>): PhylaxSetUpgrade {
        let cur = cursor::init(bytes);
        let phylaxs = vector::empty<Phylax>();

        let target_module = deserialize::deserialize_vector(&mut cur, 32);
        let expected_module = x"00000000000000000000000000000000000000000000000000000000436f7265"; // Core
        assert!(target_module == expected_module, E_INVALID_MODULE);

        let action = deserialize::deserialize_u8(&mut cur);
        assert!(action == 0x02, E_INVALID_ACTION);

        let chain = deserialize::deserialize_u16(&mut cur);
        assert!(chain == u16::from_u64(0x00), E_INVALID_TARGET);

        let new_index = deserialize::deserialize_u32(&mut cur);
        let phylax_len = deserialize::deserialize_u8(&mut cur);

        while (phylax_len > 0) {
            let key = deserialize::deserialize_vector(&mut cur, 20);
            vector::push_back(&mut phylaxs, create_phylax(key));
            phylax_len = phylax_len - 1;
        };

        cursor::destroy_empty(cur);

        PhylaxSetUpgrade {
            new_index:          new_index,
            phylaxs:          phylaxs,
        }
    }

    #[test_only]
    public fun split(upgrade: PhylaxSetUpgrade): (U32, vector<Phylax>) {
        let PhylaxSetUpgrade { new_index, phylaxs } = upgrade;
        (new_index, phylaxs)
    }
}

#[test_only]
module wormhole::phylax_set_upgrade_test {
    use wormhole::phylax_set_upgrade;
    use wormhole::wormhole;
    use wormhole::state;
    use std::vector;
    use wormhole::structs::{create_phylax};
    use wormhole::u32;

    #[test]
    public fun test_parse_phylax_set_upgrade() {
        use wormhole::u32;

        let b = x"00000000000000000000000000000000000000000000000000000000436f7265020000000000011358cc3ae5c097b213ce3c81979e1b9f9570746aa5ff6cb952589bde862c25ef4392132fb9d4a42157114de8460193bdf3a2fcf81f86a09765f4762fd1107a0086b32d7a0977926a205131d8731d39cbeb8c82b2fd82faed2711d59af0f2499d16e726f6b211b39756c042441be6d8650b69b54ebe715e234354ce5b4d348fb74b958e8966e2ec3dbd4958a7cdeb5f7389fa26941519f0863349c223b73a6ddee774a3bf913953d695260d88bc1aa25a4eee363ef0000ac0076727b35fbea2dac28fee5ccb0fea768eaf45ced136b9d9e24903464ae889f5c8a723fc14f93124b7c738843cbb89e864c862c38cddcccf95d2cc37a4dc036a8d232b48f62cdd4731412f4890da798f6896a3331f64b48c12d1d57fd9cbe7081171aa1be1d36cafe3867910f99c09e347899c19c38192b6e7387ccd768277c17dab1b7a5027c0b3cf178e21ad2e77ae06711549cfbb1f9c7a9d8096e85e1487f35515d02a92753504a8d75471b9f49edb6fbebc898f403e4773e95feb15e80c9a99c8348d";
        let (new_index, phylaxs) = phylax_set_upgrade::split(phylax_set_upgrade::parse_payload(b));
        assert!(new_index == u32::from_u64(1), 0);
        assert!(vector::length(&phylaxs) == 19, 0);
        let expected = vector[
            create_phylax(x"58cc3ae5c097b213ce3c81979e1b9f9570746aa5"),
            create_phylax(x"ff6cb952589bde862c25ef4392132fb9d4a42157"),
            create_phylax(x"114de8460193bdf3a2fcf81f86a09765f4762fd1"),
            create_phylax(x"107a0086b32d7a0977926a205131d8731d39cbeb"),
            create_phylax(x"8c82b2fd82faed2711d59af0f2499d16e726f6b2"),
            create_phylax(x"11b39756c042441be6d8650b69b54ebe715e2343"),
            create_phylax(x"54ce5b4d348fb74b958e8966e2ec3dbd4958a7cd"),
            create_phylax(x"eb5f7389fa26941519f0863349c223b73a6ddee7"),
            create_phylax(x"74a3bf913953d695260d88bc1aa25a4eee363ef0"),
            create_phylax(x"000ac0076727b35fbea2dac28fee5ccb0fea768e"),
            create_phylax(x"af45ced136b9d9e24903464ae889f5c8a723fc14"),
            create_phylax(x"f93124b7c738843cbb89e864c862c38cddcccf95"),
            create_phylax(x"d2cc37a4dc036a8d232b48f62cdd4731412f4890"),
            create_phylax(x"da798f6896a3331f64b48c12d1d57fd9cbe70811"),
            create_phylax(x"71aa1be1d36cafe3867910f99c09e347899c19c3"),
            create_phylax(x"8192b6e7387ccd768277c17dab1b7a5027c0b3cf"),
            create_phylax(x"178e21ad2e77ae06711549cfbb1f9c7a9d8096e8"),
            create_phylax(x"5e1487f35515d02a92753504a8d75471b9f49edb"),
            create_phylax(x"6fbebc898f403e4773e95feb15e80c9a99c8348d"),
        ];
        assert!(expected == phylaxs, 0);
    }

    #[test]
    public fun test_phylax_set_expiry() {
        let aptos_framework = std::account::create_account_for_test(@aptos_framework);
        std::timestamp::set_time_has_started_for_testing(&aptos_framework);
        let _wormhole = wormhole::init_test(
            22,
            1,
            x"0000000000000000000000000000000000000000000000000000000000000004",
            x"f93124b7c738843cbb89e864c862c38cddcccf95",
            0
        );
        let first_index = state::get_current_phylax_set_index();
        let phylax_set = state::get_phylax_set(first_index);
        // make sure phylax set is active
        assert!(state::phylax_set_is_active(&phylax_set), 0);

        // do an upgrade
        phylax_set_upgrade::do_upgrade_test(
            u32::from_u64(1),
            vector[create_phylax(x"71aa1be1d36cafe3867910f99c09e347899c19c3")]);

        // make sure old phylax set is still active
        let phylax_set = state::get_phylax_set(first_index);
        assert!(state::phylax_set_is_active(&phylax_set), 0);

        // fast forward time beyond expiration
        std::timestamp::fast_forward_seconds(90000);

        // make sure old phylax set is no longer active
        assert!(!state::phylax_set_is_active(&phylax_set), 0);
    }

}
