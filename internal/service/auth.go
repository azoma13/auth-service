package service

import (
	"context"
	"fmt"

	"github.com/azoma13/auth-service/config"
	"github.com/azoma13/auth-service/internal/entity"
	"github.com/azoma13/auth-service/internal/repo"
	"github.com/azoma13/auth-service/pkg/hasher"
	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId string
}

type AuthService struct {
	userRepo       repo.User
	accountRepo    repo.Account
	passwordHasher hasher.PasswordHasher
}

func NewAuthService(userRepo repo.User, accountRepo repo.Account, passwordHasher hasher.PasswordHasher) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		accountRepo:    accountRepo,
		passwordHasher: passwordHasher,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput) (string, error) {
	user := &entity.User{
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
	}

	userId, err := s.userRepo.CreateUser(ctx, *user)
	if err != nil {
		return "", fmt.Errorf("AuthService.CreateUser - s.userRepo.CreateUser: %w", err)
	}

	return userId, nil
}

func (s *AuthService) CreateAccount(ctx context.Context, input AuthCreateAccountInput) error {
	account := &entity.Account{
		RefreshToken:  input.RefreshToken,
		UserAgent:     input.UserAgent,
		XForwardedFor: input.XForwardedFor,
	}

	src := ctx.Value("source").(string)
	if src == "logInWithId" {
		account.UserId = input.UserId
	} else {
		user, err := s.userRepo.GetUserByUsername(ctx, input.Username)
		if err != nil {
			return fmt.Errorf("AuthService.CreateAccount - s.userRepo.GetUserByUsername: %w", err)
		}

		account.UserId = user.Id
	}

	err := s.accountRepo.CreateAccount(ctx, *account)
	if err != nil {
		return fmt.Errorf("AuthService.CreateAccount - s.accountRepo.CreateAccount: %w", err)
	}
	return nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	src := ctx.Value("source").(string)
	if src == "logInWithId" {
		input.TokenClaims.UserId = input.Id
	} else {
		user, err := s.userRepo.GetUserByUsernameAndPassword(ctx, input.Username, s.passwordHasher.Hash(input.Password))
		if err != nil {
			return "", fmt.Errorf("AuthService.GenerateToken: cannot get user: %v", err)
		}
		input.TokenClaims.UserId = user.Id
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, input.TokenClaims)

	tokenString, err := token.SignedString([]byte(config.Cfg.SignKey))
	if err != nil {
		return "", fmt.Errorf("AuthService.GenerateToken: cannot sign token: %v", err)
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Cfg.SignKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("error parse token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", fmt.Errorf("error parse token")
	}

	return claims.UserId, nil
}
