package domain

import (
	"time"
)

type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserId    uint      `gorm:"column:user_id" json:"user_id"`
	ProductId uint      `gorm:"column:product_id" json:"product_id"`
	SellerId  uint      `gorm:"column:seller_id" json:"seller_id"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"image_url"`
	Price     int64     `json:"price"`
	Qty       uint      `json:"qty"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}
