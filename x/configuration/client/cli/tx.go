package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CosmWasm/wasmd/x/configuration/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetTxCmd clubs together all the CLI tx commands
func GetTxCmd() *cobra.Command {
	configTxCmd := &cobra.Command{
		Use:                        types.StoreKey,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.StoreKey),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	configTxCmd.AddCommand(
		getCmdUpdateConfig(),
		getCmdUpdateFees(),
	)
	return configTxCmd
}

var defaultDuration, _ = time.ParseDuration("1h")

const defaultRegex = "^(.*?)?"
const defaultNumber = 1

func getCmdUpdateFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fees",
		Short: "update fees using a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			// get fees file
			feeFile, err := cmd.Flags().GetString("fees-file")
			if err != nil {
				return err
			}
			f, err := os.Open(feeFile)
			if err != nil {
				return fmt.Errorf("unable to open fee file: %s", err)
			}
			defer f.Close()
			newFees := new(types.Fees)
			err = json.NewDecoder(f).Decode(newFees)
			if err != nil {
				return err
			}
			msg := types.MsgUpdateFees{
				Fees:       newFees,
				Configurer: cliCtx.GetFromAddress(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid tx: %w", err)
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().String("fees-file", "fees.json", "fees file in json format")
	return cmd
}

func getCmdUpdateConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config",
		Short: "update domain configuration, provide the values you want to override in current configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err = client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			config := &types.Config{}
			if !cliCtx.GenerateOnly {
				rawCfg, _, err := cliCtx.QueryStore([]byte(types.ConfigKey), types.StoreKey)
				if err != nil {
					return err
				}
				types.ModuleCdc.MustUnmarshalBinaryBare(rawCfg, config)
			}
			var signer sdk.AccAddress
			// if tx is not generate only, use --from flag as signer, otherwise get it from signer flag
			if !cliCtx.GenerateOnly {
				signer = cliCtx.FromAddress
			} else {
				signerStr, err := cmd.Flags().GetString("signer")
				if err != nil {
					return err
				}
				signer, err = sdk.AccAddressFromBech32(signerStr)
				if err != nil {
					return err
				}
			}
			// get flags
			configurerStr, err := cmd.Flags().GetString("configurer")
			if err != nil {
				return
			}
			if configurerStr != "" {
				configurer, err := sdk.AccAddressFromBech32(configurerStr)
				if err != nil {
					return err
				}
				config.Configurer = configurer
			}
			validDomainName, err := cmd.Flags().GetString("valid-domain-name")
			if err != nil {
				return err
			}
			if validDomainName != defaultRegex {
				config.ValidDomainName = validDomainName
			}
			validAccountName, err := cmd.Flags().GetString("valid-account-name")
			if err != nil {
				return err
			}
			if validAccountName != defaultRegex {
				config.ValidAccountName = validAccountName
			}
			validURI, err := cmd.Flags().GetString("valid-uri")
			if err != nil {
				return err
			}
			if validURI != defaultRegex {
				config.ValidURI = validURI
			}
			validResource, err := cmd.Flags().GetString("valid-resource")
			if err != nil {
				return err
			}
			if validResource != defaultRegex {
				config.ValidResource = validResource
			}
			domainRenew, err := cmd.Flags().GetDuration("domain-renew-period")
			if err != nil {
				return err
			}
			if domainRenew != defaultDuration {
				config.DomainRenewalPeriod.Seconds = int64(domainRenew.Seconds())
			}
			domainRenewCountMax, err := cmd.Flags().GetUint32("domain-renew-count-max")
			if err != nil {
				return err
			}
			if domainRenewCountMax != defaultNumber {
				config.DomainRenewalCountMax = domainRenewCountMax
			}
			domainGracePeriod, err := cmd.Flags().GetDuration("domain-grace-period")
			if err != nil {
				return err
			}
			if domainGracePeriod != defaultNumber {
				config.DomainGracePeriod.Seconds = int64(domainGracePeriod.Seconds())
			}
			accountRenewPeriod, err := cmd.Flags().GetDuration("account-renew-period")
			if err != nil {
				return err
			}
			if accountRenewPeriod != defaultNumber {
				config.AccountRenewalPeriod.Seconds = int64(accountRenewPeriod.Seconds())
			}
			accountRenewCountMax, err := cmd.Flags().GetUint32("account-renew-count-max")
			if err != nil {
				return err
			}
			if accountRenewCountMax != defaultNumber {
				config.AccountRenewalCountMax = accountRenewCountMax
			}
			accountGracePeriod, err := cmd.Flags().GetDuration("account-grace-period")
			if err != nil {
				return err
			}
			if accountGracePeriod != defaultDuration {
				config.AccountGracePeriod.Seconds = int64(accountGracePeriod.Seconds())
			}
			resourceMax, err := cmd.Flags().GetUint32("resource-max")
			if err != nil {
				return err
			}
			if resourceMax != defaultNumber {
				config.ResourcesMax = resourceMax
			}
			certificateSizeMax, err := cmd.Flags().GetUint64("certificate-size-max")
			if err != nil {
				return err
			}
			if certificateSizeMax != defaultNumber {
				config.CertificateSizeMax = certificateSizeMax
			}
			certificateCountMax, err := cmd.Flags().GetUint32("certificate-count-max")
			if err != nil {
				return err
			}
			if certificateCountMax != defaultNumber {
				config.CertificateCountMax = certificateCountMax
			}
			metadataSizeMax, err := cmd.Flags().GetUint64("metadata-size-max")
			if err != nil {
				return err
			}
			if metadataSizeMax != defaultNumber {
				config.MetadataSizeMax = metadataSizeMax
			}

			if err := config.Validate(); err != nil {
				return err
			}
			// build msg
			msg := &types.MsgUpdateConfig{
				Signer:           signer,
				NewConfiguration: config,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().String("signer", "", "current configuration owner, for offline usage, otherwise --from is used")
	cmd.Flags().String("configurer", "", "new configuration owner")
	cmd.Flags().String("valid-domain-name", defaultRegex, "regexp that determines if domain name is valid or not")
	cmd.Flags().String("valid-account-name", defaultRegex, "regexp that determines if account name is valid or not")
	cmd.Flags().String("valid-uri", defaultRegex, "regexp that determines if uri is valid or not")
	cmd.Flags().String("valid-resource", defaultRegex, "regexp that determines if resource is valid or not")

	cmd.Flags().Duration("domain-renew-period", defaultDuration, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint32("domain-renew-count-max", uint32(defaultNumber), "maximum number of applicable domain renewals")
	cmd.Flags().Duration("domain-grace-period", defaultDuration, "domain grace period duration in seconds")

	cmd.Flags().Duration("account-renew-period", defaultDuration, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint32("account-renew-count-max", uint32(defaultNumber), "maximum number of applicable account renewals")
	cmd.Flags().Duration("account-grace-period", defaultDuration, "account grace period duration in seconds")

	cmd.Flags().Uint32("resource-max", uint32(defaultNumber), "maximum number of resources could be saved under an account")
	cmd.Flags().Uint64("certificate-size-max", uint64(defaultNumber), "maximum size of a certificate that could be saved under an account")
	cmd.Flags().Uint32("certificate-count-max", uint32(defaultNumber), "maximum number of certificates that could be saved under an account")
	cmd.Flags().Uint64("metadata-size-max", uint64(defaultNumber), "maximum size of metadata that could be saved under an account")
	return cmd
}
