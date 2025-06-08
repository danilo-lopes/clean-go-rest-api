// Clean Architecture - Domain Layer
// Entities and repository interfaces for User
package entity

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

type IUserRepository interface {
	Add(user User) error
	Delete(user User) error
	Update(user User) error
	GetById(id uuid.UUID) (User, error)
	Search(name string) ([]User, error)
	EmailExists(email string) bool
}
