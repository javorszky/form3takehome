package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/javorszky/form3takehome/pkg/client"
)

func TestValidateResource(t *testing.T) {
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
		{
			name: "AU is invalid when account number is fewer than 6 digits long",
			args: args{
				account: client.Resource{
					Country:       "AU",
					BIC:           bicExample,
					BankIDCode:    "AUBSB",
					AccountNumber: "12345",
				},
			},
			wantErr: true,
		},
		{
			name: "AU is invalid when account number is longer than 10 digits",
			args: args{
				account: client.Resource{
					Country:       "AU",
					BIC:           bicExample,
					BankIDCode:    "AUBSB",
					AccountNumber: "12345678901",
				},
			},
			wantErr: true,
		},
		{
			name: "AU is invalid when account number is correct length, but starts with 0",
			args: args{
				account: client.Resource{
					Country:       "AU",
					BIC:           bicExample,
					BankIDCode:    "AUBSB",
					AccountNumber: "01234567",
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
		// FR
		{
			name: "FR is valid when all fields are valid, bic, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "FR",
					BankID:        "1234567890",
					BIC:           bicExample,
					BankIDCode:    "FR",
					AccountNumber: "1234567890",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "FR is valid when all fields are valid, bic, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country:    "FR",
					BankID:     "1234567890",
					BankIDCode: "FR",
				},
			},
			wantErr: false,
		},
		{
			name: "FR is invalid when bank id is fewer than 10 digits",
			args: args{
				account: client.Resource{
					Country:    "FR",
					BankID:     "123456789",
					BankIDCode: "FR",
				},
			},
			wantErr: true,
		},
		{
			name: "FR is invalid when bank id is more than 10 digits",
			args: args{
				account: client.Resource{
					Country:    "FR",
					BankID:     "12345678901",
					BankIDCode: "FR",
				},
			},
			wantErr: true,
		},
		{
			name: "FR is invalid when bank id is 10 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:    "FR",
					BankID:     "123456789a",
					BankIDCode: "FR",
				},
			},
			wantErr: true,
		},
		{
			name: "FR is invalid when bank id code is not FR",
			args: args{
				account: client.Resource{
					Country:    "FR",
					BankID:     "1234567890",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "FR is invalid when bank id missing",
			args: args{
				account: client.Resource{
					Country:    "FR",
					BankIDCode: "FR",
				},
			},
			wantErr: true,
		},
		{
			name: "FR is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "FR",
					BankID:  "1234567890",
				},
			},
			wantErr: true,
		},
		{
			name: "FR is invalid when account is provided, but is fewer than 10 digits",
			args: args{
				account: client.Resource{
					Country:       "FR",
					BankID:        "1234567890",
					BankIDCode:    "FR",
					AccountNumber: "123456789",
				},
			},
			wantErr: true,
		},
		{
			name: "FR is invalid when account is provided, but is more than 10 digits",
			args: args{
				account: client.Resource{
					Country:       "FR",
					BankID:        "1234567890",
					BankIDCode:    "FR",
					AccountNumber: "12345678901",
				},
			},
			wantErr: true,
		},
		{
			name: "FR is invalid when account is provided, is 10 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:       "FR",
					BankID:        "1234567890",
					BankIDCode:    "FR",
					AccountNumber: "123456789a",
				},
			},
			wantErr: true,
		},
		// DE
		{
			name: "DE is valid when all fields are valid, bic, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "DE",
					BankID:        "12345678",
					BIC:           bicExample,
					BankIDCode:    "DEBLZ",
					AccountNumber: "1234567",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "DE is valid when all fields are valid, bic, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country:    "DE",
					BankID:     "12345678",
					BankIDCode: "DEBLZ",
				},
			},
			wantErr: false,
		},
		{
			name: "DE is invalid when bank id is fewer than 8 digits",
			args: args{
				account: client.Resource{
					Country:    "DE",
					BankID:     "1234567",
					BankIDCode: "DEBLZ",
				},
			},
			wantErr: true,
		},
		{
			name: "DE is invalid when bank id is more than 8 digits",
			args: args{
				account: client.Resource{
					Country:    "DE",
					BankID:     "123456789",
					BankIDCode: "DEBLZ",
				},
			},
			wantErr: true,
		},
		{
			name: "DE is invalid when bank id is 8 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:    "DE",
					BankID:     "1234567a",
					BankIDCode: "DEBLZ",
				},
			},
			wantErr: true,
		},
		{
			name: "DE is invalid when bank id is missing",
			args: args{
				account: client.Resource{
					Country:    "DE",
					BankIDCode: "DEBLZ",
				},
			},
			wantErr: true,
		},
		{
			name: "DE is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "DE",
					BankID:  "12345678",
				},
			},
			wantErr: true,
		},
		{
			name: "DE is invalid when bank id code is not the correct value",
			args: args{
				account: client.Resource{
					Country:    "DE",
					BankID:     "12345678",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "DE is invalid when account is provided, but fewer than 7 digits",
			args: args{
				account: client.Resource{
					Country:       "DE",
					BankID:        "12345678",
					BankIDCode:    "DEBLZ",
					AccountNumber: "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "DE is invalid when account is provided, but more than 7 digits",
			args: args{
				account: client.Resource{
					Country:       "DE",
					BankID:        "12345678",
					BankIDCode:    "DEBLZ",
					AccountNumber: "12345678",
				},
			},
			wantErr: true,
		},
		{
			name: "DE is invalid when account is provided, is 7 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:       "DE",
					BankID:        "12345678",
					BankIDCode:    "DEBLZ",
					AccountNumber: "123456a",
				},
			},
			wantErr: true,
		},
		// GR
		{
			name: "GR is valid when all fields are valid, bic, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "GR",
					BankID:        "1234567",
					BIC:           bicExample,
					BankIDCode:    "GRBIC",
					AccountNumber: "1234567890123456",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "GR is valid when all fields are valid, bic, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country:    "GR",
					BankID:     "1234567",
					BankIDCode: "GRBIC",
				},
			},
			wantErr: false,
		},
		{
			name: "GR is invalid when bank id is missing",
			args: args{
				account: client.Resource{
					Country:    "GR",
					BankIDCode: "GRBIC",
				},
			},
			wantErr: true,
		},
		{
			name: "GR is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "GR",
					BankID:  "1234567",
				},
			},
			wantErr: true,
		},
		{
			name: "GR is invalid when bank id is fewer than 7 digits",
			args: args{
				account: client.Resource{
					Country:    "GR",
					BankID:     "123456",
					BankIDCode: "GRBIC",
				},
			},
			wantErr: true,
		},
		{
			name: "GR is invalid when bank id is more than 7 digits",
			args: args{
				account: client.Resource{
					Country:    "GR",
					BankID:     "12345678",
					BankIDCode: "GRBIC",
				},
			},
			wantErr: true,
		},
		{
			name: "GR is invalid when bank id is 7 characters, not all of them digits",
			args: args{
				account: client.Resource{
					Country:    "GR",
					BankID:     "123456a",
					BankIDCode: "GRBIC",
				},
			},
			wantErr: true,
		},
		{
			name: "GR is invalid when bank id code is not correct value",
			args: args{
				account: client.Resource{
					Country:    "GR",
					BankID:     "1234567",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "GR is invalid when account number is provided, but fewer than 16 digits",
			args: args{
				account: client.Resource{
					Country:       "GR",
					BankID:        "1234567",
					BankIDCode:    "GRBIC",
					AccountNumber: "123456789012345",
				},
			},
			wantErr: true,
		},
		{
			name: "GR is invalid when account number is provided, but more than 16 digits",
			args: args{
				account: client.Resource{
					Country:       "GR",
					BankID:        "1234567",
					BankIDCode:    "GRBIC",
					AccountNumber: "12345678901234567",
				},
			},
			wantErr: true,
		},
		{
			name: "GR is invalid when account number is provided, is 16 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:       "GR",
					BankID:        "1234567",
					BankIDCode:    "GRBIC",
					AccountNumber: "123456789012345a",
				},
			},
			wantErr: true,
		},
		// HK
		{
			name: "HK is valid when all fields are valid, bank id, bank id code, account number provided 9 chars",
			args: args{
				account: client.Resource{
					Country:       "HK",
					BankID:        "123",
					BIC:           bicExample,
					BankIDCode:    "HKNCC",
					AccountNumber: "123456789",
				},
			},
			wantErr: false,
		},
		{
			name: "HK is valid when all fields are valid, bank id, bank id code not provided, account number provided 12 chars",
			args: args{
				account: client.Resource{
					Country:       "HK",
					BIC:           bicExample,
					AccountNumber: "123456789012",
				},
			},
			wantErr: false,
		},
		{
			name: "HK is valid when all fields are valid, bank id, bank id code, account number not provided",
			args: args{
				account: client.Resource{
					Country: "HK",
					BIC:     bicExample,
				},
			},
			wantErr: false,
		},
		{
			name: "HK is invalid when iban is provided",
			args: args{
				account: client.Resource{
					Country: "HK",
					BIC:     bicExample,
					IBAN:    ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "HK is invalid when bic is missing",
			args: args{
				account: client.Resource{
					Country: "HK",
				},
			},
			wantErr: true,
		},
		{
			name: "HK is invalid when bank id is provided, but is fewer than 3 digits",
			args: args{
				account: client.Resource{
					Country: "HK",
					BankID:  "12",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "HK is invalid when bank id is provided, but more than 3 digits",
			args: args{
				account: client.Resource{
					Country: "HK",
					BankID:  "1234",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "HK is invalid when bank id is provided, 3 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country: "HK",
					BankID:  "12a",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "HK is invalid when bank id code is provided, but is not the required value",
			args: args{
				account: client.Resource{
					Country:    "HK",
					BIC:        bicExample,
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "HK is invalid when account number is provided, but fewer than 9 digits",
			args: args{
				account: client.Resource{
					Country:       "HK",
					BIC:           bicExample,
					AccountNumber: "12345678",
				},
			},
			wantErr: true,
		},
		{
			name: "HK is invalid when account number is provided, but more than 12 digits",
			args: args{
				account: client.Resource{
					Country:       "HK",
					BIC:           bicExample,
					AccountNumber: "1234567890123",
				},
			},
			wantErr: true,
		},
		{
			name: "HK is invalid when account number is provided, is 10 characters, but not all of them digits",
			args: args{
				account: client.Resource{
					Country:       "HK",
					BIC:           bicExample,
					AccountNumber: "123456789a",
				},
			},
			wantErr: true,
		},
		// IT
		{
			name: "IT is valid when all fields are valid, bic, account number, iban are provided",
			args: args{
				account: client.Resource{
					Country:       "IT",
					BankID:        "12345678901",
					BIC:           bicExample,
					BankIDCode:    "ITNCC",
					AccountNumber: "123456789012",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "IT is valid when all fields are valid, bic, account number, iban are not provided",
			args: args{
				account: client.Resource{
					Country:    "IT",
					BankID:     "1234567890",
					BankIDCode: "ITNCC",
				},
			},
			wantErr: false,
		},
		{
			name: "IT is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "IT",
					BankID:  "1234567890",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when bank id is missing",
			args: args{
				account: client.Resource{
					Country:    "IT",
					BankIDCode: "ITNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when bank id code is not correct value",
			args: args{
				account: client.Resource{
					Country:    "IT",
					BankID:     "1234567890",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is not provided, but bank id is fewer than 10 digits",
			args: args{
				account: client.Resource{
					Country:    "IT",
					BankID:     "123456789",
					BankIDCode: "ITNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is not provided, but bank id is more than 10 digits",
			args: args{
				account: client.Resource{
					Country:    "IT",
					BankID:     "12345678901",
					BankIDCode: "ITNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is not provided, bank id is 10 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:    "IT",
					BankID:     "123456789a",
					BankIDCode: "ITNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is provided, but bank id is fewer than 11 digits",
			args: args{
				account: client.Resource{
					Country:       "IT",
					BankID:        "1234567890",
					BankIDCode:    "ITNCC",
					AccountNumber: "123456789012",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is provided, but bank id is more than 11 digits",
			args: args{
				account: client.Resource{
					Country:       "IT",
					BankID:        "123456789012",
					BankIDCode:    "ITNCC",
					AccountNumber: "123456789012",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is provided, but bank id is 11 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:       "IT",
					BankID:        "1234567890a",
					BankIDCode:    "ITNCC",
					AccountNumber: "123456789012",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is provided, but it's fewer than 12 digits",
			args: args{
				account: client.Resource{
					Country:       "IT",
					BankID:        "12345678901",
					BankIDCode:    "ITNCC",
					AccountNumber: "12345678901",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is provided, but it's more than 12 digits",
			args: args{
				account: client.Resource{
					Country:       "IT",
					BankID:        "12345678901",
					BankIDCode:    "ITNCC",
					AccountNumber: "1234567890123",
				},
			},
			wantErr: true,
		},
		{
			name: "IT is invalid when account number is provided, it's 12 characters, but not all of them digits",
			args: args{
				account: client.Resource{
					Country:       "IT",
					BankID:        "12345678901",
					BankIDCode:    "ITNCC",
					AccountNumber: "12345678901a",
				},
			},
			wantErr: true,
		},
		// LU
		{
			name: "LU is valid when all fields are valid, bic, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "LU",
					BankID:        "123",
					BIC:           bicExample,
					BankIDCode:    "LULUX",
					AccountNumber: "1234567890123",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "LU is valid when all fields are valid, bic, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country:    "LU",
					BankID:     "123",
					BankIDCode: "LULUX",
				},
			},
			wantErr: false,
		},
		{
			name: "LU is invalid when bank id is missing",
			args: args{
				account: client.Resource{
					Country:    "LU",
					BankIDCode: "LULUX",
				},
			},
			wantErr: true,
		},
		{
			name: "LU is invalid when bank id is fewer than 3 digits",
			args: args{
				account: client.Resource{
					Country:    "LU",
					BankID:     "12",
					BankIDCode: "LULUX",
				},
			},
			wantErr: true,
		},
		{
			name: "LU is invalid when bank id is more than 3 digits",
			args: args{
				account: client.Resource{
					Country:    "LU",
					BankID:     "1234",
					BankIDCode: "LULUX",
				},
			},
			wantErr: true,
		},
		{
			name: "LU is invalid when bank id is 3 characters but not all digits",
			args: args{
				account: client.Resource{
					Country:    "LU",
					BankID:     "12a",
					BankIDCode: "LULUX",
				},
			},
			wantErr: true,
		},
		{
			name: "LU is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "LU",
					BankID:  "123",
				},
			},
			wantErr: true,
		},
		{
			name: "LU is invalid when bank id code is present but wrong value",
			args: args{
				account: client.Resource{
					Country:    "LU",
					BankID:     "123",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "LU is invalid when account number is present, but fewer than 13 digits",
			args: args{
				account: client.Resource{
					Country:       "LU",
					BankID:        "123",
					BankIDCode:    "LULUX",
					AccountNumber: "123456789012",
				},
			},
			wantErr: true,
		},
		{
			name: "LU is invalid when account number is present, but more than 13 digits",
			args: args{
				account: client.Resource{
					Country:       "LU",
					BankID:        "123",
					BankIDCode:    "LULUX",
					AccountNumber: "12345678901234",
				},
			},
			wantErr: true,
		},
		{
			name: "LU is invalid when account number is present, is 13 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:       "LU",
					BankID:        "123",
					BankIDCode:    "LULUX",
					AccountNumber: "123456789012a",
				},
			},
			wantErr: true,
		},
		// NL
		{
			name: "NL is valid when all fields are valid, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "NL",
					BIC:           bicExample,
					AccountNumber: "1234567890",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "NL is valid when all fields are valid, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country: "NL",
					BIC:     bicExample,
				},
			},
			wantErr: false,
		},
		{
			name: "NL is invalid when bank id is present and not empty",
			args: args{
				account: client.Resource{
					Country: "NL",
					BankID:  "NL",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "NL is invalid when bic is not present",
			args: args{
				account: client.Resource{
					Country: "NL",
				},
			},
			wantErr: true,
		},
		{
			name: "NL is invalid when bic is empty",
			args: args{
				account: client.Resource{
					Country: "NL",
					BIC:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "NL is invalid when bank id code is present and not empty",
			args: args{
				account: client.Resource{
					Country:    "NL",
					BIC:        bicExample,
					BankIDCode: "IDCODE",
				},
			},
			wantErr: true,
		},
		{
			name: "NL is invalid when account number is provided, but fewer than 10 digits",
			args: args{
				account: client.Resource{
					Country:       "NL",
					BIC:           bicExample,
					AccountNumber: "123456789",
				},
			},
			wantErr: true,
		},
		{
			name: "NL is invalid when account number is provided, but more than 10 digits",
			args: args{
				account: client.Resource{
					Country:       "NL",
					BIC:           bicExample,
					AccountNumber: "12345678901",
				},
			},
			wantErr: true,
		},
		{
			name: "NL is invalid when account number is provided, 10 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:       "NL",
					BIC:           bicExample,
					AccountNumber: "123456789a",
				},
			},
			wantErr: true,
		},
		// PL
		{
			name: "PL is valid when all fields are valid, bic, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "PL",
					BankID:        "12345678",
					BIC:           bicExample,
					BankIDCode:    "PLKNR",
					AccountNumber: "1234567890123456",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "PL is valid when all fields are valid, bic, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country:    "PL",
					BankID:     "12345678",
					BankIDCode: "PLKNR",
				},
			},
			wantErr: false,
		},
		{
			name: "PL is invalid when bank id is not present",
			args: args{
				account: client.Resource{
					Country:    "PL",
					BankIDCode: "PLKNR",
				},
			},
			wantErr: true,
		},
		{
			name: "PL is invalid when bank id is fewer than 8 digits",
			args: args{
				account: client.Resource{
					Country:    "PL",
					BankID:     "1234567",
					BankIDCode: "PLKNR",
				},
			},
			wantErr: true,
		},
		{
			name: "PL is invalid when bank id is more than 8 digits",
			args: args{
				account: client.Resource{
					Country:    "PL",
					BankID:     "123456789",
					BankIDCode: "PLKNR",
				},
			},
			wantErr: true,
		},
		{
			name: "PL is invalid when bank id is 8 characters, but not all of them digits",
			args: args{
				account: client.Resource{
					Country:    "PL",
					BankID:     "1234567a",
					BankIDCode: "PLKNR",
				},
			},
			wantErr: true,
		},
		{
			name: "PL is invalid when bank id code is not present",
			args: args{
				account: client.Resource{
					Country: "PL",
					BankID:  "12345678",
				},
			},
			wantErr: true,
		},
		{
			name: "PL is invalid when bank id code is not the correct value",
			args: args{
				account: client.Resource{
					Country:    "PL",
					BankID:     "12345678",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "PL is invalid when account number is provided, but fewer than 16 digits",
			args: args{
				account: client.Resource{
					Country:       "PL",
					BankID:        "12345678",
					BankIDCode:    "PLKNR",
					AccountNumber: "123456789012345",
				},
			},
			wantErr: true,
		},
		{
			name: "PL is invalid when account number is provided, but more than 16 digits",
			args: args{
				account: client.Resource{
					Country:       "PL",
					BankID:        "12345678",
					BankIDCode:    "PLKNR",
					AccountNumber: "12345678901234567",
				},
			},
			wantErr: true,
		},
		{
			name: "PL is invalid when account number is provided, is 16 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:       "PL",
					BankID:        "12345678",
					BankIDCode:    "PLKNR",
					AccountNumber: "123456789012345a",
				},
			},
			wantErr: true,
		},
		// PT
		{
			name: "PT is valid when all fields are valid, bic, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "PT",
					BankID:        "12345678",
					BIC:           bicExample,
					BankIDCode:    "PTNCC",
					AccountNumber: "12345678901",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "PT is valid when all fields are valid, bic, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country:    "PT",
					BankID:     "12345678",
					BankIDCode: "PTNCC",
				},
			},
			wantErr: false,
		},
		{
			name: "PT is invalid when bank id is missing",
			args: args{
				account: client.Resource{
					Country:    "PT",
					BankIDCode: "PTNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "PT is invalid when bank id is fewer than 8 digits",
			args: args{
				account: client.Resource{
					Country:    "PT",
					BankID:     "1234567",
					BankIDCode: "PTNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "PT is invalid when bank id is more than 8 digits",
			args: args{
				account: client.Resource{
					Country:    "PT",
					BankID:     "123456789",
					BankIDCode: "PTNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "PT is invalid when bank id is 8 characters, but not all digits",
			args: args{
				account: client.Resource{
					Country:    "PT",
					BankID:     "1234567a",
					BankIDCode: "PTNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "PT is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "PT",
					BankID:  "12345678",
				},
			},
			wantErr: true,
		},
		{
			name: "PT is invalid when bank id code is wrong value",
			args: args{
				account: client.Resource{
					Country:    "PT",
					BankID:     "12345678",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "PT is invalid when account number is fewer than 11 digits",
			args: args{
				account: client.Resource{
					Country:       "PT",
					BankID:        "12345678",
					BankIDCode:    "PTNCC",
					AccountNumber: "1234567890",
				},
			},
			wantErr: true,
		},
		{
			name: "PT is invalid when account number is more than 11 digits",
			args: args{
				account: client.Resource{
					Country:       "PT",
					BankID:        "12345678",
					BankIDCode:    "PTNCC",
					AccountNumber: "123456789012",
				},
			},
			wantErr: true,
		},
		{
			name: "PT is invalid when account number is 11 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:       "PT",
					BankID:        "12345678",
					BankIDCode:    "PTNCC",
					AccountNumber: "1234567890a",
				},
			},
			wantErr: true,
		},
		// ES
		{
			name: "ES is valid when all fields are valid, bic, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "ES",
					BankID:        "12345678",
					BIC:           bicExample,
					BankIDCode:    "ESNCC",
					AccountNumber: "1234567890",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "ES is valid when all fields are valid, bic, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country:    "ES",
					BankID:     "12345678",
					BankIDCode: "ESNCC",
				},
			},
			wantErr: false,
		},
		{
			name: "ES is invalid when bank id is missing",
			args: args{
				account: client.Resource{
					Country:    "ES",
					BankIDCode: "ESNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "ES is invalid when bank id is fewer than 8 digits",
			args: args{
				account: client.Resource{
					Country:    "ES",
					BankID:     "1234567",
					BankIDCode: "ESNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "ES is invalid when bank id is more than 8 digits",
			args: args{
				account: client.Resource{
					Country:    "ES",
					BankID:     "123456789",
					BankIDCode: "ESNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "ES is invalid when bank id is 8 characters, but not all are digits",
			args: args{
				account: client.Resource{
					Country:    "ES",
					BankID:     "1234567a",
					BankIDCode: "ESNCC",
				},
			},
			wantErr: true,
		},
		{
			name: "ES is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "ES",
					BankID:  "12345678",
				},
			},
			wantErr: true,
		},
		{
			name: "ES is invalid when bank id code is wrong value",
			args: args{
				account: client.Resource{
					Country:    "ES",
					BankID:     "12345678",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "ES is invalid when account number is fewer than 10 digits",
			args: args{
				account: client.Resource{
					Country:       "ES",
					BankID:        "12345678",
					BankIDCode:    "ESNCC",
					AccountNumber: "123456789",
				},
			},
			wantErr: true,
		},
		{
			name: "ES is invalid when account number is more than 10 digits",
			args: args{
				account: client.Resource{
					Country:       "ES",
					BankID:        "12345678",
					BankIDCode:    "ESNCC",
					AccountNumber: "12345678901",
				},
			},
			wantErr: true,
		},
		{
			name: "ES is invalid when account number is 10 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:       "ES",
					BankID:        "12345678",
					BankIDCode:    "ESNCC",
					AccountNumber: "123456789a",
				},
			},
			wantErr: true,
		},
		// CH
		{
			name: "CH is valid when all fields are valid, bic, account number, iban provided",
			args: args{
				account: client.Resource{
					Country:       "CH",
					BankID:        "12345",
					BIC:           bicExample,
					BankIDCode:    "CHBCC",
					AccountNumber: "123456789012",
					IBAN:          ibanExample,
				},
			},
			wantErr: false,
		},
		{
			name: "CH is valid when all fields are valid, bic, account number, iban not provided",
			args: args{
				account: client.Resource{
					Country:    "CH",
					BankID:     "12345",
					BankIDCode: "CHBCC",
				},
			},
			wantErr: false,
		},
		{
			name: "CH is invalid when bank id is missing",
			args: args{
				account: client.Resource{
					Country:    "CH",
					BankIDCode: "CHBCC",
				},
			},
			wantErr: true,
		},
		{
			name: "CH is invalid when bank id is fewer than 5 digits",
			args: args{
				account: client.Resource{
					Country:    "CH",
					BankID:     "1234",
					BankIDCode: "CHBCC",
				},
			},
			wantErr: true,
		},
		{
			name: "CH is invalid when bank id more than 5 digits",
			args: args{
				account: client.Resource{
					Country:    "CH",
					BankID:     "123456",
					BankIDCode: "CHBCC",
				},
			},
			wantErr: true,
		},
		{
			name: "CH is invalid when bank id is 5 characters, not all of them digits",
			args: args{
				account: client.Resource{
					Country:    "CH",
					BankID:     "1234a",
					BankIDCode: "CHBCC",
				},
			},
			wantErr: true,
		},
		{
			name: "CH is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "CH",
					BankID:  "12345",
				},
			},
			wantErr: true,
		},
		{
			name: "CH is invalid when bank id code is wrong value",
			args: args{
				account: client.Resource{
					Country:    "CH",
					BankID:     "12345",
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "CH is invalid when account number is fewer than 12 digits",
			args: args{
				account: client.Resource{
					Country:       "CH",
					BankID:        "12345",
					BankIDCode:    "CHBCC",
					AccountNumber: "12345678901",
				},
			},
			wantErr: true,
		},
		{
			name: "CH is invalid when account number is more than 12 digits",
			args: args{
				account: client.Resource{
					Country:       "CH",
					BankID:        "12345",
					BankIDCode:    "CHBCC",
					AccountNumber: "1234567890123",
				},
			},
			wantErr: true,
		},
		{
			name: "CH is invalid when account number is 12 characters, not all digits",
			args: args{
				account: client.Resource{
					Country:       "CH",
					BankID:        "12345",
					BankIDCode:    "CHBCC",
					AccountNumber: "12345678901a",
				},
			},
			wantErr: true,
		},
		// US
		{
			name: "US is valid when all fields are valid, account number provided, 6 digits",
			args: args{
				account: client.Resource{
					Country:       "US",
					BankID:        "123456789",
					BIC:           bicExample,
					BankIDCode:    "USABA",
					AccountNumber: "123456",
				},
			},
			wantErr: false,
		},
		{
			name: "US is valid when all fields are valid, account number provided, 17 digits",
			args: args{
				account: client.Resource{
					Country:       "US",
					BankID:        "123456789",
					BIC:           bicExample,
					BankIDCode:    "USABA",
					AccountNumber: "12345678901234567",
				},
			},
			wantErr: false,
		},
		{
			name: "US is valid when all fields are valid, account number not provided",
			args: args{
				account: client.Resource{
					Country:    "US",
					BankID:     "123456789",
					BIC:        bicExample,
					BankIDCode: "USABA",
				},
			},
			wantErr: false,
		},
		{
			name: "US is invalid when iban is provided",
			args: args{
				account: client.Resource{
					Country:    "US",
					BankID:     "123456789",
					BIC:        bicExample,
					BankIDCode: "USABA",
					IBAN:       ibanExample,
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when bank id is missing",
			args: args{
				account: client.Resource{
					Country:    "US",
					BIC:        bicExample,
					BankIDCode: "USABA",
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when bank id is fewer than 9 digits long",
			args: args{
				account: client.Resource{
					Country:    "US",
					BankID:     "12345678",
					BIC:        bicExample,
					BankIDCode: "USABA",
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when bank id more than 9 digits long",
			args: args{
				account: client.Resource{
					Country:    "US",
					BankID:     "1234567890",
					BIC:        bicExample,
					BankIDCode: "USABA",
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when bank id is 9 characters long, not all digits",
			args: args{
				account: client.Resource{
					Country:    "US",
					BankID:     "12345678a",
					BIC:        bicExample,
					BankIDCode: "USABA",
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when bic is missing",
			args: args{
				account: client.Resource{
					Country:    "US",
					BankID:     "123456789",
					BankIDCode: "USABA",
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when bank id code is missing",
			args: args{
				account: client.Resource{
					Country: "US",
					BankID:  "123456789",
					BIC:     bicExample,
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when bank id code is not correct value",
			args: args{
				account: client.Resource{
					Country:    "US",
					BankID:     "123456789",
					BIC:        bicExample,
					BankIDCode: "NO",
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when account number is fewer than 6 digits long",
			args: args{
				account: client.Resource{
					Country:       "US",
					BankID:        "123456789",
					BIC:           bicExample,
					BankIDCode:    "USABA",
					AccountNumber: "12345",
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when account number is more than 17 digits long",
			args: args{
				account: client.Resource{
					Country:       "US",
					BankID:        "123456789",
					BIC:           bicExample,
					BankIDCode:    "USABA",
					AccountNumber: "123456789012345678",
				},
			},
			wantErr: true,
		},
		{
			name: "US is invalid when account number is correct length, but not all are digits",
			args: args{
				account: client.Resource{
					Country:       "US",
					BankID:        "123456789",
					BIC:           bicExample,
					BankIDCode:    "USABA",
					AccountNumber: "123456789a",
				},
			},
			wantErr: true,
		},
		// unknown
		{
			name: "HU is invalid because it's not in the list of countries served",
			args: args{
				account: client.Resource{
					Country: "HU",
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
