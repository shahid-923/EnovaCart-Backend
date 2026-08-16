package handlers

import (
	"ecommerce-app/internal/api/rest"
	"ecommerce-app/internal/helper"
	"ecommerce-app/internal/repository"
	"ecommerce-app/internal/service"
	"ecommerce-app/pkg/payment"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	svc           *service.TransactionService
	userSvc       *service.UserService
	paymentClient payment.PaymentClient
}

func initializeTransactionService(db *gorm.DB, auth helper.Auth, paymentClient payment.PaymentClient) *service.TransactionService {
	return service.NewTransactionService(
		repository.NewTransactionRepository(db),
		auth,
		paymentClient,
	)
}

func SetupTransactionRoutes(as *rest.RestHandler) {

	app := as.App
	svc := initializeTransactionService(as.DB, as.Auth, as.Pc)

	useSvc := service.UserService{
		Repo:   repository.NewUserRepository(as.DB),
		CRepo:  repository.NewCatalogueRepository(as.DB),
		Auth:   as.Auth,
		Config: as.Config,
	}

	handler := TransactionHandler{
		svc:           svc,
		paymentClient: as.Pc,
		userSvc:       &useSvc,
	}

	secRoute := app.Group("/", as.Auth.Authorize())
	secRoute.Get("/payment", handler.MakePayment)
	secRoute.Post("/payment/cancel", handler.CancelPayment)

	sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller())
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx fiber.Ctx) error {

	// 1. Get the currently authenticated user

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"error":   err.Error(),
		})
	}

	
	// 2. Check whether the user already has an initial payment
	
	activePayment, err := h.svc.GetActivePayment(user.ID)

	if err == nil && activePayment != nil && activePayment.ID > 0 {
	
		// 3. The payment exists in our DB.
		//    Now verify its actual status with Stripe.

		stripeSession, stripeErr := h.paymentClient.GetPaymentStatus(
			activePayment.PaymentId,
		)

	
		// 4. If Stripe session is still open, reuse it.
		//    This prevents creating multiple payment sessions.

		if stripeErr == nil &&
			stripeSession != nil &&
			stripeSession.Status == "open" {

			return ctx.Status(http.StatusOK).JSON(fiber.Map{
				"message":     "payment session already exists",
				"payment_url": activePayment.PaymentUrl,
			})
		}

		// 5. If Stripe session is not open anymore, we don't
		//    reuse the old URL.
		//    Continue below and create a new payment session.
	
	}

	// 6. Get the user's cart and calculate the total amount
	_, amount, err := h.userSvc.FindCart(user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch cart",
			"error":   err.Error(),
		})
	}

	// 7. Do not create a Stripe payment for an empty cart
	if amount <= 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "cart is empty",
		})
	}

	// 8. Generate a unique order reference.
	//    Your GenerateOrderID() returns string.
	orderID, err := helper.GenerateOrderID()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate order id",
			"error":   err.Error(),
		})
	}

	// 9. Create a new Stripe Checkout Session
	sessionResult, err := h.paymentClient.CreatePayment(
		amount,
		user.ID,
		orderID,
	)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to create payment session",
			"error":   err.Error(),
		})
	}

	// 10. Save the newly created payment session in DB
	err = h.svc.StoreCreatedPayment(
		user.ID,
		sessionResult,
		amount,
	)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to store payment session",
			"error":   err.Error(),
		})
	}

	// 11. Return the new Stripe Checkout URL
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":     "payment session created successfully",
		"payment_url": sessionResult.URL,
	})
}

func (h *TransactionHandler) CancelPayment(ctx fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	err = h.svc.CancelActivePayment(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to cancel payment",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "payment cancelled successfully",
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
