// Clean Architecture - Use Case Layer
// User use case interface and implementation
package usecase

import (
	"clean-go-rest-api/internal/domain/dto"
	"clean-go-rest-api/internal/domain/entity"
	"errors"

	"github.com/google/uuid"
)

type IUserUseCase interface {
	Add(req dto.CreateUserRequest) (uuid.UUID, error)
	Delete(req dto.DeleteUserRequest) error
	Update(req dto.UpdateUserRequest) error
	GetById(id uuid.UUID) (entity.User, error)
	Search(name string) ([]entity.User, error)
}

type UserUseCase struct {
	repo entity.IUserRepository
}

func NewUserUseCase(repo entity.IUserRepository) IUserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) Add(req dto.CreateUserRequest) (uuid.UUID, error) {
	if u.repo.EmailExists(req.Email) {
		return uuid.Nil, errors.New("user already exists")
	}
	user := entity.User{
		ID:    uuid.New(),
		Name:  req.Name,
		Email: req.Email,
	}
	if err := u.repo.Add(user); err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (u *UserUseCase) Delete(req dto.DeleteUserRequest) error {
	user, err := u.repo.GetById(req.ID)
	if err != nil {
		return err
	}
	if user.ID == uuid.Nil {
		return errors.New("user not found")
	}
	return u.repo.Delete(user)
}

func (u *UserUseCase) Update(req dto.UpdateUserRequest) error {
	user, err := u.repo.GetById(req.ID)
	if err != nil {
		return err
	}
	if user.ID == uuid.Nil {
		return errors.New("user not found")
	}
	user.Name = req.Name
	user.Email = req.Email
	return u.repo.Update(user)
}

func (u *UserUseCase) GetById(id uuid.UUID) (entity.User, error) {
	user, err := u.repo.GetById(id)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u *UserUseCase) Search(name string) ([]entity.User, error) {
	return u.repo.Search(name)
}
