package cli

import (
	"fmt"

	"github.com/CosmWasm/wasmd/pkg/queries"
	"github.com/CosmWasm/wasmd/x/configuration/types"
	"github.com/cosmos/cosmos-sdk/client"
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
	return &cobra.Command{
		Use:   "get-config",
		Short: "gets the current configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			path := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryConfig)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return err
			}
			var jsonResp types.QueryConfigResponse
			err = queries.DefaultQueryDecode(resp, &jsonResp)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(jsonResp.Config)
		},
	}
}

func getCmdQueryFees() *cobra.Command {
	return &cobra.Command{
		Use:   "get-fees",
		Short: "gets the current fees",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			path := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryFees)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return err
			}
			var jsonResp types.QueryFeesResponse
			err = queries.DefaultQueryDecode(resp, &jsonResp)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(jsonResp.Fees)
		},
	}
}
