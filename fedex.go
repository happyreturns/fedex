// History: Nov 20 13 tcolar Creation

// fedex provides access to (some) FedEx Soap API's and unmarshall answers into Go structures
package fedex

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/happyreturns/fedex/models"
)

const (
	// Convenience constants for standard Fedex API url's
	FedexAPIURL               = "https://ws.fedex.com:443/web-services"
	FedexAPITestURL           = "https://wsbeta.fedex.com:443/web-services"
	CarrierCodeExpress        = "FDXE"
	CarrierCodeGround         = "FDXG"
	CarrierCodeFreight        = "FXFR"
	CarrierCodeSmartPost      = "FXSP"
	CarrierCodeCustomCritical = "FXCC"
)

// Fedex : Utility to retrieve data from Fedex API
// Bypassing painful proper SOAP implementation and just crafting minimal XML messages to get the data we need.
// Fedex WSDL docs here: http://images.fedex.com/us/developer/product/WebServices/MyWebHelp/DeveloperGuide2012.pdf
type Fedex struct {
	Key, Password, Account, Meter, SmartPostHubID string
	FedexURL                                      string
}

func (f Fedex) wrapSoapRequest(body string, namespace string) string {
	return fmt.Sprintf(`
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:q0="%s">
		<soapenv:Body>
			%s
		</soapenv:Body>
		</soapenv:Envelope>
	`, namespace, body)
}

func (f Fedex) soapCreds(serviceID, majorVersion string) string {
	return fmt.Sprintf(`
		<q0:WebAuthenticationDetail>
			<q0:UserCredential>
				<q0:Key>%s</q0:Key>
				<q0:Password>%s</q0:Password>
			</q0:UserCredential>
		</q0:WebAuthenticationDetail>
		<q0:ClientDetail>
			<q0:AccountNumber>%s</q0:AccountNumber>
			<q0:MeterNumber>%s</q0:MeterNumber>
		</q0:ClientDetail>
		<q0:Version>
			<q0:ServiceId>%s</q0:ServiceId>
			<q0:Major>%s</q0:Major>
			<q0:Intermediate>0</q0:Intermediate>
			<q0:Minor>0</q0:Minor>
		</q0:Version>
	`, f.Key, f.Password, f.Account, f.Meter, serviceID, majorVersion)
}

// TrackByNumber : Returns tracking info for a specific Fedex tracking number
func (f Fedex) TrackByNumber(carrierCode string, trackingNo string) (reply models.TrackReply, err error) {
	// Create request body
	reqXML, err := xml.Marshal(f.trackByNumberSOAPRequest(carrierCode, trackingNo))
	if err != nil {
		return reply, fmt.Errorf("marshal request xml: %s", err)
	}

	// Post XML
	content, err := f.PostXML(f.FedexURL+"/trck", string(reqXML))
	if err != nil {
		return reply, fmt.Errorf("post xml: %s", err)
	}

	// Parse response
	reply, err = f.ParseTrackReply(content)
	if err != nil {
		return reply, fmt.Errorf("parse xml: %s", err)
	}
	return reply, nil
}

// TrackByShipperRef : Return tracking info for a specific shipper reference
// ShipperRef is usually an order ID or other unique identifier
// ShipperAccountNumber is the Fedex account number of the shipper
func (f Fedex) TrackByShipperRef(carrierCode string, shipperRef string,
	shipperAccountNumber string) (reply models.TrackReply, err error) {
	reqXML := soapRefTracking(f, carrierCode, shipperRef, shipperAccountNumber)
	content, err := f.PostXML(f.FedexURL+"/trck", reqXML)
	if err != nil {
		return reply, err
	}
	return f.ParseTrackReply(content)
}

// TrackByPo : Returns tracking info for a specific Purchase Order (often the OrderId)
// Note that Fedex requires the Destination Postal Code & country
//   to match when making PO queries
func (f Fedex) TrackByPo(carrierCode string, po string, postalCode string,
	countryCode string) (reply models.TrackReply, err error) {
	reqXML := soapPoTracking(f, carrierCode, po, postalCode, countryCode)
	content, err := f.PostXML(f.FedexURL+"/trck", reqXML)
	if err != nil {
		return reply, err
	}
	return f.ParseTrackReply(content)
}

func (f Fedex) ShipGround(fromAddress models.Address, toAddress models.Address, fromContact models.Contact, toContact models.Contact) (reply models.ProcessShipmentReply, err error) {
	// Create request body
	reqXML, err := xml.Marshal(f.shipGroundSOAPRequest(fromAddress, toAddress, fromContact, toContact))
	if err != nil {
		return reply, fmt.Errorf("marshal request xml: %s", err)
	}

	// Post XML
	content, err := f.PostXML(f.FedexURL+"/ship/v23", string(reqXML))
	if err != nil {
		return reply, fmt.Errorf("post xml: %s", err)
	}

	// Parse response
	reply, err = f.ParseProcessShipmentReply(content)
	if err != nil {
		return reply, fmt.Errorf("parse xml: %s", err)
	}
	return reply, nil
}

func (f Fedex) ShipSmartPost(fromAddress models.Address, toAddress models.Address, fromContact models.Contact, toContact models.Contact) (reply models.ProcessShipmentReply, err error) {
	// Create request body
	reqXML, err := xml.Marshal(f.shipSmartPostSOAPRequest(fromAddress, toAddress, fromContact, toContact))
	if err != nil {
		return reply, fmt.Errorf("marshal request xml: %s", err)
	}

	// Post XML
	content, err := f.PostXML(f.FedexURL+"/ship/v23", string(reqXML))
	if err != nil {
		return reply, fmt.Errorf("post xml: %s", err)
	}

	// Parse response
	reply, err = f.ParseProcessShipmentReply(content)
	if err != nil {
		return reply, fmt.Errorf("parse xml: %s", err)
	}
	return reply, nil
}

func (f Fedex) Rate(fromAddress models.Address, toAddress models.Address, fromContact models.Contact, toContact models.Contact) (reply models.RateReply, err error) {

	// Create request body
	reqXML, err := xml.Marshal(f.rateSOAPRequest(fromAddress, toAddress, fromContact, toContact))
	if err != nil {
		return reply, fmt.Errorf("marshal request xml: %s", err)
	}

	// Post XML
	content, err := f.PostXML(f.FedexURL+"/rate/v24", string(reqXML))
	if err != nil {
		return reply, fmt.Errorf("post xml: %s", err)
	}

	// Parse response
	reply, err = f.ParseRateReply(content)
	if err != nil {
		return reply, fmt.Errorf("parse xml: %s", err)
	}
	return reply, nil
}

// TODO
// func (f Fedex) Pickup() (reply models.CreatePickupReply, err error) {
// 	return nil, nil
// // Create request body
// reqXML, err := xml.Marshal(rateSOAPRequest(f, fromAddress, toAddress, fromContact, toContact))
// if err != nil {
// 	return reply, fmt.Errorf("marshal request xml: %s", err)
// }
//
// // Post XML
// content, err := f.PostXML(f.FedexURL+"/rate/v24", string(reqXML))
// if err != nil {
// 	return reply, fmt.Errorf("post xml: %s", err)
// }
// fmt.Println(string(content))
//
// // Parse response
// reply, err = f.ParseRateReply(content)
// if err != nil {
// 	return reply, fmt.Errorf("parse xml: %s", err)
// }
// return reply, nil
// }

// Unmarshal XML SOAP response into a TrackReply
func (f Fedex) ParseTrackReply(xmlResp []byte) (reply models.TrackReply, err error) {
	data := struct {
		Reply models.TrackReply `xml:"Body>TrackReply"`
	}{}
	err = xml.Unmarshal(xmlResp, &data)
	return data.Reply, err
}

func (f Fedex) ParseProcessShipmentReply(xmlResp []byte) (reply models.ProcessShipmentReply, err error) {
	data := struct {
		Reply models.ProcessShipmentReply `xml:"Body>ProcessShipmentReply"`
	}{}
	err = xml.Unmarshal(xmlResp, &data)
	return data.Reply, err
}

func (f Fedex) ParseRateReply(xmlResp []byte) (reply models.RateReply, err error) {
	data := struct {
		Reply models.RateReply `xml:"Body>RateReply"`
	}{}
	err = xml.Unmarshal(xmlResp, &data)
	return data.Reply, err
}

// Post Xml to Fedex API and return response
func (f Fedex) PostXML(url string, xml string) (content []byte, err error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
