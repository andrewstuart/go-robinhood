package robinhood

// TimeInForce is the time in force for an order.
type TimeInForce int

//go:generate stringer -type=TimeInForce
// Well-known values for TimeInForce
const (
	GTC TimeInForce = iota
	GFD
	IOC
	OPG
)
