package main

import (
	"io"
	"log"
	"os"
)

func main() {
	creds := Creds{
		Username: "andrewstuart",
		Password: os.Getenv("ROBINHOOD_PASSWORD"),
	}

	c, err := Dial(creds)
	if err != nil {
		log.Fatal(err)
	}
	a, err := c.GetAccounts()
	if err != nil {
		log.Fatal(err)
	}

	pos, err := c.GetPositions(a[0])
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range pos {
		res, err := c.Get(p.Instrument)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		io.Copy(os.Stdout, res.Body)
		log.Println("Next")
	}
}
