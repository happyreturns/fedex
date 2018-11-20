package fedex

import (
	"time"

	"github.com/happyreturns/fedex/models"
)

func (f Fedex) shipmentEnvelope(shipmentType string, fromLocation, toLocation models.Address, fromContact, toContact models.Contact) models.Envelope {

	req := models.ProcessShipmentRequest{
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
			Version: models.Version{
				ServiceID: "ship",
				Major:     23,
			},
		},
		RequestedShipment: models.RequestedShipment{
			ShipTimestamp: models.Timestamp(time.Now()),
			DropoffType:   "REGULAR_PICKUP",
			// ServiceType:   "FEDEX_GROUND",
			PackagingType: "YOUR_PACKAGING",
			Shipper: models.Shipper{
				AccountNumber: f.Account,
				Address:       fromLocation,
				Contact:       fromContact,
			},
			Recipient: models.Shipper{
				AccountNumber: f.Account,
				Address:       toLocation,
				Contact:       toContact,
			},
			ShippingChargesPayment: models.Payment{
				PaymentType: "SENDER",
				Payor: models.Payor{
					ResponsibleParty: models.ResponsibleParty{
						AccountNumber: f.Account,
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
					SequenceNumber:    1,
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
	}

	switch shipmentType {
	case "SMART_POST":
		req.RequestedShipment.ServiceType = "SMART_POST"
		req.RequestedShipment.SmartPostDetail = &models.SmartPostDetail{
			Indicia:              "PRESORTED_STANDARD",
			AncillaryEndorsement: "ADDRESS_CORRECTION",
			HubID:                f.SmartPostHubID,
		}
		req.RequestedShipment.RequestedPackageLineItems[0].Weight = models.Weight{
			Units: "LB",
			Value: 0.99,
		}
		req.RequestedShipment.RequestedPackageLineItems[0].Dimensions = models.Dimensions{
			Length: 6,
			Width:  4,
			Height: 1,
			Units:  "IN",
		}
	default:
		req.RequestedShipment.ServiceType = "FEDEX_GROUND"
		req.RequestedShipment.RequestedPackageLineItems[0].Weight = models.Weight{
			Units: "LB",
			Value: 40,
		}
		req.RequestedShipment.RequestedPackageLineItems[0].Dimensions = models.Dimensions{
			Length: 5,
			Width:  5,
			Height: 5,
			Units:  "IN",
		}
	}

	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/ship/v23",
		Body: struct {
			ProcessShipmentRequest models.ProcessShipmentRequest `xml:"q0:ProcessShipmentRequest"`
		}{
			ProcessShipmentRequest: req,
		},
	}
}

func (f Fedex) shipGroundSOAPRequest(fromLocation, toLocation models.Address, fromContact, toContact models.Contact) models.Envelope {
	return f.shipmentEnvelope("FEDEX_GROUND", fromLocation, toLocation, fromContact, toContact)
}

func (f Fedex) shipSmartPostSOAPRequest(fromLocation, toLocation models.Address, fromContact, toContact models.Contact) models.Envelope {
	return f.shipmentEnvelope("SMART_POST", fromLocation, toLocation, fromContact, toContact)
}
