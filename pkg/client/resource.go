package client

import "time"

type Resource struct {
	Country                 string    `json:"country"`
	BaseCurrency            string    `json:"base_currency"`
	BankID                  string    `json:"bank_id"`
	BankIDCode              string    `json:"bank_id_code"`
	AccountNumber           string    `json:"account_number"`
	BIC                     string    `json:"bic"`
	IBAN                    string    `json:"iban"`
	CustomerID              string    `json:"customer_id"`
	Name                    [4]string `json:"name"`
	AlternativeNames        [3]string `json:"alternative_names"`
	AccountClassification   string    `json:"account_classification"`
	JointAccount            bool      `json:"joint_account"`
	AccountMatchingOptOut   bool      `json:"account_matching_opt_out"`
	SecondaryIdentification string    `json:"secondary_identification"`
	Switched                bool      `json:"switched"`
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
