package service

import (
	"go-todo-api/model"
	"go-todo-api/repository"
)

type TodoService interface {
	GetAll(userId int) ([]model.Todo, error)
	GetById(id int) (model.Todo, error)
	Create(todo model.Todo) error
	Update(todo model.Todo) error
	Delete(id int) error
}

type todoService struct {
	todoRepository repository.TodoRepository
}

func NewTodoService(todoRepository repository.TodoRepository) TodoService {
	return &todoService{
		todoRepository: todoRepository,
	}
}

func (s *todoService) GetAll(userId int) ([]model.Todo, error) {
	return s.todoRepository.GetAll(userId)
}

func (s *todoService) GetById(id int) (model.Todo, error) {
	return s.todoRepository.GetById(id)
}

func (s *todoService) Create(todo model.Todo) error {
	return s.todoRepository.Create(todo)
}

func (s *todoService) Update(todo model.Todo) error {
	return s.todoRepository.Update(todo)
}

func (s *todoService) Delete(id int) error {
	return s.todoRepository.Delete(id)
}
