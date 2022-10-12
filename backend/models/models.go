package models

import "github.com/google/uuid"

type User struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

func NewUser(username string) *User {
	user := new(User)
	user.Username = username
	user.Id = uuid.NewString()
	return user
}
