package robinhood

import (
	"encoding/json"
	"io"
	"os"
)

type Watchlist struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	User string `json:"user"`

	Client *Client `json:",ignore"`
}

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

func (w *Watchlist) GetInstruments() ([]Instrument, error) {
	res, err := w.Client.Get(w.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r struct{ Results []Instrument }
	err = json.NewDecoder(io.TeeReader(res.Body, os.Stderr)).Decode(&r)
	return r.Results, err
}
