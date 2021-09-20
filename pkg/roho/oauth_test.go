package roho

import (
	"os"
	"testing"
)

func TestOauthPassword(t *testing.T) {
	if os.Getenv("ROBINHOOD_USERNAME") == "" {
		t.Skip("No username set")
		return
	}

	o := &CredsCacher{
		Creds: &OAuth{
			Username: os.Getenv("ROBINHOOD_USERNAME"),
			Password: os.Getenv("ROBINHOOD_PASSWORD"),
		},
	}

	tok, err := o.Token()
	if err != nil {
		t.Errorf("token failed: %v", err)
	}

	if tok == nil {
		t.Errorf("got nil token back")
	}
}
