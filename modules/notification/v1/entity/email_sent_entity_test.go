package entity_test

import (
	"testing"

	"grpc-starter/common/tools"
	"grpc-starter/modules/notification/v1/entity"
)

func TestNewEmailSentEntity(t *testing.T) {
	t.Log("TestNewEmailSentEntity")

	mID := "ID"
	from := "from"
	to := "to"
	subject := "subject"
	body := "body"
	creator := "creator"
	status := "status"
	category := "category"
	emailSentEntity := entity.NewEmailSent(mID, from, to, subject, body, status, category, creator)
	if emailSentEntity == nil {
		t.Error("emailSentEntity is nil")
	} else {
		if emailSentEntity.MId != mID {
			t.Error("emailSentEntity.Id is not equal to mId")
		}
		if emailSentEntity.From != from {
			t.Error("emailSentEntity.From != from")
		}
		if emailSentEntity.To != to {
			t.Error("emailSentEntity.To != to")
		}
		if emailSentEntity.Subject != subject {
			t.Error("emailSentEntity.Subject != subject")
		}
		if emailSentEntity.Content != body {
			t.Error("emailSentEntity.Body != body")
		}
		if emailSentEntity.CreatedBy != tools.StringToNullString(creator) {
			t.Error("emailSentEntity.Creator != creator")
		}
		if emailSentEntity.Status != status {
			t.Error("emailSentEntity.Status != status")
		}
		if emailSentEntity.Category != category {
			t.Error("emailSentEntity.Category != category")
		}
		if emailSentEntity.UpdatedBy != tools.StringToNullString(creator) {
			t.Error("emailSentEntity.UpdatedBy != creator")
		}
	}
}
