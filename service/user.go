package service

import (
	"fmt"
	"go-todo-api/model"
	"go-todo-api/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user model.User) error
	Login(username, password string) (string, error)
}

type userService struct {
	userRepository repository.UserRespository
}

func NewUserService(userRepository repository.UserRespository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Register(user model.User) error {
	if user.Username == "" || user.Password == "" {
		return fmt.Errorf("username or password cannot be empty")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	if err := s.userRepository.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepository.GetByUsername(username)
	if err != nil {
		return "", err
	}

	if !checkPasswordHash(user.Password, password) {
		return "", fmt.Errorf("username or password incorrect!")
	}

	token, err := generateToken(user.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}
