package usecase

import (
	"context"
	"errors"
	"test-crud-user-orders/internal/entity"
	"test-crud-user-orders/internal/repository"
)

type UserUseCase interface {
	Create(ctx context.Context, fullName string) (*entity.User, error)
	Update(ctx context.Context, id int, fullName string) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.User, error)
	CountData(ctx context.Context) int64
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo}
}

func (uc *userUseCase) Create(ctx context.Context, fullName string) (*entity.User, error) {
	user := &entity.User{
		FullName: fullName,
	}
	return uc.userRepo.Create(ctx, user)
}

func (uc *userUseCase) Update(ctx context.Context, id int, fullName string) error {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	user.FullName = fullName
	return uc.userRepo.Update(ctx, user)
}

func (uc *userUseCase) Delete(ctx context.Context, id int) error {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return uc.userRepo.SoftDelete(ctx, id)
}

func (uc *userUseCase) GetByID(ctx context.Context, id int) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *userUseCase) GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return uc.userRepo.GetAllPagination(ctx, limit, offset)
}

func (uc *userUseCase) CountData(ctx context.Context) int64 {
	return uc.userRepo.CountData(ctx)
}
