package repo

import (
	"context"

	"github.com/azoma13/auth-service/internal/entity"
	"github.com/azoma13/auth-service/internal/repo/pgdb"
	"github.com/azoma13/auth-service/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (string, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
}

type Account interface {
	CreateAccount(ctx context.Context, account entity.Account) error
	DeleteAccount(ctx context.Context, userId, refreshToken string) error
	GetAccountByIdAndRefToken(ctx context.Context, userId, refreshToken string) (entity.Account, error)
	UpdateRefreshToken(ctx context.Context, id int, refreshToken, xForwardedFor string) error
}

type Repositories struct {
	User
	Account
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:    pgdb.NewUserRepo(pg),
		Account: pgdb.NewAccountRepo(pg),
	}
}
