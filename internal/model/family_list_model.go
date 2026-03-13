package model

type FamilyListResponse struct {
	FlID       int    `json:"flId"`
	CstID      int    `json:"cstId"`
	FlRelation string `json:"flRelation"`
	FlName     string `json:"flName"`
	FlDob      string `json:"flDob"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type CreateFamilyListRequest struct {
	CstID      int    `json:"cstId"`
	FlRelation string `json:"flRelation" validate:"required,max=50"`
	FlName     string `json:"flName" validate:"required,max=50"`
	FlDob      string `json:"flDob" validate:"required,max=50"`
}

type UpdateFamilyListRequest struct {
	CstID      int    `json:"cstId"`
	FlID       int    `json:"-" validate:"required,lte=100"`
	FlRelation string `json:"flRelation,omitempty" validate:"max=50"`
	FlName     string `json:"flName,omitempty" validate:"max=50"`
	FlDob      string `json:"flDob,omitempty" validate:"max=50"`
}

type DeleteFamilyListRequest struct {
	CstID int `json:"-"`
	FlID  int `json:"-" validate:"required,lte=100"`
}
