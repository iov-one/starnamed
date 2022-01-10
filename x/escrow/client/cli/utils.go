package cli

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"

	"github.com/iov-one/starnamed/x/escrow/types"
)

// NewMsgCreateEscrow creates a types.MsgCreateEscrow including the seller, the fee payer, the price, the deadline
// and the given object. This method calls ValidateBasic on the created message.
// The AddCreateEscrowFlags function has to be called on the same cmd object before.
func NewMsgCreateEscrow(ctx client.Context, cmd *cobra.Command, obj types.TransferableObject) (*types.MsgCreateEscrow, error) {
	seller := ctx.GetFromAddress().String()
	feePayer, err := cmd.Flags().GetString(FlagFeePayer)
	if err != nil {
		return nil, err
	}

	priceStr, err := verifyErrAndNonEmpty(cmd, FlagPrice)
	if err != nil {
		return nil, err
	}
	price, err := sdk.ParseCoinsNormalized(priceStr)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid price : %v", priceStr)
	}

	deadlineStr, err := verifyErrAndNonEmpty(cmd, FlagDeadline)
	if err != nil {
		return nil, err
	}

	deadline, err := parseDeadline(deadlineStr)
	if err != nil {
		return nil, err
	}

	isAuction, err := cmd.Flags().GetBool(FlagIsAuction)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgCreateEscrow(seller, feePayer, obj, price, deadline, isAuction)

	// check if valid
	if err = msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return &msg, nil
}

// AddCreateEscrowFlags adds the flags used by NewMsgCreateEscrow to the given cmd.Flag() flag set
func AddCreateEscrowFlags(cmd *cobra.Command) {
	addCommonFlags(cmd.Flags())
	cmd.Flags().Bool(FlagIsAuction, false, "Specifies if the escrow should act as an auction, with the deadline being the end of the auction and the price the initial price")
	cmd.Flags().String(FlagPrice, "", "Price of the object")
	cmd.Flags().String(FlagDeadline, "", "Expiration date of the escrow, in the RFC3339 time format")
}

func verifyErrAndNonEmpty(cmd *cobra.Command, flag string) (string, error) {
	val, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", err
	}
	if len(val) == 0 {
		return "", fmt.Errorf("you must provide a %v", flag)
	}
	return val, nil
}

func parseDeadline(date string) (uint64, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return 0, sdkerrors.Wrapf(err, "The expiration date is not in RFC3339 format : %v", date)
	}

	return uint64(t.Unix()), nil
}
