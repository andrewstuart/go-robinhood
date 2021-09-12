package robinhood

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Historical struct {
	Symbol   string             `json:"symbol"`
	Interval string             `json:"interval"`
	Bounds   string             `json:"bounds"`
	Span     string             `json:"span"`
	Records  []HistoricalRecord `json:"historicals"`
}

type HistoricalRecord struct {
	BeginsAt     time.Time `json:"begins_at"`
	OpenPrice    float64   `json:"open_price,string"`
	ClosePrice   float64   `json:"close_price,string"`
	HighPrice    float64   `json:"high_price,string"`
	LowPrice     float64   `json:"low_price,string"`
	Volume       int64     `json:"volume"`
	Session      string    `json:"session"`
	Interpolated bool      `json:"interpolated"`
}

// GetHistoricals returns historical data for the list of stocks provided.
func (c *Client) GetHistoricals(ctx context.Context, interval string, span string, stocks ...string) ([]Historical, error) {
	url := fmt.Sprintf("%s?interval=%s&span=%s&symbols=%s", EPHistoricals, interval, span, strings.Join(stocks, ","))
	var r struct{ Results []Historical }
	err := c.GetAndDecode(ctx, url, &r)
	return r.Results, err
}
