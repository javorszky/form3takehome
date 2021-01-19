package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/javorszky/form3takehome/pkg/client"
	"github.com/javorszky/form3takehome/pkg/config"
)

const timeOutExampleMs = 500

func main() {
	// this is an example implementation
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("failed to get config: %s", err)
	}

	httpClient := http.Client{
		Timeout: timeOutExampleMs * time.Millisecond,
	}

	gmtLoc, err := time.LoadLocation("GMT")
	if err != nil {
		log.Fatalf("failed to load gmt timezone: %s", err)
	}

	c := client.New(cfg, httpClient, gmtLoc)

	p, err := c.Create(client.Resource{
		Country:    "GB",
		BankIDCode: "GBDSC",
		BIC:        "BARCGB22XXX",
		BankID:     "123456",
	})

	fmt.Printf("inserted a resource and this is what the service said with no error: %s: %#v", err, p)
}
