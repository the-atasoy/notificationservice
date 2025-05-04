package handlers

import (
	"fmt"
	"notificationservice/internal/models"
)

type RestHandler struct{}

func NewRestHandler() IHandler {
	return &RestHandler{}
}

func (h *RestHandler) Deliver(notification *models.Notification) error {
	// TODO: Logic for storing or delivering REST fallback notifications
	fmt.Println("Saving REST fallback notification for user", notification.UserID)
	return nil
}