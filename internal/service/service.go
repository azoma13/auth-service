package service

import (
	"context"

	"github.com/azoma13/auth-service/internal/entity"
	"github.com/azoma13/auth-service/internal/repo"
	"github.com/azoma13/auth-service/pkg/hasher"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AuthGenerateTokenInput struct {
	Id          string
	Username    string
	Password    string
	TokenClaims TokenClaims
}

type AuthCreateAccountInput struct {
	UserId        string
	Username      string
	RefreshToken  string
	UserAgent     string
	XForwardedFor string
}

type AuthDeleteAccountInput struct {
	UserId       string
	RefreshToken string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (string, error)
	CreateAccount(ctx context.Context, input AuthCreateAccountInput) error
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (string, error)
}

type AccountGetInput struct {
	UserId        string
	RefreshToken  string
	UserAgent     string
	XForwardedFor string
}

type AccountGenerateTokenInput struct {
	Id          string
	TokenClaims TokenClaims
}

type AccountUpdateInput struct {
	Id            int
	RefreshToken  string
	XForwardedFor string
}

type Account interface {
	GetAccount(ctx context.Context, input AccountGetInput) (entity.Account, error)
	GenerateToken(ctx context.Context, tokenClaims TokenClaims) (string, error)
	UpdateRefreshToken(ctx context.Context, input AccountUpdateInput) error
	DeleteAccount(ctx context.Context, input AuthDeleteAccountInput) error
}

type Services struct {
	Auth    Auth
	Account Account
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:    NewAuthService(deps.Repos.User, deps.Repos.Account, deps.Hasher),
		Account: NewAccountService(deps.Repos.Account),
	}
}
