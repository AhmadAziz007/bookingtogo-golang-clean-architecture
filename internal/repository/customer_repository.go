package repository

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Repository[entity.Customer]
	Log *logrus.Logger
}

func NewCustomerRepository(log *logrus.Logger) *CustomerRepository {
	return &CustomerRepository{
		Log: log,
	}
}

func (r *CustomerRepository) FindByCstIdWithRelation(db *gorm.DB, customer *entity.Customer, cstId int) error {
	return db.Preload("FamilyLists").Preload("Nationality").
		Where("cst_id = ?", cstId).
		Take(customer).Error
}

func (r *CustomerRepository) FindByCstId(db *gorm.DB, customer *entity.Customer, cstId int) error {
	return db.Where("cst_id = ?", cstId).Take(customer).Error
}

func (r *CustomerRepository) Search(db *gorm.DB, request *model.SearchCustomerRequest) ([]entity.Customer, int64, error) {
	var customers []entity.Customer
	if err := db.Scopes(r.FilterCustomer(request)).Preload("Nationality").Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Customer{}).Scopes(r.FilterCustomer(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return customers, total, nil
}

func (r *CustomerRepository) FilterCustomer(request *model.SearchCustomerRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if request.NationalityID > 0 {
			tx = tx.Where("nationality_id = ?", request.NationalityID)
		}

		if cstName := request.CstName; cstName != "" {
			cstName = "%" + cstName + "%"
			tx = tx.Where("cst_name LIKE ? ", cstName)
		}

		if cstPhoneNum := request.CstPhoneNum; cstPhoneNum != "" {
			cstPhoneNum = "%" + cstPhoneNum + "%"
			tx = tx.Where("cst_phoneNum LIKE ?", cstPhoneNum)
		}

		if cstEmail := request.CstEmail; cstEmail != "" {
			cstEmail = "%" + cstEmail + "%"
			tx = tx.Where("cst_email LIKE ?", cstEmail)
		}

		return tx
	}
}
