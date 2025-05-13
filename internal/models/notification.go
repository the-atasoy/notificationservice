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

type NotificationMessage struct {
	UserID     uuid.UUID        `json:"userId"`
	ExternalID uuid.UUID        `json:"externalId"`
	Subject    string           `json:"subject"`
	Body       string           `json:"body"`
	Type       NotificationType `json:"type"`
	MailInfo   *MailDetails     `json:"mailInfo,omitempty"`
}

func (msg *NotificationMessage) ToNotification() *Notification {
	now := time.Now()
	return &Notification{
		UserID:     msg.UserID,
		ExternalID: msg.ExternalID,
		Subject:    msg.Subject,
		Body:       msg.Body,
		Type:       msg.Type,
		MailInfo:   msg.MailInfo,
		DeliveryStatus: DeliveryStatus{
			NotificationStatus: Pending,
			UpdatedAt:  now,
		},
		CreatedAt:  now,
	}
}

type DeliveryStatus struct {
    NotificationStatus     NotificationStatus `bson:"notificationStatus" json:"notificationStatus"`
    UpdatedAt  time.Time          `bson:"updatedAt" json:"-"`
    Error      string             `bson:"error,omitempty" json:"error,omitempty"`
}

type MailDetails struct {
	To      string   `bson:"to" json:"to"`
	CC      []string `bson:"cc,omitempty" json:"cc,omitempty"`
	BCC     []string `bson:"bcc,omitempty" json:"bcc,omitempty"`
}

type Notification struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID         uuid.UUID           `bson:"userId" json:"userId"`
	ExternalID     uuid.UUID           `bson:"externalId" json:"externalId"`
	Subject        string              `bson:"subject" json:"subject"`
	Body           string              `bson:"body" json:"body"`
	Type           NotificationType    `bson:"type" json:"type"`
	DeliveryStatus DeliveryStatus      `bson:"deliveryStatus" json:"deliveryStatus"`
	MailInfo       *MailDetails        `bson:"mailInfo,omitempty" json:"mailInfo,omitempty"`
	CreatedAt      time.Time           `bson:"createdAt" json:"-"`
	ReceivedAt     *time.Time          `bson:"receivedAt,omitempty" json:"-"`
}