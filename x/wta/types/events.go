package types

// DONTCOVER

const (
	EventTypeBuyTicket     = "buy_ticket"
	EventTypePrizeIncrease = "prize_increase"
	EventTypeWinnerDrawn   = "winner_drawn"
	EventTypeNewDraw       = "new_draw"

	AttributeKeyTicketID        = "ticket_id"
	AttributeKeyTicketBuyer     = "ticket_buyer"
	AttributeKeyTicketTimestamp = "ticket_timestamp"
	AttributeKeyPrizeAmount     = "prize_amount"
	AttributeKeyWinnerAddress   = "winner_address"
	AttributeKeyWonAmount       = "won_amount"
	AttributeKeyDrawClosing     = "draw_closing"
)
