package repository

import (
	"clean-go-rest-api/internal/domain/entity"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	testName  string
	repoSetup func(*DBExecutorMock)
	input     entity.User
	expected  error
}

func TestUserRepository_Add(t *testing.T) {
	tests_scenarios := []testCase{
		{
			testName: "Valid User Creation",
			repoSetup: func(repo *DBExecutorMock) {
				repo.ExecFunc = func(query string, args ...interface{}) (sql.Result, error) {
					return nil, nil
				}
			},
			input: entity.User{
				ID:    uuid.New(),
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			expected: nil,
		},
		{
			testName: "Error on User Creation",
			repoSetup: func(repo *DBExecutorMock) {
				repo.ExecFunc = func(query string, args ...interface{}) (sql.Result, error) {
					return nil, errors.New("database error")
				}
			},
			input: entity.User{
				ID:    uuid.New(),
				Name:  "Jane Mary",
				Email: "jane.mary@example.com",
			},
			expected: errors.New("database error"),
		},
	}

	for _, tt := range tests_scenarios {
		t.Run(tt.testName, func(t *testing.T) {
			dbExecutor := &DBExecutorMock{}
			tt.repoSetup(dbExecutor)

			repo := NewPostgresUserRepository(dbExecutor)
			err := repo.Add(tt.input)

			switch tt.testName {
			case tests_scenarios[0].testName:
				assert.NoError(t, err, "Expected no error for valid user creation")
				assert.Equal(t, tt.expected, err, "Expected error equal to nil")
			case tests_scenarios[1].testName:
				assert.Error(t, err, "Expected an error for user creation failure")
				assert.EqualError(t, err, tt.expected.Error(), "Expected error message to match")
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	tests_scenarios := []testCase{
		{
			testName: "Valid User Deletion",
			repoSetup: func(repo *DBExecutorMock) {
				repo.BeginFunc = func() (TxExecutor, error) {
					return &TxMock{
						ExecFunc: func(query string, args ...interface{}) (sql.Result, error) {
							return nil, nil
						},
						RollbackFunc: func() error { return nil },
						CommitFunc:   func() error { return nil },
					}, nil
				}
			},
			input: entity.User{
				ID:    uuid.New(),
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			expected: nil,
		},
		{
			testName: "Error on User Deletion",
			repoSetup: func(repo *DBExecutorMock) {
				repo.BeginFunc = func() (TxExecutor, error) {
					return &TxMock{
						ExecFunc: func(query string, args ...interface{}) (sql.Result, error) {
							return nil, errors.New("database error")
						},
						RollbackFunc: func() error { return nil },
						CommitFunc:   func() error { return nil },
					}, nil
				}
			},
			input: entity.User{
				ID:    uuid.New(),
				Name:  "Jane Mary",
				Email: "jane.mary@example.com",
			},
			expected: errors.New("database error"),
		},
	}
	for _, tt := range tests_scenarios {
		t.Run(tt.testName, func(t *testing.T) {
			dbExecutor := &DBExecutorMock{}
			tt.repoSetup(dbExecutor)

			repo := NewPostgresUserRepository(dbExecutor)
			err := repo.Delete(tt.input)

			switch tt.testName {
			case tests_scenarios[0].testName:
				assert.NoError(t, err, "Expected no error for valid user deletion")
				assert.Equal(t, tt.expected, err, "Expected error equal to nil")
			case tests_scenarios[1].testName:
				assert.Error(t, err, "Expected an error for user deletion failure")
				assert.EqualError(t, err, tt.expected.Error(), "Expected error message to match")
			}
		})
	}
}
