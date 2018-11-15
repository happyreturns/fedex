// History: Nov 20 13 tcolar Creation
package fedex

// Structures to unmarshall the Fedex SOAP answer into

type Reply struct {
	HighestSeverity string
	Notifications   []Notification
	Version         Version
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
	JobId                   string
	TransactionDetail       TransactionDetail
	CompletedShipmentDetail CompletedShipmentDetail
}

type TransactionDetail struct {
	CustomerTransactionId string
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

type Rating struct {
	ActualRateType       string
	EffectiveNetDiscount Charge
	ShipmentRateDetails  []RateDetail
}

type Charge struct {
	Currency string
	Amount   string
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
	TotalRebates                     Charge
	TotalDutiesAndTaxes              Charge
	TotalAncillaryFeesAndTaxes       Charge
	TotalDutiesTaxesAndFees          Charge
	TotalNetChargeWithDutiesAndTaxes Charge
	Surcharges                       []Surcharge
}

type Version struct {
	ServiceId    string
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
	ShipTimestamp                          string
	ActualDeliveryTimestamp                string
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
	StreetLines         []string
	City                string
	StateOrProvinceCode string
	PostalCode          string
	CountryCode         string
	CountryName         string
	Residential         bool
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

type Weight struct {
	Units string
	Value float64
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

type Contact struct {
	PersonName   string
	CompanyName  string
	PhoneNumber  string
	EMailAddress string
}
