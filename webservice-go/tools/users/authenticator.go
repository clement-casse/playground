package users

import (
	"context"
	"errors"
)

// ErrAuthenticationFailure indicates that the user cannot be identified with the provided information.
var ErrAuthenticationFailure = errors.New("authentication failure")

// Authenticator allows to check whether a user is allowed in the system or not.
type Authenticator interface {
	Authenticate(ctx context.Context, userID string, params ...string) (*User, error)
}
