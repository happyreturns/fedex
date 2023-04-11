package api

import "net/http"

type API struct {
	Key      string `json:"key"`
	Password string `json:"password"`
	Account  string `json:"account"`
	Meter    string `json:"meter"`
	HubID    string `json:"hubID"` // for SmartPost

	HttpClient *http.Client `json:"-" yaml:"-"`

	FedExURL string `json:"fedexURL"`
}
