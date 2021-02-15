package types

const (
	// ModuleName is the name of the wta module
	ModuleName = "wta"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the wta module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the wta module
	RouterKey = ModuleName

	PrizeCollectorName = "wta_prize_collector"
	PrizeBurnerName    = "wta_prize_burner"
)

var (
	DrawStoreKey       = []byte{0x1}
	TicketsStorePrefix = []byte("wta_tickets")
)

// TicketsStoreKey returns the store key used to save the ticket with the given id
func TicketsStoreKey(id string) []byte {
	return append(TicketsStorePrefix, []byte(id)...)
}
