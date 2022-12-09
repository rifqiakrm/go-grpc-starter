package entity_test

import (
	"testing"

	"grpc-starter/modules/notification/v1/entity"
)

func TestNewEmailPayloadEntity(t *testing.T) {
	t.Log("TestNewEmailPayloadEntity")

	to := "to"
	subject := "subject"
	body := "body"
	category := "category"
	e := entity.NewEmailPayload(to, subject, body, category)
	if e == nil {
		t.Error("NewEmailPayloadEntity() returned nil")
	} else {
		if e.To != to {
			t.Error("NewEmailPayloadEntity() returned incorrect To")
		}
		if e.Subject != subject {
			t.Error("NewEmailPayloadEntity() returned incorrect Subject")
		}
		if e.Content != body {
			t.Error("NewEmailPayloadEntity() returned incorrect Body")
		}
		if e.Category != category {
			t.Error("NewEmailPayloadEntity() returned incorrect Category")
		}
	}
}
