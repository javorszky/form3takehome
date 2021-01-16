package client

import "time"

// Resource in this case encodes the Organisation.Account resource as the API only deals with this.
type Resource struct {
	Country                 string    `json:"country"`
	BaseCurrency            string    `json:"base_currency,omitempty"`
	BankID                  string    `json:"bank_id,omitempty"`
	BankIDCode              string    `json:"bank_id_code,omitempty"`
	AccountNumber           string    `json:"account_number,omitempty"`
	BIC                     string    `json:"bic,omitempty"`
	IBAN                    string    `json:"iban,omitempty"`
	CustomerID              string    `json:"customer_id,omitempty"`
	Name                    [4]string `json:"name"`
	AlternativeNames        [3]string `json:"alternative_names,omitempty"`
	AccountClassification   string    `json:"account_classification,omitempty"`
	JointAccount            bool      `json:"joint_account,omitempty"`
	AccountMatchingOptOut   bool      `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string    `json:"secondary_identification,omitempty"`
	Switched                bool      `json:"switched,omitempty"`
	Status                  string    `json:"status"`
}

// Data struct encodes the the data part of a create request, and data part of the responses.
type Data struct {
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	Type           string    `json:"type"`
	Version        int       `json:"version"`
	CreatedOn      time.Time `json:"created_on"`
	ModifiedOn     time.Time `json:"modified_on"`
	Attributes     Resource  `json:"attributes"`
}

// Links is used to encode the links section from the responses.
type Links struct {
	Self  string `json:"self"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}

// Payload struct is used to encode json requests and responses where there are only one of Resource being sent or
// received, such as the Organisation.Accounts.Fetch and Organisation.Accounts.Create endpoints.
type Payload struct {
	Data  Data  `json:"data"`
	Links Links `json:"links,omitempty"`
}

// MultiPayload struct is used to encode json requests and responses where they contain an array of data objects, such
// as the Organisation.Accounts.List API endpoint.
type MultiPayload struct {
	Data  []Data `json:"data"`
	Links Links  `json:"links"`
}
