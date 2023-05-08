// History: Nov 20 13 tcolar Creation

// Package fedex provides access to () FedEx Soap API's and unmarshal answers into Go structures
package fedex

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/happyreturns/fedex/api"
	"github.com/happyreturns/fedex/clock"
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

var laTimeZone *time.Location

func init() {
	var err error
	laTimeZone, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// CreatePickup creates a pickup with retry logic to try pickups on different days/times
func (f Fedex) CreatePickup(pickup *models.Pickup) (*models.PickupSuccess, error) {
	for delay := 0; delay <= 5; delay++ {
		// Attempt pickup at normal time (1:00 PM - 5:00 PM)
		earliestPickupTime := &models.PickupOffset{Days: delay, Hours: 13, Minutes: 0}
		windowHoursDuration := 4
		pickupResponse := f.bookPickup(pickup, earliestPickupTime, windowHoursDuration)
		if pickupResponse != nil {
			return pickupResponse, nil
		}
		log.Warn(fmt.Sprintf("Unable to schedule fedex pickup for 1:00 PM - 5:00 PM delay %d", delay))

		// Attempt pickup at later time (1:00 PM - 6:00 PM)
		earliestPickupTime = &models.PickupOffset{Days: delay, Hours: 13, Minutes: 0}
		windowHoursDuration = 5
		pickupResponse = f.bookPickup(pickup, earliestPickupTime, windowHoursDuration)
		if pickupResponse != nil {
			return pickupResponse, nil
		}
		log.Warn(fmt.Sprintf("Unable to schedule fedex pickup for 1:00 PM - 6:00 PM delay %d", delay))
	}

	return nil, fmt.Errorf("unable to schedule a fedex pickup")
}

func (f Fedex) bookPickup(pickup *models.Pickup, pickupOffset *models.PickupOffset, windowHoursDuration int) *models.PickupSuccess {
	fields := log.Fields{"pickup": pickup}

	window, err := pickupTimeWindow(clock.NewClock(), pickup.PickupLocation.Address, pickupOffset, windowHoursDuration)
	if err != nil {
		log.WithFields(fields).Error("calculate pickup time", err)
		return nil
	}
	fields["window"] = window

	reply, err := f.API.CreatePickup(pickup, window)
	switch err.(type) {
	case nil:
		fields["reply"] = reply
		log.WithFields(fields).Info("made pickup")
		return &models.PickupSuccess{
			ConfirmationNumber: reply.PickupConfirmationNumber,
			Window:             *window,
		}

	case models.PickupAlreadyExistsError:
		fields["reply"] = reply
		log.WithFields(fields).Info("pickup already exists")
		return nil

	default:
		fields["err"] = err
		log.WithFields(fields).Info("failed pickup")
	}

	return nil
}

func pickupTimeWindow(clock clock.Clock, pickupAddress models.Address, pickupOffset *models.PickupOffset, windowHoursDuration int) (*models.PickupTimeWindow, error) {
	location, err := toLocation(pickupAddress)
	if err != nil {
		location = laTimeZone
	}

	readyTime := clock.Now().In(location).Add(time.Duration(pickupOffset.Days*24) * time.Hour)

	// If it's past the ready time of the current day, ship the next day, not today
	if readyTime.After(timeForReadyPickup(readyTime, pickupOffset)) {
		readyTime = readyTime.Add(24 * time.Hour)
	}
	readyTime = timeForReadyPickup(readyTime, pickupOffset)

	// Don't schedule pickups for Sunday
	if readyTime.Weekday() == time.Sunday {
		return nil, fmt.Errorf("no pickups on sunday %d", pickupOffset.Days)
	}

	return &models.PickupTimeWindow{
		ReadyTime: readyTime,
		CloseTime: readyTime.Add(time.Duration(windowHoursDuration) * time.Hour),
	}, nil
}

func timeForReadyPickup(t time.Time, pickupOffset *models.PickupOffset) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), pickupOffset.Hours, pickupOffset.Minutes, 0, 0, t.Location())
}

// toLocation attempts to return the timezone based on state, returning los
// angeles if unable to
func toLocation(pickupAddress models.Address) (*time.Location, error) {
	tzDatabaseName := ""
	switch strings.ToUpper(pickupAddress.StateOrProvinceCode) {
	case "AK":
		tzDatabaseName = "America/Anchorage"
	case "HI":
		tzDatabaseName = "Pacific/Honolulu"
	case "AL", "AR", "IL", "IA", "KS", "KY", "LA", "MN", "MS", "MO", "NE", "ND", "OK", "SD", "TN", "TX", "WI":
		tzDatabaseName = "America/Chicago"
	case "AZ", "CO", "ID", "MT", "NM", "UT", "WY":
		tzDatabaseName = "America/Denver"
	case "CT", "DC", "DE", "FL", "GA", "IN", "ME", "MD", "MA", "MI", "NH", "NJ", "NY", "NC", "OH", "PA", "RI", "SC", "VT", "VA", "WV":
		tzDatabaseName = "America/New_York"
	default:
		return laTimeZone, nil
	}

	timeZone, err := time.LoadLocation(tzDatabaseName)
	if err != nil {
		return nil, fmt.Errorf("load location from time zone %s: %s", tzDatabaseName, err)
	}
	return timeZone, nil
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
