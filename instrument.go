package robinhood

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Instrument is a type to represent the "instrument" API type in the
// unofficial robinhood API.
type Instrument struct {
	BloombergUnique    string      `json:"bloomberg_unique"`
	Country            string      `json:"country"`
	DayTradeRatio      string      `json:"day_trade_ratio"`
	Fundamentals       string      `json:"fundamentals"`
	ID                 string      `json:"id"`
	ListDate           string      `json:"list_date"`
	MaintenanceRatio   string      `json:"maintenance_ratio"`
	MarginInitialRatio string      `json:"margin_initial_ratio"`
	Market             string      `json:"market"`
	MinTickSize        interface{} `json:"min_tick_size"`
	Name               string      `json:"name"`
	Quote              string      `json:"quote"`
	SimpleName         interface{} `json:"simple_name"`
	Splits             string      `json:"splits"`
	State              string      `json:"state"`
	Symbol             string      `json:"symbol"`
	Tradeable          bool        `json:"tradeable"`
	URL                string      `json:"url"`

	c Client
}

// GetInstrument returns an Instrument given a URL
func (c Client) GetInstrument(instURL string) (*Instrument, error) {
	var i Instrument
	err := c.GetAndDecode(instURL, &i)
	if err != nil {
		return nil, err
	}
	return &i, err
}

// GetInstrumentForSymbol returns an Instrument given a ticker symbol
func (c Client) GetInstrumentForSymbol(sym string) (*Instrument, error) {
	var i Instrument
	err := c.GetAndDecode(epInstruments+"?symbol="+sym, &i)
	return &i, err
}

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
	Account       string  `json:"account"`
	Instrument    string  `json:"instrument"`
	Symbol        string  `json:"symbol"`
	Type          string  `json:"type"`
	TimeInForce   string  `json:"time_in_force"`
	Trigger       string  `json:"trigger"`
	Price         float64 `json:"price"`
	StopPrice     float64 `json:"stop_price"`
	Quantity      uint64  `json:"quantity"`
	Side          string  `json:"side"`
	ExtendedHours bool    `json:"extended_hours"`

	OverrideDayTradeChecks bool `json:"override_day_trade_checks"`
	OverrideDtbpChecks     bool `json:"override_dtbp_checks"`
}

// Order places an order for a given instrument
func (c Client) Order(i Instrument, o OrderOpts) (*OrderOutput, error) {
	a := apiOrder{
		Account:       c.Account.URL,
		Instrument:    i.URL,
		Symbol:        i.Symbol,
		Type:          o.Type.String(),
		TimeInForce:   o.TimeInForce.String(),
		Quantity:      o.Quantity,
		Side:          o.Side.String(),
		ExtendedHours: o.ExtendedHours,
	}

	if o.Stop {
		a.StopPrice = o.Price
		a.Trigger = "stop"
	} else {
		a.Price = o.Price
		a.Trigger = "immediate"
	}

	bs, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	post, err := http.NewRequest("POST", epBase+"/orders/", bytes.NewReader(bs))
	post.Header.Add("Content-Type", "application/json")

	var out OrderOutput
	err = c.DoAndDecode(post, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
