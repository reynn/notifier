package notifier

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	notifierv1 "github.com/reynn/notifier/gen/proto/notifier/v1"
	v1 "github.com/reynn/notifier/gen/proto/notifier/v1"
	"github.com/reynn/notifier/gen/proto/notifier/v1/notifierv1connect"
	"github.com/reynn/notifier/internal/notifiers"
	"github.com/reynn/notifier/internal/retrievers"
	"github.com/reynn/notifier/internal/types"
)

type (
	Service struct {
		retriever retrievers.Retriever
		notifiers map[v1.NotificationType]notifiers.Sender
	}
	Config struct {
		Retriever retrievers.Retriever
		Notifiers map[v1.NotificationType]notifiers.Sender
	}
)

func NewService(c Config) *Service {
	return &Service{
		retriever: c.Retriever,
		notifiers: c.Notifiers,
	}
}

func (s *Service) Register(mux *http.ServeMux) {
	oc, _ := otelconnect.NewInterceptor()
	mux.Handle(notifierv1connect.NewNotificationServiceHandler(
		s, connect.WithInterceptors(oc),
	))
}

func (s *Service) SendNotification(ctx context.Context, req *connect.Request[v1.SendNotificationRequest]) (*connect.Response[v1.SendNotificationResponse], error) {
	if notifier, ok := s.notifiers[req.Msg.Type]; ok {
		notif := &types.Notification{
			Recipients: req.Msg.Recipients,
			Message:    []byte(req.Msg.Message),
			Tags:       req.Msg.Tags,
			Priority:   req.Msg.Priority,
			Type:       req.Msg.Type,
			Status:     notifierv1.NotificationStatus_SUBMITTED,
			CreatedAt:  time.Now(),
		}
		id, e := notifier.Send(ctx, notif)
		if e != nil {
			return nil, connect.NewError(connect.CodeInternal, e)
		}
		if e := s.retriever.Store(ctx, id, notif); e != nil {
			return nil, connect.NewError(connect.CodeInternal, e)
		}
		res := connect.NewResponse(&v1.SendNotificationResponse{
			Id: id.String(),
		})
		return res, nil
	}
	return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("%s is not a configured notifier", req.Msg.Type.String()))
}

func (s *Service) GetNotification(ctx context.Context, req *connect.Request[v1.GetNotificationRequest]) (*connect.Response[v1.GetNotificationResponse], error) {
	notif, err := s.retriever.ByID(ctx, uuid.MustParse(req.Msg.Id))
	if err != nil {
		if errors.Is(err, retrievers.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("notification not found: %w", err))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to retrieve notification: %w", err))
	}
	res := connect.NewResponse(&v1.GetNotificationResponse{
		Id:         notif.ID,
		Recipients: notif.Recipients,
		Message:    string(notif.Message),
		Tags:       notif.Tags,
		CreatedAt:  timestamppb.New(notif.CreatedAt),
		Type:       notif.Type,
		Priority:   notif.Priority,
	})
	return res, nil
}

func (s *Service) DeleteNotification(context.Context, *connect.Request[v1.DeleteNotificationRequest]) (*connect.Response[v1.DeleteNotificationResponse], error) {
	res := connect.NewResponse(&v1.DeleteNotificationResponse{})
	return res, nil
}
