package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// RandomizedGenState sets into the given simState a randomly generated genesis state
func RandomizedGenState(simState *module.SimulationState) {
	// Create a random genesis state and serialize that
	genesisState := types.NewGenesisState(
		RandomDraw(simState.Rand, time.Now().Add(time.Minute*1)),
		RandTicketsSlice(simState.Rand, 20, simState.Accounts),
		RandHistoricalDrawsData(simState.Rand, 50, simState.Accounts),
		RandomParams(simState.Rand),
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)

	// Log the params
	bz, err := json.MarshalIndent(&genesisState.Params, "", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	// Update the coins supply and the prize collector balance based on the generated draw prize
	var bankState banktypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[banktypes.ModuleName], &bankState)

	bankState.Supply = bankState.Supply.Add(genesisState.Draw.Prize...)
	bankState.Balances = append(bankState.Balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(types.PrizeCollectorName).String(),
		Coins:   genesisState.Draw.Prize,
	})

	simState.GenState[banktypes.ModuleName] = simState.Cdc.MustMarshalJSON(&bankState)
}
