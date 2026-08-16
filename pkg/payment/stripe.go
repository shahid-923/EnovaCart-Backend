package payment

import (
	"errors"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v86"
	stripeSession "github.com/stripe/stripe-go/v86/checkout/session"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error)
	GetPaymentStatus(pId string) (*stripe.CheckoutSession, error)
	CancelPayment(pId string) error
}

type payment struct {
	stripeSecretKey string
	successUrl      string
	cancelUrl       string
}

func NewPaymentClient(stripeSecretKey, successUrl, cancelUrl string) PaymentClient { // constructor function to create a new payment client
	return &payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		cancelUrl:       cancelUrl,
	}
}

func (p payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey

	amountInCents := int64(amount * 100)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					UnitAmount: stripe.Int64(amountInCents),
					Currency:   stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Electronics"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(p.successUrl),
		CancelURL:  stripe.String(p.cancelUrl),
	}

	// Add metadata here
	params.AddMetadata("order_id", orderId)
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))

	// Create Stripe session
	session, err := stripeSession.New(params)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return nil, errors.New("failed to create payment session")
	}

	return session, nil
}

func (p payment) GetPaymentStatus(pId string) (*stripe.CheckoutSession, error) {

	stripe.Key = p.stripeSecretKey
	session, err := stripeSession.Get(pId, nil)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		return nil, errors.New("payment get status failed")
	}
	return session, nil
}
func (p payment) CancelPayment(pId string) error {
	stripe.Key = p.stripeSecretKey

	_, err := stripeSession.Expire(pId, nil)
	if err != nil {
		return errors.New("failed to cancel payment session")
	}

	return nil
}