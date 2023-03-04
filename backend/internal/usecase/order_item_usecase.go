package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"test-crud-user-orders/internal/entity"
	"test-crud-user-orders/internal/repository"
)

type OrderItemUseCase interface {
	GetByID(ctx context.Context, id int) (*entity.OrderItem, error)
	Create(ctx context.Context, orderItem *entity.OrderItem) error
	Update(ctx context.Context, orderItem *entity.OrderItem) error
	Delete(ctx context.Context, id int) error
	GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.OrderItem, error)
	CountData(ctx context.Context) int64
}

type orderItemUseCase struct {
	orderItemRepo repository.OrderItemRepository
	redisClient   *redis.Client
}

func NewOrderItemUseCase(orderItemRepo repository.OrderItemRepository, redisClient *redis.Client) OrderItemUseCase {
	return &orderItemUseCase{
		orderItemRepo: orderItemRepo,
		redisClient:   redisClient,
	}
}

func (uc *orderItemUseCase) GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.OrderItem, error) {
	return uc.orderItemRepo.GetAllPagination(ctx, limit, offset)
}

func (uc *orderItemUseCase) GetByID(ctx context.Context, id int) (*entity.OrderItem, error) {
	return uc.orderItemRepo.GetByID(ctx, id)
}

func (uc *orderItemUseCase) Create(ctx context.Context, orderItem *entity.OrderItem) error {
	if err := uc.orderItemRepo.Create(ctx, orderItem); err != nil {
		return fmt.Errorf("error creating order item: %s", err.Error())
	}
	return nil
}

func (uc *orderItemUseCase) Update(ctx context.Context, orderItem *entity.OrderItem) error {
	orderItemDB, err := uc.orderItemRepo.GetByID(ctx, orderItem.ID)
	if err != nil {
		return err
	}
	if orderItemDB == nil {
		return errors.New("order item not found")
	}

	orderItemDB.Name = orderItem.Name
	orderItemDB.Price = orderItem.Price
	orderItemDB.ExpiredAt = orderItem.ExpiredAt

	if err := uc.orderItemRepo.Update(ctx, orderItemDB); err != nil {
		return fmt.Errorf("error updating order item with ID %d: %s", orderItem.ID, err.Error())
	}

	// Check Redis Connection with method Ping()
	_, err = uc.redisClient.Ping(ctx).Result()
	if err == nil {
		// Delete Redis Data of this ID
		key := "order_items:all"
		if err := uc.redisClient.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("error deleting data from Redis cache: %s", err.Error())
		}
	}

	return nil
}

func (uc *orderItemUseCase) Delete(ctx context.Context, id int) error {
	orderItemDB, err := uc.orderItemRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if orderItemDB == nil {
		return errors.New("order item not found")
	}

	if err := uc.orderItemRepo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("error deleting order item with ID %d: %s", id, err.Error())
	}

	// Check Redis Connection with method Ping()
	_, err = uc.redisClient.Ping(ctx).Result()
	if err == nil {
		// Delete the cached data since it has been deleted
		key := "order_items:all"
		if err := uc.redisClient.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("error deleting data from Redis cache: %s", err.Error())
		}
	}

	return nil
}

func (uc *orderItemUseCase) CountData(ctx context.Context) int64 {
	return uc.orderItemRepo.CountData(ctx)
}
