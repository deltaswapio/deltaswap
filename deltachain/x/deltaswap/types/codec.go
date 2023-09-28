package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgExecuteGovernanceVAA{}, "deltaswap/ExecuteGovernanceVAA", nil)
	cdc.RegisterConcrete(&MsgRegisterAccountAsPhylax{}, "deltaswap/RegisterAccountAsPhylax", nil)
	cdc.RegisterConcrete(&MsgStoreCode{}, "deltaswap/StoreCode", nil)
	cdc.RegisterConcrete(&MsgInstantiateContract{}, "deltaswap/InstantiateContract", nil)
	cdc.RegisterConcrete(&MsgMigrateContract{}, "deltaswap/MigrateContract", nil)
	cdc.RegisterConcrete(&MsgCreateAllowlistEntryRequest{}, "deltaswap/CreateAllowlistEntryRequest", nil)
	cdc.RegisterConcrete(&MsgDeleteAllowlistEntryRequest{}, "deltaswap/DeleteAllowlistEntryRequest", nil)
	cdc.RegisterConcrete(&MsgAddWasmInstantiateAllowlist{}, "deltaswap/AddWasmInstantiateAllowlist", nil)
	cdc.RegisterConcrete(&MsgDeleteWasmInstantiateAllowlist{}, "deltaswap/DeleteWasmInstantiateAllowlist", nil)
	cdc.RegisterConcrete(&MsgExecuteGatewayGovernanceVaa{}, "deltaswap/ExecuteGatewayGovernanceVaa", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgExecuteGovernanceVAA{},
		&MsgStoreCode{},
		&MsgInstantiateContract{},
		&MsgMigrateContract{},
		&MsgCreateAllowlistEntryRequest{},
		&MsgDeleteAllowlistEntryRequest{},
		&MsgExecuteGatewayGovernanceVaa{},
	)
	registry.RegisterImplementations((*gov.Content)(nil),
		&GovernanceDeltaswapMessageProposal{},
		&PhylaxSetUpdateProposal{})
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterAccountAsPhylax{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
