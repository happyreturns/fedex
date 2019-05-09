package fedex

import (
	"fmt"
	"time"

	"github.com/happyreturns/fedex/models"
)

func (f Fedex) createProcessShipmentRequest(shipment *Shipment) (models.Envelope, error) {

	packageCount := 1
	customsClearanceDetail, err := f.customsClearanceDetail(shipment)
	if err != nil {
		return models.Envelope{}, fmt.Errorf("customs clearance detail: %s", err)
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
			DropoffType:   shipment.dropoffType(),
			ServiceType:   shipment.serviceType(),
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
			ShippingChargesPayment: &models.Payment{
				PaymentType: "SENDER",
				Payor: models.Payor{
					ResponsibleParty: models.ResponsibleParty{
						AccountNumber: f.Account,
					},
				},
			},
			SmartPostDetail:               f.smartPostDetail(shipment),
			SpecialServicesRequested:      shipment.specialServicesRequested(),
			CustomsClearanceDetail:        customsClearanceDetail,
			LabelSpecification:            shipment.labelSpecification(),
			ShippingDocumentSpecification: shipment.shippingDocumentSpecification(),
			PackageCount:                  &packageCount,
			RequestedPackageLineItems:     shipment.requestedPackageLineItems(),
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
	if shipment.serviceType() == "SMART_POST" {
		return &models.SmartPostDetail{
			Indicia:              "PARCEL_RETURN",
			AncillaryEndorsement: "ADDRESS_CORRECTION",
			HubID:                f.HubID,
		}
	}
	return nil
}

func (s *Shipment) serviceType() string {
	switch {
	case s.Service == "fedex_smart_post",
		s.Service == "return" && !s.IsInternational():
		// TODO throw error if smart_post account, international?
		return "SMART_POST"
	default:
		return "FEDEX_GROUND"
	}
}

func (s *Shipment) shippingDocumentSpecification() *models.ShippingDocumentSpecification {
	if s.serviceType() != "SMART_POST" && s.IsInternational() {
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

func (s *Shipment) labelSpecification() *models.LabelSpecification {
	if s.serviceType() == "FEDEX_GROUND" && s.IsInternational() {
		stockType := "PAPER_4X6"
		return &models.LabelSpecification{
			LabelFormatType: "COMMON2D",
			ImageType:       "PDF",
			LabelStockType:  &stockType,
		}

	}
	return &models.LabelSpecification{
		LabelFormatType: "COMMON2D",
		ImageType:       "PNG",
	}
}

func (s *Shipment) dropoffType() string {
	if s.IsInternational() {
		return "BUSINESS_SERVICE_CENTER"
	}
	return "REGULAR_PICKUP"
}

func (s *Shipment) weight() models.Weight {
	commoditiesWeight := s.Commodities.Weight()
	if !commoditiesWeight.IsZero() {
		return commoditiesWeight
	}

	switch s.serviceType() {
	case "SMART_POST":
		return models.Weight{Units: "LB", Value: 0.99}
	default:
		return models.Weight{Units: "LB", Value: 13}
	}
}

func (s *Shipment) dimensions() models.Dimensions {
	switch s.serviceType() {
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

func (s *Shipment) specialServicesRequested() *models.SpecialServicesRequested {
	var specialServicesRequested *models.SpecialServicesRequested
	switch s.serviceType() {
	case "SMART_POST":
		specialServicesRequested = &models.SpecialServicesRequested{
			SpecialServiceTypes: []string{"RETURN_SHIPMENT"},
			ReturnShipmentDetail: &models.ReturnShipmentDetail{
				ReturnType: "PRINT_RETURN_LABEL",
			},
		}

		if s.NotificationEmail != "" {
			specialServicesRequested.EventNotificationDetail = defaultEventNotificationDetail(s.NotificationEmail)
		}
	default:
		if s.IsInternational() {
			// TODO notifications for international shipments?
			specialServicesRequested = &models.SpecialServicesRequested{
				SpecialServiceTypes: []string{"ELECTRONIC_TRADE_DOCUMENTS"},
				EtdDetail: &models.EtdDetail{
					RequestedDocumentCopies: "COMMERCIAL_INVOICE",
				},
			}
		} else if s.NotificationEmail != "" {
			specialServicesRequested = &models.SpecialServicesRequested{
				SpecialServiceTypes:     []string{"EVENT_NOTIFICATION"},
				EventNotificationDetail: defaultEventNotificationDetail(s.NotificationEmail),
			}
		}
	}
	return specialServicesRequested
}

func (s *Shipment) customerReference() models.CustomerReference {
	switch s.serviceType() {
	case "SMART_POST":
		return models.CustomerReference{
			CustomerReferenceType: "RMA_ASSOCIATION",
			Value: s.Reference,
		}
	default:
		return models.CustomerReference{
			CustomerReferenceType: "CUSTOMER_REFERENCE",
			Value: s.Reference,
		}
	}
}

func (f Fedex) customsClearanceDetail(shipment *Shipment) (*models.CustomsClearanceDetail, error) {
	if !shipment.IsInternational() {
		return nil, nil // TODO is this weird
	}

	customsValue, err := shipment.Commodities.CustomsValue()
	if err != nil {
		return nil, fmt.Errorf("got error: %s", err)
	}

	return &models.CustomsClearanceDetail{
		Brokers: []models.Broker{{
			Type: "IMPORT",
		}},
		DutiesPayment: models.Payment{
			PaymentType: "SENDER",
			Payor: models.Payor{
				ResponsibleParty: models.ResponsibleParty{
					AccountNumber: f.Account,
				},
			},
		},
		CustomsValue:                   &customsValue,
		Commodities:                    shipment.Commodities,
		PartiesToTransactionAreRelated: false,
		CommercialInvoice: &models.CommercialInvoice{
			Purpose: "REPAIR_AND_RETURN",
		},
	}, nil
}

func (s *Shipment) requestedPackageLineItems() []models.RequestedPackageLineItem {
	return []models.RequestedPackageLineItem{{
		SequenceNumber:     1,
		PhysicalPackaging:  "BAG",
		ItemDescription:    "ItemDescription",
		CustomerReferences: []models.CustomerReference{s.customerReference()},
		Weight:             s.weight(),
		Dimensions:         s.dimensions(),
	}}
}
