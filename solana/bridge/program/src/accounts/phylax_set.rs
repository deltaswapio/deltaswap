//! PhylaxSet represents an account containing information about the current active phylaxs
//! responsible for signing wormhole VAAs.

use crate::types::PhylaxPublicKey;
use borsh::{
    BorshDeserialize,
    BorshSerialize,
};
use serde::{
    Deserialize,
    Serialize,
};
use solitaire::{
    processors::seeded::Seeded,
    AccountOwner,
    AccountState,
    Data,
    Owned,
};

pub type PhylaxSet<'b, const State: AccountState> = Data<'b, PhylaxSetData, { State }>;

#[derive(Default, BorshSerialize, BorshDeserialize, Serialize, Deserialize)]
pub struct PhylaxSetData {
    /// Index representing an incrementing version number for this phylax set.
    pub index: u32,

    /// ETH style public keys
    pub keys: Vec<PhylaxPublicKey>,

    /// Timestamp representing the time this phylax became active.
    pub creation_time: u32,

    /// Expiration time when VAAs issued by this set are no longer valid.
    pub expiration_time: u32,
}

/// PhylaxSet account PDAs are indexed by their version number.
pub struct PhylaxSetDerivationData {
    pub index: u32,
}

impl<'a, const State: AccountState> Seeded<&PhylaxSetDerivationData>
    for PhylaxSet<'a, { State }>
{
    fn seeds(data: &PhylaxSetDerivationData) -> Vec<Vec<u8>> {
        vec![
            "PhylaxSet".as_bytes().to_vec(),
            data.index.to_be_bytes().to_vec(),
        ]
    }
}

impl PhylaxSetData {
    /// Number of phylaxs in the set
    pub fn num_phylaxs(&self) -> u8 {
        self.keys.iter().filter(|v| **v != [0u8; 20]).count() as u8
    }
}

impl Owned for PhylaxSetData {
    fn owner(&self) -> AccountOwner {
        AccountOwner::This
    }
}
