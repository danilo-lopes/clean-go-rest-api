package repository

import (
	"database/sql"
)

type DBExecutorMock struct {
	ExecFunc     func(query string, args ...interface{}) (sql.Result, error)
	QueryFunc    func(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowFunc func(query string, args ...interface{}) *sql.Row
	BeginFunc    func() (TxExecutor, error)
}

func (m *DBExecutorMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil, nil
}

func (m *DBExecutorMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return nil, nil
}

func (m *DBExecutorMock) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(query, args...)
	}
	return &sql.Row{}
}

func (m *DBExecutorMock) Begin() (TxExecutor, error) {
	if m.BeginFunc != nil {
		return m.BeginFunc()
	}
	return nil, nil
}
