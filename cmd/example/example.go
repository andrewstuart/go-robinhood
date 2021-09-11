package main

import (
	"context"
	"log"
	"os"
	"time"

	"astuart.co/go-robinhood/v2"
	"golang.org/x/oauth2"
)

func main() {
	// RH_TOKEN should be set to the value of a Bearer token (a 512+ character alphanumberic string)
	// See https://github.com/sanko/Robinhood/blob/master/Authentication.md
	// Alternatively, look at the value of the 'Authorization' header when visiting the RobinHood website
	token := os.Getenv("RH_TOKEN")

	if token == "" {
		log.Fatal("RH_TOKEN environment variable must be set.")
	}

	ctx := context.Background()

	log.Printf("Logging into Robinhood ...")
	c, err := robinhood.Dial(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}

	log.Printf("Getting portfolios ...")
	ps, err := c.GetPortfolios(ctx)
	if err != nil {
		log.Fatalf("get portfolios failed: %v", err)
	}

	for _, p := range ps {
		log.Printf("found position: %+v", p)

	}

	log.Printf("Checking SPY ...")
	i, err := c.GetInstrumentForSymbol(ctx, "SPY")
	if err != nil {
		log.Fatalf("get instrument failed: %v", err)
	}
	log.Printf("SPY: %+v", i)

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
