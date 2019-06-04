package models

type PickupAvailabilityBody struct {
	PickupAvailabilityRequest PickupAvailabilityRequest `xml:"q0:PickupAvailabilityRequest"`
}

type PickupAvailabilityRequest struct {
	Request
	PickupAddress        Address  `xml:"q0:PickupAddress"`
	PickupRequestType    []string `xml:"q0:PickupRequestType"`
	NumberOfBusinessDays int      `xml:"q0:NumberOfBusinessDays"`
	Carriers             []string `xml:"q0:Carriers"`
}

type PickupAvailabilityResponseEnvelope struct {
	Reply PickupAvailabilityReply `xml:"Body>PickupAvailabilityReply"`
}

func (c *PickupAvailabilityResponseEnvelope) Error() error {
	return c.Reply.Error()
}

// PickupAvailabilityReply : PickupAvailability reply root (`xml:"Body>PickupAvailabilityReply"`)
type PickupAvailabilityReply struct {
	Reply
	Options []PickupScheduleOption
}
