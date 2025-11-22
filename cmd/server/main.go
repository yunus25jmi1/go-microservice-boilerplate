package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/kelseyhightower/envconfig"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/jackc/pgx/v5/pgxpool"
    "sync"

    "github.com/yourorg/go-microservice-boilerplate/encoding"
    "github.com/yourorg/go-microservice-boilerplate/internal/domain"
    "github.com/yourorg/go-microservice-boilerplate/internal/infra"
)

func main() {
    // Logger setup
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    log.Info().Msg("starting service")

    // Initialise infrastructure (DB and optional NATS) concurrently
    var wg sync.WaitGroup
    wg.Add(2)

    var db *pgxpool.Pool
    var natsClient *infra.NATS
    var initErr error

    // DB init
    go func() {
        defer wg.Done()
        var err error
        db, err = infra.NewPostgres(cfg.PostgresURL, cfg)
        if err != nil {
            initErr = fmt.Errorf("postgres init: %w", err)
        }
    }()

    // NATS init (optional)
    go func() {
        defer wg.Done()
        var err error
        natsClient, err = infra.NewNATS(cfg)
        if err != nil {
            initErr = fmt.Errorf("nats init: %w", err)
        }
    }()

    wg.Wait()
    if initErr != nil {
        log.Fatal().Err(initErr).Msg("failed to initialise infrastructure")
    }
    defer db.Close()
    if natsClient != nil {
        defer natsClient.Conn.Drain()
    }

    // Initialise domain services
    userRepo := infra.NewUserRepository(db)
    userService := domain.NewUserService(userRepo)

    // HTTP router
    r := chi.NewRouter()
    r.Get("/healthz", healthHandler)
    r.Handle("/metrics", promhttp.Handler())
    r.Get("/users/{id}", getUserHandler(userService))

    srv := &http.Server{Addr: ":8080", Handler: r}

    // Graceful shutdown handling
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal().Err(err).Msg("http server error")
        }
    }()
    log.Info().Msg("server listening on :8080")

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info().Msg("shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Error().Err(err).Msg("server forced to shutdown")
    }
    log.Info().Msg("server exited")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}

func getUserHandler(svc domain.UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        user, err := svc.GetUser(r.Context(), id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        // Use zero‑allocation JSON encoder
        resp := struct {
            ID   string `json:"id"`
            Name string `json:"name"`
        }{ID: user.ID, Name: user.Name}
        data, err := encoding.JSON(resp)
        if err != nil {
            http.Error(w, "internal error", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(data)
    }
}
