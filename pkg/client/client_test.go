package client_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/javorszky/form3takehome/pkg/client"
	"github.com/javorszky/form3takehome/pkg/config"
)

const testTimeoutMs = 500

func TestNew(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load GMT location")
	}

	testClient := http.Client{
		Timeout: 30 * time.Second,
	}

	type args struct {
		cfg config.Config
		gmt *time.Location
	}

	tests := []struct {
		name string
		args args
		want client.Client
	}{
		{
			name: "constructs a new client based on incoming data",
			args: args{
				cfg: config.Config{
					AccountsAPIURL: "https://testurl",
					OrganisationID: "orgid",
				},
				gmt: gmtLoc,
			},
			want: client.Client{
				BaseURL:        "https://testurl",
				OrganisationID: "orgid",
				DateLocation:   gmtLoc,
				HttpClient:     testClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, client.New(tt.args.cfg, testClient, tt.args.gmt))
		})
	}
}

func TestClient_Create(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	type args struct {
		account client.Resource
	}

	tests := []struct {
		name        string
		handlerFunc http.HandlerFunc
		args        args
		want        client.Resource
		wantErr     bool
	}{
		{
			name: "correctly responds to server returning 201 Created",
			handlerFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = fmt.Fprint(w, returnCompactFile(t, "./testdata/payload.json"))
			}),
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want: client.Resource{
				Country:       "GB",
				BaseCurrency:  "GBP",
				BankID:        "89282dd",
				BankIDCode:    "12221",
				AccountNumber: "12345678",
				BIC:           "bic1234",
				IBAN:          "iban1234",
				CustomerID:    "anuuidv4again",
				Name: [4]string{
					"line1",
					"line2",
					"line3",
					"line4",
				},
				AlternativeNames: [3]string{
					"altname1",
					"altname2",
					"altname3",
				},
				AccountClassification:   "cop",
				JointAccount:            false,
				AccountMatchingOptOut:   false,
				SecondaryIdentification: "some custom name",
				Switched:                false,
				Status:                  "confirmed",
			},
			wantErr: false,
		},
		{
			name: "returns error if response from server is anything other than a 201",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if server takes longer to respond than the configured timeout",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep((testTimeoutMs + 100) * time.Millisecond)
				w.WriteHeader(http.StatusCreated)
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if server responds with non json body",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, "not a json")
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if server responds with json that would result in an empty Data attribute",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, `{"error":"a string, but not something that can become a Payload"}`)
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if server responds with json that would result in an empty Data.Attributes field",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, `{"data":{"id": "a string, but not something that can become a Payload"}}`)
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if request Resource fails validation",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusSeeOther)
				// not used here
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Resource{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handlerFunc)
			defer ts.Close()

			c := client.New(
				config.Config{
					AccountsAPIURL: ts.URL,
					OrganisationID: "orgid",
				},
				http.Client{
					Timeout: testTimeoutMs * time.Millisecond,
				},
				gmtLoc,
			)

			got, err := c.Create(tt.args.account)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_CreateBadURL(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	type args struct {
		account client.Resource
	}

	tests := []struct {
		name        string
		handlerFunc http.HandlerFunc
		args        args
		want        client.Resource
		wantErr     bool
	}{
		{
			name: "returns error if client is misconfigured and a new http request can't be created",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusSeeOther)
				// not used here
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Resource{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.New(
				config.Config{
					AccountsAPIURL: "htt@ps://bla",
					OrganisationID: "orgid",
				},
				http.Client{
					Timeout: testTimeoutMs * time.Millisecond,
				},
				gmtLoc,
			)

			got, err := c.Create(tt.args.account)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			fmt.Printf("%v", err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_Fetch(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	type args struct {
		accountID string
	}

	tests := []struct {
		name        string
		handlerFunc http.HandlerFunc
		args        args
		want        client.Resource
		wantErr     bool
	}{
		{
			name: "correctly returns a resource by id",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, returnCompactFile(t, "./testdata/payload.json"))
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want: client.Resource{
				Country:       "GB",
				BaseCurrency:  "GBP",
				BankID:        "89282dd",
				BankIDCode:    "12221",
				AccountNumber: "12345678",
				BIC:           "bic1234",
				IBAN:          "iban1234",
				CustomerID:    "anuuidv4again",
				Name: [4]string{
					"line1",
					"line2",
					"line3",
					"line4",
				},
				AlternativeNames: [3]string{
					"altname1",
					"altname2",
					"altname3",
				},
				AccountClassification:   "cop",
				JointAccount:            false,
				AccountMatchingOptOut:   false,
				SecondaryIdentification: "some custom name",
				Switched:                false,
				Status:                  "confirmed",
			},
			wantErr: false,
		},
		{
			name: "returns error if response takes longer than configured timeout",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep((testTimeoutMs + 100) * time.Millisecond)
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, returnCompactFile(t, "./testdata/payload.json"))
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if response is a non-200 code",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but not a json",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "not a json")
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but there is no body",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but json can't be populated into a meaningful resource",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"error":"not data"}`)
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but json can't be populated into a meaningful resource",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"data":{"randomKey": "not data"}}`)
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Resource{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but json can't be populated into a meaningful resource",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"data":{"id": "uuidve-missingattributes"}}`)
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Resource{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handlerFunc)
			defer ts.Close()

			c := client.Client{
				BaseURL:        ts.URL,
				OrganisationID: "orgid",
				HttpClient: http.Client{
					Timeout: testTimeoutMs * time.Millisecond,
				},
				DateLocation: gmtLoc,
			}

			got, err := c.Fetch(tt.args.accountID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_FetchBadURL(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	type args struct {
		accountID string
	}

	tests := []struct {
		name        string
		handlerFunc http.HandlerFunc
		args        args
		want        client.Resource
		wantErr     bool
	}{
		{
			name: "returns error if client is misconfigured and a new http request can't be created",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusSeeOther)
				// not used here
			},
			args: args{
				accountID: "uuidv4", // does not matter for this test.
			},
			want:    client.Resource{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.New(
				config.Config{
					AccountsAPIURL: "htt@ps://bla",
					OrganisationID: "orgid",
				},
				http.Client{
					Timeout: testTimeoutMs * time.Millisecond,
				},
				gmtLoc,
			)

			got, err := c.Fetch(tt.args.accountID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			fmt.Printf("%v", err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func returnCompactFile(t *testing.T, filename string) string {
	t.Helper()

	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("could not open file: %s", err)
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read file contents: %s", err)
	}

	var b bytes.Buffer

	err = json.Compact(&b, content)
	if err != nil {
		t.Fatalf("failed to compact json data: %s", err)
	}

	return b.String()
}
