package store

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/reynn/notifier/internal/types"
)

type (
	InMemory struct {
		items map[uuid.UUID]*types.Notification
	}
)

func NewInMemory() *InMemory {
	uuid1 := uuid.MustParse("0191b114-7386-78d1-8efd-786fb1db4138")
	return &InMemory{
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

func (i *InMemory) Write(ctx context.Context, n *types.Notification) error {
	n.ID = uuid.New()
	i.items[n.ID] = n
	return nil
}

func (i *InMemory) ByID(ctx context.Context, id uuid.UUID) (*types.Notification, error) {
	if n, ok := i.items[id]; ok {
		return n, nil
	}
	return nil, ErrNotFound
}
