package v1

import (
	"context"

	"gorm.io/gorm"

	"grpc-starter/common/config"
	"grpc-starter/modules/notification/v1/internal/builder"
	pubsubSDK "grpc-starter/sdk/pubsub"
)

// InitSendEmailSubscription initialize subscription for sending email
func InitSendEmailSubscription(ctx context.Context, db *gorm.DB, config config.Config) pubsubSDK.Subscriber {
	sendEmailHandler := builder.BuildSendEmailPubSubHandler(db, config)
	return sendEmailHandler
}
