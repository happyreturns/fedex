// History: Nov 20 13 tcolar Creation

// Package fedex provides access to () FedEx Soap API's and unmarshal answers into Go structures
package fedex

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/happyreturns/fedex/api"
	"github.com/happyreturns/fedex/models"
	log "github.com/sirupsen/logrus"
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
	api.API
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// Retry logic is expected to be implemented by the caller
func (f Fedex) CreatePickup(pickup *models.Pickup, startTime time.Time, endTime time.Time) (*models.PickupSuccess, error) {
	window := &models.PickupTimeWindow{ReadyTime: startTime, CloseTime: endTime}
	fields := log.Fields{"pickup": pickup, "window": window}

	reply, err := f.API.CreatePickup(pickup, window)
	switch err.(type) {
	case nil:
		fields["reply"] = reply
		log.WithFields(fields).Info("made pickup")
		return &models.PickupSuccess{
			ConfirmationNumber: reply.PickupConfirmationNumber,
			Window:             *window,
		}, nil

	case models.PickupAlreadyExistsError:
		fields["reply"] = reply
		log.WithFields(fields).Info("pickup already exists")
		return nil, fmt.Errorf("pickup already exists %w", err)

	default:
		fields["err"] = err
		log.WithFields(fields).Info("failed pickup %w", err)
	}

	return nil, fmt.Errorf("unable to schedule a fedex pickup")
}

func (f Fedex) Ship(shipment *models.Shipment) (*models.ProcessShipmentReply, error) {
	if f.isSmartPost() && shipment.IsInternational() {
		return nil, errors.New("do not ship internationally with smartpost")
	}

	// Don't use non-smartpost accounts for returns
	if !f.isSmartPost() {
		shipment.Service = "default"
	}

	reply, err := f.API.ProcessShipment(shipment)
	if err != nil {
		return nil, fmt.Errorf("api process shipment: %s", err)
	}

	return reply, nil
}

func (f Fedex) isSmartPost() bool {
	return f.API.HubID != ""
}
