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
)

type TokenGetter interface {
	GetToken() (string, error)
}

type Creds struct {
	Username, Password string
}

func (c *Creds) GetToken() (string, error) {
	res, err := http.Post(epLogin, "application/x-www-form-urlencoded", strings.NewReader(c.Values().Encode()))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var cli Client
	err = json.NewDecoder(res.Body).Decode(&cli)
	return cli.Token, err
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
	Creds
	Path string
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

	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		return "", err
	}
	defer f.Close()

	tok, err := c.Creds.GetToken()
	if err != nil {
		return "", err
	}

	if tok == "" {
		return "", fmt.Errorf("Empty token")
	}

	_, err = f.Write([]byte(tok))
	return tok, err
}

type Token string

func (t *Token) GetToken() (string, error) {
	return string(*t), nil
}
