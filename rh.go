package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	base     = "https://api.robinhood.com/"
	login    = base + "api-token-auth/"
	accounts = base + "accounts/"
)

type Creds struct {
	Username, Password string
}

func (c Creds) Values() url.Values {
	return url.Values{
		"username": []string{c.Username},
		"password": []string{c.Password},
	}
}

type Client struct {
	Token   string
	Account *Account
	c       *http.Client
}

type Meta struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	URL                        string      `json:"url"`
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

type tokenRoundtripper string

func (t tokenRoundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", t))
	return http.DefaultClient.Do(req)
}

func (c *Client) GetAccounts() ([]Account, error) {
	var r struct{ Results []Account }
	res, err := c.c.Get(accounts)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&r)

	return r.Results, err
}

func Dial(c Creds) (*Client, error) {
	res, err := http.Post(login, "application/x-www-form-urlencoded", strings.NewReader(c.Values().Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var cli Client
	err = json.NewDecoder(res.Body).Decode(&cli)

	cli.c = &http.Client{
		Transport: tokenRoundtripper(cli.Token),
	}

	return &cli, err
}

func main() {
	creds := Creds{
		Username: "andrewstuart",
		Password: os.Getenv("ROBINHOOD_PASSWORD"),
	}

	c, err := Dial(creds)
	if err != nil {
		log.Fatal(err)
	}
	a, err := c.GetAccounts()
	if err != nil {
		log.Fatal(err)
	}

	pos, err := c.GetPositions(a[0])
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range pos {
		res, err := c.c.Get(p.Instrument)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		io.Copy(os.Stdout, res.Body)
		log.Println("Next")
	}
}
