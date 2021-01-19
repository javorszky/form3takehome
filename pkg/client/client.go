package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/javorszky/form3takehome/pkg/config"
)

const (
	acceptHeaderValue = "application/vnd.api+json"
	createEndpoint    = "/v1/organisation/accounts"
	listEndpoint      = "/v1/organisation/accounts?page[number]=%d&page[size]=%d"
	fetchEndpoint     = "/v1/organisation/accounts/%s"
	deleteEndpoint    = "/v1/organisation/accounts/%s?version=%d"
	typeAccounts      = "accounts"
)

type Client struct {
	BaseURL        string
	OrganisationID string
	HttpClient     http.Client
	DateLocation   *time.Location
}

// New returns a configured Client struct.
func New(cfg config.Config, c http.Client, gmt *time.Location) Client {
	return Client{
		BaseURL:        cfg.AccountsAPIURL,
		OrganisationID: cfg.OrganisationID,
		HttpClient:     c,
		DateLocation:   gmt,
	}
}

// Create will create a Resource that belongs to organisation ID set on the Client if the Resource passes validation for
// the given dataset.
func (c Client) Create(account Resource) (Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Payload{}, fmt.Errorf("client.Create new uuid: %w", err)
	}

	err = ValidateResource(account)
	if err != nil {
		return Payload{}, fmt.Errorf("client.Create: %w", err)
	}

	requestPayload := Payload{
		Data: Data{
			ID:             id.String(),
			OrganisationID: c.OrganisationID,
			Type:           typeAccounts,
			Attributes:     account,
		},
	}

	jsonPayload, err := marshalPayload(requestPayload)
	if err != nil {
		return Payload{}, fmt.Errorf("client.Create: %w", err)
	}

	resp, err := c.do(http.MethodPost, createEndpoint, jsonPayload)
	if err != nil {
		return Payload{}, fmt.Errorf("client.Create: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return Payload{}, fmt.Errorf("client.Create response unexpected response code: %d", resp.StatusCode)
	}

	p, err := unmarshalPayload(resp.Body)
	if err != nil {
		return Payload{}, fmt.Errorf("client.Create: %w", err)
	}

	return p, nil
}

// List will list all the Resources that belong to given organisation ID, pageSize per request, and if multi paged, on
// the given pageNumber.
func (c Client) List(pageNumber, pageSize uint) (MultiPayload, error) {
	requestPath := fmt.Sprintf(listEndpoint, pageNumber, pageSize)

	resp, err := c.do(http.MethodGet, requestPath, nil)
	if err != nil {
		return MultiPayload{}, fmt.Errorf("client.List: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return MultiPayload{}, fmt.Errorf("client.List unexpected http response status: %d", resp.StatusCode)
	}

	mp, err := unmarshalMultiPayload(resp.Body)
	if err != nil {
		return MultiPayload{}, fmt.Errorf("client.List: %w", err)
	}

	return mp, nil
}

// Fetch will return a Resource struct identified by an ID, if exists.
func (c Client) Fetch(accountID string) (Payload, error) {
	requestPath := fmt.Sprintf(fetchEndpoint, accountID)

	resp, err := c.do(http.MethodGet, requestPath, nil)
	if err != nil {
		return Payload{}, fmt.Errorf("client.Fetch httpClient.Do: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Payload{}, fmt.Errorf("client.Fetch unexpected response code: %d", resp.StatusCode)
	}

	p, err := unmarshalPayload(resp.Body)
	if err != nil {
		return Payload{}, fmt.Errorf("client.Fetch: %w", err)
	}

	return p, nil
}

// Delete will remove a Resource with given ID if version that's requested to be deleted and current version of Resource
// matches.
func (c Client) Delete(accountID string, version uint) error {
	requestPath := fmt.Sprintf(deleteEndpoint, accountID, version)

	resp, err := c.do(http.MethodDelete, requestPath, nil)
	if err != nil {
		return fmt.Errorf("client.Delete: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("client.Delete unexpected response code: %d", resp.StatusCode)
	}

	return nil
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

	bodyReadCloser, _ := r.GetBody()

	body, err := ioutil.ReadAll(bodyReadCloser)
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

// marshalPayload will turn a Payload struct to its json representation.
func marshalPayload(r Payload) (io.Reader, error) {
	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(r)
	if err != nil {
		return nil, fmt.Errorf("marshalPayload: %w", err)
	}

	return b, nil
}

// unmarshalPayload will turn a json in an io.Reader into a Payload struct.
func unmarshalPayload(r io.Reader) (Payload, error) {
	var p Payload

	err := json.NewDecoder(r).Decode(&p)
	if err != nil {
		return Payload{}, fmt.Errorf("unmarshalPayload: %w", err)
	}

	if p.Data == (Data{}) {
		return Payload{}, errors.New("unmarshalPayload: Data is empty on the decoded Payload")
	}

	if p.Data.Attributes == (Resource{}) {
		return Payload{}, errors.New("unmarshalPayload: Data.Attributes is empty on the decoded Payload")
	}

	return p, nil
}

// unmarshalMultiPayload will turn a json with an array of payloads in the data part into a MultiPayload struct.
func unmarshalMultiPayload(r io.Reader) (MultiPayload, error) {
	var mp MultiPayload

	err := json.NewDecoder(r).Decode(&mp)
	if err != nil {
		return MultiPayload{}, fmt.Errorf("unmarshalMultiPayload: %w", err)
	}

	if mp.Data == nil {
		return MultiPayload{}, errors.New("unmarshalMultiPayload: there is no Data on the decoded MultiPayload")
	}

	for _, d := range mp.Data {
		if d.Attributes == (Resource{}) {
			return MultiPayload{}, errors.New("unmarshalMultiPayload: Data structs are missing required fields")
		}
	}

	return mp, nil
}

// do is a generic method to handle network calls.
func (c Client) do(method, endpoint string, payload io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		method,
		fmt.Sprintf("%s%s", c.BaseURL, endpoint),
		payload,
	)
	if err != nil {
		return nil, fmt.Errorf("client.do http.NewRequestWithContext: %w", err)
	}

	req = c.addHeaders(req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.do httpClient.Do: %w", err)
	}

	return resp, nil
}
