package service

import (
	"context"
	"log"

	"grpc-starter/common/config"
	commonError "grpc-starter/common/errors"
	"grpc-starter/modules/user/v1/entity"
	"grpc-starter/modules/user/v1/internal/repository"
)

// UserUpdater responsible for updating user
type UserUpdater struct {
	cfg                  config.Config
	updateUserRepository repository.UserUpdaterRepositoryUseCase
}

// UserUpdaterUseCase is use case for updating existing user
type UserUpdaterUseCase interface {
	// Update update user by user id
	Update(ctx context.Context, user *entity.User) error
}

// NewUserUpdater constructs new instance of UserUpdater
func NewUserUpdater(
	cfg config.Config,
	updateUserRepository repository.UserUpdaterRepositoryUseCase,
) *UserUpdater {
	return &UserUpdater{
		cfg:                  cfg,
		updateUserRepository: updateUserRepository,
	}
}

// Update updates user
func (svc *UserUpdater) Update(ctx context.Context, user *entity.User) error {
	if err := svc.updateUserRepository.Update(ctx, user); err != nil {
		log.Println("[UserUpdater - Update] Error while updating user data :", err)
		return commonError.ErrInternalServerError.Error()
	}

	return nil
}
