package api

import (
	"fmt"

	"time"

	"github.com/happyreturns/fedex/models"
)

func (a API) ProcessShipment(shipment *models.Shipment) (*models.ProcessShipmentReply, error) {
	request, err := a.processShipmentRequest(shipment)
	if err != nil {
		return nil, fmt.Errorf("create process shipment request: %s", err)
	}

	response := &models.ShipResponseEnvelope{}
	if err := a.makeRequestAndUnmarshalResponse("/ship/v23", request, response); err != nil {
		return nil, fmt.Errorf("make process shipment request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) processShipmentRequest(shipment *models.Shipment) (models.Envelope, error) {
	customsClearanceDetail, err := a.customsClearanceDetail(shipment)
	if err != nil {
		return models.Envelope{}, fmt.Errorf("customs clearance detail: %s", err)
	}

	packageCount := 1
	req := models.ProcessShipmentRequest{
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
				ServiceID: "ship",
				Major:     23,
			},
		},
		RequestedShipment: models.RequestedShipment{
			ShipTimestamp: models.Timestamp(time.Now()),
			DropoffType:   shipment.DropoffType(),
			ServiceType:   shipment.ServiceType(),
			PackagingType: "YOUR_PACKAGING",
			Shipper: models.Shipper{
				AccountNumber: a.Account,
				Address:       shipment.FromAddress,
				Contact:       shipment.FromContact,
			},
			Recipient: models.Shipper{
				AccountNumber: a.Account,
				Address:       shipment.ToAddress,
				Contact:       shipment.ToContact,
			},
			ShippingChargesPayment: &models.Payment{
				PaymentType: "SENDER",
				Payor: models.Payor{
					ResponsibleParty: models.ResponsibleParty{
						AccountNumber: a.Account,
					},
				},
			},
			SmartPostDetail:               a.SmartPostDetail(shipment),
			SpecialServicesRequested:      shipment.SpecialServicesRequested(),
			CustomsClearanceDetail:        customsClearanceDetail,
			LabelSpecification:            shipment.LabelSpecification(),
			ShippingDocumentSpecification: shipment.ShippingDocumentSpecification(),
			PackageCount:                  &packageCount,
			RequestedPackageLineItems:     shipment.RequestedPackageLineItems(),
		},
	}

	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/ship/v23",
		Body: models.ProcessShipmentBody{
			ProcessShipmentRequest: req,
		},
	}, nil
}

func (a API) SmartPostDetail(shipment *models.Shipment) *models.SmartPostDetail {
	if shipment.ServiceType() == "SMART_POST" {
		return &models.SmartPostDetail{
			Indicia:              "PARCEL_RETURN",
			AncillaryEndorsement: "ADDRESS_CORRECTION",
			HubID:                a.HubID,
		}
	}
	return nil
}

func (a API) customsClearanceDetail(shipment *models.Shipment) (*models.CustomsClearanceDetail, error) {
	if !shipment.IsInternational() {
		return nil, nil // TODO is this weird
	}

	customsValue, err := shipment.Commodities.CustomsValue()
	if err != nil {
		return nil, fmt.Errorf("got error: %s", err)
	}

	return &models.CustomsClearanceDetail{
		DutiesPayment: models.Payment{
			PaymentType: "SENDER",
			Payor: models.Payor{
				ResponsibleParty: models.ResponsibleParty{
					AccountNumber: a.Account,
				},
			},
		},
		CustomsValue:                   &customsValue,
		Commodities:                    shipment.Commodities,
		PartiesToTransactionAreRelated: false,
		CommercialInvoice: &models.CommercialInvoice{
			Purpose:        "REPAIR_AND_RETURN",
			OriginatorName: shipment.OriginatorName,
		},
	}, nil
}
