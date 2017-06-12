package main

import (
	"encoding/json"
	"io"
	"os"
)

type Portfolio struct {
	Account                                string  `json:"account"`
	AdjustedEquityPreviousClose            float64 `json:"adjusted_equity_previous_close,string"`
	Equity                                 float64 `json:"equity,string"`
	EquityPreviousClose                    float64 `json:"equity_previous_close,string"`
	ExcessMaintenance                      float64 `json:"excess_maintenance,string"`
	ExcessMaintenanceWithUnclearedDeposits float64 `json:"excess_maintenance_with_uncleared_deposits,string"`
	ExcessMargin                           float64 `json:"excess_margin,string"`
	ExcessMarginWithUnclearedDeposits      float64 `json:"excess_margin_with_uncleared_deposits,string"`
	ExtendedHoursEquity                    float64 `json:"extended_hours_equity,string"`
	ExtendedHoursMarketValue               float64 `json:"extended_hours_market_value,string"`
	LastCoreEquity                         float64 `json:"last_core_equity,string"`
	LastCoreMarketValue                    float64 `json:"last_core_market_value,string"`
	MarketValue                            float64 `json:"market_value,string"`
	StartDate                              string  `json:"start_date"`
	UnwithdrawableDeposits                 float64 `json:"unwithdrawable_deposits,string"`
	UnwithdrawableGrants                   float64 `json:"unwithdrawable_grants,string"`
	URL                                    string  `json:"url"`
	WithdrawableAmount                     float64 `json:"withdrawable_amount,string"`
}

func (c *Client) GetPortfolios(acc *Account) ([]Portfolio, error) {
	ep := epPortfolios
	if acc != nil {
		ep = acc.Portfolio
	}

	res, err := c.Get(ep)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if acc == nil {
		var p struct{ Results []Portfolio }
		err = json.NewDecoder(io.TeeReader(res.Body, os.Stderr)).Decode(&p)
		return p.Results, err
	}

	var p Portfolio
	err = json.NewDecoder(res.Body).Decode(&p)
	return []Portfolio{p}, err
}
