package roho

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
	// GTC means Good 'Til Cancelled.
	GTC TimeInForce = iota
	// GFD means Good For Day.
	GFD
	// IOC means Immediate Or Cancel.
	IOC
	// OPG means Opening (of market).
	OPG
	// FOK means Fill Or Kill.
	FOK
)
