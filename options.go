package robinhood

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetOptionChains(is ...*Instrument) ([]*OptionChain, error) {
	s := []string{}
	for _, inst := range is {
		s = append(s, inst.ID)
	}

	var res struct{ Results []*OptionChain }

	err := c.GetAndDecode(EPOptions+"chains/?equity_instrument_ids="+strings.Join(s, ","), &res)
	if err != nil {
		return nil, err
	}

	for i := range res.Results {
		res.Results[i].c = c
	}

	return res.Results, nil
}

type OptionChain struct {
	CanOpenPosition        bool                    `json:"can_open_position"`
	CashComponent          interface{}             `json:"cash_component"`
	ExpirationDates        []string                `json:"expiration_dates"`
	ID                     string                  `json:"id"`
	MinTicks               MinTicks                `json:"min_ticks"`
	Symbol                 string                  `json:"symbol"`
	TradeValueMultiplier   float64                 `json:"trade_value_multiplier,string"`
	UnderlyingInstrumentss []UnderlyingInstruments `json:"underlying_instruments"`

	c *Client
}

// type OptionTradeType string

const (
	Call = "call"
	Put  = "put"
)

func (o *OptionChain) Quote(tradeType string, dates ...string) ([]*OptionQuote, error) {
	u := EPOptions + fmt.Sprintf(
		"instruments/?chain_id=%s&expiration_dates=%s&state=active&tradability=tradable&type=%s",
		o.ID,
		strings.Join(dates, ","),
		tradeType,
	)

	var out struct{ Results []*OptionQuote }
	err := o.c.GetAndDecode(u, &out)
	if err != nil {
		return nil, err
	}
	return out.Results, nil
}

type MinTicks struct {
	AboveTick   float64 `json:"above_tick,string"`
	BelowTick   float64 `json:"below_tick,string"`
	CutoffPrice float64 `json:"cutoff_price,string"`
}

type UnderlyingInstruments struct {
	ID         string `json:"id"`
	Instrument string `json:"instrument"`
	Quantity   int    `json:"quantity"`
}

type OptionQuote struct {
	ChainID        string   `json:"chain_id"`
	ChainSymbol    string   `json:"chain_symbol"`
	CreatedAt      string   `json:"created_at"`
	ExpirationDate string   `json:"expiration_date"`
	ID             string   `json:"id"`
	IssueDate      string   `json:"issue_date"`
	MinTicks       MinTicks `json:"min_ticks"`
	RHSTradability string   `json:"rhs_tradability"`
	State          string   `json:"state"`
	StrikePrice    float64  `json:"strike_price,string"`
	Tradability    string   `json:"tradability"`
	Type           string   `json:"type"`
	UpdatedAt      string   `json:"updated_at"`
	URL            string   `json:"url"`
}

func (o OptionQuote) OrderURL() string {
	return o.URL
}

func (o OptionQuote) OrderSymbol() string {
	return o.ChainSymbol
}

type OptionLeg struct {
	PositionEffect string    `json:"position_effect"`
	Side           OrderSide `json:"side"`
	RatioQuantity  int       `json:"ratio_quantity"`
	Option         string    `json:"option"`
}

// OptionDirection is a type for whether an option order is opening or closing
// an option position
type OptionDirection int

//go:generate stringer -type OptionDirection
// The two directions
const (
	Debit OptionDirection = iota
	Credit
)

type OptionsOrderOpts struct {
	Quantity    float64
	Price       float64
	Direction   OptionDirection
	TimeInForce TimeInForce
	Type        OrderType
	Side        OrderSide
}

type optionBody struct {
	Account     string      `json:"account"`
	Direction   string      `json:"direction"`
	TimeInForce string      `json:"time_in_force"`
	Legs        []OptionLeg `json:"legs"`
	Type        string      `json:"type"`
	Quantity    float64     `json:"quantity,string"`
	ChainSymbol string      `json:"chain_symbol"`
	ChainID     string      `json:"chain_id"`
	Price       float64     `json:"price,string,omitempty"`
}

func (c *Client) OrderOptions(q *OptionQuote, o OptionsOrderOpts) (json.RawMessage, error) {
	b := optionBody{
		Account:     c.Account.URL,
		Direction:   strings.ToLower(o.Direction.String()),
		TimeInForce: strings.ToLower(o.TimeInForce.String()),
		ChainSymbol: q.ChainSymbol,
		ChainID:     q.ChainID,
		Legs: []OptionLeg{{
			Option:         q.URL,
			RatioQuantity:  1,
			Side:           o.Side,
			PositionEffect: "open",
		}},
		Type:     strings.ToLower(o.Type.String()),
		Quantity: o.Quantity,
		Price:    o.Price,
	}

	if o.Side != Buy {
		b.Legs[0].PositionEffect = "close"
	}

	bs, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", EPOptions+"orders/", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var out json.RawMessage
	err = c.DoAndDecode(req, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetOptionsOrders() (json.RawMessage, error) {
	var o json.RawMessage
	err := c.GetAndDecode(EPOptions+"orders/", &o)
	if err != nil {
		return nil, err
	}

	return o, nil

}
