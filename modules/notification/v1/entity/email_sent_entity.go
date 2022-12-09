package entity

import (
	"database/sql"

	"grpc-starter/common/entity"
)

const (
	// EmailSentStatusNoRecipient is a constant for no recipient email status
	EmailSentStatusNoRecipient = "NO_RECIPIENT"
	// EmailSentStatusOutgoing is a constant for outgoing email status
	EmailSentStatusOutgoing = "OUTGOING"
	// EmailSentStatusSuccess is a constant for success email status
	EmailSentStatusSuccess = "SUCCESS"
	// EmailSentStatusFailed is a constant for failed email status
	EmailSentStatusFailed = "FAILED"
)

// EmailSent represents table on db
type EmailSent struct {
	ID          int
	MId         string
	From        string
	To          string
	Subject     string
	Content     string
	Status      string
	StatusNotes sql.NullString
	Category    string
	entity.Auditable
}

const (
	// EmailSentTableName represents table name on db
	EmailSentTableName = "notification.email_sent"
)

// NewEmailSent is a constructor for EmailSent
func NewEmailSent(mID, from, to, subject, content, status, category, creator string) *EmailSent {
	return &EmailSent{
		MId:       mID,
		From:      from,
		To:        to,
		Subject:   subject,
		Content:   content,
		Status:    status,
		Category:  category,
		Auditable: entity.NewAuditable(creator),
	}
}

// TableName represents table name on db, need to define it because the db has multi schema
func (e *EmailSent) TableName() string {
	return EmailSentTableName
}
