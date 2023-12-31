package bindings

import (
	"github.com/CosmWasm/wasmd/x/wasm"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	tokenfactorykeeper "github.com/deltaswapio/deltachain/x/tokenfactory/keeper"
)

func RegisterCustomPlugins(
	bank *bankkeeper.BaseKeeper,
	tokenFactory *tokenfactorykeeper.Keeper,
) []wasmkeeper.Option {
	// Disabling tokenfactory custom querier because deltachain custom querier exists
	//wasmQueryPlugin := NewQueryPlugin(bank, tokenFactory)

	//queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
	//	Custom: CustomQuerier(wasmQueryPlugin),
	//})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, tokenFactory),
	)

	return []wasm.Option{
		//	queryPluginOpt,
		messengerDecoratorOpt,
	}
}
