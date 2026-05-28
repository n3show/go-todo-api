package repository

import (
	"database/sql"
	"go-todo-api/model"
)

type TodoRepository interface {
	GetAll(userId int) ([]model.Todo, error)
	Create(todo model.Todo) error
	Update(todo model.Todo) error
	Delete(id int) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) GetAll(userId int) ([]model.Todo, error) {
	rows, err := r.db.Query("SELECT id, user_id, title, done, created_at FROM todos WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo

	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.Id, &todo.UserId, &todo.Title, &todo.Done, &todo.CreatedAt); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *todoRepository) Create(todo model.Todo) error {
	_, err := r.db.Exec("INSERT INTO todos (user_id, title, done, created_at) VALUES ($1, $2, $3, $4)", todo.UserId, todo.Title, todo.Done, todo.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *todoRepository) Update(todo model.Todo) error {
	_, err := r.db.Exec("UPDATE todos SET title=$1, done=$2 WHERE id=$3", todo.Title, todo.Done, todo.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *todoRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
