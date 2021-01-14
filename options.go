package robinhood

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
)

const dateFormat = "2006-01-02"

// Date is a specific json time format for dates only
type Date struct {
	time.Time
}

// NewDate returns a new Date in the local time zone
func NewDate(y, m, d int) Date {
	return Date{time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)}
}

// NewZonedDate returns a date with a zone.
func NewZonedDate(y, m, d int, z *time.Location) Date {
	return Date{time.Date(y, time.Month(m), d, 0, 0, 0, 0, z)}
}

func (d Date) String() string {
	return d.Format(dateFormat)
}

// MarshalJSON implements json.Marshaler
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
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
func (c *Client) GetOptionChains(ctx context.Context, is ...*Instrument) ([]*OptionChain, error) {
	s := []string{}
	for _, inst := range is {
		s = append(s, inst.ID)
	}

	var res struct{ Results []*OptionChain }

	err := c.GetAndDecode(ctx, EPOptions+"chains/?equity_instrument_ids="+strings.Join(s, ","), &res)
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

type Pager struct {
	Next, Previous string
}

func (p Pager) HasMore() bool {
	return p.Next != ""
}

func (p *Pager) GetNext(ctx context.Context, c *Client, out interface{}) error {
	if p.Next == "" {
		return io.EOF
	}

	return c.GetAndDecode(ctx, p.Next, out)
}

// GetInstrument returns a list of option-typed instruments given a list of
// expiration dates for a given trade type. The request will continue until the
// provided context is cancelled. This is done to mimic the way the web UI
// fetches many, many options instruments repeatedly, since I haven't yet
// figured out how/when they decide to stop.
func (o *OptionChain) GetInstrument(ctx context.Context, tradeType string, date Date) ([]*OptionInstrument, error) {
	u := fmt.Sprintf(
		"%sinstruments/?chain_id=%s&expiration_dates=%s&state=active&tradability=tradable&type=%s",
		EPOptions,
		o.ID,
		date,
		tradeType,
	)

	var rs []*OptionInstrument
	var out struct {
		Results []*OptionInstrument
		Pager
	}
	err := o.c.GetAndDecode(ctx, u, &out)
	if err != nil {
		return nil, err
	}

	rs = append(rs, out.Results...)

	for out.HasMore() {
		select {
		case <-ctx.Done():
			return rs, ctx.Err()
		default:
		}

		err := out.GetNext(ctx, o.c, &out)
		if err != nil {
			return rs, err
		}
		rs = append(rs, out.Results...)
	}
	return rs, nil
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
	ExpirationDate Date     `json:"expiration_date"`
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

// MarketData is the current pricing data and greeks for a given option at a
// given time.
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

// OIsForDate filters OptionInstruments for expiration date.
func OIsForDate(os []*OptionInstrument, d Date) []*OptionInstrument {
	out := make([]*OptionInstrument, 0, len(os)/6)
	for i := range os {
		if os[i].ExpirationDate.Time.Equal(d.Time) {
			out = append(out, os[i])
		}
	}
	return out
}

// MarketData returns market data for all the listed Option instruments
func (c *Client) MarketData(ctx context.Context, opts ...*OptionInstrument) ([]*MarketData, error) {
	is := make([]string, len(opts))

	for i, o := range opts {
		is[i] = o.URL
	}

	u, err := url.Parse(EPOptionQuote)
	if err != nil {
		return nil, shameWrap(err, "couldn't parse URL const EPOptionQuote")
	}

	// Number of instruments to request at once
	num := 30
	// Number of requests this will require to be made
	n := len(is) / num
	if len(is)%num != 0 {
		n++
	}

	rs := []*MarketData{}

	for i := 0; i < n; i++ {
		end := (i+1)*num + 1
		if end > len(is) {
			end = len(is)
		}

		q := url.Values{"instruments": []string{strings.Join(is[i*num:end], ",")}}

		u.RawQuery = q.Encode()

		var r struct{ Results []*MarketData }
		if e := c.GetAndDecode(ctx, u.String(), &r); err != nil {
			err = multierror.Append(err, e)
			continue
		}
		for _, res := range r.Results {
			if res != nil {
				rs = append(rs, res)
			}
		}
	}

	return rs, err
}
