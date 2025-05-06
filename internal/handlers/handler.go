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

func (handler *Handler) ProcessMessage(data []byte, messageID string) error {
	notification, err := handler.validateMessage(data)
	if err != nil {
		return err
	}

	existingNotification, err := handler.repo.GetNotificationByMessageID(messageID)
  if err != nil{
  	return errors.NewRetriableError("database query failed", err)
  }

	if existingNotification == nil {
		if err := handler.repo.SaveNotification(notification); err != nil {
				return errors.NewRetriableError("database operation failed", err)
		}
	} else {
			notification = existingNotification
	}

	var deliveryErr error
    switch notification.Type {
    case models.EmailNotification:
        deliveryErr = handler.emailHandler.Deliver(notification)
    case models.InAppNotification:
        deliveryErr = handler.websocketHandler.Deliver(notification)
    default:
        return errors.NewValidationError(
            fmt.Sprintf("unknown notification type: %s", notification.Type),
            nil,
        )
    }

		if deliveryErr == nil {
				notification.DeliveryStatus = models.DeliveryStatus{
						NotificationStatus:    models.Sent,
				}

				if err := handler.repo.UpdateNotificationStatus(notification.ID, notification.DeliveryStatus); err != 	nil {
						return errors.NewRetriableError("failed to update notification status", err)
				}
		} else if errors.IsRetriableError(deliveryErr) {
				notification.DeliveryStatus = models.DeliveryStatus{
						NotificationStatus:    models.Pending,
						Error:     deliveryErr.Error(),
				}

				if err := handler.repo.UpdateNotificationStatus(notification.ID, notification.DeliveryStatus); err != 	nil {
						log.Printf("Failed to update retry status: %v", err)
				}
				return deliveryErr
		} else {
				notification.DeliveryStatus = models.DeliveryStatus{
						NotificationStatus:    models.Failed,
						Error:     deliveryErr.Error(),
				}

				if err := handler.repo.UpdateNotificationStatus(notification.ID, notification.DeliveryStatus); err != 	nil {
						log.Printf("Failed to update failed status: %v", err)
				}
				return deliveryErr
		}

		log.Printf("Notification processed and delivered: ID=%v, Type=%s, User=%s",
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