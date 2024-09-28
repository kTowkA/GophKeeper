package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres объект реализующий интерфейс Storager
type Postgres struct {
	*pgxpool.Pool
}

// Connect получение пула соединений Postgres
func Connect(ctx context.Context, pdsn string) (*Postgres, error) {
	p, err := pgxpool.New(ctx, pdsn)
	if err != nil {
		return nil, err
	}
	err = p.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &Postgres{p}, nil
}

// Close закрытие пула соединений
func (p *Postgres) Close() {
	p.Pool.Close()
}
