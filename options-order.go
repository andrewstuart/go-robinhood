package robinhood

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// OptionsOrderOpts encapsulates common Options order choices
type OptionsOrderOpts struct {
	Quantity    float64
	Price       float64
	Direction   OptionDirection
	TimeInForce TimeInForce
	Type        OrderType
	Side        OrderSide
}

// OrderOptions places a new order for options
func (c *Client) OrderOptions(q *OptionInstrument, o OptionsOrderOpts) (json.RawMessage, error) {
	b := optionInput{
		Account:     c.Account.URL,
		Direction:   o.Direction,
		TimeInForce: o.TimeInForce,
		Legs: []Leg{{
			Option:         q.URL,
			RatioQuantity:  1,
			Side:           o.Side,
			PositionEffect: "open",
		}},
		Trigger:  "immediate",
		Type:     o.Type,
		Quantity: o.Quantity,
		Price:    o.Price,
		RefID:    uuid.New().String(),
	}

	if o.Side != Buy {
		b.Legs[0].PositionEffect = "close"
	}

	bs, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", EPOptions+"orders/", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var out json.RawMessage
	err = c.DoAndDecode(req, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetOptionsOrders returns all outstanding options orders
func (c *Client) GetOptionsOrders() (json.RawMessage, error) {
	var o json.RawMessage
	err := c.GetAndDecode(EPOptions+"orders/", &o)
	if err != nil {
		return nil, err
	}

	return o, nil

}
