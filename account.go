package robinhood

// Account holds the basic account details relevant to the RobinHood API
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

// CashBalances reflect the amount of cash available
type CashBalances struct {
	Meta
	BuyingPower                float64 `json:"buying_power,string"`
	Cash                       float64 `json:"cash,string"`
	CashAvailableForWithdrawal float64 `json:"cash_available_for_withdrawal,string"`
	CashHeldForOrders          float64 `json:"cash_held_for_orders,string"`
	UnclearedDeposits          float64 `json:"uncleared_deposits,string"`
	UnsettledFunds             float64 `json:"unsettled_funds,string"`
}

// MarginBalances reflect the balance available in margin accounts
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

// GetAccounts returns all the accounts associated with a login/client.
func (c *Client) GetAccounts() ([]Account, error) {
	var r struct{ Results []Account }
	err := c.GetAndDecode(EPAccounts, &r)
	if err != nil {
		return nil, err
	}
	return r.Results, err
}

// CryptoAccount holds the basic account details relevant to robinhood API
type CryptoAccount struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	UserID string `json:"user_id"`
}

// GetCryptoAccounts will return associated cryto account
func (c *Client) GetCryptoAccounts() ([]CryptoAccount, error) {
	var r struct{ Results []CryptoAccount }
	err := c.GetAndDecode(EPCryptoAccount, &r)
	if err != nil {
		return nil, err
	}

	return r.Results, err
}
