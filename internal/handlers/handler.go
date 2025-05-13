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
	notification, err := handler.getNotification(data)
	if err != nil {
		return err
	}

	deliveryErr := handler.deliverNotification(notification)
	if deliveryErr != nil {
		return deliveryErr
	}
	
	log.Printf("Notification processed and delivered: ID=%v, Type=%s, User=%s",
		notification.ID, notification.Type, notification.UserID)
	return nil
}

func (handler *Handler) getNotification(data []byte) (*models.Notification, error) {
	notification, err := handler.unmarshalMessage(data)
	if err != nil {
		return nil, err
	}

	existingNotification, err := handler.repo.GetUnsentNotifications(notification.ExternalID)
	if err != nil {
		return nil, errors.NewRetriableError("database query failed", err)
	}

	if existingNotification != nil {
		return existingNotification, nil
	}
	
	if err := handler.repo.SaveNotification(notification); err != nil {
		return nil, errors.NewRetriableError("database operation failed", err)
	}
	
	return notification, nil
}

func (handler *Handler) deliverNotification(notification *models.Notification) error {
	var deliveryErr error
	switch notification.Type {
	case models.EmailNotification:
		deliveryErr = handler.emailHandler.Deliver(notification)
	case models.InAppNotification:
		deliveryErr = handler.websocketHandler.Deliver(notification)
	default:
		deliveryErr = errors.NewValidationError(
			fmt.Sprintf("unknown notification type: %s", notification.Type),
			nil,
		)
	}
	return handler.handleDeliveryStatus(notification, deliveryErr)
}

func (handler *Handler) handleDeliveryStatus(notification *models.Notification, deliveryErr error) error {
	if deliveryErr == nil {
		notification.DeliveryStatus = models.DeliveryStatus{
			NotificationStatus: models.Sent,
		}
		if err := handler.repo.UpdateNotificationStatus(notification.ID, notification.DeliveryStatus); err != nil {
			return errors.NewRetriableError("failed to update notification status", err)
		}
		return nil
	} 
	
	if errors.IsRetriableError(deliveryErr) {
		notification.DeliveryStatus = models.DeliveryStatus{
			NotificationStatus: models.Pending,
			Error:             deliveryErr.Error(),
		}
		if err := handler.repo.UpdateNotificationStatus(notification.ID, notification.DeliveryStatus); err != nil {
			log.Printf("Failed to update retry status: %v", err)
		}
		return deliveryErr
	} 
	
	notification.DeliveryStatus = models.DeliveryStatus{
		NotificationStatus: models.Failed,
		Error:             deliveryErr.Error(),
	}
	if err := handler.repo.UpdateNotificationStatus(notification.ID, notification.DeliveryStatus); err != nil {
		log.Printf("Failed to update failed status: %v", err)
	}
	return deliveryErr
}

func (handler *Handler) unmarshalMessage(data []byte) (*models.Notification, error) {
	var message models.NotificationMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return nil, errors.NewValidationError("invalid JSON format", err)
	}
	if message.UserID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.NewValidationError("userID is required", nil)
	}
	if message.Subject == "" {
		return nil, errors.NewValidationError("subject is required", nil)
	}
	if message.Body == "" {
		return nil, errors.NewValidationError("body is required", nil)
	}
	
	// Convert the message to a full notification
	notification := message.ToNotification()
	return notification, nil
}

func (handler *Handler) GetUnreadNotifications(userId string) ([]models.Notification, error) {
	notifications, err := handler.repo.GetUnreadNotifications(userId)
	if err != nil {
		return nil, errors.NewProcessingError("failed to get unread notifications", err)
	}
	return notifications, nil
}