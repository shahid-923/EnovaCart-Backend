package domain

import "time"

type Category struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	Name         string       `json:"name" gorm:"index;not null"`
	ParentId     uint         `json:"parent_id"`
	Products     []Product    `json:"products"`
	ImageUrl     string       `json:"image_url"`
	DisplayOrder int          `json:"display_order"`
	CreatedAt    time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}
