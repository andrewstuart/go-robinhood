package main

import (
	"net/http"
	"net/url"
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
	URL       url.URL   `json:"url"`
}
