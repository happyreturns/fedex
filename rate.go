package fedex

import (
	"time"

	"github.com/happyreturns/fedex/models"
)

func (f Fedex) rateRequest(rate *RateRequest) models.Envelope {
	if !rate.IsInternational() {
		return f.rateRequestDomestic(rate)
	}
	return f.rateRequestInternational(rate)
}

func (f Fedex) rateRequestDomestic(rate *RateRequest) models.Envelope {
	rateRequestTypes := "LIST"
	packageCount := 1
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/rate/v24",
		Body: struct {
			RateRequest models.RateRequest `xml:"q0:RateRequest"`
		}{
			RateRequest: models.RateRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      f.Key,
							Password: f.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: f.Account,
						MeterNumber:   f.Meter,
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
						AccountNumber: f.Account,
						Address:       rate.FromAndTo.FromAddress,
						Contact:       rate.FromAndTo.FromContact,
					},
					Recipient: models.Shipper{
						AccountNumber: f.Account,
						Address:       rate.FromAndTo.ToAddress,
						Contact:       rate.FromAndTo.ToContact,
					},
					ShippingChargesPayment: &models.Payment{
						PaymentType: "SENDER",
						Payor: models.Payor{
							ResponsibleParty: models.ResponsibleParty{
								AccountNumber: f.Account,
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

func (f Fedex) rateRequestInternational(rate *RateRequest) models.Envelope {

	// TODO check 800 or make different explicit function

	documentContent := "NON_DOCUMENTS"
	customsValue, err := rate.Commodities.CustomsValue()
	if err != nil {
		// TODO do something
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
		Body: struct {
			RateRequest models.RateRequest `xml:"q0:RateRequest"`
		}{
			RateRequest: models.RateRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      f.Key,
							Password: f.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: f.Account,
						MeterNumber:   f.Meter,
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
					ServiceType:   "FEDEX_GROUND", // TODO needed?
					PackagingType: "YOUR_PACKAGING",
					Shipper: models.Shipper{
						AccountNumber: f.Account,
						Address:       rate.FromAndTo.FromAddress,
						Contact:       rate.FromAndTo.FromContact,
					},
					Recipient: models.Shipper{
						AccountNumber: f.Account,
						Address:       rate.FromAndTo.ToAddress,
						Contact:       rate.FromAndTo.ToContact,
					},
					CustomsClearanceDetail: &models.CustomsClearanceDetail{
						DutiesPayment: models.Payment{
							PaymentType: "SENDER",
							Payor: models.Payor{
								ResponsibleParty: models.ResponsibleParty{
									AccountNumber: f.Account,
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
