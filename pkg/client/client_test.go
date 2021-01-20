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

	testTime, err := time.Parse(time.RFC3339, "2020-05-06T09:28:13.843Z")
	if err != nil {
		t.Fatalf("could not parse test time: %s", err)
	}

	type args struct {
		account client.Resource
	}

	tests := []struct {
		name        string
		handlerFunc http.HandlerFunc
		args        args
		want        client.Payload
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
			want: client.Payload{
				Data: client.Data{
					ID:             "a6c1a721-bb1b-41ef-bd11-800a1309ff9b",
					OrganisationID: "7442ea6b-164a-4818-b470-d98abfbc24ae",
					Type:           "accounts",
					Version:        0,
					CreatedOn:      testTime,
					ModifiedOn:     testTime,
					Attributes: client.Resource{
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
				},
				Links: client.Links{
					Self:  "https://selflink.com/resource",
					First: "https://firstlink.com/resource",
					Next:  "https://nextlink.com/resource",
					Last:  "https://lastlink.com/resource",
				},
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
			want:    client.Payload{},
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
			want:    client.Payload{},
			wantErr: true,
		},
		{
			name: "returns error if server responds with non json body",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = fmt.Fprint(w, "not a json")
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Payload{},
			wantErr: true,
		},
		{
			name: "returns error if server responds with json that would result in an empty Data attribute",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = fmt.Fprint(w, `{"error":"a string, but not something that can become a Payload"}`)
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Payload{},
			wantErr: true,
		},
		{
			name: "returns error if server responds with json that would result in an empty Data.Attributes field",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = fmt.Fprint(w, `{"data":{"id": "a string, but not something that can become a Payload"}}`)
			},
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankIDCode: "GBDSC",
					BIC:        "bic",
					BankID:     "123456",
				}, // it doesn't matter what we send to the handlerFunc as long as it passes validation.
			},
			want:    client.Payload{},
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
			want:    client.Payload{},
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
		want        client.Payload
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
			want:    client.Payload{},
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

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_Fetch(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	testTime, err := time.Parse(time.RFC3339, "2020-05-06T09:28:13.843Z")
	if err != nil {
		t.Fatalf("could not parse test time: %s", err)
	}

	type args struct {
		accountID string
	}

	tests := []struct {
		name        string
		handlerFunc http.HandlerFunc
		args        args
		want        client.Payload
		wantErr     bool
	}{
		{
			name: "correctly returns a resource by id",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, returnCompactFile(t, "./testdata/payload.json"))
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want: client.Payload{
				Data: client.Data{
					ID:             "a6c1a721-bb1b-41ef-bd11-800a1309ff9b",
					OrganisationID: "7442ea6b-164a-4818-b470-d98abfbc24ae",
					Type:           "accounts",
					Version:        0,
					CreatedOn:      testTime,
					ModifiedOn:     testTime,
					Attributes: client.Resource{
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
				},
				Links: client.Links{
					Self:  "https://selflink.com/resource",
					First: "https://firstlink.com/resource",
					Next:  "https://nextlink.com/resource",
					Last:  "https://lastlink.com/resource",
				},
			},
			wantErr: false,
		},
		{
			name: "returns error if response takes longer than configured timeout",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep((testTimeoutMs + 100) * time.Millisecond)
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, returnCompactFile(t, "./testdata/payload.json"))
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Payload{},
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
			want:    client.Payload{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but not a json",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, "not a json")
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Payload{},
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
			want:    client.Payload{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but json can't be populated into a meaningful resource",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, `{"error":"not data"}`)
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Payload{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but json can't be populated into a meaningful resource",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, `{"data":{"randomKey": "not data"}}`)
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Payload{},
			wantErr: true,
		},
		{
			name: "returns error if response is 200 code but json can't be populated into a meaningful resource",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, `{"data":{"id": "uuidve-missingattributes"}}`)
			},
			args: args{
				accountID: "uuidv4accountid", // doesn't matter what we pass in here for the time being.
			},
			want:    client.Payload{},
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
		want        client.Payload
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
			want:    client.Payload{},
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

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_Delete(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	type args struct {
		accountID string
		version   uint
	}

	tests := []struct {
		name        string
		handlerFunc http.HandlerFunc
		args        args
		wantErr     bool
	}{
		{
			name: "correctly returns response from request",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			args: args{
				accountID: "uuidv4accountid",
				version:   3,
			}, // does not matter what we pass in for these tests.
			wantErr: false,
		},
		{
			name: "returns error if the response is not a 204 no content",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusConflict)
			},
			args: args{
				accountID: "uuidv4accountid",
				version:   3,
			}, // does not matter what we pass in for these tests.
			wantErr: true,
		},
		{
			name: "returns error if the response takes longer than the timeout",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep((testTimeoutMs + 100) * time.Millisecond)
				w.WriteHeader(http.StatusNoContent)
			},
			args: args{
				accountID: "uuidv4accountid",
				version:   3,
			}, // does not matter what we pass in for these tests.
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

			err := c.Delete(tt.args.accountID, tt.args.version)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_DeleteBadURL(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	type args struct {
		accountID string
		version   uint
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "returns error if the base url is wrongly configured",
			args: args{
				accountID: "uuidv4accountid",
				version:   3,
			}, // does not matter what we pass in for these tests.
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.Client{
				BaseURL:        "htt@ps:bla//",
				OrganisationID: "orgid",
				HttpClient: http.Client{
					Timeout: testTimeoutMs * time.Millisecond,
				},
				DateLocation: gmtLoc,
			}

			err := c.Delete(tt.args.accountID, tt.args.version)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_List(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	testTime, err := time.Parse(time.RFC3339, "2020-05-06T09:28:13.843Z")
	if err != nil {
		t.Fatalf("could not parse test time: %s", err)
	}

	testTime2, err := time.Parse(time.RFC3339, "2020-08-06T09:28:13.843Z")
	if err != nil {
		t.Fatalf("could not parse test time2: %s", err)
	}

	type args struct {
		pageNumber uint
		pageSize   uint
	}

	tests := []struct {
		name        string
		handlerFunc http.HandlerFunc
		args        args
		want        client.MultiPayload
		wantErr     bool
	}{
		{
			name: "correctly returns list of resources",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, returnCompactFile(t, "./testdata/multipayload.json"))
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want: client.MultiPayload{
				Data: []client.Data{
					{
						ID:             "a6c1a721-bb1b-41ef-bd11-800a1309ff9b",
						OrganisationID: "7442ea6b-164a-4818-b470-d98abfbc24ae",
						Type:           "accounts",
						Version:        0,
						CreatedOn:      testTime,
						ModifiedOn:     testTime,
						Attributes: client.Resource{
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
					},
					{
						ID:             "ffa7706b-d8fc-40b2-be6b-67d2a628cadf",
						OrganisationID: "7442ea6b-164a-4818-b470-d98abfbc24ae",
						Type:           "accounts",
						Version:        0,
						CreatedOn:      testTime2,
						ModifiedOn:     testTime2,
						Attributes: client.Resource{
							Country:       "GB",
							BaseCurrency:  "GBP",
							BankID:        "89282dd",
							BankIDCode:    "999999",
							AccountNumber: "87654321",
							BIC:           "bic5678",
							IBAN:          "iban5678",
							CustomerID:    "anuuidv4again",
							Name: [4]string{
								"line1-2",
								"line2-2",
								"line3-2",
								"line4-2",
							},
							AlternativeNames: [3]string{
								"altname1-2",
								"altname2-2",
								"altname3-2",
							},
							AccountClassification:   "cop",
							JointAccount:            true,
							AccountMatchingOptOut:   true,
							SecondaryIdentification: "another custom name",
							Switched:                true,
							Status:                  "confirmed",
						},
					},
				},
				Links: client.Links{
					Self:  "https://selflink.com/resource",
					First: "https://firstlink.com/resource",
					Next:  "https://nextlink.com/resource",
					Last:  "https://lastlink.com/resource",
				},
			},
			wantErr: false,
		},
		{
			name: "returns error if the response code is not 200",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
			wantErr: true,
		},
		{
			name: "returns error if the response takes longer than timeout",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep((testTimeoutMs + 100) * time.Millisecond)
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, returnCompactFile(t, "./testdata/multipayload.json"))
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
			wantErr: true,
		},
		{
			name: "returns error if the response is not a json",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, "not a json")
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
			wantErr: true,
		},
		{
			name: "returns error if the response is json, but can't be unmarshaled into a multipayload (no data key)",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, `{"error":"not payload"}`)
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
			wantErr: true,
		},
		{
			name: "returns error if the response is json, but can't be unmarshaled into a multipayload (data is not array)",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, `{"data":"not a json array"}`)
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
			wantErr: true,
		},
		{
			name: "error when response can't be unmarshaled into multipayload (data is not array of objects)",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, `{"data":["not an object"]}`)
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
			wantErr: true,
		},
		{
			name: "error when response can't be unmarshaled into multipayload (data objects emtpy)",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, `{"data":[{"randomkey":"notdata"}]}`)
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
			wantErr: true,
		},
		{
			name: "error when response can't be unmarshaled into multipayload (data objects missing Attributes)",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, `{"data":[{"id":"no attributes yet"}]}`)
			},
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
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

			got, err := c.List(tt.args.pageNumber, tt.args.pageSize)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_ListBadURL(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	type args struct {
		pageNumber uint
		pageSize   uint
	}

	tests := []struct {
		name    string
		args    args
		want    client.MultiPayload
		wantErr bool
	}{
		{
			name: "error when request can't be constructed due to bad base url",
			args: args{
				pageNumber: 1,
				pageSize:   2,
			}, // does not matter what these are.
			want:    client.MultiPayload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.Client{
				BaseURL:        "htt@ps://bla",
				OrganisationID: "orgid",
				HttpClient: http.Client{
					Timeout: testTimeoutMs * time.Millisecond,
				},
				DateLocation: gmtLoc,
			}

			got, err := c.List(tt.args.pageNumber, tt.args.pageSize)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

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

	defer func() {
		_ = f.Close()
	}()

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
