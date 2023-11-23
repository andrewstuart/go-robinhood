package robinhood

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// Endpoints for the Robinhood API
const (
	EPBase                = "https://api.robinhood.com/"
	EPCryptoBase          = "https://nummus.robinhood.com/"
	EPCryptoOrders        = EPCryptoBase + "orders/"
	EPCryptoAccount       = EPCryptoBase + "accounts/"
	EPCryptoCurrencyPairs = EPCryptoBase + "currency_pairs/"
	EPCryptoHoldings      = EPCryptoBase + "holdings/"
	EPCryptoPortfolio     = EPCryptoBase + "portfolios/"
	EPLogin               = EPBase + "oauth2/token/"
	EPAccounts            = EPBase + "accounts/"
	EPQuotes              = EPBase + "quotes/"
	EPPortfolios          = EPBase + "portfolios/"
	EPPositions           = EPBase + "positions/"
	EPWatchlists          = EPBase + "watchlists/"
	EPInstruments         = EPBase + "instruments/"
	EPFundamentals        = EPBase + "fundamentals/"
	EPOptionOrders        = EPBase + "options/orders/"
	EPOrders              = EPBase + "orders/"
	EPOptions             = EPBase + "options/"
	EPMarket              = EPBase + "marketdata/"
	EPOptionQuote         = EPMarket + "options/"
)

// A Client is a helpful abstraction around some common metadata required for
// API operations.
type Client struct {
	Token         string
	Account       *Account
	CryptoAccount *CryptoAccount
	*http.Client
}

// Dial returns a client given a TokenGetter. TokenGetter implementations are
// available in this package, including a Cookie-based cache.
func Dial(ctx context.Context, s oauth2.TokenSource) (*Client, error) {
	c := &Client{
		Client: oauth2.NewClient(context.Background(), s),
	}

	a, err := c.GetAccounts(ctx)
	if len(a) > 0 {
		c.Account = &a[0]
	}
	if err != nil {
		return nil, fmt.Errorf("error getting accounts: %w", err)
	}

	ca, err := c.GetCryptoAccounts(ctx)
	if len(ca) > 0 {
		c.CryptoAccount = &ca[0]
	}

	return c, err
}

// GetAndDecode retrieves from the endpoint and unmarshals resulting json into
// the provided destination interface, which must be a pointer.
func (c *Client) GetAndDecode(ctx context.Context, url string, dest interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	return c.DoAndDecode(ctx, req, dest)
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
func (c *Client) DoAndDecode(ctx context.Context, req *http.Request, dest interface{}) error {
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		b := &bytes.Buffer{}
		var e ErrorMap
		err = json.NewDecoder(io.TeeReader(res.Body, b)).Decode(&e)
		if err != nil {
			return fmt.Errorf("got response %q and could not decode error body %q", res.Status, b.String())
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
