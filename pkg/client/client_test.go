package client_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/javorszky/form3takehome/pkg/client"
	"github.com/javorszky/form3takehome/pkg/config"
)

func TestNew(t *testing.T) {
	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load GMT location")
	}

	type args struct {
		cfg config.Config
		gmt *time.Location
	}

	tests := []struct {
		name string
		args args
		want client.Client
	}{
		{
			name: "constructs a new client based on incoming data",
			args: args{
				cfg: config.Config{
					AccountsAPIURL: "https://testurl",
				},
				gmt: gmtLoc,
			},
			want: client.Client{
				BaseURL:      "https://testurl",
				DateLocation: gmtLoc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, client.New(tt.args.cfg, tt.args.gmt))
		})
	}
}
