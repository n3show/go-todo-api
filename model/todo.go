package model

import "time"

type Todo struct {
	Id        int
	UserId    int
	Title     string
	Done      bool
	CreatedAt time.Time
}
