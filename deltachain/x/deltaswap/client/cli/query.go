package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deltaswapio/deltachain/x/deltaswap/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group deltaswap queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdListPhylaxSet())
	cmd.AddCommand(CmdShowPhylaxSet())
	cmd.AddCommand(CmdShowConfig())
	cmd.AddCommand(CmdListReplayProtection())
	cmd.AddCommand(CmdShowReplayProtection())
	cmd.AddCommand(CmdListSequenceCounter())
	cmd.AddCommand(CmdShowSequenceCounter())
	cmd.AddCommand(CmdShowConsensusPhylaxSetIndex())
	cmd.AddCommand(CmdListPhylaxValidator())
	cmd.AddCommand(CmdShowPhylaxValidator())
	cmd.AddCommand(CmdLatestPhylaxSetIndex())
	cmd.AddCommand(CmdListAllowlists())
	cmd.AddCommand(CmdShowAllowlist())
	cmd.AddCommand(CmdShowIbcComposabilityMwContract())
	cmd.AddCommand(CmdListWasmInstantiateAllowlist())

	// this line is used by starport scaffolding # 1

	return cmd
}
