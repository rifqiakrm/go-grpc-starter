package service

import (
	"context"

	"go.opencensus.io/trace"

	"grpc-starter/modules/notification/v1/entity"
)

// EmailUpdaterUsecase is use case for creating new ptk
type EmailUpdaterUsecase interface {
	// UpdateStatus update email status
	UpdateStatus(ctx context.Context, emailSent *entity.EmailSent) error
}

// EmailUpdaterRepository is use case for creating new ptk
type EmailUpdaterRepository interface {
	UpdateStatus(ctx context.Context, emailSent *entity.EmailSent) error
}

// EmailUpdater is use case for creating new ptk
type EmailUpdater struct {
	repo EmailUpdaterRepository
}

// NewEmailUpdater is constructor for EmailUpdater
func NewEmailUpdater(repository EmailUpdaterRepository) *EmailUpdater {
	return &EmailUpdater{
		repo: repository,
	}
}

// UpdateStatus update email status
func (s *EmailUpdater) UpdateStatus(ctx context.Context, emailSent *entity.EmailSent) error {
	ctxSpan, span := trace.StartSpan(context.Background(), "Notification-EmailUpdaterService-UpdateStatus")
	defer span.End()

	err := s.repo.UpdateStatus(ctxSpan, emailSent)
	if err != nil {
		return err
	}

	return nil
}
