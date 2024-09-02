package notifiers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/reynn/notifier/internal/types"
)

type (
	Email struct{}
)

func NewEmailNotifier() *Email {
	return &Email{}
}

func (e *Email) Send(context.Context, *types.Notification) (uuid.UUID, error) {
	fmt.Println("Sending email notification") // Replace with actual email sending logic")
	u, _ := uuid.NewV7()
	return u, nil
}
