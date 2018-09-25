package robinhood

// Instrument is a type to represent the "instrument" API type in the
// unofficial robinhood API.
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

	c Client
}

// GetInstrument returns an Instrument given a URL
func (c Client) GetInstrument(instURL string) (*Instrument, error) {
	var i Instrument
	err := c.GetAndDecode(instURL, &i)
	if err != nil {
		return nil, err
	}
	return &i, err
}

// GetInstrumentForSymbol returns an Instrument given a ticker symbol
func (c Client) GetInstrumentForSymbol(sym string) (*Instrument, error) {
	var i struct {
		Results []Instrument
	}
	err := c.GetAndDecode(EPInstruments+"?symbol="+sym, &i)
	return &i.Results[0], err
}
