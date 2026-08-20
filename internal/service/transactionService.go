package service

import (
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"
	"ecommerce-app/internal/helper"
	"ecommerce-app/internal/repository"
	"ecommerce-app/pkg/payment"
	"errors"

	"github.com/stripe/stripe-go/v86"
)

type TransactionService struct {
	Repo          repository.TransactionRepository
	Auth          helper.Auth
	PaymentClient payment.PaymentClient
}

func NewTransactionService(r repository.TransactionRepository, auth helper.Auth, paymentClient payment.PaymentClient) *TransactionService {
	return &TransactionService{
		Repo:          r,
		Auth:          auth,
		PaymentClient: paymentClient,
	}
}

func (s *TransactionService) GetActivePayment(uId uint) (*domain.Payment, error) {
	return s.Repo.FindInitialPayment(uId)
}

// StoreCreatedPayment stores payment info (legacy method for backward compatibility)
func (s *TransactionService) StoreCreatedPayment(uId uint, ps *stripe.CheckoutSession, amount float64) error {
	payment := domain.Payment{
		UserId:     uId,
		Amount:     amount,
		Status:     domain.PaymentStatusInitial,
		PaymentUrl: ps.URL,
		PaymentId:  ps.ID,
	}
	return s.Repo.CreatePayment(&payment)
}

// StoreCreatedPaymentWithOrder links payment to an existing order and updates order with payment ID
// This is the NEW method that properly links payment to order
func (s *TransactionService) StoreCreatedPaymentWithOrder(uId uint, order *domain.Order, ps *stripe.CheckoutSession, amount float64, userRepo interface{}) error {
	// Create payment record with order reference
	payment := domain.Payment{
		UserId:     uId,
		Amount:     amount,
		Status:     domain.PaymentStatusInitial,
		PaymentUrl: ps.URL,
		PaymentId:  ps.ID,
		OrderId:    order.OrderRefNumber, // Link payment to order using orderRefNumber
	}

	// Store payment in DB
	err := s.Repo.CreatePayment(&payment)
	if err != nil {
		return err
	}

	// Update order with Stripe payment ID
	order.PaymentId = ps.ID
	order.Status = "pending_payment" // Order waiting for payment completion

	// Type assert to UserRepository to update order
	ur, ok := userRepo.(interface{ UpdateOrder(domain.Order) error })
	if !ok {
		return errors.New("invalid user repository type")
	}

	return ur.UpdateOrder(*order)
}

func (s *TransactionService) CancelActivePayment(uId uint) error {

	paymentRecord, err := s.Repo.FindInitialPayment(uId)
	if err != nil {
		return err
	}

	if paymentRecord == nil {
		return errors.New("no active payment found")
	}

	// Cancel/expire Stripe session
	err = s.PaymentClient.CancelPayment(paymentRecord.PaymentId)
	if err != nil {
		return err
	}

	// Update DB payment status
	paymentRecord.Status = domain.PaymentStatusCancelled

	return s.Repo.UpdatePayment(paymentRecord)
}

func (s *TransactionService) GetOrders(u domain.User) ([]domain.Order, error) {
	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *TransactionService) GetOrderDetails(u domain.User, id uint) (dto.SellerOrderDetails, error) {

	order, err := s.Repo.FindOrderByID(u.ID, id)
	if err != nil {
		return dto.SellerOrderDetails{}, err
	}

	return order, nil
}
