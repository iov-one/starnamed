package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	IDParam = "id"
)

// RegisterHandlers defines routes that get registered by the main application
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// UpdateEscrowReq defines the properties of an escrow update request's body.
type UpdateEscrowReq struct {
	BaseReq              rest.BaseReq `json:"base_req" yaml:"base_req"`
	Sender               string       `json:"sender" yaml:"sender"`
	To                   string       `json:"to" yaml:"to"`
	ReceiverOnOtherChain string       `json:"receiver_on_other_chain" yaml:"receiver_on_other_chain"`
	SenderOnOtherChain   string       `json:"sender_on_other_chain" yaml:"sender_on_other_chain"`
	Amount               sdk.Coins    `json:"amount" yaml:"amount"`
	HashLock             string       `json:"hash_lock" yaml:"hash_lock"`
	TimeLock             uint64       `json:"time_lock" yaml:"time_lock"`
	Timestamp            uint64       `json:"timestamp" yaml:"timestamp"`
	Transfer             bool         `json:"transfer" yaml:"transfer"`
}

// TransferToEscrowReq defines the properties of a transfer to an escrow request's body.
type TransferToEscrowReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Id  string       `json:"id" yaml:"id"`
	Amount  sdk.Coins       `json:"id" yaml:"id"`
}

type
