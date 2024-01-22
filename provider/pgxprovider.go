package provider

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type PgxProvider struct {
	connectionPool *pgxpool.Pool
}

type txContainer struct {
	tx pgx.Tx
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

	_, err = con.Exec(ctx, query, params...)
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

	return con.QueryRow(ctx, query, params...), nil
}

func (p *PgxProvider) ExecuteQueryRows(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error) {
	con, err := p.connectionPool.Acquire(ctx)

	if err != nil {
		return nil, err
	}
	defer con.Release()

	return con.Query(ctx, query, params...)
}

func (p *PgxProvider) PerformTx(ctx context.Context) (*txContainer, error) {
	con, err := p.connectionPool.Acquire(ctx)

	if err != nil {
		return nil, err
	}

	tx, err := con.Begin(ctx)

	if err != nil {
		return nil, err
	}

	return &txContainer{tx: tx}, nil
}

func (t *txContainer) ExecuteQuery(ctx context.Context, query string, params ...interface{}) error {
	tag, err := t.tx.Exec(ctx, query, params...)

	logrus.Infoln("Execute query: ", tag)

	if err != nil {
		t.tx.Rollback(ctx)
		t.tx.Conn().Close(ctx)
		return err
	}

	return nil
}

func (t *txContainer) ExecuteQueryWithRow(ctx context.Context, query string, params ...interface{}) pgx.Row {
	row := t.tx.QueryRow(ctx, query, params...)
	logrus.Infoln("Row from db: ", row)

	return row

}

func (t *txContainer) ExecuteQueryRow(ctx context.Context, query string, params ...interface{}) (*pgx.Rows, error) {
	rows, err := t.tx.Query(ctx, query, params...)

	if err != nil {
		t.tx.Rollback(ctx)
		t.tx.Conn().Close(ctx)
		return nil, err
	}

	return &rows, nil
}

func (t *txContainer) CommitTx(ctx context.Context) error {
	err := t.tx.Commit(ctx)

	if err != nil {
		t.tx.Rollback(ctx)
		return err
	}

	logrus.Infoln("Tx was commited")

	return nil
}
