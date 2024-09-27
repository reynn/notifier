package notifiers

import (
	"context"

	"github.com/google/uuid"

	"github.com/reynn/notifier/internal/types"
)

type (
	// Sender represents a notification sender.
	Sender interface {
		Send(context.Context, *types.Notification) (uuid.UUID, error)
	}
)
