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
	reAUAccountNumber = regexp.MustCompile(`^[1-9]\d{5,9}$`)
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
	case "FR":
	case "DE":
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
