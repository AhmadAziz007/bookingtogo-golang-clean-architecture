package converter

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"time"
)

func CustomerToResponse(customer *entity.Customer) *model.CustomerResponse {
	response := &model.CustomerResponse{
		CstID:         customer.CstID,
		NationalityID: customer.NationalityID,
		CstName:       customer.CstName,
		CstDob:        customer.CstDob.Format("2006-01-02"),
		CstPhoneNum:   customer.CstPhoneNum,
		CstEmail:      customer.CstEmail,
		CreatedAt:     customer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     customer.UpdatedAt.Format(time.RFC3339),
	}

	if customer.Nationality != nil {
		response.Nationality = NationalityToResponse(customer.Nationality)
	}

	if len(customer.FamilyLists) > 0 {
		for _, fl := range customer.FamilyLists {
			response.FamilyLists = append(response.FamilyLists, *FamilyListToResponse(&fl))
		}
	}

	return response
}

func CustomerToEvent(customer *entity.Customer) *model.CustomerEvent {
	return &model.CustomerEvent{
		CstID:         customer.CstID,
		NationalityID: customer.NationalityID,
		CstName:       customer.CstName,
		CstDob:        customer.CstDob,
		CstEmail:      customer.CstEmail,
		CstPhoneNum:   customer.CstPhoneNum,
		CreatedAt:     customer.CreatedAt,
		UpdatedAt:     customer.UpdatedAt,
	}
}
