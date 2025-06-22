package repository

import "database/sql"

type TxMock struct {
	ExecFunc     func(query string, args ...interface{}) (sql.Result, error)
	RollbackFunc func() error
	CommitFunc   func() error
}

func (m *TxMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil, nil
}

func (m *TxMock) Rollback() error {
	if m.RollbackFunc != nil {
		return m.RollbackFunc()
	}
	return nil
}

func (m *TxMock) Commit() error {
	if m.CommitFunc != nil {
		return m.CommitFunc()
	}
	return nil
}
