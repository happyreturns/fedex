package fedex

import (
	"fmt"
	"time"
)

func processShipmentRequest(fedex Fedex, body string) string {
	return fedex.wrapSoapRequest(fmt.Sprintf(`
		<q0:ProcessShipmentRequest>
			%s
			%s
		</q0:ProcessShipmentRequest>
	`, fedex.soapCreds("ship", "23"), body), "http://fedex.com/ws/ship/v23")
}

func soapShipGround(fedex Fedex, fromLocation, toLocation Address, fromContact, toContact Contact) string {
	return processShipmentRequest(fedex, fmt.Sprintf(`
		<q0:RequestedShipment>
			<q0:ShipTimestamp>%s</q0:ShipTimestamp>
			<q0:DropoffType>REGULAR_PICKUP</q0:DropoffType>
			<q0:ServiceType>FEDEX_GROUND</q0:ServiceType>
			<q0:PackagingType>YOUR_PACKAGING</q0:PackagingType>
			<q0:Shipper>
				 <q0:AccountNumber>%s</q0:AccountNumber>
				 %s
				 %s
			</q0:Shipper>
			<q0:Recipient>
				 <q0:AccountNumber>%s</q0:AccountNumber>
				 %s
				 %s
			</q0:Recipient>
			<q0:ShippingChargesPayment>
				 <q0:PaymentType>SENDER</q0:PaymentType>
				 <q0:Payor>
						<q0:ResponsibleParty>
							 <q0:AccountNumber>%s</q0:AccountNumber>
						</q0:ResponsibleParty>
				 </q0:Payor>
			</q0:ShippingChargesPayment>
			<q0:LabelSpecification>
				 <q0:LabelFormatType>COMMON2D</q0:LabelFormatType>
				 <q0:ImageType>PNG</q0:ImageType>
			</q0:LabelSpecification>
			<q0:RateRequestTypes>LIST</q0:RateRequestTypes>
			<q0:PackageCount>1</q0:PackageCount>
			<q0:RequestedPackageLineItems>
				 <q0:SequenceNumber>1</q0:SequenceNumber>
				 <q0:Weight>
						<q0:Units>LB</q0:Units>
						<q0:Value>40</q0:Value>
				 </q0:Weight>
				 <q0:Dimensions>
						<q0:Length>5</q0:Length>
						<q0:Width>5</q0:Width>
						<q0:Height>5</q0:Height>
						<q0:Units>IN</q0:Units>
				 </q0:Dimensions>
				 <q0:PhysicalPackaging>BAG</q0:PhysicalPackaging>
				 <q0:ItemDescription>Stuff</q0:ItemDescription>
				 <q0:CustomerReferences>
						<q0:CustomerReferenceType>CUSTOMER_REFERENCE</q0:CustomerReferenceType>
						<q0:Value>NAFTA_COO</q0:Value>
				 </q0:CustomerReferences>
			</q0:RequestedPackageLineItems>
		</q0:RequestedShipment>
	`, time.Now().Format(time.RFC3339), fedex.Account, contactToString(fromContact), addressToString(fromLocation), fedex.Account, contactToString(toContact), addressToString(fromLocation), fedex.Account))
}

func streetLines(lines []string) string {
	l := ""
	for _, line := range lines {
		l += fmt.Sprintf("<q0:StreetLines>%s</q0:StreetLines>\n", line)
	}
	return l
}

func addressToString(a Address) string {
	residentialAsInt := 0
	if a.Residential {
		residentialAsInt = 1
	}

	return fmt.Sprintf(`
		<q0:Address>
			%s
			<q0:City>%s</q0:City>
			<q0:StateOrProvinceCode>%s</q0:StateOrProvinceCode>
			<q0:PostalCode>%s</q0:PostalCode>
			<q0:CountryCode>%s</q0:CountryCode>
			<q0:Residential>%d</q0:Residential>
		</q0:Address>
	 `, streetLines(a.StreetLines), a.City, a.StateOrProvinceCode, a.PostalCode, a.CountryCode, residentialAsInt)
}

func contactToString(c Contact) string {
	return fmt.Sprintf(`
		<q0:Contact>
			<q0:PersonName>%s</q0:PersonName>
			<q0:CompanyName>%s</q0:CompanyName>
			<q0:PhoneNumber>%s</q0:PhoneNumber>
			<q0:EMailAddress>%s</q0:EMailAddress>
		</q0:Contact>
	`, c.PersonName, c.CompanyName, c.PhoneNumber, c.EMailAddress)
}
