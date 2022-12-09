package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.opencensus.io/trace"

	"grpc-starter/common/config"
	"grpc-starter/modules/notification/v1/entity"
)

// EmailSenderUsecase is use case for creating new ptk
type EmailSenderUsecase interface {
	// SendWithMailgunAPI send email using mailgun api
	SendWithMailgunAPI(ctx context.Context, mID, from, to, subject, message, category, creator string, pubsubMessage *pubsub.Message) error
	// SendWithSendgridAPI send email using sendgrid api
	SendWithSendgridAPI(ctx context.Context, mID, senderName, senderEmail, to, subject, message, category, creator string, pubsubMessage *pubsub.Message) error
}

// EmailSentRepository is use case for creating new ptk
type EmailSentRepository interface {
	// Insert insert email sent log to database
	Insert(ctx context.Context, ent *entity.EmailSent) error
	// UpdateStatus update status email sent
	UpdateStatus(ctx context.Context, emailSent *entity.EmailSent) error
}

// EmailSender is use case for creating new ptk
type EmailSender struct {
	emailSentRepo  EmailSentRepository
	mailgunConfig  config.Mailgun
	sendgridConfig config.Sendgrid
}

// NewEmailSender is constructor for EmailSender
func NewEmailSender(repository EmailSentRepository, mailgunConfig config.Mailgun, sendgridConfig config.Sendgrid) *EmailSender {
	return &EmailSender{
		emailSentRepo:  repository,
		mailgunConfig:  mailgunConfig,
		sendgridConfig: sendgridConfig,
	}
}

// SendWithMailgunAPI send email using mailgun api and save to database
func (s *EmailSender) SendWithMailgunAPI(ctx context.Context, mID, from, to, subject, message, category, creator string, pubsubMessage *pubsub.Message) error {
	ctxSpan, span := trace.StartSpan(ctx, "Notification-EmailSenderService-SendWithMailgunAPI")
	defer span.End()

	var status string
	if len(strings.TrimSpace(to)) == 0 {
		status = entity.EmailSentStatusNoRecipient
	} else {
		status = entity.EmailSentStatusOutgoing
	}

	// save sent message to repository
	emailSent := entity.NewEmailSent(mID, from, to, subject, message,
		status, category, creator)

	err := s.emailSentRepo.Insert(ctxSpan, emailSent)

	if err != nil {
		failedEmailSent := &entity.EmailSent{
			MId:      mID,
			Status:   entity.EmailSentStatusFailed,
			Category: category,
		}
		_ = s.emailSentRepo.UpdateStatus(ctxSpan, failedEmailSent)

		pubsubMessage.Nack()
		log.Print()
	}

	if status == entity.EmailSentStatusNoRecipient {
		pubsubMessage.Ack()
		return errors.New("no recipient")
	}

	mg := mailgun.NewMailgun(s.mailgunConfig.Domain, s.mailgunConfig.APIKey)

	m := mg.NewMessage(from, subject, "", to)
	m.SetHtml(message)

	_, _, err = mg.Send(ctxSpan, m)

	if err != nil {
		failedEmailSent := &entity.EmailSent{
			MId:         mID,
			Status:      entity.EmailSentStatusFailed,
			StatusNotes: sql.NullString{String: err.Error(), Valid: true},
			Category:    category,
		}
		_ = s.emailSentRepo.UpdateStatus(ctxSpan, failedEmailSent)

		pubsubMessage.Nack()
		return err
	}

	pubsubMessage.Ack()

	updateEmailSent := &entity.EmailSent{
		MId:      mID,
		Status:   entity.EmailSentStatusSuccess,
		Category: category,
	}
	return s.emailSentRepo.UpdateStatus(ctxSpan, updateEmailSent)
}

// SendWithSendgridAPI send email using mailgun api and save to database
func (s *EmailSender) SendWithSendgridAPI(ctx context.Context, mID, senderName, senderEmail, to, subject, message, category, creator string, pubsubMessage *pubsub.Message) error {
	ctxSpan, span := trace.StartSpan(ctx, "Notification-EmailSenderService-SendWithSendgridAPI")
	defer span.End()

	var status string
	if len(strings.TrimSpace(to)) == 0 {
		status = entity.EmailSentStatusNoRecipient
	} else {
		status = entity.EmailSentStatusOutgoing
	}

	// save sent message to repository
	from := fmt.Sprintf("%s <%s>", senderName, senderEmail)
	emailSent := entity.NewEmailSent(mID, from, to, subject, message,
		status, category, creator)

	err := s.emailSentRepo.Insert(ctxSpan, emailSent)

	if err != nil {
		failedEmailSent := &entity.EmailSent{
			MId:      mID,
			Status:   entity.EmailSentStatusFailed,
			Category: category,
		}
		_ = s.emailSentRepo.UpdateStatus(ctxSpan, failedEmailSent)

		pubsubMessage.Nack()
		log.Print("failed to send email", err)
	}

	if status == entity.EmailSentStatusNoRecipient {
		pubsubMessage.Ack()
		return errors.New("no recipient")
	}

	sender := mail.NewEmail(senderName, senderEmail)
	receiver := mail.NewEmail(to, to)
	content := mail.NewContent("text/html", message)
	m := mail.NewV3MailInit(sender, subject, receiver, content)

	request := sendgrid.GetRequest(s.sendgridConfig.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err = sendgrid.API(request)

	if err != nil {
		failedEmailSent := &entity.EmailSent{
			MId:         mID,
			Status:      entity.EmailSentStatusFailed,
			StatusNotes: sql.NullString{String: err.Error(), Valid: true},
			Category:    category,
		}
		_ = s.emailSentRepo.UpdateStatus(ctxSpan, failedEmailSent)

		pubsubMessage.Nack()
		log.Print("failed to send email", err)
	}

	pubsubMessage.Ack()

	updateEmailSent := &entity.EmailSent{
		MId:      mID,
		Status:   entity.EmailSentStatusSuccess,
		Category: category,
	}
	return s.emailSentRepo.UpdateStatus(ctxSpan, updateEmailSent)
}
