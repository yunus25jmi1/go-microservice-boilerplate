package infra

import (
    "context"
    "crypto/tls"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/rs/zerolog/log"
)

// NewPostgres creates a pgx connection pool using the provided URL.
// It also wraps a health‑check ping with a circuit‑breaker.
func NewPostgres(url string, cfg Config) (*pgxpool.Pool, error) {
    pcfg, err := pgxpool.ParseConfig(url)
    if err != nil {
        return nil, fmt.Errorf("parse pgx config: %w", err)
    }
    // Pool sizing
    pcfg.MaxConns = int32(cfg.PostgresMaxConns)
    pcfg.MinConns = int32(cfg.PostgresMinConns)
    // TLS enforcement
    if cfg.PostgresTLS {
        pcfg.ConnConfig.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
    }
    pool, err := pgxpool.NewWithConfig(context.Background(), pcfg)
    if err != nil {
        return nil, fmt.Errorf("create pgx pool: %w", err)
    }
    // Circuit‑breaker for health‑check
    cb := NewBreaker("PostgresPing", cfg.ShutdownTimeout)
    _, err = cb.Execute(func() (interface{}, error) {
        return nil, pool.Ping(context.Background())
    })
    if err != nil {
        pool.Close()
        return nil, fmt.Errorf("postgres health check failed: %w", err)
    }
    log.Info().Msg("connected to postgres")
    return pool, nil
}
