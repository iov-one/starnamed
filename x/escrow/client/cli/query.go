package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/iov-one/starnamed/x/escrow/types"
)

// GetQueryCmd returns the cli query commands for the module.
func GetQueryCmd() *cobra.Command {
	escrowQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Escrow query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	escrowQueryCmd.AddCommand(
		getCmdQueryEscrow(),
		getCmdQueryEscrows(),
	)

	return escrowQueryCmd
}

func getCmdQueryEscrow() *cobra.Command {
	escrowQueryCmd := &cobra.Command{
		Use:                        "escrow [id]",
		Short:                      "Query an escrow",
		Long:                       "Query details of an escrow with the specified id.",
		Example:                    fmt.Sprintf("%s query escrow escrow <id>", version.AppName),
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

			return clientCtx.PrintProto(response)
		},
	}

	flags.AddQueryFlagsToCmd(escrowQueryCmd)

	return escrowQueryCmd
}

func getCmdQueryEscrows() *cobra.Command {
	escrowQueryCmd := &cobra.Command{
		Use:                        "escrows",
		Short:                      "Do a query over all the escrows",
		Long:                       "Query details of a list of escrows, with the possibility to filter by seller, object and/or state.",
		Example:                    fmt.Sprintf("%s query escrow escrows --seller <seller>", version.AppName),
		Args:                       cobra.ExactArgs(0),
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			seller, err := cmd.Flags().GetString(FlagSeller)
			if err != nil {
				return sdkerrors.Wrap(err, "Invalid seller address")
			}

			state, err := cmd.Flags().GetString(FlagState)
			if err != nil {
				return sdkerrors.Wrap(err, "Invalid escrow state")
			}

			objectKey, err := cmd.Flags().GetString(FlagObjectKey)
			if err != nil {
				return sdkerrors.Wrap(err, "Invalid object key")
			}

			paginationStart, err := cmd.Flags().GetUint64(FlagPaginationStart)
			if err != nil {
				return sdkerrors.Wrap(err, "Invalid pagination starting index")
			}
			paginationLength, err := cmd.Flags().GetUint64(FlagPaginationLength)
			if err != nil {
				return sdkerrors.Wrap(err, "Invalid pagination length")
			}

			queryClient := types.NewQueryClient(clientCtx)
			param := types.QueryEscrowsRequest{
				Seller:           seller,
				State:            state,
				ObjectKey:        objectKey,
				PaginationStart:  paginationStart,
				PaginationLength: paginationLength,
			}
			response, err := queryClient.Escrows(context.Background(), &param)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}

	escrowQueryCmd.Flags().AddFlagSet(FsQueryEscrows)
	flags.AddQueryFlagsToCmd(escrowQueryCmd)

	return escrowQueryCmd
}
