package domain

import "time"

type Address struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2"`
	City         string    `json:"city"`
	PostCode     string    `json:"post_code"`
	Country      string    `json:"country"`
	UserId       uint      `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
