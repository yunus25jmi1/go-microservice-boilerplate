package infra

import (
    "time"
)

type Config struct {
    // Core services
    PostgresURL      string `envconfig:"POSTGRES_URL" required:"true"`
    PostgresMaxConns int    `envconfig:"POSTGRES_MAX_CONNS" default:"10"`
    PostgresMinConns int    `envconfig:"POSTGRES_MIN_CONNS" default:"2"`
    PostgresTLS      bool   `envconfig:"POSTGRES_TLS" default:"false"`

    // NATS (optional)
    NATSEnabled      bool   `envconfig:"NATS_ENABLED" default:"false"`
    NATSURL          string `envconfig:"NATS_URL"`
    NATSMaxPending   int    `envconfig:"NATS_MAX_PENDING" default:"65536"`

    // Feature toggles
    EnableGRPC  bool          `envconfig:"ENABLE_GRPC" default:"false"`
    EnableCache bool          `envconfig:"ENABLE_CACHE" default:"false"`
    ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}
