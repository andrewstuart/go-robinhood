package robinhood

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

type TokenGetter interface {
	GetToken() (string, error)
}

// Creds are the typical Username/Password
type Creds struct {
	Username, Password string
}

type OIDC struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func (c *Creds) GetToken() (string, error) {
	req, err := http.NewRequest("POST", EPLogin+"?grant_type=password&scope=internal&client_id=c82SH0WZOsabOXGP2sxqcj34FxkvfnWRZBKlBjFS", strings.NewReader(c.Values().Encode()))
	if err != nil {
		return "", errors.Wrap(err, "could not create request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "could not post login")
	}
	defer res.Body.Close()

	var o OIDC
	err = json.NewDecoder(res.Body).Decode(&o)
	return o.AccessToken, errors.Wrap(err, "could not decode json")
}

func (c *Creds) Refresh() (string, error) {
	return c.GetToken()
}

func (c Creds) Values() url.Values {
	return url.Values{
		"username": []string{c.Username},
		"password": []string{c.Password},
	}
}

// A CredsCacher takes user credentials and a file path. The token obtained
// from the RobinHood API will be cached at the file path, and a new token will
// not be obtained.
type CredsCacher struct {
	Creds TokenGetter
	Path  string
}

// GetToken implements TokenGetter. It may fail if an error is encountered
// checking the file path provided, or if the underlying creds return an error
// when retrieving their token.
func (c *CredsCacher) GetToken() (string, error) {
	mustLogin := false

	err := os.MkdirAll(path.Dir(c.Path), 0750)
	if err != nil {
		return "", fmt.Errorf("error creating path for token: %s", err)
	}

	_, err = os.Stat(c.Path)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			mustLogin = true
		} else {
			return "", err
		}
	}

	if !mustLogin {
		bs, err := ioutil.ReadFile(c.Path)
		return string(bs), err
	}

	tok, err := c.Creds.GetToken()
	if err != nil {
		return "", err
	}

	if tok == "" {
		return "", fmt.Errorf("Empty token")
	}

	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write([]byte(tok))
	return tok, err
}

type Token string

func (t *Token) GetToken() (string, error) {
	return string(*t), nil
}
