package repository

import (
	"database/sql"
)

type dbExecutorAdapter struct {
	db *sql.DB
}

func NewDBExecutorAdapter(db *sql.DB) DBExecutor {
	return &dbExecutorAdapter{db: db}
}

func (a *dbExecutorAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.db.Exec(query, args...)
}

func (a *dbExecutorAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.Query(query, args...)
}

func (a *dbExecutorAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.db.QueryRow(query, args...)
}

func (a *dbExecutorAdapter) Begin() (TxExecutor, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return nil, err
	}
	return &txExecutorAdapter{tx: tx}, nil
}

type txExecutorAdapter struct {
	tx *sql.Tx
}

func (t *txExecutorAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t *txExecutorAdapter) Rollback() error {
	return t.tx.Rollback()
}

func (t *txExecutorAdapter) Commit() error {
	return t.tx.Commit()
}
