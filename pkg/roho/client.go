package roho

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func baseURL(s string) string {
	return "https://api.robinhood.com/" + s + "/"
}

func cryptoURL(s string) string {
	return "https://nummus.robinhood.com/" + s + "/"
}

// call retrieves from the endpoint and unmarshals resulting json into
// the provided destination interface, which must be a pointer.
func (c *Client) get(ctx context.Context, url string, dest interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	return c.call(ctx, req, dest)
}

// ErrorMap encapsulates the helpful error messages returned by the API server
type ErrorMap map[string]interface{}

func (e ErrorMap) Error() string {
	es := make([]string, 0, len(e))
	for k, v := range e {
		es = append(es, fmt.Sprintf("%s: %q", k, v))
	}
	return "Error returned from API: " + strings.Join(es, ", ")
}

// call provides useful abstractions around common errors and decoding issues.
func (c *Client) call(ctx context.Context, req *http.Request, dest interface{}) error {
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		b := &bytes.Buffer{}
		var e ErrorMap
		err = json.NewDecoder(io.TeeReader(res.Body, b)).Decode(&e)
		if err != nil {
			return fmt.Errorf("got response %q and could not decode error body %q", res.Status, b.String())
		}
		return e
	}

	return json.NewDecoder(res.Body).Decode(dest)
}

// Meta holds metadata common to many RobinHood types.
type Meta struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
}
