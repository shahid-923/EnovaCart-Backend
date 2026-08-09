package handlers

import (
	"ecommerce-app/internal/api/rest"
	"ecommerce-app/internal/helper"
	"ecommerce-app/internal/repository"
	"ecommerce-app/internal/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	svc *service.TransactionService
}

func initializeTransactionService(db *gorm.DB, auth helper.Auth) *service.TransactionService {
	return &service.TransactionService{
		Repo: repository.NewTransactionRepository(db),
		Auth: auth,
	}
}

func SetupTransactionRoutes(as *rest.RestHandler) {
	app := as.App

	svc := initializeTransactionService(as.DB, as.Auth)

	handler := TransactionHandler{
		svc: svc,
	}

	secRoute := app.Group("/", as.Auth.Authorize())
	secRoute.Get("/payment", handler.MakePayment)

	sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller())
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *TransactionHandler) GetOrders(ctx fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *TransactionHandler) GetOrderDetails(ctx fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}
