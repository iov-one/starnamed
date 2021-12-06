package keeper

import (
	"github.com/iov-one/starnamed/x/wasm/types"
)

var _ types.QueryServer = Keeper{}
