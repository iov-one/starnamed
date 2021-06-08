package cli

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/iov-one/starnamed/x/escrow/types"
	"github.com/spf13/cobra"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// GetQueryCmd returns the cli query commands for the module.
func GetQueryCmd() *cobra.Command {
	escrowQueryCmd := &cobra.Command{
		Use:                        "escrow [id]",
		Short:                      "Query an escrow",
		Long:                       "Query details of an escrow with the specified id.",
		Example:                    fmt.Sprintf("%s query escrow <id>", version.AppName),
		Args:                       cobra.ExactArgs(1),
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			param := types.QueryEscrowRequest{Id: tmbytes.HexBytes(id).String()}
			response, err := queryClient.Escrow(context.Background(), &param)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response.Escrow)
		},
	}

	flags.AddQueryFlagsToCmd(escrowQueryCmd)

	return escrowQueryCmd
}
