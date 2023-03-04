package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"

	"test-crud-user-orders/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.User, error)
	SoftDelete(ctx context.Context, id int) error
	CountData(ctx context.Context) int64
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	execDB := r.db.WithContext(ctx).Create(&user)
	if execDB.Error != nil {
		return user, execDB.Error
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Model(user).Updates(&user).Error
}

func (r *userRepository) SoftDelete(ctx context.Context, id int) error {
	user := &entity.User{ID: id}

	err := r.db.Delete(user).Error
	if err != nil {
		return fmt.Errorf("error soft-deleting user with ID %d: %s", id, err.Error())
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.WithContext(ctx).First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAllPagination(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User

	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err.Error())
	}

	return users, nil
}

func (r *userRepository) CountData(ctx context.Context) int64 {
	var count int64
	var users []*entity.User

	r.db.Find(&users).Count(&count)
	return count
}
