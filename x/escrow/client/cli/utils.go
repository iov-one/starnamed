package cli

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/starnamed/x/escrow/types"
	"github.com/spf13/cobra"
)

func NewMsgCreateEscrow(ctx client.Context, cmd *cobra.Command, obj types.TransferableObject) (*types.MsgCreateEscrow, error) {
	seller := ctx.GetFromAddress().String()

	buyer, err := verifyErrAndNonEmpty(cmd, FlagBuyer)
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

	msg := types.NewMsgCreateEscrow(seller, buyer, obj, price, deadline)
	return &msg, nil
}

func AddCreateEscrowFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagBuyer, "", "TODO: desc")
	cmd.Flags().String(FlagPrice, "", "TODO: desc")
	cmd.Flags().String(FlagDeadline, "", "TODO: desc")
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
