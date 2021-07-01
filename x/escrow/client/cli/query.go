package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/iov-one/starnamed/x/escrow/types"
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

			id := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			param := types.QueryEscrowRequest{Id: id}
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
