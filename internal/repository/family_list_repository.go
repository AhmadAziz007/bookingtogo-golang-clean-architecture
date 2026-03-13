package repository

import (
	"golang-clean-architecture/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FamilyListRepository struct {
	Repository[entity.FamilyList]
	Log *logrus.Logger
}

func NewFamilyListRepository(log *logrus.Logger) *FamilyListRepository {
	return &FamilyListRepository{
		Log: log,
	}
}

func (r *FamilyListRepository) FindByFlIDWithRelationCstID(tx *gorm.DB, customerID int) ([]entity.FamilyList, error) {
	var familyLists []entity.FamilyList
	err := tx.Where("cst_id = ?", customerID).Find(&familyLists).Error
	return familyLists, err
}

func (r *FamilyListRepository) FindByFlId(tx *gorm.DB, familyList *entity.FamilyList, flId int) error {
	return tx.Where("fl_id = ?", flId).Take(familyList).Error
}

func (r *FamilyListRepository) DeleteByCstId(tx *gorm.DB, cstId int) error {
	return tx.Where("cst_id = ?", cstId).Delete(&entity.FamilyList{}).Error
}
