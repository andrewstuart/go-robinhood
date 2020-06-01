package robinhood

import (
	"strings"
)

// A Quote is a representation of the data returned by the Robinhood API for
// current stock quotes
type Quote struct {
	AdjustedPreviousClose       float64 `json:"adjusted_previous_close,string"`
	AskPrice                    float64 `json:"ask_price,string"`
	AskSize                     int     `json:"ask_size"`
	BidPrice                    float64 `json:"bid_price,string"`
	BidSize                     int     `json:"bid_size"`
	LastExtendedHoursTradePrice float64 `json:"last_extended_hours_trade_price,string"`
	LastTradePrice              float64 `json:"last_trade_price,string"`
	PreviousClose               float64 `json:"previous_close,string"`
	PreviousCloseDate           string  `json:"previous_close_date"`
	Symbol                      string  `json:"symbol"`
	TradingHalted               bool    `json:"trading_halted"`
	UpdatedAt                   string  `json:"updated_at"`
}

// GetQuote returns all the latest stock quotes for the list of stocks provided
func (c *Client) GetQuote(stocks ...string) ([]Quote, error) {
	url := EPQuotes + "?symbols=" + strings.Join(stocks, ",")
	var r struct{ Results []Quote }
	err := c.GetAndDecode(url, &r)
	return r.Results, err
}

// Price returns the proper stock price even after hours
func (q Quote) Price() float64 {
	if IsRegularTradingTime() {
		return q.LastTradePrice
	}

	if q.LastExtendedHoursTradePrice != 0 {
		return q.LastExtendedHoursTradePrice
	} else {
		return q.AdjustedPreviousClose
	}
}
