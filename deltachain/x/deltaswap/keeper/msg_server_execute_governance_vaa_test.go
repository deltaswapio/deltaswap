package keeper_test

import (
	"crypto/ecdsa"
	"encoding/binary"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/deltaswapio/deltachain/testutil/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/keeper"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
	"github.com/stretchr/testify/assert"
)

func createExecuteGovernanceVaaPayload(k *keeper.Keeper, ctx sdk.Context, num_phylaxs byte) ([]byte, []*ecdsa.PrivateKey) {
	phylaxs, privateKeys := createNPhylaxValidator(k, ctx, int(num_phylaxs))
	next_index := k.GetPhylaxSetCount(ctx)
	set_update := make([]byte, 4)
	binary.BigEndian.PutUint32(set_update, next_index)
	set_update = append(set_update, num_phylaxs)
	// Add keys to set_update
	for _, phylax := range phylaxs {
		set_update = append(set_update, phylax.PhylaxKey...)
	}
	// governance message with sha3 of wasmBytes as the payload
	module := [32]byte{}
	copy(module[:], vaa.CoreModule)
	gov_msg := types.NewGovernanceMessage(module, byte(vaa.ActionPhylaxSetUpdate), uint16(vaa.ChainIDDeltachain), set_update)

	return gov_msg.MarshalBinary(), privateKeys
}

func TestExecuteGovernanceVAA(t *testing.T) {
	k, ctx := keepertest.DeltaswapKeeper(t)
	phylaxs, privateKeys := createNPhylaxValidator(k, ctx, 10)
	_ = privateKeys
	k.SetConfig(ctx, types.Config{
		GovernanceEmitter:   vaa.GovernanceEmitter[:],
		GovernanceChain:     uint32(vaa.GovernanceChain),
		ChainId:             uint32(vaa.ChainIDDeltachain),
		PhylaxSetExpiration: 86400,
	})
	signer_bz := [20]byte{}
	signer := sdk.AccAddress(signer_bz[:])

	set := createNewPhylaxSet(k, ctx, phylaxs)
	k.SetConsensusPhylaxSetIndex(ctx, types.ConsensusPhylaxSetIndex{Index: set.Index})

	context := sdk.WrapSDKContext(ctx)
	msgServer := keeper.NewMsgServerImpl(*k)

	// create governance to update phylax set with extra phylax
	payload, newPrivateKeys := createExecuteGovernanceVaaPayload(k, ctx, 11)
	v := generateVaa(set.Index, privateKeys, vaa.ChainID(vaa.GovernanceChain), payload)
	vBz, _ := v.Marshal()
	_, err := msgServer.ExecuteGovernanceVAA(context, &types.MsgExecuteGovernanceVAA{
		Signer: signer.String(),
		Vaa:    vBz,
	})
	assert.NoError(t, err)

	// we should have a new set with 11 phylaxs now
	new_index := k.GetLatestPhylaxSetIndex(ctx)
	assert.Equal(t, set.Index+1, new_index)
	new_set, _ := k.GetPhylaxSet(ctx, new_index)
	assert.Len(t, new_set.Keys, 11)

	// Submitting another change with the old set doesn't work
	v = generateVaa(set.Index, privateKeys, vaa.ChainID(vaa.GovernanceChain), payload)
	vBz, _ = v.Marshal()
	_, err = msgServer.ExecuteGovernanceVAA(context, &types.MsgExecuteGovernanceVAA{
		Signer: signer.String(),
		Vaa:    vBz,
	})
	assert.ErrorIs(t, err, types.ErrPhylaxSetNotSequential)

	// Invalid length
	v = generateVaa(set.Index, privateKeys, vaa.ChainID(vaa.GovernanceChain), payload[:len(payload)-1])
	vBz, _ = v.Marshal()
	_, err = msgServer.ExecuteGovernanceVAA(context, &types.MsgExecuteGovernanceVAA{
		Signer: signer.String(),
		Vaa:    vBz,
	})
	assert.ErrorIs(t, err, types.ErrInvalidGovernancePayloadLength)

	// Include a phylax address twice in an update
	payload_bad, _ := createExecuteGovernanceVaaPayload(k, ctx, 11)
	copy(payload_bad[len(payload_bad)-20:], payload_bad[len(payload_bad)-40:len(payload_bad)-20])
	v = generateVaa(set.Index, privateKeys, vaa.ChainID(vaa.GovernanceChain), payload_bad)
	vBz, _ = v.Marshal()
	_, err = msgServer.ExecuteGovernanceVAA(context, &types.MsgExecuteGovernanceVAA{
		Signer: signer.String(),
		Vaa:    vBz,
	})
	assert.ErrorIs(t, err, types.ErrDuplicatePhylaxAddress)

	// Change set again with new set update
	payload, _ = createExecuteGovernanceVaaPayload(k, ctx, 12)
	v = generateVaa(new_set.Index, newPrivateKeys, vaa.ChainID(vaa.GovernanceChain), payload)
	vBz, _ = v.Marshal()
	_, err = msgServer.ExecuteGovernanceVAA(context, &types.MsgExecuteGovernanceVAA{
		Signer: signer.String(),
		Vaa:    vBz,
	})
	assert.NoError(t, err)
	new_index2 := k.GetLatestPhylaxSetIndex(ctx)
	assert.Equal(t, new_set.Index+1, new_index2)
}
