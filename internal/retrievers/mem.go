package retrievers

import (
	"context"

	"github.com/google/uuid"

	"github.com/reynn/notifier/internal/types"
)

type (
	InMemoryRetriever struct {
		items map[uuid.UUID]*types.Notification
	}
)

func NewInMemoryRetriever() *InMemoryRetriever {
	return &InMemoryRetriever{
		items: make(map[uuid.UUID]*types.Notification),
	}
}

func (i *InMemoryRetriever) Store(ctx context.Context, id uuid.UUID, n *types.Notification) error {
	i.items[id] = n
	return nil
}

func (i *InMemoryRetriever) ByID(ctx context.Context, id uuid.UUID) (*types.Notification, error) {
	if n, ok := i.items[id]; ok {
		return n, nil
	}
	return nil, ErrNotFound
}
