package domain

import "time"

type Payment struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	UserId        uint          `json:"user_id"`
	CaptureMethod string        `json:"capture_method"`
	Amount        float64       `json:"amount"`
	TransactionId string        `json:"transaction_id"`
	CustomerId    string        `json:"customer_id"` // for stripe customer
	PaymentId     string        `json:"payment_id"`
	OrderId       string        `json:"order_id"`
	PaymentUrl    string        `json:"payment_url"`
	Status        PaymentStatus `json:"status" gorm:"default:initial"`
	Response      string        `json:"response"`
	CreatedAt     time.Time     `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time     `gorm:"default:current_timestamp"`
}

type PaymentStatus string

const (
	PaymentStatusInitial   PaymentStatus = "initial"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)
