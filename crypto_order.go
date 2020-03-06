package robinhood

import (
	"bytes"
	"math"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"encoding/json"

	"net/http"
)

// CryptoOrder is the payload to create a crypto currency order
type CryptoOrder struct {
	AccountID      string  `json:"account_id,omitempty"`
	CurrencyPairID string  `json:"currency_pair_id,omitempty"`
	Price          float64 `json:"price,omitempty"`
	RefID          string  `json:"ref_id,omitempty"`
	Side           string  `json:"side,omitempty"`
	TimeInForce    string  `json:"time_in_force,omitempty"`
	Quantity       float64 `json:"quantity,omitempty"`
	Type           string  `json:"type,omitempty"`
}

// CryptoOrderOutput holds the response from api
type CryptoOrderOutput struct {
	Meta
	Account            string        `json:"account"`
	AveragePrice       float64       `json:"average_price,string"`
	CancelURL          string        `json:"cancel"`
	CreatedAt          string        `json:"created_at"`
	CumulativeQuantity string        `json:"cumulative_quantity"`
	CurrencyPairID     string        `json:"currency_pair_id"`
	Executions         []interface{} `json:"executions"`
	ID                 string        `json:"id"`
	LastTransactionAt  string        `json:"last_transaction_at"`
	Price              float64       `json:"price,string"`
	Quantity           string        `json:"quantity"`
	RejectReason       string        `json:"reject_reason"`
	Side               string        `json:"side"`
	State              string        `json:"state"`
	StopPrice          float64       `json:"stop_price,string"`
	TimeInForce        string        `json:"time_in_force"`
	Type               string        `json:"type"`

	client *Client
}

// CryptoOrderOpts encapsulates differences between order types
type CryptoOrderOpts struct {
	Side            OrderSide
	Type            OrderType
	AmountInDollars float64
	Quantity        float64
	Price           float64
	TimeInForce     TimeInForce
	ExtendedHours   bool
	Stop, Force     bool
}

// CryptoOrder will actually place the order
func (c *Client) CryptoOrder(cryptoPair CryptoCurrencyPair, o CryptoOrderOpts) (*CryptoOrderOutput, error) {
	var quantity = math.Round(o.AmountInDollars / o.Price)
	a := CryptoOrder{
		AccountID:      c.CryptoAccount.ID,
		CurrencyPairID: cryptoPair.ID,
		Quantity:       quantity,
		Price:          o.Price,
		RefID:          uuid.New().String(),
		Side:           o.Side.String(),
		TimeInForce:    o.TimeInForce.String(),
		Type:           o.Type.String(),
	}

	payload, err := json.Marshal(a)

	if err != nil {
		return nil, err
	}

	post, err := http.NewRequest("POST", EPCryptoOrders, bytes.NewReader(payload))
	post.Header.Add("Content-Type", "application/json")

	var out CryptoOrderOutput
	err = c.DoAndDecode(post, &out)

	if err != nil {
		return nil, err
	}

	out.client = c
	return &out, nil
}

// Cancel will cancel the order
func (o CryptoOrderOutput) Cancel() error {
	post, err := http.NewRequest("POST", o.CancelURL, nil)
	if err != nil {
		return err
	}

	var output CryptoOrderOutput
	err = o.client.DoAndDecode(post, &output)

	if err != nil {
		return errors.Wrap(err, "could not decode response")
	}

	if output.RejectReason != "" {
		return errors.New(output.RejectReason)
	}

	return nil
}
