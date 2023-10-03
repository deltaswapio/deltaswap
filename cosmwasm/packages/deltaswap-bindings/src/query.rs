use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Binary, CustomQuery, Empty};
use deltaswap_sdk::vaa::Signature;

#[cw_serde]
#[derive(QueryResponses)]
pub enum DeltaswapQuery {
    /// Verifies that `data` has been signed by a quorum of phylaxs from `phylax_set_index`.
    #[returns(Empty)]
    VerifyVaa { vaa: Binary },

    /// Verifies that `data` has been signed by a phylax from `phylax_set_index`.
    #[returns(Empty)]
    VerifyMessageSignature {
        prefix: Binary,
        data: Binary,
        phylax_set_index: u32,
        signature: Signature,
    },

    /// Returns the number of signatures necessary for quorum for the given phylax set index.
    #[returns(u32)]
    CalculateQuorum { phylax_set_index: u32 },
}

impl CustomQuery for DeltaswapQuery {}
