package cli

import (
	"context"

	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/spf13/cobra"
)

func CmdListPhylaxValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-guardian-validator",
		Short: "list all guardian-validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPhylaxValidatorRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PhylaxValidatorAll(context.Background(), params)
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

func CmdShowPhylaxValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-guardian-validator [guardian-key]",
		Short: "shows a guardian-validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPhylaxKey, err := hex.DecodeString(args[0])

			if err != nil {
				return err
			}

			params := &types.QueryGetPhylaxValidatorRequest{
				PhylaxKey: argPhylaxKey,
			}

			res, err := queryClient.PhylaxValidator(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
