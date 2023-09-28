package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/deltaswapio/deltachain/x/wormhole/types"
	"github.com/spf13/cobra"
)

func CmdListWasmInstantiateAllowlist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-wasm-instantiate-allowlist",
		Short: "list all allowed contract address and code IDs pairs for wasm instantiate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllWasmInstantiateAllowlist{
				Pagination: pageReq,
			}

			res, err := queryClient.WasmInstantiateAllowlistAll(context.Background(), params)
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
