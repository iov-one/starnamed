package cli

import (
	"fmt"
	"io/fs"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/iov-one/starnamed/x/offchain/types"
)

// GetTxCmd clubs together all the CLI tx commands
func GetTxCmd() *cobra.Command {
	domainTxCmd := &cobra.Command{
		Use:                        "offchain",
		Short:                      "Offchain transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	domainTxCmd.AddCommand(
		getCmdVerifyMsg(),
		getCmdSignMsg(),
	)
	return domainTxCmd
}
func getCmdVerifyMsg() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify file",
		Aliases: []string{"v", "verify-msg", "verify-message", "vm"},
		Short:   "verify an arbitrary message signature",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			verifier := types.NewVerifier(clientCtx.TxConfig.SignModeHandler())

			jsonOutput, err := cmd.Flags().GetBool("json")
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "you must provide the file containing the signed transacation")
			}

			file := args[0]

			data, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			txData, err := clientCtx.TxConfig.TxJSONDecoder()(data)

			err = verifier.Verify(txData)
			if err != nil {
				return err
			}

			msgs := make([]*types.MsgSignData, len(txData.GetMsgs()))

			for i, msg := range txData.GetMsgs() {
				msgSign, ok := msg.(*types.MsgSignData)
				if !ok {
					return fmt.Errorf("invalid message type %T at index %v", msg, i)
				}
				msgs[i] = msgSign
			}

			if jsonOutput {
				data = clientCtx.JSONMarshaler.MustMarshalJSON(&types.ListOfMsgSignData{
					Msgs: msgs,
				})

				fmt.Println(string(data))
			} else {
				fmt.Printf("%v messages in transaction\n", len(msgs))
				for index, msg := range msgs {
					fmt.Printf("%v\n\tSigner: %s\n\tData: %s\n", index, msg.Signer, msg.Data)
				}
			}
			return nil
		},
	}

	cmd.Flags().Bool("json", false, "If that flag is present, the output will be in JSON format. "+
		"The JSON output is a ListOfMsgSignData serialized object, with base64 encoded data and bech32 signed address")

	return cmd
}

func getCmdSignMsg() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sign destination-file",
		Aliases: []string{"s", "sign-message", "sign-msg"},
		Short:   "sign an arbitrary message",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signerAddr := clientCtx.GetFromAddress()
			if signerAddr == nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "you must provide the --from flag")
			}

			if len(args) != 1 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "you must provide a destination file for the signed message")
			}
			destFile := args[0]

			file, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}

			stringData, err := cmd.Flags().GetString("text")
			if err != nil {
				return err
			}

			var data []byte
			if len(file) != 0 && len(stringData) != 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "you must provide only one of --file or --text flags")
			} else if len(file) != 0 {
				data, err = ioutil.ReadFile(file)
				if err != nil {
					return err
				}
			} else if len(stringData) != 0 {
				data = []byte(stringData)
			} else {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "you must provide at least one of --file or --text flags")
			}

			msg := types.NewMsgSignData(signerAddr, data)

			signer, err := types.NewSignerFromClientContext(clientCtx)
			if err != nil {
				return sdkerrors.Wrapf(err, "error while creating the signer object")
			}

			txObj, err := signer.Sign([]sdk.Msg{msg})
			if err != nil {
				return err
			}

			txData, err := clientCtx.TxConfig.TxJSONEncoder()(txObj)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(destFile, txData, fs.ModePerm)
		},
	}
	// add flags
	addKeyFlags(cmd.Flags())
	cmd.Flags().StringP("file", "f", "", "File to sign (only one of file or text must be provided)")
	cmd.Flags().StringP("text", "t", "", "Text data to sign (only one of file or text must be provided)")
	return cmd
}

// addKeyFlags adds all flags relating to authentication and private key management
func addKeyFlags(fs *pflag.FlagSet) {
	fs.String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	fs.String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test|memory)")
	fs.Bool(flags.FlagUseLedger, false, "Use a connected Ledger device")
	fs.String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
}
