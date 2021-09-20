package roho

import (
	"context"
	"net/url"
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

type OptionPostion struct {
	Chain                    string        `json:"chain"`
	AverageOpenPrice         string        `json:"average_open_price"`
	Symbol                   string        `json:"symbol"`
	Quantity                 string        `json:"quantity"`
	Direction                string        `json:"direction"`
	IntradayDirection        string        `json:"intraday_direction"`
	TradeValueMultiplier     string        `json:"trade_value_multiplier"`
	Account                  string        `json:"account"`
	Strategy                 string        `json:"strategy"`
	Legs                     []LegPosition `json:"legs"`
	IntradayQuantity         string        `json:"intraday_quantity"`
	UpdatedAt                string        `json:"updated_at"`
	Id                       string        `json:"id"`
	IntradayAverageOpenPrice string        `json:"intraday_average_open_price"`
	CreatedAt                string        `json:"created_at"`
}

type LegPosition struct {
	Id             string `json:"id"`
	Position       string `json:"position"`
	PositionType   string `json:"position_type"`
	Option         string `json:"option"`
	RatioQuantity  string `json:"ratio_quantity"`
	ExpirationDate string `json:"expiration_date"`
	StrikePrice    string `json:"strike_price"`
	OptionType     string `json:"option_type"`
}

type Unknown interface{}

// GetPositions returns all the positions associated with an account.
func (c *Client) OptionPositions(ctx context.Context) ([]OptionPostion, error) {
	return c.OptionPositionsParams(ctx, PositionParams{NonZero: true})
}

// GetPositions returns all the positions associated with an account.
func (c *Client) Positions(ctx context.Context) ([]Position, error) {
	return c.PositionsParams(ctx, PositionParams{NonZero: true})
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
func (c *Client) PositionsParams(ctx context.Context, p PositionParams) ([]Position, error) {
	u, err := url.Parse(baseURL("positions"))
	if err != nil {
		return nil, err
	}
	u.RawQuery = p.encode()

	var r struct{ Results []Position }
	return r.Results, c.get(ctx, u.String(), &r)
}

// GetPositionsParams returns all the positions associated with a count, but
// passes the encoded PositionsParams object along to the RobinHood API as part
// of the query string.
func (c *Client) OptionPositionsParams(ctx context.Context, p PositionParams) ([]OptionPostion, error) {
	u, err := url.Parse(baseURL("options") + "aggregate_positions/")
	if err != nil {
		return nil, err
	}
	u.RawQuery = p.encode()

	var r struct{ Results []OptionPostion }
	return r.Results, c.get(ctx, u.String(), &r)
}
