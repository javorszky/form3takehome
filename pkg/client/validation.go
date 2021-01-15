package client

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	reThreeDigits     = regexp.MustCompile(`^\d{3}$`)
	reFiveDigits      = regexp.MustCompile(`^\d{5}$`)
	reSixDigits       = regexp.MustCompile(`^\d{6}$`)
	reSevenDigits     = regexp.MustCompile(`^\d{7}$`)
	reEightDigits     = regexp.MustCompile(`^\d{8}$`)
	reNineDigits      = regexp.MustCompile(`^\d{9}$`)
	reTenDigits       = regexp.MustCompile(`^\d{10}$`)
	reElevenDigits    = regexp.MustCompile(`^\d{11}$`)
	reTwelveDigits    = regexp.MustCompile(`^\d{12}$`)
	reThirteenDigits  = regexp.MustCompile(`^\d{13}$`)
	reSixteenDigits   = regexp.MustCompile(`^\d{16}$`)
	reAUAccountNumber = regexp.MustCompile(`^[1-9]\d{5,9}$`)
	reCARoutingNumber = regexp.MustCompile(`^0\d{8}$`)
	reCAAccountNumber = regexp.MustCompile(`^\d{7,12}$`)
	reHKAccountNumber = regexp.MustCompile(`^\d{9,12}$`)
	reUSAccountNumber = regexp.MustCompile(`^\d{6,17}$`)
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
		return validateGR(account)
	case "HK":
		return validateHK(account)
	case "IT":
		return validateIT(account)
	case "LU":
		return validateLU(account)
	case "NL":
		return validateNL(account)
	case "PL":
		return validatePL(account)
	case "PT":
		return validatePT(account)
	case "ES":
		return validateES(account)
	case "CH":
		return validateCH(account)
	case "US":
		return validateUS(account)
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

func validateGR(account Resource) error {
	errs := make([]string, 0)
	// required, 7 characters, HEBIC (Hellenic Bank Identification Code)
	if !reSevenDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("GR bank id is not correct format, (7 digits): '%s'", account.BankID))
	}

	// Bank ID code is required, has to be GRBIC.
	if account.BankIDCode != "GRBIC" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not GRBIC, got '%s'", account.BankIDCode))
	}

	// Account number optional, 16 characters, generated if not provided.
	if account.AccountNumber != "" && !reSixteenDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 16 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateHK(account Resource) error {
	errs := make([]string, 0)
	// optional, 3 characters, Bank Code or Institution ID
	if !reThreeDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("UK bank id is not correct format: '%s'", account.BankID))
	}

	// BIC required
	if account.BIC == "" {
		errs = append(errs, "BIC is required, got empty")
	}

	// Bank ID code is required, has to be HKNCC
	if account.BankIDCode != "HKNCC" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not HKNCC, got '%s'", account.BankIDCode))
	}

	// Account number optional, 9-12 characters, generated if not provided
	if account.AccountNumber != "" && !reHKAccountNumber.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 9-12 numbers: '%s'", account.AccountNumber))
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

func validateIT(account Resource) error {
	errs := make([]string, 0)
	accountPresent := false
	reAccount := reTenDigits

	// Account number optional, 12 characters, generated if not provided
	if account.AccountNumber != "" && !reTwelveDigits.MatchString(account.AccountNumber) {
		accountPresent = true
		errs = append(errs, fmt.Sprintf("account number was provided, but not 8 numbers: '%s'", account.AccountNumber))
	}

	if accountPresent {
		reAccount = reElevenDigits
	}

	// required, national bank code (ABI) + branch code (CAB), 10 characters if account number is not present,
	// 11 characters with added check digit as first character if account number is present
	if !reAccount.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("IT bank id is not correct format: '%s'", account.BankID))
	}

	// Bank ID code is required, has to be ITNCC
	if account.BankIDCode != "ITNCC" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not ITNCC, got '%s'", account.BankIDCode))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateLU(account Resource) error {
	errs := make([]string, 0)
	// required, 3 characters, IBAN Bank Identifier
	if !reThreeDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("LU bank id is not correct format: '%s'", account.BankID))
	}

	// Bank ID code is required, has to be LULUX
	if account.BankIDCode != "LULUX" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not LULUX, got '%s'", account.BankIDCode))
	}

	// Account number optional, 13 characters, generated if not provided.
	if account.AccountNumber != "" && !reThirteenDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 13 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateNL(account Resource) error {
	errs := make([]string, 0)
	// not supported, has to be empty
	if account.BankID != "" {
		errs = append(errs, fmt.Sprintf("NL Bank id is not supported, has to be empty. Got '%s'", account.BankID))
	}

	// BIC required
	if account.BIC == "" {
		errs = append(errs, "BIC is required, got empty")
	}

	// Bank ID code not supported, has to be empty
	if account.BankIDCode != "" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not supported, has to be empty. Got '%s'", account.BankIDCode))
	}

	// Account number optional, 8 characters, generated if not provided
	if account.AccountNumber != "" && !reTenDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 10 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validatePL(account Resource) error {
	errs := make([]string, 0)
	// required, 8 characters, national bank code + branch code + national check digit
	if !reEightDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("PL bank id is not correct format: '%s'", account.BankID))
	}

	// Bank ID code is required, has to be PLKNR
	if account.BankIDCode != "PLKNR" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not PLKNR, got '%s'", account.BankIDCode))
	}

	// Account number optional, 16 characters, generated if not provided
	if account.AccountNumber != "" && !reSixteenDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 16 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validatePT(account Resource) error {
	errs := make([]string, 0)
	// required, 8 characters, bank identifier + PSP reference number
	if !reEightDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("PT bank id is not correct format: '%s'", account.BankID))
	}

	// Bank ID code is required, has to be PTNCC
	if account.BankIDCode != "PTNCC" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not PTNCC, got '%s'", account.BankIDCode))
	}

	// Account number optional, 11 characters, generated if not provided
	if account.AccountNumber != "" && !reElevenDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 11 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateES(account Resource) error {
	errs := make([]string, 0)
	// required, 8 characters, Código de entidad + Código de oficina
	if !reEightDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("ES bank id is not correct format: '%s'", account.BankID))
	}

	// Bank ID code is required, has to be ESNCC
	if account.BankIDCode != "ESNCC" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not ESNCC, got '%s'", account.BankIDCode))
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

func validateCH(account Resource) error {
	errs := make([]string, 0)
	// required, 5 characters
	if !reFiveDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("CH bank id is not correct format: '%s'", account.BankID))
	}

	// Bank ID code is required, has to be CHBCC
	if account.BankIDCode != "CHBCC" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not CHBCC, got '%s'", account.BankIDCode))
	}

	// Account number optional, 12 characters, generated if not provided
	if account.AccountNumber != "" && !reTwelveDigits.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 12 numbers: '%s'", account.AccountNumber))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, ", "))
}

func validateUS(account Resource) error {
	errs := make([]string, 0)
	// required, 9 characters, ABA routing number
	if !reNineDigits.MatchString(account.BankID) {
		errs = append(errs, fmt.Sprintf("US bank id is not correct format: '%s'", account.BankID))
	}

	// BIC required
	if account.BIC == "" {
		errs = append(errs, "BIC is required, got empty")
	}

	// Bank ID code is required, has to be USABA
	if account.BankIDCode != "USABA" {
		errs = append(errs, fmt.Sprintf("Bank ID Code is not USABA, got '%s'", account.BankIDCode))
	}

	// Account number optional, 6-17 characters, generated if not provided
	if account.AccountNumber != "" && !reUSAccountNumber.MatchString(account.AccountNumber) {
		errs = append(errs, fmt.Sprintf("account number was provided, but not 6-17 numbers: '%s'", account.AccountNumber))
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
