package health

import (
	"context"
	"log/slog"
	"maps"
	"slices"

	"connectrpc.com/grpchealth"
	"golang.org/x/sync/errgroup"
)

type (
	Pinger interface {
		Ping(context.Context) error
	}

	Checker struct {
		pingables map[string]Pinger
		logger    *slog.Logger
	}
)

func New(logger *slog.Logger, pingables map[string]Pinger) *Checker {
	return &Checker{
		pingables: pingables,
		logger:    logger,
	}
}

func (h *Checker) Check(ctx context.Context, req *grpchealth.CheckRequest) (*grpchealth.CheckResponse, error) {
	if req.Service != "" {
		return checks(ctx, []string{req.Service}, h.pingables)
	}
	return checks(ctx, slices.Sorted(maps.Keys(h.pingables)), h.pingables)
}

func checks(ctx context.Context, services []string, pingables map[string]Pinger) (*grpchealth.CheckResponse, error) {
	wg, gCtx := errgroup.WithContext(ctx)
	for _, service := range services {
		wg.Go(func() error {
			if e := pingables[service].Ping(gCtx); e != nil {
				return e
			}
			return nil
		})
	}
	if e := wg.Wait(); e != nil {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusNotServing}, e
	}
	return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
}
