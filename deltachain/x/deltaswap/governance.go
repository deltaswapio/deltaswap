package deltaswap

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

// NewDeltaswapGovernanceProposalHandler creates a governance handler to manage new proposal types.
// It enables PhylaxSetProposal to update the phylax set and GenericDeltaswapMessageProposal to emit a generic wormhole
// message from the governance emitter.
func NewDeltaswapGovernanceProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.PhylaxSetUpdateProposal:
			return handlePhylaxSetUpdateProposal(ctx, k, c)

		case *types.GovernanceDeltaswapMessageProposal:
			return handleGovernanceDeltaswapMessageProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized deltaswap proposal content type: %T", c)
		}
	}
}

func handlePhylaxSetUpdateProposal(ctx sdk.Context, k keeper.Keeper, proposal *types.PhylaxSetUpdateProposal) error {
	err := k.UpdatePhylaxSet(ctx, types.PhylaxSet{
		Index:          proposal.NewPhylaxSet.Index,
		Keys:           proposal.NewPhylaxSet.Keys,
		ExpirationTime: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to update phylax set: %w", err)
	}

	config, ok := k.GetConfig(ctx)
	if !ok {
		return types.ErrNoConfig
	}

	// Post a deltaswap phylax set update governance message
	message := &bytes.Buffer{}

	// Header
	message.Write(vaa.CoreModule)
	MustWrite(message, binary.BigEndian, uint8(2))
	MustWrite(message, binary.BigEndian, uint16(0))

	// Body
	MustWrite(message, binary.BigEndian, proposal.NewPhylaxSet.Index)
	MustWrite(message, binary.BigEndian, uint8(len(proposal.NewPhylaxSet.Keys)))
	for _, key := range proposal.NewPhylaxSet.Keys {
		message.Write(key)
	}

	emitterAddress, err := types.EmitterAddressFromBytes32(config.GovernanceEmitter)
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	err = k.PostMessage(ctx, emitterAddress, 0, message.Bytes())
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}

func handleGovernanceDeltaswapMessageProposal(ctx sdk.Context, k keeper.Keeper, proposal *types.GovernanceDeltaswapMessageProposal) error {
	config, ok := k.GetConfig(ctx)
	if !ok {
		return types.ErrNoConfig
	}

	// Post a deltaswap governance message
	message := &bytes.Buffer{}
	message.Write(proposal.Module)
	MustWrite(message, binary.BigEndian, uint8(proposal.Action))
	MustWrite(message, binary.BigEndian, uint16(proposal.TargetChain))
	message.Write(proposal.Payload)

	emitterAddress, err := types.EmitterAddressFromBytes32(config.GovernanceEmitter)
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	err = k.PostMessage(ctx, emitterAddress, 0, message.Bytes())
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}

// MustWrite calls binary.Write and panics on errors
func MustWrite(w io.Writer, order binary.ByteOrder, data interface{}) {
	if err := binary.Write(w, order, data); err != nil {
		panic(fmt.Errorf("failed to write binary data: %v", data).Error())
	}
}
