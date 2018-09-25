package robinhood

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"astuart.co/clyde"
)

// Endpoints for the Robinhood API
const (
	EPBase         = "https://api.robinhood.com/"
	EPLogin        = EPBase + "api-token-auth/"
	EPAccounts     = EPBase + "accounts/"
	EPQuotes       = EPBase + "quotes/"
	EPPortfolios   = EPBase + "portfolios/"
	EPWatchlists   = EPBase + "watchlists/"
	EPInstruments  = EPBase + "instruments/"
	EPFundamentals = EPBase + "fundamentals/"
	EPOrders       = EPBase + "orders/"
	EPOptions      = EPBase + "options/"
)

// A Client is a helpful abstraction around some common metadata required for
// API operations.
type Client struct {
	Token   string
	Account *Account
	*http.Client
}

// Dial returns a client given a TokenGetter. TokenGetter implementations are
// available in this package, including a Cookie-based cache.
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

// GetAndDecode retrieves from the endpoint and unmarshals resulting json into
// the provided destination interface, which must be a pointer.
func (c Client) GetAndDecode(url string, dest interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	return c.DoAndDecode(req, dest)
}

// ErrorMap encapsulates the helpful error messages returned by the API server
type ErrorMap map[string]interface{}

func (e ErrorMap) Error() string {
	es := make([]string, 0, len(e))
	for k, v := range e {
		es = append(es, fmt.Sprintf("%s: %q", k, v))
	}
	return "Error returned from API: " + strings.Join(es, ", ")
}

// DoAndDecode provides useful abstractions around common errors and decoding
// issues.
func (c Client) DoAndDecode(req *http.Request, dest interface{}) error {
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		var e ErrorMap
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return fmt.Errorf("got response %q and could not decode error body", res.Status)
		}
		return e
	}

	return json.NewDecoder(res.Body).Decode(dest)
}

// Meta holds metadata common to many RobinHood types.
type Meta struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
}
