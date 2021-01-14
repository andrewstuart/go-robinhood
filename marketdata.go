package robinhood

import (
	"context"
	"fmt"
)

type EntryPrice struct {
	Amount        string
	Currency_code string
}

type PriceBookEntry struct {
	Side     string
	Price    EntryPrice
	Quantity float64
}

type PriceBookData struct {
	Asks []PriceBookEntry `json:"asks"`
	Bids []PriceBookEntry `json:"bids"`

	InstrumentID string `json:"instrument_id"`
	UpdatedAt    string `json:"updated_at"`
}

// Pricebook get the current snapshot of the pricebook data
func (c *Client) Pricebook(ctx context.Context, instrumentID string) (*PriceBookData, error) {
	var out PriceBookData
	err := c.GetAndDecode(ctx, fmt.Sprintf("%spricebook/snapshots/%s/", EPMarket, instrumentID), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
