package repository

import (
	"database/sql"
	"go-todo-api/model"
)

type UserRespository interface {
	Create(user model.User) error
	GetByUsername(username string) (model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRespository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user model.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByUsername(username string) (model.User, error) {
	var user model.User

	if err := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.Id, &user.Username, &user.Password); err != nil {
		return model.User{}, err
	}

	return user, nil
}
