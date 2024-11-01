package users

import (
	"context"
	"errors"
)

//go:generate go run go.uber.org/mock/mockgen -destination=../../test/mocks/users/authenticator_gen.go -package=users . Authenticator

// ErrAuthenticationFailure indicates that the user cannot be identified with the provided information.
var ErrAuthenticationFailure = errors.New("authentication failure")

// Authenticator allows to check whether a user is allowed in the system or not.
type Authenticator interface {
	Authenticate(ctx context.Context, userID string, params ...string) (*User, error)
}
