package models

type User struct {
	ID       int
	Name     string
	Password string
	IsActive bool
}

func NewUser(id int, name, password string ) *User {
	return &User{
		ID:   id,
		Name: name,
		Password: password,
		IsActive: true,
	}
}
