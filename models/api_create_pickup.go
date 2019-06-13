package models

// Pickup wraps all the Fedex API fields needed for creating a pickup
type Pickup struct {
	PickupLocation PickupLocation
	ToAddress      Address
}

type CreatePickupBody struct {
	CreatePickupRequest CreatePickupRequest `xml:"q0:CreatePickupRequest"`
}

type CreatePickupRequest struct {
	Request
	OriginDetail         OriginDetail        `xml:"q0:OriginDetail"`
	FreightPickupDetail  FreightPickupDetail `xml:"q0:FreightPickupDetail"`
	PackageCount         int                 `xml:"q0:PackageCount"`
	CarrierCode          string              `xml:"q0:CarrierCode"`
	Remarks              string              `xml:"q0:Remarks"`
	CommodityDescription string              `xml:"q0:CommodityDescription"`
}

type CreatePickupResponseEnvelope struct {
	Reply CreatePickupReply `xml:"Body>CreatePickupReply"`
}

func (c *CreatePickupResponseEnvelope) Error() error {
	err := c.Reply.Error()

	// switch e := err.(type) {
	// case nil:
	// 	log.WithFields(log.Fields{
	// 		"pickupConfirmationNumber": reply.PickupConfirmationNumber,
	// 		"delay":                    delay,
	// 		"streetLines":              pickup.PickupLocation.Address.StreetLines,
	// 	}).Info("made pickup")
	// 	break
	// case models.ReplyError:
	// 	if  {
	// 		log.WithFields(log.Fields{
	// 			"pickupConfirmationNumber": reply.PickupConfirmationNumber,
	// 			"delay":                    delay,
	// 			"streetLines":              pickup.PickupLocation.Address.StreetLines,
	// 		}).Info("made pickup")
	// 		break
	// 	}
	// 		fallthrough
	// 	}

	return err
}

// CreatePickupReply : CreatePickup reply root (`xml:"Body>CreatePickupReply"`)
type CreatePickupReply struct {
	Reply
	PickupConfirmationNumber string
	Location                 string
}
