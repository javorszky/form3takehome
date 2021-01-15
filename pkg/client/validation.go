//nolint:dupl
package client

import (
	"fmt"
	"regexp"
)

const (
	GBBankID = "GBDSC"
	AUBankID = "AUBSB"
	BEBankID = "BE"
	CABankID = "CACPA"
	FRBankID = "FR"
	DEBankID = "DEBLZ"
	GRBankID = "GRBIC"
	HKBankID = "HKNCC"
	ITBankID = "ITNCC"
	LUBankID = "LULUX"
	PLBankID = "PLKNR"
	PTBankID = "PTNCC"
	ESBankID = "ESNCC"
	CHBankID = "CHBCC"
	USBankID = "USABA"
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

	return fmt.Errorf("unsupported country code: %s", account.Country)
}

func validateGB(account Resource) error {
	r, err := bicRequired(account, nil)
	r, err = bankIDRequiredMust(r, err, reSixDigits)
	r, err = bankIDCodeMust(r, err, GBBankID)
	_, err = accountNumberOptionalMust(r, err, reEightDigits)

	return err
}

func validateAU(account Resource) error {
	r, err := ibanNotSupported(bicRequired(account, nil))
	r, err = bankIDOptionalMust(r, err, reSixDigits)
	r, err = bankIDCodeMust(r, err, AUBankID)
	_, err = accountNumberOptionalMust(r, err, reAUAccountNumber)

	return err
}

func validateBE(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reThreeDigits)
	r, err = bankIDCodeMust(r, err, BEBankID)
	_, err = accountNumberOptionalMust(r, err, reSevenDigits)

	return err
}

func validateCA(account Resource) error {
	r, err := ibanNotSupported(bicRequired(account, nil))
	r, err = bankIDOptionalMust(r, err, reCARoutingNumber)
	r, err = bankIDCodeOptionalMust(r, err, CABankID)
	_, err = accountNumberOptionalMust(r, err, reCAAccountNumber)

	return err
}

func validateFR(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reTenDigits)
	r, err = bankIDCodeMust(r, err, FRBankID)
	_, err = accountNumberOptionalMust(r, err, reTenDigits)

	return err
}

func validateDE(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reEightDigits)
	r, err = bankIDCodeMust(r, err, DEBankID)
	_, err = accountNumberOptionalMust(r, err, reSevenDigits)

	return err
}

func validateGR(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reSevenDigits)
	r, err = bankIDCodeMust(r, err, GRBankID)
	_, err = accountNumberOptionalMust(r, err, reSixteenDigits)

	return err
}

func validateHK(account Resource) error {
	r, err := ibanNotSupported(bicRequired(account, nil))
	r, err = bankIDOptionalMust(r, err, reThreeDigits)
	r, err = bankIDCodeOptionalMust(r, err, HKBankID)
	_, err = accountNumberOptionalMust(r, err, reHKAccountNumber)

	return err
}

func validateIT(account Resource) error {
	var accErr error

	accountPresent := false
	reAccount := reTenDigits

	// Account number optional, 12 characters, generated if not provided
	if account.AccountNumber != "" {
		accountPresent = true

		if !reTwelveDigits.MatchString(account.AccountNumber) {
			accErr = fmt.Errorf("account number was provided, but not 12 numbers: '%s'", account.AccountNumber)
		}
	}

	if accountPresent {
		reAccount = reElevenDigits
	}

	r, err := bankIDRequiredMust(account, accErr, reAccount)
	_, err = bankIDCodeMust(r, err, ITBankID)

	return err
}

func validateLU(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reThreeDigits)
	r, err = bankIDCodeMust(r, err, LUBankID)
	_, err = accountNumberOptionalMust(r, err, reThirteenDigits)

	return err
}

func validateNL(account Resource) error {
	r, err := bicRequired(bankIDNotSupported(account, nil))
	r, err = bankIDCodeMust(r, err, "")
	_, err = accountNumberOptionalMust(r, err, reTenDigits)

	return err
}

func validatePL(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reEightDigits)
	r, err = bankIDCodeMust(r, err, PLBankID)
	_, err = accountNumberOptionalMust(r, err, reSixteenDigits)

	return err
}

func validatePT(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reEightDigits)
	r, err = bankIDCodeMust(r, err, PTBankID)
	_, err = accountNumberOptionalMust(r, err, reElevenDigits)

	return err
}

func validateES(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reEightDigits)
	r, err = bankIDCodeMust(r, err, ESBankID)
	_, err = accountNumberOptionalMust(r, err, reTenDigits)

	return err
}

func validateCH(account Resource) error {
	r, err := bankIDRequiredMust(account, nil, reFiveDigits)
	r, err = bankIDCodeMust(r, err, CHBankID)
	_, err = accountNumberOptionalMust(r, err, reTwelveDigits)

	return err
}

func validateUS(account Resource) error {
	r, err := ibanNotSupported(bicRequired(account, nil))
	r, err = accountNumberOptionalMust(r, err, reUSAccountNumber)
	r, err = bankIDCodeMust(r, err, USBankID)
	_, err = bankIDRequiredMust(r, err, reNineDigits)

	return err
}

func bicRequired(r Resource, e error) (Resource, error) {
	if r.BIC == "" {
		return r, fmt.Errorf("BIC is required, was empty: %w", e)
	}

	return r, e
}

func ibanNotSupported(r Resource, e error) (Resource, error) {
	if r.IBAN != "" {
		return r, fmt.Errorf("IBAN is not supported, got '%s': %w", r.IBAN, e)
	}

	return r, e
}

func bankIDNotSupported(r Resource, e error) (Resource, error) {
	if r.BankID != "" {
		return r, fmt.Errorf("bank ID is not supported, has to be empty. Got '%s': %w", r.BankID, e)
	}

	return r, e
}

func bankIDCodeMust(r Resource, e error, bankIDCode string) (Resource, error) {
	if r.BankIDCode != bankIDCode {
		return r, fmt.Errorf("bank ID Code is not '%s', got %s: %w", bankIDCode, r.BankIDCode, e)
	}

	return r, e
}

func bankIDCodeOptionalMust(r Resource, e error, bankIDCode string) (Resource, error) {
	if r.BankIDCode != "" && r.BankIDCode != bankIDCode {
		return r, fmt.Errorf("bank ID Code is not '%s', got '%s': %w", bankIDCode, r.BankIDCode, e)
	}

	return r, e
}

func bankIDRequiredMust(r Resource, e error, pattern *regexp.Regexp) (Resource, error) {
	if !pattern.MatchString(r.BankID) {
		return r, fmt.Errorf("%s bank id is not in correct format. '%s': %w", r.Country, r.BankID, e)
	}

	return r, e
}

func bankIDOptionalMust(r Resource, e error, pattern *regexp.Regexp) (Resource, error) {
	if r.BankID != "" && !pattern.MatchString(r.BankID) {
		return r, fmt.Errorf("%s bank id is not in correct format. '%s': %w", r.Country, r.BankID, e)
	}

	return r, e
}

func accountNumberOptionalMust(r Resource, e error, pattern *regexp.Regexp) (Resource, error) {
	if r.AccountNumber != "" && !pattern.MatchString(r.AccountNumber) {
		return r, fmt.Errorf("%s account number is not in correct format. '%s': %w", r.Country, r.AccountNumber, e)
	}

	return r, e
}
