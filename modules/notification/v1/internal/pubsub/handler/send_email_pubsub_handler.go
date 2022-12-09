package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"grpc-starter/common/config"
	"grpc-starter/common/logger"
	"grpc-starter/modules/notification/v1/entity"
	"grpc-starter/modules/notification/v1/service"
)

// SendEmailPubSubHandler struct
type SendEmailPubSubHandler struct {
	cfg            config.Config
	emailSenderSvc service.EmailSenderUsecase
}

const (
	// SendEmailSubName is a subscriber name for SendEmailSub
	SendEmailSubName = "send-email-sub"
)

// NewSendEmailPubSubHandler create ptk pubsub handler
func NewSendEmailPubSubHandler(
	cfg config.Config,
	emailSenderSvc service.EmailSenderUsecase,
) *SendEmailPubSubHandler {
	return &SendEmailPubSubHandler{
		cfg:            cfg,
		emailSenderSvc: emailSenderSvc,
	}
}

// SubscriptionName is a function for getting subscription name
func (pubsub *SendEmailPubSubHandler) SubscriptionName() string {
	return SendEmailSubName
}

// ProcessMessage is a function for processing message from pubsub
func (pubsub *SendEmailPubSubHandler) ProcessMessage(ctx context.Context, m *pubsub.Message) {
	// log message from pubsub
	logger.Info(fmt.Sprintf("Received message: %s", m.Data))

	ctxSpan, span := trace.StartSpan(ctx, "Notification-SendEmailPubSubHandler-ProcessMessage")
	defer span.End()

	var payload entity.EmailPayload

	// parsing json payload
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		log.Print(errors.Wrap(err, fmt.Sprintf("[SendEmailPubSubHandler-ProcessMessage] error unmarshal: %s", m.Attributes)))
		m.Ack()
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeInternal,
			Message: err.Error(),
		})
		return
	}

	// send email
	err := pubsub.emailSenderSvc.SendWithSendgridAPI(
		ctxSpan, m.ID, pubsub.cfg.SMTP.FromName, pubsub.cfg.SMTP.FromEmail,
		payload.To, payload.Subject, payload.Content, payload.Category, pubsub.SubscriptionName(), m)
	if err != nil {
		log.Print(errors.Wrap(err, fmt.Sprintf("[SendEmailPubSubHandler-ProcessMessage] error send email svc: %s", m.Attributes)))
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeInternal,
			Message: err.Error(),
		})
		return
	}

	span.SetStatus(trace.Status{
		Code: trace.StatusCodeOK,
	})
}
