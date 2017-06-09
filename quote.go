package main

import (
	"encoding/json"
	"strings"
)

type Quote struct {
	AdjustedPreviousClose       string  `json:"adjusted_previous_close"`
	AskPrice                    string  `json:"ask_price"`
	AskSize                     int     `json:"ask_size"`
	BidPrice                    string  `json:"bid_price"`
	BidSize                     int     `json:"bid_size"`
	LastExtendedHoursTradePrice *string `json:"last_extended_hours_trade_price"`
	LastTradePrice              string  `json:"last_trade_price"`
	PreviousClose               string  `json:"previous_close"`
	PreviousCloseDate           string  `json:"previous_close_date"`
	Symbol                      string  `json:"symbol"`
	TradingHalted               bool    `json:"trading_halted"`
	UpdatedAt                   string  `json:"updated_at"`
}

func (c Client) GetQuote(stocks ...string) ([]Quote, error) {
	url := epQuotes + "?symbols=" + strings.Join(stocks, ",")
	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r struct{ Results []Quote }

	err = json.NewDecoder(res.Body).Decode(&r)
	return r.Results, err
}
