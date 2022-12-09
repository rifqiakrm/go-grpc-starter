package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"

	"grpc-starter/common/tools"
)

// Auditable define entity for auditable
type Auditable struct {
	CreatedBy sql.NullString `json:"created_by"`
	UpdatedBy sql.NullString `json:"updated_by"`
	DeletedBy sql.NullString `json:"deleted_by"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not_null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not_null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// NewAuditable create new auditable
func NewAuditable(createdBy string) Auditable {
	return Auditable{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: tools.StringToNullString(createdBy),
		UpdatedBy: tools.StringToNullString(createdBy),
	}
}
