package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// ParamChanges returns a randomly generated set or parameter changes
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreDistributionParamsKey),
			func(r *rand.Rand) string {
				params := RandomDistributionParams(r)
				bz, _ := json.Marshal(&params)
				return string(bz)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreDrawParamsKey),
			func(r *rand.Rand) string {
				params := RandomDrawParams(r)
				return fmt.Sprintf(`{"duration":"%d"}`, params.Duration)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreTicketParamsKey),
			func(r *rand.Rand) string {
				params := RandomTicketParams(r)
				bz, _ := json.Marshal(&params)
				return string(bz)
			},
		),
	}
}
