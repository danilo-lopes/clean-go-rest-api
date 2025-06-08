// Clean Architecture - Domain Layer
// Use case input/output DTOs
package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateUserRequest struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type ErrorResponse struct {
	Reason string `json:"reason"`
}
