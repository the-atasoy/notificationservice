package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationType string
const (
	EmailNotification     NotificationType = "email"
	WebSocketNotification NotificationType = "websocket"
	RestNotification      NotificationType = "rest"
)

type EmailRecipients struct {
	To  []string `bson:"to" json:"to"`
	Cc  []string `bson:"cc,omitempty" json:"cc,omitempty"`
	Bcc []string `bson:"bcc,omitempty" json:"bcc,omitempty"`
}

type EmailContent struct {
	Subject     string          `bson:"subject" json:"subject"`
	Body        string          `bson:"body" json:"body"`
	Recipients  EmailRecipients `bson:"recipients" json:"recipients"`
	Attachments []string        `bson:"attachments,omitempty" json:"attachments,omitempty"`
}

type Notification struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     string             `bson:"userId" json:"userId"`
	Message    string             `bson:"message" json:"message"`
	Type       NotificationType   `bson:"type" json:"type"`
	IsReceived bool               `bson:"isReceived" json:"isReceived"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	ReceivedAt *time.Time         `bson:"receivedAt,omitempty" json:"receivedAt,omitempty"`
	EmailData  *EmailContent      `bson:"emailData,omitempty" json:"emailData,omitempty"`
}
