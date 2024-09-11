package types

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	Notification struct {
		ID         uuid.UUID
		Recipients []string
		Message    []byte
		Tags       []string
		Priority   NotificationPriority
		Type       NotificationType
		Status     NotificationStatus
		CreatedAt  time.Time
		UpdatedAt  *time.Time
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
