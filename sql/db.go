package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/agurinov/gopl/diag/trace"
)

type (
	DB struct {
		sqlxClient *sqlx.DB
	}
)

func (db DB) NamedQueryContext(
	ctx context.Context,
	query Query,
	args map[string]any,
) (
	*sqlx.Rows,
	error,
) {
	ctx, span := trace.StartSpan(ctx, "db.NamedQueryContext")
	defer span.End()

	//nolint:sqlclosecheck
	rows, err := db.sqlxClient.NamedQueryContext(ctx, query.WithSpan(span), args)
	if err != nil {
		return nil, trace.CatchError(span, err)
	}

	return rows, nil
}

func (db DB) NamedExecContext(
	ctx context.Context,
	query Query,
	args map[string]any,
) (
	sql.Result,
	error,
) {
	ctx, span := trace.StartSpan(ctx, "db.NamedExecContext")
	defer span.End()

	result, err := db.sqlxClient.NamedExecContext(ctx, query.WithSpan(span), args)
	if err != nil {
		return nil, trace.CatchError(span, err)
	}

	return result, nil
}

func (db DB) MultiNamedExecContext(
	ctx context.Context,
	query Query,
	args []map[string]any,
) (
	sql.Result,
	error,
) {
	ctx, span := trace.StartSpan(ctx, "db.MultiNamedExecContext")
	defer span.End()

	result, err := db.sqlxClient.NamedExecContext(ctx, query.WithSpan(span), args)
	if err != nil {
		return nil, trace.CatchError(span, err)
	}

	return result, nil
}
