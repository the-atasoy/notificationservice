package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"notificationservice/internal/models"
	"notificationservice/internal/repository"
)

// NotificationHandler handles the processing of notifications from message queue
type NotificationHandler struct {
    repo *repository.MongoRepository
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(repo *repository.MongoRepository) *NotificationHandler {
    return &NotificationHandler{
        repo: repo,
    }
}

func (handler *NotificationHandler) ProcessMessage(data []byte) error {
    var notification models.Notification
    if err := json.Unmarshal(data, &notification); err != nil {
        return fmt.Errorf("failed to unmarshal notification: %w", err)
    }

    if notification.Status.Status == "" {
        notification.Status.Status = models.NotificationStatus("Pending")
    }

    if err := handler.repo.SaveNotification(&notification); err != nil {
        return fmt.Errorf("failed to save notification: %w", err)
    }

    // Send(&notification)
    log.Printf("Notification processed and saved: ID=%v, Type=%s, User=%s", 
        notification.ID, notification.Type, notification.UserID)

    return nil
}

// func Send(notification *models.Notification) error {}

func(handler *NotificationHandler) GetUnreadNotifications(userId string) ([]models.Notification, error) {
    notifications, err := handler.repo.GetUnreadNotifications(userId)
    if err != nil {
        return nil, fmt.Errorf("failed to get unread notifications: %w", err)
    }
    return notifications, nil
}