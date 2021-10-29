package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/iov-one/starnamed/x/escrow/types"
)

const (
	CreateRoute     = "create"
	UpdateRoute     = "update"
	TransferToRoute = "transfer"
	RefundRoute     = "refund"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	escrowRouteTpl := fmt.Sprintf("/%s/{%s}/", types.ModuleName, IDParam)
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.ModuleName, CreateRoute), createEscrowHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(escrowRouteTpl+UpdateRoute, updateEscrowHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(escrowRouteTpl+TransferToRoute, transferToEscrowHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(escrowRouteTpl+RefundRoute, refundEscrowHandlerFn(cliCtx)).Methods("POST")
}

func updateEscrowHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getVar(r, w, IDParam)
		if len(id) == 0 {
			return
		}

		var req UpdateEscrowReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.MsgUpdateEscrow{
			Id:       id,
			Updater:  req.Updater,
			FeePayer: req.FeePayer,
			Seller:   req.Seller,
			Price:    req.Price,
			Deadline: req.Deadline,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func transferToEscrowHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getVar(r, w, IDParam)
		if len(id) == 0 {
			return
		}

		var req TransferToEscrowReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.MsgTransferToEscrow{
			Id:       id,
			Sender:   req.Sender,
			FeePayer: req.FeePayer,
			Amount:   req.Amount,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func refundEscrowHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getVar(r, w, IDParam)
		if len(id) == 0 {
			return
		}

		var req RefundEscrowReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.MsgRefundEscrow{
			Id:       id,
			Sender:   req.Sender,
			FeePayer: req.FeePayer,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func createEscrowHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req CreateEscrowReq
		if !rest.ReadRESTReq(writer, request, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(writer) {
			return
		}

		msg := types.NewMsgCreateEscrow(req.Seller, req.FeePayer, req.Object, req.Price, req.Deadline, req.IsAuction)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, writer, req.BaseReq, &msg)
	}
}

func getVar(r *http.Request, w http.ResponseWriter, name string) string {
	vars := mux.Vars(r)
	variable, present := vars[name]
	if !present {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "You must provide the escrow "+name)
		return ""
	}
	return variable
}
