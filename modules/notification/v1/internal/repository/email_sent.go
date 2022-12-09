package repository

import (
	"context"

	"go.opencensus.io/trace"
	"gorm.io/gorm"

	"grpc-starter/modules/notification/v1/entity"
)

// EmailSent struct
type EmailSent struct {
	gormDB *gorm.DB
}

// NewEmailSent will create new email sent repository
func NewEmailSent(db *gorm.DB) *EmailSent {
	return &EmailSent{db}
}

// Insert will insert notification email sent to database
func (r *EmailSent) Insert(ctx context.Context, emailSent *entity.EmailSent) error {
	ctxSpan, span := trace.StartSpan(ctx, "Notification-EmailSentRepository-Insert")
	defer span.End()

	// check if exist
	exist := &entity.EmailSent{}
	if err := r.gormDB.
		WithContext(ctxSpan).
		Where("m_id = ?", emailSent.MId).
		First(exist).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// insert into database
			if err := r.gormDB.
				WithContext(ctxSpan).
				Model(&entity.EmailSent{}).
				Create(emailSent).
				Error; err != nil {
				return err
			}
		} else {
			// update status
			if err := r.gormDB.
				WithContext(ctxSpan).
				Model(&entity.EmailSent{}).
				Where("m_id = ?", emailSent.MId).
				Update("status", "OUTGOING").
				Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateStatus will update status of email sent
func (r *EmailSent) UpdateStatus(ctx context.Context, emailSent *entity.EmailSent) error {
	ctxSpan, span := trace.StartSpan(ctx, "Notification-EmailSentRepository-UpdateStatus")
	defer span.End()

	return r.gormDB.
		WithContext(ctxSpan).
		Model(&entity.EmailSent{}).
		Where("m_id = ?", emailSent.MId).
		Update("status", emailSent.Status).
		Error
}
