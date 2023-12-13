package domain

import "context"

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	FindById(ctx context.Context, id int) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
	Store(ctx context.Context, user *User) error
}
