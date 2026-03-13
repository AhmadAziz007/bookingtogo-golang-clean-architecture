package model

type NationalityResponse struct {
	NationalityID   int    `json:"nationalityId"`
	NationalityName string `json:"nationalityName"`
	NationalityCode string `json:"nationalityCode"`
}
