package api

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

func (a API) TrackByNumber(carrierCode, trackingNo string) (*models.TrackReply, error) {
	fmt.Printf("*** TrackByNumber, carrierCode: %+v\n", carrierCode)
	fmt.Printf("*** TrackByNumber, trackingNo: %+v\n", trackingNo)

	request := a.trackByNumberRequest(carrierCode, trackingNo)
	fmt.Printf("*** TrackByNumber, request: %+v\n", request)

	response := &models.TrackResponseEnvelope{}
	fmt.Printf("*** TrackByNumber, response: %+v\n", response)

	err := a.makeRequestAndUnmarshalResponse("/trck", request, response)
	if err != nil {
		fmt.Printf("*** TrackByNumber, err: %+v\n", err)
		return nil, fmt.Errorf("make track request and unmarshal: %s", err)
	}
	fmt.Printf("*** TrackByNumber, &response.Reply: %+v\n", &response.Reply)
	return &response.Reply, nil
}

func (a API) trackByNumberRequest(carrierCode string, trackingNo string) *models.Envelope {
	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/track/v16",
		Body: models.TrackBody{
			TrackRequest: models.TrackRequest{
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
						ServiceID: "trck",
						Major:     16,
					},
				},
				ProcessingOptions: "INCLUDE_DETAILED_SCANS",
				SelectionDetails: models.SelectionDetails{
					CarrierCode: carrierCode,
					PackageIdentifier: models.PackageIdentifier{
						Type:  "TRACKING_NUMBER_OR_DOORTAG",
						Value: trackingNo,
					},
				},
			},
		},
	}
}
