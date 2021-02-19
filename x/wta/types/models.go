package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTicket allows to build a new Ticket instance.
func NewTicket(id string, timestamp time.Time, owner string) Ticket {
	return Ticket{
		Id:        id,
		Owner:     owner,
		Timestamp: timestamp,
	}
}

// Validate returns an error if there is something wrong inside t
func (t *Ticket) Validate() error {
	if t.Id == "" {
		return fmt.Errorf("invalid ticket id: %s", t.Id)
	}

	if t.Timestamp.IsZero() {
		return fmt.Errorf("invalid ticket creation time: %s", t.Timestamp.Format(time.RFC3339))
	}

	if t.Owner == "" {
		return fmt.Errorf("invalid ticket owner: %s", t.Owner)
	}

	return nil
}

// MarshalTicket marshals the given ticket to a slice of bytes
func MarshalTicket(cdc codec.BinaryMarshaler, ticket Ticket) ([]byte, error) {
	return cdc.MarshalBinaryBare(&ticket)
}

// MustMarshalTicket marshals the given ticket into a slice of bytes, and panics on error
func MustMarshalTicket(cdc codec.BinaryMarshaler, ticket Ticket) []byte {
	bz, err := MarshalTicket(cdc, ticket)
	if err != nil {
		panic(err)
	}
	return bz
}

// UnmarshalTicket reads the provided byte array as a Ticket object
func UnmarshalTicket(cdc codec.BinaryMarshaler, bz []byte) (Ticket, error) {
	var ticket Ticket
	err := cdc.UnmarshalBinaryBare(bz, &ticket)
	return ticket, err
}

// MustUnmarshalTicket unmarshals the given byte slice into a Ticket object, and panics on error
func MustUnmarshalTicket(cdc codec.BinaryMarshaler, bz []byte) Ticket {
	ticket, err := UnmarshalTicket(cdc, bz)
	if err != nil {
		panic(err)
	}

	return ticket
}

// IsTicketIDDuplicated tells whether or not the given id is duplicated inside the provided slice
func IsTicketIDDuplicated(id string, slice []Ticket) bool {
	var count = 0
	for _, ticket := range slice {
		if ticket.Id == id {
			count++
		}
	}
	return count > 1
}

// ------------------------------------------------------------------------------------------------------------------

// EmptyDraw returns a new Draw that will end on the given time
func EmptyDraw(endTime time.Time) Draw {
	return Draw{
		EndTime: endTime,
	}
}

// NewDraw allows to build a new Draw instance
func NewDraw(participants, ticketsSold uint32, prize sdk.Coins, endTime time.Time) Draw {
	return Draw{
		Participants: participants,
		TicketsSold:  ticketsSold,
		Prize:        prize,
		EndTime:      endTime,
	}
}

// Validate returns an error if there is something wrong with the provided Draw
func (d Draw) Validate() error {
	err := d.Prize.Validate()
	if err != nil {
		return err
	}

	if d.EndTime.IsZero() {
		return fmt.Errorf("invalid draw end time")
	}

	return nil
}

// Equal tells whether d and e contain the same data
func (d Draw) Equal(e Draw) bool {
	return d.Participants == e.Participants &&
		d.TicketsSold == d.TicketsSold &&
		d.Prize.IsEqual(e.Prize) &&
		d.EndTime.Equal(e.EndTime)
}

// MarshalDraw marshals the given Draw as a byte array
func MarshalDraw(cdc codec.BinaryMarshaler, draw Draw) ([]byte, error) {
	return cdc.MarshalBinaryBare(&draw)
}

// MustMarshalDraw marshals the given Draw as a byte array and panics on error
func MustMarshalDraw(cdc codec.BinaryMarshaler, draw Draw) []byte {
	bz, err := MarshalDraw(cdc, draw)
	if err != nil {
		panic(err)
	}
	return bz
}

// UnmarshalDraw reads the given byte slice as a Draw object
func UnmarshalDraw(cdc codec.BinaryMarshaler, bz []byte) (Draw, error) {
	var draw Draw
	err := cdc.UnmarshalBinaryBare(bz, &draw)
	return draw, err
}

// MustUnmarshalDraw unmarshals the given byte slice into a Draw object
func MustUnmarshalDraw(cdc codec.BinaryMarshaler, bz []byte) Draw {
	draw, err := UnmarshalDraw(cdc, bz)
	if err != nil {
		panic(err)
	}

	return draw
}

// ------------------------------------------------------------------------------------------------------------------

// NewHistoricalDrawData creates a new HistoricalDrawData
func NewHistoricalDrawData(draw Draw, winningTicket *Ticket) HistoricalDrawData {
	return HistoricalDrawData{
		Draw:          draw,
		WinningTicket: winningTicket,
	}
}

func (h *HistoricalDrawData) Validate() error {
	err := h.Draw.Validate()
	if err != nil {
		return err
	}

	if h.WinningTicket != nil {
		err = h.WinningTicket.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// MarshalHistoricalDraw marshals the given historical draw as a byte array
func MarshalHistoricalDraw(cdc codec.BinaryMarshaler, draw HistoricalDrawData) ([]byte, error) {
	return cdc.MarshalBinaryBare(&draw)
}

// MustMarshalHistoricalDraw marshals the given draws as a byte array
func MustMarshalHistoricalDraw(cdc codec.BinaryMarshaler, draw HistoricalDrawData) []byte {
	bz, err := MarshalHistoricalDraw(cdc, draw)
	if err != nil {
		panic(err)
	}

	return bz
}

// UnmarshalHistoricalDraw unmarshals the given byte array as a HistoricalDrawData object
func UnmarshalHistoricalDraw(cdc codec.BinaryMarshaler, bz []byte) (HistoricalDrawData, error) {
	var draws HistoricalDrawData
	err := cdc.UnmarshalBinaryBare(bz, &draws)
	if err != nil {
		return HistoricalDrawData{}, err
	}

	return draws, nil
}

// MustUnmarshalHistoricalDrawData unmarshals the given byte array as a HistoricalDrawData object and panics on errors
func MustUnmarshalHistoricalDrawData(cdc codec.BinaryMarshaler, bz []byte) HistoricalDrawData {
	data, err := UnmarshalHistoricalDraw(cdc, bz)
	if err != nil {
		panic(err)
	}
	return data
}
