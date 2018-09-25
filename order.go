package robinhood

import (
	"net/http"

	"github.com/pkg/errors"
)

// OrderOutput is the response from the Order api
type OrderOutput struct {
	Account                string   `json:"account"`
	AveragePrice           float64  `json:"average_price"`
	CancelURL              string   `json:"cancel"`
	CreatedAt              string   `json:"created_at"`
	CumulativeQuantity     string   `json:"cumulative_quantity"`
	Executions             []string `json:"executions"`
	ExtendedHours          bool     `json:"extended_hours"`
	Fees                   string   `json:"fees"`
	ID                     string   `json:"id"`
	Instrument             string   `json:"instrument"`
	LastTransactionAt      string   `json:"last_transaction_at"`
	OverrideDayTradeChecks bool     `json:"override_day_trade_checks"`
	OverrideDtbpChecks     bool     `json:"override_dtbp_checks"`
	Position               string   `json:"position"`
	Price                  float64  `json:"price"`
	Quantity               string   `json:"quantity"`
	RejectReason           string   `json:"reject_reason"`
	Side                   string   `json:"side"`
	State                  string   `json:"state"`
	StopPrice              float64  `json:"stop_price"`
	TimeInForce            string   `json:"time_in_force"`
	Trigger                string   `json:"trigger"`
	Type                   string   `json:"type"`
	UpdatedAt              string   `json:"updated_at"`
	URL                    string   `json:"url"`

	client *Client
}

// Cancel attempts to cancel an odrer
func (o OrderOutput) Cancel() error {
	post, err := http.NewRequest("POST", o.CancelURL, nil)
	if err != nil {
		return err
	}

	var o2 OrderOutput
	err = o.client.DoAndDecode(post, &o2)
	if err != nil {
		return errors.Wrap(err, "could not decode response")
	}

	if o2.RejectReason != "" {
		return errors.New(o2.RejectReason)
	}
	return nil
}
