package repository

import (
	"golang-clean-architecture/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NationalityRepository struct {
	Repository[entity.Nationality]
	Log *logrus.Logger
}

func NewNationalityRepository(log *logrus.Logger) *NationalityRepository {
	return &NationalityRepository{
		Log: log,
	}
}

func (r *NationalityRepository) FindAll(tx *gorm.DB) ([]entity.Nationality, error) {
	var nationalities []entity.Nationality
	if err := tx.Find(&nationalities).Error; err != nil {
		return nil, err
	}
	return nationalities, nil
}

func (r *NationalityRepository) FindByID(tx *gorm.DB, nationality *entity.Nationality, id int) error {
	return tx.Where("nationality_id = ?", id).First(nationality).Error
}

func (r *NationalityRepository) FindByCode(tx *gorm.DB, nationality *entity.Nationality, code string) error {
	return tx.Where("nationality_code = ?", code).First(nationality).Error
}
