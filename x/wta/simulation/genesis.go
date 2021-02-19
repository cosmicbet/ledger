package simulation

import (
	"fmt"
	"github.com/cosmicbet/ledger/x/wta/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"time"
)

// RandomizedGenState sets into the given simState a randomly generated genesis state
func RandomizedGenState(simState *module.SimulationState) {
	genesisState := types.NewGenesisState(
		RandomDraw(simState.Rand, time.Now().Add(time.Minute*1)),
		RandTicketsSlice(simState.Rand, 20, simState.Accounts),
		RandHistoricalDrawsData(simState.Rand, 50, simState.Accounts),
		RandomParams(simState.Rand),
	)

	bz, err := simState.Cdc.MarshalJSON(genesisState)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)
}
