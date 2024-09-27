package store

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/reynn/notifier/internal/types"
)

var ErrNotFound = errors.New("notification not found")

type (
	Writer interface {
		Write(ctx context.Context, id uuid.UUID, n *types.Notification) error
	}
	Reader interface {
		ByID(ctx context.Context, id uuid.UUID) (*types.Notification, error)
	}
	ReaderWriter interface {
		Reader
		Writer
	}
)
