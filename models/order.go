package models

import "time"

type Order struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time
	ProductRefer int    `json:"product_id"`
	SerialNumber string `gorm:"foreignKey:ProductRefer"`
	UserRefer    int    `json:"user_id"`
	User         User   `gorm:"foreignKey:UserRefer"`
}
