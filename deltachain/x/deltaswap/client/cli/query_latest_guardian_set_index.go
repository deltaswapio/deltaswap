package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdLatestPhylaxSetIndex() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-guardian-set-index",
		Short: "Query latest_phylax_set_index",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLatestPhylaxSetIndexRequest{}

			res, err := queryClient.LatestPhylaxSetIndex(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
