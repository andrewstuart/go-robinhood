package roho

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMarketData(t *testing.T) {
	if os.Getenv("ROBINHOOD_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	o := &OAuth{
		Username: os.Getenv("ROBINHOOD_USERNAME"),
		Password: os.Getenv("ROBINHOOD_PASSWORD"),
	}
	ctx := context.Background()

	c, err := Dial(ctx, &CredsCacher{Creds: o})
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	if c == nil {
		t.Errorf("dial returned nil client")
	}

	i, err := c.Lookup(ctx, "SPY")
	if err != nil {
		t.Errorf("lookup failed: %v", err)
	}

	ch, err := c.OptionChains(ctx, i)
	if err != nil {
		t.Errorf("option chains failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insts, err := ch[0].Instrument(ctx, "call", NewDate(2019, 2, 1))
	if err != nil {
		t.Errorf("get instrument failed: %v", err)
	}

	fmt.Printf("len(insts) = %+v\n", len(insts))

	is, err := c.MarketData(context.Background(), insts...)
	if err != nil {
		t.Errorf("market failed: %v", err)
	}

	fmt.Printf("len(is) = %+v\n", len(is))
}
