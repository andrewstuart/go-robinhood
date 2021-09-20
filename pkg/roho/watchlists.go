package roho

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// A Watchlist is a list of stock Instruments that an investor is tracking in
// his Robinhood portfolio/app.
type Watchlist struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	User string `json:"user"`

	c *Client
}

// GetWatchlists retrieves the watchlists for a given set of credentials/accounts.
func (c *Client) Watchlists(ctx context.Context) ([]Watchlist, error) {
	var r struct{ Results []Watchlist }
	err := c.get(ctx, baseURL("watchlists"), &r)
	if err != nil {
		return nil, err
	}
	if r.Results != nil {
		for i := range r.Results {
			r.Results[i].c = c
		}
	}
	return r.Results, nil
}

// GetInstruments returns the list of Instruments associated with a Watchlist.
func (w *Watchlist) Instruments(ctx context.Context) ([]Instrument, error) {
	var r struct {
		Results []struct {
			Instrument, URL string
		}
	}

	err := w.c.get(ctx, w.URL, &r)
	if err != nil {
		return nil, err
	}

	insts := make([]*Instrument, len(r.Results))
	eg, ctx := errgroup.WithContext(ctx)

	for i := range r.Results {
		// shadow for safe closure access
		i := i
		eg.Go(func() error {
			inst, err := w.c.Lookup(ctx, r.Results[i].Instrument)
			insts[i] = inst
			return err
		})
	}

	err = eg.Wait()

	// Filter slice for empties (if error)
	retInsts := make([]Instrument, 0, len(r.Results))
	for _, inst := range insts {
		if inst != nil {
			retInsts = append(retInsts, *inst)
		}
	}
	return retInsts, err
}
