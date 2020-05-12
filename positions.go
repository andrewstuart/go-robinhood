package robinhood

import "net/url"

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

// GetPositions returns all the positionss associated with a client's
// credentials and accounts
func (c *Client) GetPositions() ([]Position, error) {
	var p struct{ Results []Position }
	err := c.GetAndDecode(EPPositions, &p)
	return p.Results, err
}

// PositionParams encapsulates parameters known to the RobinHood positions API
// endpoint.
type PositionParams struct {
	NonZero bool
}

// Encode returns the query string associated with the requested parameters
func (p PositionParams) encode() string {
	v := url.Values{}
	if p.NonZero {
		v.Set("nonzero", "true")
	}
	return v.Encode()
}

// GetPositionsParams returns all the positions associated with a count, but
// passes the encoded PositionsParams object along to the RobinHood API as part
// of the query string.
func (c *Client) GetPositionsParams(a Account, p PositionParams) ([]Position, error) {
	u, err := url.Parse(a.Positions)
	if err != nil {
		return nil, err
	}
	u.RawQuery = p.encode()

	var r struct{ Results []Position }
	return r.Results, c.GetAndDecode(u.String(), &r)
}
