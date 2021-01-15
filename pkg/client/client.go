package client

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/javorszky/form3takehome/pkg/config"
)

const acceptHeaderValue = "application/vnd.api+json"

type Client struct {
	BaseURL      string
	DateLocation *time.Location
}

// New returns a configured Client struct.
func New(cfg config.Config, gmt *time.Location) Client {
	return Client{
		BaseURL:      cfg.AccountsAPIURL,
		DateLocation: gmt,
	}
}

func (c Client) Create(account Resource) (Resource, error) {
	return Resource{}, nil
}

func (c Client) List() {
}

// addHeaders will decorate a header with the needed key/value pairs. If the body is not empty, it also adds the
// Content-Type header.
//
// Ideally the Content-Length header should not be added automatically, because per the http documentation of go, if
// it's not set, and the total size of all written data is under a few KB and there are no Flush calls, the
// Content-Length header will be added automatically on A Write call. I' adding it here regardless so it shows up for
// the tests.
//
// Authorization headers are not added per the spec of the take home exercise.
func (c Client) addHeaders(r *http.Request) *http.Request {
	r.Header.Add("Host", c.BaseURL)
	r.Header.Add("Date", c.currentHTTPDate())
	r.Header.Add("Accept", acceptHeaderValue)

	if r.Body == nil {
		return r
	}

	body, err := ioutil.ReadAll(r.Body)
	if err == nil && len(body) > 0 {
		r.Header.Add("Content-Type", acceptHeaderValue)
		r.Header.Add("Content-Length", strconv.Itoa(len(body)))
	}

	return r
}

// currentHTTPDate returns the current date time in GMT, per RFC 7231/7.1.1.1.
func (c Client) currentHTTPDate() string {
	return time.Now().In(c.DateLocation).Format(time.RFC1123)
}
