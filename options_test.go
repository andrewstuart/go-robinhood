package robinhood

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarketData(t *testing.T) {
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

	insts, err := ch[0].GetInstrument("put", Date{time.Date(2019, 01, 01, 0, 0, 0, 0, time.Local)})
	asrt.NoError(err)

	_, err = c.MarketData(insts[:10]...)
	asrt.NoError(err)
}
