package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"notificationservice/internal/errors"
	"notificationservice/internal/models"
	"notificationservice/internal/repository"
)

type Handler struct {
	repo             *repository.MongoRepository
	emailHandler     IHandler
	websocketHandler IHandler
}

func NewHandler(repo *repository.MongoRepository) *Handler {
	return &Handler{
		repo:             repo,
		emailHandler:     NewEmailHandler(),
		websocketHandler: NewWebSocketHandler(),
	}
}

func (handler *Handler) ProcessMessage(data []byte) error {
	notification, err := handler.validateMessage(data)
	if err != nil {
		return err
	}

	if err := handler.repo.SaveNotification(notification); err != nil {
		return errors.NewRetriableError("database operation failed", err)
	}

	switch notification.Type {
	case models.EmailNotification:
		if err := handler.emailHandler.Deliver(notification); err != nil {
			return err
		}
	case models.InAppNotification:
		if err := handler.websocketHandler.Deliver(notification); err != nil {
			return err
		}
	default:
		return errors.NewValidationError(
			fmt.Sprintf("unknown notification type: %s", notification.Type),
			nil,
		)
	}

	log.Printf("Notification processed and saved: ID=%v, Type=%s, User=%s",
		notification.ID, notification.Type, notification.UserID)
	return nil
}

func (handler *Handler) validateMessage(data []byte) (*models.Notification, error) {
	var notification models.Notification
	if err := json.Unmarshal(data, &notification); err != nil {
		return nil, errors.NewValidationError("invalid JSON format", err)
	}
	if notification.UserID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.NewValidationError("userID is required", nil)
	}
	if notification.Subject == "" {
		return nil, errors.NewValidationError("subject is required", nil)
	}
	if notification.Body == "" {
		return nil, errors.NewValidationError("body is required", nil)
	}
	return &notification, nil
}

func (handler *Handler) GetUnreadNotifications(userId string) ([]models.Notification, error) {
	notifications, err := handler.repo.GetUnreadNotifications(userId)
	if err != nil {
		return nil, errors.NewProcessingError("failed to get unread notifications", err)
	}
	return notifications, nil
}