package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TxFunc func(*sqlx.Tx) error

type Database struct {
	DB *sqlx.DB
}

func (d *Database) WithTx(ctx context.Context, fn TxFunc) error {

	tx, err := d.DB.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	defer func() {
		// panic safety
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
