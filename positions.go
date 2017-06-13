package robinhood

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

// GetPositions returns all the positions associated with an account.
func (c Client) GetPositions(a Account) ([]Position, error) {
	var r struct{ Results []Position }
	err := c.GetAndDecode(a.Positions, &r)
	return r.Results, err
}
