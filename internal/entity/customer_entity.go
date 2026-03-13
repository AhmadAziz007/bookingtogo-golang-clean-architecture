package entity

import (
	"time"
)

type Customer struct {
	CstID         int       `gorm:"column:cst_id;primaryKey;autoIncrement"`
	NationalityID int       `gorm:"column:nationality_id;not null"`
	CstName       string    `gorm:"column:cst_name;type:char(50);not null"`
	CstDob        time.Time `gorm:"column:cst_dob;type:date;not null"`
	CstPhoneNum   string    `gorm:"column:cst_phoneNum;type:varchar(20);not null"`
	CstEmail      string    `gorm:"column:cst_email;type:varchar(50);not null"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Nationality *Nationality `gorm:"foreignKey:NationalityID;references:nationality_id"`
	FamilyLists []FamilyList `gorm:"foreignKey:CstID;references:cst_id"`
}

func (c *Customer) TableName() string {
	return "customers"
}
