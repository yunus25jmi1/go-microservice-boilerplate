package infra

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

// startPostgresContainer spins up a temporary PostgreSQL container for testing.
func startPostgresContainer(t *testing.T) (dsn string, cleanup func()) {
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15-alpine",
        Env:          map[string]string{"POSTGRES_PASSWORD": "pwd", "POSTGRES_USER": "postgres", "POSTGRES_DB": "postgres"},
        ExposedPorts: []string{"5432/tcp"},
        WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
    }
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
    if err != nil {
        t.Fatalf("failed to start postgres container: %v", err)
    }
    host, err := container.Host(ctx)
    if err != nil {
        t.Fatalf("failed to get container host: %v", err)
    }
    port, err := container.MappedPort(ctx, "5432")
    if err != nil {
        t.Fatalf("failed to get mapped port: %v", err)
    }
    dsn = fmt.Sprintf("postgres://postgres:pwd@%s:%s/postgres?sslmode=disable", host, port.Port())
    cleanup = func() {
        _ = container.Terminate(ctx)
    }
    return dsn, cleanup
}

// TestPostgresIntegration verifies that a temporary DB can be created, a schema migrated, and cleaned up.
func TestPostgresIntegration(t *testing.T) {
    dsn, cleanup := startPostgresContainer(t)
    defer cleanup()

    // Create a unique schema for this test.
    schema := "test_" + uuid.New().String()
    cfg := Config{
        PostgresURL:      dsn,
        PostgresMaxConns: 5,
        PostgresMinConns: 1,
        PostgresTLS:      false,
        ShutdownTimeout: 5 * time.Second,
    }
    db, err := NewPostgres(cfg.PostgresURL, cfg)
    if err != nil {
        t.Fatalf("failed to init postgres: %v", err)
    }
    defer db.Close()

    // Create schema.
    if _, err := db.Exec(context.Background(), fmt.Sprintf("CREATE SCHEMA %s", schema)); err != nil {
        t.Fatalf("failed to create schema: %v", err)
    }
    // Run a tiny migration inside the schema.
    migration := fmt.Sprintf(`
        SET search_path TO %s;
        CREATE TABLE users (id TEXT PRIMARY KEY, name TEXT);
        INSERT INTO users (id, name) VALUES ('u1', 'Alice');
    `, schema)
    if _, err := db.Exec(context.Background(), migration); err != nil {
        t.Fatalf("migration failed: %v", err)
    }
    // Verify data.
    var name string
    row := db.QueryRow(context.Background(), fmt.Sprintf("SELECT name FROM %s.users WHERE id=$1", schema), "u1")
    if err := row.Scan(&name); err != nil {
        t.Fatalf("query failed: %v", err)
    }
    if name != "Alice" {
        t.Fatalf("expected Alice, got %s", name)
    }
    // Clean up schema.
    if _, err := db.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA %s CASCADE", schema)); err != nil {
        t.Fatalf("failed to drop schema: %v", err)
    }
}
