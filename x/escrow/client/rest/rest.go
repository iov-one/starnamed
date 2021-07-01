package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/iov-one/starnamed/x/escrow/types"
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

// CreateEscrowReq defines the properties of a escrow creation request's body.
type CreateEscrowReq struct {
	BaseReq  rest.BaseReq             `json:"base_req" yaml:"base_req"`
	Seller   string                   `json:"seller" yaml:"seller"`
	FeePayer string                   `json:"fee_payer" yaml:"fee_payer"`
	Price    sdk.Coins                `json:"price" yaml:"price"`
	Deadline uint64                   `json:"expiration" yaml:"expiration"`
	Object   types.TransferableObject `json:"object" yaml:"object"`
}

// UpdateEscrowReq defines the properties of an escrow update request's body.
type UpdateEscrowReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Updater  string       `json:"updater" yaml:"updater"`
	FeePayer string       `json:"fee_payer" yaml:"fee_payer"`
	Seller   string       `json:"seller" yaml:"seller"`
	Price    sdk.Coins    `json:"price" yaml:"price"`
	Deadline uint64       `json:"expiration" yaml:"expiration"`
}

// TransferToEscrowReq defines the properties of a transfer to an escrow request's body.
type TransferToEscrowReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Sender   string       `json:"sender" yaml:"sender"`
	FeePayer string       `json:"fee_payer" yaml:"fee_payer"`
	Amount   sdk.Coins    `json:"amount" yaml:"amount"`
}

// RefundEscrowReq defines the properties of a escrow refund request's body.
type RefundEscrowReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Sender   string       `json:"sender" yaml:"sender"`
	FeePayer string       `json:"fee_payer" yaml:"fee_payer"`
}
