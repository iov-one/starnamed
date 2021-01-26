package starname

import (
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/iov-one/starnamed/x/starname/keeper"
	"github.com/iov-one/starnamed/x/starname/types"
)

// NewHandler builds the tx requests handler for the starname module
func NewHandler(k *Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	f := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		/* TODO: FIXME
		switch msg := msg.(type) {
		// domain handlers
		case *types.MsgTransferDomain:
			return handlerMsgTransferDomain(ctx, *k, msg)
		// account handlers
		case *types.MsgRenewAccount:
			return handlerMsgRenewAccount(ctx, *k, msg)
		case *types.MsgDeleteAccountCertificate:
			return handlerMsgDeleteAccountCertificate(ctx, *k, msg)
		case *types.MsgDeleteAccount:
			return handlerMsgDeleteAccount(ctx, *k, msg)
		case *types.MsgReplaceAccountResources:
			return handlerMsgReplaceAccountResources(ctx, *k, msg)
		case *types.MsgTransferAccount:
			return handlerMsgTransferAccount(ctx, *k, msg)
		case *types.MsgReplaceAccountMetadata:
			return handlerMsgReplaceAccountMetadata(ctx, *k, msg)
		}
		*/

		var (
			res proto.Message
			err error
		)
		switch msg := msg.(type) {
		// domain msgs
		case *types.MsgDeleteDomain:
			res, err = msgServer.DeleteDomain(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgRegisterDomain:
			res, err = msgServer.RegisterDomain(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgRenewDomain:
			res, err = msgServer.RenewDomain(sdk.WrapSDKContext(ctx), msg)
		// account msgs
		case *types.MsgAddAccountCertificate:
			res, err = msgServer.AddAccountCertificate(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgRegisterAccount:
			res, err = msgServer.RegisterAccount(sdk.WrapSDKContext(ctx), msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unregonized request: %T", msg))
		}

		return sdk.WrapServiceResult(ctx, res, err)
	}

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		/*
			TODO
			remove when cosmos sdk decides that you are allowed to panic on errors that should not happen
			instead of returning random internal errors that mean actually nothing to a developer without
			a stacktrace or at least the error string of the panic itself, and also substitute 'log' stdlib
			with cosmos sdk logger when they make clear how you can use it and how to set up env to achieve so
		*/
		defer func() {
			if r := recover(); r != nil {
				log.Printf("FATAL-PANIC while executing message: %#v\nReason: %v", msg, r)
				// and lets panic again to throw it back to cosmos sdk yikes.
				panic(r)
			}
		}()
		resp, err := f(ctx, msg)
		if err != nil {
			msg := fmt.Sprintf("tx rejected %T: %s", msg, err)
			k.Logger(ctx).With("module", types.ModuleName).Info(msg)
		}
		return resp, err
	}
}
