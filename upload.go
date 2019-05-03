package fedex

import (
	"github.com/happyreturns/fedex/models"
)

func (f Fedex) uploadImagesRequest(images []models.Image) models.Envelope {
	// body
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/cdus/v12",
		Body: struct {
			UploadImagesRequest models.UploadImagesRequest `xml:"q0:UploadImagesRequest"`
		}{
			UploadImagesRequest: models.UploadImagesRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      f.Key,
							Password: f.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: f.Account,
						MeterNumber:   f.Meter,
					},
					Version: models.Version{
						ServiceID: "cdus",
						Major:     12,
					},
				},
				Images: images,
			},
		},
	}
}
