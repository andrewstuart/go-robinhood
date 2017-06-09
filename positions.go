package main

import "encoding/json"

type Position struct {
	Meta
	Account                 string `json:"account"`
	AverageBuyPrice         string `json:"average_buy_price"`
	Instrument              string `json:"instrument"`
	IntradayAverageBuyPrice string `json:"intraday_average_buy_price"`
	IntradayQuantity        string `json:"intraday_quantity"`
	Quantity                string `json:"quantity"`
	SharesHeldForBuys       string `json:"shares_held_for_buys"`
	SharesHeldForSells      string `json:"shares_held_for_sells"`
}

func (c Client) GetPositions(a Account) ([]Position, error) {
	res, err := c.Get(a.Positions)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r struct{ Results []Position }
	err = json.NewDecoder(res.Body).Decode(&r)

	return r.Results, err
}
