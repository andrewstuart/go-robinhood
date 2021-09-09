package robinhood

import (
	"fmt"
	"strings"
)

// OptionDirection is a type for whether an option order is opening or closing
// an option position.
type OptionDirection int

// MarshalJSON implements json.Marshaler.
func (o OptionDirection) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", strings.ToLower(o.String()))), nil
}

//go:generate stringer -type OptionDirection
// The two directions.
const (
	Debit OptionDirection = iota
	Credit
)
