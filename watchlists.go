package robinhood

import (
	"sync"
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
func (c *Client) GetWatchlists() ([]Watchlist, error) {
	var r struct{ Results []Watchlist }
	err := c.GetAndDecode(epWatchlists, &r)
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
func (w *Watchlist) GetInstruments() ([]Instrument, error) {
	var r struct {
		Results []struct {
			Instrument, URL string
		}
	}

	err := w.Client.GetAndDecode(w.URL, &r)
	if err != nil {
		return nil, err
	}

	insts := make([]*Instrument, len(r.Results))
	wg := &sync.WaitGroup{}
	wg.Add(len(r.Results))

	for i := range r.Results {
		go func(i int) {
			defer wg.Done()

			inst, err := w.Client.GetInstrument(r.Results[i].Instrument)
			if err != nil {
				return
			}

			insts[i] = inst
		}(i)
	}

	wg.Wait()

	// Filter slice for empties (if error)
	retInsts := []Instrument{}
	for _, inst := range insts {
		if inst != nil {
			retInsts = append(retInsts, *inst)
		}
	}

	return retInsts, err
}
