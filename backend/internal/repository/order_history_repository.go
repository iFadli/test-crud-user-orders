package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"test-crud-user-orders/internal/entity"
)

type OrderHistoryRepository interface {
	Create(ctx context.Context, orderHistory *entity.OrderHistory) (*entity.OrderHistory, error)
	Update(ctx context.Context, orderHistory *entity.OrderHistory) error
	GetByID(ctx context.Context, id int) (*entity.OrderHistory, error)
	GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.OrderHistory, error)
	GetByUserID(ctx context.Context, userID, limit, offset int) ([]*entity.OrderHistory, error)
	SoftDelete(ctx context.Context, id int) error
	CountData(ctx context.Context, userId int) int64
}

type orderHistoryRepository struct {
	db *gorm.DB
}

func NewOrderHistoryRepository(db *gorm.DB) OrderHistoryRepository {
	return &orderHistoryRepository{db}
}

func (r *orderHistoryRepository) Create(ctx context.Context, orderHistory *entity.OrderHistory) (*entity.OrderHistory, error) {
	execDB := r.db.WithContext(ctx).Create(&orderHistory)
	if execDB.Error != nil {
		return orderHistory, execDB.Error
	}

	var count int64
	r.db.Model(&entity.User{}).Where("id = ? AND first_order IS NULL", orderHistory.UserID).Count(&count)
	if count == 1 {
		now := time.Now()
		r.db.Model(&entity.User{}).Where("id = ? AND first_order IS NULL", orderHistory.UserID).Update("first_order", &now)
	}

	return orderHistory, nil
}

func (r *orderHistoryRepository) Update(ctx context.Context, orderHistory *entity.OrderHistory) error {
	// Check if the related User is not soft-deleted
	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", orderHistory.UserID).First(&user).Error; err != nil {
		return fmt.Errorf("user data not found")
	}

	// Check if the related OrderItem is not soft-deleted
	var orderItem entity.OrderItem
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", orderHistory.OrderItemID).First(&orderItem).Error; err != nil {
		return fmt.Errorf("order item data not found")
	}

	// Update the OrderHistory if the related data is not soft-deleted
	if err := r.db.WithContext(ctx).Model(orderHistory).Updates(&orderHistory).Error; err != nil {
		return err
	}

	return nil
}

func (r *orderHistoryRepository) SoftDelete(ctx context.Context, id int) error {
	orderHistory := &entity.OrderHistory{ID: id}
	err := r.db.Delete(orderHistory).Error
	if err != nil {
		return fmt.Errorf("error soft-deleting order history with ID %d: %s", id, err.Error())
	}

	return nil
}

func (r *orderHistoryRepository) GetByID(ctx context.Context, id int) (*entity.OrderHistory, error) {
	orderHistory := &entity.OrderHistory{}
	err := r.db.
		Preload("OrderItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		WithContext(ctx).First(orderHistory, id).Error
	if err != nil {
		return nil, err
	}

	return orderHistory, nil
}

func (r *orderHistoryRepository) GetByUserID(ctx context.Context, userID, limit, offset int) ([]*entity.OrderHistory, error) {
	var orderHistories []*entity.OrderHistory
	err := r.db.
		Preload("OrderItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&orderHistories).Error
	if err != nil {
		return nil, err
	}
	return orderHistories, nil
}

func (r *orderHistoryRepository) GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.OrderHistory, error) {
	var orderHistory []*entity.OrderHistory

	err := r.db.
		Preload("OrderItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Limit(limit).Offset(offset).Find(&orderHistory).Error
	if err != nil {
		return nil, err
	}

	return orderHistory, nil
}

func (r *orderHistoryRepository) CountData(ctx context.Context, userID int) int64 {
	var count int64
	var orderHistory []*entity.OrderHistory

	if userID < 1 {
		r.db.Find(&orderHistory).Count(&count)
	} else {
		r.db.Where("user_id = ?", userID).Find(&orderHistory).Count(&count)
	}
	return count
}
