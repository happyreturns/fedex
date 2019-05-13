package api

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

func (a API) UploadImages(images []models.Image) error {
	request := a.uploadImagesRequest(images)

	response := &models.UploadImagesResponseEnvelope{}
	if err := a.makeRequestAndUnmarshalResponse("/uploaddocument/v11", request, response); err != nil {
		return fmt.Errorf("make upload images request and unmarshal: %s", err)
	}

	return nil
}

func (a API) uploadImagesRequest(images []models.Image) models.Envelope {
	// body
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/uploaddocument/v11",
		Body: struct {
			UploadImagesRequest models.UploadImagesRequest `xml:"q0:UploadImagesRequest"`
		}{
			UploadImagesRequest: models.UploadImagesRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      a.Key,
							Password: a.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: a.Account,
						MeterNumber:   a.Meter,
					},
					Version: models.Version{
						ServiceID: "cdus",
						Major:     11,
					},
				},
				Images: images,
			},
		},
	}
}
