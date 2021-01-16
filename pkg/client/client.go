package client

import (
	"bytes"
	"context"
	"encoding/json"
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

func (c Client) Create(account Resource) (Resource, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Resource{}, fmt.Errorf("client.Create new uuid: %w", err)
	}

	err = ValidateResource(account)
	if err != nil {
		return Resource{}, fmt.Errorf("client.Create: %w", err)
	}

	payload := Payload{
		Data: Data{
			ID:             id.String(),
			OrganisationID: c.OrganisationID,
			Type:           typeAccounts,
			Attributes:     account,
		},
	}

	jsonPayload, err := marshalPayload(payload)
	if err != nil {
		return Resource{}, fmt.Errorf("client.Create: %w", err)
	}

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		fmt.Sprintf("%s%s", c.BaseURL, createEndpoint),
		jsonPayload,
	)
	if err != nil {
		return Resource{}, fmt.Errorf("client.Create: newRequestWithContext: %w", err)
	}

	req = c.addHeaders(req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return Resource{}, fmt.Errorf("client.Create c.HttpClient.Do: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return Resource{}, fmt.Errorf("client.Create response unexpected response code: %d", resp.StatusCode)
	}

	p, err := unmarshalPayload(resp.Body)
	if err != nil {
		return Resource{}, fmt.Errorf("client.Create: %w", err)
	}

	return p.Data.Attributes, nil
}

func (c Client) List(pageNumber, pageSize uint) ([]Resource, error) {
	requestPath := fmt.Sprintf(listEndpoint, pageNumber, pageSize)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		fmt.Sprintf("%s%s", c.BaseURL, requestPath),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("client.List http.NewRequestWithContext: %w", err)
	}

	req = c.addHeaders(req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.List httpClient.Do: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client.List unexpected http response status: %d", resp.StatusCode)
	}

	mp, err := unmarshalMultiPayload(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("client.List: %w", err)
	}

	resources := make([]Resource, 0)
	for _, d := range mp.Data {
		resources = append(resources, d.Attributes)
	}

	return resources, nil
}

func (c Client) Fetch(accountID string) (Resource, error) {
	requestPath := fmt.Sprintf(fetchEndpoint, accountID)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		fmt.Sprintf("%s%s", c.BaseURL, requestPath),
		nil,
	)
	if err != nil {
		return Resource{}, fmt.Errorf("client.Fetch http.NewRequestWithContext: %w", err)
	}

	req = c.addHeaders(req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return Resource{}, fmt.Errorf("client.Fetch httpClient.Do: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Resource{}, fmt.Errorf("client.Fetch unexpected response code: %d", resp.StatusCode)
	}

	p, err := unmarshalPayload(resp.Body)
	if err != nil {
		return Resource{}, fmt.Errorf("client.Fetch: %w", err)
	}

	return p.Data.Attributes, nil
}

func (c Client) Delete(accountID string, version uint) error {
	requestPath := fmt.Sprintf(deleteEndpoint, accountID, version)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		fmt.Sprintf("%s%s", c.BaseURL, requestPath),
		nil,
	)
	if err != nil {
		return fmt.Errorf("client.Delete http.NewRequestWithContext: %w", err)
	}

	req = c.addHeaders(req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client.Delete httpClient.Do: %w", err)
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

	return p, nil
}

// unmarshalMultiPayload will turn a json with an array of payloads in the data part into a MultiPayload struct.
func unmarshalMultiPayload(r io.Reader) (MultiPayload, error) {
	var mp MultiPayload

	err := json.NewDecoder(r).Decode(&mp)
	if err != nil {
		return MultiPayload{}, fmt.Errorf("unmarshalMultiPayload: %w", err)
	}

	return mp, nil
}
