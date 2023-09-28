package cli

import (
	"fmt"
	"strconv"

	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/deltaswapio/deltachain/x/deltaswap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRegisterAccountAsPhylax() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-account-as-guardian [signature]",
		Short: "Register a guardian public key with a deltaswap chain address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argSignature, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("malformed signature: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterAccountAsPhylax(
				clientCtx.GetFromAddress().String(),
				argSignature,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
