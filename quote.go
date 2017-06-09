package main

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
)

type Quote struct {
	AdjustedPreviousClose       string   `json:"adjusted_previous_close"`
	AskPrice                    float64  `json:"ask_price"`
	AskSize                     int      `json:"ask_size"`
	BidPrice                    float64  `json:"bid_price"`
	BidSize                     int      `json:"bid_size"`
	LastExtendedHoursTradePrice *float64 `json:"last_extended_hours_trade_price"`
	LastTradePrice              float64  `json:"last_trade_price"`
	PreviousClose               float64  `json:"previous_close"`
	PreviousCloseDate           float64  `json:"previous_close_date"`
	Symbol                      string   `json:"symbol"`
	TradingHalted               bool     `json:"trading_halted"`
	UpdatedAt                   string   `json:"updated_at"`
}

func (c Client) GetQuote(stocks ...string) ([]Quote, error) {
	v := url.Values{
		"symbols": []string{strings.Join(stocks, ",")},
	}
	res, err := c.Get(epQuotes + v.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var r struct{ Results []Quote }

	err = json.NewDecoder(res.Body).Decode(&r)
	return r.Results, err
}
