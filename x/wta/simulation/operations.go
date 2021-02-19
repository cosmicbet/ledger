package simulation

// DONTCOVER

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmicbet/ledger/app/params"
	"github.com/cosmicbet/ledger/x/wta/keeper"
	"github.com/cosmicbet/ledger/x/wta/types"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	OpWeightBuyTickets = "op_weight_buy_tickets"
	DefaultGasValue    = 200000
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONMarshaler,
	k keeper.Keeper, ak authkeeper.AccountKeeper, bk bankkeeper.Keeper,
) sim.WeightedOperations {
	var weightBuyTickets int
	appParams.GetOrGenerate(cdc, OpWeightBuyTickets, &weightBuyTickets, nil,
		func(_ *rand.Rand) {
			weightBuyTickets = params.DefaultWeightMsgBuyTickets
		},
	)

	return sim.WeightedOperations{
		sim.NewWeightedOperation(
			weightBuyTickets,
			SimulateMsgBuyTickets(k, ak, bk),
		),
	}
}

// SimulateMsgBuyTickets generates a random types.MsgBuyTickets and sends it to the chain.
func SimulateMsgBuyTickets(k keeper.Keeper, ak authkeeper.AccountKeeper, bk bankkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accounts []simtypes.Account, chainID string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {

		// Get random message data and build the message
		acc, ticketsQuantity, ticketsCost, skip := randomBuyTicketsData(r, ctx, accounts, k, bk)
		if skip {
			return simtypes.NoOpMsg(types.RouterKey, types.ModuleName, "skipping after generating data"), nil, nil
		}
		msg := types.NewMsgBuyTickets(ticketsQuantity, acc.Address.String())

		// Send the message
		err = sendMsgBuyTickets(r, app, ak, bk, msg, ticketsCost, ctx, chainID, []cryptotypes.PrivKey{acc.PrivKey})
		if err != nil {
			return simtypes.NoOpMsg(types.RouterKey, types.ModuleName, ""), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// randomBuyTicketsData generates random parameters that can be used to create a types.MsgBuyTickets.
// It returns a random amount of tickets to buy, as well as the account that should buy them and
// the overall cost of the operation
func randomBuyTicketsData(
	r *rand.Rand, ctx sdk.Context, accounts []simtypes.Account, k keeper.Keeper, bk bankkeeper.Keeper,
) (account simtypes.Account, ticketsAmt uint32, ticketsCost sdk.Coin, skip bool) {
	// Get a random account
	account, _ = simtypes.RandomAcc(r, accounts)

	// Get a random number of tickets (min 1, max 10 tickets)
	ticketsAmt = uint32(r.Int31n(10)) + 1

	// Compute the ticket cost based on the params
	wtaParams := k.GetParams(ctx)
	ticketsCost = sdk.NewCoin(wtaParams.TicketPrice.Denom, wtaParams.TicketPrice.Amount.MulRaw(int64(ticketsAmt)))

	// Make sure the account has enough balance to pay for the tickets
	balance := bk.SpendableCoins(ctx, account.Address)
	if balance.IsZero() || sdk.NewCoins(ticketsCost).IsAnyGT(balance) {
		return simtypes.Account{}, 0, sdk.Coin{}, true
	}

	return account, ticketsAmt, ticketsCost, false
}

// sendMsgBuyTickets sends a transaction with a types.MsgBuyTickets from a provided random profile.
func sendMsgBuyTickets(
	r *rand.Rand, app *baseapp.BaseApp, ak authkeeper.AccountKeeper, bk bankkeeper.Keeper,
	msg *types.MsgBuyTickets, ticketsCost sdk.Coin, ctx sdk.Context, chainID string, privkeys []cryptotypes.PrivKey,
) error {
	addr, _ := sdk.AccAddressFromBech32(msg.Buyer)
	account := ak.GetAccount(ctx, addr)

	// Compute the amount of fees that the account can spend based on the amount of money it will spend on tickets
	coins := bk.SpendableCoins(ctx, account.GetAddress())
	fees, err := simtypes.RandomFees(r, ctx, coins.Sub(sdk.NewCoins(ticketsCost)))
	if err != nil {
		return err
	}

	txGen := simappparams.MakeTestEncodingConfig().TxConfig
	tx, err := helpers.GenTx(
		txGen,
		[]sdk.Msg{msg},
		fees,
		DefaultGasValue,
		chainID,
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		privkeys...,
	)
	if err != nil {
		return err
	}

	_, _, err = app.Deliver(txGen.TxEncoder(), tx)
	if err != nil {
		return err
	}

	return nil
}
