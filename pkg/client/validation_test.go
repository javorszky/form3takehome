package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateResource(t *testing.T) {
	const (
		bicExample  = "BARCGB22XXX"
		ibanExample = "GB33BUKB20201555555555"
	)
	type args struct {
		account Resource
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GB is valid when all fields are valid, account number and iban provided",
			args: args{
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
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
				account: Resource{
					Country:    "AU",
					BankIDCode: "AUBSB",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Error(t, ValidateResource(tt.args.account))
			} else {
				assert.NoError(t, ValidateResource(tt.args.account))
			}
		})
	}
}
