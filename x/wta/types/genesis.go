package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState returns a new GenesisState containing the provided data
func NewGenesisState(draw Draw, tickets []Ticket, params Params) *GenesisState {
	return &GenesisState{
		Draw:    draw,
		Tickets: tickets,
		Params:  params,
	}
}

// DefaultGenesisState returns a default GenesisState
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		NewDraw(sdk.NewCoins(), time.Now().Add(time.Hour*24)),
		[]Ticket{},
		DefaultParams(),
	)
}

// ValidateGenesis validates the given genesis state and returns an error if something is invalid
func ValidateGenesis(state *GenesisState) error {
	// Validate the draw
	if !state.Draw.EndTime.IsZero() {
		err := state.Draw.Validate()
		if err != nil {
			return err
		}
	}

	// Validate the tickets
	for _, t := range state.Tickets {
		err := t.Validate()
		if err != nil {
			return err
		}

		// Check id duplicates
		if IsTicketIDDuplicated(t.Id, state.Tickets) {
			return fmt.Errorf("ticket id duplicated: %s", t.Id)
		}
	}

	return nil
}
