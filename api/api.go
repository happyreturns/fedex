package api

type API struct {
	Key      string `json:"key"`
	Password string `json:"password"`
	Account  string `json:"account"`
	Meter    string `json:"meter"`
	HubID    string `json:"hubID"` // for SmartPost

	HrEnv string `json:"hrEnv"` // for logging

	FedExURL string `json:"fedexURL"`
}
