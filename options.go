package robinhood

import (
	"fmt"
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
