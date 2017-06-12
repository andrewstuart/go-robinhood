package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
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
	tw := tabwriter.NewWriter(os.Stdout, 22, 2, 1, ' ', 0)

	updatePos()

	fmt.Fprint(tw, strings.Join(syms, "\t"))
	fmt.Fprint(tw, "\tTotal\n")
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

		tot := 0.0

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

			tot += qty * price

			fmt.Fprintf(tw, "%.0f@$%.3f ", qty, price)

			// Print delta
			c := color.Reset
			if delta > 0.0 {
				c = color.FgGreen
				math.Tru
			} else if delta < -0.0 {
				c = color.FgRed
			}

			color.New(c).Fprintf(tw, "%.2f", delta)
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprintf(tw, "%.2f", tot)
		fmt.Fprintln(tw)
		tw.Flush()

		time.Sleep(2 * time.Second)
	}
}
