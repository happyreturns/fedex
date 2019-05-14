package api

import (
	"fmt"
	"time"

	"github.com/happyreturns/fedex/models"
)

func (a API) Rate(rate *models.Rate) (*models.RateReply, error) {

	request := a.rateRequest(rate)
	response := &models.RateResponseEnvelope{}

	err := a.makeRequestAndUnmarshalResponse("/rate/v24", request, response)
	if err != nil {
		return nil, fmt.Errorf("make rate request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) rateRequest(rate *models.Rate) models.Envelope {
	rateRequestTypes := "LIST"
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
					ShippingChargesPayment: &models.Payment{
						PaymentType: "SENDER",
						Payor: models.Payor{
							ResponsibleParty: models.ResponsibleParty{
								AccountNumber: a.Account,
							},
						},
					},
					LabelSpecification: &models.LabelSpecification{
						LabelFormatType: "COMMON2D",
						ImageType:       "PDF",
					},
					RateRequestTypes: &rateRequestTypes,
					PackageCount:     &packageCount,
					RequestedPackageLineItems: []models.RequestedPackageLineItem{
						{
							SequenceNumber:    1,
							GroupPackageCount: 1,
							Weight: models.Weight{
								Units: "LB",
								Value: 40,
							},
							Dimensions: models.Dimensions{
								Length: 5,
								Width:  5,
								Height: 5,
								Units:  "IN",
							},
							PhysicalPackaging: "BAG",
							ItemDescription:   "Stuff",
							CustomerReferences: []models.CustomerReference{
								{
									CustomerReferenceType: "CUSTOMER_REFERENCE",
									Value: "NAFTA_COO",
								},
							},
						},
					},
				},
			},
		},
	}
}
