package robinhood

import (
	"fmt"
	"strings"
)

// TimeInForce is the time in force for an order.
type TimeInForce int

// MarshalJSON implements json.Marshaler
func (t TimeInForce) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", strings.ToLower(t.String()))), nil
}

//go:generate stringer -type=TimeInForce
// Well-known values for TimeInForce
const (
	GTC TimeInForce = iota
	GFD
	IOC
	OPG
)
