package models

import "errors"

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
