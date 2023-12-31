package cli

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/spf13/cobra"
)

const FlagPhylaxSetKeys = "phylax-set-keys"
const FlagPhylaxSetIndex = "phylax-set-index"

// NewCmdSubmitPhylaxSetUpdateProposal implements a command handler for submitting a phylax set update governance
// proposal.
func NewCmdSubmitPhylaxSetUpdateProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-phylax-set [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Submit a phylax set update proposal",
		Long:  "Submit a proposal to update the current phylax set to a new one",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			keyStrings, err := cmd.Flags().GetStringArray(FlagPhylaxSetKeys)
			if err != nil {
				return err
			}

			newIndex, err := cmd.Flags().GetUint32(FlagPhylaxSetIndex)
			if err != nil {
				return err
			}

			keys := make([][]byte, len(keyStrings))
			for i, keyString := range keyStrings {
				keyBytes, err := hex.DecodeString(keyString)
				if err != nil {
					return err
				}
				keys[i] = keyBytes
			}

			content := types.NewPhylaxSetUpdateProposal(title, description, types.PhylaxSet{
				Index:          newIndex,
				Keys:           keys,
				ExpirationTime: 0,
			})
			err = content.ValidateBasic()
			if err != nil {
				return err
			}

			msg, err := gov.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	cmd.Flags().StringArray(FlagPhylaxSetKeys, []string{}, "list of phylax keys (hex encoded without 0x)")
	cmd.Flags().Uint32(FlagPhylaxSetIndex, 0, "index of the new phylax set")
	cmd.MarkFlagRequired(cli.FlagTitle)
	cmd.MarkFlagRequired(cli.FlagDescription)
	cmd.MarkFlagRequired(FlagPhylaxSetKeys)
	cmd.MarkFlagRequired(FlagPhylaxSetIndex)

	return cmd
}

const FlagAction = "action"
const FlagModule = "module"
const FlagTargetChainID = "target-chain-id"
const FlagPayload = "payload"

// NewCmdSubmitDeltaswapGovernanceMessageProposal implements a command handler for submitting a generic Deltaswap
// governance message.
func NewCmdSubmitDeltaswapGovernanceMessageProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wormhole-governance-message [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Submit a deltaswap governance message proposal",
		Long:  "Submit a proposal to emit a generic deltaswap governance message",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			action, err := cmd.Flags().GetUint8(FlagAction)
			if err != nil {
				return err
			}

			targetChain, err := cmd.Flags().GetUint16(FlagTargetChainID)
			if err != nil {
				return err
			}

			module, err := cmd.Flags().GetBytesHex(FlagModule)
			if err != nil {
				return err
			}

			payload, err := cmd.Flags().GetBytesHex(FlagPayload)
			if err != nil {
				return err
			}

			content := types.NewGovernanceDeltaswapMessageProposal(title, description, action, targetChain, module, payload)
			err = content.ValidateBasic()
			if err != nil {
				return err
			}

			msg, err := gov.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	cmd.Flags().Uint8(FlagAction, 0, "target chain of the message (0 for all)")
	cmd.Flags().Uint16(FlagTargetChainID, 0, "target chain of the message (0 for all)")
	cmd.Flags().BytesHex(FlagModule, []byte{}, "module identifier of the message")
	cmd.Flags().BytesHex(FlagPayload, []byte{}, "payload of the message")
	cmd.MarkFlagRequired(cli.FlagTitle)
	cmd.MarkFlagRequired(cli.FlagDescription)
	cmd.MarkFlagRequired(FlagAction)
	cmd.MarkFlagRequired(FlagTargetChainID)
	cmd.MarkFlagRequired(FlagModule)
	cmd.MarkFlagRequired(FlagPayload)

	return cmd
}
