package robinhood

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

	"golang.org/x/oauth2"
)

var defaultPath = ""

func init() {
	u, err := user.Current()
	if err == nil {
		defaultPath = path.Join(u.HomeDir, ".config", "robinhood.token")
	}
}

// A CredsCacher takes user credentials and a file path. The token obtained
// from the RobinHood API will be cached at the file path, and a new token will
// not be obtained.
type CredsCacher struct {
	Creds oauth2.TokenSource
	Path  string
}

// Token implements TokenSource. It may fail if an error is encountered
// checking the file path provided, or if the underlying creds return an error
// when retrieving their token.
func (c *CredsCacher) Token() (*oauth2.Token, error) {
	if c.Path == "" {
		c.Path = defaultPath
	}

	mustLogin := false

	err := os.MkdirAll(path.Dir(c.Path), 0750)
	if err != nil {
		return nil, fmt.Errorf("error creating path for token: %s", err)
	}

	_, err = os.Stat(c.Path)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			mustLogin = true
		} else {
			return nil, err
		}
	}

	if !mustLogin {
		bs, err := ioutil.ReadFile(c.Path)
		if err != nil {
			return nil, err
		}

		if len(bs) > 0 {
			var o oauth2.Token
			if err := json.Unmarshal(bs, &o); err == nil && o.Valid() {
				return &o, err
			}
		}
	}

	tok, err := c.Creds.Token()
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(tok)
	return tok, err
}
