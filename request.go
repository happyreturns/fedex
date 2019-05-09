package fedex

import "github.com/happyreturns/fedex/models"

type FromAndTo struct {
	FromAddress models.Address
	ToAddress   models.Address
	FromContact models.Contact
	ToContact   models.Contact
}

func (ft FromAndTo) IsInternational() bool {
	fromCountryCode := ft.FromAddress.CountryCode
	if fromCountryCode == "" {
		fromCountryCode = "US"
	}

	toCountryCode := ft.ToAddress.CountryCode
	if toCountryCode == "" {
		toCountryCode = "US"
	}

	return fromCountryCode != toCountryCode
}

// RateRequest wraps all the Fedex API fields needed for getting a rate
type RateRequest struct {
	FromAndTo

	// Only used for international ground shipments
	Commodities models.Commodities
}

// Shipment wraps all the Fedex API fields needed for creating a shipment
type Shipment struct {
	FromAndTo

	NotificationEmail string
	Reference         string
	Service           string

	// Only used for international ground shipments
	Commodities models.Commodities
}
