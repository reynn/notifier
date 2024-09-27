package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/reynn/notifier/internal/types"
)

type (
	postgresNotification struct{}
	Config               struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	Postgres struct {
		db *sqlx.DB
	}
)

func NewPostgres(cfg Config) *Postgres {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	return &Postgres{
		db: sqlx.MustConnect("postgres", dbInfo),
	}
}

func (s *Postgres) Write(ctx context.Context, n *types.Notification) error {
	_, err := s.db.NamedExec(postgresWrite, n)
	if err != nil {
		return fmt.Errorf("failed to write notification: %w", err)
	}
	return nil
}

func (s *Postgres) ByID(ctx context.Context, id uuid.UUID) (*types.Notification, error) {
	notif := &types.Notification{}
	err := s.db.Get(notif, postgresRetrieveByIDSQL, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to retrieve notification: %w", err)
	}
	return notif, nil
}
