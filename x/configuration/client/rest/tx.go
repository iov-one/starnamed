package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/iov-one/iovns/x/configuration/types"
)

// handleTxRequest is a helper function that takes care of checking base requests, sdk messages, after verifying
// requests it forwards an error to the client in case of error, otherwise it will return a transaction to sign
// and send to the /tx endpoint to do a request
func handleTxRequest(cliCtx client.Context, baseReq rest.BaseReq, msg sdk.Msg, writer http.ResponseWriter) {
	baseReq = baseReq.Sanitize()
	if !baseReq.ValidateBasic(writer) {
		return
	}
	// validate request
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
	}
	// write tx
	tx.WriteGeneratedTxResponse(cliCtx, writer, baseReq, &msg)
}

type updateConfig struct {
	BaseReq rest.BaseReq           `json:"base_req"`
	Message *types.MsgUpdateConfig `json:"message"`
}

func updateConfigHandler(cliCtx client.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req updateConfig
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, "failed to parse request")
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

type updateFees struct {
	BaseReq rest.BaseReq         `json:"base_req"`
	Message *types.MsgUpdateFees `json:"message"`
}

func updateFeesHandler(cliCtx client.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req updateFees
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}
