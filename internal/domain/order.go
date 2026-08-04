package domain

import (
	"time"
)

type Order struct {
	ID             uint        `gorm:"PrimaryKey json:"id"`
	UserId         uint        `json:"user_id`
	Status         string      `json:"status"`
	Amount         float64     `json:"amount"`
	TransactionId  uint        `json:transaction_id"`
	OrderRefNumber int         `json:"order_ref_number"`
	PaymentId      string      `json:"payment_id"`
	Items          []OrderItem `json:"items"`
	CreatedAt      time.Time   `gorm:"default:current_timestamp"`
	UpdatedAt      time.Time   `gorm:"default:current_timestamp"`
}
