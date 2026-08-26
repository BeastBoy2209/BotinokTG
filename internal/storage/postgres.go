package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const pingTimeout = 5 * time.Second

func NewPostgres(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	slog.Info("pool was created... \n")
	pingCtx, pingCancel := context.WithTimeout(ctx, pingTimeout)
	defer pingCancel()

	if err := dbpool.Ping(pingCtx); err != nil {
		dbpool.Close()

		return nil, err
	}
	slog.Info("database connection established successfully...")

	return dbpool, nil
}
