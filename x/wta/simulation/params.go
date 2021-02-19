package simulation

// DONTCOVER

import (
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
				return params.PrizePercentage.String()
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.CommunityPoolPercentageParamKey),
			func(r *rand.Rand) string {
				return params.CommunityPoolPercentage.String()
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.BurnPercentageParamKey),
			func(r *rand.Rand) string {
				return params.BurnPercentage.String()
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.DrawDurationParamKey),
			func(r *rand.Rand) string {
				return params.DrawDuration.String()
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.TicketPriceParamKey),
			func(r *rand.Rand) string {
				return params.TicketPrice.String()
			},
		),
	}
}
