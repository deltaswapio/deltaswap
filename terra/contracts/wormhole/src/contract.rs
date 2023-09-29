use cosmwasm_std::{
    has_coins, to_binary, BankMsg, Binary, Coin, CosmosMsg, Deps, DepsMut, Env, MessageInfo,
    Response, StdError, StdResult, Storage, WasmMsg,
};

#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;

use crate::{
    byte_utils::{extend_address_to_32, ByteUtils},
    error::ContractError,
    msg::{
        ExecuteMsg, GetAddressHexResponse, GetStateResponse, PhylaxSetInfoResponse,
        InstantiateMsg, MigrateMsg, QueryMsg,
    },
    state::{
        config, config_read, phylax_set_get, phylax_set_set, sequence_read, sequence_set,
        vaa_archive_add, vaa_archive_check, ConfigInfo, ContractUpgrade, GovernancePacket,
        PhylaxAddress, PhylaxSetInfo, PhylaxSetUpgrade, ParsedVAA, SetFee, TransferFee,
    },
};

use k256::{
    ecdsa::{
        recoverable::{Id as RecoverableId, Signature as RecoverableSignature},
        Signature, VerifyingKey,
    },
    EncodedPoint,
};
use sha3::{Digest, Keccak256};

use generic_array::GenericArray;
use std::convert::TryFrom;

type HumanAddr = String;

// Chain ID of Terra
const CHAIN_ID: u16 = 3;

// Lock assets fee amount and denomination
const FEE_AMOUNT: u128 = 0;
pub const FEE_DENOMINATION: &str = "uluna";

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn migrate(_deps: DepsMut, _env: Env, _msg: MigrateMsg) -> StdResult<Response> {
    Ok(Response::default())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: InstantiateMsg,
) -> StdResult<Response> {
    // Save general wormhole info
    let state = ConfigInfo {
        gov_chain: msg.gov_chain,
        gov_address: msg.gov_address.as_slice().to_vec(),
        phylax_set_index: 0,
        phylax_set_expirity: msg.phylax_set_expirity,
        fee: Coin::new(FEE_AMOUNT, FEE_DENOMINATION), // 0.01 Luna (or 10000 uluna) fee by default
    };
    config(deps.storage).save(&state)?;

    // Add initial phylax set to storage
    phylax_set_set(
        deps.storage,
        state.phylax_set_index,
        &msg.initial_phylax_set,
    )?;

    Ok(Response::default())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(deps: DepsMut, env: Env, info: MessageInfo, msg: ExecuteMsg) -> StdResult<Response> {
    match msg {
        ExecuteMsg::PostMessage { message, nonce } => {
            handle_post_message(deps, env, info, message.as_slice(), nonce)
        }
        ExecuteMsg::SubmitVAA { vaa } => handle_submit_vaa(deps, env, info, vaa.as_slice()),
    }
}

/// Process VAA message signed by quardians
fn handle_submit_vaa(
    deps: DepsMut,
    env: Env,
    _info: MessageInfo,
    data: &[u8],
) -> StdResult<Response> {
    let state = config_read(deps.storage).load()?;

    let vaa = parse_and_verify_vaa(deps.storage, data, env.block.time.seconds())?;
    vaa_archive_add(deps.storage, vaa.hash.as_slice())?;

    if state.gov_chain == vaa.emitter_chain && state.gov_address == vaa.emitter_address {
        if state.phylax_set_index != vaa.phylax_set_index {
            return Err(StdError::generic_err(
                "governance VAAs must be signed by the current phylax set",
            ));
        }
        return handle_governance_payload(deps, env, &vaa.payload);
    }

    ContractError::InvalidVAAAction.std_err()
}

fn handle_governance_payload(deps: DepsMut, env: Env, data: &[u8]) -> StdResult<Response> {
    let gov_packet = GovernancePacket::deserialize(data)?;

    let module = String::from_utf8(gov_packet.module).unwrap();
    let module: String = module.chars().filter(|c| c != &'\0').collect();

    if module != "Core" {
        return Err(StdError::generic_err("this is not a valid module"));
    }

    if gov_packet.chain != 0 && gov_packet.chain != CHAIN_ID {
        return Err(StdError::generic_err(
            "the governance VAA is for another chain",
        ));
    }

    match gov_packet.action {
        1u8 => vaa_update_contract(deps, env, &gov_packet.payload),
        2u8 => vaa_update_phylax_set(deps, env, &gov_packet.payload),
        3u8 => handle_set_fee(deps, env, &gov_packet.payload),
        4u8 => handle_transfer_fee(deps, env, &gov_packet.payload),
        _ => ContractError::InvalidVAAAction.std_err(),
    }
}

/// Parses raw VAA data into a struct and verifies whether it contains sufficient signatures of an
/// active phylax set i.e. is valid according to Wormhole consensus rules
fn parse_and_verify_vaa(
    storage: &dyn Storage,
    data: &[u8],
    block_time: u64,
) -> StdResult<ParsedVAA> {
    let vaa = ParsedVAA::deserialize(data)?;

    if vaa.version != 1 {
        return ContractError::InvalidVersion.std_err();
    }

    // Check if VAA with this hash was already accepted
    if vaa_archive_check(storage, vaa.hash.as_slice()) {
        return ContractError::VaaAlreadyExecuted.std_err();
    }

    // Load and check phylax set
    let phylax_set = phylax_set_get(storage, vaa.phylax_set_index);
    let phylax_set: PhylaxSetInfo =
        phylax_set.or_else(|_| ContractError::InvalidPhylaxSetIndex.std_err())?;

    if phylax_set.expiration_time != 0 && phylax_set.expiration_time < block_time {
        return ContractError::PhylaxSetExpired.std_err();
    }
    if (vaa.len_signers as usize) < phylax_set.quorum() {
        return ContractError::NoQuorum.std_err();
    }

    // Verify phylax signatures
    let mut last_index: i32 = -1;
    let mut pos = ParsedVAA::HEADER_LEN;

    for _ in 0..vaa.len_signers {
        if pos + ParsedVAA::SIGNATURE_LEN > data.len() {
            return ContractError::InvalidVAA.std_err();
        }
        let index = data.get_u8(pos) as i32;
        if index <= last_index {
            return ContractError::WrongPhylaxIndexOrder.std_err();
        }
        last_index = index;

        let signature = Signature::try_from(
            &data[pos + ParsedVAA::SIG_DATA_POS
                ..pos + ParsedVAA::SIG_DATA_POS + ParsedVAA::SIG_DATA_LEN],
        )
        .or_else(|_| ContractError::CannotDecodeSignature.std_err())?;
        let id = RecoverableId::new(data.get_u8(pos + ParsedVAA::SIG_RECOVERY_POS))
            .or_else(|_| ContractError::CannotDecodeSignature.std_err())?;
        let recoverable_signature = RecoverableSignature::new(&signature, id)
            .or_else(|_| ContractError::CannotDecodeSignature.std_err())?;

        let verify_key = recoverable_signature
            .recover_verify_key_from_digest_bytes(GenericArray::from_slice(vaa.hash.as_slice()))
            .or_else(|_| ContractError::CannotRecoverKey.std_err())?;

        let index = index as usize;
        if index >= phylax_set.addresses.len() {
            return ContractError::TooManySignatures.std_err();
        }
        if !keys_equal(&verify_key, &phylax_set.addresses[index]) {
            return ContractError::PhylaxSignatureError.std_err();
        }
        pos += ParsedVAA::SIGNATURE_LEN;
    }

    Ok(vaa)
}

fn vaa_update_phylax_set(deps: DepsMut, env: Env, data: &[u8]) -> StdResult<Response> {
    /* Payload format
    0   uint32 new_index
    4   uint8 len(keys)
    5   [][20]uint8 phylax addresses
    */

    let mut state = config_read(deps.storage).load()?;

    let PhylaxSetUpgrade {
        new_phylax_set_index,
        new_phylax_set,
    } = PhylaxSetUpgrade::deserialize(data)?;

    if new_phylax_set_index != state.phylax_set_index + 1 {
        return ContractError::PhylaxSetIndexIncreaseError.std_err();
    }

    let old_phylax_set_index = state.phylax_set_index;

    state.phylax_set_index = new_phylax_set_index;

    phylax_set_set(deps.storage, state.phylax_set_index, &new_phylax_set)?;

    config(deps.storage).save(&state)?;

    let mut old_phylax_set = phylax_set_get(deps.storage, old_phylax_set_index)?;
    old_phylax_set.expiration_time = env.block.time.seconds() + state.phylax_set_expirity;
    phylax_set_set(deps.storage, old_phylax_set_index, &old_phylax_set)?;

    Ok(Response::new()
        .add_attribute("action", "phylax_set_change")
        .add_attribute("old", old_phylax_set_index.to_string())
        .add_attribute("new", state.phylax_set_index.to_string()))
}

fn vaa_update_contract(_deps: DepsMut, env: Env, data: &[u8]) -> StdResult<Response> {
    /* Payload format
    0   [][32]uint8 new_contract
    */

    let ContractUpgrade { new_contract } = ContractUpgrade::deserialize(data)?;

    Ok(Response::new()
        .add_message(CosmosMsg::Wasm(WasmMsg::Migrate {
            contract_addr: env.contract.address.to_string(),
            new_code_id: new_contract,
            msg: to_binary(&MigrateMsg {})?,
        }))
        .add_attribute("action", "contract_upgrade"))
}

pub fn handle_set_fee(deps: DepsMut, _env: Env, data: &[u8]) -> StdResult<Response> {
    let set_fee_msg = SetFee::deserialize(data)?;

    // Save new fees
    let mut state = config_read(deps.storage).load()?;
    state.fee = set_fee_msg.fee;
    config(deps.storage).save(&state)?;

    Ok(Response::new()
        .add_attribute("action", "fee_change")
        .add_attribute("new_fee.amount", state.fee.amount)
        .add_attribute("new_fee.denom", state.fee.denom))
}

pub fn handle_transfer_fee(deps: DepsMut, _env: Env, data: &[u8]) -> StdResult<Response> {
    let transfer_msg = TransferFee::deserialize(data)?;

    Ok(Response::new().add_message(CosmosMsg::Bank(BankMsg::Send {
        to_address: deps.api.addr_humanize(&transfer_msg.recipient)?.to_string(),
        amount: vec![transfer_msg.amount],
    })))
}

fn handle_post_message(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    message: &[u8],
    nonce: u32,
) -> StdResult<Response> {
    let state = config_read(deps.storage).load()?;
    let fee = state.fee;

    // Check fee
    if fee.amount.u128() > 0 && !has_coins(info.funds.as_ref(), &fee) {
        return ContractError::FeeTooLow.std_err();
    }

    let emitter = extend_address_to_32(&deps.api.addr_canonicalize(info.sender.as_str())?);
    let sequence = sequence_read(deps.storage, emitter.as_slice());
    sequence_set(deps.storage, emitter.as_slice(), sequence + 1)?;

    Ok(Response::new()
        .add_attribute("message.message", hex::encode(message))
        .add_attribute("message.sender", hex::encode(emitter))
        .add_attribute("message.chain_id", CHAIN_ID.to_string())
        .add_attribute("message.nonce", nonce.to_string())
        .add_attribute("message.sequence", sequence.to_string())
        .add_attribute("message.block_time", env.block.time.seconds().to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::PhylaxSetInfo {} => to_binary(&query_phylax_set_info(deps)?),
        QueryMsg::VerifyVAA { vaa, block_time } => to_binary(&query_parse_and_verify_vaa(
            deps,
            vaa.as_slice(),
            block_time,
        )?),
        QueryMsg::GetState {} => to_binary(&query_state(deps)?),
        QueryMsg::QueryAddressHex { address } => to_binary(&query_address_hex(deps, &address)?),
    }
}

pub fn query_phylax_set_info(deps: Deps) -> StdResult<PhylaxSetInfoResponse> {
    let state = config_read(deps.storage).load()?;
    let phylax_set = phylax_set_get(deps.storage, state.phylax_set_index)?;
    let res = PhylaxSetInfoResponse {
        phylax_set_index: state.phylax_set_index,
        addresses: phylax_set.addresses,
    };
    Ok(res)
}

pub fn query_parse_and_verify_vaa(
    deps: Deps,
    data: &[u8],
    block_time: u64,
) -> StdResult<ParsedVAA> {
    parse_and_verify_vaa(deps.storage, data, block_time)
}

// returns the hex of the 32 byte address we use for some address on this chain
pub fn query_address_hex(deps: Deps, address: &HumanAddr) -> StdResult<GetAddressHexResponse> {
    Ok(GetAddressHexResponse {
        hex: hex::encode(extend_address_to_32(&deps.api.addr_canonicalize(address)?)),
    })
}

pub fn query_state(deps: Deps) -> StdResult<GetStateResponse> {
    let state = config_read(deps.storage).load()?;
    let res = GetStateResponse { fee: state.fee };
    Ok(res)
}

fn keys_equal(a: &VerifyingKey, b: &PhylaxAddress) -> bool {
    let mut hasher = Keccak256::new();

    let point = if let Some(p) = EncodedPoint::from(a).decompress() {
        p
    } else {
        return false;
    };

    hasher.update(&point.as_bytes()[1..]);
    let a = &hasher.finalize()[12..];

    let b = &b.bytes;
    if a.len() != b.len() {
        return false;
    }
    for (ai, bi) in a.iter().zip(b.as_slice().iter()) {
        if ai != bi {
            return false;
        }
    }
    true
}
