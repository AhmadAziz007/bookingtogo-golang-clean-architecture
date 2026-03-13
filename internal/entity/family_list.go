package entity

import (
	"time"
)

type FamilyList struct {
	FlID       int       `gorm:"column:fl_id;primaryKey;autoIncrement"`
	CstID      int       `gorm:"column:cst_id;not null"`
	FlRelation string    `gorm:"column:fl_relation;type:varchar(50);not null"`
	FlName     string    `gorm:"column:fl_name;type:varchar(50);not null"`
	FlDob      string    `gorm:"column:fl_dob;type:varchar(50);not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Customer *Customer `gorm:"foreignKey:CstID;references:cst_id"`
}

func (f *FamilyList) TableName() string {
	return "family_list"
}
