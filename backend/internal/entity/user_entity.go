package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             int            `json:"id" gorm:"primaryKey"`
	FullName       string         `json:"name" gorm:"size:100;not null"`
	FirstOrder     *time.Time     `json:"first_order" gorm:"index;null"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	OrderHistories []OrderHistory `json:"order_histories,omitempty" gorm:"foreignkey:UserID"`
}

type CreateUser struct {
	FullName string `json:"name" validate:"required"`
}
