package client

import "time"

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

type Data struct {
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	Type           string    `json:"type"`
	Version        int       `json:"version"`
	CreatedOn      time.Time `json:"created_on"`
	ModifiedOn     time.Time `json:"modified_on"`
	Attributes     Resource  `json:"attributes"`
}

type Links struct {
	Self  string `json:"self"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}

type Payload struct {
	Data  Data  `json:"data"`
	Links Links `json:"links"`
}
