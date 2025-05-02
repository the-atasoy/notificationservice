package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/mail"

	"notificationservice/internal/errors"
	"notificationservice/internal/models"
	"notificationservice/internal/repository"
)

type NotificationHandler struct {
    repo *repository.MongoRepository
}

func NewNotificationHandler(repo *repository.MongoRepository) *NotificationHandler {
    return &NotificationHandler{
        repo: repo,
    }
}

func (handler *NotificationHandler) ProcessMessage(data []byte) error {
    notification, err := handler.validateMessage(data)
    if err != nil {
        return err
    }

    if err := handler.repo.SaveNotification(notification); err != nil {
        return errors.NewRetriableError("database operation failed", err)
    }

    switch notification.Type {
    case models.EmailNotification:
        if err := handler.processEmailNotification(notification); err != nil {
            return err
        }
    case models.WebSocketNotification:
        if err := handler.processWebSocketNotification(notification); err != nil {
            return err
        }
    case models.RestNotification:
        if err := handler.processRestNotification(notification); err != nil {
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

func (handler *NotificationHandler) processEmailNotification(notification *models.Notification) error {
    if err := validateEmailNotification(notification); err != nil {
        return err
    }
    // To-Do: Implement email sending logic
    return nil
}

func (handler *NotificationHandler) processWebSocketNotification(notification *models.Notification) error {
    // To-Do: Implement WebSocket notification logic
    return nil
}

func (handler *NotificationHandler) processRestNotification(notification *models.Notification) error {
    //To-Do: Implement REST notification logic
    return nil
}

func (handler *NotificationHandler) validateMessage(data []byte) (*models.Notification, error) {
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

func validateEmailNotification(notification *models.Notification) error {
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

func(handler *NotificationHandler) GetUnreadNotifications(userId string) ([]models.Notification, error) {
    notifications, err := handler.repo.GetUnreadNotifications(userId)
    if err != nil {
        return nil, errors.NewProcessingError("failed to get unread notifications", err)
    }
    return notifications, nil
}