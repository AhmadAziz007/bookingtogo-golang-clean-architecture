package model

import "time"

type CustomerEvent struct {
	CstID         int       `json:"cstId"`
	NationalityID int       `json:"nationalityId"`
	CstName       string    `json:"cstName"`
	CstDob        time.Time `json:"cstDob"`
	CstPhoneNum   string    `json:"cstPhoneNum"`
	CstEmail      string    `json:"cstEmail"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (c *CustomerEvent) GetCstID() int {
	return c.CstID
}
