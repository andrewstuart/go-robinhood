package robinhood

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestMarketData(t *testing.T) {
	if os.Getenv("ROBINHOOD_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	o := &OAuth{
		Username: os.Getenv("ROBINHOOD_USERNAME"),
		Password: os.Getenv("ROBINHOOD_PASSWORD"),
	}

	c, err := Dial(&CredsCacher{Creds: o})

	asrt.NoError(err)
	asrt.NotNil(c)

	i, err := c.GetInstrumentForSymbol("SPY")
	asrt.NoError(err)

	ch, err := c.GetOptionChains(i)
	asrt.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insts, err := ch[0].GetInstrument(ctx, "call", NewDate(2019, 2, 1))
	asrt.NoError(err)

	fmt.Printf("len(insts) = %+v\n", len(insts))

	is, err := c.MarketData(insts...)
	asrt.NoError(err)

	spew.Dump(is)
	fmt.Printf("len(is) = %+v\n", len(is))
}
