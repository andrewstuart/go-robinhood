// This is an example robinhood client written in Go
package main

// usage:
//
// RH_USER=email@example.org RH_PASS=password go run .

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

func main() {
	user := os.Getenv("RH_USER")
	pass := os.Getenv("RH_PASS")
	if user == "" || pass == "" {
		log.Fatalf("RH_USER and RH_PASS environment variables must be defined.")
	}

	ctx := context.Background()
	c, err := roho.New(ctx, &roho.Config{Username: user, Password: pass})
	if err != nil {
		log.Fatalf("new failed: %v", err)
	}

	log.Printf("Getting portfolios ...")
	ps, err := c.Portfolios(ctx)
	if err != nil {
		log.Fatalf("get portfolios failed: %v", err)
	}

	for _, p := range ps {
		log.Printf("portfolio value: $%.2f buying power: $%.2f", p.Equity, p.WithdrawableAmount)
	}

	sym := "SPY"
	log.Printf("Looking up %s ...", sym)
	i, err := c.Lookup(ctx, sym)
	if err != nil {
		log.Fatalf("get instrument failed: %v", err)
	}
	log.Printf("SPY is %s", i.Name)

	fs, err := c.Fundamentals(ctx, sym)
	if err != nil {
		log.Fatalf("get fundamentals failed: %v", err)
	}
	log.Printf("SPY opening price was $%.2f (52 week high: $%.2f)", fs[0].Open, fs[0].High52Weeks)

	/*	hs, err := c.Historicals(ctx, "5minute", "week", sym)
		if err != nil {
			log.Fatalf("get historicals failed: %v", err)
		}
		for _, r := range hs[0].Records {
			log.Printf("  %s: opened at %.2f, closed at %.2f", r.BeginsAt, r.OpenPrice, r.ClosePrice)
		}
	*/
	qs, err := c.Quote(ctx, "SPY")
	if err != nil {
		log.Fatalf("get quote failed: %v", err)
	}
	log.Printf("SPY current price is $%.2f", qs[0].Price())

	if len(os.Args) == 1 {
		return
	}

	switch os.Args[1] {
	case "buy":
		log.Printf("Buying 1 share of %s ...", i.Symbol)
		o, err := c.Buy(ctx, i, roho.OrderOpts{
			Price:    1.0,
			Quantity: 1,
		})

		if err != nil {
			log.Fatalf("buy failed: %v", err)
		}

		time.Sleep(5 * time.Millisecond)

		log.Printf("Need to buy groceries - cancelling buy of %s ...", i.Symbol)
		err = o.Cancel(ctx)
		if err != nil {
			log.Fatalf("buy failed: %v", err)
		}
	case "sell":
		log.Printf("Selling 1 share of %s ...", i.Symbol)
		_, err := c.Sell(ctx, i, roho.OrderOpts{
			Price:    1.0,
			Quantity: 1,
		})
		if err != nil {
			log.Fatalf("sell failed: %v", err)
		}
	default:
		log.Fatalf("%q is an unknown verb", os.Args[1])
	}
}
