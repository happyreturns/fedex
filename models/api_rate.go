package models

import (
	"errors"
	"fmt"
	"strings"
)

// Rate wraps all the Fedex API fields needed for getting a rate
type Rate struct {
	FromAndTo

	// Only used for international ground shipments
	Commodities Commodities
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

func (rr *RateReply) TotalDutiesAndTaxes() (Charge, error) {
	rateDetail, err := rr.firstRatedShipmentDetails()
	if err != nil {
		return Charge{}, fmt.Errorf("first rated shipment details: %s", err)
	}

	return rateDetail.TotalDutiesAndTaxes, nil
}

// TODO not 100% sure what we want: the Amount of dutyAndTax, or the
// TaxableValue of dutyAndTax. I think we want TaxableValue
func (rr *RateReply) DutiesAndTaxesByItem() ([]Charge, error) {
	rateDetail, err := rr.firstRatedShipmentDetails()
	if err != nil {
		return nil, fmt.Errorf("first rated shipment details: %s", err)
	}

	charges, err := rateDetail.TaxByItem()
	if err != nil {
		return nil, fmt.Errorf("tax by item: %s", err)
	}

	return charges, nil
}

func (rr *RateReply) TaxableValues() ([]Charge, error) {
	rateDetail, err := rr.firstRatedShipmentDetails()
	if err != nil {
		return nil, fmt.Errorf("first rated shipment details: %s", err)
	}

	charges := make([]Charge, len(rateDetail.DutiesAndTaxes))
	for idx, dutyAndTax := range rateDetail.DutiesAndTaxes {
		if len(dutyAndTax.Taxes) == 0 {
			return nil, errors.New("dutyAndTax has length 0")
		}
		// Assume the customs value is the first taxable value of the first tax,
		// even though there may be many taxes with different taxable values
		charges[idx] = Charge{Currency: dutyAndTax.Taxes[0].TaxableValue.Currency}
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
