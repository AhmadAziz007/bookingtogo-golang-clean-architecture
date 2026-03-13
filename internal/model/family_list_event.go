package model

import "time"

type FamilyListEvent struct {
	FlID       int       `json:"flId"`
	CstID      int       `json:"cstId"`
	FlRelation string    `json:"flRelation"`
	FlName     string    `json:"flName"`
	FlDob      string    `json:"flDob"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (f *FamilyListEvent) GetFlID() int {
	return f.FlID
}
