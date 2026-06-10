package handler

import (
	"encoding/json"
	"fmt"
	"go-todo-api/model"
	"go-todo-api/service"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type TodoHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type todoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) TodoHandler {
	return &todoHandler{
		todoService: todoService,
	}
}

func (h *todoHandler) List(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	todos, err := h.todoService.GetAll(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *todoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.UserId = userId

	if err := h.todoService.Create(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *todoHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if todo.UserId != userId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var req model.Todo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.Title = req.Title
	todo.Done = req.Done

	if err := h.todoService.Update(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *todoHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func getUserIdFromToken(r *http.Request) (int, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return 0, fmt.Errorf("no token provided")
	}

	token, err := service.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userId := int(claims["user_id"].(float64))
	return userId, nil
}
