package models

// TODO is this needed? i don't think so if the rate itself doesn't decide to do stuff differently for international (rather f Fedex does)
type FromAndTo struct {
	FromAddress Address
	ToAddress   Address
	FromContact Contact
	ToContact   Contact
}

func (ft FromAndTo) IsInternational() bool {
	fromCountryCode := ft.FromAddress.CountryCode
	if fromCountryCode == "" {
		fromCountryCode = "US"
	}

	toCountryCode := ft.ToAddress.CountryCode
	if toCountryCode == "" {
		toCountryCode = "US"
	}

	return fromCountryCode != toCountryCode
}

// Rate wraps all the Fedex API fields needed for getting a rate
type Rate struct {
	FromAndTo

	// Only used for international ground shipments
	Commodities Commodities
}

// Shipment wraps all the Fedex API fields needed for creating a shipment
type Shipment struct {
	FromAndTo

	NotificationEmail string
	Reference         string
	Service           string

	// Only used for international ground shipments
	Commodities Commodities
}

func (s *Shipment) ServiceType() string {
	switch {
	case s.Service == "fedex_smart_post",
		s.Service == "return" && !s.IsInternational():
		// TODO throw error if smart_post account, international?
		return "SMART_POST"
	default:
		return "FEDEX_GROUND"
	}
}

func (s *Shipment) ShippingDocumentSpecification() *ShippingDocumentSpecification {
	if s.ServiceType() != "SMART_POST" && s.IsInternational() {
		return &ShippingDocumentSpecification{
			ShippingDocumentTypes: []string{"COMMERCIAL_INVOICE"},
			CommercialInvoiceDetail: []CommercialInvoiceDetail{
				{
					Format: Format{
						ImageType: "PDF",
						StockType: "PAPER_LETTER",
					},
					CustomerImageUsages: []CustomerImageUsage{
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

func (s *Shipment) LabelSpecification() *LabelSpecification {
	if s.ServiceType() == "FEDEX_GROUND" && s.IsInternational() {
		stockType := "PAPER_4X6"
		return &LabelSpecification{
			LabelFormatType: "COMMON2D",
			ImageType:       "PDF",
			LabelStockType:  &stockType,
		}

	}
	return &LabelSpecification{
		LabelFormatType: "COMMON2D",
		ImageType:       "PNG",
	}
}

func (s *Shipment) DropoffType() string {
	if s.IsInternational() {
		return "BUSINESS_SERVICE_CENTER"
	}
	return "REGULAR_PICKUP"
}

func (s *Shipment) Weight() Weight {
	commoditiesWeight := s.Commodities.Weight()
	if !commoditiesWeight.IsZero() {
		return commoditiesWeight
	}

	switch s.ServiceType() {
	case "SMART_POST":
		return Weight{Units: "LB", Value: 0.99}
	default:
		return Weight{Units: "LB", Value: 13}
	}
}

func (s *Shipment) Dimensions() Dimensions {
	switch s.ServiceType() {
	case "SMART_POST":
		return Dimensions{Length: 6, Width: 5, Height: 5, Units: "IN"}
	default:
		return Dimensions{Length: 13, Width: 13, Height: 13, Units: "IN"}
	}
}

func (s *Shipment) SpecialServicesRequested() *SpecialServicesRequested {
	var specialServicesRequested *SpecialServicesRequested
	switch s.ServiceType() {
	case "SMART_POST":
		specialServicesRequested = &SpecialServicesRequested{
			SpecialServiceTypes: []string{"RETURN_SHIPMENT"},
			ReturnShipmentDetail: &ReturnShipmentDetail{
				ReturnType: "PRINT_RETURN_LABEL",
			},
		}

		if s.NotificationEmail != "" {
			specialServicesRequested.EventNotificationDetail = defaultEventNotificationDetail(s.NotificationEmail)
		}
	default:
		if s.IsInternational() {
			// TODO notifications for international shipments?
			specialServicesRequested = &SpecialServicesRequested{
				SpecialServiceTypes: []string{"ELECTRONIC_TRADE_DOCUMENTS"},
				EtdDetail: &EtdDetail{
					RequestedDocumentCopies: "COMMERCIAL_INVOICE",
				},
			}
		} else if s.NotificationEmail != "" {
			specialServicesRequested = &SpecialServicesRequested{
				SpecialServiceTypes:     []string{"EVENT_NOTIFICATION"},
				EventNotificationDetail: defaultEventNotificationDetail(s.NotificationEmail),
			}
		}
	}
	return specialServicesRequested
}

func (s *Shipment) CustomerReference() CustomerReference {
	switch s.ServiceType() {
	case "SMART_POST":
		return CustomerReference{
			CustomerReferenceType: "RMA_ASSOCIATION",
			Value: s.Reference,
		}
	default:
		return CustomerReference{
			CustomerReferenceType: "CUSTOMER_REFERENCE",
			Value: s.Reference,
		}
	}
}

func defaultEventNotificationDetail(notificationEmail string) *EventNotificationDetail {
	return &EventNotificationDetail{
		AggregationType: "PER_SHIPMENT",
		EventNotifications: []EventNotification{{
			Role: "SHIPPER",
			Events: []string{
				"ON_DELIVERY",
				"ON_ESTIMATED_DELIVERY",
				"ON_EXCEPTION",
				"ON_SHIPMENT",
				"ON_TENDER",
			},
			NotificationDetail: NotificationDetail{
				NotificationType: "EMAIL",
				EmailDetail: EmailDetail{
					EmailAddress: notificationEmail,
					Name:         "TEST NAME",
				},
				Localization: Localization{
					LanguageCode: "en",
				},
			},
			FormatSpecification: FormatSpecification{
				Type: "HTML",
			},
		}},
	}
}

func (s *Shipment) RequestedPackageLineItems() []RequestedPackageLineItem {
	return []RequestedPackageLineItem{{
		SequenceNumber:     1,
		PhysicalPackaging:  "BAG",
		ItemDescription:    "ItemDescription",
		CustomerReferences: []CustomerReference{s.CustomerReference()},
		Weight:             s.Weight(),
		Dimensions:         s.Dimensions(),
	}}
}
