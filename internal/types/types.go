package types

import (
	"time"

	notifierv1 "github.com/reynn/notifier/gen/proto/notifier/v1"
)

type (
	Notification struct {
		ID         string
		Recipients []string
		Message    []byte
		Tags       []string
		Priority   notifierv1.NotificationPriority
		Type       notifierv1.NotificationType
		Status     notifierv1.NotificationStatus
		CreatedAt  time.Time
		UpdatedAt  *time.Time
	}
)
