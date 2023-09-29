use crate::{
    accounts::{
        Bridge,
        BridgeConfig,
        FeeCollector,
        PhylaxSet,
        PhylaxSetDerivationData,
    },
    error::Error::TooManyPhylaxs,
    MAX_LEN_GUARDIAN_KEYS,
};
use solana_program::sysvar::clock::Clock;
use solitaire::{
    CreationLamports::Exempt,
    *,
};

type Payer<'a> = Signer<Info<'a>>;

#[derive(FromAccounts)]
pub struct Initialize<'b> {
    /// Bridge config.
    pub bridge: Mut<Bridge<'b, { AccountState::Uninitialized }>>,

    /// Location the new phylax set will be allocated at.
    pub phylax_set: Mut<PhylaxSet<'b, { AccountState::Uninitialized }>>,

    /// Location of the fee collector that users will need to pay.
    pub fee_collector: Mut<FeeCollector<'b>>,

    /// Payer for account creation.
    pub payer: Mut<Payer<'b>>,

    /// Clock used for recording the initialization time.
    pub clock: Sysvar<'b, Clock>,
}

#[derive(BorshDeserialize, BorshSerialize, Default)]
pub struct InitializeData {
    /// Period for how long a phylax set is valid after it has been replaced by a new one.  This
    /// guarantees that VAAs issued by that set can still be submitted for a certain period.  In
    /// this period we still trust the old phylax set.
    pub phylax_set_expiration_time: u32,

    /// Amount of lamports that needs to be paid to the protocol to post a message
    pub fee: u64,

    /// Initial Phylax Set
    pub initial_phylaxs: Vec<[u8; 20]>,
}

pub fn initialize(
    ctx: &ExecutionContext,
    accs: &mut Initialize,
    data: InitializeData,
) -> Result<()> {
    let index = 0;

    if data.initial_phylaxs.len() > MAX_LEN_GUARDIAN_KEYS {
        return Err(TooManyPhylaxs.into());
    }

    // Allocate initial phylax set with the provided keys.
    accs.phylax_set.index = index;
    accs.phylax_set.creation_time = accs.clock.unix_timestamp as u32;
    accs.phylax_set.keys.extend(&data.initial_phylaxs);

    // Initialize Phylax Set
    accs.phylax_set.create(
        &PhylaxSetDerivationData { index },
        ctx,
        accs.payer.key,
        Exempt,
    )?;

    // Initialize the Bridge state for the first time.
    accs.bridge.create(ctx, accs.payer.key, Exempt)?;
    accs.bridge.phylax_set_index = index;
    accs.bridge.config = BridgeConfig {
        phylax_set_expiration_time: data.phylax_set_expiration_time,
        fee: data.fee,
    };

    // Initialize the fee collector account so it's rent exempt and will keep funds
    accs.fee_collector.create(
        ctx,
        accs.payer.key,
        Exempt,
        0,
        &solana_program::system_program::id(),
    )?;
    accs.bridge.last_lamports = accs.fee_collector.lamports();

    Ok(())
}
