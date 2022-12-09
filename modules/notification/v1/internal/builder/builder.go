// Package builder is used to build the handler.
package builder

import (
	"gorm.io/gorm"

	"grpc-starter/common/config"
	"grpc-starter/modules/notification/v1/internal/pubsub/handler"
	"grpc-starter/modules/notification/v1/internal/repository"
	"grpc-starter/modules/notification/v1/service"
)

// BuildSendEmailPubSubHandler is used to build the pubsub handler.
func BuildSendEmailPubSubHandler(db *gorm.DB, config config.Config) *handler.SendEmailPubSubHandler {
	repo := repository.NewEmailSent(db)
	svc := service.NewEmailSender(repo, config.Mailgun, config.Sendgrid)

	return handler.NewSendEmailPubSubHandler(config, svc)
}
