package main

type Instrument struct {
	BloombergUnique    string `json:"bloomberg_unique"`
	Country            `json:"country"`
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
	Splits             string      `json:"splits"`
	State              string      `json:"state"`
	Symbol             string      `json:"symbol"`
	Tradeable          bool        `json:"tradeable"`
	URL                string      `json:"url"`
}

type Country struct {
	Alpha3        string      `json:"alpha3"`
	Code          string      `json:"code"`
	Flag          string      `json:"flag"`
	FlagURL       interface{} `json:"flag_url"`
	IocCode       string      `json:"ioc_code"`
	Name          string      `json:"name"`
	Numeric       interface{} `json:"numeric"`
	NumericPadded interface{} `json:"numeric_padded"`
}
