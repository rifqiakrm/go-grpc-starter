package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"grpc-starter/common/cache"
	"grpc-starter/modules/user/v1/entity"
)

// UserFinderRepository defines dependencies for UserFinder
type UserFinderRepository struct {
	db    *gorm.DB
	cache cache.Cacheable
}

// NewUserFinderRepository creates a new UserFinder repository
func NewUserFinderRepository(
	db *gorm.DB,
	cache cache.Cacheable,
) *UserFinderRepository {
	return &UserFinderRepository{
		db:    db,
		cache: cache,
	}
}

// UserFinderRepositoryUseCase is use case for finding in user table
type UserFinderRepositoryUseCase interface {
	// FindByID finds user
	FindByID(ctx context.Context, refID uuid.UUID) (*entity.User, error)
	// FindByEmail finds user by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

// FindByID finds user
func (r *UserFinderRepository) FindByID(ctx context.Context, refID uuid.UUID) (*entity.User, error) {
	var result *entity.User
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", refID).First(&result).Error; err != nil {
		return nil, errors.Wrap(err, "[UserFinderRepository - FindByID] Error while finding user data")
	}

	return result, nil
}

// FindByEmail finds user by email
func (r *UserFinderRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var result *entity.User
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).First(&result).Error; err != nil {
		return nil, errors.Wrap(err, "[UserFinderRepository - FindByEmail] Error while finding user data")
	}

	return result, nil
}
