package payment

import (
	"errors"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v86"
	stripeSession "github.com/stripe/stripe-go/v86/checkout/session"
)

// LineItemData represents a single product line item for Stripe
type LineItemData struct {
	ProductName  string
	Description  string
	ImageURL     string
	Quantity     int64
	PriceInCents int64
}

// PaymentResponse contains the Stripe session details and order metadata
type PaymentResponse struct {
	SessionID  string
	OrderID    string
	PaymentURL string
	Amount     float64
}

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error)
	CreatePaymentWithItems(amount float64, userId uint, orderId string, items []LineItemData) (*stripe.CheckoutSession, error)
	GetPaymentStatus(pId string) (*stripe.CheckoutSession, error)
	CancelPayment(pId string) error
}

type payment struct {
	stripeSecretKey string
	successUrl      string
	cancelUrl       string
}

func NewPaymentClient(stripeSecretKey, successUrl, cancelUrl string) PaymentClient {
	return &payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		cancelUrl:       cancelUrl,
	}
}

// CreatePayment creates a basic Stripe session (legacy method for backward compatibility)
func (p payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error) {
	// Fallback to single line item if no items are provided
	items := []LineItemData{
		{
			ProductName:  fmt.Sprintf("Order %s", orderId),
			Description:  "Product from order",
			Quantity:     1,
			PriceInCents: int64(amount * 100),
		},
	}
	return p.CreatePaymentWithItems(amount, userId, orderId, items)
}

// CreatePaymentWithItems creates a Stripe session with multiple line items (BEST PRACTICE)
// This method accepts actual product information from the order
func (p payment) CreatePaymentWithItems(amount float64, userId uint, orderId string, items []LineItemData) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey

	// Validate inputs
	if len(items) == 0 {
		return nil, errors.New("at least one line item is required")
	}

	if orderId == "" {
		return nil, errors.New("order ID is required")
	}

	// Build line items array from actual product data
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(items))
	totalAmount := int64(0)

	for _, item := range items {
		if item.ProductName == "" || item.Quantity <= 0 || item.PriceInCents <= 0 {
			return nil, errors.New("invalid line item data: missing product name, quantity, or price")
		}

		// Create line item with actual product information
		lineItem := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				UnitAmount: stripe.Int64(item.PriceInCents),
				Currency:   stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(item.ProductName),
					Description: stripe.String(item.Description),
				},
			},
			Quantity: stripe.Int64(item.Quantity),
		}

		// Add image URL if available (improves customer experience)
		if item.ImageURL != "" {
			lineItem.PriceData.ProductData.Images = stripe.StringSlice([]string{item.ImageURL})
		}

		lineItems = append(lineItems, lineItem)
		totalAmount += item.PriceInCents * item.Quantity
	}

	// Verify total amount matches
	expectedTotal := int64(amount * 100)
	if totalAmount != expectedTotal {
		log.Printf("Warning: calculated total (%d cents) doesn't match provided amount (%d cents)", totalAmount, expectedTotal)
	}

	// Create Stripe session with comprehensive metadata
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(p.successUrl),
		CancelURL:          stripe.String(p.cancelUrl),
	}

	// Add metadata for order reconciliation - CRITICAL for tracking
	params.AddMetadata("order_id", orderId)
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))
	params.AddMetadata("created_at", fmt.Sprintf("%d", 0)) // Timestamp will be added by Stripe
	params.AddMetadata("total_items", fmt.Sprintf("%d", len(items)))

	// Log the payment session details for debugging
	log.Printf("Creating Stripe session - OrderID: %s, UserID: %d, Amount: %.2f USD, Items: %d",
		orderId, userId, amount, len(items))

	// Create Stripe session
	session, err := stripeSession.New(params)
	if err != nil {
		log.Printf("Error creating Stripe session: %v", err)
		return nil, errors.New("failed to create payment session")
	}

	// Log successful session creation with session ID (THIS IS HOW YOU SEE THE SESSION ID)
	log.Printf("✓ Stripe session created successfully - SessionID: %s, OrderID: %s, URL: %s",
		session.ID, orderId, session.URL)

	return session, nil
}

func (p payment) GetPaymentStatus(pId string) (*stripe.CheckoutSession, error) {
	if pId == "" {
		return nil, errors.New("payment ID is required")
	}

	stripe.Key = p.stripeSecretKey

	log.Printf("Fetching payment status for SessionID: %s", pId)

	session, err := stripeSession.Get(pId, nil)
	if err != nil {
		log.Printf("Error fetching session %s: %v", pId, err)
		return nil, errors.New("payment status fetch failed")
	}

	log.Printf("✓ Payment status retrieved - SessionID: %s, Status: %s, PaymentStatus: %s",
		pId, session.Status, session.PaymentStatus)

	return session, nil
}

func (p payment) CancelPayment(pId string) error {
	if pId == "" {
		return errors.New("payment ID is required")
	}

	stripe.Key = p.stripeSecretKey

	log.Printf("Cancelling payment session: %s", pId)

	_, err := stripeSession.Expire(pId, nil)
	if err != nil {
		log.Printf("Error cancelling session %s: %v", pId, err)
		return errors.New("failed to cancel payment session")
	}

	log.Printf("✓ Payment session cancelled successfully - SessionID: %s", pId)

	return nil
}
