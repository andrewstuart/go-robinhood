package robinhood

import (
	"encoding/json"
	"net/http"
	"time"

	"astuart.co/clyde"
)

const (
	epBase         = "https://api.robinhood.com/"
	epLogin        = epBase + "api-token-auth/"
	epAccounts     = epBase + "accounts/"
	epQuotes       = epBase + "quotes/"
	epPortfolios   = epBase + "portfolios/"
	epWatchlists   = epBase + "watchlists/"
	epInstruments  = epBase + "instruments/"
	epFundamentals = epBase + "fundamentals/"
)

type Client struct {
	Token   string
	Account *Account
	*http.Client
}

func Dial(t TokenGetter) (*Client, error) {
	tkn, err := t.GetToken()
	if err != nil {
		return nil, err
	}

	c := &Client{
		Token:  tkn,
		Client: &http.Client{Transport: clyde.HeaderRoundTripper{"Authorization": "Token " + tkn}},
	}

	a, _ := c.GetAccounts()
	if len(a) > 0 {
		c.Account = &a[0]
	}

	return c, nil
}

func (c Client) GetAndDecode(url string, dest interface{}) error {
	res, err := c.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(dest)
}

func (c Client) DoAndDecode(req *http.Request, dest interface{}) error {
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(dest)
}

type Meta struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
}
