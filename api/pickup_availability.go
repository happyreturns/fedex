package api

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

func (a API) PickupAvailability(pickup *models.Pickup) (*models.PickupAvailabilityReply, error) {
	endpoint := fmt.Sprintf("/pickup/%s", createPickupVersion)
	request := a.pickupAvailabilityRequest(pickup)
	response := &models.PickupAvailabilityResponseEnvelope{}

	if err := a.makeRequestAndUnmarshalResponse(endpoint, request, response); err != nil {
		return nil, fmt.Errorf("make pickup availability request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) pickupAvailabilityRequest(pickup *models.Pickup) *models.Envelope {
	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: fmt.Sprintf("http://fedex.com/ws/pickup/%s", createPickupVersion),
		Body: models.PickupAvailabilityBody{
			PickupAvailabilityRequest: models.PickupAvailabilityRequest{
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
				PickupRequestType:    []string{"FUTURE_DAY", "SAME_DAY"},
				PickupAddress:        pickup.PickupLocation.Address,
				NumberOfBusinessDays: 3,
				Carriers:             []string{"FDXG"},
			},
		},
	}
}
