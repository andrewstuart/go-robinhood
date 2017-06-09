package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type ip struct {
	Instrument
	Position
	last *Quote
}

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

	syms := []string{}
	positions := map[string]*ip{}

	updatePos := func() {
		pos, err := c.GetPositions(a[0])
		if err != nil {
			log.Fatal(err)
		}

		for _, p := range pos {
			inst, err := c.GetInstrument(p.Instrument)
			if err != nil {
				log.Fatal(err)
			}

			if positions[inst.Symbol] == nil {
				positions[inst.Symbol] = &ip{Instrument: *inst, Position: p}

				q, _ := strconv.ParseFloat(p.Quantity, 64)

				if q > 0.0 || inst.Symbol == "NVDA" || inst.Symbol == "AMD" {
					syms = append(syms, inst.Symbol)
				}
			} else {
				positions[inst.Symbol].Instrument = *inst
				positions[inst.Symbol].Position = p
			}

		}
	}
	tw := tabwriter.NewWriter(os.Stdout, 30, 2, 1, ' ', 0)

	updatePos()

	for _, s := range syms {
		fmt.Fprint(tw, s+"\t")
	}
	fmt.Fprint(tw, "\n")
	tw.Flush()

	i := 0
	for {
		i++
		if i%5 == 0 {
			updatePos()
		}

		q, err := c.GetQuote(syms...)
		if err != nil {
			log.Fatal("Quotes error", err)
		}

		for _, q := range q {
			qip := positions[q.Symbol]
			if qip.last == nil {
				qip.last = &Quote{}
				*qip.last = q
			}
			last := qip.last

			qty, _ := strconv.ParseFloat(qip.Quantity, 64)
			price, _ := strconv.ParseFloat(q.LastTradePrice, 64)

			lastp, _ := strconv.ParseFloat(last.LastTradePrice, 64)

			delta := (price - lastp) / price * 100.0

			*qip.last = q

			fmt.Fprintf(tw, "%.0f@$%.3f ($%.2f) %.2f\t", qty, price, qty*price, delta)
		}
		fmt.Fprintln(tw)
		tw.Flush()

		time.Sleep(10 * time.Second)
	}
}
