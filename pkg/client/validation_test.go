package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/javorszky/form3takehome/pkg/client"
)

func TestValidateResource(t *testing.T) {
	const (
		bicExample  = "BARCGB22XXX"            // from https://www.iban.com/search-bic
		ibanExample = "GB33BUKB20201555555555" // from https://www.iban.com/structure
	)

	type args struct {
		account client.Resource
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GB is valid when all fields are valid, account number and iban provided",
			args: args{
				account: client.Resource{
					Country:       "GB",
					BankID:        "123456",
					BIC:           bicExample,
					BankIDCode:    "GBDSC",
					AccountNumber: "12345678",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "GB is valid when all fields are valid, account number and iban are not provided",
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankID:     "123456",
					BIC:        bicExample,
					BankIDCode: "GBDSC",
				},
			},
			wantErr: false,
		},
		{
			name: "GB is invalid when BIC is not provided",
			args: args{
				account: client.Resource{
					Country:       "GB",
					BankID:        "123456",
					BankIDCode:    "GBDSC",
					AccountNumber: "12345678",
					IBAN:          ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "GB is invalid when bank id is fewer than 4 digits",
			args: args{
				account: client.Resource{
					Country:       "GB",
					BankID:        "1234",
					BIC:           bicExample,
					BankIDCode:    "GBDSC",
					AccountNumber: "12345678",
					IBAN:          ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "GB is invalid when bank id is 6 letters",
			args: args{
				account: client.Resource{
					Country:       "GB",
					BankID:        "abcdef",
					BIC:           bicExample,
					BankIDCode:    "GBDSC",
					AccountNumber: "12345678",
					IBAN:          ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "GB is invalid when bank id is more than 6 digits",
			args: args{
				account: client.Resource{
					Country:       "GB",
					BankID:        "1234567",
					BIC:           bicExample,
					BankIDCode:    "GBDSC",
					AccountNumber: "12345678",
					IBAN:          ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "GB is invalid when account number is provided, but fewer than 8 digits",
			args: args{
				account: client.Resource{
					Country:       "GB",
					BankID:        "123456",
					BIC:           bicExample,
					BankIDCode:    "GBDSC",
					AccountNumber: "1234567",
					IBAN:          ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "GB is invalid when account number is provided, but more than 8 digits",
			args: args{
				account: client.Resource{
					Country:       "GB",
					BankID:        "123456",
					BIC:           bicExample,
					BankIDCode:    "GBDSC",
					AccountNumber: "123456789",
					IBAN:          ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "GB is invalid when account number is provided, but not all 8 characters are digits",
			args: args{
				account: client.Resource{
					Country:       "GB",
					BankID:        "123456",
					BIC:           bicExample,
					BankIDCode:    "GBDSC",
					AccountNumber: "1234567a",
					IBAN:          ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "GB is invalid when bank id code is not GBDSC",
			args: args{
				account: client.Resource{
					Country:    "GB",
					BankID:     "123456",
					BIC:        bicExample,
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		// AU
		{
			name: "AU is valid when all fields are valid, bank id, account number provided, 6 chars",
			args: args{
				account: client.Resource{
					Country:       "AU",
					BankID:        "123456",
					BIC:           bicExample,
					BankIDCode:    "AUBSB",
					AccountNumber: "123456",
				},
			},
			wantErr: false,
		},
		{
			name: "AU is valid when all fields are valid, bank id not provided, account number provided, 10 chars",
			args: args{
				account: client.Resource{
					Country:       "AU",
					BIC:           bicExample,
					BankIDCode:    "AUBSB",
					AccountNumber: "1234567890",
				},
			},
			wantErr: false,
		},
		{
			name: "AU is valid when all fields are valid, bank id account number not provided",
			args: args{
				account: client.Resource{
					Country:    "AU",
					BIC:        bicExample,
					BankIDCode: "AUBSB",
				},
			},
			wantErr: false,
		},
		{
			name: "AU is invalid when iban is provided and not empty",
			args: args{
				account: client.Resource{
					Country:    "AU",
					BIC:        bicExample,
					BankIDCode: "AUBSB",
					IBAN:       ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "AU is invalid when bank id is provided, but fewer than 6 digits",
			args: args{
				account: client.Resource{
					Country:    "AU",
					BankID:     "12345",
					BIC:        bicExample,
					BankIDCode: "AUBSB",
				},
			},
			wantErr: true,
		},
		{
			name: "AU is invalid when bank id is provided, but more than 6 digits",
			args: args{
				account: client.Resource{
					Country:    "AU",
					BankID:     "1234567",
					BIC:        bicExample,
					BankIDCode: "AUBSB",
				},
			},
			wantErr: true,
		},
		{
			name: "AU is invalid when bank id is provided, 6 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:    "AU",
					BankID:     "12345a",
					BIC:        bicExample,
					BankIDCode: "AUBSB",
				},
			},
			wantErr: true,
		},
		{
			name: "AU is invalid when bank id code is not correct",
			args: args{
				account: client.Resource{
					Country:    "AU",
					BIC:        bicExample,
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "AU is invalid when bic is missing",
			args: args{
				account: client.Resource{
					Country:    "AU",
					BankIDCode: "AUBSB",
				},
			},
			wantErr: true,
		},
		// BE
		{
			name: "BE is valid when all fields are valid, bic, account number provided, 7 chars",
			args: args{
				account: client.Resource{
					Country:       "BE",
					BankID:        "123",
					BIC:           bicExample,
					BankIDCode:    "BE",
					AccountNumber: "1234567",
				},
			},
			wantErr: false,
		},
		{
			name: "BE is valid when all fields are valid, bic, account number not provided",
			args: args{
				account: client.Resource{
					Country:    "BE",
					BankID:     "123",
					BankIDCode: "BE",
				},
			},
			wantErr: false,
		},
		{
			name: "BE is invalid when bank id is fewer than 3 digits",
			args: args{
				account: client.Resource{
					Country:    "BE",
					BankID:     "12",
					BankIDCode: "BE",
				},
			},
			wantErr: true,
		},
		{
			name: "BE is invalid when bank id is more than 3 digits",
			args: args{
				account: client.Resource{
					Country:    "BE",
					BankID:     "1234",
					BankIDCode: "BE",
				},
			},
			wantErr: true,
		},
		{
			name: "BE is invalid when bank id is 3 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:    "BE",
					BankID:     "12a",
					BankIDCode: "BE",
				},
			},
			wantErr: true,
		},
		{
			name: "BE is invalid when bank id code is not present",
			args: args{
				account: client.Resource{
					Country: "BE",
					BankID:  "123",
				},
			},
			wantErr: true,
		},
		{
			name: "BE is invalid when bank id code is not BE",
			args: args{
				account: client.Resource{
					Country:    "BE",
					BankID:     "123",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "BE is invalid when account number is provided, but fewer than 7 digits",
			args: args{
				account: client.Resource{
					Country:       "BE",
					BankID:        "123",
					BankIDCode:    "BE",
					AccountNumber: "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "BE is invalid when account number is provided, but more than 7 digits",
			args: args{
				account: client.Resource{
					Country:       "BE",
					BankID:        "123",
					BankIDCode:    "BE",
					AccountNumber: "12345678",
				},
			},
			wantErr: true,
		},
		{
			name: "BE is invalid when account number is provided, is 7 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:       "BE",
					BankID:        "123",
					BankIDCode:    "BE",
					AccountNumber: "123456a",
				},
			},
			wantErr: true,
		},
		// CA
		{
			name: "CA is valid when all fields are valid, bank id, bank id code, account number provided, 7 chars",
			args: args{
				account: client.Resource{
					Country:       "CA",
					BankID:        "012345678",
					BIC:           bicExample,
					BankIDCode:    "CACPA",
					AccountNumber: "1234567",
				},
			},
			wantErr: false,
		},
		{
			name: "CA is valid when all fields are valid, bank id, bank id code not provided, account number provided, 12 chars",
			args: args{
				account: client.Resource{
					Country:       "CA",
					BIC:           bicExample,
					AccountNumber: "123456789012",
				},
			},
			wantErr: false,
		},
		{
			name: "CA is invalid when iban is provided",
			args: args{
				account: client.Resource{
					Country: "CA",
					BIC:     bicExample,
					IBAN:    ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when bic is not present",
			args: args{
				account: client.Resource{
					Country: "CA",
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when bank id is provided, 9 digits, but starts with something other than 0",
			args: args{
				account: client.Resource{
					Country: "CA",
					BankID:  "123456789",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when bank id is provided, fewer than 9 digits, starts with 0",
			args: args{
				account: client.Resource{
					Country: "CA",
					BankID:  "012",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when bank id is provided, more than 9 digits, starts with 0",
			args: args{
				account: client.Resource{
					Country: "CA",
					BankID:  "0123456789",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when bank id is provided, 9 digits, starts with 0, has non-digit in it",
			args: args{
				account: client.Resource{
					Country: "CA",
					BankID:  "01234567a",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when bank id code is provided, but is not CACPA",
			args: args{
				account: client.Resource{
					Country:    "CA",
					BIC:        bicExample,
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when account number is provided, but fewer than 7 digits",
			args: args{
				account: client.Resource{
					Country:       "CA",
					BIC:           bicExample,
					BankIDCode:    "CACPA",
					AccountNumber: "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when account number is provided, but more than 12 digits",
			args: args{
				account: client.Resource{
					Country:       "CA",
					BIC:           bicExample,
					BankIDCode:    "CACPA",
					AccountNumber: "1234567890123",
				},
			},
			wantErr: true,
		},
		{
			name: "CA is invalid when account number is provided is 10 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:       "CA",
					BIC:           bicExample,
					BankIDCode:    "CACPA",
					AccountNumber: "123456789a",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.ValidateResource(tt.args.account)
			if tt.wantErr {
				assert.Error(t, got)
			} else {
				assert.NoError(t, got)
			}
		})
	}
}
