package config

import (
	"fmt"
	"os"
)

const (
	AccountsAPIURLKey = "ACCOUNTS_ADDRESS"
	OrganisationIDKey = "ORGANISATION_ID"
)

type Config struct {
	AccountsAPIURL string
	OrganisationID string
}

type validationFunc func(string) error

func Get() (Config, error) {
	for key, f := range map[string]validationFunc{
		AccountsAPIURLKey: stringNotEmpty,
		OrganisationIDKey: stringNotEmpty,
	} {
		err := f(key)
		if err != nil {
			return Config{}, fmt.Errorf("config.Get: %s failed validation: %w", key, err)
		}
	}

	return Config{
		AccountsAPIURL: os.Getenv(AccountsAPIURLKey),
		OrganisationID: os.Getenv(OrganisationIDKey),
	}, nil
}

func stringNotEmpty(key string) error {
	setting := os.Getenv(key)
	if setting == "" {
		return fmt.Errorf("setting with key '%s' is empty", setting)
	}

	return nil
}
