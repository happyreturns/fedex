package fedex

import (
	"testing"

	"github.com/happyreturns/fedex/models"
)

var f Fedex = Fedex{
	Key:      "IfRnoRdbpEvBPbEn",
	Password: "A4dpGK2dPW4P2sSba9suwOCpo",
	Account:  "510087780",
	Meter:    "119090332",
	FedexURL: FedexAPITestURL,
}

func TestTrack(t *testing.T) {
	reply, err := f.TrackByNumber(CarrierCodeExpress, "123456789012")
	if err != nil {
		t.Fatal(err)
	}
	if reply.Failed() {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "trck" ||
		reply.Notifications[0].Code != "0" ||
		reply.Notifications[0].Message != "Request was successfully processed." ||
		reply.Notifications[0].LocalizedMessage != "Request was successfully processed." ||
		reply.Version.ServiceID != "trck" ||
		reply.Version.Major != 16 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		len(reply.CompletedTrackDetails) != 1 ||
		!reply.CompletedTrackDetails[0].DuplicateWaybill ||
		reply.CompletedTrackDetails[0].MoreData ||
		len(reply.CompletedTrackDetails[0].TrackDetails) != 16 ||
		reply.CompletedTrackDetails[0].TrackDetails[0].OperatingCompanyOrCarrierDescription != "FedEx Express" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumber != "123456789012" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumberUniqueIdentifier != "2458115001~123456789012~FX" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].CarrierCode != "FDXE" {
		t.Fatal("output not correct")
	}
}

func TestRate(t *testing.T) {
	reply, err := f.Rate(models.Address{
		StreetLines:         []string{"1511 15th Street", "#205"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90404",
		CountryCode:         "US",
	}, models.Address{
		StreetLines:         []string{"1106 Broadway"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90401",
		CountryCode:         "US",
	}, models.Contact{
		PersonName:   "Joachim Valdez",
		PhoneNumber:  "213 867 5309",
		EmailAddress: "joachim@happyreturns.com",
	}, models.Contact{
		CompanyName:  "Happy Returns",
		PhoneNumber:  "424 325 9510",
		EmailAddress: "joachim@happyreturns.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	if reply.Failed() {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "crs" ||
		reply.Notifications[0].Code != "0" ||
		reply.Version.ServiceID != "crs" ||
		reply.Version.Major != 24 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		len(reply.RateReplyDetails) != 1 ||
		reply.RateReplyDetails[0].ServiceType != "FEDEX_GROUND" ||
		reply.RateReplyDetails[0].ServiceDescription.ServiceType != "FEDEX_GROUND" ||
		reply.RateReplyDetails[0].ServiceDescription.Code != "92" ||
		reply.RateReplyDetails[0].ServiceDescription.AstraDescription != "FXG" ||
		reply.RateReplyDetails[0].PackagingType != "YOUR_PACKAGING" ||
		reply.RateReplyDetails[0].DestinationAirportID != "YOUR_PACKAGING" ||
		reply.RateReplyDetails[0].IneligibleForMoneyBackGuarantee ||
		reply.RateReplyDetails[0].SignatureOption != "SERVICE_DEFAULT" ||
		reply.RateReplyDetails[0].ActualRateType != "PAYOR_ACCOUNT_PACKAGE" ||
		len(reply.RateReplyDetails[0].RatedShipmentDetails) != 2 ||
		reply.RateReplyDetails[0].RatedShipmentDetails[0].EffectiveNetDiscount.Amount != "USD" ||
		len(reply.RateReplyDetails[0].RatedShipmentDetails[0].RatedPackages) != 1 ||
		reply.RateReplyDetails[0].RatedShipmentDetails[0].RatedPackages[0].PackageRateDetail.NetCharge.Amount != "0.0" ||

		len(reply.RateReplyDetails[1].RatedShipmentDetails) != 2 ||
		reply.RateReplyDetails[1].RatedShipmentDetails[0].EffectiveNetDiscount.Amount != "USD" ||
		len(reply.RateReplyDetails[1].RatedShipmentDetails[0].RatedPackages) != 1 ||
		reply.RateReplyDetails[1].RatedShipmentDetails[0].RatedPackages[0].PackageRateDetail.NetCharge.Amount != "0.0" ||

		reply.RateReplyDetails[1].RatedShipmentDetails[0].EffectiveNetDiscount.Amount != "USD" ||
		len(reply.RateReplyDetails[1].RatedShipmentDetails[0].RatedPackages) != 1 ||
		reply.RateReplyDetails[1].RatedShipmentDetails[0].RatedPackages[0].PackageRateDetail.NetCharge.Amount != "0.0" ||
		reply.RateReplyDetails[1].RatedShipmentDetails[0].EffectiveNetDiscount.Amount != "USD" {
		t.Fatal("output not correct")
	}
}

func TestShipGround(t *testing.T) {
	reply, err := f.ShipGround(models.Address{
		StreetLines:         []string{"1511 15th Street", "#205"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90404",
		CountryCode:         "US",
	}, models.Address{
		StreetLines:         []string{"1106 Broadway"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90401",
		CountryCode:         "US",
	}, models.Contact{
		PersonName:   "Joachim Valdez",
		PhoneNumber:  "213 867 5309",
		EmailAddress: "joachim@happyreturns.com",
	}, models.Contact{
		CompanyName:  "Happy Returns",
		PhoneNumber:  "424 325 9510",
		EmailAddress: "joachim@happyreturns.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	if reply.Failed() {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "ship" ||
		reply.Notifications[0].Code != "0000" ||
		reply.Notifications[0].Message != "Success" ||
		reply.Notifications[0].LocalizedMessage != "Success" ||
		reply.Version.ServiceID != "ship" ||
		reply.Version.Major != 23 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		reply.JobID == "" ||
		reply.CompletedShipmentDetail.UsDomestic != "true" ||
		reply.CompletedShipmentDetail.CarrierCode != "FDXG" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingIdType != "FEDEX" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingNumber == "" ||
		reply.CompletedShipmentDetail.ServiceTypeDescription != "FXG" ||
		reply.CompletedShipmentDetail.ServiceDescription.ServiceType != "FEDEX_GROUND" ||
		reply.CompletedShipmentDetail.ServiceDescription.Code != "92" ||
		// skip ServiceDescription.Names
		reply.CompletedShipmentDetail.PackagingDescription != "YOUR_PACKAGING" ||
		reply.CompletedShipmentDetail.OperationalDetail.OriginLocationNumber != "901" ||
		reply.CompletedShipmentDetail.OperationalDetail.DestinationLocationNumber != "901" ||
		reply.CompletedShipmentDetail.OperationalDetail.TransitTime != "ONE_DAY" ||
		reply.CompletedShipmentDetail.OperationalDetail.IneligibleForMoneyBackGuarantee != "false" ||
		reply.CompletedShipmentDetail.OperationalDetail.DeliveryEligibilities != "SATURDAY_DELIVERY" ||
		reply.CompletedShipmentDetail.OperationalDetail.ServiceCode != "92" ||
		reply.CompletedShipmentDetail.OperationalDetail.PackagingCode != "01" ||
		reply.CompletedShipmentDetail.ShipmentRating.ActualRateType != "PAYOR_ACCOUNT_PACKAGE" ||
		reply.CompletedShipmentDetail.ShipmentRating.EffectiveNetDiscount.Currency != "USD" ||
		reply.CompletedShipmentDetail.ShipmentRating.EffectiveNetDiscount.Amount != "0.0" ||
		len(reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails) != 2 ||
		// skip most ShipmentRateDetails fields
		reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails[0].RateType != "PAYOR_ACCOUNT_PACKAGE" ||
		reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails[1].RateType != "PAYOR_LIST_PACKAGE" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds) != 1 ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds[0].TrackingIdType != "FEDEX" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Type != "OUTBOUND_LABEL" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.ImageType != "PNG" ||

		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts) != 1 ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image == "" {
		t.Fatal("output not correct")
	}
}
