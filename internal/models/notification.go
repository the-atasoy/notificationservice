package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationType string

const (
    EmailNotification     NotificationType = "Mail"
    InAppNotification NotificationType = "InApp"
)

type NotificationStatus string

const (
    Pending NotificationStatus = "Pending"
    Sent    NotificationStatus = "Sent"
    Failed  NotificationStatus = "Failed"
)

type DeliveryStatus struct {
    NotificationStatus     NotificationStatus `bson:"status" json:"status"`
    UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
    Error      string             `bson:"error,omitempty" json:"error,omitempty"`
}

type MailDetails struct {
	To      string `bson:"to" json:"to"`
	CC      []string `bson:"cc,omitempty" json:"cc,omitempty"`
	BCC     []string `bson:"bcc,omitempty" json:"bcc,omitempty"`
}

type Notification struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID         uuid.UUID           `bson:"userId" json:"userId"`
	MessageID      string              `bson:"messageId" json:"messageId"`
	Subject        string              `bson:"subject" json:"subject"`
	Body           string              `bson:"body" json:"body"`
	Type           NotificationType    `bson:"type" json:"type"`
	DeliveryStatus DeliveryStatus      `bson:"status" json:"status"`
	MailInfo       *MailDetails        `bson:"mailInfo,omitempty" json:"mailInfo,omitempty"`
	CreatedAt      time.Time           `bson:"createdAt" json:"createdAt"`
	ReceivedAt     *time.Time          `bson:"receivedAt,omitempty" json:"receivedAt,omitempty"`
}