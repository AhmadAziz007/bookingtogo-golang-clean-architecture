package converter

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
)

func NationalityToResponse(nationality *entity.Nationality) *model.NationalityResponse {
	return &model.NationalityResponse{
		NationalityID:   nationality.NationalityID,
		NationalityName: nationality.NationalityName,
		NationalityCode: nationality.NationalityCode,
	}
}

func NationalitiesToResponses(nationalities []entity.Nationality) []model.NationalityResponse {
	var responses []model.NationalityResponse
	for _, nationality := range nationalities {
		responses = append(responses, *NationalityToResponse(&nationality))
	}
	return responses
}
