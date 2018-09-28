package robinhood

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// DefaultClientID is used by the website.
const DefaultClientID = "c82SH0WZOsabOXGP2sxqcj34FxkvfnWRZBKlBjFS"

// PasswordToken implements oauth2 using the robinhood implementation
type PasswordToken struct {
	ClientID, Username, Password string
}

// Token implements TokenSource
func (p *PasswordToken) Token() (*oauth2.Token, error) {
	v := url.Values{
		"username": []string{p.Username},
		"password": []string{p.Password},
	}

	cli := p.ClientID
	if cli == "" {
		cli = DefaultClientID
	}

	u, _ := url.Parse(EPLogin)
	q := u.Query()
	q.Add("client_id", cli)
	q.Add("grant_type", "password")
	q.Add("scope", "internal")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(
		"POST",
		u.String(),
		strings.NewReader(v.Encode()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not post login")
	}
	defer res.Body.Close()

	var o struct {
		oauth2.Token
		ExpiresIn int `json:"expires_in"`
	}
	err = json.NewDecoder(res.Body).Decode(&o)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode token")
	}

	o.Token.Expiry = time.Now().Add(time.Duration(o.ExpiresIn) * time.Second)

	return &o.Token, nil
}
