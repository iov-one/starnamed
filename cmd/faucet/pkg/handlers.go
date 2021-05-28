package pkg

import (
	"encoding/json"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/prometheus/common/log"
)

// keeps tx manager and mutex locks sequence bump
type FaucetHandler struct {
	tm *TxManager
}

func NewFaucetHandler(tm *TxManager) *FaucetHandler {
	return &FaucetHandler{
		tm: tm,
	}
}

func jsonErr(w http.ResponseWriter, status int, msg string) {
	errJson := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(errJson)
	if err != nil {
		panic("Error while encoding JSON response")
	}
	return
}

func (f *FaucetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addrStr := r.URL.Query().Get("address")
	if addrStr == "" {
		jsonErr(w, http.StatusBadRequest, "provide a bech32 address")
		return
	}
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		log.Error(errors.Wrap(err, "incorrect bech32 address"))
		jsonErr(w, http.StatusBadRequest, "provide a bech32 address")
		return
	}

	tx, err := f.tm.BuildAndSignTx(addr)
	if err != nil {
		log.Error(errors.Wrap(err, "tx signing failed"))
		jsonErr(w, http.StatusInternalServerError, "internal error")
		return
	}

	res, err := f.tm.BroadcastTx(tx)
	if err != nil {
		log.Error(errors.Wrap(err, "broadcast tx failed"))
		jsonErr(w, http.StatusInternalServerError, "internal error")
		return
	}

	if res.Code != errors.SuccessABCICode {
		log.Error(errors.Wrap(err, "broadcast tx failed"))
		jsonErr(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := struct {
		Msg  string `json:"msg"`
		Hash string `json:"hash"`
	}{
		Msg:  "check your balance :-)",
		Hash: res.TxHash,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		panic("Error while encoding JSON response")
	}
	return
}
