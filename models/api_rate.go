package models

import (
	"errors"
	"fmt"
	"strings"
)

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

func (rr *RateReply) TotalDutiesAndTaxes() (Charge, error) {
	rateDetail, err := rr.firstRatedShipmentDetails()
	if err != nil {
		return Charge{}, fmt.Errorf("first rated shipment details: %s", err)
	}

	return rateDetail.TotalDutiesAndTaxes, nil
}

func (rr *RateReply) DutiesAndTaxesByItem() ([]Charge, error) {
	rateDetail, err := rr.firstRatedShipmentDetails()
	if err != nil {
		return nil, fmt.Errorf("first rated shipment details: %s", err)
	}

	charges := make([]Charge, len(rateDetail.DutiesAndTaxes))
	for idx, dutyAndTax := range rateDetail.DutiesAndTaxes {
		// sum up the all the taxes for this item
		// TODO make me its own function
		if len(dutyAndTax.Taxes) == 0 {
			return nil, errors.New("dutyAndTax has length 0")
		}
		charges[idx] = Charge{Currency: dutyAndTax.Taxes[0].Amount.Currency}
		for _, tax := range dutyAndTax.Taxes {
			charges[idx].Amount += tax.Amount.Amount
		}
	}

	return charges, nil
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
