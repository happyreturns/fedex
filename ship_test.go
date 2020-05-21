package fedex

import (
	"testing"

	"github.com/happyreturns/fedex/models"
	. "github.com/onsi/gomega"
)

func TestFedexShipment(t *testing.T) {
	t.Run("heavier-packages-are-more-expensive", func(t *testing.T) {
		// set up test cases
		type testCase struct {
			name        string
			fedex       Fedex
			getShipment func() *models.Shipment
		}

		// try rates with different services (smartpost, ground) and
		// destinations (international, domestic)
		testCases := []testCase{
			{
				name:        "international-ground",
				fedex:       testFedex,
				getShipment: exampleInternationalGroundShipment,
			},
			{
				name:        "domestic-ground",
				fedex:       testFedex,
				getShipment: exampleDomesticGroundShipment,
			},
			{
				name:        "domestic-smartpost",
				fedex:       laSmartPostFedex,
				getShipment: exampleDomesticSmartpostShipment,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				g := NewWithT(t)

				// Create shipment with light weight
				lightRequest := testCase.getShipment()
				lightReply, err := testCase.fedex.Ship(lightRequest)
				g.Expect(err).NotTo(HaveOccurred())

				// Create shipment with much heavier weight (10 times heavier)
				heavyRequest := testCase.getShipment()
				for idx := range heavyRequest.Commodities {
					heavyRequest.Commodities[idx].Weight.Value *= 10
				}
				heavyReply, err := testCase.fedex.Ship(heavyRequest)
				g.Expect(err).NotTo(HaveOccurred())

				// Parse reply details for the light request
				lightReplyDetails := lightReply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails
				g.Expect(lightReplyDetails).To(HaveLen(1))

				// Parse reply details for the heavy request
				heavyReplyDetails := heavyReply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails
				g.Expect(heavyReplyDetails).To(HaveLen(1))

				// Verify rate for the heavier request is much greater (at least 5 times
				// more expensive) than the light request
				lightCost := lightReplyDetails[0].TotalNetChargeWithDutiesAndTaxes
				heavyCost := heavyReplyDetails[0].TotalNetChargeWithDutiesAndTaxes
				g.Expect(lightCost.Amount).To(BeNumerically(">", 0))
				g.Expect(heavyCost.Amount).To(BeNumerically(">", 0))
				g.Expect(heavyCost.Amount).To(BeNumerically(">", lightCost.Amount*5.0))
			})
		}

	})

}

func exampleInternationalGroundShipment() *models.Shipment {
	shipment := exampleDomesticGroundShipment()
	shipment.FromAndTo.FromAddress = models.Address{
		StreetLines:         []string{"1234 Main Street", "Suite 200"},
		City:                "Winnipeg",
		StateOrProvinceCode: "MB",
		PostalCode:          "R2M4B5",
		CountryCode:         "CA",
	}
	return shipment
}

func exampleDomesticGroundShipment() *models.Shipment {
	return &models.Shipment{
		FromAndTo: models.FromAndTo{
			FromAddress: models.Address{
				StreetLines:         []string{"1517 Lincoln Blvd"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90401",
				CountryCode:         "US",
			},
			ToAddress: models.Address{
				StreetLines:         []string{"1106 Broadway"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90401",
				CountryCode:         "US",
			},
			FromContact: models.Contact{
				PersonName:   "Jenny",
				PhoneNumber:  "213 867 5309",
				EmailAddress: "jenny@jenny.com",
			},
			ToContact: models.Contact{
				CompanyName:  "Some Company",
				PhoneNumber:  "214 867 5309",
				EmailAddress: "somecompany@somecompany.com",
			},
		},
		NotificationEmail: "dev-notifications@happyreturns.com",
		References:        []string{"My ship ground reference - rothy's", "order number blah"},
		Commodities: []models.Commodity{
			{
				NumberOfPieces:       1,
				Description:          "Computer Keyboard",
				Quantity:             1,
				QuantityUnits:        "unit",
				CountryOfManufacture: "US",
				Weight:               models.Weight{Units: "LB", Value: 10.0},
				UnitPrice:            &models.Money{Currency: "USD", Amount: 25.00},
				CustomsValue:         &models.Money{Currency: "USD", Amount: 30.00},
			},
			{
				NumberOfPieces:       1,
				Description:          "Computer Monitor",
				Quantity:             1,
				QuantityUnits:        "unit",
				CountryOfManufacture: "US",
				Weight:               models.Weight{Units: "LB", Value: 5.0},
				UnitPrice:            &models.Money{Currency: "USD", Amount: 214.42},
				CustomsValue:         &models.Money{Currency: "USD", Amount: 381.12},
			},
		},
		Service: "default",
	}
}

func exampleDomesticSmartpostShipment() *models.Shipment {
	return &models.Shipment{
		FromAndTo: models.FromAndTo{
			FromAddress: models.Address{
				StreetLines:         []string{"1517 Lincoln Blvd"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90401",
				CountryCode:         "US",
			},
			ToAddress: models.Address{},
			FromContact: models.Contact{
				PersonName:   "Jenny",
				PhoneNumber:  "213 867 5309",
				EmailAddress: "jenny@jenny.com",
			},
			ToContact: models.Contact{
				CompanyName:  "Some Company",
				PhoneNumber:  "214 867 5309",
				EmailAddress: "somecompany@somecompany.com",
			},
		},
		NotificationEmail: "dev-notifications@happyreturns.com",
		References:        []string{"REF", "ORDER_NUM some string greater than 20 chars"},
		Service:           "return",
		Commodities: []models.Commodity{
			{
				NumberOfPieces:       1,
				Description:          "Computer Keyboard",
				Quantity:             1,
				QuantityUnits:        "unit",
				CountryOfManufacture: "US",
				Weight:               models.Weight{Units: "LB", Value: 1.0},
				UnitPrice:            &models.Money{Currency: "USD", Amount: 25.00},
				CustomsValue:         &models.Money{Currency: "USD", Amount: 30.00},
			},
			{
				NumberOfPieces:       1,
				Description:          "Computer Monitor",
				Quantity:             1,
				QuantityUnits:        "unit",
				CountryOfManufacture: "US",
				Weight:               models.Weight{Units: "LB", Value: 1.0},
				UnitPrice:            &models.Money{Currency: "USD", Amount: 214.42},
				CustomsValue:         &models.Money{Currency: "USD", Amount: 381.12},
			},
		},
	}
}
