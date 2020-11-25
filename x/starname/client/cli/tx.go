package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/starname/types"
	"github.com/spf13/cobra"
)

// GetTxCmd clubs together all the CLI tx commands
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	domainTxCmd := &cobra.Command{
		Use:                        storeKey,
		Short:                      fmt.Sprintf("%s transactions subcommands", storeKey),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		Aliases:                    []string{"starname"},
	}

	domainTxCmd.AddCommand(flags.PostCommands(
		getCmdRegisterDomain(cdc),
		getCmdAddAccountCerts(cdc),
		getCmdTransferAccount(cdc),
		getCmdTransferDomain(cdc),
		getmCmdReplaceAccountResources(cdc),
		getCmdDelDomain(cdc),
		getCmdDelAccount(cdc),
		getCmdRenewDomain(cdc),
		getCmdRenewAccount(cdc),
		getCmdDelAccountCerts(cdc),
		getCmdRegisterAccount(cdc),
		getCmdSetAccountMetadata(cdc),
	)...)
	return domainTxCmd
}

func getCmdTransferDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-domain",
		Short: "transfer a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				FeePayerAddr: feePayer,
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
	cmd.Flags().String("domain", "", "the domain name to transfer")
	cmd.Flags().String("new-owner", "", "the new owner address in bech32 format")
	cmd.Flags().Int("transfer-flag", types.TransferResetNone, fmt.Sprintf("transfer flags for a domain"))
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	return cmd
}

func getCmdTransferAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-account",
		Short: "transfer an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				Domain:       domain,
				Name:         name,
				Owner:        clientCtx.GetFromAddress(),
				NewOwner:     newOwnerAddr,
				ToReset:      resetBool,
				FeePayerAddr: feePayer,
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
	cmd.Flags().String("domain", "", "the domain name of account")
	cmd.Flags().String("name", "", "the name of the account you want to transfer")
	cmd.Flags().String("new-owner", "", "the new owner address in bech32 format")
	cmd.Flags().String("reset", "false", "true: reset all data associated with the account, false: preserves the data")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	return cmd
}

func getmCmdReplaceAccountResources(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace-resources",
		Short: "replace account resources",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				FeePayerAddr: feePayer,
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
	cmd.Flags().String("domain", "", "the domain name of account")
	cmd.Flags().String("name", "", "the name of the account whose resources you want to replace")
	cmd.Flags().String("src", "resources.json", "the file containing the new resources in json format")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	// return cmd
	return cmd
}

func getCmdDelDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del-domain",
		Short: "delete a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				Domain:       domain,
				Owner:        clientCtx.GetFromAddress(),
				FeePayerAddr: feePayer,
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
	cmd.Flags().String("domain", "", "name of the domain you want to delete")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	return cmd
}

func getCmdDelAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del-account",
		Short: "delete an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				Domain:       domain,
				Name:         name,
				Owner:        clientCtx.GetFromAddress(),
				FeePayerAddr: feePayer,
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
	cmd.Flags().String("domain", "", "the domain name of account")
	cmd.Flags().String("name", "", "the name of the account you want to delete")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	return cmd
}

func getCmdRenewDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew-domain",
		Short: "renew a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				Domain:       domain,
				Signer:       clientCtx.GetFromAddress(),
				FeePayerAddr: feePayer,
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
	cmd.Flags().String("domain", "", "name of the domain you want to renew")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	// return
	return cmd
}

func getCmdRenewAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew-account",
		Short: "renew an account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				Domain:       domain,
				Name:         name,
				Signer:       clientCtx.GetFromAddress(),
				FeePayerAddr: feePayer,
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
	cmd.Flags().String("domain", "", "domain name of the account")
	cmd.Flags().String("name", "", "account name you want to renew")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	return cmd
}

func getCmdDelAccountCerts(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del-certs",
		Short: "delete certificates of an account",
		Long:  "delete certificates of an account. Either use cert or cert-file flags",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			cert, err := cmd.Flags().GetBytesBase64("cert")
			if err != nil {
				return err
			}
			certFile, err := cmd.Flags().GetString("cert-file")
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				FeePayerAddr:      feePayer,
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
	cmd.Flags().String("domain", "", "domain name of the account")
	cmd.Flags().String("name", "", "account name")
	cmd.Flags().BytesBase64("cert", []byte{}, "certificate you want to add in base64 encoded format")
	cmd.Flags().String("cert-file", "", "directory of certificate file")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	// return cmd
	return cmd
}

func getCmdAddAccountCerts(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-certs",
		Short: "add certificates to account",
		Long:  "add certificates of an account. Either use cert or cert-file flags",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			cert, err := cmd.Flags().GetBytesBase64("cert")
			if err != nil {
				return err
			}
			certFile, err := cmd.Flags().GetString("cert-file")
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
			msg := &types.MsgAddAccountCertificates{
				Domain:         domain,
				Name:           name,
				Owner:          clientCtx.GetFromAddress(),
				NewCertificate: c,
				FeePayerAddr:   feePayer,
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
	cmd.Flags().String("domain", "", "domain of the account")
	cmd.Flags().String("name", "", "name of the account")
	cmd.Flags().BytesBase64("cert", []byte{}, "certificate json you want to add in base64 encoded format")
	cmd.Flags().String("cert-file", "", "directory of certificate file in json format")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	return cmd
}

// getCmdRegisterAccount is the cli command to register accounts
func getCmdRegisterAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "register-account",
		Short:                      "register an account",
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				Domain:       domain,
				Name:         name,
				Owner:        ownerAddr,
				Registerer:   clientCtx.GetFromAddress(),
				FeePayerAddr: feePayer,
				Broker:       broker,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("domain", "", "the existing domain name for your account")
	cmd.Flags().String("name", "", "the name of your account")
	cmd.Flags().String("owner", "", "the address of the owner, if no owner provided signer is the owner")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	cmd.Flags().String("broker", "", "address of the broker, optional")
	return cmd
}

func getCmdRegisterDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-domain",
		Short: "register a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			feePayerStr, err := cmd.Flags().GetString("fee-payer")
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
				Name:         domain,
				Admin:        clientCtx.GetFromAddress(),
				DomainType:   types.DomainType(dType),
				Broker:       broker,
				FeePayerAddr: feePayer,
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
	cmd.Flags().String("domain", "", "name of the domain you want to register")
	cmd.Flags().String("type", types.ClosedDomain, "type of the domain")
	cmd.Flags().String("fee-payer", "", "address of the fee payer, optional")
	cmd.Flags().String("broker", "", "address of the broker, optional")
	return cmd
}

func getCmdSetAccountMetadata(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-account-metadata",
		Short: "sets account metadata",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
			feePayerStr, err := cmd.Flags().GetString("fee_payer")
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
				FeePayerAddr:   feePayer,
				NewMetadataURI: metadata,
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
	cmd.Flags().String("domain", "", "the domain name of account")
	cmd.Flags().String("name", "", "the name of the account whose resources you want to replace")
	cmd.Flags().String("metadata", "", "the new metadata URI, leave empty to unset")
	cmd.Flags().String("fee_payer", "", "address of the fee payer, optional")
	// return cmd
	return cmd
}
