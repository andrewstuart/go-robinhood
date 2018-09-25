package robinhood

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// OrderSide is which side of the trade an order is on
type OrderSide int

//go:generate stringer -type OrderSide
// Buy/Sell
const (
	Sell OrderSide = iota
	Buy
)

// OrderType represents a Limit or Market order
type OrderType int

//go:generate stringer -type OrderType
// Well-known order types. Default is Market.
const (
	Market OrderType = iota
	Limit
)

// OrderOpts encapsulates differences between order types
type OrderOpts struct {
	Side          OrderSide
	Type          OrderType
	Quantity      uint64
	Price         float64
	TimeInForce   TimeInForce
	ExtendedHours bool
	Stop, Force   bool
}

type apiOrder struct {
	Account       string  `json:"account,omitempty"`
	Instrument    string  `json:"instrument,omitempty"`
	Symbol        string  `json:"symbol,omitempty"`
	Type          string  `json:"type,omitempty"`
	TimeInForce   string  `json:"time_in_force,omitempty"`
	Trigger       string  `json:"trigger,omitempty"`
	Price         float64 `json:"price,omitempty"`
	StopPrice     float64 `json:"stop_price,omitempty"`
	Quantity      uint64  `json:"quantity,omitempty"`
	Side          string  `json:"side,omitempty"`
	ExtendedHours bool    `json:"extended_hours,omitempty"`

	OverrideDayTradeChecks bool `json:"override_day_trade_checks,omitempty"`
	OverrideDtbpChecks     bool `json:"override_dtbp_checks,omitempty"`
}

// Order places an order for a given instrument
func (c *Client) Order(i *Instrument, o OrderOpts) (*OrderOutput, error) {
	a := apiOrder{
		Account:       c.Account.URL,
		Instrument:    i.URL,
		Symbol:        i.Symbol,
		Type:          strings.ToLower(o.Type.String()),
		TimeInForce:   strings.ToLower(o.TimeInForce.String()),
		Quantity:      o.Quantity,
		Side:          strings.ToLower(o.Side.String()),
		ExtendedHours: o.ExtendedHours,
		Price:         o.Price,
		Trigger:       "immediate",
	}

	if o.Stop {
		a.StopPrice = o.Price
		a.Trigger = "stop"
	}

	bs, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	post, err := http.NewRequest("POST", epOrders, bytes.NewReader(bs))
	post.Header.Add("Content-Type", "application/json")

	var out OrderOutput
	err = c.DoAndDecode(post, &out)
	if err != nil {
		return nil, err
	}

	out.client = c
	return &out, nil
}

// OrderOutput is the response from the Order api
type OrderOutput struct {
	Meta
	Account                string   `json:"account"`
	AveragePrice           float64  `json:"average_price"`
	CancelURL              string   `json:"cancel"`
	CreatedAt              string   `json:"created_at"`
	CumulativeQuantity     string   `json:"cumulative_quantity"`
	Executions             []string `json:"executions"`
	ExtendedHours          bool     `json:"extended_hours"`
	Fees                   string   `json:"fees"`
	ID                     string   `json:"id"`
	Instrument             string   `json:"instrument"`
	LastTransactionAt      string   `json:"last_transaction_at"`
	OverrideDayTradeChecks bool     `json:"override_day_trade_checks"`
	OverrideDtbpChecks     bool     `json:"override_dtbp_checks"`
	Position               string   `json:"position"`
	Price                  float64  `json:"price,string"`
	Quantity               string   `json:"quantity"`
	RejectReason           string   `json:"reject_reason"`
	Side                   string   `json:"side"`
	State                  string   `json:"state"`
	StopPrice              float64  `json:"stop_price"`
	TimeInForce            string   `json:"time_in_force"`
	Trigger                string   `json:"trigger"`
	Type                   string   `json:"type"`

	client *Client
}

// Cancel attempts to cancel an odrer
func (o OrderOutput) Cancel() error {
	post, err := http.NewRequest("POST", o.CancelURL, nil)
	if err != nil {
		return err
	}

	var o2 OrderOutput
	err = o.client.DoAndDecode(post, &o2)
	if err != nil {
		return errors.Wrap(err, "could not decode response")
	}

	if o2.RejectReason != "" {
		return errors.New(o2.RejectReason)
	}
	return nil
}
