package handlers

import "notificationservice/internal/models"

type IHandler interface {
	Deliver(notification *models.Notification) error
}