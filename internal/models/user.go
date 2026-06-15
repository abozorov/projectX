package models

import (
	"strings"
	"time"
)

type User struct {
	ID        int
	Name      string
	Login     string
	Password  string
	IsActive  bool
	CreatedAt time.Time
}

func NewUser(name, login, password string) *User {
	return &User{
		Name:      name,
		Login:     login,
		Password:  password,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
}

func (u *User) Validate(update bool) bool {
	if update && u.ID <= 0 {
		return false
	}
	u.Name = strings.TrimSpace(u.Name)
	u.Login = strings.TrimSpace(u.Login)
	u.Password = strings.TrimSpace(u.Password)

	return u.Name != "" && u.Login != "" && (len([]rune(u.Password)) > 7 || update)
}