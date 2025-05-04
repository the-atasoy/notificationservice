package handlers

import (
	"fmt"
	"net/mail"

	"notificationservice/internal/errors"
	"notificationservice/internal/models"
)

type EmailHandler struct{}

func NewEmailHandler() IHandler {
	return &EmailHandler{}
}

func (h *EmailHandler) Deliver(notification *models.Notification) error {
	if err := h.validate(notification); err != nil {
		return err
	}
	// TODO: implement email sending logic
	fmt.Println("Sending email notification to", notification.MailInfo.To)
	return nil
}

func (h *EmailHandler) validate(notification *models.Notification) error {
	if notification.MailInfo == nil {
		return errors.NewValidationError("mailInfo is required for email notifications", nil)
	}
	if notification.MailInfo.To == "" {
		return errors.NewValidationError("recipient email is required", nil)
	}
	if _, err := mail.ParseAddress(notification.MailInfo.To); err != nil {
		return errors.NewValidationError("invalid email address format", err)
	}
	for _, cc := range notification.MailInfo.CC {
		if _, err := mail.ParseAddress(cc); err != nil {
			return errors.NewValidationError("invalid CC email address format", err)
		}
	}
	for _, bcc := range notification.MailInfo.BCC {
		if _, err := mail.ParseAddress(bcc); err != nil {
			return errors.NewValidationError("invalid BCC email address format", err)
		}
	}
	return nil
}