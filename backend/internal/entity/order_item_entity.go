package entity

import (
	"gorm.io/gorm"
	"time"
)

type OrderItem struct {
	ID             int            `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:100;not null" validate:"required"`
	Price          int            `json:"price" gorm:"not null" validate:"required"`
	ExpiredAt      time.Time      `json:"expired_at" gorm:"index" validate:"required"`
	CreatedAt      time.Time      `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	OrderHistories []OrderHistory `json:"order_histories,omitempty" gorm:"foreignkey:OrderItemID"`
}

type CreateOrderItem struct {
	Name       string `json:"name" validate:"required"`
	Price      int    `json:"price" validate:"required"`
	ExpiredDay int    `json:"expired_days" validate:"required"`
}
