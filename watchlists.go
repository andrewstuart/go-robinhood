package robinhood

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

	Client *Client `json:",ignore"`
}

// GetWatchlists retrieves the watchlists for a given set of credentials/accounts.
func (c *Client) GetWatchlists(ctx context.Context) ([]Watchlist, error) {
	var r struct{ Results []Watchlist }
	err := c.GetAndDecode(ctx, EPWatchlists, &r)
	if err != nil {
		return nil, err
	}
	if r.Results != nil {
		for i := range r.Results {
			r.Results[i].Client = c
		}
	}
	return r.Results, nil
}

// GetInstruments returns the list of Instruments associated with a Watchlist.
func (w *Watchlist) GetInstruments(ctx context.Context) ([]Instrument, error) {
	var r struct {
		Results []struct {
			Instrument, URL string
		}
	}

	err := w.Client.GetAndDecode(ctx, w.URL, &r)
	if err != nil {
		return nil, err
	}

	insts := make([]*Instrument, len(r.Results))
	eg, ctx := errgroup.WithContext(ctx)

	for i := range r.Results {
		// shadow for safe closure access
		i := i
		eg.Go(func() error {
			inst, err := w.Client.GetInstrument(ctx, r.Results[i].Instrument)
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
