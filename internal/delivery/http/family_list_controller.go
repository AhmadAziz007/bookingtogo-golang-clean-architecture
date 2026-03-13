package http

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type FamilyListController struct {
	UseCase *usecase.FamilyListUseCase
	Log     *logrus.Logger
}

func NewFamilyListController(useCase *usecase.FamilyListUseCase, log *logrus.Logger) *FamilyListController {
	return &FamilyListController{
		UseCase: useCase,
		Log:     log,
	}
}

func (f *FamilyListController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateFamilyListRequest)
	if err := ctx.BodyParser(request); err != nil {
		f.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := f.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		f.Log.WithError(err).Error("failed to create family list")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.FamilyListResponse]{Data: response})
}

func (f *FamilyListController) GetByFlIdWithRelationCstId(ctx *fiber.Ctx) error {
	request, err := ctx.ParamsInt("customerId")
	if err != nil {
		f.Log.WithError(err).Error("failed to parse customer id")
		return fiber.ErrBadRequest
	}

	response, err := f.UseCase.GetByFlIDWithRelationCstId(ctx.UserContext(), request)
	if err != nil {
		f.Log.WithError(err).Error("failed to get family lists")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.FamilyListResponse]{Data: response})
}

func (f *FamilyListController) Update(ctx *fiber.Ctx) error {
	flId, err := ctx.ParamsInt("familyListId")
	if err != nil {
		f.Log.WithError(err).Error("failed to parse family list id")
		return fiber.ErrBadRequest
	}

	request := new(model.UpdateFamilyListRequest)
	if err := ctx.BodyParser(request); err != nil {
		f.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.FlID = flId

	response, err := f.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		f.Log.WithError(err).Error("failed to update family list")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.FamilyListResponse]{Data: response})
}

func (f *FamilyListController) Delete(ctx *fiber.Ctx) error {
	flId, err := ctx.ParamsInt("familyListId")
	if err != nil {
		f.Log.WithError(err).Error("failed to parse family list id")
		return fiber.ErrBadRequest
	}

	request := &model.DeleteFamilyListRequest{
		FlID: flId,
	}

	if err := f.UseCase.Delete(ctx.UserContext(), request); err != nil {
		f.Log.WithError(err).Error("failed to delete family list")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
