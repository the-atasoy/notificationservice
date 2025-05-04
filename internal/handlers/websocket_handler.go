package handlers

import (
	"fmt"
	"notificationservice/internal/models"
)

type WebSocketHandler struct{}

func NewWebSocketHandler() IHandler {
	return &WebSocketHandler{}
}

func (h *WebSocketHandler) Deliver(notification *models.Notification) error {
	// TODO: implement WebSocket push logic
	fmt.Println("Sending WebSocket notification to user", notification.UserID)
	return nil
}