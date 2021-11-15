package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/iov-one/starnamed/x/escrow/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// query an escrow
	r.HandleFunc(fmt.Sprintf("/%s/escrow/{%s}", types.ModuleName, IDParam), queryEscrowHandlerFn(cliCtx)).Methods("GET")
	// do a query over all the escrows
	r.HandleFunc(fmt.Sprintf("/%s/escrows", types.ModuleName), queryEscrowsHandlerFn(cliCtx)).Methods("GET")

}

func queryEscrowHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryEscrowParams{
			Id: vars[IDParam],
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryEscrow)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryEscrowsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var paginationStart, paginationLength uint64
		startParam := r.URL.Query().Get("pagination_start")
		lengthParam := r.URL.Query().Get("pagination_length")
		if len(startParam) != 0 {
			var err error
			paginationStart, err = strconv.ParseUint(startParam, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid pagination_start value : "+err.Error())
				return
			}
		}
		if len(lengthParam) != 0 {
			var err error
			paginationLength, err = strconv.ParseUint(lengthParam, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid pagination_length value : "+err.Error())
				return
			}
		}

		params := types.QueryEscrowsParams{
			Seller:           r.URL.Query().Get("seller"),
			State:            r.URL.Query().Get("state"),
			ObjectKey:        r.URL.Query().Get("object"),
			PaginationStart:  paginationStart,
			PaginationLength: paginationLength,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryEscrows)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
