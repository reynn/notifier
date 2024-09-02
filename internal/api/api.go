package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"

	"github.com/reynn/notifier/gen/proto/notifier/v1/notifierv1connect"
	"github.com/reynn/notifier/internal/api/health"
	"github.com/reynn/notifier/internal/api/notifier"
)

type (
	Config struct {
		Port     int
		Logger   *slog.Logger
		Notifier *notifier.Service
	}
	Server struct {
		s        *http.Server
		port     int
		logger   *slog.Logger
		notifier *notifier.Service
	}
)

func New(c *Config) *Server {
	mux := http.NewServeMux()

	// setup gRPC reflection handler to allow clients to discover the available services and methods.
	reflector := grpcreflect.NewStaticReflector(notifierv1connect.NotificationServiceName)
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle(grpcreflect.NewHandlerV1(reflector))

	c.Notifier.Register(mux)

	// setup health check handler to allow clients/k8s to check the status of the server.
	hc := health.New(c.Logger, map[string]health.Pinger{
		"google": health.NewHTTPChecker("https://google.com", http.DefaultClient),
		"bing":   health.NewHTTPChecker("https://bing.com", http.DefaultClient),
	})
	mux.Handle(grpchealth.NewHandler(hc))

	return &Server{
		s: &http.Server{
			Addr:    fmt.Sprintf(":%d", c.Port),
			Handler: h2c.NewHandler(mux, &http2.Server{}),
		},
		port:   c.Port,
		logger: c.Logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("server started", slog.Int("port", s.port))
	wg, gCtx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		if e := s.s.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			return e
		}
		return nil
	})
	wg.Go(func() error {
		<-gCtx.Done()
		toCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.s.Shutdown(toCtx)
	})

	if e := wg.Wait(); e != nil {
		return e
	}

	return nil
}
