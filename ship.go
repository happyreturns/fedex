package fedex

import (
	"fmt"
	"time"

	"github.com/happyreturns/fedex/models"
)

func processShipmentRequest(fedex Fedex, body string) string {
	return fedex.wrapSoapRequest(fmt.Sprintf(`
		<q0:ProcessShipmentRequest>
			%s
			%s
		</q0:ProcessShipmentRequest>
	`, fedex.soapCreds("ship", "23"), body), "http://fedex.com/ws/ship/v23")
}

func soapShipGround(fedex Fedex, fromLocation, toLocation models.Address, fromContact, toContact models.Contact) models.Envelope {
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/ship/v23",
		Body: struct {
			ProcessShipmentRequest models.ProcessShipmentRequest `xml:"q0:ProcessShipmentRequest"`
		}{
			ProcessShipmentRequest: models.ProcessShipmentRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      fedex.Key,
							Password: fedex.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: fedex.Account,
						MeterNumber:   fedex.Meter,
					},
					Version: models.Version{
						ServiceID: "ship",
						Major:     23,
					},
				},
				RequestedShipment: models.RequestedShipment{
					ShipTimestamp: models.Timestamp(time.Now()),
					DropoffType:   "REGULAR_PICKUP",
					ServiceType:   "FEDEX_GROUND",
					PackagingType: "YOUR_PACKAGING",
					Shipper: models.Shipper{
						AccountNumber: fedex.Account,
						Address:       fromLocation,
						Contact:       fromContact,
					},
					Recipient: models.Shipper{
						AccountNumber: fedex.Account,
						Address:       toLocation,
						Contact:       toContact,
					},
					ShippingChargesPayment: models.ShippingChargesPayment{
						PaymentType: "SENDER",
						Payor: models.Payor{
							ResponsibleParty: models.ResponsibleParty{
								AccountNumber: fedex.Account,
							},
						},
					},
					LabelSpecification: models.LabelSpecification{
						LabelFormatType: "COMMON2D",
						ImageType:       "PNG",
					},
					RateRequestTypes: "LIST",
					PackageCount:     1,
					RequestedPackageLineItems: []models.RequestedPackageLineItem{
						{
							SequenceNumber: 1,
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
