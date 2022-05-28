package robinhood

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// DefaultClientID is used by the website.
const DefaultClientID = "c82SH0WZOsabOXGP2sxqcj34FxkvfnWRZBKlBjFS"

// OAuth implements oauth2 using the robinhood implementation
type OAuth struct {
	Endpoint, ClientID, Username, Password, MFA string
}

// ErrMFARequired indicates the MFA was required but not provided.
var ErrMFARequired = fmt.Errorf("Two Factor Auth code required and not supplied")

// Token implements TokenSource
func (p *OAuth) Token() (*oauth2.Token, error) {
	cliID := p.ClientID
	if cliID == "" {
		cliID = DefaultClientID
	}

	u, _ := url.Parse(EPLogin)
	q := map[string]interface{}{
		"expires_in":   86400,
		"client_id":    cliID,
		"device_token": "34898bf2-3aad-4540-8401-bea572ab8c09",
		"grant_type":   "password",
		"scope":        "internal",
		"username":     p.Username,
		"password":     p.Password,
	}
	if p.MFA != "" {
		q["mfa_code"] = p.MFA
	}

	bs, _ := json.Marshal(q)

	req, err := http.NewRequest(
		"POST",
		u.String(),
		bytes.NewReader(bs),
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Robinhood-API-Version", "1.431.4")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not post login")
	}
	defer res.Body.Close()

	var o struct {
		oauth2.Token
		ExpiresIn   int    `json:"expires_in"`
		MFARequired bool   `json:"mfa_required"`
		MFAType     string `json:"mfa_type"`
	}

	err = json.NewDecoder(res.Body).Decode(&o)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode token")
	}

	if o.MFARequired {
		return nil, ErrMFARequired
	}

	o.Token.Expiry = time.Now().Add(time.Duration(o.ExpiresIn) * time.Second)

	return &o.Token, nil
}
