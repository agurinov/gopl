package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect(
	ctx context.Context,
	driver string,
	dsn string,
) (
	DB,
	error,
) {
	var sqlxClient *sqlx.DB

	switch driver {
	case "noop":
	case "pgx":
		driverConfig, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return DB{}, err
		}

		pool, err := pgxpool.NewWithConfig(ctx, driverConfig)
		if err != nil {
			return DB{}, err
		}

		sqlDB := stdlib.OpenDBFromPool(pool)
		if pingErr := sqlDB.PingContext(ctx); pingErr != nil {
			return DB{}, pingErr
		}

		sqlxClient = sqlx.NewDb(sqlDB, driver)
	default:
		c, err := sqlx.ConnectContext(ctx, driver, dsn)
		if err != nil {
			return DB{}, err
		}

		sqlxClient = c
	}

	db := DB{
		sqlxClient: sqlxClient,
	}

	return db, nil
}
