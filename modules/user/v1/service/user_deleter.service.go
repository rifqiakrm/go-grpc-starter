package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	"grpc-starter/common/config"
	commonError "grpc-starter/common/errors"
	"grpc-starter/modules/user/v1/internal/repository"
)

// UserDeleter responsible for deleting user
type UserDeleter struct {
	cfg                   config.Config
	userDeleterRepository repository.UserDeleterRepositoryUseCase
}

// UserDeleterUseCase is use case for deleting existing user
type UserDeleterUseCase interface {
	// Delete delete user by user id
	Delete(ctx context.Context, refID uuid.UUID) error
}

// NewUserDeleter constructs new instance of UserDeleter
func NewUserDeleter(
	cfg config.Config,
	userDeleterRepository repository.UserDeleterRepositoryUseCase,
) *UserDeleter {
	return &UserDeleter{
		cfg:                   cfg,
		userDeleterRepository: userDeleterRepository,
	}
}

// Delete deletes user
func (svc *UserDeleter) Delete(ctx context.Context, refID uuid.UUID) error {
	err := svc.userDeleterRepository.Delete(ctx, refID)

	if err != nil {
		log.Print("[UserDeleter - Delete] Error while deleting user data :", err)
		return commonError.ErrInternalServerError.Error()
	}

	return nil
}
