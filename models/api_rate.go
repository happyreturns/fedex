package models

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/happyreturns/fedex/conv"
)

// Rate wraps all the Fedex API fields needed for getting a rate
type Rate struct {
	FromAndTo

	Service     string
	Commodities Commodities
}

func (r *Rate) ServiceType() string {
	serviceType := ServiceType(r.FromAndTo, r.Service, r.Service)
	if serviceType == ServiceTypeSmartPost {
		// This is necessary. We can't get back smartpost rates. So using ground
		// instead here.
		// Per Page 239 of the Dev Guide: "Estimated shipping rates are not
		// available for SmartPost Returns"
		serviceType = ServiceTypeFedexGround
	}
	return serviceType
}

func (r *Rate) SpecialServicesRequested() *SpecialServicesRequested {
	var (
		specialServiceTypes []string

		etdDetail               *EtdDetail
		eventNotificationDetail *EventNotificationDetail
		returnShipmentDetail    *ReturnShipmentDetail
	)

	if r.ServiceType() == ServiceTypeSmartPost {
		specialServiceTypes = append(specialServiceTypes, SpecialServiceTypeReturnShipment)
		returnShipmentDetail = &ReturnShipmentDetail{
			ReturnType: ReturnTypePrintReturnLabel,
		}
	}

	if r.IsInternational() {
		specialServiceTypes = append(specialServiceTypes, SpecialServiceTypeElectronicTradeDocuments)
		etdDetail = &EtdDetail{
			RequestedDocumentCopies: DocumentTypeCommercialInvoice,
		}
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

func (r *Rate) Weight() Weight {
	commoditiesSumWeight := r.Commodities.Weight()

	if !commoditiesSumWeight.IsZero() {
		commoditiesSumWeight.Value = math.Min(commoditiesSumWeight.Value, MaximumWeightInLbs)

		return commoditiesSumWeight
	}

	return Weight{Units: WeightUnitsLB, Value: conv.WeightInLbs(SafeGuardForZeroWeightOz, "oz")}
}

type RateBody struct {
	RateRequest RateRequest `xml:"q0:RateRequest"`
}

type RateRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type RateResponseEnvelope struct {
	Reply RateReply `xml:"Body>RateReply"`
}

func (r *RateResponseEnvelope) Error() error {
	return r.Reply.Error()
}

// RateReply : Process shipment reply root (`xml:"Body>RateReply"`)
type RateReply struct {
	Reply
	TransactionDetail TransactionDetail
	RateReplyDetails  []RateReplyDetail
}

// TotalCost returns the sum of any charges in the reply
func (rr *RateReply) TotalCost() (Charge, error) {
	rateDetail, err := rr.firstRatedShipmentDetails()
	if err != nil {
		return Charge{}, fmt.Errorf("first rated shipment details: %s", err)
	}

	return rateDetail.TotalNetChargeWithDutiesAndTaxes, nil
}

func (rr *RateReply) firstRatedShipmentDetails() (RateDetail, error) {

	// Find the rated shipment detail of type "PREFERRED_ACCOUNT_PACKAGE"
	for _, rateReplyDetail := range rr.RateReplyDetails {
		for _, ratedShipmentDetail := range rateReplyDetail.RatedShipmentDetails {
			if ratedShipmentDetail.ShipmentRateDetail.RateType == RateTypePreferredAccountPackage {
				return ratedShipmentDetail.ShipmentRateDetail, nil
			}
		}
	}

	// We prefer the rated shipment detail of type "PREFERRED_ACCOUNT_PACKAGE",
	// but if that isn't found, return the rated shipment detail with RateType
	// equal to `PAYOR_ACCOUNT_PACKAGE` or `PAYOR_ACCOUNT_SHIPMENT`
	for _, rateReplyDetail := range rr.RateReplyDetails {
		for _, ratedShipmentDetail := range rateReplyDetail.RatedShipmentDetails {
			if strings.HasPrefix(ratedShipmentDetail.ShipmentRateDetail.RateType, "PAYOR_") {
				return ratedShipmentDetail.ShipmentRateDetail, nil
			}
		}
	}

	return RateDetail{}, errors.New("no RatedShipmentDetails found")
}
