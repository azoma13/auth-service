package service

import (
	"context"
	"fmt"

	"github.com/azoma13/auth-service/config"
	"github.com/azoma13/auth-service/internal/entity"
	"github.com/azoma13/auth-service/internal/repo"
	"github.com/golang-jwt/jwt"
)

type AccountService struct {
	accountRepo repo.Account
}

func NewAccountService(accountRepo repo.Account) *AccountService {
	return &AccountService{
		accountRepo: accountRepo}
}

func (s *AccountService) GetAccount(ctx context.Context, input AccountGetInput) (entity.Account, error) {
	account, err := s.accountRepo.GetAccountByIdAndRefToken(ctx, input.UserId, input.RefreshToken)
	if err != nil {
		return entity.Account{}, fmt.Errorf("AccountService.RefreshTokens - s.userRepo.GetUserByUsername: %w", err)
	}
	if account.UserAgent != input.UserAgent {
		return entity.Account{}, ErrDifferentUserAgent
	}
	if account.XForwardedFor != input.XForwardedFor {
		return account, ErrDifferentXForwardedFor
	}

	return account, nil
}

func (s *AccountService) GenerateToken(ctx context.Context, tokenClaims TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	tokenString, err := token.SignedString([]byte(config.Cfg.SignKey))
	if err != nil {
		return "", fmt.Errorf("AuthService.GenerateToken: cannot sign token: %v", err)
	}

	return tokenString, nil
}

func (s *AccountService) UpdateRefreshToken(ctx context.Context, input AccountUpdateInput) error {
	err := s.accountRepo.UpdateRefreshToken(ctx, input.Id, input.RefreshToken, input.XForwardedFor)

	return err
}

func (s *AccountService) DeleteAccount(ctx context.Context, input AuthDeleteAccountInput) error {
	err := s.accountRepo.DeleteAccount(ctx, input.UserId, input.RefreshToken)
	if err != nil {
		return fmt.Errorf("AuthService.DeleteAccount - s.accountRepo.DeleteAccount: %w", err)
	}

	return nil
}
