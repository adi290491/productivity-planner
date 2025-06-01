package models

type Repository interface {
	CreateUser(user *User) (*User, error)
	GetUser(user *User) (*User, error)
}
