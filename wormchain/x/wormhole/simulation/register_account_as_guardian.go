package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/deltaswapio/deltachain/x/wormhole/keeper"
	"github.com/deltaswapio/deltachain/x/wormhole/types"
)

func SimulateMsgRegisterAccountAsPhylax(
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRegisterAccountAsPhylax{
			Signer: simAccount.Address.String(),
		}

		// TODO: Handling the RegisterAccountAsPhylax simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "RegisterAccountAsPhylax simulation not implemented"), nil, nil
	}
}
