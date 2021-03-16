package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/starnamed/x/starname/types"
	"github.com/spf13/cobra"
)

// GetTxCmd clubs together all the CLI tx commands
func GetTxCmd() *cobra.Command {
	domainTxCmd := &cobra.Command{
		Use:                        types.DomainStoreKey,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.DomainStoreKey),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	domainTxCmd.AddCommand(
		getCmdRegisterDomain(),
		getCmdAddAccountCertificate(),
		getCmdTransferAccount(),
		getCmdTransferDomain(),
		getmCmdSetAccountResources(),
		getCmdDeleteDomain(),
		getCmdDeleteAccount(),
		getCmdRenewDomain(),
		getCmdRenewAccount(),
		getCmdDeleteAccountCertificate(),
		getCmdRegisterAccount(),
		getCmdSetAccountMetadata(),
	)
	return domainTxCmd
}

func getCmdTransferDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain-transfer",
		Aliases: []string{"dt", "transfer-domain", "td"},
		Short:   "transfer a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			newOwner, err := cmd.Flags().GetString("new-owner")
			if err != nil {
				return err
			}
			// get transfer flag
			transferFlag, err := cmd.Flags().GetInt("transfer-flag")
			if err != nil {
				return err
			}
			// get sdk.AccAddress from string
			newOwnerAddr, err := sdk.AccAddressFromBech32(newOwner)
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgTransferDomain{
				Domain:       domain,
				Owner:        clientCtx.GetFromAddress(),
				NewAdmin:     newOwnerAddr,
				TransferFlag: types.TransferFlag(transferFlag),
				Payer:        feePayer,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "the domain name to transfer")
	cmd.Flags().StringP("new-owner", "o", "", "the new owner address in bech32 format")
	cmd.Flags().IntP("transfer-flag", "t", types.TransferResetNone, fmt.Sprintf("transfer flags for a domain"))
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdTransferAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-transfer",
		Aliases: []string{"at", "transfer-account", "ta"},
		Short:   "transfer an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			newOwner, err := cmd.Flags().GetString("new-owner")
			if err != nil {
				return err
			}
			// get sdk.AccAddress from string
			newOwnerAddr, err := sdk.AccAddressFromBech32(newOwner)
			if err != nil {
				return err
			}

			reset, err := cmd.Flags().GetString("reset")
			if err != nil {
				return err
			}
			var resetBool bool
			if resetBool, err = strconv.ParseBool(reset); err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgTransferAccount{
				Domain:   domain,
				Name:     name,
				Owner:    clientCtx.GetFromAddress(),
				NewOwner: newOwnerAddr,
				ToReset:  resetBool,
				Payer:    feePayer,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "the domain name of account")
	cmd.Flags().StringP("name", "n", "", "the name of the account you want to transfer")
	cmd.Flags().StringP("new-owner", "o", "", "the new owner address in bech32 format")
	cmd.Flags().StringP("reset", "r", "false", "true: reset all data associated with the account, false: preserves the data")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getmCmdSetAccountResources() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-resources-set",
		Aliases: []string{"ars", "account-set-resources", "asr", "set-resources", "sr", "replace-resources", "rr", "account-resources"},
		Short:   "set resources for an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			resourcesPath, err := cmd.Flags().GetString("src")
			if err != nil {
				return err
			}
			// open resources file
			f, err := os.Open(resourcesPath)
			if err != nil {
				return err
			}
			defer f.Close()
			// unmarshal resources
			var resources []*types.Resource
			err = json.NewDecoder(f).Decode(&resources)
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgReplaceAccountResources{
				Domain:       domain,
				Name:         name,
				NewResources: resources,
				Owner:        clientCtx.GetFromAddress(),
				Payer:        feePayer,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "the domain name of account")
	cmd.Flags().StringP("name", "n", "", "the name of the account whose resources you want to replace")
	cmd.Flags().StringP("src", "r", "resources.json", "the file containing the new resources in json format")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdDeleteDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain-delete",
		Aliases: []string{"dd", "delete-domain", "del-domain"},
		Short:   "delete a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgDeleteDomain{
				Domain: domain,
				Owner:  clientCtx.GetFromAddress().String(),
				Payer:  feePayer.String(),
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "name of the domain you want to delete")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdDeleteAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-delete",
		Aliases: []string{"ad", "account-del", "delete-account", "da", "del-account", "starname-delete", "delete-starname"},
		Short:   "delete an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgDeleteAccount{
				Domain: domain,
				Name:   name,
				Owner:  clientCtx.GetFromAddress(),
				Payer:  feePayer,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "the domain name of account")
	cmd.Flags().StringP("name", "n", "", "the name of the account you want to delete")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdRenewDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain-renew",
		Aliases: []string{"dn", "renew-domain", "nd"},
		Short:   "renew a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgRenewDomain{
				Domain: domain,
				Signer: clientCtx.GetFromAddress(),
				Payer:  feePayer,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "name of the domain you want to renew")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdRenewAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-renew",
		Aliases: []string{"an", "renew-account", "na"},
		Short:   "renew an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgRenewAccount{
				Domain: domain,
				Name:   name,
				Signer: clientCtx.GetFromAddress(),
				Payer:  feePayer,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "domain name of the account")
	cmd.Flags().StringP("name", "n", "", "account name you want to renew")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdDeleteAccountCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-certificate-delete",
		Aliases: []string{"acd", "account-delete-certificate", "adc", "delete-certificate", "dc", "del-certs", "certificate-delete", "cd"},
		Short:   "delete a certificate from an account",
		Long:    "delete a certificate from an account; either use the --certificate or --certificate-file flag",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			cert, err := cmd.Flags().GetBytesBase64("certificate")
			if err != nil {
				return err
			}
			certFile, err := cmd.Flags().GetString("certificate-file")
			if err != nil {
				return err
			}

			var c []byte
			switch {
			case len(cert) == 0 && len(certFile) == 0:
				return ErrCertificateNotProvided
			case len(cert) != 0 && len(certFile) != 0:
				return ErrCertificateProvideOnlyOne
			case len(cert) != 0 && len(certFile) == 0:
				c = cert
			case len(cert) == 0 && len(certFile) != 0:
				cf, err := os.Open(certFile)
				if err != nil {
					return err
				}
				cfb, err := ioutil.ReadAll(cf)
				if err != nil {
					return err
				}
				var j json.RawMessage
				if err := json.Unmarshal(cfb, &j); err != nil {
					return err
				}
				c = j
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgDeleteAccountCertificate{
				Domain:            domain,
				Name:              name,
				Owner:             clientCtx.GetFromAddress(),
				DeleteCertificate: c,
				Payer:             feePayer,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "domain name of the account")
	cmd.Flags().StringP("name", "n", "", "account name")
	cmd.Flags().BytesBase64P("certificate", "c", []byte{}, "certificate you want to add in base64 encoded format")
	cmd.Flags().StringP("certificate-file", "f", "", "directory of certificate file")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdAddAccountCertificate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-certificate-add",
		Aliases: []string{"aca", "account-add-certificate", "aac", "add-certificate", "ac", "add-certs", "certificate-add", "ca"},
		Short:   "add a certificate to an account",
		Long:    "add a certificate of an account; either use the --certificate or --certificate-file flag",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			cert, err := cmd.Flags().GetBytesBase64("certificate")
			if err != nil {
				return err
			}
			certFile, err := cmd.Flags().GetString("certificate-file")
			if err != nil {
				return err
			}

			var c json.RawMessage
			switch {
			case len(cert) == 0 && len(certFile) == 0:
				return ErrCertificateNotProvided
			case len(cert) != 0 && len(certFile) != 0:
				return ErrCertificateProvideOnlyOne
			case len(cert) != 0 && len(certFile) == 0:
				c = cert
			case len(cert) == 0 && len(certFile) != 0:
				cf, err := os.Open(certFile)
				if err != nil {
					return sdkerrors.Wrapf(ErrInvalidCertificate, "err: %s", err)
				}
				cfb, err := ioutil.ReadAll(cf)
				if err != nil {
					return sdkerrors.Wrapf(ErrInvalidCertificate, "err: %s", err)
				}
				if err := json.Unmarshal(cfb, &c); err != nil {
					return sdkerrors.Wrapf(ErrInvalidCertificate, "err: %s", err)
				}
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgAddAccountCertificate{
				Domain:         domain,
				Name:           name,
				Owner:          clientCtx.GetFromAddress(),
				NewCertificate: c,
				Payer:          feePayer,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "domain of the account")
	cmd.Flags().StringP("name", "n", "", "name of the account")
	cmd.Flags().BytesBase64P("certificate", "c", []byte{}, "certificate json you want to add in base64 encoded format")
	cmd.Flags().StringP("certificate-file", "f", "", "directory of certificate file in json format")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// getCmdRegisterAccount is the cli command to register accounts
func getCmdRegisterAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-register",
		Aliases: []string{"ar", "register-account", "ra"},
		Short:   "register an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			owner, err := cmd.Flags().GetString("owner")
			if err != nil {
				return err
			}
			var ownerAddr sdk.AccAddress
			if owner == "" {
				ownerAddr = clientCtx.GetFromAddress()
			} else {
				// get sdk.AccAddress from string
				ownerAddr, err = sdk.AccAddressFromBech32(owner)
				if err != nil {
					return err
				}
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			brokerStr, err := cmd.Flags().GetString("broker")
			if err != nil {
				return err
			}
			var broker sdk.AccAddress
			if brokerStr != "" {
				broker, err = sdk.AccAddressFromBech32(brokerStr)
				if err != nil {
					return err
				}
			}
			// build msg
			msg := &types.MsgRegisterAccount{
				Domain:     domain,
				Name:       name,
				Owner:      ownerAddr,
				Registerer: clientCtx.GetFromAddress(),
				Payer:      feePayer,
				Broker:     broker,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().StringP("domain", "d", "", "the existing domain for your account")
	cmd.Flags().StringP("name", "n", "", "the name of your account")
	cmd.Flags().StringP("owner", "o", "", "the address of the owner, if no owner provided signer is the owner")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	cmd.Flags().StringP("broker", "r", "", "address of the broker, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdRegisterDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain-register",
		Aliases: []string{"dr", "register-domain", "rd"},
		Short:   "register a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			dType, err := cmd.Flags().GetString("type")
			if err != nil {
				return err
			}

			if err := types.ValidateDomainType(types.DomainType(dType)); err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			brokerStr, err := cmd.Flags().GetString("broker")
			if err != nil {
				return err
			}
			var broker sdk.AccAddress
			if brokerStr != "" {
				broker, err = sdk.AccAddressFromBech32(brokerStr)
				if err != nil {
					return err
				}
			}
			msg := &types.MsgRegisterDomain{
				Name:       domain,
				Admin:      clientCtx.GetFromAddress().String(),
				DomainType: types.DomainType(dType),
				Broker:     broker.String(),
				Payer:      feePayer.String(),
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// add flags
	cmd.Flags().StringP("domain", "d", "", "name of the domain you want to register")
	cmd.Flags().StringP("type", "t", types.ClosedDomain, "type of the domain")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	cmd.Flags().StringP("broker", "r", "", "address of the broker, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdSetAccountMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-metadata-set",
		Aliases: []string{"ams", "account-set-metadata", "asm", "set-metadata", "sm", "set-account-metadata", "sam", "account-metadata"},
		Short:   "set metadata for an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			metadata, err := cmd.Flags().GetString("metadata")
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("payer")
			if err != nil {
				return err
			}
			var feePayer sdk.AccAddress
			if feePayerStr != "" {
				feePayer, err = sdk.AccAddressFromBech32(feePayerStr)
				if err != nil {
					return err
				}
			}
			msg := &types.MsgReplaceAccountMetadata{
				Domain:         domain,
				Name:           name,
				Owner:          clientCtx.GetFromAddress(),
				Payer:          feePayer,
				NewMetadataURI: metadata,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// add flags
	cmd.Flags().StringP("domain", "d", "", "the domain name of account")
	cmd.Flags().StringP("name", "n", "", "the name of the account whose resources you want to replace")
	cmd.Flags().StringP("metadata", "m", "", "the new metadata, leave empty to unset")
	cmd.Flags().StringP("payer", "p", "", "address of the fee payer, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
