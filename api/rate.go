package api

import (
	"fmt"
	"time"

	"github.com/happyreturns/fedex/models"
)

const (
	rateVersion = "v24"
)

func (a API) Rate(rate *models.Rate) (*models.RateReply, error) {

	endpoint := fmt.Sprintf("/rate/%s", rateVersion)
	request := a.rateRequest(rate)
	response := &models.RateResponseEnvelope{}

	err := a.makeRequestAndUnmarshalResponse(endpoint, request, response)
	if err != nil {
		return nil, fmt.Errorf("make rate request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) rateRequest(rate *models.Rate) *models.Envelope {
	rateRequestTypes := "PREFERRED"
	packageCount := 1

	// When the service type is smartpost, getting rates from FedEx API doesn't
	// work
	serviceType := rate.ServiceType()
	serviceTypeInRequest := serviceType
	if serviceType == "SMART_POST" {
		// TODO figure out why this is necessary. We aren't getting back smartpost
		// rates. So using ground instead here.
		serviceTypeInRequest = "FEDEX_GROUND"
	}

	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: fmt.Sprintf("http://fedex.com/ws/rate/%s", rateVersion),
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
						CustomerTransactionID: "Rate Request",
					},
					Version: models.Version{
						ServiceID: "crs",
						Major:     24,
					},
				},
				RequestedShipment: models.RequestedShipment{
					ShipTimestamp:     models.Timestamp(time.Now()),
					DropoffType:       "REGULAR_PICKUP",
					ServiceType:       serviceTypeInRequest,
					PackagingType:     "YOUR_PACKAGING",
					PreferredCurrency: "USD",
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
							ResponsibleParty: models.Shipper{
								AccountNumber: a.Account,
							},
						},
					},
					SmartPostDetail: a.SmartPostDetail(serviceType),
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
							Weight:            rate.Weight(),
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
									Value:                 "NAFTA_COO",
								},
							},
						},
					},
				},
			},
		},
	}
}
