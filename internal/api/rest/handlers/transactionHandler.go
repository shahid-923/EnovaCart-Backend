package handlers

import (
	"ecommerce-app/internal/api/rest"
	"ecommerce-app/internal/helper"
	"ecommerce-app/internal/repository"
	"ecommerce-app/internal/service"
	"ecommerce-app/pkg/payment"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
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

	// PUBLIC: Stripe calls this directly
	app.Post("/payment/webhook", handler.StripeWebhook)

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

	// 2. Get the latest order for this user (should be created by POST /order)
	orders, err := h.userSvc.GetOrders(user)
	if err != nil || len(orders) == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "no order found - please create an order first",
			"error":   "use POST /order endpoint before making payment",
		})
	}

	// Get the most recent order
	order := orders[0]
	for _, o := range orders {
		if o.CreatedAt.After(order.CreatedAt) {
			order = o
		}
	}

	// 3. Check if this order already has an active payment
	if order.PaymentId != "" {
		// Order already has a payment ID, verify it with Stripe
		stripeSession, stripeErr := h.paymentClient.GetPaymentStatus(order.PaymentId)

		if stripeErr == nil && stripeSession != nil && stripeSession.Status == "open" {
			// Stripe session is still valid, reuse it
			return ctx.Status(http.StatusOK).JSON(fiber.Map{
				"message":     "payment session already exists for this order",
				"order_id":    order.OrderRefNumber,
				"payment_url": stripeSession.URL,
				"session_id":  stripeSession.ID,
				"status":      "pending_payment",
			})
		}
		// If Stripe session is not open, continue to create a new one
	}

	// 4. Calculate total amount from order items
	amount := order.Amount
	if amount <= 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "order amount is invalid",
		})
	}

	// 5. Create a new Stripe Checkout Session
	sessionResult, err := h.paymentClient.CreatePayment(
		amount,
		user.ID,
		order.OrderRefNumber, // Use the order reference number
	)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to create payment session",
			"error":   err.Error(),
		})
	}

	// 6. Save the payment and link it to the order
	// This updates both Payment table and Order.PaymentId
	err = h.svc.StoreCreatedPaymentWithOrder(
		user.ID,
		&order,
		sessionResult,
		amount,
		h.userSvc.Repo, // Pass repository to update order
	)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to store payment session",
			"error":   err.Error(),
		})
	}

	// 7. Return comprehensive payment response
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":     "payment session created successfully",
		"order_id":    order.OrderRefNumber,
		"session_id":  sessionResult.ID,
		"payment_url": sessionResult.URL,
		"amount":      amount,
		"status":      "pending_payment",
		"instruction": "redirect user to payment_url to complete payment",
	})
}

func (h *TransactionHandler) StripeWebhook(ctx fiber.Ctx) error {

	payload := ctx.Body()
	signature := ctx.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(
		payload,
		signature,
		os.Getenv("STRIPE_WEBHOOK_SECRET"),
	)

	if err != nil {
		log.Printf("Stripe webhook signature verification failed: %v", err)

		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid stripe webhook signature",
		})
	}

	switch event.Type {

	case stripe.EventTypeCheckoutSessionCompleted:

		var session stripe.CheckoutSession

		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Failed to parse checkout session: %v", err)

			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid checkout session",
			})
		}

		orderID := session.Metadata["order_id"]
		userID := session.Metadata["user_id"]

		log.Printf(
			"Stripe payment completed - SessionID: %s, OrderID: %s, UserID: %s",
			session.ID,
			orderID,
			userID,
		)

	default:
		log.Printf("Unhandled Stripe event: %s", event.Type)
	}

	return ctx.SendStatus(http.StatusOK)
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
