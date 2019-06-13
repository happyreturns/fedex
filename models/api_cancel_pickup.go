package models

type CancelPickupBody struct {
	CancelPickupRequest CancelPickupRequest `xml:"q0:CancelPickupRequest"`
}

type CancelPickupRequest struct {
	Request
	CarrierCode              string   `xml:"q0:CarrierCode"`
	PickupConfirmationNumber string   `xml:"q0:PickupConfirmationNumber"`
	Remarks                  []string `xml:"q0:Remarks"`
}

type CancelPickupResponseEnvelope struct {
	Reply CancelPickupReply `xml:"Body>CancelPickupReply"`
}

func (c *CancelPickupResponseEnvelope) Error() error {
	return c.Reply.Error()
}

// CancelPickupReply : CancelPickup reply root (`xml:"Body>CancelPickupReply"`)
type CancelPickupReply struct {
	Reply
	Message string
}
