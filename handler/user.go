package handler

import (
	"encoding/json"
	"go-todo-api/model"
	"go-todo-api/service"
	"net/http"
)

type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userServie service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userServie: userService,
	}
}

type AuthRequest struct {
	Username string
	Password string
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := model.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := h.userServie.Register(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.userServie.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
