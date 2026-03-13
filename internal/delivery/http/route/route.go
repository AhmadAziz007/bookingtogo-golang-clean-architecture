package route

import (
	"golang-clean-architecture/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                   *fiber.App
	UserController        *http.UserController
	ContactController     *http.ContactController
	AddressController     *http.AddressController
	CustomerController    *http.CustomerController
	FamilyListController  *http.FamilyListController
	NationalityController *http.NationalityController
	AuthMiddleware        fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/users/_login", c.UserController.Login)

	c.App.Get("/api/nationalities", c.NationalityController.List)
	c.App.Get("/api/nationalities/code/:code", c.NationalityController.GetByCode)
	c.App.Get("/api/nationalities/:nationalityId", c.NationalityController.GetByID)

	c.App.Post("/api/customers", c.CustomerController.Create)
	c.App.Get("/api/customers", c.CustomerController.List)
	c.App.Get("/api/customers/:customerId", c.CustomerController.Get)
	c.App.Put("/api/customers/:customerId", c.CustomerController.Update)
	c.App.Delete("/api/customers/:customerId", c.CustomerController.Delete)

	c.App.Post("/api/family-lists", c.FamilyListController.Create)
	c.App.Get("/api/customers/:customerId/family-lists", c.FamilyListController.GetByFlIdWithRelationCstId)
	c.App.Put("/api/family-lists/:familyListId", c.FamilyListController.Update)
	c.App.Delete("/api/family-lists/:familyListId", c.FamilyListController.Delete)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Delete("/api/users", c.UserController.Logout)
	c.App.Patch("/api/users/_current", c.UserController.Update)
	c.App.Get("/api/users/_current", c.UserController.Current)

	c.App.Get("/api/contacts", c.ContactController.List)
	c.App.Post("/api/contacts", c.ContactController.Create)
	c.App.Put("/api/contacts/:contactId", c.ContactController.Update)
	c.App.Get("/api/contacts/:contactId", c.ContactController.Get)
	c.App.Delete("/api/contacts/:contactId", c.ContactController.Delete)

	c.App.Get("/api/contacts/:contactId/addresses", c.AddressController.List)
	c.App.Post("/api/contacts/:contactId/addresses", c.AddressController.Create)
	c.App.Put("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Update)
	c.App.Get("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Get)
	c.App.Delete("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Delete)
}
