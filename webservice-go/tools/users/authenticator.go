package users

import (
	"context"
	"errors"
)

//go:generate mockgen -destination=../../test/mocks/users/authenticator_gen.go -package=users github.com/clement-casse/playground/webservice-go/tools/users Authenticator

// ErrAuthenticationFailure indicates that the user cannot be identified with the provided information.
var ErrAuthenticationFailure = errors.New("authentication failure")

// Authenticator allows to check whether a user is allowed in the system or not.
type Authenticator interface {
	Authenticate(ctx context.Context, userID string, params ...string) (*User, error)
}
