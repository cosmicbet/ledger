package types

import (
	"fmt"
	"time"
)

// NewGenesisState returns a new GenesisState containing the provided data
func NewGenesisState(drawEndTime time.Time, tickets []Ticket, pastDraws []HistoricalDrawData, params Params) *GenesisState {
	return &GenesisState{
		DrawEndTime: drawEndTime,
		Tickets:     tickets,
		PastDraws:   pastDraws,
		Params:      params,
	}
}

// DefaultGenesisState returns a default GenesisState
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		time.Now().Add(time.Hour*24),
		[]Ticket{},
		[]HistoricalDrawData{},
		DefaultParams(),
	)
}

// ValidateGenesis validates the given genesis state and returns an error if something is invalid
func ValidateGenesis(state *GenesisState) error {
	// Validate the draw
	if state.DrawEndTime.IsZero() || time.Now().After(state.DrawEndTime) {
		return fmt.Errorf("invalid draw end time: %s", state.DrawEndTime.Format(time.RFC3339))
	}

	// Validate the tickets
	for _, t := range state.Tickets {
		err := t.Validate()
		if err != nil {
			return err
		}

		// Check that the timestamp is not after the current draw
		if t.Timestamp.After(state.DrawEndTime) {
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
