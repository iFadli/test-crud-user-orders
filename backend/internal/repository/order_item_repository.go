package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"test-crud-user-orders/internal/entity"
)

type OrderItemRepository interface {
	Create(ctx context.Context, orderItem *entity.OrderItem) error
	Update(ctx context.Context, orderItem *entity.OrderItem) error
	GetByID(ctx context.Context, id int) (*entity.OrderItem, error)
	GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.OrderItem, error)
	SoftDelete(ctx context.Context, id int) error
	CountData(ctx context.Context) int64
}

type orderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepository{db}
}

func (r *orderItemRepository) Create(ctx context.Context, orderItem *entity.OrderItem) error {
	execDB := r.db.WithContext(ctx).Create(&orderItem)
	if execDB.Error != nil {
		return execDB.Error
	}
	return nil
}

func (r *orderItemRepository) Update(ctx context.Context, orderItem *entity.OrderItem) error {
	return r.db.WithContext(ctx).Model(orderItem).Updates(&orderItem).Error
}

func (r *orderItemRepository) SoftDelete(ctx context.Context, id int) error {
	orderItem := &entity.OrderItem{ID: id}
	err := r.db.Delete(orderItem).Error
	if err != nil {
		return fmt.Errorf("error soft-deleting order item with ID %d: %s", id, err.Error())
	}
	return nil
}

func (r *orderItemRepository) GetByID(ctx context.Context, id int) (*entity.OrderItem, error) {
	orderItem := &entity.OrderItem{}
	err := r.db.WithContext(ctx).First(orderItem, id).Error
	if err != nil {
		return nil, err
	}
	return orderItem, nil
}

func (r *orderItemRepository) GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.OrderItem, error) {
	var orderItems []*entity.OrderItem

	err := r.db.Limit(limit).Offset(offset).Find(&orderItems).Error
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return orderItems, nil
}

func (r *orderItemRepository) CountData(ctx context.Context) int64 {
	var count int64
	var orderItems []*entity.OrderItem

	r.db.Find(&orderItems).Count(&count)
	return count
}
