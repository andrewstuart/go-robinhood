package robinhood

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestOauthPassword(t *testing.T) {
	asrt := assert.New(t)

	o := OAuth{
		Username: os.Getenv("ROBINHOOD_USERNAME"),
		Password: os.Getenv("ROBINHOOD_PASSWORD"),
	}

	tok, err := o.Token()
	asrt.NoError(err)
	asrt.NotNil(tok)
	spew.Dump(tok)
}
