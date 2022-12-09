package entity

import (
	"database/sql"

	"github.com/google/uuid"

	"grpc-starter/common/entity"
	"grpc-starter/common/tools"
)

// SMSSent represents table on db
type SMSSent struct {
	ID          int
	MId         string
	ClientMId   uuid.UUID
	To          string
	Content     string
	Status      string
	StatusNotes sql.NullString
	Category    sql.NullString
	entity.Auditable
}

const (
	// SMSSentTableName represents table name on db
	SMSSentTableName = "notification.sms_sent"
)

// NewSMSSent is a constructor for SMSSent
func NewSMSSent(mID string, clientMessageID uuid.UUID, to, content, status, statusNotes, category, creator string) *SMSSent {
	var statusNotesConv, categoryConv sql.NullString

	if statusNotes != "" {
		statusNotesConv = tools.StringToNullString(statusNotes)
	}

	if category != "" {
		categoryConv = tools.StringToNullString(category)
	}

	return &SMSSent{
		MId:         mID,
		ClientMId:   clientMessageID,
		To:          to,
		Content:     content,
		Status:      status,
		StatusNotes: statusNotesConv,
		Category:    categoryConv,
		Auditable:   entity.NewAuditable(creator),
	}
}

// TableName represents table name on db, need to define it because the db has multi schema
func (e *SMSSent) TableName() string {
	return SMSSentTableName
}
