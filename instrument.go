package robinhood

import "encoding/json"

type Instrument struct {
	BloombergUnique    string      `json:"bloomberg_unique"`
	Country            string      `json:"country"`
	DayTradeRatio      string      `json:"day_trade_ratio"`
	Fundamentals       string      `json:"fundamentals"`
	ID                 string      `json:"id"`
	ListDate           string      `json:"list_date"`
	MaintenanceRatio   string      `json:"maintenance_ratio"`
	MarginInitialRatio string      `json:"margin_initial_ratio"`
	Market             string      `json:"market"`
	MinTickSize        interface{} `json:"min_tick_size"`
	Name               string      `json:"name"`
	Quote              string      `json:"quote"`
	SimpleName         interface{} `json:"simple_name"`
	Splits             string      `json:"splits"`
	State              string      `json:"state"`
	Symbol             string      `json:"symbol"`
	Tradeable          bool        `json:"tradeable"`
	URL                string      `json:"url"`
}

func (c Client) GetInstrument(instURL string) (*Instrument, error) {
	res, err := c.Get(instURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var i Instrument
	err = json.NewDecoder(res.Body).Decode(&i)
	return &i, err
}
