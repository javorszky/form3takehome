package config

import (
	"fmt"
	"os"
)

const (
	AccountsAPIURLKey = "ACCOUNTS_ADDRESS"
)

type Config struct {
	AccountsAPIURL string
}

func Get() (Config, error) {
	if notEmpty(os.Getenv(AccountsAPIURLKey)) {
		return Config{
			AccountsAPIURL: os.Getenv(AccountsAPIURLKey),
		}, nil
	}

	return Config{}, fmt.Errorf("config.Get: required setting of %s is empty", AccountsAPIURLKey)
}

func notEmpty(setting string) bool {
	return setting != ""
}
