package robinhood

import (
	"encoding/json"
)

type Position struct {
	Meta
	Account                 string  `json:"account"`
	AverageBuyPrice         float64 `json:"average_buy_price,string"`
	Instrument              string  `json:"instrument"`
	IntradayAverageBuyPrice float64 `json:"intraday_average_buy_price,string"`
	IntradayQuantity        float64 `json:"intraday_quantity,string"`
	Quantity                float64 `json:"quantity,string"`
	SharesHeldForBuys       float64 `json:"shares_held_for_buys,string"`
	SharesHeldForSells      float64 `json:"shares_held_for_sells,string"`
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
