package repository

import (
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.Order, error)
	FindOrderByID(uId, id uint) (dto.SellerOrderDetails, error)
	FindInitialPayment(uId uint) (*domain.Payment, error)
	UpdatePayment(payment *domain.Payment) error
}

type transactionStorage struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionStorage{
		db: db,
	}
}

func (t *transactionStorage) CreatePayment(payment *domain.Payment) error {
	return t.db.Create(payment).Error
}

func (t *transactionStorage) FindInitialPayment(uId uint) (*domain.Payment, error) {
	var payment domain.Payment

	err := t.db.
		Where("user_id = ? AND status = ?", uId, domain.PaymentStatusInitial).
		Order("created_at DESC").
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (t *transactionStorage) UpdatePayment(payment *domain.Payment) error {
	return t.db.Save(payment).Error
}

func (t *transactionStorage) CreateOrder(order domain.Order) error {
	return t.db.Create(&order).Error
}

func (t *transactionStorage) FindOrders(userID uint) ([]domain.Order, error) {
	var orders []domain.Order

	err := t.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (t *transactionStorage) FindOrderByID(uId uint, id uint) (dto.SellerOrderDetails, error) {

	var order dto.SellerOrderDetails

	// your query here

	return order, nil
}
