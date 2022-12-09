package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"grpc-starter/common/cache"
	"grpc-starter/modules/user/v1/entity"
)

// UserDeleterRepository defines dependencies for UserDeleter
type UserDeleterRepository struct {
	db    *gorm.DB
	cache cache.Cacheable
}

// NewUserDeleterRepository creates a new UserDeleter repository
func NewUserDeleterRepository(
	db *gorm.DB,
	cache cache.Cacheable,
) *UserDeleterRepository {
	return &UserDeleterRepository{
		db:    db,
		cache: cache,
	}
}

// UserDeleterRepositoryUseCase is use case for deleting in user table
type UserDeleterRepositoryUseCase interface {
	// Delete deletes user
	Delete(ctx context.Context, refID uuid.UUID) error
}

// Delete deletes user
func (r *UserDeleterRepository) Delete(ctx context.Context, refID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Delete(&entity.User{}, "ref_id = ?", refID).Error; err != nil {
		return errors.Wrap(err, "[UserDeleterRepository - Delete] Error while deleting user data")
	}

	return nil
}
