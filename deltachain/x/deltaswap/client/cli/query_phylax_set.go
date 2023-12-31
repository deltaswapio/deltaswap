package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/spf13/cobra"
)

func CmdListPhylaxSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-phylax-set",
		Short: "list all PhylaxSet",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPhylaxSetRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PhylaxSetAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowPhylaxSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-phylax-set [id]",
		Short: "shows a PhylaxSet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			params := &types.QueryGetPhylaxSetRequest{
				Index: uint32(id),
			}

			res, err := queryClient.PhylaxSet(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
