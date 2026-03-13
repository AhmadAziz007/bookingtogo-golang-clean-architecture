package http

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type NationalityController struct {
	UseCase *usecase.NationalityUseCase
	Log     *logrus.Logger
}

func NewNationalityController(useCase *usecase.NationalityUseCase, log *logrus.Logger) *NationalityController {
	return &NationalityController{
		UseCase: useCase,
		Log:     log,
	}
}

func (n *NationalityController) List(ctx *fiber.Ctx) error {
	responses, err := n.UseCase.List(ctx.UserContext())
	if err != nil {
		n.Log.WithError(err).Error("failed to list nationalities")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.NationalityResponse]{Data: responses})
}

func (n *NationalityController) GetByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("nationalityId")
	if err != nil {
		n.Log.WithError(err).Error("failed to parse nationality id")
		return fiber.ErrBadRequest
	}

	response, err := n.UseCase.GetByID(ctx.UserContext(), id)
	if err != nil {
		n.Log.WithError(err).Error("failed to get nationality by id")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.NationalityResponse]{Data: response})
}

func (n *NationalityController) GetByCode(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	if code == "" {
		return fiber.ErrBadRequest
	}

	response, err := n.UseCase.GetByCode(ctx.UserContext(), code)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.NationalityResponse]{Data: response})
}

func (n *NationalityController) Search(ctx *fiber.Ctx) error {
	name := ctx.Query("name", "")
	code := ctx.Query("code", "")

	responses, err := n.UseCase.List(ctx.UserContext())
	if err != nil {
		n.Log.WithError(err).Error("failed to search nationalities")
		return err
	}

	var filtered []model.NationalityResponse
	if name != "" || code != "" {
		for _, nationality := range responses {
			if (name == "" || containsIgnoreCase(nationality.NationalityName, name)) &&
				(code == "" || containsIgnoreCase(nationality.NationalityCode, code)) {
				filtered = append(filtered, nationality)
			}
		}
	} else {
		filtered = responses
	}

	return ctx.JSON(model.WebResponse[[]model.NationalityResponse]{Data: filtered})
}

func containsIgnoreCase(s, substr string) bool {
	return len(substr) == 0 ||
		len(s) >= len(substr) &&
			strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
