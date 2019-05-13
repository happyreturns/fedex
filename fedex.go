// History: Nov 20 13 tcolar Creation

// Package fedex provides access to () FedEx Soap API's and unmarshal answers into Go structures
package fedex

import (
	"fmt"

	"github.com/happyreturns/fedex/api"
	"github.com/happyreturns/fedex/models"
)

// Convenience constants for standard Fedex API URLs
const (
	FedexAPIURL               = "https://ws.fedex.com:443/web-services"
	FedexAPITestURL           = "https://wsbeta.fedex.com:443/web-services"
	CarrierCodeExpress        = "FDXE"
	CarrierCodeGround         = "FDXG"
	CarrierCodeFreight        = "FXFR"
	CarrierCodeSmartPost      = "FXSP"
	CarrierCodeCustomCritical = "FXCC"
)

// Fedex : Utility to retrieve data from Fedex API
// Bypassing painful proper SOAP implementation and just crafting minimal XML messages to get the data we need.
// Fedex WSDL docs here: http://images.fedex.com/us/developer/product/WebServices/MyWebHelp/DeveloperGuide2012.pdf
type Fedex struct {
	API api.API
}

// TrackByNumber returns tracking info for a specific Fedex tracking number
func (f Fedex) TrackByNumber(carrierCode string, trackingNo string) (*models.TrackReply, error) {
	reply, err := f.API.TrackByNumber(carrierCode, trackingNo)
	if err != nil {
		return nil, fmt.Errorf("api track by number: %s", err)
	}
	return reply, nil

}

// Rate : Gets the estimated rates for a shipment
func (f Fedex) Rate(rate *models.Rate) (*models.RateReply, error) {
	reply, err := f.API.Rate(rate)
	if err != nil {
		return nil, fmt.Errorf("api rate: %s", err)
	}
	return reply, nil
}

// CreatePickup creates a pickup
func (f Fedex) CreatePickup(pickup *models.Pickup) (*models.CreatePickupReply, error) {
	var (
		reply *models.CreatePickupReply
		err   error
	)

	for delay := 0; delay <= 5; delay++ {
		reply, err = f.API.CreatePickup(pickup, delay)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("api create pickup: %s", err)
	}
	return reply, nil
}

// SendNotifications gets notifications sent to an email
func (f Fedex) SendNotifications(trackingNo, email string) (*models.SendNotificationsReply, error) {
	reply, err := f.API.SendNotifications(trackingNo, email)
	if err != nil {
		return nil, fmt.Errorf("api send notifications: %s", err)
	}
	return reply, nil
}

func (f Fedex) Ship(shipment *models.Shipment) (*models.ProcessShipmentReply, error) {
	commodities, err := f.commoditiesWithCustoms(shipment)
	if err != nil {
		return nil, fmt.Errorf("commodities with customs: %s", err)
	}
	shipment.Commodities = commodities

	reply, err := f.API.ProcessShipment(shipment)
	if err != nil {
		return nil, fmt.Errorf("api process shipment: %s", err)
	}

	return reply, nil
}

func (f Fedex) UploadImages(images []models.Image) error {
	err := f.API.UploadImages(images)
	if err != nil {
		return fmt.Errorf("upload images: %s", err)
	}
	return nil
}

// TODO unit price or customs value on shipment.Commodities
func (f Fedex) commoditiesWithCustoms(shipment *models.Shipment) (models.Commodities, error) {
	needsInvoice, err := needsCustomCommercialInvoice(shipment)
	if err != nil {
		return nil, fmt.Errorf("needs custom commercial invoice: %s", err)
	}
	if !needsInvoice {
		return shipment.Commodities, nil
	}

	rateReply, err := f.API.RateForCustoms(&models.Rate{
		FromAndTo:   shipment.FromAndTo,
		Commodities: shipment.Commodities,
	})
	if err != nil {
		return nil, fmt.Errorf("rate for customs: %s", err)
	}

	charges, err := rateReply.TaxableValues()
	if err != nil {
		return nil, fmt.Errorf("taxable values: %s", err)
	}
	if len(charges) != len(shipment.Commodities) {
		return nil, fmt.Errorf("charges should match commodities length %d %d", len(charges), len(shipment.Commodities))
	}

	// TODO not 100% sure what to do with this, or if this is right
	newCommodities := make([]models.Commodity, len(shipment.Commodities))
	for idx, commodity := range shipment.Commodities {
		newCommodities[idx] = commodity
		newCommodities[idx].CustomsValue = &models.Money{
			Currency: charges[idx].Currency,
			Amount:   charges[idx].Amount,
		}
		newCommodities[idx].UnitPrice = nil
	}

	return newCommodities, nil
}

func needsCustomCommercialInvoice(shipment *models.Shipment) (bool, error) {
	if !shipment.IsInternational() {
		return false, nil
	}

	customsValue, err := shipment.Commodities.CustomsValue()
	if err != nil {
		return false, fmt.Errorf("customs value: %s", err)
	}
	return customsValue.Currency == "USD" && customsValue.Amount >= 800, nil
}
