package robinhood

import (
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const dateFormat = "2006-01-02"

// Date is a specific json time format for dates only
type Date struct {
	time.Time
}

// MarshalJSON implements json.Marshaler
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", d.Format(dateFormat))), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Date) UnmarshalJSON(bs []byte) error {
	t, err := time.Parse(dateFormat, strings.Trim(strings.TrimSpace(string(bs)), "\""))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// GetOptionChains returns options for the given instruments
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

// OptionChain represents the data the RobinHood API holds behind options chains
type OptionChain struct {
	CanOpenPosition       bool                   `json:"can_open_position"`
	CashComponent         interface{}            `json:"cash_component"`
	ExpirationDates       []string               `json:"expiration_dates"`
	ID                    string                 `json:"id"`
	MinTicks              MinTicks               `json:"min_ticks"`
	Symbol                string                 `json:"symbol"`
	TradeValueMultiplier  float64                `json:"trade_value_multiplier,string"`
	UnderlyingInstruments []UnderlyingInstrument `json:"underlying_instruments"`

	c *Client
}

// GetInstruments returns a list of option-typed instruments given a list of
// expiration dates for a given trade type.
func (o *OptionChain) GetInstruments(tradeType string, dates ...string) ([]*OptionInstrument, error) {
	u := fmt.Sprintf(
		"%sinstruments/?chain_id=%s&expiration_dates=%s&state=active&tradability=tradable&type=%s",
		EPOptions,
		o.ID,
		strings.Join(dates, ","),
		tradeType,
	)

	var out struct{ Results []*OptionInstrument }
	err := o.c.GetAndDecode(u, &out)
	if err != nil {
		return nil, err
	}
	return out.Results, nil
}

// MinTicks probably is important.
type MinTicks struct {
	AboveTick   float64 `json:"above_tick,string"`
	BelowTick   float64 `json:"below_tick,string"`
	CutoffPrice float64 `json:"cutoff_price,string"`
}

// UnderlyingInstrument is the type that represents a link from an option back
// to its standard financial instrument (stock)
type UnderlyingInstrument struct {
	ID         string `json:"id"`
	Instrument string `json:"instrument"`
	Quantity   int    `json:"quantity"`
}

// An OptionInstrument can have a quote
type OptionInstrument struct {
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

	c *Client
}

type MarketData struct {
	AdjustedMarkPrice   float64 `json:"adjusted_mark_price,string"`
	AskPrice            float64 `json:"ask_price,string"`
	AskSize             int     `json:"ask_size"`
	BidPrice            float64 `json:"bid_price,string"`
	BidSize             int     `json:"bid_size"`
	BreakEvenPrice      float64 `json:"break_even_price,string"`
	ChanceOfProfitLong  float64 `json:"chance_of_profit_long,string"`
	ChanceOfProfitShort float64 `json:"chance_of_profit_short,string"`
	Delta               float64 `json:"delta,string"`
	Gamma               float64 `json:"gamma,string"`
	HighPrice           float64 `json:"high_price,string"`
	ImpliedVolatility   string  `json:"implied_volatility"`
	Instrument          string  `json:"instrument"`
	LastTradePrice      float64 `json:"last_trade_price,string"`
	LastTradeSize       int     `json:"last_trade_size"`
	LowPrice            float64 `json:"low_price,string"`
	MarkPrice           float64 `json:"mark_price,string"`
	OpenInterest        int     `json:"open_interest"`
	PreviousCloseDate   Date    `json:"previous_close_date"`
	PreviousClosePrice  float64 `json:"previous_close_price,string"`
	Rho                 string  `json:"rho"`
	Theta               string  `json:"theta"`
	Vega                string  `json:"vega"`
	Volume              int     `json:"volume"`
}

func (c *Client) MarketData(os ...*OptionInstrument) ([]MarketData, error) {
	is := make([]string, len(os))

	for i, o := range os {
		is[i] = o.URL
	}

	var r struct{ Results []MarketData }
	spew.Dump(EPOptionQuote + "?instruments=" + strings.Join(is, ","))
	err := c.GetAndDecode(EPOptionQuote+"?instruments="+strings.Join(is, ","), &r)
	if err != nil {
		return nil, err
	}

	return r.Results, nil

}
