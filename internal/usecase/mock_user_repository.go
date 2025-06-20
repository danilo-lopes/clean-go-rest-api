package usecase

import (
	"clean-go-rest-api/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepositoryMock struct {
	users      map[string]entity.User
	emailExist bool
	addErr     error
	getByIdErr error
	updateErr  error
	deleteErr  error
}

func SetupMockRepo() *UserRepositoryMock {
	return &UserRepositoryMock{users: make(map[string]entity.User)}
}

func (m *UserRepositoryMock) Add(user entity.User) error {
	if m.getByIdErr != nil {
		return m.getByIdErr
	}
	m.users[user.ID.String()] = user
	return nil
}

func (m *UserRepositoryMock) Delete(user entity.User) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.users, user.ID.String())
	return nil
}

func (m *UserRepositoryMock) Update(user entity.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.users[user.ID.String()] = user
	return nil
}

func (m *UserRepositoryMock) GetById(id uuid.UUID) (entity.User, error) {
	if m.getByIdErr != nil {
		return entity.User{}, m.getByIdErr
	}
	user, ok := m.users[id.String()]
	if !ok {
		return entity.User{ID: uuid.Nil}, nil
	}
	return user, nil
}

func (m *UserRepositoryMock) Search(name string) ([]entity.User, error) {
	var result []entity.User
	for _, u := range m.users {
		if u.Name == name {
			result = append(result, u)
		}
	}
	return result, nil
}

func (m *UserRepositoryMock) EmailExists(email string) bool {
	return m.emailExist
}
