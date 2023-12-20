package handler

import (
	"awesomeProject/internal/models/user"
	"awesomeProject/internal/repository"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	UserRepo repository.UserRepository
}

func NewUserHandler(user repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: user}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserReq CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	us := user.User{
		Username:  createUserReq.Username,
		Email:     createUserReq.Email,
		Password:  createUserReq.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := u.UserRepo.Create(context.Background(), us)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "User created successfully",
	}
	json.NewEncoder(w).Encode(response)
}
func (u *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := u.UserRepo.GetByID(context.Background(), userID)
	if err != nil {

		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func (u *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var updatedUser user.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := u.UserRepo.Update(context.Background(), updatedUser); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (u *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := u.UserRepo.Delete(context.Background(), id); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
func (u *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	users, err := u.UserRepo.List(context.Background(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
