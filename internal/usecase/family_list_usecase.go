package usecase

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FamilyListUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	FamilyListRepository *repository.FamilyListRepository
	CustomerRepository   *repository.CustomerRepository
	Redis                *redis.Client
}

func NewFamilyListUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	familyListRepository *repository.FamilyListRepository, customerRepository *repository.CustomerRepository,
	redis *redis.Client) *FamilyListUseCase {
	return &FamilyListUseCase{
		DB:                   db,
		Log:                  logger,
		Validate:             validate,
		FamilyListRepository: familyListRepository,
		CustomerRepository:   customerRepository,
		Redis:                redis,
	}
}

func (f *FamilyListUseCase) Create(ctx context.Context, request *model.CreateFamilyListRequest) (*model.FamilyListResponse, error) {
	tx := f.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := f.Validate.Struct(request); err != nil {
		f.Log.WithError(err).Error("failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	familylist := &entity.FamilyList{
		CstID:      request.CstID,
		FlName:     request.FlName,
		FlRelation: request.FlRelation,
		FlDob:      request.FlDob,
	}

	if err := f.FamilyListRepository.Create(tx, familylist); err != nil {
		f.Log.WithError(err).Error("failed to create family list")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		f.Log.WithError(err).Error("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.FamilyListToResponse(familylist), nil
}

func (f *FamilyListUseCase) Update(ctx context.Context, request *model.UpdateFamilyListRequest) (*model.FamilyListResponse, error) {
	tx := f.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := f.Validate.Struct(request); err != nil {
		f.Log.WithError(err).Error("failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	familyList := new(entity.FamilyList)
	if err := f.FamilyListRepository.FindByFlId(tx, familyList, request.FlID); err != nil {
		f.Log.WithError(err).Error("failed to find family list")
		return nil, fiber.ErrNotFound
	}

	if request.CstID > 0 {
		familyList.CstID = request.CstID
	}

	if request.FlRelation != "" {
		familyList.FlRelation = request.FlRelation
	}
	if request.FlName != "" {
		familyList.FlName = request.FlName
	}
	if request.FlDob != "" {
		familyList.FlDob = request.FlDob
	}

	if err := f.FamilyListRepository.Update(tx, familyList); err != nil {
		f.Log.WithError(err).Error("failed to update family list")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		f.Log.WithError(err).Error("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.FamilyListToResponse(familyList), nil
}

func (f *FamilyListUseCase) Delete(ctx context.Context, request *model.DeleteFamilyListRequest) error {
	tx := f.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := f.Validate.Struct(request); err != nil {
		f.Log.WithError(err).Error("failed to validate request body")
		return fiber.ErrBadRequest
	}

	familyList := new(entity.FamilyList)
	if err := f.FamilyListRepository.FindByFlId(tx, familyList, request.FlID); err != nil {
		f.Log.WithError(err).Error("failed to find family list")
		return fiber.ErrNotFound
	}

	if err := f.FamilyListRepository.Delete(tx, familyList); err != nil {
		f.Log.WithError(err).Error("failed to delete family list")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		f.Log.WithError(err).Error("failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (f *FamilyListUseCase) GetByFlIDWithRelationCstId(ctx context.Context, customerID int) ([]model.FamilyListResponse, error) {
	db := f.DB.WithContext(ctx)

	familyLists, err := f.FamilyListRepository.FindByFlIDWithRelationCstID(db, customerID)
	if err != nil {
		f.Log.WithError(err).Error("failed to load family lists")
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.FamilyListResponse
	for _, familyList := range familyLists {
		responses = append(responses, *converter.FamilyListToResponse(&familyList))
	}

	return responses, nil
}
