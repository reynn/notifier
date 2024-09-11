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
			Priority:   types.ParseNotificationPriority(req.Msg.Priority.String()),
			Type:       types.ParseNotificationType(req.Msg.Type.String()),
			Status:     types.NotificationStatusSubmitted,
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
		Id:         notif.ID.String(),
		Recipients: notif.Recipients,
		Message:    string(notif.Message),
		Tags:       notif.Tags,
		CreatedAt:  timestamppb.New(notif.CreatedAt),
		Type:       NotificationTypeToProto(notif.Type),
		Priority:   NotificationPriorityToProto(notif.Priority),
		Status:     NotificationStatusToProto(notif.Status),
	})
	return res, nil
}

func (s *Service) DeleteNotification(context.Context, *connect.Request[v1.DeleteNotificationRequest]) (*connect.Response[v1.DeleteNotificationResponse], error) {
	res := connect.NewResponse(&v1.DeleteNotificationResponse{})
	return res, nil
}

func NotificationTypeToProto(t types.NotificationType) v1.NotificationType {
	switch t {
	case types.NotificationTypeEmail:
		return v1.NotificationType_EMAIL
	case types.NotificationTypeSMS:
		return v1.NotificationType_SMS
	case types.NotificationTypePush:
		return v1.NotificationType_PUSH
	case types.NotificationTypeWebhook:
		return v1.NotificationType_WEBHOOK
	default:
		return v1.NotificationType_UNSET
	}
}

func NotificationPriorityToProto(p types.NotificationPriority) v1.NotificationPriority {
	switch p {
	case types.NotificationPriorityHigh:
		return v1.NotificationPriority_HIGH
	case types.NotificationPriorityLow:
		return v1.NotificationPriority_LOW
	default:
		return v1.NotificationPriority_DEFAULT
	}
}

func NotificationStatusToProto(s types.NotificationStatus) v1.NotificationStatus {
	switch s {
	case types.NotificationStatusCompleted:
		return v1.NotificationStatus_COMPLETED
	case types.NotificationStatusFailed:
		return v1.NotificationStatus_FAILED
	default:
		return v1.NotificationStatus_SUBMITTED
	}
}
