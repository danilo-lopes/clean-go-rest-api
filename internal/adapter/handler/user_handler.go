// Clean Architecture - Interface Adapter Layer
// HTTP Handlers for User
package handler

import (
	"clean-go-rest-api/internal/crosscutting/logger"
	"clean-go-rest-api/internal/domain/dto"
	"clean-go-rest-api/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	useCase usecase.IUserUseCase
	logger  logger.ILogger
}

func NewUserHandler(UseCase usecase.IUserUseCase, Logger logger.ILogger) *UserHandler {
	return &UserHandler{useCase: UseCase, logger: logger.NewLogger()}
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", h.Add).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/users/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/users", h.Search).Methods(http.MethodGet)
}

func (h *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: err.Error()})
		h.logger.Error("Error decoding request body: " + err.Error())
		return
	}

	h.logger.Info("Received request to create user")
	id, err := h.useCase.Add(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: err.Error()})
		h.logger.Error("Error creating user: " + err.Error())
		return
	}

	h.logger.Info(fmt.Sprintf("User created successfully with id: %s", id.String()))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateUserResponse{ID: id})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: "invalid ID format"})
		h.logger.Error("Error parsing ID: " + err.Error())
		return
	}

	h.logger.Info(fmt.Sprintf("Received request to delete user with ID: %s", id))
	err = h.useCase.Delete(dto.DeleteUserRequest{ID: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: err.Error()})
		h.logger.Error("Error creating user: " + err.Error())
		return
	}

	h.logger.Info(fmt.Sprintf("User with ID %s deleted successfully", id))
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: "invalid ID format"})
		h.logger.Error("Error parsing ID: " + err.Error())
		return
	}

	h.logger.Info(fmt.Sprintf("Received request to update user with ID: %s", id))
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: err.Error()})
		h.logger.Error("Error decoding request body: " + err.Error())
		return
	}
	req.ID = id
	err = h.useCase.Update(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: err.Error()})
		h.logger.Error("Error updating user: " + err.Error())
		return
	}

	h.logger.Info(fmt.Sprintf("User with ID %s updated successfully", id))
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: "invalid ID format"})
		h.logger.Error("Error parsing ID: " + err.Error())
		return
	}

	h.logger.Info(fmt.Sprintf("Received request to get user with ID: %s", id))
	user, err := h.useCase.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: err.Error()})
		h.logger.Error("Error getting user: " + err.Error())
		return
	}
	if user.ID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: "user not found"})
		h.logger.Error(fmt.Sprintf("User with ID %s not found", id))
		return
	}

	h.logger.Info(fmt.Sprintf("User with ID %s retrieved successfully", id))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: "name parameter is required"})
		h.logger.Error("Error: name parameter is required")
		return
	}

	h.logger.Info(fmt.Sprintf("Received request to search users with name: %s", name))
	users, err := h.useCase.Search(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Reason: err.Error()})
		h.logger.Error("Error searching users: " + err.Error())
		return
	}

	h.logger.Info(fmt.Sprintf("Found %d users with name: %s", len(users), name))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
