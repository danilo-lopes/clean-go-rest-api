package usecase

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"clean-go-rest-api/internal/domain/dto"
	"clean-go-rest-api/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type createUserTestCase struct {
	testName  string
	repoSetup func(*UserRepositoryMock)
	input     dto.CreateUserRequest
	expected  uuid.UUID
}

func TestUserUseCase_Add(t *testing.T) {
	tests_scenarios := []createUserTestCase{
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
				fmt.Println("result: ", result)
				assert.Equal(t, tt.expected, result, "should return uuid.Nil when adding user fails")
				assert.EqualError(t,
					err, "database error", "should return the correct error message",
				)
			}
		})
	}
}

type deleteUserTestCase struct {
	testName  string
	repoSetup func(*UserRepositoryMock)
	input     dto.DeleteUserRequest
	expected  error
}

func TestUserUseCase_Delete(t *testing.T) {
	tests_scenarios := []deleteUserTestCase{
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

type updateUserTestCase struct {
	testName  string
	repoSetup func(*UserRepositoryMock)
	input     dto.UpdateUserRequest
	expected  interface{}
}

func TestUserUseCase_Update(t *testing.T) {
	tests_scenarios := []updateUserTestCase{
		{
			testName: "Valid Update",
			repoSetup: func(repo *UserRepositoryMock) {
				userID := uuid.New()
				repo.users[userID.String()] = entity.User{
					ID:    userID,
					Name:  "John Doe",
					Email: "john.doe@xample.com",
				}
				repo.getByIdErr = nil
				repo.updateErr = nil
			},
		},
		{
			testName: "Error GetById",
			repoSetup: func(repo *UserRepositoryMock) {
				repo.getByIdErr = errors.New("get error")
				repo.updateErr = nil
			},
			expected: errors.New("get error"),
		},
		{
			testName: "User not found",
			repoSetup: func(repo *UserRepositoryMock) {
				repo.getByIdErr = nil
				repo.updateErr = nil
			},
			expected: errors.New("user not found"),
		},
	}

	for _, tt := range tests_scenarios {
		t.Run(tt.testName, func(t *testing.T) {
			repo := SetupMockRepo()
			tt.repoSetup(repo)

			if tt.testName == tests_scenarios[0].testName {
				for id := range repo.users {
					uuidVal, _ := uuid.Parse(id)
					tt.input.ID = uuidVal
					tt.input.Name = "John new"
					tt.input.Email = "john.new@example.com"
				}
			}

			useCase := NewUserUseCase(repo)
			err := useCase.Update(tt.input)

			switch tt.testName {
			case tests_scenarios[0].testName:
				assert.NoError(t, err, "should not return an error for valid update")
			case tests_scenarios[1].testName:
				assert.Error(t, err, "should return an error for mocked error")
				assert.Equal(t, err, tt.expected)
			case tests_scenarios[2].testName:
				assert.Error(t, err, "should return an error for unextisting user")
				assert.Equal(t, err, tt.expected)
			}
		})
	}
}

type getByIdUserTestCase struct {
	testName  string
	repoSetup func(*UserRepositoryMock)
	input     dto.UpdateUserRequest
	expected  interface{}
}

func TestUserUseCase_GetById(t *testing.T) {
	tests_scenarios := []getByIdUserTestCase{
		{
			testName: "User found",
			repoSetup: func(repo *UserRepositoryMock) {
				userId := uuid.New()
				repo.users[userId.String()] = entity.User{
					ID:    userId,
					Name:  "Jane Mary",
					Email: "jane.mary@example.com",
				}
				repo.getByIdErr = nil
			},
		},
		{
			testName: "User not found",
			repoSetup: func(repo *UserRepositoryMock) {
				repo.getByIdErr = errors.New("user not found")
			},
		},
	}

	for _, tt := range tests_scenarios {
		t.Run(tt.testName, func(t *testing.T) {
			repo := SetupMockRepo()
			tt.repoSetup(repo)

			if tt.testName == tests_scenarios[0].testName {
				for id := range repo.users {
					uuidVal, _ := uuid.Parse(id)
					tt.input.ID = uuidVal
				}
			}

			useCase := NewUserUseCase(repo)
			user, err := useCase.GetById(tt.input.ID)

			switch tt.testName {
			case tests_scenarios[0].testName:
				tt.expected = repo.users[tt.input.ID.String()]
				assert.NoError(t, err, "should not return an error for existent user")
				assert.Equal(t, tt.expected, user, "user must be the expected")
			case tests_scenarios[1].testName:
				tt.expected = errors.New("user not found")
				assert.Error(t, err, "should return an error for unexisting user")
				assert.Equal(t, err, tt.expected)
				assert.Equal(t, user, entity.User{})
			}
		})
	}
}

type searchUserTestCase struct {
	testName  string
	repoSetup func(*UserRepositoryMock)
	input     string
	expected  []entity.User
}

func TestUserUseCase_Search(t *testing.T) {
	tests_scenarios := []searchUserTestCase{
		{
			testName: "User found",
			repoSetup: func(repo *UserRepositoryMock) {
				userId := uuid.New()
				repo.users[userId.String()] = entity.User{
					ID:    userId,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}
			},
			input: "John",
		},
		{
			testName: "User Not Found",
			repoSetup: func(repo *UserRepositoryMock) {
			},
			input: "John",
		},
	}

	for _, tt := range tests_scenarios {
		t.Run(tt.testName, func(t *testing.T) {
			repo := SetupMockRepo()
			tt.repoSetup(repo)

			if tt.testName == tests_scenarios[0].testName {
				for id := range repo.users {
					tt.expected = append(tt.expected, repo.users[id])
				}
			}

			useCase := NewUserUseCase(repo)
			users, err := useCase.Search(tt.input)

			switch tt.testName {
			case tests_scenarios[0].testName:
				assert.NoError(t, err, "should not return an error for existent user")
				assert.Equal(t, users, tt.expected)
			case tests_scenarios[1].testName:
				assert.NoError(t, err, "should not return an error to search")
				assert.Equal(t, users, tt.expected)
			}
		})
	}
}
