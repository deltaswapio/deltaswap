package keeper

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/deltaswapio/deltachain/app"
	"github.com/deltaswapio/deltachain/app/wasm_handlers"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"

	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/spm/cosmoscmd"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func DeltaswapKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	k, _, _, ctx := DeltaswapKeeperAndWasmd(t)
	return k, ctx
}

func DeltaswapKeeperAndWasmd(t testing.TB) (*keeper.Keeper, wasmkeeper.Keeper, *wasmkeeper.PermissionedKeeper, sdk.Context) {
	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		paramstypes.StoreKey,
		capabilitytypes.StoreKey,
		types.StoreKey,
		wasmtypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, types.MemStoreKey)
	maccPerms := map[string][]string{}

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(keys[authtypes.StoreKey], sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys[paramstypes.StoreKey], sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys[capabilitytypes.StoreKey], sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys[types.StoreKey], sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keys[wasmtypes.StoreKey], sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memKeys[types.MemStoreKey], sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(tkeys[paramstypes.TStoreKey], sdk.StoreTypeTransient, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	encodingConfig := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)
	appCodec := encodingConfig.Marshaler
	amino := encodingConfig.Amino

	paramsKeeper := paramskeeper.NewKeeper(appCodec, amino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	paramsKeeper.Subspace(types.ModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)

	paramsKeeper.Subspace(authtypes.ModuleName)
	subspace_auth, _ := paramsKeeper.GetSubspace(authtypes.ModuleName)
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], subspace_auth, authtypes.ProtoBaseAccount, maccPerms,
	)
	// this line is used by starport scaffolding # stargate/app/paramSubspace

	subspaceWasmd, _ := paramsKeeper.GetSubspace(wasmtypes.ModuleName)

	bApp := baseapp.NewBaseApp("deltachain", log.NewNopLogger(), db, encodingConfig.TxConfig.TxDecoder())
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)

	appapp := &app.App{
		BaseApp: bApp,
	}

	k := keeper.NewKeeper(
		appCodec,
		keys[types.StoreKey],
		memKeys[types.MemStoreKey],
		accountKeeper,
		nil,
	)

	supportedFeatures := "iterator,staking,stargate,deltaswap"
	appapp.DeltaswapKeeper = *k

	appapp.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedWasmKeeper := appapp.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	wasmDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	wasmKeeper := wasm.NewKeeper(
		appCodec,
		keys[wasmtypes.StoreKey],
		subspaceWasmd,
		accountKeeper,
		&wasm_handlers.BankKeeperHandler{},
		&wasm_handlers.StakingKeeperHandler{},
		&wasm_handlers.DistributionKeeperHandler{},
		&wasm_handlers.ChannelKeeperHandler{},
		&wasm_handlers.PortKeeperHandler{},
		scopedWasmKeeper,
		&wasm_handlers.ICS20TransferPortSourceHandler{},
		appapp.DeltaswapKeeper,
		appapp.MsgServiceRouter(),
		appapp.GRPCQueryRouter(),
		wasmDir,
		wasm.DefaultWasmConfig(),
		supportedFeatures,
		wasmkeeper.WithQueryPlugins(keeper.NewCustomQueryHandler(appapp.DeltaswapKeeper)),
	)
	ctx := sdk.NewContext(stateStore, tmproto.Header{
		Time: time.Now(),
		// The height should be at least 1, because the allowlist antehandler
		// passes everything at height 0 for gen tx's.
		Height: 1,
	}, false, log.NewNopLogger())
	appapp.MountKVStores(keys)
	appapp.MountTransientStores(tkeys)
	appapp.MountMemoryStores(memKeys)

	wasmGenState := wasmtypes.GenesisState{}
	wasmGenState.Params.CodeUploadAccess = wasmtypes.DefaultUploadAccess
	wasmGenState.Params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
	wasmKeeper.SetParams(ctx, wasmGenState.Params)
	permissionedWasmKeeper := wasmkeeper.NewDefaultPermissionKeeper(wasmKeeper)
	appapp.DeltaswapKeeper.SetWasmdKeeper(permissionedWasmKeeper)
	k.SetWasmdKeeper(permissionedWasmKeeper)

	return k, wasmKeeper, permissionedWasmKeeper, ctx
}
