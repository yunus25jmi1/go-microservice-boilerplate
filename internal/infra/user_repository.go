package infra

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/yourorg/go-microservice-boilerplate/internal/domain"
    "github.com/rs/zerolog/log"
)

type PostgresUserRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    row := r.db.QueryRow(ctx, "SELECT id, name FROM users WHERE id=$1", id)
    if err := row.Scan(&user.ID, &user.Name); err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    return &user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
    _, err := r.db.Exec(ctx, "INSERT INTO users (id, name) VALUES ($1, $2)", user.ID, user.Name)
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }
    log.Info().Msgf("created user %s", user.ID)
    return nil
}
