package robinhood

type Position struct {
	Account                 string `json:"account"`
	AverageBuyPrice         string `json:"average_buy_price"`
	CreatedAt               string `json:"created_at"`
	Instrument              string `json:"instrument"`
	IntradayAverageBuyPrice string `json:"intraday_average_buy_price"`
	IntradayQuantity        string `json:"intraday_quantity"`
	Quantity                string `json:"quantity"`
	SharesHeldForBuys       string `json:"shares_held_for_buys"`
	SharesHeldForSells      string `json:"shares_held_for_sells"`
	UpdatedAt               string `json:"updated_at"`
	URL                     string `json:"url"`
}

func (c Client) GetPositions(a Account) ([]Position, error) {
	var r struct{ Response []Position }
	res, err := c.c.Get(a.Positions)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
}
