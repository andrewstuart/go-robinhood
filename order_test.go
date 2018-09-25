package robinhood

import (
	"fmt"
	"log"
	"os"
)

func ExampleClient_Order() {
	creds := &Creds{
		Username: "andrewstuart",
		Password: os.Getenv("ROBINHOOD_PASSWORD"),
	}

	c, err := Dial(creds)
	if err != nil {
		log.Fatal(err)
	}

	i, err := c.GetInstrumentForSymbol("SPY")
	if err != nil {
		log.Fatal(err)
	}

	out, err := c.Order(i, OrderOpts{
		Price:    100.0,
		Type:     Limit,
		Side:     Buy,
		Quantity: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out.CancelURL)

	// Output: https://api.robinhood.com/orders/foo/cancel
}
