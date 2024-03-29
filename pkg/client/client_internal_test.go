package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_addHeaders(t *testing.T) {
	const (
		testURL                 = "https://atesturl"
		testJSONBody            = `{data:{key:"value"}}`
		testContentType         = "application/vnd.api+json"
		testHeaderDateThreshold = 15
	)

	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		assert.FailNowf(t, "could not load GMT location", "error: %s", err)
	}

	requestNoBody, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, testURL, nil)
	if err != nil {
		assert.FailNowf(t, "could not create test request with no body", "error: %s", err)
	}

	requestBody, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		testURL,
		strings.NewReader(testJSONBody),
	)
	if err != nil {
		assert.FailNowf(t, "could not create test request with body", "error: %s", err)
	}

	requestEmptyBody, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		testURL,
		strings.NewReader(""),
	)
	if err != nil {
		assert.FailNowf(t, "could not create test request with empty body", "error: %s", err)
	}

	type fields struct {
		BaseURL      string
		DateLocation *time.Location
	}

	type args struct {
		r *http.Request
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantHeaders map[string]string
	}{
		{
			name: "decorates request with headers with no body present",
			fields: fields{
				BaseURL:      testURL,
				DateLocation: gmtLoc,
			},
			args: args{
				r: requestNoBody,
			},
			wantHeaders: map[string]string{
				"Accept": testContentType,
				"Host":   testURL,
			},
		},
		{
			name: "decorates request with headers with empty body present",
			fields: fields{
				BaseURL:      testURL,
				DateLocation: gmtLoc,
			},
			args: args{
				r: requestEmptyBody,
			},
			wantHeaders: map[string]string{
				"Accept": testContentType,
				"Host":   testURL,
			},
		},
		{
			name: "decorates request with headers with body present",
			fields: fields{
				BaseURL:      testURL,
				DateLocation: gmtLoc,
			},
			args: args{
				r: requestBody,
			},
			wantHeaders: map[string]string{
				"Accept":         testContentType,
				"Host":           testURL,
				"Content-Type":   testContentType,
				"Content-Length": strconv.Itoa(len(testJSONBody)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				BaseURL:      tt.fields.BaseURL,
				DateLocation: tt.fields.DateLocation,
			}
			got := c.addHeaders(tt.args.r)
			for k, v := range tt.wantHeaders {
				assert.Equal(t, v, got.Header.Get(k))
			}

			// Check the Date header separately
			headerDate := got.Header.Get("Date")
			if !strings.HasSuffix(headerDate, "GMT") {
				assert.FailNowf(t, "header date should end with GMT. It doesn't", "error: %s", err)
			}
			parsedHeaderDate, err := time.Parse(time.RFC1123, headerDate)
			if err != nil {
				assert.FailNowf(
					t,
					"could not parse header date into a time.Time struct",
					"error: %s",
					err,
				)
			}
			assert.WithinDuration(t, parsedHeaderDate, time.Now(), testHeaderDateThreshold*time.Second)
		})
	}
}

func Test_unmarshalPayload(t *testing.T) {
	f, err := os.Open("./testdata/payload.json")
	if err != nil {
		assert.FailNowf(t, "could not open file", "error: %s", err)
	}

	defer func() {
		_ = f.Close()
	}()

	testTime, err := time.Parse(time.RFC3339, "2020-05-06T09:28:13.843Z")
	if err != nil {
		assert.FailNowf(t, "could not parse test time", "error: %s", err)
	}

	type args struct {
		r io.Reader
	}

	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr bool
	}{
		{
			name: "unmarshals payload json",
			args: args{
				r: f,
			},
			want: Payload{
				Data: Data{
					ID:             "a6c1a721-bb1b-41ef-bd11-800a1309ff9b",
					OrganisationID: "7442ea6b-164a-4818-b470-d98abfbc24ae",
					Type:           "accounts",
					Version:        0,
					CreatedOn:      testTime,
					ModifiedOn:     testTime,
					Attributes: Resource{
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
				Links: Links{
					Self:  "https://selflink.com/resource",
					First: "https://firstlink.com/resource",
					Next:  "https://nextlink.com/resource",
					Last:  "https://lastlink.com/resource",
				},
			},
		},
		{
			name: "returns error on not a json",
			args: args{
				r: strings.NewReader("notajson"),
			},
			want:    Payload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalPayload(tt.args.r)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_unmarshalMultiPayload(t *testing.T) {
	f, err := os.Open("./testdata/multipayload.json")
	if err != nil {
		assert.FailNowf(t, "could not open file", "error: %s", err)
	}

	defer func() {
		_ = f.Close()
	}()

	testTime, err := time.Parse(time.RFC3339, "2020-05-06T09:28:13.843Z")
	if err != nil {
		assert.FailNowf(t, "could not parse test time", "error: %s", err)
	}

	testTime2, err := time.Parse(time.RFC3339, "2020-08-06T09:28:13.843Z")
	if err != nil {
		assert.FailNowf(t, "could not parse test time2", "error: %s", err)
	}

	type args struct {
		r io.Reader
	}

	tests := []struct {
		name    string
		args    args
		want    MultiPayload
		wantErr bool
	}{
		{
			name: "correctly unmarshals multipayload",
			args: args{
				r: f,
			},
			want: MultiPayload{
				Data: []Data{
					{
						ID:             "a6c1a721-bb1b-41ef-bd11-800a1309ff9b",
						OrganisationID: "7442ea6b-164a-4818-b470-d98abfbc24ae",
						Type:           "accounts",
						Version:        0,
						CreatedOn:      testTime,
						ModifiedOn:     testTime,
						Attributes: Resource{
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
						Attributes: Resource{
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
				Links: Links{
					Self:  "https://selflink.com/resource",
					First: "https://firstlink.com/resource",
					Next:  "https://nextlink.com/resource",
					Last:  "https://lastlink.com/resource",
				},
			},
		},
		{
			name: "returns error on unintelligible json",
			args: args{
				r: strings.NewReader("notajson"),
			},
			want:    MultiPayload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalMultiPayload(tt.args.r)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_marshalPayload(t *testing.T) {
	testTime, err := time.Parse(time.RFC3339, "2020-05-06T09:28:13.843Z")
	if err != nil {
		assert.FailNowf(t, "could not parse test time", "error: %s", err)
	}

	f, err := os.Open("./testdata/payload.json")
	if err != nil {
		assert.FailNowf(t, "could not open file", "error: %s", err)
	}

	defer func() {
		_ = f.Close()
	}()

	jsonPayload, err := ioutil.ReadAll(f)
	if err != nil {
		assert.FailNowf(t, "could not read file contents", "error: %s", err)
	}

	var b bytes.Buffer

	err = json.Compact(&b, jsonPayload)
	if err != nil {
		assert.FailNowf(t, "could not compact json data", "error: %s", err)
	}

	type args struct {
		r Payload
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "marshals payload correctly",
			args: args{
				r: Payload{
					Data: Data{
						ID:             "a6c1a721-bb1b-41ef-bd11-800a1309ff9b",
						OrganisationID: "7442ea6b-164a-4818-b470-d98abfbc24ae",
						Type:           "accounts",
						Version:        0,
						CreatedOn:      testTime,
						ModifiedOn:     testTime,
						Attributes: Resource{
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
					Links: Links{
						Self:  "https://selflink.com/resource",
						First: "https://firstlink.com/resource",
						Next:  "https://nextlink.com/resource",
						Last:  "https://lastlink.com/resource",
					},
				},
			},
			want:    b.String() + "\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshalPayload(tt.args.r)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			gotJson, err := ioutil.ReadAll(got)
			if err != nil {
				assert.FailNowf(t, "could not read ioutil.readall contents", "error: %s", err)
			}

			assert.Equal(t, tt.want, string(gotJson))
		})
	}
}
