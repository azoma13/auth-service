package pgdb

import (
	"context"
	"fmt"

	"github.com/azoma13/auth-service/internal/entity"
	"github.com/azoma13/auth-service/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (string, error) {
	query := `
		INSERT INTO users
			(username, password)
			VALUES ($1, $2)
		RETURNING id
	`

	var id string
	err := r.Pool.QueryRow(ctx, query, user.Username, user.Password).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %w", err)
	}

	return id, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	query := `
			SELECT id, username, password, created_at
				FROM users
			WHERE username = $1
		`
	var employee entity.User
	err := r.Pool.QueryRow(ctx, query, username).Scan(
		&employee.Id,
		&employee.Username,
		&employee.Password,
		&employee.CreatedAt,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
	}

	return employee, nil
}

func (r *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
	query := `
			SELECT id, username, password, created_at
				FROM users
			WHERE username = $1 AND password = $2
		`
	var user entity.User
	err := r.Pool.QueryRow(ctx, query, username, password).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}
