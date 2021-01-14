package robinhood

import (
	"context"
	"errors"

	"fmt"
)

// CryptoCurrencyPair represent all availabe crypto currencies and whether they are tradeable or not
type CryptoCurrencyPair struct {
	CyrptoAssetCurrency    AssetCurrency `json:"asset_currency"`
	ID                     string        `json:"id"`
	MaxOrderSize           float64       `json:"max_order_size,string"`
	MinOrderPriceIncrement float64       `json:"min_order_price_increment,string"`
	MinOrderSize           float64       `json:"min_order_size,string"`
	Name                   string        `json:"name"`
	CrytoQuoteCurrency     QuoteCurrency `json:"quote_currency"`
	Symbol                 string        `json:"symbol"`
	Tradability            string        `json:"tradability"`
}

// QuoteCurrency holds info about currency you can use to buy the cyrpto currency
type QuoteCurrency struct {
	Code      string  `json:"code"`
	ID        string  `json:"id"`
	Increment float64 `json:"increment,string"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
}

// AssetCurrency has code and id of cryptocurrency
type AssetCurrency struct {
	BrandColor string  `json:"brand_color"`
	Code       string  `json:"code"`
	ID         string  `json:"id"`
	Increment  float64 `json:"increment,string"`
	Name       string  `json:"name"`
}

// GetCryptoCurrencyPairs will give which crypto currencies are tradeable and corresponding ids
func (c *Client) GetCryptoCurrencyPairs(ctx context.Context) ([]CryptoCurrencyPair, error) {
	var r struct{ Results []CryptoCurrencyPair }
	err := c.GetAndDecode(ctx, EPCryptoCurrencyPairs, &r)
	return r.Results, err
}

// GetCryptoInstrument will take standard crypto symbol and return usable information
// to place the order
func (c *Client) GetCryptoInstrument(ctx context.Context, symbol string) (*CryptoCurrencyPair, error) {
	allPairs, err := c.GetCryptoCurrencyPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("call failed with error: %v", err.Error())
	}

	for _, pair := range allPairs {
		if pair.CyrptoAssetCurrency.Code == symbol {
			return &pair, nil
		}
	}

	return nil, errors.New("could not find given symbol")
}
