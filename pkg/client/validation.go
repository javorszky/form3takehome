package client

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	reThreeDigits     = regexp.MustCompile(`^\d{3}$`)
	reSixDigits       = regexp.MustCompile(`^\d{6}$`)
	reSevenDigits     = regexp.MustCompile(`^\d{7}$`)
	reEightDigits     = regexp.MustCompile(`^\d{8}$`)
	reTenDigits       = regexp.MustCompile(`^\d{10}$`)
	reAUAccountNumber = regexp.MustCompile(`^[1-9]\d{5,9}$`)
	reCARoutingNumber = regexp.MustCompile(`^0\d{8}$`)
	reCAAccountNumber = regexp.MustCompile(`^\d{7,12}$`)
)

//nolint:gocyclo
func ValidateResource(account Resource) error {
	switch account.Country {
	case "GB":
		return validateGB(account)
	case "AU":
		return validateAU(account)
	case "BE":
		return validateBE(account)
	case "CA":
		return validateCA(account)
	case "FR":
		return validateFR(account)
	case "DE":
		return validateDE(account)
	case "GR":
	case "HK":
	case "IT":
	case "LU":
	case "NL":
	case "PL":
	case "PT":
	case "ES":
	case "CH":
	case "US":
	}

	return nil
}

func validateGB(account Resource) error {
	errs := make([]string, 0)
	// required, 6 characters, UK sort code
	if !reSixDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("UK bank id is not correct format: '%s'", account.BankID))
	}

	// BIC required
	if account.BIC == "" {
		errs = append(errs, "BIC is required, got empty")
	}

	// Bank ID code is required, has to be GBDSC
	if account.BankIDCode != "GBDSC" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not GBDSC, got '%s'", account.BankIDCode))
	}

	// Account number optional, 8 characters, generated if not provided
	if account.AccountNumber != "" && !reEightDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 8 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateAU(account Resource) error {
	errs := make([]string, 0)
	// optional, 6 characters, Australian Bank State Branch (BSB) code
	if account.BankID != "" && !reSixDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("AU bank id was provided, but is not correct format: '%s'", account.BankID))
	}

	// BIC required
	if account.BIC == "" {
		errs = append(errs, "BIC is required, got empty")
	}

	// Bank ID code is required, has to be AUBSB
	if account.BankIDCode != "AUBSB" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not AUBSB, got '%s'", account.BankIDCode))
	}

	// Account number optional, 6-10 characters, first character cannot be 0, generated if not provided.
	if account.AccountNumber != "" && !reAUAccountNumber.MatchString(account.AccountNumber) {
		errs = append(
			errs,
			fmt.Sprintf(
				"account number was provided, but not correct format: between 6-10 digits, first is not 0: '%s'",
				account.AccountNumber,
			),
		)
	}

	// IBAN has to be empty
	if account.IBAN != "" {
		errs = append(errs, fmt.Sprintf("IBAN has to be empty, got '%s'", account.IBAN))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateBE(account Resource) error {
	errs := make([]string, 0)
	// required, 3 characters
	if !reThreeDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("BE bank id is not correct format *3 digits): '%s'", account.BankID))
	}

	// Bank ID code is required, has to be BE
	if account.BankIDCode != "BE" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not BE, got '%s'", account.BankIDCode))
	}

	// Account number optional, 7 characters, generated if not provided
	if account.AccountNumber != "" && !reSevenDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 7 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateCA(account Resource) error {
	errs := make([]string, 0)
	// optional, 9 characters starting with zero, Routing Number for Electronic Funds Transfers
	if account.BankID != "" && !reCARoutingNumber.MatchString(account.BankID) {
		errs = append(
			errs,
			fmt.Sprintf(
				"CA electronic funds routing number is not correct format (9 digits, leading 0): '%s'",
				account.BankID,
			),
		)
	}

	// BIC required
	if account.BIC == "" {
		errs = append(errs, "BIC is required, got empty")
	}

	// Bank ID code is optional, if provided has to be CACPA
	if account.BankIDCode != "" && account.BankIDCode != "CACPA" {
		errs = append(
			errs,
			fmt.Sprintf(
				"Bank ID Code was provided, is not CACPA, got '%s'",
				account.BankIDCode,
			),
		)
	}

	// Account number optional, 8 characters, generated if not provided
	if account.AccountNumber != "" && !reCAAccountNumber.MatchString(account.AccountNumber) {
		errs = append(
			errs,
			fmt.Sprintf(
				"account number was provided, but not 8 numbers: '%s'",
				account.AccountNumber,
			),
		)
	}

	// IBAN: not supported, has to be empty
	if account.IBAN != "" {
		errs = append(errs, fmt.Sprintf("IBAN is not supported, has to be empty. Got '%s'", account.IBAN))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateFR(account Resource) error {
	errs := make([]string, 0)
	// required, 10 characters, national bank code + branch code (code guichet)
	if !reTenDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("FR bank id is not correct format, needs to be 10 digits, got: '%s'", account.BankID))
	}

	// Bank ID code is required, has to be FR
	if account.BankIDCode != "FR" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not FR, got '%s'", account.BankIDCode))
	}

	// Account number optional, 10 characters, generated if not provided
	if account.AccountNumber != "" && !reTenDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 10 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateDE(account Resource) error {
	errs := make([]string, 0)
	// required, 8 characters, Bankleitzahl (BLZ)
	if !reEightDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("DE bank id is not correct format (8 digits): '%s'", account.BankID))
	}

	// Bank ID code is required, has to be DEBLZ
	if account.BankIDCode != "DEBLZ" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not DEBLZ, got '%s'", account.BankIDCode))
	}

	// Account number optional, 7 characters, generated if not provided
	if account.AccountNumber != "" && !reSevenDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 7 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}
