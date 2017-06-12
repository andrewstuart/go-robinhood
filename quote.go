package robinhood

import (
	"encoding/json"
	"strings"
)

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
