package repository

import (
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId, id uint) (dto.SellerOrderDetails, error)
}

type transactionStorage struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionStorage{
		db: db,
	}
}

func (r *transactionStorage) CreatePayment(payment *domain.Payment) error {

}

func (r *transactionStorage) FindOrders(uId uint) ([]domain.OrderItem, error) {

}

func (r *transactionStorage) FindOrderById(uId, id uint) (dto.SellerOrderDetails, error) {

}
