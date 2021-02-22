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
	params := RandomParams(r)

	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.PrizePercentageParamKey),
			func(r *rand.Rand) string {
				return fmt.Sprintf(`"%s"`, params.PrizePercentage)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.CommunityPoolPercentageParamKey),
			func(r *rand.Rand) string {
				return fmt.Sprintf(`"%s"`, params.CommunityPoolPercentage)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.BurnPercentageParamKey),
			func(r *rand.Rand) string {
				return fmt.Sprintf(`"%s"`, params.BurnPercentage)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.DrawDurationParamKey),
			func(r *rand.Rand) string {
				return fmt.Sprintf(`"%d"`, params.DrawDuration)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.TicketPriceParamKey),
			func(r *rand.Rand) string {
				bz, err := json.Marshal(params.TicketPrice)
				if err != nil {
					panic(err)
				}
				return string(bz)
			},
		),
	}
}
