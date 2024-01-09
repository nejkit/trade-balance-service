package provider

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxProvider struct {
	connectionPool *pgxpool.Pool
}

func NewPgxProvider(connectionPool *pgxpool.Pool) PgxProvider {
	return PgxProvider{connectionPool: connectionPool}
}

func (p *PgxProvider) ExecuteQuery(ctx context.Context, query string, params ...interface{}) error {
	con, err := p.connectionPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer con.Release()

	_, err = con.Exec(ctx, query, params)
	if err != nil {
		return err
	}

	return nil
}

func (p *PgxProvider) ExecuteQueryWithRow(ctx context.Context, query string, params ...interface{}) (pgx.Row, error) {
	con, err := p.connectionPool.Acquire(ctx)

	if err != nil {
		return nil, err
	}

	defer con.Release()

	return con.QueryRow(ctx, query, params), nil
}

func (p *PgxProvider) ExecuteQueryRows(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error) {
	con, err := p.connectionPool.Acquire(ctx)

	if err != nil {
		return nil, err
	}
	defer con.Release()

	return con.Query(ctx, query, params)
}
