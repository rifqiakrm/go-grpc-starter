package repository

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"grpc-starter/common/cache"
	"grpc-starter/modules/user/v1/entity"
)

// UserCreatorRepository defines dependencies for UserCreator
type UserCreatorRepository struct {
	db    *gorm.DB
	cache cache.Cacheable
}

// NewUserCreatorRepository creates a new UserCreator repository
func NewUserCreatorRepository(
	db *gorm.DB,
	cache cache.Cacheable,
) *UserCreatorRepository {
	return &UserCreatorRepository{
		db:    db,
		cache: cache,
	}
}

// UserCreatorRepositoryUseCase is use case for creating in user table
type UserCreatorRepositoryUseCase interface {
	// Create creates user
	Create(ctx context.Context, user *entity.User) error
}

// Create creates user
func (r *UserCreatorRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.Wrap(err, "[UserCreatorRepository - Create] Error while creating user data")
	}

	return nil
}
