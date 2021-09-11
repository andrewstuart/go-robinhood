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

	"astuart.co/go-robinhood/v2"
	"golang.org/x/oauth2"
)

func main() {
	user := os.Getenv("RH_USER")
	pass := os.Getenv("RH_PASS")
	if user == "" || pass == "" {
		log.Fatalf("RH_USER and RH_PASS environment variables must be defined.")
	}

	o := &robinhood.CredsCacher{
		Creds: &robinhood.OAuth{
			Username: os.Getenv("RH_USER"),
			Password: os.Getenv("RH_PASS"),
		},
	}

	token, err := o.Token()
	if err != nil {
		log.Fatalf("Unable to get token: %v", err)
	}

	ctx := context.Background()

	log.Printf("Logging into Robinhood as %s ...", user)
	c, err := robinhood.Dial(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}

	log.Printf("Getting portfolios ...")
	ps, err := c.GetPortfolios(ctx)
	if err != nil {
		log.Fatalf("get portfolios failed: %v", err)
	}

	for _, p := range ps {
		log.Printf("portfolio value: $%.2f buying power: $%.2f", p.Equity, p.WithdrawableAmount)
	}

	sym := "SPY"
	log.Printf("Looking up %s ...", sym)
	i, err := c.GetInstrumentForSymbol(ctx, sym)
	if err != nil {
		log.Fatalf("get instrument failed: %v", err)
	}
	log.Printf("SPY is %s", i.Name)

	fs, err := c.GetFundamentals(ctx, sym)
	if err != nil {
		log.Fatalf("get fundamentals failed: %v", err)
	}
	log.Printf("SPY opening price was $%.2f (52 week high: $%.2f)", fs[0].Open, fs[0].High52Weeks)

	qs, err := c.GetQuote(ctx, "SPY")
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
		o, err := c.Order(ctx, i, robinhood.OrderOpts{
			Price:    1.0,
			Side:     robinhood.Buy,
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
		_, err := c.Order(ctx, i, robinhood.OrderOpts{
			Price:    1.0,
			Side:     robinhood.Sell,
			Quantity: 1,
		})
		if err != nil {
			log.Fatalf("sell failed: %v", err)
		}
	default:
		log.Fatalf("%q is an unknown verb", os.Args[1])
	}
}
