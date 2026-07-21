package database

import (
	"context"
	"fmt"
	"time"

	"github.com/bcpriok/pantas/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 1
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.HealthCheckPeriod = 30 * time.Second
	// Aman untuk Supavisor session maupun transaction pooler; tidak memakai
	// prepared statement lintas transaksi.
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolConfig.ConnConfig.RuntimeParams["application_name"] = "pantas"

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(checkCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect database: %w", err)
	}
	return pool, nil
}

func BootstrapAdmin(ctx context.Context, pool *pgxpool.Pool, username, name, initialPassword string) error {
	if username == "" {
		return nil
	}
	var id string
	if err := pool.QueryRow(ctx, "select public.pantas_bootstrap_admin($1, $2, $3)", username, name, initialPassword).Scan(&id); err != nil {
		return fmt.Errorf("bootstrap admin: %w", err)
	}
	return nil
}

func BootstrapTreasuryAdmin(ctx context.Context, pool *pgxpool.Pool, username, name, initialPassword string) error {
	if username == "" {
		return nil
	}
	var id string
	if err := pool.QueryRow(ctx, "select public.pantas_bootstrap_treasury_admin($1, $2, $3)", username, name, initialPassword).Scan(&id); err != nil {
		return fmt.Errorf("bootstrap treasury admin: %w", err)
	}
	return nil
}
