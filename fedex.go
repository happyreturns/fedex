// History: Nov 20 13 tcolar Creation

// Package fedex provides access to (some) FedEx Soap API's and unmarshal answers into Go structures
package fedex

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	Key      string
	Password string
	Account  string
	Meter    string
	HubID    string // for SmartPost

	FedexURL string
}

// TrackByNumber returns tracking info for a specific Fedex tracking number
func (f Fedex) TrackByNumber(carrierCode string, trackingNo string) (*models.TrackReply, error) {

	request := f.trackByNumberRequest(carrierCode, trackingNo)
	response := &models.TrackResponseEnvelope{}

	err := f.makeRequestAndUnmarshalResponse("/trck", request, response)
	if err != nil {
		return nil, fmt.Errorf("make track request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

// Rate : Gets the estimated rates for a shipment
func (f Fedex) Rate(rate *RateRequest) (*models.RateReply, error) {

	request := f.rateRequest(rate)
	response := &models.RateResponseEnvelope{}

	err := f.makeRequestAndUnmarshalResponse("/rate/v24", request, response)
	if err != nil {
		return nil, fmt.Errorf("make rate request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

// CreatePickup creates a pickup
func (f Fedex) CreatePickup(pickupLocation models.PickupLocation, toAddress models.Address) (*models.CreatePickupReply, error) {

	request := f.createPickupRequest(pickupLocation, toAddress)
	response := &models.CreatePickupResponseEnvelope{}

	err := f.makeRequestAndUnmarshalResponse("/pickup/v17", request, response)
	if err != nil {
		return nil, fmt.Errorf("make create pickup request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

// SendNotifications gets notifications sent to an email
func (f Fedex) SendNotifications(trackingNo, email string) (*models.SendNotificationsReply, error) {

	request := f.notificationsRequest(trackingNo, email)
	response := &models.SendNotificationsResponseEnvelope{}

	err := f.makeRequestAndUnmarshalResponse("/track/v16", request, response)
	if err != nil {
		return nil, fmt.Errorf("make send notifications request: %s", err)
	}
	return &response.Reply, nil
}

func (f Fedex) Ship(shipment *Shipment) (*models.ProcessShipmentReply, error) {
	commodities, err := f.commoditiesWithCustoms(shipment)
	if err != nil {
		return nil, fmt.Errorf("commodities with customs: %s", err)
	}
	shipment.Commodities = commodities

	request, err := f.createProcessShipmentRequest(shipment)
	if err != nil {
		return nil, fmt.Errorf("create process shipment request: %s", err)
	}

	response := &models.ShipResponseEnvelope{}
	if err := f.makeRequestAndUnmarshalResponse("/ship/v23", request, response); err != nil {
		return nil, fmt.Errorf("make process shipment request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (f Fedex) commoditiesWithCustoms(shipment *Shipment) (models.Commodities, error) {
	// TODO new struct
	if !shipment.IsInternational() {
		return shipment.Commodities, nil // TODO is this weird
	}

	customsValue, err := shipment.Commodities.CustomsValue()
	if err != nil {
		return nil, fmt.Errorf("customs value: %s", err)
	}
	if customsValue.Amount < 800 {
		return shipment.Commodities, nil
	}

	rateReply, err := f.Rate(&RateRequest{
		FromAndTo:   shipment.FromAndTo,
		Commodities: shipment.Commodities,
	})
	if err != nil {
		return nil, fmt.Errorf("get rate: %s", err)
	}

	charges, err := rateReply.DutiesAndTaxesByItem()
	if err != nil {
		return nil, fmt.Errorf("duties and taxes by item: %s", err)
	}
	if len(charges) != len(shipment.Commodities) {
		return nil, errors.New("charges should match commodities length")
	}

	// TODO not 100% sure what to do with this, or if this is right
	newCommodities := make([]models.Commodity, len(shipment.Commodities))
	for idx, commodity := range shipment.Commodities {
		newCommodities[idx] = commodity
		newCommodities[idx].CustomsValue = &models.Money{
			Currency: charges[idx].Currency,
			Amount:   charges[idx].Amount,
		}
	}

	return newCommodities, nil
}

// TODO get me to work
func (f Fedex) UploadImages(images []models.Image) error {
	request := f.uploadImagesRequest(images)

	response := &models.UploadImagesResponseEnvelope{}
	if err := f.makeRequestAndUnmarshalResponse("/uploaddocument/v11", request, response); err != nil {
		return fmt.Errorf("make upload images request and unmarshal: %s", err)
	}

	return nil

}

func (f Fedex) makeRequestAndUnmarshalResponse(url string, request models.Envelope,
	response models.Response) error {
	// Create request body
	reqXML, err := xml.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request xml: %s", err)
	}

	// Post XML
	content, err := f.postXML(f.FedexURL+url, string(reqXML))
	if err != nil {
		return fmt.Errorf("post xml: %s", err)
	}

	// Parse response
	err = xml.Unmarshal(content, response)
	if err != nil {
		return fmt.Errorf("parse xml: %s", err)
	}

	// Check if reply failed (FedEx responds with 200 even though it failed)
	err = response.Error()
	if err != nil {
		return fmt.Errorf("response error: %s", err)
	}

	return nil
}

// postXML to Fedex API and return response
func (f Fedex) postXML(url string, xml string) ([]byte, error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all bytes: %s", err)
	}
	return content, nil
}
