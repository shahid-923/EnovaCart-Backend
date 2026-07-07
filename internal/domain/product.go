package domain

import "time"

type Product struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	Name         string       `json:"name" gorm:"index;not null"`
	Description  string       `json:"description"` 
	CategoryId   uint         `json:"category_id"`
	ImageUrl     string       `json:"image_url"`
	Price        float64      `json:"price"` 
	UserId       int          `json:"user_id"`    // belongs to seller
	Stock        uint         `json:"stock"`
	CreatedAt    time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}