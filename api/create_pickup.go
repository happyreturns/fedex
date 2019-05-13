package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/happyreturns/fedex/models"
)

var laTimeZone *time.Location

func init() {
	var err error
	laTimeZone, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
}

func (a API) CreatePickup(pickup *models.Pickup, numDaysToDelay int) (*models.CreatePickupReply, error) {
	request := a.createPickupRequest(pickup, numDaysToDelay)
	response := &models.CreatePickupResponseEnvelope{}

	err := a.makeRequestAndUnmarshalResponse("/pickup/v17", request, response)
	if err != nil {
		return nil, fmt.Errorf("make create pickup request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) createPickupRequest(pickup *models.Pickup, numDaysToDelay int) models.Envelope {
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/pickup/v17",
		Body: struct {
			CreatePickupRequest models.CreatePickupRequest `xml:"q0:CreatePickupRequest"`
		}{
			CreatePickupRequest: models.CreatePickupRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      a.Key,
							Password: a.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: a.Account,
						MeterNumber:   a.Meter,
					},
					Version: models.Version{
						ServiceID: "disp",
						Major:     17,
					},
				},
				OriginDetail: models.OriginDetail{
					UseAccountAddress:       false,
					PickupLocation:          pickup.PickupLocation,
					PackageLocation:         "NONE",
					BuildingPart:            "SUITE",
					BuildingPartDescription: "",
					ReadyTimestamp:          models.Timestamp(pickupTime(pickup.PickupLocation.Address, numDaysToDelay)),
					CompanyCloseTime:        "16:00:00", // TODO not necessarily true
				},
				FreightPickupDetail: models.FreightPickupDetail{
					ApprovedBy:  pickup.PickupLocation.Contact,
					Payment:     "SENDER",
					Role:        "SHIPPER",
					SubmittedBy: models.Contact{},
					LineItems: []models.FreightPickupLineItem{
						{
							Service:        "INTERNATIONAL_ECONOMY_FREIGHT",
							SequenceNumber: 1,
							Destination:    pickup.ToAddress,
							Packaging:      "BAG",
							Pieces:         1,
							Weight: models.Weight{
								Units: "LB",
								Value: 1,
							},
							TotalHandlingUnits: 1,
							JustOneMore:        false,
							Description:        "",
						},
					},
				},
				PackageCount:         1,
				CarrierCode:          "FDXE",
				Remarks:              "",
				CommodityDescription: "",
			},
		},
	}
}

func pickupTime(pickupAddress models.Address, numDaysToDelay int) time.Time {
	location, err := toLocation(pickupAddress)
	if err != nil {
		location = laTimeZone
	}

	pickupTime := time.Now().In(location)

	if pickupTime.Hour() >= 12 {
		// If it's past 12pm, ship the next day, not today
		pickupTime.Add(24 * time.Hour)
	}

	pickupTime.Add(time.Duration(numDaysToDelay*24) * time.Hour)

	// Don't schedule pickups for Saturday or Sunday
	if pickupTime.Weekday() == time.Saturday {
		pickupTime.Add(48 * time.Hour)
	} else if pickupTime.Weekday() == time.Sunday {
		pickupTime.Add(24 * time.Hour)
	}

	year, month, day := pickupTime.Date()
	return time.Date(year, month, day, 12, 0, 0, 0, location)
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
	case "CT", "DE", "FL", "GA", "IN", "ME", "MD", "MA", "MI", "NH", "NJ", "NY", "NC", "OH", "PA", "RI", "SC", "VT", "VA", "WV":
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
