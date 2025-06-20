package usecase

import (
	"errors"
	"reflect"
	"testing"

	"clean-go-rest-api/internal/domain/dto"
	"clean-go-rest-api/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type UserUsecaseCreateUserTestSuite struct {
	testName  string
	repoSetup func(*UserRepositoryMock)
	input     dto.CreateUserRequest
	expected  uuid.UUID
}

type UserUsecaseDeleteUserTestSuite struct {
	testName  string
	repoSetup func(*UserRepositoryMock)
	input     dto.DeleteUserRequest
	expected  error
}

func TestUserUseCase_Add(t *testing.T) {
	tests_scenarios := []UserUsecaseCreateUserTestSuite{
		{
			testName: "Valid User Creation",
			repoSetup: func(repo *UserRepositoryMock) {
				repo.emailExist = false
				repo.addErr = nil
			},
			input: dto.CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expected: uuid.UUID{},
		},
		{
			testName: "User Already Exists",
			repoSetup: func(repo *UserRepositoryMock) {
				repo.emailExist = true
				repo.addErr = nil
			},
			input: dto.CreateUserRequest{
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
			expected: uuid.Nil,
		},
		{
			testName: "Error Adding User",
			repoSetup: func(repo *UserRepositoryMock) {
				repo.emailExist = false
				repo.addErr = errors.New("database error")
			},
			input: dto.CreateUserRequest{
				Name:  "Jonas Doe",
				Email: "jonas@example.com",
			},
			expected: uuid.Nil,
		},
	}

	for _, tt := range tests_scenarios {
		t.Run(tt.testName, func(t *testing.T) {
			repo := SetupMockRepo()
			tt.repoSetup(repo)

			useCase := NewUserUseCase(repo)
			result, err := useCase.Add(tt.input)

			switch tt.testName {
			case tests_scenarios[0].testName:
				assert.NoError(t, err, "should not return an error for user creation")
				assert.NotEqual(t, uuid.Nil, result, "should return a valid UUID")
				assert.Equal(t,
					reflect.TypeOf(tt.expected), reflect.TypeOf(result),
					"the type should be `uuid.UUID`",
				)
			case tests_scenarios[1].testName:
				assert.Error(t, err, "should return an error for existing user")
				assert.Equal(t, uuid.Nil, result, "should return uuid.Nil for existing user")
			case tests_scenarios[2].testName:
				assert.Error(t, err, "should return an error when adding user fails")
				assert.Equal(t, uuid.Nil, result, "should return uuid.Nil when adding user fails")
				assert.EqualError(t,
					err, "database error", "should return the correct error message",
				)
			}
		})
	}
}

func TestUserUseCase_Delete(t *testing.T) {
	tests_scenarios := []UserUsecaseDeleteUserTestSuite{
		{
			testName: "Valid User Deletion",
			repoSetup: func(repo *UserRepositoryMock) {
				userID := uuid.New()
				repo.users[userID.String()] = entity.User{
					ID:    userID,
					Name:  "John Doe",
					Email: "john@example.com",
				}
				repo.deleteErr = nil
				repo.getByIdErr = nil
			},
			input: dto.DeleteUserRequest{
				ID: uuid.New(),
			},
			expected: nil,
		},
		{
			testName: "User Not Found",
			repoSetup: func(repo *UserRepositoryMock) {
				repo.getByIdErr = nil
			},
			input: dto.DeleteUserRequest{
				ID: uuid.New(),
			},
			expected: errors.New("user not found"),
		},
		{
			testName: "Error Deleting User",
			repoSetup: func(repo *UserRepositoryMock) {
				userID := uuid.New()
				repo.users[userID.String()] = entity.User{
					ID:    userID,
					Name:  "Jane Doe",
					Email: "jane@example.com",
				}
				repo.deleteErr = errors.New("delete error")
				repo.getByIdErr = nil
			},
			input: dto.DeleteUserRequest{
				ID: uuid.New(),
			},
			expected: errors.New("delete error"),
		},
		{
			testName: "Error Getting User",
			repoSetup: func(repo *UserRepositoryMock) {
				repo.getByIdErr = errors.New("get error")
			},
			input: dto.DeleteUserRequest{
				ID: uuid.New(),
			},
			expected: errors.New("get error"),
		},
	}

	for _, tt := range tests_scenarios {
		t.Run(tt.testName, func(t *testing.T) {
			repo := SetupMockRepo()
			tt.repoSetup(repo)

			if tt.testName == tests_scenarios[0].testName ||
				tt.testName == tests_scenarios[2].testName {
				for id := range repo.users {
					uuidVal, _ := uuid.Parse(id)
					tt.input.ID = uuidVal
				}
			}

			useCase := NewUserUseCase(repo)
			err := useCase.Delete(tt.input)

			switch tt.testName {
			case tests_scenarios[0].testName:
				assert.NoError(t, err, "should not return an error for valid deletion")
			case tests_scenarios[1].testName:
				assert.Error(t, err, "should return an error for user not found")
				assert.EqualError(t, err, "user not found", "should return the correct error message")
			case tests_scenarios[2].testName:
				assert.Error(t, err, "should return an error when delete fails")
				assert.EqualError(t, err, "delete error", "should return the correct error message")
			case tests_scenarios[3].testName:
				assert.Error(t, err, "should return an error when get fails")
				assert.EqualError(t, err, "get error", "should return the correct error message")
			}
		})
	}
}
