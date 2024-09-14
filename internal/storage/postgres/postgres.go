package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	*pgxpool.Pool
}

func Connect(ctx context.Context, pdsn string) (*Postgres, error) {
	p, err := pgxpool.New(ctx, pdsn)
	if err != nil {
		return nil, err
	}
	return &Postgres{p}, nil
}
func (p *Postgres) Close() {
	p.Pool.Close()
}
