package simulation

import (
	"github.com/cosmicbet/ledger/x/wta/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"math/rand"
	"time"
)

var (
	hexLetters = []rune("abcdef0123456789")
)

// RandDate generates a new Date that does not exceed the one set
func RandDate(r *rand.Rand, max time.Time) time.Time {
	return time.Date(
		r.Intn(max.Year()),
		time.Month(r.Intn(int(max.Month()))),
		r.Intn(max.Day()),
		r.Intn(max.Hour()),
		r.Intn(max.Minute()),
		r.Intn(max.Second()),
		r.Intn(max.Nanosecond()),
		max.Location(),
	)
}

// RandHexString generates a random hex string of given length
func RandHexString(r *rand.Rand, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = hexLetters[r.Intn(len(hexLetters))]
	}
	return string(b)
}

func RandCoint(r *rand.Rand, amountLimit int64) sdk.Coin {
	return sdk.NewCoin(sdk.DefaultBondDenom, simtypes.RandomAmount(r, sdk.NewInt(amountLimit)))
}

// -------------------------------------------------------------------------------------------------------------------

// RandomDraw generates a new random types.Draw object that has an ending not going above the limit provided
func RandomDraw(r *rand.Rand, limitTime time.Time) types.Draw {
	return types.NewDraw(
		r.Uint32(),
		r.Uint32(),
		sdk.NewCoins(RandCoint(r, 1000000)),
		RandDate(r, limitTime),
	)
}

// -------------------------------------------------------------------------------------------------------------------

// RandTicket generates a random ticket for the given address
func RandTicket(r *rand.Rand, owner sdk.AccAddress) types.Ticket {
	return types.NewTicket(
		RandHexString(r, 20),
		RandDate(r, time.Now()),
		owner.String(),
	)
}

// RandTicketsSlice generates a slice of random tickets of the given length
func RandTicketsSlice(r *rand.Rand, length int, accounts []simtypes.Account) []types.Ticket {
	tickets := make([]types.Ticket, length)
	for i := range tickets {
		owner := accounts[r.Intn(len(accounts))]
		tickets[i] = RandTicket(r, owner.Address)
	}
	return tickets
}

// -------------------------------------------------------------------------------------------------------------------

// RandHistoricalDrawData returns a randomly generated HistoricalDrawData
func RandHistoricalDrawData(r *rand.Rand, accounts []simtypes.Account) types.HistoricalDrawData {
	// 50% chance of not having any ticket
	var winningTicket *types.Ticket
	if r.Intn(100) > 50 {
		randTicket := RandTicket(r, accounts[r.Intn(len(accounts))].Address)
		winningTicket = &randTicket
	}

	return types.NewHistoricalDrawData(
		RandomDraw(r, time.Now().Add(-time.Minute*10)),
		winningTicket,
	)
}

// RandHistoricalDrawsData returns a randomly generated slice of types.HistoricalDrawData of the given length
func RandHistoricalDrawsData(r *rand.Rand, length int, accounts []simtypes.Account) []types.HistoricalDrawData {
	data := make([]types.HistoricalDrawData, length)
	for i := range data {
		data[i] = RandHistoricalDrawData(r, accounts)
	}
	return data
}

// -------------------------------------------------------------------------------------------------------------------

// RandomParams returns a randomly generated parameters
func RandomParams(r *rand.Rand) types.Params {
	prizePercentage := r.Int63n(100)
	poolPercentage := r.Int63n(100 - prizePercentage)
	burnPercentage := (100 - prizePercentage) - poolPercentage

	return types.NewParams(
		sdk.NewInt(prizePercentage),
		sdk.NewInt(poolPercentage),
		sdk.NewInt(burnPercentage),
		time.Second*time.Duration(r.Int63n(60)),
		RandCoint(r, 1000),
	)
}
