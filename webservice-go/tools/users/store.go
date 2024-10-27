package users

import "context"

//go:generate go run go.uber.org/mock/mockgen -destination=../../test/mocks/users/store_gen.go -package=users github.com/clement-casse/playground/webservice-go/tools/users Store

type Store interface {
	Authenticator
	CreateUser(ctx context.Context, name, email, password string) (*User, error)
	DeleteUser(ctx context.Context, user *User) error
}
