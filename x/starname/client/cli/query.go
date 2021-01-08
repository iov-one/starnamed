package cli

import (
	"context"
	"errors"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/iov-one/starnamed/x/starname/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd builds the commands for queries in the domain module
func GetQueryCmd() *cobra.Command {
	domainQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "query commands for the starname module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	domainQueryCmd.AddCommand(
		getQueryResolveDomain(),
		getQueryResolveAccount(),
		getQueryDomainAccounts(),
		getQueryOwnerAccounts(),
		getQueryOwnerDomains(),
		getQueryResourceAccounts(),
	)
	return domainQueryCmd
}

func getQueryResolveDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resolve-domain",
		Aliases: []string{"rd", "domain-resolve", "dr", "domain-info", "di"},
		Short:   "resolve a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Domain(
				context.Background(),
				&types.QueryDomainRequest{
					Name: domain,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "the name of the domain that you want to resolve")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func getQueryDomainAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "accounts-by-domain",
		Aliases: []string{"abd", "domain-accounts", "da"},
		Short:   "get accounts in a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).DomainAccounts(
				context.Background(),
				&types.QueryDomainAccountsRequest{
					Domain:     domain,
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "the domain of interest")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "domain accounts")
	return cmd
}

func getQueryOwnerAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "accounts-by-owner",
		Aliases: []string{"abo", "owner-accounts", "oa"},
		Short:   "get accounts owned by an address",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			owner, err := cmd.Flags().GetString("address")
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).OwnerAccounts(
				context.Background(),
				&types.QueryOwnerAccountsRequest{
					Owner:      owner,
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	// add flags
	cmd.Flags().StringP("address", "a", "", "the bech32 address of the owner")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "owner accounts")
	return cmd
}

func getQueryOwnerDomains() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domains-by-owner",
		Aliases: []string{"dbo", "owner-domains", "od"},
		Short:   "get domains owned by an address",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			owner, err := cmd.Flags().GetString("address")
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).OwnerDomains(
				context.Background(),
				&types.QueryOwnerDomainsRequest{
					Owner:      owner,
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	// add flags
	cmd.Flags().StringP("address", "a", "", "the bech32 address of the owner")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "owner domains")
	return cmd
}

func getQueryResolveAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resolve-account",
		Aliases: []string{"ra", "account-resolve", "ar", "resolve", "resolve-starname", "starname-resolve", "sr"},
		Short:   "resolve an account by providing either --starname or --name and --domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			starname, err := cmd.Flags().GetString("starname")
			if err != nil {
				return err
			}
			if len(domain) > 0 {
				if len(starname) > 0 {
					return errors.New("either specify starname or name and domain")
				}

				starname = strings.Join([]string{name, domain}, types.StarnameSeparator)
			}
			// TODO: Validate() that starname is well formed
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Starname(
				context.Background(),
				&types.QueryStarnameRequest{
					Starname: starname,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	// add flags
	cmd.Flags().StringP("starname", "s", "", "the starname representation of the account, eg antoine*iov")
	cmd.Flags().StringP("domain", "d", "", "the domain of the account")
	cmd.Flags().StringP("name", "n", "", "the name of the account")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "resolve account")
	return cmd
}

func getQueryResourceAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "accounts-by-resource",
		Aliases: []string{"abr", "resource-accounts", "ra", "resolve-resource", "rr", "resource-resolve"},
		Short:   "get accounts associated with a resource",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			uri, err := cmd.Flags().GetString("uri")
			if err != nil {
				return err
			}
			resource, err := cmd.Flags().GetString("resource")
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).ResourceAccounts(
				context.Background(),
				&types.QueryResourceAccountsRequest{
					Uri:        uri,
					Resource:   resource,
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	// add flags
	cmd.Flags().String("uri", "", "the resource uri")
	cmd.Flags().String("resource", "", "resource")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "resource accounts")
	return cmd
}
