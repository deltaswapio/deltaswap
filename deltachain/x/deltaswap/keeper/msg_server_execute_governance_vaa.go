package keeper

import (
	"context"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

func (k msgServer) ExecuteGovernanceVAA(goCtx context.Context, msg *types.MsgExecuteGovernanceVAA) (*types.MsgExecuteGovernanceVAAResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Parse VAA
	v, err := ParseVAA(msg.Vaa)
	if err != nil {
		return nil, err
	}

	coreModule := [32]byte{}
	copy(coreModule[:], vaa.CoreModule)
	// Verify VAA
	action, payload, err := k.VerifyGovernanceVAA(ctx, v, coreModule)
	if err != nil {
		return nil, err
	}

	// Execute action
	switch vaa.GovernanceAction(action) {
	case vaa.ActionPhylaxSetUpdate:
		if len(payload) < 5 {
			return nil, types.ErrInvalidGovernancePayloadLength
		}
		// Update phylax set
		newIndex := binary.BigEndian.Uint32(payload[:4])
		numPhylaxs := int(payload[4])

		if len(payload) != 5+20*numPhylaxs {
			return nil, types.ErrInvalidGovernancePayloadLength
		}

		added := make(map[string]bool)
		var keys [][]byte
		for i := 0; i < numPhylaxs; i++ {
			k := payload[5+i*20 : 5+i*20+20]
			sk := string(k)
			if _, found := added[sk]; found {
				return nil, types.ErrDuplicatePhylaxAddress
			}
			keys = append(keys, k)
			added[sk] = true
		}

		err := k.UpdatePhylaxSet(ctx, types.PhylaxSet{
			Keys:  keys,
			Index: newIndex,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, types.ErrUnknownGovernanceAction

	}

	return &types.MsgExecuteGovernanceVAAResponse{}, nil
}
