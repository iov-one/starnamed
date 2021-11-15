package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"

	"github.com/iov-one/starnamed/x/configuration/types"
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
const defaultString = ""

func getCmdUpdateFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fees",
		Short: "update fees using a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return fmt.Errorf("unable to get context: %s", err)
			}
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
				Configurer: cliCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid tx: %w", err)
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().String("fees-file", "fees.json", "fees file in json format")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdUpdateConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config",
		Short: "update domain configuration, provide the values you want to override in current configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return fmt.Errorf("unable to get context: %s", err)
			}
			config := &types.Config{}
			if !cliCtx.GenerateOnly {
				rawCfg, _, err := cliCtx.QueryStore([]byte(types.ConfigKey), types.StoreKey)
				if err != nil {
					return err
				}
				interfaceRegistry := cdctypes.NewInterfaceRegistry()
				interfaceRegistry.RegisterInterface("cosmos.base.v1beta1.Msg",
					(*sdk.Msg)(nil),
					&types.MsgUpdateConfig{},
					&types.MsgUpdateFees{},
				)
				marshaler := codec.NewProtoCodec(interfaceRegistry)
				marshaler.MustUnmarshalBinaryBare(rawCfg, config)
			}
			var signer string
			// if tx is not generate only, use --from flag as signer, otherwise get it from signer flag
			if !cliCtx.GenerateOnly {
				signer = cliCtx.GetFromAddress().String()
			} else {
				signer, err = cmd.Flags().GetString("signer")
				if err != nil {
					return err
				}
			}
			// get flags
			configurerStr, err := cmd.Flags().GetString("configurer")
			if err != nil {
				return
			}
			if configurerStr != defaultString {
				config.Configurer = configurerStr
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
				config.DomainRenewalPeriod = domainRenew
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
				config.DomainGracePeriod = domainGracePeriod
			}
			accountRenewPeriod, err := cmd.Flags().GetDuration("account-renew-period")
			if err != nil {
				return err
			}
			if accountRenewPeriod != defaultNumber {
				config.AccountRenewalPeriod = accountRenewPeriod
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
				config.AccountGracePeriod = accountGracePeriod
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

			escrowBroker, err := cmd.Flags().GetString("escrow-broker")
			if err != nil {
				return err
			}
			if escrowBroker != defaultString {
				config.EscrowBroker = escrowBroker
			}

			escrowCommission, err := cmd.Flags().GetString("escrow-commission")
			if err != nil {
				return err
			}
			if escrowCommission != defaultString {
				config.EscrowCommission, err = sdk.NewDecFromStr(escrowCommission)
				if err != nil {
					return sdkerrors.Wrap(err, "invalid escrow commission")
				}
			}

			escrowMaxPeriod, err := cmd.Flags().GetDuration("escrow-max-period")
			if err != nil {
				return err
			}
			if escrowMaxPeriod != defaultNumber {
				config.EscrowMaxPeriod = escrowMaxPeriod
			}

			if err := config.Validate(); err != nil {
				return err
			}
			// build msg
			msg := types.MsgUpdateConfig{
				Signer:           signer,
				NewConfiguration: config,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	// add flags
	cmd.Flags().String("signer", defaultString, "current configuration owner, for offline usage, otherwise --from is used")
	cmd.Flags().String("configurer", defaultString, "new configuration owner")
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

	cmd.Flags().Duration("escrow-max-period", defaultDuration, "maximum allowed duration for an escrow")
	cmd.Flags().String("escrow-commission", defaultString, "commission that will be received by the broker. The number represent the fraction of the price that will be sent to the broker account, it must be between 0 and 1.")
	cmd.Flags().String("escrow-broker", defaultString, "bech32 encoded address of the broker account")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
