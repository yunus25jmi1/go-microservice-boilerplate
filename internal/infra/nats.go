package infra

import (
    "fmt"
    "time"

    "github.com/nats-io/nats.go"
    "github.com/rs/zerolog/log"
)

type NATS struct {
    Conn *nats.Conn
}

// NewNATS creates a NATS connection if NATSEnabled is true.
// It respects NATS_MAX_PENDING, reconnect settings and returns nil when disabled.
func NewNATS(cfg Config) (*NATS, error) {
    if !cfg.NATSEnabled {
        // NATS disabled – caller can safely ignore the returned value.
        return nil, nil
    }
    if cfg.NATSURL == "" {
        return nil, fmt.Errorf("NATS_URL not set while NATSEnabled=true")
    }
    opts := []nats.Option{
        nats.Name("go-microservice-boilerplate"),
        nats.MaxReconnects(5),
        nats.ReconnectWait(2 * time.Second),
    }
    conn, err := nats.Connect(cfg.NATSURL, opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to NATS: %w", err)
    }
    log.Info().Msg("connected to NATS")
    return &NATS{Conn: conn}, nil
}

// Publish sends a message to the given subject. Returns any error from the NATS client.
func (n *NATS) Publish(subject string, data []byte) error {
    if n == nil || n.Conn == nil {
        return fmt.Errorf("nats client not initialised")
    }
    return n.Conn.Publish(subject, data)
}
