package roho

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type Config struct {
	Username string
	Password string
}

func New(ctx context.Context, c *Config) (*Client, error) {
	user := c.Username
	if user == "" {
		user = os.Getenv("RH_USER")
	}

	pass := c.Password
	if pass == "" {
		pass = os.Getenv("RH_PASS")
	}
	o := &CredsCacher{Creds: &OAuth{Username: user, Password: pass}}

	token, err := o.Token()
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	log.Printf("Logging into Robinhood as %s ...", user)
	return Dial(ctx, oauth2.StaticTokenSource(token))
}

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

	a, err := c.Accounts(ctx)
	log.Printf("Found %d accounts", len(a))
	if len(a) > 0 {
		c.Account = &a[0]
	}

	ca, err := c.CryptoAccounts(ctx)
	log.Printf("Found %d crypto accounts", len(ca))
	if len(ca) > 0 {
		c.CryptoAccount = &ca[0]
	}

	return c, err
}
