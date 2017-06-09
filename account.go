package main

import "encoding/json"

type Account struct {
	Meta
	AccountNumber              string `json:"account_number"`
	BuyingPower                string `json:"buying_power"`
	Cash                       string `json:"cash"`
	CashAvailableForWithdrawal string `json:"cash_available_for_withdrawal"`
	CashBalances               `json:"cash_balances"`
	CashHeldForOrders          string      `json:"cash_held_for_orders"`
	Deactivated                bool        `json:"deactivated"`
	DepositHalted              bool        `json:"deposit_halted"`
	MarginBalances             interface{} `json:"margin_balances"`
	MaxAchEarlyAccessAmount    string      `json:"max_ach_early_access_amount"`
	OnlyPositionClosingTrades  bool        `json:"only_position_closing_trades"`
	Portfolio                  string      `json:"portfolio"`
	Positions                  string      `json:"positions"`
	Sma                        interface{} `json:"sma"`
	SmaHeldForOrders           interface{} `json:"sma_held_for_orders"`
	SweepEnabled               bool        `json:"sweep_enabled"`
	Type                       string      `json:"type"`
	UnclearedDeposits          string      `json:"uncleared_deposits"`
	UnsettledFunds             string      `json:"unsettled_funds"`
	UpdatedAt                  string      `json:"updated_at"`
	User                       string      `json:"user"`
	WithdrawalHalted           bool        `json:"withdrawal_halted"`
}

type CashBalances struct {
	Meta
	BuyingPower                string `json:"buying_power"`
	Cash                       string `json:"cash"`
	CashAvailableForWithdrawal string `json:"cash_available_for_withdrawal"`
	CashHeldForOrders          string `json:"cash_held_for_orders"`
	UnclearedDeposits          string `json:"uncleared_deposits"`
	UnsettledFunds             string `json:"unsettled_funds"`
}

func (c *Client) GetAccounts() ([]Account, error) {
	var r struct{ Results []Account }
	res, err := c.Get(epAccounts)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&r)

	return r.Results, err
}
