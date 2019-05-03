package fedex

import (
	"fmt"
	"time"

	"github.com/happyreturns/fedex/models"
)

func (f Fedex) createProcessShipmentRequest(shipment *Shipment) (models.Envelope, error) {

	customsClearanceDetail, err := f.customsClearanceDetail(shipment)
	if err != nil {
		return models.Envelope{}, fmt.Errorf("get customs clearance detail: %s", err) // TODO test this error
	}

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
			DropoffType:   dropoffType(shipment),
			ServiceType:   serviceType(shipment),
			PackagingType: "YOUR_PACKAGING",
			Shipper: models.Shipper{
				AccountNumber: f.Account,
				Address:       shipment.FromAddress,
				Contact:       shipment.FromContact,
			},
			Recipient: models.Shipper{
				AccountNumber: f.Account,
				Address:       shipment.ToAddress,
				Contact:       shipment.ToContact,
			},
			ShippingChargesPayment: models.Payment{
				PaymentType: "SENDER",
				Payor: models.Payor{
					ResponsibleParty: models.ResponsibleParty{
						AccountNumber: f.Account,
					},
				},
			},
			SmartPostDetail:               f.smartPostDetail(shipment),
			SpecialServicesRequested:      specialServicesRequested(shipment),
			CustomsClearanceDetail:        customsClearanceDetail,
			LabelSpecification:            labelSpecification(shipment),
			ShippingDocumentSpecification: shippingDocumentSpecification(shipment),
			PackageCount:                  1,
			RequestedPackageLineItems:     requestedPackageLineItems(shipment),
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

func (f Fedex) smartPostDetail(shipment *Shipment) *models.SmartPostDetail {
	if serviceType(shipment) == "SMART_POST" {
		return &models.SmartPostDetail{
			Indicia:              "PARCEL_RETURN",
			AncillaryEndorsement: "ADDRESS_CORRECTION",
			HubID:                f.HubID,
		}
	}
	return nil
}

func serviceType(shipment *Shipment) string {
	switch {
	case shipment.Service == "fedex_smart_post",
		shipment.Service == "return" && !isInternational(shipment):
		return "SMART_POST"
	default:
		return "FEDEX_GROUND"
	}
}

func shippingDocumentSpecification(shipment *Shipment) *models.ShippingDocumentSpecification {
	if serviceType(shipment) != "SMART_POST" && isInternational(shipment) {
		return &models.ShippingDocumentSpecification{
			ShippingDocumentTypes: []string{"COMMERCIAL_INVOICE"},
			CommercialInvoiceDetail: []models.CommercialInvoiceDetail{
				{
					Format: models.Format{
						ImageType: "PDF",
						StockType: "PAPER_LETTER",
					},
					CustomerImageUsages: []models.CustomerImageUsage{
						{
							Type: "LETTER_HEAD",
							ID:   "IMAGE_1", // TODO
						},
						{
							Type: "SIGNATURE",
							ID:   "IMAGE_2", // TODO
						},
					},
				},
			},
		}
	}
	return nil
}

func labelSpecification(shipment *Shipment) models.LabelSpecification {
	if serviceType(shipment) == "FEDEX_GROUND" && isInternational(shipment) {
		stockType := "PAPER_4X6"
		return models.LabelSpecification{
			LabelFormatType: "COMMON2D",
			ImageType:       "PDF",
			LabelStockType:  &stockType,
		}

	}
	return models.LabelSpecification{
		LabelFormatType: "COMMON2D",
		ImageType:       "PNG",
	}
}

func dropoffType(shipment *Shipment) string {
	if isInternational(shipment) {
		return "BUSINESS_SERVICE_CENTER"
	}
	return "REGULAR_PICKUP"
}

func isInternational(shipment *Shipment) bool {
	fromCountryCode := shipment.FromAddress.CountryCode
	if fromCountryCode == "" {
		fromCountryCode = "US"
	}

	toCountryCode := shipment.ToAddress.CountryCode
	if toCountryCode == "" {
		toCountryCode = "US"
	}

	return fromCountryCode != toCountryCode
}

func weight(shipment *Shipment) models.Weight {
	if isInternational(shipment) && len(shipment.Commodities) > 0 {
		weight := models.Weight{
			Units: shipment.Commodities[0].Weight.Units,
			Value: 0.0,
		}
		for _, commodity := range shipment.Commodities {
			weight.Value += commodity.Weight.Value
		}
		return weight
	}

	switch serviceType(shipment) {
	case "SMART_POST":
		return models.Weight{Units: "LB", Value: 0.99}
	default:
		return models.Weight{Units: "LB", Value: 13}
	}
}

func dimensions(shipment *Shipment) models.Dimensions {
	switch serviceType(shipment) {
	case "SMART_POST":
		return models.Dimensions{Length: 6, Width: 5, Height: 5, Units: "IN"}
	default:
		return models.Dimensions{Length: 13, Width: 13, Height: 13, Units: "IN"}
	}
}

func defaultEventNotificationDetail(notificationEmail string) *models.EventNotificationDetail {
	return &models.EventNotificationDetail{
		AggregationType: "PER_SHIPMENT",
		EventNotifications: []models.EventNotification{{
			Role: "SHIPPER",
			Events: []string{
				"ON_DELIVERY",
				"ON_ESTIMATED_DELIVERY",
				"ON_EXCEPTION",
				"ON_SHIPMENT",
				"ON_TENDER",
			},
			NotificationDetail: models.NotificationDetail{
				NotificationType: "EMAIL",
				EmailDetail: models.EmailDetail{
					EmailAddress: notificationEmail,
					Name:         "TEST NAME",
				},
				Localization: models.Localization{
					LanguageCode: "en",
				},
			},
			FormatSpecification: models.FormatSpecification{
				Type: "HTML",
			},
		}},
	}
}

func specialServicesRequested(shipment *Shipment) *models.SpecialServicesRequested {
	var specialServicesRequested *models.SpecialServicesRequested
	switch serviceType(shipment) {
	case "SMART_POST":
		specialServicesRequested = &models.SpecialServicesRequested{
			SpecialServiceTypes: []string{"RETURN_SHIPMENT"},
			ReturnShipmentDetail: &models.ReturnShipmentDetail{
				ReturnType: "PRINT_RETURN_LABEL",
			},
		}

		if shipment.NotificationEmail != "" {
			specialServicesRequested.EventNotificationDetail = defaultEventNotificationDetail(shipment.NotificationEmail)
		}
	default:
		if isInternational(shipment) {
			// TODO notifications for international shipments?
			specialServicesRequested = &models.SpecialServicesRequested{
				SpecialServiceTypes: []string{"ELECTRONIC_TRADE_DOCUMENTS"},
				EtdDetail: &models.EtdDetail{
					RequestedDocumentCopies: "COMMERCIAL_INVOICE",
				},
			}
		} else if shipment.NotificationEmail != "" {
			specialServicesRequested = &models.SpecialServicesRequested{
				SpecialServiceTypes:     []string{"EVENT_NOTIFICATION"},
				EventNotificationDetail: defaultEventNotificationDetail(shipment.NotificationEmail),
			}
		}
	}
	return specialServicesRequested
}

func customerReference(shipment *Shipment) models.CustomerReference {
	switch serviceType(shipment) {
	case "SMART_POST":
		return models.CustomerReference{
			CustomerReferenceType: "RMA_ASSOCIATION",
			Value: shipment.Reference,
		}
	default:
		return models.CustomerReference{
			CustomerReferenceType: "CUSTOMER_REFERENCE",
			Value: shipment.Reference,
		}
	}
}

func (f Fedex) customsClearanceDetail(shipment *Shipment) (*models.CustomsClearanceDetail, error) {
	if !isInternational(shipment) {
		return nil, nil
	}

	customsValue, err := customsValue(shipment.Commodities)
	if err != nil {
		return nil, fmt.Errorf("customs value: %s", err)
	}

	return &models.CustomsClearanceDetail{
		DutiesPayment: models.Payment{
			PaymentType: "SENDER",
			Payor: models.Payor{
				ResponsibleParty: models.ResponsibleParty{
					AccountNumber: f.Account,
				},
			},
		},
		CustomsValue: customsValue,
		Commodities:  shipment.Commodities,
	}, nil
}

func customsValue(commodities []models.Commodity) (models.Money, error) {
	// TODO eventually we will call an endpoint to calculate the
	// customsValueAmount when the sum value of all items becomes greater than
	// $800
	customsValue := models.Money{Currency: "USD"}

	if len(commodities) == 0 {
		return customsValue, nil
	}

	customsValue.Currency = commodities[0].CustomsValue.Currency
	for _, commodity := range commodities {
		if commodity.CustomsValue.Currency != customsValue.Currency {
			return customsValue, fmt.Errorf("mismatching customs currencies: %s %s", commodity.CustomsValue.Currency, customsValue.Currency)
		}
		customsValue.Amount += commodity.CustomsValue.Amount
	}
	return customsValue, nil

}

func requestedPackageLineItems(shipment *Shipment) []models.RequestedPackageLineItem {
	return []models.RequestedPackageLineItem{{
		SequenceNumber:     1,
		PhysicalPackaging:  "BAG",
		ItemDescription:    "ItemDescription",
		CustomerReferences: []models.CustomerReference{customerReference(shipment)},
		Weight:             weight(shipment),
		Dimensions:         dimensions(shipment),
	}}
}
