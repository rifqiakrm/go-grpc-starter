package repository

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"grpc-starter/common/cache"
	"grpc-starter/modules/user/v1/entity"
)

// UserUpdaterRepository defines dependencies for UserUpdater
type UserUpdaterRepository struct {
	db    *gorm.DB
	cache cache.Cacheable
}

// NewUserUpdaterRepository creates a new UserUpdater repository
func NewUserUpdaterRepository(
	db *gorm.DB,
	cache cache.Cacheable,
) *UserUpdaterRepository {
	return &UserUpdaterRepository{
		db:    db,
		cache: cache,
	}
}

// UserUpdaterRepositoryUseCase is use case for UserUpdaterRepository
type UserUpdaterRepositoryUseCase interface {
	// Update updates user
	Update(ctx context.Context, user *entity.User) error
}

// Update updates user
func (r *UserUpdaterRepository) Update(ctx context.Context, user *entity.User) error {
	oldTime := user.UpdatedAt
	user.UpdatedAt = time.Now()
	if err := r.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			sourceModel := new(entity.User)
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sourceModel, user.ID).Error; err != nil {
				log.Println("[UserRepository - Update]", err)
				return err
			}
			if err := tx.Model(&entity.User{}).
				Where(`id`, user.ID).
				UpdateColumns(sourceModel.MapUpdateFrom(user)).Error; err != nil {
				log.Println("[UserRepository - Update]", err)
				return err
			}
			return nil
		}); err != nil {
		user.UpdatedAt = oldTime
	}

	return nil
}
