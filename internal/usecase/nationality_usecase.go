package usecase

import (
	"context"

	"errors"
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

type NationalityUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	NationalityRepository *repository.NationalityRepository
	Redis                 *redis.Client
}

func NewNationalityUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	nationalityRepository *repository.NationalityRepository, redis *redis.Client) *NationalityUseCase {
	return &NationalityUseCase{
		DB:                    db,
		Log:                   logger,
		Validate:              validate,
		NationalityRepository: nationalityRepository,
		Redis:                 redis,
	}
}

func (n *NationalityUseCase) List(ctx context.Context) ([]model.NationalityResponse, error) {
	db := n.DB.WithContext(ctx)

	nationalities, err := n.NationalityRepository.FindAll(db)
	if err != nil {
		n.Log.WithError(err).Error("failed to list nationalities")
		return nil, fiber.ErrInternalServerError
	}

	return converter.NationalitiesToResponses(nationalities), nil
}

func (n *NationalityUseCase) GetByID(ctx context.Context, id int) (*model.NationalityResponse, error) {
	db := n.DB.WithContext(ctx)

	nationality := new(entity.Nationality)
	if err := n.NationalityRepository.FindByID(db, nationality, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		n.Log.WithError(err).Error("failed to get nationality by id")
		return nil, fiber.ErrInternalServerError
	}

	return converter.NationalityToResponse(nationality), nil
}

func (n *NationalityUseCase) GetByCode(ctx context.Context, code string) (*model.NationalityResponse, error) {
	db := n.DB.WithContext(ctx)

	nationality := new(entity.Nationality)
	if err := n.NationalityRepository.FindByCode(db, nationality, code); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		n.Log.WithError(err).Error("failed to get nationality by code")
		return nil, fiber.ErrInternalServerError
	}

	return converter.NationalityToResponse(nationality), nil
}
