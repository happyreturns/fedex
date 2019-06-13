package api

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

func (a API) CancelPickup(pickupNumber string) (*models.CancelPickupReply, error) {
	request, err := a.cancelPickupRequest(pickupNumber)
	if err != nil {
		return nil, fmt.Errorf("cancel pickup request: %s", err)
	}

	endpoint := fmt.Sprintf("/pickup/%s", createPickupVersion)
	response := &models.CancelPickupResponseEnvelope{}
	err = a.makeRequestAndUnmarshalResponse(endpoint, request, response)
	if err != nil {
		return nil, fmt.Errorf("make cancel pickup request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) cancelPickupRequest(pickupNumber string) (*models.Envelope, error) {
	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: fmt.Sprintf("http://fedex.com/ws/pickup/%s", createPickupVersion),
		Body: models.CancelPickupBody{
			CancelPickupRequest: models.CancelPickupRequest{
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
				PickupConfirmationNumber: pickupNumber,
				CarrierCode:              "FDXG",
				Remarks:                  []string{"Accidentally made a pickup."},
			},
		},
	}, nil
}
