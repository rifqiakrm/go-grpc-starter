// Package entity is an entity definer
package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	commonentity "grpc-starter/common/entity"
	"grpc-starter/common/tools"
)

const (
	// UserTableName represents table name on db
	UserTableName = "users.users"
)

// User defines table for user
type User struct {
	ID          uuid.UUID      `json:"id"`
	Username    sql.NullString `json:"username"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	PhoneNumber sql.NullString `json:"phone_number"`
	commonentity.Auditable
}

// NewUser creates new NewUser
func NewUser(
	id uuid.UUID,
	username string,
	email string,
	password string,
	phoneNumber string,
	createdBy string,
) *User {
	encrypted, _ := tools.BcryptEncrypt(password)

	return &User{
		ID:          id,
		Username:    tools.StringToNullString(username),
		Email:       email,
		Password:    encrypted,
		PhoneNumber: tools.StringToNullString(phoneNumber),
		Auditable:   commonentity.NewAuditable(createdBy),
	}
}

// MapUpdateFrom mapping from model
func (u *User) MapUpdateFrom(from *User) *map[string]interface{} {
	if from == nil {
		return &map[string]interface{}{
			"username":     u.Username,
			"password":     u.Password,
			"phone_number": u.PhoneNumber,
		}
	}

	mapped := make(map[string]interface{})
	if u.Username != from.Username {
		mapped["username"] = from.Username
	}
	if u.Password != from.Password {
		mapped["password"] = from.Password
	}
	if u.PhoneNumber != from.PhoneNumber {
		mapped["phone_number"] = from.PhoneNumber
	}

	mapped["updated_at"] = time.Now()
	return &mapped
}

// TableName represents table name on db, need to define it because the db has multi schema
func (u *User) TableName() string {
	return UserTableName
}
