package client_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/javorszky/form3takehome/pkg/client"
	"github.com/javorszky/form3takehome/pkg/config"
)

const (
	integrationTestURL   = "http://accountapi:8080"
	integrationTestOrgID = "0e1445e5-2047-4a98-ad4d-55068b25359a" // generated by https://www.uuidgenerator.net/version4
	bicExample           = "BARCGB22XXX"                          // from https://www.iban.com/search-bic
	ibanExample          = "GB33BUKB20201555555555"               // from https://www.iban.com/structure
)

func TestClient_IntegrationCreateFetchListDelete(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load gmt location: %s", err)
	}

	type args struct {
		accounts []client.Resource
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "correctly creates, fetches, lists, and deletes resources against the service in docker",
			args: args{
				accounts: []client.Resource{
					{
						Country:    "GB",
						BankIDCode: "GBDSC",
						BIC:        bicExample,
						BankID:     "123456",
					},
					{
						Country:       "AU",
						BankID:        "123456",
						BIC:           bicExample,
						BankIDCode:    "AUBSB",
						AccountNumber: "123456",
					},
					{
						Country:       "BE",
						BankID:        "123",
						BIC:           bicExample,
						BankIDCode:    "BE",
						AccountNumber: "1234567",
					},
					{
						Country:       "CA",
						BankID:        "012345678",
						BIC:           bicExample,
						BankIDCode:    "CACPA",
						AccountNumber: "1234567",
					},
					{
						Country:       "FR",
						BankID:        "1234567890",
						BIC:           bicExample,
						BankIDCode:    "FR",
						AccountNumber: "1234567890",
						IBAN:          ibanExample,
					},
					{
						Country:       "DE",
						BankID:        "12345678",
						BIC:           bicExample,
						BankIDCode:    "DEBLZ",
						AccountNumber: "1234567",
						IBAN:          ibanExample,
					},
					{
						Country:       "GR",
						BankID:        "1234567",
						BIC:           bicExample,
						BankIDCode:    "GRBIC",
						AccountNumber: "1234567890123456",
						IBAN:          ibanExample,
					},
					{
						Country:       "HK",
						BankID:        "123",
						BIC:           bicExample,
						BankIDCode:    "HKNCC",
						AccountNumber: "123456789",
					},
					{
						Country:       "IT",
						BankID:        "12345678901",
						BIC:           bicExample,
						BankIDCode:    "ITNCC",
						AccountNumber: "123456789012",
						IBAN:          ibanExample,
					},
					{
						Country:       "LU",
						BankID:        "123",
						BIC:           bicExample,
						BankIDCode:    "LULUX",
						AccountNumber: "1234567890123",
						IBAN:          ibanExample,
					},
					{
						Country:       "NL",
						BIC:           bicExample,
						AccountNumber: "1234567890",
						IBAN:          ibanExample,
					},
					{
						Country:       "PL",
						BankID:        "12345678",
						BIC:           bicExample,
						BankIDCode:    "PLKNR",
						AccountNumber: "1234567890123456",
						IBAN:          ibanExample,
					},
					{
						Country:       "PT",
						BankID:        "12345678",
						BIC:           bicExample,
						BankIDCode:    "PTNCC",
						AccountNumber: "12345678901",
						IBAN:          ibanExample,
					},
					{
						Country:       "ES",
						BankID:        "12345678",
						BIC:           bicExample,
						BankIDCode:    "ESNCC",
						AccountNumber: "1234567890",
						IBAN:          ibanExample,
					},
					{
						Country:       "CH",
						BankID:        "12345",
						BIC:           bicExample,
						BankIDCode:    "CHBCC",
						AccountNumber: "123456789012",
						IBAN:          ibanExample,
					},
					{
						Country:       "US",
						BankID:        "123456789",
						BIC:           bicExample,
						BankIDCode:    "USABA",
						AccountNumber: "123456",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.New(
				config.Config{
					AccountsAPIURL: integrationTestURL,
					OrganisationID: integrationTestOrgID,
				},
				http.Client{
					Timeout: testTimeoutMs * time.Millisecond,
				},
				gmtLoc,
			)

			payloadsHelper := make([]client.Payload, 0)

			// First let's store all of the payloads one by one. Every supported country is present, and the Resources
			// have the same data as the ones I used for the validation tests. All of these should be created.
			for _, r := range tt.args.accounts {
				got, err := c.Create(r)
				if err != nil {
					assert.FailNowf(t, "create encountered an error", "resource %#v: %s", r, err)
				}
				payloadsHelper = append(payloadsHelper, got)
			}

			// Then let's fetch them one by one to make sure that they are actually present in the service and compare/
			// with what we have.
			for _, stored := range payloadsHelper {
				got, err := c.Fetch(stored.Data.ID)
				if err != nil {
					assert.FailNowf(t,
						"fetching resource encountered an error",
						"resource %#v with id %s: %s",
						stored.Data.Attributes,
						stored.Data.ID, err,
					)
				}
				assert.Equal(t, stored, got)
			}

			// Then let's list them, and compare them with the payloadsHelper slice
			l, err := c.List(0, 100)
			if err != nil {
				assert.FailNowf(t, "list encountered an error", "error message: %s", err)
			}

			// The list should be the same length as the payloadshelper. If not, we're either bleeding data, or
			// something is wrong in our code.
			assert.Equal(t, len(l.Data), len(payloadsHelper))

			listHelper := make(map[string]client.Data)

			for _, listItem := range l.Data {
				listHelper[listItem.ID] = listItem
			}

			// Check that all data that we put in the service also came out. This, and the same length check earlier
			// will mean that all of them are accounted for.
			for _, payloadItem := range payloadsHelper {
				_, ok := listHelper[payloadItem.Data.ID]
				assert.Truef(t, ok, "could not find account with id %s", payloadItem.Data.ID)
			}

			// now delete all of them
			for _, payloadItemToDelete := range payloadsHelper {
				c.Delete(payloadItemToDelete.Data.ID, uint(payloadItemToDelete.Data.Version))
			}

			// and check that they are indeed missing in two different ways
			for _, payloadItemToCheckAfterDelete := range payloadsHelper {
				_, errChecked := c.Fetch(payloadItemToCheckAfterDelete.Data.ID)
				assert.Errorf(
					t,
					errChecked,
					"fetch should have returned an error for account with id %s",
					payloadItemToCheckAfterDelete.Data.ID,
				)
			}

			// and with list
			deletedList, err := c.List(0, 100)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(deletedList.Data))
		})
	}
}
