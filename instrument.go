package robinhood

import (
	"context"
	"fmt"
)

// Instrument is a type to represent the "instrument" API type in the
// unofficial robinhood API.
type Instrument struct {
	BloombergUnique       string      `json:"bloomberg_unique"`
	Country               string      `json:"country"`
	DayTradeRatio         string      `json:"day_trade_ratio"`
	DefaultCollarFraction string      `json:"default_collar_fraction"`
	FractionalTradability string      `json:"fractional_tradability"`
	Fundamentals          string      `json:"fundamentals"`
	ID                    string      `json:"id"`
	ListDate              string      `json:"list_date"`
	MaintenanceRatio      string      `json:"maintenance_ratio"`
	MarginInitialRatio    string      `json:"margin_initial_ratio"`
	Market                string      `json:"market"`
	MinTickSize           interface{} `json:"min_tick_size"`
	Name                  string      `json:"name"`
	Quote                 string      `json:"quote"`
	RhsTradability        string      `json:"rhs_tradability"`
	SimpleName            interface{} `json:"simple_name"`
	Splits                string      `json:"splits"`
	State                 string      `json:"state"`
	Symbol                string      `json:"symbol"`
	Tradeable             bool        `json:"tradeable"`
	Tradability           string      `json:"tradability"`
	TradableChainID       string      `json:"tradable_chain_id"`
	Type                  string      `json:"type"`
	URL                   string      `json:"url"`
}

func (i Instrument) OrderURL() string {
	return i.URL
}

func (i Instrument) OrderSymbol() string {
	return i.Symbol
}

// GetInstrument returns an Instrument given a URL
func (c *Client) GetInstrument(ctx context.Context, instURL string) (*Instrument, error) {
	var i Instrument
	err := c.GetAndDecode(ctx, instURL, &i)
	if err != nil {
		return nil, err
	}
	return &i, err
}

// GetInstrumentForSymbol returns an Instrument given a ticker symbol
func (c *Client) GetInstrumentForSymbol(ctx context.Context, sym string) (*Instrument, error) {
	var i struct {
		Results []Instrument
	}
	err := c.GetAndDecode(ctx, EPInstruments+"?symbol="+sym, &i)
	if err != nil {
		return nil, err
	}
	if len(i.Results) < 1 {
		return nil, fmt.Errorf("no results")
	}
	return &i.Results[0], err
}
