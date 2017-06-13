package robinhood

import (
	"encoding/json"
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

// Watchlists retrieves the watchlists for a given set of credentials/accounts.
func (c *Client) Watchlists() ([]Watchlist, error) {
	res, err := c.Get(epBase + "watchlists/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r struct{ Results []Watchlist }
	err = json.NewDecoder(res.Body).Decode(&r)
	if r.Results != nil {
		for i := range r.Results {
			r.Results[i].Client = c
		}
	}
	return r.Results, err
}

// GetInstruments returns the list of Instruments associated with a Watchlist.
func (w *Watchlist) GetInstruments() ([]Instrument, error) {
	res, err := w.Client.Get(w.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r struct {
		Results []struct {
			Instrument, URL string
		}
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	insts := make([]Instrument, len(r.Results))
	wg := &sync.WaitGroup{}
	wg.Add(len(r.Results))

	for i := range r.Results {
		go func(i int) {
			defer wg.Done()

			url := r.Results[i].Instrument
			res, err := w.Client.Get(url)
			if err != nil {
				return
			}
			defer res.Body.Close()

			var inst Instrument
			err = json.NewDecoder(res.Body).Decode(&inst)
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
		if (inst != Instrument{}) {
		}
		{
			retInsts = append(retInsts, inst)
		}
	}

	return retInsts, err
}
