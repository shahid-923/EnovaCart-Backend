package handlers

import (
	"ecommerce-app/internal/api/rest"
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/repository"
	"ecommerce-app/internal/service"
	"ecommerce-app/pkg/payment"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
)

type WebhookHandler struct {
	svc           *service.TransactionService
	userSvc       *service.UserService
	paymentClient payment.PaymentClient
}

func SetupWebhookRoutes(as *rest.RestHandler) {
	svc := initializeTransactionService(as.DB, as.Auth, as.Pc)

	useSvc := service.UserService{
		Repo:   repository.NewUserRepository(as.DB),
		CRepo:  repository.NewCatalogueRepository(as.DB),
		Auth:   as.Auth,
		Config: as.Config,
	}

	handler := WebhookHandler{
		svc:           svc,
		paymentClient: as.Pc,
		userSvc:       &useSvc,
	}

	// Public endpoint for Stripe webhooks (no auth required)
	as.App.Post("/webhook/stripe", handler.HandleStripeWebhook)
}

// HandleStripeWebhook processes Stripe webhook events
func (h *WebhookHandler) HandleStripeWebhook(ctx fiber.Ctx) error {
	// Get Stripe secret for webhook verification
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Printf("STRIPE_WEBHOOK_SECRET not set")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "webhook secret not configured",
		})
	}

	// Get request body using Fiber's Body() method
	body := ctx.Body()

	// Get signature from header
	signatureHeader := ctx.Get("Stripe-Signature")

	// Verify webhook signature
	event, err := webhook.ConstructEvent(body, signatureHeader, webhookSecret)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid signature",
		})
	}

	log.Printf("✓ Webhook signature verified - Event Type: %s", event.Type)

	// Handle different event types
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Error parsing session data: %v", err)
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid session data",
			})
		}
		return h.handleCheckoutSessionCompleted(ctx, &session)

	case "checkout.session.expired":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Error parsing session data: %v", err)
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid session data",
			})
		}
		return h.handleCheckoutSessionExpired(ctx, &session)

	case "charge.refunded":
		var charge stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &charge); err != nil {
			log.Printf("Error parsing charge data: %v", err)
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid charge data",
			})
		}
		return h.handleChargeRefunded(ctx, &charge)

	default:
		log.Printf("Unhandled event type: %s", event.Type)
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"received": true,
		})
	}
}

// handleCheckoutSessionCompleted handles successful payment
func (h *WebhookHandler) handleCheckoutSessionCompleted(ctx fiber.Ctx, session *stripe.CheckoutSession) error {
	sessionID := session.ID
	orderID := session.Metadata["order_id"]
	userIDStr := session.Metadata["user_id"]

	log.Printf("Payment completed - SessionID: %s, OrderID: %s, Amount: %d", sessionID, orderID, session.AmountTotal)

	if orderID == "" || userIDStr == "" {
		log.Printf("Missing metadata in session: %s", sessionID)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "missing order/user metadata",
		})
	}

	// Parse user ID
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id",
		})
	}

	// Find order by reference number and user ID
	order, err := h.userSvc.Repo.FindOrderByRefNumber(orderID, uint(userID))
	if err != nil {
		log.Printf("Order not found - OrderID: %s, UserID: %d, Error: %v", orderID, userID, err)
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "order not found",
		})
	}

	// Update order status to completed
	order.Status = "completed"
	order.TransactionId = session.ID

	err = h.userSvc.Repo.UpdateOrder(order)
	if err != nil {
		log.Printf("Error updating order status: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update order",
		})
	}

	// Update payment status to success
	paymentRecord, err := h.svc.Repo.FindPaymentBySessionID(sessionID)
	if err == nil && paymentRecord != nil {
		paymentRecord.Status = domain.PaymentStatusSuccess
		paymentRecord.TransactionId = session.ID
		err = h.svc.Repo.UpdatePayment(paymentRecord)
		if err != nil {
			log.Printf("Error updating payment status: %v", err)
		}
	}

	log.Printf("✓ Order updated successfully - OrderID: %s, Status: completed", orderID)

	// TODO: Send confirmation email to user and seller
	// TODO: Delete cart items for this user after successful payment

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"received": true,
		"status":   "payment_confirmed",
		"order_id": orderID,
	})
}

// handleCheckoutSessionExpired handles expired payment sessions
func (h *WebhookHandler) handleCheckoutSessionExpired(ctx fiber.Ctx, session *stripe.CheckoutSession) error {
	orderID := session.Metadata["order_id"]
	log.Printf("Payment session expired - SessionID: %s, OrderID: %s", session.ID, orderID)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"received": true,
		"status":   "session_expired",
	})
}

// handleChargeRefunded handles refunds
func (h *WebhookHandler) handleChargeRefunded(ctx fiber.Ctx, charge *stripe.Charge) error {
	log.Printf("Payment refunded - ChargeID: %s, Amount: %d", charge.ID, charge.Amount)

	// TODO: Update order status to refunded
	// TODO: Restore cart items or notify user

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"received": true,
		"status":   "refund_processed",
	})
}
