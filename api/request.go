package api

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/happyreturns/fedex/models"
)

func (a API) makeRequestAndUnmarshalResponse(url string, request *models.Envelope,
	response models.Response) error {
	// Create request body
	reqXML, err := xml.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request xml: %s", err)
	}
	fmt.Println(string(reqXML))

	// Post XML
	content, err := postXML(a.FedExURL+url, string(reqXML))
	if err != nil {
		return fmt.Errorf("post xml: %s", err)
	}
	fmt.Println(string(content))

	// Parse response
	err = xml.Unmarshal(content, response)
	if err != nil {
		return fmt.Errorf("parse xml: %s", err)
	}

	// Check if reply failed (FedEx responds with 200 even though it failed)
	err = response.Error()
	if err != nil {
		return fmt.Errorf("response error: %s", err)
	}

	return nil
}

// postXML to Fedex API and return response
func postXML(url string, xml string) ([]byte, error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all bytes: %s", err)
	}
	return content, nil
}
