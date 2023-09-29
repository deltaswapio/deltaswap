package deltaswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the phylaxSet
	for _, elem := range genState.PhylaxSetList {
		if _, err := k.AppendPhylaxSet(ctx, elem); err != nil {
			panic(err)
		}
	}

	// Set if defined
	if genState.Config != nil {
		k.SetConfig(ctx, *genState.Config)
	}
	// Set all the replayProtection
	for _, elem := range genState.ReplayProtectionList {
		k.SetReplayProtection(ctx, elem)
	}
	// Set all the sequenceCounter
	for _, elem := range genState.SequenceCounterList {
		k.SetSequenceCounter(ctx, elem)
	}
	// Set if defined
	if genState.ConsensusPhylaxSetIndex != nil {
		k.SetConsensusPhylaxSetIndex(ctx, *genState.ConsensusPhylaxSetIndex)
	}
	// Set all the phylaxValidator
	for _, elem := range genState.PhylaxValidatorList {
		k.SetPhylaxValidator(ctx, elem)
	}
	for _, elem := range genState.AllowedAddresses {
		k.SetValidatorAllowedAddress(ctx, elem)
	}
	// Set all the contract/code_id pairs for the wasm instantiate allowlist
	for _, elem := range genState.WasmInstantiateAllowlist {
		k.SetWasmInstantiateAllowlist(ctx, elem)
	}
	k.StoreIbcComposabilityMwContract(ctx, genState.IbcComposabilityMwContract)
	// this line is used by starport scaffolding # genesis/module/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.PhylaxSetList = k.GetAllPhylaxSet(ctx)

	// Get all config
	config, found := k.GetConfig(ctx)
	if found {
		genesis.Config = &config
	}
	genesis.ReplayProtectionList = k.GetAllReplayProtection(ctx)
	genesis.SequenceCounterList = k.GetAllSequenceCounter(ctx)
	// Get all consensusPhylaxSetIndex
	consensusPhylaxSetIndex, found := k.GetConsensusPhylaxSetIndex(ctx)
	if found {
		genesis.ConsensusPhylaxSetIndex = &consensusPhylaxSetIndex
	}
	genesis.PhylaxValidatorList = k.GetAllPhylaxValidator(ctx)
	genesis.AllowedAddresses = k.GetAllAllowedAddresses(ctx)
	genesis.WasmInstantiateAllowlist = k.GetAllWasmInstiateAllowedAddresses(ctx)
	genesis.IbcComposabilityMwContract = k.GetIbcComposabilityMwContract(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
