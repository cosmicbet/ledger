package simulation

// DONTCOVER

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmicbet/ledger/x/wta/types"
)

// RandomizedGenState sets into the given simState a randomly generated genesis state
func RandomizedGenState(simState *module.SimulationState) {
	// Create a random genesis state and serialize that
	genesisState := types.NewGenesisState(
		RandDate(simState.Rand, time.Now().Add(time.Minute*1)),
		RandTicketsSlice(simState.Rand, 20, simState.Accounts),
		RandHistoricalDrawsData(simState.Rand, 50, simState.Accounts),
		RandomDistributionParams(simState.Rand),
		RandomDrawParams(simState.Rand),
		RandomTicketParams(simState.Rand),
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)

	// Update the coins supply and the prize collector balance based on the generated draw prize
	prize := RandCoin(simState.Rand, 100000)

	var bankState banktypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[banktypes.ModuleName], &bankState)

	bankState.Supply = bankState.Supply.Add(prize)
	bankState.Balances = append(bankState.Balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(types.PrizeCollectorName).String(),
		Coins:   sdk.NewCoins(prize),
	})

	simState.GenState[banktypes.ModuleName] = simState.Cdc.MustMarshalJSON(&bankState)
}
