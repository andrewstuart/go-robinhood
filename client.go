package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"astuart.co/clyde"
)

const (
	epBase       = "https://api.robinhood.com/"
	epLogin      = epBase + "api-token-auth/"
	epAccounts   = epBase + "accounts/"
	epQuotes     = epBase + "quotes/"
	epPortfolios = epBase + "portfolios/"
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
	*http.Client
}

func Dial(c Creds) (*Client, error) {
	res, err := http.Post(epLogin, "application/x-www-form-urlencoded", strings.NewReader(c.Values().Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var cli Client
	err = json.NewDecoder(res.Body).Decode(&cli)

	cli.Client = &http.Client{
		Transport: clyde.HeaderRoundTripper{"Authorization": "Token " + cli.Token},
	}

	return &cli, err
}

type Meta struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
}
