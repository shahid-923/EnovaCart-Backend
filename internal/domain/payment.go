package domain

import "time"

type Payment struct {
	ID            uint      `json: "id" gorm: "primaryKey"`
	UserId        uint      `json: "user_id"`
	CaptureMethod string    `json:capture_methdod"`
	Amount        float64   `json:amount"`
	TransactionId string    `json:"transaction_id"`
	CustomerId    string    `json:"customer_id"` // if stripe customer
	PaymentId     string    `json:"payment_id"`  // payment id
	Status        string    `json:"status"`      // initially success fail
	Response      string    `json:"response"`
	CreatedAt     time.Time `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time `gorm:"default:current_timestamp"`
}
