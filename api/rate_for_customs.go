package api

import (
	"fmt"
	"time"

	"github.com/happyreturns/fedex/models"
)

func (a API) RateForCustoms(rate *models.Rate) (*models.RateReply, error) {

	request := a.rateForCustomsRequest(rate)
	response := &models.RateResponseEnvelope{}

	err := a.makeRequestAndUnmarshalResponse("/rate/v24", request, response)
	if err != nil {
		return nil, fmt.Errorf("make rate request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) rateForCustomsRequest(rate *models.Rate) models.Envelope {
	documentContent := "NON_DOCUMENTS"
	customsValue, err := rate.Commodities.CustomsValue()
	if err != nil {
		customsValue = models.Money{Currency: "USD"}
	}

	weight := rate.Commodities.Weight()
	if weight.IsZero() {
		weight = models.Weight{
			Units: "LB",
			Value: 40,
		}
	}

	edtRequestType := "ALL"
	packageCount := 1

	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/rate/v24",
		Body: models.RateBody{
			RateRequest: models.RateRequest{
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
					TransactionDetail: &models.TransactionDetail{
						CustomerTransactionID: "RAS Example",
					},
					Version: models.Version{
						ServiceID: "crs",
						Major:     24,
					},
				},
				RequestedShipment: models.RequestedShipment{
					ShipTimestamp: models.Timestamp(time.Now()),
					DropoffType:   "REGULAR_PICKUP",
					ServiceType:   "FEDEX_GROUND",
					PackagingType: "YOUR_PACKAGING",
					Shipper: models.Shipper{
						AccountNumber: a.Account,
						Address:       rate.FromAndTo.FromAddress,
						Contact:       rate.FromAndTo.FromContact,
					},
					Recipient: models.Shipper{
						AccountNumber: a.Account,
						Address:       rate.FromAndTo.ToAddress,
						Contact:       rate.FromAndTo.ToContact,
					},
					CustomsClearanceDetail: &models.CustomsClearanceDetail{
						DutiesPayment: models.Payment{
							PaymentType: "SENDER",
							Payor: models.Payor{
								ResponsibleParty: models.ResponsibleParty{
									AccountNumber: a.Account,
								},
							},
						},
						DocumentContent: &documentContent,
						CustomsValue:    &customsValue,
						Commodities:     rate.Commodities,
					},
					EdtRequestType: &edtRequestType,
					PackageCount:   &packageCount,
					RequestedPackageLineItems: []models.RequestedPackageLineItem{
						{
							SequenceNumber:    1,
							GroupPackageCount: 1,
							Weight:            weight,
							Dimensions: models.Dimensions{
								Length: 6,
								Width:  5,
								Height: 5,
								Units:  "IN",
							},
							PhysicalPackaging: "BAG",
							ItemDescription:   "Stuff",
						},
					},
				},
			},
		},
	}
}
