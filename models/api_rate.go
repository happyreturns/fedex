package models

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Rate wraps all the Fedex API fields needed for getting a rate
type Rate struct {
	FromAndTo

	Service     string
	Commodities Commodities
}

func (r *Rate) ServiceType() string {
	return ServiceType(r.FromAndTo, r.Service)
}

func (r *Rate) SpecialServicesRequested() *SpecialServicesRequested {
	var (
		specialServiceTypes []string

		etdDetail               *EtdDetail
		eventNotificationDetail *EventNotificationDetail
		returnShipmentDetail    *ReturnShipmentDetail
	)

	if r.ServiceType() == "SMART_POST" {
		specialServiceTypes = append(specialServiceTypes, "RETURN_SHIPMENT")
		returnShipmentDetail = &ReturnShipmentDetail{
			ReturnType: "PRINT_RETURN_LABEL",
		}
	}

	if r.IsInternational() {
		specialServiceTypes = append(specialServiceTypes, "ELECTRONIC_TRADE_DOCUMENTS")
		etdDetail = &EtdDetail{
			RequestedDocumentCopies: "COMMERCIAL_INVOICE",
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
	commoditiesWeight := r.Commodities.Weight()
	if !commoditiesWeight.IsZero() {
		// Assume the weight must be between than 13 and 150 lbs.
		// If the weight is less than 13 lbs, assume a weight of 13 lbs, which is
		// heavy enough that the destination will matter when choosing between two
		// fedex ground rates
		commoditiesWeight.Value = math.Min(commoditiesWeight.Value, 150.0)
		commoditiesWeight.Value = math.Max(commoditiesWeight.Value, 13.0)
		return commoditiesWeight
	}

	return Weight{Units: "LB", Value: 13}
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
	// TODO We find the first RatedshipmentDetail for figuring out the cost of
	// the total shipment, taxes, etc. There can be other RatedshipmentDetails (
	// From what I can tell online, the ones RateType equal to
	// `PAYOR_ACCOUNT_PACKAGE` or `PAYOR_ACCOUNT_SHIPMENT` are the ones we should
	// pay attention.
	for _, rateReplyDetail := range rr.RateReplyDetails {
		for _, ratedShipmentDetail := range rateReplyDetail.RatedShipmentDetails {
			if strings.HasPrefix(ratedShipmentDetail.ShipmentRateDetail.RateType, "PAYOR_") {
				return ratedShipmentDetail.ShipmentRateDetail, nil
			}
		}
	}

	return RateDetail{}, errors.New("no RatedShipmentDetails with PAYOR_ prefix found")
}
