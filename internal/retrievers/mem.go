package retrievers

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/reynn/notifier/internal/types"
)

type (
	InMemoryRetriever struct {
		items map[uuid.UUID]*types.Notification
	}
)

func NewInMemoryRetriever() *InMemoryRetriever {
	uuid1 := uuid.MustParse("0191b114-7386-78d1-8efd-786fb1db4138")
	return &InMemoryRetriever{
		items: map[uuid.UUID]*types.Notification{
			uuid1: {
				ID:         uuid1,
				Recipients: []string{"nic@reynn.dev"},
				Message:    []byte(`hello world`),
				Tags:       []string{"example"},
				Priority:   types.NotificationPriorityHigh,
				Type:       types.NotificationTypeEmail,
				Status:     types.NotificationStatusSubmitted,
				CreatedAt:  time.Now(),
			},
		},
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
