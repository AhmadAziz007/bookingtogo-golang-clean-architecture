package usecase

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	CustomerRepository   *repository.CustomerRepository
	FamilyListRepository *repository.FamilyListRepository
	Redis                *redis.Client
}

func NewCustomerUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	customerRepository *repository.CustomerRepository, familyListRepository *repository.FamilyListRepository,
	redis *redis.Client) *CustomerUseCase {
	return &CustomerUseCase{
		DB:                   db,
		Log:                  logger,
		Validate:             validate,
		CustomerRepository:   customerRepository,
		FamilyListRepository: familyListRepository,
		Redis:                redis,
	}
}

func (c *CustomerUseCase) Create(ctx context.Context, request *model.CreateCustomerRequest) (*model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	dob, err := time.Parse("2006-01-02", request.CstDob)
	if err != nil {
		c.Log.WithError(err).Error("failed to parse date of birth")
		return nil, fiber.ErrBadRequest
	}

	customer := &entity.Customer{
		NationalityID: request.NationalityID,
		CstName:       request.CstName,
		CstDob:        dob,
		CstEmail:      request.CstEmail,
		CstPhoneNum:   request.CstPhoneNum,
	}

	if err := c.CustomerRepository.Create(tx, customer); err != nil {
		c.Log.WithError(err).Error("failed to create customer")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	// After commit: run non-critical side-effects concurrently with errgroup.
	// Examples: update cache and publish event. These should not affect the
	// main request success; we log on error.
	// g, _ := errgroup.WithContext(ctx)

	// // cache update
	// g.Go(func() error {
	// 	c.Log.WithFields(logrus.Fields{
	// 		"cst_id": customer.CstID,
	// 		"action": "cache_update_skipped",
	// 	}).Info("background cache update skipped (no external cache configured)")
	// 	return nil
	// })

	// // publish event (placeholder)
	// g.Go(func() error {
	// 	// Replace with actual publisher (kafka, rabbitmq, etc.)
	// 	// here we just simulate by logging — return nil to indicate non-fatal.
	// 	c.Log.WithFields(logrus.Fields{
	// 		"cst_id": customer.CstID,
	// 		"action": "customer_created",
	// 	}).Info("publish customer created event (placeholder)")
	// 	return nil
	// })

	// if err := g.Wait(); err != nil {
	// 	c.Log.WithError(err).Warn("one or more background tasks failed")
	// }

	// background tasks using WaitGroup + error channel (placeholders)
	// var wg sync.WaitGroup
	// errCh := make(chan error, 2)

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	c.Log.WithFields(logrus.Fields{
	// 		"cst_id": customer.CstID,
	// 		"action": "cache_update_skipped",
	// 	}).Info("background cache update skipped (no external cache configured)")
	// 	errCh <- nil
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	c.Log.WithFields(logrus.Fields{
	// 		"cst_id": customer.CstID,
	// 		"action": "customer_created",
	// 	}).Info("publish customer created event (placeholder)")
	// 	errCh <- nil
	// }()

	// wg.Wait()
	// close(errCh)
	// for e := range errCh {
	// 	if e != nil {
	// 		c.Log.WithError(e).Warn("background task error")
	// 	}
	// }

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) Update(ctx context.Context, request *model.UpdateCustomerRequest) (*model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByCstId(tx, customer, request.CstID); err != nil {
		c.Log.WithError(err).Error("failed to get customer")
		return nil, fiber.ErrNotFound
	}

	if request.CstName != "" {
		customer.CstName = request.CstName
	}

	if request.NationalityID > 0 {
		customer.NationalityID = request.NationalityID
	}

	if request.CstDob != "" {
		dob, err := time.Parse("2006-01-02", request.CstDob)
		if err != nil {
			c.Log.WithError(err).Error("failed to parse date of birth")
			return nil, fiber.ErrBadRequest
		}
		customer.CstDob = dob
	}

	if request.CstEmail != "" {
		customer.CstEmail = request.CstEmail
	}

	if request.CstPhoneNum != "" {
		customer.CstPhoneNum = request.CstPhoneNum
	}

	if err := c.CustomerRepository.Update(tx, customer); err != nil {
		c.Log.WithError(err).Error("failed to update customer")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	// After commit: run non-critical side-effects concurrently with errgroup.
	// For now these are placeholders (no Redis/Kafka integration).
	// g, _ := errgroup.WithContext(ctx)

	// g.Go(func() error {
	// 	c.Log.WithFields(logrus.Fields{
	// 		"cst_id": customer.CstID,
	// 		"action": "cache_update_skipped",
	// 	}).Info("background cache update skipped (no external cache configured)")
	// 	return nil
	// })

	// g.Go(func() error {
	// 	c.Log.WithFields(logrus.Fields{
	// 		"cst_id": customer.CstID,
	// 		"action": "customer_updated",
	// 	}).Info("publish customer updated event (placeholder)")
	// 	return nil
	// })

	// if err := g.Wait(); err != nil {
	// 	c.Log.WithError(err).Warn("one or more background tasks failed")
	// }

	// background tasks using WaitGroup + error channel (placeholders)
	// var wg sync.WaitGroup
	// errCh := make(chan error, 2)

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	c.Log.WithFields(logrus.Fields{
	// 		"cst_id": customer.CstID,
	// 		"action": "cache_update_skipped",
	// 	}).Info("background cache update skipped (no external cache configured)")
	// 	errCh <- nil
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	c.Log.WithFields(logrus.Fields{
	// 		"cst_id": customer.CstID,
	// 		"action": "customer_updated",
	// 	}).Info("publish customer updated event (placeholder)")
	// 	errCh <- nil
	// }()

	// wg.Wait()
	// close(errCh)
	// for e := range errCh {
	// 	if e != nil {
	// 		c.Log.WithError(e).Warn("background task error")
	// 	}
	// }

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) Get(ctx context.Context, request *model.GetCustomerRequest) (*model.CustomerResponse, error) {
	db := c.DB.WithContext(ctx)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByCstIdWithRelation(db, customer, request.CstID); err != nil {
		c.Log.WithError(err).Error("failed to get customer")
		return nil, fiber.ErrNotFound
	}

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) Delete(ctx context.Context, request *model.DeleteCustomerRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return fiber.ErrBadRequest
	}

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByCstIdWithRelation(tx, customer, request.CstID); err != nil {
		c.Log.WithError(err).Error("failed to get customer")
		return fiber.ErrNotFound
	}

	if len(customer.FamilyLists) > 0 {
		if err := c.FamilyListRepository.DeleteByCstId(tx, request.CstID); err != nil {
			c.Log.WithError(err).Error("error deleting family lists")
			return fiber.ErrInternalServerError
		}
	}

	if err := c.CustomerRepository.Delete(tx, customer); err != nil {
		c.Log.WithError(err).Error("error deleting customer")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error deleting customer")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *CustomerUseCase) Search(ctx context.Context, request *model.SearchCustomerRequest) ([]model.CustomerResponse, int64, error) {
	db := c.DB.WithContext(ctx)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	customer, total, err := c.CustomerRepository.Search(db, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting customer")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.CustomerResponse, len(customer))
	for i, customer := range customer {
		responses[i] = *converter.CustomerToResponse(&customer)
	}

	return responses, total, nil
}
