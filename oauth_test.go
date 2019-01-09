package robinhood

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOauthPassword(t *testing.T) {
	if os.Getenv("ROBINHOOD_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)

	o := &CredsCacher{
		Creds: &OAuth{
			Username: os.Getenv("ROBINHOOD_USERNAME"),
			Password: os.Getenv("ROBINHOOD_PASSWORD"),
		},
	}

	tok, err := o.Token()
	asrt.NoError(err)
	asrt.NotNil(tok)
}
