package conv

import "strings"

func WeightInLbs(weight float64, unit string) float64 {
	switch strings.ToLower(unit) {
	case "lb":
		return weight
	case "oz":
		return weight / 16.0
	case "kg":
		return weight * 2.205
	case "g":
		return weight * 0.002205
	default:
		return 0.0
	}
}
