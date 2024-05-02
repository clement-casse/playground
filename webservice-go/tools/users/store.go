package users

import "context"

type Store interface {
	Authenticator
	CreateUser(ctx context.Context, name, email string) (*User, error)
	DeleteUser(ctx context.Context, user *User) error
}
