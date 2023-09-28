package rest

import (
	"encoding/hex"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

type (
	// PhylaxSetUpdateProposalReq defines a guardian set update proposal request body.
	PhylaxSetUpdateProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title          string         `json:"title" yaml:"title"`
		Description    string         `json:"description" yaml:"description"`
		PhylaxSetIndex uint32         `json:"guardianSetIndex" yaml:"guardianSetIndex"`
		PhylaxSetKeys  []string       `json:"guardianSetKeys" yaml:"guardianSetKeys"`
		Proposer       sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit        sdk.Coins      `json:"deposit" yaml:"deposit"`
	}

	// DeltaswapGovernanceMessageProposalReq defines a deltaswap governance message proposal request body.
	DeltaswapGovernanceMessageProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		TargetChain uint16         `json:"targetChain" yaml:"targetChain"`
		Action      uint8          `json:"action" yaml:"action"`
		Module      []byte         `json:"module" yaml:"module"`
		Payload     []byte         `json:"payload" yaml:"payload"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)

// ProposalPhylaxSetUpdateRESTHandler returns a ProposalRESTHandler that exposes the guardian set update
// REST handler with a given sub-route.
func ProposalPhylaxSetUpdateRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wormhole_phylax_update",
		Handler:  postProposalPhylaxSetUpdateHandlerFn(clientCtx),
	}
}

func postProposalPhylaxSetUpdateHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PhylaxSetUpdateProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		keys := make([][]byte, len(req.PhylaxSetKeys))
		for i, keyString := range req.PhylaxSetKeys {
			keyBytes, err := hex.DecodeString(keyString)
			if rest.CheckBadRequestError(w, err) {
				return
			}
			keys[i] = keyBytes
		}

		content := types.NewPhylaxSetUpdateProposal(req.Title, req.Description, types.PhylaxSet{
			Index:          req.PhylaxSetIndex,
			Keys:           keys,
			ExpirationTime: 0,
		})

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// ProposalDeltaswapGovernanceMessageRESTHandler returns a ProposalRESTHandler that exposes the deltaswap governance message
// REST handler with a given sub-route.
func ProposalDeltaswapGovernanceMessageRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "wormhole_governance_message",
		Handler:  postProposalDeltaswapGovernanceMessageHandlerFn(clientCtx),
	}
}

func postProposalDeltaswapGovernanceMessageHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DeltaswapGovernanceMessageProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewGovernanceDeltaswapMessageProposal(req.Title, req.Description, req.Action, req.TargetChain, req.Module, req.Payload)

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
