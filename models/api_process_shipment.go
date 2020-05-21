package models

import (
	"errors"
	"regexp"
	"time"
)

// Shipment wraps all the Fedex API fields needed for creating a shipment
type Shipment struct {
	FromAndTo

	NotificationEmail string
	References        []string
	Service           string
	Dimensions        Dimensions

	// Only used for international ground shipments
	OriginatorName    string
	Commodities       Commodities
	LetterheadImageID string
}

var (
	nonAlphanumericRegex *regexp.Regexp
)

func init() {
	nonAlphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")
}

func ServiceType(fromAndTo FromAndTo, service string) string {
	// TODO This is confusing. If the service is marked as "fedex_smart_post" or
	// "fedex_international_economy" (this is done through the CMS), then
	// explicitly set the service type as "SMART_POST" or "INTERNATIONAL_ECONOMY"
	// respectively. Otherwise, we deduce the service type based on whether
	// the service was "return", whether the return is international, and where
	// the return is coming from. In the future, we should just not allow using
	// anything other than "fedex_smart_post", "fedex_international_economy",
	// "fedex_ground" and make the CMS user to be explicit. However currently
	// there are many shipping methods that depend on this deduction logic.
	switch {
	case s.Service == "fedex_smart_post",
		s.Service == "return" && !s.IsInternational():
		return "SMART_POST"
	case s.Service == "fedex_international_economy" ||
		(s.IsInternational() && s.FromAddress.ShipsOutWithInternationalEconomy()):
		return "INTERNATIONAL_ECONOMY"
	default:
		return "FEDEX_GROUND"
	}
}

func (s *Shipment) Broker() string {
	switch s.ServiceType() {
	case "INTERNATIONAL_ECONOMY":
		return "FedEx Express"
	default:
		return "FedEx Logistics"
	}
}

func (s *Shipment) ShipTime() time.Time {
	t := time.Now()
	if s.IsInternational() {
		t = t.Add(9 * 24 * time.Hour)
	}

	return t
}

func (s *Shipment) ShippingDocumentSpecification() *ShippingDocumentSpecification {
	if s.ServiceType() == "SMART_POST" || !s.IsInternational() {
		return nil
	}

	letterheadImageID := s.LetterheadImageID
	if s.LetterheadImageID == "" {
		letterheadImageID = "IMAGE_1"
	}

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
						ID:   letterheadImageID,
					},
					{
						Type: "SIGNATURE",
						ID:   "IMAGE_2",
					},
				},
			},
		},
	}
}

func (s *Shipment) LabelSpecification() *LabelSpecification {
	if s.IsInternational() {
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
	// TODO ask if this is the intended behavior.
	// So for getting rates, we sup
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

func (s *Shipment) ValidatedDimensions() Dimensions {
	if s.Dimensions.IsValid() {
		return s.Dimensions
	}

	switch s.ServiceType() {
	case "SMART_POST":
		return Dimensions{Length: 6, Width: 5, Height: 5, Units: "IN"}
	default:
		return Dimensions{Length: 13, Width: 13, Height: 13, Units: "IN"}
	}
}

func (s *Shipment) SpecialServicesRequested() *SpecialServicesRequested {
	var (
		specialServiceTypes []string

		etdDetail               *EtdDetail
		eventNotificationDetail *EventNotificationDetail
		returnShipmentDetail    *ReturnShipmentDetail
	)

	if s.ServiceType() == "SMART_POST" {
		specialServiceTypes = append(specialServiceTypes, "RETURN_SHIPMENT")
		returnShipmentDetail = &ReturnShipmentDetail{
			ReturnType: "PRINT_RETURN_LABEL",
		}
	}

	if s.IsInternational() {
		specialServiceTypes = append(specialServiceTypes, "ELECTRONIC_TRADE_DOCUMENTS")
		etdDetail = &EtdDetail{
			RequestedDocumentCopies: "COMMERCIAL_INVOICE",
		}
	}

	if s.NotificationEmail != "" {
		specialServiceTypes = append(specialServiceTypes, "EVENT_NOTIFICATION")
		eventNotificationDetail = defaultEventNotificationDetail(s.NotificationEmail)
	}

	if len(specialServiceTypes) == 0 {
		return nil
	}
	return &SpecialServicesRequested{
		SpecialServiceTypes: specialServiceTypes,

		EtdDetail:               etdDetail,
		EventNotificationDetail: eventNotificationDetail,
		ReturnShipmentDetail:    returnShipmentDetail,
	}
}

func (s *Shipment) CustomerReferences() []CustomerReference {
	customerReferences := make([]CustomerReference, len(s.References))
	for idx, reference := range s.References {
		switch s.ServiceType() {
		case "SMART_POST":
			customerReferences[idx] = CustomerReference{
				CustomerReferenceType: "RMA_ASSOCIATION",
				Value:                 sanitizeReferenceForFedexAPI(reference),
			}
		default:
			customerReferences[idx] = CustomerReference{
				CustomerReferenceType: "CUSTOMER_REFERENCE",
				Value:                 sanitizeReferenceForFedexAPI(reference),
			}
		}
	}
	return customerReferences
}

func sanitizeReferenceForFedexAPI(reference string) string {

	// Remove non-alphanumeric chars
	validatedReference := nonAlphanumericRegex.ReplaceAllString(reference, "")

	// Trim length
	if len(validatedReference) > 20 {
		validatedReference = validatedReference[0:20]
	}

	return validatedReference
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
					Name:         "Happy Returns dev team",
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
		CustomerReferences: s.CustomerReferences(),
		Weight:             s.Weight(),
		Dimensions:         s.ValidatedDimensions(),
	}}
}

type ProcessShipmentBody struct {
	ProcessShipmentRequest ProcessShipmentRequest `xml:"q0:ProcessShipmentRequest"`
}

type ProcessShipmentRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type ShipResponseEnvelope struct {
	Reply ProcessShipmentReply `xml:"Body>ProcessShipmentReply"`
}

func (s *ShipResponseEnvelope) Error() error {
	return s.Reply.Error()
}

// ProcessShipReply : Process shipment reply root (`xml:"Body>ProcessShipmentReply"`)
type ProcessShipmentReply struct {
	Reply
	TransactionDetail       TransactionDetail
	CompletedShipmentDetail CompletedShipmentDetail
	Events                  []Event
}

func (p *ProcessShipmentReply) LabelDataAndImageType() ([]byte, string, error) {
	if label := p.CompletedShipmentDetail.CompletedPackageDetails.Label; len(label.Parts) > 0 {
		return []byte(label.Parts[0].Image), label.ImageType, nil
	}
	return nil, "", errors.New("no label")
}

func (p *ProcessShipmentReply) CommercialInvoiceDataAndImageType() ([]byte, string, error) {
	for _, document := range p.CompletedShipmentDetail.ShipmentDocuments {
		if document.Type == "COMMERCIAL_INVOICE" && len(document.Parts) > 0 {
			return []byte(document.Parts[0].Image), document.ImageType, nil
		}
	}
	return nil, "", errors.New("no commercial invoice")
}
