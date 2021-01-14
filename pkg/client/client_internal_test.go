package client

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_addHeaders(t *testing.T) {
	const (
		testURL                 = "https://atesturl"
		testJSONBody            = `{data:{key:"value"}}`
		testContentType         = "application/vnd.api+json"
		testHeaderDateThreshold = 15
	)

	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatalf("could not load GMT location: %s", err)
	}

	requestNoBody, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, testURL, nil)
	if err != nil {
		t.Fatalf("could not create a test request with no body: %s", err)
	}

	requestBody, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		testURL,
		strings.NewReader(testJSONBody),
	)
	if err != nil {
		t.Fatalf("could not create a test request with body: %s", err)
	}

	type fields struct {
		BaseURL      string
		DateLocation *time.Location
	}

	type args struct {
		r *http.Request
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantHeaders map[string]string
	}{
		{
			name: "decorates request with headers with no body present",
			fields: fields{
				BaseURL:      testURL,
				DateLocation: gmtLoc,
			},
			args: args{
				r: requestNoBody,
			},
			wantHeaders: map[string]string{
				"Accept": testContentType,
				"Host":   testURL,
			},
		},
		{
			name: "decorates request with headers with body present",
			fields: fields{
				BaseURL:      testURL,
				DateLocation: gmtLoc,
			},
			args: args{
				r: requestBody,
			},
			wantHeaders: map[string]string{
				"Accept":         testContentType,
				"Host":           testURL,
				"Content-Type":   testContentType,
				"Content-Length": strconv.Itoa(len(testJSONBody)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				BaseURL:      tt.fields.BaseURL,
				DateLocation: tt.fields.DateLocation,
			}
			got := c.addHeaders(tt.args.r)
			for k, v := range tt.wantHeaders {
				assert.Equal(t, v, got.Header.Get(k))
			}

			// Check the Date header separately
			headerDate, err := time.Parse(time.RFC1123, got.Header.Get("Date"))
			if err != nil {
				t.Fatalf("could not parse the header date into a time.Time struct: %s", err)
			}
			assert.WithinDuration(t, headerDate, time.Now(), testHeaderDateThreshold*time.Second)
		})
	}
}
