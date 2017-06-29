# robinhood
[![GoDoc](https://godoc.org/astuart.co/go-robinhood?status.svg)](https://godoc.org/astuart.co/go-robinhood)
--
    import "astuart.co/go-robinhood"


## Usage

```go
const (
	HrExtendedOpen    = 4
	HrRHExtendedOpen  = 9
	HrClose           = 12 + 4
	HrRHExtendedClose = 12 + 6
	HrExtendedClose   = 12 + 8

	MinExtendedOpen    = HrExtendedOpen * 60
	MinRHExtendedOpen  = HrRHExtendedOpen * 60
	MinOpen            = 9*60 + 30
	MinClose           = HrClose * 60
	MinRHExtendedClose = HrRHExtendedClose * 60
	MinExtendedClose   = HrExtendedClose * 60
)
```
Common constants for hours and minutes from midnight at which market events
occur.

#### func  IsExtendedTradingTime

```go
func IsExtendedTradingTime() bool
```
IsExtendedTradingTime returns whether or not extended hours equity will be
updated because extended-hours trades may still be allowed in the markets.

#### func  IsRegularTradingTime

```go
func IsRegularTradingTime() bool
```
IsRegularTradingTime returns whether or not the markets are currently open for
regular trading.

#### func  IsRobinhoodExtendedTradingTime

```go
func IsRobinhoodExtendedTradingTime() bool
```
IsRobinhoodExtendedTradingTime returns whether or not trades can still be placed
during the robinhood gold extended trading hours.

#### func  IsWeekDay

```go
func IsWeekDay(t time.Time) bool
```
IsWeekDay returns whether the given time is a regular weekday

#### func  MinuteOfDay

```go
func MinuteOfDay(t time.Time) int
```
MinuteOfDay returns the minute of the day for a given time.Time (hr * 60 + min).

#### func  NextMarketClose

```go
func NextMarketClose() time.Time
```
NextMarketClose returns the time of the next market close.

#### func  NextMarketExtendedClose

```go
func NextMarketExtendedClose() time.Time
```
NextMarketExtendedClose returns the time of the next extended market close, when
stock equity numbers will stop being updated until the next extended open.

#### func  NextMarketExtendedOpen

```go
func NextMarketExtendedOpen() time.Time
```
NextMarketExtendedOpen returns the time of the next extended opening time, when
stock equity may begin to fluctuate again.

#### func  NextMarketOpen

```go
func NextMarketOpen() time.Time
```
NextMarketOpen returns the time of the next opening bell, when regular trading
begins.

#### func  NextRobinhoodExtendedClose

```go
func NextRobinhoodExtendedClose() time.Time
```
NextRobinhoodExtendedClose returns the time of the next robinhood extended
closing time, when robinhood users must place their last extended-hours trade.

#### func  NextRobinhoodExtendedOpen

```go
func NextRobinhoodExtendedOpen() time.Time
```
NextRobinhoodExtendedOpen returns the time of the next robinhood extended
opening time, when robinhood users can make trades.

#### func  NextWeekday

```go
func NextWeekday() time.Time
```
NextWeekday returns the next weekday.

#### type Account

```go
type Account struct {
	Meta
	AccountNumber              string         `json:"account_number"`
	BuyingPower                float64        `json:"buying_power,string"`
	Cash                       float64        `json:"cash,string"`
	CashAvailableForWithdrawal float64        `json:"cash_available_for_withdrawal,string"`
	CashBalances               CashBalances   `json:"cash_balances"`
	CashHeldForOrders          float64        `json:"cash_held_for_orders,string"`
	Deactivated                bool           `json:"deactivated"`
	DepositHalted              bool           `json:"deposit_halted"`
	MarginBalances             MarginBalances `json:"margin_balances"`
	MaxAchEarlyAccessAmount    string         `json:"max_ach_early_access_amount"`
	OnlyPositionClosingTrades  bool           `json:"only_position_closing_trades"`
	Portfolio                  string         `json:"portfolio"`
	Positions                  string         `json:"positions"`
	Sma                        interface{}    `json:"sma"`
	SmaHeldForOrders           interface{}    `json:"sma_held_for_orders"`
	SweepEnabled               bool           `json:"sweep_enabled"`
	Type                       string         `json:"type"`
	UnclearedDeposits          float64        `json:"uncleared_deposits,string"`
	UnsettledFunds             float64        `json:"unsettled_funds,string"`
	User                       string         `json:"user"`
	WithdrawalHalted           bool           `json:"withdrawal_halted"`
}
```


#### type CashBalances

```go
type CashBalances struct {
	Meta
	BuyingPower                float64 `json:"buying_power,string"`
	Cash                       float64 `json:"cash,string"`
	CashAvailableForWithdrawal float64 `json:"cash_available_for_withdrawal,string"`
	CashHeldForOrders          float64 `json:"cash_held_for_orders,string"`
	UnclearedDeposits          float64 `json:"uncleared_deposits,string"`
	UnsettledFunds             float64 `json:"unsettled_funds,string"`
}
```


#### type Client

```go
type Client struct {
	Token   string
	Account *Account
	*http.Client
}
```


#### func  Dial

```go
func Dial(t TokenGetter) (*Client, error)
```

#### func (*Client) GetAccounts

```go
func (c *Client) GetAccounts() ([]Account, error)
```

#### func (Client) GetAndDecode

```go
func (c Client) GetAndDecode(url string, dest interface{}) error
```

#### func (Client) GetInstrument

```go
func (c Client) GetInstrument(instURL string) (*Instrument, error)
```

#### func (Client) GetInstrumentForSymbol

```go
func (c Client) GetInstrumentForSymbol(sym string) (*Instrument, error)
```

#### func (*Client) GetPortfolios

```go
func (c *Client) GetPortfolios() ([]Portfolio, error)
```
GetPortfolios returns all the portfolios associated with a client's credentials
and accounts

#### func (Client) GetPositions

```go
func (c Client) GetPositions(a Account) ([]Position, error)
```
GetPositions returns all the positions associated with an account.

#### func (Client) GetQuote

```go
func (c Client) GetQuote(stocks ...string) ([]Quote, error)
```
GetQuote returns all the latest stock quotes for the list of stocks provided

#### func (*Client) GetWatchlists

```go
func (c *Client) GetWatchlists() ([]Watchlist, error)
```
GetWatchlists retrieves the watchlists for a given set of credentials/accounts.

#### type Creds

```go
type Creds struct {
	Username, Password string
}
```


#### func (*Creds) GetToken

```go
func (c *Creds) GetToken() (string, error)
```

#### func (Creds) Values

```go
func (c Creds) Values() url.Values
```

#### type CredsCacher

```go
type CredsCacher struct {
	Creds TokenGetter
	Path  string
}
```

A CredsCacher takes user credentials and a file path. The token obtained from
the RobinHood API will be cached at the file path, and a new token will not be
obtained.

#### func (*CredsCacher) GetToken

```go
func (c *CredsCacher) GetToken() (string, error)
```
GetToken implements TokenGetter. It may fail if an error is encountered checking
the file path provided, or if the underlying creds return an error when
retrieving their token.

#### type Instrument

```go
type Instrument struct {
	BloombergUnique    string      `json:"bloomberg_unique"`
	Country            string      `json:"country"`
	DayTradeRatio      string      `json:"day_trade_ratio"`
	Fundamentals       string      `json:"fundamentals"`
	ID                 string      `json:"id"`
	ListDate           string      `json:"list_date"`
	MaintenanceRatio   string      `json:"maintenance_ratio"`
	MarginInitialRatio string      `json:"margin_initial_ratio"`
	Market             string      `json:"market"`
	MinTickSize        interface{} `json:"min_tick_size"`
	Name               string      `json:"name"`
	Quote              string      `json:"quote"`
	SimpleName         interface{} `json:"simple_name"`
	Splits             string      `json:"splits"`
	State              string      `json:"state"`
	Symbol             string      `json:"symbol"`
	Tradeable          bool        `json:"tradeable"`
	URL                string      `json:"url"`
}
```


#### type MarginBalances

```go
type MarginBalances struct {
	Meta
	Cash                              float64 `json:"cash,string"`
	CashAvailableForWithdrawal        float64 `json:"cash_available_for_withdrawal,string"`
	CashHeldForOrders                 float64 `json:"cash_held_for_orders,string"`
	DayTradeBuyingPower               float64 `json:"day_trade_buying_power,string"`
	DayTradeBuyingPowerHeldForOrders  float64 `json:"day_trade_buying_power_held_for_orders,string"`
	DayTradeRatio                     float64 `json:"day_trade_ratio,string"`
	MarginLimit                       float64 `json:"margin_limit,string"`
	MarkedPatternDayTraderDate        string  `json:"marked_pattern_day_trader_date"`
	OvernightBuyingPower              float64 `json:"overnight_buying_power,string"`
	OvernightBuyingPowerHeldForOrders float64 `json:"overnight_buying_power_held_for_orders,string"`
	OvernightRatio                    float64 `json:"overnight_ratio,string"`
	UnallocatedMarginCash             float64 `json:"unallocated_margin_cash,string"`
	UnclearedDeposits                 float64 `json:"uncleared_deposits,string"`
	UnsettledFunds                    float64 `json:"unsettled_funds,string"`
}
```


#### type Meta

```go
type Meta struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
}
```


#### type Portfolio

```go
type Portfolio struct {
	Account                                string  `json:"account"`
	AdjustedEquityPreviousClose            float64 `json:"adjusted_equity_previous_close,string"`
	Equity                                 float64 `json:"equity,string"`
	EquityPreviousClose                    float64 `json:"equity_previous_close,string"`
	ExcessMaintenance                      float64 `json:"excess_maintenance,string"`
	ExcessMaintenanceWithUnclearedDeposits float64 `json:"excess_maintenance_with_uncleared_deposits,string"`
	ExcessMargin                           float64 `json:"excess_margin,string"`
	ExcessMarginWithUnclearedDeposits      float64 `json:"excess_margin_with_uncleared_deposits,string"`
	ExtendedHoursEquity                    float64 `json:"extended_hours_equity,string"`
	ExtendedHoursMarketValue               float64 `json:"extended_hours_market_value,string"`
	LastCoreEquity                         float64 `json:"last_core_equity,string"`
	LastCoreMarketValue                    float64 `json:"last_core_market_value,string"`
	MarketValue                            float64 `json:"market_value,string"`
	StartDate                              string  `json:"start_date"`
	UnwithdrawableDeposits                 float64 `json:"unwithdrawable_deposits,string"`
	UnwithdrawableGrants                   float64 `json:"unwithdrawable_grants,string"`
	URL                                    string  `json:"url"`
	WithdrawableAmount                     float64 `json:"withdrawable_amount,string"`
}
```


#### type Position

```go
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
```


#### type Quote

```go
type Quote struct {
	AdjustedPreviousClose       float64 `json:"adjusted_previous_close,string"`
	AskPrice                    float64 `json:"ask_price,string"`
	AskSize                     int     `json:"ask_size"`
	BidPrice                    float64 `json:"bid_price,string"`
	BidSize                     int     `json:"bid_size"`
	LastExtendedHoursTradePrice float64 `json:"last_extended_hours_trade_price,string"`
	LastTradePrice              float64 `json:"last_trade_price,string"`
	PreviousClose               float64 `json:"previous_close,string"`
	PreviousCloseDate           string  `json:"previous_close_date"`
	Symbol                      string  `json:"symbol"`
	TradingHalted               bool    `json:"trading_halted"`
	UpdatedAt                   string  `json:"updated_at"`
}
```

A Quote is a representation of the data returned by the Robinhood API for
current stock quotes

#### type Token

```go
type Token string
```


#### func (*Token) GetToken

```go
func (t *Token) GetToken() (string, error)
```

#### type TokenGetter

```go
type TokenGetter interface {
	GetToken() (string, error)
}
```


#### type Watchlist

```go
type Watchlist struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	User string `json:"user"`

	Client *Client `json:",ignore"`
}
```

A Watchlist is a list of stock Instruments that an investor is tracking in his
Robinhood portfolio/app.

#### func (*Watchlist) GetInstruments

```go
func (w *Watchlist) GetInstruments() ([]Instrument, error)
```
GetInstruments returns the list of Instruments associated with a Watchlist.
