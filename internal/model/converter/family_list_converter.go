package converter

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"time"
)

func FamilyListToResponse(familyList *entity.FamilyList) *model.FamilyListResponse {
	return &model.FamilyListResponse{
		FlID:       familyList.FlID,
		CstID:      familyList.CstID,
		FlRelation: familyList.FlRelation,
		FlName:     familyList.FlName,
		FlDob:      familyList.FlDob,
		CreatedAt:  familyList.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  familyList.UpdatedAt.Format(time.RFC3339),
	}
}

func FamilyListToEvent(request *entity.FamilyList) *model.FamilyListEvent {
	return &model.FamilyListEvent{
		CstID:      request.CstID,
		FlRelation: request.FlRelation,
		FlName:     request.FlName,
		FlDob:      request.FlDob,
	}
}
