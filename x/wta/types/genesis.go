package types

import (
	"fmt"
	"time"
)

// NewGenesisState returns a new GenesisState containing the provided data
func NewGenesisState(draw Draw, tickets []Ticket, pastDraws []HistoricalDrawData, params Params) *GenesisState {
	return &GenesisState{
		Draw:      draw,
		Tickets:   tickets,
		PastDraws: pastDraws,
		Params:    params,
	}
}

// DefaultGenesisState returns a default GenesisState
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		EmptyDraw(time.Now().Add(time.Hour*24)),
		[]Ticket{},
		[]HistoricalDrawData{},
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

		// Check that the timestamp is not after the current draw
		if t.Timestamp.After(state.Draw.EndTime) {
			return fmt.Errorf("ticket with id %s has creation date after the draw end time ", t.Id)
		}

		// Check that the timestamp is not in the future
		if t.Timestamp.After(time.Now()) {
			return fmt.Errorf("ticket with id %s has creation date set in the future", t.Id)
		}

		// Check id duplicates
		if IsTicketIDDuplicated(t.Id, state.Tickets) {
			return fmt.Errorf("ticket id duplicated: %s", t.Id)
		}
	}

	// Validate the historical draws data
	for _, data := range state.PastDraws {
		err := data.Validate()
		if err != nil {
			return err
		}
	}

	// Validate the params
	err := state.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
