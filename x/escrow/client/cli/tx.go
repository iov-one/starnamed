package cli

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/iov-one/starnamed/x/escrow/types"
	"github.com/spf13/cobra"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	escrowTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Escrow transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	escrowTxCmd.AddCommand(
		GetCmdUpdateEscrow(),
		GetCmdTransferToEscrow(),
		GetCmdRefundEscrow(),
	)

	return escrowTxCmd
}

// GetCmdUpdateEscrow implements updating an escrow command
func GetCmdUpdateEscrow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Updates an escrow",
		Long: "Updates the fields of an escrow. Object is not modifiable, buyer address is modifiable by the buyer" +
			" and all the other fields are modifiable by the seller.",
		Example: fmt.Sprintf("$ %s tx escrow update <id> --price 5atom", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			updater := clientCtx.GetFromAddress().String()

			seller, err := cmd.Flags().GetString(FlagSeller)
			if err != nil {
				return err
			}
			buyer, err := cmd.Flags().GetString(FlagBuyer)
			if err != nil {
				return err
			}
			price, err := cmd.Flags().GetString(FlagPrice)
			if err != nil {
				return err
			}

			var priceCoins sdk.Coins
			if len(price) > 0 {
				priceCoins, err = sdk.ParseCoinsNormalized(price)
				if err != nil {
					return sdkerrors.Wrap(err, "Incorrect price format")
				}
			}

			deadlineStr, err := cmd.Flags().GetString(FlagDeadline)
			if err != nil {
				return err
			}

			var deadline uint64
			if len(deadlineStr) != 0 {
				deadline, err = parseDeadline(deadlineStr)
				if err != nil {
					return err
				}
			}

			msg := types.MsgUpdateEscrow{
				Id:       args[0],
				Updater:  updater,
				Seller:   seller,
				Buyer:    buyer,
				Price:    priceCoins,
				Deadline: deadline,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsEscrow)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferToEscrow implements transfering to an escrow command
func GetCmdTransferToEscrow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [id] [amount]",
		Short: "Transfers coins to an escrow",
		Long: "Transfer coins to an escrow, if the minimum price is not reached, the transaction will fail." +
			"Otherwise, an amount equal to the escrow price will be sent to the escrow and the exchange will" +
			"be done",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "Invalid amount format")
			}

			msg := types.MsgTransferToEscrow{
				Id:     args[0],
				Sender: sender,
				Amount: amount,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRefundEscrow implements refunding an escrow command
func GetCmdRefundEscrow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund [id]",
		Short: "Refund the engaged assets in an escrow",
		Long:  "Refund the engaged assets in an escrow.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			seller := clientCtx.GetFromAddress().String()

			msg := types.MsgRefundEscrow{
				Id:     args[0],
				Seller: seller,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
