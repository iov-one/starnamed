package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/iov-one/starnamed/x/configuration/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd builds all the query commands for the module
func GetQueryCmd() *cobra.Command {
	// group config queries under a sub-command
	configQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	// add queries
	configQueryCmd.AddCommand(
		getCmdQueryConfig(),
		getCmdQueryFees(),
	)
	// return cmd list
	return configQueryCmd
}

// getCmdQueryConfig returns the command to get the configuration
func getCmdQueryConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-config",
		Short: "gets the current configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			path := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryConfig)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return err
			}
			fmt.Println(string(resp)) // TODO: here and others: handle output format proper. See clientCtx.PrintOutput(string(res))
			return nil
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func getCmdQueryFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-fees",
		Short: "gets the current fees",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			path := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryFees)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return err
			}
			fmt.Println(string(resp)) // TODO: here and others: handle output format proper. See clientCtx.PrintOutput(string(res))
			return nil
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
