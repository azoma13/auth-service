package pgdb

import (
	"context"
	"fmt"

	"github.com/azoma13/auth-service/internal/entity"
	"github.com/azoma13/auth-service/pkg/postgres"
)

type AccountRepo struct {
	*postgres.Postgres
}

func NewAccountRepo(pg *postgres.Postgres) *AccountRepo {
	return &AccountRepo{pg}
}

func (r *AccountRepo) CreateAccount(ctx context.Context, account entity.Account) error {
	exec := `
		INSERT INTO accounts
			(user_id, refresh_token, user_agent, x_forwarded_for)
		VALUES ($1, crypt($2, gen_salt('bf', 10)), $3, $4)
	`
	_, err := r.Pool.Exec(ctx, exec, account.UserId, account.RefreshToken, account.UserAgent, account.XForwardedFor)
	if err != nil {
		return fmt.Errorf("AccountRepo.CreateAccount - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *AccountRepo) DeleteAccount(ctx context.Context, userId, refreshToken string) error {
	exec := `
		DELETE FROM accounts
		WHERE user_id = $1 AND refresh_token = crypt($2, refresh_token);
	`

	_, err := r.Pool.Exec(ctx, exec, userId, refreshToken)
	if err != nil {
		return fmt.Errorf("AccountRepo.DeleteAccount - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *AccountRepo) GetAccountByIdAndRefToken(ctx context.Context, userId, refreshToken string) (entity.Account, error) {
	query := `
			SELECT id, user_id, refresh_token, user_agent, x_forwarded_for, created_at
				FROM accounts
			WHERE user_id = $1 AND refresh_token = crypt($2, refresh_token);
		`
	var account entity.Account
	err := r.Pool.QueryRow(ctx, query, userId, refreshToken).Scan(
		&account.Id,
		&account.UserId,
		&account.RefreshToken,
		&account.UserAgent,
		&account.XForwardedFor,
		&account.CreatedAt,
	)

	if err != nil {
		return entity.Account{}, fmt.Errorf("AccountRepo.GetAccountByIdAndRefToken - r.Pool.QueryRow: %v", err)
	}

	return account, nil
}

func (r *AccountRepo) UpdateRefreshToken(ctx context.Context, id int, refreshToken, xForwardedFor string) error {
	exec := `
		UPDATE accounts
			SET refresh_token = crypt($1, gen_salt('bf', 10)), x_forwarded_for = $2
		WHERE id = $3;
	`
	_, err := r.Pool.Exec(ctx, exec, refreshToken, xForwardedFor, id)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}
