package models

import (
	shippingModels "github.com/happyreturns/shipping/models"
)

var (
	ShippingStatusMap = map[string]string{
		"AA":  shippingModels.StatusInTransit,
		"AC":  shippingModels.StatusInTransit,
		"AF":  shippingModels.StatusInTransit,
		"AP":  shippingModels.StatusInTransit,
		"AR":  shippingModels.StatusInTransit,
		"AX":  shippingModels.StatusInTransit,
		"AXA": shippingModels.StatusInTransit,
		"CC":  shippingModels.StatusInTransit,
		"CP":  shippingModels.StatusInTransit,
		"DP":  shippingModels.StatusInTransit,
		"DS":  shippingModels.StatusInTransit,
		"ED":  shippingModels.StatusInTransit,
		"EO":  shippingModels.StatusInTransit,
		"FD":  shippingModels.StatusInTransit,
		"HL":  shippingModels.StatusInTransit,
		"IT":  shippingModels.StatusInTransit,
		"LO":  shippingModels.StatusInTransit,
		"OF":  shippingModels.StatusInTransit,
		"PF":  shippingModels.StatusInTransit,
		"PL":  shippingModels.StatusInTransit,
		"SE":  shippingModels.StatusInTransit,
		"SF":  shippingModels.StatusInTransit,
		"SP":  shippingModels.StatusInTransit,
		"TR":  shippingModels.StatusInTransit,
		"CA":  shippingModels.StatusException,
		"DD":  shippingModels.StatusException,
		"DE":  shippingModels.StatusException,
		"CD":  shippingModels.StatusException,
		"AD":  shippingModels.StatusOutForDelivery,
		"OD":  shippingModels.StatusOutForDelivery,
		"DL":  shippingModels.StatusDelivered,
	}
)
