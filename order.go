package robinhood

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// OrderSide is which side of the trade an order is on
type OrderSide int

// MarshalJSON implements json.Marshaler
func (o OrderSide) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strings.ToLower(o.String()) + "\""), nil
}

//go:generate stringer -type OrderSide
// Buy/Sell
const (
	Sell OrderSide = iota + 1
	Buy
)

// OrderType represents a Limit or Market order
type OrderType int

// MarshalJSON implements json.Marshaler
func (o OrderType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", strings.ToLower(o.String()))), nil
}

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
	Account       string    `json:"account,omitempty"`
	Instrument    string    `json:"instrument,omitempty"`
	Symbol        string    `json:"symbol,omitempty"`
	Type          string    `json:"type,omitempty"`
	TimeInForce   string    `json:"time_in_force,omitempty"`
	Trigger       string    `json:"trigger,omitempty"`
	Price         float64   `json:"price,omitempty"`
	StopPrice     float64   `json:"stop_price,omitempty"`
	Quantity      uint64    `json:"quantity,omitempty"`
	Side          OrderSide `json:"side,omitempty"`
	ExtendedHours bool      `json:"extended_hours,omitempty"`

	OverrideDayTradeChecks bool `json:"override_day_trade_checks,omitempty"`
	OverrideDtbpChecks     bool `json:"override_dtbp_checks,omitempty"`
}

// Order places an order for a given instrument. Cancellation of the given
// context cancels only the _http request_ and not any orders that may have
// been created regardless of the cancellation.
func (c *Client) Order(ctx context.Context, i *Instrument, o OrderOpts) (*OrderOutput, error) {
	a := apiOrder{
		Account:       c.Account.URL,
		Instrument:    i.URL,
		Symbol:        i.Symbol,
		Type:          strings.ToLower(o.Type.String()),
		TimeInForce:   strings.ToLower(o.TimeInForce.String()),
		Quantity:      o.Quantity,
		Side:          o.Side,
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

	post, err := http.NewRequest("POST", EPOrders, bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("error creating POST http.Request: %w", err)
	}

	post.Header.Add("Content-Type", "application/json")

	out := OrderOutput{}
	err = c.DoAndDecode(ctx, post, &out)
	if err != nil {
		return &out, err
	}

	out.client = c
	return &out, nil
}

// OrderOutput is the response from the Order api
type OrderOutput struct {
	Meta
	Account                string        `json:"account"`
	AveragePrice           float64       `json:"average_price,string"`
	CancelURL              string        `json:"cancel"`
	CreatedAt              string        `json:"created_at"`
	CumulativeQuantity     string        `json:"cumulative_quantity"`
	Executions             []interface{} `json:"executions"`
	ExtendedHours          bool          `json:"extended_hours"`
	Fees                   string        `json:"fees"`
	ID                     string        `json:"id"`
	Instrument             string        `json:"instrument"`
	LastTransactionAt      string        `json:"last_transaction_at"`
	OverrideDayTradeChecks bool          `json:"override_day_trade_checks"`
	OverrideDtbpChecks     bool          `json:"override_dtbp_checks"`
	Position               string        `json:"position"`
	Price                  float64       `json:"price,string"`
	Quantity               string        `json:"quantity"`
	RejectReason           string        `json:"reject_reason"`
	Side                   string        `json:"side"`
	State                  string        `json:"state"`
	StopPrice              float64       `json:"stop_price,string"`
	TimeInForce            string        `json:"time_in_force"`
	Trigger                string        `json:"trigger"`
	Type                   string        `json:"type"`

	client *Client
}

// Update returns any errors and updates the item with any recent changes.
func (o *OrderOutput) Update(ctx context.Context) error {
	return o.client.GetAndDecode(ctx, o.URL, o)
}

// Cancel attempts to cancel an odrer
func (o OrderOutput) Cancel(ctx context.Context) error {
	post, err := http.NewRequest("POST", o.CancelURL, nil)
	if err != nil {
		return err
	}

	var o2 OrderOutput
	err = o.client.DoAndDecode(ctx, post, &o2)
	if err != nil {
		return errors.Wrap(err, "could not decode response")
	}

	if o2.RejectReason != "" {
		return errors.New(o2.RejectReason)
	}
	return nil
}

// RecentOrders returns any recent orders made by this client.
func (c *Client) RecentOrders(ctx context.Context) ([]OrderOutput, error) {
	var o struct {
		Results []OrderOutput
	}
	err := c.GetAndDecode(ctx, EPOrders, &o)
	if err != nil {
		return o.Results, err
	}

	for i := range o.Results {
		o.Results[i].client = c
	}

	return o.Results, nil
}

// AllOrders returns all orders made by this client.
func (c *Client) AllOrders(ctx context.Context) ([]OrderOutput, error) {
	var o struct {
		Results []OrderOutput
	}

	urls := []string{EPOrders, EPOptionOrders}
	for _, url := range urls {
		for {
			select {
			case <-ctx.Done():
				return o.Results, ctx.Err()
			default:
			}
	
			var tmp struct {
				Results []OrderOutput
				Next    string
			}
			err := c.GetAndDecode(ctx, url, &tmp)
	
			if err != nil {
				return o.Results, err
			}
	
			url = tmp.Next
			o.Results = append(o.Results, tmp.Results...)
	
			if url == "" {
				break
			}
		}
	}

	return o.Results, nil
}
