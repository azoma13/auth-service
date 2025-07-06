package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

func New(url string, options ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, option := range options {
		option(pg)
	}
	log.Println(url)
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("error pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for ; pg.connAttempts > 0; pg.connAttempts-- {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
	}

	if err != nil {
		return nil, fmt.Errorf("error pgxpool.NewWithConfig: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
