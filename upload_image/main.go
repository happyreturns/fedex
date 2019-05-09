// main uploads letterhead.png and signature.png to FedEx prod
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/happyreturns/fedex"
	"github.com/happyreturns/fedex/models"
)

func main() {

	letterhead := base64EncodedFile("letterhead.png")
	signature := base64EncodedFile("signature.png")

	credData, err := ioutil.ReadFile("../creds.json")
	if err != nil {
		panic(err)
	}

	creds := map[string]fedex.Fedex{}
	if err := json.Unmarshal(credData, &creds); err != nil {
		panic(err)
	}

	///////////////////////////////////////////////////// COMM INVOICE STUFF
	commInvoiceData, err := ioutil.ReadFile("../commercialencoded.txt")
	if err != nil {
		panic(err)
	}
	data, err := base64.StdEncoding.DecodeString(string(commInvoiceData))
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("commercial-invoice.pdf"), data, 0644)
	if err != nil {
		panic(err)
	}
	///////////////////////////////////////////////////

	prodFedex := creds["test"]

	err = prodFedex.UploadImages([]models.Image{
		{
			ID:    "IMAGE_1",
			Image: letterhead,
		},
		{
			ID:    "IMAGE_2",
			Image: signature,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("DONE")
}

//  ...
func base64EncodedFile(filename string) string {
	readBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(readBytes)
}
