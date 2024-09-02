package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"

	notifierv1 "github.com/reynn/notifier/gen/proto/notifier/v1"
	"github.com/reynn/notifier/internal/api"
	"github.com/reynn/notifier/internal/api/notifier"
	"github.com/reynn/notifier/internal/config"
	"github.com/reynn/notifier/internal/constants"
	"github.com/reynn/notifier/internal/notifiers"
	"github.com/reynn/notifier/internal/retrievers"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
		Level: slog.LevelDebug.Level(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	})).With(
		slog.String("version", constants.AppVersion),
		slog.String("app-name", constants.AppName),
		slog.String("app-module", constants.AppModule),
	)

	cfg := config.Load()

	logger.Info("config file", slog.String("path", cfg.ConfigFilePath))

	apiServer := api.New(&api.Config{
		Port:   cfg.HTTPPort,
		Logger: logger,
		Notifier: notifier.NewService(notifier.Config{
			Retriever: retrievers.NewInMemoryRetriever(),
			Notifiers: map[notifierv1.NotificationType]notifiers.Sender{
				notifierv1.NotificationType_EMAIL: notifiers.NewEmailNotifier(),
			},
		}),
	})

	wg, gCtx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return apiServer.Start(gCtx)
	})

	if err := wg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("application error: %v", err)
	}
}
