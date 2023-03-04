package usecase

import (
	"context"
	"errors"
	"time"

	"test-crud-user-orders/internal/entity"
	"test-crud-user-orders/internal/repository"
)

type OrderHistoryUseCase interface {
	Create(ctx context.Context, userID int, orderItemID int, descriptions string) (*entity.OrderHistory, error)
	Update(ctx context.Context, id int, userID int, orderItemID int, descriptions string) error
	GetByID(ctx context.Context, id int) (*entity.OrderHistory, error)
	GetByUserID(ctx context.Context, userID, limit, offset int) ([]*entity.OrderHistory, error)
	GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.OrderHistory, error)
	CountData(ctx context.Context, userId int) int64
}

type orderHistoryUseCase struct {
	orderHistoryRepo repository.OrderHistoryRepository
	orderItemRepo    repository.OrderItemRepository
	userRepo         repository.UserRepository
}

func NewOrderHistoryUseCase(
	orderHistory repository.OrderHistoryRepository,
	orderItem repository.OrderItemRepository,
	user repository.UserRepository,
) OrderHistoryUseCase {
	return &orderHistoryUseCase{orderHistory, orderItem, user}
}

func (uc *orderHistoryUseCase) Create(ctx context.Context, userID int, orderItemID int, descriptions string) (*entity.OrderHistory, error) {
	userData, errUser := uc.userRepo.GetByID(ctx, userID)
	if errUser != nil {
		return nil, errors.New("user not found")
	}
	orderItemData, errOrderItem := uc.orderItemRepo.GetByID(ctx, orderItemID)
	if errOrderItem != nil {
		return nil, errors.New("order item not found")
	}

	orderHistory := &entity.OrderHistory{
		UserID:       userID,
		OrderItemID:  orderItemID,
		Descriptions: descriptions,
		CreatedAt:    time.Now(),
		User:         userData,
		OrderItem:    orderItemData,
	}
	return uc.orderHistoryRepo.Create(ctx, orderHistory)
}

func (uc *orderHistoryUseCase) Update(ctx context.Context, id int, userID int, orderItemID int, descriptions string) error {
	orderHistory, err := uc.orderHistoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if orderHistory == nil {
		return errors.New("order history not found")
	}
	orderHistory.UserID = userID
	orderHistory.OrderItemID = orderItemID
	orderHistory.Descriptions = descriptions
	return uc.orderHistoryRepo.Update(ctx, orderHistory)
}

func (uc *orderHistoryUseCase) GetByID(ctx context.Context, id int) (*entity.OrderHistory, error) {
	return uc.orderHistoryRepo.GetByID(ctx, id)
}

func (uc *orderHistoryUseCase) GetByUserID(ctx context.Context, userID, limit, offset int) ([]*entity.OrderHistory, error) {
	return uc.orderHistoryRepo.GetByUserID(ctx, userID, limit, offset)
}

func (uc *orderHistoryUseCase) GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.OrderHistory, error) {
	return uc.orderHistoryRepo.GetAllPagination(ctx, limit, offset)
}

func (uc *orderHistoryUseCase) CountData(ctx context.Context, userId int) int64 {
	return uc.orderHistoryRepo.CountData(ctx, userId)
}
