package entity

type Nationality struct {
	NationalityID   int    `gorm:"column:nationality_id;primaryKey"`
	NationalityName string `gorm:"column:nationality_name;type:varchar(50);not null"`
	NationalityCode string `gorm:"column:nationality_code;type:char(2);not null"`
}

func (n *Nationality) TableName() string {
	return "nationality"
}
