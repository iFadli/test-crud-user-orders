package entity

import (
	"time"
)

type OrderHistory struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	UserID       int        `json:"-" gorm:"not null;foreignkey:UserID"`
	OrderItemID  int        `json:"-" gorm:"not null;foreignkey:OrderItemID"`
	Descriptions string     `json:"descriptions" gorm:"size:255"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	User         *User      `json:"user,omitempty" gorm:"foreignkey:UserID`
	OrderItem    *OrderItem `json:"order_item,omitempty" gorm:"foreignkey:OrderItemID`
}

type CreateOrderHistory struct {
	UserID       int    `json:"user_id" validate:"required"`
	OrderItemID  int    `json:"order_item_id" validate:"required"`
	Descriptions string `json:"descriptions" validate:"required"`
}

func (OrderHistory) TableName() string {
	return "order_histories"
}
