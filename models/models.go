package models

// Structures to unmarshall the Fedex SOAP answer into

// Envelope is the soap wrapper for all requests
type Envelope struct {
	XMLName   string      `xml:"soapenv:Envelope"`
	Body      interface{} `xml:"soapenv:Body"`
	Soapenv   string      `xml:"xmlns:soapenv,attr"`
	Namespace string      `xml:"xmlns:q0,attr"`
}

// Request has just the default auth fields on all requests
type Request struct {
	WebAuthenticationDetail WebAuthenticationDetail `xml:"q0:WebAuthenticationDetail"`
	ClientDetail            ClientDetail            `xml:"q0:ClientDetail"`
	TransitionDetail        *TransactionDetail      `xml:"q0:TransitionDetail,omitempty"`
	Version                 Version                 `xml:"q0:Version"`
}

type WebAuthenticationDetail struct {
	UserCredential UserCredential `xml:"q0:UserCredential"`
}

type ClientDetail struct {
	AccountNumber string `xml:"q0:AccountNumber"`
	MeterNumber   string `xml:"q0:MeterNumber"`
}

type UserCredential struct {
	Key      string `xml:"q0:Key"`
	Password string `xml:"q0:Password"`
}

type RateRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type TrackRequest struct {
	Request
	SelectionDetails  SelectionDetails `xml:"q0:SelectionDetails"`
	ProcessingOptions string           `xml:"q0:ProcessingOptions"`
}

type CreatePickupRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type SelectionDetails struct {
	CarrierCode       string            `xml:"q0:CarrierCode"`
	PackageIdentifier PackageIdentifier `xml:"q0:PackageIdentifier"`
	// Destination           Destination
	// ShipmentAccountNumber string
}

type PackageIdentifier struct {
	Type  string `xml:"q0:Type"`
	Value string `xml:"q0:Value"`
}

type Version struct {
	ServiceID    string `xml:"q0:ServiceId"`
	Major        int    `xml:"q0:Major"`
	Intermediate int    `xml:"q0:Intermediate"`
	Minor        int    `xml:"q0:Minor"`
}

type ProcessShipmentRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type RequestedShipment struct {
	ShipTimestamp Timestamp `xml:"q0:ShipTimestamp"`
	DropoffType   string    `xml:"q0:DropoffType"`
	ServiceType   string    `xml:"q0:ServiceType"`
	PackagingType string    `xml:"q0:PackagingType"`

	Shipper   Shipper `xml:"q0:Shipper"`
	Recipient Shipper `xml:"q0:Recipient"`

	ShippingChargesPayment    Payment                    `xml:"q0:ShippingChargesPayment"`
	SmartPostDetail           *SmartPostDetail           `xml:"q0:SmartPostDetail,omitempty"`
	LabelSpecification        LabelSpecification         `xml:"q0:LabelSpecification"`
	RateRequestTypes          string                     `xml:"q0:RateRequestTypes"`
	PackageCount              int                        `xml:"q0:PackageCount"`
	RequestedPackageLineItems []RequestedPackageLineItem `xml:"q0:RequestedPackageLineItems"`
}

type SmartPostDetail struct {
	Indicia              string `xml:"q0:Indicia"`
	AncillaryEndorsement string `xml:"q0:AncillaryEndorsement"`
	HubID                string `xml:"q0:HubId"`
}
type RequestedPackageLineItem struct {
	SequenceNumber     int                 `xml:"q0:SequenceNumber"`
	GroupPackageCount  int                 `xml:"q0:GroupPackageCount,omitempty"`
	Weight             Weight              `xml:"q0:Weight"`
	Dimensions         Dimensions          `xml:"q0:Dimensions"`
	PhysicalPackaging  string              `xml:"q0:PhysicalPackaging"`
	ItemDescription    string              `xml:"q0:ItemDescription"`
	CustomerReferences []CustomerReference `xml:"q0:CustomerReferences"`
}

type CustomerReference struct {
	CustomerReferenceType string `xml:"q0:CustomerReferenceType"`
	Value                 string `xml:"q0:Value"`
}

type Weight struct {
	Units string  `xml:"q0:Units"`
	Value float64 `xml:"q0:Value"`
}

type Contact struct {
	PersonName   string `xml:"q0:PersonName"`
	CompanyName  string `xml:"q0:CompanyName"`
	PhoneNumber  string `xml:"q0:PhoneNumber"`
	EmailAddress string `xml:"q0:EMailAddress"`
}

type Dimensions struct {
	Length int    `xml:"q0:Length"`
	Width  int    `xml:"q0:Width"`
	Height int    `xml:"q0:Height"`
	Units  string `xml:"q0:Units"`
}

type Payment struct {
	PaymentType string `xml:"q0:PaymentType"`
	Payor       Payor  `xml:"q0:Payor"`
}

type Payor struct {
	ResponsibleParty ResponsibleParty `xml:"q0:ResponsibleParty"`
}

type ResponsibleParty struct {
	AccountNumber string `xml:"q0:AccountNumber"`
}

type LabelSpecification struct {
	LabelFormatType string `xml:"q0:LabelFormatType"`
	ImageType       string `xml:"q0:ImageType"`
}

type Shipper struct {
	AccountNumber string  `xml:"q0:AccountNumber"`
	Contact       Contact `xml:"q0:Contact"`
	Address       Address `xml:"q0:Address"`
}

// Reply has common stuff on all responses from FedEx API
type Reply struct {
	HighestSeverity string
	Notifications   []Notification
	Version         VersionResponse
	JobID           string `xml:"JobId"`
}

func (r Reply) Failed() bool {
	return r.HighestSeverity != "SUCCESS"
}

// TrackReply : Track reply root (`xml:"Body>TrackReply"`)
type TrackReply struct {
	Reply
	CompletedTrackDetails []CompletedTrackDetail
}

// ProcessShipReply : Process shipment reply root (`xml:"Body>ProcessShipmentReply"`)
type ProcessShipmentReply struct {
	Reply
	TransactionDetail       TransactionDetail
	CompletedShipmentDetail CompletedShipmentDetail
}

// RateReply : Process shipment reply root (`xml:"Body>RateReply"`)
type RateReply struct {
	Reply
	TransactionDetail TransactionDetail
	RateReplyDetails  []RateReplyDetail
}

// CreatePickupReply : CreatePickup reply root (`xml:"Body>CreatePickupReply"`)
type CreatePickupReply struct {
	Reply
	PickupConfirmationNumber string
	// TransactionDetail       TransactionDetail
	// CompletedShipmentDetail CompletedShipmentDetail
}

type RateReplyDetail struct {
	ServiceType                     string
	ServiceDescription              ServiceDescription
	PackagingType                   string
	DestinationAirportID            string `xml:"DestinationAirportId"`
	IneligibleForMoneyBackGuarantee bool
	SignatureOption                 string
	ActualRateType                  string
	RatedShipmentDetails            []Rating // TODO
}

type TransactionDetail struct {
	CustomerTransactionID string `xml:"q0:CustomerTransactionId,omitempty"`
}

type CompletedShipmentDetail struct {
	UsDomestic              string
	CarrierCode             string
	MasterTrackingId        TrackingID
	ServiceTypeDescription  string
	ServiceDescription      ServiceDescription
	PackagingDescription    string
	OperationalDetail       OperationalDetail
	ShipmentRating          Rating
	CompletedPackageDetails CompletedPackageDetails
}

type Part struct {
	DocumentPartSequenceNumber string
	Image                      string
}

type Label struct {
	Type                        string
	ShippingDocumentDisposition string
	ImageType                   string
	Resolution                  string
	CopiesToPrint               string
	Parts                       []Part
}

type CompletedPackageDetails struct {
	SequenceNumber string
	TrackingIds    []TrackingID
	Label          Label
}

type TrackingID struct {
	TrackingIdType string
	TrackingNumber string
}

type Name struct {
	Type     string
	Encoding string
	Value    string
}

type ServiceDescription struct {
	ServiceType      string
	Code             string
	Names            []Name
	Description      string
	AstraDescription string
}

type Surcharge struct {
	SurchargeType string
	Level         string
	Description   string
	Amount        Charge
}

type OperationalDetail struct {
	OriginLocationNumber            string
	DestinationLocationNumber       string
	TransitTime                     string
	IneligibleForMoneyBackGuarantee string
	DeliveryEligibilities           string
	ServiceCode                     string
	PackagingCode                   string
}

type RatedShipmentDetail struct {
	EffectiveNetDiscount Charge
	ShipmentRateDetail   RateDetail
	RatedPackages        []RatedPackage
}

type Rating struct {
	ActualRateType       string
	GroupNumber          string
	EffectiveNetDiscount Charge
	ShipmentRateDetails  []RateDetail
	RatedPackages        []RatedPackage
}

type Charge struct {
	Currency string
	Amount   string
}

type RatedPackage struct {
	GroupNumber          string
	EffectiveNetDiscount Charge
	PackageRateDetail    RateDetail
}

type RateDetail struct {
	RateType                         string
	RateZone                         string
	RatedWeightMethod                string
	DimDivisor                       string
	FuelSurchargePercent             string
	TotalBillingWeight               Weight
	TotalBaseCharge                  Charge
	TotalFreightDiscounts            Charge
	TotalNetFreight                  Charge
	TotalSurcharges                  Charge
	TotalNetFedExCharge              Charge
	TotalTaxes                       Charge
	TotalNetCharge                   Charge
	NetCharge                        Charge
	TotalRebates                     Charge
	TotalDutiesAndTaxes              Charge
	TotalAncillaryFeesAndTaxes       Charge
	TotalDutiesTaxesAndFees          Charge
	TotalNetChargeWithDutiesAndTaxes Charge
	Surcharges                       []Surcharge
}

type VersionResponse struct {
	ServiceID    string `xml:"ServiceId"`
	Major        int
	Intermediate int
	Minor        int
}

type CompletedTrackDetail struct {
	HighestSeverity  string
	Notifications    []Notification
	DuplicateWaybill bool
	MoreData         bool
	TrackDetails     []TrackDetail
}

type TrackDetail struct {
	TrackingNumber                         string
	TrackingNumberUniqueIdentifier         string
	Notification                           Notification
	StatusDetail                           StatusDetail
	CarrierCode                            string
	OperatingCompanyOrCarrierDescription   string
	OtherIdentifiers                       []OtherIdentifier
	Service                                Service
	PackageWeight                          Weight
	ShipmentWeight                         Weight
	Packaging                              string
	PackagingType                          string
	PackageSequenceNumber                  int
	PackageCount                           int
	SpecialHandlings                       []SpecialHandling
	ShipTimestamp                          Timestamp
	ActualDeliveryTimestamp                Timestamp
	DestinationAddress                     Address
	ActualDeliveryAddress                  Address
	DeliveryLocationType                   string
	DeliveryLocationDescription            string
	DeliveryAttempts                       int
	DeliverySignatureName                  string
	TotalUniqueAddressCountInConsolidation int
	NotificationEventsAvailable            string
	RedirectToHoldEligibility              string
	Events                                 []Event
}

type Notification struct {
	Severity         string
	Source           string
	Code             string
	Message          string
	LocalizedMessage string
}

type StatusDetail struct {
	CreationTime     string
	Code             string
	Description      string
	Location         Address
	AncillaryDetails []AncillaryDetail
}

type Address struct {
	StreetLines         []string `xml:"q0:StreetLines"`
	City                string   `xml:"q0:City"`
	StateOrProvinceCode string   `xml:"q0:StateOrProvinceCode"`
	PostalCode          string   `xml:"q0:PostalCode"`
	CountryCode         string   `xml:"q0:CountryCode"`
	// CountryName         string   `xml:"q0:CountryName"`
	Residential Bool `xml:"q0:Residential"`
}

type AncillaryDetail struct {
	Reason            string
	ReasonDescription string
}

type OtherIdentifier struct {
	PackageIdentifier Identifier
}

type Service struct {
	Type             string
	Description      string
	ShortDescription string
}

type Identifier struct {
	Type  string
	Value string
}

type SpecialHandling struct {
	Type        string
	Description string
	PaymentType string
}

type Event struct {
	Timestamp                  string
	EventType                  string
	EventDescription           string
	StatusExceptionCode        string
	StatusExceptionDescription string
	Address                    Address
	ArrivalLocation            string
}
