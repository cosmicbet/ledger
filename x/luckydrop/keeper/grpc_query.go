package keeper

import (
	"github.com/cosmicbet/ledger/x/luckydrop/types"
)

var _ types.QueryServer = Keeper{}
