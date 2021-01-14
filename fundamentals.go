package robinhood

import (
	"context"
	"strings"
)

type Fundamental struct {
	Open          float64 `json:"open,string"`
	High          float64 `json:"high,string"`
	Low           float64 `json:"low,string"`
	Volume        float64 `json:"volume,string"`
	AverageVolume float64 `json:"average_volume,string"`
	High52Weeks   float64 `json:"high_52_weeks,string"`
	DividendYield float64 `json:"dividend_yield,string"`
	Low52Weeks    float64 `json:"low_52_weeks,string"`
	MarketCap     float64 `json:"market_cap,string"`
	PERatio       float64 `json:"pe_ratio,string"`
	Description   string  `json:"description"`
	Instrument    string  `json:"instrument"`
}

// GetFundamentals returns fundamental data for the list of stocks provided.
func (c *Client) GetFundamentals(ctx context.Context, stocks ...string) ([]Fundamental, error) {
	url := EPFundamentals + "?symbols=" + strings.Join(stocks, ",")
	var r struct{ Results []Fundamental }
	err := c.GetAndDecode(ctx, url, &r)
	return r.Results, err
}
