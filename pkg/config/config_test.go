package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/javorszky/form3takehome/pkg/config"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		want    config.Config
		wantErr bool
	}{
		{
			name: "correctly returns config struct based on existing non empty env var",
			setup: func() {
				_ = os.Setenv(config.AccountsAPIURLKey, "anurl")
				_ = os.Setenv(config.OrganisationIDKey, "an-uuidv4")
			},
			want: config.Config{
				AccountsAPIURL: "anurl",
				OrganisationID: "an-uuidv4",
			},
			wantErr: false,
		},
		{
			name: "correctly returns error and empty config on existing but empty environment variable",
			setup: func() {
				_ = os.Setenv(config.AccountsAPIURLKey, "")
				_ = os.Setenv(config.OrganisationIDKey, "")
			},
			want:    config.Config{},
			wantErr: true,
		},
		{
			name: "correctly returns error and empty config on non-existing environment variables",
			setup: func() {
				// Do not set an environment variable.
			},
			want:    config.Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables so there's no bleedthrough between test cases.
			os.Clearenv()

			// Call the setup function which should set the required environment variables to be used in this specific
			// test.
			tt.setup()

			got, err := config.Get()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
