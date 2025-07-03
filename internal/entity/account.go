package entity

import "time"

type Account struct {
	Id            int       `db:"db"`
	UserId        string    `db:"user_id"`
	RefreshToken  string    `db:"refresh_token"`
	UserAgent     string    `db:"user_agent"`
	XForwardedFor string    `db:"x_forwarded_for"`
	CreatedAt     time.Time `db:"created_at"`
}
