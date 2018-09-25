package robinhood

import (
	"fmt"
	"log"
	"os"

	robinhood "astuart.co/go-robinhood"
)

func ExampleOrder() {
	creds := robinhood.CredsCacher{
		Creds: &robinhood.Creds{
			Username: "andrewstuart",
			Password: os.Getenv("ROBINHOOD_PASSWORD"),
		},
		Path: "/home/andrewstuart/.config/myrh.token",
	}

	c, err := robinhood.Dial(&creds)
	if err != nil {
		log.Fatal(err)
	}

	i, err := c.GetInstrumentForSymbol("SPY")
	if err != nil {
		log.Fatal(err)
	}

	out, err := c.Order(i, robinhood.OrderOpts{
		Price:    100.0,
		Type:     robinhood.Limit,
		Side:     robinhood.Buy,
		Quantity: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out.CancelURL)

	// Output: https://api.robinhood.com/orders/foo/cancel
}
