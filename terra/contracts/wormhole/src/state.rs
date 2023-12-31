use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{Binary, CanonicalAddr, Coin, StdResult, Storage, Uint128};
use cosmwasm_storage::{
    bucket, bucket_read, singleton, singleton_read, Bucket, ReadonlyBucket, ReadonlySingleton,
    Singleton,
};

use crate::{byte_utils::ByteUtils, error::ContractError};

use sha3::{Digest, Keccak256};

type HumanAddr = String;

pub static CONFIG_KEY: &[u8] = b"config";
pub static PHYLAX_SET_KEY: &[u8] = b"phylax_set";
pub static SEQUENCE_KEY: &[u8] = b"sequence";
pub static WRAPPED_ASSET_KEY: &[u8] = b"wrapped_asset";
pub static WRAPPED_ASSET_ADDRESS_KEY: &[u8] = b"wrapped_asset_address";

// Phylax set information
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct ConfigInfo {
    // Current active phylax set
    pub phylax_set_index: u32,

    // Period for which a phylax set stays active after it has been replaced
    pub phylax_set_expirity: u64,

    // governance contract details
    pub gov_chain: u16,
    pub gov_address: Vec<u8>,

    // Message sending fee
    pub fee: Coin,
}

// Validator Action Approval(VAA) data
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct ParsedVAA {
    pub version: u8,
    pub phylax_set_index: u32,
    pub timestamp: u32,
    pub nonce: u32,
    pub len_signers: u8,

    pub emitter_chain: u16,
    pub emitter_address: Vec<u8>,
    pub sequence: u64,
    pub consistency_level: u8,
    pub payload: Vec<u8>,

    pub hash: Vec<u8>,
}

impl ParsedVAA {
    /* VAA format:

    header (length 6):
    0   uint8   version (0x01)
    1   uint32  phylax set index
    5   uint8   len signatures

    per signature (length 66):
    0   uint8       index of the signer (in phylax keys)
    1   [65]uint8   signature

    body:
    0   uint32      timestamp (unix in seconds)
    4   uint32      nonce
    8   uint16      emitter_chain
    10  [32]uint8   emitter_address
    42  uint64      sequence
    50  uint8       consistency_level
    51  []uint8     payload
    */

    pub const HEADER_LEN: usize = 6;
    pub const SIGNATURE_LEN: usize = 66;

    pub const PHYLAX_SET_INDEX_POS: usize = 1;
    pub const LEN_SIGNER_POS: usize = 5;

    pub const VAA_NONCE_POS: usize = 4;
    pub const VAA_EMITTER_CHAIN_POS: usize = 8;
    pub const VAA_EMITTER_ADDRESS_POS: usize = 10;
    pub const VAA_SEQUENCE_POS: usize = 42;
    pub const VAA_CONSISTENCY_LEVEL_POS: usize = 50;
    pub const VAA_PAYLOAD_POS: usize = 51;

    // Signature data offsets in the signature block
    pub const SIG_DATA_POS: usize = 1;
    // Signature length minus recovery id at the end
    pub const SIG_DATA_LEN: usize = 64;
    // Recovery byte is last after the main signature
    pub const SIG_RECOVERY_POS: usize = Self::SIG_DATA_POS + Self::SIG_DATA_LEN;

    pub fn deserialize(data: &[u8]) -> StdResult<Self> {
        let version = data.get_u8(0);

        // Load 4 bytes starting from index 1
        let phylax_set_index: u32 = data.get_u32(Self::PHYLAX_SET_INDEX_POS);
        let len_signers = data.get_u8(Self::LEN_SIGNER_POS) as usize;
        let body_offset: usize = Self::HEADER_LEN + Self::SIGNATURE_LEN * len_signers;

        // Hash the body
        if body_offset >= data.len() {
            return ContractError::InvalidVAA.std_err();
        }
        let body = &data[body_offset..];
        let mut hasher = Keccak256::new();
        hasher.update(body);
        let hash = hasher.finalize().to_vec();

        // Rehash the hash
        let mut hasher = Keccak256::new();
        hasher.update(hash);
        let hash = hasher.finalize().to_vec();

        // Signatures valid, apply VAA
        if body_offset + Self::VAA_PAYLOAD_POS > data.len() {
            return ContractError::InvalidVAA.std_err();
        }

        let timestamp = data.get_u32(body_offset);
        let nonce = data.get_u32(body_offset + Self::VAA_NONCE_POS);
        let emitter_chain = data.get_u16(body_offset + Self::VAA_EMITTER_CHAIN_POS);
        let emitter_address = data
            .get_bytes32(body_offset + Self::VAA_EMITTER_ADDRESS_POS)
            .to_vec();
        let sequence = data.get_u64(body_offset + Self::VAA_SEQUENCE_POS);
        let consistency_level = data.get_u8(body_offset + Self::VAA_CONSISTENCY_LEVEL_POS);
        let payload = data[body_offset + Self::VAA_PAYLOAD_POS..].to_vec();

        Ok(ParsedVAA {
            version,
            phylax_set_index,
            timestamp,
            nonce,
            len_signers: len_signers as u8,
            emitter_chain,
            emitter_address,
            sequence,
            consistency_level,
            payload,
            hash,
        })
    }
}

// Phylax address
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PhylaxAddress {
    pub bytes: Binary, // 20-byte addresses
}

use crate::contract::FEE_DENOMINATION;
#[cfg(test)]
use hex;

#[cfg(test)]
impl PhylaxAddress {
    pub fn from(string: &str) -> PhylaxAddress {
        PhylaxAddress {
            bytes: hex::decode(string).expect("Decoding failed").into(),
        }
    }
}

// Phylax set information
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct PhylaxSetInfo {
    pub addresses: Vec<PhylaxAddress>,
    // List of phylax addresses
    pub expiration_time: u64, // Phylax set expiration time
}

impl PhylaxSetInfo {
    pub fn quorum(&self) -> usize {
        // allow quorum of 0 for testing purposes...
        if self.addresses.is_empty() {
            return 0;
        }
        ((self.addresses.len() * 10 / 3) * 2) / 10 + 1
    }
}

// Wormhole contract generic information
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct WormholeInfo {
    // Period for which a phylax set stays active after it has been replaced
    pub phylax_set_expirity: u64,
}

pub fn config(storage: &mut dyn Storage) -> Singleton<ConfigInfo> {
    singleton(storage, CONFIG_KEY)
}

pub fn config_read(storage: &dyn Storage) -> ReadonlySingleton<ConfigInfo> {
    singleton_read(storage, CONFIG_KEY)
}

pub fn phylax_set_set(
    storage: &mut dyn Storage,
    index: u32,
    data: &PhylaxSetInfo,
) -> StdResult<()> {
    bucket(storage, PHYLAX_SET_KEY).save(&index.to_be_bytes(), data)
}

pub fn phylax_set_get(storage: &dyn Storage, index: u32) -> StdResult<PhylaxSetInfo> {
    bucket_read(storage, PHYLAX_SET_KEY).load(&index.to_be_bytes())
}

pub fn sequence_set(storage: &mut dyn Storage, emitter: &[u8], sequence: u64) -> StdResult<()> {
    bucket(storage, SEQUENCE_KEY).save(emitter, &sequence)
}

pub fn sequence_read(storage: &dyn Storage, emitter: &[u8]) -> u64 {
    bucket_read(storage, SEQUENCE_KEY)
        .load(emitter)
        .unwrap_or(0)
}

pub fn vaa_archive_add(storage: &mut dyn Storage, hash: &[u8]) -> StdResult<()> {
    bucket(storage, PHYLAX_SET_KEY).save(hash, &true)
}

pub fn vaa_archive_check(storage: &dyn Storage, hash: &[u8]) -> bool {
    bucket_read(storage, PHYLAX_SET_KEY)
        .load(hash)
        .unwrap_or(false)
}

pub fn wrapped_asset(storage: &mut dyn Storage) -> Bucket<HumanAddr> {
    bucket(storage, WRAPPED_ASSET_KEY)
}

pub fn wrapped_asset_read(storage: &dyn Storage) -> ReadonlyBucket<HumanAddr> {
    bucket_read(storage, WRAPPED_ASSET_KEY)
}

pub fn wrapped_asset_address(storage: &mut dyn Storage) -> Bucket<Vec<u8>> {
    bucket(storage, WRAPPED_ASSET_ADDRESS_KEY)
}

pub fn wrapped_asset_address_read(storage: &dyn Storage) -> ReadonlyBucket<Vec<u8>> {
    bucket_read(storage, WRAPPED_ASSET_ADDRESS_KEY)
}

pub struct GovernancePacket {
    pub module: Vec<u8>,
    pub action: u8,
    pub chain: u16,
    pub payload: Vec<u8>,
}

impl GovernancePacket {
    pub fn deserialize(data: &[u8]) -> StdResult<Self> {
        let module = data.get_bytes32(0).to_vec();
        let action = data.get_u8(32);
        let chain = data.get_u16(33);
        let payload = data[35..].to_vec();

        Ok(GovernancePacket {
            module,
            action,
            chain,
            payload,
        })
    }
}

// action 1
pub struct ContractUpgrade {
    pub new_contract: u64,
}

// action 2
pub struct PhylaxSetUpgrade {
    pub new_phylax_set_index: u32,
    pub new_phylax_set: PhylaxSetInfo,
}

impl ContractUpgrade {
    pub fn deserialize(data: &[u8]) -> StdResult<Self> {
        let new_contract = data.get_u64(24);
        Ok(ContractUpgrade { new_contract })
    }
}

impl PhylaxSetUpgrade {
    pub fn deserialize(data: &[u8]) -> StdResult<Self> {
        const ADDRESS_LEN: usize = 20;

        let new_phylax_set_index = data.get_u32(0);

        let n_phylaxs = data.get_u8(4);

        let mut addresses = vec![];

        for i in 0..n_phylaxs {
            let pos = 5 + (i as usize) * ADDRESS_LEN;
            if pos + ADDRESS_LEN > data.len() {
                return ContractError::InvalidVAA.std_err();
            }

            addresses.push(PhylaxAddress {
                bytes: data[pos..pos + ADDRESS_LEN].to_vec().into(),
            });
        }

        let new_phylax_set = PhylaxSetInfo {
            addresses,
            expiration_time: 0,
        };

        Ok(PhylaxSetUpgrade {
            new_phylax_set_index,
            new_phylax_set,
        })
    }
}

// action 3
pub struct SetFee {
    pub fee: Coin,
}

impl SetFee {
    pub fn deserialize(data: &[u8]) -> StdResult<Self> {
        let (_, amount) = data.get_u256(0);
        let fee = Coin {
            denom: String::from(FEE_DENOMINATION),
            amount: Uint128::new(amount),
        };
        Ok(SetFee { fee })
    }
}

// action 4
pub struct TransferFee {
    pub amount: Coin,
    pub recipient: CanonicalAddr,
}

impl TransferFee {
    pub fn deserialize(data: &[u8]) -> StdResult<Self> {
        let recipient = data.get_address(0);

        let (_, amount) = data.get_u256(32);
        let amount = Coin {
            denom: String::from(FEE_DENOMINATION),
            amount: Uint128::new(amount),
        };
        Ok(TransferFee { amount, recipient })
    }
}
