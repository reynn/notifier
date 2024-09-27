package types

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	Notification struct {
		ID         uuid.UUID            `db:"id" json:"id"`
		Recipients []string             `db:"recipients" json:"recipients"`
		Message    []byte               `db:"message" json:"message"`
		Tags       []string             `db:"tags" json:"tags"`
		Priority   NotificationPriority `db:"priority" json:"priority"`
		Type       NotificationType     `db:"type" json:"type"`
		Status     NotificationStatus   `db:"status" json:"status"`
		CreatedAt  time.Time            `db:"created_at" json:"created_at"`
		UpdatedAt  *time.Time           `db:"updated_at" json:"updated_at"`
	}
)

type (
	NotificationPriority string
	NotificationType     string
	NotificationStatus   string
)

const (
	NotificationPriorityHigh    NotificationPriority = "HIGH"
	NotificationPriorityLow     NotificationPriority = "LOW"
	NotificationPriorityDefault NotificationPriority = "DEFAULT"
)

const (
	NotificationTypeEmail   NotificationType = "EMAIL"
	NotificationTypeUnset   NotificationType = "UNSET"
	NotificationTypeWebhook NotificationType = "WEBHOOK"
	NotificationTypeSMS     NotificationType = "SMS"
	NotificationTypePush    NotificationType = "PUSH"
)

const (
	NotificationStatusSubmitted NotificationStatus = "SUBMITTED"
	NotificationStatusCompleted NotificationStatus = "COMPLETED"
	NotificationStatusFailed    NotificationStatus = "FAILED"
)

func ParseNotificationPriority(p string) NotificationPriority {
	switch strings.ToUpper(p) {
	case "HIGH":
		return NotificationPriorityHigh
	case "LOW":
		return NotificationPriorityLow
	default:
		return NotificationPriorityDefault
	}
}

func ParseNotificationType(t string) NotificationType {
	switch strings.ToUpper(t) {
	case "EMAIL":
		return NotificationTypeEmail
	case "SMS":
		return NotificationTypeSMS
	case "PUSH":
		return NotificationTypePush
	case "WEBHOOK":
		return NotificationTypeWebhook
	default:
		return NotificationTypeUnset
	}
}

func ParseNotificationStatus(s string) NotificationStatus {
	switch strings.ToUpper(s) {
	case "COMPLETED":
		return NotificationStatusCompleted
	case "FAILED":
		return NotificationStatusFailed
	default:
		return NotificationStatusSubmitted
	}
}
