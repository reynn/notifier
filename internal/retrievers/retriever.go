package retrievers

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/reynn/notifier/internal/types"
)

var ErrNotFound = errors.New("notification not found")

type (
	// Retriever represents a notification retriever.
	Retriever interface {
		Store(ctx context.Context, id uuid.UUID, n *types.Notification) error
		ByID(ctx context.Context, id uuid.UUID) (*types.Notification, error)
	}
)
